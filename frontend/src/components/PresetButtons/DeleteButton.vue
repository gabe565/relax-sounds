<template>
  <v-col class="flex-grow-0">
    <v-btn elevation="0" icon variant="plain" aria-label="Share" @click.stop="show = true">
      <v-icon :icon="TrashIcon" aria-hidden="true" />
    </v-btn>

    <v-dialog v-model="show" width="400">
      <v-card>
        <v-card-title class="text-h5">Confirm</v-card-title>
        <v-card-text>Delete the preset "{{ preset.name }}"?</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="show = false">
            <v-icon aria-hidden="true">$close</v-icon>
            Close
          </v-btn>
          <v-btn color="red" variant="text" @click="remove">
            <v-icon :icon="TrashIcon" aria-hidden="true" />
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-col>
</template>

<script setup>
import { ref } from "vue";
import TrashIcon from "~icons/solar/trash-bin-2-bold";
import { wait } from "../../util/helpers";
import { usePresetsStore } from "../../plugins/store/presets";

const props = defineProps({
  preset: {
    type: Object,
    required: true,
  },
});

const show = ref(false);

const remove = async () => {
  show.value = false;
  await wait(300);
  usePresetsStore().remove({ preset: props.preset });
};
</script>
