<template>
  <div v-if="uploadedFiles.length > 0" class="files-section">
    <q-separator class="q-my-md" />

    <q-card-section>
      <div class="files-header">
        <q-icon name="attach_file" size="sm" color="primary" />
        <span class="text-subtitle2 q-ml-sm">Attached Files</span>
        <span class="text-caption q-ml-sm text-grey-6"
          >({{ uploadedFiles.length }} items)</span
        >
      </div>

      <div class="files-grid q-mt-md">
        <div
          v-for="(file, index) in uploadedFiles"
          :key="file.id"
          class="file-item"
          @click="handleFileClick(file, index)"
        >
          <div class="file-icon" :class="{ loading: file.loading }">
            <!-- 이미지인 경우 썸네일 표시 -->
            <img
              v-if="isImageFile(file)"
              :src="getFileUrl(file)"
              :alt="file.name"
              class="file-thumbnail"
              @load="handleImageLoad(file)"
              @error="handleImageError"
            />
            <!-- 비디오인 경우 썸네일 표시 -->
            <div
              v-else-if="isVideoFile(file)"
              class="video-thumbnail-container"
            >
              <img
                v-if="file.thumbnail"
                :src="file.thumbnail"
                :alt="file.name"
                class="file-thumbnail"
                @error="handleVideoThumbnailError(file)"
              />
              <video
                v-else
                :src="getFileUrl(file)"
                class="video-thumbnail"
                preload="metadata"
                @loadeddata="handleVideoLoad(file, $event)"
                @error="handleVideoError(file)"
                muted
                playsinline
              />
              <!-- 비디오 재생 버튼 오버레이 -->
              <div class="video-play-overlay">
                <q-icon name="play_arrow" size="sm" color="white" />
              </div>
            </div>
            <!-- 이미지/비디오가 아닌 경우 아이콘 표시 -->
            <q-icon
              v-else
              :name="getFileIcon(file.contentType)"
              size="lg"
              :color="getFileColor(file.contentType)"
            />
          </div>

          <div class="file-info">
            <div class="file-name" :title="file.name">
              {{ truncateFileName(file.name) }}
            </div>
            <div class="file-meta">
              <span class="file-size">{{ formatFileSize(file.size) }}</span>
              <span class="file-date">{{ formatDate(file.created_at) }}</span>
            </div>
          </div>

          <div class="file-actions">
            <q-btn
              flat
              round
              size="sm"
              icon="download"
              color="primary"
              @click.stop="downloadFile(file)"
            >
              <q-tooltip>Download File</q-tooltip>
            </q-btn>
          </div>
        </div>
      </div>
    </q-card-section>

    <!-- Lightbox Dialog -->
    <q-dialog v-model="lightboxOpen" full-width full-height>
      <q-card class="lightbox-card">
        <q-card-section class="lightbox-header">
          <div class="lightbox-title">
            {{ currentFile?.name }}
          </div>
          <q-btn
            flat
            round
            icon="close"
            @click="closeLightbox"
            class="lightbox-close"
          />
        </q-card-section>

        <q-card-section class="lightbox-content">
          <!-- 이미지 표시 -->
          <div v-if="isImageFile(currentFile)" class="media-container">
            <img
              :src="getFileUrl(currentFile)"
              :alt="currentFile?.name"
              class="media-content"
              @click="closeLightbox"
            />
          </div>

          <!-- 비디오 표시 -->
          <div v-else-if="isVideoFile(currentFile)" class="media-container">
            <video
              :src="getFileUrl(currentFile)"
              controls
              class="media-content"
              @click.stop
            >
              Your browser does not support the video tag.
            </video>
          </div>

          <!-- PDF 표시 -->
          <div v-else-if="isPdfFile(currentFile)" class="media-container">
            <iframe
              :src="getFileUrl(currentFile)"
              class="media-content"
              frameborder="0"
            ></iframe>
          </div>
        </q-card-section>

        <!-- Navigation Buttons -->
        <div v-if="mediaFiles.length > 1" class="lightbox-navigation">
          <q-btn
            flat
            round
            icon="chevron_left"
            size="lg"
            color="white"
            class="nav-btn nav-prev"
            @click="previousFile"
          />
          <q-btn
            flat
            round
            icon="chevron_right"
            size="lg"
            color="white"
            class="nav-btn nav-next"
            @click="nextFile"
          />
        </div>

        <!-- File Counter -->
        <div v-if="mediaFiles.length > 1" class="file-counter">
          {{ currentIndex + 1 }} / {{ mediaFiles.length }}
        </div>
      </q-card>
    </q-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, onBeforeMount, computed, onMounted, onBeforeUnmount } from 'vue';
