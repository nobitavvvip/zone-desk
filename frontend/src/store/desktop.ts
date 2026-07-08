import { defineStore } from 'pinia'
import { ref, reactive } from 'vue'

interface WindowState {
  x: number
  y: number
  width: number
  height: number
  isMaximized: boolean
  isMinimized: boolean
  _prevX?: number
  _prevY?: number
  _prevWidth?: number
  _prevHeight?: number
}

export const useDesktopStore = defineStore('desktop', () => {
  const activeWindow = ref<string | null>(null)
  const windowList = ref<string[]>([])
  const windows = reactive<Record<string, WindowState>>({})

  function getWindowState(name: string): WindowState {
    if (!windows[name]) {
      windows[name] = {
        x: 50,
        y: 50,
        width: Math.floor(window.innerWidth * 0.9),
        height: Math.floor(window.innerHeight * 0.86),
        isMaximized: false,
        isMinimized: false,
      }
    }
    return windows[name]
  }

  function openWindow(name: string) {
    activeWindow.value = name
    const state = getWindowState(name)
    state.isMinimized = false
    if (!windowList.value.includes(name)) {
      windowList.value.push(name)
    }
  }

  function closeWindow(name: string) {
    activeWindow.value = null
    const idx = windowList.value.indexOf(name)
    if (idx !== -1) {
      windowList.value.splice(idx, 1)
    }
    delete windows[name]
  }

  function minimizeWindow(name: string) {
    const state = getWindowState(name)
    state.isMinimized = true
  }

  function toggleWindow(name: string) {
    const state = getWindowState(name)
    if (state.isMinimized || activeWindow.value !== name) {
      openWindow(name)
    } else {
      minimizeWindow(name)
    }
  }

  function toggleMaximize(name: string) {
    const state = getWindowState(name)
    if (state.isMaximized) {
      state.x = state._prevX ?? 50
      state.y = state._prevY ?? 50
      state.width = state._prevWidth ?? Math.floor(window.innerWidth * 0.9)
      state.height = state._prevHeight ?? Math.floor(window.innerHeight * 0.86)
      state.isMaximized = false
    } else {
      state._prevX = state.x
      state._prevY = state.y
      state._prevWidth = state.width
      state._prevHeight = state.height
      state.x = 0
      state.y = 0
      state.width = window.innerWidth
      state.height = window.innerHeight
      state.isMaximized = true
    }
  }

  function updatePosition(name: string, x: number, y: number) {
    const state = getWindowState(name)
    state.x = x
    state.y = y
  }

  return {
    activeWindow,
    windowList,
    windows,
    getWindowState,
    openWindow,
    closeWindow,
    minimizeWindow,
    toggleWindow,
    toggleMaximize,
    updatePosition,
  }
})
