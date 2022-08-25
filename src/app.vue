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
        <i class="fas fa-user-secret" /> OTS - One Time Secrets
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
          <b-nav-form class="ml-2">
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
              <b-form-input
                ref="secretUrl"
                :value="secretUrl"
                readonly
                @focus="$refs.secretUrl.select()"
              />
            </b-form-group>
            <p v-html="$t('text-burn-hint')" />
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

      <b-row class="mt-5">
        <b-col class="footer">
          {{ $t('text-powered-by') }} <a href="https://github.com/Luzifer/ots"><i class="fab fa-github" /> Luzifer/OTS</a> {{ $root.version }}
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import axios from 'axios'
import AES from 'gibberish-aes/src/gibberish-aes'

const passwordCharset = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ'
const passwordLength = 20

export default {
  name: 'App',

  data() {
    return {
      error: '',
      explanationShown: false,
      mode: 'create',
      secret: '',
      securePassword: '',
      secretId: '',
      showError: false,
      darkTheme: false,
    }
  },

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

  watch: {
    darkTheme(to) {
      window.setTheme(to ? 'dark' : 'light')
    },
  },

  // Trigger initialization functions
  mounted() {
    this.darkTheme = window.getTheme() === 'dark'
    window.onhashchange = this.hashLoad
    this.hashLoad()
  },

  methods: {
    // createSecret executes the secret creation after encrypting the secret
    createSecret() {
      this.securePassword = [...window.crypto.getRandomValues(new Uint8Array(passwordLength))]
        .map(n => passwordCharset[n % passwordCharset.length])
        .join('')
      const secret = AES.enc(this.secret, this.securePassword)

      axios.post('api/create', { secret })
        .then(resp => {
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
        })

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
      axios.get(`api/get/${this.secretId}`)
        .then(resp => {
          let secret = resp.data.secret
          if (this.securePassword) {
            secret = AES.dec(secret, this.securePassword)
          }
          this.secret = secret
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
}
</script>
