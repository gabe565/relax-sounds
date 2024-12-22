<template>
  <div class="d-flex align-center">
    <span v-if="!loading">New content available, click on reload button to update.</span>
    <span v-else>Refreshing. Please wait...</span>
    <div class="flex-grow-1" />
    <button
      class="v-btn v-btn--slim v-btn--density-default v-btn--size-default v-btn--variant-text v-btn--icon"
      :class="{ 'v-btn--disabled': loading }"
      title="Refresh"
      @click="refresh"
    >
      <span class="v-btn__overlay" />
      <span class="v-btn__underlay" />
      <span class="v-btn__content">
        <i class="v-icon v-icon--size-default" aria-hidden="true">
          <refresh-icon />
        </i>
      </span>
    </button>
  </div>
</template>

<script setup>
import { ref } from "vue";
import RefreshIcon from "~icons/material-symbols/refresh-rounded";

const emit = defineEmits(["refresh"]);

const loading = ref(false);

const refresh = () => {
  emit("refresh");
  loading.value = true;
};
</script>
