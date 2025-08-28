<template>
  <div v-html="renderedString"></div>
</template>

<script setup lang="ts">
import katex from 'katex';
import 'katex/dist/katex.min.css';

import { ref, defineProps, onMounted, watch } from 'vue';

const props = defineProps(['text']);
const text = ref(props.text);

const renderedString = ref('');

const renderString = (text) => {
  const regex = /\$(.*?)\$/g;
  let match;
  let lastIndex = 0;
  let result = '';

  while ((match = regex.exec(text)) !== null) {
    const [fullMatch, formula] = match;
    const startIndex = match.index;
    const endIndex = regex.lastIndex;

    // Append text before the formula
    result += text.slice(lastIndex, startIndex);

    // Render the formula using KaTeX
    const renderedFormula = katex.renderToString(formula, {
      throwOnError: false,
    });

    // Append the rendered formula
    result += renderedFormula;

    lastIndex = endIndex;
  }

  // Append any remaining text after the last formula
  result += text.slice(lastIndex);

  return result;
};

onMounted(() => {
  renderedString.value = renderString(text.value);
});

watch(
  () => props.text,
  (newValue, oldValue) => {
    renderedString.value = renderString(newValue);
  }
);
</script>
