import http from './http'

export interface DesktopSettings {
  wallpaper: string
  blur: number
  mask: number
  accent: string
}

export async function getSettings() {
  const res = await http.get('/settings')
  return res.data.data as DesktopSettings
}

export async function saveSettings(settings: DesktopSettings) {
  const res = await http.put('/settings', settings)
  return res.data.data as DesktopSettings
}
