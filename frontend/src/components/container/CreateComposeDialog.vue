<template>
  <teleport to="body">
    <div class="dialog-overlay" @click.self="$emit('close')">
      <div class="dialog">
        <div class="dialog-header">
          <span>创建 Compose 项目</span>
          <button class="dialog-close" @click="$emit('close')">✕</button>
        </div>
        <div class="dialog-body">
          <div class="form-row">
            <label>项目目录 <span class="hint">存放 compose.yaml 的路径</span></label>
            <input v-model="projectDir" placeholder="." class="form-input" />
          </div>
          <div class="form-row">
            <label>文件名</label>
            <input v-model="filename" placeholder="compose.yaml" class="form-input" />
          </div>
          <div class="form-row">
            <label>Compose 内容 *</label>
            <textarea v-model="content" class="form-input form-textarea" rows="12" placeholder="version: '3'&#10;services:&#10;  app:&#10;    image: nginx"></textarea>
          </div>
        </div>
        <div class="dialog-footer">
          <label class="start-label">
            <input type="checkbox" v-model="startImmediately" /> 创建后启动
          </label>
          <div class="footer-btns">
            <button class="btn" @click="$emit('close')">取消</button>
            <button class="btn btn-primary" @click="handleCreate" :disabled="!content">创建</button>
          </div>
        </div>
      </div>
    </div>
  </teleport>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { createCompose } from '@/api/container'

const emit = defineEmits<{ close: [] }>()
const message = useMessage()

const projectDir = ref('.')
const filename = ref('compose.yaml')
const content = ref('')
const startImmediately = ref(true)

async function handleCreate() {
  if (!content.value) return
  try {
    await createCompose(projectDir.value, filename.value, content.value, startImmediately.value)
    message.success('Compose 项目已创建')
    emit('close')
  } catch (e: any) {
    message.error(e.message)
  }
}
</script>

<style scoped>
.dialog-overlay {
  position: fixed; inset: 0; background: rgba(0,0,0,0.5);
  display: flex; align-items: center; justify-content: center; z-index: 10000;
}
.dialog {
  width: 520px; max-height: 80vh; background: var(--bg-primary);
  border: 1px solid var(--border-color); border-radius: 10px;
  display: flex; flex-direction: column; box-shadow: 0 12px 40px rgba(0,0,0,0.3);
}
.dialog-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 12px 16px; border-bottom: 1px solid var(--border-color);
  font-size: 14px; font-weight: 600; color: var(--text-primary);
}
.dialog-close {
  background: none; border: none; color: var(--text-muted);
  cursor: pointer; font-size: 14px; padding: 4px; border-radius: 4px;
}
.dialog-close:hover { background: var(--bg-hover); color: var(--text-primary); }
.dialog-body {
  padding: 16px; overflow-y: auto; display: flex; flex-direction: column; gap: 10px;
}
.dialog-body::-webkit-scrollbar { width: 6px; }
.dialog-body::-webkit-scrollbar-track { background: transparent; }
.dialog-body::-webkit-scrollbar-thumb {
  background: var(--text-muted); border-radius: 3px;
  border: 2px solid transparent; background-clip: content-box; opacity: 0.4;
}
.dialog-body::-webkit-scrollbar-thumb:hover { opacity: 0.7; }
.form-row { display: flex; flex-direction: column; gap: 4px; }
.form-row label { font-size: 12px; font-weight: 500; color: var(--text-secondary); }
.hint { font-weight: 400; color: var(--text-muted); font-size: 11px; }
.form-input {
  padding: 6px 10px; border: 1px solid var(--border-input);
  border-radius: 4px; background: var(--bg-input); color: var(--text-primary);
  font-size: 13px; outline: none; transition: border-color 0.15s;
}
.form-input:focus { border-color: var(--accent); }
.form-textarea {
  min-height: 200px; line-height: 1.5; resize: vertical;
  font-family: 'SF Mono', Consolas, monospace; font-size: 12px;
}
.dialog-footer {
  display: flex; align-items: center; justify-content: space-between;
  padding: 12px 16px; border-top: 1px solid var(--border-color);
}
.start-label {
  display: flex; align-items: center; gap: 6px;
  font-size: 12px; color: var(--text-secondary); cursor: pointer;
}
.start-label input { accent-color: var(--accent); }
.footer-btns { display: flex; gap: 8px; }
.btn {
  padding: 6px 14px; border: 1px solid var(--border-color);
  border-radius: 4px; background: var(--bg-primary);
  color: var(--text-primary); font-size: 12px; cursor: pointer;
  transition: all 0.1s;
}
.btn:hover:not(:disabled) { background: var(--bg-hover); }
.btn:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-primary { background: var(--accent); color: #fff; border-color: var(--accent); }
.btn-primary:hover:not(:disabled) { opacity: 0.9; }
</style>
