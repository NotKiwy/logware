type _Kid = Node | string | null | undefined | false

type _Attr = Record<string, unknown>

export const h = <T extends keyof HTMLElementTagNameMap>(
  __tag: T,
  __attr: _Attr = {},
  ...__kids: _Kid[]
): HTMLElementTagNameMap[T] => {
  const _e = document.createElement(__tag)
  for (const [_n, _v] of Object.entries(__attr)) {
    if (_n.startsWith('on') && typeof _v === 'function') _e.addEventListener(_n.slice(2), _v as EventListener)
    else if (_n === 'class') _e.className = String(_v)
    else if (_v === true) _e.setAttribute(_n, '')
    else if (_v !== false && _v != null) _e.setAttribute(_n, String(_v))
  }
  for (const _k of __kids) if (_k) _e.append(_k)
  return _e
}
