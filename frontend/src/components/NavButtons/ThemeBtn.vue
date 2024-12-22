<template>
  <v-list-item
    v-if="listItem"
    :prepend-icon="config.icon"
    :title="config.text"
    aria-label="Change theme"
    @click="next"
  />
  <v-tooltip v-else :text="config.text" location="top">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        :icon="config.icon"
        color="transparent"
        variant="flat"
        aria-label="Change theme"
        @click="next"
      />
    </template>
  </v-tooltip>
</template>

<script setup>
import { shallowRef, watch } from "vue";
import { Theme, usePreferencesStore } from "../../plugins/store/preferences";
import DarkIcon from "~icons/material-symbols/brightness-4";
import LightIcon from "~icons/material-symbols/brightness-7";
import AutoIcon from "~icons/material-symbols/brightness-auto";

const props = defineProps({
  listItem: {
    type: Boolean,
    default: false,
  },
});

const preferences = usePreferencesStore();

let config = shallowRef({});

const getConfig = () => {
  switch (preferences.theme) {
    case Theme.light:
      config.value = {
        icon: LightIcon,
        text: "Switch to dark mode",
      };
      break;
    case Theme.dark:
      config.value = {
        icon: DarkIcon,
        text: "Switch to system preference",
      };
      break;
    default:
      config.value = {
        icon: AutoIcon,
        text: "Switch to light mode",
      };
      break;
  }
};

if (props.listItem) {
  getConfig();
} else {
  watch(() => preferences.theme, getConfig, { immediate: true });
}

const next = () => {
  switch (preferences.theme) {
    case Theme.light:
      preferences.theme = Theme.dark;
      break;
    case Theme.dark:
      preferences.theme = Theme.auto;
      break;
    default:
      preferences.theme = Theme.light;
      break;
  }
};
</script>

<style scoped lang="scss"></style>
