<template>
  <q-page class="flex flex-center board-page">
    <q-table
      :rows="rows"
      :title="board.name"
      :columns="columns"
      row-key="id"
      :loading="!dataReady"
      v-model:pagination="queryParams"
      @request="onRequest"
    >
      <template v-slot:top-right>
        <q-btn
          @click="showCreateForm"
          dense
          size="12px"
          color="primary"
          icon="add"
          class="q-ml-sm"
      /></template>
      <template #body-cell-read="props">
        <q-td>
          <q-td>{{
            props.row.settings.find((s) => s.key === 'read')?.v[0].toUpperCase()
          }}</q-td>
        </q-td>
      </template>
      <template #body-cell-write="props">
        <q-td>
          <q-td>{{
            props.row.settings
              .find((s) => s.key === 'write')
              ?.v[0].toUpperCase()
          }}</q-td>
        </q-td>
      </template>
      <template #body-cell-title="props">
        <q-td>
          <router-link :to="`/c/p/${props.row.slug}`">{{
            props.row.title
          }}</router-link>
        </q-td>
      </template>
      <template #body-cell-active="props">
        <q-td>
          <q-icon name="check" v-if="props.row.active === 'Y'" color="green" />
          <q-icon name="close" v-else color="red" />
        </q-td>
      </template>
      <!-- <template #body-cell-edit="props">
        <q-td @click="showEditForm(props.row)" class="cursor-pointer">
          <q-icon name="edit" />
        </q-td>
      </template>
      <template #body-cell-info="props">
        <q-td @click="showReplyForm(props.row)" class="cursor-pointer">
          <q-icon name="add" />
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
      </template> -->
    </q-table>
  </q-page>

  <q-dialog v-model="showForm" maximized>
    <q-card style="width: 800px">
      <q-card-section>
        <div class="text-h6">Create a board</div>
      </q-card-section>
      <q-card-section>
        <FormBuilder
          :data="payload"
          :settings="form"
          :handler="submitHandler"
        />
      </q-card-section>
    </q-card>
  </q-dialog>

  <q-dialog v-model="rowToDelete.confirm" persistent>
    <q-card>
      <q-card-section class="row items-center">
        <q-avatar icon="delete" color="primary" text-color="white" />
        <span class="q-ml-sm text-h6"
          >Do you want to delete "{{ rowToDelete.row.name }}" board?</span
        >
      </q-card-section>
      <q-card-section>
        <p>
          If you delete the board, all of the board activity data/logs will also
          be deleted. Once deleted, the data cannot be recovered.
        </p>
      </q-card-section>
      <q-card-actions align="right">
        <q-btn
          @click="rowToDelete = { confirm: false, row: {} }"
          label="Cancel"
          color="primary"
          v-close-popup
        />
        <q-btn
          flat
          label="Delete"
          color="red"
          @click="() => deleteRow()"
          v-close-popup
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
  <q-dialog v-model="showPreviewDialog" maximized>
    <q-card class="full-height">
      <q-card-section class="q-pa-none">
        <div class="row items-center q-pa-md">
          <div class="text-h6">preview</div>
          <q-space />
          <q-btn icon="close" flat round v-close-popup />
        </div>
        <q-separator />
        <div
          class="full-height"
          style="overflow-y: auto; max-height: calc(100vh - 80px)"
        >
          <PostWithMultiComments :post="post" :board="board" />
        </div>
      </q-card-section>
    </q-card>
  </q-dialog>
  <VerificationComponent
    :data="verificationPrompt"
    @close="closedVerification"
  />
</template>

<script setup lang="ts">
import { ref, reactive, onBeforeMount, watch } from 'vue';
import FormBuilder from 'src/components/form/FormBuilderComponent.vue';
import { api } from 'src/boot/axios';
import { formatDistanceToNow } from 'date-fns';
import PostWithComments from 'src/components/board/PostWithCommentsComponent.vue';
import PostWithMultiComments from 'src/components/board/PostWithMultiCommentsComponent.vue';
import VerificationComponent from 'src/components/form/VerificationComponent.vue';

import { rawPayload } from 'src/service/utils/form-builder';
import { Payload, StringAnyType } from 'src/service/types/payload';
import { Response } from 'src/types';
import { queryBuilder } from 'src/service/utils/http';
import { useRoute, useRouter } from 'vue-router';
import { FormInput, NameValueData, VerificationInstruction } from 'src/types';
import { useUserStore } from 'src/stores/user-store';

