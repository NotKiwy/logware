import * as api from '../api/api'
import type { Me, Mode } from '../api/api'
import { go } from '../app/router'
import { drop, load } from '../app/store'
import { canCloud, canSecure } from '../app/tg'
import { b64, pair, ready } from '../crypto/box'
import { forget, put } from '../keys/keys'
import { h } from '../ui/dom'
import { lot } from '../ui/lottie'
import { skFeat } from '../ui/skel'
import { C, I } from '../ui/icons'
import { K, cta, empty, group, note, orb, page, row, sheet, svg, toast, wait } from '../ui/kit'
import { stories } from './intro'

const make = async (__root: HTMLElement, __mode: Mode, __pass = ''): Promise<void> => {
  wait(__root, 'Создаю твой ключ…')
  try {
    await ready()
    await forget()
    const _k = pair()
    const _vault = await put(__mode, _k.sec, __pass)
    await api.reg(__mode, b64(_k.pub), _vault)
    drop()
    go('/')
    toast('Ключ создан')
  } catch {
    toast('Не получилось создать ключ', 'err')
    void reg(__root)
  }
}

const pass = (__root: HTMLElement): void => {
  const _p1 = h('input', { class: 'field', type: 'password', placeholder: 'Пароль, минимум 8 символов', autocomplete: 'new-password' })
  const _p2 = h('input', { class: 'field', type: 'password', placeholder: 'Повтори пароль', autocomplete: 'new-password' })
  const _err = note('', true)
  const _btn = cta('Создать ключ', () => {
    if (_p1.value.length < 8) return void (_err.textContent = 'Пароль слишком короткий')
    if (_p1.value !== _p2.value) return void (_err.textContent = 'Пароли не совпадают')
    _s.close()
    void make(__root, 'pass', _p1.value)
  })
  const _s = sheet('Придумай пароль', h('p', { class: 'sheet-text' }, 'Запомни его или запиши. Без пароля расшифровать сохранённое не получится, восстановить его нельзя.'), _p1, _p2, _err, _btn)
  setTimeout(() => _p1.focus(), 350)
}

const choose = (__root: HTMLElement, __me: Me): void => {
  const _cloud = canCloud()
  const _secure = canSecure()
  const _info = h('div', { class: 'outline' }, h('span', {}, __me.reg
    ? 'У тебя уже есть ключ. Новый ключ не расшифрует старые сообщения, поэтому всё сохранённое удалится.'
    : 'Ключ создаётся на твоём устройстве. Закрытая часть никогда не попадает на сервер.'))
  const _page = page('bare',
    h('div', { class: 'title' },
      lot('lock', 120),
      h('h1', {}, 'Где хранить твой ключ?'),
      h('p', {}, 'Закрытый ключ расшифровывает сохранённые сообщения. Без него их не прочитать.')),
    group('Проще всего', null,
      row({ lead: orb(I.cloud, C.blue), title: 'В Telegram', sub: _cloud ? 'Облачное хранилище Telegram, на всех устройствах' : 'Нужна новая версия Telegram', chev: true, cls: _cloud ? '' : 'dim', onclick: _cloud ? () => void make(__root, 'cloud') : undefined })),
    group('Надёжнее', null,
      row({ lead: orb(I.lock, C.green), title: 'В этом телефоне', sub: _secure ? 'Защищённое хранилище, только на этом устройстве' : 'Работает только в Telegram на телефоне', chev: true, cls: _secure ? '' : 'dim', onclick: _secure ? () => void make(__root, 'secure') : undefined }),
      row({ lead: orb(I.key, C.orange), title: 'Под паролем', sub: 'Ключ шифруется паролем, его знаешь только ты', chev: true, onclick: () => pass(__root) })),
    _info,
  )
  if (__me.reg) {
    const _back = h('button', { class: 'outline', type: 'button' }, h('span', {}, 'Оставить старый ключ'))
    _back.append(svg(I.chev))
    _back.addEventListener('click', () => go('/'))
    _page.append(_back)
  }
  __root.replaceChildren(_page)
  scrollTo(0, 0)
}

export const reg = async (__root: HTMLElement): Promise<void> => {
  __root.replaceChildren(skFeat())
  const _me = await load(true)
  if (!_me) return empty(__root, I.cloud, K.gray, 'Нет связи', 'Не получилось связаться с сервером. Закрой и открой ещё раз.')
  if (_me.reg) return choose(__root, _me)
  __root.replaceChildren()
  stories(() => choose(__root, _me))
}
