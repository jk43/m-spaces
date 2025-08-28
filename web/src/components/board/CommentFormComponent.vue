<template>
  <div class="text-subtitle2 q-mb-sm">Write a comment</div>
  <FormBuilder :data="payload" :settings="form" :handler="submitHandler" />
</template>

<script setup lang="ts">
/*
<template>
  <div class="text-subtitle2 q-mb-sm">Write a comment</div>
  <q-input
    v-model="payload.text.value"
    type="textarea"
    placeholder="Enter your comment..."
    rows="3"
    outlined
    dense
  />
  <div class="row justify-end q-mt-sm">
    <q-btn
      label="Submit Comment"
      color="primary"
      @click="submitHandler"
      :loading="submitting"
      :disable="!payload.text.value.trim()"
    />
  </div>
</template>
*/

import { api } from 'src/boot/axios';
import { Response, NameValueData, FormInput } from 'src/types';
import { Payload, StringAnyType } from 'src/service/types/payload';

import { ref, defineEmits, reactive } from 'vue';
import { formatDistanceToNow } from 'date-fns';
import { rawPayload } from 'src/service/utils/form-builder';
import { useRoute, useRouter } from 'vue-router';
import { queryBuilder } from 'src/service/utils/http';
import FormBuilder from 'src/components/form/FormBuilderComponent.vue';

//const router = useRouter();
// const route = useRoute();

const getUrl = '/board/posts/comments';
const postUrl = '/board/posts';

const props = defineProps<{
  parentPost: any;
  parentBoard: any;
}>();

console.log('Debugging - props: ', props);

const emit = defineEmits(['comment-added']);

const slug = props.parentPost.slug;
const boardSlug = props.parentBoard.slug;

const post = ref(props.parentPost);
const rows = ref<any[]>([]);
const newComment = ref('');
const submitting = ref(false);
const dataReady = ref(false);

const queryParams = ref<{ [key: string]: string | number | boolean }>({
  boardSlug: boardSlug as string,
  postSlug: slug as string,
  sortBy: 'created_at',
  descending: true,
  page: 1,
  rowsPerPage: 25,
  rowsNumber: 0,
});

// Comment submission function
const submitHandler = async () => {
  const p: Payload<StringAnyType> = {
    data: { ...rawPayload(payload.value) },
  };
  const uuid = crypto.randomUUID();
  p.data.boardSlug = boardSlug;
  p.data.parentSlug = slug;
  p.data.title = uuid;
  p.data.slug = uuid;
  await api().post(postUrl, p);
  emit('comment-added');
  payload.value.text.value = '';
};

const form = reactive<FormInput[]>([
  {
    key: 'text',
    name: 'Text',
    discription: '',
    rules: ['Required'],
    type: 'Advanced Editor',
    editable: true,
  },
]);

const payload = ref<{ [key: string]: NameValueData }>({
  //id: { name: '', value: '' },
  boardSlug: { name: '', value: '' },
  parentSlug: { name: '', value: '' }, //parent slug
  title: { name: '', value: '' },
  slug: { name: '', value: '' },
  text: { name: '', value: '' },
});
</script>

<style scoped>
.post-card {
  max-width: 800px;
  margin: 0 auto;
}

.comment-card {
  background-color: #f8f9fa;
}

.comment-item {
  margin-bottom: 1rem;
}

.comment-item:last-child {
  margin-bottom: 0;
}
</style>
