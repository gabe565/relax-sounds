<template>
  <v-row class="filters">
    <v-col>
      <v-combobox
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
import { onMounted, ref } from "vue";
import { useToast } from "vue-toastification";
import SearchIcon from "~icons/material-symbols/search-rounded";
import { getTags } from "@/data/tags";
import { useFiltersStore } from "@/plugins/store/filters";

const toast = useToast();
const tags = ref([]);
const loading = ref(true);

const filters = useFiltersStore();

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
