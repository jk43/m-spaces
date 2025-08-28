<template>
  <q-page class="flex flex-center">
    <div class="full-width" style="max-width: 400px">
      <q-card class="">
        <q-card-section>
          <h5 class="q-my-md text-center">{{ storeOrg.name }}</h5>
          <GeneralFormErrorComponent :errors="errors" class="q-mb-md" />
          <q-form ref="loginForm">
            <q-input
              hide-bottom-space
              :error-message="errors.inputErrors.email.message"
              :error="errors.inputErrors.email.error"
              ref="email"
              v-model="credentials.email"
              outlined
              label="Email"
              class="q-mb-md"
              :rules="[
                (val) => (val && val.length > 0) || 'Please type something',
              ]"
            />
            <q-input
              hide-bottom-space
              :error-message="errors.inputErrors.password.message"
              :error="errors.inputErrors.password.error"
              ref="password"
              v-model="credentials.password"
              outlined
              label="Password"
              type="password"
              class="q-mb-sm"
              :rules="[
                (val) => (val && val.length > 0) || 'Please type something',
              ]"
            />
          </q-form>
        </q-card-section>
        <q-btn
          @click="login"
          label="login"
          class="full-width bg-black text-white"
          size="lg"
          square
        ></q-btn>
        <q-card-section>
          <router-link to="/oauth/auth/google">Google Login</router-link><br />
          <router-link to="/oauth/auth/microsoft">Microsoft Login</router-link
          ><br />
          <router-link to="/oauth/auth/facebook">Facebook Login</router-link
          ><br />
        </q-card-section>
        <q-card-section>
          <router-link to="/auth/forgotpassword">Forgot Password?</router-link>
          <div v-if="storeOrg.allowSelfMemberRegistration">
            Don't have a account?
            <router-link to="/signup">Sign up now!</router-link>
          </div>
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { reactive, ref, onBeforeMount } from 'vue';
import { api } from 'src/boot/axios';
import { useUserStore } from 'src/stores/user-store';
import {
  LoginResponse,
  LoginSuccessResponse,
  FormErrors,
  Payload,
} from 'src/types';
import { QForm } from 'quasar';
import { handleServerErrors, resetServerErrors } from 'src/service/utils/error';
import { refreshTokenInterval } from 'src/service/user/auth';
import { getOrganizationSettings } from 'src/service/organization/setting';
import { useRouter } from 'vue-router';

import GeneralFormErrorComponent from 'components/GeneralFormErrorComponent.vue';
import { useOrgStore } from 'src/stores/org-store';

const storeUser = useUserStore();
const router = useRouter();
const storeOrg = useOrgStore();

let credentials = reactive<{
  email: string;
  password: string;
}>({
  email: '',
  password: '',
});

onBeforeMount(() => {
  if (storeOrg.otp) {
    router.push('/auth/otp');
  }
});

let errors = reactive<FormErrors>({
  inputErrors: {
    email: {
      error: false,
      message: '',
    },
    password: {
      error: false,
      message: '',
    },
    general: {
      error: false,
      message: '',
    },
  },
  formErrors: [],
});

const loginForm = ref<QForm>(null);

let login = async () => {
  resetServerErrors(errors);
  const isValid = await loginForm.value.validate();
  if (!isValid) {
    return;
  }
  try {
    const payload: Payload<any> = {
      data: credentials,
    };
    const result = await api().post<Payload<LoginResponse>>(
      '/user/login',
      payload
    );
    if (result.data.data.mfa) {
      router.push(`/auth/mfa/${result.data.data.mfaToken}`);
      return;
    }
    storeUser.setUserInfo(result.data.data as LoginSuccessResponse);
    refreshTokenInterval(5, null);
    getOrganizationSettings();
    router.push('/user/dashboard');
    return;
  } catch (err) {
    handleServerErrors(err, errors);
  }
};
</script>

<style></style>
