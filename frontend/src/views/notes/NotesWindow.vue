<template>
  <div class="notes-window" :style="windowStyle" @keydown="handleKeyDown" @contextmenu.prevent.stop tabindex="0">
    <div class="window-header" @mousedown="startDrag" @dblclick="desktopStore.toggleMaximize(windowName)">
      <div class="header-left">
        <span class="window-title">📝</span>
        <div class="tabs">
          <div
            v-for="file in notesStore.openFiles"
            :key="file.path"
            class="tab"
            :class="{ active: file.path === notesStore.activeFile?.path }"
            @click="notesStore.activeFile = file"
          >
            <span class="tab-name">{{ file.name }}</span>
            <span v-if="file.modified" class="modified-dot">●</span>
            <button class="tab-close" @click.stop="notesStore.closeFile(file.path)">×</button>
          </div>
        </div>
      </div>
      <div class="header-actions" @click.stop @dblclick.stop>
        <button class="action-btn" @click="notesStore.saveCurrentFile()" :disabled="!notesStore.activeFile?.modified" title="保存 (Ctrl+S)">
          <span>💾</span>
        </button>
        <div class="action-sep"></div>
        <button class="action-btn" @click="notesStore.sidebarCollapsed = !notesStore.sidebarCollapsed" :class="{ active: !notesStore.sidebarCollapsed }" title="侧边栏">
          <span>📋</span>
        </button>
        <button class="action-btn" @click="toggleOutline" title="大纲">
          <span>📑</span>
        </button>
        <div class="action-sep"></div>
        <div class="width-control" @mouseenter="showWidthPopup" @mouseleave="hideWidthPopup">
          <button class="action-btn" :class="{ active: notesStore.editorPadding > 40 }" title="编辑区宽度">
            <span>↔</span>
          </button>
          <div v-if="showWidthSlider" class="width-popup" @mouseenter="showWidthPopup" @mouseleave="hideWidthPopup">
            <div class="width-label">边距 {{ notesStore.editorPadding }}px</div>
            <input
              type="range"
              class="width-slider"
              min="0"
              max="300"
              step="10"
              :value="notesStore.editorPadding"
              @input="notesStore.editorPadding = +($event.target as HTMLInputElement).value"
            />
            <div class="width-presets">
              <button class="preset-btn" :class="{ active: notesStore.editorPadding === 0 }" @click="notesStore.editorPadding = 0">紧凑</button>
              <button class="preset-btn" :class="{ active: notesStore.editorPadding === 40 }" @click="notesStore.editorPadding = 40">默认</button>
              <button class="preset-btn" :class="{ active: notesStore.editorPadding === 120 }" @click="notesStore.editorPadding = 120">阅读</button>
              <button class="preset-btn" :class="{ active: notesStore.editorPadding === 200 }" @click="notesStore.editorPadding = 200">沉浸</button>
            </div>
          </div>
        </div>
        <div class="action-sep"></div>
        <div class="autosave-control" @mouseenter="showAutoSavePopup" @mouseleave="hideAutoSavePopup">
          <button class="action-btn" :class="{ active: notesStore.autoSave }" title="自动保存">
            <span>⏱</span>
          </button>
          <div v-if="showAutoSaveSettings" class="autosave-popup" @mouseenter="showAutoSavePopup" @mouseleave="hideAutoSavePopup">
            <div class="autosave-header">自动保存</div>
            <label class="autosave-toggle">
              <input type="checkbox" :checked="notesStore.autoSave" @change="toggleAutoSave" />
              <span>内容变化时自动保存</span>
            </label>
          </div>
        </div>
        <div class="action-sep"></div>
        <button class="control-btn minimize" @click="desktopStore.minimizeWindow(windowName)">─</button>
        <button class="control-btn maximize" @click="desktopStore.toggleMaximize(windowName)">□</button>
        <button class="control-btn close" @click="handleClose">✕</button>
      </div>
    </div>

    <div class="workspace">
      <NoteSidebar
        v-if="!notesStore.sidebarCollapsed"
        :tree-nodes="notesStore.treeNodes"
        :loading="notesStore.loading"
        :root-dir="notesStore.rootDir"
        @toggle="handleToggle"
        @open-file="handleOpenFile"
        @create-note="handleCreateNote"
        @create-folder="handleCreateFolder"
        @delete-item="handleDeleteItem"
        @rename-item="handleRenameItem"
        @refresh-tree="handleRefreshTree"
      />

      <div class="editor-container">
        <NoteEditor
          v-if="notesStore.activeFile"
          ref="editorRef"
          :file="notesStore.activeFile"
          :padding="notesStore.editorPadding"
          @update:content="handleContentUpdate"
          @blur="handleEditorBlur"
        />
        <div v-else class="empty-state">
          <div class="empty-icon">📝</div>
          <div class="empty-text">选择或创建一个笔记开始编辑</div>
          <div class="empty-hint">支持 Markdown 语法，即时渲染模式</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useDesktopStore } from '@/store/desktop'
