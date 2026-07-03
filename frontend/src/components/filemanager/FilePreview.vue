<template>
  <div class="preview-overlay" @click.self="$emit('close')">
    <div class="preview-card">
      <div class="preview-header">
        <span class="preview-title">{{ file.name }}</span>
        <button class="preview-close" @click="$emit('close')">&times;</button>
      </div>
      <div class="preview-body">
        <img
          v-if="isImage"
          :src="imageUrl"
          class="preview-image"
          alt="preview"
        />
        <pre v-else-if="textContent !== null" class="preview-text">{{ textContent }}</pre>
        <div v-else class="preview-loading">加载中...</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { FileItem } from '@/api/file'
import { readFile, getDownloadUrl } from '@/api/file'

const props = defineProps<{
  file: FileItem
}>()

defineEmits<{
  close: []
}>()

const textContent = ref<string | null>(null)
const imageUrl = computed(() => getDownloadUrl(props.file.path))

const isImage = computed(() => {
  const ext = props.file.name.toLowerCase().split('.').pop()
  return ['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg'].includes(ext || '')
})

onMounted(async () => {
  if (!isImage.value) {
    try {
      const result = await readFile(props.file.path)
      textContent.value = result.content
    } catch {
      textContent.value = '无法预览此文件'
    }
  }
})
</script>

<style scoped>
.preview-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.preview-card {
  background: var(--bg-primary);
  border: 1px solid var(--border-input);
  border-radius: 12px;
  width: 80vw;
  height: 80vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
}

.preview-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  border-bottom: 1px solid var(--border-color);
}

.preview-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.preview-close {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 20px;
  cursor: pointer;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-close:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.preview-body {
  flex: 1;
  overflow: auto;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}

.preview-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.preview-text {
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 12px;
  font-family: 'Cascadia Code', 'Fira Code', monospace;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-secondary);
  white-space: pre-wrap;
  word-break: break-all;
  overflow: auto;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 8px;
}

.preview-loading {
  color: var(--text-muted);
  font-size: 14px;
}
</style>
