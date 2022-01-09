<template>
  <Page>
    I hate to be the bearer of bad news, but the requested file could not be located.
    <router-link to="/">Click here to go home.</router-link>
  </Page>
</template>

<script>
import Page from '../layouts/Page.vue';

export default {
  name: 'NotFoundPage',
  components: { Page },
};
</script>
