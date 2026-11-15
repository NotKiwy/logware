import * as api from '../api/api'
import type { Dev, DevOne, DevPlug, DevReq, Me, PlugIn } from '../api/api'
import { go } from '../app/router'
import { back, buzz, docs, web, who } from '../app/tg'
import { ava } from '../ui/ava'
import { h } from '../ui/dom'
import { C, I } from '../ui/icons'
import { K, box, empty, group, page, picIc, row, sheet, sq, svg, toast } from '../ui/kit'
import { lot } from '../ui/lottie'
import { skDev } from '../ui/skel'

const __errs__: Record<string, string> = {
  'bad nick': 'Ник: латиница, цифры и _, от 3 до 24 символов, начинается с буквы',
  'nick taken': 'Этот ник уже занят',
  'bio too long': 'Слишком длинное описание',
  'no nick': 'Сначала придумай ник в профиле',
  'too many': 'На проверке уже пять заявок. Дождись ответа',
  'busy': 'У плагина уже есть заявка на проверке',
  'bad slug': 'ID: латиница в нижнем регистре, цифры и дефис, от 2 до 24 символов',
  'slug taken': 'Такой ID уже занят',
  'bad pic': 'Аватарка должна быть картинкой PNG или JPG',
  'same ver': 'Укажи новую версию',
  'bad repo': 'Нужна ссылка вида github.com/ник/плагин',
  'bad fields': 'Проверь название и описание',
}

const err = (__e: unknown, __def: string): string => __errs__[(__e as Error).message.trim()] ?? __def

const __st__: Record<DevReq['status'], [string, string]> = {
  wait: ['На проверке', 'wait'],
  ok: ['Одобрено', 'ok'],
  no: ['Отклонено', 'no'],
}

const __kinds__: Record<DevReq['kind'], string> = { new: 'Первая версия', update: 'Обновление', edit: 'Изменение' }

const pill = (__t: string, __c: string): HTMLElement => h('span', { class: 'dv-st ' + __c }, __t)

const state = (__p: DevPlug): [string, string] =>
  __p.last === 'wait' ? ['На проверке', 'wait'] : __p.live ? ['Опубликован', 'ok'] : __p.last === 'no' ? ['Отклонён', 'no'] : ['Черновик', 'wait']

const date = (__t: number, __time = false): string =>
  new Date(__t * 1000).toLocaleString('ru', __time ? { day: 'numeric', month: 'long', hour: '2-digit', minute: '2-digit' } : { day: 'numeric', month: 'long', year: 'numeric' })

const plural = (__n: number, __a: string, __b: string, __c: string): string => {
  const _m = __n % 100
  const _d = __n % 10
  return _m > 10 && _m < 20 ? __c : _d === 1 ? __a : _d > 1 && _d < 5 ? __b : __c
}

const plugIc = (__p: { pic: string }, __big = false): HTMLElement => picIc(__p.pic, I.code, C.indigo, __big)

let __cache: Dev | null = null

const edit = (__d: Dev, __done: () => void): void => {
  const _nick = h('input', { class: 'field', type: 'text', maxlength: '24', placeholder: 'Ник латиницей', value: __d.nick, enterkeyhint: 'done', autocapitalize: 'off', autocomplete: 'off', spellcheck: 'false' }) as HTMLInputElement
  const _bio = h('textarea', { class: 'field area', maxlength: '120', rows: '3', placeholder: 'Пара слов о себе' }) as HTMLTextAreaElement
  _bio.value = __d.bio
  const _e = h('p', { class: 'note err' })
  _e.hidden = true
  const _btn = h('button', { class: 'v-btn main', type: 'button' }, 'Сохранить')
  const _s = sheet('Профиль', h('p', { class: 'sheet-text' }, 'Ник видят все в магазине рядом с твоими плагинами. Латиница, цифры и _, от 3 до 24 символов.'), _nick, _bio, _e, _btn)
  const save = async (): Promise<void> => {
    if (_btn.classList.contains('busy')) return
    const _n = _nick.value.trim().replace(/^@+/, '')
    if (!_n) {
      buzz('error')
      _e.textContent = 'Сначала придумай ник'
      _e.hidden = false
      toast('Сначала придумай ник', 'err')
      return _nick.focus()
    }
    ;(document.activeElement as HTMLElement | null)?.blur()
    _btn.classList.add('busy')
    try {
      const _r = await api.setDev(_n, _bio.value.trim())
      __d.nick = _r.nick
      __d.bio = _r.bio
      buzz('success')
      toast('Сохранено')
      _s.close()
      __done()
    } catch (_x) {
      buzz('error')
      _e.textContent = err(_x, 'Не получилось сохранить')
      _e.hidden = false
      toast(_e.textContent, 'err')
    } finally {
      _btn.classList.remove('busy')
    }
  }
  _btn.addEventListener('click', () => void save())
  _nick.addEventListener('keydown', (_k) => _k.key === 'Enter' && void save())
  _nick.addEventListener('input', () => (_e.hidden = true))
  setTimeout(() => _nick.focus(), 300)
}

