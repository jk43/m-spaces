<template>
  <q-form @submit="onSubmit" class="q-gutter-xs">
    <q-input
      filled
      v-model="value.firstName"
      label="First Name"
      lazy-rules
      :rules="[(val) => (val && val.length > 0) || 'Please type something']"
    />
    <q-input
      filled
      v-model="value.lastName"
      label="Last Name"
      lazy-rules
      :rules="[(val) => (val && val.length > 0) || 'Please type something']"
    />
    <q-input
      filled
      v-model="value.email"
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
import { defineProps, ref, defineEmits, watch } from 'vue';

const prop = defineProps(['user']);
const emit = defineEmits(['emitted']);
const value = ref(prop.user);
let clone: any = {};

watch(
  () => prop.user,
  (params) => {
    value.value = params;
    clone = JSON.parse(JSON.stringify(params));
  }
);

const emitData = () => {
  emit('emitted', { name: 'josef', email: 'jk@jktech.net' });
};

const onReset = () => {
  console.log('clonecloneclone: ', clone);
  value.value.firstName = clone.firstName;
  value.value.lastName = clone.lastName;
  value.value.email = clone.email;
};

const onSubmit = () => {
  console.log('onsubmit: ');
  emit('reset', value.value);
};
</script>

<style scoped>
/* Your component-specific styles go here */
</style>
