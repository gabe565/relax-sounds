import { computedAsync } from "@vueuse/core";
import { defineStore } from "pinia";
import PocketBase from "pocketbase";
import { computed, onScopeDispose, ref } from "vue";
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

export const usePocketBase = defineStore("pocketbase", () => {
  const client = new PocketBase(ApiPath());
  const user = ref(client.authStore.record);
  const isAuthenticated = computed(() => !!user.value?.verified);

  const unsubscribeAuth = client.authStore.onChange((token, record) => {
    user.value = record;
  });
  onScopeDispose(unsubscribeAuth);

  const refreshAuth = async () => {
    if (client.authStore.isValid) {
      await client.collection("users").authRefresh();
    }
  };
  refreshAuth().catch(console.error);

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
      fields: "collectionId,id,old_id,name,icon,file,expand.tags.name",
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
    authEnabled,
    authMethods,
    avatarURL,
    loadSounds,
    loadTags,
    logout,
  };
});
