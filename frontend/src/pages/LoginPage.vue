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
        <template v-if="pb.authMethods.loading">
          <div class="flex justify-center py-8">
            <v-progress-circular indeterminate color="primary" />
          </div>
        </template>
        <template v-else>
          <v-form
            v-if="pb.authMethods.password?.enabled"
            @submit.prevent="props.register ? registerWithPassword() : loginWithPassword()"
          >
            <v-alert v-if="alert.text" v-bind="alert" variant="tonal" class="mb-6" />
            <v-text-field
              v-model="email"
              label="Email"
              type="email"
              variant="outlined"
              density="comfortable"
              rounded="lg"
              :prepend-inner-icon="MailIcon"
              autocomplete="email"
              class="mb-2"
              :rules="[(v) => !!v || 'Email is required']"
              required
            />
            <v-text-field
              v-model="password"
              label="Password"
              :type="showPassword ? 'text' : 'password'"
              variant="outlined"
              density="comfortable"
              rounded="lg"
              :prepend-inner-icon="LockIcon"
              :append-inner-icon="showPassword ? VisibilityOffIcon : VisibilityIcon"
              :autocomplete="props.register ? 'new-password' : 'current-password'"
              class="mb-2"
              :rules="[(v) => !!v || 'Password is required']"
              required
              @click:append-inner="showPassword = !showPassword"
            />
            <v-text-field
              v-if="props.register"
              v-model="passwordConfirm"
              label="Confirm Password"
              :type="showPassword ? 'text' : 'password'"
              variant="outlined"
              density="comfortable"
              rounded="lg"
              :prepend-inner-icon="LockIcon"
              autocomplete="new-password"
              class="mb-4"
              :rules="[
                (v) => !!v || 'Password confirmation is required',
                (v) => v === password || 'Passwords do not match',
              ]"
              required
            />
            <div v-else class="text-right mb-4">
              <v-btn variant="text" size="small" to="/reset-password" class="text-none">
                Reset password
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
              {{ props.register ? "Log in" : "Create account" }}
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

          <div
            v-if="pb.authMethods.password?.enabled && pb.authMethods.oauth2?.providers?.length"
            class="flex items-center my-6"
          >
            <v-divider />
            <span class="text-on-surface-variant mx-4 text-xs">OR</span>
            <v-divider />
          </div>

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
import { onMounted, reactive, ref, watchEffect } from "vue";
import { useRoute, useRouter } from "vue-router";
import { toast } from "vue-sonner";
import LockIcon from "~icons/material-symbols/lock-rounded";
import MailIcon from "~icons/material-symbols/mail-rounded";
import VisibilityOffIcon from "~icons/material-symbols/visibility-off-rounded";
import VisibilityIcon from "~icons/material-symbols/visibility-rounded";
import PageLayout from "@/layouts/PageLayout.vue";
import { SessionExpiredToast, getErrorMessage, usePocketBase } from "@/plugins/store/pocketbase.js";

const props = defineProps({
  register: {
    type: Boolean,
    default: false,
  },
});

const router = useRouter();
const pb = usePocketBase();

onMounted(() => toast.dismiss(SessionExpiredToast));

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
const showPassword = ref(false);
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
    alert.type = "error";
  } else {
    alert.text = getErrorMessage(error);
    alert.type = "error";
  }
};

const resendVerification = async () => {
  if (!email.value) return;
  isResending.value = true;
  try {
    await pb.client.collection("users").requestVerification(email.value);
    alert.text = "Verification email sent.";
    alert.type = "success";
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
    alert.type = "success";
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
