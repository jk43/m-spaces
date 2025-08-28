<template>
  <div class="column">
    <div class="col">
      <q-input
        filled
        :error="error.error"
        :errorMessage="error.message"
        :label="settings.name"
        v-model="data.value"
        :rules="rules"
        hide-bottom-space
        :hide="!settings.editable"
        :readonly="!settings.editable"
        :class="{ 'no-hint': !settings.description }"
        :type="isPwd ? 'password' : 'text'"
      >
        <template v-slot:append>
          <q-icon
            :name="isPwd ? 'visibility_off' : 'visibility'"
            class="cursor-pointer"
            @click="isPwd = !isPwd"
          />
        </template>
      </q-input>
      <q-input
        filled
        :error="error.error"
        :errorMessage="error.message"
        label="Confirm Password"
        v-model="data.value"
        :rules="rules"
        hide-bottom-space
        :hide="!settings.editable"
        :readonly="!settings.editable"
        :class="{ 'no-hint': !settings.description }"
        :type="isPwd ? 'password' : 'text'"
        class="q-mt-sm"
      >
        <template v-slot:append>
          <q-icon
            :name="isPwd ? 'visibility_off' : 'visibility'"
            class="cursor-pointer"
            @click="isPwd = !isPwd"
          />
        </template>
      </q-input>
    </div>
    <div
      v-if="settings.description"
      class="col text-caption text-grey-8 setting-description"
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
const isPwd = ref(true);
let rules: ((val: string) => string | boolean)[] = [];
onBeforeMount(() => {
  rules = setupRules(settings.rules);
});
if (settings.callback !== undefined) {
  settings.callback(settings, data);
}
</script>
