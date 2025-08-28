<template>
  <q-page class="flex flex-center">
    <div v-if="formData.length === 0">
      <q-card class="q-pa-md">
        <q-card-section>
          <div class="text-h6">No Form Input Found</div>
        </q-card-section>
        <q-card-section>
          <div class="text-subtitle2">Please add form input to get started</div>
        </q-card-section>
        <q-card-section>
          <div class="text-subtitle2">
            <q-btn
              @click="editFormInput(null)"
              color="primary"
              label="Add input"
            />
          </div>
        </q-card-section>
      </q-card>
    </div>
    <q-markup-table v-else>
      <thead>
        <tr>
          <th class="text-left"></th>
          <th class="text-left">Name</th>
          <th class="text-left">Type</th>
          <th class="text-center">Shareable</th>
          <th class="text-center">View</th>
          <th class="text-center">Edit</th>
          <th class="text-center">Default Value</th>
          <th class="text-center">
            <q-btn
              @click="editFormInput(null)"
              color="primary"
              icon="add"
              dense
              size="12px"
            />
          </th>
        </tr>
      </thead>
      <draggable
        v-model="formData"
        tag="tbody"
        item-key="id"
        @end="onEnd"
        @start="onStart"
      >
        <template #item="{ element }">
          <tr @click="element.key;">
            <td><q-icon name="drag_handle" /></td>
            <td @click="editFormInput(element.key)">{{ element.name }}</td>
            <td>{{ element.type }}</td>
            <td class="text-center">
              <q-icon v-if="element.shareable" name="done" />
            </td>
            <td class="text-center">
              {{ element.view }}
            </td>
            <td class="text-center">
              {{ element.edit }}
            </td>
            <td>{{ element.defaultValue }}</td>
            <td class="text-center">
              <q-btn
                @click="deleteInput(element)"
                icon="clear"
                dense
                size="8px"
              />
            </td>
          </tr>
        </template>
      </draggable>
      <q-inner-loading
        :showing="!dataReady"
        label="Please wait..."
        label-class="text-teal"
        label-style="font-size: 1.1em"
      />
    </q-markup-table>
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
          {{ selectedData.targetKey.value ? 'Edit' : 'Add' }} Input
        </div>
      </q-card-section>
      <q-card-section>
        <FormBuilder
          v-if="dataReadyForForm"
          :data="selectedData"
          :settings="formInputElem"
          :handler="submitHandler"
          @onSubmit="formbuilderSubmitted"
          @onReset="formbuilderResetted"
        />
      </q-card-section>
    </q-card>
  </q-dialog>
  {{ selectedData }}
</template>
<script setup lang="ts">
import { onBeforeMount, reactive, ref, watch, toRaw, camelize } from 'vue';
import draggable from 'vuedraggable';

import { api } from 'src/boot/axios';
import { useOrgStore } from 'src/stores/org-store';
import { FormInput, Payload, StringAnyType, DeletePayload } from 'src/types';
import FormBuilder from 'src/components/form/FormBuilderComponent.vue';
import {
  labelValueOption,
  valueOption,
  formBuilderRules,
  rawPayload,
} from 'src/service/utils/form-builder';
import _ from 'lodash';
import { is } from 'quasar';

const storeOrg = useOrgStore();

const dataReady = ref(false);
const formData = ref([]);
const isEdit = ref(false);

onBeforeMount(async () => {
  loadData();
});

const loadData = async () => {
  dataReady.value = false;
  const res = await api().get('/admin/organization/form/userMetadata');
  try {
    formData.value = res.data.data.userMetadata.sort(
      (a, b) => a.order - b.order
    );
  } catch (err) {
    formData.value = [];
  }
  dataReady.value = true;
};

const onEnd = async (evt) => {
  if (evt.newIndex === evt.oldIndex) {
    return;
  }
  const newOrder = {};
  for (const [index, element] of formData.value.entries()) {
    newOrder[element.key] = index + 1;
  }
  const payload: Payload<StringAnyType> = {
    data: newOrder,
  };
  const res = await api().put('/admin/organization/form-order', payload);
  loadData();
};

const onStart = (evt) => {
  console.log('Debugging - onStart: ', evt);
};

//to hide or show the options
const hideOption = (settings, data) => {
  if (labelValueOption.includes(selectedData.value.type.value)) {
    settings.show = { showLabelValue: true, showOption: false };
  } else {
    settings.show = { showLabelValue: false, showOption: false };
  }
  watch(selectedData.value, (newVal) => {
    if (labelValueOption.includes(selectedData.value.type.value)) {
      settings.show = { showLabelValue: true, showOption: false };
    } else {
      settings.show = { showLabelValue: false, showOption: false };
    }
  });
};

const hideTree = (settings, data) => {
  console.log('Debugging - settings: ', selectedData.value.type.value);
  if (selectedData.value.type.value === 'Tree') {
    settings.show = true;
  } else {
    settings.show = false;
  }
  watch(selectedData.value, (newVal) => {
    if (selectedData.value.type.value === 'Tree') {
      settings.show = true;
    } else {
      settings.show = false;
    }
  });
  console.log('Debugging - settings.show : ', settings.show);
};

