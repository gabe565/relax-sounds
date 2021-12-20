import base64 from 'base64-url';

export function encode(playlist) {
  const name = encodeURIComponent(playlist.name).replace(/%20/g, '+');
  const sounds = base64.encode(JSON.stringify(playlist.sounds.reduce((acc, curr) => {
    acc.push([curr.id, Math.round(curr.volume * 1000) / 1000]);
    return acc;
  }, [])));

  return { name, sounds };
}

export function decode(params) {
  const name = params.name.replace(/\+/g, ' ');
  const sounds = JSON.parse(base64.decode(params.songs)).map(
    (song) => ({ id: song[0], volume: song[1] }),
  );

  return { name, sounds };
}
