<template>
  <v-card @click.stop="playPause" :light="this.sound.playing" raised outlined>
    <v-list-item>
      <v-list-item-content>
        <v-list-item-title class="headline">
          {{ sound.name }}
        </v-list-item-title>
      </v-list-item-content>
    </v-list-item>
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
      this.sound.playing = true;
      this.sound.player.once('load', () => {
        this.play();
      });
      this.sound.player.load();
    },
    play() {
      this.sound.playing = true;
      this.sound.player.play();
      this.sound.player.fade(0, 1, 1000);
    },
    stop() {
      this.sound.playing = false;
      this.sound.player.once('fade', async () => {
        this.sound.player.stop();
      });
      this.sound.player.fade(this.sound.player.volume(), 0, 500);
    },
  },
};
</script>

<style scoped>

</style>
