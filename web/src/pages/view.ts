import * as api from '../api/api'
import type { Ev, Item } from '../api/api'
import { buzz, goChat } from '../app/tg'
import { open, ready, str, unb64 } from '../crypto/box'
import { key } from '../keys/lock'
import * as msg from '../render/msg'
import { skView } from '../ui/skel'
import { h } from '../ui/dom'
import { I } from '../ui/icons'
import { K, empty, page, skel, svg, toast } from '../ui/kit'

const __files__ = new Set(['photo', 'video', 'round', 'voice', 'anim', 'audio', 'doc', 'sticker'])


const copy = async (__t: string): Promise<void> => {
  try {
    await navigator.clipboard.writeText(__t)
  } catch {
    const _a = h('textarea', { style: 'position:fixed;opacity:0' }) as HTMLTextAreaElement
    _a.value = __t
    document.body.append(_a)
    _a.select()
    document.execCommand('copy')
    _a.remove()
  }
  buzz('success')
  toast('Скопировано')
}

type _Msg = { it: Item; m: Record<string, unknown>; text: string; blob: Blob | null; name: string; sec: Uint8Array; ready?: Promise<void>; show?: () => void }

const grab = (__x: _Msg): Promise<void> =>
  (__x.ready ??= api.media(__x.it.id).then((_raw) => {
    __x.blob = new Blob([open(_raw, __x.sec) as BlobPart], { type: msg.mime(__x.m, __x.it.kind) })
  }))

const decode = (__items: Item[], __sec: Uint8Array): _Msg[] => {
  const _xs: _Msg[] = []
  for (const _it of __items) {
    try {
      const _m = JSON.parse(str(open(unb64(_it.body), __sec))) as Record<string, unknown>
      _xs.push({ it: _it, m: _m, text: msg.text(_m), blob: null, name: msg.fname(_m, _it.kind), sec: __sec })
    } catch {
      continue
    }
  }
  return _xs
}

const __eye__ = new IntersectionObserver((_es) => {
  for (const _e of _es) {
    if (!_e.isIntersecting) continue
    __eye__.unobserve(_e.target)
    ;(_e.target as HTMLElement & { fire?: () => void }).fire?.()
  }
}, { rootMargin: '800px 0px' })

const watch = (__el: HTMLElement, __fn: () => void): void => {
  ;(__el as HTMLElement & { fire?: () => void }).fire = __fn
  __eye__.observe(__el)
}

const line = (__x: _Msg): string => {
  const _d = new Date((typeof __x.m.date === 'number' ? __x.m.date : __x.it.at) * 1000)
  const _t = _d.toLocaleString('ru', { day: '2-digit', month: '2-digit', year: '2-digit', hour: '2-digit', minute: '2-digit' })
  const _body = __x.text || msg.label(__x.it.kind) || 'сообщение'
  return `[${_t}] ${msg.who(__x.m, __x.it.out)}: ${_body}`
}

const push = async (__x: _Msg): Promise<void> => {
  const _st = msg.sticker(__x.m)
  if (_st?.id) return api.send('sticker', '', undefined, '', _st.id)
  if (__x.it.media) await grab(__x)
  if (!__x.text && !__x.blob) return
  await api.send(__x.it.kind, __x.text, __x.blob ?? undefined, __x.name)
}

const diff = (__a: string, __b: string): [HTMLElement, HTMLElement] => {
  const _x = __a.split(/(\s+)/)
  const _y = __b.split(/(\s+)/)
  const _t: number[][] = Array.from({ length: _x.length + 1 }, () => new Array(_y.length + 1).fill(0))
  for (let _i = _x.length - 1; _i >= 0; _i--) {
    for (let _j = _y.length - 1; _j >= 0; _j--) {
      _t[_i][_j] = _x[_i] === _y[_j] ? _t[_i + 1][_j + 1] + 1 : Math.max(_t[_i + 1][_j], _t[_i][_j + 1])
    }
  }
  const _old = h('p', { class: 'bub-text' })
  const _new = h('p', { class: 'bub-text' })
  let _i = 0
  let _j = 0
  while (_i < _x.length || _j < _y.length) {
    if (_i < _x.length && _j < _y.length && _x[_i] === _y[_j]) {
      _old.append(_x[_i++])
      _new.append(_y[_j++])
    } else if (_j < _y.length && (_i >= _x.length || _t[_i][_j + 1] >= _t[_i + 1][_j])) {
      const _w = _y[_j++]
      _new.append(/^\s+$/.test(_w) ? _w : h('ins', {}, _w))
    } else {
      const _w = _x[_i++]
      _old.append(/^\s+$/.test(_w) ? _w : h('del', {}, _w))
    }
  }
  return [_old, _new]
}

