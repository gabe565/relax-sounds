import { computedAsync, useEventListener } from "@vueuse/core";
import { defineStore } from "pinia";
import PocketBase from "pocketbase";
import { computed, nextTick, onScopeDispose, ref } from "vue";
import { useRouter } from "vue-router";
import { toast } from "vue-sonner";
import AuthentikIcon from "@/assets/authentik.svg";
import { ApiPath } from "@/config/api.js";
import { Sound } from "@/util/Sound";
import { once } from "@/util/helpers";

export function getErrorMessage(error) {
  const response = error.response;
  if (response?.data && typeof response.data === "object") {
    const messages = [];
    for (const key in response.data) {
      if (response.data[key]?.message) {
        messages.push(`${key}: ${response.data[key].message}`);
      }
    }
    if (messages.length > 0) {
      return messages.join("\n");
    }
  }
  return response?.message || error.message || "An unexpected error occurred";
}

export function getProviderIconURL(provider) {
  if (provider.displayName?.toLowerCase() === "authentik") {
    return AuthentikIcon;
  }
  return ApiPath(`/_/images/oauth2/${provider.name.toLowerCase()}.svg`);
}

const AuthRefreshInterval = 60 * 60 * 1000;

export const SessionExpiredToast = "session-expired";

export const usePocketBase = defineStore("pocketbase", () => {
  const client = new PocketBase(ApiPath());
  const user = ref(client.authStore.record);
  const hasValidToken = ref(client.authStore.isValid);
  const isAuthenticated = computed(() => hasValidToken.value && !!user.value?.verified);

  const unsubscribeAuth = client.authStore.onChange((token, record) => {
    user.value = record;
    hasValidToken.value = client.authStore.isValid;
  });
  onScopeDispose(unsubscribeAuth);

  let lastRefresh = 0;

  const router = useRouter();

  const expireSession = async () => {
    const hadSession = !!client.authStore.token;
    client.authStore.clear();
    if (!hadSession) {
      return;
    }
    await nextTick();
    toast.info("Your session has expired.", {
      id: SessionExpiredToast,
      duration: Infinity,
      action: {
        label: "Log in",
        onClick: () => router.push({ name: "Log in" }),
      },
    });
  };

  const refreshAuth = async () => {
    lastRefresh = Date.now();
    if (!client.authStore.isValid) {
      await expireSession();
      return;
    }
    try {
      await client.collection("users").authRefresh();
    } catch (err) {
      if (err.status === 401 || err.status === 403 || err.status === 404) {
        await expireSession();
      } else {
        console.error("Failed to refresh auth:", err);
      }
    }
  };

  const authReady = refreshAuth();

  useEventListener(window, "focus", () => {
    if (client.authStore.token && Date.now() - lastRefresh > AuthRefreshInterval) {
      refreshAuth();
    }
  });

  const authMethods = computedAsync(
    async () => {
      const res = await client.collection("users").listAuthMethods();
      if (res.oauth2?.providers) {
        res.oauth2.providers = res.oauth2?.providers.map((provider) => ({
          ...provider,
          icon: getProviderIconURL(provider),
        }));
      }
      return res;
    },
    { loading: true },
    { lazy: true },
  );

  const authEnabled = computed(() => {
    return (
      isAuthenticated.value ||
      authMethods.value.password?.enabled ||
      authMethods.value.oauth2?.enabled
    );
  });

  const avatarURL = computed(() => {
    if (user.value?.avatar) {
      return client.files.getURL(user.value, user.value.avatar);
    }
    return null;
  });

  const loadSounds = once(async () => {
    const data = await client.collection("sounds").getFullList({
      fields: "collectionId,id,short_id,name,icon,file,expand.tags.name",
      expand: "tags",
      sort: "name",
    });
    return data.map((sound) => {
      sound.tags = sound.expand?.tags?.map((tag) => tag.name);
      delete sound.expand;
      return new Sound(sound);
    });
  });

  const loadTags = once(() =>
    client.collection("tags").getFullList({
      fields: "icon,name",
    }),
  );

  const logout = () => {
    client.authStore.clear();
    toast.success("Logged out.");
  };

  return {
    client,
    user,
    isAuthenticated,
    refreshAuth,
    authReady,
    authEnabled,
    authMethods,
    avatarURL,
    loadSounds,
    loadTags,
    logout,
  };
});
