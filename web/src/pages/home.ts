import type { Me, Prefs } from '../api/api'
import { go } from '../app/router'
import { buzz, path, who } from '../app/tg'
import { ava } from '../ui/ava'
import { h } from '../ui/dom'
import { C, I } from '../ui/icons'
import { box, group, page, row, sq } from '../ui/kit'
import { keepLabel } from './feat'
import { connect } from './help'
import { __tiers__, card } from './plans'
import { mods, pro } from '../app/store'
import { proTag } from './mods'

const __flags__: { key: keyof Prefs; ic: string; color: string; name: string }[] = [
  { key: 'del', ic: I.trash, color: C.red, name: 'Удалённые' },
  { key: 'edit', ic: I.pencil, color: C.orange, name: 'Изменённые' },
  { key: 'media', ic: I.photo, color: C.purple, name: 'Фото и видео' },
  { key: 'once', ic: I.fire, color: C.pink, name: 'Одноразовые' },
]

export const heroTop = (__me: Me): string => (__tiers__[__me.plan] ?? __tiers__.free).hero[0]

const plan = (__me: Me): HTMLElement => {
  const _card = card(__me.plan in __tiers__ ? __me.plan : 'free', __me)
  _card.addEventListener('click', () => {
    buzz('tap')
    go('/f/plans')
  })
  return _card
}

const setup = (__me: Me): HTMLElement | null => {
  const _c = __me.conn
  if (_c?.active && _c.reply && _c.delAll) return null
  const _one = !!_c?.active
  return box(row({
    lead: sq(I.plug, C.blue),
    title: _one ? 'Выдай боту права' : 'Подключи бота',
    sub: path(__me.bot),
    chev: true,
    onclick: () => connect(__me.bot, _one),
  }))
}

export const homePage = (__me: Me): HTMLElement => {
  const _p = __me.prefs
  const _u = who()
  const _h = (__tiers__[__me.plan] ?? __tiers__.free).hero
  const _rows = __flags__.map((_f) =>
    row({ lead: sq(_f.ic, _f.color), title: _f.name, value: _p[_f.key] ? 'Вкл' : 'Выкл', chev: true, onclick: () => go('/f/' + _f.key) }))
  const _clock = mods().find((_m) => _m.id === 'clock')
  if (_clock) {
    const _r = row({ lead: sq(I.clock, C.blue), title: 'Время в имени', value: !pro() ? undefined : _clock.on ? 'Вкл' : 'Выкл', chev: true, onclick: () => go('/f/clock') })
    const _tag = proTag()
    const _t = _r.querySelector('.row-t')
    if (_tag && _t) {
      _t.classList.add('tagged')
      _t.append(_tag)
    }
    _rows.push(_r)
  }
  const _keep = row({ lead: sq(I.clock, C.indigo), title: 'Сколько хранить', value: keepLabel(_p.keep), chev: true, onclick: () => go('/f/keep') })
  const _set = setup(__me)
  if (_set) _set.style.marginTop = '20px'
  return page('',
    h('div', { class: 'blue', style: `--h1:${_h[0]};--h2:${_h[1]};--h3:${_h[2]}` },
      h('div', { class: 'rays' }),
      h('div', { class: 'topbar' }, ava(_u.id, _u.name, 'ava'), h('div', {}, h('b', {}, _u.name), _u.nick && h('span', {}, _u.nick))),
      plan(__me)),
    _set,
    group('Основное', null, ..._rows),
    group('Ещё', null,
      row({ lead: sq(I.hist, C.blue), title: 'История', chev: true, onclick: () => go('/f/hist') }),
      row({ lead: sq(I.tag, C.indigo), title: 'Ключевые слова', chev: true, onclick: () => go('/f/kw') }),
      row({ lead: sq(I.eyeOff, C.gray), title: 'Чаты без слежки', chev: true, onclick: () => go('/f/ign') })),
    box(_keep),
  )
}