const bubble = (__x: _Msg, __n: number, __body?: HTMLElement): HTMLElement => {
  const _it = __x.it
  const _date = typeof __x.m.date === 'number' ? __x.m.date : _it.at
  const _b = h('div', { class: 'bub tail' + (_it.media ? ' media' : '') })
  const _time = h('time', {}, msg.clock(_date))
  if (_it.media) {
    const _slot = h('div', { class: 'bub-media' }, skel(260, 300, 16))
    _b.append(_slot)
    __x.show = () => void grab(__x)
      .then(() => _slot.replaceChildren(msg.player(_it.kind, URL.createObjectURL(__x.blob!), __x.name, __x.blob!)))
      .catch(() => _slot.replaceChildren(h('p', { class: 'bub-text dim' }, 'Медиа не загрузилось')))
  } else if (!__x.text && msg.label(_it.kind)) {
    _b.append(h('p', { class: 'bub-text dim' }, msg.label(_it.kind) + (__files__.has(_it.kind) ? ', файл не сохранялся' : '')))
  }
  if (__x.text) _b.addEventListener('click', () => void copy(__x.text))
  if (__body) {
    __body.append(_time)
    _b.append(__body)
  } else if (__x.text) _b.append(h('p', { class: 'bub-text' }, __x.text, _time))
  else _b.append(_time)
  const _row = h('div', { class: 'bub-row' + (_it.out ? ' mine' : '') }, _b)
  if (__x.show) watch(_row, __x.show)
  return _row
}

const pill = (__t: string, __cls = ''): HTMLElement => h('div', { class: 'pill-s ' + __cls }, h('span', {}, __t))

const btn = (__text: string, __fn: () => Promise<void>, __main: boolean, __done: string): HTMLButtonElement => {
  const _b = h('button', { class: 'v-btn' + (__main ? ' main' : ''), type: 'button' }) as HTMLButtonElement
  const _set = (__s: 'idle' | 'busy' | 'done') => {
    _b.disabled = __s !== 'idle'
    _b.classList.toggle('done', __s === 'done')
    _b.replaceChildren(...(__s === 'busy' ? [h('span', { class: 'spin' })] : []), __s === 'done' ? __done : __text)
  }
  _set('idle')
  _b.addEventListener('click', async () => {
    if (_b.disabled) return
    buzz('tap')
    _set('busy')
    try {
      await __fn()
      _set('done')
      setTimeout(() => _set('idle'), 2500)
    } catch {
      _set('idle')
      toast('Не получилось', 'err')
    }
  })
  return _b
}

const title = (__kind: Ev['kind'], __xs: _Msg[], __n: number, __name: string): string => {
  if (__kind === 'edit') return 'Изменено · ' + __name
  if (__kind === 'kw') return 'Ключевое слово · ' + __name
  if (__kind === 'once') return 'Одноразовое ' + (__xs[0]?.it.kind === 'video' ? 'видео' : 'фото') + ' · ' + __name
  return (__n > 1 ? 'Удалено ' + msg.many(__n, ['сообщение', 'сообщения', 'сообщений']) : 'Удалено') + ' · ' + __name
}

const __few__ = 10

