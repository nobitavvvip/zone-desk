import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as api from '@/api/container'
import type { Container, ComposeProject, DockerImage } from '@/api/container'

export type TabName = 'containers' | 'compose' | 'images'

export const useContainerStore = defineStore('container', () => {
  const containers = ref<Container[]>([])
  const composeProjects = ref<ComposeProject[]>([])
  const images = ref<DockerImage[]>([])
  const loading = ref(false)
  const activeTab = ref<TabName>('containers')
  const containerLogs = ref('')
  const composeLogs = ref('')
  const showAll = ref(true)
  const searchQuery = ref('')
  const selectedIds = ref<Set<string>>(new Set())
  const detailContainer = ref<Container | null>(null)
  const detailInfo = ref('')
  const error = ref('')

  const filteredContainers = computed(() => {
    let list = containers.value
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase()
      list = list.filter(c =>
        c.name.toLowerCase().includes(q) ||
        c.image.toLowerCase().includes(q) ||
        c.id.toLowerCase().includes(q)
      )
    }
    return list
  })

  const runningContainers = computed(() =>
    containers.value.filter(c => c.state === 'running')
  )

  const hasSelection = computed(() => selectedIds.value.size > 0)

  async function fetchContainers() {
    loading.value = true
    try {
      containers.value = await api.listContainers(showAll.value)
      selectedIds.value = new Set()
    } finally {
      loading.value = false
    }
  }

  async function fetchCompose() {
    loading.value = true
    error.value = ''
    try {
      composeProjects.value = await api.listCompose()
    } catch (e: any) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  async function fetchImages() {
    loading.value = true
    try {
      images.value = await api.listImages()
    } finally {
      loading.value = false
    }
  }

  async function start(id: string) {
    await api.startContainer(id)
    await fetchContainers()
  }

  async function stop(id: string) {
    await api.stopContainer(id)
    await fetchContainers()
  }

  async function restart(id: string) {
    await api.restartContainer(id)
    await fetchContainers()
  }

  async function remove(id: string, force: boolean = false) {
    await api.removeContainer(id, force)
    await fetchContainers()
  }

  async function batchStart() {
    for (const id of selectedIds.value) {
      await api.startContainer(id)
    }
    await fetchContainers()
  }

  async function batchStop() {
    for (const id of selectedIds.value) {
      await api.stopContainer(id)
    }
    await fetchContainers()
  }

  async function batchRemove() {
    for (const id of selectedIds.value) {
      await api.removeContainer(id, true)
    }
    await fetchContainers()
  }

  async function fetchLogs(id: string, tail: number = 100) {
    containerLogs.value = await api.getContainerLogs(id, tail)
  }

  async function fetchDetail(id: string) {
    const c = containers.value.find(c => c.id === id)
    if (c) detailContainer.value = c
    try {
      detailInfo.value = await api.execContainer(id, ['cat', '/etc/hostname'], false)
    } catch {
      detailInfo.value = '(unavailable)'
    }
  }

  async function pullImage(image: string) {
    await api.pullImage(image)
    await fetchImages()
  }

  async function removeImage(image: string, force: boolean = false) {
    await api.removeImage(image, force)
    await fetchImages()
  }

  async function prune() {
    await api.pruneContainers()
    await fetchContainers()
  }

  async function composeUp(projectDir: string, file: string) {
    await api.composeUp(projectDir, file, true)
    await fetchCompose()
  }

  async function composeDown(projectDir: string, file: string) {
    await api.composeDown(projectDir, file)
    await fetchCompose()
  }

  async function composeStart(projectDir: string, file: string) {
    await api.composeStart(projectDir, file)
    await fetchCompose()
  }

  async function composeStop(projectDir: string, file: string) {
    await api.composeStop(projectDir, file)
    await fetchCompose()
  }

  async function composeRestart(projectDir: string, file: string) {
    await api.composeRestart(projectDir, file)
    await fetchCompose()
  }

  async function fetchComposeLogs(projectDir: string, file: string, tail: number = 100) {
    composeLogs.value = await api.getComposeLogs(projectDir, file, tail)
  }

  async function readComposeFile(projectDir: string, file: string) {
    return await api.readComposeFile(projectDir, file)
  }

  async function updateComposeFile(projectDir: string, filename: string, content: string) {
    await api.updateComposeFile(projectDir, filename, content)
    await fetchCompose()
  }

  async function deleteComposeProject(projectDir: string, filename: string) {
    await api.deleteComposeProject(projectDir, filename)
    await fetchCompose()
  }

  function setActiveTab(tab: TabName) {
    activeTab.value = tab
  }

  function toggleShowAll() {
    showAll.value = !showAll.value
  }

  function toggleSelect(id: string) {
    const s = new Set(selectedIds.value)
    if (s.has(id)) s.delete(id)
    else s.add(id)
    selectedIds.value = s
  }

  function selectAll() {
    selectedIds.value = new Set(containers.value.map(c => c.id))
  }

  function clearSelection() {
    selectedIds.value = new Set()
  }

  return {
    containers,
    composeProjects,
    images,
    loading,
    activeTab,
    containerLogs,
    composeLogs,
    showAll,
    searchQuery,
    selectedIds,
    detailContainer,
    detailInfo,
    filteredContainers,
    runningContainers,
    hasSelection,
    error,
    fetchContainers,
    fetchCompose,
    fetchImages,
    start,
    stop,
    restart,
    remove,
    batchStart,
    batchStop,
    batchRemove,
    fetchLogs,
    fetchDetail,
    pullImage,
    removeImage,
    prune,
    composeUp,
    composeDown,
    composeStart,
    composeStop,
    composeRestart,
    fetchComposeLogs,
    readComposeFile,
    updateComposeFile,
    deleteComposeProject,
    setActiveTab,
    toggleShowAll,
    toggleSelect,
    selectAll,
    clearSelection,
  }
})
