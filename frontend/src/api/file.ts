import http from './http'

export interface FileItem {
  name: string
  path: string
  type: 'file' | 'dir' | 'symlink' | 'other'
  size: number
  mode: string
  modTime: string
  hidden: boolean
}

export interface FileListResult {
  path: string
  parent: string | null
  items: FileItem[]
}

export async function listFiles(path: string) {
  const res = await http.get('/files/list', { params: { path } })
  return res.data.data as FileListResult
}

export async function statFile(path: string) {
  const res = await http.get('/files/stat', { params: { path } })
  return res.data.data as FileItem
}

export async function readFile(path: string) {
  const res = await http.get('/files/read', { params: { path } })
  return res.data.data as { path: string; content: string; encoding: string }
}

export function getDownloadUrl(path: string) {
  return `/api/files/download?path=${encodeURIComponent(path)}`
}

export async function uploadFile(path: string, file: File) {
  const form = new FormData()
  form.append('path', path)
  form.append('file', file)
  const res = await http.post('/files/upload', form)
  return res.data.data as { path: string }
}

export async function mkdir(path: string) {
  const res = await http.post('/files/mkdir', { path })
  return res.data.data
}

export async function renameFile(oldPath: string, newPath: string) {
  const res = await http.post('/files/rename', { oldPath, newPath })
  return res.data.data
}

export async function deleteFile(path: string) {
  const res = await http.post('/files/delete', { path })
  return res.data.data
}

export async function copyFile(source: string, destination: string): Promise<void> {
  const res = await http.post('/files/copy', { source, destination })
  if (!res.data.success) {
    throw new Error(res.data.message || '复制失败')
  }
}

export async function moveFile(source: string, destination: string): Promise<void> {
  const res = await http.post('/files/move', { source, destination })
  if (!res.data.success) {
    throw new Error(res.data.message || '移动失败')
  }
}
