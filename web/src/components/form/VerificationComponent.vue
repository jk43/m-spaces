<template>
  <q-dialog v-model="show" @hide="handleExit" persistent>
    <q-card style="min-width: 350px">
      <q-banner
        v-if="error.error"
        dense
        inline-actions
        class="text-white bg-red"
      >
        {{ error.message }}
      </q-banner>
      <q-banner v-if="message" dense inline-actions class="text-white bg-blue">
        {{ message }}
      </q-banner>
      <div></div>
      <q-card-section>
        <div class="text-h5">{{ data.data.title }}</div>
        <div class="text-h7">{{ data.data.message }}</div>
      </q-card-section>

      <q-card-section class="q-pt-none">
        <q-input
          v-if="!codeInputDisabled"
          v-model="code"
          dense
          autofocus
          :disabled="true"
          @keyup.enter="data.show = false"
          :rules="[
            (val) => (val && val.length > 0) || 'Please type verification code',
          ]"
        />
      </q-card-section>

      <q-card-actions align="right" class="text-primary">
        <q-btn flat label="Close" v-close-popup />
        <q-btn
          v-if="!resendBtnDisabled"
          flat
          label="Resend"
          @click="resend"
          :disabled="resendBtnDisabled"
        />
        <q-btn
          v-if="!verifyBtnDisabled"
          flat
          label="Verify"
          @click="runAction"
          :disabled="verifyBtnDisabled"
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
<script setup lang="ts">
import { defineProps, ref, watch, reactive, defineEmits } from 'vue';
import { api } from 'src/boot/axios';
import { Payload, StringAnyType, DialogClosedType } from 'src/types';

const emit = defineEmits(['close']);
const props = defineProps(['data']);
const data = reactive(props.data);
const show = ref(false);
const code = ref();
const error = reactive({ error: false, message: '' });
const message = ref('');
const resendBtnDisabled = ref(false);
const verifyBtnDisabled = ref(false);
const codeInputDisabled = ref(false);
let closeType: DialogClosedType = 'cancel';

const handleExit = () => {
  emit('close', data.data.verificationName, closeType);
  data.result = '';
  message.value = '';
  resendBtnDisabled.value = false;
  verifyBtnDisabled.value = false;
  codeInputDisabled.value = false;
};

const resend = async () => {
  try {
    error.error = false;
    const method = data.data.resend.method.toLowerCase();
    console.log('hello: ');
    const response = await api()[method](
      data.data.resend.URL,
      data.data.resend.payload
    );
    message.value = 'Verification code resent';
    disableResendBtn();
  } catch (error) {
    console.error(error);
  }
};

const disableResendBtn = () => {
  resendBtnDisabled.value = true;
  setTimeout(() => {
    resendBtnDisabled.value = false;
  }, 10000);
};

const runAction = async (event: any) => {
  try {
    error.error = false;
    error.message = '';
    message.value = '';
    const payload: Payload<StringAnyType> = {
      data: { [data.data.keyName]: code.value, ...data.data.payload.data },
    };
    const response = await api().post(data.data.url, payload);
    message.value = response.data.data.title;
    data.data.message = response.data.data.message;
    resendBtnDisabled.value = true;
    verifyBtnDisabled.value = true;
    codeInputDisabled.value = true;
    code.value = '';
    closeType = 'ok';
    emit('close', data.data.verificationName);
  } catch (err) {
    error.error = true;
    error.message = err.response.data.data.error;
  }
  event.stopPropagation();
};

watch(
  () => data.result,
  (value) => {
    show.value = false;
    if (value === 'verification_required') {
      show.value = true;
    }
  }
);
</script>
