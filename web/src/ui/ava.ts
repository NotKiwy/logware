import { pic } from '../api/api'
import { h } from './dom'

const __tints__ = ['#ff6b5e', '#ffa63d', '#a77bf3', '#4b8ff7', '#34c3a0', '#f25c9b', '#5ac8fa']

const __urls__ = new Map<string, Promise<string | null>>()

const url = (__id: number, __base?: string): Promise<string | null> => {
  const _k = (__base ?? '') + __id
  let _p = __urls__.get(_k)
  if (!_p) {
    _p = pic(__id, __base).then((_b) => URL.createObjectURL(_b)).catch(() => null)
    __urls__.set(_k, _p)
  }
  return _p
}

export const ava = (__id: number, __name: string, __cls = 'm-ava', __base?: string): HTMLElement => {
  const _c = __tints__[[...__name].reduce((_a, _x) => _a + _x.charCodeAt(0), 0) % __tints__.length]
  const _el = h('span', { class: __cls, style: __cls.startsWith('m-ava') ? `background-color:${_c}` : null }, [...__name][0]?.toUpperCase() ?? '?')
  if (__id > 0) {
    void url(__id, __base).then((_u) => {
      if (!_u) return
      _el.textContent = ''
      _el.style.backgroundImage = `url(${_u})`
      _el.classList.add('pic')
    })
  }
  return _el
}
