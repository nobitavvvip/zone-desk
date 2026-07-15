<template>
  <div class="container-window" :style="windowStyle" @keydown="handleKeyDown" tabindex="0">
    <div class="window-header" @mousedown="startDrag" @dblclick="desktopStore.toggleMaximize(windowName)">
      <span class="window-title">🐳 容器管理</span>
      <div class="window-controls">
        <button class="control-btn minimize" @click.stop="desktopStore.minimizeWindow(windowName)">─</button>
        <button class="control-btn maximize" @click.stop="desktopStore.toggleMaximize(windowName)">□</button>
        <button class="control-btn close" @click="desktopStore.closeWindow(windowName)">✕</button>
      </div>
    </div>

    <div class="window-body">
      <div class="tab-bar">
        <button class="tab-btn" :class="{ active: store.activeTab === 'containers' }" @click="switchTab('containers')">
          容器 <span class="tab-badge">{{ store.runningContainers.length }}/{{ store.containers.length }}</span>
        </button>
        <button class="tab-btn" :class="{ active: store.activeTab === 'compose' }" @click="switchTab('compose')">
          Compose <span class="tab-badge">{{ store.composeProjects.length }}</span>
        </button>
        <button class="tab-btn" :class="{ active: store.activeTab === 'images' }" @click="switchTab('images')">
          镜像 <span class="tab-badge">{{ store.images.length }}</span>
        </button>
      </div>

      <!-- ============ 容器标签页 ============ -->
      <div v-if="store.activeTab === 'containers'" class="tab-content">
        <div class="toolbar">
          <div class="toolbar-left">
            <button class="btn" @click="store.fetchContainers()" :disabled="store.loading">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="btn-icon"><path d="M23 4v6h-6M1 20v-6h6"/><path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15"/></svg> 刷新
            </button>
            <input v-model="store.searchQuery" class="search-input" placeholder="搜索名称/镜像/ID..." />
          </div>
          <div class="toolbar-center">
            <template v-if="store.hasSelection">
              <button class="btn btn-sm btn-success" @click="store.batchStart()">▶ 批量启动</button>
              <button class="btn btn-sm btn-warning" @click="store.batchStop()">■ 批量停止</button>
              <button class="btn btn-sm btn-danger" @click="handleBatchRemove">🗑 批量删除</button>
              <button class="btn btn-sm" @click="store.clearSelection()">取消选择</button>
            </template>
          </div>
          <div class="toolbar-right">
            <button class="btn" @click="showCreateDialog = true">+ 创建</button>
            <label class="toggle-label">
              <input type="checkbox" :checked="store.showAll" @change="toggleShowAll" /> 全部
            </label>
            <button class="btn btn-danger" @click="handlePrune">清理</button>
          </div>
        </div>

        <div class="main-area">
          <div class="table-wrapper" :class="{ 'with-detail': store.detailContainer }">
            <table class="container-table">
              <thead>
                <tr>
                  <th class="col-check"><input type="checkbox" :checked="allChecked" @change="toggleAll" /></th>
                  <th class="col-status"></th>
                  <th class="col-name">名称</th>
                  <th class="col-image">镜像</th>
                  <th class="col-status-text">状态</th>
                  <th class="col-ports">端口</th>
                  <th class="col-actions">操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-if="!store.loading && store.filteredContainers.length === 0">
                  <td colspan="7" class="empty-row">{{ store.error || '无容器' }}</td>
                </tr>
                <tr v-for="c in store.filteredContainers" :key="c.id"
                  :class="{ 'row-running': c.state === 'running', 'row-selected': store.selectedIds.has(c.id) }"
                  @dblclick="showDetail(c)"
                >
                  <td class="col-check"><input type="checkbox" :checked="store.selectedIds.has(c.id)" @change="store.toggleSelect(c.id)" /></td>
                  <td class="col-status"><span class="status-dot" :class="c.state === 'running' ? 'dot-green' : c.state === 'exited' ? 'dot-red' : 'dot-yellow'"></span></td>
                  <td class="col-name" :title="c.id">{{ c.name }}</td>
                  <td class="col-image">{{ c.image }}</td>
                  <td class="col-status-text"><span class="status-label" :class="'status-' + c.state">{{ c.status }}</span></td>
                  <td class="col-ports">{{ c.ports || '-' }}</td>
                  <td class="col-actions">
                    <button v-if="c.state !== 'running'" class="btn btn-xs btn-success" @click="store.start(c.id)" title="启动">▶ 启动</button>
                    <button v-if="c.state === 'running'" class="btn btn-xs btn-warning" @click="store.stop(c.id)" title="停止">■ 停止</button>
                    <button class="btn btn-xs" @click="store.restart(c.id)" title="重启">↻ 重启</button>
                    <button class="btn btn-xs" @click="showLogs(c.id, c.name)" title="日志">📋 日志</button>
                    <button class="btn btn-xs" @click="showDetail(c)" title="详情">ℹ 详情</button>
                    <button class="btn btn-xs btn-danger" @click="handleRemove(c)" title="删除">🗑 删除</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <transition name="slide">
            <div v-if="store.detailContainer" class="detail-panel">
              <div class="detail-header">
                <span>📋 {{ store.detailContainer.name }}</span>
                <button class="btn btn-xs" @click="store.detailContainer = null">✕</button>
              </div>
              <div class="detail-body">
                <div class="detail-row"><span class="detail-label">ID</span><span class="detail-val mono">{{ store.detailContainer.id }}</span></div>
                <div class="detail-row"><span class="detail-label">镜像</span><span class="detail-val">{{ store.detailContainer.image }}</span></div>
                <div class="detail-row"><span class="detail-label">状态</span><span class="detail-val"><span class="status-label" :class="'status-' + store.detailContainer.state">{{ store.detailContainer.status }}</span></span></div>
                <div class="detail-row"><span class="detail-label">端口</span><span class="detail-val">{{ store.detailContainer.ports || '-' }}</span></div>
                <div class="detail-row"><span class="detail-label">创建时间</span><span class="detail-val">{{ store.detailContainer.created }}</span></div>
                <div class="detail-row"><span class="detail-label">Hostname</span><span class="detail-val mono">{{ store.detailInfo }}</span></div>
              </div>
              <div class="detail-actions">
                <button v-if="store.detailContainer.state !== 'running'" class="btn btn-sm btn-success" @click="store.start(store.detailContainer!.id); store.detailContainer = null">▶ 启动</button>
                <button v-if="store.detailContainer.state === 'running'" class="btn btn-sm btn-warning" @click="store.stop(store.detailContainer!.id); store.detailContainer = null">■ 停止</button>
                <button class="btn btn-sm" @click="store.restart(store.detailContainer!.id); store.detailContainer = null">↻ 重启</button>
              </div>
            </div>
          </transition>
        </div>

        <div v-if="store.containerLogs" class="log-panel">
          <div class="log-header">
            <div class="log-header-left">
              <span class="log-title">📋 {{ logName }}</span>
              <input v-model="logSearch" class="log-search" placeholder="过滤日志..." />
            </div>
            <div class="log-header-right">
              <button class="btn btn-xs" :class="{ active: logFollow }" @click="logFollow = !logFollow" title="自动滚动">▼ 跟随</button>
              <button class="btn btn-xs" @click="store.containerLogs = ''">✕ 关闭</button>
            </div>
          </div>
          <div class="log-body" ref="logBodyRef">
            <div v-for="(line, i) in filteredLogLines" :key="i" class="log-line">
              <span class="log-num">{{ i + 1 }}</span>
              <span class="log-text" :class="logLineClass(line)">{{ line }}</span>
            </div>
            <div v-if="filteredLogLines.length === 0" class="log-empty">无匹配日志</div>
          </div>
        </div>
      </div>

      <!-- ============ Compose 标签页 ============ -->
      <div v-else-if="store.activeTab === 'compose'" class="tab-content">
        <div class="toolbar">
          <div class="toolbar-left">
            <button class="btn" @click="store.fetchCompose()" :disabled="store.loading"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="btn-icon"><path d="M23 4v6h-6M1 20v-6h6"/><path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15"/></svg> 刷新</button>
          </div>
          <div class="toolbar-right">
            <button class="btn" @click="showCreateComposeDialog = true">+ 创建</button>
          </div>
        </div>
        <div class="table-wrapper">
          <table class="container-table">
            <thead><tr><th>项目名称</th><th>状态</th><th>配置文件</th><th>操作</th></tr></thead>
            <tbody>
              <tr v-if="!store.loading && store.composeProjects.length === 0"><td colspan="4" class="empty-row">{{ store.error || '无 Compose 项目' }}</td></tr>
              <tr v-for="p in store.composeProjects" :key="p.name">
                <td>{{ p.name }}</td>
                <td><span class="status-label" :class="'status-' + (p.status === 'running' ? 'running' : 'exited')">{{ p.status }}</span></td>
                <td class="cell-path" :title="p.configFiles">{{ p.configFiles || '-' }}</td>
                <td class="col-actions">
                  <button class="btn btn-xs btn-success" @click="store.composeUp(p.workingDir, p.configFiles)" title="启动">▶ 启动</button>
                  <button class="btn btn-xs btn-warning" @click="store.composeStop(p.workingDir, p.configFiles)" title="停止">■ 停止</button>
                  <button class="btn btn-xs" @click="store.composeRestart(p.workingDir, p.configFiles)" title="重启">↻ 重启</button>
                  <button class="btn btn-xs btn-danger" @click="store.composeDown(p.workingDir, p.configFiles)" title="下线">↓ 下线</button>
                  <button class="btn btn-xs" @click="showComposeLogs(p)" title="日志">📋 日志</button>
                  <button class="btn btn-xs" @click="openComposeEditor(p)" title="编辑 YAML">✏ 编辑</button>
                  <button class="btn btn-xs btn-danger" @click="handleDeleteCompose(p)" title="删除">🗑 删除</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-if="store.composeLogs" class="log-panel">
          <div class="log-header">
            <div class="log-header-left">
              <span class="log-title">📋 {{ composeLogName }}</span>
              <input v-model="composeLogSearch" class="log-search" placeholder="过滤日志..." />
            </div>
            <div class="log-header-right">
              <button class="btn btn-xs" :class="{ active: logFollow }" @click="logFollow = !logFollow" title="自动滚动">▼ 跟随</button>
              <button class="btn btn-xs" @click="store.composeLogs = ''">✕ 关闭</button>
            </div>
          </div>
          <div class="log-body" ref="composeLogBodyRef">
            <div v-for="(line, i) in filteredComposeLogLines" :key="i" class="log-line">
              <span class="log-num">{{ i + 1 }}</span>
              <span class="log-text" :class="logLineClass(line)">{{ line }}</span>
            </div>
            <div v-if="filteredComposeLogLines.length === 0" class="log-empty">无匹配日志</div>
          </div>
        </div>
        <div v-if="composeEditor" class="log-panel editor-panel">
          <div class="log-header">
            <span>✏ 编辑 {{ composeEditor.filename }} ({{ composeEditor.projectDir }})</span>
            <div class="log-header-actions">
              <button class="btn btn-xs" @click="saveComposeEditor">💾 保存</button>
              <button class="btn btn-xs btn-danger" @click="composeEditor = null">✕</button>
            </div>
          </div>
          <textarea class="editor-content" v-model="composeEditor.content" spellcheck="false"></textarea>
        </div>
      </div>

      <!-- ============ 镜像标签页 ============ -->
      <div v-else class="tab-content">
        <div class="toolbar">
          <div class="toolbar-left">
            <button class="btn" @click="store.fetchImages()" :disabled="store.loading"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="btn-icon"><path d="M23 4v6h-6M1 20v-6h6"/><path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15"/></svg> 刷新</button>
          </div>
          <div class="toolbar-right">
            <div class="pull-row">
              <input v-model="pullImageName" class="search-input" placeholder="nginx:latest" @keydown.enter="handlePullImage" />
              <button class="btn" @click="handlePullImage" :disabled="!pullImageName">拉取</button>
            </div>
          </div>
        </div>
        <div class="table-wrapper">
          <table class="container-table">
            <thead><tr><th>仓库</th><th>标签</th><th>ID</th><th>大小</th><th>创建时间</th><th>操作</th></tr></thead>
            <tbody>
              <tr v-if="!store.loading && store.images.length === 0"><td colspan="6" class="empty-row">无镜像</td></tr>
              <tr v-for="img in store.images" :key="img.id">
                <td>{{ img.repository || '&lt;none&gt;' }}</td>
                <td><span class="tag-pill">{{ img.tag || 'latest' }}</span></td>
                <td class="mono">{{ img.id.substring(0, 12) }}</td>
                <td>{{ img.size }}</td>
                <td>{{ img.created }}</td>
                <td class="col-actions">
                  <button class="btn btn-xs btn-danger" @click="handleRemoveImage(img)">🗑</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <CreateContainerDialog v-if="showCreateDialog" @close="handleCreateClose" />
    <CreateComposeDialog v-if="showCreateComposeDialog" @close="handleCreateComposeClose" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useDesktopStore } from '@/store/desktop'
