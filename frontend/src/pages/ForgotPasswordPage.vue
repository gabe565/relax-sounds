<template>
  <page-layout>
    <v-card max-width="400" class="mx-auto mt-8" color="cardBackground" variant="flat">
      <v-card-title class="text-center pt-6">Reset Password</v-card-title>
      <v-card-text>
        <v-form @submit.prevent="requestReset">
          <v-alert v-if="alert.text" v-bind="alert" class="my-6" />

          <p class="mb-4 text-medium-emphasis">
            Enter your email address and we'll send you a link to reset your password.
          </p>

          <v-text-field
            v-model="email"
            label="Email"
            type="email"
            variant="outlined"
            density="comfortable"
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
          <v-btn variant="text" size="small" to="/login">Back to Login</v-btn>
        </div>
      </v-card-text>
    </v-card>
  </page-layout>
</template>

<script setup>
import { reactive, ref, watchEffect } from "vue";
import { useRouter } from "vue-router";
import { useAuth } from "@/composables/useAuth.js";
import PageLayout from "@/layouts/PageLayout.vue";
import { getErrorMessage, pb } from "@/plugins/pocketbase";

const email = ref("");
const isLoading = ref(false);
const router = useRouter();
const { authMethods, isAuthenticated } = useAuth();
const alert = reactive({});

watchEffect(async () => {
  if (
    isAuthenticated.value ||
    (!authMethods.value.loading && !authMethods.value.password?.enabled)
  ) {
    await router.replace("/");
  }
});

const requestReset = async () => {
  if (!email.value) return;

  isLoading.value = true;
  try {
    await pb.collection("users").requestPasswordReset(email.value);
    alert.text = "Password reset link sent. Please check your email.";
    alert.color = "success";
  } catch (error) {
    console.error(error);
    alert.text = getErrorMessage(error);
    alert.color = "error";
  } finally {
    isLoading.value = false;
  }
};
</script>
