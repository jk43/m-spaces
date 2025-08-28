<template>
  <q-page class="flex flex-center">
    <div class="full-width" style="max-width: 400px">
      <q-card class="">
        <q-card-section>
          <h5 class="q-my-md text-center">{{ storeOrg.name }}</h5>
          <GeneralFormErrorComponent :errors="errors" class="q-mb-md" />
          <q-form ref="loginForm" @submit="login">
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

            <!-- 로그인 버튼을 폼 안으로 이동 -->
            <q-btn
              type="submit"
              label="login"
              class="full-width bg-black text-white"
              size="lg"
              square
            />
          </q-form>
        </q-card-section>

        <q-card-section>
          <router-link to="/oauth/auth/google">Google Login</router-link><br />
          <router-link to="/oauth/auth/microsoft">Microsoft Login</router-link
          ><br />
          <router-link to="/oauth/auth/facebook">Facebook Login</router-link>
        </q-card-section>
        <q-card-section>
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
import { useRouter, useRoute } from 'vue-router';

import GeneralFormErrorComponent from 'components/GeneralFormErrorComponent.vue';
import { useOrgStore } from 'src/stores/org-store';

const storeUser = useUserStore();
const router = useRouter();
const route = useRoute();
const storeOrg = useOrgStore();

const token = ref(route.params.token as string);

let credentials = reactive<{
  email: string;
}>({
  email: '',
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

const loginForm = ref<QForm>();

// 중복 호출 방지를 위한 플래그
let isSubmitting = false;

let login = async () => {
  // 중복 호출 방지
  if (isSubmitting) {
    return;
  }

  isSubmitting = true;
  resetServerErrors(errors);

  try {
    const isValid = await loginForm.value?.validate();
    if (!isValid) {
      isSubmitting = false;
      return;
    }

    const payload: Payload<any> = {
      data: credentials,
    };
    const result = await api().post<Payload<LoginResponse>>(
      '/auth/otp',
      payload
    );

    // LoginSuccessResponse 타입 체크
    if ('token' in result.data.data) {
      router.push('/auth/otp/' + result.data.data.token);
    } else {
      console.error('Unexpected response format');
    }
  } catch (err) {
    handleServerErrors(err, errors);
  } finally {
    isSubmitting = false;
  }
};
</script>

<style></style>
