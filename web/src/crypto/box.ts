import _s from 'libsodium-wrappers-sumo'

export const ready = (): Promise<void> => _s.ready

export type Pair = { pub: Uint8Array; sec: Uint8Array }

export const pair = (): Pair => {
  const _k = _s.crypto_box_keypair()
  return { pub: _k.publicKey, sec: _k.privateKey }
}

export const open = (__ct: Uint8Array, __sec: Uint8Array): Uint8Array =>
  _s.crypto_box_seal_open(__ct, _s.crypto_scalarmult_base(__sec), __sec)

export const b64 = (__b: Uint8Array): string => _s.to_base64(__b, _s.base64_variants.ORIGINAL)

export const unb64 = (__t: string): Uint8Array => _s.from_base64(__t, _s.base64_variants.ORIGINAL)

export const str = (__b: Uint8Array): string => _s.to_string(__b)

export const pubOf = (__sec: Uint8Array): Uint8Array => _s.crypto_scalarmult_base(__sec)

export const seal = (__data: Uint8Array, __pub: Uint8Array): Uint8Array => _s.crypto_box_seal(__data, __pub)

export const sha = (__t: string): Uint8Array => _s.crypto_hash_sha256(_s.from_string(__t))

export const bytes = (__t: string): Uint8Array => _s.from_string(__t)
