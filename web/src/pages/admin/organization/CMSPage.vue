<template>
  <q-page class="flex flex-center">
    <q-table
      :rows="rows"
      title="Boards"
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
      <template #body-cell-active="props">
        <q-td>
          <q-icon name="check" v-if="props.row.active === 'Y'" color="green" />
          <q-icon name="close" v-else color="red" />
        </q-td>
      </template>
      <template #body-cell-edit="props">
        <q-td @click="showEditForm(props.row)" class="cursor-pointer">
          <q-icon name="edit" />
        </q-td>
      </template>
      <template #body-cell-info="props">
        <q-td @click="goToBoard(props.row)" class="cursor-pointer">
          <q-icon name="pageview" />
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

  <q-dialog v-model="showForm">
    <q-card style="width: 400px">
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
</template>

<script setup lang="ts">
import { ref, reactive, onBeforeMount, watch } from 'vue';
import FormBuilder from 'src/components/form/FormBuilderComponent.vue';
import { api } from 'src/boot/axios';

import { rawPayload } from 'src/service/utils/form-builder';
import { Payload, StringAnyType } from 'src/service/types/payload';
import { Response } from 'src/types';
import { queryBuilder } from 'src/service/utils/http';
import { useRoute, useRouter } from 'vue-router';

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

interface FormField {
  key: string;
  name: string;
  description?: string;
  discription?: string;
  type: string;
  editable: boolean;
  rules: string[];
  options?: {
    labelValue: Array<{ label: string; value: string }>;
  };
  order: number;
}

interface PayloadField {
  name: string;
  value: string | number | null;
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

const router = useRouter();
const route = useRoute();

const getUrl = '/admin/boards';
const postUrl = '/admin/board';
const putUrl = '/admin/board';
const deleteUrl = '/admin/board';
const formUrl = '/admin/board/form';

const columns: TableColumn[] = [
  {
    name: 'name',
    label: 'Name',
    field: 'name',
    align: 'left',
    headerClasses: 'bg-grey-1 text-black',
    sortable: true,
  },
  {
    name: 'slug',
    label: 'Slug',
    field: 'slug',
    align: 'left',
    headerClasses: 'bg-grey-1 text-black',
  },
  {
    name: 'write',
    label: 'Write',
    field: (row: Board) => row.settings.find((s) => s.key === 'write')?.v,
    align: 'center',
    headerClasses: 'bg-grey-1 text-black',
  },
  {
    name: 'read',
    label: 'Read',
    field: (row: Board) => row.settings.find((s) => s.key === 'read')?.v,
    align: 'center',
    headerClasses: 'bg-grey-1 text-black',
  },
  {
    name: 'active',
    label: 'Active',
    field: 'active',
    align: 'center',
    headerClasses: 'bg-grey-1 text-black',
    sortable: true,
  },
  {
    name: 'created_at',
    label: 'Created At',
    field: (row: Board) => row.CreatedAt.split('T')[0],
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

onBeforeMount(async () => {
  for (const [k, v] of Object.entries(route.query)) {
    queryParams.value[k] = v as string | number | boolean;
  }
  loadForm();
  loadData();
});

const dataReady = ref<boolean>(false);

// Query parameters to send to the API
const queryParams = ref<{ [key: string]: string | number | boolean }>({
  sortBy: 'created_at',
  descending: false,
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
const rows = ref<Board[]>([]);

// Fetch boards from API
const loadData = async (): Promise<void> => {
  const q = queryBuilder(router, queryParams);
  const res: Response<{ [key: string]: any }> = await api().get(getUrl + q);
  rows.value = res.data.data.rows;
  queryParams.value.rowsNumber = res.data.data.total;
  dataReady.value = true;
};

const payload = ref<{ [key: string]: PayloadField }>({
  id: { name: '', value: '' },
  name: { name: '', value: '' },
  slug: { name: '', value: '' },
  active: { name: '', value: '' },
});

const showForm = ref<boolean>(false);
const form = reactive<FormField[]>([
  {
    key: 'name',
    name: 'Name',
    description: '',
    type: 'Input',
    editable: true,
    rules: ['Required'],
  },
  {
    key: 'active',
    name: 'Active',
    discription: '',
    rules: ['Required'],
    type: 'Select',
    editable: true,
    options: {
      labelValue: [
        { label: 'Active', value: 'Y' },
        { label: 'Inactive', value: 'N' },
      ],
    },
  },
]);

const loadForm = async (): Promise<void> => {
  const res: Response<{ [key: string]: any }> = await api().get(formUrl);
  const s = Object.entries(res.data.data);
  s.sort((a, b) => a[1].order - b[1].order);
  for (const [k, v] of s) {
    form.push(v as FormField);
  }
};

const showCreateForm = (): void => {
  for (const key of Object.keys(payload.value)) {
    payload.value[key].value = null;
  }
  showForm.value = true;
};

const showEditForm = (row: Board): void => {
  console.log('Debugging - row: ', row);
  payload.value = {
    id: { name: 'id', value: row.id },
    name: { name: 'name', value: row.name },
    slug: { name: 'slug', value: row.slug },
    active: { name: 'active', value: row.active },
  };
  for (const v of row.settings) {
    payload.value[v.key] = { name: v.key, value: v.v };
  }
  showForm.value = true;
};

const submitHandler = async (): Promise<void> => {
  const p: Payload<StringAnyType> = {
    data: { ...rawPayload(payload.value) },
  };
  if (payload.value.id.value) {
    await api().put(putUrl, p);
  } else {
    await api().post(postUrl, p);
  }
  showForm.value = false;
  loadData();
};

const rowToDelete = ref<RowToDelete>({
  confirm: false,
  row: {
    id: 0,
    name: '',
    slug: '',
    active: '',
    CreatedAt: '',
    settings: [],
  },
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
        id: rowToDelete.value.row.id,
      },
    },
  };
  await api().delete(deleteUrl, deletePayload);
  rowToDelete.value = { confirm: false, row: {} as Board };
  loadData();
};

watch(
  () => payload.value.name.value,
  (newValue: string | number | null) => {
    if (newValue && typeof newValue === 'string') {
      const slug = newValue
        .trim()
        .replace(/[^a-zA-Z0-9\s]/g, '') // Remove special characters
        .replace(/\s+/g, '-') // Convert spaces to hyphens
        .toLowerCase(); // Convert to lowercase

      payload.value.slug.value = slug;
    }
  }
);

const goToBoard = (row: Board): void => {
  router.push(`/admin/board/${row.slug}`);
};

const showPreview = (row: Board): void => {
  console.log('Debugging - row: ', row);
};
</script>

<style scoped></style>
