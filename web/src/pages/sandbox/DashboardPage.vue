<template>
  <q-page padding class="colum flex-center" width="400px">
    <q-input v-model="inputText" filled placeholder="Type something..." />
    {{ inputText }}
  </q-page>
</template>

<script setup>
import { ref } from 'vue';

const inputText = ref('');
</script>

<style lang="sass"></style>
