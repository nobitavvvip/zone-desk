<template>
  <div
    class="file-list"
    :class="[`view-${viewMode}`]"
    @click="handleBackgroundClick"
    @contextmenu.prevent="handleBackgroundContextMenu"
    @dragenter.prevent="handleDragEnter"
    @dragover.prevent="handleDragOver"
    @dragleave="handleDragLeave"
    @drop.prevent="handleDrop"
  >
    <!-- 详细信息视图 -->
    <div v-if="viewMode === 'details'" class="details-view">
      <div class="details-header">
        <div class="col-name" @click="handleSort('name')">
          名称
          <span v-if="sortBy === 'name'" class="sort-arrow">
            {{ sortOrder === 'asc' ? '↑' : '↓' }}
          </span>
        </div>
        <div class="col-type" @click="handleSort('type')">
          类型
          <span v-if="sortBy === 'type'" class="sort-arrow">
            {{ sortOrder === 'asc' ? '↑' : '↓' }}
          </span>
        </div>
        <div class="col-size" @click="handleSort('size')">
          大小
          <span v-if="sortBy === 'size'" class="sort-arrow">
            {{ sortOrder === 'asc' ? '↑' : '↓' }}
          </span>
        </div>
        <div class="col-modified" @click="handleSort('modified')">
          修改时间
          <span v-if="sortBy === 'modified'" class="sort-arrow">
            {{ sortOrder === 'asc' ? '↑' : '↓' }}
          </span>
        </div>
      </div>
      <div class="details-body">
        <transition name="fade" mode="out-in">
          <div v-if="loading" class="loading-state" :key="'loading'">
            <n-spin size="large" />
          </div>
          <div v-else-if="items.length === 0" class="empty-state" :key="'empty'">
            <span class="empty-icon">📂</span>
            <span class="empty-text">此文件夹为空</span>
          </div>
          <div v-else :key="'list'">
            <div
          v-for="item in items"
          :key="item.path"
          class="details-row"
          :class="{ selected: selectedItems.has(item.path) }"
          @click.stop="handleClick($event, item)"
          @dblclick.stop="handleDoubleClick(item)"
          @contextmenu.prevent.stop="handleContextMenu($event, item)"
        >
          <div class="col-name">
            <span class="file-icon">
              <svg v-if="item.type === 'dir'" class="folder" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M2 6a2 2 0 012-2h5l2 2h9a2 2 0 012 2v10a2 2 0 01-2 2H4a2 2 0 01-2-2V6z"/>
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/>
              </svg>
            </span>
            <span class="file-name">{{ item.name }}</span>
          </div>
          <div class="col-type">{{ getTypeLabel(item.type) }}</div>
          <div class="col-size">{{ item.type === 'dir' ? '-' : formatSize(item.size) }}</div>
          <div class="col-modified">{{ formatDate(item.modTime) }}</div>
        </div>
          </div>
        </transition>
      </div>
    </div>

    <!-- 列表视图 -->
    <div v-else-if="viewMode === 'list'" class="list-view">
      <transition name="fade" mode="out-in">
        <div v-if="loading" class="loading-state" :key="'loading'">
          <n-spin size="large" />
        </div>
        <div v-else-if="items.length === 0" class="empty-state" :key="'empty'">
          <span class="empty-icon">📂</span>
          <span class="empty-text">此文件夹为空</span>
        </div>
        <div v-else :key="'list'">
          <div
        v-for="item in items"
        :key="item.path"
        class="list-item"
        :class="{ selected: selectedItems.has(item.path) }"
        @click.stop="handleClick($event, item)"
        @dblclick.stop="handleDoubleClick(item)"
        @contextmenu.prevent.stop="handleContextMenu($event, item)"
      >
        <span class="file-icon">
          <svg v-if="item.type === 'dir'" class="folder" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M2 6a2 2 0 012-2h5l2 2h9a2 2 0 012 2v10a2 2 0 01-2 2H4a2 2 0 01-2-2V6z"/>
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/>
          </svg>
        </span>
        <span class="file-name">{{ item.name }}</span>
          </div>
        </div>
      </transition>
    </div>

    <!-- 大图标视图 -->
    <div v-else-if="viewMode === 'icons'" class="icons-view">
      <transition name="fade" mode="out-in">
        <div v-if="loading" class="loading-state" :key="'loading'">
          <n-spin size="large" />
        </div>
        <div v-else-if="items.length === 0" class="empty-state" :key="'empty'">
          <span class="empty-icon">📂</span>
          <span class="empty-text">此文件夹为空</span>
        </div>
        <div v-else :key="'list'" class="icons-grid">
          <div
        v-for="item in items"
        :key="item.path"
        class="icon-item"
        :class="{ selected: selectedItems.has(item.path) }"
        @click.stop="handleClick($event, item)"
        @dblclick.stop="handleDoubleClick(item)"
        @contextmenu.prevent.stop="handleContextMenu($event, item)"
      >
        <span class="file-icon-large">
          <svg v-if="item.type === 'dir'" class="folder" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M2 6a2 2 0 012-2h5l2 2h9a2 2 0 012 2v10a2 2 0 01-2 2H4a2 2 0 01-2-2V6z"/>
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/>
          </svg>
        </span>
        <span class="file-name">{{ item.name }}</span>
          </div>
        </div>
      </transition>
    </div>

    <!-- 平铺视图 -->
    <div v-else-if="viewMode === 'tile'" class="tile-view">
      <transition name="fade" mode="out-in">
        <div v-if="loading" class="loading-state" :key="'loading'">
          <n-spin size="large" />
        </div>
        <div v-else-if="items.length === 0" class="empty-state" :key="'empty'">
          <span class="empty-icon">📂</span>
          <span class="empty-text">此文件夹为空</span>
        </div>
        <div v-else :key="'list'" class="tile-grid">
          <div
        v-for="item in items"
        :key="item.path"
        class="tile-item"
        :class="{ selected: selectedItems.has(item.path) }"
        @click.stop="handleClick($event, item)"
        @dblclick.stop="handleDoubleClick(item)"
        @contextmenu.prevent.stop="handleContextMenu($event, item)"
      >
        <span class="file-icon-large">
          <svg v-if="item.type === 'dir'" class="folder" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M2 6a2 2 0 012-2h5l2 2h9a2 2 0 012 2v10a2 2 0 01-2 2H4a2 2 0 01-2-2V6z"/>
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/>
          </svg>
        </span>
        <div class="tile-info">
          <span class="file-name">{{ item.name }}</span>
          <span class="file-meta">{{ item.type === 'dir' ? '文件夹' : formatSize(item.size) }}</span>
        </div>
          </div>
        </div>
      </transition>
    </div>

    <!-- 拖拽上传提示 -->
    <div v-if="isDragging" class="drag-overlay">
      <div class="drag-content">
        <span class="drag-icon">📁</span>
        <span class="drag-text">拖拽文件到此处上传</span>
      </div>
    </div>

    <!-- 右键菜单 -->
    <teleport to="body">
      <div
        v-if="contextMenu.visible"
        class="context-menu-overlay"
        @click="closeContextMenu"
        @contextmenu.prevent="closeContextMenu"
      >
        <div
          class="context-menu"
          :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
          @click.stop
        >
          <template v-if="contextMenu.item">
            <div class="menu-item" @click="handleMenuOpen">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M2 6a2 2 0 012-2h5l2 2h9a2 2 0 012 2v10a2 2 0 01-2 2H4a2 2 0 01-2-2V6z"/></svg>
              打开 <span class="shortcut">Enter</span>
            </div>
            <div v-if="contextMenu.item.type !== 'dir'" class="menu-item" @click="handleMenuPreview">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
              预览
            </div>
            <div v-if="contextMenu.item.type !== 'dir'" class="menu-item" @click="handleMenuDownload">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
              下载
            </div>
            <div class="menu-divider"></div>
            <div class="menu-item" @click="handleMenuCopy">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/></svg>
              复制 <span class="shortcut">Ctrl+C</span>
            </div>
            <div class="menu-item" @click="handleMenuCut">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="6" cy="6" r="3"/><circle cx="6" cy="18" r="3"/><line x1="20" y1="4" x2="8.12" y2="15.88"/><line x1="14.47" y1="14.48" x2="20" y2="20"/><line x1="8.12" y1="8.12" x2="12" y2="12"/></svg>
              剪切 <span class="shortcut">Ctrl+X</span>
            </div>
            <div class="menu-divider"></div>
            <div class="menu-item" @click="handleMenuRename">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
              重命名 <span class="shortcut">F2</span>
            </div>
            <div class="menu-item danger" @click="handleMenuDelete">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
              删除 <span class="shortcut">Delete</span>
            </div>
            <div class="menu-divider"></div>
            <div class="menu-item" @click="handleMenuProperties">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z"/></svg>
              属性
            </div>
          </template>
          <template v-else>
            <div class="menu-item" @click="handleMenuNewFolder">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M2 6a2 2 0 012-2h5l2 2h9a2 2 0 012 2v10a2 2 0 01-2 2H4a2 2 0 01-2-2V6z"/><path d="M12 10v6m-3-3h6"/></svg>
              新建目录 <span class="shortcut">Ctrl+N</span>
            </div>
            <div class="menu-item" @click="handleMenuUpload">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
              上传文件
            </div>
            <div class="menu-divider"></div>
            <div class="menu-item" :class="{ disabled: !hasClipboard }" @click="handleMenuPaste">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M16 4h2a2 2 0 012 2v14a2 2 0 01-2 2H6a2 2 0 01-2-2V6a2 2 0 012-2h2"/><rect x="8" y="2" width="8" height="4" rx="1" ry="1"/></svg>
              粘贴 <span class="shortcut">Ctrl+V</span>
            </div>
            <div class="menu-divider"></div>
            <div class="menu-item" @click="handleMenuRefresh">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M23 4v6h-6M1 20v-6h6"/><path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15"/></svg>
              刷新 <span class="shortcut">F5</span>
            </div>
            <div class="menu-divider"></div>
            <div class="menu-item" @click="handleMenuCopyPath">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/></svg>
              复制当前路径 <span class="shortcut">Ctrl+C</span>
            </div>
            <div class="menu-item" @click="handleMenuDirProperties">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z"/></svg>
              属性
            </div>
          </template>
        </div>
      </div>
    </teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { NSpin } from 'naive-ui'
