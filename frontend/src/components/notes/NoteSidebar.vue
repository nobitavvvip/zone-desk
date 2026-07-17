<template>
  <div class="note-sidebar" @contextmenu.prevent="showBgContextMenu">
    <div class="sidebar-header">
      <span class="header-title">文件</span>
      <button class="refresh-btn" @click="handleRefresh" title="刷新">↻</button>
    </div>

    <div class="sidebar-content" @contextmenu.prevent="showBgContextMenu">
      <div v-if="loading && treeNodes.length === 0" class="loading-state">
        <span>加载中...</span>
      </div>

      <div v-else-if="treeNodes.length === 0" class="empty-sidebar">
        <span>暂无笔记</span>
      </div>

      <div v-else class="tree">
        <TreeNodeItem
          v-for="node in treeNodes"
          :key="node.path"
          :node="node"
          :depth="0"
          :selected-path="selectedPath"
          @select="handleSelect"
          @toggle="handleToggle"
          @open-file="handleOpenFile"
          @context-menu="showItemContextMenu"
        />
      </div>
    </div>

    <div
      v-if="contextMenu.visible"
      class="context-menu"
      :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
    >
      <template v-if="contextMenu.node">
        <template v-if="contextMenu.node.type === 'dir'">
          <div class="ctx-item" @click="handleContextAction('new-note')">📄 新建笔记</div>
          <div class="ctx-item" @click="handleContextAction('new-folder')">📁 新建文件夹</div>
          <div class="ctx-sep"></div>
        </template>
        <div class="ctx-item" @click="handleContextAction('rename')">重命名</div>
        <div class="ctx-sep"></div>
        <div class="ctx-item danger" @click="handleContextAction('delete')">删除</div>
      </template>
      <template v-else>
        <div class="ctx-item" @click="handleContextAction('new-note')">📄 新建笔记</div>
        <div class="ctx-item" @click="handleContextAction('new-folder')">📁 新建文件夹</div>
        <div class="ctx-sep"></div>
        <div class="ctx-item" @click="handleContextAction('refresh')">🔄 刷新</div>
      </template>
    </div>

    <NoteNewDialog
      v-if="showCreateDialog"
      :title="createDialogTitle"
      :placeholder="createDialogPlaceholder"
      @close="showCreateDialog = false"
      @confirm="handleCreateConfirm"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { TreeNode } from '@/store/notes'
import TreeNodeItem from './TreeNodeItem.vue'
import NoteNewDialog from './NoteNewDialog.vue'

interface Props {
  treeNodes: TreeNode[]
  loading: boolean
  rootDir?: string
}

const props = withDefaults(defineProps<Props>(), {
  rootDir: '/',
})

const emit = defineEmits<{
  (e: 'toggle', node: TreeNode): void
  (e: 'open-file', node: TreeNode): void
  (e: 'create-note', dirPath: string, name: string): void
  (e: 'create-folder', dirPath: string, name: string): void
  (e: 'delete-item', path: string): void
  (e: 'rename-item', oldPath: string, newName: string): void
  (e: 'refresh-tree'): void
}>()

interface ContextMenuState {
  visible: boolean
  x: number
  y: number
  node: TreeNode | null
}

const contextMenu = ref<ContextMenuState>({
  visible: false,
  x: 0,
  y: 0,
  node: null,
})

const selectedPath = ref('')
const showCreateDialog = ref(false)
const createDialogTitle = ref('')
const createDialogPlaceholder = ref('')
const createType = ref<'note' | 'folder'>('note')
const createTargetDir = ref('')

function handleSelect(node: TreeNode) {
  selectedPath.value = node.path
}

function handleToggle(node: TreeNode) {
  emit('toggle', node)
}

function handleOpenFile(node: TreeNode) {
  emit('open-file', node)
}

function handleRefresh() {
  emit('refresh-tree')
}

function showBgContextMenu(e: MouseEvent) {
  contextMenu.value = {
    visible: true,
    x: e.clientX,
    y: e.clientY,
    node: null,
  }
}

function showItemContextMenu(e: MouseEvent, node: TreeNode) {
  contextMenu.value = {
    visible: true,
    x: e.clientX,
    y: e.clientY,
    node,
  }
}

function getTargetDir(): string {
  const { node } = contextMenu.value
  if (node) {
    return node.type === 'dir' ? node.path : getParentDir(node.path)
  }
  return props.rootDir
}

function getParentDir(path: string): string {
  const parts = path.split('/')
  parts.pop()
  return parts.join('/') || '/'
}

function openCreateDialog(type: 'note' | 'folder') {
  createType.value = type
  createTargetDir.value = getTargetDir()
  if (type === 'note') {
    createDialogTitle.value = '新建笔记'
    createDialogPlaceholder.value = '输入笔记名称'
  } else {
    createDialogTitle.value = '新建文件夹'
    createDialogPlaceholder.value = '输入文件夹名称'
  }
  showCreateDialog.value = true
}

function handleCreateConfirm(name: string) {
  showCreateDialog.value = false
  if (createType.value === 'note') {
    const fileName = name.endsWith('.md') ? name : `${name}.md`
    emit('create-note', createTargetDir.value, fileName)
  } else {
    emit('create-folder', createTargetDir.value, name)
  }
}

function handleContextAction(action: string) {
  const { node } = contextMenu.value
  contextMenu.value.visible = false

  switch (action) {
    case 'new-note':
      openCreateDialog('note')
      break
    case 'new-folder':
      openCreateDialog('folder')
      break
    case 'rename': {
      if (!node) return
      const newName = prompt('重命名为：', node.name)
      if (newName && newName !== node.name) {
        emit('rename-item', node.path, newName)
      }
      break
    }
    case 'delete': {
      if (!node) return
      if (confirm(`确定删除「${node.name}」？`)) {
        emit('delete-item', node.path)
      }
      break
    }
    case 'refresh':
      emit('refresh-tree')
      break
  }
}

function handleClickOutside() {
  contextMenu.value.visible = false
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.note-sidebar {
  width: 240px;
  min-width: 180px;
  max-width: 400px;
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  user-select: none;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border-bottom: 1px solid var(--border-color);
}

.header-title {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.refresh-btn {
  width: 24px;
  height: 24px;
  background: transparent;
  border: none;
  border-radius: 4px;
  color: var(--text-secondary);
  font-size: 16px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.refresh-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.sidebar-content {
  flex: 1;
  overflow-y: auto;
  padding: 4px 0;
}

.sidebar-content::-webkit-scrollbar {
  width: 4px;
}

.sidebar-content::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
}

.loading-state,
.empty-sidebar {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  color: var(--text-muted);
  font-size: 12px;
}

.tree {
  padding: 2px 0;
}

.context-menu {
  position: fixed;
  width: 150px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 4px;
  z-index: 1000;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
}

.ctx-item {
  padding: 6px 10px;
  font-size: 12px;
  color: var(--text-primary);
  border-radius: 4px;
  cursor: pointer;
  transition: background 0.1s ease;
}

.ctx-item:hover {
  background: var(--bg-hover);
}

.ctx-item.danger {
  color: var(--error);
}

.ctx-item.danger:hover {
  background: rgba(228, 77, 106, 0.15);
}

.ctx-sep {
  height: 1px;
  background: var(--border-color);
  margin: 4px 8px;
}
</style>
