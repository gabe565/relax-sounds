import { Sound } from "../util/Sound";
import pb from "../plugins/pocketbase";

let sounds;

export const getSounds = async (force = false) => {
  if (!force && sounds) {
    return sounds;
  }

  const data = await pb.collection("sounds").getFullList({ expand: "tags" });
  sounds = data
    .map((sound) => ({
      ...sound,
      tags: sound.expand.tags?.map((tag) => tag.name),
    }))
    .sort((left, right) => left.name.localeCompare(right.name))
    .map((sound) => new Sound(sound));
  return sounds;
};
