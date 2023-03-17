const opensslBanner = new Uint8Array(new TextEncoder('utf8').encode('Salted__'))
const pbkdf2Params = { hash: 'SHA-512', iterations: 300000, name: 'PBKDF2' }

function decrypt(passphrase, encData) {
  const data = new Uint8Array(atob(encData).split('')
    .map(c => c.charCodeAt(0)))

  return deriveKey(passphrase, data.slice(8, 16))
    .then(({ iv, key }) => window.crypto.subtle.decrypt({ iv, name: 'AES-CBC' }, key, data.slice(16)))
    .then(data => new TextDecoder('utf8').decode(data))
}

function deriveKey(passphrase, salt) {
  return window.crypto.subtle.importKey('raw', new TextEncoder('utf8').encode(passphrase), 'PBKDF2', false, ['deriveBits'])
    .then(passwordKey => window.crypto.subtle.deriveBits({ ...pbkdf2Params, salt }, passwordKey, 384))
    .then(key => window.crypto.subtle.importKey('raw', key.slice(0, 32), { name: 'AES-CBC' }, false, ['encrypt', 'decrypt'])
      .then(aesKey => ({ iv: key.slice(32, 48), key: aesKey })))
}

function encrypt(passphrase, salt, plainData) {
  return deriveKey(passphrase, salt)
    .then(({ iv, key }) => window.crypto.subtle.encrypt({ iv, name: 'AES-CBC' }, key, new TextEncoder('utf8').encode(plainData)))
    .then(encData => new Uint8Array([...opensslBanner, ...salt, ...new Uint8Array(encData)]))
    .then(data => btoa(String.fromCharCode.apply(null, data)))
}

function generateSalt() {
  const salt = new Uint8Array(8) // Salt MUST consist of 8 byte
  return window.crypto.getRandomValues(salt)
}

export default {
  dec: (cipherText, passphrase) => decrypt(passphrase, cipherText),
  enc: (plainText, passphrase) => encrypt(passphrase, generateSalt(), plainText),
}
