<template>
  <div>
    <div v-for="(form, index) in forms" :key="index">
      <component :is="form.type" :label="form.name" :options="form.options" />
    </div>

    <!-- Submit button or other form controls here -->
  </div>
</template>

<script setup>
import { QInput, QSelect, QRadio } from 'quasar';
//import QInputComponent from 'src/components/form/QInputComponent.vue';

const selectOptions = [
  { label: 'Option 1', value: 'option1' },
  { label: 'Option 2', value: 'option2' },
];
const radioOptions = [
  { label: 'True', value: true },
  { label: 'False', value: false },
];

// const forms = [
//   { type: QInputComponent, name: 'gaga', data: '' },
//   { type: QSelect, name: 'age', options: selectOptions },
//   { type: QRadio, name: 'true', options: radioOptions },
// ];

// export default {
//   data() {
//     var selectOptions: [
//       { label: 'Option 1', value: 'option1' },
//       { label: 'Option 2', value: 'option2' }
//     ];
//     var radioOptions: [
//       { label: 'True', value: true },
//       { label: 'False', value: false }
//     ];
//     return {
//       forms: [
//         { type: QInput, name: 'name' },
//         { type: QSelect, name: 'age', options: selectOptions },
//         { type: QRadio, name: 'true', options: radioOptions },
//       ],
//       formData: {},
//     };
//   },
// };
</script>
