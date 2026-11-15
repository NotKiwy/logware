import { buzz, dim } from '../app/tg'
import { h } from './dom'
import { I } from './icons'

type _Kid = Node | string | null | undefined | false

export const svg = (__html: string): DocumentFragment => {
  const _t = document.createElement('template')
  _t.innerHTML = __html
  return _t.content
}

export const page = (__cls: string, ...__kids: _Kid[]): HTMLElement => h('div', { class: 'app ' + __cls }, ...__kids)

export const orb = (__ic: string, __color: string, __cls = ''): HTMLElement => {
  const _o = h('span', { class: 'orb ' + __cls, style: `background:${__color};--ic-ink:${__color}` })
  _o.append(svg(__ic))
  return _o
}

export const coin = (__ic: string, __k: [string, string, string]): HTMLElement => {
  const _c = h('div', { class: 'coin', style: `--k1:${__k[0]};--k2:${__k[1]};--k3:${__k[2]};--ic-ink:${__k[1]}` })
  _c.append(svg(__ic))
  return _c
}

export const K: Record<string, [string, string, string]> = {
  blue: ['#9cd0ff', '#3f8cf4', '#2a63c4'],
  green: ['#a4efc4', '#2fae74', '#1f7f54'],
  violet: ['#cdbfff', '#7460e6', '#4f3fb4'],
  orange: ['#ffd9a1', '#f2932e', '#bd6a12'],
  gray: ['#c9c9cf', '#8a8a90', '#5f5f65'],
}

export const link = (__text: string, __click: () => void): HTMLElement => {
  const _b = h('button', { class: 'link', type: 'button' }, __text + ' ›')
  _b.addEventListener('click', () => {
    buzz('tick')
    __click()
  })
  return _b
}

export const box = (...__kids: _Kid[]): HTMLElement => h('div', { class: 'box' }, ...__kids)

export const group = (__title: string, __link: HTMLElement | null, ...__kids: _Kid[]): DocumentFragment => {
  const _f = document.createDocumentFragment()
  _f.append(h('div', { class: 'gh' }, h('h2', {}, __title), __link), box(...__kids))
  return _f
}

export const sec = (__title: string, __sub?: string, __cls = ''): HTMLElement =>
  h('div', { class: 'sec ' + __cls }, h('h2', {}, __title), __sub && h('p', {}, __sub))

export const lic = (__ic: string): HTMLElement => {
  const _l = h('span', { class: 'lic' })
  if (__ic.startsWith('<')) _l.append(svg(__ic))
  else _l.textContent = __ic
  return _l
}

export const sq = (__ic: string, __color: string): HTMLElement => {
  const _q = h('span', { class: 'sq', style: `background:${__color};--ic-ink:${__color}` })
  _q.append(svg(__ic))
  return _q
}

export const appIc = (__ic: string, __color: string, __big = false): HTMLElement => {
  const _a = h('span', { class: 'app-ic' + (__big ? ' big' : ''), style: `background:${__color};--ic-ink:${__color}` })
  _a.append(svg(__ic))
  return _a
}

export const picIc = (__url: string | Promise<string | null> | null | undefined, __ic: string, __color: string, __big = false): HTMLElement => {
  const _a = appIc(__ic, __color, __big)
  const put = (__u: string | null): void => {
    if (!__u) return
    _a.replaceChildren()
    _a.classList.add('pic')
    _a.style.backgroundImage = `url("${__u}")`
  }
  if (typeof __url === 'string') put(__url)
  else if (__url) void __url.then(put)
  return _a
}

export type Add = HTMLElement & { set: (_v: boolean) => void }

export const add = (__on: boolean, __change: (_v: boolean) => void): Add => {
  const _b = h('button', { class: 'add', type: 'button' }) as unknown as Add
  let _v = __on
  _b.set = (__x) => {
    _v = __x
    _b.classList.toggle('on', _v)
    _b.replaceChildren(svg(_v ? I.oOk : I.oPlus))
  }
  _b.set(__on)
  _b.addEventListener('click', (_e) => {
    _e.stopPropagation()
    buzz('tick')
    _b.set(!_v)
    __change(_v)
  })
  return _b
}

