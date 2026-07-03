<template>
  <n-config-provider :theme="themeStore.isDark ? darkTheme : null">
    <n-message-provider>
      <n-dialog-provider>
        <DesktopView />
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { darkTheme, NMessageProvider, NDialogProvider } from 'naive-ui'
import { useThemeStore } from '@/store/theme'
import DesktopView from './views/desktop/DesktopView.vue'

const themeStore = useThemeStore()
</script>
