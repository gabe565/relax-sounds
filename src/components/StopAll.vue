<template>
  <v-btn icon @click="stopAll" :disabled="state === 'stopped'">
    <v-icon>fas fa-stop-circle</v-icon>
  </v-btn>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';

export default {
  name: 'StopAll',

  computed: mapGetters('player', [
    'state',
  ]),

  methods: mapActions('player', [
    'stopAll',
  ]),
};
</script>

<style scoped>

</style>
