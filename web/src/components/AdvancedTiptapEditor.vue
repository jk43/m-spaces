<template>
  <div class="advanced-tiptap-editor">
    <!-- 메뉴바 -->

    <!-- 툴바 -->
    <div class="toolbar" v-if="!readonly">
      <!-- 텍스트 스타일 -->
      <q-btn-group flat>
        <q-btn-dropdown flat label="Heading" color="primary">
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

      <!-- 텍스트 서식 -->
      <q-btn-group flat>
        <q-btn
          flat
          round
          :icon="editor?.isActive('bold') ? 'format_bold' : 'format_bold'"
          :color="editor?.isActive('bold') ? 'primary' : 'grey'"
          @click="editor?.chain().focus().toggleBold().run()"
          size="sm"
        />
        <q-btn
          flat
          round
          :icon="editor?.isActive('italic') ? 'format_italic' : 'format_italic'"
          :color="editor?.isActive('italic') ? 'primary' : 'grey'"
          @click="editor?.chain().focus().toggleItalic().run()"
          size="sm"
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
        />
        <q-btn
          flat
          round
          :icon="editor?.isActive('code') ? 'code' : 'code'"
          :color="editor?.isActive('code') ? 'primary' : 'grey'"
          @click="editor?.chain().focus().toggleCode().run()"
          size="sm"
        />
      </q-btn-group>

      <q-separator vertical />

      <!-- 텍스트 색상 -->
      <q-btn-group flat>
        <q-btn-dropdown flat icon="format_color_text" color="primary">
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

        <q-btn-dropdown flat icon="format_color_fill" color="primary">
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

      <!-- 정렬 -->
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
        />
      </q-btn-group>

      <q-separator vertical />

      <!-- 목록 -->
      <q-btn-group flat>
        <q-btn
          flat
          round
          :icon="
            editor?.isActive('bulletList')
              ? 'format_list_bulleted'
              : 'format_list_bulleted'
          "
          :color="editor?.isActive('bulletList') ? 'primary' : 'grey'"
          @click="editor?.chain().focus().toggleBulletList().run()"
          size="sm"
        />
        <q-btn
          flat
          round
          :icon="
            editor?.isActive('orderedList')
              ? 'format_list_numbered'
              : 'format_list_numbered'
          "
          :color="editor?.isActive('orderedList') ? 'primary' : 'grey'"
          @click="editor?.chain().focus().toggleOrderedList().run()"
          size="sm"
        />
        <q-btn
          flat
          round
          :icon="editor?.isActive('taskList') ? 'checklist' : 'checklist'"
          :color="editor?.isActive('taskList') ? 'primary' : 'grey'"
          @click="editor?.chain().focus().toggleTaskList().run()"
          size="sm"
        />
      </q-btn-group>

      <q-separator vertical />

      <!-- 인용문 & 코드 블록 -->
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
        />
        <q-btn
          flat
          round
          :icon="editor?.isActive('codeBlock') ? 'code' : 'code'"
          :color="editor?.isActive('codeBlock') ? 'primary' : 'grey'"
          @click="editor?.chain().focus().toggleCodeBlock().run()"
          size="sm"
        />
      </q-btn-group>

      <q-separator vertical />

      <!-- 링크 & 이미지 -->
      <q-btn-group flat>
        <q-btn
          flat
          round
          icon="link"
          :color="editor?.isActive('link') ? 'primary' : 'grey'"
          @click="setLink"
          size="sm"
        />
        <q-btn flat round icon="image" @click="addImage" size="sm" />
        <q-btn flat round icon="table_chart" @click="insertTable" size="sm" />
      </q-btn-group>

      <q-separator vertical />

      <!-- 기타 도구 -->
      <q-btn-group flat>
        <q-btn
          flat
          round
          icon="format_clear"
          @click="editor?.chain().focus().clearNodes().unsetAllMarks().run()"
          size="sm"
        />
        <q-btn
          flat
          round
          icon="undo"
          @click="editor?.chain().focus().undo().run()"
          size="sm"
        />
        <q-btn
          flat
          round
          icon="redo"
          @click="editor?.chain().focus().redo().run()"
          size="sm"
        />
      </q-btn-group>
    </div>

    <!-- 에디터 영역 -->
    <div
      class="editor-content"
      :class="{ fullscreen: isFullscreen, 'drag-over': isDragOver }"
      @dragenter="handleDragEnter"
      @dragleave="handleDragLeave"
      @dragover="handleDragOver"
      @drop="handleDrop"
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

    <!-- 링크 입력 다이얼로그 -->
    <q-dialog v-model="linkDialog" persistent>
      <q-card style="min-width: 350px">
        <q-card-section>
          <div class="text-h6">Add Link</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-input
            v-model="linkUrl"
            label="URL"
            dense
            outlined
            @keyup.enter="confirmLink"
          />
        </q-card-section>

        <q-card-actions align="right" class="text-primary">
          <q-btn flat label="Cancel" v-close-popup />
          <q-btn flat label="Confirm" @click="confirmLink" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- 이미지 URL 입력 다이얼로그 -->
    <q-dialog v-model="imageDialog" persistent>
      <q-card style="min-width: 350px">
        <q-card-section>
          <div class="text-h6">Add Image</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-input
            v-model="imageUrl"
            label="이미지 URL"
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

    <!-- 테이블 삽입 다이얼로그 -->
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue';
import { useEditor, EditorContent } from '@tiptap/vue-3';
import StarterKit from '@tiptap/starter-kit';
import Placeholder from '@tiptap/extension-placeholder';
import Link from '@tiptap/extension-link';
import Image from '@tiptap/extension-image';
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
// tiptap-extension-resize-image import
import ResizeImage from 'tiptap-extension-resize-image';

