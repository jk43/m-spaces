<template>
  <div class="comment-tree">
    <div v-for="comment in comments" :key="comment.id" class="comment-item">
      <q-card flat bordered class="comment-card">
        <q-card-section class="q-pa-sm">
          <div class="row items-center q-mb-xs">
            <q-avatar size="24px" class="q-mr-sm">
              <q-icon name="person" size="16px" />
            </q-avatar>
            <div class="text-caption text-weight-medium">
              {{ comment.user_name?.first_name }}
              {{ comment.user_name?.last_name }}
            </div>
            <q-space />
            <div class="text-caption text-grey-6">
              {{ formatDate(comment.created_at) }}
            </div>
          </div>
          <div
            class="text-body2 tiptap-content"
            style="white-space: pre-wrap"
            v-html="comment.text"
            @click="handleContentClick"
          ></div>
          <FilesComponent :slug="comment.slug" />

          <!-- Edit/Reply button -->
          <div class="row justify-end q-mt-xs">
            <q-btn
              flat
              dense
              size="sm"
              label="Edit"
              icon="edit"
              color="primary"
              @click="editComment(comment)"
            />
            <q-btn
              flat
              dense
              size="sm"
              label="Reply"
              icon="reply"
              color="primary"
              @click="showReplyForm(comment)"
            />
          </div>
        </q-card-section>
      </q-card>

      <!-- Reply form -->
      <div
        v-if="replyingTo.slug === comment.slug"
        class="reply-form q-ml-lg q-mt-sm"
      >
        <q-card flat bordered class="q-mt-md">
          <q-card-section>
            <CommentForm
              :parent-board="parentBoard"
              :parent-post="selectedPost"
              @comment-added="loadChildComments(selectedPost.slug)"
            />
          </q-card-section>
        </q-card>
      </div>

      <!-- Child comments -->
      <div
        v-if="commentChildren[comment.slug]?.length > 0"
        class="child-comments q-ml-lg q-mt-sm"
      >
        <CommentTree
          :replying-to="replyingTo"
          :comments="commentChildren[comment.slug]"
          :parent-board="parentBoard"
          :parent-post="selectedPost"
          @comment-added="loadChildComments(selectedPost.slug)"
        />
      </div>
    </div>

    <!-- Image Lightbox Dialog -->
    <q-dialog v-model="imageLightboxOpen" full-width full-height>
      <q-card class="lightbox-card">
        <q-card-section class="lightbox-header">
          <div class="lightbox-title">
            {{ currentImage?.alt || 'Image' }}
          </div>
          <q-btn
            flat
            round
            icon="close"
            @click="closeImageLightbox"
            class="lightbox-close"
          />
        </q-card-section>

        <q-card-section class="lightbox-content">
          <div class="media-container">
            <img
              :src="currentImage?.src"
              :alt="currentImage?.alt"
              class="media-content"
              @click="closeImageLightbox"
            />
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </div>
  <q-dialog v-model="editMode" maximized>
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
</template>

<script setup lang="ts">
import { api } from 'src/boot/axios';
import { Response, NameValueData, FormInput } from 'src/types';
import CommentForm from './CommentFormComponent.vue';
import FilesComponent from './FilesComponent.vue';
import FormBuilder from 'src/components/form/FormBuilderComponent.vue';
import { Payload, StringAnyType } from 'src/service/types/payload';
import { rawPayload } from 'src/service/utils/form-builder';

import { ref, reactive, onMounted, watch, defineAsyncComponent } from 'vue';
import { formatDistanceToNow } from 'date-fns';
import { useUserStore } from 'src/stores/user-store';

// Define self-referencing component for recursive comments
const CommentTree = defineAsyncComponent(
  () => import('./CommentTreeComponent.vue')
);

const getUrl = '/board/posts/comments';
const postUrl = '/board/posts';

const userStore = useUserStore();

const props = defineProps<{
  comments: any[];
  parentBoard: any;
  parentPost: any;
  replyingTo: { slug: string };
}>();

// to hide reply form.
const replyingTo = reactive(props.replyingTo);

const emit = defineEmits<{
  'comment-added': [];
}>();

const putUrl = '/board/posts';

const replyText = ref('');
const submitting = ref(false);
const commentChildren = ref<{ [key: number]: any[] }>({});

// Date format function
const formatDate = (dateString: string): string => {
  if (!dateString) return '';
  try {
    return formatDistanceToNow(new Date(dateString), { addSuffix: true });
  } catch {
    return new Date(dateString).toLocaleDateString('en-US');
  }
};

