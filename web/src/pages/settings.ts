import * as api from '../api/api'
import type { Me } from '../api/api'
import { load } from '../app/store'
import { go } from '../app/router'
import { buzz, docs, getTheme, goChat, setTheme, web } from '../app/tg'
import type { Theme } from '../app/tg'
import { __tiers__ } from './plans'
import { connect } from './help'
import { C, I } from '../ui/icons'
import { ask, box, group, page, row, sheet, sq } from '../ui/kit'
import { lot } from '../ui/lottie'
import { h } from '../ui/dom'
import { keepLabel } from './feat'
import { stories } from './intro'
import { wipeScreen } from './wipe'

const __modes__: Record<string, { name: string; ic: string; color: string }> = {
  cloud: { name: 'В Telegram', ic: I.cloud, color: C.blue },
  secure: { name: 'В телефоне', ic: I.lock, color: C.green },
  pass: { name: 'Под паролем', ic: I.key, color: C.orange },
}

const __themes__: { v: Theme; label: string }[] = [
  { v: 'auto', label: 'Как в Telegram' },
  { v: 'light', label: 'Светлая' },
  { v: 'dark', label: 'Тёмная' },
  { v: 'deep', label: 'Глубокая тёмная' },
]

const __errs__: Record<string, string> = {
  none: 'Такого промокода нет',
  used: 'Ты уже активировал этот промокод',
  over: 'Этот промокод закончился',
  'too many requests': 'Слишком много попыток, подожди минуту',
}

const promoSheet = (): void => {
  const _in = h('input', { class: 'field promo-in', placeholder: 'КОД', maxlength: '32', autocapitalize: 'characters', autocomplete: 'off', autocorrect: 'off', spellcheck: 'false', enterkeyhint: 'go' }) as HTMLInputElement
  const _err = h('p', { class: 'promo-err' })
  const _go = h('button', { class: 'v-btn main', type: 'button', disabled: true }, 'Активировать') as HTMLButtonElement
  const _body = h('div', { class: 'promo' },
    h('div', { class: 'sh-lot' }, lot('stars', 120), h('h2', {}, 'Промокод')),
    h('p', { class: 'sheet-text' }, 'Введи код, и подписка подключится сразу. Если она уже есть, дни добавятся к ней.'),
    _in, _err, _go)
  const _s = sheet(null, _body)
  let _busy = false
  _in.addEventListener('input', () => {
    _in.value = _in.value.toUpperCase().replace(/\s+/g, '')
    _go.disabled = !_in.value
    _err.textContent = ''
    _in.classList.remove('bad')
  })
  const use = async (): Promise<void> => {
    const _c = _in.value.trim()
    if (!_c || _busy) return
    _busy = true
    _in.blur()
    _go.disabled = true
    _go.replaceChildren(h('span', { class: 'spin' }), 'Проверяю')
    try {
      const _r = await api.promo(_c)
      buzz('success')
      const _d = new Date(_r.until * 1000).toLocaleDateString('ru', { day: 'numeric', month: 'long', year: 'numeric' })
      const _ok = h('button', { class: 'v-btn main', type: 'button' }, 'Отлично')
      _ok.addEventListener('click', () => _s.close())
      _body.replaceChildren(
        h('div', { class: 'sh-lot' }, lot('surprise', 130), h('h2', {}, `${__tiers__[_r.plan]?.name ?? 'Подписка'} подключён`)),
        h('p', { class: 'sheet-text' }, `Действует до ${_d}.`), _ok)
      const _h = _s.done
      await load(true).catch(() => null)
      await _h
      go(location.pathname)
      return
    } catch (_e) {
      buzz('error')
      _err.textContent = __errs__[(_e as Error).message.trim()] ?? 'Не получилось, попробуй ещё раз'
      _in.classList.remove('bad')
      void _in.offsetWidth
      _in.classList.add('bad')
      _go.replaceChildren('Активировать')
      _go.disabled = false
    } finally {
      _busy = false
    }
  }
  _go.addEventListener('click', () => void use())
  _in.addEventListener('keydown', (_e) => _e.key === 'Enter' && void use())
}

