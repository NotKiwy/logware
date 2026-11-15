import { devForm, devPlug } from '../pages/dev'
import { feat } from '../pages/feat'
import { outside } from '../pages/out'
import { plugin } from '../pages/plugin'
import { reg } from '../pages/reg'
import { shell } from '../pages/shell'
import { view } from '../pages/view'
import { back, inside } from './tg'

export const go = (__to: string = location.pathname + location.search): void => {
  if (__to !== location.pathname + location.search) history.pushState(null, '', __to)
  const __path = __to.split('?')[0]
  const _root = document.getElementById('app')!
  if (!inside()) return outside(_root)
  back(null)
  const _f = __path.match(/^\/f\/([a-z]+)$/)
  if (_f) return void feat(_root, _f[1])
  const _p = __path.match(/^\/p\/([a-z0-9-]+)$/)
  if (_p) return void plugin(_root, _p[1])
  const _ev = __path.match(/^\/e\/([a-f0-9]+)$/)
  if (_ev) return void view(_root, _ev[1])
  if (__path === '/reg') return void reg(_root)
  if (__path === '/dev/new') return void devForm(_root, 'new')
  const _dp = __path.match(/^\/dev\/p\/(\d+)(?:\/(update|edit))?$/)
  if (_dp) return void (_dp[2] ? devForm(_root, _dp[2] as 'update' | 'edit', Number(_dp[1])) : devPlug(_root, Number(_dp[1])))
  void shell(_root, __path)
}

addEventListener('popstate', () => go())

for (const _e of ['gesturestart', 'gesturechange', 'gestureend']) document.addEventListener(_e, (_x) => _x.preventDefault(), { passive: false })

document.addEventListener('touchmove', (_x) => {
  if ((_x as TouchEvent).touches.length > 1) _x.preventDefault()
}, { passive: false })