import { useContainerStore } from '@/store/container'
import { useDialog, useMessage } from 'naive-ui'
import type { Container, ComposeProject, DockerImage } from '@/api/container'
import CreateContainerDialog from '@/components/container/CreateContainerDialog.vue'
import CreateComposeDialog from '@/components/container/CreateComposeDialog.vue'

const desktopStore = useDesktopStore()
const store = useContainerStore()
const dialog = useDialog()
const message = useMessage()

const windowName = 'container'
const isDragging = ref(false)
const dragOffset = ref({ x: 0, y: 0 })
const logName = ref('')
const composeLogName = ref('')
const logSearch = ref('')
const composeLogSearch = ref('')
const logFollow = ref(true)
const logBodyRef = ref<HTMLElement | null>(null)
const composeLogBodyRef = ref<HTMLElement | null>(null)
const showCreateDialog = ref(false)
const showCreateComposeDialog = ref(false)
const pullImageName = ref('')
const composeEditor = ref<{ projectDir: string; filename: string; content: string } | null>(null)

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

const allChecked = computed(() =>
  store.filteredContainers.length > 0 && store.filteredContainers.every(c => store.selectedIds.has(c.id))
)

const filteredLogLines = computed(() => {
  const lines = store.containerLogs.split('\n').filter(Boolean)
  if (!logSearch.value) return lines
  const q = logSearch.value.toLowerCase()
  return lines.filter(l => l.toLowerCase().includes(q))
})

