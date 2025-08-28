<template>
  <div class="advanced-tiptap-editor">
    <!-- Menu Bar -->

    <!-- Toolbar -->
    <div class="toolbar" v-if="!readonly">
      <!-- Text Style -->
      <q-btn-group flat>
        <q-btn-dropdown flat label="Heading" color="primary" size="md" dense>
          <q-list>
            <q-item
              clickable
              v-close-popup
              @click="editor?.chain().focus().toggleHeading({ level: 1 }).run()"
            >
              <q-item-section>
                <q-item-label>Heading 1</q-item-label>
              </q-item-section>
            </q-item>
            <q-item
              clickable
              v-close-popup
              @click="editor?.chain().focus().toggleHeading({ level: 2 }).run()"
            >
              <q-item-section>
                <q-item-label>Heading 2</q-item-label>
              </q-item-section>
            </q-item>
            <q-item
              clickable
              v-close-popup
              @click="editor?.chain().focus().toggleHeading({ level: 3 }).run()"
            >
              <q-item-section>
                <q-item-label>Heading 3</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-btn-dropdown>
      </q-btn-group>

      <q-separator vertical />

      <!-- Text Formatting -->
      <q-btn-group flat>
        <q-btn
          flat
          round
          :icon="editor?.isActive('bold') ? 'format_bold' : 'format_bold'"
          :color="editor?.isActive('bold') ? 'primary' : 'grey'"
          @click="editor?.chain().focus().toggleBold().run()"
          size="sm"
          dense
        />
        <q-btn
          flat
          round
          :icon="editor?.isActive('italic') ? 'format_italic' : 'format_italic'"
          :color="editor?.isActive('italic') ? 'primary' : 'grey'"
          @click="editor?.chain().focus().toggleItalic().run()"
          size="sm"
          dense
        />
        <q-btn
          flat
          round
          :icon="
            editor?.isActive('underline')
              ? 'format_underline'
              : 'format_underline'
          "
          :color="editor?.isActive('underline') ? 'primary' : 'grey'"
          @click="editor?.chain().focus().toggleUnderline().run()"
          size="sm"
          dense
        />
        <q-btn
          flat
          round
          :icon="
            editor?.isActive('strike') ? 'strikethrough_s' : 'strikethrough_s'
          "
          :color="editor?.isActive('strike') ? 'primary' : 'grey'"
          @click="editor?.chain().focus().toggleStrike().run()"
          size="sm"
          dense
        />
        <q-btn
          flat
          round
          :icon="editor?.isActive('code') ? 'code' : 'code'"
          :color="editor?.isActive('code') ? 'primary' : 'grey'"
          @click="editor?.chain().focus().toggleCode().run()"
          size="sm"
          dense
        />
      </q-btn-group>

      <q-separator vertical />

      <!-- Text Color -->
      <q-btn-group flat>
        <q-btn-dropdown
          flat
          icon="format_color_text"
          color="primary"
          size="sm"
          dense
        >
          <q-card style="min-width: 200px">
            <q-card-section>
              <div class="text-subtitle2">Text Color</div>
            </q-card-section>
            <q-card-section class="q-pt-none">
              <div class="row q-gutter-xs">
                <q-btn
                  v-for="color in textColors"
                  :key="color.name"
                  :style="{ backgroundColor: color.value }"
                  size="sm"
                  @click="setTextColor(color.value)"
                  class="color-btn"
                />
              </div>
              <q-separator class="q-my-sm" />
              <q-btn
                flat
                label="Reset Color"
                icon="clear"
                @click="resetTextColor"
                size="sm"
                class="full-width"
              />
            </q-card-section>
          </q-card>
        </q-btn-dropdown>

        <q-btn-dropdown
          flat
          icon="format_color_fill"
          color="primary"
          size="sm"
          dense
        >
          <q-card style="min-width: 200px">
            <q-card-section>
              <div class="text-subtitle2">Background Color</div>
            </q-card-section>
            <q-card-section class="q-pt-none">
              <div class="row q-gutter-xs">
                <q-btn
                  v-for="color in highlightColors"
                  :key="color.name"
                  :style="{ backgroundColor: color.value }"
                  size="sm"
                  @click="setHighlight(color.value)"
                  class="color-btn"
                />
              </div>
              <q-separator class="q-my-sm" />
              <q-btn
                flat
                label="Reset Background"
                icon="clear"
                @click="resetHighlight"
                size="sm"
                class="full-width"
              />
            </q-card-section>
          </q-card>
        </q-btn-dropdown>
      </q-btn-group>

      <q-separator vertical />

      <!-- Alignment -->
      <q-btn-group flat>
        <q-btn
          flat
          round
          :icon="
            editor?.isActive({ textAlign: 'left' })
              ? 'format_align_left'
              : 'format_align_left'
          "
          :color="editor?.isActive({ textAlign: 'left' }) ? 'primary' : 'grey'"
          @click="editor?.chain().focus().setTextAlign('left').run()"
          size="sm"
          dense
        />
        <q-btn
          flat
          round
          :icon="
            editor?.isActive({ textAlign: 'center' })
              ? 'format_align_center'
              : 'format_align_center'
          "
          :color="
            editor?.isActive({ textAlign: 'center' }) ? 'primary' : 'grey'
          "
          @click="editor?.chain().focus().setTextAlign('center').run()"
          size="sm"
          dense
        />
        <q-btn
          flat
          round
          :icon="
            editor?.isActive({ textAlign: 'right' })
              ? 'format_align_right'
              : 'format_align_right'
          "
          :color="editor?.isActive({ textAlign: 'right' }) ? 'primary' : 'grey'"
          @click="editor?.chain().focus().setTextAlign('right').run()"
          size="sm"
          dense
        />
        <q-btn
          flat
          round
          :icon="
            editor?.isActive({ textAlign: 'justify' })
              ? 'format_align_justify'
              : 'format_align_justify'
          "
          :color="
            editor?.isActive({ textAlign: 'justify' }) ? 'primary' : 'grey'
          "
          @click="editor?.chain().focus().setTextAlign('justify').run()"
          size="sm"
          dense
        />
      </q-btn-group>

      <q-separator vertical />

      <!-- Lists -->
      <q-btn-group flat>
        <q-btn
          flat
          icon="format_list_bulleted"
          color="primary"
          @click="editor?.chain().focus().toggleBulletList().run()"
          :class="{ 'bg-grey-3': editor?.isActive('bulletList') }"
          size="sm"
          dense
        >
          <q-tooltip>Bullet List</q-tooltip>
        </q-btn>
        <q-btn
          flat
          icon="format_list_numbered"
          color="primary"
          @click="editor?.chain().focus().toggleOrderedList().run()"
          :class="{ 'bg-grey-3': editor?.isActive('orderedList') }"
          size="sm"
          dense
        >
          <q-tooltip>Numbered List</q-tooltip>
        </q-btn>
        <q-btn
          flat
          icon="checklist"
          color="primary"
          @click="toggleTaskList"
          :class="{ 'bg-grey-3': editor?.isActive('taskList') }"
          size="sm"
          dense
        >
          <q-tooltip>Task List</q-tooltip>
        </q-btn>
      </q-btn-group>

      <q-separator vertical />

      <!-- Blockquote & Code Block -->
      <q-btn-group flat>
        <q-btn
          flat
          round
          :icon="
            editor?.isActive('blockquote') ? 'format_quote' : 'format_quote'
          "
          :color="editor?.isActive('blockquote') ? 'primary' : 'grey'"
          @click="editor?.chain().focus().toggleBlockquote().run()"
          size="sm"
          dense
        />
        <q-btn
          flat
          round
          :icon="editor?.isActive('codeBlock') ? 'code' : 'code'"
          :color="editor?.isActive('codeBlock') ? 'primary' : 'grey'"
          @click="editor?.chain().focus().toggleCodeBlock().run()"
          size="sm"
          dense
        />
      </q-btn-group>

      <q-separator vertical />

      <!-- Link & Image & YouTube -->
      <q-btn-group flat>
        <q-btn
          flat
          icon="link"
          color="primary"
          @click="openLinkDialog"
          :class="{ 'bg-grey-3': editor?.isActive('link') }"
          size="sm"
          dense
        >
          <q-tooltip>Insert Link</q-tooltip>
        </q-btn>

        <q-btn
          flat
          icon="image"
          color="primary"
          @click="addImage"
          size="sm"
          dense
        >
          <q-tooltip>Insert Image</q-tooltip>
        </q-btn>

        <q-btn
          flat
          icon="play_circle"
          color="primary"
          @click="openYouTubeDialog"
          size="sm"
          dense
        >
          <q-tooltip>Insert YouTube Video</q-tooltip>
        </q-btn>
      </q-btn-group>

      <q-separator vertical />

      <!-- Other Tools -->
      <q-btn-group flat>
        <q-btn
          flat
          round
          icon="format_clear"
          @click="editor?.chain().focus().clearNodes().unsetAllMarks().run()"
          size="sm"
          dense
        />
      </q-btn-group>
    </div>

    <!-- Editor Area -->
    <div
      class="editor-content"
      :class="{ fullscreen: isFullscreen, 'drag-over': isDragOver }"
      @dragenter="handleDragEnter"
      @dragleave="handleDragLeave"
      @dragover="handleDragOver"
      @drop="handleEditorDrop"
    >
      <editor-content :editor="editor" />

      <!-- Drag overlay -->
      <div v-if="isDragOver" class="drag-overlay">
        <div class="drag-message">
          <q-icon name="cloud_upload" size="48px" color="primary" />
          <div class="text-h6 q-mt-md">Drop files here</div>
          <div class="text-caption">
            Images will be inserted, text files will be added as content
          </div>
        </div>
      </div>
    </div>

    <!-- Link Input Dialog -->
    <q-dialog v-model="linkDialog" persistent>
      <q-card style="min-width: 350px">
        <q-card-section>
          <div class="text-h6">Insert Link</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-input
            v-model="linkText"
            label="Link Text"
            dense
            outlined
            placeholder="Enter link text"
            class="q-mb-sm"
          />
          <q-input
            v-model="linkUrl"
            label="URL"
            dense
            outlined
            @keyup.enter="setLink"
          />
        </q-card-section>

        <q-card-actions align="right" class="text-primary">
          <q-btn flat label="Cancel" @click="linkDialog = false" />
          <q-btn flat label="Insert" @click="setLink" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Image URL Input Dialog -->
    <q-dialog v-model="imageDialog" persistent>
      <q-card style="min-width: 350px">
        <q-card-section>
          <div class="text-h6">Add Image</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-input
            v-model="imageUrl"
            label="Image URL"
            dense
            outlined
            @keyup.enter="confirmImage"
          />
        </q-card-section>

        <q-card-actions align="right" class="text-primary">
          <q-btn flat label="Cancel" v-close-popup />
          <q-btn flat label="Confirm" @click="confirmImage" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Table Insert Dialog -->
    <q-dialog v-model="tableDialog" persistent>
      <q-card style="min-width: 300px">
        <q-card-section>
          <div class="text-h6">Insert Table</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <q-input
                v-model.number="tableRows"
                label="Rows"
                type="number"
                dense
                outlined
                min="1"
                max="10"
              />
            </div>
            <div class="col-6">
              <q-input
                v-model.number="tableCols"
                label="Columns"
                type="number"
                dense
                outlined
                min="1"
                max="10"
              />
            </div>
          </div>
        </q-card-section>

        <q-card-actions align="right" class="text-primary">
          <q-btn flat label="Cancel" v-close-popup />
          <q-btn flat label="Insert" @click="confirmTable" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Uploaded Files List -->
    <div v-if="uploadedFiles.length > 0" class="uploaded-files-section">
      <div class="uploaded-files-header">
        <q-icon name="attach_file" size="sm" />
        <span class="text-subtitle2 q-ml-sm">Uploaded Files</span>
        <span class="text-caption q-ml-sm"
          >({{ uploadedFiles.length }} items)</span
        >
      </div>

      <div class="uploaded-files-list">
        <div
          v-for="file in uploadedFiles"
          :key="file.id"
          class="uploaded-file-item"
          draggable="true"
          @dragstart="onFileDragStart(file)"
        >
          <div class="file-info">
            <q-icon
              :name="getFileIcon(file.contentType)"
              size="sm"
              :color="getFileColor(file.contentType)"
            />
            <div class="file-details">
              <div class="file-name">
                {{ file.name
                }}<span class="file-size"
                  >({{ formatFileSize(file.size) }})</span
                >
              </div>
            </div>
          </div>

          <div class="file-actions">
            <!-- Loading spinner -->
            <q-spinner
              v-if="file.loading"
              size="sm"
              color="primary"
              class="q-mr-sm"
            />

            <!-- Show delete button only after loading is complete -->
            <q-btn
              v-if="!file.loading"
              flat
              round
              size="sm"
              icon="delete"
              color="negative"
              @click="removeFile(file)"
            >
              <q-tooltip>Delete File</q-tooltip>
            </q-btn>
          </div>
        </div>
      </div>
    </div>

    <!-- YouTube URL Input Dialog -->
    <q-dialog v-model="youtubeDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">Insert YouTube Video</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-input
            v-model="youtubeUrl"
            label="YouTube URL"
            dense
            outlined
            placeholder="https://www.youtube.com/watch?v=..."
            @keyup.enter="insertYouTubeFromUrl"
          />
          <div class="text-caption q-mt-sm text-grey-6">
            Supported formats: youtube.com/watch?v=..., youtu.be/...,
            youtube.com/embed/...
          </div>
        </q-card-section>

        <q-card-actions align="right" class="text-primary">
          <q-btn flat label="Cancel" @click="youtubeDialog = false" />
          <q-btn flat label="Insert" @click="insertYouTubeFromUrl" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup lang="ts">
