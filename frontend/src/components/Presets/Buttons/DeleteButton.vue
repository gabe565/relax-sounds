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
        <v-icon :icon="TrashIcon" />
      </v-btn>
    </template>
  </v-tooltip>
</template>

<script setup>
import { useToast } from "vue-toastification";
import TrashIcon from "~icons/material-symbols/delete-rounded";
import DeleteToast from "@/components/Presets/Buttons/DeleteToast.vue";
import { usePresetsStore } from "@/plugins/store/presets";

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
