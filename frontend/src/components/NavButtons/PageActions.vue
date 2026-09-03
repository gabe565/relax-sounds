<template>
  <v-menu
    v-if="isMobile || hasMenu || DebugEnabled"
    location="bottom right"
    transition="slide-y-transition"
  >
    <template #activator="{ props: menuProps }">
      <v-btn icon variant="flat" color="transparent" v-bind="menuProps" aria-label="Menu">
        <v-icon :icon="MenuIcon" />
      </v-btn>
    </template>

    <v-list>
      <slot name="menu" />
      <v-divider v-if="hasMenu && (isMobile || DebugEnabled)" />
      <theme-btn v-if="isMobile" list-item />
      <debug-button v-if="DebugEnabled" />
    </v-list>
  </v-menu>

  <v-btn
    v-if="isMobile && pb.authEnabled && !pb.isAuthenticated && !route.meta.hideLogin"
    to="/login"
    icon
    aria-label="Log in"
  >
    <v-icon :icon="LoginIcon" />
  </v-btn>
</template>

<script setup>
import { useRoute } from "vue-router";
import { useDisplay } from "vuetify";
import LoginIcon from "~icons/material-symbols/login-rounded";
import MenuIcon from "~icons/material-symbols/more-horiz";
import ThemeBtn from "@/components/NavButtons/ThemeBtn.vue";
import DebugButton from "@/components/Presets/Buttons/DebugButton.vue";
import { DebugEnabled } from "@/config/debug";
import { usePocketBase } from "@/plugins/store/pocketbase.js";

defineProps({
  hasMenu: {
    type: Boolean,
    default: false,
  },
});

const { smAndDown: isMobile } = useDisplay();
const route = useRoute();
const pb = usePocketBase();
</script>
