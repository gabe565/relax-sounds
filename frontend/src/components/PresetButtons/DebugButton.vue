<template>
  <v-col class="flex-grow-0">
    <v-btn
      elevation="0"
      icon
      variant="plain"
      aria-label="Debug"
      @click.stop="show = true"
    >
      <v-icon aria-hidden="true">
        fas fa-fw fa-bug
      </v-icon>
    </v-btn>

    <v-dialog
      v-model="show"
      max-width="400"
    >
      <v-card>
        <v-card-title class="text-h5">
          Debug
        </v-card-title>
        <v-card-text>
          <v-btn
            :href="preset.mixUrlAs('mp3')"
            target="_blank"
            class="mr-2"
          >
            Mix MP3
          </v-btn>
          <v-btn
            :href="preset.mixUrlAs('wav')"
            target="_blank"
            class="mr-2"
          >
            Mix WAV
          </v-btn>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            variant="text"
            @click="show = false"
          >
            <v-icon aria-hidden="true">
              fal fa-times fa-fw
            </v-icon>
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-col>
</template>

<script>
import { Preset } from '../../util/Preset';

export default {
  name: 'DebugButton',

  props: {
    preset: {
      type: Preset,
      required: true,
    },
  },

  emits: ['update:modelValue'],

  data: () => ({
    show: false,
  }),

  watch: {
    value: {
      handler(val) {
        this.show = val;
      },
      immediate: true,
    },
    show(val) {
      this.$emit('update:modelValue', val);
    },
  },
};
</script>