import { setupRules } from 'src/service/utils/form-builder';
import { api } from 'src/boot/axios';
import { useOrgStore } from 'src/stores/org-store';
import { Payload, StringAnyType } from 'src/service/types/payload';
import type { FileItem } from 'src/types';

import { reactive, ref, onBeforeMount, onBeforeUnmount, watch } from 'vue';
import { useEditor, EditorContent } from '@tiptap/vue-3';
import StarterKit from '@tiptap/starter-kit';
import Placeholder from '@tiptap/extension-placeholder';
import Link from '@tiptap/extension-link';
import TextAlign from '@tiptap/extension-text-align';
import Underline from '@tiptap/extension-underline';
import TextStyle from '@tiptap/extension-text-style';
import Color from '@tiptap/extension-color';
import Highlight from '@tiptap/extension-highlight';
import TaskList from '@tiptap/extension-task-list';
import TaskItem from '@tiptap/extension-task-item';
import Table from '@tiptap/extension-table';
import TableRow from '@tiptap/extension-table-row';
import TableCell from '@tiptap/extension-table-cell';
import TableHeader from '@tiptap/extension-table-header';
import CodeBlock from '@tiptap/extension-code-block';
import ResizeImage from 'tiptap-extension-resize-image';
import Video from 'tiptap-extension-video';
import { Node } from '@tiptap/core';

