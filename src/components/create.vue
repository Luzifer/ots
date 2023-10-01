<!-- eslint-disable vue/no-v-html -->
<template>
  <!-- Creation disabled -->
  <div
    v-if="!canWrite"
    class="card border-info-subtle mb-3"
  >
    <div
      class="card-header bg-info-subtle"
      v-html="$t('title-secret-create-disabled')"
    />
    <div
      class="card-body"
      v-html="$t('text-secret-create-disabled')"
    />
  </div>

  <!-- Creation possible -->
  <div
    v-else
    class="card border-primary-subtle mb-3"
  >
    <div
      class="card-header bg-primary-subtle"
      v-html="$t('title-new-secret')"
    />
    <div class="card-body">
      <form
        class="row"
        @submit.prevent="createSecret"
      >
        <div class="col-12 mb-3">
          <label for="createSecretData">{{ $t('label-secret-data') }}</label>
          <textarea
            id="createSecretData"
            v-model="secret"
            class="form-control"
            rows="5"
          />
        </div>
        <div class="col-12 mb-3">
          <label for="createSecretFiles">{{ $t('label-secret-files') }}</label>
          <input
            id="createSecretFiles"
            ref="createSecretFiles"
            class="form-control"
            type="file"
            multiple
          >
        </div>
        <div class="col-md-6 col-12 order-2 order-md-1">
          <button
            type="submit"
            class="btn btn-success"
            :disabled="secret.trim().length < 1"
          >
            {{ $t('btn-create-secret') }}
          </button>
        </div>
        <div
          v-if="!$root.customize.disableExpiryOverride"
          class="col-md-6 col-12 order-1 order-md-2"
        >
          <div class="row mb-3 justify-content-end">
            <label
              class="col-md-6 col-form-label text-md-end"
              for="createSecretExpiry"
            >{{ $t('label-expiry') }}</label>
            <div class="col-md-6">
              <select
                v-model="selectedExpiry"
                class="form-select"
              >
                <option
                  v-for="opt in expiryChoices"
                  :key="opt.value"
                  :value="opt.value"
                >
                  {{ opt.text }}
                </option>
              </select>
            </div>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>
<script>
/* global maxSecretExpire */

import appCrypto from '../crypto.js'
import OTSMeta from '../ots-meta'

const defaultExpiryChoices = [
  90 * 86400, // 90 days
  30 * 86400, // 30 days
  7 * 86400, // 7 days
  3 * 86400, // 3 days
  24 * 3600, // 1 day
  12 * 3600, // 12 hours
  4 * 3600, // 4 hours
  60 * 60, // 1 hour
  30 * 60, // 30 minutes
  5 * 60, // 5 minutes
]

const passwordCharset = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ'
const passwordLength = 20

export default {
  computed: {
    expiryChoices() {
      const choices = [{ text: this.$t('expire-default'), value: null }]
      for (const choice of this.$root.customize.expiryChoices || defaultExpiryChoices) {
        if (maxSecretExpire > 0 && choice > maxSecretExpire) {
          continue
        }

        const option = { value: choice }
        if (choice >= 86400) {
          option.text = this.$tc('expire-n-days', Math.round(choice / 86400))
        } else if (choice >= 3600) {
          option.text = this.$tc('expire-n-hours', Math.round(choice / 3600))
        } else if (choice >= 60) {
          option.text = this.$tc('expire-n-minutes', Math.round(choice / 60))
        } else {
          option.text = this.$tc('expire-n-seconds', choice)
        }

        choices.push(option)
      }

      return choices
    },
  },

  created() {
    this.checkWriteAccess()
  },

  data() {
    return {
      canWrite: null,
      secret: '',
      securePassword: null,
      selectedExpiry: null,
    }
  },

  methods: {
    checkWriteAccess() {
      fetch('api/isWritable', {
        credentials: 'same-origin',
        method: 'GET',
        redirect: 'error',
      })
        .then(resp => {
          if (resp.status !== 204) {
            throw new Error(`unexpected status: ${resp.status}`)
          }
          this.canWrite = true
        })
        .catch(() => {
          this.canWrite = false
        })
    },

    // createSecret executes the secret creation after encrypting the secret
    createSecret() {
      if (this.secret.trim().length < 1) {
        return false
      }

      this.securePassword = [...window.crypto.getRandomValues(new Uint8Array(passwordLength))]
        .map(n => passwordCharset[n % passwordCharset.length])
        .join('')

      const meta = new OTSMeta()
      meta.secret = this.secret

      if (this.$refs.createSecretFiles) {
        for (const f of [...this.$refs.createSecretFiles.files]) {
          meta.files.push(f)
        }
      }

      meta.serialize()
        .then(secret => appCrypto.enc(secret, this.securePassword))
        .then(secret => {
          let reqURL = 'api/create'
          if (this.selectedExpiry !== null) {
            reqURL = `api/create?expire=${this.selectedExpiry}`
          }

          return fetch(reqURL, {
            body: JSON.stringify({ secret }),
            headers: {
              'content-type': 'application/json',
            },
            method: 'POST',
          })
            .then(resp => {
              if (resp.status !== 201) {
                // Server says "no"
                this.$emit('error', this.$t('alert-something-went-wrong'))
                return
              }

              resp.json()
                .then(data => {
                  this.$root.navigate({
                    path: '/display-secret-url',
                    query: {
                      expiresAt: data.expires_at,
                      secretId: data.secret_id,
                      securePassword: this.securePassword,
                    },
                  })
                })
            })
            .catch(() => {
              // Network error
              this.$emit('error', this.$t('alert-something-went-wrong'))
            })
        })

      return false
    },
  },

  name: 'AppCreate',
}
</script>
