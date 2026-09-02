<template>
  <div
    class="bg-primary rounded-circle flex shrink-0 items-center justify-center"
    :style="{ width: `${size}px`, height: `${size}px` }"
  >
    <v-icon v-if="isIconifyName" :size="glyphSize" class="text-on-primary">
      <icon :icon="icon" />
    </v-icon>
    <v-icon v-else :icon="icon" :size="glyphSize" class="text-on-primary" />
  </div>
</template>

<script setup>
import { Icon } from "@iconify/vue";
import { computed } from "vue";

const props = defineProps({
  icon: {
    type: [String, Object, Function],
    required: true,
  },
  size: {
    type: Number,
    default: 40,
  },
});

const isIconifyName = computed(() => typeof props.icon === "string");
const glyphSize = computed(() => Math.round(props.size * 0.55));
</script>
