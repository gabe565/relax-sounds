<template>
  <v-dialog v-model="show" max-width="500">
    <template #activator="{ props }">
      <v-list-item title="Restore" :prepend-icon="RestoreIcon" v-bind="props" />
    </template>

    <v-form @submit.prevent="restore">
      <v-card title="Restore">
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
        <template #actions>
          <v-spacer />
          <v-btn variant="text" @click="show = false">
            <v-icon aria-hidden="true">$close</v-icon>
            Close
          </v-btn>
          <v-btn variant="text" :disabled="!file" type="submit">
            <v-icon aria-hidden="true">$complete</v-icon>
            Import
          </v-btn>
        </template>
      </v-card>
    </v-form>
  </v-dialog>
</template>

<script setup>
import { ref, watch } from "vue";
import { Preset } from "../../../util/Preset";
import { usePresetsStore } from "../../../plugins/store/presets";
import { useToast } from "vue-toastification";
import RestoreIcon from "~icons/material-symbols/backup";

const props = defineProps({
  modelValue: Boolean,
});

const emit = defineEmits(["update:modelValue"]);

const toast = useToast();
const show = ref(false);
const file = ref(null);
const error = ref(false);

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
    try {
      await Promise.all(
        presets.map(async (preset) => {
          preset = new Preset(preset);
          await preset.migrate();
          usePresetsStore().add({ preset, playing: false });
        }),
      );
      toast.success(`Imported ${presets.length} preset${presets.length !== 1 ? "s" : ""}.`, {
        icon: RestoreIcon,
      });
    } catch (err) {
      console.error(err);
      toast.error("Failed to import presets.");
    } finally {
      show.value = false;
    }
  } catch (e) {
    console.error(e);
    error.value = true;
  }
};
</script>
