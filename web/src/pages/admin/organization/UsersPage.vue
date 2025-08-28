<template>
  <q-page class="flex flex-center">
    <q-table
      :rows="rows"
      title="Users"
      :columns="columns"
      row-key="id"
      :visibleColumns="visibleColumns"
      :loading="!dataReady"
      v-model:pagination="queryParams"
      @request="onRequest"
    >
      <template v-slot:no-data="{ icon, message, filter }">
        <div class="full-width row flex-center text-accent q-gutter-sm">
          <span> {{ message }} </span>
          <q-icon size="2em" :name="filter ? 'filter_b_and_w' : icon" />
        </div>
      </template>
      <template v-slot:top-right>
        <SavedSerches
          v-if="searchReady"
          :data="savedSearchesToComponent"
          :model="selectedValue"
          @selected="handleSavedSearch"
          class="q-mr-sm"
        />
        <span v-if="searchMode"
          ><q-btn dense flat
            ><q-badge outline color="secondary" label="In Search" />
            <q-popup-proxy ref="searchInfo">
              <q-card>
                <q-card-section class="text-h6" v-if="searchingBySaved">
                  {{ savedSearchValue.label }}
                </q-card-section>
                <q-separator inset />

                <q-item
                  v-ripple
                  v-for="[k, v] of Object.entries(searchValues)"
                  :key="k"
                >
                  <q-item-section>
                    <q-item-label caption>{{ v.name }}</q-item-label>
                    <q-item-label
                      ><span class="text-capitalize">{{
                        v.value
                      }}</span></q-item-label
                    >
                  </q-item-section>
                </q-item>
                <q-card-section>
                  <q-btn
                    v-if="!searchingBySaved"
                    @click="saveSearch = true"
                    flat
                    >Save</q-btn
                  >
                  <q-btn v-if="searchingBySaved" @click="deleteSavedSerch" flat
                    >Delete</q-btn
                  >
                  <q-btn @click="clearSearch" flat>Clear</q-btn>
                </q-card-section>
              </q-card>
            </q-popup-proxy>
          </q-btn>
        </span>
        <q-btn @click="search = true" dense size="12px" icon="search" flat />
        <q-btn
          @click="createUser"
          dense
          size="12px"
          color="primary"
          icon="add"
          class="q-ml-sm"
        />
      </template>
      <template #body-cell-status="props">
        <q-td>
          <q-badge v-if="props.row.status === 'active'" color="primary"
            >{{ props.row.status[0].toUpperCase()
            }}<q-tooltip>{{ props.row.status }}</q-tooltip></q-badge
          >
          <q-badge v-if="props.row.status === 'waiting'" color="grey"
            >{{ props.row.status[0].toUpperCase()
            }}<q-tooltip>{{ props.row.status }}</q-tooltip></q-badge
          >
          <q-badge v-if="props.row.status === 'inactive'" color="red"
            >{{ props.row.status[0].toUpperCase()
            }}<q-tooltip>{{ props.row.status }}</q-tooltip></q-badge
          >
          <q-badge v-if="props.row.status === '-'" color="primary">{{
            props.row.status[0].toUpperCase()
          }}</q-badge>
        </q-td>
      </template>

      <template #body-cell-edit="props">
        <q-td @click="editUser(props.row)" class="cursor-pointer">
          <q-icon name="edit" />
        </q-td>
      </template>
      <template #body-cell-info="props">
        <q-td @click="showPopup(props.row)" class="cursor-pointer">
          <q-icon name="info" />
        </q-td>
      </template>
      <template #body-cell-delete="props">
        <q-td
          @click="deleteRow(props.row)"
          class="cursor-pointer"
          style="color: red"
        >
          <q-icon name="delete" />
        </q-td>
      </template>
    </q-table>
  </q-page>
  <q-dialog v-model="edit">
    <q-card>
      <q-card-section class="text-h6 q-pb-none">
        <span>Edit</span>
      </q-card-section>
      <q-card-section class="row items-center q-pt-none">
        <div class="" style="width: 400px">
          <FormBuilder
            v-if="dataReadyForForm"
            :data="currentRowConverted"
            :settings="userEditForm"
            :handler="editUserHandler"
          />
        </div>
      </q-card-section>
    </q-card>
  </q-dialog>
  <q-dialog v-model="create">
    <q-card>
      <q-card-section class="text-h6 q-pb-none">
        <span>Create User</span>
      </q-card-section>
      <q-card-section class="row items-center q-pt-none">
        <div class="" style="width: 400px">
          <FormBuilder
            v-if="dataReadyForForm"
            :data="currentRowConverted"
            :settings="userCreateForm"
            :handler="createUserHandler"
          />
        </div>
      </q-card-section>
    </q-card>
  </q-dialog>
  <q-dialog v-model="search">
    <q-card>
      <q-card-section class="text-h6 q-pb-none">
        <span>Search</span>
      </q-card-section>
      <q-card-section class="row items-center q-pt-none">
        <div class="" style="width: 400px">
          <FormBuilder
            :data="searchValues"
            :settings="searchForm"
            :handler="searchHandler"
          >
            <template #buttons>
              <div class="col">
                <q-btn label="Submit" type="submit" color="primary" />
                <q-btn
                  label="Clear"
                  type="reset"
                  color="primary"
                  flat
                  class="q-ml-sm"
                  @click="clearSearch"
                /></div></template
          ></FormBuilder>
        </div>
      </q-card-section>
    </q-card>
  </q-dialog>
  <q-dialog v-model="saveSearch">
    <q-card>
      <q-card-section> Search </q-card-section>
      <q-card-section class="row items-center">
        <div class="q-pa-md" style="width: 400px">
          <FormBuilder
            :data="saveSearchValues"
            :settings="saveSearchForm"
            :handler="saveSearchHandler"
          ></FormBuilder>
        </div>
      </q-card-section>
    </q-card>
  </q-dialog>

  <q-dialog v-model="info.show">
    <q-card style="width: 500px">
      <q-card-section>
        <div class="text-h6 text-capitalize">
          {{ info.data.firstName }}
          {{ info.data.lastName }} (<span>{{ info.data.role }}</span
          >)
          <q-badge color="blue">
            {{ info.data.status }}
          </q-badge>
        </div>
        <div class="text-subtitle2">{{ info.data.email }}</div>
      </q-card-section>
      <q-separator inset />
      <q-card-section>
        <q-item
          v-ripple
          v-for="[k, v] of Object.entries(info.data.metadata.userMetadata)"
          :key="k"
        >
          <q-item-section>
            <q-item-label caption>{{ metadataFieldName[k] }}</q-item-label>
            <q-item-label
              ><span class="text-capitalize">{{ v }}</span></q-item-label
            >
          </q-item-section>
        </q-item>
      </q-card-section>
    </q-card>
  </q-dialog>
  <q-dialog v-model="deleteUser.confirm" persistent>
    <q-card>
      <q-card-section class="row items-center">
        <q-avatar icon="delete" color="primary" text-color="white" />
        <span class="q-ml-sm text-h6"
          >Do you want to delete {{ deleteUser.row.firstName }}
          {{ deleteUser.row.lastName }}?</span
        >
      </q-card-section>
      <q-card-section>
        <p>
          If you delete a user, all of the user's activity logs will also be
          deleted. Once deleted, the data cannot be recovered.
        </p>
      </q-card-section>
      <q-card-actions align="right">
        <q-btn
          @click="deleteUser = {}"
          label="Cancel"
          color="primary"
          v-close-popup
        />
        <q-btn
          flat
          label="Delete"
          color="red"
          @click="deleteRow"
          v-close-popup
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import { ref, onBeforeMount, reactive } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { api } from 'src/boot/axios';

