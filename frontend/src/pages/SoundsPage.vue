<template>
  <PageLayout :alert="alert" :actions="actions">
    <FilterSection>
      <v-row>
        <template v-if="loading">
          <v-overlay v-model="loading" class="align-center justify-center" persistent>
            <v-progress-circular color="primary" indeterminate size="64" />
          </v-overlay>
        </template>
        <template v-else>
          <v-col v-for="(sound, key) of sounds" :key="key" cols="12" md="6" lg="4">
            <SoundCard :sound="sound" />
          </v-col>
        </template>
      </v-row>
    </FilterSection>
  </PageLayout>
</template>

<script>
import { mapActions, mapGetters } from "vuex";
import PageLayout from "../layouts/PageLayout.vue";
import SoundCard from "../components/SoundCard.vue";
import FilterSection from "../components/FilterSection.vue";

export default {
  name: "SoundsPage",

  components: {
    FilterSection,
    PageLayout,
    SoundCard,
  },

  props: {
    alert: {
      type: Object,
      default: null,
    },
  },

  data: () => ({
    loading: true,
    error: null,
    page: 1,
  }),

  computed: {
    actions() {
      return [
        {
          title: "Preload All",
          icon: "fas fa-sync",
          on: {
            click: () => {
              this.$store.dispatch("player/prefetch");
            },
          },
        },
      ];
    },
    ...mapGetters("filters", ["sounds"]),
  },

  async created() {
    await this.initSounds();
    this.loading = false;
  },
  methods: mapActions("player", ["initSounds"]),
};
</script>
