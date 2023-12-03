<template>
  <v-row class="filters">
    <v-col>
      <v-combobox
        v-model="filters.filters.word"
        label="Search"
        :prepend-inner-icon="SearchIcon"
        :loading="loading"
        clearable
        rounded
        density="comfortable"
        hide-details
        chips
        hide-no-data
        auto-select-first
        :items="tags"
        item-title="name"
        item-value="name"
        :return-object="false"
      >
        <template #chip="{ props, item }">
          <v-chip v-bind="props" :title="item.raw.name">
            <template #prepend>
              <v-icon v-if="item.raw.icon" class="mr-2">
                <Icon :icon="item.raw.icon" />
              </v-icon>
            </template>
          </v-chip>
        </template>

        <template #item="{ props, item }">
          <v-list-item v-bind="props" :title="item.raw.name">
            <template #prepend>
              <v-icon v-if="item.raw.icon">
                <Icon :icon="item.raw.icon" />
              </v-icon>
            </template>
          </v-list-item>
        </template>
      </v-combobox>
    </v-col>
  </v-row>
  <slot />
</template>

<script setup>
import { onMounted, ref } from "vue";
import SearchIcon from "~icons/material-symbols/search-rounded";
import { useFiltersStore } from "../../plugins/store/filters";
import { getTags } from "../../data/tags";
import { Icon } from "@iconify/vue";

const tags = ref([]);
const loading = ref(true);

const filters = useFiltersStore();

onMounted(async () => {
  tags.value = await getTags();
  loading.value = false;
});
</script>
