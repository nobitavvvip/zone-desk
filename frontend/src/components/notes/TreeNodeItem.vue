<template>
  <div class="tree-node-item">
    <div
      class="node-row"
      :class="{ 'is-dir': node.type === 'dir', 'is-selected': isSelected }"
      :style="{ paddingLeft: depth * 16 + 8 + 'px' }"
      @click.stop="handleClick"
      @dblclick.stop="handleDblClick"
      @contextmenu.stop.prevent="handleContextMenu"
    >
      <span
        v-if="node.type === 'dir'"
        class="expand-arrow"
        :class="{ expanded: node.expanded }"
        @click.stop="handleToggle"
      >
        ▶
      </span>
      <span v-else class="expand-arrow-placeholder"></span>

      <span class="node-icon">{{ node.type === 'dir' ? (node.expanded ? '📂' : '📁') : '📄' }}</span>
      <span class="node-name">{{ node.name }}</span>
    </div>

    <div v-if="node.type === 'dir' && node.expanded && node.children.length > 0" class="node-children">
      <TreeNodeItem
        v-for="child in node.children"
        :key="child.path"
        :node="child"
        :depth="depth + 1"
        :selected-path="selectedPath"
        @select="$emit('select', $event)"
        @toggle="$emit('toggle', $event)"
        @open-file="$emit('open-file', $event)"
        @context-menu="(evt, n) => $emit('context-menu', evt, n)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { TreeNode } from '@/store/notes'

interface Props {
  node: TreeNode
  depth?: number
  selectedPath?: string
}

const props = withDefaults(defineProps<Props>(), {
  depth: 0,
  selectedPath: '',
})

const emit = defineEmits<{
  (e: 'select', node: TreeNode): void
  (e: 'toggle', node: TreeNode): void
  (e: 'open-file', node: TreeNode): void
  (e: 'context-menu', event: MouseEvent, node: TreeNode): void
}>()

const isSelected = computed(() => props.selectedPath === props.node.path)

function handleClick() {
  emit('select', props.node)
  if (props.node.type === 'file') {
    emit('open-file', props.node)
  }
}

function handleDblClick() {
  if (props.node.type === 'dir') {
    emit('toggle', props.node)
  }
}

function handleToggle() {
  emit('toggle', props.node)
}

function handleContextMenu(e: MouseEvent) {
  emit('select', props.node)
  emit('context-menu', e, props.node)
}
</script>

<style scoped>
.tree-node-item {
  width: 100%;
}

.node-row {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 5px 8px;
  margin: 1px 4px;
  font-size: 13px;
  color: var(--text-primary);
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.1s ease;
  white-space: nowrap;
}

.node-row:hover {
  background: var(--bg-hover);
}

.node-row.is-selected {
  background: var(--accent);
  color: white;
}

.node-row.is-selected:hover {
  background: var(--accent-hover, var(--accent));
}

.expand-arrow {
  width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 8px;
  color: var(--text-secondary);
  transition: transform 0.15s ease;
  flex-shrink: 0;
}

.node-row.is-selected .expand-arrow {
  color: white;
}

.expand-arrow.expanded {
  transform: rotate(90deg);
}

.expand-arrow:hover {
  color: var(--text-primary);
}

.expand-arrow-placeholder {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.node-icon {
  font-size: 14px;
  flex-shrink: 0;
}

.node-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  min-width: 0;
}

.node-children {
  width: 100%;
}
</style>
