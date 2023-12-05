<template>
  <v-container>
    <v-row class="pa-5">
      <v-col>
        <h1 class="text-h4">
          <slot name="title">
            {{ $route.name }}
          </slot>
        </h1>
      </v-col>

      <v-spacer />

      <v-col v-if="actions" class="flex-grow-0">
        <v-menu location="bottom right" transition="slide-y-transition">
          <template #activator="{ props }">
            <v-btn
              icon
              variant="flat"
              color="transparent"
              v-bind="props"
              aria-label="Menu"
              density="comfortable"
            >
              <v-icon :icon="MenuIcon" size="large" />
            </v-btn>
          </template>

          <v-list>
            <slot name="actions" />
          </v-list>
        </v-menu>
      </v-col>
    </v-row>

    <slot />
  </v-container>
</template>

<script setup>
import MenuIcon from "~icons/material-symbols/more-horiz";

defineProps({
  actions: {
    type: Boolean,
    default: false,
  },
});
</script>
