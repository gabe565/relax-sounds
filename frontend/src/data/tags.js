import pb from "../plugins/pocketbase";

let tags;

export const getTags = async (force = false) => {
  if (!force && tags) {
    return tags;
  }

  tags = await pb.collection("tags").getFullList({
    fields: "icon,name",
  });
  return tags;
};
