import * as api from '../api/api'
import type { Me } from '../api/api'
import { load } from '../app/store'
import { back2, buzz, docs, invoice, tgLink, tint, who } from '../app/tg'
import { h } from '../ui/dom'
import { I } from '../ui/icons'
import { box, group, page, row, sheet, sq, svg, toast } from '../ui/kit'
import { lot } from '../ui/lottie'

export type Plan = Me['plan']

type _Tier = { name: string; tag: string; ic: string; hero: [string, string, string]; price: number; rub: number; usd: string; items: [string, string, boolean][] }

export const __tiers__: Record<Plan, _Tier> = {
  free: {
    name: 'Free',
    tag: 'ДЛЯ ВСЕХ',
    ic: I.shield,
    hero: ['#2a7af5', '#2369ea', '#1d5fdf'],
    price: 0,
    rub: 0,
    usd: '',
    items: [
      [I.trash, 'Удалённые сообщения', false],
      [I.pencil, 'Изменённые сообщения: было и стало', false],
      [I.photo, 'Фото и видео из удалённых, до 20 МБ', false],
      [I.fire, 'Одноразовые фото и видео', false],
      [I.lock, 'Шифрование твоим ключом', false],
    ],
  },
  pro: {
    name: 'Pro',
    tag: 'ВСЕ ПЛАГИНЫ',
    ic: I.bolt,
    hero: ['#8a4dff', '#7a3cf0', '#6a2ee0'],
    price: 149,
    rub: 149,
    usd: '1.79',
    items: [
      [I.shield, 'Всё, что есть в Free', false],
      [I.grid, 'Все плагины из магазина', false],
      [I.clock, 'Время в имени', false],
    ],
  },
  dev: {
    name: 'Developer',
    tag: 'ДЛЯ РАЗРАБОТЧИКОВ',
    ic: I.grid,
    hero: ['#1fae6a', '#18995c', '#12854f'],
    price: 299,
    rub: 299,
    usd: '3.59',
    items: [
      [I.bolt, 'Всё, что есть в Pro', false],
      [I.code, 'Кабинет разработчика', false],
      [I.grid, 'Свои плагины в магазине', false],
      [I.hist, 'Версии, обновления и статистика установок', false],
    ],
  },
}

const __order__: Plan[] = ['free', 'pro', 'dev']

export const card = (__p: Plan, __me: Me | null, __tag: 'button' | 'div' = 'button'): HTMLElement => {
  const _t = __tiers__[__p]
  const _u = who()
  const _paid = !!__me && __me.plan === __p && __p !== 'free' && __me.until > 0
  const _brand = h('span', { class: 'card-brand', 'aria-label': 'Logware' })
  _brand.append(svg(I.word))
  const _logo = h('span', { class: 'card-logo' })
  _logo.append(svg(I.logo))
  const _bot = __me && __me.plan === __p
    ? [h('div', {}, h('small', {}, 'ВЛАДЕЛЕЦ'), h('b', {}, _u.full)),
      h('div', {}, h('small', {}, _paid ? 'ДЕЙСТВУЕТ ДО' : 'СРОК'), h('b', {}, _paid ? new Date(__me.until * 1000).toLocaleDateString('ru', { day: '2-digit', month: '2-digit', year: '2-digit' }) : 'БЕССРОЧНО'))]
    : [h('div', {}, h('small', {}, 'ПЛАН'), h('b', {}, _t.tag))]
  return h(__tag, { class: 'card', type: __tag === 'button' ? 'button' : null, style: `--h1:${_t.hero[0]};--h2:${_t.hero[1]};--h3:${_t.hero[2]}` },
    _logo,
    h('div', { class: 'card-top' }, _brand, h('span', { class: 'card-plan' }, _t.name.toUpperCase())),
    h('div', { class: 'card-bot' }, ..._bot))
}

