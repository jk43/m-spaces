<template>
  <q-page padding class="colum flex-center" width="400px">
    <Slot>
      <div>
        <sapn>regular slot</sapn
        ><q-btn color="primary" label="Click me"></q-btn>
      </div>
      <template #btn="prop">
        <div>Val from slot : {{ prop.slotVal }} {{ prop.text }}</div>
        <div>This is btn slot</div>
        <q-btn color="primary" label="Click me"></q-btn>
        <q-btn color="primary" label="Click me"></q-btn>
      </template>
    </Slot>
  </q-page>
</template>

<script setup>
import { ref } from 'vue';
import Slot from 'src/components/sandbox/SlotComponent.vue';

const inputText = ref('');
</script>

<style lang="sass"></style>
