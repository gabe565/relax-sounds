<template>
  <v-row>
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
        color="secondary"
        @keydown.esc="
          filters.filters.word = '';
          combobox.blur();
        "
        @keydown.enter="combobox.blur()"
      >
        <template #item="{ props, item }">
          <v-list-item v-bind="props" :title="item?.name ?? item">
            <template #prepend>
              <v-icon v-if="item?.icon">
                <icon :icon="item.icon" />
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
import { useTemplateRef, watch } from "vue";
import { toast } from "vue-sonner";
import DropdownIcon from "~icons/material-symbols/arrow-drop-down-rounded";
import CloseIcon from "~icons/material-symbols/close-rounded";
import SearchIcon from "~icons/material-symbols/search-rounded";
import { useFilters } from "@/plugins/store/filters";
import { usePocketBase } from "@/plugins/store/pocketbase.js";

const combobox = useTemplateRef("combobox");
const pb = usePocketBase();
const filters = useFilters();

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
} = useAsyncState(() => pb.loadTags(), [], {
  onError(err) {
    console.error(err);
    toast.error(`Failed to fetch tags:\n${err}`);
  },
});
</script>

<style scoped>
.v-input:deep(.v-field__outline) {
  display: none;
}
</style>
