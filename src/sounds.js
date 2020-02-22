import { Howl } from 'howler';

import defaultSounds from './assets/sounds.json';

export default defaultSounds.map((sound) => ({
  ...sound,
  state: 'stopped',
  volume: 1,
  loading: false,
  player: new Howl({
    src: [`/audio/${sound.id}.ogg`],
    loop: true,
    preload: false,
    volume: 0,
  }),
})).sort(
  (left, right) => left.name.localeCompare(right.name),
);
