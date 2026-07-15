<template>
  <div class="desktop" @contextmenu.prevent="showContextMenu">
    <div class="wallpaper" :class="'wp-' + wallpaper" :style="{ filter: 'blur(' + blur + 'px)' }"></div>
    <div class="mask" :style="{ backgroundColor: 'rgba(10,12,24,' + mask + ')' }"></div>

    <div class="topbar">
      <div class="logo" :style="{ color: iconTextColor, textShadow: iconTextShadow }">ZoneDesk</div>
      <div class="clock" :style="{ color: iconTextColor, textShadow: iconTextShadow }">{{ clockStr }}</div>
    </div>

    <div class="icon-grid">
      <div
        v-for="app in apps"
        :key="app.id"
        class="app-icon"
        @dblclick="openApp(app.id)"
      >
        <div class="icon-box">{{ app.emoji }}</div>
        <div class="label" :style="{ color: iconTextColor, textShadow: iconTextShadow }">{{ app.name }}</div>
      </div>
    </div>

    <FileManagerWindow v-if="store.activeWindow === 'filemanager'" />
    <ContainerWindow v-if="store.activeWindow === 'container'" />

    <div class="settings-panel" :class="{ show: showSettings }" id="settingsPanel">
      <h3>⚙️ 个性化</h3>
      <div class="section">
        <div class="section-title">外观</div>
        <div class="sp-row">
          <span class="sp-label">强调色</span>
          <div class="accent-row">
            <div
              v-for="c in accentColors" :key="c"
              class="accent-dot"
              :class="{ active: accent === c }"
              :style="{ background: c }"
              @click="setAccent(c)"
            ></div>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">桌面背景</div>
        <div class="sp-row">
          <span class="sp-label">壁纸</span>
          <div class="sp-ctrl" style="gap:4px;flex-wrap:wrap">
            <button
              v-for="wp in wallpaperOptions" :key="wp.id"
              class="btn"
              :class="{ active: wallpaper === wp.id }"
              @click="wallpaper = wp.id"
            >{{ wp.label }}</button>
          </div>
        </div>
        <div class="sp-row">
          <span class="sp-label">背景模糊</span>
          <div class="sp-ctrl">
            <input type="range" class="sp-slider" min="0" max="20" v-model.number="blur">
            <span class="sp-val">{{ blur }}px</span>
          </div>
        </div>
        <div class="sp-row">
          <span class="sp-label">遮罩浓度</span>
          <div class="sp-ctrl">
            <input type="range" class="sp-slider" min="0" max="100" :value="Math.round(mask * 100)" @input="mask = clamp(($event.target as HTMLInputElement).valueAsNumber / 100, 0, 1)">
            <span class="sp-val">{{ mask.toFixed(2) }}</span>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="section-title">应用管理</div>
        <div class="sp-hint">在设置面板中添加 / 删除桌面应用</div>
      </div>
    </div>

    <div class="context-menu" :class="{ show: ctxVisible }" :style="ctxStyle" id="contextMenu">
      <div class="ctx-item">打开</div>
      <div class="ctx-item">新窗口打开</div>
      <div class="ctx-sep"></div>
      <div class="ctx-item">重命名</div>
      <div class="ctx-item">删除</div>
      <div class="ctx-sep"></div>
      <div class="ctx-item">属性</div>
    </div>

    <Taskbar @toggle-settings="showSettings = !showSettings" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useDesktopStore } from '@/store/desktop'
import { getSettings, saveSettings } from '@/api/settings'
import Taskbar from '@/components/desktop/Taskbar.vue'
import FileManagerWindow from '@/views/filemanager/FileManagerWindow.vue'
import ContainerWindow from '@/views/container/ContainerWindow.vue'

const store = useDesktopStore()

const apps = [
  { id: 'filemanager', name: '文件管理器', emoji: '📁' },
  { id: 'container', name: '容器管理', emoji: '🐳' },
  { id: 'status', name: '系统状态', emoji: '📊' },
  { id: 'notes', name: '日报笔记', emoji: '📝' },
  { id: 'settings', name: '系统设置', emoji: '⚙️' },
  { id: 'terminal', name: '远程终端', emoji: '🖥️' },
]

const wallpaperOptions = [
  { id: 'none', label: '无' },
  { id: '3', label: '极光' },
  { id: 'cat', label: '窗边小猫' },
]

const accentColors = ['#3b82f6', '#18a058', '#f0a020', '#e44d6a']

const wallpaper = ref('cat')
const blur = ref(0)
const mask = ref(0.35)
const accent = ref('#3b82f6')
const loaded = ref(false)
const showSettings = ref(false)
const ctxVisible = ref(false)
const ctxX = ref(0)
const ctxY = ref(0)
const clockStr = ref('')
let clockTimer: ReturnType<typeof setInterval> | null = null

const iconTextColor = '#fff'
const iconTextShadow = '0 1px 8px rgba(0,0,0,0.6)'

const ctxStyle = computed(() => ({
  left: ctxX.value + 'px',
  top: ctxY.value + 'px',
}))

function clamp(v: number, min: number, max: number) {
  return Math.min(max, Math.max(min, v))
}

function setAccent(c: string) {
  accent.value = c
  document.documentElement.style.setProperty('--accent', c)
}

function openApp(id: string) {
  if (id === 'filemanager' || id === 'container') {
    store.openWindow(id)
  }
}

function showContextMenu(e: MouseEvent) {
  ctxX.value = e.clientX
  ctxY.value = e.clientY
  ctxVisible.value = true
  showSettings.value = false
}

