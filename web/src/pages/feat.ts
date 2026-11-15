import * as api from '../api/api'
import type { Prefs } from '../api/api'
import { go } from '../app/router'
import { load, mods } from '../app/store'
import { back, docs, tint } from '../app/tg'
import { h } from '../ui/dom'
import { C, I } from '../ui/icons'
import { K, box, empty, group, page, pick, row, sq, sw, toast } from '../ui/kit'
import { lot } from '../ui/lottie'
import { skFeat, skList, skPlans } from '../ui/skel'
import { clockPage } from './clock'
import { plansPage } from './plans'
import { more } from './more'

type _Flag = keyof Prefs

type _Opt = { key: _Flag; ic: string; color: string; title: string; sub?: string; pick?: { v: number; label: string }[] }

type _Feat = { key: _Flag; lot: string; name: string; text: string; extra: { title: string; opts: _Opt[] }[]; how: string }

export const __keeps__: { v: number; label: string }[] = [
  { v: 7, label: '7 дней' },
  { v: 30, label: '30 дней' },
  { v: 90, label: '90 дней' },
  { v: 365, label: '1 год' },
  { v: 0, label: 'Всегда' },
]

export const keepLabel = (__v: number): string => __keeps__.find((_k) => _k.v === __v)?.label ?? ''

export const __feats__: Record<string, _Feat> = {
  del: {
    key: 'del',
    lot: 'deleted',
    name: 'Удалённые',
    text: 'Когда собеседник удалит сообщение, я пришлю его тебе в чат со мной.',
    extra: [
      { title: 'Что присылать', opts: [{ key: 'delMe', ic: I.person, color: C.gray, title: 'Удалённые мной', sub: 'Сообщать и о том, что удалил я' }] },
      { title: 'Уведомления', opts: [{ key: 'delMute', ic: I.bell, color: C.red, title: 'Без звука', sub: 'Сообщение придёт тихо' }] },
      {
        title: 'Массовое удаление',
        opts: [
          { key: 'panic', ic: I.bolt, color: C.orange, title: 'Предупреждать', sub: 'Громко сообщу, если удалят много сразу' },
          { key: 'panicN', ic: I.trash, color: C.gray, title: 'Сколько это «много»', pick: [3, 5, 10, 20, 50].map((_v) => ({ v: _v, label: 'от ' + _v })) },
        ],
      },
    ],
    how: '/features/deleted',
  },
  edit: {
    key: 'edit',
    lot: 'edited',
    name: 'Изменённые',
    text: 'Когда сообщение изменят, покажу, что в нём было написано раньше.',
    extra: [
      { title: 'Что присылать', opts: [{ key: 'editMe', ic: I.person, color: C.gray, title: 'Изменённые мной', sub: 'Сообщать о том, что изменил я' }] },
      { title: 'Уведомления', opts: [{ key: 'editMute', ic: I.bell, color: C.orange, title: 'Без звука', sub: 'Сообщение придёт тихо' }] },
    ],
    how: '/features/edited',
  },
  media: {
    key: 'media',
    lot: 'image',
    name: 'Фото и видео',
    text: 'Сохраню фото, видео, голосовые и файлы до 20 МБ вместе с сообщением.',
    extra: [
      {
        title: 'Что сохранять',
        opts: [
          { key: 'photo', ic: I.photo, color: C.purple, title: 'Фото' },
          { key: 'video', ic: I.video, color: C.blue, title: 'Видео и кружки' },
          { key: 'voice', ic: I.mic, color: C.orange, title: 'Голосовые и музыка' },
          { key: 'files', ic: I.grid, color: C.indigo, title: 'Файлы и стикеры' },
        ],
      },
    ],
    how: '/features/media',
  },
  once: {
    key: 'once',
    lot: 'chameleon',
    name: 'Одноразовые',
    text: 'Фото и видео, которые можно открыть только один раз. Я сохраню их для тебя.',
    extra: [
      {
        title: 'Как сохранять',
        opts: [
          { key: 'onceAuto', ic: I.fire, color: C.pink, title: 'При ответе', sub: 'Ответь на фото чем угодно, и я его сохраню' },
          { key: 'onceHide', ic: I.eye, color: C.gray, title: 'Прятать .s', sub: 'Твоя команда сразу исчезнет из чата' },
        ],
      },
      { title: 'Уведомления', opts: [{ key: 'onceMute', ic: I.bell, color: C.pink, title: 'Без звука', sub: 'Сообщение придёт тихо' }] },
    ],
    how: '/features/once',
  },
}