import { api } from 'src/boot/axios';
import { useOrgStore } from 'src/stores/org-store';
import { formatDistanceToNow } from 'date-fns';
import type { FileItem } from 'src/types';

const props = defineProps<{
  slug: string;
}>();

const storeOrg = useOrgStore();
const cdnAddr = storeOrg.cdnAddr;
const slug = props.slug;
const uploadedFiles = ref<FileItem[]>([]);

// Lightbox state
const lightboxOpen = ref(false);
const currentIndex = ref(0);
const currentFile = ref<FileItem | null>(null);

// Filter only media files
const mediaFiles = computed(() => {
  return uploadedFiles.value.filter(
    (file) => isImageFile(file) || isVideoFile(file) || isPdfFile(file)
  );
});

onBeforeMount(() => {
  loadFiles();
});

const loadFiles = async (): Promise<void> => {
  try {
    const response = await api().get('/board/posts/files?slug=' + slug);
    // Add loading state for image and video files
    uploadedFiles.value = response.data.data.map((file: FileItem) => ({
      ...file,
      loading: isImageFile(file) || isVideoFile(file),
      thumbnail: null,
    }));
  } catch (error) {
    console.error('Failed to load files:', error);
  }
};

// File type check functions
const isImageFile = (file: FileItem): boolean => {
  return file?.contentType?.startsWith('image/') || false;
};

const isVideoFile = (file: FileItem): boolean => {
  return file?.contentType?.startsWith('video/') || false;
};

const isPdfFile = (file: FileItem): boolean => {
  return file?.contentType?.includes('pdf') || false;
};

// Get file URL
const getFileUrl = (file: FileItem): string => {
  return cdnAddr + '/' + file.s3Path;
};

// Handle file click
const handleFileClick = (file: FileItem, index: number): void => {
  if (isImageFile(file) || isVideoFile(file) || isPdfFile(file)) {
    // Open lightbox for media files
    const mediaIndex = mediaFiles.value.findIndex((f) => f.id === file.id);
    if (mediaIndex !== -1) {
      currentIndex.value = mediaIndex;
      currentFile.value = file;
      lightboxOpen.value = true;
    }
  } else {
    // Download other files
    downloadFile(file);
  }
};

// Close lightbox
const closeLightbox = (): void => {
  lightboxOpen.value = false;
  currentFile.value = null;
};

// Previous file (infinite loop)
const previousFile = (): void => {
  if (currentIndex.value > 0) {
    currentIndex.value--;
  } else {
    // Go to last slide when clicking previous on first slide
    currentIndex.value = mediaFiles.value.length - 1;
  }
  currentFile.value = mediaFiles.value[currentIndex.value];
};

// Next file (infinite loop)
const nextFile = (): void => {
  if (currentIndex.value < mediaFiles.value.length - 1) {
    currentIndex.value++;
  } else {
    // Go to first slide when clicking next on last slide
    currentIndex.value = 0;
  }
  currentFile.value = mediaFiles.value[currentIndex.value];
};

// Handle keyboard events (infinite loop support)
const handleKeydown = (event: KeyboardEvent): void => {
  if (!lightboxOpen.value) return;

  switch (event.key) {
    case 'Escape':
      closeLightbox();
      break;
    case 'ArrowLeft':
      previousFile();
      break;
    case 'ArrowRight':
      nextFile();
      break;
  }
};

// Add/remove keyboard event listeners
onMounted(() => {
  document.addEventListener('keydown', handleKeydown);
});

onBeforeUnmount(() => {
  document.removeEventListener('keydown', handleKeydown);
});

