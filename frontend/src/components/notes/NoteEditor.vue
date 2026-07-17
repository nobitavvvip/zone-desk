<template>
  <div class="note-editor">
    <div ref="editorRef" class="vditor-container"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import type { OpenFile } from '@/store/notes'

interface Props {
  file: OpenFile
  padding?: number
}

const props = withDefaults(defineProps<Props>(), {
  padding: 40,
})

const emit = defineEmits<{
  (e: 'update:content', content: string): void
  (e: 'blur'): void
}>()

const editorRef = ref<HTMLDivElement>()
let editor: Vditor | null = null
let isUpdating = false

function applyPadding(value: number) {
  const container = editorRef.value
  if (!container) return
  container.style.setProperty('--editor-padding-h', value + 'px')
  container.style.setProperty('--editor-padding-v', value > 0 ? '20px' : '4px')
  container.classList.toggle('compact', value === 0)
}

function createEditor() {
  if (!editorRef.value) return

  editor = new Vditor(editorRef.value, {
    mode: 'ir',
    toolbar: [
      'emoji', 'headings', 'bold', 'italic', 'strike', '|',
      'line', 'quote', 'list', 'ordered-list', 'check', '|',
      'code', 'inline-code', 'table', '|',
      'undo', 'redo', '|',
      'outline', 'edit-preview',
    ],
    outline: {
      enable: true,
      position: 'right',
    },
    tab: '    ',
    cache: { enable: false },
    theme: 'classic',
    width: '100%',
    height: '100%',
    input: (value: string) => {
      if (!isUpdating) {
        emit('update:content', value)
      }
    },
    blur: () => {
      emit('blur')
    },
    after: () => {
      if (editor && props.file) {
        isUpdating = true
        editor.setValue(props.file.content)
        isUpdating = false
        applyPadding(props.padding)
      }
    },
  })
}

function destroyEditor() {
  if (editor) {
    editor.destroy()
    editor = null
  }
}

function toggleOutline() {
  if (!editorRef.value) return
  const btn = editorRef.value.querySelector('[data-type="outline"]') as HTMLElement
  if (btn) btn.click()
}

defineExpose({ toggleOutline })

watch(
  () => props.file?.path,
  (newPath, oldPath) => {
    if (newPath !== oldPath && editor) {
      isUpdating = true
      editor.setValue(props.file?.content || '')
      isUpdating = false
    }
  }
)

watch(
  () => props.padding,
  (val) => applyPadding(val)
)

onMounted(() => {
  createEditor()
})

onBeforeUnmount(() => {
  destroyEditor()
})
</script>

<style scoped>
.note-editor {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.vditor-container {
  flex: 1;
  overflow: hidden;
}

:deep(.vditor) {
  border: none !important;
  border-radius: 0 !important;
}

:deep(.vditor-toolbar) {
  border-bottom: 1px solid var(--border-color) !important;
  background: var(--bg-secondary) !important;
  padding: 4px 8px !important;
  min-height: 36px !important;
}

:deep(.vditor-toolbar-item) {
  background: transparent !important;
}

:deep(.vditor-toolbar-item:hover) {
  background: var(--bg-hover) !important;
}

:deep(.vditor-content) {
  background: var(--bg-primary) !important;
  padding: 0 !important;
}

:deep(.vditor-ir) {
  background: var(--bg-primary) !important;
  color: var(--text-primary) !important;
  font-size: 15px !important;
  line-height: 1.8 !important;
  padding: var(--editor-padding-v, 20px) var(--editor-padding-h, 40px) !important;
}

:deep(.vditor-ir:focus) {
  outline: none !important;
}

:deep(.vditor-ir pre.vditor-reset) {
  padding: var(--editor-padding-v, 20px) var(--editor-padding-h, 40px) !important;
}

.compact :deep(.vditor-ir p),
.compact :deep(.vditor-ir h1),
.compact :deep(.vditor-ir h2),
.compact :deep(.vditor-ir h3),
.compact :deep(.vditor-ir h4),
.compact :deep(.vditor-ir h5),
.compact :deep(.vditor-ir h6),
.compact :deep(.vditor-ir li),
.compact :deep(.vditor-ir blockquote),
.compact :deep(.vditor-ir pre),
.compact :deep(.vditor-ir table),
.compact :deep(.vditor-ir hr) {
  margin-left: 0 !important;
  margin-right: 0 !important;
}

:deep(.vditor-outline) {
  border-left: 1px solid var(--border-color) !important;
  background: var(--bg-secondary) !important;
  font-size: 12px !important;
}

:deep(.vditor-outline__item) {
  color: var(--text-secondary) !important;
}

:deep(.vditor-outline__item:hover) {
  color: var(--text-primary) !important;
  background: var(--bg-hover) !important;
}

:deep(.vditor-outline__item--current) {
  color: var(--accent) !important;
}

:deep(*::-webkit-scrollbar) {
  width: 6px !important;
  height: 6px !important;
}

:deep(*::-webkit-scrollbar-track) {
  background: transparent !important;
}

:deep(*::-webkit-scrollbar-thumb) {
  background: rgba(128, 128, 128, 0.3) !important;
  border-radius: 3px !important;
}

:deep(*::-webkit-scrollbar-thumb:hover) {
  background: rgba(128, 128, 128, 0.5) !important;
}
</style>
