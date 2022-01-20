<template>
  <v-container>
    <v-row class="pa-5">
      <h1 class="display-2">
        <slot name="title">{{ $route.name }}</slot>
      </h1>
    </v-row>

    <slot/>

  </v-container>
</template>

<script>
export default {
  name: 'Page',
};
</script>
