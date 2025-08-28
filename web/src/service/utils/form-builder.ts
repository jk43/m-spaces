import { watch, ref, reactive } from 'vue';
import { FormInput, NameValueData, StringAnyType } from 'src/types';

export const helper = (ci) => {
  // const inputValue = reactive(ci.props.data);
  // console.log('inputValue: ', inputValue);
  // watch(inputValue, (value) => {
  //   // console.log(':Helper ');
  //   // ci.emit('updateInputData', {
  //   //   key: ci.props.settings.key,
  //   //   value: value.value,
  //   // });
  // });
  // return inputValue;
};

export const filterNonEditable = (
  data: StringAnyType,
  form: FormInput[]
): StringAnyType => {
  const newData: StringAnyType = {};
  for (const f of form) {
    if (!f.editable) {
      continue;
    }
    newData[f.key] = data[f.key];
  }
  return newData;
};

export const getFieldName = (form: FormInput[]): { [key: string]: string } => {
  const output: { [key: string]: string } = {};
  for (const v of form) {
    output[v.key] = v.name;
  }
  return output;
};

export const convertToNameValueData = (
  data: StringAnyType,
  prefixKey: string | null
): { [key: string]: NameValueData } => {
  const output: { [key: string]: NameValueData } = {};
  for (const [k, v] of Object.entries(data)) {
    console.log('Debugging - prefixKey + k: ', prefixKey + k);
    output[prefixKey + k] = { name: '', value: v };
  }
  return output;
};

export const formBuilderRules: {
  [key: string]: (val: string) => string | boolean;
} = {
  Required: (value: string) => !!value || 'Required.',
  Email: (value: string) => {
    const pattern = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/;
    return pattern.test(value) || 'Invalid e-mail.';
  },
  Number: (value: string) => {
    const pattern = /^\d+$/;
    return pattern.test(value) || 'Invalid number.';
  },
  Password: (value: string) => {
    const pattern = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/;
    return pattern.test(value) || 'Invalid password.';
  },
  LetterAndSpaceOnly: (value: string) => {
    const pattern = /^[a-zA-Z\s]*$/;
    return pattern.test(value) || 'Only letters are allowed.';
  },
  NoSpace: (value: string) => {
    const pattern = /^[a-zA-Z_]*$/;
    return pattern.test(value) || 'Only letters are allowed.';
  },
};

export const setupRules = (
  rules: string[] | undefined
): ((val: string) => string | boolean)[] => {
  if (!rules) {
    return [];
  }
  const r = [];
  for (const rule of rules) {
    r.push(formBuilderRules[rule]);
  }
  return r;
};

export const labelValueOption: string[] = [
  'Radio',
  'Button Toggle',
  'Option Group',
  'Checkbox',
  'Select',
];
// If input support only value option
export const valueOption: string[] = [];

export const rawPayload = (data: StringAnyType): { [key: string]: any } => {
  const output = {};
  for (const [k, v] of Object.entries(data)) {
    output[k] = v.value;
  }
  return output;
};
