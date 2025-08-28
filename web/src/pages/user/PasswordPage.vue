<template>
  <div class="q-pa-md" style="max-width: 400px">
    <FormBuilder
      :data="passwordData"
      :settings="passwordSettings"
      :handler="passwordHandler"
      @onSubmit="formbuilderSubmitted"
      @onReset="formbuilderResetted"
    />
  </div>
  <VerificationComponent
    :data="verificationPrompt"
    @close="closedVerification"
  />
</template>

<script setup lang="ts">
import { ref, watch, onBeforeMount, reactive, defineEmits } from 'vue';
import { Component } from 'vue';
import { useRoute } from 'vue-router';
import { api } from 'src/boot/axios';
import FormBuilder from 'src/components/form/FormBuilderComponent.vue';
import VerificationComponent from 'src/components/form/VerificationComponent.vue';
//import { useOrgStore } from 'src/stores/org-store';
import { useUserStore } from 'src/stores/user-store';
import {
  User,
  UserSetting,
  Payload,
  Response,
  StringAnyType,
  VerificationInstruction,
  DialogCloseType,
} from 'src/types';

const route = useRoute();
const userStore = useUserStore();
const closedVerification = (name: string, closeType: DialogCloseType) => {
  passwordData.currentPassword.value = null;
  passwordData.password.value = null;
  passwordData.confirmPassword.value = null;
};

let verificationPrompt: Response<VerificationInstruction> = reactive({
  result: '',
  data: {
    verificationName: '',
    URL: '',
    keyName: '',
    method: '',
    message: '',
    resend: {
      payload: {},
      method: '',
      URL: '',
    },
  },
});

const passwordSettings = [
  {
    discription: '',
    key: 'currentPassword',
    name: 'Current Password',
    options: [],
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'Password',
    editable: !userStore.oauth,
    rules: ['Required'],
  },
  {
    discription: '',
    key: 'password',
    name: 'Password',
    options: [],
    //rules: "[(val) => (val && val.length > 0) || 'Last Name cannot be empty']",
    type: 'Password',
    editable: !userStore.oauth,
    rules: ['Required'],
  },
  {
    discription: '',
    key: 'confirmPassword',
    name: 'Confirm Password',
    options: [],
    //rules: "[(val) => (val && val.length > 0) || 'Email cannot be empty']",
    type: 'Password',
    editable: !userStore.oauth,
    rules: ['Required'],
  },
];

const passwordData = reactive<StringAnyType>({
  currentPassword: { name: '', value: null },
  password: { name: '', value: null },
  confirmPassword: { name: '', value: null },
});

const passwordHandler = async (data: StringAnyType) => {
  const payload: Payload<StringAnyType> = {
    data: data,
  };
  try {
    const res = await api().put('/user/password', payload);
    verificationPrompt.result = res.data.result;
    verificationPrompt.data = res.data.data;
    verificationPrompt.data.verificationName = 'Password';
    verificationPrompt.data.resend = {
      payload: payload,
      method: 'put',
      URL: '/user/password',
    };
  } catch (error) {
    console.log(':errr ');
    throw error;
  }
};
const formbuilderSubmitted = (data: any) => {
  return;
};

const formbuilderResetted = (data: any) => {
  return;
};
</script>

<style lang="sass"></style>