import FormBuilder from 'src/components/form/FormBuilderComponent.vue';
import SavedSerches from 'src/components/SavedSearchComponent.vue';
import { useOrgStore } from 'src/stores/org-store';
import { Payload, Response, StringAnyType, FormInput } from 'src/types';

import { getFieldName, rawPayload } from 'src/service/utils/form-builder';
import { queryBuilder, separatePagination } from 'src/service/utils/http';
import { date } from 'quasar';
import _ from 'lodash';

const router = useRouter();
const route = useRoute();
const storeOrg = useOrgStore();

const dataReady = ref<boolean>(false);

const getAddr = '/admin/users';
const putAddr = '/admin/user';
const postAddr = '/admin/user';

const rows = reactive([]);

const dataReadyForForm = ref<boolean>(false);
//for detail page
let metadataFieldName: { [key: string]: string } = {};
let userMetadataInputs = [];

const queryParams = ref<{ [key: string]: string | number | boolean }>({
  sortBy: 'createdAt',
  descending: true,
  page: 1,
  rowsPerPage: 10,
  rowsNumber: 0,
});

onBeforeMount(async () => {
  const formRes: Response<FormInput> = (
    await api().get('/admin/organization/form/userMetadata')
  ).data;
  metadataFieldName = getFieldName(formRes.data.userMetadata);
  for (const i of formRes.data.userMetadata) {
    i.key = `metadata.userMetadata.${i.key}`;
    userEditForm.push(i);
    userCreateForm.push(i);
    const forSearch = { ...i };
    forSearch.rules = [];
    forSearch.defaultValue = '';
    forSearch.clearable = true;
    forSearch.multiple = true;
    searchForm.push(forSearch);
    userMetadataInputs.push(i);
  }
  let i = 0;
  for (const [k, v] of Object.entries(route.query)) {
    queryParams.value[k] = v;
    i++;
    if (i > 5) {
      searchMode.value = true;
    }
    // build search value from query
    for (const [ac, av] of Object.entries(userEditForm)) {
      if (av.key === k) {
        searchValues[k] = { name: av.name, value: v };
      }
    }
  }
  const page = parseInt(route.query.page);
  const rowsPerPage = parseInt(route.query.rowsPerPage);
  queryParams.value.sortBy = route.query.sortBy
    ? route.query.sortBy
    : 'createdAt';
  queryParams.value.descending = route.query.descending
    ? Boolean(route.query.descending)
    : true;
  queryParams.value.page = page ? page : 1;
  queryParams.value.rowsPerPage = rowsPerPage ? rowsPerPage : 10;

  loadSavedSearches();
  loadData();
});

