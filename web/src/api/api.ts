import { init } from '../app/tg'

export type Prefs = {
  del: boolean
  delMe: boolean
  delMute: boolean
  edit: boolean
  editMe: boolean
  editMute: boolean
  media: boolean
  photo: boolean
  video: boolean
  voice: boolean
  files: boolean
  once: boolean
  onceAuto: boolean
  onceHide: boolean
  onceMute: boolean
  panic: boolean
  panicN: number
  keep: number
}

export type Mode = 'cloud' | 'secure' | 'pass'

export type Me = {
  plan: 'free' | 'pro' | 'dev'
  until: number
  stats: { caught: number; today: number; kept: number }
  reg: boolean
  mode?: Mode
  vault?: string
  prefs: Prefs
  conn: { active: boolean; reply: boolean; read?: boolean; delOwn?: boolean; delAll: boolean; name?: boolean } | null
  bot: string
  dev: boolean
  via: ('stars' | 'cb' | 'xr')[]
  stars: boolean
  support: string
}

export type DevPlug = { id: number; slug: string; name: string; about: string; ver: string; repo: string; live: boolean; last: '' | 'wait' | 'ok' | 'no'; why: string; kind: string; users: number; pic: string }

export type DevReq = { id: number; kind: 'new' | 'update' | 'edit'; name: string; about: string; ver: string; note: string; repo: string; pic: boolean; status: 'wait' | 'ok' | 'no'; reason: string; created: number; done: number }

export type Dev = { nick: string; bio: string; plugs: DevPlug[]; gh: string }

export type DevOne = DevPlug & { reqs: DevReq[] }

export type PlugIn = { kind?: 'update' | 'edit'; slug?: string; name: string; about: string; ver: string; note: string; repo: string; pic: string }

export type Item = {
  id: number
  ver: number
  out: boolean
  kind: string
  body: string
  media: boolean
  at: number
}

export type Ev = {
  kind: 'del' | 'edit' | 'once' | 'kw'
  at: number
  total: number
  more: boolean
  items: Item[]
}

const head = (): HeadersInit => ({
  'X-Init': init(),
  'Content-Type': 'application/json',
  'ngrok-skip-browser-warning': '1',
})

export const call = async <T>(__path: string, __body?: unknown): Promise<T> => {
  const _r = await fetch('/api/' + __path, {
    method: __body === undefined ? 'GET' : 'POST',
    headers: head(),
    body: __body === undefined ? undefined : JSON.stringify(__body),
  })
  if (!_r.ok) throw new Error(_r.status === 404 ? 'not found' : (await _r.text()) || String(_r.status))
  return _r.json() as Promise<T>
}

export const me = () => call<Me>('me')

export const promo = (__code: string) => call<{ plan: 'pro' | 'dev'; until: number }>('promo', { code: __code })

export const pay = (__plan: string, __via = 'stars') => call<{ url: string }>('pay', { plan: __plan, via: __via })

export const reg = (__mode: Mode, __pub: string, __vault: string | null) =>
  call<{ ok: boolean }>('reg', { mode: __mode, pub: __pub, vault: __vault })

export const prefs = (__p: Prefs) => call<Prefs>('prefs', __p)

export const wipe = () => call<{ ok: boolean }>('wipe', {})

export const ev = (__id: string, __from = 0, __n = 50) => call<Ev>(`ev/${__id}?from=${__from}&n=${__n}`)

export const send = async (__kind: string, __text: string, __file?: Blob, __name = 'file', __id = ''): Promise<void> => {
  const _f = new FormData()
  _f.append('kind', __kind)
  _f.append('text', __text)
  if (__id) _f.append('id', __id)
  else if (__file) _f.append('file', __file, __name)
  const _h = head() as Record<string, string>
  delete _h['Content-Type']
  const _r = await fetch('/api/send', { method: 'POST', headers: _h, body: _f })
  if (!_r.ok) throw new Error('send ' + _r.status)
}

export const pic = async (__id: number, __base = '/api/ava/'): Promise<Blob> => {
  const _r = await fetch(__base + __id, { headers: head() })
  if (!_r.ok) throw new Error('ava ' + _r.status)
  return _r.blob()
}

export const media = async (__id: number): Promise<Uint8Array> => {
  const _r = await fetch('/api/media/' + __id, { headers: head() })
  if (!_r.ok) throw new Error('media ' + _r.status)
  return new Uint8Array(await _r.arrayBuffer())
}

export type Opt = { key: string; name: string; def: string; choices: { v: string; label: string }[] | null; text?: boolean; max?: number; hide?: boolean }

export type Mod = {
  id: string
  pic?: string
  name: string
  desc: string
  long?: string
  cat?: string
  ver?: string
  check?: boolean
  perms?: { right: string; why: string }[] | null
  icon: string
  color: string
  author: string
  soon: boolean
  hidden?: boolean
  on: boolean
  ready: boolean
  cfg: Record<string, string>
  opts: Opt[]
}

export const mods = async (): Promise<Mod[]> =>
  (await call<Mod[]>('mods')).map((_m) => ({ ..._m, opts: _m.opts ?? [], perms: _m.perms ?? [], cfg: _m.cfg ?? {} }))

export const setMod = (__id: string, __on: boolean, __cfg: Record<string, string>) =>
  call<{ on: boolean; cfg: Record<string, string> }>('mods/' + __id, { on: __on, cfg: __cfg })

export type Hist = { id: string; kind: Ev['kind']; at: number; n: number; body: string; type: string }

export const hist = (__from = 0, __n = 40) => call<{ items: Hist[]; more: boolean }>(`hist?from=${__from}&n=${__n}`)

export const clear = () => call<{ ok: boolean }>('clear', {})

export type Chat = { id: number; body: string; n: number; off: boolean }

export const chats = () => call<Chat[]>('chats')

export const setChat = (__id: number, __off: boolean) => call<{ ok: boolean }>('chats', { id: __id, off: __off })

export type Kw = { hash: string; box: string }

export const kws = () => call<Kw[]>('kws')

export const addKw = (__k: Kw) => call<{ ok: boolean }>('kws', __k)

export const delKw = (__hash: string) => call<{ ok: boolean }>('kws/del', { hash: __hash, box: '' })

export const dev = () => call<Dev>('dev')

export const setDev = (__nick: string, __bio: string) => call<{ nick: string; bio: string }>('dev', { nick: __nick, bio: __bio })

export const devPlug = (__id: number) => call<DevOne>(`dev/plug/${__id}`)

export const newPlug = (__p: PlugIn) => call<{ id: number }>('dev/plug', __p)

export const changePlug = (__id: number, __p: PlugIn) => call<{ id: number }>(`dev/plug/${__id}`, __p)

const __pics__ = new Map<number, Promise<string | null>>()

export const reqPic = (__id: number): Promise<string | null> => {
  let _p = __pics__.get(__id)
  if (!_p) {
    _p = fetch(`/api/dev/req/${__id}/pic`, { headers: head() }).then((_r) => (_r.ok ? _r.blob() : null)).then((_b) => (_b ? URL.createObjectURL(_b) : null)).catch(() => null)
    __pics__.set(__id, _p)
  }
  return _p
}