const makeKey = (settings, data) => {
  if (isEdit.value) {
    return;
  }
  watch(data, (newVal) => {
    selectedData.value.key.value = _.camelCase(newVal.value);
  });
};

const formInputElem = reactive([
  {
    key: 'name',
    name: 'Name',
    description: '',
    options: [],
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'Input',
    editable: true,
    rules: ['Required', 'LetterAndSpaceOnly'],
    callback: makeKey,
  },
  {
    key: 'key',
    name: 'Key',
    description: 'This is the Key used when saving to the DB.',
    options: [],
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'Input',
    editable: false,
    rules: ['Required', 'NoSpace'],
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
    key: 'type',
    name: 'Type',
    description: '',
    options: {
      labelValue: [
        { label: 'Input', value: 'Input' },
        { label: 'Select', value: 'Select' },
        { label: 'Radio', value: 'Radio' },
        { label: 'Checkbox', value: 'Checkbox' },
        { label: 'Toggle', value: 'Toggle' },
        { label: 'Option Group', value: 'Option Group' },
        { label: 'Password', value: 'Password' },
        { label: 'Hierarchy', value: 'Tree' },
      ],
    },
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'Select',
    editable: true,
    rules: ['Required'],
  },
  {
    key: 'options',
    name: 'Options',
    description: '',
    options: [],
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'Options',
    callback: hideOption,
  },
  {
    key: 'slug',
    name: 'Hierarchy',
    description: '',
    options: [],
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'TreeSelect',
    callback: hideTree,
  },
  {
    key: 'rules',
    name: 'Rules',
    description: '',
    options: [],
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'Rules',
    editable: true,
  },
  {
    key: 'edit',
    name: 'Can be edit by',
    description: 'Role with the ability to modify data',
    options: {
      labelValue: storeOrg.getRoles(),
    },
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'Select',
    editable: true,
    defaultValue: 'admin',
  },
  {
    key: 'view',
    name: 'View',
    description: 'Role with the ability to modify data',
    options: {
      labelValue: storeOrg.getRoles(),
    },
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'Select',
    editable: true,
    defaultValue: 'admin',
  },
  {
    key: 'defaultValue',
    name: 'Default Value',
    description: '',
    options: [],
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'Input',
    editable: true,
  },
  {
    key: 'onTable',
    name: 'Display on Table',
    description: '',
    options: [],
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'Toggle',
    editable: true,
    defaultValue: false,
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
// flag to show the form dialog
const edit = ref(false);
const dataReadyForForm = ref(false);
const selectedData = ref({
  targetKey: { name: '', value: '' },
  type: { name: '', value: '' },
  name: { name: '', value: '' },
  description: { name: '', value: '' },
  key: { name: '', value: '' },
  order: { name: '', value: '' },
  options: { name: '', value: { labelValue: [], options: [] } },
  slug: { name: '', value: '' },
  rules: { name: '', value: [] },
  shareable: { name: '', value: false },
  defaultValue: { name: '', value: '' },
  edit: { name: '', value: '' },
  view: { name: '', value: '' },
});

const editFormInput = (key) => {
  // add a new input
  isEdit.value = true;
  if (key === null) {
    isEdit.value = false;
    selectedData.value.targetKey = { name: '', value: '' };
    selectedData.value.type = { name: '', value: '' };
    selectedData.value.name = { name: '', value: '' };
    selectedData.value.key = { name: '', value: '' };
    selectedData.value.description = { name: '', value: '' };
    selectedData.value.order = { name: '', value: 100 };
    selectedData.value.options = {
      name: '',
      value: { labelValue: [], options: [] },
    };
    selectedData.value.slug = { name: '', value: '' };
    selectedData.value.rules = { name: '', value: [] };
    selectedData.value.shareable = { name: '', value: false };
    selectedData.value.edit = {
      name: '',
      value: storeOrg.getRoles()[storeOrg.getRoles().length - 1].value,
    };
    selectedData.value.view = {
      name: '',
      value: storeOrg.getRoles()[storeOrg.getRoles().length - 1].value,
    };
    selectedData.value.defaultValue = { name: '', value: '' };
  }
  edit.value = true;
  dataReadyForForm.value = true;
  for (const v of formData.value) {
    if (v.key !== key) {
      continue;
    }
    // convert to utils.FormBuilderFields(golang) format
    for (const [key, value] of Object.entries(v)) {
      if (key === 'key') {
        selectedData.value['targetKey'] = { name: v.name, value: value };
      }
      selectedData.value[key] = { name: v.name, value: value };
    }
  }
};

const submitHandler = async () => {
  const payload: Payload<StringAnyType> = {
    data: { ...rawPayload(selectedData.value) },
  };
  //add new
  if (!selectedData.value.targetKey.value) {
    // create key
    await api().post('/admin/organization/form', payload);
  } else {
    await api().put('/admin/organization/form', payload);
  }
  //edit.value = false;
  loadData();
};

const formbuilderSubmitted = (data: any) => {
  return;
};

const formbuilderResetted = (data: any) => {
  return;
};
</script>