const userStore = useUserStore();

// Type definitions
interface Board {
  id: number;
  name: string;
  slug: string;
  active: string;
  CreatedAt: string;
  settings: BoardSetting[];
}

interface BoardSetting {
  key: string;
  v: string;
}

interface TableColumn {
  name: string;
  label: string;
  field: string | ((row: Board) => string);
  align: 'left' | 'center' | 'right';
  headerClasses: string;
  sortable?: boolean;
}

interface PaginationProps {
  pagination: {
    sortBy: string;
    descending: boolean;
    page: number;
    rowsPerPage: number;
  };
  filter?: any;
  getCellValue?: (col: any, row: any) => any;
}

interface RowToDelete {
  confirm: boolean;
  row: Board;
}

let verificationPrompt: Response<VerificationInstruction> = reactive({
  result: '',
  data: {
    verificationName: '',
    URL: '',
    keyName: '',
    method: '',
    message: '',
    payload: {} as Payload<StringAnyType>,
    resend: {
      payload: {} as Payload<StringAnyType>,
      method: '',
      URL: '',
    },
  },
});

const router = useRouter();
const route = useRoute();

const slug = route.params.slug;

const getUrl = '/board/posts';
const postUrl = '/board/posts';
const putUrl = '/board/posts';
const deleteUrl = '/board/posts';
const formUrl = '/board/form';

const columns: TableColumn[] = [
  {
    name: 'title',
    label: 'Title',
    field: 'title',
    align: 'left',
    headerClasses: 'bg-grey-1 text-black',
    sortable: true,
  },
  {
    name: 'writer',
    label: 'Writer',
    field: (row: Board) =>
      row.user_name.first_name + ' ' + row.user_name.last_name,
    align: 'center',
    headerClasses: 'bg-grey-1 text-black',
  },
  {
    name: 'created_at',
    label: 'Created At',
    field: (row: Board) => {
      return formatDistanceToNow(new Date(row.created_at), {
        addSuffix: true,
        includeSeconds: true,
      });
    },
    align: 'left',
    headerClasses: 'bg-grey-1 text-black',
    sortable: true,
  },
  // {
  //   name: 'info',
  //   label: '',
  //   field: 'info',
  //   align: 'center',
  //   headerClasses: 'bg-grey-1 text-black',
  //   sortable: false,
  // },
  // {
  //   name: 'edit',
  //   label: '',
  //   field: 'edit',
  //   align: 'center',
  //   headerClasses: 'bg-grey-1 text-black',
  //   sortable: false,
  // },
  // {
  //   name: 'delete',
  //   label: '',
  //   field: 'delete',
  //   align: 'center',
  //   headerClasses: 'bg-grey-1 text-black',
  //   sortable: false,
  // },
];

onBeforeMount(async () => {
  for (const [k, v] of Object.entries(route.query)) {
    queryParams.value[k] = v as string | number | boolean;
  }
  //loadForm();
  loadBoard();
  loadData();
});

const dataReady = ref<boolean>(false);

// Query parameters to send to the API
const queryParams = ref<{ [key: string]: string | number | boolean }>({
  slug: slug as string,
  sortBy: 'created_at',
  descending: true,
  page: 1,
  rowsPerPage: 25,
  rowsNumber: 0,
});

// Pagination handler
const onRequest = (props: PaginationProps) => {
  queryParams.value = props.pagination;
  queryBuilder(router, queryParams);
  loadData();
};

// Content from the API
const rows = ref([]);
const board = ref<Board>({
  id: 0,
  name: '',
  slug: '',
  active: '',
  CreatedAt: '',
  settings: [],
});

// Fetch board from API
const loadBoard = async (): Promise<void> => {
  const res: Response<{ [key: string]: any }> = await api().get(
    `/board?slug=${slug}`
  );
  board.value = res.data.data;
};

// Fetch boards from API
const loadData = async (): Promise<void> => {
  const q = queryBuilder(router, queryParams);
  const res: Response<{ [key: string]: any }> = await api().get(getUrl + q);
  console.log('Debugging - res.data.data.rows: ', res.data.data.rows);
  rows.value = res.data.data.rows;
  queryParams.value.rowsNumber = res.data.data.total;
  dataReady.value = true;
};

const payload = ref<{ [key: string]: NameValueData }>({
  //id: { name: '', value: '' },
  boardSlug: { name: '', value: '' },
  parentSlug: { name: '', value: '' }, //parent slug
  title: { name: '', value: '' },
  slug: { name: '', value: '' },
  text: { name: '', value: '' },
});