// Unused imports removed

// form builder
const props = defineProps([
  'modelValue',
  'placeholder',
  'readonly',
  'height',
  'data',
  'settings',
  'error',
  'globalData',
  'globalSettings',
  'submitClicked',
]);

const storeOrg = useOrgStore();
const cdnAddr = storeOrg.cdnAddr;
const formId = props.globalData.formId;
const slug = props.globalData.slug;

const data = reactive(props.data);
const settings = reactive(props.settings);

onBeforeMount(() => {
  loadFiles();
});
if (settings.callback !== undefined) {
  settings.callback(settings, data);
}

watch(
  () => props.submitClicked,
  (newVal) => {
    uploadedFiles.value = [];
  }
);

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
const loadFiles = async () => {
  let endpoint: string;
  if (slug.value) {
    endpoint = '/board/posts/files?slug=' + slug.value;
  } else {
    endpoint = '/board/posts/files?formId=' + formId.value;
  }
  const response = await api().get(endpoint);
  uploadedFiles.value = response.data.data;
};
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
const removeFile = async (file: FileItem) => {
  const deletePayload: Payload<{ [key: string]: string }> = {
    data: {
      data: {
        id: file.id,
      },
    },
  };
  await api().delete('/board/posts/file', deletePayload);
  uploadedFiles.value = uploadedFiles.value.filter((f) => f.id !== file.id);
};
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
const handleFile = async (file: File) => {
  try {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('service', 'cms_service');
    formData.append('serviceCtx', 'posts');
    formData.append('formId', formId.value);
    // Send the image data to the server
    const response = await api().post('/files', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        const percentCompleted = Math.round(
          (progressEvent.loaded * 100) / progressEvent.total
        );
        console.log('Upload progress:', percentCompleted + '%');
      },
    });
    return response.data.data[0];
  } catch (error) {
    console.error('Error uploading file:', error);
    throw error;
  }
};
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// form builder end

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
  (e: 'change', value: string): void;
}>();

