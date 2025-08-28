<template>
  <div>
    <span>This is a slot component</span>
    <slot></slot>
    <slot v-if="$slots.btn" name="btn" :slotVal="slotVal" text="yoyo"></slot>
    <div v-if="!$slots.btn">
      <span>NO SLOT</span>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, useSlots } from 'vue';

const slots = useSlots();
const slotVal = ref('Slot val from component');

onMounted(() => {
  console.log('Btn slot provided:', !!slots.btn);
});
</script>

<style lang="sass"></style>
