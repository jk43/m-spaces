<template>
  <q-layout class="my-layout">
    <q-page-container>
      <main>
        <router-view />
      </main>
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
// import { ref, onBeforeMount, defineAsyncComponent } from 'vue';
// import { useRoute } from 'vue-router';

// const route = useRoute();
// let myView = ref(null);

// // onBeforeMount(() => {
// //   myView.value = defineAsyncComponent(
// //     () => import(`src/pages/sandbox/TimyMCEPage.vue`)
// //   );
// // });
</script>

<style lang="sass">
.my-layout
  background: rgba(#999,.15)
  border: 10px solid #000000
</style>