const linkDialog = ref(false);
const imageDialog = ref(false);
const tableDialog = ref(false);
const linkUrl = ref('');
const imageUrl = ref('');
const tableRows = ref(3);
const tableCols = ref(3);
const isFullscreen = ref(false);
const linkText = ref('');

// Color palette
const textColors = [
  { name: 'Black', value: '#000000' },
  { name: 'Red', value: '#ff0000' },
  { name: 'Blue', value: '#0000ff' },
  { name: 'Green', value: '#008000' },
  { name: 'Purple', value: '#800080' },
  { name: 'Orange', value: '#ffa500' },
  { name: 'Gray', value: '#808080' },
  { name: 'Brown', value: '#a52a2a' },
];

const highlightColors = [
  { name: 'Yellow', value: '#ffff00' },
  { name: 'Light Green', value: '#90ee90' },
  { name: 'Light Blue', value: '#87ceeb' },
  { name: 'Light Orange', value: '#ffd700' },
  { name: 'Light Purple', value: '#dda0dd' },
  { name: 'Light Gray', value: '#d3d3d3' },
];

// Custom Download File Node
const DownloadFile = Node.create({
  name: 'downloadFile',
  group: 'block',
  content: 'inline*',

  parseHTML() {
    return [
      {
        tag: 'div.download-file',
      },
    ];
  },

  renderHTML({ HTMLAttributes }) {
    return ['div', { class: 'download-file', ...HTMLAttributes }, 0];
  },
});