const prof = (__d: Dev): HTMLElement => {
  const _u = who()
  const _name = h('b', {}, __d.nick || 'Без ника')
  if (__d.nick) _name.append(svg(I.badge))
  return h('div', { class: 'dv-prof' },
    ava(_u.id, _u.name, 'm-ava dv-ava'),
    _name,
    h('span', {}, __d.bio || 'Разработчик Logware'))
}

const stats = (__d: Dev): HTMLElement => {
  const _users = __d.plugs.reduce((_a, _p) => _a + _p.users, 0)
  const _live = __d.plugs.filter((_p) => _p.live).length
  const _wait = __d.plugs.filter((_p) => _p.last === 'wait').length
  const cell = (__n: number, __t: string): HTMLElement => h('div', {}, h('b', {}, String(__n)), h('span', {}, __t))
  return h('div', { class: 'dv-stats' }, cell(_live, 'в магазине'), cell(_users, 'установок'), cell(_wait, 'на проверке'))
}

const devPage = (__d: Dev): HTMLElement => {
  const _root = h('div', {})
  const paint = (): void => {
    const _plugs = __d.plugs.length
      ? group('Мои плагины', null, ...__d.plugs.map((_p) => {
        const [_t, _c] = state(_p)
        return row({ lead: plugIc(_p), title: _p.name, sub: `${_p.ver ? 'v' + _p.ver + ' · ' : ''}${_p.users} ${plural(_p.users, 'установка', 'установки', 'установок')}`, right: pill(_t, _c), onclick: () => go('/dev/p/' + _p.id) })
      }))
      : h('div', { class: 'pg-none dv-empty' }, h('b', {}, 'Пока нет плагинов'), h('p', {}, 'Нажми «Новый плагин», чтобы отправить первый. Здесь будут его статус, версии и история проверок.'))
    const _new = h('button', { class: 'v-btn main dv-new', type: 'button' }, 'Новый плагин')
    _new.addEventListener('click', () => {
      buzz('tap')
      if (!__d.nick) return edit(__d, paint)
      go('/dev/new')
    })
    _root.replaceChildren(
      prof(__d),
      stats(__d),
      _new,
      _plugs,
      group('Профиль', null,
        row({ lead: sq(I.tag, C.blue), title: 'Ник', value: __d.nick || 'Не задан', chev: true, onclick: () => edit(__d, paint) }),
        row({ lead: sq(I.pencil, C.purple), title: 'О себе', value: __d.bio ? 'Есть' : 'Пусто', chev: true, onclick: () => edit(__d, paint) })),
      group('Помощь', null,
        row({ lead: sq(I.code, C.indigo), title: 'Как написать плагин', chev: true, onclick: () => docs('/plugins/dev') }),
        row({ lead: sq(I.shield, C.green), title: 'Как мы проверяем', chev: true, onclick: () => docs('/plugins/publish#review') })),
    )
  }
  paint()
  return page('dv', _root)
}

export const devTab = async (__me: Me): Promise<HTMLElement> => {
  if (!__me.dev) return page('dv')
  try {
    __cache = await api.dev()
  } catch {
    const _w = page('dv')
    empty(_w, I.cloud, K.gray, 'Нет связи', 'Не получилось загрузить кабинет. Попробуй ещё раз.')
    return _w
  }
  return devPage(__cache)
}

