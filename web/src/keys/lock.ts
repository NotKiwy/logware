import type { Me } from '../api/api'
import { buzz } from '../app/tg'
import { h } from '../ui/dom'
import { cta, note, page } from '../ui/kit'
import { byPass, load } from './keys'

const password = (__root: HTMLElement, __vault: string): Promise<Uint8Array> =>
  new Promise((_ok) => {
    const _in = h('input', { class: 'field', type: 'password', placeholder: 'Пароль', autocomplete: 'current-password' })
    const _err = note('', true)
    const _btn = cta('Открыть', () => {
      _btn.busy(true, 'Проверяю')
      setTimeout(() => {
        try {
          _ok(byPass(__vault, _in.value))
        } catch {
          buzz('error')
          _err.textContent = 'Неверный пароль'
          _btn.busy(false)
        }
      }, 30)
    })
    _in.addEventListener('keydown', (_e) => _e.key === 'Enter' && _btn.click())
    __root.replaceChildren(page('bare short',
      h('div', { class: 'title', style: 'padding-bottom:20px' }, h('h1', {}, 'Введи пароль'), h('p', {}, 'Он расшифровывает сохранённые сообщения.')),
      _in, _err, h('div', { style: 'height:8px' }), _btn))
    _in.focus()
  })

export const key = async (__root: HTMLElement, __me: Me): Promise<Uint8Array | null> => {
  if (__me.mode === 'pass' && __me.vault) return password(__root, __me.vault)
  return load(__me.mode!).catch(() => null)
}
