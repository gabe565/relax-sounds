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
        variant="solo"
        flat
        bg-color="surface-container-highest"
        rounded="pill"
        density="comfortable"
        hide-details
        hide-no-data
        :items="tags"
        :error="error"
        item-title="name"
        item-value="name"
        :return-object="false"
        color="primary"
        @keydown.esc="
          filters.filters.word = '';
          combobox.blur();
        "
        @keydown.enter="combobox.blur()"
      >
        <template #item="{ props, item }">
          <v-list-item v-bind="props" :title="item?.name ?? item">
            <template #prepend>
              <sound-icon
                v-if="item?.icon"
                :icon="item.icon"
                :size="32"
                :color="item.color ? `tag-${item.color}` : 'primary'"
                class="mr-3"
              />
            </template>
          </v-list-item>
        </template>
      </v-combobox>
    </v-col>
  </v-row>
</template>

<script setup>
import { useAsyncState, useMagicKeys } from "@vueuse/core";
import { useTemplateRef, watch } from "vue";
import { toast } from "vue-sonner";
import SearchIcon from "~icons/material-symbols/search-rounded";
import SoundIcon from "@/components/Sounds/SoundIcon.vue";
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