// Custom YouTube Video Node
const YouTubeVideo = Node.create({
  name: 'youtubeVideo',
  group: 'block',

  addAttributes() {
    return {
      videoId: {
        default: null,
      },
    };
  },

  parseHTML() {
    return [
      {
        tag: 'div.youtube-video-container',
      },
    ];
  },

  renderHTML({ HTMLAttributes }) {
    const videoId = HTMLAttributes.videoId;
    const embedUrl = `https://www.youtube.com/embed/${videoId}`;

    return [
      'div',
      {
        class: 'youtube-video-container',
        style:
          'position: relative; padding-bottom: 56.25%; height: 0; overflow: hidden; max-width: 100%; margin: 20px 0; border-radius: 8px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);',
      },
      [
        'iframe',
        {
          src: embedUrl,
          style:
            'position: absolute; top: 0; left: 0; width: 100%; height: 100%; border: 0; border-radius: 8px;',
          allowfullscreen: true,
          title: 'YouTube video player',
        },
      ],
    ];
  },
});

// Custom Download Link Node
const DownloadLink = Link.extend({
  addAttributes() {
    return {
      ...this.parent?.(),
      download: {
        default: null,
        parseHTML: (element) => element.getAttribute('download'),
        renderHTML: (attributes) => {
          if (!attributes.download) {
            return {};
          }

          return {
            download: attributes.download,
          };
        },
      },
    };
  },
});

const editor = useEditor({
  content: props.modelValue || (data && data.value ? data.value : ''),
  editable: !props.readonly,
  extensions: [
    StarterKit.configure({
      codeBlock: false,
    }),
    Placeholder.configure({
      placeholder: props.placeholder,
    }),
    // Remove Link extension and use DownloadLink only
    DownloadLink.configure({
      openOnClick: false,
      HTMLAttributes: {
        class: 'text-primary',
      },
    }),
    ResizeImage.configure({
      HTMLAttributes: {
        class: 'max-w-full h-auto',
        style: 'display: block; margin: 0;',
      },
      // Remove div wrapper and render only image
      renderHTML({ HTMLAttributes }) {
        return ['img', HTMLAttributes];
      },
    }),
    TextAlign.configure({
      types: ['heading', 'paragraph'],
    }),
    Underline,
    TextStyle,
    Color,
    Highlight,
    Table.configure({
      resizable: true,
    }),
    TableRow,
    TableCell,
    TableHeader,
    CodeBlock,
    Video.configure({
      HTMLAttributes: {
        class: 'max-w-full h-auto',
        controls: true,
        crossorigin: 'anonymous',
      },
      addAttributes() {
        return {
          ...this.parent?.(),
          src: {
            default: null,
          },
          type: {
            default: null,
          },
          controls: {
            default: true,
          },
          autoplay: {
            default: false,
          },
          muted: {
            default: false,
          },
        };
      },
    }),
    TaskList.configure({
      HTMLAttributes: {
        class: 'task-list',
      },
    }),
    TaskItem.configure({
      HTMLAttributes: {
        class: 'task-item',
      },
    }),
    DownloadFile,
    YouTubeVideo,
    // DownloadLink is already registered above, so remove it here
  ],
  onUpdate: ({ editor }) => {
    const html = editor.getHTML();
    emit('update:modelValue', html);
    emit('change', html);

    if (data && data.value !== undefined) {
      data.value = html;
    }
  },
});

