<template>
  <Teleport to="body">
    <div class="properties-overlay" @click="$emit('close')">
      <div class="properties-dialog" @click.stop>
        <div class="dialog-header">
          <span class="dialog-title">属性 - {{ file.name }}</span>
          <button class="close-btn" @click="$emit('close')">✕</button>
        </div>
        <div class="dialog-body">
          <div class="property-row">
            <span class="property-label">名称:</span>
            <span class="property-value">{{ file.name }}</span>
          </div>
          <div class="property-row">
            <span class="property-label">类型:</span>
            <span class="property-value">{{ getTypeLabel(file.type) }}</span>
          </div>
          <div class="property-row">
            <span class="property-label">路径:</span>
            <span class="property-value">{{ file.path }}</span>
          </div>
          <div class="property-row">
            <span class="property-label">大小:</span>
            <span class="property-value">{{ file.type === 'dir' ? '-' : formatSize(file.size) }}</span>
          </div>
          <div class="property-row">
            <span class="property-label">权限:</span>
            <span class="property-value">{{ file.mode }}</span>
          </div>
          <div class="property-row">
            <span class="property-label">修改时间:</span>
            <span class="property-value">{{ formatDate(file.modTime) }}</span>
          </div>
        </div>
        <div class="dialog-footer">
          <button class="btn btn-primary" @click="$emit('close')">确定</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import type { FileItem } from '@/api/file'

defineProps<{
  file: FileItem
}>()

defineEmits<{
  close: []
}>()

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
</script>

<style scoped>
.properties-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.properties-dialog {
  background: var(--bg-primary);
  border-radius: 4px;
  box-shadow: var(--shadow);
  width: 400px;
  max-width: 90vw;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
}

.dialog-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.close-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 16px;
  cursor: pointer;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
}

.close-btn:hover {
  background: var(--bg-hover);
}

.dialog-body {
  padding: 16px;
}

.property-row {
  display: flex;
  padding: 8px 0;
  border-bottom: 1px solid var(--border-color);
}

.property-row:last-child {
  border-bottom: none;
}

.property-label {
  width: 100px;
  font-size: 13px;
  color: var(--text-secondary);
}

.property-value {
  flex: 1;
  font-size: 13px;
  color: var(--text-primary);
  word-break: break-all;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  padding: 12px 16px;
  border-top: 1px solid var(--border-color);
}

.btn {
  padding: 6px 16px;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-primary {
  background: var(--accent);
  color: white;
  border: 1px solid var(--accent);
}

.btn-primary:hover {
  background: var(--accent-hover);
}
</style>
