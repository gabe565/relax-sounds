<template>
  <v-app-bar color="surface" flat :title="route.name">
    <template #prepend>
      <v-btn v-if="isMobile" to="/" icon size="small">
        <v-icon :icon="AppIcon" size="28" color="primary" aria-label="Relax Sounds" />
      </v-btn>
    </template>

    <v-spacer />

    <slot name="actions" />

    <v-menu
      v-if="isMobile || $slots.menu || DebugEnabled"
      location="bottom right"
      transition="slide-y-transition"
    >
      <template #activator="{ props }">
        <v-btn icon variant="flat" color="transparent" v-bind="props" aria-label="Menu">
          <v-icon :icon="MenuIcon" />
        </v-btn>
      </template>

      <v-list>
        <slot name="menu" />
        <v-divider v-if="$slots.menu && (isMobile || DebugEnabled)" />
        <template v-if="isMobile">
          <theme-btn list-item />
        </template>
        <debug-button v-if="DebugEnabled" />
      </v-list>
    </v-menu>

    <v-btn
      v-if="pb.authEnabled && !pb.isAuthenticated && !route.meta.hideLogin"
      to="/login"
      icon
      aria-label="Login"
    >
      <v-icon :icon="LoginIcon" />
    </v-btn>

    <v-menu v-if="pb.isAuthenticated" location="bottom right" transition="slide-y-transition">
      <template #activator="{ props: menuProps }">
        <v-btn icon v-bind="menuProps" :loading="isLoading">
          <v-avatar size="32">
            <v-img v-if="pb.avatarURL" :src="pb.avatarURL" :alt="pb.user.name || pb.user.email" />
            <v-icon v-else :icon="PersonIcon" />
          </v-avatar>
        </v-btn>
      </template>

      <template #default="{ isActive }">
        <v-list width="250">
          <v-list-item
            :title="pb.user.name || pb.user.username"
            :subtitle="pb.user.email"
            :prepend-avatar="pb.avatarURL"
            class="pb-2"
          >
            <template v-if="!pb.avatarURL" #prepend>
              <v-icon :icon="PersonIcon" />
            </template>
          </v-list-item>
          <v-divider class="mt-2" />
          <v-list-item
            title="Edit Profile"
            :prepend-icon="EditIcon"
            @click="
              isActive.value = false;
              profileOpen = true;
            "
          />
          <v-list-item
            to="/logout"
            title="Logout"
            :prepend-icon="LogoutIcon"
            class="text-error"
            @click.prevent="pb.logout"
          />
        </v-list>
      </template>
    </v-menu>

    <profile-dialog v-if="pb.isAuthenticated" v-model="profileOpen" :user="pb.user" />
  </v-app-bar>

  <v-container class="pt-6">
    <slot />
  </v-container>
</template>

<script setup>
import { computed, ref } from "vue";
import { useRoute } from "vue-router";
import { useDisplay } from "vuetify";
import EditIcon from "~icons/material-symbols/edit-rounded";
import LoginIcon from "~icons/material-symbols/login-rounded";
import LogoutIcon from "~icons/material-symbols/logout-rounded";
import MenuIcon from "~icons/material-symbols/more-horiz";
import PersonIcon from "~icons/material-symbols/person-rounded";
import AppIcon from "~icons/relax-sounds/icon";
import ThemeBtn from "@/components/NavButtons/ThemeBtn.vue";
import DebugButton from "@/components/Presets/Buttons/DebugButton.vue";
import ProfileDialog from "@/components/Profile/ProfileDialog.vue";
import { DebugEnabled } from "@/config/debug";
import { usePocketBase } from "@/plugins/store/pocketbase.js";
import { usePresets } from "@/plugins/store/presets";

const { smAndDown: isMobile } = useDisplay();
const route = useRoute();
const pb = usePocketBase();
const presets = usePresets();
const profileOpen = ref(false);

const isLoading = computed(() => presets.isSyncing || !pb.user?.id);
</script>
