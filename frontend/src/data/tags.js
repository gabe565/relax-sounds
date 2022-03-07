import axios from 'axios';

let tags;

export const getTags = async () => {
  if (tags) {
    return tags;
  }

  const { data } = await axios.get('/api/tags');
  tags = data;
  return tags;
};