// Handle video load
const handleVideoLoad = async (file: FileItem, event: Event): Promise<void> => {
  try {
    const video = event.target as HTMLVideoElement;

    // Capture first frame of video to generate thumbnail
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');

    if (ctx) {
      canvas.width = 48;
      canvas.height = 48;

      // Draw video to canvas
      ctx.drawImage(video, 0, 0, canvas.width, canvas.height);

      // Convert canvas to data URL
      const thumbnailUrl = canvas.toDataURL('image/jpeg', 0.8);
      file.thumbnail = thumbnailUrl;
      file.loading = false;
    }
  } catch (error) {
    console.error('Failed to generate video thumbnail:', error);
    handleVideoError(file);
  }
};

// Handle video error
const handleVideoError = (file: FileItem): void => {
  file.loading = false;
  // Replace with default icon when video load fails
  const fileIcon = document.querySelector(
    `[data-file-id="${file.id}"] .file-icon`
  );
  if (fileIcon) {
    fileIcon.innerHTML = `
      <q-icon name="movie" size="lg" color="deep-orange" />
    `;
  }
};

// Handle video thumbnail error
const handleVideoThumbnailError = (file: FileItem): void => {
  // Replace with default icon when thumbnail load fails
  file.thumbnail = null;
  const fileIcon = document.querySelector(
    `[data-file-id="${file.id}"] .file-icon`
  );
  if (fileIcon) {
    fileIcon.innerHTML = `
      <q-icon name="movie" size="lg" color="deep-orange" />
    `;
  }
};

// Handle image load completion
const handleImageLoad = (file: FileItem): void => {
  file.loading = false;
};

// Handle image load error
const handleImageError = (event: Event): void => {
  const img = event.target as HTMLImageElement;
  const fileIcon = img.parentElement;

  if (fileIcon) {
    // Remove loading state
    fileIcon.classList.remove('loading');

    // Hide image and replace with icon
    img.style.display = 'none';

    // Add icon if not exists
    if (!fileIcon.querySelector('.material-icons')) {
      const icon = document.createElement('i');
      icon.className = 'material-icons';
      icon.textContent = 'image';
      icon.style.color = '#1976d2';
      icon.style.fontSize = '24px';
      icon.style.position = 'absolute';
      icon.style.top = '50%';
      icon.style.left = '50%';
      icon.style.transform = 'translate(-50%, -50%)';
      fileIcon.appendChild(icon);
    }
  }
};

// Get file icon
const getFileIcon = (contentType: string): string => {
  if (contentType.startsWith('image/')) return 'image';
  if (contentType.startsWith('video/')) return 'movie';
  if (contentType.startsWith('audio/')) return 'audiotrack';
  if (contentType.includes('pdf')) return 'picture_as_pdf';
  if (contentType.includes('word') || contentType.includes('document'))
    return 'description';
  if (contentType.includes('excel') || contentType.includes('spreadsheet'))
    return 'table_chart';
  if (
    contentType.includes('powerpoint') ||
    contentType.includes('presentation')
  )
    return 'slideshow';
  if (
    contentType.includes('zip') ||
    contentType.includes('rar') ||
    contentType.includes('archive')
  )
    return 'archive';
  return 'insert_drive_file';
};

// Get file color
const getFileColor = (contentType: string): string => {
  if (contentType.startsWith('image/')) return 'primary';
  if (contentType.startsWith('video/')) return 'deep-orange';
  if (contentType.startsWith('audio/')) return 'purple';
  if (contentType.includes('pdf')) return 'red';
  if (contentType.includes('word') || contentType.includes('document'))
    return 'blue';
  if (contentType.includes('excel') || contentType.includes('spreadsheet'))
    return 'green';
  if (
    contentType.includes('powerpoint') ||
    contentType.includes('presentation')
  )
    return 'orange';
  if (
    contentType.includes('zip') ||
    contentType.includes('rar') ||
    contentType.includes('archive')
  )
    return 'brown';
  return 'grey';
};

// Format file size
const formatFileSize = (size: number): string => {
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
  if (size < 1024 * 1024 * 1024)
    return `${(size / (1024 * 1024)).toFixed(1)} MB`;
  return `${(size / (1024 * 1024 * 1024)).toFixed(1)} GB`;
};

