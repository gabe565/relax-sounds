import axios from 'axios';
import tagFile from './tags.json';

let tags;

export const getTags = async () => {
  if (tags) {
    return tags;
  }

  const { data } = await axios.get(tagFile);
  tags = data;
  return tags;
};