import { useNotesStore } from '@/store/notes'
import { useMessage } from 'naive-ui'
import NoteSidebar from '@/components/notes/NoteSidebar.vue'
import NoteEditor from '@/components/notes/NoteEditor.vue'

const desktopStore = useDesktopStore()
const notesStore = useNotesStore()
const message = useMessage()

const windowName = 'notes'
const isDragging = ref(false)
const dragOffset = ref({ x: 0, y: 0 })
const showWidthSlider = ref(false)
const showAutoSaveSettings = ref(false)
const editorRef = ref()

const windowStyle = computed(() => {
  const state = desktopStore.getWindowState(windowName)
  if (state.isMinimized) return { display: 'none' }
  return {
    left: state.x + 'px',
    top: state.y + 'px',
    width: state.width + 'px',
    height: state.height + 'px',
  }
})

function startDrag(event: MouseEvent) {
  if (desktopStore.getWindowState(windowName).isMaximized) return
  isDragging.value = true
  const state = desktopStore.getWindowState(windowName)
  dragOffset.value = {
    x: event.clientX - state.x,
    y: event.clientY - state.y,
  }
  document.addEventListener('mousemove', onDrag)
  document.addEventListener('mouseup', stopDrag)
}

function onDrag(event: MouseEvent) {
  if (!isDragging.value) return
  const newX = event.clientX - dragOffset.value.x
  const newY = event.clientY - dragOffset.value.y
  desktopStore.updatePosition(windowName, newX, newY)
}

function stopDrag() {
  isDragging.value = false
  document.removeEventListener('mousemove', onDrag)
  document.removeEventListener('mouseup', stopDrag)
}

let widthHideTimer: ReturnType<typeof setTimeout> | null = null

function showWidthPopup() {
  if (widthHideTimer) {
    clearTimeout(widthHideTimer)
    widthHideTimer = null
  }
  showWidthSlider.value = true
}

function hideWidthPopup() {
  widthHideTimer = setTimeout(() => {
    showWidthSlider.value = false
  }, 200)
}

function toggleOutline() {
  editorRef.value?.toggleOutline()
}

function handleContentUpdate(content: string) {
  if (notesStore.activeFile) {
    notesStore.updateContent(notesStore.activeFile.path, content)
  }
}

function handleEditorBlur() {
  notesStore.saveCurrentFile()
}

async function handleToggle(node: any) {
  await notesStore.toggleDirExpand(node)
}

async function handleOpenFile(node: any) {
  if (node.type === 'dir') return
  const err = await notesStore.openFile(node)
  if (err) {
    message.error(err)
  }
}

async function handleRefreshTree() {
  await notesStore.loadTreeRoot()
  message.success('已刷新')
}

async function handleCreateNote(dirPath: string, name: string) {
  const err = await notesStore.createNote(dirPath, name)
  if (err) {
    message.error(err)
  } else {
    message.success('笔记已创建')
  }
}

async function handleCreateFolder(dirPath: string, name: string) {
  const err = await notesStore.createFolder(dirPath, name)
  if (err) {
    message.error(err)
  } else {
    message.success('文件夹已创建')
  }
}

async function handleDeleteItem(path: string) {
  try {
    await notesStore.deleteItem(path)
    message.success('已删除')
  } catch (e: any) {
    message.error(e.message)
  }
}

async function handleRenameItem(oldPath: string, newName: string) {
  try {
    await notesStore.renameItem(oldPath, newName)
    message.success('重命名成功')
  } catch (e: any) {
    message.error(e.message)
  }
}

async function handleClose() {
  const modifiedCount = notesStore.openFiles.filter(f => f.modified).length
  if (modifiedCount > 0) {
    if (!confirm(`有 ${modifiedCount} 个未保存的笔记，确定关闭？`)) {
      return
    }
  }
  desktopStore.closeWindow(windowName)
}

function handleKeyDown(event: KeyboardEvent) {
  if ((event.ctrlKey || event.metaKey) && event.key === 's') {
    event.preventDefault()
    notesStore.saveCurrentFile()
  }
  if ((event.ctrlKey || event.metaKey) && event.key === 'w') {
    event.preventDefault()
    if (notesStore.activeFile) {
      notesStore.closeFile(notesStore.activeFile.path)
    }
  }
}

