<template>
  <page-layout>
    <template #actions>
      <cast-icon v-if="isMobile" button />
    </template>
    <template #menu>
      <v-list-item title="Backup" :prepend-icon="BackupIcon" @click="exportPresets" />
      <restore-presets />
      <v-list-item title="Remove All" :prepend-icon="RemoveAllIcon" @click="removeAll" />
    </template>

    <v-overlay v-model="isLoading" class="align-center justify-center" persistent>
      <v-progress-circular color="primary" indeterminate size="64" />
    </v-overlay>

    <v-row>
      <v-fade-transition group leave-absolute>
        <template v-if="presets.active.length !== 0">
          <v-col v-for="preset of presets.active" :key="preset.name" cols="12" md="6" lg="4" xl="3">
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
import { useAsyncState } from "@vueuse/core";
import { saveAs } from "file-saver/src/FileSaver";
import { useToast } from "vue-toastification";
import { useDisplay } from "vuetify";
import BackupIcon from "~icons/material-symbols/cloud-download-rounded";
import RemoveAllIcon from "~icons/material-symbols/delete-rounded";
import InfoIcon from "~icons/material-symbols/info-rounded";
import CastIcon from "@/components/NavButtons/CastIcon.vue";
import RemoveAllToast from "@/components/Presets/Actions/RemoveAllToast.vue";
import RestorePresets from "@/components/Presets/Actions/RestorePresets.vue";
import PresetCard from "@/components/Presets/PresetCard.vue";
import PageLayout from "@/layouts/PageLayout.vue";
import { usePlayerStore } from "@/plugins/store/player";
import { usePresetsStore } from "@/plugins/store/presets";

const toast = useToast();
const { smAndDown: isMobile } = useDisplay();
const presets = usePresetsStore();

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

const { isLoading } = useAsyncState(
  async () => {
    await usePlayerStore().initSounds();
    await usePresetsStore().migrate();
  },
  undefined,
  {
    onError(e) {
      toast.error(`Failed to load: ${e}`);
    },
  },
);

const removeAll = () => {
  presets.hideAll();
  toast.success(RemoveAllToast, {
    icon: RemoveAllIcon,
    timeout: 10000,
    closeOnClick: false,
    onClose: presets.removeHidden,
  });
};
</script>