interface Props {
  modelValue?: string;
  placeholder?: string;
  readonly?: boolean;
  height?: string;
}

interface Emits {
  (e: 'update:modelValue', value: string): void;
  (e: 'change', value: string): void;
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: 'Enter your content here...',
  readonly: false,
  height: '400px',
});

const emit = defineEmits<Emits>();

const linkDialog = ref(false);
const imageDialog = ref(false);
const tableDialog = ref(false);
const linkUrl = ref('');
const imageUrl = ref('');
const tableRows = ref(3);
const tableCols = ref(3);
const isFullscreen = ref(false);

// Drag and drop state
const isDragOver = ref(false);

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

const editor = useEditor({
  content: props.modelValue,
  editable: !props.readonly,
  extensions: [
    StarterKit.configure({
      codeBlock: false,
    }),
    Placeholder.configure({
      placeholder: props.placeholder,
    }),
    Link.configure({
      openOnClick: false,
      HTMLAttributes: {
        class: 'text-primary',
      },
    }),
    // Image extension 대신 ResizeImage 사용
    ResizeImage.configure({
      HTMLAttributes: {
        class: 'max-w-full h-auto',
      },
      // 옵션: minWidth, maxWidth, minHeight, maxHeight 등 필요시 추가
    }),
    TextAlign.configure({
      types: ['heading', 'paragraph'],
    }),
    Underline,
    TextStyle,
    Color,
    Highlight.configure({
      multicolor: true,
    }),
    TaskList,
    TaskItem.configure({
      nested: true,
    }),
    Table.configure({
      resizable: true,
    }),
    TableRow,
    TableHeader,
    TableCell,
    CodeBlock,
  ],
  onUpdate: ({ editor }) => {
    const html = editor.getHTML();
    emit('update:modelValue', html);
    emit('change', html);
  },
});

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

const setLink = () => {
  if (editor.value?.isActive('link')) {
    editor.value.chain().focus().unsetLink().run();
    return;
  }

  linkDialog.value = true;
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
    editor.value?.chain().focus().setImage({ src: imageUrl.value }).run();
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

const newDocument = () => {
  editor.value?.commands.setContent('');
};

const exportHTML = () => {
  const html = editor.value?.getHTML();
  const blob = new Blob([html || ''], { type: 'text/html' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = 'document.html';
  a.click();
  URL.revokeObjectURL(url);
};

const exportJSON = () => {
  const json = editor.value?.getJSON();
  const blob = new Blob([JSON.stringify(json, null, 2)], {
    type: 'application/json',
  });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = 'document.json';
  a.click();
  URL.revokeObjectURL(url);
};

const selectAll = () => {
  editor.value?.chain().focus().selectAll().run();
};

const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value;
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

const handleDrop = async (event: DragEvent) => {
  console.log('Debugging -handleDrop: ');
  event.preventDefault();
  isDragOver.value = false;

  if (!event.dataTransfer?.files) return;

  const files = Array.from(event.dataTransfer.files);

  for (const file of files) {
    await processDroppedFile(file);
  }
};

const processDroppedFile = async (file: File) => {
  console.log('Debugging -processDroppedFile: ');
  // Check if it's an image or video file
  if (file.type.startsWith('image/') || file.type.startsWith('video/')) {
    await handleImageFile(file);
  } else {
    // Handle other file types as download links
    await handleDownloadFile(file);
  }
};

const handleImageFile = async (file: File) => {
  console.log('Debugging -handleImageFile: ');
  try {
    // Create a preview URL for the image/video
    const fileUrl = URL.createObjectURL(file);

    if (file.type.startsWith('image/')) {
      // Insert the image into the editor
      editor.value
        ?.chain()
        .focus()
        .setImage({
          src: fileUrl,
          alt: file.name,
          title: file.name,
        })
        .run();
    } else if (file.type.startsWith('video/')) {
      // Insert video element
      editor.value
        ?.chain()
        .focus()
        .insertContent(
          `
        <video controls style="max-width: 100%; height: auto;">
          <source src="${fileUrl}" type="${file.type}">
          Your browser does not support the video tag.
        </video>
      `
        )
        .run();
    }

    // Clean up the object URL after a delay
    setTimeout(() => {
      URL.revokeObjectURL(fileUrl);
    }, 1000);
  } catch (error) {
    console.error('Error processing image/video file:', error);
  }
};

const handleDownloadFile = async (file: File) => {
  console.log('Debugging -handleDownloadFile: ');
  try {
    const fileUrl = URL.createObjectURL(file);

    // Insert download link with file icon
    editor.value
      ?.chain()
      .focus()
      .insertContent(
        `
      <a href="${fileUrl}" download="${
          file.name
        }" style="display: inline-flex; align-items: center; text-decoration: none; color: #1976d2; padding: 4px 8px; border: 1px solid #ddd; border-radius: 4px; background-color: #f8f9fa;">
        <i class="material-icons" style="margin-right: 8px; font-size: 18px;">download</i>
        <span>${file.name}</span>
        <span style="margin-left: 8px; font-size: 12px; color: #666;">(${(
          file.size / 1024
        ).toFixed(1)} KB)</span>
      </a>
    `
      )
      .run();

    // Clean up the object URL after a delay
    setTimeout(() => {
      URL.revokeObjectURL(fileUrl);
    }, 1000);
  } catch (error) {
    console.error('Error processing download file:', error);
  }
};

watch(
  () => props.modelValue,
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
</style>
