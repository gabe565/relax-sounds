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
          <v-alert v-if="alert.text" v-bind="alert" variant="tonal" class="mb-6" />

          <v-alert
            v-if="mfaId"
            type="info"
            variant="tonal"
            :icon="ShieldLockIcon"
            class="mb-6"
            title="Verify it's you"
            text="Your account requires a second sign-in step."
          />

          <v-form
            v-if="props.register && pb.authMethods.password?.enabled"
            @submit.prevent="registerWithPassword"
          >
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
              autocomplete="new-password"
              class="mb-2"
              :rules="[(v) => !!v || 'Password is required']"
              required
              @click:append-inner="showPassword = !showPassword"
            />
            <v-text-field
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
            <v-btn
              type="submit"
              color="primary"
              block
              size="large"
              :loading="isLoading"
              variant="flat"
            >
              Register
            </v-btn>
          </v-form>

          <v-form v-else-if="step === 'otp'" @submit.prevent="loginWithOTP">
            <p class="text-on-surface-variant mb-4 text-sm">
              Enter the one-time code sent to <strong>{{ email }}</strong
              >.
            </p>
            <v-text-field
              v-model="otpCode"
              label="One-time code"
              type="text"
              inputmode="numeric"
              variant="outlined"
              density="comfortable"
              rounded="lg"
              :prepend-inner-icon="PinIcon"
              autocomplete="one-time-code"
              autofocus
              class="mb-4"
              :rules="[(v) => !!v || 'Code is required']"
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
              Verify code
            </v-btn>
            <div class="flex justify-between mt-4">
              <v-btn variant="text" size="small" :loading="isSendingOTP" @click="requestOTP">
                Resend code
              </v-btn>
              <v-btn v-if="!mfaId" variant="text" size="small" @click="goBackFromOTP">Back</v-btn>
            </div>
          </v-form>

          <v-form
            v-else-if="step === 'password' || (showPasswordForm && !emailFirst)"
            @submit.prevent="loginWithPassword"
          >
            <v-text-field
              v-model="email"
              label="Email"
              type="email"
              variant="outlined"
              density="comfortable"
              rounded="lg"
              :prepend-inner-icon="MailIcon"
              :append-inner-icon="emailFirst ? EditIcon : undefined"
              autocomplete="username"
              :readonly="emailFirst"
              :autofocus="!emailFirst"
              class="mb-2"
              :rules="[(v) => !!v || 'Email is required']"
              required
              @click:append-inner="changeEmail"
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
              autocomplete="current-password"
              :autofocus="emailFirst"
              class="mb-2"
              :rules="[(v) => !!v || 'Password is required']"
              required
              @click:append-inner="showPassword = !showPassword"
            />
            <div class="flex mb-4" :class="showOTP ? 'justify-between' : 'justify-end'">
              <v-btn
                v-if="showOTP"
                variant="text"
                size="small"
                class="text-none"
                :loading="isSendingOTP"
                @click="requestOTP"
              >
                Email me a code instead
              </v-btn>
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
              Log in
            </v-btn>
          </v-form>

          <v-form v-else-if="showPasswordForm || showOTP" @submit.prevent="continueWithEmail">
            <v-text-field
              v-model="email"
              label="Email"
              type="email"
              variant="outlined"
              density="comfortable"
              rounded="lg"
              :prepend-inner-icon="MailIcon"
              autocomplete="username"
              autofocus
              class="mb-4"
              :rules="[(v) => !!v || 'Email is required']"
              required
            />
            <v-btn
              type="submit"
              color="primary"
              block
              size="large"
              :loading="isSendingOTP"
              variant="flat"
            >
              {{ showPasswordForm ? "Continue" : "Email me a code" }}
            </v-btn>
          </v-form>

          <div v-if="pb.authMethods.password?.enabled && !mfaId" class="text-center mt-4">
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

          <template v-if="showOAuth && (props.register || onFirstStep)">
            <div
              v-if="
                (props.register && pb.authMethods.password?.enabled) || showPasswordForm || showOTP
              "
              class="flex items-center my-6"
            >
              <v-divider />
              <span class="text-on-surface-variant mx-4 text-xs">OR</span>
              <v-divider />
            </div>

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

          <v-alert
            v-if="mfaId && !showPasswordForm && !showOTP && !showOAuth"
            type="warning"
            variant="tonal"
            text="No second sign-in method is available. Please contact an administrator."
          />

          <div v-if="mfaId" class="text-center mt-2">
            <v-btn variant="text" size="small" @click="resetMFA">Cancel</v-btn>
          </div>
        </template>
      </v-card-text>
    </v-card>
  </page-layout>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watchEffect } from "vue";
