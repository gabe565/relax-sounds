<template>
  <v-app-bar v-if="isMobile" color="surface" flat :title="route.name">
    <template #prepend>
      <v-btn to="/" icon size="small">
        <v-icon :icon="AppIcon" size="28" color="primary" aria-label="Relax Sounds" />
      </v-btn>
    </template>

    <v-spacer />

    <slot name="actions" />
    <page-actions :has-menu="!!$slots.menu">
      <template #menu><slot name="menu" /></template>
    </page-actions>
    <account-menu />
  </v-app-bar>

  <v-container :class="isMobile ? 'pt-6' : 'pt-8'">
    <div v-if="!isMobile" class="mb-4 flex items-center gap-1">
      <h1 class="text-2xl">{{ route.name }}</h1>
      <v-spacer />
      <slot name="actions" />
      <page-actions :has-menu="!!$slots.menu">
        <template #menu><slot name="menu" /></template>
      </page-actions>
    </div>

    <slot />
  </v-container>
</template>

<script setup>
import { useRoute } from "vue-router";
import { useDisplay } from "vuetify";
import AppIcon from "~icons/relax-sounds/icon";
import PageActions from "@/components/NavButtons/PageActions.vue";
import AccountMenu from "@/components/Profile/AccountMenu.vue";

const { smAndDown: isMobile } = useDisplay();
const route = useRoute();
</script>
