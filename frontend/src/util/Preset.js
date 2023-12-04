import base64 from "base64-url";
import { nanoid } from "nanoid";
import { Filetype } from "./filetype";
import { getSounds } from "../data/sounds";

export const toShorthand = (sounds) => {
  return sounds.map((sound) => {
    const entry = [sound.id, Math.round(sound.volume * 100) / 100];
    if (sound.rate && sound.rate !== 1) {
      entry.push(sound.rate);
    }
    return entry;
  });
};

export const fromShorthand = (shorthand) =>
  shorthand.map((song) => {
    const entry = { id: song[0], volume: song[1] };
    if (song.length === 3) {
      entry.rate = song[3];
    }
    return entry;
  });

export class Preset {
  constructor(obj) {
    this.name = "Unnamed Preset";
    this.sounds = [];
    this.new = false;
    Object.assign(this, obj);
  }

  get shorthand() {
    return toShorthand(this.sounds);
  }

  set shorthand(val) {
    this.sounds = fromShorthand(val);
  }

  get encodedName() {
    return encodeURIComponent(this.name).replace(/%20/g, "+");
  }

  set encodedName(val) {
    this.name = val.replace(/\+/g, " ");
  }

  get encodedShorthand() {
    return base64.encode(JSON.stringify(this.shorthand));
  }

  set encodedShorthand(val) {
    this.shorthand = JSON.parse(base64.decode(val));
  }

  get shareUrl() {
    return `${window.location.origin}/import/${this.encodedName}/${this.encodedShorthand}`;
  }

  mixUrlAs(filetype = Filetype.Mp3) {
    let uuid = sessionStorage.getItem("uuid");
    if (!uuid) {
      uuid = nanoid();
      sessionStorage.setItem("uuid", uuid);
    }
    const apiAddress = import.meta.env.VITE_API_ADDRESS || window.location.origin;
    return `${apiAddress}/api/mix/${uuid}/${this.encodedShorthand}.${filetype}`;
  }

  get mixUrl() {
    return this.mixUrlAs(Filetype.Mp3);
  }

  set mixUrl(val) {
    [, this.encodedShorthand] = val.match(/\/api\/mix\/.+?\/(.+?)(\..+)?$/);
  }

  async migrate() {
    await Promise.all(
      this.sounds.map(async (sound) => {
        if (sound.id.length <= 3) {
          const sounds = await getSounds();
          const found = sounds.find((e) => `${e.old_id}` === sound.id);
          if (found) {
            sound.id = found.id;
          }
        }
      }),
    );
  }
}
