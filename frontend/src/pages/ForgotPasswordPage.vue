<template>
  <page-layout>
    <v-card
      max-width="400"
      class="border-outline-variant mx-auto mt-8 border"
      color="surface-container-high"
      variant="flat"
      rounded="xl"
    >
      <v-card-text class="pt-6">
        <v-form @submit.prevent="requestReset">
          <v-alert v-if="alert.text" v-bind="alert" variant="tonal" class="my-6" />

          <p class="text-on-surface-variant mb-4">
            Enter your email address and we'll send you a link to reset your password.
          </p>

          <v-text-field
            v-model="email"
            label="Email"
            type="email"
            variant="outlined"
            density="comfortable"
            rounded="lg"
            :prepend-inner-icon="MailIcon"
            autocomplete="email"
            class="mb-4"
            :rules="[(v) => !!v || 'Email is required']"
            required
          />

          <v-btn
            type="submit"
            color="primary"
            block
            size="large"
            :loading="isLoading"
            variant="flat"
          >
            Send Reset Link
          </v-btn>
        </v-form>

        <div class="text-center mt-4">
          <v-btn variant="text" size="small" to="/login">Back to log in</v-btn>
        </div>
      </v-card-text>
    </v-card>
  </page-layout>
</template>

<script setup>
import { reactive, ref, watchEffect } from "vue";
import { useRouter } from "vue-router";
import MailIcon from "~icons/material-symbols/mail-rounded";
import PageLayout from "@/layouts/PageLayout.vue";
import { getErrorMessage, usePocketBase } from "@/plugins/store/pocketbase.js";

const email = ref("");
const isLoading = ref(false);
const router = useRouter();
const pb = usePocketBase();
const alert = reactive({});

watchEffect(async () => {
  if (pb.isAuthenticated || (!pb.authMethods.loading && !pb.authMethods.password?.enabled)) {
    await router.replace("/");
  }
});

const requestReset = async () => {
  if (!email.value) return;

  isLoading.value = true;
  try {
    await pb.client.collection("users").requestPasswordReset(email.value);
    alert.text = "Password reset link sent. Please check your email.";
    alert.type = "success";
  } catch (error) {
    console.error(error);
    alert.text = getErrorMessage(error);
    alert.type = "error";
  } finally {
    isLoading.value = false;
  }
};
</script>