export const plansPage = (__me: Me): HTMLElement => {
  const _mine = __order__.indexOf(__me.plan in __tiers__ ? __me.plan : 'free')
  const _to = new URLSearchParams(location.search).get('to') as Plan | null
  let _i = _to && __order__.includes(_to) ? __order__.indexOf(_to) : _mine
  let _dir = 0
  const _hold = h('div', { class: 'pv-hold' })
  const _dots = h('div', { class: 'pv-dots' }, ...__order__.map(() => h('span', {})))
  const _block = h('div', { class: 'blue pv-block' }, h('div', { class: 'rays' }), _hold, _dots)
  const _perks = h('div', {})
  let _tick = 0
  const _cta = h('button', { class: 'v-btn main', type: 'button' })
  const _bar = h('div', { class: 'v-bar pv-bar' }, _cta)
  let _busy = false
  const done = (__p: Plan, __m: Me, __text: string): void => {
    _pg.replaceWith(plansPage(__m))
    const _ok = h('button', { class: 'v-btn main', type: 'button' }, 'Отлично')
    const _s = sheet(null, h('div', { class: 'sh-lot' }, lot('surprise', 130), h('h2', {}, `${__tiers__[__p].name} подключён`)),
      h('p', { class: 'sheet-text' }, __text), _ok)
    _ok.addEventListener('click', () => _s.close())
  }
  const wait = async (__p: Plan, __tries: number, __text: string): Promise<boolean> => {
    const _was = __me.until
    for (let _n = 0; _n < __tries; _n++) {
      const _m = await load(true)
      if (_m && _m.plan === __p && _m.until !== _was) {
        done(__p, _m, __text)
        return true
      }
      await new Promise((_r) => setTimeout(_r, 2000))
    }
    return false
  }
  const stars = async (__p: Plan): Promise<void> => {
    const { url } = await api.pay(__p)
    const _s = await invoice(url)
    if (_s === 'paid') {
      buzz('success')
      toast('Оплачено, подключаю ' + __tiers__[__p].name)
      await wait(__p, 10, 'Спасибо за поддержку. Всё уже работает, подписка продлевается сама раз в месяц.')
    } else if (_s === 'failed') {
      buzz('error')
      toast('Оплата не прошла', 'err')
    }
  }
  const coin = async (__p: Plan, __via: 'cb' | 'xr'): Promise<void> => {
    const { url } = await api.pay(__p, __via)
    const _back = back2()
    tgLink(url)
    await _back
    toast('Проверяю оплату')
    if (!(await wait(__p, 30, 'Спасибо за поддержку. Подписка действует 30 дней, за 3 дня до конца бот напомнит продлить.'))) toast('Оплата пока не пришла. Как только придёт, бот напишет', 'err')
  }
  const choose = (__p: Plan): Promise<'stars' | 'cb' | 'xr' | null> => {
    const _t = __tiers__[__p]
    let _res: 'stars' | 'cb' | 'xr' | null = null
    const pickV = (__v: 'stars' | 'cb' | 'xr'): void => {
      _res = __v
      _s.close()
    }
    const _rows: HTMLElement[] = [row({ lead: sq(I.star, '#f5a13c'), title: 'Telegram Stars', sub: `${_t.price} ★ в месяц, продлевается само`, chev: true, onclick: () => pickV('stars') })]
    if (__me.via.includes('cb')) _rows.push(row({ lead: sq(I.bolt, '#2ea6ff'), title: 'CryptoBot', sub: `${_t.rub} ₽ в USDT, TON, BTC и других`, chev: true, onclick: () => pickV('cb') }))
    if (__me.via.includes('xr')) _rows.push(row({ lead: sq(I.bolt, '#7a5cf0'), title: 'xRocket', sub: `${_t.usd} USDT`, chev: true, onclick: () => pickV('xr') }))
    const _s = sheet('Как оплатить', box(..._rows), h('p', { class: 'note' }, 'Криптовалютой платишь за 30 дней. Продлить можно заранее, срок сложится.'))
    return _s.done.then(() => _res)
  }
  _cta.addEventListener('click', async () => {
    if (_busy) return
    const _p = __order__[_i]
    buzz('tap')
    const _v = (__me.via ?? ['stars']).length > 1 && !__me.stars ? await choose(_p) : 'stars'
    if (!_v) return
    _busy = true
    _cta.classList.add('busy')
    try {
      if (_v === 'stars') await stars(_p)
      else await coin(_p, _v)
    } catch (_e) {
      buzz('error')
      const _m = String(_e)
      toast(_m.includes('already paid') ? 'Этот план уже подключён' : _m.includes('has stars') ? 'Подписка оплачена звёздами и продлевается сама' : 'Не получилось открыть оплату', 'err')
    } finally {
      _busy = false
      _cta.classList.remove('busy')
    }
  })
  const paint = (): void => {
    const _p = __order__[_i]
    const _t = __tiers__[_p]
    _block.style.cssText = `--h1:${_t.hero[0]};--h2:${_t.hero[1]};--h3:${_t.hero[2]}`
    tint(_t.hero[0])
    const _c = card(_p, __me, 'div')
    if (_dir) _c.classList.add(_dir > 0 ? 'from-r' : 'from-l')
    _hold.replaceChildren(_c)
    ;[..._dots.children].forEach((_d, _n) => _d.classList.toggle('on', _n === _i))
    _perks.replaceChildren(group('Что входит', null, ..._t.items.map(([_ic, _txt, _soon]) =>
      row({ lead: sq(_ic, _t.hero[0]), title: _txt, right: _soon ? h('span', { class: 'pv-soon' }, 'Скоро') : undefined }))))
    const _renew = _i === _mine && __me.plan !== 'free' && !__me.stars
    _bar.hidden = _i < _mine || (_i === _mine && !_renew)
    _cta.style.backgroundColor = _t.hero[0]
    const _stars = h('span', { class: 'roll-i' }, String(_t.price), h('span', { class: 'tg-star' }))
    _stars.lastElementChild!.append(svg(I.star))
    const _faces = [_stars, h('span', { class: 'roll-i' }, _t.rub + ' ₽'), h('span', { class: 'roll-i' }, '$' + _t.usd)]
    const _roll = h('span', { class: 'roll' }, _faces[0])
    _cta.replaceChildren(_renew ? `Продлить ${_t.name} за` : `Подключить ${_t.name} за`, _roll)
    clearInterval(_tick)
    let _f = 0
    _tick = window.setInterval(() => {
      if (!_cta.isConnected) return clearInterval(_tick)
      const _old = _faces[_f]
      _f = (_f + 1) % _faces.length
      const _new = _faces[_f]
      _old.className = 'roll-i out'
      _new.className = 'roll-i in'
      _roll.append(_new)
      _roll.style.width = _new.offsetWidth + 'px'
      setTimeout(() => _old !== _faces[_f] && _old.replaceWith(), 720)
    }, 2000)
    requestAnimationFrame(() => (_roll.style.width = _faces[0].offsetWidth + 'px'))
  }
  const go = (__d: number): void => {
    const _n = Math.max(0, Math.min(__order__.length - 1, _i + __d))
    if (_n === _i) return buzz('error')
    buzz('tick')
    _dir = __d
    _i = _n
    paint()
  }
  let _x0 = 0
  let _y0 = 0
  let _t0 = 0
  let _axis: 'x' | 'y' | null = null
  let _down = false
  const move = (__x: number, __y: number): void => {
    if (!_down) return
    const _mx = __x - _x0
    const _my = __y - _y0
    if (!_axis && Math.hypot(_mx, _my) > 8) _axis = Math.abs(_mx) > Math.abs(_my) ? 'x' : 'y'
    if (_axis !== 'x') return
    const _edge = (_mx > 0 && _i === 0) || (_mx < 0 && _i === __order__.length - 1)
    const _dx = _edge ? _mx * 0.25 : _mx
    _hold.style.transition = 'none'
    _hold.style.transform = `translateX(${_dx}px) rotate(${_dx / 40}deg)`
    _hold.style.opacity = String(1 - Math.min(0.5, Math.abs(_dx) / 400))
  }
  const end = (__x: number): void => {
    if (!_down) return
    _down = false
    _hold.style.transition = ''
    _hold.style.transform = ''
    _hold.style.opacity = ''
    if (_axis !== 'x') return
    const _mx = __x - _x0
    const _v = Math.abs(_mx) / Math.max(1, Date.now() - _t0)
    if (Math.abs(_mx) > 50 || (Math.abs(_mx) > 18 && _v > 0.35)) go(_mx < 0 ? 1 : -1)
  }
  _block.addEventListener('touchstart', (_e) => {
    _down = true
    _axis = null
    _x0 = _e.touches[0].clientX
    _y0 = _e.touches[0].clientY
    _t0 = Date.now()
  }, { passive: true })
  _block.addEventListener('touchmove', (_e) => move(_e.touches[0].clientX, _e.touches[0].clientY), { passive: true })
  _block.addEventListener('touchend', (_e) => end(_e.changedTouches[0].clientX))
  _block.addEventListener('touchcancel', () => end(_x0))
  _dots.addEventListener('click', (_e) => {
    const _n = [..._dots.children].indexOf(_e.target as Element)
    if (_n >= 0 && _n !== _i) go(_n - _i)
  })
  paint()
  const _pg = page('plan-pg',
    _block,
    _perks,
    h('p', { class: 'note pv-note' }, 'Оплата звёздами Telegram или криптовалютой. Звёзды продлеваются сами, криптой платишь за каждые 30 дней. Основные функции бесплатны навсегда. Оплачивая, ты принимаешь ', (() => { const _a = h('a', { href: 'https://docs.logware.dev/legal/terms', class: 'inl' }, 'пользовательское соглашение'); _a.addEventListener('click', (_e) => { _e.preventDefault(); docs('/legal/terms') }); return _a })(), '.'),
    _bar,
  )
  return _pg
}
