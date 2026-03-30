<template>
  <v-dialog v-model="dialog" max-width="500">
    <template #activator="{ props: dialogProps }">
      <v-list-item
        title="Edit Profile"
        :prepend-icon="EditIcon"
        v-bind="dialogProps"
        @click="openDialog"
      />
    </template>

    <v-card color="cardBackground" variant="flat">
      <v-card-title class="pt-6 px-6">Edit Profile</v-card-title>
      <v-card-text>
        <v-form @submit.prevent="saveProfile">
          <div class="d-flex justify-center mb-6">
            <v-avatar
              size="96"
              class="elevation-2 cursor-pointer text-center"
              @click="fileInput.click()"
            >
              <v-img v-if="avatarPreview || avatarURL" :src="avatarPreview || avatarURL" />
              <v-icon v-else :icon="PersonIcon" size="64" />
              <v-overlay
                activator="parent"
                class="align-center justify-center"
                scrim="black"
                contained
                open-on-hover
              >
                <v-icon :icon="EditIcon" color="white" />
              </v-overlay>
            </v-avatar>
            <input
              ref="fileInput"
              type="file"
              accept="image/*"
              class="d-none"
              @change="onFileChange"
            />
          </div>

          <v-text-field
            v-model="name"
            label="Name"
            variant="outlined"
            density="comfortable"
            class="mb-2"
          />

          <v-text-field
            v-model="username"
            label="Username"
            variant="outlined"
            density="comfortable"
            class="mb-2"
          />

          <v-btn
            v-if="authMethods.password?.enabled"
            variant="outlined"
            block
            class="mb-6"
            :prepend-icon="LockResetIcon"
            @click="resetPassword"
          >
            Reset Password
          </v-btn>

          <template v-if="authMethods.oauth2?.providers?.length">
            <div class="text-title-medium mb-2">Linked Accounts</div>
            <v-alert
              v-if="externalAuthsError"
              type="error"
              density="compact"
              variant="tonal"
              class="mb-4"
            >
              {{ getErrorMessage(externalAuthsError) }}
            </v-alert>
            <v-list v-else class="bg-transparent pa-0 mb-4">
              <v-list-item v-for="provider in authMethods.oauth2?.providers" :key="provider.name">
                <template #prepend>
                  <v-avatar size="24" rounded="0" class="mr-3">
                    <v-img :src="provider.icon" :cover="false" />
                  </v-avatar>
                </template>
                <v-list-item-title class="text-body-2">
                  {{ provider.displayName }}
                </v-list-item-title>
                <template #append>
                  <v-btn
                    v-if="isLinked(provider)"
                    size="small"
                    color="error"
                    :loading="linkingProvider === provider"
                    @click="unlinkProvider(provider)"
                  >
                    Unlink
                  </v-btn>
                  <v-btn
                    v-else
                    size="small"
                    variant="outlined"
                    :loading="linkingProvider === provider"
                    @click="linkProvider(provider)"
                  >
                    Link
                  </v-btn>
                </template>
              </v-list-item>
            </v-list>
          </template>

          <div class="d-flex justify-end">
            <v-btn variant="text" class="mr-2" @click="dialog = false">Cancel</v-btn>
            <v-btn type="submit" color="primary" :loading="isLoading" variant="flat">Save</v-btn>
          </div>
        </v-form>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { useAsyncState } from "@vueuse/core";
import { ref, useTemplateRef } from "vue";
import { useToast } from "vue-toastification";
import EditIcon from "~icons/material-symbols/edit-rounded";
import LockResetIcon from "~icons/material-symbols/lock-reset-rounded";
import PersonIcon from "~icons/material-symbols/person-rounded";
import { useAuth } from "@/composables/useAuth.js";
import { getErrorMessage, pb } from "@/plugins/pocketbase";

const props = defineProps({
  user: {
    type: Object,
    required: true,
  },
});

const fileInput = useTemplateRef("fileInput");
const toast = useToast();
const { authMethods, avatarURL } = useAuth();

const dialog = ref(false);
const isLoading = ref(false);
const name = ref(props.user?.name ?? "");
const username = ref(props.user?.username ?? "");
const avatarFile = ref(null);
const avatarPreview = ref(null);
const linkingProvider = ref(null);

const isLinked = (provider) => {
  return externalAuths.value.some((auth) => auth.provider === provider.name);
};

const {
  state: externalAuths,
  execute: fetchExternalAuths,
  error: externalAuthsError,
} = useAsyncState(
  async () =>
    await pb.collection("_externalAuths").getFullList({ filter: `recordRef = "${props.user.id}"` }),
  [],
  { immediate: false },
);

const openDialog = async () => {
  name.value = props.user?.name ?? "";
  username.value = props.user?.username ?? "";
  avatarFile.value = null;
  avatarPreview.value = null;
  await fetchExternalAuths();
};

const linkProvider = async (provider) => {
  if (!props.user) return;

  linkingProvider.value = provider;
  try {
    await pb.collection("users").authWithOAuth2({ provider: provider.name });
    toast.success(`Linked ${provider.displayName}`);
    await fetchExternalAuths();
  } catch (error) {
    console.error(error);
    toast.error(getErrorMessage(error));
  } finally {
    linkingProvider.value = null;
  }
};

const unlinkProvider = async (provider) => {
  if (!props.user) return;

  const record = externalAuths.value.find((a) => a.provider === provider.name);
  if (!record) return;

  linkingProvider.value = provider;
  try {
    await pb.collection("_externalAuths").delete(record.id);
    toast.success(`Unlinked ${provider.displayName}`);
    await fetchExternalAuths();
  } catch (error) {
    console.error(error);
    toast.error(getErrorMessage(error));
  } finally {
    linkingProvider.value = null;
  }
};

const onFileChange = (e) => {
  const file = e.target.files[0];
  if (file) {
    avatarFile.value = file;
    avatarPreview.value = URL.createObjectURL(file);
  }
};

const saveProfile = async () => {
  if (!props.user) return;

  isLoading.value = true;
  try {
    const formData = new FormData();
    formData.append("name", name.value);
    formData.append("username", username.value);
    if (avatarFile.value) {
      formData.append("avatar", avatarFile.value);
    }

    await pb.collection("users").update(props.user.id, formData);
    toast.success("Profile updated");
    dialog.value = false;
  } catch (error) {
    console.error(error);
    toast.error(getErrorMessage(error));
  } finally {
    isLoading.value = false;
  }
};

const resetPassword = async () => {
  if (!props.user) return;

  try {
    await pb.collection("users").requestPasswordReset(props.user.email);
    toast.success("Password reset email sent.");
  } catch (e) {
    console.error(e);
    toast.error(getErrorMessage(e));
  }
};
</script>
