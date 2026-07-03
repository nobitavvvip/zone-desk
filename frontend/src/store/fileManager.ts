import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { listFiles, type FileListResult, type FileItem } from '@/api/file'

export type ViewMode = 'details' | 'list' | 'icons' | 'tile'
export type SortBy = 'name' | 'size' | 'type' | 'modified'
export type SortOrder = 'asc' | 'desc'

export const useFileManagerStore = defineStore('fileManager', () => {
  const currentPath = ref('/')
  const currentDir = ref<FileListResult | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const history = ref<string[]>(['/'])
  const historyIndex = ref(0)
  const viewMode = ref<ViewMode>('details')
  const selectedItems = ref<Set<string>>(new Set())
  const clipboard = ref<{ files: FileItem[], action: 'copy' | 'cut' } | null>(null)
  const sortBy = ref<SortBy>('name')
  const sortOrder = ref<SortOrder>('asc')

  const sortedItems = computed(() => {
    if (!currentDir.value?.items) return []
    
    const items = [...currentDir.value.items]
    
    items.sort((a, b) => {
      // 目录始终在文件前面
      if (a.type !== b.type) {
        if (a.type === 'dir') return -1
        if (b.type === 'dir') return 1
      }

      let result: number
      switch (sortBy.value) {
        case 'size':
          result = a.size - b.size
          break
        case 'type':
          result = a.type.localeCompare(b.type)
          break
        case 'modified':
          result = new Date(a.modTime).getTime() - new Date(b.modTime).getTime()
          break
        default: // 'name'
          result = a.name.localeCompare(b.name)
      }

      return sortOrder.value === 'desc' ? -result : result
    })

    return items
  })

  async function navigate(path: string) {
    loading.value = true
    error.value = null
    currentDir.value = null  // 先清空，触发加载动画
    const startTime = Date.now()
    try {
      const result = await listFiles(path)
      // 最小加载时间 150ms，避免快速请求时闪烁
      const elapsed = Date.now() - startTime
      if (elapsed < 150) {
        await new Promise(resolve => setTimeout(resolve, 150 - elapsed))
      }
      currentPath.value = result.path
      currentDir.value = result
      clearSelection()
      updateUrlHash(result.path)
      return result
    } catch (e: any) {
      error.value = e.message
      return null
    } finally {
      loading.value = false
    }
  }

  function updateUrlHash(path: string) {
    if (path === '/') {
      window.history.replaceState(null, '', window.location.pathname)
    } else {
      window.history.replaceState(null, '', `#/filemanager?path=${encodeURIComponent(path)}`)
    }
  }

  function getPathFromUrl(): string {
    const hash = window.location.hash
    if (hash.startsWith('#/filemanager')) {
      const queryStart = hash.indexOf('?')
      if (queryStart === -1) return '/'
      const params = new URLSearchParams(hash.slice(queryStart + 1))
      return params.get('path') || '/'
    }
    return '/'
  }

  async function goTo(path: string) {
    const result = await navigate(path)
    if (result) {
      history.value = history.value.slice(0, historyIndex.value + 1)
      history.value.push(path)
      historyIndex.value = history.value.length - 1
    }
    return result
  }

  async function goBack() {
    if (historyIndex.value > 0) {
      historyIndex.value--
      await navigate(history.value[historyIndex.value])
    }
  }

  async function goForward() {
    if (historyIndex.value < history.value.length - 1) {
      historyIndex.value++
      await navigate(history.value[historyIndex.value])
    }
  }

  async function goUp() {
    if (currentDir.value?.parent) {
      await goTo(currentDir.value.parent)
    }
  }

  async function refresh() {
    await navigate(currentPath.value)
  }

  function selectItem(path: string, ctrl = false, shift = false) {
    if (shift && selectedItems.value.size > 0) {
      // 范围选择
      const items = sortedItems.value
      const lastSelected = Array.from(selectedItems.value).pop()
      const lastIndex = items.findIndex(item => item.path === lastSelected)
      const currentIndex = items.findIndex(item => item.path === path)
      
      const start = Math.min(lastIndex, currentIndex)
      const end = Math.max(lastIndex, currentIndex)
      
      for (let i = start; i <= end; i++) {
        selectedItems.value.add(items[i].path)
      }
    } else if (ctrl) {
      // 追加选择
      if (selectedItems.value.has(path)) {
        selectedItems.value.delete(path)
      } else {
        selectedItems.value.add(path)
      }
    } else {
      // 单选
      selectedItems.value.clear()
      selectedItems.value.add(path)
    }
  }

  function selectAll() {
    sortedItems.value.forEach(item => {
      selectedItems.value.add(item.path)
    })
  }

  function clearSelection() {
    selectedItems.value.clear()
  }

  function copy() {
    const files = sortedItems.value.filter(item => selectedItems.value.has(item.path))
    if (files.length > 0) {
      clipboard.value = { files, action: 'copy' }
    }
  }

  function cut() {
    const files = sortedItems.value.filter(item => selectedItems.value.has(item.path))
    if (files.length > 0) {
      clipboard.value = { files, action: 'cut' }
    }
  }

  function setSortBy(field: SortBy) {
    if (sortBy.value === field) {
      sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
    } else {
      sortBy.value = field
      sortOrder.value = 'asc'
    }
  }

  function setSortOrder(order: SortOrder) {
    sortOrder.value = order
  }

  return {
    currentPath,
    currentDir,
    loading,
    error,
    history,
    historyIndex,
    viewMode,
    selectedItems,
    clipboard,
    sortBy,
    sortOrder,
    sortedItems,
    navigate,
    goTo,
    goBack,
    goForward,
    goUp,
    refresh,
    selectItem,
    selectAll,
    clearSelection,
    copy,
    cut,
    setSortBy,
    setSortOrder,
    setViewMode: (mode: ViewMode) => { viewMode.value = mode },
    getPathFromUrl,
  }
})
