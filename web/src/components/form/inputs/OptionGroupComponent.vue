<template>
  <div class="column q-pl-xs bg-grey-2 rounded-borders q-pa-sm">
    <div class="col text-body1 setting-name">{{ settings.name }}</div>
    <div
      v-if="settings.description"
      class="col text-caption setting-description"
    >
      {{ settings.description }}
    </div>
    <div class="col q-mt-xs">
      <q-option-group
        :error="error.error"
        :errorMessage="error.message"
        v-model="data.value"
        :rules="rules"
        :options="settings.options.labelValue"
        color="primary"
        inline
        :disable="!settings.editable"
        dense
        size="xs"
      />
    </div>
    <div class="col text-red-10 text-caption" v-if="error.error">
      <q-icon name="error" size="18px" class="q-mr-xs" />{{ error.message }}
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
