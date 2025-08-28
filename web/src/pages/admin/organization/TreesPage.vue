<template>
  <q-page class="flex flex-center">
    <div v-if="dataReady && trees.length === 0">
      <q-card class="q-pa-md">
        <q-card-section>
          <div class="text-h6">No hierarchy Found</div>
        </q-card-section>
        <q-card-section>
          <div class="text-subtitle2">Please add hierarchy to get started</div>
        </q-card-section>
        <q-card-section>
          <div class="text-subtitle2">
            <q-btn @click="showForm" color="primary" label="Add hierarchy" />
          </div>
          <div class="q-mt-sm">
            <q-btn @click="bulk = true" color="secondary">Bulk upload</q-btn>
          </div>
        </q-card-section>
      </q-card>
    </div>
    <q-table
      v-else
      :rows="trees"
      title="Hierarchies"
      :columns="columns"
      row-key="id"
      :loading="!dataReady"
    >
      <template v-slot:top-right>
        <q-btn @click="bulk = true" dense color="secondary" size="12px"
          >Bulk upload</q-btn
        >
        <q-btn
          @click="showForm"
          dense
          size="12px"
          color="primary"
          icon="add"
          class="q-ml-sm"
      /></template>
      <template #body-cell-label="props">
        <q-td class="cursor-pointer">
          <router-link
            :to="`/admin/organization/tree/${props.row.attributes.slug}`"
          >
            {{ props.row.attributes.label }}</router-link
          >
        </q-td>
      </template>
    </q-table>
  </q-page>
  <q-dialog v-model="showDeleteDialog" persistent>
    <q-card>
      <q-card-section class="row items-center">
        <q-avatar icon="delete" color="primary" text-color="white" />
        <span class="q-ml-sm"
          >Do you want to delete {{ elementToDelete.name }}?</span
        >
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat label="Cancel" color="primary" v-close-popup />
        <q-btn
          flat
          label="Delete"
          color="red"
          @click="deleteInput(null)"
          v-close-popup
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
  <q-dialog v-model="edit">
    <q-card style="width: 400px">
      <q-card-section>
        <div class="text-h6">
          {{ seletedData.id.value ? 'Edit' : 'Add' }} Input
        </div>
      </q-card-section>
      <q-card-section>
        <FormBuilder
          v-if="dataReadyForForm"
          :data="seletedData"
          :settings="formInputElem"
          :handler="submitHandler"
          @onSubmit="formbuilderSubmitted"
          @onReset="formbuilderResetted"
        />
      </q-card-section>
    </q-card>
  </q-dialog>
  <q-dialog v-model="bulk">
    <q-card style="width: 400px">
      <q-card-section>
        <div class="text-h6">Bulk Upload</div>
      </q-card-section>
      <q-card-section>
        <FormBuilder
          v-if="dataReadyForForm"
          :data="bulkData"
          :settings="bulkDataElem"
          :handler="bulkDataHandler"
          @onSubmit="formbuilderSubmitted"
          @onReset="formbuilderResetted"
        />
      </q-card-section>
    </q-card>
  </q-dialog>
</template>
<script setup lang="ts">
import { onBeforeMount, reactive, ref, watch, toRaw, camelize } from 'vue';
import draggable from 'vuedraggable';

import { api } from 'src/boot/axios';
import { useOrgStore } from 'src/stores/org-store';
import { Payload, StringAnyType, DeletePayload } from 'src/types';
import FormBuilder from 'src/components/form/FormBuilderComponent.vue';
import { rawPayload } from 'src/service/utils/form-builder';
import _ from 'lodash';

const storeOrg = useOrgStore();

const dataReady = ref(false);
const trees = ref([]);
const isEdit = ref(false);

onBeforeMount(async () => {
  loadData();
  dataReadyForForm.value = true;
});

const loadData = async () => {
  dataReady.value = false;
  const res = await api().get('/admin/trees');
  trees.value = res.data.data;
  dataReady.value = true;
};

const columns = [
  {
    name: 'id',
    label: 'id',
    field: 'id',
    align: 'Left',
    headerClasses: 'bg-grey-1 text-black',
  },
  {
    name: 'label',
    label: 'Label',
    field: (row) => row.attributes.label,
    align: 'left',
    headerClasses: 'bg-grey-1 text-black',
    sortable: true,
  },
  {
    name: 'alias',
    label: 'Alias',
    field: (row) => row.attributes.slug,
    align: 'left',
    headerClasses: 'bg-grey-1 text-black',
    sortable: true,
  },
  {
    name: 'view',
    label: 'View',
    field: (row) => _.capitalize(row.attributes.view),
    align: 'left',
    headerClasses: 'bg-grey-1 text-black',
  },
];

const formInputElem = reactive([
  {
    key: 'label',
    name: 'Label',
    description: '',
    options: [],
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'Input',
    editable: true,
    rules: ['Required'],
  },
  {
    key: 'description',
    name: 'description',
    description: '',
    options: [],
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'Input',
    editable: true,
  },
  {
    key: 'view',
    name: 'View',
    description: 'Role with the ability to view hierarchy',
    options: {
      labelValue: storeOrg.getRoles(),
    },
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'Select',
    editable: true,
    defaultValue: 'admin',
    rules: ['Required'],
  },
]);

const elementToDelete = ref({});
const showDeleteDialog = ref(false);

const deleteInput = async (elem) => {
  if (elem !== null) {
    elementToDelete.value = elem;
  }
  if (!showDeleteDialog.value) {
    showDeleteDialog.value = true;
    return;
  }
  const payload: DeletePayload<any> = {
    data: {
      data: elementToDelete.value,
    },
  };
  const res = await api().delete('/admin/organization/form', payload);
  loadData();
};

const edit = ref(false);
const bulk = ref(false);

const bulkDataElem = reactive([
  {
    key: 'json',
    name: 'JSON',
    description: '',
    options: [],
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'TextArea',
    editable: true,
    rules: ['Required'],
  },
]);

const bulkData = ref({
  json: { name: '', value: '' },
});

const bulkDataHandler = async () => {
  const data = JSON.parse(bulkData.value.json.value);
  const payload: Payload<StringAnyType> = {
    data,
  };
  await api().post('/admin/trees', payload);
  bulkData.value.json.value = '';
  bulk.value = false;
  loadData();
};

const dataReadyForForm = ref(false);
const seletedData = ref({
  title: { name: '', value: '' },
  description: { name: '', value: '' },
  view: { name: '', value: '' },
  parentId: { name: '', value: '' },
  rootId: { name: '', value: '' },
  id: { name: '', value: '' },
});

const showForm = (key) => {
  edit.value = true;
};

const submitHandler = async () => {
  const payload: Payload<StringAnyType> = {
    data: { ...rawPayload(seletedData.value) },
  };
  //add new
  if (!seletedData.value.id.value) {
    // create key
    await api().post('/admin/tree', payload);
  } else {
    await api().put('/admin/tree', payload);
  }
  //edit.value = false;
  loadData();
  edit.value = false;
};

const formbuilderSubmitted = (data: any) => {
  return;
};

const formbuilderResetted = (data: any) => {
  return;
};
</script>
