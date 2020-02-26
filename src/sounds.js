import { Howl } from 'howler';

import defaultSounds from './assets/sounds.json';

export default defaultSounds.map((sound) => {
  const src = `/audio/${sound.id}.ogg`;
  return {
    ...sound,
    state: 'stopped',
    volume: 1,
    loading: false,
    src,
    player: new Howl({
      src: [src],
      loop: true,
      preload: false,
      volume: 0,
    }),
  };
}).sort(
  (left, right) => left.name.localeCompare(right.name),
);
