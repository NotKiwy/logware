import * as api from '../api/api'
import { go } from '../app/router'
import { drop } from '../app/store'
import { back, buzz } from '../app/tg'
import { forget } from '../keys/keys'
import { h } from '../ui/dom'
import { C, I } from '../ui/icons'
import { box, row, sq, svg, toast } from '../ui/kit'
import { lot } from '../ui/lottie'

const slider = (__go: () => Promise<boolean>): HTMLElement => {
  const _knob = h('span', { class: 'wp-knob' })
  _knob.append(svg(I.trash))
  const _txt = h('span', { class: 'wp-txt' }, 'Проведи, чтобы удалить')
  const _fill = h('span', { class: 'wp-fill' })
  const _track = h('div', { class: 'wp-track', role: 'slider', 'aria-label': 'Проведи, чтобы удалить' }, _fill, _txt, _knob)
  let _x0 = 0
  let _dx = 0
  let _on = false
  let _done = false
  let _max = 0
  const set = (__v: number): void => {
    _dx = Math.max(0, Math.min(_max, __v))
    _knob.style.transform = `translateX(${_dx}px)`
    _fill.style.width = `${_dx + _knob.offsetWidth}px`
    _txt.style.opacity = String(1 - Math.min(1, _dx / (_max * 0.6)))
  }
  const reset = (): void => {
    _track.classList.add('back')
    set(0)
    setTimeout(() => _track.classList.remove('back'), 300)
  }
  _knob.addEventListener('pointerdown', (_e) => {
    if (_done) return
    _on = true
    _max = _track.clientWidth - _knob.offsetWidth - 8
    _x0 = _e.clientX - _dx
    _knob.setPointerCapture(_e.pointerId)
    buzz('tick')
  })
  _knob.addEventListener('pointermove', (_e) => {
    if (!_on) return
    set(_e.clientX - _x0)
  })
  const end = async (): Promise<void> => {
    if (!_on) return
    _on = false
    if (_dx < _max * 0.92) {
      buzz('error')
      return reset()
    }
    _done = true
    set(_max)
    buzz('warning')
    _track.classList.add('busy')
    _knob.replaceChildren(h('span', { class: 'spin' }))
    _txt.textContent = 'Удаляем'
    _txt.style.opacity = '1'
    if (!(await __go())) {
      _done = false
      _knob.replaceChildren()
      _knob.append(svg(I.trash))
      _txt.textContent = 'Проведи, чтобы удалить'
      _track.classList.remove('busy')
      reset()
    }
  }
  _knob.addEventListener('pointerup', () => void end())
  _knob.addEventListener('pointercancel', () => void end())
  return _track
}

export const wipeScreen = (): void => {
  const close = (): void => {
    _w.classList.add('out')
    back(null)
    setTimeout(() => _w.remove(), 260)
  }
  const _w = h('div', { class: 'wp', role: 'dialog' },
    h('div', { class: 'wp-body' },
      h('div', { class: 'feat-h' }, lot('delete', 132), h('h1', {}, 'Удалить всё?'), h('p', {}, 'Это нельзя отменить. Восстановить данные не сможем ни мы, ни ты.')),
      h('div', { class: 'gh' }, h('h2', {}, 'Что удалится')),
      box(
        row({ lead: sq(I.lock, C.indigo), title: 'Сообщения, фото и видео' }),
        row({ lead: sq(I.hist, C.blue), title: 'История и уведомления' }),
        row({ lead: sq(I.key, C.gray), title: 'Ключ шифрования' }),
        row({ lead: sq(I.star, C.purple), title: 'Подписка и автопродление' }),
        row({ lead: sq(I.grid, C.orange), title: 'Плагины и настройки' }),
        row({ lead: sq(I.code, C.green), title: 'Профиль разработчика' })),
      h('p', { class: 'note' }, 'Звёзды за подписку не вернутся. Бот останется подключён, но перестанет что-либо сохранять.')),
    h('div', { class: 'wp-bar' }, slider(async () => {
      try {
        await api.wipe()
        await forget()
        drop()
        buzz('success')
        toast('Данные удалены')
        close()
        go('/reg')
        return true
      } catch {
        buzz('error')
        toast('Не получилось удалить', 'err')
        return false
      }
    })))
  back(close)
  document.body.append(_w)
}
