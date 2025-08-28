<template>
  <div class="column bg-grey-2 rounded-borders q-pa-sm">
    <div class="col">
      <q-toggle
        :error="error.error"
        :errorMessage="error.message"
        :label="settings.name"
        v-model="data.value"
        :rules="rules"
        hide-bottom-space
        left-label
        :disable="!settings.editable"
      />
    </div>
    <div
      v-if="settings.description"
      class="col text-caption setting-description"
    >
      {{ settings.description }}
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
