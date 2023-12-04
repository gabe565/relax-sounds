<template>
  <v-btn elevation="0" icon color="transparent" aria-label="Share" @click.stop="shareOrShow">
    <v-icon :icon="ShareIcon" aria-hidden="true" />
  </v-btn>

  <v-dialog v-model="show" width="400">
    <v-card>
      <v-card-title class="text-h5">Share</v-card-title>
      <v-card-text>
        <v-text-field
          readonly
          :model-value="preset.shareUrl"
          label="Share URL"
          @focus="select"
          @click="select"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="show = false">
          <v-icon aria-hidden="true">$close</v-icon>
          Close
        </v-btn>
        <v-btn variant="text" @click="copy">
          <v-icon :icon="CopyIcon" aria-hidden="true" />
          Copy
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { computed, nextTick, ref } from "vue";
import ShareIcon from "~icons/material-symbols/share";
import CopyIcon from "~icons/material-symbols/content-copy-rounded";
import { Preset } from "../../../util/Preset";
import { toast } from "vue3-toastify";

const props = defineProps({
  preset: {
    type: Preset,
    required: true,
  },
});

const show = ref(false);

const shareData = computed(() => {
  return {
    title: "Relax Sounds",
    text: `Import my Relax Sounds preset called "${props.preset.name}"`,
    url: props.preset.shareUrl,
  };
});

const canShare = computed(() => navigator.canShare && navigator.canShare(shareData.value));

const shareOrShow = async () => {
  if (canShare.value) {
    return share();
  } else {
    show.value = true;
  }
};

const select = async (event) => {
  await nextTick();
  event.target.select();
  event.target.scrollLeft = 0;
};

const copy = async () => {
  try {
    await navigator.clipboard.writeText(props.preset.shareUrl);
    toast.success("Copied to clipboard.", { icon: CopyIcon });
  } catch (err) {
    console.error(err);
    toast.error("Failed to copy to clipboard.");
  } finally {
    show.value = false;
  }
};

const share = async () => {
  try {
    await navigator.share(shareData.value);
    show.value = false;
  } catch (err) {
    console.error(err);
    show.value = true;
  }
};
</script>