function updateClock() {
  const now = new Date()
  const days = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  clockStr.value = days[now.getDay()] + ' ' + now.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

function handleClick(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (!target.closest('.settings-panel') && !target.closest('.gear')) {
    showSettings.value = false
  }
  if (!target.closest('.context-menu') && !target.closest('.desktop')) {
    ctxVisible.value = false
  }
}

let saveTimer: ReturnType<typeof setTimeout> | null = null
function queueSave() {
  if (saveTimer) clearTimeout(saveTimer)
  saveTimer = setTimeout(() => {
    saveSettings({
      wallpaper: wallpaper.value,
      blur: clamp(blur.value, 0, 20),
      mask: clamp(mask.value, 0, 1),
      accent: accent.value,
    }).catch(() => {})
  }, 500)
}

watch([wallpaper, blur, mask, accent], () => {
  if (!loaded.value) return
  queueSave()
})

onMounted(async () => {
  updateClock()
  clockTimer = setInterval(updateClock, 10000)
  document.addEventListener('click', handleClick)

  try {
    const s = await getSettings()
    wallpaper.value = s.wallpaper
    blur.value = s.blur
    mask.value = s.mask
    setAccent(s.accent)
  } catch {
    // use defaults
  }

  loaded.value = true

  const hash = window.location.hash
  if (hash.startsWith('#/filemanager')) {
    store.openWindow('filemanager')
  }
})

onUnmounted(() => {
  if (clockTimer) clearInterval(clockTimer)
  if (saveTimer) clearTimeout(saveTimer)
  document.removeEventListener('click', handleClick)
})
</script>

<style scoped>
.desktop {
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  position: relative;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  background: #0a0a0f;
  color: #e8e8ec;
}

.wallpaper {
  position: absolute;
  inset: -3%;
  width: 106%;
  height: 106%;
  background-size: cover;
  background-position: center;
  transition: filter 0.3s ease;
}

.wp-none {
  background-image: linear-gradient(135deg, #0f1123 0%, #1a1b2e 50%, #252840 100%);
}

.wp-3 {
  background-image: url('https://images.unsplash.com/photo-1519681393784-d120267933ba?w=1200');
}

.wp-cat {
  background-image: url('https://oss.puboot.com/zonedesk/cat.jpg');
}

.mask {
  position: absolute;
  inset: 0;
  transition: background-color 0.3s ease;
  z-index: 1;
}

.topbar {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  z-index: 10;
  pointer-events: none;
}

.logo {
  font-size: 18px;
  font-weight: 700;
}

.clock {
  font-size: 13px;
  opacity: 0.9;
}

.icon-grid {
  position: absolute;
  top: 64px;
  left: 24px;
  display: grid;
  grid-template-columns: repeat(2, 82px);
  gap: 16px 12px;
  z-index: 10;
}

.app-icon {
  text-align: center;
  cursor: pointer;
  transition: all 0.2s ease;
  padding: 6px;
  border-radius: 10px;
}

.app-icon:hover {
  background: rgba(96, 165, 250, 0.12);
}

.icon-box {
  width: 60px;
  height: 60px;
  margin: 0 auto;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  transition: all 0.2s ease;
}

.app-icon:hover .icon-box {
  background: rgba(96, 165, 250, 0.2);
  border-color: var(--accent);
  transform: scale(1.05);
}

.label {
  font-size: 11px;
  margin-top: 6px;
  font-weight: 500;
  line-height: 1.3;
}

.settings-panel {
  position: absolute;
  right: 16px;
  bottom: 56px;
  width: 320px;
  background: rgba(30, 32, 48, 0.96);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 18px;
  z-index: 50;
  display: none;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.5);
}

.settings-panel.show {
  display: block;
}

.settings-panel h3 {
  font-size: 13px;
  font-weight: 600;
  color: #e8e8ec;
  margin-bottom: 14px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.settings-panel .section {
  margin-bottom: 14px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.settings-panel .section:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.section-title {
  font-size: 11px;
  color: #8b9bbf;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 10px;
}

.sp-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 0;
}

.sp-label {
  font-size: 12px;
  color: #c8d0e0;
}

.sp-ctrl {
  display: flex;
  align-items: center;
  gap: 6px;
}

.sp-slider {
  width: 100px;
  height: 4px;
  -webkit-appearance: none;
  appearance: none;
  background: rgba(255, 255, 255, 0.12);
  border-radius: 2px;
  outline: none;
}

.sp-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: var(--accent);
  cursor: pointer;
}

.sp-val {
  font-size: 10px;
  color: #8b9bbf;
  min-width: 30px;
  text-align: right;
}

.sp-hint {
  font-size: 10px;
  color: #8b9bbf;
  margin-top: 4px;
}

.accent-row {
  display: flex;
  gap: 6px;
}

.accent-dot {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.15s ease;
}

.accent-dot:hover {
  transform: scale(1.1);
}

.accent-dot.active {
  border-color: #fff;
  box-shadow: 0 0 8px var(--accent);
}

.btn {
  padding: 4px 10px;
  border-radius: 4px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  background: rgba(255, 255, 255, 0.06);
  color: #c8d0e0;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn:hover {
  background: rgba(96, 165, 250, 0.15);
  border-color: var(--accent);
}

.btn.active {
  background: var(--accent);
  color: #fff;
  border-color: var(--accent);
}

.context-menu {
  position: fixed;
  width: 180px;
  background: rgba(30, 32, 48, 0.96);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  padding: 6px;
  z-index: 60;
  display: none;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
}

.context-menu.show {
  display: block;
}

.ctx-item {
  padding: 7px 12px;
  font-size: 12px;
  color: #c8d0e0;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.1s ease;
}

.ctx-item:hover {
  background: rgba(96, 165, 250, 0.15);
}

.ctx-sep {
  height: 1px;
  background: rgba(255, 255, 255, 0.06);
  margin: 4px 8px;
}
</style>
