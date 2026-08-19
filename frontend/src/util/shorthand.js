import { decode, encode } from "@msgpack/msgpack";

export const SHORTHAND_VERSION = 3;

const toPercent = (value) => Math.round(value * 100);

const fromPercent = (percent) => percent / 100;

export const encodeShorthand = (tracks, indexById = new Map()) => {
  const rows = tracks.map((track) => {
    const volume = track.volume ?? 1;
    const rate = track.rate ?? 1;
    const pan = track.pan ?? 0;

    const index = indexById.get(track.id);
    const row = [Number.isInteger(index) && index >= 0 ? index : track.id];

    if (pan !== 0) {
      row.push(toPercent(volume), toPercent(rate), toPercent(pan));
    } else if (rate !== 1) {
      row.push(toPercent(volume), toPercent(rate));
    } else if (volume !== 1) {
      row.push(toPercent(volume));
    }
    return row;
  });

  const packed = encode(rows);
  const out = new Uint8Array(packed.length + 1);
  out[0] = SHORTHAND_VERSION;
  out.set(packed, 1);
  return out;
};

export const decodeShorthand = (bytes, idByIndex = new Map()) => {
  if (bytes[0] !== SHORTHAND_VERSION) {
    throw new Error(`unsupported shorthand version: ${bytes[0]}`);
  }

  return decode(bytes.subarray(1)).map((row) => {
    const [ref, volume, rate, pan] = row;

    let id = ref;
    if (typeof ref === "number") {
      id = idByIndex.get(ref);
      if (id === undefined) {
        throw new Error(`unknown sound index: ${ref}`);
      }
    }

    const track = { id };
    if (volume !== undefined) track.volume = fromPercent(volume);
    if (rate !== undefined) track.rate = fromPercent(rate);
    if (pan !== undefined) track.pan = fromPercent(pan);
    return track;
  });
};
