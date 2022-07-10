import axios from 'axios';
import { Sound } from '../util/Sound';

let sounds;

export const getSounds = async () => {
  if (sounds) {
    return sounds;
  }

  let { data } = await axios.get('/api/sounds');
  data = data
    .sort((left, right) => left.name.localeCompare(right.name))
    .map((sound) => new Sound(sound));
  sounds = data;
  return sounds;
};

export const prefetch = async () => {
  const cache = await window.caches.open('data-cache');
  const soundConfig = await getSounds();
  await Promise.all(soundConfig.map(async (sound) => {
    sound.isLoading = true;
    await cache.add(sound.src);
    sound.isLoading = false;
  }));
};
