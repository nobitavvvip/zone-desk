import http from './http'

export interface Shortcut {
  id: string
  name: string
  path: string
}

export async function listShortcuts() {
  const res = await http.get('/shortcuts')
  return res.data.data as Shortcut[]
}

export async function addShortcut(name: string, path: string) {
  const res = await http.post('/shortcuts', { name, path })
  return res.data.data as Shortcut
}

export async function deleteShortcut(id: string) {
  await http.delete(`/shortcuts/${id}`)
}
