<template>
  <div class="d-flex align-center">
    <span>Removed {{ presets.presets.length }} presets.</span>
    <div class="flex-grow-1" />
    <button
      class="v-btn v-btn--slim v-btn--density-default v-btn--size-default v-btn--variant-text v-btn--icon"
      title="Undo"
      @click="undo"
    >
      <span class="v-btn__overlay" />
      <span class="v-btn__underlay" />
      <span class="v-btn__content">
        <i class="v-icon v-icon--size-default" aria-hidden="true">
          <undo-icon />
        </i>
      </span>
    </button>
  </div>
</template>

<script setup>
import { usePresetsStore } from "../../../plugins/store/presets";
import UndoIcon from "~icons/material-symbols/undo-rounded";

const emit = defineEmits(["close-toast"]);

const presets = usePresetsStore();

const undo = () => {
  presets.unhideAll();
  emit("close-toast");
};
</script>
