import base64 from "base64-url";
import { Filetype } from "./filetype";

export const toShorthand = (sounds) =>
  sounds.map((sound) => [sound.id, Math.round(sound.volume * 1000) / 1000]);

export const fromShorthand = (shorthand) =>
  shorthand.map((song) => ({ id: song[0], volume: song[1] }));

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
    return `${window.location.origin}/api/mix/${this.encodedShorthand}.${filetype}`;
  }

  get mixUrl() {
    return this.mixUrlAs(Filetype.Mp3);
  }

  set mixUrl(val) {
    [, this.encodedShorthand] = val.match(/\/api\/mix\/(.+?)(\..+)?$/);
  }
}
