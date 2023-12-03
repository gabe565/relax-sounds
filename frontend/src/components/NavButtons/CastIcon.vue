<template>
  <v-tooltip v-if="castEnabled" v-model="showTooltip" location="bottom" aria-label="Cast">
    <template #activator="{ props }">
      <v-btn elevation="0" icon aria-label="Cast" @click.stop="element.click()">
        <v-icon v-bind="props">
          <google-cast-launcher ref="element" />
        </v-icon>
      </v-btn>
    </template>
    <span>Cast</span>
  </v-tooltip>
</template>

<script setup>
import { ref } from "vue";
import { castEnabled } from "../../cast";

const showTooltip = ref(false);
const element = ref();
</script>

<style scoped lang="scss">
google-cast-launcher {
  --disconnected-color: #fff;
}
</style>
