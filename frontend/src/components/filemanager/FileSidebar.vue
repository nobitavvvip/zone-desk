<template>
  <div class="file-sidebar" :style="{ width: sidebarWidth + 'px' }">
    <div class="sidebar-header">
      <span class="sidebar-title">快速访问</span>
    </div>
    <div class="sidebar-content">
      <div
        v-for="shortcut in shortcuts"
        :key="shortcut.id"
        class="sidebar-item"
        :class="{ active: currentPath === shortcut.path }"
        @click="$emit('navigate', shortcut.path)"
        @contextmenu.prevent="showContextMenu($event, shortcut)"
      >
        <span class="item-icon">📁</span>
        <span class="item-name">{{ shortcut.name }}</span>
      </div>
    </div>
    <div class="sidebar-footer">
      <button class="add-btn" @click="$emit('addShortcut')">
        <span class="add-icon">+</span>
        <span class="add-text">添加快捷目录</span>
      </button>
    </div>
    <div class="resize-handle" @mousedown="startResize" />

    <teleport to="body">
      <div
        v-if="contextMenu.visible"
        class="context-overlay"
        @click="closeContextMenu"
        @contextmenu.prevent="closeContextMenu"
      >
        <div
          class="context-menu"
          :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
          @click.stop
        >
          <div class="menu-item danger" @click="handleDelete">
            删除快捷目录
          </div>
        </div>
      </div>
    </teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { listShortcuts, type Shortcut } from '@/api/shortcut'

const props = defineProps<{
  currentPath: string
}>()

const emit = defineEmits<{
  navigate: [path: string]
  addShortcut: []
  deleteShortcut: [id: string]
}>()

const shortcuts = ref<Shortcut[]>([])
const sidebarWidth = ref(200)
const contextMenu = ref<{ visible: boolean; x: number; y: number; shortcut: Shortcut | null }>({
  visible: false, x: 0, y: 0, shortcut: null,
})

onMounted(async () => {
  await loadShortcuts()
})

async function loadShortcuts() {
  try {
    shortcuts.value = await listShortcuts()
  } catch {
    shortcuts.value = []
  }
}

function showContextMenu(event: MouseEvent, shortcut: Shortcut) {
  contextMenu.value = { visible: true, x: event.clientX, y: event.clientY, shortcut }
}

function closeContextMenu() {
  contextMenu.value = { visible: false, x: 0, y: 0, shortcut: null }
}

function handleDelete() {
  if (contextMenu.value.shortcut) {
    emit('deleteShortcut', contextMenu.value.shortcut.id)
  }
  closeContextMenu()
}

function startResize(event: MouseEvent) {
  event.preventDefault()
  const startX = event.clientX
  const startWidth = sidebarWidth.value

  function onMouseMove(e: MouseEvent) {
    const newWidth = startWidth + (e.clientX - startX)
    sidebarWidth.value = Math.max(120, Math.min(400, newWidth))
  }

  function onMouseUp() {
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
  }

  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
}
</script>

<style scoped>
.file-sidebar {
  min-width: 120px;
  max-width: 400px;
  background: var(--bg-sidebar);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  position: relative;
}

.resize-handle {
  position: absolute;
  right: -3px;
  top: 0;
  bottom: 0;
  width: 6px;
  cursor: col-resize;
  z-index: 10;
}

.resize-handle:hover {
  background: var(--accent);
  opacity: 0.3;
}

.sidebar-header {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
}

.sidebar-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
}

.sidebar-content {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.sidebar-content::-webkit-scrollbar {
  width: 6px;
}

.sidebar-content::-webkit-scrollbar-track {
  background: transparent;
}

.sidebar-content::-webkit-scrollbar-thumb {
  background: var(--text-muted);
  border-radius: 3px;
  border: 2px solid transparent;
  background-clip: content-box;
  opacity: 0.4;
}

.sidebar-content::-webkit-scrollbar-thumb:hover {
  opacity: 0.7;
}

.sidebar-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  cursor: pointer;
  transition: background 0.15s ease;
}

.sidebar-item:hover {
  background: var(--bg-hover);
}

.sidebar-item.active {
  background: var(--bg-selected);
}

.item-icon {
  font-size: 16px;
}

.item-name {
  flex: 1;
  font-size: 13px;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.sidebar-footer {
  padding: 8px 16px;
  border-top: 1px solid var(--border-color);
}

.add-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 8px 12px;
  background: transparent;
  border: 1px dashed var(--border-color);
  border-radius: 4px;
  color: var(--text-secondary);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.add-btn:hover {
  background: var(--bg-hover);
  border-color: var(--accent);
  color: var(--accent);
}

.add-icon {
  font-size: 16px;
}
</style>

<style>
.context-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
}

.context-menu {
  position: fixed;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 4px;
  min-width: 140px;
  box-shadow: var(--shadow);
  z-index: 10000;
}

.menu-item {
  padding: 8px 12px;
  font-size: 13px;
  color: var(--text-primary);
  cursor: pointer;
  border-radius: 4px;
  transition: background 0.1s ease;
}

.menu-item:hover {
  background: var(--bg-hover);
}

.menu-item.danger {
  color: var(--error);
}

.menu-item.danger:hover {
  background: var(--hover-bg);
}
</style>
