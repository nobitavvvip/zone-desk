<template>
  <div class="taskbar">
    <div class="taskbar-left"></div>
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
      <button class="gear" @click="$emit('toggleSettings')" title="个性化设置">⚙</button>
      <button class="theme-btn" @click="themeStore.toggle" :title="themeStore.isDark ? '切换亮色模式' : '切换暗色模式'">
        <svg v-if="themeStore.isDark" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
        </svg>
        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useDesktopStore } from '@/store/desktop'
import { useThemeStore } from '@/store/theme'

defineEmits<{
  toggleSettings: []
}>()

const desktopStore = useDesktopStore()
const themeStore = useThemeStore()

function getWindowName(id: string) {
  const names: Record<string, string> = { filemanager: '文件管理器', container: '容器管理' }
  return names[id] || id
}
</script>

<style scoped>
.taskbar {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 44px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  font-size: 12px;
  color: #e8e8ec;
  z-index: 20;
}

.taskbar::before {
  content: '';
  position: absolute;
  inset: 0;
  background: rgba(20, 22, 34, 0.45);
  opacity: 0.6;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  pointer-events: none;
}

.taskbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
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
  gap: 10px;
}

.taskbar-window {
  padding: 5px 14px;
  border-radius: 6px;
  background: rgba(96, 165, 250, 0.12);
  color: #e8e8ec;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.1s ease;
}

.taskbar-window.active {
  background: var(--accent);
  color: #fff;
}

.taskbar-window.minimized {
  opacity: 0.6;
  font-style: italic;
}

.taskbar-window:hover {
  opacity: 0.9;
}

.gear {
  width: 30px;
  height: 30px;
  border-radius: 7px;
  background: rgba(96, 165, 250, 0.18);
  color: var(--accent);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 16px;
  border: none;
  transition: all 0.2s ease;
}

.gear:hover {
  background: var(--accent);
  color: #fff;
}

.theme-btn {
  background: none;
  border: none;
  color: #c8d0e0;
  cursor: pointer;
  width: 30px;
  height: 30px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.theme-btn svg {
  width: 18px;
  height: 18px;
}

.theme-btn:hover {
  background: rgba(96, 165, 250, 0.15);
}
</style>
