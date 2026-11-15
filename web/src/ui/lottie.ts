import { h } from './dom'

type _Anim = { play: () => void; goToAndPlay: (_f: number, _isFrame: boolean) => void; destroy: () => void }

const __cache__ = new Map<string, Promise<unknown>>()

const load = (__name: string): Promise<unknown> => {
  let _p = __cache__.get(__name)
  if (!_p) {
    _p = fetch(`/lottie/${__name}.tgs`).then(async (_r) => {
      if (!_r.ok) throw new Error('lottie ' + _r.status)
      return unzip(await _r.blob())
    })
    __cache__.set(__name, _p)
  }
  return _p
}

const play = (__data: Promise<unknown>, __size: number, __loop: boolean): HTMLElement => {
  const _el = h('div', { class: 'lot', style: `width:${__size}px;height:${__size}px`, 'aria-hidden': 'true' })
  let _a: _Anim | null = null
  Promise.all([import('lottie-web/build/player/lottie_light'), __data])
    .then(([_m, _data]) => {
      _a = _m.default.loadAnimation({ container: _el, renderer: 'svg', loop: __loop, autoplay: true, animationData: _data }) as unknown as _Anim
    })
    .catch(() => void 0)
  if (!__loop) _el.addEventListener('click', () => _a?.goToAndPlay(0, true))
  return _el
}

const unzip = async (__b: Blob): Promise<unknown> =>
  JSON.parse(await new Response(__b.stream().pipeThrough(new DecompressionStream('gzip'))).text())

export const lot = (__name: string, __size = 120): HTMLElement => play(load(__name), __size, false)

export const lotBlob = (__b: Blob, __size = 200): HTMLElement => play(unzip(__b), __size, true)
