<template>
  <v-dialog v-model="show" max-width="500">
    <v-card title="Import Preset">
      <v-card-text>
        Add <strong>{{ presets.pendingImport?.name }}</strong> to your presets?
      </v-card-text>
      <v-card-text v-if="sounds.length" class="pt-0">
        <ul class="flex max-h-80 flex-col gap-0.5 overflow-y-auto">
          <li
            v-for="{ sound, volume } in sounds"
            :key="sound.id"
            class="bg-surface-container-high flex min-h-14 items-center gap-3 rounded-md px-4 first:rounded-t-3xl last:rounded-b-3xl"
          >
            <sound-icon
              :icon="sound.icon"
              :size="32"
              :color="sound.tagColor ? `tag-${sound.tagColor}` : 'primary'"
            />
            <span class="grow truncate">{{ sound.name }}</span>
            <span class="text-on-surface-variant text-sm tabular-nums">{{ volume }}%</span>
          </li>
        </ul>
      </v-card-text>
      <template #actions>
        <v-spacer />
        <v-btn variant="text" @click="show = false">
          <v-icon :icon="CloseIcon" />
          Cancel
        </v-btn>
        <v-btn variant="text" :loading="importing" @click="confirm">
          <v-icon :icon="CheckIcon" />
          Import
        </v-btn>
      </template>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { computed, ref } from "vue";
import { toast } from "vue-sonner";
import CheckIcon from "~icons/material-symbols/check-rounded";
import CloseIcon from "~icons/material-symbols/close-rounded";
import SoundIcon from "@/components/Sounds/SoundIcon.vue";
import { usePlayer } from "@/plugins/store/player";
import { getErrorMessage } from "@/plugins/store/pocketbase.js";
import { usePresets } from "@/plugins/store/presets";

const presets = usePresets();
const player = usePlayer();
const importing = ref(false);

const show = computed({
  get: () => !!presets.pendingImport,
  set: (val) => {
    if (!val) presets.pendingImport = null;
  },
});

const sounds = computed(() =>
  (presets.pendingImport?.sounds ?? []).flatMap(({ id, volume }) => {
    const sound = player.soundById(id);
    if (!sound) return [];
    return [{ sound, volume: Math.round((volume ?? 1) * 100) }];
  }),
);

const confirm = async () => {
  const preset = presets.pendingImport;
  importing.value = true;
  try {
    await presets.add({ preset });
    toast.success(`Imported ${preset.name}.`);
  } catch (err) {
    console.error(err);
    toast.error(`Failed to import preset:\n${getErrorMessage(err)}`);
  } finally {
    importing.value = false;
    presets.pendingImport = null;
  }
};
</script>
