import { Howl } from 'howler';

import soundConfig from '../data/sounds.json';

class Sound {
  constructor(obj) {
    Object.assign(this, obj);
    this.state = 'stopped';
    this._volume = 1;
    this.isLoading = false;
    this.howl = new Howl({
      src: [this.src],
      loop: true,
      preload: false,
      volume: 0,
    });
  }

  get src() {
    return `${process.env.BASE_URL}audio/${this.id}.ogg`;
  }

  load() {
    if (this.isUnloaded) {
      this.isLoading = true;
      return new Promise((resolve, reject) => {
        this.howl.once('load', () => {
          this.isLoading = false;
          resolve();
        });
        this.howl.once('loaderror', () => {
          this.isLoading = false;
          reject();
        });
        this.howl.load();
      });
    }
    return true;
  }

  play(local = true, fade = 500) {
    this.state = 'playing';
    if (local) {
      this.howl.play();
      if (fade) {
        this.howl.fade(0, this._volume, fade);
      } else {
        this.howl.volume(this._volume);
      }
    }
  }

  pause(local = true) {
    this.state = 'paused';
    if (local) {
      this.howl.pause();
    }
  }

  stop(local = true, fade = 500) {
    this.state = 'stopped';
    if (local) {
      if (fade) {
        this.howl.once('fade', async () => {
          this.howl.stop();
          this.howl.unload();
        });
        this.howl.fade(this.howl.volume(), 0, fade);
      } else {
        this.howl.stop();
        this.howl.unload();
      }
    }
  }

  set volume(v) {
    this.howl.volume(v);
    this._volume = v;
  }

  get volume() {
    return this._volume;
  }

  get isPlaying() {
    return this.state === 'playing';
  }

  get isStopped() {
    return this.state === 'stopped';
  }

  get isPaused() {
    return this.state === 'paused';
  }

  get isUnloaded() {
    return this.howl.state() === 'unloaded';
  }
}

export const sounds = soundConfig.map((sound) => new Sound(sound)).sort(
  (left, right) => left.name.localeCompare(right.name),
);

export async function prefetch() {
  const cache = await window.caches.open('audio-cache');
  await Promise.all(sounds.map(async (sound) => {
    sound.isLoading = true;
    await cache.add(sound.src);
    sound.isLoading = false;
  }));
}