const loadData = async () => {
  rows.splice(0, rows.length);
  const q = queryBuilder(router, queryParams);
  dataReady.value = false;
  const res: Response<{ [key: string]: any }> = (await api().get(getAddr + q))
    .data;
  if (res.data.data === null) {
    //noData.value = true;
    dataReady.value = true;
    return;
  }
  for (const v of res.data.data) {
    rows.push(v);
  }
  queryParams.value.rowsNumber = res.data.total;
  dataReady.value = true;
};

const loadSavedSearches = async () => {
  // get saved searches
  //savedSearchesToComponent.value = [];
  savedSearchesToComponent.value.splice(
    0,
    savedSearchesToComponent.value.length
  );
  const searchs: Response<any> = (
    await api().get('/user/store/adminUserSearch')
  ).data;
  // hide saved search input if no saved search
  if (searchs.data === null || searchs.data.length === 0) {
    searchReady.value = false;
    return;
  }
  savedSearches.value = searchs.data;
  searchReady.value = true;
  for (const i of searchs.data) {
    savedSearchesToComponent.value.push({ label: i.label, value: i.key });
  }
};

// delete

const deleteUser = ref({
  row: {},
  confirm: false,
});

const deleteRow = async (row) => {
  if (!deleteUser.value.confirm) {
    deleteUser.value.row = row;
    deleteUser.value.confirm = true;
    return;
  }
  const payload: Payload<{ [key: string]: string }> = {
    data: {
      data: {
        id: deleteUser.value.row.id,
      },
    },
  };
  await api().delete(putAddr, payload);
  deleteUser.value = {};
  loadData();
};

// saved search

// ref contains saved searches
const savedSearches = ref([]);
const savedSearchesToComponent = ref([]);
// values to send backend
const saveSearchValues = reactive({});
const saveSearch = ref(false);
const search = ref(false);
const searchMode = ref(false);
const searchReady = ref(false);
const searchingBySaved = ref(false);
const savedSearchValue = ref({});
// v-model for component
const selectedValue = reactive({ value: null });
const searchInfo = ref(null);

const saveSearchHandler = async () => {
  const pData = rawPayload(saveSearchValues);
  route.query.page = '1';
  pData['value'] = route.query;
  pData['ctx'] = 'adminUserSearch';
  pData['key'] = _.camelCase(pData.name);
  const payload: Payload<StringAnyType> = {
    data: pData,
  };
  await api().post('/user/store', payload);
  saveSearch.value = false;
  //reset saveSearchValues
  saveSearchValues.name = {};
  loadSavedSearches();
  selectedValue.value = { label: pData.name, value: pData.key };
  savedSearchValue.value.label = pData.name;
  savedSearchValue.value.value = pData.key;
  searchingBySaved.value = true;
  searchInfo.value.hide();
};

//const noData = ref(false);

const handleSavedSearch = (value) => {
  savedSearchValue.value = value;
  const key = value.value;
  for (const i of savedSearches.value) {
    if (i.key !== key) {
      continue;
    }
    const parts = separatePagination(i.value);
    queryParams.value = i.value;
    for (const [k, v] of Object.entries(parts.others)) {
      for (const s of searchForm) {
        if (s.key === k) {
          searchValues[k] = { name: s.name, value: v };
        }
      }
    }
  }
  //turn on "in search" mode
  searchMode.value = true;
  searchingBySaved.value = true;
  loadData();
};