export const view = async (__root: HTMLElement, __id: string): Promise<void> => {
  __root.replaceChildren(skView())
  const [_me, _ev] = await Promise.all([api.me().catch(() => null), api.ev(__id).catch((_e: Error) => _e)])
  if (!_me) return empty(__root, I.cloud, K.gray, 'Нет связи', 'Не получилось связаться с сервером.')
  if (_ev instanceof Error) return empty(__root, I.clock, K.gray, 'Уже удалено', 'Срок хранения этого сообщения вышел, либо ты сменил ключ.')
  if (!_me.reg) return empty(__root, I.key, K.orange, 'Нет ключа', 'Открой бота и создай ключ заново.')
  await ready()
  const _sec = await key(__root, _me)
  if (!_sec) {
    return empty(__root, I.lock, K.gray, 'Ключ не найден', _me.mode === 'secure'
      ? 'Ключ лежит в другом телефоне. Открой сообщение на нём.'
      : 'Не получилось достать ключ из Telegram.')
  }
  if (_ev.items.length === 0) return empty(__root, I.clock, K.gray, 'Уже удалено', 'Срок хранения этих сообщений вышел.')
  const _xs = decode(_ev.items, _sec)
  if (_xs.length === 0) return empty(__root, I.lock, K.gray, 'Не открылось', 'Эти сообщения закрыты другим ключом.')
  const _peer = _xs.find((_x) => !_x.it.out)
  const _name = _peer ? msg.who(_peer.m, false) : 'Ты'
  const _stage = h('div', { class: 'stage' }, pill(title(_ev.kind, _xs, _ev.total, _name), 'head'))
  const _edit = _ev.kind === 'edit' && _xs.length > 1
  if (_edit) {
    const _a = _xs[0]
    const _b = _xs[_xs.length - 1]
    const [_old, _new] = diff(_a.text, _b.text)
    const _ar = h('div', { class: 'arrow' })
    _ar.append(svg(I.chev))
    _stage.append(pill('Было', 'lbl'), h('div', { class: 'was' }, bubble(_a, 0, _old)), _ar, pill('Стало', 'lbl'), bubble(_b, 1, _new))
  } else {
    _xs.forEach((_x, _i) => _stage.append(bubble(_x, _i)))
  }
  let _from = _ev.items.length
  let _more = _ev.more
  const next = (): void => {
    if (!_more) return
    const _tail = h('div', { class: 'more' }, h('span', { class: 'spin' }))
    _stage.append(_tail)
    watch(_tail, async () => {
      try {
        const _pg = await api.ev(__id, _from, 50)
        _from += _pg.items.length
        _more = _pg.more && _pg.items.length > 0
        const _new = decode(_pg.items, _sec)
        _xs.push(..._new)
        _tail.replaceWith(..._new.map((_x, _i) => bubble(_x, _i)))
        next()
      } catch {
        _tail.replaceChildren(h('span', { class: 'bub-text dim' }, 'Не получилось загрузить'))
      }
    })
  }
  next()
  const all = async (): Promise<_Msg[]> => {
    const _out = [..._xs]
    let _at = _from
    let _go = _more
    while (_go) {
      const _pg = await api.ev(__id, _at, 200)
      _at += _pg.items.length
      _go = _pg.more && _pg.items.length > 0
      _out.push(...decode(_pg.items, _sec))
    }
    return _out
  }
  const send = async () => {
    if (_edit) return void (await api.send('text', _xs[0].text + '\n\n' + _xs[_xs.length - 1].text))
    if (_ev.total <= __few__) {
      for (const _x of _xs) await push(_x)
      return
    }
    const _list = await all()
    const _file = new Blob([_list.map(line).join('\n') + '\n'], { type: 'text/plain' })
    await api.send('doc', 'Переписка с ' + _name + ' · ' + msg.many(_list.length, ['сообщение', 'сообщения', 'сообщений']), _file, 'Переписка с ' + _name + '.txt')
  }
  const _chat = (_peer ?? _xs[0]).m.chat as { id?: number; username?: string } | undefined
  const _bar = h('div', { class: 'v-bar' })
  _bar.append(btn(_ev.total > __few__ ? 'Сохранить файлом' : 'Сохранить', async () => {
    await send()
    buzz('success')
    toast('Сохранил в чат с ботом')
  }, true, 'Сохранено'))
  if (_chat?.id) {
    const _go = h('button', { class: 'v-link', type: 'button' }, 'Перейти в чат')
    _go.addEventListener('click', () => {
      buzz('tick')
      void goChat(_chat.id!, _chat.username)
    })
    _bar.append(_go)
  }
  __root.replaceChildren(page('bare view', _stage, _bar))
}
