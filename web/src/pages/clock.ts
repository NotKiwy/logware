import type { Mod } from '../api/api'
import { buzz, who } from '../app/tg'
import { ava } from '../ui/ava'
import { h } from '../ui/dom'
import { box, group, page, row, sw } from '../ui/kit'
import { head, howRow } from './feat'
import { saveMod } from './mods'

type _Cfg = Record<string, string>

const __days__ = ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб']

const __full__ = ['воскресенье', 'понедельник', 'вторник', 'среда', 'четверг', 'пятница', 'суббота']

const __mons__ = ['янв', 'фев', 'мар', 'апр', 'мая', 'июн', 'июл', 'авг', 'сен', 'окт', 'ноя', 'дек']

const __months__ = ['января', 'февраля', 'марта', 'апреля', 'мая', 'июня', 'июля', 'августа', 'сентября', 'октября', 'ноября', 'декабря']

const pad = (__n: number): string => String(__n).padStart(2, '0')

const __codes__: { code: string; name: string; get: (_d: Date) => string }[] = [
  { code: '%H', name: 'Часы', get: (_d) => pad(_d.getHours()) },
  { code: '%M', name: 'Минуты', get: (_d) => pad(_d.getMinutes()) },
  { code: '%h', name: 'Часы до 12', get: (_d) => String(_d.getHours() % 12 || 12) },
  { code: '%P', name: 'AM или PM', get: (_d) => (_d.getHours() < 12 ? 'AM' : 'PM') },
  { code: '%D', name: 'День недели', get: (_d) => __full__[_d.getDay()] },
  { code: '%DS', name: 'День недели коротко', get: (_d) => __days__[_d.getDay()] },
  { code: '%d', name: 'Число', get: (_d) => pad(_d.getDate()) },
  { code: '%m', name: 'Месяц числом', get: (_d) => pad(_d.getMonth() + 1) },
  { code: '%B', name: 'Месяц словом', get: (_d) => __months__[_d.getMonth()] },
  { code: '%BS', name: 'Месяц коротко', get: (_d) => __mons__[_d.getMonth()] },
  { code: '%y', name: 'Год, две цифры', get: (_d) => String(_d.getFullYear()).slice(2) },
  { code: '%Y', name: 'Год', get: (_d) => String(_d.getFullYear()) },
  { code: '%E', name: 'Часы-эмодзи', get: (_d) => String.fromCodePoint((_d.getMinutes() >= 30 ? 0x1f55c : 0x1f550) + ((_d.getHours() + 11) % 12)) },
]

const __ready__ = ['[%H:%M]', '%h:%M %P', '%DS %H:%M', '%d.%m %H:%M', '%E %H:%M', '%D, %d %B']

const __re__ = /%(DS|BS|[HMhPDdmByYE])/g

const stamp = (__tpl: string, __d: Date): string => {
  const _out = __tpl.replace(__re__, (_m) => __codes__.find((_c) => _c.code === _m)?.get(__d) ?? _m)
  return _out.trim() ? [..._out].slice(0, 40).join('') : `${pad(__d.getHours())}:${pad(__d.getMinutes())}`
}

const preview = (__c: _Cfg): string => {
  const _tag = stamp(__c.tpl ?? '', new Date())
  const _f = (__c.first || who().name).trim()
  const _l = (__c.last ?? '').trim()
  if (__c.pos === 'before') return `${_tag} ${_f} ${_l}`.trim()
  return `${_f} ${_l} ${_tag}`.replace(/ +/g, ' ').trim()
}

export const clockPage = (__m: Mod): HTMLElement => {
  const _u = who()
  const _now = new Date()
  const _name = h('b', {})
  const _pv = h('div', { class: 'pv' }, ava(_u.id, _u.name), h('div', {}, _name, h('span', {}, 'так тебя увидят в Telegram')))
  const paint = (): void => void (_name.textContent = preview(__m.cfg))
  const put = (__k: string, __v: string, __save = true): void => {
    __m.cfg = { ...__m.cfg, [__k]: __v }
    paint()
    if (__save) void saveMod(__m)
  }
  const _tpl = h('input', { class: 'field tpl', type: 'text', maxlength: '40', placeholder: '[%H:%M]', value: __m.cfg.tpl ?? '', autocapitalize: 'off', autocomplete: 'off', spellcheck: 'false' }) as HTMLInputElement
  _tpl.addEventListener('input', () => put('tpl', _tpl.value, false))
  _tpl.addEventListener('change', () => void saveMod(__m))
  const insert = (__t: string): void => {
    const _a = _tpl.selectionStart ?? _tpl.value.length
    const _b = _tpl.selectionEnd ?? _a
    const _v = (_tpl.value.slice(0, _a) + __t + _tpl.value.slice(_b)).slice(0, 40)
    _tpl.value = _v
    const _pos = Math.min(_a + __t.length, _v.length)
    _tpl.setSelectionRange(_pos, _pos)
    buzz('tick')
    put('tpl', _v)
  }
  const _codes = __codes__.map((_c) =>
    row({ title: _c.name, value: _c.get(_now), right: h('code', { class: 'tpl-code' }, _c.code), onclick: () => insert(_c.code) }))
  const _chips = h('div', { class: 'tpl-chips' }, ...__ready__.map((_r) => {
    const _b = h('button', { type: 'button' }, stamp(_r, _now))
    _b.addEventListener('click', () => {
      buzz('tick')
      _tpl.value = _r
      put('tpl', _r)
    })
    return _b
  }))
  const field = (__k: string, __title: string): HTMLElement => {
    const _i = h('input', { class: 'opt-in', type: 'text', maxlength: '64', placeholder: 'как в Telegram', value: __m.cfg[__k] ?? '' }) as HTMLInputElement
    _i.addEventListener('input', () => put(__k, _i.value, false))
    _i.addEventListener('change', () => void saveMod(__m))
    return h('label', { class: 'line' }, h('span', { class: 'row-main' }, h('span', { class: 'row-t' }, __title)), _i)
  }
  const _where: HTMLElement[] = (__m.opts.find((_o) => _o.key === 'pos')?.choices ?? []).map((_c) => {
    const _r = row({ title: _c.label, check: true, cls: __m.cfg.pos === _c.v ? 'on' : '', onclick: () => {
      for (const _x of _where) _x.classList.toggle('on', _x === _r)
      put('pos', _c.v)
    } })
    return _r
  })
  const _in = h('div', {},
    h('div', { class: 'gh' }, h('h2', {}, 'Шаблон')),
    _tpl,
    _chips,
    group('Коды · нажми, чтобы вставить', null, ..._codes),
    group('Где показывать', null, ..._where),
    group('Имя', null, field('first', 'Имя'), field('last', 'Фамилия')))
  const _extra = h('div', { class: 'feat-x' + (__m.on ? ' on' : '') }, _in)
  const _main = sw(__m.on, async (_v) => {
    __m.on = _v
    _extra.classList.toggle('on', _v)
    if (!(await saveMod(__m)) && _v) {
      __m.on = false
      _main.classList.remove('on')
      _extra.classList.remove('on')
    }
    paint()
  })
  paint()
  return page('feat',
    head('clock', 'Время в имени', 'Часы рядом с твоим именем. Обновляются каждую минуту.'),
    _pv,
    box(row({ title: 'Включено', right: _main, onclick: () => _main.flip() })),
    _extra,
    howRow('/features/clock'),
  )
}
