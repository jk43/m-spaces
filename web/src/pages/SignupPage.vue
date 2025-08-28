<template>
  <q-page class="flex flex-center">
    <div class="full-width" style="max-width: 500px">
      <q-card v-if="!registrationSuccess">
        <q-form ref="singUpForm">
          <q-card-section>
            <h5 class="q-my-md text-center">{{ storeOrg.name }}</h5>
            <GeneralFormErrorComponent class="q-mb-sm" :errors="errors" />
            <div class="row wrap justify-start q-mb-sm q-col-gutter-sm">
              <q-input
                hide-bottom-space
                outlined
                v-model="user.firstName"
                label="First Name"
                class="col-grow q-col-gutter-sm"
                :rules="[
                  (val) => (val && val.length > 0) || 'Please type something',
                ]"
              />
              <q-input
                hide-bottom-space
                outlined
                v-model="user.lastName"
                label="Last Name"
                class="col-grow q-col-gutter-sm"
                :rules="[
                  (val) => (val && val.length > 0) || 'Please type something',
                ]"
              />
            </div>
            <q-input
              hide-bottom-space
              outlined
              v-model="user.email"
              label="Email"
              class="q-mb-sm"
              type="email"
              :error-message="errors.inputErrors.email.message"
              :error="errors.inputErrors.email.error"
              :rules="[
                (val) => (val && val.length > 0) || 'Please type something',
              ]"
            />
            <q-input
              v-if="!storeOrg.otp"
              hide-bottom-space
              outlined
              v-model="user.password"
              label="Password"
              type="password"
              class="q-mb-sm"
              :error-message="errors.inputErrors.password.message"
              :error="errors.inputErrors.password.error"
              :rules="[
                (val) => (val && val.length > 0) || 'Please type something',
              ]"
            />
            <q-input
              v-if="!storeOrg.otp"
              hide-bottom-space
              outlined
              v-model="user.confirmPassword"
              label="Confirm Password"
              type="password"
              class="q-mb-sm"
              :error-message="errors.inputErrors.confirmPassword.message"
              :error="errors.inputErrors.confirmPassword.error"
              :rules="[
                (val) => (val && val.length > 0) || 'Please type something',
                (val) => val == user.password || 'Password is not matching',
              ]"
            />
          </q-card-section>
          <q-btn
            @click="register"
            label="Create Account"
            class="full-width bg-black text-white"
            size="lg"
            square
          ></q-btn>
          <q-card-section>
            <span
              >Have a account?
              <router-link to="/auth/login">Sign in now!</router-link></span
            >
          </q-card-section>
        </q-form>
      </q-card>
      <q-card v-else>
        <q-card-section>
          <h5 class="text-center">Registration Success</h5>
          <p class="text-center" v-if="!storeOrg.otp">
            Please check your email to verify your account
          </p>
          <p class="text-center" v-else>
            Please click the button below to sign in
          </p>
        </q-card-section>
        <q-card-section>
          <router-link to="/auth/login">
            <q-btn
              label="Sign in"
              class="full-width bg-black text-white"
              size="lg"
              square
            ></q-btn>
          </router-link>
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { api } from 'src/boot/axios';
import { reactive, ref, onBeforeMount } from 'vue';
import { useRouter } from 'vue-router';
import { FormErrors, Payload } from 'src/types';
import { QForm } from 'quasar';
import { handleServerErrors, resetServerErrors } from 'src/service/utils/error';

import GeneralFormErrorComponent from 'components/GeneralFormErrorComponent.vue';
import { useOrgStore } from 'src/stores/org-store';

const storeOrg = useOrgStore();
const router = useRouter();
const singUpForm = ref<QForm>(null);
const user = reactive({
  firstName: '',
  lastName: '',
  email: '',
  password: '',
  confirmPassword: '',
});

onBeforeMount(() => {
  if (!storeOrg.allowSelfMemberRegistration) {
    router.push('/auth/login');
  }
});

const errors = reactive<FormErrors>({
  inputErrors: {
    email: {
      error: false,
      message: '',
    },
    password: {
      error: false,
      message: '',
    },
    confirmPassword: {
      error: false,
      message: '',
    },
  },
  formErrors: [],
});
const registrationSuccess = ref<boolean>(false);

let register = async () => {
  resetServerErrors(errors);
  const isValid = await singUpForm.value.validate();
  if (!isValid) {
    return;
  }
  try {
    const payload: Payload<any> = {
      data: user,
    };
    const res = await api().post('/user/save', payload);
    registrationSuccess.value = true;
  } catch (err) {
    handleServerErrors(err, errors);
  }
  console.log('errors: ', errors);
};
</script>

<style></style>
