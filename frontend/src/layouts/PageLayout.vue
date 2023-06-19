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
            <v-btn icon variant="plain" v-bind="props" aria-label="Menu">
              <v-icon :icon="MenuIcon" size="large" />
            </v-btn>
          </template>

          <v-list>
            <v-list-item v-for="(item, index) in actions" :key="index" v-on="item.on">
              <template #prepend>
                <v-icon :icon="item.icon" aria-hidden="true" />
              </template>
              <v-list-item-title>{{ item.title }}</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-col>
    </v-row>

    <v-fade-transition>
      <v-row v-if="alert.show">
        <v-col>
          <v-alert v-model="alert.show" closable prominent text :type="alert.type">
            {{ alert.message }}
          </v-alert>
        </v-col>
      </v-row>
    </v-fade-transition>

    <slot />
  </v-container>
</template>

<script setup>
import MenuIcon from "~icons/material-symbols/more-horiz";
import { useAlertStore } from "../plugins/store/alert";

defineProps({
  actions: {
    type: Array,
    default: null,
  },
});

const alert = useAlertStore();
</script>
