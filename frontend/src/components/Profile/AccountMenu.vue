<template>
  <v-list-group
    v-if="pb.isAuthenticated && collapsible"
    value="account"
    expand-icon="$collapse"
    collapse-icon="$expand"
  >
    <template #activator="{ props: groupProps }">
      <v-list-item
        v-bind="groupProps"
        :title="pb.user.name || pb.user.username"
        :prepend-avatar="pb.avatarURL"
      >
        <template v-if="!pb.avatarURL" #prepend>
          <v-icon :icon="PersonIcon" />
        </template>
      </v-list-item>
    </template>

    <v-list-item title="Edit Profile" :prepend-icon="EditIcon" @click="profileOpen = true" />
    <v-list-item
      to="/logout"
      title="Log out"
      :prepend-icon="LogoutIcon"
      class="text-error"
      @click.prevent="pb.logout"
    />
  </v-list-group>

  <v-menu v-else-if="pb.isAuthenticated" location="bottom right" transition="slide-y-transition">
    <template #activator="{ props: menuProps }">
      <v-list-item
        v-if="listItem"
        v-bind="menuProps"
        :title="pb.user.name || pb.user.username"
        :prepend-avatar="pb.avatarURL"
        aria-label="Account"
      >
        <template v-if="!pb.avatarURL" #prepend>
          <v-icon :icon="PersonIcon" />
        </template>
      </v-list-item>

      <v-btn v-else icon v-bind="menuProps" :loading="isLoading" aria-label="Account">
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
          title="Log out"
          :prepend-icon="LogoutIcon"
          class="text-error"
          @click.prevent="pb.logout"
        />
      </v-list>
    </template>
  </v-menu>

  <profile-dialog v-if="pb.isAuthenticated" v-model="profileOpen" :user="pb.user" />
</template>

<script setup>
import { computed, ref } from "vue";
import EditIcon from "~icons/material-symbols/edit-rounded";
import LogoutIcon from "~icons/material-symbols/logout-rounded";
import PersonIcon from "~icons/material-symbols/person-rounded";
import ProfileDialog from "@/components/Profile/ProfileDialog.vue";
import { usePocketBase } from "@/plugins/store/pocketbase.js";
import { usePresets } from "@/plugins/store/presets";

defineProps({
  listItem: {
    type: Boolean,
    default: false,
  },
  collapsible: {
    type: Boolean,
    default: false,
  },
});

const pb = usePocketBase();
const presets = usePresets();
const profileOpen = ref(false);
const isLoading = computed(() => presets.isSyncing || !pb.user?.id);
</script>

<style scoped>
.v-list-group {
  display: flex;
  flex-direction: column-reverse;
}
</style>
