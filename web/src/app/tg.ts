type _Cb = (_e: unknown, _v?: unknown) => void

type _Store = {
  getItem: (_k: string, _cb: _Cb) => void
  setItem: (_k: string, _v: string, _cb?: _Cb) => void
  removeItem: (_k: string, _cb?: _Cb) => void
}

type _Tg = {
  initData: string
  initDataUnsafe?: { user?: { id?: number; first_name?: string; last_name?: string; username?: string; photo_url?: string } }
  platform: string
  colorScheme?: 'light' | 'dark'
  themeParams?: Record<string, string>
  ready: () => void
  expand: () => void
  isVersionAtLeast: (_v: string) => boolean
  setHeaderColor?: (_c: string) => void
  setBackgroundColor?: (_c: string) => void
  setBottomBarColor?: (_c: string) => void
  onEvent?: (_e: string, _cb: () => void) => void
  HapticFeedback?: {
    impactOccurred: (_s: 'light' | 'medium' | 'soft') => void
    notificationOccurred: (_t: 'success' | 'error' | 'warning') => void
    selectionChanged: () => void
  }
  CloudStorage?: _Store
  SecureStorage?: _Store
  openTelegramLink?: (_url: string) => void
  openLink?: (_url: string) => void
  openInvoice?: (_url: string, _cb: (_s: 'paid' | 'cancelled' | 'failed' | 'pending') => void) => void
  BackButton?: { show: () => void; hide: () => void; onClick: (_cb: () => void) => void; offClick: (_cb: () => void) => void }
}

export const __tg__: _Tg | undefined = (window as unknown as { Telegram?: { WebApp?: _Tg } }).Telegram?.WebApp

export const init = (): string => __tg__?.initData ?? ''

export const inside = (): boolean => init() !== ''

const ver = (__v: string): boolean => !!__tg__?.isVersionAtLeast?.(__v)

export const canCloud = (): boolean => ver('6.9') && !!__tg__?.CloudStorage

export const canSecure = (): boolean =>
  ver('9.0') && !!__tg__?.SecureStorage && ['ios', 'android'].includes(__tg__?.platform ?? '')

export type Box = {
  get: (_k: string) => Promise<string | null>
  set: (_k: string, _v: string) => Promise<void>
  del: (_k: string) => Promise<void>
}

const wrap = (__s: _Store | undefined): Box => ({
  get: (__k) =>
    new Promise((_ok, _no) => {
      if (!__s) return _no(new Error('storage unavailable'))
      __s.getItem(__k, (_e, _v) => (_e ? _no(_e) : _ok(typeof _v === 'string' && _v !== '' ? _v : null)))
    }),
  set: (__k, __v) =>
    new Promise((_ok, _no) => {
      if (!__s) return _no(new Error('storage unavailable'))
      __s.setItem(__k, __v, (_e) => (_e ? _no(_e) : _ok()))
    }),
  del: (__k) =>
    new Promise((_ok, _no) => {
      if (!__s) return _no(new Error('storage unavailable'))
      __s.removeItem(__k, (_e) => (_e ? _no(_e) : _ok()))
    }),
})

export const cloud = (): Box => wrap(__tg__?.CloudStorage)

export const secure = (): Box => wrap(__tg__?.SecureStorage)

export const buzz = (__t: 'success' | 'error' | 'warning' | 'tick' | 'tap'): void => {
  const _h = ver('6.1') ? __tg__?.HapticFeedback : undefined
  if (!_h) return
  if (__t === 'tick') _h.selectionChanged()
  else if (__t === 'tap') _h.impactOccurred('light')
  else _h.notificationOccurred(__t)
}

let __dims = 0

let __tint: string | null = null

const hex = (__rgb: string): string | null => {
  const _m = __rgb.match(/\d+(\.\d+)?/g)
  if (!_m || _m.length < 3 || (_m.length > 3 && Number(_m[3]) === 0)) return null
  return '#' + _m.slice(0, 3).map((_v) => Math.round(Number(_v)).toString(16).padStart(2, '0')).join('')
}

const base = (): string =>
  hex(getComputedStyle(document.body).backgroundColor) ??
  __tg__?.themeParams?.secondary_bg_color ??
  (document.documentElement.dataset.scheme === 'dark' ? '#000000' : '#efeff4')

const shade = (__hex: string, __a: number): string => {
  const _m = /^#?([0-9a-f]{6})$/i.exec(__hex.trim())
  if (!_m) return __hex
  const _n = parseInt(_m[1], 16)
  const _ch = (__v: number) => Math.round(__v * (1 - __a)).toString(16).padStart(2, '0')
  return `#${_ch(_n >> 16)}${_ch((_n >> 8) & 255)}${_ch(_n & 255)}`
}