const filteredComposeLogLines = computed(() => {
  const lines = store.composeLogs.split('\n').filter(Boolean)
  if (!composeLogSearch.value) return lines
  const q = composeLogSearch.value.toLowerCase()
  return lines.filter(l => l.toLowerCase().includes(q))
})

function toggleAll() {
  allChecked.value ? store.clearSelection() : store.selectAll()
}

function startDrag(event: MouseEvent) {
  if (desktopStore.getWindowState(windowName).isMaximized) return
  isDragging.value = true
  const state = desktopStore.getWindowState(windowName)
  dragOffset.value = { x: event.clientX - state.x, y: event.clientY - state.y }
  document.addEventListener('mousemove', onDrag)
  document.addEventListener('mouseup', stopDrag)
}

function onDrag(event: MouseEvent) {
  if (!isDragging.value) return
  desktopStore.updatePosition(windowName, event.clientX - dragOffset.value.x, event.clientY - dragOffset.value.y)
}

function stopDrag() {
  isDragging.value = false
  document.removeEventListener('mousemove', onDrag)
  document.removeEventListener('mouseup', stopDrag)
}

function handleKeyDown(event: KeyboardEvent) {
  if (event.key === 'F5') { event.preventDefault(); refresh() }
}

function switchTab(tab: 'containers' | 'compose' | 'images') {
  store.setActiveTab(tab)
  store.containerLogs = ''
  store.composeLogs = ''
  store.detailContainer = null
  store.error = ''
  refresh()
}

