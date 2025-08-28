<template>
  <q-layout view="hhh lpR fff">
    <q-header elevated class="bg-white text-black">
      <div class="row my-header q-gutter-lg">
        <div class="col section">
          <q-toolbar-title>
            <q-btn
              dense
              flat
              round
              icon="menu"
              @click="toggleLeftDrawer"
              v-if="!$q.screen.gt.xs"
            />
            <q-avatar v-if="$q.screen.gt.xs">
              <img
                src="https://cdn.quasar.dev/logo-v2/svg/logo-mono-black.svg"
              />
            </q-avatar>
            {{ storeOrg.name }}
          </q-toolbar-title>
        </div>
        <div class="col section" v-if="$q.screen.gt.xs">
          <!-- <q-tabs v-model="tab"> 더 알아보기 -->
          <q-tabs>
            <q-tab
              v-for="item in storeOrg.items.topMenu"
              :key="item.label"
              :label="item.label"
              :name="item.label"
            />
          </q-tabs>
        </div>
        <div class="col section" v-if="$q.screen.gt.xs">
          <div class="row left-menu q-col-gutter-lg">
            <div class="col-11 right-menu-div">
              <q-input
                square
                outlined
                input-class="search-in-layout"
                placeholder="Search"
                dense
              ></q-input>

              <q-btn
                color="primary"
                @click="toggleSearch"
                icon="search"
                class="search-button-in-right-menu"
              />
            </div>
            <div class="col-1 right-menu-div">
              <q-btn round>
                <q-avatar
                  :color="storeUser.getProfileImage ? 'white' : 'primary'"
                  text-color="white"
                  size="42px"
                >
                  <img
                    v-if="storeUser.getProfileImage"
                    :src="storeUser.getProfileImage"
                  />
                  <div v-else>{{ userInitials }}</div>
                </q-avatar>
                <q-menu :offset="[0, 10]">
                  <q-list style="min-width: 200px">
                    <template
                      v-for="(item, index) in storeOrg.items.userDropdownMenu"
                      :key="index"
                    >
                      <q-item clickable v-close-popup>
                        <q-item-section avatar>
                          <q-icon :name="item.icon" />
                        </q-item-section>
                        <q-item-section
                          ><router-link
                            :to="item.to"
                            class="router-link-exact-active"
                            >{{ item.label }}</router-link
                          ></q-item-section
                        >
                      </q-item>
                    </template>
                    <q-separator />
                    <q-item clickable v-close-popup>
                      <q-item-section avatar>
                        <q-icon name="logout" />
                      </q-item-section>
                      <q-item-section @click="logout">Logout</q-item-section>
                    </q-item>
                  </q-list>
                </q-menu>
              </q-btn>
            </div>
          </div>
        </div>
      </div>
    </q-header>
    <q-drawer show-if-above v-model="leftDrawerOpen" side="left" bordered>
      <q-scroll-area class="fit">
        <MenuListComponent
          v-if="storeOrg.itemReady"
          :ctxSource="route.path"
          :ctxParser="menuCtxParser"
          :menuData="storeOrg.items"
        />
      </q-scroll-area>
    </q-drawer>

    <q-page-container>
      <router-view v-if="storeOrg.itemReady" />
      <!-- <router-view v-if="storeOrg.itemReady" @onLoad="handleOnLoad" /> -->
    </q-page-container>

    <q-footer elevated class="bg-grey-8 text-white">
      <q-toolbar>
        <q-toolbar-title>
          <q-avatar>
            <img src="https://cdn.quasar.dev/logo-v2/svg/logo-mono-white.svg" />
          </q-avatar>
          <div>Title</div>
        </q-toolbar-title>
      </q-toolbar>
    </q-footer>
  </q-layout>
</template>

<script setup lang="ts">
import { ref, onBeforeMount, reactive, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { api } from 'src/boot/axios';
import { OrganizationItems, Response, OrganizationItem } from 'src/types';
import { useUserStore } from 'src/stores/user-store';
import { useOrgStore } from 'src/stores/org-store';
import MenuListComponent from 'src/components/MenuListComponent.vue';

const storeUser = useUserStore();
const storeOrg = useOrgStore();
const router = useRouter();
const route = useRoute();

const leftDrawerOpen = ref(false);
const searchMode = ref(false);
const toggleSearch = () => {
  searchMode.value = !searchMode.value;
};
const toggleLeftDrawer = () => {
  leftDrawerOpen.value = !leftDrawerOpen.value;
};

// let orgItems = reactive<OrganizationItems>({});
let userInitials = ref<string>('');
// let profileImage = ref<string>('');
// let profileImageKey = ref<string>(storeUser.profileImage);

onBeforeMount(async () => {
  const res = await api().get<Response<OrganizationItems>>(
    '/organization/items'
  );
  await storeOrg.setItems(res.data.data);
  userInitials.value = storeUser.firstName[0] + storeUser.lastName[0];
});

const menuCtxParser = (newCtx: string, oldCtx: string): string => {
  newCtx = newCtx.split('/')[1];
  oldCtx = oldCtx.split('/')[1];
  if (newCtx === oldCtx) {
    return '';
  }
  return newCtx + 'LeftNavMenu';
};

let menuItems = reactive<OrganizationItem[]>([]);
const handleOnLoad = (data: OrganizationItem[]) => {
  menuItems.splice(0);
  for (const i in data) {
    menuItems.push(data[i]);
  }
};

let logout = async () => {
  try {
    const result = await api().post('/user/logout', {});
  } catch (err) {
    console.log('err: ', err);
  }
  storeUser.setUserInfo(null);
  router.push('/auth/login');
  console.log('logout: ');
};
</script>

<style lang="scss">
.my-header {
  .section {
    // background-color: red;
    padding: 20px;
  }
}
.left-menu {
  .right-menu-div {
    display: flex;
    justify-content: right; /* 수평 가운데 정렬 */
    align-items: center; /* 수직 가운데 정렬 */
  }
}
.search-in-layout {
  width: 200px;
}

.search-button-in-right-menu {
  margin-left: 10px;
  margin-right: 10px;
}
</style>
