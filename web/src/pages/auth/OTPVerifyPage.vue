<template>
  <q-page class="flex flex-center">
    <div class="full-width" style="max-width: 400px">
      <q-card class="">
        <q-card-section>
          <h5 class="q-my-md text-center">{{ storeOrg.name }}</h5>
          <q-banner dense inline-actions class="text-black bg-white">
            <h7> We just emailed your passcode. </h7>
          </q-banner>
          <GeneralFormErrorComponent :errors="errors" class="q-mb-md" />
          <q-form ref="loginForm">
            <q-input
              hide-bottom-space
              :error-message="errors.inputErrors.email.message"
              :error="errors.inputErrors.email.error"
              ref="otp-code"
              v-model="credentials.code"
              outlined
              label="One Time Passcode"
              class="q-mb-md"
              :rules="[
                (val) => (val && val.length > 0) || 'Please type something',
              ]"
            />
          </q-form>
        </q-card-section>
        <q-btn
          @click="login"
          label="Continue"
          class="full-width bg-black text-white"
          size="lg"
          square
        ></q-btn>
        <q-card-section>
          <a @click="resendCode"> Resend One Time Passcode </a>
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
import { useRouter, useRoute } from 'vue-router';

import GeneralFormErrorComponent from 'components/GeneralFormErrorComponent.vue';
import { useOrgStore } from 'src/stores/org-store';

const storeUser = useUserStore();
const router = useRouter();
const route = useRoute();
const storeOrg = useOrgStore();

const token = ref(route.params.token);

let credentials = reactive<{
  code: string;
  token: string;
}>({
  code: '',
  token: '',
});

let errors = reactive<FormErrors>({
  inputErrors: {
    email: {
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
    credentials.token = token.value;
    const payload: Payload<any> = {
      data: credentials,
    };
    const result = await api().post<Payload<LoginResponse>>(
      '/auth/verify-otp',
      payload
    );
    // if (result.data.data.token) {
    //   router.push(`/auth/mfa/${result.data.data.mfaToken}`);
    //   return;
    // }
    // console.log('Debugging - result.data.data: ', result.data.data);
    storeUser.setUserInfo(result.data.data as LoginSuccessResponse);
    refreshTokenInterval(5, null);
    getOrganizationSettings();
    router.push('/user/dashboard');
    return;
  } catch (err) {
    handleServerErrors(err, errors);
  }
};

let resendCode = async () => {
  const result = await api().post('/auth/otp', {
    data: { token: token.value },
  });
  token.value = result.data.data.token;
};
</script>

<style></style>