function toggleShowAll() {
  store.toggleShowAll()
  store.fetchContainers()
}

function refresh() {
  if (store.activeTab === 'containers') store.fetchContainers()
  else if (store.activeTab === 'compose') store.fetchCompose()
  else store.fetchImages()
}

function showLogs(id: string, name: string) {
  logName.value = name
  logSearch.value = ''
  store.fetchLogs(id, 200)
}

function showComposeLogs(p: ComposeProject) {
  composeLogName.value = p.name
  composeLogSearch.value = ''
  store.fetchComposeLogs(p.workingDir, p.configFiles, 200)
}

function logLineClass(line: string): string {
  if (/error|exception|traceback|fail|warn/i.test(line)) return 'line-warn'
  if (/^\d{4}[-\/]\d{2}[-\/]\d{2}[T ]\d{2}:\d{2}:\d{2}/.test(line)) return 'line-ts'
  return ''
}

watch(() => store.containerLogs, () => {
  if (logFollow.value && logBodyRef.value) {
    requestAnimationFrame(() => { logBodyRef.value!.scrollTop = logBodyRef.value!.scrollHeight })
  }
})

watch(() => store.composeLogs, () => {
  if (logFollow.value && composeLogBodyRef.value) {
    requestAnimationFrame(() => { composeLogBodyRef.value!.scrollTop = composeLogBodyRef.value!.scrollHeight })
  }
})

