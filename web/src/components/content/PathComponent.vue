<template>
  <q-breadcrumbs>
    <q-breadcrumbs-el v-for="(v, i) of path" :key="i" :label="v.label" />
  </q-breadcrumbs>
</template>

<script setup lang="ts">
import { defineProps, onBeforeMount, reactive, ref, computed } from 'vue';

const props = defineProps(['path']);
const path = computed(() => {
  return props.path.slice(1);
});
</script>
