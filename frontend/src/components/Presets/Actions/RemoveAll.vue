<template>
  <div>
    <v-dialog v-model="show" max-width="400">
      <v-card>
        <v-card-title class="text-h5">Confirm</v-card-title>
        <v-card-text>Delete {{ count }} presets?</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="show = false">
            <v-icon aria-hidden="true">$close</v-icon>
            Close
          </v-btn>
          <v-btn variant="text" :disabled="countdown > 0" color="red" @click="remove">
            <v-icon aria-hidden="true">$complete</v-icon>
            Delete
            <template v-if="countdown > 0">({{ countdown }})</template>
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import { usePresetsStore } from "../../../plugins/store/presets";
import { useToast } from "vue-toastification";
import { wait } from "../../../util/helpers";
import TrashIcon from "~icons/material-symbols/delete-rounded";

let timeout;

const props = defineProps({
  modelValue: Boolean,
});

const emit = defineEmits(["update:modelValue"]);

const toast = useToast();
const presets = usePresetsStore();
const show = ref(false);
const countdown = ref(0);
const count = computed(() => presets.presets.length || 0);

watch(
  () => props.modelValue,
  (val) => {
    show.value = val;
  },
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
  const prevCount = count.value;
  presets.removeAll();
  toast.success(`Removed ${prevCount} presets.`, { icon: TrashIcon });
};
</script>