const showForm = ref<boolean>(false);
const form = reactive<FormInput[]>([
  {
    key: 'title',
    name: 'Title',
    description: '',
    type: 'Input',
    editable: true,
    rules: ['Required'],
  },
  {
    key: 'text',
    name: 'Text',
    discription: '',
    rules: ['Required'],
    type: 'Advanced Editor',
    editable: true,
  },
]);

if (!userStore.email) {
  form.push({
    key: 'name',
    name: 'Your Name',
    description: '',
    type: 'Input',
    editable: true,
    rules: ['Required'],
  });
  form.push({
    key: 'email',
    name: 'Your Email',
    description: '',
    type: 'Input',
    editable: true,
    rules: ['Required', 'Email'],
  });
}

const loadForm = async (): Promise<void> => {
  const res: Response<{ [key: string]: any }> = await api().get(formUrl);
  const s = Object.entries(res.data.data);
  s.sort((a, b) => a[1].order - b[1].order);
  for (const [k, v] of s) {
    form.push(v as FormInput);
  }
};

const showCreateForm = (): void => {
  isEdit.value = false;
  for (const key of Object.keys(payload.value)) {
    payload.value[key].value = null;
  }
  payload.value.parentSlug.value = '';
  showForm.value = true;
};

const isEdit = ref<boolean>(false);

const showEditForm = (row: Board): void => {
  console.log('Debugging - row: ', row);
  payload.value.title.value = row.title;
  payload.value.text.value = row.text;
  payload.value.slug.value = row.slug;
  showForm.value = true;
  isEdit.value = true;
};

const showReplyForm = (row: Board): void => {
  payload.value.parentSlug.value = row.slug;
  payload.value.boardSlug.value = slug as string;
  payload.value.title.value = '';
  payload.value.text.value = '';
  showForm.value = true;
  isEdit.value = false;
};

const verificationCodeDialog = ref<boolean>(false);

const submitHandler = async (): Promise<void> => {
  const p: Payload<StringAnyType> = {
    data: { ...rawPayload(payload.value) },
  };
  p.data.boardSlug = board.value.slug;
  if (isEdit.value) {
    await api().put(putUrl, p);
  } else {
    const res: Response<{ [key: string]: any }> = await api().post(postUrl, p);
    if (res.data.result === 'verification_required') {
      verificationPrompt.result = res.data.result;
      verificationPrompt.data = res.data.data;
      verificationPrompt.data.verificationName = 'Board';
      verificationPrompt.data.payload = p;
      verificationPrompt.data.resend = {
        payload: p,
        method: 'post',
        URL: postUrl,
      };
      return;
    }
  }
  showForm.value = false;
  isEdit.value = false;
  loadData();
};

const closedVerification = (): void => {
  showForm.value = false;
  isEdit.value = false;
  loadData();
};

const rowToDelete = ref<RowToDelete>({
  confirm: false,
  row: {},
});

const deleteRow = async (row?: Board): Promise<void> => {
  if (!rowToDelete.value.confirm) {
    if (row) {
      rowToDelete.value.row = row;
      rowToDelete.value.confirm = true;
    }
    return;
  }
  const deletePayload: Payload<{ [key: string]: string }> = {
    data: {
      data: {
        slug: rowToDelete.value.row.slug,
        boardSlug: slug,
      },
    },
  };
  console.log('Debugging - deletePayload: ', deletePayload);
  await api().delete(deleteUrl, deletePayload);
  rowToDelete.value = { confirm: false, row: {} as Board };
  loadData();
};

watch(
  () => payload.value.title.value,
  (newValue: string | number | null) => {
    if (isEdit.value) {
      return;
    }
    if (newValue && typeof newValue === 'string') {
      const slug = newValue
        .trim()
        .replace(/[/&<>:"\\|?*]/g, '') // URL에서 사용할 수 없는 특수문자 제거
        .replace(/\s+/g, '-') // 연속된 공백을 하이픈 하나로 변환
        .toLowerCase(); // 소문자로 변환

      payload.value.slug.value = slug;
    }
  }
);

const showPreviewDialog = ref<boolean>(false);
const post = ref({});

const showPreview = (row: Board): void => {
  post.value = row;
  showPreviewDialog.value = true;
};
</script>

<style scoped>
.board-page {
  width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}
</style>
