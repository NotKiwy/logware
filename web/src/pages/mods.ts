import * as api from '../api/api'
import type { Mod } from '../api/api'
import { go } from '../app/router'
import { pro } from '../app/store'
import { buzz, docs } from '../app/tg'
import { h } from '../ui/dom'
import { I } from '../ui/icons'
import { box, page, picIc, row, seg, sheet, svg, toast } from '../ui/kit'
import { head, howRow } from './feat'
import { lot } from '../ui/lottie'

export const icon = (__name: string): string => (I as Record<string, string>)[__name] ?? I.bolt

export const proTag = (): HTMLElement | null => (pro() ? null : h('span', { class: 'pro-tag' }, 'PRO'))

export const needPro = (): void => {
  buzz('warning')
  const _b = h('button', { class: 'v-btn main', type: 'button' }, 'Подключить Pro')
  const _s = sheet(null,
    h('div', { class: 'sh-lot' }, lot('diamond', 120), h('h2', {}, 'Доступно в Pro')),
    h('p', { class: 'sheet-text' }, 'Плагины и время в имени работают с подпиской Pro. Удалённые, изменённые и одноразовые сообщения остаются бесплатными.'),
    _b)
  _b.addEventListener('click', () => {
    buzz('tap')
    _s.close()
    go('/f/plans?to=pro')
  })
}

export const saveMod = async (__m: Mod): Promise<boolean> => {
  if (__m.on && !pro()) {
    __m.on = false
    needPro()
    return false
  }
  if (__m.on && __m.opts.some((_o) => _o.key === 'tz')) __m.cfg = { ...__m.cfg, tz: Intl.DateTimeFormat().resolvedOptions().timeZone }
  try {
    const _r = await api.setMod(__m.id, __m.on, __m.cfg)
    __m.on = _r.on
    __m.cfg = _r.cfg
    return true
  } catch (_e) {
    const _t = (_e as Error).message
    if (_t.trim() === 'need pro') {
      __m.on = false
      needPro()
      return false
    }
    toast(_t && _t.length < 120 && !/^\d+$/.test(_t) && _t !== 'internal error' ? _t.trim() : 'Не получилось сохранить', 'err')
    return false
  }
}

export const by = (__m: Mod): HTMLElement => {
  const _b = h('span', { class: 'pg-by' }, __m.author)
  if (__m.check) _b.append(svg(I.badge))
  return _b
}

const item = (__m: Mod): HTMLElement => {
  const _tag = __m.soon ? h('span', { class: 'pg-tag' }, 'Скоро') : null
  const _r = h('button', { class: 'line app pg-item', type: 'button' },
    picIc(__m.pic, icon(__m.icon), __m.color),
    h('span', { class: 'row-main' }, h('span', { class: 'pg-name' }, __m.name, proTag()), by(__m), h('span', { class: 'row-s' }, __m.desc)),
    _tag)
  _r.append(svg(I.chev))
  _r.addEventListener('click', () => {
    buzz('tick')
    go('/p/' + __m.id)
  })
  return _r
}

let __tab: 'own' | 'store' = 'store'

export const modsPage = (__list: Mod[]): HTMLElement => {
  const _body = h('div', {})
  let _q = ''
  let _f = 'all'
  const _filters = [{ v: 'all', label: 'Все' }, { v: 'check', label: 'Проверенные авторы' }, ...(__list.some((_m) => _m.soon) ? [{ v: 'soon', label: 'Скоро' }] : [])]
  const _search = h('input', { class: 'pg-q', type: 'search', placeholder: 'Поиск плагинов', enterkeyhint: 'search' }) as HTMLInputElement
  const _sbox = h('label', { class: 'pg-search' })
  _sbox.append(svg(I.search), _search)
  const _chips = h('div', { class: 'pg-chips' }, ..._filters.map((_x) => {
    const _b = h('button', { type: 'button', class: _x.v === _f ? 'on' : '' }, _x.label)
    _b.addEventListener('click', () => {
      buzz('tick')
      _f = _x.v
      for (const _c of _chips.children) _c.classList.toggle('on', _c === _b)
      list()
    })
    return _b
  }))
  const _list = h('div', {})
  const list = (): void => {
    const _q2 = _q.trim().toLowerCase()
    const _xs = __list.filter((_m) =>
      (_f === 'all' || (_f === 'check' && _m.check) || (_f === 'soon' && _m.soon)) &&
      (!_q2 || (_m.name + ' ' + _m.desc + ' ' + _m.author).toLowerCase().includes(_q2)))
    _list.replaceChildren(_xs.length
      ? box(..._xs.map(item))
      : h('p', { class: 'pg-empty' }, 'Ничего не нашлось. Попробуй другой запрос.'))
  }
  _search.addEventListener('input', () => {
    _q = _search.value
    list()
  })
  const store = (): Node[] => {
    list()
    return [_sbox, _chips, _list, h('p', { class: 'note pg-note' }, 'Скоро здесь появятся плагины от других людей. Мы проверим каждый перед публикацией.')]
  }
  const own = (): Node[] => {
    const _xs = __list.filter((_m) => _m.on)
    if (_xs.length === 0) {
      const _go = h('button', { class: 'link', type: 'button' }, 'Открыть магазин')
      _go.addEventListener('click', () => (_seg.querySelectorAll('button')[1] as HTMLButtonElement).click())
      return [h('div', { class: 'pg-none' },
        h('b', {}, 'Пока ничего не установлено'),
        h('p', {}, 'Загляни в магазин: там есть полезные плагины, которые ставятся в один клик.'),
        _go)]
    }
    return [box(..._xs.map((_m) => {
      return row({ lead: picIc(_m.pic, icon(_m.icon), _m.color), title: _m.name, sub: _m.author, chev: true, onclick: () => go('/p/' + _m.id) })
    }))]
  }
  const show = (__t: 'own' | 'store'): void => {
    __tab = __t
    _body.replaceChildren(...(__t === 'own' ? own() : store()))
  }
  const _n = __list.filter((_m) => _m.on).length
  const _seg = seg([{ v: 'own', label: _n ? `Установленные · ${_n}` : 'Установленные' }, { v: 'store', label: 'Магазин' }], __tab, (_v) => show(_v as 'own' | 'store'))
  show(__tab)
  return page('feat',
    head('bot', 'Плагины', 'Новые возможности бота в один клик. Включай только то, что нужно.'),
    _seg,
    _body,
    howRow('/plugins/'),
  )
}

export const __docs__: Record<string, string> = { time: '/plugins/example#what' }

export const docsOf = (__m: Mod): void => docs(__docs__[__m.id] ?? '/plugins/')