// Format date
const formatDate = (dateString: string): string => {
  if (!dateString) return '';
  try {
    return formatDistanceToNow(new Date(dateString), { addSuffix: true });
  } catch {
    return new Date(dateString).toLocaleDateString('en-US');
  }
};

// Truncate file name
const truncateFileName = (fileName: string): string => {
  if (fileName.length <= 25) return fileName;
  const extension = fileName.split('.').pop();
  const name = fileName.substring(0, fileName.lastIndexOf('.'));
  return `${name.substring(0, 20)}...${extension ? '.' + extension : ''}`;
};

// Download file
const downloadFile = async (file: FileItem): Promise<void> => {
  try {
    const fileUrl = cdnAddr + '/' + file.s3Path;

    const response = await fetch(fileUrl);
    const blob = await response.blob();

    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = file.name;
    link.style.display = 'none';

    document.body.appendChild(link);
    link.click();

    setTimeout(() => {
      document.body.removeChild(link);
      URL.revokeObjectURL(link.href);
    }, 100);
  } catch (error) {
    console.error('Download failed:', error);
    window.open(cdnAddr + '/' + file.s3Path, '_blank');
  }
};
</script>

<style scoped>
.files-section {
  margin-top: 16px;
}

.files-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.files-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.file-item {
  display: flex;
  align-items: center;
  padding: 12px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  background-color: #fafafa;
  transition: all 0.3s ease;
  cursor: pointer;
}

.file-item:hover {
  background-color: #f5f5f5;
  border-color: #1976d2;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.file-icon {
  margin-right: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background-color: #f0f0f0;
  border-radius: 8px;
  overflow: hidden;
  position: relative;
}

.file-thumbnail {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 8px;
  transition: transform 0.3s ease;
}

.file-thumbnail:hover {
  transform: scale(1.05);
}

/* Video thumbnail container */
.video-thumbnail-container {
  position: relative;
  width: 100%;
  height: 100%;
  border-radius: 8px;
  overflow: hidden;
}

.video-thumbnail {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 8px;
  transition: transform 0.3s ease;
}

.video-thumbnail:hover {
  transform: scale(1.05);
}

/* Video play button overlay */
.video-play-overlay {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 24px;
  height: 24px;
  background-color: rgba(0, 0, 0, 0.7);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
  transition: all 0.3s ease;
}

.video-thumbnail-container:hover .video-play-overlay {
  background-color: rgba(0, 0, 0, 0.8);
  transform: translate(-50%, -50%) scale(1.1);
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #666;
}

.file-size {
  font-weight: 500;
}

.file-date {
  color: #999;
}

.file-actions {
  margin-left: 12px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.file-item:hover .file-actions {
  opacity: 1;
}

/* Lightbox Styles */
.lightbox-card {
  background-color: rgba(0, 0, 0, 0.9);
  color: white;
}

.lightbox-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background-color: rgba(0, 0, 0, 0.8);
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
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 120px);
  padding: 0;
  position: relative;
}

.media-container {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  max-width: 90vw;
  max-height: 80vh;
}

.media-content {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

/* Navigation Buttons */
.lightbox-navigation {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 100%;
  display: flex;
  justify-content: space-between;
  padding: 0 24px;
  pointer-events: none;
}

.nav-btn {
  pointer-events: auto;
  background-color: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
}

.nav-btn:hover {
  background-color: rgba(0, 0, 0, 0.7);
  transform: scale(1.1);
}

/* File Counter */
.file-counter {
  position: absolute;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  background-color: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 500;
}

/* Loading spinner for images */
.file-icon.loading::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 20px;
  height: 20px;
  margin: -10px 0 0 -10px;
  border: 2px solid #f3f3f3;
  border-top: 2px solid #1976d2;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  opacity: 1;
  z-index: 1;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

/* Responsive design */
@media (max-width: 600px) {
  .file-item {
    padding: 8px;
  }

  .file-icon {
    width: 40px;
    height: 40px;
    margin-right: 8px;
  }

  .file-meta {
    flex-direction: column;
    gap: 2px;
  }

  .file-actions {
    opacity: 1;
  }

  .lightbox-navigation {
    padding: 0 12px;
  }

  .nav-btn {
    width: 40px;
    height: 40px;
  }
}
</style>
