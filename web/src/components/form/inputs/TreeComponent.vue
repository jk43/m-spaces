<template>
  <div class="column">
    <div class="col">
      <q-select
        class="q-mt-md"
        v-for="(node, index) in nodes"
        :key="index"
        :clearable="settings.clearable"
        :error="error.error"
        :errorMessage="error.message"
        :label="labels[index]"
        v-model="data.value[index]"
        :rules="rules"
        hide-bottom-space
        :options="nodes[index]"
        :readonly="!settings.editable"
        :multiple="settings.multiple"
        filled
        emit-value
        map-options
        @update:model-value="update(index, $event)"
      />
    </div>
    <div
      v-if="settings.description"
      class="col text-caption setting-description"
    ></div>
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
const nodes = ref([]);
const tree = ref({});
const labels = ref([]);

let rules: ((val: string) => string | boolean)[] = [];
onBeforeMount(async () => {
  rules = setupRules(settings.rules);
  const res = await api().get('/tree/' + settings.slug);
  tree.value = res.data.data;
  labels.value.push(tree.value.attributes.label);
  if (data.value === undefined || data.value === null || data.value === '') {
    data.value = [];
  }
  if (data.value.length === 0) {
    getSiblings(tree.value.children);
  } else {
    populateNodes();
  }
});

const populateNodes = () => {
  let current = tree.value.children;
  for (const [i, d] of data.value.entries()) {
    if (i === 0) {
      getSiblings(tree.value.children);
    }
    let found = false;
    for (const [j, n] of current.entries()) {
      if (n.attributes.slug === d) {
        labels.value.push(n.attributes.label);
        getSiblings(n.children);
        current = n.children;
        found = true;
        break;
      }
    }
    if (!found) {
      data.value[i] = null;
      spliceValues(i);
    }
  }
};

const getSiblings = (elem) => {
  let siblings = [];
  if (elem === undefined || elem === null) {
    return;
  }
  for (const [i, s] of elem.entries()) {
    siblings.push({
      label: s.attributes.label,
      value: s.attributes.slug,
      s: s,
    });
  }
  nodes.value.push(siblings);
};

const spliceValues = (i) => {
  nodes.value.splice(i + 1);
  labels.value.splice(i + 1);
  data.value.splice(i + 1);
};

const update = (i, e) => {
  spliceValues(i);
  for (const n of nodes.value[i]) {
    if (n.value === e) {
      labels.value.push(n.label);
      getSiblings(n.s.children);
      break;
    }
  }
};

if (settings.callback !== undefined) {
  settings.callback(settings, data);
}
</script>
