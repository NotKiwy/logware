import { __docs__ } from '../app/tg'
import { h } from '../ui/dom'
import { I } from '../ui/icons'
import { svg } from '../ui/kit'
import { lot } from '../ui/lottie'

export const outside = (__root: HTMLElement): void => {
  const _word = h('span', { class: 'out-word', 'aria-label': 'Logware' })
  _word.append(svg(I.word))
  const _open = h('a', { class: 'v-btn main out-btn', href: 'https://t.me/logwareapp' }, 'Открыть в Telegram') as HTMLAnchorElement
  const _docs = h('a', { class: 'out-link', href: __docs__ }, 'Как это работает')
  __root.replaceChildren(h('div', { class: 'out' },
    h('div', { class: 'out-glow' }),
    h('div', { class: 'out-in' },
      lot('hello', 150),
      _word,
      h('p', {}, 'Сохраняет удалённые и изменённые сообщения из твоих чатов в Telegram. Под замком, ключ от которого есть только у тебя.'),
      _open,
      _docs),
    h('small', { class: 'out-foot' }, 'Приложение работает внутри Telegram')))
  void fetch('/api/bot').then((_r) => _r.json()).then((_j: { bot?: string }) => {
    if (_j.bot) _open.href = 'https://t.me/' + _j.bot
  }).catch(() => void 0)
}
