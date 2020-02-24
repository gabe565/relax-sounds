<template>
  <v-btn icon @click="stopAll" :disabled="state === 'stopped'">
    <v-icon>fal fa-stop-circle</v-icon>
  </v-btn>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';

export default {
  name: 'StopAll',

  computed: mapGetters('sounds', [
    'state',
  ]),

  methods: mapActions('sounds', [
    'stopAll',
  ]),
};
</script>

<style scoped>

</style>
