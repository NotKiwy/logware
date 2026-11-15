import { cloud, secure } from '../app/tg'
import type { Mode } from '../api/api'
import { b64, unb64 } from '../crypto/box'
import { lock, unlock } from '../crypto/vault'

const __key__ = 'lw_sk'

let __sec: Uint8Array | null = null

export const cached = (): Uint8Array | null => __sec

export const put = async (__mode: Mode, __s: Uint8Array, __pass?: string): Promise<string | null> => {
  __sec = __s
  if (__mode === 'cloud') await cloud().set(__key__, b64(__s))
  if (__mode === 'secure') await secure().set(__key__, b64(__s))
  if (__mode === 'pass') return b64(lock(__s, __pass ?? ''))
  return null
}

export const load = async (__mode: Mode): Promise<Uint8Array | null> => {
  if (__sec) return __sec
  const _v = __mode === 'cloud' ? await cloud().get(__key__) : __mode === 'secure' ? await secure().get(__key__) : null
  if (_v) __sec = unb64(_v)
  return __sec
}

export const byPass = (__vault: string, __pass: string): Uint8Array => {
  __sec = unlock(unb64(__vault), __pass)
  return __sec
}

export const forget = async (): Promise<void> => {
  __sec = null
  await cloud().del(__key__).catch(() => {})
  await secure().del(__key__).catch(() => {})
}