import type { FileItem } from '@/api/file'
import type { ViewMode, SortBy, SortOrder } from '@/store/fileManager'

const props = defineProps<{
  items: FileItem[]
  selectedItems: Set<string>
  viewMode: ViewMode
  sortBy: SortBy
  sortOrder: SortOrder
  hasClipboard: boolean
  loading: boolean
  currentPath: string
}>()

const emit = defineEmits<{
  select: [path: string, ctrl: boolean, shift: boolean]
  open: [item: FileItem]
  preview: [file: FileItem]
  copy: []
  cut: []
  paste: []
  delete: []
  rename: [item: FileItem]
  properties: [item: FileItem]
  sort: [field: SortBy]
  upload: [files: FileList]
  newFolder: []
  refresh: []
  copyPath: [path: string]
}>()

const isDragging = ref(false)
const dragCounter = ref(0)
const contextMenu = ref({
  visible: false,
  x: 0,
  y: 0,
  item: null as FileItem | null,
})

function getTypeLabel(type: string): string {
  const labels: Record<string, string> = {
    file: '文件',
    dir: '目录',
    symlink: '链接',
  }
  return labels[type] || type
}

function formatSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  let i = 0
  let size = bytes
  while (size >= 1024 && i < units.length - 1) {
    size /= 1024
    i++
  }
  return i === 0 ? `${size} ${units[i]}` : `${size.toFixed(1)} ${units[i]}`
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

