<template>
  <page-layout>
    <template #menu>
      <v-list-item title="Backup" :prepend-icon="BackupIcon" @click="exportPresets" />
      <restore-presets />
      <v-list-item title="Remove All" :prepend-icon="RemoveAllIcon" @click="removeAll" />
    </template>

    <v-overlay v-model="isLoading" class="self-center justify-center" persistent>
      <v-progress-circular color="primary" indeterminate size="64" />
    </v-overlay>

    <v-row>
      <v-fade-transition group leave-absolute>
        <template v-if="presets.active.length !== 0">
          <v-col v-for="preset of presets.active" :key="preset.id" cols="12" md="6" lg="4" xl="3">
            <preset-card :preset="preset" />
          </v-col>
        </template>
        <v-col v-else>
          <v-alert prominent type="info" :icon="InfoIcon">No presets saved yet</v-alert>
        </v-col>
      </v-fade-transition>
    </v-row>
  </page-layout>
</template>

<script setup>
import { saveAs } from "file-saver/src/FileSaver";
import { onActivated, ref } from "vue";
import { toast } from "vue-sonner";
import BackupIcon from "~icons/material-symbols/cloud-download-rounded";
import RemoveAllIcon from "~icons/material-symbols/delete-rounded";
import InfoIcon from "~icons/material-symbols/info-rounded";
import RestorePresets from "@/components/Presets/Actions/RestorePresets.vue";
import PresetCard from "@/components/Presets/PresetCard.vue";
import PageLayout from "@/layouts/PageLayout.vue";
import { usePlayerStore } from "@/plugins/store/player";
import { getErrorMessage } from "@/plugins/store/pocketbase.js";
import { usePresetsStore } from "@/plugins/store/presets";

const player = usePlayerStore();
const presets = usePresetsStore();

onActivated(async () => {
  try {
    await presets.sync();
  } catch (err) {
    console.error("Failed to sync presets:", err);
    toast.error(`Failed to sync presets:\n${getErrorMessage(err)}`);
  }
});

const exportPresets = () => {
  const blob = new Blob([JSON.stringify(presets.presets)], {
    type: "application/json;charset=utf-8",
  });
  const offset = new Date().getTimezoneOffset() * 60000; // Offset in milliseconds
  const localISOTime = new Date(Date.now() - offset)
    .toISOString()
    .slice(0, -5) // Remove ".000Z"
    .replaceAll(":", "");
  saveAs(blob, `relax-sounds-presets-${localISOTime}.json`);
  toast.success(`Downloaded ${presets.presets.length} presets.`, { icon: BackupIcon });
};

const isLoading = ref(true);

(async () => {
  try {
    await player.loadSounds();
    await presets.migrate();
  } catch (err) {
    console.error(err);
    toast.error(`Failed to load:\n${getErrorMessage(err)}`);
  } finally {
    isLoading.value = false;
  }
})();

const removeAll = () => {
  const count = presets.presets.length;
  presets.hideAll();
  const closeHandler = async () => {
    try {
      await presets.removeHidden();
    } catch (err) {
      console.error("Failed to remove hidden presets:", err);
      toast.error(`Failed to remove all presets:\n${getErrorMessage(err)}`);
    }
  };
  toast.success(`Removed ${count} presets.`, {
    icon: RemoveAllIcon,
    duration: 10000,
    action: {
      label: "Undo",
      onClick: () => presets.unhideAll(),
    },
    onDismiss: closeHandler,
    onAutoClose: closeHandler,
  });
};
</script>
