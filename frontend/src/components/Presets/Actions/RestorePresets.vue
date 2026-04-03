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
            :prepend-icon="AttachmentIcon"
          />
        </v-card-text>
        <template #actions>
          <v-spacer />
          <v-btn variant="text" @click="show = false">
            <v-icon :icon="CloseIcon" />
            Close
          </v-btn>
          <v-btn variant="text" :disabled="!file" type="submit">
            <v-icon :icon="CheckIcon" />
            Import
          </v-btn>
        </template>
      </v-card>
    </v-form>
  </v-dialog>
</template>

<script setup>
import { nanoid } from "nanoid";
import { ref, watch } from "vue";
import { toast } from "vue-sonner";
import AttachmentIcon from "~icons/material-symbols/attach-file-rounded";
import RestoreIcon from "~icons/material-symbols/backup";
import CheckIcon from "~icons/material-symbols/check-rounded";
import CloseIcon from "~icons/material-symbols/close-rounded";
import { getErrorMessage } from "@/plugins/pocketbase.js";
import { usePresetsStore } from "@/plugins/store/presets";
import { Preset } from "@/util/Preset";

const props = defineProps({
  modelValue: Boolean,
});

const emit = defineEmits(["update:modelValue"]);

const presets = usePresetsStore();
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
    const presetsData = JSON.parse(await file.value.text());
    try {
      await Promise.all(
        presetsData.map(async (data) => {
          const preset = new Preset(data);
          preset.id = nanoid();
          preset.synced = false;
          await preset.migrate();
          await presets.add({ preset, playing: false, sync: false });
        }),
      );
      await presets.sync();
      toast.success(
        `Imported ${presetsData.length} preset${presetsData.length !== 1 ? "s" : ""}.`,
        {
          icon: RestoreIcon,
        },
      );
    } catch (err) {
      console.error(err);
      toast.error(`Failed to import presets:\n${getErrorMessage(err)}`);
    } finally {
      show.value = false;
    }
  } catch (e) {
    console.error(e);
    error.value = true;
  }
};
</script>
