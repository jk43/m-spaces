<template>
  <q-form @submit="onSubmit" @reset="onReset" class="q-gutter-md">
    <!-- <q-input
      filled
      v-for="item in value.settings"
      v-model="item.value"
      :key="item.key"
      :label="item.name"
      lazy-rules
      :rules="[(val) => (val && val.length > 0) || 'Please type something']"
    />
      <component
      v-for="(form, index) in value.settings"
      :key="index"
      :is="getComponent(1111)"
      :label="form.name"
      placeholder="Enter something"
    ></component>
   -->

    <!-- <q-toggle v-model="accept" label="I accept the license and terms" /> -->
    <div>
      <q-btn label="Submit" type="submit" color="primary" />
      <q-btn label="Reset" type="reset" color="primary" flat class="q-ml-sm" />
    </div>
  </q-form>
  <div @click="emitData">emit</div>
</template>

<script setup lang="ts">
import { defineProps, ref, defineEmits, watch } from 'vue';

const props = defineProps(['user', 'formData']);
const user = ref(props.user);
const formData = ref(props.formData);
console.log('formData: ', props.formData);
const emit = defineEmits(['emitted']);
let clone: any = {};

watch(
  () => props.user,
  (params) => {
    user.value = params;
  }
);

watch(
  () => props.formData,
  (params) => {
    formData.value = params;
  }
);

const emitData = () => {
  emit('emitted', { name: 'josef', email: 'jk@jktech.net' });
};

const onReset = () => {
  //console.log('onReset', value.value.settings[0].key);
};

const onSubmit = () => {
  console.log('user', user.value);
  console.log('formData', formData.value.forms.userSettings);
};
const getComponent = (type: string) => {
  return 'q-input';
};
</script>

<style scoped>
/* Your component-specific styles go here */
</style>
