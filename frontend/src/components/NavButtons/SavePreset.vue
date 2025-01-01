<template>
  <v-dialog v-model="showDialog" max-width="500">
    <template #activator="{ props: dialogProps }">
      <v-tooltip v-if="button" text="Save preset" location="bottom">
        <template #activator="{ props: tooltipProps }">
          <v-btn
            v-if="button"
            v-bind="{ ...tooltipProps, ...dialogProps }"
            icon
            title="Save preset"
            :disabled="player.isStopped"
            aria-label="Save preset"
          >
            <v-icon :icon="AddIcon" aria-hidden="true" />
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
      <v-card title="Preset Name">
        <template #text>
          <v-text-field v-model="name" required autofocus label="Name" />
        </template>
        <template #actions>
          <v-spacer />
          <v-btn variant="text" @click="cancel">
            <v-icon aria-hidden="true">$close</v-icon>
            Cancel
          </v-btn>
          <v-btn color="green" variant="text" type="submit">
            <v-icon aria-hidden="true">$complete</v-icon>
            Save Preset
          </v-btn>
        </template>
      </v-card>
    </v-form>
  </v-dialog>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useToast } from "vue-toastification";
import { VBtn, VListItem } from "vuetify/components";
import AddIcon from "~icons/material-symbols/add-circle-rounded";
import SaveIcon from "~icons/material-symbols/save-rounded";
import { usePlayerStore } from "@/plugins/store/player";
import { usePresetsStore } from "@/plugins/store/presets";
import { wait } from "@/util/helpers";

defineProps({
  button: {
    type: Boolean,
    default: false,
  },
});

const toast = useToast();
const showDialog = ref(false);
const name = ref("");
const player = usePlayerStore();
const router = useRouter();

const cancel = () => {
  showDialog.value = false;
  name.value = "";
};

const save = async () => {
  showDialog.value = false;
  await wait(300);

  try {
    usePresetsStore().savePlaying({ name: name.value });
    toast.success(`Saved "${name.value}".`, { icon: SaveIcon });
    name.value = "";
    return router.push({ name: "Presets" });
  } catch (error) {
    console.error(error);
    toast.error(`Failed to save preset:\n${error}`);
  }
};
</script>
