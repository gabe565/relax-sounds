<template>
  <v-container>
    <v-row class="pa-5">
      <v-col>
        <h1 class="display-1">
          <slot name="title">{{ $route.name }}</slot>
        </h1>
      </v-col>

      <v-spacer/>

      <v-col class="shrink" v-if="actions">
        <v-menu left bottom>
          <template #activator="{ on, attrs }">
            <v-btn icon v-bind="attrs" v-on="on">
              <v-icon>fas fa-ellipsis-v</v-icon>
            </v-btn>
          </template>

          <v-list>
            <v-list-item v-for="(item, index) in actions" :key="index" v-on="item.on">
              <v-list-item-icon>
                <v-icon aria-hidden="true">{{ item.icon }}</v-icon>
              </v-list-item-icon>
              <v-list-item-title>{{ item.title }}</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-col>
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
    actions: Array,
  },
  data: () => ({
    showAlert: false,
  }),
  async created() {
    if (this.alert) this.showAlert = true;
  },
};
</script>
