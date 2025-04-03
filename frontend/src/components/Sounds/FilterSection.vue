<template>
  <v-row class="filters">
    <v-col>
      <v-combobox
        ref="combobox"
        v-model="filters.filters.word"
        label="Search"
        :prepend-inner-icon="SearchIcon"
        :loading="isLoading"
        clearable
        persistent-clear
        rounded
        density="comfortable"
        hide-details
        hide-no-data
        :items="tags"
        :error="error"
        item-title="name"
        item-value="name"
        :return-object="false"
        :menu-icon="DropdownIcon"
        :clear-icon="CloseIcon"
        @keydown.esc="
          filters.filters.word = '';
          combobox.blur();
        "
        @keydown.enter="combobox.blur()"
      >
        <template #item="{ props, item }">
          <v-list-item v-bind="props" :title="item.raw.name">
            <template #prepend>
              <v-icon v-if="item.raw.icon">
                <icon :icon="item.raw.icon" />
              </v-icon>
            </template>
          </v-list-item>
        </template>
      </v-combobox>
    </v-col>
  </v-row>
</template>

<script setup>
import { Icon } from "@iconify/vue";
import { useAsyncState, useMagicKeys } from "@vueuse/core";
import { ref, watch } from "vue";
import { useToast } from "vue-toastification";
import DropdownIcon from "~icons/material-symbols/arrow-drop-down-rounded";
import CloseIcon from "~icons/material-symbols/close-rounded";
import SearchIcon from "~icons/material-symbols/search-rounded";
import { getTags } from "@/data/tags";
import { useFiltersStore } from "@/plugins/store/filters";

const combobox = ref();
const toast = useToast();
const filters = useFiltersStore();

const { Cmd_K, Ctrl_K } = useMagicKeys();
watch([Cmd_K, Ctrl_K], (v) => {
  if (v) {
    combobox.value?.focus();
  }
});

const {
  state: tags,
  isLoading,
  error,
} = useAsyncState(getTags, [], {
  onError(e) {
    toast.error(`Failed to fetch tags:\n${e}`);
  },
});
</script>

<style scoped lang="scss">
.v-input:deep(.v-field__outline) {
  display: none;
}
</style>
