import * as api from '../api/api'
import type { Hist, Me } from '../api/api'
import { go } from '../app/router'
import { buzz, who } from '../app/tg'
import { b64, bytes, open, pubOf, ready, seal, sha, str, unb64 } from '../crypto/box'
import { key } from '../keys/lock'
import * as msg from '../render/msg'
import { ava } from '../ui/ava'
import { h } from '../ui/dom'
import { C, I } from '../ui/icons'
import { K, ask, box, empty, group, page, row, sq, svg, sw, toast } from '../ui/kit'
import { head, howRow } from './feat'

type _Obj = Record<string, unknown>

const __kinds__: Record<string, [string, string]> = {
  del: ['🗑', 'Удалено'],
  edit: ['✏️', 'Изменено'],
  once: ['🔥', 'Одноразовое'],
  kw: ['🔔', 'Ключевое слово'],
}

const chatId = (__m: _Obj | null): number => Number((__m?.chat as _Obj | undefined)?.id) || 0

const peek = (__body: string, __sec: Uint8Array): _Obj | null => {
  try {
    return JSON.parse(str(open(unb64(__body), __sec))) as _Obj
  } catch {
    return null
  }
}

const day = (__unix: number): string => {
  const _d = new Date(__unix * 1000)
  const _n = new Date()
  const _y = new Date(_n.getFullYear(), _n.getMonth(), _n.getDate() - 1)
  if (_d.toDateString() === _n.toDateString()) return 'Сегодня'
  if (_d.toDateString() === _y.toDateString()) return 'Вчера'
  return _d.toLocaleDateString('ru', { day: 'numeric', month: 'long' })
}

