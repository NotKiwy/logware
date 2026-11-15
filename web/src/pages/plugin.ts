import type { Me, Mod } from '../api/api'
import { go } from '../app/router'
import { load, mods } from '../app/store'
import { back, buzz, tint } from '../app/tg'
import { h } from '../ui/dom'
import { C, I } from '../ui/icons'
import { K, ask, box, empty, group, page, picIc, row, sq, svg } from '../ui/kit'
import { skPlugin } from '../ui/skel'
import { connect } from './help'
import { by, docsOf, icon, proTag, saveMod } from './mods'

const __rights__: Record<string, { name: string; ic: string; has: (_c: NonNullable<Me['conn']>) => boolean }> = {
  reply: { name: 'Отвечать на сообщения', ic: I.chat, has: (_c) => _c.reply },
  read: { name: 'Читать сообщения', ic: I.eye, has: (_c) => !!_c.read },
  delete: { name: 'Удалять сообщения', ic: I.trash, has: (_c) => _c.delAll },
  name: { name: 'Изменять имя', ic: I.person, has: (_c) => !!_c.name },
}

const options = (__m: Mod): HTMLElement[] => __m.opts.filter((_o) => !_o.text && !_o.hide).map((_o) => {
  const _rows: HTMLElement[] = []
  for (const _c of _o.choices ?? []) {
    const _r = row({ title: _c.label, check: true, cls: __m.cfg[_o.key] === _c.v ? 'on' : '', onclick: () => {
      __m.cfg = { ...__m.cfg, [_o.key]: _c.v }
      for (const _x of _rows) _x.classList.toggle('on', _x === _r)
      void saveMod(__m)
    } })
    _rows.push(_r)
  }
  return h('div', {}, group(_o.name, null, ..._rows))
})

const rights = (__m: Mod, __me: Me): Node[] => {
  const _ps = __m.perms ?? []
  const _safe = 'У плагинов нет доступа к твоему ключу и сохранённым сообщениям.'
  if (_ps.length === 0) return [group('Права бота', null, row({ lead: sq(I.shield, C.green), title: 'Дополнительные права не нужны' })), h('p', { class: 'note pg-note' }, _safe)]
  let _miss = 0
  const _rows = _ps.map((_p) => {
    const _r = __rights__[_p.right]
    const _ok = !!__me.conn && !!_r && _r.has(__me.conn)
    if (!_ok) _miss++
    let _end: HTMLElement
    if (_ok) {
      _end = h('span', { class: 'pg-yes', 'aria-label': 'Выдано' })
      _end.append(svg(I.check))
    } else {
      _end = h('button', { class: 'pg-give', type: 'button' }, 'Выдать')
      _end.addEventListener('click', () => {
        buzz('tap')
        connect(__me.bot, !!__me.conn?.active)
      })
    }
    return row({ lead: sq(_r?.ic ?? I.shield, _ok ? C.green : C.orange), title: _r?.name ?? _p.right, sub: _p.why, right: _end })
  })
  return [h('div', { class: 'pg-rights' }, group('Права бота', null, ..._rows)), h('p', { class: 'note pg-note' }, _miss ? 'Без этих прав плагин не сможет работать. ' + _safe : _safe)]
}

const pluginPage = (__m: Mod, __me: Me): HTMLElement => {
  const _act = h('button', { class: 'v-btn pg-act', type: 'button' })
  const _st = h('span', { class: 'pg-state' })
  _st.append(svg(I.check), 'Установлен')
  const _set = h('div', { class: 'feat-x' + (__m.on ? ' on' : '') }, h('div', {}, ...options(__m)))
  const paint = (): void => {
    _act.replaceChildren(__m.soon ? 'Скоро' : __m.on ? 'Удалить' : 'Установить')
    _act.classList.toggle('ghost', !!__m.soon)
    _act.classList.toggle('main', !__m.on && !__m.soon)
    _act.classList.toggle('danger', __m.on)
    _act.disabled = !!__m.soon || (!__m.ready && !__m.on)
    _st.hidden = !__m.on
    _set.classList.toggle('on', __m.on)
  }
  const flip = async (__v: boolean): Promise<void> => {
    _act.classList.add('busy')
    _act.replaceChildren(h('span', { class: 'spin' }))
    const _was = __m.on
    __m.on = __v
    const _ok = await saveMod(__m)
    if (!_ok) __m.on = _was
    _act.classList.remove('busy')
    paint()
    if (_ok) {
      buzz(__v ? 'success' : 'tick')
      _act.classList.remove('pop')
      void _act.offsetWidth
      _act.classList.add('pop')
    }
  }
  _act.addEventListener('click', async () => {
    if (_act.disabled || _act.classList.contains('busy')) return
    buzz('tap')
    if (!__m.on) return void flip(true)
    if (!(await ask('Удалить плагин?', `«${__m.name}» перестанет работать. Его можно будет снова установить из магазина.`, 'Удалить', true))) return
    void flip(false)
  })
  paint()
  const _info = [
    row({ title: 'Разработчик', right: by(__m) }),
    __m.ver ? row({ title: 'Версия', value: __m.ver }) : null,
  ].filter(Boolean) as HTMLElement[]
  return page('feat',
    h('div', { class: 'pg-head' },
      picIc(__m.pic, icon(__m.icon), __m.color, true),
      h('div', {}, h('span', { class: 'pg-name' }, __m.name, __m.on ? null : proTag()), by(__m), _st)),
    _act,
    __m.ready || __m.soon ? null : h('p', { class: 'note err' }, 'Плагин временно недоступен'),
    group('Описание', null, h('p', { class: 'pg-long' }, __m.long || __m.desc)),
    ...rights(__m, __me),
    _set,
    group('Информация', null, ..._info),
    box(row({ lead: sq(I.help, C.blue), title: 'Как это работает', chev: true, onclick: () => docsOf(__m) })),
  )
}

export const plugin = async (__root: HTMLElement, __id: string): Promise<void> => {
  __root.replaceChildren(skPlugin())
  const _me = await load()
  if (!_me) return empty(__root, I.cloud, K.gray, 'Нет связи', 'Не получилось связаться с сервером. Закрой и открой ещё раз.')
  const _m = mods().find((_x) => _x.id === __id && !_x.hidden)
  if (!_m) return go('/mods')
  tint(null)
  back(() => (history.length > 1 ? history.back() : go('/mods')))
  __root.replaceChildren(pluginPage(_m, _me))
  scrollTo(0, 0)
}
