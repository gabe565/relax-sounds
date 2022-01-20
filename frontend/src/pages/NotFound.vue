<template>
  <Page>
    <v-row>
      <v-col>
        I hate to be the bearer of bad news, but the requested file could not be located.
        <router-link to="/">Click here to go home.</router-link>
      </v-col>
    </v-row>
  </Page>
</template>

<script>
import Page from '../layouts/Page.vue';

export default {
  name: 'NotFoundPage',
  components: { Page },
};
</script>
