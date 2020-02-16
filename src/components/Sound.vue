<template>
  <v-card raised
          outlined
  >
    <v-progress-linear
      v-model="volumePercentage"
      absolute
      height="100%"
      color="secondary"
      v-if="sound.state !== 'stopped'"
    />
    <v-row align="center" justify="center" dense>
      <v-col class="grow">
        <v-card-title class="headline">
          <v-icon class="pr-2">{{ sound.icon }}</v-icon>
          {{ sound.name }}
        </v-card-title>
      </v-col>
      <v-col class="shrink pr-4">
        <v-btn @click.stop="playPause" elevation="0" outlined :loading="loading">
          <v-icon>{{ icon }}</v-icon>
        </v-btn>
      </v-col>
    </v-row>
  </v-card>
</template>

<script>
import { Howl } from 'howler';

export default {
  name: 'Sound',

  props: {
    sound: {
      type: Object,
    },
  },

  data: () => ({
    loading: false,
  }),

  computed: {
    volume() {
      return this.sound.volume;
    },
    state() {
      return this.sound.state;
    },
    volumePercentage: {
      get() {
        return this.sound.volume * 100;
      },
      set(newValue) {
        this.sound.volume = newValue / 100;
      },
    },
    icon() {
      return this.sound.state === 'playing' ? 'mdi-stop' : 'mdi-play';
    },
  },

  created() {
    if (!this.sound.player) {
      this.sound.player = new Howl({
        src: [`/audio/${this.sound.id}.ogg`],
        loop: true,
        preload: false,
        volume: 0,
      });
    }
  },

  methods: {
    async playPause() {
      if (this.sound.player.playing()) {
        this.stop();
      } else if (this.sound.player.state() === 'loaded') {
        this.play();
      } else {
        this.load();
      }
    },
    load() {
      this.loading = true;
      this.sound.player.once('load', () => {
        this.sound.state = 'playing';
        this.loading = false;
        this.play();
      });
      this.sound.player.load();
    },
    play(fade = 500) {
      this.sound.state = 'playing';
      this.sound.player.play();
      if (fade) {
        this.sound.player.fade(0, this.sound.volume, fade);
      }
      this.$emit('play');
    },
    stop() {
      this.sound.state = 'stopped';
      this.sound.player.once('fade', async () => {
        this.sound.player.stop();
      });
      this.sound.player.fade(this.sound.player.volume(), 0, 500);
    },
    pause() {
      this.sound.player.pause();
    },
  },

  watch: {
    volume(newValue) {
      this.sound.player.volume(newValue);
    },
  },
};
</script>

<style scoped>
  .v-card {
    overflow: hidden;
  }
  .v-progress-linear {
    z-index: 0;
    left: 0;
    right: 0;
    height: 100%;
    width: 100%;
  }
  .row {
    position: relative;
    z-index: 1;
    pointer-events: none;
  }
  .v-btn {
    pointer-events: all;
  }
</style>
