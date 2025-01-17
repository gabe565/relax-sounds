<template>
  <v-row class="filters">
    <v-col>
      <v-combobox
        ref="combobox"
        v-model="filters.filters.word"
        label="Search"
        :prepend-inner-icon="SearchIcon"
        :loading="loading"
        clearable
        persistent-clear
        rounded
        density="comfortable"
        hide-details
        hide-no-data
        :items="tags"
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
import { useMagicKeys } from "@vueuse/core";
import { onMounted, ref, watch } from "vue";
import { useToast } from "vue-toastification";
import DropdownIcon from "~icons/material-symbols/arrow-drop-down-rounded";
import CloseIcon from "~icons/material-symbols/close-rounded";
import SearchIcon from "~icons/material-symbols/search-rounded";
import { getTags } from "@/data/tags";
import { useFiltersStore } from "@/plugins/store/filters";

const combobox = ref(null);
const tags = ref([]);
const loading = ref(true);
const toast = useToast();
const filters = useFiltersStore();

const { Cmd_K, Ctrl_K } = useMagicKeys();
watch([Cmd_K, Ctrl_K], (v) => {
  if (v) {
    combobox.value?.focus();
  }
});

onMounted(async () => {
  try {
    tags.value = await getTags();
  } catch (err) {
    console.error(err);
    toast.error(`Failed to fetch tags:\n${err}`);
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped lang="scss">
.v-input:deep(.v-field__outline) {
  display: none;
}
</style>