async function openComposeEditor(p: ComposeProject) {
  try {
    const res = await store.readComposeFile(p.workingDir, p.configFiles)
    composeEditor.value = { projectDir: res.projectDir, filename: res.filename, content: res.content }
  } catch (e: any) {
    message.error(e.message)
  }
}

async function saveComposeEditor() {
  if (!composeEditor.value) return
  try {
    await store.updateComposeFile(composeEditor.value.projectDir, composeEditor.value.filename, composeEditor.value.content)
    message.success('已保存')
    composeEditor.value = null
  } catch (e: any) {
    message.error(e.message)
  }
}

function handleDeleteCompose(p: ComposeProject) {
  dialog.warning({
    title: '删除 Compose 项目',
    content: `确定要删除 "${p.name}" 的 Compose 文件吗？`,
    positiveText: '删除', negativeText: '取消',
    positiveButtonProps: { type: 'error' },
    onPositiveClick: async () => {
      try { await store.deleteComposeProject(p.workingDir, p.configFiles); message.success('已删除') }
      catch (e: any) { message.error(e.message) }
    },
  })
}

function showDetail(c: Container) {
  store.fetchDetail(c.id)
}

function handleRemove(c: Container) {
  dialog.warning({
    title: '删除容器',
    content: `确定要删除容器 "${c.name}" 吗？`,
    positiveText: '删除', negativeText: '取消',
    positiveButtonProps: { type: 'error' },
    onPositiveClick: async () => {
      try { await store.remove(c.id, true); message.success('已删除') }
      catch (e: any) { message.error(e.message) }
    },
  })
}

function handleBatchRemove() {
  dialog.warning({
    title: '批量删除',
    content: `确定要删除选中的 ${store.selectedIds.size} 个容器吗？`,
    positiveText: '删除', negativeText: '取消',
    positiveButtonProps: { type: 'error' },
    onPositiveClick: async () => {
      try { await store.batchRemove(); message.success('已删除') }
      catch (e: any) { message.error(e.message) }
    },
  })
}

