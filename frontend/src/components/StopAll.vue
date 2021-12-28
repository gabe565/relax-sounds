<template>
  <v-tooltip bottom>
    <template #activator="{ on, attrs }">
      <v-btn
        v-bind="attrs"
        v-on="on"
        @click="stopAll"
        :disabled="state === 'stopped'"
        icon
      >
        <v-icon>fas fa-stop-circle</v-icon>
      </v-btn>
    </template>
    <span>Stop All</span>
  </v-tooltip>
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
