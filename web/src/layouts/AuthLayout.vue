<template>
  <q-layout view="hHh lpR fFf" class="bg-image">
    <q-page-container>
      <router-view />
    </q-page-container>

    <q-footer elevated class="bg-black text-white">
      <q-toolbar>
        <q-toolbar-title class="fit row justify-between">
          <div>Moly Space</div>
          <div>
            <span class="text-caption"
              >Copyright 2022 Hotdev.com | Background image source:
              {{ getRandomImageUrl }}
            </span>
          </div>
        </q-toolbar-title>
      </q-toolbar>
    </q-footer>
  </q-layout>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const getRandomImageUrl = computed<string>((): string => {
  const providers: string[] = [
    'https://picsum.photos',
    'https://random.responsiveimages.io',
  ];
  const i = Math.floor(Math.random() * providers.length);
  return providers[i];
});

const randomImageUrl = `url(${getRandomImageUrl.value}/1920/1080)`;
</script>

<style scoped>
.bg-image {
  background-image: v-bind(randomImageUrl);
  -webkit-background-size: cover;
  -moz-background-size: cover;
  -o-background-size: cover;
  background-size: cover;
}
</style>
