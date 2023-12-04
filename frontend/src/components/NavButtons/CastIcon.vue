<template>
  <template v-if="castEnabled">
    <v-btn v-if="button" icon title="Cast" aria-label="Cast" @click.stop="castLauncher.click()">
      <v-icon v-bind="props">
        <google-cast-launcher ref="castLauncher" />
      </v-icon>
    </v-btn>
    <v-list-item v-else title="Cast" aria-label="Cast" @click.stop="castLauncher.click()">
      <template #prepend>
        <v-icon v-bind="props">
          <google-cast-launcher ref="castLauncher" />
        </v-icon>
      </template>
    </v-list-item>
  </template>
</template>

<script setup>
import { ref } from "vue";
import { castEnabled } from "../../cast";
import { VBtn, VListItem } from "vuetify/components";

const castLauncher = ref();

const props = defineProps({
  button: {
    type: Boolean,
    default: false,
  },
});
</script>

<style scoped lang="scss">
google-cast-launcher {
  --disconnected-color: #fff;
}
</style>
