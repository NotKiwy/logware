import { h } from '../ui/dom'
import { lotBlob } from '../ui/lottie'

type _Obj = Record<string, unknown>

const __mimes__: Record<string, string> = {
  photo: 'image/jpeg',
  sticker: 'image/webp',
  video: 'video/mp4',
  round: 'video/mp4',
  anim: 'video/mp4',
  voice: 'audio/ogg',
  audio: 'audio/mpeg',
  doc: 'application/octet-stream',
}

const __fields__: Record<string, string> = {
  video: 'video',
  round: 'video_note',
  anim: 'animation',
  voice: 'voice',
  audio: 'audio',
  doc: 'document',
  sticker: 'sticker',
}

const __labels__: Record<string, string> = {
  photo: '🖼 Фото',
  video: '🎬 Видео',
  round: '⭕️ Кружок',
  anim: '🎞 GIF',
  voice: '🎙 Голосовое',
  audio: '🎵 Аудио',
  doc: '📎 Файл',
  sticker: '🌟 Стикер',
  rich: '📰 Rich-сообщение',
  other: '📦 Сообщение другого типа',
  empty: '👻 Пустое или неподдерживаемое сообщение',
}

const flat = (__v: unknown, __out: string[]): void => {
  if (Array.isArray(__v)) return void __v.forEach((_x) => flat(_x, __out))
  if (!__v || typeof __v !== 'object') return
  for (const [_k, _x] of Object.entries(__v as _Obj)) {
    if (_k === 'text' && typeof _x === 'string') __out.push(_x)
    else flat(_x, __out)
  }
}

export const text = (__m: _Obj): string => {
  if (typeof __m.text === 'string') return __m.text
  if (typeof __m.caption === 'string') return __m.caption
  if (__m.rich_message) {
    const _out: string[] = []
    flat((__m.rich_message as _Obj).blocks, _out)
    return _out.join('\n')
  }
  return ''
}

export const chat = (__m: _Obj): string => {
  const _c = (__m.chat ?? {}) as _Obj
  return [_c.first_name, _c.last_name].filter(Boolean).join(' ') || (typeof _c.title === 'string' && _c.title) || (_c.username ? '@' + _c.username : 'Чат')
}

export const who = (__m: _Obj, __out: boolean): string => {
  if (__out) return 'Ты'
  const _f = (__m.from ?? {}) as _Obj
  return [_f.first_name, _f.last_name].filter(Boolean).join(' ') || (_f.username ? '@' + _f.username : 'Собеседник')
}

export const sticker = (__m: _Obj): { anim: boolean; vid: boolean; id: string } | null => {
  const _s = __m.sticker as _Obj | undefined
  if (!_s) return null
  return { anim: _s.is_animated === true, vid: _s.is_video === true, id: typeof _s.file_id === 'string' ? _s.file_id : '' }
}

export const mime = (__m: _Obj, __kind: string): string => {
  const _st = sticker(__m)
  if (_st?.anim) return 'application/x-tgsticker'
  if (_st?.vid) return 'video/webm'
  const _f = __m[__fields__[__kind] ?? ''] as _Obj | undefined
  return (typeof _f?.mime_type === 'string' && _f.mime_type) || __mimes__[__kind] || 'application/octet-stream'
}

const __exts__: Record<string, string> = { photo: 'jpg', sticker: 'webp', video: 'mp4', round: 'mp4', anim: 'mp4', voice: 'ogg', audio: 'mp3' }

export const fname = (__m: _Obj, __kind = 'doc'): string => {
  const _st = sticker(__m)
  if (_st) return 'sticker.' + (_st.anim ? 'tgs' : _st.vid ? 'webm' : 'webp')
  const _d = (__m.document ?? __m.audio) as _Obj | undefined
  return (typeof _d?.file_name === 'string' && _d.file_name) || (__kind + '.' + (__exts__[__kind] ?? 'bin'))
}

export const player = (__kind: string, __url: string, __name: string, __blob?: Blob): HTMLElement => {
  if (__kind === 'sticker' && __name.endsWith('.tgs') && __blob) return lotBlob(__blob, 200)
  if (__kind === 'sticker' && __name.endsWith('.webm')) return h('video', { src: __url, class: 'media sticker', autoplay: true, loop: true, muted: true, playsinline: true })
  if (__kind === 'sticker') return h('img', { src: __url, class: 'media sticker' })
  if (__kind === 'photo') return h('img', { src: __url, class: 'media' })
  if (__kind === 'video' || __kind === 'round') return h('video', { src: __url, class: 'media' + (__kind === 'round' ? ' round' : ''), controls: true, playsinline: true })
  if (__kind === 'anim') return h('video', { src: __url, class: 'media', autoplay: true, loop: true, muted: true, playsinline: true })
  if (__kind === 'voice' || __kind === 'audio') return h('audio', { src: __url, controls: true })
  return h('a', { href: __url, download: __name, class: 'file' }, '📎 ' + __name)
}

export const label = (__kind: string): string => __labels__[__kind] ?? ''

export const when = (__unix: number): string =>
  new Date(__unix * 1000).toLocaleString('ru', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })

export const clock = (__unix: number): string => {
  const _d = new Date(__unix * 1000)
  const _t = _d.toLocaleTimeString('ru', { hour: '2-digit', minute: '2-digit' })
  return _d.toDateString() === new Date().toDateString() ? _t : _d.toLocaleDateString('ru', { day: 'numeric', month: 'short' }) + ', ' + _t
}

export const many = (__n: number, __forms: [string, string, string]): string => {
  const _a = __n % 10
  const _b = __n % 100
  return `${__n} ${__forms[_a === 1 && _b !== 11 ? 0 : _a >= 2 && _a <= 4 && (_b < 12 || _b > 14) ? 1 : 2]}`
}

const __plain__: Record<string, string> = {
  photo: 'Фото',
  video: 'Видео',
  round: 'Кружок',
  anim: 'GIF',
  voice: 'Голосовое',
  audio: 'Аудио',
  doc: 'Файл',
  sticker: 'Стикер',
  rich: 'Сообщение',
  other: 'Сообщение',
  empty: 'Сообщение',
}

export const plain = (__kind: string): string => __plain__[__kind] ?? 'Сообщение'
