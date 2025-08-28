<template>
  <div class="q-pa-md" style="max-width: 400px">
    <FormBuilder
      v-if="dataReady"
      :data="data"
      :settings="form"
      :handler="submitHandler"
      @onSubmit="formbuilderSubmitted"
      @onReset="formbuilderResetted"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onBeforeMount, reactive } from 'vue';
import { api } from 'src/boot/axios';

import FormBuilder from 'src/components/form/FormBuilderComponent.vue';
import { filterNonEditable } from 'src/service/utils/form-builder';
import { useOrgStore } from 'src/stores/org-store';
// import { useUserStore } from 'src/stores/user-store';
import { Payload, Response, StringAnyType, FormInput } from 'src/types';

// const storeUser = useUserStore();
const storeOrg = useOrgStore();
const dataReady = ref<boolean>(false);

const getAddr = '/user/user';
const putAddr = '/user/settings';
const form: FormInput[] = storeOrg.settings.forms.userMetadata;

const data = reactive<StringAnyType>({});

//handler to send form builder
const submitHandler = async (data: StringAnyType) => {
  const payload: Payload<StringAnyType> = {
    data: filterNonEditable(data, form),
  };
  return await api().put(putAddr, payload);
};

const formbuilderSubmitted = (data: any) => {
  return;
};

const formbuilderResetted = (data: any) => {
  return;
};

onBeforeMount(async () => {
  const res: Response<StringAnyType> = (await api().get(getAddr)).data;
  for (const [key, value] of Object.entries(res.data.metadata.userMetadata)) {
    data[key] = { name: '', value: value };
  }
  dataReady.value = true;
});
</script>

<style lang="sass"></style>
