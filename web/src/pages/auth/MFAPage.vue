<template>
  <q-page class="flex flex-center">
    <div class="full-width" style="max-width: 600px">
      <q-card class="">
        <q-card-section>
          <h5 class="q-my-md text-center">{{ storeOrg.name }}</h5>
          <q-banner dense inline-actions class="text-black bg-white">
            <h7>
              Your account is protected with two-factor authentication. We've
              sent you a email with a code. Please enter the code below.
            </h7>
          </q-banner>
          <GeneralFormErrorComponent :errors="errors" class="q-mb-sm" />
          <q-input
            hide-bottom-space
            v-model="code"
            outlined
            label="Two-Factor Code"
            type="text"
            class="q-mb-sm"
            :rules="[
              (val) => (val && val.length > 0) || 'Please type something',
            ]"
          />
        </q-card-section>
        <q-btn
          @click="send"
          label="Continue"
          class="full-width bg-black text-white"
          size="lg"
          square
        ></q-btn>
        <q-card-section>
          <span> <a @click="resendCode">Resend Code</a></span>
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { api } from 'src/boot/axios';
import { handleServerErrors, resetServerErrors } from 'src/service/utils/error';
import {
  StringAnyType,
  Payload,
  FormErrors,
  LoginSuccessResponse,
} from 'src/types';
import GeneralFormErrorComponent from 'components/GeneralFormErrorComponent.vue';
import { refreshTokenInterval } from 'src/service/user/auth';
import { getOrganizationSettings } from 'src/service/organization/setting';
import { useOrgStore } from 'src/stores/org-store';
import { useUserStore } from 'src/stores/user-store';

const storeUser = useUserStore();
const router = useRouter();
const route = useRoute();
const storeOrg = useOrgStore();

const code = ref('');
const token = ref(route.params.token);

let errors = reactive<FormErrors>({
  inputErrors: {},
  formErrors: [],
});

const send = async () => {
  resetServerErrors(errors);
  const payload: Payload<StringAnyType> = {
    data: {
      token: token.value,
      code: code.value,
    },
  };
  try {
    const result = await api().post('/auth/verify-mfa-code', payload);
    storeUser.setUserInfo(result.data.data as LoginSuccessResponse);
    refreshTokenInterval(5, null);
    getOrganizationSettings();
    router.push('/user/dashboard');
  } catch (error) {
    handleServerErrors(error, errors);
  }
};

const resendCode = async () => {
  resetServerErrors(errors);
  const payload: Payload<StringAnyType> = {
    data: {
      token: token.value,
    },
  };
  try {
    const result = await api().post('/auth/resend-mfa-code', payload);
    token.value = result.data.data.mfaToken;
  } catch (error) {
    handleServerErrors(error, errors);
  }
};
</script>

<style></style>
