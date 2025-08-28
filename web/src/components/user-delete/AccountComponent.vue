<template>
  <q-form @submit="onSubmit" class="q-gutter-xs">
    <q-input
      filled
      v-model="user.firstName"
      label="First Name"
      lazy-rules
      :rules="[(val) => (val && val.length > 0) || 'Please type something']"
    />
    <q-input
      filled
      v-model="user.lastName"
      label="Last Name"
      lazy-rules
      :rules="[(val) => (val && val.length > 0) || 'Please type something']"
    />
    <q-input
      filled
      v-model="user.email"
      label="Email"
      lazy-rules
      :rules="[(val) => (val && val.length > 0) || 'Please type something']"
    />
    <div>
      <q-btn label="Submit" type="submit" color="primary" />
      <q-btn
        label="Reset"
        type="reset"
        color="primary"
        flat
        class="q-ml-sm"
        @click="onReset"
      />
    </div>
  </q-form>
  <div @click="emitData">emit</div>
</template>

<script setup lang="ts">
import {
  defineProps,
  reactive,
  defineEmits,
  watch,
  onUnmounted,
  onMounted,
  ref,
} from 'vue';

import { api } from 'src/boot/axios';
import { User } from 'src/types';

const prop = defineProps(['user']);
const emit = defineEmits(['emitted']);
let user = reactive<User>(prop.user);
let resetData: User = {};

onMounted(() => {
  resetData = JSON.parse(JSON.stringify(user));
});

const onReset = () => {
  user.firstName = resetData.firstName;
  user.lastName = resetData.lastName;
  user.email = resetData.email;
};

const onSubmit = async () => {
  let res;
  try {
    res = await api().put('/user/user', user);
  } catch (err) {
    console.log('err: ', err);
    return err;
  }
  console.log('prop.user', prop.user);
};
</script>

<style scoped>
/* Your component-specific styles go here */
</style>
