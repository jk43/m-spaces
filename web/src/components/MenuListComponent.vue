<template>
  <q-list>
    <template v-for="(item, index) in menuItem" :key="index">
      <q-item clickable v-ripple>
        <q-item-section avatar>
          <q-icon :name="item.icon" />
        </q-item-section>
        <q-item-section>
          <router-link :to="item.to" class="router-link-exact-active">{{
            item.label
          }}</router-link>
        </q-item-section>
      </q-item>
      <!-- <q-separator :key="'sep' + index" v-if="menuItem.options.separator" /> -->
      <template v-for="(sub, subIdx) in item.subMenu" :key="subIdx">
        <q-item clickable v-ripple>
          <q-item-section avatar>
            <q-icon :name="sub.icon" />
          </q-item-section>
          <q-item-section>
            <router-link :to="sub.to" class="router-link-exact-active">{{
              sub.label
            }}</router-link>
          </q-item-section>
        </q-item>
      </template>
    </template>
  </q-list>
</template>

<script setup lang="ts">
import { defineProps, watch, reactive, onBeforeMount } from 'vue';
import { useRoute } from 'vue-router';
import { OrganizationItems, OrganizationItem } from 'src/types';

const prop = defineProps(['ctxSource', 'ctxParser', 'menuData']) as {
  ctxSource: any;
  ctxParser: (newCtx: string, oldCtx: string) => string;
  menuData: OrganizationItems;
};

const route = useRoute();
let menuItem = reactive<OrganizationItem[]>([]);

const getFirstPart = (path: string): string => {
  return path.split('/')[1];
};
let prevPath = getFirstPart(route.path);

const setMenuItem = (ctx: string) => {
  menuItem.splice(0);
  for (const i in prop.menuData[ctx]) {
    menuItem[i] = prop.menuData[ctx][i];
  }
};

watch(
  () => prop.ctxSource,
  (newCtx, oldCtx) => {
    const ctx = prop.ctxParser(newCtx, oldCtx);
    if (ctx === '') {
      return;
    }
    setMenuItem(ctx);
  }
);

onBeforeMount(async () => {
  const ctx = prop.ctxParser(prop.ctxSource, '');
  setMenuItem(ctx);
});
</script>
<style></style>
