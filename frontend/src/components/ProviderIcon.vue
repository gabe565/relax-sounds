<template>
  <v-img :src="src" :cover="false" />
</template>

<script setup>
import { computed } from "vue";
import { useTheme } from "vuetify";

const props = defineProps({
  icon: { type: [String, Object], required: true },
});

const theme = useTheme();

const src = computed(() => {
  if (typeof props.icon === "string") {
    return props.icon;
  }
  return theme.current.value.dark ? props.icon.dark : props.icon.light;
});
</script>
