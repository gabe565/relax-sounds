import base64 from 'base64-url';

export function encodeSounds(sounds) {
  return base64.encode(JSON.stringify(
    sounds.map((sound) => [sound.id, Math.round(sound.volume * 1000) / 1000]),
  ));
}

export function encode(preset) {
  const name = encodeURIComponent(preset.name).replace(/%20/g, '+');
  const sounds = encodeSounds(preset.sounds);
  return { name, sounds };
}

export function decodeSounds(param) {
  return JSON.parse(base64.decode(param)).map(
    (song) => ({ id: song[0], volume: song[1] }),
  );
}

export function decode(params) {
  const name = params.name.replace(/\+/g, ' ');
  const sounds = decodeSounds(params.songs);
  return { name, sounds };
}