const esc = (__s: string): string => __s.replace(/[&<>"]/g, (_c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;' })[_c]!)

type _Line = { at: number; who: string; text: string; out: boolean }

type _Block = { kind: string; at: number; chat: string; lines: _Line[] }

const dump = async (__sec: Uint8Array): Promise<_Block[]> => {
  const _out: _Block[] = []
  let _from = 0
  for (;;) {
    const _pg = await api.hist(_from, 200)
    _from += _pg.items.length
    for (const _e of _pg.items) {
      const _b: _Block = { kind: _e.kind, at: _e.at, chat: '', lines: [] }
      let _at = 0
      for (;;) {
        const _ev = await api.ev(_e.id, _at, 200)
        _at += _ev.items.length
        for (const _it of _ev.items) {
          const _m = peek(_it.body, __sec)
          if (!_m) continue
          _b.chat ||= msg.chat(_m)
          _b.lines.push({ at: typeof _m.date === 'number' ? _m.date : _it.at, who: msg.who(_m, _it.out), text: msg.text(_m) || msg.label(_it.kind) || 'сообщение', out: _it.out })
        }
        if (!_ev.more || _ev.items.length === 0) break
      }
      _out.push(_b)
    }
    if (!_pg.more || _pg.items.length === 0) break
  }
  return _out
}

const stamp = (__unix: number): string =>
  new Date(__unix * 1000).toLocaleString('ru', { day: '2-digit', month: '2-digit', year: '2-digit', hour: '2-digit', minute: '2-digit' })

const txt = (__bs: _Block[]): string =>
  __bs.map((_b) => `== ${__kinds__[_b.kind]?.[1] ?? _b.kind} · ${_b.chat} · ${stamp(_b.at)} ==\n` + _b.lines.map((_l) => `[${stamp(_l.at)}] ${_l.who}: ${_l.text}`).join('\n')).join('\n\n') + '\n'

const html = (__bs: _Block[]): string => `<!doctype html><html lang="ru"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>История Logware</title>
<style>body{margin:0;padding:24px 16px;background:#1c1c1d;color:#fff;font:16px/1.4 -apple-system,system-ui,sans-serif}h1{font-size:24px;margin:0 0 20px}section{max-width:640px;margin:0 auto 28px}h2{font-size:14px;color:#98989f;font-weight:600;margin:0 0 10px;text-align:center}.m{max-width:80%;margin:4px 0;padding:8px 12px;border-radius:18px;background:#2c2c2e;white-space:pre-wrap;word-break:break-word}.m.o{margin-left:auto;background:#4b8ff7}.m small{display:block;font-size:12px;opacity:.6;margin-top:2px}.w{font-size:13px;font-weight:600;opacity:.75}</style></head><body><section><h1>История Logware</h1></section>` +
  __bs.map((_b) => `<section><h2>${esc(__kinds__[_b.kind]?.[1] ?? _b.kind)} · ${esc(_b.chat)} · ${stamp(_b.at)}</h2>` +
    _b.lines.map((_l) => `<div class="m${_l.out ? ' o' : ''}"><div class="w">${esc(_l.who)}</div>${esc(_l.text)}<small>${stamp(_l.at)}</small></div>`).join('') + '</section>').join('') + '</body></html>'

const histPage = (__root: HTMLElement, __sec: Uint8Array): HTMLElement => {
  const _list = h('div', {})
  const _foot = h('div', {})
  let _from = 0
  let _last = ''
  let _cur: HTMLElement | null = null
  const add = (__e: Hist): void => {
    const _m = peek(__e.body, __sec)
    const _name = _m ? msg.chat(_m) : 'Не открылось'
    const _txt = _m ? msg.text(_m).replace(/\s+/g, ' ').trim() : ''
    const _what = _txt ? `«${_txt}»` : msg.plain(__e.type).toLowerCase()
    const _sub = __e.kind === 'del'
      ? (__e.n > 1 ? 'Удалено ' + msg.many(__e.n, ['сообщение', 'сообщения', 'сообщений']) : 'Удалено: ' + _what)
      : __e.kind === 'edit' ? 'Изменено: ' + _what
      : __e.kind === 'once' ? 'Одноразовое ' + (__e.type === 'video' ? 'видео' : 'фото')
      : 'Ключевое слово: ' + _what
    const _d = day(__e.at)
    if (_d !== _last) {
      _last = _d
      _cur = box()
      _list.append(h('div', { class: 'gh' }, h('h2', {}, _d)), _cur)
    }
    _cur!.append(row({ lead: ava(chatId(_m), _name, 'm-ava sm'), title: _name, sub: _sub, value: new Date(__e.at * 1000).toLocaleTimeString('ru', { hour: '2-digit', minute: '2-digit' }), chev: true, onclick: () => go('/e/' + __e.id) }))
  }
  const load = async (): Promise<void> => {
    _foot.replaceChildren(h('div', { class: 'more' }, h('span', { class: 'spin' })))
    try {
      const _pg = await api.hist(_from, 40)
      _from += _pg.items.length
      _pg.items.forEach(add)
      if (_from === 0) _list.append(h('div', { class: 'center', style: 'min-height:auto;padding:28px 24px' }, h('p', {}, 'Пока пусто. Здесь появится всё, что я сохраню.')))
      _foot.replaceChildren()
      if (_pg.more) {
        const _b = h('button', { class: 'v-link', type: 'button', style: 'display:block;margin:6px auto' }, 'Показать ещё')
        _b.addEventListener('click', () => void load())
        _foot.append(_b)
      }
    } catch {
      _foot.replaceChildren(h('p', { class: 'note err' }, 'Не получилось загрузить'))
    }
  }
  void load()
  const exp = (__kind: 'txt' | 'html') => async () => {
    toast('Собираю историю…')
    try {
      const _bs = await dump(__sec)
      if (_bs.length === 0) return toast('История пустая', 'err')
      const _file = __kind === 'txt' ? new Blob([txt(_bs)], { type: 'text/plain' }) : new Blob([html(_bs)], { type: 'text/html' })
      await api.send('doc', 'История Logware', _file, 'logware.' + __kind)
      buzz('success')
      toast('Прислал файл в чат с ботом')
    } catch {
      toast('Не получилось выгрузить', 'err')
    }
  }
  const _wipe = row({ title: 'Очистить историю', cls: 'danger', onclick: async () => {
    if (!(await ask('Очистить историю?', 'Все сохранённые сообщения удалятся навсегда. Ключ и настройки останутся.', 'Очистить', true))) return
    try {
      await api.clear()
      toast('История очищена')
      __root.replaceChildren(histPage(__root, __sec))
    } catch {
      toast('Не получилось очистить', 'err')
    }
  } })
  return page('feat',
    head('deleted', 'История', 'Всё, что я сохранил: удалённые, изменённые и одноразовые.'),
    group('Выгрузить', null,
      row({ lead: sq(I.copy, C.blue), title: 'Текстом (.txt)', chev: true, onclick: exp('txt') }),
      row({ lead: sq(I.grid, C.purple), title: 'Страницей (.html)', chev: true, onclick: exp('html') })),
    _list,
    _foot,
    h('div', { class: 'gap' }),
    box(_wipe),
    howRow('/features/history'),
  )
}

const norm = (__s: string): string[] =>
  __s.toLowerCase().replace(/ё/g, 'е').split(/[^\p{L}\p{N}]+/u).filter(Boolean)

const kwPage = (__sec: Uint8Array): HTMLElement => {
  const _pub = pubOf(__sec)
  const _uid = who().id
  const _box = box()
  const _count = h('span', {})
  const _in = h('input', { class: 'kw-in', type: 'text', placeholder: 'Слово или фраза', maxlength: '50', enterkeyhint: 'done' }) as HTMLInputElement
  let _n = 0
  const draw = (__l: api.Kw[]): void => {
    _n = __l.length
    _count.textContent = `${_n} из 20`
    _box.replaceChildren(...__l.map((_k) => {
      let _w = '•••'
      try {
        _w = str(open(unb64(_k.box), __sec))
      } catch {
        void 0
      }
      const _x = h('button', { class: 'kw-x', type: 'button', 'aria-label': 'Удалить' })
      _x.append(svg(I.close))
      const _r = row({ lead: sq(I.tag, C.indigo), title: _w, right: _x })
      _x.addEventListener('click', async () => {
        buzz('tick')
        try {
          await api.delKw(_k.hash)
          void refresh()
        } catch {
          toast('Не получилось удалить', 'err')
        }
      })
      return _r
    }))
    _box.style.display = _n ? '' : 'none'
  }
  const refresh = async () => draw(await api.kws().catch(() => []))
  const add = async (): Promise<void> => {
    const _raw = _in.value.trim()
    const _ws = norm(_raw)
    if (_ws.length === 0) return
    if (_ws.length > 3) return toast('Не больше трёх слов', 'err')
    if (_n >= 20) return toast('Можно до 20 слов', 'err')
    try {
      await api.addKw({ hash: b64(sha(`lw:${_uid}:${_ws.join(' ')}`)), box: b64(seal(bytes(_raw), _pub)) })
      _in.value = ''
      buzz('success')
      await refresh()
    } catch {
      toast('Не получилось добавить', 'err')
    }
  }
  const _addBtn = h('button', { class: 'kw-go', type: 'button' }, 'Добавить')
  _addBtn.addEventListener('click', () => void add())
  _in.addEventListener('keydown', (_e) => _e.key === 'Enter' && void add())
  void refresh()
  return page('feat',
    head('ring', 'Ключевые слова', 'Напишу тебе, если собеседник упомянет одно из этих слов.'),
    box(h('label', { class: 'line kw-add' }, _in, _addBtn)),
    h('div', { class: 'gh' }, h('h2', {}, 'Мои слова'), _count),
    _box,
    howRow('/features/keywords'),
  )
}

const ignPage = (__sec: Uint8Array): HTMLElement => {
  const _body = h('div', {}, h('div', { class: 'more' }, h('span', { class: 'spin' })))
  void api.chats().then((_cs) => {
    if (_cs.length === 0) return _body.replaceChildren(h('div', { class: 'center', style: 'min-height:auto;padding:28px 24px' }, h('p', {}, 'Чаты появятся здесь, когда тебе кто-нибудь напишет.')))
    _body.replaceChildren(group('Чаты', null, ..._cs.map((_c) => {
      const _m = peek(_c.body, __sec)
      const _name = _m ? msg.chat(_m) : 'Чат'
      const _s = sw(!_c.off, async (_v) => {
        try {
          await api.setChat(_c.id, !_v)
        } catch {
          toast('Не получилось сохранить', 'err')
        }
      })
      return row({ lead: ava(_c.id, _name, 'm-ava sm'), title: _name, sub: 'Сохранено ' + msg.many(_c.n, ['сообщение', 'сообщения', 'сообщений']), right: _s, onclick: () => _s.flip() })
    })))
  }).catch(() => _body.replaceChildren(h('p', { class: 'note err' }, 'Не получилось загрузить')))
  return page('feat',
    head('rest', 'Чаты без слежки', 'Выключи чат, и я перестану сохранять его и присылать уведомления.'),
    _body,
    howRow('/features/ignore'),
  )
}

export const more = async (__root: HTMLElement, __key: string, __me: Me): Promise<void> => {
  await ready()
  const _sec = await key(__root, __me)
  if (!_sec) return empty(__root, I.lock, K.gray, 'Ключ не найден', 'Не получилось достать ключ.')
  __root.replaceChildren(__key === 'hist' ? histPage(__root, _sec) : __key === 'kw' ? kwPage(_sec) : ignPage(_sec))
  scrollTo(0, 0)
}
