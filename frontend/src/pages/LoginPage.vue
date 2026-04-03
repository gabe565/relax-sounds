<template>
  <page-layout>
    <v-card max-width="400" class="mx-auto mt-8" color="cardBackground" variant="flat">
      <v-card-title class="text-center py-6">
        {{ route.name }}
      </v-card-title>
      <v-card-text>
        <template v-if="pb.authMethods.loading">
          <div class="flex justify-center py-8">
            <v-progress-circular indeterminate color="primary" />
          </div>
        </template>
        <template v-else>
          <!-- Password Form -->
          <v-form
            v-if="pb.authMethods.password?.enabled"
            @submit.prevent="props.register ? registerWithPassword() : loginWithPassword()"
          >
            <v-alert v-if="alert.text" v-bind="alert" class="mb-6" />
            <v-text-field
              v-model="email"
              label="Email"
              type="email"
              variant="outlined"
              density="comfortable"
              class="mb-2"
              :rules="[(v) => !!v || 'Email is required']"
              required
            />
            <v-text-field
              v-model="password"
              label="Password"
              type="password"
              variant="outlined"
              density="comfortable"
              class="mb-2"
              :rules="[(v) => !!v || 'Password is required']"
              required
            />
            <v-text-field
              v-if="props.register"
              v-model="passwordConfirm"
              label="Confirm Password"
              type="password"
              variant="outlined"
              density="comfortable"
              class="mb-4"
              :rules="[
                (v) => !!v || 'Password confirmation is required',
                (v) => v === password || 'Passwords do not match',
              ]"
              required
            />
            <div v-else class="text-right mb-4">
              <v-btn
                variant="text"
                size="x-small"
                to="/forgot-password"
                class="text-caption text-none pa-0"
              >
                Forgot Password?
              </v-btn>
            </div>
            <v-btn
              type="submit"
              color="primary"
              block
              size="large"
              :loading="isLoading"
              variant="flat"
            >
              {{ route.name }}
            </v-btn>
          </v-form>

          <div v-if="pb.authMethods.password?.enabled" class="text-center mt-4">
            <v-btn variant="text" size="small" :to="props.register ? '/login' : '/register'">
              {{ props.register ? "Already have an account?" : "Don't have an account?" }}
            </v-btn>
          </div>

          <div v-if="showResend" class="text-center mt-2">
            <v-btn
              variant="outlined"
              color="primary"
              size="small"
              :loading="isResending"
              @click="resendVerification"
            >
              Resend verification email
            </v-btn>
          </div>

          <!-- Divider -->
          <div
            v-if="pb.authMethods.password?.enabled && pb.authMethods.oauth2?.providers?.length"
            class="flex items-center my-6"
          >
            <v-divider />
            <span class="mx-4 text-caption text-medium-emphasis">OR</span>
            <v-divider />
          </div>

          <!-- OAuth2 Providers -->
          <template v-if="pb.authMethods.oauth2?.providers?.length">
            <v-btn
              v-for="provider in pb.authMethods.oauth2.providers"
              :key="provider.name"
              variant="outlined"
              block
              size="large"
              class="mb-3"
              :loading="providerLoading === provider.name"
              @click="loginWithProvider(provider)"
            >
              <template #prepend>
                <v-avatar size="24" rounded="0" variant="text">
                  <v-img :src="provider.icon" :cover="false" />
                </v-avatar>
              </template>
              Continue with {{ provider.displayName }}
            </v-btn>
          </template>
        </template>
      </v-card-text>
    </v-card>
  </page-layout>
</template>

<script setup>
import { reactive, ref, watchEffect } from "vue";
import { useRoute, useRouter } from "vue-router";
import PageLayout from "@/layouts/PageLayout.vue";
import { getErrorMessage, usePocketBase } from "@/plugins/store/pocketbase.js";

const props = defineProps({
  register: {
    type: Boolean,
    default: false,
  },
});

const router = useRouter();
const pb = usePocketBase();

watchEffect(async () => {
  if (pb.isAuthenticated || (!pb.authMethods.loading && !pb.authEnabled)) {
    await router.replace("/");
  }
});

const route = useRoute();

const email = ref("");
const password = ref("");
const passwordConfirm = ref("");
const isLoading = ref(false);
const providerLoading = ref(null);
const alert = reactive({});
const showResend = ref(false);
const isResending = ref(false);

const handleAuthError = (error) => {
  console.error(error);
  const response = error.response;
  if (
    (error.status === 400 || error.status === 403) &&
    response?.message?.includes("satisfy the collection requirements")
  ) {
    showResend.value = true;
    alert.text = "Please verify your email address before logging in.";
    alert.color = "error";
  } else {
    alert.text = getErrorMessage(error);
    alert.color = "error";
  }
};

const resendVerification = async () => {
  if (!email.value) return;
  isResending.value = true;
  try {
    await pb.client.collection("users").requestVerification(email.value);
    alert.text = "Verification email sent.";
    alert.color = "success";
    showResend.value = false;
  } catch (error) {
    handleAuthError(error);
  } finally {
    isResending.value = false;
  }
};

const loginWithPassword = async () => {
  if (!email.value || !password.value) return;

  isLoading.value = true;
  try {
    await pb.client.collection("users").authWithPassword(email.value, password.value);
    await router.push("/");
  } catch (error) {
    handleAuthError(error);
  } finally {
    isLoading.value = false;
  }
};

const registerWithPassword = async () => {
  if (!email.value || !password.value || password.value !== passwordConfirm.value) return;

  isLoading.value = true;
  try {
    await pb.client.collection("users").create({
      email: email.value,
      password: password.value,
      passwordConfirm: passwordConfirm.value,
    });
    await pb.client.collection("users").requestVerification(email.value);
    alert.text = "Account created. Please check your email for verification.";
    alert.color = "success";
    await router.push("/login");
  } catch (error) {
    handleAuthError(error);
  } finally {
    isLoading.value = false;
  }
};

const loginWithProvider = async (provider) => {
  providerLoading.value = provider.name;
  try {
    await pb.client.collection("users").authWithOAuth2({ provider: provider.name });
    await router.push("/");
  } catch (error) {
    handleAuthError(error);
  } finally {
    providerLoading.value = null;
  }
};
</script>
