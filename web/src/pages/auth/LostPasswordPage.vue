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
            <h7
              >An email for password change has been sent. Please check your
              email and follow the next steps.</h7
            >
          </q-banner>
          <q-form v-if="show" ref="loginForm">
            <q-input
              hide-bottom-space
              v-model="data.email.value"
              outlined
              label="Email"
              class="q-mb-md"
              :rules="[
                (val) => (val && val.length > 0) || 'Please type something',
              ]"
            />
          </q-form>
        </q-card-section>
        <q-btn
          v-if="show"
          @click="send"
          label="Send"
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
import { api } from 'src/boot/axios';
import { StringAnyType, Payload } from 'src/types';
import { useOrgStore } from 'src/stores/org-store';

const storeOrg = useOrgStore();

const data = reactive<StringAnyType>({ email: { value: '', name: 'Email' } });
const banner = ref(false);
const show = ref(true);

const send = async () => {
  const payload: Payload<StringAnyType> = {
    data: data,
  };
  const result = await api().post('/auth/forgotpassword', payload);
  banner.value = true;
  data.email.value = '';
  show.value = false;
};
</script>

<style></style>
