<template>
  <q-select
    filled
    v-model="data.value"
    multiple
    :options="options"
    label="Ruels"
    style="width: 250px"
  />
</template>
<script setup>
import { reactive, defineProps, ref, onBeforeMount } from 'vue';
import { formBuilderRules } from 'src/service/utils/form-builder';
import { setupRules } from 'src/service/utils/form-builder';
const props = defineProps(['data', 'settings', 'error']);
const data = reactive(props.data);
const error = reactive(props.error);
const settings = reactive(props.settings);
const multiple = ref([]);
const options = [];
let rules = [];
onBeforeMount(() => {
  for (const [i, r] of Object.entries(formBuilderRules)) {
    options.push(i);
  }
  rules = setupRules(props.settings.rules);
});
if (settings.callback !== undefined) {
  settings.callback(settings, data);
}
</script>
