<template>
  <q-form @submit="onSubmit">
    <div class="column q-gutter-sm">
      <div class="col"><GeneralFormErrorComponent :errors="errors" /></div>
      <div v-for="(item, index) in formSettings" :key="index" class="col">
        <component
          :is="item.type"
          :data="item.data"
          :settings="item.settings"
          :error="item.error"
          :globalData="reactiveData"
          :globalSettings="formSettings"
          :submitClicked="submitClicked"
        ></component>
      </div>
      <div v-if="hasEditable && !$slots.buttons" class="col">
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
      <slot name="buttons"></slot>
    </div>
  </q-form>
</template>

<script setup lang="ts">
import { defineProps, ref, shallowRef, reactive, onMounted } from 'vue';
import { v4 as uuidv4 } from 'uuid';
import GeneralFormErrorComponent from 'components/GeneralFormErrorComponent.vue';
import Input from './inputs/InputComponent.vue';
import DynamicLabelValue from './inputs/DynamicLabelAndValue.vue';
import Rules from './inputs/RulesComponent.vue';
import Select from './inputs/SelectComponent.vue';
import Radio from './inputs/RadioComponent.vue';
import Checkbox from './inputs/CheckboxComponent.vue';
import Toggle from './inputs/ToggleComponent.vue';
import OptionGroup from './inputs/OptionGroupComponent.vue';
import ToggleReversed from './inputs/ToggleReversedComponent.vue';
import Password from './inputs/PasswordComponent.vue';
import PasswordWithConfirm from './inputs/PasswordWithConfirmComponent.vue';
import TextArea from './inputs/TextAreaComponent.vue';
import Tree from './inputs/TreeComponent.vue';
import TreeSelect from './inputs/TreeSelectComponent.vue';
import TipTap from './inputs/TipTapComponent.vue';

import {
  FormBuilderSetting,
  FormErrors,
  ErrorResponses,
  FormError,
  StringAnyType,
} from 'src/types';
import { resetServerErrors } from 'src/service/utils/error';

const props = defineProps(['data', 'settings', 'handler']);
const reactiveData = reactive<StringAnyType>(props.data);
const formSettings = reactive<FormBuilderSetting[]>([]);
const loaded = ref(false);
const hasEditable = ref(false);
const submitClicked = ref<number>(0);

const emit = defineEmits(['onSubmit', 'onReset']);

const errors = reactive<FormErrors>({
  inputErrors: {},
  formErrors: [],
});

let resetData: StringAnyType = {};

onMounted(() => {
  reactiveData.formId = { name: 'formId', value: uuidv4() };
  formatSettings();
  resetData = JSON.parse(JSON.stringify(reactiveData));
});

const formatSettings = () => {
  const settings = props.settings;
  let index = 0;
  for (const s of settings) {
    if (s.show === false) {
      continue;
    }
    let qElement: any;
    let value;
    //select custom element
    switch (s.type) {
      case 'Input':
        qElement = Input;
        break;
      case 'Options':
        qElement = DynamicLabelValue;
        break;
      case 'Rules':
        qElement = Rules;
        break;
      case 'Select':
        qElement = Select;
        break;
      case 'Radio':
        qElement = Radio;
        break;
      case 'Checkbox':
        qElement = Checkbox;
        break;
      case 'Toggle':
        qElement = Toggle;
        break;
      case 'Option Group':
        qElement = OptionGroup;
        break;
      case 'Toggle Reversed':
        qElement = ToggleReversed;
        break;
      case 'Password':
        qElement = Password;
        break;
      case 'Password With Confirm':
        qElement = PasswordWithConfirm;
        break;
      case 'TextArea':
        qElement = TextArea;
        break;
      case 'Tree':
        qElement = Tree;
        break;
      case 'TreeSelect':
        qElement = TreeSelect;
        break;
      case 'Advanced Editor':
        qElement = TipTap;
        break;
      default:
        qElement = Input;
    }
    if (s.editable) {
      hasEditable.value = true;
    }
    let inData = s.key in reactiveData;
    let valueInReactiveData: any;
    if (s.multiple) {
      const valueInReactiveData = [[s.defaultValue]];
    } else {
      const valueInReactiveData = '';
    }
    if (!inData) {
      reactiveData[s.key] = {
        name: s.name,
        value: valueInReactiveData,
      };
    } else {
      reactiveData[s.key].name = s.name;
    }
    //to hide or show the input
    if (s.show === undefined) {
      s.show = true;
    }
    const formSetting: FormBuilderSetting = {
      type: shallowRef(qElement),
      settings: s,
      data: reactiveData[s.key],
      error: { error: null, message: '' },
    };
    if (formSettings[index]) {
      formSettings[s] = formSetting;
    } else {
      formSettings.push(formSetting);
    }
    index++;
  }
};

const onSubmit = async () => {
  resetErrors();
  try {
    const ret = await props.handler(reactiveData);
    // to reset formId to reload files
    submitClicked.value = ++submitClicked.value;
    reactiveData.formId = { name: 'formId', value: uuidv4() };
    emit('onSubmit', ret);
  } catch (err: any) {
    console.log('errerr: ', err);
    if (err?.response?.data?.data) {
      processErrors(err.response.data.data);
    } else if (err?.response?.data) {
      processErrors(err.response.data);
    } else if (err?.message) {
      errors.formErrors.push({
        message: err.message,
        error: true,
      } as FormError);
    } else {
      errors.formErrors.push({
        message: 'Unknown error occurred.',
        error: true,
      } as FormError);
    }
  }
};

const processErrors = (err: ErrorResponses) => {
  if (!Array.isArray(err)) {
    errors.formErrors.push({
      message: err.error,
      error: true,
    } as FormError);
    return;
  }
  for (const i in err) {
    if (!err[i].field) {
      errors.formErrors.push({
        message: err[i].error,
        error: true,
      } as FormError);
      continue;
    }
    const field = err[i].field;
    for (const j in formSettings) {
      if (formSettings[j].settings.key === field) {
        formSettings[j].error.error = true;
        formSettings[j].error.message = err[i].error;
      }
    }
  }
};

const resetErrors = () => {
  for (const i in formSettings) {
    formSettings[i].error.error = false;
    formSettings[i].error.message = '';
  }
  errors.formErrors = [];
};

const onReset = () => {
  for (const i in resetData) {
    reactiveData[i].value = resetData[i].value;
  }
  emit('onReset');
};
</script>

<style scoped></style>
