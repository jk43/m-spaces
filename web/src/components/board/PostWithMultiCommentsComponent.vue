<template>
  <div class="q-pa-md">
    <q-card class="post-card">
      <q-card-section class="q-pb-none">
        <div class="text-h5 text-weight-bold q-mb-md">{{ post.title }}</div>
        <div
          class="text-body1 q-mb-lg tiptap-content"
          v-html="post.text"
          @click="handleContentClick"
        ></div>

        <!-- Post meta information -->
        <div class="row items-center q-mb-md">
          <q-avatar size="32px" class="q-mr-sm">
            <q-icon name="person" />
          </q-avatar>
          <div class="text-caption text-grey-6">
            Author: {{ post.user_name?.first_name }}
            {{ post.user_name?.last_name }}
          </div>
          <q-space />
          <div class="text-caption text-grey-6">
            {{ formatDate(post.created_at) }}
          </div>
        </div>
      </q-card-section>
      <FilesComponent :slug="slug" />
      <q-separator />

      <!-- Comments section -->
      <q-card-section>
        <div class="text-h6 q-mb-md">Comments</div>

        <!-- Main comment form -->
        <q-card flat bordered class="q-mt-md">
          <q-card-section>
            <CommentForm
              :parent-board="board"
              :parent-post="post"
              @comment-added="loadData"
            />
          </q-card-section>
        </q-card>

        <!-- Comments tree -->
        <div v-if="comments.length > 0" class="comments-tree q-mt-md">
          <CommentTree
            :replying-to="replyingTo"
            :comments="comments"
            :parent-board="board"
            :parent-post="post"
            @comment-added="loadData"
          />
        </div>

        <!-- When no comments -->
        <div v-else class="text-center q-pa-lg">
          <q-icon name="chat_bubble_outline" size="48px" color="grey-4" />
          <div class="text-grey-6 q-mt-sm">No comments yet.</div>
        </div>
      </q-card-section>
    </q-card>

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
  {{ comments }}
</template>

<script setup lang="ts">
import { api } from 'src/boot/axios';
import { Response } from 'src/types';

import { ref, reactive, onMounted, nextTick, onBeforeUnmount } from 'vue';
import { formatDistanceToNow } from 'date-fns';
import CommentTree from './CommentTreeComponent.vue';
import FilesComponent from './FilesComponent.vue';
import CommentForm from './CommentFormComponent.vue';

const getUrl = '/board/posts/comments';

const props = defineProps<{
  post: Record<string, unknown>;
  board: Record<string, unknown>;
}>();

const slug = props.post.slug as string;

const replyingTo = reactive({ slug: '' });
const post = ref(props.post);
const comments = ref<Record<string, unknown>[]>([]);
const dataReady = ref(false);
const contentRef = ref<HTMLElement>();

// Image lightbox state
const imageLightboxOpen = ref(false);
const currentImage = ref<{ src: string; alt: string } | null>(null);

const queryParams = ref<{ [key: string]: string | number | boolean }>({
  slug: slug,
  sortBy: 'created_at',
  descending: true,
  page: 1,
  rowsPerPage: 25,
  rowsNumber: 0,
});

// Date format function
const formatDate = (dateString: string): string => {
  if (!dateString) return '';
  try {
    return formatDistanceToNow(new Date(dateString), { addSuffix: true });
  } catch {
    return new Date(dateString).toLocaleDateString('en-US');
  }
};

// Load comments on component mount
onMounted(() => {
  loadData();

  // nextTick is used to wait for the component to be mounted before adding the event listener.
  nextTick(() => {
    if (contentRef.value) {
      contentRef.value.addEventListener('click', handleContentClick);
    }
  });
});

// remove event listener when the component is unmounted
onBeforeUnmount(() => {
  if (contentRef.value) {
    contentRef.value.removeEventListener('click', handleContentClick);
  }
});

const loadData = async () => {
  try {
    const postRes = await api().get(`/board/post?slug=${slug}`);
    post.value = postRes.data.data.post;
    //comments
    const q = new URLSearchParams();
    Object.entries(queryParams.value).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        q.append(key, String(value));
      }
    });
    const commentsRes: Response<{ [key: string]: unknown }> = await api().get(
      getUrl + '?' + q.toString()
    );
    comments.value = commentsRes.data.data.rows || [];
    queryParams.value.rowsNumber = commentsRes.data.data.total || 0;
    dataReady.value = true;
  } catch (error) {
    console.error('Failed to load comments:', error);
  }
};

// handle content click
const handleContentClick = (event: Event) => {
  const target = event.target as HTMLElement;

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

  // handle file download
  if (target.tagName === 'A' && target.closest('.download-file')) {
    event.preventDefault();

    const link = target as HTMLAnchorElement;
    const href = link.getAttribute('href');
    const download = link.getAttribute('download');

    if (href) {
      handleFileDownload(href, download || 'file');
    }
  }
};

// close image lightbox
const closeImageLightbox = () => {
  imageLightboxOpen.value = false;
  currentImage.value = null;
};

// handle file download
const handleFileDownload = async (fileUrl: string, fileName: string) => {
  try {
    console.log('Downloading file:', fileName, 'from:', fileUrl);

    // if the file name is undefined, extract the file name from the url
    let actualFileName = fileName;
    if (fileName === 'undefined' || !fileName) {
      const urlParts = fileUrl.split('/');
      actualFileName = urlParts[urlParts.length - 1] || 'download';
    }

    // check the file extension
    const fileExtension = actualFileName.split('.').pop()?.toLowerCase();

    // if the file is an image or video, open it in a new tab
    if (
      fileExtension &&
      ['jpg', 'jpeg', 'png', 'gif', 'webp', 'mp4', 'webm', 'ogg'].includes(
        fileExtension
      )
    ) {
      window.open(fileUrl, '_blank');
      return;
    }

    // download other files
    const response = await fetch(fileUrl);
    const blob = await response.blob();

    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = actualFileName;
    link.style.display = 'none';

    document.body.appendChild(link);
    link.click();

    // clean up
    setTimeout(() => {
      document.body.removeChild(link);
      URL.revokeObjectURL(link.href);
    }, 100);
  } catch (error) {
    console.error('Download failed:', error);
    // fallback: open in a new tab
    window.open(fileUrl, '_blank');
  }
};
</script>

<style scoped>
.post-card {
  max-width: 800px;
  margin: 0 auto;
}

.comments-tree {
  margin-top: 1rem;
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
}

/* make image clickable */
.media-content {
  cursor: pointer;
}
</style>
