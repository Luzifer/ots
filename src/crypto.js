const opensslBanner = new Uint8Array(new TextEncoder('utf8').encode('Salted__'))
const pbkdf2Params = { hash: 'SHA-512', iterations: 300000, name: 'PBKDF2' }

/**
 * @param {ArrayBuffer} data Data to encode to base64
 * @returns String
 */
function abToB64(data) {
  const outdata = []
  const bytes = new Uint8Array(data)
  for (let i = 0; i < bytes.byteLength; i++) {
    outdata.push(String.fromCodePoint(bytes[i]))
  }
  return btoa(outdata.join(''))
}

/**
 *
 * @param {String} encoded Base64 encoded data
 * @returns ArrayBuffer
 */
function b64ToAb(encoded) {
  const binary = atob(encoded)
  return Uint8Array.from(binary, c => c.codePointAt(0)).buffer
}

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
  const data = new Uint8Array(b64ToAb(encData))

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
    .then(data => abToB64(data.buffer))
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

export default { abToB64, b64ToAb, dec, enc }