let autoSaveHideTimer: ReturnType<typeof setTimeout> | null = null

function showAutoSavePopup() {
  if (autoSaveHideTimer) {
    clearTimeout(autoSaveHideTimer)
    autoSaveHideTimer = null
  }
  showAutoSaveSettings.value = true
}

function hideAutoSavePopup() {
  autoSaveHideTimer = setTimeout(() => {
    showAutoSaveSettings.value = false
  }, 200)
}

function toggleAutoSave() {
  notesStore.setAutoSave(!notesStore.autoSave)
}

onMounted(async () => {
  await notesStore.init()
})

onUnmounted(() => {
  document.removeEventListener('mousemove', onDrag)
  document.removeEventListener('mouseup', stopDrag)
})
</script>

<style scoped>
.notes-window {
  position: absolute;
  left: 5vw;
  top: 5vh;
  width: 90vw;
  height: 86vh;
  border-radius: 4px;
  overflow: hidden;
  background: var(--bg-primary);
  box-shadow: var(--shadow);
  display: flex;
  flex-direction: column;
  z-index: 100;
  border: 1px solid var(--border-color);
  outline: none;
}

.window-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 8px;
  height: 40px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.window-title {
  font-size: 16px;
  flex-shrink: 0;
}

.tabs {
  display: flex;
  align-items: center;
  gap: 2px;
  flex: 1;
  min-width: 0;
  overflow-x: auto;
  scrollbar-width: none;
}

.tabs::-webkit-scrollbar {
  display: none;
}

.tab {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  font-size: 12px;
  color: var(--text-secondary);
  background: transparent;
  border-radius: 4px;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.15s ease;
  max-width: 120px;
}

.tab:hover {
  background: var(--bg-hover);
}

.tab.active {
  background: var(--bg-primary);
  color: var(--text-primary);
}

.tab-name {
  overflow: hidden;
  text-overflow: ellipsis;
}

.modified-dot {
  color: var(--accent);
  font-size: 10px;
}

.tab-close {
  width: 14px;
  height: 14px;
  background: transparent;
  border: none;
  color: var(--text-secondary);
  font-size: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 3px;
  opacity: 0;
  transition: all 0.1s ease;
  padding: 0;
}

.tab:hover .tab-close {
  opacity: 1;
}

.tab-close:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-primary);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 2px;
  flex-shrink: 0;
}

.action-btn {
  width: 28px;
  height: 28px;
  background: transparent;
  border: none;
  border-radius: 4px;
  color: var(--text-secondary);
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.action-btn:hover:not(:disabled) {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.action-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.action-btn.active {
  background: var(--bg-hover);
  color: var(--accent);
}

.action-sep {
  width: 1px;
  height: 16px;
  background: var(--border-color);
  margin: 0 4px;
}

.control-btn {
  width: 28px;
  height: 28px;
  background: transparent;
  border: none;
  color: var(--text-secondary);
  font-size: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.1s ease;
}

.control-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.control-btn.close:hover {
  background: var(--error);
  color: white;
}

.workspace {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.editor-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: var(--text-muted);
}

.empty-icon {
  font-size: 48px;
  opacity: 0.5;
}

.empty-text {
  font-size: 16px;
  font-weight: 500;
}

.empty-hint {
  font-size: 12px;
  opacity: 0.6;
}

.width-control {
  position: relative;
}

.width-popup {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  width: 220px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
  z-index: 100;
}

.width-label {
  font-size: 11px;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.width-slider {
  width: 100%;
  height: 4px;
  -webkit-appearance: none;
  appearance: none;
  background: rgba(255, 255, 255, 0.12);
  border-radius: 2px;
  outline: none;
}

.width-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: var(--accent);
  cursor: pointer;
}

.width-presets {
  display: flex;
  gap: 4px;
  margin-top: 10px;
}

.preset-btn {
  flex: 1;
  padding: 4px 0;
  font-size: 11px;
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.15s ease;
}

.preset-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.preset-btn.active {
  background: var(--accent);
  color: white;
  border-color: var(--accent);
}

.autosave-control {
  position: relative;
}

.autosave-popup {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  width: 200px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
  z-index: 100;
}

.autosave-header {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 12px;
}

.autosave-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: var(--text-secondary);
  cursor: pointer;
}

.autosave-toggle input {
  cursor: pointer;
}
</style>
