<template>
  <q-page class="flex flex-center">
    <div class="full-width" style="max-width: 400px">
      <q-card class="">
        <q-card-section>
          <h5 class="q-my-md text-center">{{ storeOrg.name }}</h5>
          <q-banner
            v-if="banner"
            dense
            inline-actions
            class="text-black bg-white"
          >
            <h7>
              You have successfully changed your password. Please log in with
              your newly changed password on the login page.
            </h7>
          </q-banner>
          <GeneralFormErrorComponent :errors="errors" class="q-mb-sm" />
          <q-input
            v-if="!banner"
            hide-bottom-space
            v-model="password"
            outlined
            label="New Password"
            type="password"
            class="q-mb-sm"
            :rules="[
              (val) => (val && val.length > 0) || 'Please type something',
            ]"
          />
          <q-input
            v-if="!banner"
            hide-bottom-space
            v-model="password1"
            outlined
            label="Confirm Password"
            type="password"
            class="q-mb-sm"
            :rules="[
              (val) => (val && val.length > 0) || 'Please type something',
            ]"
          />
        </q-card-section>
        <q-btn
          v-if="!banner"
          @click="send"
          label="Reset Password"
          class="full-width bg-black text-white"
          size="lg"
          square
        ></q-btn>
        <q-card-section>
          <span> <router-link to="/auth/login">Login page</router-link></span>
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRoute } from 'vue-router';
import { api } from 'src/boot/axios';
import { handleServerErrors, resetServerErrors } from 'src/service/utils/error';
import { StringAnyType, Payload, FormErrors } from 'src/types';
import { useOrgStore } from 'src/stores/org-store';
import GeneralFormErrorComponent from 'components/GeneralFormErrorComponent.vue';

const storeOrg = useOrgStore();
const route = useRoute();

const email = ref('');
const password = ref('');
const password1 = ref('');
const banner = ref(false);

let errors = reactive<FormErrors>({
  inputErrors: {},
  formErrors: [],
});

const send = async () => {
  resetServerErrors(errors);
  const payload: Payload<StringAnyType> = {
    data: {
      token: { value: route.params.token, name: 'Token' },
      password: { value: password.value, name: 'Password' },
    },
  };
  try {
    const result = await api().put('/auth/password', payload);
    banner.value = true;
  } catch (error) {
    handleServerErrors(error, errors);
  }
};
</script>

<style></style>
