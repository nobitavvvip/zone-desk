import http from './http'

export interface Container {
  id: string
  name: string
  image: string
  status: string
  state: string
  ports: string
  created: string
}

export interface ContainerStats {
  cpuPercent: number
  memoryUsageMB: number
  memoryLimitMB: number
  networkRxMB: number
  networkTxMB: number
  blockRxMB: number
  blockWxMB: number
  pids: number
}

export interface ComposeProject {
  name: string
  status: string
  configFiles: string
  workingDir: string
}

export interface DockerImage {
  id: string
  repository: string
  tag: string
  size: string
  created: string
}

export async function listContainers(all: boolean = false) {
  const res = await http.get('/containers/list', { params: { all } })
  return res.data.data as Container[]
}

export async function startContainer(id: string) {
  await http.post('/containers/start', { id })
}

export async function stopContainer(id: string) {
  await http.post('/containers/stop', { id })
}

export async function restartContainer(id: string) {
  await http.post('/containers/restart', { id })
}

export async function removeContainer(id: string, force: boolean = false) {
  await http.post('/containers/remove', { id, force })
}

export async function getContainerLogs(id: string, tail: number = 100) {
  const res = await http.get('/containers/logs', { params: { id, tail } })
  return res.data.data.logs as string
}

export async function getContainerStats(id: string) {
  const res = await http.get('/containers/stats', { params: { id } })
  return res.data.data as ContainerStats
}

export async function pruneContainers() {
  await http.post('/containers/prune')
}

export async function execContainer(id: string, cmd: string[], detach: boolean = false) {
  const res = await http.post('/containers/exec', { id, cmd, detach, interactive: !detach })
  return res.data.data.output as string
}

export async function createContainer(params: {
  name?: string
  image: string
  ports?: string[]
  env?: string[]
  volumes?: string[]
  network?: string
  entrypoint?: string
  cmd?: string[]
  restart?: string
}) {
  const res = await http.post('/containers/create', params)
  return res.data.data.containerId as string
}

export async function listImages() {
  const res = await http.get('/containers/images')
  const raw = res.data.data.output as string
  const lines = raw.split('\n').filter(Boolean)
  const images: DockerImage[] = []
  for (const line of lines) {
    try {
      const item = JSON.parse(line)
      const repoTag = (item.Repository || '') + ':' + (item.Tag || '')
      images.push({
        id: item.ID || '',
        repository: item.Repository || '',
        tag: item.Tag || '',
        size: item.Size || '',
        created: item.CreatedAt || '',
      })
    } catch {}
  }
  return images
}

export async function pullImage(image: string) {
  await http.post('/containers/pull', { image })
}

export async function removeImage(image: string, force: boolean = false) {
  await http.post('/containers/rmi', { image, force })
}

export async function listCompose() {
  const res = await http.get('/compose/list')
  return res.data.data as ComposeProject[]
}

export async function composeUp(projectDir: string, file: string, detached: boolean = true) {
  await http.post('/compose/up', { projectDir, file, detached })
}

export async function composeDown(projectDir: string, file: string) {
  await http.post('/compose/down', { projectDir, file })
}

export async function composeStart(projectDir: string, file: string) {
  await http.post('/compose/start', { projectDir, file })
}

export async function composeStop(projectDir: string, file: string) {
  await http.post('/compose/stop', { projectDir, file })
}

export async function composeRestart(projectDir: string, file: string) {
  await http.post('/compose/restart', { projectDir, file })
}

export async function getComposeLogs(projectDir: string, file: string, tail: number = 100) {
  const res = await http.get('/compose/logs', { params: { projectDir, file, tail } })
  return res.data.data.logs as string
}

export async function createCompose(projectDir: string, filename: string, content: string, start: boolean = true) {
  await http.post('/compose/create', { projectDir, filename, content, start })
}

export async function readComposeFile(projectDir: string, file: string) {
  const res = await http.get('/compose/read', { params: { projectDir, file } })
  return res.data.data as { content: string; projectDir: string; filename: string }
}

export async function updateComposeFile(projectDir: string, filename: string, content: string) {
  await http.put('/compose/update', { projectDir, filename, content })
}

export async function deleteComposeProject(projectDir: string, filename: string) {
  await http.post('/compose/delete', { projectDir, filename })
}
