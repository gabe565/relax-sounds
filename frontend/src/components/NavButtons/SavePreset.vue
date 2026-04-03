<template>
  <v-dialog v-model="showDialog" max-width="500">
    <template #activator="{ props: dialogProps }">
      <v-tooltip v-if="button" text="Save preset" location="top">
        <template #activator="{ props: tooltipProps }">
          <v-btn
            v-bind="{ ...tooltipProps, ...dialogProps }"
            icon
            title="Save preset"
            :disabled="player.isStopped"
            aria-label="Save preset"
          >
            <v-icon :icon="AddIcon" />
          </v-btn>
        </template>
      </v-tooltip>
      <v-list-item
        v-else
        v-bind="dialogProps"
        :prepend-icon="AddIcon"
        title="Save Preset"
        :disabled="player.isStopped"
        aria-label="Save Preset"
      />
    </template>

    <v-form @submit.prevent="save">
      <v-card color="card-background" variant="flat">
        <v-card-title class="pt-6 px-6">Save Preset</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="name"
            required
            autofocus
            label="Name"
            variant="outlined"
            density="comfortable"
          />
        </v-card-text>
        <v-card-actions class="mr-4 mb-4">
          <v-btn variant="text" class="mr-2" @click="cancel">Cancel</v-btn>
          <v-btn type="submit" color="primary" variant="flat">Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-form>
  </v-dialog>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { toast } from "vue-sonner";
import AddIcon from "~icons/material-symbols/add-sharp";
import SaveIcon from "~icons/material-symbols/save-sharp";
import { usePlayer } from "@/plugins/store/player";
import { getErrorMessage } from "@/plugins/store/pocketbase.js";
import { usePresets } from "@/plugins/store/presets";
import { wait } from "@/util/helpers";

defineProps({
  button: {
    type: Boolean,
    default: false,
  },
});

const showDialog = ref(false);
const name = ref("");
const player = usePlayer();
const router = useRouter();
const presets = usePresets();

const cancel = () => {
  showDialog.value = false;
  name.value = "";
};

const save = async () => {
  showDialog.value = false;
  await wait(300);

  try {
    await presets.savePlaying({ name: name.value });
    toast.success(`Saved "${name.value}".`, { icon: SaveIcon });
    name.value = "";
    return router.push({ name: "Presets" });
  } catch (err) {
    console.error(err);
    toast.error(`Failed to save preset:\n${getErrorMessage(err)}`);
  }
};
</script>
