<template>
  <div class="taskbar">
    <div class="taskbar-left">
      <span class="taskbar-item" @click="showAbout = !showAbout">ZoneDesk</span>
    </div>
    <div class="taskbar-center">
      <span
        v-for="win in desktopStore.windowList"
        :key="win"
        class="taskbar-window"
        :class="{
          active: desktopStore.activeWindow === win,
          minimized: desktopStore.getWindowState(win).isMinimized,
        }"
        @click="desktopStore.toggleWindow(win)"
      >
        {{ getWindowName(win) }}
      </span>
    </div>
    <div class="taskbar-right">
      <button class="theme-btn" @click="themeStore.toggle" :title="themeStore.isDark ? '切换亮色模式' : '切换暗色模式'">
        <svg v-if="themeStore.isDark" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
        </svg>
        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/>
        </svg>
      </button>
      <span class="taskbar-item">{{ time }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useDesktopStore } from '@/store/desktop'
import { useThemeStore } from '@/store/theme'

const desktopStore = useDesktopStore()
const themeStore = useThemeStore()
const showAbout = ref(false)
const time = ref('')
let timer: ReturnType<typeof setInterval> | null = null

function getWindowName(id: string) {
  const names: Record<string, string> = { filemanager: '文件管理器' }
  return names[id] || id
}

function updateTime() {
  const now = new Date()
  time.value = now.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  updateTime()
  timer = setInterval(updateTime, 10000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.taskbar {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 48px;
  background: var(--taskbar-bg);
  backdrop-filter: blur(12px);
  display: flex;
  align-items: center;
  padding: 0 12px;
  gap: 8px;
  border-top: 1px solid var(--border-color);
}

.taskbar-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.taskbar-center {
  flex: 1;
  display: flex;
  justify-content: center;
  gap: 8px;
}

.taskbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.taskbar-item {
  padding: 6px 14px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.taskbar-item:hover {
  background: var(--hover-bg);
  color: var(--text-primary);
}

.theme-btn {
  background: none;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  width: 36px;
  height: 36px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.theme-btn svg {
  width: 20px;
  height: 20px;
}

.theme-btn:hover {
  background: var(--hover-bg);
  color: var(--text-primary);
}

.taskbar-window {
  padding: 6px 16px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  border-radius: 4px;
  background: var(--hover-bg-subtle);
  transition: all 0.2s ease;
}

.taskbar-window.active {
  background: var(--hover-bg);
  color: var(--text-primary);
  font-weight: 500;
}

.taskbar-window.minimized {
  opacity: 0.6;
  font-style: italic;
}

.taskbar-window:hover {
  background: var(--hover-bg);
}
</style>
