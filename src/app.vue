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
        <span v-if="!customize.disableAppTitle">{{ customize.appTitle || 'OTS - One Time Secrets' }}</span>
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
            v-if="mode == 'create' && !secretId"
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
            <b-button
              :disabled="secret.trim().length < 1"
              variant="success"
              @click="createSecret"
            >
              {{ $t('btn-create-secret') }}
            </b-button>
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
                <b-input-group-append
                  v-if="!customize.disableQRSupport"
                >
                  <b-button
                    id="secret-url-qrcode"
                    :disabled="!secretQRDataURL"
                    variant="primary"
                  >
                    <i class="fas fa-qrcode" />
                  </b-button>
                </b-input-group-append>
              </b-input-group>
            </b-form-group>
            <p v-html="$t('text-burn-hint')" />

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
              <b-form-group>
                <b-form-textarea
                  max-rows="25"
                  readonly
                  rows="5"
                  :value="secret"
                />
              </b-form-group>
              <p v-html="$t('text-hint-burned')" />
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

const passwordCharset = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ'
const passwordLength = 20

export default {
  computed: {
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
      customize: {},
      darkTheme: false,
      error: '',
      explanationShown: false,
      mode: 'create',
      secret: '',
      secretId: '',
      secretQRDataURL: '',
      securePassword: '',
      showError: false,
    }
  },

  methods: {
    // createSecret executes the secret creation after encrypting the secret
    createSecret() {
      if (this.secret.trim().length < 1) {
        return false
      }

      this.securePassword = [...window.crypto.getRandomValues(new Uint8Array(passwordLength))]
        .map(n => passwordCharset[n % passwordCharset.length])
        .join('')
      crypto.enc(this.secret, this.securePassword)
        .then(secret => fetch('api/create', {
          body: JSON.stringify({ secret }),
          headers: {
            'content-type': 'application/json',
          },
          method: 'POST',
        })
          .then(resp => resp.json())
          .then(data => ({ data }))
          .then(resp => {
            console.warn(resp)
            this.secretId = resp.data.secret_id
            this.secret = ''

            // Give the interface a moment to transistion and focus
            window.setTimeout(() => this.$refs.secretUrl.focus(), 100)
          })
          .catch(err => {
            switch (err.response.status) {
            case 404:
              // Mock for interface testing
              this.secretId = 'foobar'
              break
            default:
              this.error = this.$t('alert-something-went-wrong')
              this.showError = true
            }
          }))

      return false
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
        .then(resp => resp.json())
        .then(data => ({ data }))
        .then(resp => {
          const secret = resp.data.secret
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
        .catch(err => {
          switch (err.response.status) {
          case 404:
            this.error = this.$t('alert-secret-not-found')
            this.showError = true
            break
          default:
            this.error = this.$t('alert-something-went-wrong')
            this.showError = true
          }
        })
    },
  },

  // Trigger initialization functions
  mounted() {
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
