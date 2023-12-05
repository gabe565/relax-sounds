<template>
  <page-layout>
    <v-row>
      <v-col>
        I hate to be the bearer of bad news, but the requested file could not be located.
        <router-link to="/">Click here to go home.</router-link>
      </v-col>
    </v-row>
  </page-layout>
</template>

<script setup>
import PageLayout from "../layouts/PageLayout.vue";
</script>
