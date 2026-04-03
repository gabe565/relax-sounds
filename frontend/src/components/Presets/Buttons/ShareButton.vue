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
            <v-icon :icon="ShareIcon" />
          </v-btn>
        </template>
      </v-tooltip>
    </template>

    <v-card color="cardBackground" variant="flat">
      <v-card-title class="pt-6 px-6">Share</v-card-title>
      <v-card-text>
        <v-text-field
          readonly
          :model-value="url"
          label="Share URL"
          hide-details
          @focus="select"
          @click="select"
        />
      </v-card-text>
      <v-card-actions class="mr-4 mb-4">
        <v-btn variant="text" @click="show = false">Close</v-btn>
        <v-btn color="primary" variant="flat" @click="copy">Copy</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { computed, nextTick, ref, watch } from "vue";
import { toast } from "vue-sonner";
import CopyIcon from "~icons/material-symbols/content-copy-rounded";
import ShareIcon from "~icons/material-symbols/share";
import { Preset } from "@/util/Preset";

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
