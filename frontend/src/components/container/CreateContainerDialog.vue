<template>
  <teleport to="body">
    <div class="dialog-overlay" @click.self="$emit('close')">
      <div class="dialog">
        <div class="dialog-header">
          <span>创建容器</span>
          <button class="dialog-close" @click="$emit('close')">✕</button>
        </div>
        <div class="dialog-body">
          <div class="form-row">
            <label>镜像 *</label>
            <input v-model="form.image" placeholder="nginx:latest" class="form-input" />
          </div>
          <div class="form-row">
            <label>容器名称</label>
            <input v-model="form.name" placeholder="my-container" class="form-input" />
          </div>
          <div class="form-row">
            <label>端口映射 <span class="hint">宿主机端口:容器端口</span></label>
            <textarea v-model="portStr" placeholder="8080:80" class="form-input form-textarea" rows="3"></textarea>
          </div>
          <div class="form-row">
            <label>环境变量 <span class="hint">每行一个 KEY=VALUE</span></label>
            <textarea v-model="envStr" placeholder="APP_ENV=production" class="form-input form-textarea" rows="3"></textarea>
          </div>
          <div class="form-row">
            <label>卷挂载 <span class="hint">每行一个 宿主机路径:容器路径</span></label>
            <textarea v-model="volStr" placeholder="/data:/app/data" class="form-input form-textarea" rows="3"></textarea>
          </div>
          <div class="form-row">
            <label>网络模式</label>
            <input v-model="form.network" placeholder="bridge / host / none" class="form-input" />
          </div>
          <div class="form-row">
            <label>重启策略</label>
            <select v-model="form.restart" class="form-input">
              <option value="">默认</option>
              <option value="no">no</option>
              <option value="always">always</option>
              <option value="on-failure">on-failure</option>
              <option value="unless-stopped">unless-stopped</option>
            </select>
          </div>
          <div class="form-row">
            <label>命令</label>
            <input v-model="cmdStr" placeholder="参数用空格分隔" class="form-input" />
          </div>
        </div>
        <div class="dialog-footer">
          <button class="btn" @click="$emit('close')">取消</button>
          <button class="btn btn-primary" @click="handleCreate" :disabled="!form.image">创建</button>
        </div>
      </div>
    </div>
  </teleport>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useMessage } from 'naive-ui'
import { createContainer } from '@/api/container'

const emit = defineEmits<{ close: [] }>()
const message = useMessage()

const form = reactive({
  image: '',
  name: '',
  network: '',
  restart: '',
})
const portStr = ref('')
const envStr = ref('')
const volStr = ref('')
const cmdStr = ref('')
const creating = ref(false)

async function handleCreate() {
  if (!form.image) return
  creating.value = true
  try {
    await createContainer({
      image: form.image,
      name: form.name || undefined,
      ports: portStr.value.split('\n').filter(Boolean),
      env: envStr.value.split('\n').filter(Boolean),
      volumes: volStr.value.split('\n').filter(Boolean),
      network: form.network || undefined,
      restart: form.restart || undefined,
      cmd: cmdStr.value ? cmdStr.value.split(/\s+/) : undefined,
    })
    message.success('容器创建成功')
    emit('close')
  } catch (e: any) {
    message.error(e.message)
  } finally {
    creating.value = false
  }
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
  z-index: 10000;
}

.dialog {
  width: 480px;
  max-height: 80vh;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.3);
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.dialog-close {
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  font-size: 14px;
  padding: 4px;
  border-radius: 4px;
}

.dialog-close:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.dialog-body {
  padding: 16px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.dialog-body::-webkit-scrollbar {
  width: 6px;
}

.dialog-body::-webkit-scrollbar-track {
  background: transparent;
}

.dialog-body::-webkit-scrollbar-thumb {
  background: var(--text-muted);
  border-radius: 3px;
  border: 2px solid transparent;
  background-clip: content-box;
  opacity: 0.4;
}

.dialog-body::-webkit-scrollbar-thumb:hover {
  opacity: 0.7;
}

.form-row {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.form-row label {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
}

.hint {
  font-weight: 400;
  color: var(--text-muted);
  font-size: 11px;
}

.form-input {
  padding: 6px 10px;
  border: 1px solid var(--border-input);
  border-radius: 4px;
  background: var(--bg-input);
  color: var(--text-primary);
  font-size: 13px;
  outline: none;
  transition: border-color 0.15s;
}

.form-input:focus {
  border-color: var(--accent);
}

.form-textarea {
  min-height: 60px;
  line-height: 1.4;
  resize: vertical;
  font-family: inherit;
}

select.form-input {
  cursor: pointer;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid var(--border-color);
}

.btn {
  padding: 6px 14px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.1s;
}

.btn:hover:not(:disabled) {
  background: var(--bg-hover);
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: var(--accent);
  color: #fff;
  border-color: var(--accent);
}

.btn-primary:hover:not(:disabled) {
  opacity: 0.9;
}
</style>
