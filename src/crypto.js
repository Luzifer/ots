import base64 from 'base64-js'

const opensslBanner = new Uint8Array(new TextEncoder('utf8').encode('Salted__'))
const pbkdf2Params = { hash: 'SHA-512', iterations: 300000, name: 'PBKDF2' }

/**
 * @param {String} cipherText Encrypted data in base64 encoded form
 * @param {String} passphrase Encryption passphrase used for key-derivation
 * @returns String
 */
function dec(cipherText, passphrase) {
  return decrypt(passphrase, cipherText)
}
/**
 *
 * @param {String} plainText Data to encrypt
 * @param {String} passphrase Encryption passphrase used for key-derivation
 * @returns String
 */
function enc(plainText, passphrase) {
  return encrypt(passphrase, generateSalt(), plainText)
}

/**
 * @param {String} passphrase Encryption passphrase used for key-derivation
 * @param {String} encData Encrypted data in base64 encoded form
 * @returns String
 */
function decrypt(passphrase, encData) {
  const data = base64.toByteArray(encData)

  return deriveKey(passphrase, data.slice(8, 16))
    .then(({ iv, key }) => window.crypto.subtle.decrypt({ iv, name: 'AES-CBC' }, key, data.slice(16)))
    .then(data => new TextDecoder('utf8').decode(data))
}

/**
 *
 * @param {String} passphrase
 * @param {Uint8Array} salt
 * @returns Object
 */
function deriveKey(passphrase, salt) {
  return window.crypto.subtle.importKey('raw', new TextEncoder('utf8').encode(passphrase), 'PBKDF2', false, ['deriveBits'])
    .then(passwordKey => window.crypto.subtle.deriveBits({ ...pbkdf2Params, salt }, passwordKey, 384))
    .then(key => window.crypto.subtle.importKey('raw', key.slice(0, 32), { name: 'AES-CBC' }, false, ['encrypt', 'decrypt'])
      .then(aesKey => ({ iv: key.slice(32, 48), key: aesKey })))
}

/**
 * @param {String} passphrase Encryption passphrase used for key-derivation
 * @param {Uint8Array} salt Cryptographically random salt of 8 byte length
 * @param {String} plainData Data to encrypt
 * @returns String
 */
function encrypt(passphrase, salt, plainData) {
  return deriveKey(passphrase, salt)
    .then(({ iv, key }) => window.crypto.subtle.encrypt({ iv, name: 'AES-CBC' }, key, new TextEncoder('utf8').encode(plainData)))
    .then(encData => new Uint8Array([...opensslBanner, ...salt, ...new Uint8Array(encData)]))
    .then(data => base64.fromByteArray(data))
}

/**
 * Generates a cryptographically secure random salt
 *
 * @returns Uint8Array
 */
function generateSalt() {
  const salt = new Uint8Array(8) // Salt MUST consist of 8 byte
  return window.crypto.getRandomValues(salt)
}

export default { dec, enc }