const deleteSavedSerch = async () => {
  const payload: Payload<{ [key: string]: string }> = {
    data: {
      data: {
        key: savedSearchValue.value.value,
        ctx: 'adminUserSearch',
      },
    },
  };
  await api().delete(`/user/store`, payload);
  selectedValue.value = null;
  loadSavedSearches();
  clearSearch();
};

const clearSearch = () => {
  router.push({ path: route.path, query: {} });
  queryParams.value = {
    sortBy: queryParams.value.sortBy,
    descending: queryParams.value.descending,
    page: 1,
    rowsPerPage: 10,
    rowsNumber: 0,
  };
  for (const [key, value] of Object.entries(searchValues)) {
    searchValues[key] = {
      name: value.name,
    };
  }
  loadData();
  search.value = false;
  searchMode.value = false;
};

// pagination
const onRequest = (props) => {
  queryParams.value = props.pagination;
  queryBuilder(router, queryParams);
  loadData();
};

// details

const info = ref({
  show: false,
  data: {},
});

const showPopup = (row) => {
  info.value.data = row;
  info.value.show = true;
};

const formatPayloadForNewUser = (data) => {
  const raw = { ...rawPayload(data) };
  const output = {
    metadata: {
      userMetadata: {},
    },
  };
  for (const [key, value] of Object.entries(raw)) {
    if (key.includes('metadata.userMetadata.')) {
      const newKey = key.replace('metadata.userMetadata.', '');
      output.metadata.userMetadata[newKey] = value;
    } else {
      output[key] = value;
    }
  }
  return output;
};

const searchValues = reactive<{ [key: string]: any }>({});

const searchHandler = async () => {
  searchingBySaved.value = false;
  for (const [key, value] of Object.entries(rawPayload(searchValues))) {
    if (value === '' || value === null || value === undefined) {
      delete queryParams.value[key];
      continue;
    }
    searchMode.value = true;
    queryParams.value[key] = value;
  }
  queryParams.value.page = 1;
  const q = queryBuilder(router, queryParams);
  loadData();
  search.value = false;
};

// edit

const edit = ref(false);
const currentRowConverted = ref({});

const editUserHandler = async () => {
  const payload: Payload<StringAnyType> = {
    data: currentRowConverted.value,
  };
  await api().put(putAddr, payload);
  edit.value = false;
  loadData();
};

const editUser = (row: any) => {
  edit.value = true;
  const valuesToAdd = [
    'id',
    'email',
    'firstName',
    'lastName',
    'role',
    'status',
  ];
  for (const v of valuesToAdd) {
    currentRowConverted.value[v] = {
      name: '',
      value: row[v],
    };
  }
  for (const v of userMetadataInputs) {
    currentRowConverted.value[v.key] = {
      name: v.name,
      value:
        row.metadata.userMetadata[v.key.replace('metadata.userMetadata.', '')],
    };
  }
  dataReadyForForm.value = true;
};

// create

const create = ref(false);

const createUserHandler = async () => {
  const payload: Payload<StringAnyType> = {
    data: formatPayloadForNewUser(currentRowConverted.value),
  };
  await api().post(postAddr, payload);
  create.value = false;
  loadData();
};

const createUser = () => {
  create.value = true;
  currentRowConverted.value = {};
  for (const i of userCreateForm) {
    let value = '';
    if (i.key === 'role') {
      value = useOrgStore().getRoles()[0].value;
    }
    currentRowConverted.value[i.key] = {
      name: i.name,
      value: value,
    };
  }
  dataReadyForForm.value = true;
};

const saveSearchForm = [
  {
    key: 'name',
    name: 'Name',
    discription: '',
    options: [],
    rules: [],
    type: 'input',
    editable: true,
    clearable: true,
  },
];

const searchForm = [
  {
    key: 'firstName',
    name: 'First Name',
    discription: '',
    options: [],
    rules: [],
    type: 'input',
    editable: true,
    clearable: true,
  },
  {
    key: 'lastName',
    name: 'Last Name',
    discription: '',
    options: [],
    rules: [],
    type: 'input',
    editable: true,
    clearable: true,
  },
  {
    key: 'email',
    name: 'Email',
    discription: '',
    options: [],
    rules: [],
    type: 'input',
    editable: true,
    clearable: true,
  },
  {
    key: 'role',
    name: 'Role',
    discription: '',
    rules: [],
    type: 'Select',
    editable: true,
    options: {
      labelValue: storeOrg.getRoles(),
    },
    defaultValue: storeOrg.getRoles()[0].label,
    clearable: true,
  },
  {
    key: 'status',
    name: 'Status',
    discription: '',
    rules: [],
    type: 'Select',
    editable: true,
    clearable: true,
    options: {
      labelValue: [
        { label: 'Active', value: 'active' },
        { label: 'Inactive', value: 'inactive' },
      ],
    },
  },
];