const themeLabel = (): string => __themes__.find((_t) => _t.v === getTheme())?.label ?? ''

export const settingsPage = (__me: Me): HTMLElement => {
  const _t = __tiers__[__me.plan] ?? __tiers__.free
  const _sub = __me.plan !== 'free' && __me.until
    ? `${_t.name} до ${new Date(__me.until * 1000).toLocaleDateString('ru', { day: '2-digit', month: '2-digit', year: '2-digit' })}`
    : _t.name
  const _theme: HTMLElement = row({ lead: sq(I.moon, C.indigo), title: 'Тема', value: themeLabel(), chev: true, onclick: async () => {
    let _v: Theme | null = null
    const _rows = __themes__.map((_t) => row({ title: _t.label, check: true, cls: _t.v === getTheme() ? 'on' : '', onclick: () => {
      _v = _t.v
      setTheme(_t.v)
      for (const _r of _rows) _r.classList.toggle('on', _r === _rows[__themes__.indexOf(_t)])
      setTimeout(() => _s.close(), 180)
    } }))
    const _s = sheet(null, h('div', { class: 'sh-lot' }, lot('moon', 110), h('h2', {}, 'Тема')), box(..._rows))
    await _s.done
    if (!_v) return
    const _val = _theme.querySelector('.val')
    if (_val) _val.textContent = themeLabel()
  } })
  const _p = __me.prefs
  const _keep = row({ lead: sq(I.clock, C.indigo), title: 'Сколько хранить', value: keepLabel(_p.keep), chev: true, onclick: () => go('/f/keep') })
  const _m = __modes__[__me.mode ?? ''] ?? __modes__.pass
  const _wipe = row({ title: 'Удалить мои данные', cls: 'danger', onclick: async () => {
    if (!(await ask('Удалить мои данные?', 'Сообщения, ключ, подписка и все настройки будут удалены. На следующем шаге покажем, что именно.', 'Продолжить', true))) return
    wipeScreen()
  } })
  return page('',
    group('Подписка', null,
      row({ lead: sq(I.star, C.purple), title: 'Мой план', value: _sub, chev: true, onclick: () => go('/f/plans') }),
      row({ lead: sq(I.tag, C.pink), title: 'Промокод', chev: true, onclick: promoSheet })),
    group('Хранение', null, _keep),
    group('Оформление', null, _theme),
    group('Защита', null,
      row({ lead: sq(_m.ic, _m.color), title: 'Мой ключ', value: _m.name, chev: true, onclick: () => docs('/start/key') }),
      row({ lead: sq(I.key, C.gray), title: 'Сменить ключ', chev: true, onclick: () => go('/reg') }),
      row({ lead: sq(I.shield, C.green), title: 'Почему нам можно верить', chev: true, onclick: () => stories(() => void 0, true) }),
      row({ lead: sq(I.lock, C.indigo), title: 'Как устроена защита', chev: true, onclick: () => docs('/security') }),
      row({ lead: sq(I.code, C.gray), title: 'Исходный код', chev: true, onclick: () => web('https://github.com/NotKiwy/logware') })),
    group('Помощь', null,
      row({ lead: sq(I.person, C.blue), title: 'Подключить бота', chev: true, onclick: () => connect(__me.bot) }),
      row({ lead: sq(I.plug, C.indigo), title: 'Как подключить бота', chev: true, onclick: () => docs('/start/connect') }),
      row({ lead: sq(I.help, C.purple), title: 'Как это работает', chev: true, onclick: () => docs('/') }),
      row({ lead: sq(I.chat, C.green), title: 'Поддержка', chev: true, onclick: () => goChat(0, __me.support || __me.bot) })),
    group('Правовая информация', null,
      row({ lead: sq(I.shield, C.gray), title: 'Пользовательское соглашение', chev: true, onclick: () => docs('/legal/terms') }),
      row({ lead: sq(I.lock, C.gray), title: 'Политика конфиденциальности', chev: true, onclick: () => docs('/legal/privacy') })),
    box(_wipe),
  )
}
