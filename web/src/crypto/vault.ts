import _s from 'libsodium-wrappers-sumo'

const derive = (__pass: string, __salt: Uint8Array): Uint8Array =>
  _s.crypto_pwhash(
    _s.crypto_secretbox_KEYBYTES,
    __pass,
    __salt,
    _s.crypto_pwhash_OPSLIMIT_INTERACTIVE,
    _s.crypto_pwhash_MEMLIMIT_INTERACTIVE,
    _s.crypto_pwhash_ALG_ARGON2ID13,
  )

export const lock = (__sec: Uint8Array, __pass: string): Uint8Array => {
  const _salt = _s.randombytes_buf(_s.crypto_pwhash_SALTBYTES)
  const _nonce = _s.randombytes_buf(_s.crypto_secretbox_NONCEBYTES)
  const _ct = _s.crypto_secretbox_easy(__sec, _nonce, derive(__pass, _salt))
  const _out = new Uint8Array(_salt.length + _nonce.length + _ct.length)
  _out.set(_salt, 0)
  _out.set(_nonce, _salt.length)
  _out.set(_ct, _salt.length + _nonce.length)
  return _out
}

export const unlock = (__vault: Uint8Array, __pass: string): Uint8Array => {
  const _a = _s.crypto_pwhash_SALTBYTES
  const _b = _a + _s.crypto_secretbox_NONCEBYTES
  return _s.crypto_secretbox_open_easy(__vault.slice(_b), __vault.slice(_a, _b), derive(__pass, __vault.slice(0, _a)))
}
