import base64 from 'base64-js'

/**
 * OTSMeta defines the structure of (de-)serializing stored payload for secrets
 */
class OTSMeta {
  /** @type File[] */
  #files = []

  /** @type String */
  #secret = ''

  /** @type Number */
  #version = 1.0

  /**
   * @param {String | null} jsonString JSON string representation of OTSMeta created by serialize
   */
  constructor(jsonString = null) {
    if (jsonString === null) {
      return
    }

    if (!jsonString.startsWith('OTSMeta')) {
      // Looks like we got a plain string, we assume that to be a secret only
      this.#secret = jsonString
      return
    }

    const data = JSON.parse(jsonString.replace(/^OTSMeta/, ''))

    this.#secret = data.secret
    this.#version = data.v

    for (const f of data.attachments || []) {
      const content = base64.toByteArray(f.data)
      this.#files.push(new File([content], f.name, { type: f.type }))
    }
  }

  get files() {
    return this.#files
  }

  get secret() {
    return this.#secret
  }

  set secret(secret) {
    this.#secret = secret
  }

  /**
   * @returns {Promise<string>}
   */
  serialize() {
    const output = {
      secret: this.#secret,
      v: this.#version,
    }

    if (this.#files.length === 0) {
      /*
       * We got no attachments, therefore we do a simple fallback to
       * the old "just the secret"-format
       */
      return new Promise(resolve => {
        resolve(this.#secret)
      })
    }

    const encodes = []
    output.attachments = []

    for (const f of this.#files) {
      encodes.push(f.arrayBuffer()
        .then(ab => {
          const data = base64.fromByteArray(new Uint8Array(ab))
          output.attachments.push({ data, name: f.name, type: f.type })
        }))
    }

    return Promise.all(encodes)
      .then(() => `OTSMeta${JSON.stringify(output)}`)
  }
}

export default OTSMeta
