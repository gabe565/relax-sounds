<template>
  <div>
    <v-dialog v-model="show" max-width="500">
      <v-card>
        <v-card-title class="text-h5">Restore</v-card-title>
        <v-card-text>
          If you previously backed up your presets, you can restore them here.
        </v-card-text>
        <v-card-text>
          <v-file-input
            v-model="file"
            accept="application/json"
            :error="error"
            label="Preset File"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="show = false">
            <v-icon aria-hidden="true">$close</v-icon>
            Close
          </v-btn>
          <v-btn variant="text" :disabled="!file" @click="restore">
            <v-icon aria-hidden="true">$complete</v-icon>
            Import
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-snackbar
      v-model="showSnackbar"
      timeout="5000"
      location="bottom"
      content-class="mb-14 mb-md-0"
    >
      Imported {{ imported }} preset{{ imported !== 1 ? "s" : "" }}.
    </v-snackbar>
  </div>
</template>

<script setup>
import { ref, watch } from "vue";
import { Preset } from "../util/Preset";
import { usePresetsStore } from "../plugins/store/presets";

const props = defineProps({
  modelValue: Boolean,
});

const emit = defineEmits(["update:modelValue"]);

const show = ref(false);
const file = ref(null);
const error = ref(false);
const showSnackbar = ref(false);
const imported = ref(0);

watch(
  () => props.modelValue,
  (val) => {
    show.value = val;
  },
);

watch(show, (val) => emit("update:modelValue", val));

const restore = async () => {
  try {
    const presets = JSON.parse(await file.value[0].text());
    await Promise.all(
      presets.map(async (preset) => {
        preset = new Preset(preset);
        await preset.migrate();
        usePresetsStore().add({ preset, playing: false });
      }),
    );
    show.value = false;
    imported.value = presets.length;
    showSnackbar.value = true;
  } catch (e) {
    console.error(e);
    error.value = true;
  }
};
</script>
