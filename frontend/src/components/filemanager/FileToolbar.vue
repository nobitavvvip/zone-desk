<template>
  <div class="file-toolbar">
    <div class="toolbar-group">
      <button class="toolbar-btn" @click="$emit('newFolder')" title="新建文件夹 (Ctrl+Shift+N)">
        <span class="toolbar-icon">📁</span>
        <span class="toolbar-text">新建</span>
      </button>
    </div>
    
    <div class="toolbar-divider"></div>
    
    <div class="toolbar-group">
      <button 
        class="toolbar-btn" 
        :disabled="selectedCount === 0"
        @click="$emit('copy')" 
        title="复制 (Ctrl+C)"
      >
        <span class="toolbar-icon">📋</span>
        <span class="toolbar-text">复制</span>
      </button>
      <button 
        class="toolbar-btn" 
        :disabled="selectedCount === 0"
        @click="$emit('cut')" 
        title="剪切 (Ctrl+X)"
      >
        <span class="toolbar-icon">✂️</span>
        <span class="toolbar-text">剪切</span>
      </button>
      <button 
        class="toolbar-btn" 
        :disabled="!hasClipboard"
        @click="$emit('paste')" 
        title="粘贴 (Ctrl+V)"
      >
        <span class="toolbar-icon">📋</span>
        <span class="toolbar-text">粘贴</span>
      </button>
    </div>
    
    <div class="toolbar-divider"></div>
    
    <div class="toolbar-group">
      <button 
        class="toolbar-btn" 
        :disabled="selectedCount === 0"
        @click="$emit('delete')" 
        title="删除 (Delete)"
      >
        <span class="toolbar-icon">🗑️</span>
        <span class="toolbar-text">删除</span>
      </button>
    </div>
    
    <div class="toolbar-divider"></div>
    
    <div class="toolbar-group">
      <button class="toolbar-btn" @click="$emit('refresh')" title="刷新 (F5)">
        <span class="toolbar-icon">🔄</span>
        <span class="toolbar-text">刷新</span>
      </button>
    </div>
    
    <div class="toolbar-divider"></div>
    
    <div class="toolbar-group">
      <button 
        class="toolbar-btn" 
        :disabled="selectedCount !== 1"
        @click="$emit('properties')" 
        title="属性 (Alt+Enter)"
      >
        <span class="toolbar-icon">📊</span>
        <span class="toolbar-text">属性</span>
      </button>
    </div>
    
    <div class="toolbar-spacer"></div>
    
    <div class="toolbar-group view-mode-group">
      <button 
        class="toolbar-btn view-btn" 
        :class="{ active: viewMode === 'details' }"
        @click="$emit('setViewMode', 'details')"
        title="详细信息"
      >
        📋
      </button>
      <button 
        class="toolbar-btn view-btn" 
        :class="{ active: viewMode === 'list' }"
        @click="$emit('setViewMode', 'list')"
        title="列表"
      >
        📄
      </button>
      <button 
        class="toolbar-btn view-btn" 
        :class="{ active: viewMode === 'icons' }"
        @click="$emit('setViewMode', 'icons')"
        title="大图标"
      >
        🖼️
      </button>
      <button 
        class="toolbar-btn view-btn" 
        :class="{ active: viewMode === 'tile' }"
        @click="$emit('setViewMode', 'tile')"
        title="平铺"
      >
        📦
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ViewMode } from '@/store/fileManager'

defineProps<{
  selectedCount: number
  hasClipboard: boolean
  viewMode: ViewMode
}>()

defineEmits<{
  newFolder: []
  copy: []
  cut: []
  paste: []
  delete: []
  refresh: []
  goUp: []
  properties: []
  setViewMode: [mode: ViewMode]
}>()
</script>

<style scoped>
.file-toolbar {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  gap: 4px;
}

.toolbar-group {
  display: flex;
  gap: 2px;
}

.toolbar-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: transparent;
  border: 1px solid transparent;
  border-radius: 4px;
  color: var(--text-primary);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.15s ease;
  white-space: nowrap;
}

.toolbar-btn:hover:not(:disabled) {
  background: var(--bg-hover);
  border-color: var(--border-color);
}

.toolbar-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.toolbar-icon {
  font-size: 16px;
}

.toolbar-text {
  font-size: 13px;
}

.toolbar-divider {
  width: 1px;
  height: 24px;
  background: var(--border-color);
  margin: 0 4px;
}

.toolbar-spacer {
  flex: 1;
}

.view-mode-group {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 2px;
}

.view-btn {
  padding: 4px 8px;
  font-size: 14px;
}

.view-btn.active {
  background: var(--accent);
  color: white;
}
</style>
