<template>
  <div class="q-pa-md" style="max-width: 400px">
    <!-- <comp :params="user" @emitted="handleEmitted" /> -->
    <Account
      v-if="route.params.section === 'account'"
      :user="user"
      @emitted="handleReset"
    />
    <Password
      v-if="route.params.section === 'password'"
      :params="user"
      @emitted="handleEmitted"
    />
    <Settings
      v-if="route.params.section === 'settings'"
      :user="user"
      :formData="orgSettings.settings"
      @emitted="handleEmitted"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onBeforeMount, defineAsyncComponent } from 'vue';
import { Component } from 'vue';
import { useRoute } from 'vue-router';
import { api } from 'src/boot/axios';
import Account from 'src/components/sandbox/AccountComponent.vue';
import Password from 'src/components/sandbox/PasswordComponent.vue';
import Settings from 'src/components/sandbox/SettingsComponent.vue';

import { useOrganiationSettingsStore } from 'src/stores/org-store';

const orgSettings = useOrganiationSettingsStore();
console.log(
  'orgSettingsorgSettingsorgSettings: ',
  orgSettings.settings.forms.userSettings[0].key
);
const comp = ref<Component | null>(Account);
const route = useRoute();

let user: any = ref({});

const handleReset = (data: any) => {
  console.log('useruseruser: ', user.value);
};

const handleEmitted = (data: any) => {
  alert(data.name);
};

const setComponent = (section: string) => {
  if (section === 'account') {
    comp.value = Account;
  } else if (section === 'password') {
    comp.value = Password;
  } else if (section === 'settings') {
    comp.value = Settings;
  }
};

onBeforeMount(async () => {
  user.value = (await api().get('/user/info')).data;
  setComponent(route.params.section as string);
});

watch(
  () => route.params.section as string,
  async (section) => {
    setComponent(section);
  }
);

// // dynamicComponent.value = Account;
// const d = {
//   name: 'alex',
//   email: 'jk@jktech.net',
// };
//dc.value = Account;

function senddata() {
  user.value = { email: 'gogo@jktech.net' };
}
</script>

<style lang="sass"></style>