function handleClick(event: MouseEvent, item: FileItem) {
  emit('select', item.path, event.ctrlKey || event.metaKey, event.shiftKey)
}

function handleDoubleClick(item: FileItem) {
  emit('open', item)
}

function handleBackgroundClick() {
  emit('select', '', false, false)
}

function handleContextMenu(event: MouseEvent, item: FileItem) {
  if (!props.selectedItems.has(item.path)) {
    emit('select', item.path, false, false)
  }
  contextMenu.value = {
    visible: true,
    x: event.clientX,
    y: event.clientY,
    item,
  }
}

function handleBackgroundContextMenu(event: MouseEvent) {
  emit('select', '', false, false)
  contextMenu.value = {
    visible: true,
    x: event.clientX,
    y: event.clientY,
    item: null,
  }
}

function handleSort(field: SortBy) {
  emit('sort', field)
}

function handleDragEnter(event: DragEvent) {
  event.preventDefault()
  if (!event.dataTransfer?.types.includes('Files')) return
  dragCounter.value++
  isDragging.value = true
}

function handleDragLeave(event: DragEvent) {
  event.preventDefault()
  if (dragCounter.value > 0) dragCounter.value--
  if (dragCounter.value === 0) {
    isDragging.value = false
  }
}

function handleDragOver(event: DragEvent) {
  event.preventDefault()
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = 'copy'
  }
}