export const howRow = (__path: string): HTMLElement =>
  box(row({ lead: sq(I.help, C.blue), title: 'Как это работает', chev: true, onclick: () => docs(__path) }))

export const head = (__lot: string, __name: string, __text: string): HTMLElement =>
  h('div', { class: 'feat-h' }, lot(__lot, 128), h('h1', {}, __name), h('p', {}, __text))

const flagPage = (__f: _Feat, __p: Prefs, __save: () => void): HTMLElement => {
  const flip = (__k: _Flag) => (_v: boolean) => {
    ;(__p as Record<string, unknown>)[__k] = _v
    __save()
  }
  const _in = h('div', {})
  const _extra = h('div', { class: 'feat-x' + (__p[__f.key] ? ' on' : '') }, _in)
  for (const _g of __f.extra) {
    _in.append(group(_g.title, null, ..._g.opts.map((_o) => {
      if (_o.pick) {
        const _opts = _o.pick
        const _lbl = () => _opts.find((_x) => _x.v === __p[_o.key])?.label ?? ''
        const _r: HTMLElement = row({ lead: sq(_o.ic, _o.color), title: _o.title, value: _lbl(), chev: true, onclick: async () => {
          const _v = await pick(_o.title, null, _opts, __p[_o.key] as number)
          if (_v == null) return
          ;(__p as Record<string, unknown>)[_o.key] = _v
          _r.querySelector('.val')!.textContent = _lbl()
          __save()
        } })
        return _r
      }
      const _s = sw(!!__p[_o.key], flip(_o.key))
      return row({ lead: sq(_o.ic, _o.color), title: _o.title, sub: _o.sub, right: _s, onclick: () => _s.flip() })
    })))
  }
  const _main = sw(!!__p[__f.key], (_v) => {
    flip(__f.key)(_v)
    _extra.classList.toggle('on', _v)
  })
  return page('feat',
    head(__f.lot, __f.name, __f.text),
    box(row({ title: 'Включено', right: _main, onclick: () => _main.flip() })),
    _extra,
    howRow(__f.how),
  )
}

const keepPage = (__p: Prefs, __save: () => void): HTMLElement => {
  const _rows = __keeps__.map((_k) => {
    const _r = row({ title: _k.label, check: true, cls: _k.v === __p.keep ? 'on' : '', onclick: () => {
      if (_k.v === __p.keep) return
      __p.keep = _k.v
      for (const _x of _rows) _x.classList.toggle('on', _x === _r)
      __save()
    } })
    return _r
  })
  return page('feat',
    head('clock', 'Сколько хранить', 'Всё, что старше этого срока, удаляется само.'),
    group('Срок', null, ..._rows),
    howRow('/features/keep'),
  )
}

export const feat = async (__root: HTMLElement, __key: string): Promise<void> => {
  __root.replaceChildren(__key === 'plans' ? skPlans() : __key === 'hist' || __key === 'ign' ? skList() : skFeat())
  const _me = await load()
  if (!_me) return empty(__root, I.cloud, K.gray, 'Нет связи', 'Не получилось связаться с сервером. Закрой и открой ещё раз.')
  if (!_me.reg) return go('/reg')
  const _p = _me.prefs
  const save = async () => {
    try {
      Object.assign(_p, await api.prefs(_p))
    } catch {
      toast('Не получилось сохранить', 'err')
    }
  }
  if (__key === 'hist' || __key === 'kw' || __key === 'ign') {
    tint(null)
    back(() => (history.length > 1 ? history.back() : go('/')))
    return more(__root, __key, _me)
  }
  if (__key === 'plans') {
    back(() => (history.length > 1 ? history.back() : go('/')))
    __root.replaceChildren(plansPage(_me))
    return void scrollTo(0, 0)
  }
  if (__key === 'clock') {
    const _m = mods().find((_x) => _x.id === 'clock')
    if (!_m) return go('/')
    tint(null)
    back(() => (history.length > 1 ? history.back() : go('/')))
    __root.replaceChildren(clockPage(_m))
    return void scrollTo(0, 0)
  }
  const _f = __feats__[__key]
  if (!_f && __key !== 'keep') return go('/')
  tint(null)
  back(() => history.length > 1 ? history.back() : go('/'))
  __root.replaceChildren(_f ? flagPage(_f, _p, () => void save()) : keepPage(_p, () => void save()))
  scrollTo(0, 0)
}
