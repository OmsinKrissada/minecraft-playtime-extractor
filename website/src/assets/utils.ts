export const BAR_COLORS = [
  '#5D9E3F',
  '#4DD0E1',
  '#F5C518',
  '#C0C0C0',
  '#E8A97B',
  '#A78BFA',
  '#F9A8D4',
  '#86EFAC',
]

export function timeAgo(iso: string, forceRelative = false) {
  const diff = (Date.now().valueOf() - new Date(iso).valueOf()) / 1000
  if (diff < 60) return 'just now'
  if (diff < 3600) return Math.floor(diff / 60) + 'm ago'
  if (diff < 86400) return Math.floor(diff / 3600) + 'h ago'
  if (forceRelative || diff < 604800) return Math.floor(diff / 86400) + 'd ago'
  return new Date(iso).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

export function isOnline(iso: string) {
  return Date.now() - new Date(iso).valueOf() < 10 * 60 * 1000
}

export function fmtValue(h: number) {
  return h >= 1 ? h.toFixed(1) : (h * 60).toFixed(0)
}

export function fmtUnit(h: number) {
  return h >= 1 ? 'HOURS' : 'MINS'
}

export function avatarUrl(uuid: string) {
  return `https://mc-heads.net/avatar/${uuid}`
}
