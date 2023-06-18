<template>
  <v-tooltip location="bottom" aria-label="Save Preset">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        :disabled="player.isStopped"
        icon
        aria-label="Save Preset"
        @click="showDialog = !showDialog"
      >
        <v-icon :icon="AddIcon" aria-hidden="true" />
        <v-dialog v-model="showDialog" max-width="500">
          <v-card>
            <v-card-title class="text-h5">Preset Name</v-card-title>
            <v-card-text>
              <v-form @submit.prevent="save">
                <v-text-field v-model="name" required autofocus label="Name" />
              </v-form>
            </v-card-text>
            <v-card-actions>
              <v-spacer />
              <v-btn variant="text" @click="cancel">
                <v-icon :icon="CloseIcon" aria-hidden="true" />
                Cancel
              </v-btn>
              <v-btn color="green" variant="text" @click="save">
                <v-icon :icon="SaveIcon" aria-hidden="true" />
                Save Preset
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-btn>
    </template>
    <span>Save Preset</span>
  </v-tooltip>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import AddIcon from "~icons/solar/add-circle-bold";
import CloseIcon from "~icons/solar/close-circle-bold";
import SaveIcon from "~icons/solar/diskette-bold";
import { wait } from "../util/helpers";
import { usePlayerStore } from "../plugins/store/player";
import { usePresetsStore } from "../plugins/store/presets";
import { useAlertStore } from "../plugins/store/alert";

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

  const alert = useAlertStore();
  try {
    usePresetsStore().savePlaying({ name: name.value });
    alert.type = "info";
    alert.message = `Preset "${name.value}" saved successfully.`;
    name.value = "";
  } catch (error) {
    console.error(error);
    alert.type = "error";
    alert.message = "Failed to save preset. Please try again later.";
  }
  alert.show = true;
  return router.push({ name: "Presets" });
};
</script>
