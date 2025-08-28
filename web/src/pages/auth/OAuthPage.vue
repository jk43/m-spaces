<template>
  <div></div>
</template>

<script setup lang="ts">
import { ref, reactive, onBeforeMount } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { api } from 'src/boot/axios';
import { queryBuilder } from 'src/service/utils/http';

const route = useRoute();
const router = useRouter();
const type = route.params.type as string;
const provider = route.params.provider as string;

console.log(': ', route.query);

onBeforeMount(async () => {
  const res = await api().get(`/auth/oauth/${type}/${provider}`, {
    params: route.query,
  });
  window.location.replace(res.data.data.redirect);
});
</script>

<style scoped></style>
