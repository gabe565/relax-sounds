import { Sound } from "../util/Sound";
import pb from "../plugins/pocketbase";

let sounds;

export const getSounds = async (force = false) => {
  if (!force && sounds) {
    return sounds;
  }

  const data = await pb.collection("sounds").getFullList({
    fields: "collectionId,id,old_id,name,icon,file,expand.tags.name",
    expand: "tags",
    sort: "name",
  });

  sounds = data.map((sound) => {
    sound.tags = sound.expand?.tags?.map((tag) => tag.name);
    delete sound.expand;
    return new Sound(sound);
  });
  return sounds;
};
