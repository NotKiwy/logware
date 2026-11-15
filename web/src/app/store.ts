import * as api from '../api/api'
import type { Me, Mod } from '../api/api'

let __me: Me | null = null

let __mods: Mod[] = []

let __ban = false

export const banned = (): boolean => __ban

export const cached = (): boolean => __me !== null

export const mods = (): Mod[] => __mods

export const load = async (__force = false): Promise<Me | null> => {
  if (__me && !__force) return __me
  const [_m, _l] = await Promise.all([api.me().catch((_e: Error) => ((__ban = _e.message.trim() === 'banned'), null)), api.mods().catch(() => [] as Mod[])])
  __me = _m
  __mods = _l
  return __me
}

export const drop = (): void => {
  __me = null
  __mods = []
}

export const pro = (): boolean => !!__me && __me.plan !== 'free'
