<template>
  <v-tooltip text="Delete" location="bottom">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        elevation="0"
        icon
        color="transparent"
        aria-label="Delete"
        @click.stop="remove"
      >
        <v-icon :icon="TrashIcon" aria-hidden="true" />
      </v-btn>
    </template>
  </v-tooltip>
</template>

<script setup>
import TrashIcon from "~icons/material-symbols/delete-rounded";
import { usePresetsStore } from "../../../plugins/store/presets";
import { useToast } from "vue-toastification";
import DeleteToast from "./DeleteToast.vue";

const props = defineProps({
  preset: {
    type: Object,
    required: true,
  },
});

const presets = usePresetsStore();
const toast = useToast();

const remove = async () => {
  presets.hide({ preset: props.preset });
  toast.success(
    {
      component: DeleteToast,
      props: {
        preset: props.preset,
      },
    },
    {
      icon: TrashIcon,
      timeout: 10000,
      closeOnClick: false,
      onClose: () => {
        if (props.preset.hidden) {
          presets.remove({ preset: props.preset });
        }
      },
    },
  );
};
</script>
