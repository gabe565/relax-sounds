<template>
  <template v-if="castEnabled">
    <v-tooltip v-if="button" text="Cast" location="bottom">
      <template #activator="{ props }">
        <v-btn
          v-bind="props"
          icon
          title="Cast"
          aria-label="Cast"
          @click.stop="castLauncher.click()"
        >
          <v-icon>
            <google-cast-launcher ref="castLauncher" class="position-relative" />
          </v-icon>
        </v-btn>
      </template>
    </v-tooltip>
    <v-list-item v-else title="Cast" aria-label="Cast" @click.stop="castLauncher.click()">
      <template #prepend>
        <v-icon>
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

defineProps({
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
