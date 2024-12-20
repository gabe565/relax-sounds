import { nanoid } from "nanoid";
import { Filetype } from "./filetype";
import { getSounds } from "../data/sounds";
import { compress, decompress } from "./helpers";
import base64 from "base64-url";
import { ApiPath } from "../config/api";

export const legacyFromShorthand = (shorthand) =>
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
    return this.sounds.map((sound) => {
      const entry = { id: sound.id };
      if (sound.volume !== 1) {
        entry.volume = sound.volume;
      }
      if (sound.rate && sound.rate !== 1) {
        entry.rate = sound.rate;
      }
      if (sound.pan !== 0) {
        entry.pan = sound.pan;
      }
      return entry;
    });
  }

  set shorthand(val) {
    this.sounds = val;
  }

  get encodedName() {
    return encodeURIComponent(this.name).replace(/%20/g, "+");
  }

  set encodedName(val) {
    this.name = val.replace(/\+/g, " ");
  }

  async getEncodedShorthand() {
    return await compress(JSON.stringify(this.shorthand));
  }

  async setEncodedShorthand(val) {
    let raw;
    try {
      raw = await decompress(val);
      this.shorthand = JSON.parse(raw);
    } catch {
      this.shorthand = legacyFromShorthand(JSON.parse(base64.decode(val)));
    }
  }

  async getShareUrl() {
    const shorthand = await this.getEncodedShorthand();
    return `${window.location.origin}/import/${this.encodedName}/${shorthand}`;
  }

  async mixUrlAs(filetype = Filetype.Mp3) {
    let uuid = sessionStorage.getItem("uuid");
    if (!uuid) {
      uuid = nanoid();
      sessionStorage.setItem("uuid", uuid);
    }
    const encoded = await this.getEncodedShorthand();
    return ApiPath(`/api/mix/${uuid}/${encoded}.${filetype}`);
  }

  async setMixUrl(val) {
    const [, encoded] = val.match(/\/api\/mix\/.+?\/(.+?)(\..+)?$/);
    await this.setEncodedShorthand(encoded);
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
