<template>
  <v-dialog v-model="show" width="400">
    <template #activator="{ props: dialogProps }">
      <v-tooltip text="Share" location="bottom">
        <template #activator="{ props: tooltipProps }">
          <v-btn
            v-bind="{ ...tooltipProps, ...dialogProps }"
            elevation="0"
            icon
            color="transparent"
            aria-label="Share"
            @click.stop="shareOrShow"
          >
            <v-icon :icon="ShareIcon" aria-hidden="true" />
          </v-btn>
        </template>
      </v-tooltip>
    </template>
    <v-card title="Share">
      <template #text>
        <v-text-field
          readonly
          :model-value="url"
          label="Share URL"
          @focus="select"
          @click="select"
        />
      </template>
      <template #actions>
        <v-spacer />
        <v-btn variant="text" @click="show = false">
          <v-icon aria-hidden="true">$close</v-icon>
          Close
        </v-btn>
        <v-btn variant="text" @click="copy">
          <v-icon :icon="CopyIcon" aria-hidden="true" />
          Copy
        </v-btn>
      </template>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { computed, nextTick, ref, watch } from "vue";
import { useToast } from "vue-toastification";
import { Preset } from "../../../util/Preset";
import CopyIcon from "~icons/material-symbols/content-copy-rounded";
import ShareIcon from "~icons/material-symbols/share";

const props = defineProps({
  preset: {
    type: Preset,
    required: true,
  },
});

const toast = useToast();
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
    await navigator.clipboard.writeText(url.value);
    toast.success("Copied to clipboard.", { icon: CopyIcon });
  } catch (err) {
    console.error(err);
    toast.error(`Failed to copy to clipboard:\n${err}`);
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

const url = ref("");
watch(
  props.preset,
  async () => {
    url.value = await props.preset.getShareUrl();
  },
  { immediate: true },
);
</script>