export const devSkel = skDev

const detail = (__q: DevReq): void => {
  const [_t, _c] = __st__[__q.status]
  const _rows: HTMLElement[] = [row({ title: 'Отправлено', value: date(__q.created, true) })]
  if (__q.status !== 'wait') _rows.push(row({ title: 'Проверено', value: date(__q.done, true) }))
  if (__q.ver) _rows.push(row({ title: 'Версия', value: __q.ver }))
  _rows.push(row({ title: 'Название', value: __q.name }))
  const _pic = __q.pic ? h('div', { class: 'dv-pic' }, picIc(api.reqPic(__q.id), I.code, C.indigo, true), h('span', {}, 'Новая аватарка')) : null
  sheet(null,
    h('div', { class: 'ad-head' }, h('b', {}, __kinds__[__q.kind]), h('span', {}, date(__q.created, true)), pill(_t, _c)),
    __q.status === 'no' && __q.reason ? h('div', { class: 'ad-note no' }, h('small', {}, 'Причина отказа'), h('p', {}, __q.reason)) : null,
    _pic,
    box(..._rows),
    __q.about ? h('div', { class: 'ad-note' }, h('small', {}, 'Описание'), h('p', {}, __q.about)) : null,
    __q.kind !== 'edit' && __q.repo ? box(row({ lead: sq(I.code, C.indigo), title: 'Код', sub: __q.repo.replace(/^https:\/\//, ''), chev: true, onclick: () => web(__q.repo) })) : null,
    __q.note ? h('div', { class: 'ad-note' }, h('small', {}, __q.kind === 'update' ? 'Что нового' : 'Комментарий для проверки'), h('p', {}, __q.note)) : null)
}

export const devPlug = async (__root: HTMLElement, __id: number): Promise<void> => {
  back(() => (history.length > 1 ? history.back() : go('/dev')))
  __root.replaceChildren(skDev())
  let _p: DevOne
  try {
    _p = await api.devPlug(__id)
  } catch {
    return empty(__root, I.cloud, K.gray, 'Не нашли плагин', 'Попробуй открыть его ещё раз из кабинета.')
  }
  if (location.pathname !== '/dev/p/' + __id) return
  const [_t, _c] = state(_p)
  const _busy = _p.last === 'wait'
  const _acts = _busy
    ? [h('p', { class: 'note' }, 'Пока заявка на проверке, новую отправить нельзя. Ответ придёт в чат с ботом.')]
    : [group('Управление', null,
      row({ lead: sq(I.bolt, C.blue), title: _p.live ? 'Выпустить обновление' : 'Отправить снова', sub: _p.live ? 'Новая версия кода' : 'Исправленная версия', chev: true, onclick: () => go(`/dev/p/${__id}/update`) }),
      row({ lead: sq(I.pencil, C.purple), title: 'Изменить карточку', sub: 'Название, описание, аватарка', chev: true, onclick: () => go(`/dev/p/${__id}/edit`) }))]
  const _no = !_busy && _p.last === 'no' && _p.why
    ? h('div', { class: 'ad-note no' }, h('small', {}, 'Последняя заявка отклонена'), h('p', {}, _p.why))
    : null
  __root.replaceChildren(page('feat dv-plug',
    h('div', { class: 'ad-head dv-top' }, plugIc(_p, true), h('b', {}, _p.name), h('span', {}, h('code', {}, _p.slug), _p.ver ? ' · v' + _p.ver : ''), pill(_t, _c)),
    _p.about ? h('p', { class: 'ad-text' }, _p.about) : null,
    h('div', { class: 'dv-stats two' },
      h('div', {}, h('b', {}, _p.ver || '-'), h('span', {}, 'версия')),
      h('div', {}, h('b', {}, String(_p.users)), h('span', {}, plural(_p.users, 'установка', 'установки', 'установок')))),
    _no,
    ..._acts,
    _p.repo ? group('Код', null, row({ lead: sq(I.code, C.indigo), title: 'Репозиторий', sub: _p.repo.replace(/^https:\/\//, ''), chev: true, onclick: () => web(_p.repo) })) : null,
    group('История', null, ..._p.reqs.map((_q) => {
      const [_qt, _qc] = __st__[_q.status]
      return row({ lead: sq(_q.kind === 'edit' ? I.pencil : _q.kind === 'update' ? I.bolt : I.spark, _q.status === 'ok' ? C.green : _q.status === 'no' ? C.red : C.orange),
        title: __kinds__[_q.kind] + (_q.kind !== 'edit' && _q.ver ? ' ' + _q.ver : ''), sub: date(_q.created), right: pill(_qt, _qc), onclick: () => detail(_q) })
    })),
  ))
  scrollTo(0, 0)
}

const avatar = (__cur: string): { el: HTMLElement; get: () => string } => {
  let _data = ''
  const _ic = picIc(__cur, I.code, C.indigo, true)
  const _file = h('input', { type: 'file', accept: 'image/png,image/jpeg', hidden: '' }) as HTMLInputElement
  const _lbl = h('span', { class: 'link' }, __cur ? 'Сменить аватарку' : 'Добавить аватарку')
  const _el = h('button', { class: 'dv-ava-pick', type: 'button' }, _ic, _lbl, _file)
  _el.addEventListener('click', (_e) => {
    if (_e.target === _file) return
    buzz('tick')
    _file.click()
  })
  _file.addEventListener('change', () => {
    const _f = _file.files?.[0]
    if (!_f) return
    const _img = new Image()
    const _u = URL.createObjectURL(_f)
    _img.onload = () => {
      const _c = document.createElement('canvas')
      _c.width = _c.height = 256
      const _x = _c.getContext('2d')!
      const _s = Math.min(_img.width, _img.height)
      _x.drawImage(_img, (_img.width - _s) / 2, (_img.height - _s) / 2, _s, _s, 0, 0, 256, 256)
      _data = _c.toDataURL('image/jpeg', 0.9)
      URL.revokeObjectURL(_u)
      _ic.replaceChildren()
      _ic.classList.add('pic')
      _ic.style.backgroundImage = `url("${_data}")`
      _lbl.textContent = 'Сменить аватарку'
    }
    _img.onerror = () => toast('Не получилось открыть картинку', 'err')
    _img.src = _u
  })
  return { el: _el, get: () => _data }
}

type _Kind = 'new' | 'update' | 'edit'

export const devForm = async (__root: HTMLElement, __kind: _Kind, __id = 0): Promise<void> => {
  back(() => (history.length > 1 ? history.back() : go(__id ? '/dev/p/' + __id : '/dev')))
  let _p: DevOne | null = null
  if (__id) {
    __root.replaceChildren(skDev())
    try {
      _p = await api.devPlug(__id)
    } catch {
      return empty(__root, I.cloud, K.gray, 'Не нашли плагин', 'Попробуй открыть его ещё раз из кабинета.')
    }
  }
  const input = (__ph: string, __max: number, __v = '', __code = false): HTMLInputElement =>
    h('input', { class: 'opt-in', type: 'text', maxlength: String(__max), placeholder: __ph, value: __v, ...(__code ? { autocapitalize: 'off', autocomplete: 'off', spellcheck: 'false' } : {}) }) as HTMLInputElement
  const area = (__ph: string, __max: number, __v = ''): HTMLTextAreaElement => {
    const _a = h('textarea', { class: 'field area', maxlength: String(__max), rows: '3', placeholder: __ph }) as HTMLTextAreaElement
    _a.value = __v
    return _a
  }
  const line = (__t: string, __i: HTMLElement): HTMLElement => h('label', { class: 'line' }, h('span', { class: 'row-main' }, h('span', { class: 'row-t' }, __t)), __i)
  const gh = (__t: string): HTMLElement => h('div', { class: 'gh' }, h('h2', {}, __t))
  const _name = input('Time Parser', 32, _p?.name ?? '')
  const _slug = input('time-parser', 24, '', true)
  const _ver = input(_p?.ver ? 'было ' + _p.ver : '1.0', 16, '', true)
  const _about = area('Что умеет плагин', 300, _p?.about ?? '')
  const _repo = h('input', { class: 'field', type: 'url', inputmode: 'url', maxlength: '200', placeholder: 'github.com/ник/плагин', value: _p?.repo ?? '', autocapitalize: 'off', autocomplete: 'off', spellcheck: 'false' }) as HTMLInputElement
  const _note = area(__kind === 'update' ? 'Что изменилось в этой версии' : 'Например, какие права нужны и зачем', 1000)
  const _ava = avatar(_p?.pic ?? '')
  const _gh = h('p', { class: 'note' }, 'Ссылка на репозиторий GitHub с кодом плагина.')
  const _d = __cache ?? (await api.dev().catch(() => null))
  if (_d?.gh) _gh.textContent = `Если репозиторий закрытый, добавь @${_d.gh} в участники: Settings > Collaborators.`
  _slug.addEventListener('input', () => (_slug.value = _slug.value.toLowerCase().replace(/[^a-z0-9-]/g, '')))
  const __t__: Record<_Kind, [string, string, string]> = {
    new: ['Новый плагин', 'Заполни карточку и пришли ссылку на код. Мы проверим его и добавим в магазин.', 'Отправить на проверку'],
    update: [_p?.live ? 'Обновление' : 'Отправить снова', 'Новая версия кода уходит на проверку. Пока её смотрят, в магазине остаётся текущая.', 'Отправить обновление'],
    edit: ['Карточка плагина', 'Название, описание и аватарка поменяются в магазине после проверки.', 'Отправить на проверку'],
  }
  const [_title, _text, _btnT] = __t__[__kind]
  const _send = h('button', { class: 'v-btn main', type: 'button' }, _btnT)
  _send.addEventListener('click', async () => {
    if (_send.classList.contains('busy')) return
    if (_name.value.trim().length < 2) return toast('Напиши название', 'err')
    if (__kind === 'new' && !/^[a-z][a-z0-9-]{1,23}$/.test(_slug.value)) return toast('Придумай ID плагина латиницей', 'err')
    if (__kind !== 'edit' && !_ver.value.trim()) return toast('Укажи версию', 'err')
    if (__kind !== 'edit' && !/github\.com\/[^/\s]+\/[^/\s]+/i.test(_repo.value)) return toast('Вставь ссылку на репозиторий GitHub', 'err')
    ;(document.activeElement as HTMLElement | null)?.blur()
    const _in: PlugIn = { name: _name.value.trim(), about: _about.value.trim(), ver: _ver.value.trim().replace(/,/g, '.'), note: _note.value.trim(), repo: _repo.value.trim(), pic: _ava.get() }
    _send.classList.add('busy')
    try {
      if (__kind === 'new') {
        const _r = await api.newPlug({ ..._in, slug: _slug.value })
        __cache = null
        history.replaceState(null, '', '/dev/p/' + _r.id)
        void devPlug(__root, _r.id)
      } else {
        await api.changePlug(__id, { ..._in, kind: __kind })
        history.length > 1 ? history.back() : go('/dev/p/' + __id)
      }
      buzz('success')
      toast('Отправлено на проверку')
    } catch (_e) {
      buzz('error')
      toast(err(_e, 'Не получилось отправить'), 'err')
    } finally {
      _send.classList.remove('busy')
    }
  })
  const _kids: (HTMLElement | null)[] = __kind === 'update'
    ? [box(line('Версия', _ver)), gh('Репозиторий'), _repo, _gh, gh('Что нового'), _note]
    : __kind === 'edit'
      ? [_ava.el, box(line('Название', _name)), gh('Описание'), _about, gh('Комментарий для проверки'), _note]
      : [_ava.el, box(line('Название', _name), line('ID', _slug), line('Версия', _ver)),
        h('p', { class: 'note' }, 'ID должен совпадать с полем ID в коде плагина. Поменять его потом нельзя.'),
        gh('Описание'), _about, gh('Репозиторий'), _repo, _gh, gh('Комментарий для проверки'), _note]
  __root.replaceChildren(page('feat dv-form',
    h('div', { class: 'feat-h dv-fh' }, lot(__kind === 'edit' ? 'surprise' : 'stars', 112), h('h1', {}, _title), h('p', {}, _text)),
    ..._kids,
    h('p', { class: 'note' }, 'Код видим только мы. В магазин попадает только готовый плагин. Ответ придёт в чат с ботом.'),
    h('div', { class: 'v-bar' }, _send),
  ))
  scrollTo(0, 0)
}
