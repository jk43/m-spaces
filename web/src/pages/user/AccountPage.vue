<template>
  <div class="q-pa-md" style="max-width: 400px">
    <FormBuilder
      :data="accountData"
      :settings="accountSettings"
      :handler="accountHandler"
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
import { ref, onBeforeMount, reactive } from 'vue';
import { api } from 'src/boot/axios';
import FormBuilder from 'src/components/form/FormBuilderComponent.vue';
import VerificationComponent from 'src/components/form/VerificationComponent.vue';
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

const storeUser = useUserStore();

const closedVerification = (name: string, closeType: DialogCloseType) => {
  return;
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

const accountSettings = [
  {
    discription: '',
    key: 'firstName',
    name: 'First Name',
    options: [],
    rules: ['Required'],
    type: 'input',
    editable: !storeUser.oauth,
  },
  {
    discription: '',
    key: 'lastName',
    name: 'Last Name',
    options: [],
    rules: ['Required'],
    type: 'input',
    editable: !storeUser.oauth,
  },
  {
    discription: '',
    key: 'email',
    name: 'Email',
    options: [],
    rules: ['Required', 'Email'],
    type: 'input',
    editable: !storeUser.oauth,
  },
];

const accountData = reactive({
  firstName: { name: '', value: null },
  lastName: { name: '', value: null },
  email: { name: '', value: null },
});

const metadataData = reactive<StringAnyType>({});

let user = reactive<User>({});
const userDataReady = ref<boolean>(false);

const accountHandler = async (data: StringAnyType) => {
  const payload: Payload<StringAnyType> = {
    data: data,
  };
  try {
    const res = await api().put('/user/account', payload);
    //if email has changed, send verification code
    if (res.data.result === 'verification_required') {
      verificationPrompt.result = res.data.result;
      verificationPrompt.data = res.data.data;
      verificationPrompt.data.verificationName = 'Password';
      verificationPrompt.data.resend = {
        payload: payload,
        method: 'put',
        URL: '/user/account',
      };
    }
  } catch (error) {
    console.log('error');
    throw error;
  }
};

const formbuilderSubmitted = (data: any) => {
  console.log('Debugging - formbuilderSubmitted: ');
};

const formbuilderResetted = (data: any) => {
  console.log('Debugging - formbuilderResetted: ');
};

onBeforeMount(async () => {
  for (const i in accountData) {
    accountData[i].value = storeUser[i];
  }
  userDataReady.value = true;
});
</script>

<style lang="sass"></style>