function handlePrune() {
  dialog.warning({
    title: '清理容器',
    content: '确定要清理所有已停止的容器吗？此操作不可恢复。',
    positiveText: '清理', negativeText: '取消',
    positiveButtonProps: { type: 'error' },
    onPositiveClick: async () => {
      try { await store.prune(); message.success('清理完成') }
      catch (e: any) { message.error(e.message) }
    },
  })
}

async function handlePullImage() {
  if (!pullImageName.value) return
  try {
    await store.pullImage(pullImageName.value)
    message.success(`拉取完成: ${pullImageName.value}`)
    pullImageName.value = ''
  } catch (e: any) {
    message.error(e.message)
  }
}

function handleRemoveImage(img: DockerImage) {
  const tag = img.repository ? `${img.repository}:${img.tag || 'latest'}` : img.id.substring(0, 12)
  dialog.warning({
    title: '删除镜像',
    content: `确定要删除镜像 "${tag}" 吗？`,
    positiveText: '删除', negativeText: '取消',
    positiveButtonProps: { type: 'error' },
    onPositiveClick: async () => {
      try { await store.removeImage(tag, true); message.success('已删除') }
      catch (e: any) { message.error(e.message) }
    },
  })
}

function handleCreateClose() {
  showCreateDialog.value = false
  store.fetchContainers()
}

function handleCreateComposeClose() {
  showCreateComposeDialog.value = false
  store.fetchCompose()
}

onMounted(() => { store.fetchContainers() })
onUnmounted(() => {
  document.removeEventListener('mousemove', onDrag)
  document.removeEventListener('mouseup', stopDrag)
})
</script>

<style scoped>
.container-window {
  position: absolute;
  left: 6vw; top: 6vh; width: 88vw; height: 84vh;
  border-radius: 4px; overflow: hidden;
  background: var(--bg-primary); box-shadow: var(--shadow);
  display: flex; flex-direction: column; z-index: 100;
  border: 1px solid var(--border-color); outline: none;
}
.window-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 8px 12px; background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
}
.window-title { font-size: 13px; font-weight: 600; color: var(--text-primary); }
.window-controls { display: flex; gap: 2px; }
.control-btn {
  width: 40px;
  height: 30px;
  background: transparent; border: none;
  color: var(--text-primary); font-size: 14px; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
  border-radius: 2px; transition: background 0.1s ease;
}
.control-btn:hover { background: var(--bg-hover); }
.control-btn.close:hover { background: var(--error); color: white; }
.window-body { flex: 1; display: flex; flex-direction: column; overflow: hidden; }

.tab-bar {
  display: flex; background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color); padding: 0 12px;
}
.tab-btn {
  padding: 8px 16px; background: transparent; border: none;
  border-bottom: 2px solid transparent; color: var(--text-secondary);
  font-size: 13px; cursor: pointer; display: flex; align-items: center; gap: 6px;
  transition: all 0.15s ease;
}
.tab-btn:hover { color: var(--text-primary); background: var(--bg-hover); }
.tab-btn.active { color: var(--accent); border-bottom-color: var(--accent); font-weight: 600; }
.tab-badge {
  font-size: 11px; background: var(--bg-tertiary); padding: 1px 6px;
  border-radius: 8px; color: var(--text-muted);
}
.tab-btn.active .tab-badge { background: rgba(59,130,246,0.15); color: var(--accent); }

.tab-content { flex: 1; display: flex; flex-direction: column; overflow: hidden; }

