<template>
  <v-card variant="outlined" :dark="showProgress" class="pa-2">
    <v-row align="center" justify="center" dense class="flex-nowrap pl-2">
      <v-col>
        <v-card-title class="text-h5">
          <v-icon aria-hidden="true" class="mr-4" size="x-small" :color="iconColor">
            {{ iconStyle }} {{ sound.icon }} fa-fw
          </v-icon>
          {{ sound.name }}
        </v-card-title>
      </v-col>
      <v-col v-if="showProgress" class="flex-grow-0">
        <v-dialog location="bottom" location-strategy="connected" max-width="400">
          <template #activator="{ props }">
            <v-btn v-bind="props" elevation="0" icon variant="plain" aria-label="Volume">
              <v-icon aria-hidden="true">fas fa-fw fa-volume</v-icon>
            </v-btn>
          </template>

          <v-card class="pa-8">
            <v-slider
              v-model="volumePercentage"
              :min="0"
              :max="100"
              :step="1"
              thumb-label
              color="deep-orange-lighten-1"
              class="pb-1"
              hide-details
            />
          </v-card>
        </v-dialog>
      </v-col>
      <v-col class="flex-grow-0">
        <v-btn
          elevation="0"
          icon
          variant="plain"
          :aria-label="sound.isPlaying ? 'Stop' : 'Play'"
          @click.stop="playStop"
        >
          <v-icon aria-hidden="true">fas fa-fw {{ icon }}</v-icon>
        </v-btn>
      </v-col>
    </v-row>
  </v-card>
</template>

<script>
export default {
  name: "SoundCard",

  props: {
    sound: {
      type: Object,
      required: true,
    },
  },

  computed: {
    volumePercentage: {
      get() {
        return this.sound.volume * 100;
      },
      set(newValue) {
        // eslint-disable-next-line vue/no-mutating-props
        this.sound.volume = newValue / 100;
        this.$store.dispatch("player/updateCast");
      },
    },

    iconStyle() {
      return this.sound.isStopped ? "fal" : "fas";
    },

    iconColor() {
      return this.sound.isStopped ? "" : "primary";
    },

    icon() {
      if (this.sound.isLoading) {
        return "fa-spinner-third fa-spin-2x";
      }
      if (this.sound.isPlaying) {
        return "fa-stop";
      }
      return "fa-play";
    },

    showProgress() {
      return !this.sound.isStopped;
    },
  },

  methods: {
    async playStop() {
      return this.$store.dispatch("player/playStop", { sound: this.sound });
    },
  },
};
</script>
