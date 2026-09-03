<template>
  <div
    class="rounded-circle flex shrink-0 items-center justify-center"
    :class="`bg-${color}`"
    :style="{ width: `${size}px`, height: `${size}px` }"
  >
    <v-icon v-if="isIconifyName" :size="glyphSize">
      <icon :icon="icon" />
    </v-icon>
    <v-icon v-else :icon="icon" :size="glyphSize" />
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
  color: {
    type: String,
    default: "primary",
  },
});

const isIconifyName = computed(() => typeof props.icon === "string");
const glyphSize = computed(() => Math.round(props.size * 0.55));
</script>
