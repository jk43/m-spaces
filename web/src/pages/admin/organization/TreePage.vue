<template>
  <q-page class="flex flex-left q-pa-none">
    <span class="text-h4 back"
      ><q-btn icon="arrow_back" @click="goBack" flat dense size="10px"
    /></span>
    <TreeComponent
      :node="treeData"
      :deleteHandler="confirmDeleteNode"
      :onReorder="onReorder"
      @showForm="showForm"
    ></TreeComponent>
  </q-page>
  <q-dialog v-model="showFormDialog">
    <q-card style="width: 400px">
      <q-card-section class="q-pb-xs">
        <div v-if="formMode == 'edit'" class="text-h6">
          Edit {{ targetNodeLabel }}
        </div>
        <div v-else class="text-h6">
          Add new node to the {{ targetNodeLabel }}
        </div>
      </q-card-section>
      <q-card-section class="q-pt-none">
        <FormBuilder
          :data="selectedData"
          :settings="formInputElem"
          :handler="submitHandler"
          @onSubmit="formbuilderSubmitted"
          @onReset="formbuilderResetted"
        />
      </q-card-section>
    </q-card>
  </q-dialog>
  <q-dialog v-model="showDeleteDialog" persistent>
    <q-card>
      <q-card-section class="row items-center">
        <q-avatar icon="delete" color="primary" text-color="white" />
        <span class="q-ml-sm text-h6"
          >Do you want to delete "{{ nodeNameToDelete }}"?</span
        >
      </q-card-section>
      <q-card-section>
        <p>
          If you delete this "{{ nodeNameToDelete }}", all nodes connected to it
          will also be deleted. Once deleted, the data cannot be recovered.
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
          @click="deleteNode"
          v-close-popup
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
  <br />
</template>

<script setup>
import { ref, reactive, onBeforeMount } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { api } from 'src/boot/axios';
import TreeComponent from 'src/components/TreeComponent.vue';
import FormBuilder from 'src/components/form/FormBuilderComponent.vue';
import { rawPayload } from 'src/service/utils/form-builder';
import { useOrgStore } from 'src/stores/org-store';

const storeOrg = useOrgStore();
const slug = useRoute().params.slug;
const router = useRouter();
const dataReady = ref(false);

onBeforeMount(async () => {
  loadData();
});

function goBack() {
  router.back();
}

const treeData = ref([]);

const loadData = async () => {
  dataReady.value = false;
  try {
    const res = await api().get('/admin/tree?slug=' + slug);
    treeData.value[0] = res.data.data;
    //treeData.value[0].children = [];
  } catch (err) {
    treeData.value = [];
  }
  dataReady.value = true;
};

const showFormDialog = ref(false);
// node from TreeComponent
const node = ref({});
let formMode = 'add';
let targetNodeLabel = '';

const showForm = (n, index, mode) => {
  node.value = n;
  if (mode === 'edit') {
    formMode = 'edit';
    for (const key in n.attributes) {
      if (selectedData.value[key]) {
        console.log('Debugging - key: ', key);
        selectedData.value[key].value = n.attributes[key];
      }
    }
    targetNodeLabel = selectedData.value.label.value;
  } else {
    formMode = 'add';
    cleanSelectedData();
    targetNodeLabel = n.attributes.label;
    selectedData.value.slug.value = n.attributes.slug;
  }
  showFormDialog.value = true;
};

const selectedData = ref({
  id: { name: '', value: null },
  label: { name: '', value: '' },
  description: { name: '', value: '' },
  slug: { name: '', value: '' },
});

const cleanSelectedData = () => {
  for (const key in selectedData.value) {
    selectedData.value[key].value = '';
  }
};

const submitHandler = async () => {
  const payload = {
    data: { ...rawPayload(selectedData.value) },
  };
  //add new
  if (!selectedData.value.id.value) {
    // create key
    const res = (await api().post('/admin/tree', payload)).data.data;
    selectedData.value.slug.value = res.slug;
    if (node.value.children == undefined) {
      node.value.children = [];
    }
    node.value.children.push({
      // fake id
      parent_id: 1,
      id: res.id,
      attributes: { ...rawPayload(selectedData.value) },
      children: [],
    });
  } else {
    const res = (await api().put('/admin/tree', payload)).data.data;
    for (const key in payload.data) {
      node.value.attributes[key] = payload.data[key];
    }
  }
  showFormDialog.value = false;
  cleanSelectedData();
};

const formbuilderSubmitted = async () => {
  return;
};

const onReorder = (node) => {
  console.log('onReorder', node);
  const payload = {
    data: node,
  };
  api().put('/admin/tree/reorder', payload);
};

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
]);

const showDeleteDialog = ref(false);
const nodeNameToDelete = ref('');
const nodeIndexToDelete = ref(0);

const confirmDeleteNode = (n, index) => {
  node.value = n;
  nodeNameToDelete.value = node.value[index].attributes.label;
  showDeleteDialog.value = true;
  nodeIndexToDelete.value = index;
  return;
};

const deleteNode = async () => {
  const payload = {
    data: {
      data: node.value[nodeIndexToDelete.value],
    },
  };
  const res = await api().delete('/admin/tree', payload);
  const isRoot = node.value[nodeIndexToDelete.value].parent_id === null;
  node.value.splice(nodeIndexToDelete.value, 1);
  nodeNameToDelete.value = '';
  showDeleteDialog.value = false;
  nodeIndexToDelete.value = 0;
  if (isRoot) {
    goBack();
  }
};

const formbuilderResetted = () => {
  return;
};
</script>

<style scoped>
.back {
  cursor: pointer;
  margin: 0px;
  margin-top: 8px;
  margin-left: 15px;
}
</style>
