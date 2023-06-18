<template>
  <v-app>
    <v-app-bar theme="dark" color="accent" flat>
      <v-btn to="/" class="text-body-2 text-none px-2">
        <template #prepend>
          <img alt="Relax Sounds Logo" :src="icon" style="height: 30px" />
        </template>

        <v-app-bar-title>Relax Sounds</v-app-bar-title>
      </v-btn>

      <v-spacer />

      <google-cast-launcher
        v-pre
        class="v-btn v-btn--icon v-theme--dark v-btn--density-default v-btn--size-x-small mr-4"
      />
      <SavePreset />
      <PlayPauseAll />
      <StopAll />

      <template v-if="mdAndUp" #extension>
        <v-tabs align-tabs="center" class="w-100" color="primary">
          <v-tab v-for="route in routes" :key="route.path" :to="route.path" exact>
            <v-icon class="pr-2">fas {{ route.meta.icon }} fa-fw</v-icon>
            {{ route.name }}
          </v-tab>
        </v-tabs>
      </template>
    </v-app-bar>

    <UpdateSnackbar />

    <v-main>
      <router-view v-slot="{ Component }">
        <keep-alive>
          <component :is="Component" />
        </keep-alive>
      </router-view>
      <v-spacer :style="{ height: smAndDown ? '56px' : '28px' }" />
    </v-main>

    <v-bottom-navigation
      v-if="smAndDown"
      bg-color="accent"
      color="primary"
      theme="dark"
      mode="shift"
    >
      <v-btn v-for="route in routes" :key="route.path" :to="route.path" :value="route.name">
        <v-icon>fas {{ route.meta.icon }} fa-fw</v-icon>
        <span>{{ route.name }}</span>
      </v-btn>
    </v-bottom-navigation>
  </v-app>
</template>

<script setup>
import { useDisplay, useTheme } from "vuetify";
import { computed, onBeforeMount } from "vue";
import { useRouter } from "vue-router";
import SavePreset from "./components/SavePreset.vue";
import PlayPauseAll from "./components/PlayPauseAll.vue";
import StopAll from "./components/StopAll.vue";
import UpdateSnackbar from "./components/UpdateSnackbar.vue";
import icon from "./assets/icon-white.svg";

const { mdAndUp, smAndDown } = useDisplay();

const routes = computed(() => {
  return useRouter().options.routes.filter((route) => route.meta?.showInNav);
});

onBeforeMount(() => {
  // check for browser support
  if (window.matchMedia && window.matchMedia("(prefers-color-scheme)").media !== "not all") {
    // set to preferred scheme
    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
    const theme = useTheme();
    theme.name.value = mediaQuery.matches ? "dark" : "light";
    // react to changes
    mediaQuery.addEventListener("change", (e) => {
      theme.name.value = e.matches ? "dark" : "light";
    });
  }
});
</script>

<style lang="scss">
html {
  --disconnected-color: #fff;
}

.fa-spin-2x {
  animation: fa-spin 750ms infinite linear;
}

.v-card--variant-outlined,
.v-chip--variant-outlined,
.v-btn--variant-outlined:not(.text-primary) {
  border-color: rgba(255, 255, 255, 0.12);
}

.v-card--variant-flat {
  border: thin solid rgba(0, 0, 0, 0);
}
</style>
