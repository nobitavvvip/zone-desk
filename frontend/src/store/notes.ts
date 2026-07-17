import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  getNotesRoot,
  listNotes,
  readNote,
  writeNote,
  createNote as apiCreateNote,
  createNoteFolder,
  deleteNote,
  renameNote,
  moveNote,
} from '@/api/notes'
import type { NoteFileItem } from '@/api/notes'

const STORAGE_KEY_AUTO_SAVE = 'notes_auto_save'
const STORAGE_KEY_EXPANDED_DIRS = 'notes_expanded_dirs'
const DEBOUNCE_DELAY = 1000

export interface OpenFile {
  name: string
  path: string
  content: string
  modified: boolean
}

export interface TreeNode {
  name: string
  path: string
  type: 'file' | 'dir'
  children: TreeNode[]
  expanded: boolean
  loaded: boolean
}

export const useNotesStore = defineStore('notes', () => {
  const rootDir = ref('/')
  const currentDir = ref('/')
  const fileList = ref<NoteFileItem[]>([])
  const openFiles = ref<OpenFile[]>([])
  const activeFile = ref<OpenFile | null>(null)
  const loading = ref(false)
  const showPreview = ref(false)
  const sidebarCollapsed = ref(false)
  const editorPadding = ref(40)
  const autoSave = ref(localStorage.getItem(STORAGE_KEY_AUTO_SAVE) !== 'false')

  const treeNodes = ref<TreeNode[]>([])
  const expandedDirs = ref<Set<string>>(loadExpandedDirs())

  let debounceTimer: ReturnType<typeof setTimeout> | null = null
  let pendingSavePath: string | null = null

  const mdFiles = computed(() =>
    fileList.value.filter(f => f.name.endsWith('.md'))
  )

  function loadExpandedDirs(): Set<string> {
    try {
      const data = localStorage.getItem(STORAGE_KEY_EXPANDED_DIRS)
      if (data) {
        const arr: string[] = JSON.parse(data)
        return new Set(arr)
      }
    } catch {}
    return new Set()
  }

  function saveExpandedDirs() {
    localStorage.setItem(
      STORAGE_KEY_EXPANDED_DIRS,
      JSON.stringify(Array.from(expandedDirs.value))
    )
  }

  function setAutoSave(value: boolean) {
    autoSave.value = value
    localStorage.setItem(STORAGE_KEY_AUTO_SAVE, String(value))
  }

  function clearDebounceTimer() {
    if (debounceTimer) {
      clearTimeout(debounceTimer)
      debounceTimer = null
    }
    pendingSavePath = null
  }

  function scheduleDebounceSave(path: string) {
    if (!autoSave.value) return

    clearDebounceTimer()
    pendingSavePath = path
    debounceTimer = setTimeout(() => {
      const file = openFiles.value.find(f => f.path === path)
      if (file && file.modified) {
        saveFile(file)
      }
      debounceTimer = null
      pendingSavePath = null
    }, DEBOUNCE_DELAY)
  }

  function buildTreeNodes(items: NoteFileItem[], parentPath: string): TreeNode[] {
    const dirs: TreeNode[] = []
    const files: TreeNode[] = []

    for (const item of items) {
      if (item.type === 'dir') {
        dirs.push({
          name: item.name,
          path: item.path,
          type: 'dir',
          children: [],
          expanded: expandedDirs.value.has(item.path),
          loaded: false,
        })
      } else if (item.name.endsWith('.md')) {
        files.push({
          name: item.name,
          path: item.path,
          type: 'file',
          children: [],
          expanded: false,
          loaded: true,
        })
      }
    }

    dirs.sort((a, b) => a.name.localeCompare(b.name))
    files.sort((a, b) => a.name.localeCompare(b.name))

    return [...dirs, ...files]
  }

  async function loadTreeRoot() {
    try {
      const result = await listNotes(rootDir.value)
      treeNodes.value = buildTreeNodes(result.items || [], rootDir.value)

      for (const node of treeNodes.value) {
        if (node.type === 'dir' && node.expanded) {
          await loadTreeChildren(node)
        }
      }
    } catch {}
  }

  async function loadTreeChildren(node: TreeNode) {
    try {
      const result = await listNotes(node.path)
      node.children = buildTreeNodes(result.items || [], node.path)
      node.loaded = true

      for (const child of node.children) {
        if (child.type === 'dir' && child.expanded) {
          await loadTreeChildren(child)
        }
      }
    } catch {}
  }

  async function toggleDirExpand(node: TreeNode) {
    if (node.expanded) {
      node.expanded = false
      expandedDirs.value.delete(node.path)
    } else {
      if (!node.loaded) {
        await loadTreeChildren(node)
      }
      node.expanded = true
      expandedDirs.value.add(node.path)
    }
    saveExpandedDirs()
  }

  async function refreshTree() {
    const refreshNode = async (nodes: TreeNode[]) => {
      for (const node of nodes) {
        if (node.type === 'dir' && node.expanded) {
          await loadTreeChildren(node)
          await refreshNode(node.children)
        }
      }
    }
    await refreshNode(treeNodes.value)
  }

  async function init() {
    try {
      const config = await getNotesRoot()
      rootDir.value = config.rootDir
      await navigate(config.rootDir)
      await loadTreeRoot()
    } catch (e: any) {
      rootDir.value = '/'
      await navigate('/')
    }
  }

  async function navigate(path: string): Promise<string> {
    loading.value = true
    try {
      const result = await listNotes(path)
      currentDir.value = result.path
      fileList.value = (result.items || []).filter(
        f => f.type === 'dir' || f.name.endsWith('.md')
      )
      return ''
    } catch (e: any) {
      return e.message || 'failed to load directory'
    } finally {
      loading.value = false
    }
  }

  async function openFile(file: NoteFileItem | TreeNode): Promise<string> {
    if (file.type === 'dir') {
      return ''
    }

    const existing = openFiles.value.find(f => f.path === file.path)
    if (existing) {
      activeFile.value = existing
      return ''
    }

    try {
      const result = await readNote(file.path)
      const openFile: OpenFile = {
        name: file.name,
        path: file.path,
        content: result.content,
        modified: false,
      }
      openFiles.value.push(openFile)
      activeFile.value = openFile
      return ''
    } catch (e: any) {
      return e.message || 'failed to open file'
    }
  }

  async function saveFile(file: OpenFile): Promise<string> {
    try {
      await writeNote(file.path, file.content)
      file.modified = false
      return ''
    } catch (e: any) {
      return e.message || 'failed to save file'
    }
  }

  async function saveCurrentFile(): Promise<string> {
    if (activeFile.value && activeFile.value.modified) {
      clearDebounceTimer()
      return saveFile(activeFile.value)
    }
    return ''
  }

  function updateContent(path: string, content: string) {
    const file = openFiles.value.find(f => f.path === path)
    if (file) {
      file.content = content
      file.modified = true
      scheduleDebounceSave(path)
    }
  }

  function closeFile(path: string) {
    if (pendingSavePath === path) {
      const file = openFiles.value.find(f => f.path === path)
      if (file && file.modified) {
        saveFile(file)
      }
      clearDebounceTimer()
    }

    const idx = openFiles.value.findIndex(f => f.path === path)
    if (idx === -1) return

    openFiles.value.splice(idx, 1)

    if (activeFile.value?.path === path) {
      activeFile.value = openFiles.value[Math.min(idx, openFiles.value.length - 1)] || null
    }
  }

  async function createNote(dirPath: string, name: string): Promise<string> {
    try {
      await apiCreateNote(dirPath, name)
      await navigate(currentDir.value)
      await loadTreeRoot()
      return ''
    } catch (e: any) {
      return e.message || 'failed to create note'
    }
  }

  async function createFolder(dirPath: string, name: string): Promise<string> {
    try {
      await createNoteFolder(dirPath, name)
      await navigate(currentDir.value)
      await loadTreeRoot()
      return ''
    } catch (e: any) {
      return e.message || 'failed to create folder'
    }
  }

  async function deleteItem(path: string): Promise<string> {
    try {
      await deleteNote(path)
      closeFile(path)
      await navigate(currentDir.value)
      await loadTreeRoot()
      return ''
    } catch (e: any) {
      return e.message || 'failed to delete'
    }
  }

  async function renameItem(oldPath: string, newName: string): Promise<string> {
    try {
      const dir = oldPath.substring(0, oldPath.lastIndexOf('/'))
      const newPath = `${dir}/${newName}`
      await renameNote(oldPath, newPath)

      const openIdx = openFiles.value.findIndex(f => f.path === oldPath)
      if (openIdx !== -1) {
        openFiles.value[openIdx].path = newPath
        openFiles.value[openIdx].name = newName
        if (activeFile.value?.path === oldPath) {
          activeFile.value.path = newPath
          activeFile.value.name = newName
        }
      }

      await navigate(currentDir.value)
      await loadTreeRoot()
      return ''
    } catch (e: any) {
      return e.message || 'failed to rename'
    }
  }

  return {
    rootDir,
    currentDir,
    fileList,
    openFiles,
    activeFile,
    loading,
    showPreview,
    sidebarCollapsed,
    editorPadding,
    autoSave,
    mdFiles,
    treeNodes,
    expandedDirs,
    init,
    navigate,
    openFile,
    saveFile,
    saveCurrentFile,
    updateContent,
    closeFile,
    createNote,
    createFolder,
    deleteItem,
    renameItem,
    setAutoSave,
    toggleDirExpand,
    refreshTree,
    loadTreeRoot,
  }
})
