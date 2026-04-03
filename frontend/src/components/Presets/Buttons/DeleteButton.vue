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
import { toast } from "vue-sonner";
import TrashIcon from "~icons/material-symbols/delete-rounded";
import { getErrorMessage } from "@/plugins/store/pocketbase.js";
import { usePresetsStore } from "@/plugins/store/presets";

const props = defineProps({
  preset: {
    type: Object,
    required: true,
  },
});

const presets = usePresetsStore();

const remove = async () => {
  presets.hide({ preset: props.preset });
  const closeHandler = async () => {
    if (props.preset.hidden) {
      try {
        presets.remove({ preset: props.preset });
      } catch (err) {
        console.error("Failed to delete remote preset:", err);
        toast.error(`Failed to remove preset from server.\n${getErrorMessage(err)}`);
      }
    }
  };
  toast.success(`Removed "${props.preset.name}".`, {
    icon: TrashIcon,
    duration: 10000,
    action: {
      label: "Undo",
      onClick: () => presets.unhide({ preset: props.preset }),
    },
    onDismiss: closeHandler,
    onAutoClose: closeHandler,
  });
};
</script>
