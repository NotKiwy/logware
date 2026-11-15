import { buzz, tint } from '../app/tg'
import { h } from '../ui/dom'
import { lot } from '../ui/lottie'

type _Slide = { lot: string; title: string; text: string }

const __slides__: _Slide[] = [
  {
    lot: 'hello',
    title: 'Привет!',
    text: 'Logware сохраняет удалённые и изменённые сообщения из твоих чатов. Сейчас коротко расскажем, почему стоит выбрать именно его.',
  },
  {
    lot: 'lock',
    title: 'Сквозное шифрование',
    text: 'Каждое сообщение шифруется сразу при получении. Расшифровать его можно только твоим ключом, поэтому даже мы не можем прочитать твои сообщения.',
  },
  {
    lot: 'deleted',
    title: 'Открытый код',
    text: 'На сервере лежит только шифротекст. Код бота открыт, любой может проверить, как он работает.',
  },
  {
    lot: 'bot',
    title: 'Магазин плагинов',
    text: 'Новые возможности включаются в один клик. Любой может написать свой плагин и выложить его в магазин, а мы проверим каждый перед публикацией.',
  },
]

export const stories = (__done: () => void, __replay = false): void => {
  let _i = 0
  const _bars = __slides__.map(() => h('i', {}))
  const _body = h('div', { class: 'story-body' })
  const _btn = h('button', { class: 'v-btn main story-btn', type: 'button' })
  const _prev = h('button', { class: 'story-tap l', type: 'button', 'aria-label': 'Назад' })
  const _next = h('button', { class: 'story-tap r', type: 'button', 'aria-label': 'Дальше' })
  const _el = h('div', { class: 'story' }, h('div', { class: 'story-bars' }, ..._bars), _body, _prev, _next, h('div', { class: 'story-foot' }, _btn))
  const draw = (__dir: number): void => {
    const _s = __slides__[_i]
    const _box = h('div', { class: 'story-slide' + (__dir < 0 ? ' back' : '') }, lot(_s.lot, 168), h('h1', {}, _s.title), h('p', {}, _s.text))
    _body.replaceChildren(_box)
    _bars.forEach((_b, _n) => _b.classList.toggle('on', _n <= _i))
    _btn.textContent = _i < __slides__.length - 1 ? 'Дальше' : __replay ? 'Понятно' : 'Начнём'
  }
  const step = (__d: number): void => {
    const _n = _i + __d
    if (_n < 0) return
    buzz('tick')
    if (_n >= __slides__.length) {
      _el.classList.add('out')
      setTimeout(() => _el.remove(), 220)
      return __done()
    }
    _i = _n
    draw(__d)
  }
  _btn.addEventListener('click', () => step(1))
  _next.addEventListener('click', () => step(1))
  _prev.addEventListener('click', () => step(-1))
  tint(null)
  document.body.append(_el)
  draw(1)
}