export type RowOpts = {
  lead?: Node
  title: string
  badge?: boolean
  sub?: string
  subCls?: string
  subIc?: string
  meta?: Node
  value?: string
  right?: Node
  chev?: boolean
  check?: boolean
  cls?: string
  onclick?: () => void
}

export const row = (__o: RowOpts): HTMLElement => {
  const _tag = __o.onclick ? 'button' : 'div'
  const _t = h('span', { class: 'row-t' }, __o.title)
  if (__o.badge) _t.append(svg(I.badge))
  const _s = __o.sub ? h('span', { class: 'row-s' + (__o.subCls ? ' ' + __o.subCls : '') }) : null
  if (_s && __o.subIc) _s.append(svg(__o.subIc))
  if (_s) _s.append(__o.sub!)
  const _lc = !__o.lead ? '' : (__o.lead as HTMLElement).classList?.contains('app-ic') ? ' app' : ' lead'
  const _r = h(_tag, { class: 'line' + _lc + (__o.cls ? ' ' + __o.cls : ''), type: _tag === 'button' ? 'button' : null },
    __o.lead,
    h('span', { class: 'row-main' }, _t, _s, __o.meta),
    __o.value && h('span', { class: 'val' }, __o.value),
    __o.right,
  )
  if (__o.chev) _r.append(svg(I.chev))
  if (__o.check) _r.append(svg(I.check))
  if (__o.onclick) _r.addEventListener('click', () => {
    buzz('tick')
    __o.onclick!()
  })
  return _r
}

export type Sw = HTMLElement & { flip: () => void }

export const sw = (__on: boolean, __change: (_v: boolean) => void): Sw => {
  const _s = h('span', { class: 'switch' + (__on ? ' on' : ''), role: 'switch', 'aria-checked': String(__on) }, h('span', {})) as unknown as Sw
  let _v = __on
  _s.flip = () => {
    _v = !_v
    _s.classList.toggle('on', _v)
    _s.setAttribute('aria-checked', String(_v))
    buzz('tick')
    __change(_v)
  }
  _s.addEventListener('click', (_e) => {
    _e.stopPropagation()
    _s.flip()
  })
  return _s
}

export const seg = <T extends string>(__opts: { v: T; label: string }[], __cur: T, __change: (_v: T) => void): HTMLElement => {
  const _box = h('div', { class: 'segmented', role: 'tablist', style: `--n:${__opts.length};--i:${Math.max(0, __opts.findIndex((_o) => _o.v === __cur))}` }, h('span', { class: 'thumb' }))
  let _cur = __cur
  const _btns = __opts.map((_o, _n) => {
    const _b = h('button', { class: _o.v === __cur ? 'on' : '', type: 'button', role: 'tab' }, _o.label)
    _b.addEventListener('click', () => {
      if (_o.v === _cur) return
      _cur = _o.v
      for (const _x of _btns) _x.classList.toggle('on', _x === _b)
      _box.style.setProperty('--i', String(_n))
      buzz('tick')
      __change(_o.v)
    })
    return _b
  })
  _box.append(..._btns)
  return _box
}

export const list = (__items: [string, string][], __nums = false): HTMLElement =>
  h('div', { class: 'list' }, ...__items.map(([_e, _t]) => h('div', {}, h(__nums ? 'b' : 'i', {}, _e), h('span', {}, _t))))

export const note = (__text: string, __err = false): HTMLParagraphElement => h('p', { class: 'note' + (__err ? ' err' : '') }, __text)

export type Btn = HTMLButtonElement & { busy: (_v: boolean, _text?: string) => void }

export const cta = (__text: string, __click: () => void, __cls = ''): Btn => {
  const _b = h('button', { class: 'cta ' + __cls, type: 'button' }, __text) as Btn
  _b.addEventListener('click', () => {
    if (_b.disabled) return
    buzz('tap')
    __click()
  })
  _b.busy = (__v, __t) => {
    _b.disabled = __v
    _b.replaceChildren(...(__v ? [h('span', { class: 'spin' })] : []), __t ?? __text)
  }
  return _b
}

