<template>
  <div style="width: 200px">
    <q-select
      v-model="selectedValue.value"
      :options="data"
      label="Saved Searches"
      dense
      filled
      clearable
      @update:model-value="updateHander"
      @clear="updateHander"
      map-options
    />
  </div>
</template>

<script setup lang="ts">
import { ref, defineProps, defineEmits, reactive } from 'vue';

const props = defineProps(['data', 'model']);
const emit = defineEmits(['selected']);
const data = ref(props.data);
const selectedValue = reactive(props.model);
const savedSearch = reactive([]);
//const selectedValue = ref(null);

const updateHander = (value) => {
  console.log('Debugging - value: ', value);
  emit('selected', { ...value });

  selectedValue.value = value;
};

// onBeforeMount(() => {
//   for (const i of data.value) {
//     savedSearch.push({ label: i.label, value: i.key });
//   }
// });
</script>
