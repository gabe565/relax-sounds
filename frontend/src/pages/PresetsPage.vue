<template>
  <PageLayout :actions="actions">
    <v-row>
      <v-fade-transition group leave-absolute>
        <v-col v-for="preset of presets.presets" :key="preset.name" cols="12" md="6" lg="4" xl="3">
          <PresetCard :preset="preset" />
        </v-col>
      </v-fade-transition>
      <v-col v-if="presets.presets.length === 0">
        <v-alert prominent text type="info"> No Presets Saved Yet! </v-alert>
      </v-col>
    </v-row>

    <restore-presets v-model="showRestore" />
    <remove-all v-model="showRemoveAll" />
  </PageLayout>
</template>

<script setup>
import { saveAs } from "file-saver/src/FileSaver";
import { onMounted, ref } from "vue";
import BackupIcon from "~icons/material-symbols/cloud-download-rounded";
import RestoreIcon from "~icons/material-symbols/backup";
import RemoveAllIcon from "~icons/material-symbols/delete-rounded";
import PresetCard from "../components/Presets/PresetCard.vue";
import PageLayout from "../layouts/PageLayout.vue";
import RestorePresets from "../components/Presets/Actions/RestorePresets.vue";
import RemoveAll from "../components/Presets/Actions/RemoveAll.vue";
import { usePlayerStore } from "../plugins/store/player";
import { usePresetsStore } from "../plugins/store/presets";
import { useToast } from "vue-toastification";

const toast = useToast();
const showRestore = ref(false);
const showRemoveAll = ref(false);

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

const actions = [
  {
    title: "Backup",
    icon: BackupIcon,
    on: {
      click: exportPresets,
    },
  },
  {
    title: "Restore",
    icon: RestoreIcon,
    on: {
      click: () => {
        showRestore.value = true;
      },
    },
  },
  {
    title: "Remove All",
    icon: RemoveAllIcon,
    on: {
      click: () => {
        showRemoveAll.value = true;
      },
    },
  },
];

onMounted(async () => {
  try {
    await usePlayerStore().initSounds();
  } catch (err) {
    console.error(err);
    toast.error("Failed to fetch sounds.");
  }
  try {
    await usePresetsStore().migrate();
  } catch (err) {
    console.error(err);
    toast.error("Failed to migrate presets.");
  }
});
</script>
