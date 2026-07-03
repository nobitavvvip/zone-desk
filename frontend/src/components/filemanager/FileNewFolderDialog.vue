<template>
  <Teleport to="body">
    <div class="dialog-overlay" @click="$emit('close')">
      <div class="dialog" @click.stop>
        <div class="dialog-header">
          <span class="dialog-title">新建文件夹</span>
          <button class="close-btn" @click="$emit('close')">✕</button>
        </div>
        <div class="dialog-body">
          <label class="label">输入目录名称：</label>
          <input 
            ref="inputRef"
            v-model="folderName" 
            class="input" 
            placeholder="新建文件夹"
            @keydown.enter="handleConfirm"
          />
        </div>
        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="$emit('close')">取消</button>
          <button class="btn btn-primary" @click="handleConfirm">确定</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const emit = defineEmits<{
  close: []
  confirm: [name: string]
}>()

const inputRef = ref<HTMLInputElement | null>(null)
const folderName = ref('新建文件夹')

onMounted(() => {
  inputRef.value?.focus()
  inputRef.value?.select()
})

function handleConfirm() {
  if (!folderName.value) {
    emit('close')
    return
  }
  
  emit('confirm', folderName.value)
}
</script>

<style scoped>
.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog {
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

.label {
  display: block;
  margin-bottom: 8px;
  font-size: 13px;
  color: var(--text-secondary);
}

.input {
  width: 100%;
  padding: 8px 12px;
  background: var(--bg-primary);
  border: 1px solid var(--border-input);
  border-radius: 4px;
  color: var(--text-primary);
  font-size: 13px;
  outline: none;
  transition: border-color 0.15s ease;
}

.input:focus {
  border-color: var(--accent);
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
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

.btn-secondary {
  background: var(--bg-secondary);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.btn-secondary:hover {
  background: var(--bg-hover);
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
