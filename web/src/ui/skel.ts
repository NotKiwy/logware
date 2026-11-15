import { h } from './dom'

type _Lead = 'sq' | 'ava' | 'app' | null

const bar = (__w: string | number, __h = 12, __r = 6, __cls = ''): HTMLElement =>
  h('span', { class: 'sk ' + __cls, style: `width:${typeof __w === 'number' ? __w + 'px' : __w};height:${__h}px;border-radius:${__r}px` })

const __widths__ = ['46%', '62%', '38%', '55%', '70%', '42%', '58%']

const line = (__n: number, __lead: _Lead, __sub: boolean, __val: boolean): HTMLElement => {
  const _lead = __lead === 'sq' ? bar(32, 32, 9) : __lead === 'ava' ? bar(32, 32, 16) : __lead === 'app' ? bar(44, 44, 12) : null
  const _main = h('span', { class: 'sk-main' }, bar(__widths__[__n % __widths__.length], 13), __sub ? bar(__widths__[(__n + 3) % __widths__.length], 10) : null)
  return h('div', { class: 'line sk-line' + (__lead ? ' lead' : '') + (__lead === 'app' ? ' app' : '') + (__sub ? ' tall' : '') }, _lead, _main, __val ? bar(38, 12) : null)
}

const rows = (__n: number, __lead: _Lead = 'sq', __sub = false, __val = false, __from = 0): HTMLElement =>
  h('div', { class: 'box sk-box' }, ...Array.from({ length: __n }, (_, _i) => line(_i + __from, __lead, __sub, __val)))

const gh = (__w = 92): HTMLElement => h('div', { class: 'gh' }, bar(__w, 14))

const head = (): HTMLElement =>
  h('div', { class: 'feat-h sk-head' }, bar(128, 128, 64), bar(170, 24, 8), bar('72%', 12), bar('54%', 12))

const wrap = (__cls: string, ...__kids: (Node | null)[]): HTMLElement => h('div', { class: 'app sk-page ' + __cls, 'aria-busy': 'true' }, ...__kids)

export const skHome = (): HTMLElement =>
  wrap('',
    h('div', { class: 'blue sk-hero', style: '--h1:#2a7af5;--h2:#2369ea;--h3:#1d5fdf' },
      h('div', { class: 'topbar' }, bar(44, 44, 22, 'on-blue'), h('div', { class: 'sk-main' }, bar(96, 14, 6, 'on-blue'), bar(64, 11, 6, 'on-blue'))),
      bar('100%', 0, 20, 'on-blue sk-card')),
    gh(),
    rows(5, 'sq', false, true),
    gh(54),
    rows(3, 'sq', false, false, 2),
    rows(1, 'sq', false, true, 4))

export const skMods = (): HTMLElement =>
  wrap('feat',
    head(),
    bar('100%', 40, 999),
    h('div', { class: 'sk-gap' }),
    bar('100%', 42, 14),
    h('div', { class: 'sk-chips' }, bar(52, 32, 999), bar(150, 32, 999), bar(70, 32, 999)),
    rows(3, 'app', true))

export const skSet = (): HTMLElement =>
  wrap('',
    gh(80), rows(1, 'sq', false, true),
    gh(70), rows(3, 'sq', false, true, 1),
    gh(64), rows(2, 'sq', false, false, 3),
    h('div', { class: 'sk-gap' }),
    rows(1, null, false, false, 5))

export const skFeat = (): HTMLElement =>
  wrap('feat', head(), rows(1, null, false, false, 2), gh(), rows(2, 'sq', true, false, 1))

export const skList = (): HTMLElement =>
  wrap('feat', head(), gh(), rows(4, 'ava', true, true), gh(70), rows(2, 'ava', true, true, 3))

export const skPlugin = (): HTMLElement =>
  wrap('feat',
    h('div', { class: 'pg-head' }, bar(84, 84, 24), h('div', { class: 'sk-main' }, bar(170, 22, 8), bar(80, 12))),
    bar('100%', 54, 999),
    h('div', { class: 'sk-gap' }),
    gh(),
    h('div', { class: 'box sk-text' }, bar('92%', 12), bar('86%', 12), bar('95%', 12), bar('60%', 12)),
    gh(90),
    rows(1, 'sq', true),
    gh(100),
    rows(2, null, false, true, 2))

export const skPlans = (): HTMLElement =>
  wrap('plan-pg',
    h('div', { class: 'blue pv-block sk-hero', style: '--h1:#2a7af5;--h2:#2369ea;--h3:#1d5fdf' },
      h('div', { class: 'pv-hold' }, bar('100%', 0, 20, 'on-blue sk-card')),
      h('div', { class: 'pv-dots' }, bar(18, 6, 3, 'on-blue'), bar(6, 6, 3, 'on-blue'), bar(6, 6, 3, 'on-blue'))),
    gh(),
    rows(5, 'sq'))

export const skView = (): HTMLElement =>
  wrap('bare view',
    h('div', { class: 'stage' },
      h('div', { class: 'pill-s head' }, bar(190, 26, 999)),
      h('div', { class: 'bub-row' }, bar('62%', 54, 20, 'sk-bub')),
      h('div', { class: 'bub-row' }, bar('78%', 76, 20, 'sk-bub')),
      h('div', { class: 'bub-row' }, bar('48%', 40, 20, 'sk-bub'))),
    h('div', { class: 'v-bar' }, bar('100%', 54, 999), h('div', { class: 'sk-link' }, bar(120, 14))))

export const skDev = (): HTMLElement =>
  wrap('dv',
    h('div', { class: 'dv-prof sk-head' }, bar(96, 96, 48), bar(140, 22, 8), bar(110, 12)),
    bar('100%', 64, 20),
    h('div', { class: 'sk-gap' }),
    bar('100%', 52, 999),
    gh(70), rows(2, 'sq', false, true),
    gh(110), rows(1, 'app', true),
    gh(64), rows(2, 'sq', false, false, 3))

