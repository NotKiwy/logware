import { buzz, docs } from '../app/tg'
import { h } from '../ui/dom'
import { cta, list, sheet, toast } from '../ui/kit'

const copy = async (__t: string): Promise<void> => {
  try {
    await navigator.clipboard.writeText(__t)
  } catch {
    const _a = h('textarea', { style: 'position:fixed;opacity:0' }) as HTMLTextAreaElement
    _a.value = __t
    document.body.append(_a)
    _a.select()
    document.execCommand('copy')
    _a.remove()
  }
  buzz('success')
  toast('Ник бота скопирован')
}

export const connect = (__bot: string, __rights = false): void => {
  const _nick = '@' + __bot
  const _more = h('button', { class: 'link sheet-link', type: 'button' }, 'Подробнее в документации')
  _more.addEventListener('click', () => docs(__rights ? '/start/connect#rights' : '/start/connect'))
  const _s = sheet(__rights ? 'Выдай боту права' : 'Подключи бота',
    h('p', { class: 'sheet-text' }, __rights ? 'Бот уже подключён, но ему не хватает прав.' : 'Это делается один раз в настройках Telegram.'),
    list(__rights ? [
      ['1', 'Открой свой профиль в Telegram'],
      ['2', 'Нажми «Изменить»'],
      ['3', 'Зайди в «ИИ-автоматизация»'],
      ['4', `Выбери ${_nick}`],
      ['5', 'Включи «Отвечать на сообщения» и «Удалять сообщения»'],
    ] : [
      ['1', 'Открой свой профиль в Telegram'],
      ['2', 'Нажми «Изменить»'],
      ['3', 'Зайди в «ИИ-автоматизация»'],
      ['4', `Вставь ${_nick} и добавь бота`],
    ], true),
    cta(`Скопировать ${_nick}`, () => void copy(_nick)),
    _more)
  void _s
}
