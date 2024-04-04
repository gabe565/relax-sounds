import { Howl } from "howler";
import pb from "../plugins/pocketbase";

export const SoundState = {
  PLAYING: "playing",
  PAUSED: "paused",
  STOPPED: "stopped",
  UNLOADED: "unloaded",
};

export class Sound {
  constructor(obj) {
    Object.assign(this, obj);
    this.state = SoundState.STOPPED;
    this._volume = 1;
    this._rate = 1;
    this._pan = 0;
    this.isLoading = false;
    this.howl = new Howl({
      src: this.src,
      loop: true,
      preload: false,
      volume: 0,
    });
  }

  get src() {
    return this.file.map((e) => pb.getFileUrl(this, e));
  }

  load() {
    if (this.isUnloaded) {
      this.isLoading = true;
      return new Promise((resolve, reject) => {
        this.howl.once("load", () => {
          this.isLoading = false;
          resolve();
        });
        this.howl.once("loaderror", (_, err) => {
          this.isLoading = false;
          reject(err);
        });
        this.howl.load();
      });
    }
    return true;
  }

  play(local = true, fade = 250) {
    this.state = SoundState.PLAYING;
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
    this.state = SoundState.PAUSED;
    if (local) {
      this.howl.pause();
    }
  }

  stop(local = true, fade = 250) {
    this.state = SoundState.STOPPED;
    this._volume = 1;
    if (local) {
      if (fade) {
        this.howl.once("fade", () => {
          this.howl.stop();
          this.rate = 1;
          this.pan = 0;
          this.howl.unload();
        });
        this.howl.fade(this.howl.volume(), 0, fade);
      } else {
        this.howl.stop();
        this.rate = 1;
        this.pan = 0;
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

  set rate(v) {
    this.howl.rate(v);
    this._rate = v;
  }

  get rate() {
    return this._rate;
  }

  set pan(v) {
    this.howl.stereo(v);
    this._pan = v;
  }

  get pan() {
    return this._pan;
  }

  get isPlaying() {
    return this.state === SoundState.PLAYING;
  }

  get isStopped() {
    return this.state === SoundState.STOPPED;
  }

  get isPaused() {
    return this.state === SoundState.PAUSED;
  }

  get isUnloaded() {
    return this.howl.state() === SoundState.UNLOADED;
  }
}
