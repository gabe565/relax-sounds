<template>
  <template v-if="loading">
    <slot />
  </template>
  <template v-else>
    <v-row class="filters">
      <v-col class="pb-0">
        <v-text-field
          v-model="filters.filters.word"
          label="Search"
          :prepend-icon="SearchIcon"
          clearable
          variant="underlined"
        />
      </v-col>
      <v-col class="flex-grow-0 pb-0">
        <v-switch v-model="filters.filters.playing" label="Playing" inset />
      </v-col>
    </v-row>
    <v-row class="pb-5">
      <v-chip-group v-model="filters.filters.word" column>
        <v-chip
          v-for="(tag, key) in tags"
          :key="key"
          :value="tag.name"
          active-class="deep-orange"
          class="ma-2"
          filter
        >
          <v-icon v-if="tag.icon" size="x-small" class="mr-2">
            <Icon :icon="tag.icon" />
          </v-icon>
          {{ tag.name }}
        </v-chip>
      </v-chip-group>
    </v-row>
    <v-row>
      <v-col>
        <v-divider />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-pagination
          v-model="filters.filters.page"
          :length="filters.pages"
          size="small"
          active-color="primary"
        />
      </v-col>
    </v-row>
    <slot />
    <v-row>
      <v-col>
        <v-pagination
          v-model="filters.filters.page"
          :length="filters.pages"
          size="small"
          active-color="primary"
        />
      </v-col>
    </v-row>
  </template>
</template>

<script setup>
import { onMounted, ref, watch } from "vue";
import SearchIcon from "~icons/material-symbols/search-rounded";
import { useFiltersStore } from "../../plugins/store/filters";
import { getTags } from "../../data/tags";
import { Icon } from "@iconify/vue";

const tags = ref(null);
const loading = ref(true);

const filters = useFiltersStore();

watch(
  () => filters.page,
  () => {
    filters.page = 1;
  },
);

onMounted(async () => {
  tags.value = await getTags();
  loading.value = false;
});
</script>

<style scoped>
.filters :deep(.v-input__control) {
  grid-area: auto;
}
</style>