const userEditForm = [
  {
    key: 'id',
    name: 'ID',
    discription: '',
    options: [],
    rules: [],
    type: 'Input',
    editable: false,
    show: false,
  },
  {
    key: 'firstName',
    name: 'First Name',
    discription: '',
    options: [],
    rules: ['Required'],
    type: 'input',
    editable: true,
  },
  {
    key: 'lastName',
    name: 'Last Name',
    discription: '',
    options: [],
    rules: ['Required'],
    type: 'input',
    editable: true,
  },
  {
    key: 'email',
    name: 'Email',
    discription: '',
    options: [],
    rules: ['Required', 'Email'],
    type: 'input',
    editable: true,
  },
  {
    key: 'role',
    name: 'Role',
    discription: '',
    rules: ['Required'],
    type: 'Select',
    editable: true,
    options: {
      labelValue: storeOrg.getRoles(),
    },
    defaultValue: storeOrg.getRoles()[0].label,
  },
  {
    key: 'status',
    name: 'Status',
    discription: '',
    rules: ['Required'],
    type: 'Select',
    editable: true,
    options: {
      labelValue: [
        { label: 'Active', value: 'active' },
        { label: 'Inactive', value: 'inactive' },
      ],
    },
  },
  {
    key: 'password',
    name: 'Password',
    discription: '',
    options: [],
    rules: '',
    type: 'Password',
    editable: true,
  },
  {
    key: 'confirmPassword',
    name: 'Confirm Password',
    discription: '',
    options: [],
    rules: '',
    type: 'Password',
    editable: true,
  },
];

const userCreateForm = [
  {
    key: 'email',
    name: 'Email',
    discription: '',
    options: [],
    rules: ['Required', 'Email'],
    type: 'input',
    editable: true,
  },
  {
    key: 'firstName',
    name: 'First Name',
    discription: '',
    options: [],
    rules: [],
    type: 'input',
    editable: true,
  },
  {
    key: 'lastName',
    name: 'Last Name',
    discription: '',
    options: [],
    rules: [],
    type: 'input',
    editable: true,
  },
  {
    key: 'role',
    name: 'Role',
    discription: '',
    rules: ['Required'],
    type: 'Select',
    editable: true,
    options: {
      labelValue: storeOrg.getRoles(),
    },
    defaultValue: storeOrg.getRoles()[0].label,
  },
];
// for table
const visibleColumns = reactive([
  'email',
  'name',
  'role',
  'status',
  'createdAt',
  'edit',
  'info',
  'delete',
]);

const columns = [
  {
    name: 'id',
    label: 'id',
    field: 'id',
    align: 'Left',
    headerClasses: 'bg-grey-1 text-black',
  },
  {
    name: 'email',
    label: 'Email',
    field: 'email',
    align: 'left',
    headerClasses: 'bg-grey-1 text-black',
    sortable: true,
  },
  {
    name: 'name',
    label: 'Name',
    field: (row) =>
      (row.firstName === undefined ? '' : row.firstName) +
      ' ' +
      (row.lastName === undefined ? '' : row.lastName),
    align: 'left',
    headerClasses: 'bg-grey-1 text-black',
  },
  {
    name: 'role',
    label: 'Role',
    field: 'role',
    align: 'left',
    headerClasses: 'bg-grey-1 text-black',
    sortable: true,
  },
  {
    name: 'status',
    label: 'Status',
    field: 'status',
    align: 'center',
    headerClasses: 'bg-grey-1 text-black',
    sortable: true,
  },
  {
    name: 'createdAt',
    label: 'Created At',
    field: (row) => date.formatDate(row.createdAt, 'YYYY-MM-DD'),
    align: 'left',
    headerClasses: 'bg-grey-1 text-black',
    sortable: true,
  },
  {
    name: 'info',
    label: '',
    field: 'info',
    align: 'center',
    headerClasses: 'bg-grey-1 text-black',
    sortable: false,
  },
  {
    name: 'edit',
    label: '',
    field: 'edit',
    align: 'center',
    headerClasses: 'bg-grey-1 text-black',
    sortable: false,
  },
  {
    name: 'delete',
    label: '',
    field: 'delete',
    align: 'center',
    headerClasses: 'bg-grey-1 text-black',
    sortable: false,
  },
];
</script>
