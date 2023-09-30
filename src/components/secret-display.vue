<!-- eslint-disable vue/no-v-html -->
<template>
  <div class="card border-primary-subtle mb-3">
    <div
      class="card-header bg-primary-subtle"
      v-html="$t('title-reading-secret')"
    />
    <div class="card-body">
      <template v-if="!secret">
        <p v-html="$t('text-pre-reveal-hint')" />
        <button
          class="btn btn-success"
          @click="requestSecret"
        >
          {{ $t('btn-reveal-secret') }}
        </button>
      </template>
      <template v-else>
        <div class="input-group mb-3">
          <textarea
            class="form-control"
            readonly
            :value="secret"
            rows="4"
          />
          <div class="d-flex align-items-start p-0">
            <div
              class="btn-group-vertical"
              role="group"
            >
              <app-clipboard-button :content="secret" />
              <a
                class="btn btn-secondary"
                :href="secretContentBlobURL"
                download
              >
                <i class="fas fa-fw fa-download" />
              </a>
              <app-qr-button :qr-content="secret" />
            </div>
          </div>
        </div>
        <p v-html="$t('text-hint-burned')" />
      </template>
    </div>
  </div>
</template>
<script>
import appClipboardButton from './clipboard-button.vue'
import appCrypto from '../crypto.js'
import appQrButton from './qr-button.vue'

export default {
  components: { appClipboardButton, appQrButton },

  data() {
    return {
      popover: null,
      secret: null,
      secretContentBlobURL: null,
    }
  },

  methods: {
    // requestSecret requests the encrypted secret from the backend
    requestSecret() {
      window.history.replaceState({}, '', window.location.href.split('#')[0])
      fetch(`api/get/${this.secretId}`)
        .then(resp => {
          if (resp.status === 404) {
            // Secret has already been consumed
            this.$emit('error', this.$t('alert-secret-not-found'))
            return
          }

          if (resp.status !== 200) {
            // Some other non-200: Something(tm) was wrong
            this.$emit('error', this.$t('alert-something-went-wrong'))
            return
          }

          resp.json()
            .then(data => {
              const secret = data.secret
              if (!this.securePassword) {
                this.secret = secret
                return
              }

              appCrypto.dec(secret, this.securePassword)
                .then(secret => {
                  this.secret = secret
                })
                .catch(() => {
                  this.$emit('error', this.$t('alert-something-went-wrong'))
                })
            })
        })
        .catch(() => {
          // Network error
          this.$emit('error', this.$t('alert-something-went-wrong'))
        })
    },
  },

  name: 'AppSecretDisplay',
  props: {
    secretId: {
      required: true,
      type: String,
    },

    securePassword: {
      default: null,
      required: false,
      type: String,
    },
  },

  watch: {
    secret(to) {
      if (this.secretContentBlobURL) {
        window.URL.revokeObjectURL(this.secretContentBlobURL)
      }
      this.secretContentBlobURL = window.URL.createObjectURL(new Blob([to], { type: 'text/plain' }))
    },
  },
}
</script>