// Drag and drop state
const isDragOver = ref(false);

// Uploaded files management
const uploadedFiles = ref<FileItem[]>([]);
const draggedFile = ref<FileItem | null>(null);

// File drag start
const onFileDragStart = (file: FileItem) => {
  draggedFile.value = file;
};

// Editor area drop handler
const handleEditorDrop = (event: DragEvent) => {
  event.preventDefault();
  isDragOver.value = false;

  // If dragging from uploaded files
  if (draggedFile.value) {
    handleImageFile(draggedFile.value);
    handleDownloadFile(draggedFile.value);
    draggedFile.value = null;
    return;
  }

  // If external files, call handleDrop
  if (event.dataTransfer?.files && event.dataTransfer.files.length > 0) {
    console.log('External files detected');
    handleDrop(event);
  }
};

const handleDrop = async (event: DragEvent) => {
  console.log('Debugging -handleDrop: ');
  event.preventDefault();
  isDragOver.value = false;

  if (!event.dataTransfer?.files) return;

  const files = Array.from(event.dataTransfer.files);
  console.log('External files:', files);

  for (const file of files) {
    await processDroppedFile(file);
    await new Promise((resolve) => setTimeout(resolve, 1000));
    const currentPos = editor.value?.state.selection.from || 0;
    editor.value
      ?.chain()
      .focus()
      .insertContentAt(currentPos + 1, '\n')
      .run();
  }
};

const getFileIcon = (type: string): string => {
  if (type.startsWith('image/')) return 'image';
  if (type.startsWith('video/')) return 'movie';
  if (type.startsWith('application/pdf')) return 'picture_as_pdf';
  return 'insert_drive_file';
};

const getFileColor = (type: string): string => {
  if (type.startsWith('image/')) return 'primary';
  if (type.startsWith('video/')) return 'deep-orange';
  if (type.startsWith('application/pdf')) return 'red';
  return 'grey';
};

const formatFileSize = (size: number): string => {
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
  return `${(size / (1024 * 1024)).toFixed(1)} MB`;
};

const processDroppedFile = async (file: File) => {
  console.log('Debugging -processDroppedFile: ', file.name);

  uploadedFiles.value.push({
    id: '',
    name: file.name,
    size: file.size,
    contentType: file.type,
    s3Path: '',
    loading: true,
  });
  const i = uploadedFiles.value.length - 1;
  const fileInfo = await handleFile(file);
  uploadedFiles.value[i] = {
    id: fileInfo.uuid,
    name: fileInfo.name,
    size: fileInfo.size,
    contentType: fileInfo.contentType,
    s3Path: fileInfo.s3Path,
    loading: false,
  };
  console.log('Debugging - uploadedFiles.value: ', uploadedFiles.value);

  if (file.type.startsWith('image/') || file.type.startsWith('video/')) {
    await handleImageFile(fileInfo);
  } else {
    await handleDownloadFile(fileInfo);
  }
};

const handleImageFile = async (file: any) => {
  try {
    const fileUrl = cdnAddr + '/' + file.s3Path;
    if (file.contentType.startsWith('image/')) {
      const currentPos = editor.value?.state.selection.from || 0;
      // Add display: block to center the image properly
      editor.value
        ?.chain()
        .focus()
        .insertContent(
          `<img src="${fileUrl}" alt="${file.name}" title="${file.name}" style="width: 300px; height: auto; display: block; margin: 0;" />`
        )
        .run();
    } else if (file.contentType.startsWith('video/')) {
      // Save current cursor position
      const currentPos = editor.value?.state.selection.from || 0;
      editor.value
        ?.chain()
        .focus()
        .setVideo({
          src: fileUrl,
          type: file.type,
          controls: true,
          autoplay: false,
          muted: false,
        })
        .run();
    }
  } catch (error) {
    console.error('Error processing image/video file:', file.name, error);
  }
};

