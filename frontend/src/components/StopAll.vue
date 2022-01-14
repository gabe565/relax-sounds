<template>
  <v-tooltip bottom>
    <template #activator="{ on, attrs }">
      <v-btn
        v-bind="attrs"
        v-on="on"
        @click="stopAll"
        :disabled="isStopped"
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
    'isStopped',
  ]),

  methods: mapActions('player', [
    'stopAll',
  ]),
};
</script>
