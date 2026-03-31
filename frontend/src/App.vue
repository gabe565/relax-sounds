<template>
  <v-app>
    <v-navigation-drawer
      v-if="!isMobile"
      :rail="preferences.shrinkLeftPanel"
      color="surface"
      width="200"
      class="border-e-0"
    >
      <template #prepend>
        <v-list v-if="!isMobile">
          <v-list-item to="/" title="Relax Sounds">
            <template #prepend>
              <v-icon :icon="AppIcon" color="secondary" />
            </template>
          </v-list-item>
        </v-list>
        <v-divider />
      </template>

      <v-list nav>
        <v-list-item
          v-for="route in routes"
          :key="route.path"
          :to="route.path"
          exact
          link
          :title="route.name"
          :prepend-icon="route.meta.icon"
        />
      </v-list>
    </v-navigation-drawer>

    <v-main>
      <router-view v-slot="{ Component }">
        <keep-alive exclude="LoginPage,ForgotPasswordPage">
          <component :is="Component" />
        </keep-alive>
      </router-view>
      <v-spacer v-if="isMobile && showPlayerBar" class="h-16" />
    </v-main>

    <player-bar v-if="showPlayerBar" />

    <v-bottom-navigation
      v-if="isMobile && showPlayerBar"
      v-model="botnav"
      bg-color="surface"
      color="primary"
      grow
      :model-value="null"
    >
      <v-btn v-for="route in routes" :key="route.path" :to="route.path" :value="route.name">
        <v-icon v-if="route.meta.icon" :icon="route.meta.icon" />
        <span>{{ route.name }}</span>
      </v-btn>
    </v-bottom-navigation>
  </v-app>
</template>

<script setup>
import { computed, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useDisplay, useTheme } from "vuetify";
import AppIcon from "~icons/relax-sounds/icon";
import PlayerBar from "@/components/NavButtons/PlayerBar.vue";
import { useAuth } from "@/composables/useAuth";
import { registerSW } from "@/plugins/pwa";
import { Theme, usePreferencesStore } from "@/plugins/store/preferences";

const { smAndDown: isMobile } = useDisplay();
const preferences = usePreferencesStore();
const theme = useTheme();
const { isAuthenticated } = useAuth();
const route = useRoute();

const routes = computed(() => {
  return useRouter().options.routes.filter((route) => {
    if (!route.meta?.showInNav) return false;
    if (route.meta.guestOnly && isAuthenticated.value) return false;
    if (route.meta.authOnly && !isAuthenticated.value) return false;
    return true;
  });
});

const showPlayerBar = computed(() => {
  return route.meta?.showInNav !== false;
});

const autoTheme = (e) => (theme.name.value = e.matches ? "dark" : "light");
const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
watch(
  () => preferences.theme,
  (val) => {
    mediaQuery.removeEventListener("change", autoTheme);
    switch (val) {
      case Theme.dark:
        theme.name.value = "dark";
        break;
      case Theme.light:
        theme.name.value = "light";
        break;
      default:
        // check for browser support
        if (window.matchMedia && window.matchMedia("(prefers-color-scheme)").media !== "not all") {
          // set to preferred scheme
          autoTheme(mediaQuery);
          // react to changes
          mediaQuery.addEventListener("change", autoTheme);
        }
    }
  },
  { immediate: true },
);

registerSW();
</script>
