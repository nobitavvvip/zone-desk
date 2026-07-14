<template>
  <div class="file-manager-window" :style="windowStyle" @keydown="handleKeyDown" tabindex="0">
    <div class="window-header" @mousedown="startDrag" @dblclick="desktopStore.toggleMaximize(windowName)">
      <span class="window-title">文件管理器</span>
      <div class="window-controls">
        <button class="control-btn minimize" @click.stop="desktopStore.minimizeWindow(windowName)">─</button>
        <button class="control-btn maximize" @click.stop="desktopStore.toggleMaximize(windowName)">□</button>
        <button class="control-btn close" @click="desktopStore.closeWindow(windowName)">✕</button>
      </div>
    </div>
    <div class="window-body">
      <FileToolbar
        :selected-count="store.selectedItems.size"
        :has-clipboard="!!store.clipboard"
        :view-mode="store.viewMode"
        @new-folder="handleNewFolder"
        @copy="handleCopy"
        @cut="handleCut"
        @paste="handlePaste"
        @delete="handleDelete"
        @refresh="store.refresh"
        @properties="handlePropertiesFromToolbar"
        @set-view-mode="store.setViewMode"
      />
      <div class="window-content">
        <FileSidebar
          :key="sidebarKey"
          :current-path="store.currentPath"
          @navigate="handleNavigate"
          @add-shortcut="handleAddShortcut"
          @delete-shortcut="handleDeleteShortcut"
        />
        <div class="main-area">
          <div class="address-bar">
            <button class="nav-btn" @click="store.goBack" :disabled="store.historyIndex === 0" title="后退">
              ←
            </button>
            <button class="nav-btn" @click="store.goForward" :disabled="store.historyIndex >= store.history.length - 1" title="前进">
              →
            </button>
            <button class="nav-btn" @click="store.goUp" :disabled="!store.currentDir?.parent" title="上级">
              ↑
            </button>
            <input
              v-model="addressPath"
              class="address-input"
              placeholder="输入路径后按 Enter"
              @keydown.enter="handleAddressEnter"
            />
          </div>
          <FileBreadcrumb :path="store.currentPath" @navigate="handleNavigate" />
          <FileList
            :items="store.sortedItems"
            :selected-items="store.selectedItems"
            :view-mode="store.viewMode"
            :sort-by="store.sortBy"
            :sort-order="store.sortOrder"
            :has-clipboard="!!store.clipboard"
            :loading="store.loading"
            :current-path="store.currentPath"
            @select="handleSelect"
            @open="handleOpen"
            @copy="handleCopy"
            @cut="handleCut"
            @paste="handlePaste"
            @delete="handleDelete"
            @rename="handleRename"
            @properties="handleProperties"
            @sort="store.setSortBy"
            @upload="handleUpload"
            @refresh="store.refresh"
            @preview="handlePreview"
            @new-folder="handleNewFolder"
            @copy-path="handleCopyPath"
          />
        </div>
      </div>
    </div>

    <FilePreview
      v-if="previewFile"
      :file="previewFile"
      @close="previewFile = null"
    />

    <FileProperties
      v-if="propertiesFile"
      :file="propertiesFile"
      @close="propertiesFile = null"
    />

    <FileRenameDialog
      v-if="renameFile"
      :file="renameFile"
      @close="renameFile = null"
      @confirm="handleRenameConfirm"
    />

    <FileNewFolderDialog
      v-if="showNewFolderDialog"
      @close="showNewFolderDialog = false"
      @confirm="handleNewFolderConfirm"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useDesktopStore } from '@/store/desktop'
import { useFileManagerStore } from '@/store/fileManager'
import { useMessage, useDialog } from 'naive-ui'
import {
  deleteFile as apiDelete,
  renameFile as apiRename,
  mkdir as apiMkdir,
  uploadFile as apiUpload,
  copyFile as apiCopy,
  moveFile as apiMove
} from '@/api/file'
import { addShortcut as apiAddShortcut, deleteShortcut as apiDeleteShortcut } from '@/api/shortcut'
import type { FileItem } from '@/api/file'
import FileToolbar from '@/components/filemanager/FileToolbar.vue'
import FileBreadcrumb from '@/components/filemanager/FileBreadcrumb.vue'
import FileSidebar from '@/components/filemanager/FileSidebar.vue'
import FileList from '@/components/filemanager/FileList.vue'
import FilePreview from '@/components/filemanager/FilePreview.vue'
import FileProperties from '@/components/filemanager/FileProperties.vue'
import FileRenameDialog from '@/components/filemanager/FileRenameDialog.vue'
import FileNewFolderDialog from '@/components/filemanager/FileNewFolderDialog.vue'

