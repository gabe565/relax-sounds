<template>
  <v-col class="flex-grow-0">
    <v-btn elevation="0" icon variant="plain" aria-label="Share" @click.stop="show = true">
      <v-icon aria-hidden="true">fas fa-fw fa-share-alt</v-icon>
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
            <v-icon aria-hidden="true">fal fa-times fa-fw</v-icon>
            Close
          </v-btn>
          <v-btn variant="text" @click="copy">
            <v-icon aria-hidden="true">fal fa-copy fa-fw</v-icon>
            Copy
          </v-btn>
          <v-btn v-if="canShare" variant="text" @click="share">
            <v-icon aria-hidden="true">fal fa-share-alt fa-fw</v-icon>
            Share
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-snackbar
      v-model="showSnackbar"
      timeout="5000"
      location="bottom"
      content-class="mb-14 mb-md-0"
    >
      Copied to clipboard.
    </v-snackbar>
  </v-col>
</template>

<script setup>
import { computed, nextTick, ref } from "vue";
import { Preset } from "../../util/Preset";

const props = defineProps({
  preset: {
    type: Preset,
    required: true,
  },
});

const show = ref(false);
const showSnackbar = ref(false);

const shareData = computed(() => {
  return {
    title: "Relax Sounds",
    text: `Import my Relax Sounds preset called "${props.preset.name}"`,
    url: props.preset.shareUrl,
  };
});

const canShare = computed(() => navigator.canShare && navigator.canShare(shareData.value));

const select = async (event) => {
  await nextTick();
  event.target.select();
  event.target.scrollLeft = 0;
};

const copy = async () => {
  await navigator.clipboard.writeText(props.preset.shareUrl);
  if (showSnackbar.value) {
    showSnackbar.value = false;
    await this.$nextTick();
  }
  showSnackbar.value = true;
  show.value = false;
};

const share = async () => {
  await navigator.share(shareData.value);
  show.value = false;
};
</script>
