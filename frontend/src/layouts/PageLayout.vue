<template>
  <v-app-bar theme="dark" color="accent" flat :title="$route.name">
    <template #prepend>
      <v-btn v-if="isMobile" to="/" icon size="small">
        <v-icon :icon="AppIcon" size="28" style="opacity: 0.7" aria-label="Relax Sounds" />
      </v-btn>
    </template>

    <v-spacer />

    <slot name="actions" />

    <v-menu v-if="$slots.menu" location="bottom right" transition="slide-y-transition">
      <template #activator="{ props }">
        <v-btn icon variant="flat" color="transparent" v-bind="props" aria-label="Menu">
          <v-icon :icon="MenuIcon" />
        </v-btn>
      </template>

      <v-list>
        <slot name="menu" />
      </v-list>
    </v-menu>
  </v-app-bar>

  <v-container class="pt-6 pt-lg-12">
    <slot />
  </v-container>
</template>

<script setup>
import AppIcon from "~icons/relax-sounds/icon-white";
import MenuIcon from "~icons/material-symbols/more-horiz";
import { useDisplay } from "vuetify";

const { smAndDown: isMobile } = useDisplay();
</script>
