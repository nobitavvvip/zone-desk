import http from './http'

export interface NotesRootResult {
  rootDir: string
}

export interface NoteFileItem {
  name: string
  path: string
  type: 'file' | 'dir'
  size: number
  mode: string
  modTime: string
  hidden: boolean
}

export interface NotesListResult {
  path: string
  parent: string | null
  items: NoteFileItem[]
}

export interface NoteReadResult {
  path: string
  content: string
  encoding: string
}

export async function getNotesRoot() {
  const res = await http.get('/notes/root')
  return res.data.data as NotesRootResult
}

export async function listNotes(path: string) {
  const res = await http.get('/notes/list', { params: { path } })
  return res.data.data as NotesListResult
}

export async function readNote(path: string) {
  const res = await http.get('/notes/read', { params: { path } })
  return res.data.data as NoteReadResult
}

export async function writeNote(path: string, content: string) {
  const res = await http.post('/notes/write', { path, content })
  return res.data.data
}

export async function createNote(dirPath: string, name: string) {
  const res = await http.post('/notes/create', { dirPath, name })
  return res.data.data as { path: string; name: string }
}

export async function createNoteFolder(dirPath: string, name: string) {
  const res = await http.post('/notes/mkdir', { dirPath, name })
  return res.data.data as { path: string }
}

export async function deleteNote(path: string) {
  const res = await http.post('/notes/delete', { path })
  return res.data.data as { path: string }
}

export async function renameNote(oldPath: string, newPath: string) {
  const res = await http.post('/notes/rename', { oldPath, newPath })
  return res.data.data as { oldPath: string; newPath: string }
}

export async function moveNote(source: string, destination: string) {
  const res = await http.post('/notes/move', { source, destination })
  return res.data.data as { source: string; destination: string }
}
