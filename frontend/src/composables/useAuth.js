import { computedAsync } from "@vueuse/core";
import { computed, ref } from "vue";
import { useToast } from "vue-toastification";
import AuthentikIcon from "@/assets/authentik.svg";
import { ApiPath } from "@/config/api.js";
import { pb } from "@/plugins/pocketbase";

const user = ref(pb.authStore.record);
const isAuthenticated = computed(() => !!user.value?.verified);

pb.authStore.onChange((token, record) => {
  user.value = record;
});

/**
 * Returns the icon URL for a given OAuth2 provider.
 * @param {Object} provider - The provider object.
 * @returns {string} - The icon URL or path.
 */
const getProviderIconURL = (provider) => {
  if (provider.displayName?.toLowerCase() === "authentik") {
    return AuthentikIcon;
  }
  return ApiPath(`/_/images/oauth2/${provider.name.toLowerCase()}.svg`);
};

const authMethods = computedAsync(
  async () => {
    const res = await pb.collection("users").listAuthMethods();
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

const logout = () => {
  pb.authStore.clear();
  useToast().success("Logged out.");
};

const avatarURL = computed(() => {
  if (user.value?.avatar) {
    return pb.files.getURL(user.value, user.value.avatar);
  }
  return null;
});

/**
 * Composable for accessing reactive PocketBase authentication state.
 */
export function useAuth() {
  return { user, isAuthenticated, authEnabled, authMethods, logout, avatarURL };
}