const desktopStore = useDesktopStore()
const store = useFileManagerStore()
const message = useMessage()
const dialog = useDialog()

const windowName = 'filemanager'
const isDragging = ref(false)
const dragOffset = ref({ x: 0, y: 0 })

const sidebarKey = ref(0)
const addressPath = ref('/')
const previewFile = ref<FileItem | null>(null)
const propertiesFile = ref<FileItem | null>(null)
const renameFile = ref<FileItem | null>(null)
const showNewFolderDialog = ref(false)

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

function handlePreview(file: FileItem) {
  previewFile.value = file
}

function handleHashChange() {
  const path = store.getPathFromUrl()
  if (path !== store.currentPath) {
    store.navigate(path)
  }
}

onMounted(async () => {
  const initialPath = store.getPathFromUrl()
  await store.goTo(initialPath)
  addressPath.value = store.currentPath
  window.addEventListener('hashchange', handleHashChange)
})

onUnmounted(() => {
  window.removeEventListener('hashchange', handleHashChange)
  document.removeEventListener('mousemove', onDrag)
  document.removeEventListener('mouseup', stopDrag)
})

watch(() => store.currentPath, (path) => {
  addressPath.value = path
})

function handleKeyDown(event: KeyboardEvent) {
  if ((event.ctrlKey || event.metaKey) && event.key === 'c') {
    event.preventDefault()
    handleCopy()
  }
  else if ((event.ctrlKey || event.metaKey) && event.key === 'x') {
    event.preventDefault()
    handleCut()
  }
  else if ((event.ctrlKey || event.metaKey) && event.key === 'v') {
    event.preventDefault()
    handlePaste()
  }
  else if ((event.ctrlKey || event.metaKey) && event.key === 'a') {
    const target = event.target as HTMLElement
    if (target.tagName !== 'INPUT' && target.tagName !== 'TEXTAREA') {
      event.preventDefault()
      store.selectAll()
    }
  }
  else if (event.key === 'Delete') {
    event.preventDefault()
    handleDelete()
  }
  else if (event.key === 'Backspace') {
    const target = event.target as HTMLElement
    if (target.tagName !== 'INPUT' && target.tagName !== 'TEXTAREA') {
      event.preventDefault()
      store.goBack()
    }
  }
  else if (event.key === 'F2') {
    event.preventDefault()
    if (store.selectedItems.size === 1) {
      const selectedPath = Array.from(store.selectedItems)[0]
      const item = store.sortedItems.find(i => i.path === selectedPath)
      if (item) {
        handleRename(item)
      }
    }
  }
  else if (event.key === 'F5') {
    event.preventDefault()
    store.refresh()
  }
  else if (event.key === 'Enter') {
    event.preventDefault()
    if (store.selectedItems.size === 1) {
      const selectedPath = Array.from(store.selectedItems)[0]
      const item = store.sortedItems.find(i => i.path === selectedPath)
      if (item) {
        handleOpen(item)
      }
    }
  }
}

async function handleNavigate(path: string) {
  await store.goTo(path)
}

async function handleAddressEnter() {
  await handleNavigate(addressPath.value)
}

function handleSelect(path: string, ctrl: boolean, shift: boolean) {
  store.selectItem(path, ctrl, shift)
}

function handleOpen(item: FileItem) {
  if (item.type === 'dir') {
    store.goTo(item.path)
  } else {
    previewFile.value = item
  }
}

function handleCopy() {
  store.copy()
  message.info('已复制到剪贴板')
}

function handleCut() {
  store.cut()
  message.info('已剪切到剪贴板')
}

async function handlePaste() {
  if (!store.clipboard) return

  const { files, action } = store.clipboard
  const destination = store.currentPath

  try {
    for (const file of files) {
      const destPath = destination + '/' + file.name
      if (action === 'copy') {
        await apiCopy(file.path, destPath)
      } else {
        await apiMove(file.path, destPath)
      }
    }

    store.clipboard = null
    message.success('粘贴成功')
    await store.refresh()
  } catch (e: any) {
    message.error(e.message)
  }
}

