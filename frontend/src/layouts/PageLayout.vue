<template>
  <v-app-bar theme="dark" color="accent" flat :title="$route.name">
    <template #prepend>
      <v-btn v-if="isMobile" to="/" icon size="small">
        <v-icon :icon="AppIcon" size="28" style="opacity: 0.7" aria-label="Relax Sounds" />
      </v-btn>
    </template>

    <v-spacer />

    <slot name="actions" />

    <v-menu v-if="isMobile || $slots.menu" location="bottom right" transition="slide-y-transition">
      <template #activator="{ props }">
        <v-btn icon variant="flat" color="transparent" v-bind="props" aria-label="Menu">
          <v-icon :icon="MenuIcon" />
        </v-btn>
      </template>

      <v-list>
        <slot name="menu" />
        <template v-if="isMobile">
          <v-divider v-if="$slots.menu" />
          <debug-button v-if="DebugEnabled" list-item />
          <theme-btn list-item />
        </template>
      </v-list>
    </v-menu>
  </v-app-bar>

  <v-container class="pt-6 pt-lg-12">
    <slot />
  </v-container>
</template>

<script setup>
import { useDisplay } from "vuetify";
import ThemeBtn from "../components/NavButtons/ThemeBtn.vue";
import DebugButton from "../components/Presets/Buttons/DebugButton.vue";
import { DebugEnabled } from "../config/debug";
import MenuIcon from "~icons/material-symbols/more-horiz";
import AppIcon from "~icons/relax-sounds/icon-white";

const { smAndDown: isMobile } = useDisplay();
</script>
