<template>
  <div class="column">
    <div class="col">
      <q-select
        v-if="settings.show"
        :clearable="settings.clearable"
        :error="error.error"
        :errorMessage="error.message"
        :label="settings.name"
        v-model="data.value"
        :rules="rules"
        hide-bottom-space
        :options="settings.options.labelValue"
        :multiple="settings.multiple"
        filled
        emit-value
        map-options
      />
    </div>
    <div
      v-if="settings.description"
      class="col text-caption setting-description"
    >
      {{ settings.description }}
    </div>
  </div>
</template>
<script setup lang="ts">
import { reactive, defineProps, ref, onBeforeMount } from 'vue';
import { setupRules } from 'src/service/utils/form-builder';
import { api } from 'src/boot/axios';

const props = defineProps(['data', 'settings', 'error']);
const data = reactive(props.data);
const error = reactive(props.error);
const settings = reactive(props.settings);
let rules: ((val: string) => string | boolean)[] = [];
onBeforeMount(async () => {
  rules = setupRules(settings.rules);
  const res = await api().get('/admin/trees');
  settings.options.labelValue = [];
  for (const tree of res.data.data) {
    settings.options.labelValue.push({
      label: tree.attributes.label,
      value: tree.attributes.slug,
    });
  }
});
if (settings.callback !== undefined) {
  settings.callback(settings, data);
}
</script>