function handleDelete() {
  if (store.selectedItems.size === 0) return

  const paths = Array.from(store.selectedItems)
  const count = paths.length

  dialog.warning({
    title: '确认删除',
    content: `确定要删除 ${count} 个项目吗？此操作不可恢复。`,
    positiveText: '删除',
    negativeText: '取消',
    positiveButtonProps: { type: 'error' },
    onPositiveClick: async () => {
      try {
        for (const path of paths) {
          await apiDelete(path)
        }
        message.success('删除成功')
        await store.refresh()
      } catch (e: any) {
        message.error(e.message)
      }
    },
  })
}

function handleRename(item: FileItem) {
  renameFile.value = item
}

async function handleRenameConfirm(newPath: string) {
  if (!renameFile.value) return

  try {
    await apiRename(renameFile.value.path, newPath)
    message.success('重命名成功')
    renameFile.value = null
    await store.refresh()
  } catch (e: any) {
    message.error(e.message)
  }
}

function handleProperties(item: FileItem) {
  propertiesFile.value = item
}

function handlePropertiesFromToolbar() {
  if (store.selectedItems.size === 1) {
    const selectedPath = Array.from(store.selectedItems)[0]
    const item = store.sortedItems.find(i => i.path === selectedPath)
    if (item) {
      propertiesFile.value = item
    }
  }
}

function handleNewFolder() {
  showNewFolderDialog.value = true
}

async function handleNewFolderConfirm(name: string) {
  const newPath = store.currentPath.replace(/\/$/, '') + '/' + name
  try {
    await apiMkdir(newPath)
    message.success('新建成功')
    showNewFolderDialog.value = false
    await store.refresh()
  } catch (e: any) {
    message.error(e.message)
  }
}

async function handleUpload(files: FileList) {
  const fileArray = Array.from(files)
  let successCount = 0
  let failCount = 0

  await Promise.all(fileArray.map(async (file) => {
    try {
      await apiUpload(store.currentPath, file)
      successCount++
    } catch (e: any) {
      message.error(`${file.name}: ${e.message}`)
      failCount++
    }
  }))

  if (successCount > 0) {
    message.success(`成功上传 ${successCount} 个文件`)
  }
  await store.refresh()
}

function handleCopyPath(path: string) {
  navigator.clipboard.writeText(path)
  message.success('路径已复制')
}

async function handleAddShortcut() {
  const name = prompt('快捷目录名称：', store.currentPath)
  if (!name) return
  try {
    await apiAddShortcut(name, store.currentPath)
    message.success('已添加')
    sidebarKey.value++
  } catch (e: any) {
    message.error(e.message)
  }
}

async function handleDeleteShortcut(id: string) {
  try {
    await apiDeleteShortcut(id)
    sidebarKey.value++
  } catch (e: any) {
    message.error(e.message)
  }
}
</script>

<style scoped>
.file-manager-window {
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
  padding: 8px 12px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
}

.window-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
}

.window-controls {
  display: flex;
  gap: 2px;
}

.control-btn {
  width: 32px;
  height: 24px;
  background: transparent;
  border: none;
  color: var(--text-primary);
  font-size: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 2px;
  transition: background 0.1s ease;
}

.control-btn:hover {
  background: var(--bg-hover);
}

.control-btn.close:hover {
  background: var(--error);
  color: white;
}

.window-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.window-content {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.main-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.address-bar {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
}

.nav-btn {
  width: 28px;
  height: 28px;
  background: transparent;
  border: 1px solid transparent;
  border-radius: 4px;
  color: var(--text-primary);
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.1s ease;
}

.nav-btn:hover:not(:disabled) {
  background: var(--bg-hover);
  border-color: var(--border-color);
}

.nav-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.address-input {
  flex: 1;
  padding: 6px 12px;
  background: var(--bg-primary);
  border: 1px solid var(--border-input);
  border-radius: 4px;
  color: var(--text-primary);
  font-size: 13px;
  outline: none;
  transition: border-color 0.15s ease;
}

.address-input:focus {
  border-color: var(--accent);
}

.address-input::placeholder {
  color: var(--text-muted);
}
</style>