// Load child comments for a specific comment
const loadChildComments = async (commentId: number) => {
  try {
    const queryParams = {
      boardSlug: props.parentBoard.slug,
      postSlug: commentId.toString(),
      sortBy: 'created_at',
      descending: true,
      page: 1,
      rowsPerPage: 25,
    };

    const q = new URLSearchParams();
    Object.entries(queryParams).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        q.append(key, String(value));
      }
    });

    const res: Response<{ [key: string]: any }> = await api().get(
      getUrl + '?' + q.toString()
    );

    commentChildren.value[commentId] = res.data.data.rows || [];
  } catch (error) {
    console.error('Failed to load child comments:', error);
    commentChildren.value[commentId] = [];
  }
  replyingTo.slug = '';
};

// Load all child comments when comments change
watch(
  () => props.comments,
  async (newComments) => {
    if (newComments) {
      for (const comment of newComments) {
        await loadChildComments(comment.slug);
      }
    }
  },
  { immediate: true }
);

const selectedPost = ref({});
// Show reply form
const showReplyForm = (comment: any) => {
  if (replyingTo.slug !== '') {
    replyingTo.slug = '';
    return;
  }
  console.log('Debugging - comment: ', comment);
  selectedPost.value = comment;
  replyingTo.slug = comment.slug;
  replyText.value = '';
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

const editingComment = ref<any>({});
const editMode = ref(false);

const selectedComment = ref<any>({});

const editComment = (comment: any) => {
  selectedComment.value = comment;
  editMode.value = true;
  editingComment.value = comment;
  payload.value.text.value = comment.text;
  payload.value.slug.value = comment.slug;
  payload.value.boardSlug.value = props.parentBoard.slug;
  payload.value.parentSlug.value = comment.slug;
  payload.value.title.value = comment.title;
};

const submitHandler = async (): Promise<void> => {
  const p: Payload<StringAnyType> = {
    data: { ...rawPayload(payload.value) },
  };
  await api().put(putUrl, p);
  editMode.value = false;
  selectedComment.value.text = payload.value.text.value;
};

// // Cancel reply
// const cancelReply = () => {
//   replyingTo.value = null;
//   replyText.value = '';
// };

// // Submit reply
// const submitReply = async (comment: any) => {
//   if (!replyText.value.trim()) return;

//   submitting.value = true;
//   try {
//     const p = {
//       data: {
//         boardSlug: props.parentBoard.slug,
//         parentSlug: comment.slug,
//         title: crypto.randomUUID(),
//         slug: crypto.randomUUID(),
//         text: replyText.value,
//       },
//     };

//     await api().post(postUrl, p);
//     replyText.value = '';
//     replyingTo.value = null;

//     // Reload child comments for this comment
//     await loadChildComments(comment.slug);
//     emit('comment-added');
//   } catch (error) {
//     console.error('Reply submission failed:', error);
//   } finally {
//     submitting.value = false;
//   }
// };

// Image lightbox state
const imageLightboxOpen = ref(false);
const currentImage = ref<{ src: string; alt: string } | null>(null);

// handle content click
const handleContentClick = (event: Event) => {
  const target = event.target as HTMLElement;

  // handle image click
  if (target.tagName === 'IMG') {
    event.preventDefault();
    event.stopPropagation();

    const img = target as HTMLImageElement;
    currentImage.value = {
      src: img.src,
      alt: img.alt || img.title || 'Image',
    };
    imageLightboxOpen.value = true;
    return;
  }
};

// close image lightbox
const closeImageLightbox = () => {
  imageLightboxOpen.value = false;
  currentImage.value = null;
};
</script>

<style scoped>
.comment-tree {
  margin-left: 0;
}

.comment-item {
  margin-bottom: 1rem;
}

.comment-item:last-child {
  margin-bottom: 0;
}

.comment-card {
  background-color: #f8f9fa;
}

.reply-card {
  background-color: #f0f0f0;
}

.child-comments {
  border-left: 2px solid #e0e0e0;
  padding-left: 1rem;
}

.reply-form {
  margin-left: 2rem;
}

/* make image clickable */
.tiptap-content img {
  cursor: pointer;
  transition: opacity 0.2s ease;
}

.tiptap-content img:hover {
  opacity: 0.8;
}

/* lightbox style */
.lightbox-card {
  background: rgba(0, 0, 0, 0.9);
  border-radius: 0;
}

.lightbox-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 16px 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.lightbox-title {
  font-size: 16px;
  font-weight: 500;
  color: white;
}

.lightbox-close {
  color: white;
}

.lightbox-content {
  display: flex;
  justify-content: center;
  align-items: center;
  height: calc(100vh - 80px);
  padding: 0;
  background: transparent;
}

.media-container {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  padding: 20px;
}

.media-content {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  cursor: pointer;
}
</style>
