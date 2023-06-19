<template>
  <PageLayout :alert="alert" :actions="actions">
    <v-row>
      <v-col v-for="(preset, key) of presets.presets" :key="key" cols="12" md="6" xl="4">
        <PresetCard :preset="preset" />
      </v-col>
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
import PresetCard from "../components/PresetCard.vue";
import PageLayout from "../layouts/PageLayout.vue";
import RestorePresets from "../components/RestorePresets.vue";
import RemoveAll from "../components/RemoveAll.vue";
import { usePlayerStore } from "../plugins/store/player";
import { usePresetsStore } from "../plugins/store/presets";

defineProps({
  alert: {
    type: Object,
    default: null,
  },
});

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
  await usePlayerStore().initSounds();
  await usePresetsStore().migrate();
});
</script>