function handleDrop(event: DragEvent) {
  event.preventDefault()
  dragCounter.value = 0
  isDragging.value = false
  if (event.dataTransfer?.files) {
    emit('upload', event.dataTransfer.files)
  }
}

function closeContextMenu() {
  contextMenu.value.visible = false
}

function handleMenuOpen() {
  if (contextMenu.value.item) {
    emit('open', contextMenu.value.item)
  }
  closeContextMenu()
}

function handleMenuPreview() {
  if (contextMenu.value.item) {
    emit('preview', contextMenu.value.item)
  }
  closeContextMenu()
}

function handleMenuDownload() {
  if (contextMenu.value.item) {
    window.open(`/api/files/download?path=${encodeURIComponent(contextMenu.value.item.path)}`, '_blank')
  }
  closeContextMenu()
}

function handleMenuCopy() {
  emit('copy')
  closeContextMenu()
}

function handleMenuCut() {
  emit('cut')
  closeContextMenu()
}

function handleMenuPaste() {
  if (props.hasClipboard) {
    emit('paste')
  }
  closeContextMenu()
}

function handleMenuRename() {
  if (contextMenu.value.item) {
    emit('rename', contextMenu.value.item)
  }
  closeContextMenu()
}

function handleMenuDelete() {
  emit('delete')
  closeContextMenu()
}

function handleMenuProperties() {
  if (contextMenu.value.item) {
    emit('properties', contextMenu.value.item)
  }
  closeContextMenu()
}

function handleMenuNewFolder() {
  emit('newFolder')
  closeContextMenu()
}

function handleMenuUpload() {
  const input = document.createElement('input')
  input.type = 'file'
  input.multiple = true
  input.onchange = (e: Event) => {
    const target = e.target as HTMLInputElement
    if (target.files) {
      emit('upload', target.files)
    }
  }
  input.click()
  closeContextMenu()
}

function handleMenuRefresh() {
  emit('refresh')
  closeContextMenu()
}

function handleMenuCopyPath() {
  if (contextMenu.value.item) {
    emit('copyPath', contextMenu.value.item.path)
  } else {
    emit('copyPath', props.currentPath)
  }
  closeContextMenu()
}

function handleMenuDirProperties() {
  if (contextMenu.value.item) {
    emit('properties', contextMenu.value.item)
  } else {
    const dirItem: FileItem = {
      name: props.currentPath.split('/').filter(Boolean).pop() || '/',
      path: props.currentPath,
      type: 'dir',
      size: 0,
      mode: '',
      modTime: '',
      hidden: false,
    }
    emit('properties', dirItem)
  }
  closeContextMenu()
}