import { useRouter } from "vue-router";
import { toast } from "vue-sonner";
import EditIcon from "~icons/material-symbols/edit-rounded";
import LockIcon from "~icons/material-symbols/lock-rounded";
import MailIcon from "~icons/material-symbols/mail-rounded";
import PinIcon from "~icons/material-symbols/pin-rounded";
import ShieldLockIcon from "~icons/material-symbols/shield-lock-rounded";
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

const step = ref("email");
const email = ref("");
const password = ref("");
const passwordConfirm = ref("");
const isLoading = ref(false);
const providerLoading = ref(null);
const alert = reactive({});
const showPassword = ref(false);
const showResend = ref(false);
const isResending = ref(false);

const mfaId = ref(null);
const mfaMethod = ref(null);
const otpId = ref(null);
const otpCode = ref("");
const isSendingOTP = ref(false);

const otpSecondFactor = computed(
  () => !!mfaId.value && !!pb.authMethods.otp?.enabled && mfaMethod.value !== "otp",
);
const showPasswordForm = computed(
  () =>
    !!pb.authMethods.password?.enabled && mfaMethod.value !== "password" && !otpSecondFactor.value,
);
const otpAllowedFirst = computed(
  () =>
    !(
      pb.authMethods.mfa?.enabled &&
      (pb.authMethods.password?.enabled || pb.authMethods.oauth2?.providers?.length)
    ),
);
const showOTP = computed(
  () =>
    !props.register &&
    !!pb.authMethods.otp?.enabled &&
    mfaMethod.value !== "otp" &&
    (!!mfaId.value || otpAllowedFirst.value),
);
const emailFirst = computed(() => showPasswordForm.value && showOTP.value);
const onFirstStep = computed(
  () => step.value === "email" || (step.value === "password" && !emailFirst.value),
);
const showOAuth = computed(
  () =>
    !!pb.authMethods.oauth2?.providers?.length &&
    mfaMethod.value !== "oauth2" &&
    !otpSecondFactor.value,
);

const mfaOptions = () => (mfaId.value ? { mfaId: mfaId.value } : {});

const clearAlert = () => {
  alert.text = "";
  showResend.value = false;
};

const resetMFA = () => {
  mfaId.value = null;
  mfaMethod.value = null;
  otpId.value = null;
  otpCode.value = "";
  password.value = "";
  step.value = "email";
  clearAlert();
};

const continueWithEmail = () => {
  if (!email.value) return;
  clearAlert();
  if (showPasswordForm.value) {
    step.value = "password";
  } else if (showOTP.value) {
    requestOTP();
  }
};

const changeEmail = () => {
  password.value = "";
  step.value = "email";
  clearAlert();
};

const goBackFromOTP = () => {
  otpId.value = null;
  otpCode.value = "";
  step.value = showPasswordForm.value && email.value ? "password" : "email";
  clearAlert();
};

const startMFA = (id, method) => {
  mfaId.value = id;
  mfaMethod.value = method;
  otpId.value = null;
  otpCode.value = "";
  password.value = "";
  step.value = "email";
  clearAlert();
  if (otpSecondFactor.value) {
    requestOTP();
  } else if (showPasswordForm.value && email.value) {
    step.value = "password";
  }
};

const handleAuthError = (error, method) => {
  console.error(error);
  const response = error.response;
  if (error.status === 401 && response?.mfaId && method) {
    startMFA(response.mfaId, method);
    return;
  }
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
    await pb.client.collection("users").authWithPassword(email.value, password.value, mfaOptions());
    await router.push("/");
  } catch (error) {
    handleAuthError(error, "password");
  } finally {
    isLoading.value = false;
  }
};

const requestOTP = async () => {
  if (!otpSecondFactor.value && !email.value) return;

  isSendingOTP.value = true;
  try {
    let res;
    if (otpSecondFactor.value) {
      res = await pb.client.send("/api/mfa/request-otp", {
        method: "POST",
        body: { mfaId: mfaId.value },
      });
      email.value = res.email;
    } else {
      res = await pb.client.collection("users").requestOTP(email.value);
    }
    otpId.value = res.otpId;
    otpCode.value = "";
    step.value = "otp";
    clearAlert();
  } catch (error) {
    handleAuthError(error);
  } finally {
    isSendingOTP.value = false;
  }
};

const loginWithOTP = async () => {
  if (!otpId.value || !otpCode.value) return;

  isLoading.value = true;
  try {
    await pb.client.collection("users").authWithOTP(otpId.value, otpCode.value, mfaOptions());
    await router.push("/");
  } catch (error) {
    handleAuthError(error, "otp");
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
    password.value = "";
    passwordConfirm.value = "";
    step.value = "password";
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
    await pb.client
      .collection("users")
      .authWithOAuth2({ provider: provider.name, ...mfaOptions() });
    await router.push("/");
  } catch (error) {
    handleAuthError(error, "oauth2");
  } finally {
    providerLoading.value = null;
  }
};
</script>
