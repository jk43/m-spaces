<template>
  <div class="error" v-if="errors.formErrors.length">
    <q-banner rounded inline-actions class="text-white bg-red">
      <div v-for="e in errors.formErrors" :key="e.message">
        {{ e.message }}
      </div>
    </q-banner>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue';
import { FormErrors } from 'src/types';

const props = defineProps<{
  errors: FormErrors;
}>();
</script>
<style></style>