function handleClickOutside(event: MouseEvent) {
  if (contextMenu.value.visible) {
    closeContextMenu()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.file-list {
  flex: 1;
  overflow: auto;
  background: var(--bg-primary);
  position: relative;
}

.file-list::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.file-list::-webkit-scrollbar-track {
  background: transparent;
}

.file-list::-webkit-scrollbar-thumb {
  background: var(--text-muted);
  border-radius: 4px;
  border: 2px solid transparent;
  background-clip: content-box;
  opacity: 0.4;
}

.file-list::-webkit-scrollbar-thumb:hover {
  opacity: 0.7;
}

.file-list::-webkit-scrollbar-corner {
  background: transparent;
}

/* 详细信息视图 */
.details-view {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.details-header {
  display: flex;
  padding: 8px 16px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
}

.details-header > div {
  cursor: pointer;
}

.details-header > div:hover {
  color: var(--text-primary);
}

.col-name {
  flex: 2;
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.col-type, .col-size, .col-modified {
  flex: 1;
}

.sort-arrow {
  margin-left: 4px;
  color: var(--accent);
}

.details-body {
  flex: 1;
  overflow-y: auto;
}

.details-body::-webkit-scrollbar {
  width: 8px;
}

.details-body::-webkit-scrollbar-track {
  background: transparent;
}

.details-body::-webkit-scrollbar-thumb {
  background: var(--text-muted);
  border-radius: 4px;
  border: 2px solid transparent;
  background-clip: content-box;
  opacity: 0.4;
}

.details-body::-webkit-scrollbar-thumb:hover {
  opacity: 0.7;
}

.details-row {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
  transition: background 0.1s ease;
}

.details-row:hover {
  background: var(--hover-bg-subtle);
}

.details-row.selected {
  background: rgba(64, 158, 255, 0.12);
}

/* 列表视图 */
.list-view {
  display: flex;
  flex-direction: column;
  padding: 8px;
}

.list-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: background 0.1s ease;
}

.list-item:hover {
  background: var(--hover-bg-subtle);
}

.list-item.selected {
  background: rgba(64, 158, 255, 0.12);
}

/* 大图标视图 */
.icons-view {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 16px;
}

.icons-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.icon-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 90px;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.1s ease;
  user-select: none;
}

.icon-item:hover {
  background: var(--hover-bg-subtle);
}

.icon-item.selected {
  background: rgba(64, 158, 255, 0.12);
}

.icon-item .file-name {
  max-width: 100%;
  text-align: center;
  word-break: break-all;
}

/* 平铺视图 */
.tile-view {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 16px;
}

.tile-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tile-item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 200px;
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.1s ease;
}

.tile-item:hover {
  background: var(--hover-bg-subtle);
  border-color: var(--accent);
}

.tile-item.selected {
  background: rgba(64, 158, 255, 0.12);
  border-color: var(--accent);
}

.tile-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.file-meta {
  font-size: 12px;
  color: var(--text-muted);
}

/* 通用样式 */
.file-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
}

.file-icon svg {
  width: 16px;
  height: 16px;
}

.file-icon .folder {
  color: #eab308;
}

.file-icon-large {
  font-size: 48px;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.file-icon-large svg {
  width: 44px;
  height: 44px;
  color: var(--text-muted);
}

.file-icon-large .folder {
  color: #eab308;
}

.file-name {
  font-size: 13px;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 过渡动画 */
.fade-enter-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from {
  opacity: 0;
}

/* 加载和空状态 */
.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 40px;
  color: var(--text-muted);
  font-size: 15px;
}

.empty-icon {
  font-size: 48px;
}

.empty-text {
  font-size: 15px;
  color: var(--text-muted);
}

.list-loading,
.list-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 40px;
  color: var(--text-muted);
  font-size: 15px;
}

/* 拖拽上传 */
.drag-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 120, 212, 0.1);
  border: 2px dashed var(--accent);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
}

.drag-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.drag-icon {
  font-size: 48px;
}

.drag-text {
  font-size: 16px;
  color: var(--accent);
  font-weight: 600;
}
</style>

<style>
.context-menu-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
}

.context-menu {
  position: fixed;
  background: var(--bg-primary);
  border: 1px solid var(--border-input);
  border-radius: 10px;
  padding: 4px 0;
  min-width: 200px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.35);
  z-index: 10000;
}

.context-menu .menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  border-radius: 6px;
}

.context-menu .menu-item svg {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.context-menu .menu-item:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.context-menu .menu-item.danger {
  color: #ef4444;
}

.context-menu .menu-item.danger:hover {
  background: rgba(239, 68, 68, 0.12);
}

.context-menu .menu-item.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.context-menu .menu-item.disabled:hover {
  background: transparent;
}

.context-menu .menu-divider {
  height: 1px;
  background: var(--border-color);
  margin: 2px 0;
}

.context-menu .shortcut {
  margin-left: auto;
  font-size: 12px;
  color: var(--text-muted);
  opacity: 0.6;
}
</style>
