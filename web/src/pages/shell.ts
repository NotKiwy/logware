import { go } from '../app/router'
import { banned, cached, load, mods } from '../app/store'
import { buzz, tint } from '../app/tg'
import { h } from '../ui/dom'
import { I } from '../ui/icons'
import { K, empty, svg } from '../ui/kit'
import { skHome, skMods, skSet } from '../ui/skel'
import { heroTop, homePage } from './home'
import { devSkel, devTab } from './dev'
import { modsPage } from './mods'
import { settingsPage } from './settings'

const __tabs__ = [
  { id: 'home', path: '/', label: 'Главная', ic: I.tHome },
  { id: 'mods', path: '/mods', label: 'Плагины', ic: I.tMods },
  { id: 'set', path: '/set', label: 'Настройки', ic: I.tSet },
]

const __dev__ = { id: 'dev', path: '/dev', label: 'Developer', ic: I.tDev }

let __has = false

const tabs = (): typeof __tabs__ => (__has ? [...__tabs__.slice(0, 2), __dev__, __tabs__[2]] : __tabs__)

let __prev = 0

let __pop = false

const tabbar = (__i: number): HTMLElement => {
  const _ts = tabs()
  const _bar = h('nav', { class: 'tabbar', style: `--n:${_ts.length};--i:${__prev}` }, h('span', { class: 'tab-ind', 'aria-hidden': 'true' }))
  _ts.forEach((_t, _n) => {
    const _on = _n === __i
    const _ic = h('span', { class: 'ic' + (_on && __pop ? ' pop' : ''), 'data-tab': _t.id })
    _ic.append(svg(_t.ic))
    const _b = h('button', { class: _on ? 'on' : '', type: 'button', 'aria-current': String(_on) }, _ic, _t.label)
    _b.addEventListener('click', () => {
      buzz('tick')
      if (_on) {
        _ic.classList.remove('pop')
        void _ic.offsetWidth
        _ic.classList.add('pop')
        return scrollTo({ top: 0, behavior: 'smooth' })
      }
      __pop = true
      go(_t.path)
    })
    _bar.append(_b)
  })
  requestAnimationFrame(() => requestAnimationFrame(() => _bar.style.setProperty('--i', String(__i))))
  __prev = __i
  return _bar
}

export const shell = async (__root: HTMLElement, __path: string): Promise<void> => {
  const _sk: Record<string, () => HTMLElement> = { mods: skMods, set: skSet, dev: devSkel, home: skHome }
  if (!cached()) {
    const _n = Math.max(0, tabs().findIndex((_t) => _t.path === __path))
    __root.replaceChildren(_sk[tabs()[_n].id](), tabbar(_n))
  }
  const _me = await load()
  if (!_me && banned()) return empty(__root, I.shield, K.gray, 'Доступ закрыт', 'Твой аккаунт заблокирован в Logware.')
  if (!_me) return empty(__root, I.cloud, K.gray, 'Нет связи', 'Не получилось связаться с сервером. Закрой и открой ещё раз.')
  if (!_me.reg) return go('/reg')
  __has = _me.dev
  if (__path === '/dev' && !_me.dev) return go('/f/plans')
  const _i = Math.max(0, tabs().findIndex((_t) => _t.path === __path))
  const _id = tabs()[_i].id
  if (_id === 'dev' && __root.querySelector('.dv') === null) __root.replaceChildren(devSkel(), tabbar(_i))
  const _page = _id === 'mods' ? modsPage(mods().filter((_m) => !_m.hidden)) : _id === 'set' ? settingsPage(_me) : _id === 'dev' ? await devTab(_me) : homePage(_me)
  if (location.pathname !== __path) return
  tint(_id === 'home' ? heroTop(_me) : null)
  __root.replaceChildren(_page, tabbar(_i))
  scrollTo(0, 0)
}