.toolbar {
  display: flex; align-items: center; justify-content: space-between;
  padding: 8px 12px; background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color); gap: 8px; flex-wrap: wrap;
}
.toolbar-left, .toolbar-center, .toolbar-right {
  display: flex; align-items: center; gap: 6px;
}
.search-input {
  width: 200px; padding: 4px 8px; border: 1px solid var(--border-input);
  border-radius: 4px; background: var(--bg-input); color: var(--text-primary);
  font-size: 12px; outline: none;
}
.search-input:focus { border-color: var(--accent); }
.pull-row { display: flex; gap: 4px; }

.btn {
  display: inline-flex; align-items: center; gap: 4px;
  padding: 4px 10px; border: 1px solid var(--border-color);
  border-radius: 4px; background: var(--bg-primary);
  color: var(--text-primary); font-size: 12px; cursor: pointer; transition: all 0.1s;
  white-space: nowrap;
}
.btn:hover:not(:disabled) { background: var(--bg-hover); border-color: var(--border-focus); }
.btn:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-sm { padding: 5px 10px; font-size: 13px; }
.btn-xs {
  padding: 6px 10px; font-size: 13px; border: 1px solid transparent;
  background: transparent; min-width: 32px; min-height: 32px;
}
.btn-xs:hover { background: var(--bg-hover); border-color: var(--border-color); }
.btn-xs.active { background: var(--accent); color: #fff; border-color: var(--accent); }
.btn-icon { width: 14px; height: 14px; }
.btn-success { color: var(--success); }
.btn-success:hover { background: rgba(22,163,74,0.12); }
.btn-warning { color: var(--warning); }
.btn-warning:hover { background: rgba(234,88,12,0.12); }
.btn-danger { color: var(--error); }
.btn-danger:hover { background: rgba(239,68,68,0.12); }

.toggle-label {
  display: flex; align-items: center; gap: 4px;
  font-size: 12px; color: var(--text-secondary); cursor: pointer;
}
.toggle-label input { accent-color: var(--accent); }

.main-area { flex: 1; display: flex; overflow: hidden; }

.table-wrapper { flex: 1; overflow: auto; }
.table-wrapper.with-detail { flex: 1; }
.table-wrapper::-webkit-scrollbar { width: 8px; height: 8px; }
.table-wrapper::-webkit-scrollbar-track { background: transparent; }
.table-wrapper::-webkit-scrollbar-thumb {
  background: var(--text-muted); border-radius: 4px;
  border: 2px solid transparent; background-clip: content-box; opacity: 0.4;
}
.table-wrapper::-webkit-scrollbar-thumb:hover { opacity: 0.7; }

.container-table { width: 100%; border-collapse: collapse; font-size: 14px; }
.container-table th {
  padding: 8px 12px; text-align: left; font-weight: 600; color: var(--text-secondary);
  background: var(--bg-secondary); border-bottom: 1px solid var(--border-color);
  font-size: 13px; white-space: nowrap; position: sticky; top: 0; z-index: 1;
}
.container-table td { padding: 8px 12px; border-bottom: 1px solid var(--border-color); color: var(--text-primary); }
.container-table tbody tr:hover { background: var(--hover-bg-subtle); }
.row-running { background: rgba(22,163,74,0.03); }
.row-selected { background: rgba(59,130,246,0.06); }
.empty-row { text-align: center !important; padding: 40px !important; color: var(--text-muted) !important; }
.col-check { width: 28px; text-align: center; }
.col-status { width: 28px; text-align: center; }
.status-dot { display: inline-block; width: 10px; height: 10px; flex-shrink: 0; }
.dot-green { background: #16a34a; border-radius: 50%; box-shadow: 0 0 6px rgba(22,163,74,0.5); }
.dot-red { background: #ef4444; border-radius: 2px; }
.dot-yellow { background: #eab308; border-radius: 2px; transform: rotate(45deg); }
.status-label { display: inline-block; padding: 3px 10px; border-radius: 10px; font-size: 12px; font-weight: 600; }
.status-running { background: rgba(22,163,74,0.12); color: #16a34a; }
.status-exited { background: rgba(239,68,68,0.12); color: #ef4444; }
.status-paused { background: rgba(234,179,8,0.12); color: #eab308; }
.col-actions { white-space: nowrap; }
.cell-path { max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.mono { font-family: 'SF Mono', Consolas, monospace; font-size: 12px; }
.tag-pill {
  display: inline-block; padding: 2px 8px; border-radius: 4px;
  background: var(--accent); color: #fff; font-size: 12px; font-weight: 500;
}

.detail-panel {
  width: 280px; border-left: 1px solid var(--border-color);
  background: var(--bg-primary); display: flex; flex-direction: column; overflow: hidden;
}
.detail-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 8px 12px; background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color); font-size: 12px; font-weight: 600;
}
.detail-body { flex: 1; overflow: auto; padding: 12px; display: flex; flex-direction: column; gap: 10px; font-size: 13px; }
.detail-row { display: flex; flex-direction: column; gap: 2px; }
.detail-label { font-size: 12px; color: var(--text-muted); font-weight: 500; }
.detail-val { color: var(--text-primary); word-break: break-all; }
.detail-actions {
  display: flex; gap: 4px; padding: 8px 12px;
  border-top: 1px solid var(--border-color);
}

.log-panel {
  border-top: 1px solid var(--border-color);
  display: flex; flex-direction: column;
  height: 240px; min-height: 120px;
}
.log-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 6px 12px; background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
}
.log-header-left, .log-header-right {
  display: flex; align-items: center; gap: 8px;
}
.log-title {
  font-size: 12px; color: var(--text-secondary); font-weight: 600; white-space: nowrap;
}
.log-search {
  width: 180px; padding: 2px 8px;
  border: 1px solid var(--border-input); border-radius: 4px;
  background: var(--bg-input); color: var(--text-primary);
  font-size: 12px; outline: none;
}
.log-search:focus { border-color: var(--accent); }
.log-body {
  flex: 1; overflow-y: auto; padding: 0;
  font-family: 'SF Mono', Consolas, monospace; font-size: 12px; line-height: 1.7;
  background: var(--bg-primary);
}
.log-body::-webkit-scrollbar { width: 8px; }
.log-body::-webkit-scrollbar-track { background: transparent; }
.log-body::-webkit-scrollbar-thumb {
  background: var(--text-muted); border-radius: 4px;
  border: 2px solid transparent; background-clip: content-box;
}
.log-line {
  display: flex; padding: 0 12px;
  transition: background 0.1s;
}
.log-line:hover { background: var(--hover-bg-subtle); }
.log-num {
  width: 40px; flex-shrink: 0; text-align: right;
  padding-right: 10px; color: var(--text-muted); opacity: 0.5;
  user-select: none;
}
.log-text { flex: 1; white-space: pre-wrap; word-break: break-all; color: var(--text-primary); }
.log-text.line-ts { color: var(--text-secondary); }
.log-text.line-warn { color: var(--warning); background: rgba(234,88,12,0.06); }
.log-empty {
  padding: 20px; text-align: center; color: var(--text-muted); font-size: 12px;
}
.editor-panel {
  max-height: 300px;
}

.editor-content {
  flex: 1; overflow: auto; padding: 8px 12px; margin: 0;
  font-family: 'SF Mono', Consolas, monospace; font-size: 12px; line-height: 1.5;
  color: var(--text-primary); background: var(--bg-primary);
  border: none; outline: none; resize: none;
}
.editor-content::-webkit-scrollbar { width: 6px; }
.editor-content::-webkit-scrollbar-track { background: transparent; }
.editor-content::-webkit-scrollbar-thumb {
  background: var(--text-muted); border-radius: 3px;
  border: 2px solid transparent; background-clip: content-box; opacity: 0.4;
}

.slide-enter-active, .slide-leave-active { transition: width 0.2s ease, opacity 0.2s ease; }
.slide-enter-from, .slide-leave-to { width: 0 !important; opacity: 0; overflow: hidden; }
</style>
