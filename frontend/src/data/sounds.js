import axios from "axios";
import { Sound } from "../util/Sound";

let sounds;

export const getSounds = async (force = false) => {
  if (!force && sounds) {
    return sounds;
  }

  let { data } = await axios.get("/api/sounds");
  data = data
    .sort((left, right) => left.name.localeCompare(right.name))
    .map((sound) => new Sound(sound));
  sounds = data;
  return sounds;
};
