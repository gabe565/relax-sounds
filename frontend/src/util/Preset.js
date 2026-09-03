import { nanoid } from "nanoid";
import { ApiPath } from "@/config/api";
import { usePlayer } from "@/plugins/store/player.js";
import { usePocketBase } from "@/plugins/store/pocketbase.js";
import { Filetype } from "@/util/filetype";
import { compress, decompress, fromBase64, fromUrlSafeBase64, toBase64 } from "@/util/helpers";
import { SHORTHAND_VERSION, decodeShorthand, encodeShorthand } from "@/util/shorthand";

export const legacyFromShorthand = (shorthand) =>
  shorthand.map((song) => {
    const entry = { id: song[0], volume: song[1] };
    if (song.length === 3) {
      entry.rate = song[3];
    }
    return entry;
  });

const soundIndex = async () => {
  const sounds = await usePocketBase().loadSounds();
  const indexById = new Map();
  const idByIndex = new Map();
  for (const sound of sounds) {
    if (!Number.isInteger(sound.short_id)) continue;
    indexById.set(sound.id, sound.short_id);
    idByIndex.set(sound.short_id, sound.id);
  }
  return { indexById, idByIndex };
};

const streamId = ({ fresh = false } = {}) => {
  if (fresh) return nanoid(6);
  let uuid = sessionStorage.getItem("streamId");
  if (!uuid) {
    uuid = nanoid(6);
    sessionStorage.setItem("streamId", uuid);
  }
  return uuid;
};

export class Preset {
  constructor(obj) {
    this.id = nanoid();
    this.name = "Unnamed Preset";
    this.sounds = [];
    this.new = false;
    this.synced = false;
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
    return encodeURIComponent(this.name).replaceAll("%20", "+");
  }

  set encodedName(val) {
    this.name = val.replaceAll("+", " ");
  }

  async getEncodedShorthand() {
    try {
      const { indexById } = await soundIndex();
      return toBase64(encodeShorthand(this.shorthand, indexById));
    } catch (err) {
      console.warn("Falling back to compressed JSON preset encoding.", err);
      return await compress(JSON.stringify(this.shorthand));
    }
  }

  async setEncodedShorthand(val) {
    const bytes = fromBase64(val);
    if (bytes[0] === SHORTHAND_VERSION) {
      const { idByIndex } = await soundIndex();
      this.shorthand = decodeShorthand(bytes, idByIndex);
      return;
    }

    try {
      this.shorthand = JSON.parse(await decompress(val));
    } catch {
      this.shorthand = legacyFromShorthand(JSON.parse(atob(fromUrlSafeBase64(val))));
    }
  }

  async getShareUrl() {
    const shorthand = await this.getEncodedShorthand();
    return `${globalThis.location.origin}/import/${this.encodedName}/${shorthand}`;
  }

  async mixUrlAs(filetype = Filetype.Mp3, { fresh = false } = {}) {
    const encoded = await this.getEncodedShorthand();
    return ApiPath(`/api/mix/${streamId({ fresh })}/${encoded}.${filetype}`);
  }

  async hlsUrl({ fresh = false } = {}) {
    const encoded = await this.getEncodedShorthand();
    return ApiPath(`/api/mix/${streamId({ fresh })}/${encoded}.m3u8`);
  }

  async setMixUrl(val) {
    const [, encoded] = val.match(/\/api\/mix\/[^/]+\/([^./]+)/);
    await this.setEncodedShorthand(encoded);
  }

  async migrate() {
    const pb = usePocketBase();
    const player = usePlayer();
    await pb.loadSounds();
    await Promise.all(
      this.sounds.map(async (sound) => {
        if (sound.id.length <= 3) {
          const found = player.sounds.find((e) => `${e.short_id}` === sound.id);
          if (found) {
            sound.id = found.id;
          }
        }
      }),
    );
  }
}
