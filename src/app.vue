<template>
  <div id="app">
    <b-navbar
      toggleable="lg"
      type="dark"
      variant="primary"
    >
      <b-navbar-brand
        href="#"
        @click="newSecret"
      >
        <i
          v-if="!customize.appIcon"
          class="fas fa-user-secret mr-1"
        />
        <img
          v-else
          class="mr-1"
          :src="customize.appIcon"
        >
        <span v-if="!customize.disableAppTitle">{{ customize.appTitle }}</span>
      </b-navbar-brand>

      <b-navbar-toggle target="nav-collapse" />

      <b-collapse
        id="nav-collapse"
        is-nav
      >
        <b-navbar-nav class="ml-auto">
          <b-nav-item @click="explanationShown = !explanationShown">
            <i class="fas fa-question" /> {{ $t('btn-show-explanation') }}
          </b-nav-item>
          <b-nav-item @click="newSecret">
            <i class="fas fa-plus" /> {{ $t('btn-new-secret') }}
          </b-nav-item>
          <b-nav-form
            v-if="!customize.disableThemeSwitcher"
            class="ml-2"
          >
            <b-form-checkbox
              v-model="darkTheme"
              switch
            >
              <i class="fas fa-moon" />&ZeroWidthSpace;
            </b-form-checkbox>
          </b-nav-form>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>

    <b-container class="mt-4">
      <b-row class="justify-content-center">
        <b-col md="8">
          <b-alert
            v-model="showError"
            variant="danger"
            dismissible
            v-html="error"
          />
        </b-col>
      </b-row>

      <b-row>
        <b-col>
          <!-- Explanation -->
          <b-card
            v-if="explanationShown"
            class="mb-3"
            border-variant="primary"
            header-bg-variant="primary"
            header-text-variant="white"
          >
            <span
              slot="header"
              v-html="$t('title-explanation')"
            />
            <ul>
              <li v-for="explanation in $t('items-explanation')">
                {{ explanation }}
              </li>
            </ul>
          </b-card>

          <!-- Creation dialog -->
          <b-card
            v-if="mode == 'create' && !secretId && canWrite"
            border-variant="primary"
            header-bg-variant="primary"
            header-text-variant="white"
          >
            <span
              slot="header"
              v-html="$t('title-new-secret')"
            />
            <b-form-group :label="$t('label-secret-data')">
              <b-form-textarea
                id="secret"
                v-model="secret"
                max-rows="25"
                rows="5"
              />
            </b-form-group>
            <b-row>
              <b-col
                cols="12"
                md="6"
                order="2"
                order-md="1"
              >
                <b-button
                  :disabled="secret.trim().length < 1"
                  variant="success"
                  @click="createSecret"
                >
                  {{ $t('btn-create-secret') }}
                </b-button>
              </b-col>
              <b-col
                v-if="!customize.disableExpiryOverride"
                cols="12"
                md="6"
                order="1"
                order-md="2"
              >
                <b-form-group
                  :label="$t('label-expiry')"
                  label-for="expiry"
                  label-align-md="right"
                  label-cols-md
                >
                  <b-form-select
                    id="expiry"
                    v-model="selectedExpiry"
                    :options="expiryChoices()"
                  />
                </b-form-group>
              </b-col>
            </b-row>
          </b-card>

          <!-- Creation disabled -->
          <b-card
            v-if="mode == 'create' && !secretId && canWrite === false"
            border-variant="info"
            header-bg-variant="info"
            header-text-variant="white"
          >
            <span
              slot="header"
              v-html="$t('title-secret-create-disabled')"
            />
            <p v-html="$t('text-secret-create-disabled')" />
          </b-card>

          <!-- Secret created, show secret URL -->
          <b-card
            v-if="mode == 'create' && secretId"
            border-variant="success"
            header-bg-variant="success"
            header-text-variant="white"
          >
            <span
              slot="header"
              v-html="$t('title-secret-created')"
            />
            <p v-html="$t('text-pre-url')" />
            <b-form-group>
              <b-input-group>
                <b-form-input
                  ref="secretUrl"
                  :value="secretUrl"
                  readonly
                  @focus="$refs.secretUrl.select()"
                />
                <b-input-group-append>
                  <b-button
                    v-if="hasClipboard"
                    :disabled="!secretUrl"
                    :variant="copyToClipboardSuccess ? 'success' : 'primary'"
                    @click="copySecretUrl"
                  >
                    <i class="fas fa-clipboard" />
                  </b-button>
                  <b-button
                    v-if="!customize.disableQRSupport"
                    id="secret-url-qrcode"
                    :disabled="!secretQRDataURL"
                    variant="secondary"
                  >
                    <i class="fas fa-qrcode" />
                  </b-button>
                </b-input-group-append>
              </b-input-group>
            </b-form-group>
            <p v-html="$t('text-burn-hint')" />
            <p v-if="secretExpiry">
              {{ $t('text-burn-time') }}
              <strong>{{ secretExpiry.toLocaleString() }}</strong>
            </p>

            <b-popover
              v-id="!customize.disableQRSupport"
              target="secret-url-qrcode"
              triggers="focus"
              placement="leftbottom"
            >
              <b-img
                height="200"
                :src="secretQRDataURL"
                width="200"
              />
            </b-popover>
          </b-card>


          <!-- Display secret -->
          <b-card
            v-if="mode == 'view'"
            border-variant="primary"
            header-bg-variant="primary"
            header-text-variant="white"
          >
            <span
              slot="header"
              v-html="$t('title-reading-secret')"
            />
            <template v-if="!secret">
              <p v-html="$t('text-pre-reveal-hint')" />
              <b-button
                variant="success"
                @click="requestSecret"
              >
                {{ $t('btn-reveal-secret') }}
              </b-button>
            </template>
            <template v-else>
              <b-input-group>
                <b-form-textarea
                  max-rows="25"
                  readonly
                  rows="4"
                  :value="secret"
                />
                <b-input-group-text class="d-flex align-items-start p-0">
                  <b-button-group vertical>
                    <b-button
                      v-if="hasClipboard"
                      :disabled="!secretUrl"
                      :variant="copyToClipboardSuccess ? 'success' : 'primary'"
                      title="Copy Secret to Clipboard"
                      @click="copySecret"
                    >
                      <i class="fas fa-fw fa-clipboard" />
                    </b-button>
                    <b-button
                      :href="`data:text/plain;charset=UTF-8,${secret}`"
                      download
                      title="Download Secret as Text File"
                    >
                      <i class="fas fa-fw fa-download" />
                    </b-button>
                    <b-button
                      v-if="!customize.disableQRSupport && secretContentQRDataURL"
                      id="secret-data-qrcode"
                      variant="secondary"
                      title="Display Content as QR-Code"
                    >
                      <i class="fas fa-fw fa-qrcode" />
                    </b-button>
                  </b-button-group>
                </b-input-group-text>
              </b-input-group>
              <p v-html="$t('text-hint-burned')" />

              <b-popover
                v-id="!customize.disableQRSupport"
                target="secret-data-qrcode"
                triggers="focus"
                placement="leftbottom"
              >
                <b-img
                  height="200"
                  :src="secretContentQRDataURL"
                  width="200"
                />
              </b-popover>
            </template>
          </b-card>
        </b-col>
      </b-row>

      <b-row
        v-if="!customize.disablePoweredBy"
        class="mt-5"
      >
        <b-col class="footer">
          {{ $t('text-powered-by') }} <a href="https://github.com/Luzifer/ots"><i class="fab fa-github" /> OTS</a> {{ $root.version }}
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import crypto from './crypto.js'
import qrcode from 'qrcode'

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
    hasClipboard() {
      return Boolean(navigator.clipboard && navigator.clipboard.writeText)
    },

    secretUrl() {
      return [
        window.location.href,
        encodeURIComponent([
          this.secretId,
          this.securePassword,
        ].join('|')),
      ].join('#')
    },
  },

  data() {
    return {
      canWrite: null,
      copyToClipboardSuccess: false,
      customize: {},
      darkTheme: false,
      error: '',
      explanationShown: false,
      mode: 'create',
      secret: '',
      secretExpiry: null,
      secretId: '',
      secretQRDataURL: '',
      secretContentQRDataURL: '',
      securePassword: '',
      selectedExpiry: null,
      showError: false,
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

    copySecret() {
      navigator.clipboard.writeText(this.secret)
        .then(() => {
          this.copyToClipboardSuccess = true
          window.setTimeout(() => {
            this.copyToClipboardSuccess = false
          }, 500)
        })
    },

    copySecretUrl() {
      navigator.clipboard.writeText(this.secretUrl)
        .then(() => {
          this.copyToClipboardSuccess = true
          window.setTimeout(() => {
            this.copyToClipboardSuccess = false
          }, 500)
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
      crypto.enc(this.secret, this.securePassword)
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
                this.error = this.$t('alert-something-went-wrong')
                this.showError = true
                return
              }

              resp.json()
                .then(data => {
                  this.secretId = data.secret_id
                  this.secret = ''

                  if (data.expires_at) {
                    this.secretExpiry = new Date(data.expires_at)
                  }

                  // Give the interface a moment to transistion and focus
                  window.setTimeout(() => this.$refs.secretUrl.focus(), 100)
                })
            })
            .catch(err => {
            // Network error
              this.error = this.$t('alert-something-went-wrong')
              this.showError = true
            })
        })

      return false
    },

    expiryChoices() {
      const choices = [{ text: this.$t('expire-default'), value: null }]
      for (const choice of this.customize.expiryChoices || defaultExpiryChoices) {
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

    // hashLoad reacts on a changed window hash an starts the diplaying of the secret
    hashLoad() {
      const hash = decodeURIComponent(window.location.hash)
      if (hash.length === 0) {
        return
      }

      const parts = hash.substring(1).split('|')
      if (parts.length === 2) {
        this.securePassword = parts[1]
      }
      this.secretId = parts[0]
      this.mode = 'view'
    },

    // newSecret removes the window hash and therefore returns to "new secret" mode
    newSecret() {
      location.href = location.href.split('#')[0]
    },

    // requestSecret requests the encrypted secret from the backend
    requestSecret() {
      fetch(`api/get/${this.secretId}`)
        .then(resp => {
          if (resp.status === 404) {
            // Secret has already been consumed
            this.error = this.$t('alert-secret-not-found')
            this.showError = true
            return
          }

          if (resp.status !== 200) {
            // Some other non-200: Something(tm) was wrong
            this.error = this.$t('alert-something-went-wrong')
            this.showError = true
            return
          }

          resp.json()
            .then(data => {
              const secret = data.secret
              if (!this.securePassword) {
                this.secret = secret
                return
              }

              crypto.dec(secret, this.securePassword)
                .then(secret => {
                  this.secret = secret
                })
                .catch(() => {
                  this.error = this.$t('alert-something-went-wrong')
                  this.showError = true
                })
            })
        })
        .catch(err => {
          // Network error
          this.error = this.$t('alert-something-went-wrong')
          this.showError = true
        })
    },
  },

  // Trigger initialization functions
  mounted() {
    this.checkWriteAccess()
    this.customize = window.OTSCustomize
    this.darkTheme = window.getTheme() === 'dark'
    window.onhashchange = this.hashLoad
    this.hashLoad()
  },

  name: 'App',

  watch: {
    darkTheme(to) {
      window.setTheme(to ? 'dark' : 'light')
    },

    secret(to) {
      if (this.customize.disableQRSupport || !to) {
        return
      }

      qrcode.toDataURL(to, { width: 200 })
        .then(url => {
          this.secretContentQRDataURL = url
        })
        .catch(() => this.secretContentQRDataURL = null)
    },

    secretUrl(to) {
      if (this.customize.disableQRSupport) {
        return
      }

      qrcode.toDataURL(to, { width: 200 })
        .then(url => {
          this.secretQRDataURL = url
        })
    },
  },
}
</script>