let __toast: HTMLElement | null = null

export const toast = (__text: string, __kind: 'ok' | 'err' = 'ok'): void => {
  __toast?.remove()
  const _i = h('i', { class: __kind === 'err' ? 'err' : '' })
  _i.append(svg(__kind === 'err' ? I.err : I.ok))
  const _t = h('div', { class: 'toast', role: 'status' }, _i, h('span', {}, __text))
  __toast = _t
  document.body.append(_t)
  if (__kind === 'err') buzz('error')
  setTimeout(() => _t.classList.add('out'), 2200)
  setTimeout(() => _t.remove(), 2450)
}

export const sheet = (__title: string | null, ...__kids: _Kid[]): { close: () => void; done: Promise<void> } => {
  let _end: () => void = () => {}
  const _done = new Promise<void>((_ok) => (_end = _ok))
  const _x = h('button', { class: 'round', type: 'button', 'aria-label': 'Закрыть' })
  _x.append(svg(I.close))
  const _s = h('div', { class: 'sheet', role: 'dialog' }, h('div', { class: 'sheet-head' }, _x, __title && h('h2', {}, __title)), ...__kids)
  const _w = h('div', { class: 'sheet-wrap' }, _s)
  let _open = true
  const _close = () => {
    if (!_open) return
    _open = false
    dim(false)
    _w.classList.add('out')
    setTimeout(() => {
      _w.remove()
      _end()
    }, 240)
  }
  _x.addEventListener('click', () => {
    buzz('tap')
    _close()
  })
  _w.addEventListener('click', _close)
  _s.addEventListener('click', (_e) => _e.stopPropagation())
  dim(true)
  document.body.append(_w)
  return { close: _close, done: _done }
}

export const pick = <T>(__title: string, __text: string | null, __opts: { v: T; label: string }[], __cur: T): Promise<T | null> => {
  let _res: T | null = null
  const _rows = __opts.map((_o) =>
    row({ title: _o.label, check: true, cls: _o.v === __cur ? 'on' : '', onclick: () => {
      _res = _o.v
      _s.close()
    } }),
  )
  const _s = sheet(__title, __text && h('p', { class: 'sheet-text' }, __text), box(..._rows))
  return _s.done.then(() => _res)
}

export const ask = (__title: string, __text: string, __yes: string, __danger = false): Promise<boolean> => {
  let _res = false
  const _b = cta(__yes, () => {
    _res = true
    _s.close()
  }, __danger ? 'danger' : '')
  const _s = sheet(__title, h('p', { class: 'sheet-text' }, __text), _b)
  return _s.done.then(() => _res)
}

export const info = (__title: string, __text: string | null, ...__kids: _Kid[]): void => {
  const _b = cta('Понятно', () => _s.close())
  const _s = sheet(__title, __text && h('p', { class: 'sheet-text' }, __text), ...__kids, _b)
}

export const skel = (__w: string | number, __hgt: number, __r = 10): HTMLElement =>
  h('span', { class: 'skel', style: `width:${typeof __w === 'number' ? __w + 'px' : __w};height:${__hgt}px;border-radius:${__r}px` })

export const loading = (__root: HTMLElement): void => {
  __root.replaceChildren(
    page('',
      h('div', { class: 'hi', style: 'gap:18px' }, skel(112, 112, 56), skel(200, 34, 12)),
      skel('100%', 240, 26),
    ),
  )
}

export const empty = (__root: HTMLElement, __ic: string, __k: [string, string, string], __title: string, __text: string): void => {
  __root.replaceChildren(h('div', { class: 'center' }, coin(__ic, __k), h('h1', {}, __title), h('p', {}, __text)))
}

export const wait = (__root: HTMLElement, __text: string): void => {
  __root.replaceChildren(h('div', { class: 'center' }, h('span', { class: 'spin big' }), h('p', {}, __text)))
}
