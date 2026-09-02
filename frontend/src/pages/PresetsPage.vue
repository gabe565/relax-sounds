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

    <div
      v-if="presets.active.length !== 0"
      class="card-group card-cols-3 grid grid-cols-1 gap-0.5 sm:gap-2 sm:grid-cols-2 lg:grid-cols-3"
    >
      <v-fade-transition group leave-absolute>
        <div v-for="preset of presets.active" :key="preset.id">
          <preset-card :preset="preset" />
        </div>
      </v-fade-transition>
    </div>

    <v-row v-else>
      <v-fade-transition group leave-absolute>
        <v-col>
          <div class="text-on-surface-variant flex flex-col items-center gap-3 py-16 text-center">
            <v-icon :icon="EmptyIcon" size="48" class="opacity-60" />
            <span>No presets saved yet</span>
            <span class="text-sm opacity-70">
              Start some sounds, then save the mix from the player bar.
            </span>
            <v-btn to="/sounds" variant="tonal" rounded="pill">Browse sounds</v-btn>
          </div>
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
import EmptyIcon from "~icons/material-symbols/playlist-play-rounded";
import RestorePresets from "@/components/Presets/Actions/RestorePresets.vue";
import PresetCard from "@/components/Presets/PresetCard.vue";
import PageLayout from "@/layouts/PageLayout.vue";
import { getErrorMessage, usePocketBase } from "@/plugins/store/pocketbase.js";
import { usePresets } from "@/plugins/store/presets";

const pb = usePocketBase();
const presets = usePresets();

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
    await pb.loadSounds();
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