const handleDownloadFile = async (file: File) => {
  if (
    file.contentType.startsWith('image/') ||
    file.contentType.startsWith('video/')
  ) {
    return;
  }
  console.log('Debugging -handleDownloadFile: ', file.name);
  try {
    const fileUrl = cdnAddr + '/' + file.s3Path;
    console.log('Debugging - file: ', file);
    // Insert download link with file icon
    editor.value
      ?.chain()
      .focus()
      .insertContent(
        `<p><div class="download-file"><a href="${fileUrl}" download="${
          file.name
        }">${file.name} (${formatFileSize(file.size)})</a></div></p>`
      )
      .run();
  } catch (error) {
    console.error('Error processing download file:', file.name, error);
  }
};

// Link button click handler
const openLinkDialog = () => {
  // Set link text with selected text
  const selectedText = editor.value?.state.doc.textBetween(
    editor.value.state.selection.from,
    editor.value.state.selection.to
  );
  linkText.value = selectedText || '';
  linkDialog.value = true;
};

const setLink = () => {
  if (linkUrl.value) {
    // YouTube URL processing
    if (isYouTubeUrl(linkUrl.value)) {
      const videoId = extractYouTubeVideoId(linkUrl.value);
      if (videoId) {
        insertYouTubeVideo(videoId);
        linkDialog.value = false;
        linkUrl.value = '';
        linkText.value = '';
        return;
      }
    }

    // Regular link processing
    if (linkText.value) {
      // Insert link text with text
      editor.value
        ?.chain()
        .focus()
        .insertContent(`<a href="${linkUrl.value}">${linkText.value}</a>`)
        .run();
    } else {
      // Insert link text with URL
      editor.value
        ?.chain()
        .focus()
        .insertContent(`<a href="${linkUrl.value}">${linkUrl.value}</a>`)
        .run();
    }

    linkDialog.value = false;
    linkUrl.value = '';
    linkText.value = '';
  }
};

const removeLink = () => {
  editor.value?.chain().focus().unsetLink().run();
};

// YouTube URL check function
const isYouTubeUrl = (url: string): boolean => {
  const youtubeRegex = /^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.be)\/.+/;
  return youtubeRegex.test(url);
};

