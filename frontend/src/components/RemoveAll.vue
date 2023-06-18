<template>
  <div>
    <v-dialog v-model="show" max-width="400">
      <v-card>
        <v-card-title class="text-h5">Confirm</v-card-title>
        <v-card-text>Delete {{ count }} presets?</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="show = false">
            <v-icon :icon="CloseIcon" aria-hidden="true" />
            Close
          </v-btn>
          <v-btn variant="text" :disabled="countdown > 0" color="red" @click="remove">
            <v-icon :icon="RemoveAllIcon" aria-hidden="true" />
            Delete
            <template v-if="countdown > 0">({{ countdown }})</template>
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
      All presets have been removed.
    </v-snackbar>
  </div>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import CloseIcon from "~icons/solar/close-circle-bold";
import RemoveAllIcon from "~icons/solar/trash-bin-2-bold";
import { wait } from "../util/helpers";
import { usePresetsStore } from "../plugins/store/presets";

let timeout;

const props = defineProps({
  modelValue: Boolean,
});

const emit = defineEmits(["update:modelValue"]);

const show = ref(false);
const showSnackbar = ref(false);
const countdown = ref(0);

const count = computed(() => usePresetsStore().presets.length || 0);

watch(
  () => props.modelValue,
  (val) => {
    show.value = val;
  }
);

const doCountdown = () => {
  countdown.value -= 1;
  if (countdown.value > 0) {
    timeout = setTimeout(doCountdown, 1000);
  }
};

watch(show, (val) => {
  emit("update:modelValue", val);

  if (val) {
    clearTimeout(timeout);
    countdown.value = 5;
    timeout = setTimeout(doCountdown, 1000);
  }
});

const remove = async () => {
  show.value = false;
  await wait(300);
  usePresetsStore().removeAll();
  showSnackbar.value = true;
};
</script>
