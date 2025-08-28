<template>
  <q-page class="flex flex-center">
    <q-table
      :rows="rows"
      title="Questions"
      :columns="columns"
      row-key="id"
      :loading="!dataReady"
      v-model:pagination="queryParams"
      @request="onRequest"
    >
      <template v-slot:top-right>
        <div v-if="totalProgresses > 0">In queue: {{ totalProgresses }}</div>
        <q-btn
          @click="showCreateForm = true"
          dense
          size="12px"
          color="primary"
          icon="add"
          class="q-ml-sm"
      /></template>
      <template #body-cell-active="props">
        <q-td>
          <div v-if="props.row.is_active">
            <q-icon name="check" />
          </div>
          <div v-else>
            <q-icon name="close" color="red" />
          </div>
        </q-td>
      </template>
      <template #body-cell-error="props">
        <q-td>
          <div v-if="props.row.error">
            <q-icon
              @click="showError(props.row.error)"
              name="error"
              color="red"
            />
          </div>
        </q-td>
      </template>
      <template #body-cell-view="props">
        <q-td>
          <q-icon
            @click="viewQuestion(props.row.id)"
            name="pageview"
            size="xs"
          />
        </q-td>
      </template>
    </q-table>
  </q-page>

  <q-dialog v-model="showCreateForm">
    <q-card style="width: 400px">
      <q-card-section>
        <div class="text-h6">Create a content</div>
      </q-card-section>
      <q-card-section>
        <FormBuilder
          :data="questionPayload"
          :settings="questionRequestForm"
          :handler="submitHandler"
        />
      </q-card-section>
    </q-card>
  </q-dialog>

  <q-dialog v-model="showErrorDialog">
    <q-card>
      <q-card-section> Error </q-card-section>
      <q-card-section class="row items-center">
        {{ errorMessage }}
      </q-card-section>
    </q-card>
  </q-dialog>

  <q-dialog v-model="showQuesion">
    <q-card>
      <q-card-section>
        <QuestionDetailsComponent :data="question" />
      </q-card-section>
    </q-card>
  </q-dialog>
</template>
<script setup lang="ts">
import { ref, reactive, onBeforeMount } from 'vue';
import FormBuilder from 'src/components/form/FormBuilderComponent.vue';
import { api } from 'src/boot/axios';

import { rawPayload } from 'src/service/utils/form-builder';
import { Payload, StringAnyType } from 'src/service/types/payload';
import { Response } from 'src/types';
import { GetWebSocket, queryBuilder } from 'src/service/utils/http';
import { useRoute, useRouter } from 'vue-router';

import QuestionDetailsComponent from 'src/components/content/QuestionDetailsComponent.vue';

const router = useRouter();
const route = useRoute();

const showCreateForm = ref(false);

///////////////////////////////

const message = ref('');
const reload = ref(false);
const progresses = ref(null);
const totalProgresses = ref(0);
let websocket = null;

const connectWebSocket = () => {
  websocket = GetWebSocket('post-question-progress');

  websocket.onopen = () => {
    console.log('WebSocket connection opened');
  };

  websocket.onmessage = (event) => {
    console.log('WebSocket message received:', event.data);
    const message = JSON.parse(event.data);
    if (message.reload) {
      loadData();
      console.log('loading data: ');
    }
    if (message.progress) {
      progresses.value = message.progress;
      totalProgresses.value = Object.keys(message.progress).length;
    }
  };

  websocket.onerror = (error) => {
    console.error('WebSocket error:', error);
  };

  websocket.onclose = () => {
    console.log('WebSocket connection closed');
  };
};

///////////////////////////////

onBeforeMount(async () => {
  for (const [k, v] of Object.entries(route.query)) {
    queryParams.value[k] = v;
  }
  loadData();
  connectWebSocket();
});

const dataReady = ref(false);
// query to send to the API
const queryParams = ref<{ [key: string]: string | number | boolean }>({
  sortBy: 'created_at',
  descending: false,
  page: 1,
  rowsPerPage: 25,
  rowsNumber: 0,
});
// pagination
const onRequest = (props) => {
  queryParams.value = props.pagination;
  queryBuilder(router, queryParams);
  loadData();
};
// Questions from the API
let rows = ref([]);
// q-tqble column configuration
const columns = [
  {
    name: 'id',
    label: 'ID',
    field: 'id',
    align: 'Left',
    headerClasses: 'bg-grey-1 text-black',
  },
  {
    name: 'topic',
    label: 'Topic',
    field: (row) => row.slugs.topic,
    align: 'left',
    headerClasses: 'bg-grey-1 text-black',
  },
  {
    name: 'is_active',
    label: 'Active',
    field: 'is_active',
    align: 'left',
    headerClasses: 'bg-grey-1 text-black',
    sortable: true,
  },
  {
    name: 'rate',
    label: 'Rate',
    field: (row) =>
      isNaN(row.rate_total / row.rate_count)
        ? 'N/A'
        : row.rate_total / row.rate_count,
    align: 'center',
    headerClasses: 'bg-grey-1 text-black',
    sortable: true,
  },
  {
    name: 'error',
    label: 'Error',
    field: (row) => (row.error ? true : false),
    align: 'center',
    headerClasses: 'bg-grey-1 text-black',
    sortable: true,
  },
  {
    name: 'created_at',
    label: 'Created At',
    field: (row) => row.created_at.split('T')[0],
    align: 'left',
    headerClasses: 'bg-grey-1 text-black',
    sortable: true,
  },
  {
    name: 'view',
    label: '',
    field: '',
    align: 'left',
    headerClasses: 'bg-grey-1 text-black',
    sortable: false,
  },
];

// fetch question data
const loadData = async () => {
  let q = queryBuilder(router, queryParams);
  const res: Response<{ [key: string]: any }> = await api().get(
    '/admin/math/questions' + q
  );
  rows.value = res.data.data;
  console.log('Debugging - res.data.data: ', res.data.data);
  console.log('Debugging - rows.value: ', rows);
  dataReady.value = true;
  queryParams.value.rowsNumber = res.data.total;
};

// error dialog
const showErrorDialog = ref(false);
const errorMessage = ref('');
const showError = (error) => {
  errorMessage.value = error;
  showErrorDialog.value = true;
};

const questionPayload = ref({
  tree: {
    name: '',
    value: [],
  },
  llm: { name: '', value: '' },
});

const questionRequestForm = reactive([
  {
    key: 'llm',
    name: 'LLM Model',
    description: '',
    options: {
      labelValue: [
        { label: 'OpenAI', value: 'openai' },
        { label: 'Claude', value: 'claude' },
      ],
    },
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'Select',
    editable: true,
    rules: ['Required'],
  },
  {
    key: 'numberOfQuestions',
    name: 'How many questions?',
    description: '',
    options: {
      labelValue: [
        { label: 'OpenAI', value: 'openai' },
        { label: 'Claude', value: 'claude' },
      ],
    },
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'Input',
    editable: true,
    rules: ['Required'],
  },
  {
    key: 'tree',
    name: 'Tree',
    description: '',
    options: [],
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'Tree',
    editable: true,
    slug: 'mathematics',
    rules: [],
  },
]);

const showQuesion = ref(false);
let question = reactive({});

const viewQuestion = async (id) => {
  const res = await api().get('/admin/math/question/' + id);
  question = res.data.data;
  showQuesion.value = true;
  //router.push({ name: 'admin-question', params: { id: id } });
};

const submitHandler = async () => {
  const payload: Payload<StringAnyType> = {
    data: { ...rawPayload(questionPayload.value) },
  };
  await api().post('/admin/math', payload);
};
</script>