// YouTube video ID extraction function
const extractYouTubeVideoId = (url: string): string | null => {
  const patterns = [
    /(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/)([^&\n?#]+)/,
    /youtube\.com\/v\/([^&\n?#]+)/,
  ];

  for (const pattern of patterns) {
    const match = url.match(pattern);
    if (match) {
      return match[1];
    }
  }

  return null;
};

// YouTube video insertion function
const insertYouTubeVideo = (videoId: string) => {
  editor.value
    ?.chain()
    .focus()
    .insertContent({
      type: 'youtubeVideo',
      attrs: {
        videoId: videoId,
      },
    })
    .run();
};

const addImage = () => {
  imageDialog.value = true;
};

const insertTable = () => {
  tableDialog.value = true;
};

const confirmLink = () => {
  if (linkUrl.value) {
    editor.value?.chain().focus().setLink({ href: linkUrl.value }).run();
  }
  linkDialog.value = false;
  linkUrl.value = '';
};

const confirmImage = () => {
  if (imageUrl.value) {
    editor.value
      ?.chain()
      .focus()
      .setImage({ src: imageUrl.value, width: '200px' })
      .run();
  }
  imageDialog.value = false;
  imageUrl.value = '';
};

const confirmTable = () => {
  editor.value
    ?.chain()
    .focus()
    .insertTable({ rows: tableRows.value, cols: tableCols.value })
    .run();
  tableDialog.value = false;
};

// Drag and drop handlers
const handleDragEnter = (event: DragEvent) => {
  event.preventDefault();
  isDragOver.value = true;
};

const handleDragLeave = (event: DragEvent) => {
  event.preventDefault();
  // Only set to false if we're leaving the editor area completely
  const currentTarget = event.currentTarget as HTMLElement;
  const relatedTarget = event.relatedTarget as Node;
  if (!currentTarget?.contains(relatedTarget)) {
    isDragOver.value = false;
  }
};

const handleDragOver = (event: DragEvent) => {
  console.log('Debugging -handleDragOver: ');
  event.preventDefault();
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = 'copy';
  }
};

// TaskList toggle function
const toggleTaskList = () => {
  editor.value?.chain().focus().toggleTaskList().run();
};

watch(
  () => props.modelValue,
  (newValue) => {
    if (editor.value && editor.value.getHTML() !== newValue) {
      editor.value.commands.setContent(newValue || '');
    }
  }
);

// data.value change detection
watch(
  () => data?.value,
  (newValue) => {
    if (editor.value && editor.value.getHTML() !== newValue) {
      editor.value.commands.setContent(newValue || '');
    }
  }
);

watch(
  () => props.readonly,
  (newValue) => {
    if (editor.value) {
      editor.value.setEditable(!newValue);
    }
  }
);

// Add color setting functions
const setTextColor = (color: string) => {
  editor.value?.chain().focus().setColor(color).run();
};

const setHighlight = (color: string) => {
  editor.value?.chain().focus().setHighlight({ color }).run();
};

const resetTextColor = () => {
  editor.value?.chain().focus().unsetColor().run();
};

const resetHighlight = () => {
  editor.value?.chain().focus().unsetHighlight().run();
};

const youtubeDialog = ref(false);
const youtubeUrl = ref('');

// Open YouTube dialog
const openYouTubeDialog = () => {
  youtubeDialog.value = true;
};

// Insert YouTube video from URL
const insertYouTubeFromUrl = () => {
  if (youtubeUrl.value) {
    if (isYouTubeUrl(youtubeUrl.value)) {
      const videoId = extractYouTubeVideoId(youtubeUrl.value);
      if (videoId) {
        insertYouTubeVideo(videoId);
        youtubeDialog.value = false;
        youtubeUrl.value = '';
      } else {
        // Show error message
        console.error('Invalid YouTube URL');
      }
    } else {
      // Show error message
      console.error('Not a valid YouTube URL');
    }
  }
};

onBeforeUnmount(() => {
  editor.value?.destroy();
});
</script>

<style scoped>
.advanced-tiptap-editor {
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
}

.menubar {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  background-color: #f5f5f5;
  border-bottom: 1px solid #ddd;
}

.toolbar {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 1px 5px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #ddd;
  flex-wrap: wrap;
}

.editor-content {
  min-height: v-bind(height);
  max-height: 600px;
  overflow-y: auto;
  transition: all 0.3s ease;
}

.editor-content.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 9999;
  background: white;
  border: none;
  border-radius: 0;
}

.editor-content.drag-over {
  border: 2px dashed #1976d2;
  background-color: #e3f2fd;
}

.drag-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: rgba(255, 255, 255, 0.9);
  z-index: 10;
  border-radius: 8px;
}

.drag-message {
  text-align: center;
  color: #333;
  padding: 20px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.color-btn {
  width: 10px;
  height: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.uploaded-files-section {
  margin-top: 16px;
  padding: 12px;
  border: 0px solid #eee;
  border-radius: 0px;
  background: #fafbfc;
}
.uploaded-files-header {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}
.uploaded-files-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.uploaded-file-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 0;
  border-bottom: 1px solid #f0f0f0;
}
.file-info {
  display: flex;
  align-items: center;
  gap: 8px;
}
.file-details {
  display: flex;
  flex-direction: column;
}
.file-name {
  font-weight: 500;
}
.file-size {
  font-size: 12px;
  color: #888;
  margin-left: 5px;
}
.file-actions {
  display: flex;
  align-items: center;
}

.youtube-video-container {
  position: relative;
  padding-bottom: 56.25%; /* 16:9 ratio */
  height: 0;
  overflow: hidden;
  max-width: 100%;
  margin: 20px 0;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.youtube-video-container iframe {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border: 0;
  border-radius: 8px;
}

/* TaskList styles */
:deep(.task-list) {
  list-style: none;
  padding: 0;
  margin: 0;
}

:deep(.task-item) {
  display: flex;
  align-items: flex-start;
  margin: 0.5em 0;
}

:deep(.task-item input[type='checkbox']) {
  margin-right: 0.5em;
  margin-top: 0.25em;
}

:deep(.task-item p) {
  margin: 0;
  flex: 1;
}

/* TipTap 에디터와 프리뷰에서 이미지 왼쪽 정렬 */
.ProseMirror img,
.tiptap-content img,
[contenteditable='true'] img {
  display: block;
  margin: 0;
}

/* 또는 더 구체적으로 */
.ProseMirror img[style*='margin: 0px auto'],
.tiptap-content img[style*='margin: 0px auto'],
[contenteditable='true'] img[style*='margin: 0px auto'] {
  display: block !important;
  margin: 0 !important;
}
</style>