const paint = (): void => {
  const _c = __tint ?? base()
  try {
    __tg__?.setHeaderColor?.(__dims ? shade(_c, 0.35) : _c)
  } catch {
    return
  }
}

export const tint = (__c: string | null): void => {
  __tint = __c
  paint()
}

export const who = (): { id: number; name: string; full: string; nick: string; photo: string } => {
  const _u = __tg__?.initDataUnsafe?.user
  const _name = _u?.first_name || _u?.username || 'друг'
  return { id: _u?.id ?? 0, name: _name, full: [_u?.first_name, _u?.last_name].filter(Boolean).join(' ') || _name, nick: _u?.username ? '@' + _u.username : '', photo: _u?.photo_url ?? '' }
}

export const dim = (__on: boolean): void => {
  __dims = Math.max(0, __dims + (__on ? 1 : -1))
  paint()
}

export type Theme = 'auto' | 'light' | 'dark' | 'deep'

export const getTheme = (): Theme => {
  try {
    const _t = localStorage.getItem('theme')
    return _t === 'light' || _t === 'dark' || _t === 'deep' ? _t : 'auto'
  } catch {
    return 'auto'
  }
}

export const setTheme = (__t: Theme): void => {
  try {
    localStorage.setItem('theme', __t)
  } catch {
    void 0
  }
  theme()
}

const theme = (): void => {
  const _root = document.documentElement
  const _t = getTheme()
  const _sys = __tg__?.colorScheme ?? (matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light')
  _root.dataset.scheme = _t === 'auto' ? _sys : _t === 'light' ? 'light' : 'dark'
  if (_t === 'deep') _root.dataset.deep = ''
  else delete _root.dataset.deep
  try {
    const _c = base()
    __tg__?.setHeaderColor?.(__tint ?? _c)
    __tg__?.setBackgroundColor?.(_c)
    if (ver('7.10')) __tg__?.setBottomBarColor?.(_c)
  } catch {
    return
  }
}

export const boot = (__full: boolean): void => {
  const _root = document.documentElement
  const _p = __tg__?.platform ?? 'web'
  _root.dataset.platform = _p === 'android' || _p === 'android_x' ? 'android' : _p
  const _apple = _p === 'ios' || _p === 'macos'
  const _blur = CSS.supports?.('backdrop-filter', 'blur(1px)') || CSS.supports?.('-webkit-backdrop-filter', 'blur(1px)')
  if ((__tg__ && !_apple) || !_blur) _root.dataset.glass = 'solid'
  theme()
  addEventListener('load', theme)
  if (!__tg__) return
  __tg__.ready()
  if (__full) __tg__.expand()
  __tg__.onEvent?.('themeChanged', () => requestAnimationFrame(theme))
}

let __back: (() => void) | null = null

export const back = (__cb: (() => void) | null): void => {
  const _b = ver('6.1') ? __tg__?.BackButton : undefined
  if (!_b) return
  if (__back) _b.offClick(__back)
  __back = __cb
  if (!__cb) return _b.hide()
  _b.onClick(__cb)
  _b.show()
}


export const goChat = async (__id: number, __nick?: string): Promise<void> => {
  if (__nick && __tg__?.openTelegramLink) return __tg__.openTelegramLink('https://t.me/' + __nick)
  location.href = 'tg://user?id=' + __id
}

export const __docs__ = 'https://docs.logware.dev'

export const web = (__url: string): void => {
  if (__tg__?.openLink) return __tg__.openLink(__url)
  open(__url, '_blank', 'noopener')
}

export const docs = (__path: string): void => web(__docs__ + __path)

export const path = (__bot: string): string => `Профиль > Изменить > ИИ-автоматизация > @${__bot}`

export const invoice = (__url: string): Promise<'paid' | 'cancelled' | 'failed' | 'pending'> =>
  new Promise((_ok) => {
    if (!__tg__?.openInvoice) {
      location.href = __url
      return _ok('pending')
    }
    __tg__.openInvoice(__url, _ok)
  })

export const tgLink = (__url: string): void => {
  if (__tg__?.openTelegramLink && __url.startsWith('https://t.me/')) return __tg__.openTelegramLink(__url)
  if (__tg__?.openLink) return __tg__.openLink(__url)
  open(__url, '_blank', 'noopener')
}

export const back2 = (): Promise<void> =>
  new Promise((_ok) => {
    const _f = (): void => {
      if (document.visibilityState !== 'visible') return
      document.removeEventListener('visibilitychange', _f)
      _ok()
    }
    document.addEventListener('visibilitychange', _f)
  })
