<template>
  <q-page class="flex flex-center">
    <div
      v-if="data.resetPassword == 'N'"
      class="full-width success-verify"
      style="max-width: 500px"
    >
      <q-card class="">
        <q-card-section>
          <h5 class="q-my-md text-center">
            Hi <span class="text-capitalize">{{ data.firstName }}</span>
          </h5>
        </q-card-section>
        <q-separator inset />
        <q-card-section>
          <p class="text-subtitle1">
            Your email address has been successfully verified.
          </p>
          <p v-if="data.status === 'active'">
            Please sign in to use the services provided by {{ storeOrg.name }}.
          </p>
          <p v-else>
            Please wait for the admin to approve your account to use the
            services provided by {{ storeOrg.name }}.
          </p>
        </q-card-section>
        <q-card-section v-if="data.status === 'active'">
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
    <div
      v-if="data.resetPassword == 'Y'"
      class="full-width reset-password"
      style="max-width: 500px"
    >
      <q-card v-if="!resetPasswordSuccess">
        <q-card-section>
          <h5 class="q-my-md text-center">Hi {{ data.firstName }}</h5>
        </q-card-section>
        <q-separator inset />
        <q-card-section class="q-pa-none q-mt-md">
          <p class="flex flex-center q-ma-none">
            <span
              >Please fill out the following information to register on the "{{
                storeOrg.name
              }}".</span
            >
          </p>
        </q-card-section>
        <q-form ref="setPasswordForm">
          <q-card-section>
            <GeneralFormErrorComponent class="q-mb-sm" :errors="errors" />
            <div class="row wrap justify-start q-mb-sm q-col-gutter-sm">
              <q-input
                v-model="data.firstName"
                type="input"
                label="First Name"
                outlined
                :error="errors.inputErrors.password.error"
                :error-message="errors.inputErrors.password.message"
                :rules="[
                  (val) => (val && val.length > 0) || 'Please type something',
                ]"
                class="col-grow q-col-gutter-sm"
                hide-bottom-space
              />
              <q-input
                v-model="data.lastName"
                type="input"
                label="Last Name"
                outlined
                :error="errors.inputErrors.password.error"
                :error-message="errors.inputErrors.password.message"
                :rules="[
                  (val) => (val && val.length > 0) || 'Please type something',
                ]"
                class="col-grow q-col-gutter-sm"
                hide-bottom-space
              />
            </div>
            <q-input
              v-if="!isOTP"
              v-model="data.password"
              type="password"
              label="Password"
              outlined
              :error="errors.inputErrors.password.error"
              :error-message="errors.inputErrors.password.message"
              :rules="[
                (val) => (val && val.length > 0) || 'Please type something',
              ]"
              hide-bottom-space
              class="q-mb-sm"
            />
            <q-input
              v-if="!isOTP"
              v-model="data.confirmPassword"
              type="password"
              label="Confirm Password"
              outlined
              :error="errors.inputErrors.confirmPassword.error"
              :error-message="errors.inputErrors.confirmPassword.message"
              :rules="[
                (val) => (val && val.length > 0) || 'Please type something',
                (val) => data.password == val || 'Passwords do not match',
              ]"
              hide-bottom-space
            />
          </q-card-section>
          <q-card-section>
            <q-btn
              v-if="!isOTP"
              label="Set Password"
              class="full-width bg-black text-white"
              size="lg"
              square
              @click="setPassword"
            ></q-btn>
            <q-btn
              v-if="isOTP"
              label="Continue"
              class="full-width bg-black text-white"
              size="lg"
              square
              @click="setPassword"
            ></q-btn>
          </q-card-section>
        </q-form>
      </q-card>
      <q-card v-if="resetPasswordSuccess">
        <q-card-section>
          <h5 class="q-my-md text-center">
            Hi <span class="text-capitalize">{{ data.firstName }}</span>
          </h5>
        </q-card-section>
        <q-separator inset />
        <q-card-section>
          <p>
            You have successfully set up your password. Please log in with your
            <b>{{ data.email }}</b> email.
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
    <div v-if="error" class="full-width fail-verify" style="max-width: 500px">
      <q-card class="">
        <q-card-section>
          <h5 class="q-my-md text-center">Invalid token!</h5>
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
import axios from 'axios';
import { ref, reactive, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { api } from 'src/boot/axios';
import { Payload, FormErrors } from 'src/types';
import { useOrgStore } from 'src/stores/org-store';
import { handleServerErrors, resetServerErrors } from 'src/service/utils/error';
import GeneralFormErrorComponent from 'components/GeneralFormErrorComponent.vue';

const storeOrg = useOrgStore();

const isOTP = storeOrg.otp;

const route = useRoute();
const token = route.params.token;
const data = ref<{ [key: string]: any }>({});
const resetPassword = ref<boolean>(false);
const resetPasswordSuccess = ref<boolean>(false);
const loading = ref<boolean>(true);
const error = ref<boolean>(false);
const setPasswordForm = ref({});
const nameMissing = ref<boolean>(false);

const errors = reactive<FormErrors>({
  inputErrors: {
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

onMounted(async () => {
  loading.value = true;
  const payload: Payload<any> = {
    data: { token: token },
  };
  try {
    const res = await api().post('/user/verifyemail', payload);
    data.value = res.data.data;
    if (res.data.data.resetPassword === 'Y') {
      resetPassword.value = true;
    }
    if (res.data.data.firstName === '' && res.data.data.lastName === '') {
      console.log('gogo: ');
      nameMissing.value = true;
    }
  } catch (err) {
    error.value = true;
  }
  loading.value = false;
});

const setPassword = async () => {
  resetServerErrors(errors);
  const isValid = await setPasswordForm.value.validate();
  if (!isValid) {
    return;
  }
  const payload: Payload<any> = {
    data: {
      token: token,
      firstName: data.value.firstName,
      lastName: data.value.lastName,
      password: data.value.password,
      confirmPassword: data.value.confirmPassword,
    },
  };
  try {
    const res = await api().post('/user/setpassword', payload);
    resetPasswordSuccess.value = true;
  } catch (err) {
    handleServerErrors(err, errors);
  }
};
</script>

<style></style>
