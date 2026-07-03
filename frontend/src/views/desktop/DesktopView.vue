<template>
  <div class="desktop" :class="{ light: !themeStore.isDark }" @contextmenu.prevent>
    <div class="desktop-icons">
      <DesktopIcon
        v-for="app in apps"
        :key="app.id"
        :app="app"
        @dblclick="openApp(app.id)"
      />
    </div>
    <FileManagerWindow v-if="store.activeWindow === 'filemanager'" />
    <Taskbar />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useDesktopStore } from '@/store/desktop'
import { useThemeStore } from '@/store/theme'
import DesktopIcon from '@/components/desktop/DesktopIcon.vue'
import Taskbar from '@/components/desktop/Taskbar.vue'
import FileManagerWindow from '@/views/filemanager/FileManagerWindow.vue'

const store = useDesktopStore()
const themeStore = useThemeStore()

const apps = [
  { id: 'filemanager', name: '文件管理器', icon: 'folder' },
]

function openApp(id: string) {
  store.openWindow(id)
}

onMounted(() => {
  const hash = window.location.hash
  if (hash.startsWith('#/filemanager')) {
    store.openWindow('filemanager')
  }
})
</script>

<style scoped>
.desktop {
  width: 100vw;
  height: 100vh;
  background: linear-gradient(135deg, #faf8f5 0%, #f5f0eb 50%, #ebe5de 100%);
  color: var(--text-primary);
  overflow: hidden;
  position: relative;
}

.desktop.light {
  background: linear-gradient(135deg, #faf8f5 0%, #f5f0eb 50%, #ebe5de 100%);
  color: var(--text-primary);
}

.desktop-icons {
  display: flex;
  flex-direction: column;
  flex-wrap: wrap;
  gap: 8px;
  padding: 16px;
  height: calc(100vh - 48px);
  align-content: flex-start;
}
</style>
