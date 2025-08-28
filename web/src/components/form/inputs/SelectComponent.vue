<template>
  <div class="column">
    <div class="col">
      <q-select
        :clearable="settings.clearable"
        :error="error.error"
        :errorMessage="error.message"
        :label="settings.name"
        v-model="data.value"
        :rules="rules"
        hide-bottom-space
        :options="settings.options.labelValue"
        :readonly="!settings.editable"
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
const props = defineProps(['data', 'settings', 'error']);
const data = reactive(props.data);
const error = reactive(props.error);
const settings = reactive(props.settings);
let rules: ((val: string) => string | boolean)[] = [];
onBeforeMount(() => {
  rules = setupRules(settings.rules);
});
if (settings.callback !== undefined) {
  settings.callback(settings, data);
}
</script>
