<template>
  <v-container>
    <v-row class="pa-5">
      <h1 class="display-1">
        <slot name="title">{{ $route.name }}</slot>
      </h1>
    </v-row>

    <v-fade-transition>
      <v-row v-if="showAlert">
        <v-col>
          <v-alert
            dismissible
            prominent
            text
            :type="alert.type"
            v-model="showAlert"
          >
            {{ alert.text }}
          </v-alert>
        </v-col>
      </v-row>
    </v-fade-transition>

    <slot/>

  </v-container>
</template>

<script>
export default {
  name: 'Page',
  props: {
    alert: Object,
  },
  data: () => ({
    showAlert: false,
  }),
  async created() {
    if (this.alert) this.showAlert = true;
  },
};
</script>
