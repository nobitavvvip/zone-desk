import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const saved = localStorage.getItem('zonedesk-theme')
  const isDark = ref(saved !== 'light')

  watch(isDark, (val) => {
    localStorage.setItem('zonedesk-theme', val ? 'dark' : 'light')
    document.documentElement.setAttribute('data-theme', val ? 'dark' : 'light')
  })

  function toggle() {
    isDark.value = !isDark.value
  }

  document.documentElement.setAttribute('data-theme', isDark.value ? 'dark' : 'light')

  return { isDark, toggle }
})
