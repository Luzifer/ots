<template>
  <div id="app">
    <b-navbar toggleable="lg" type="dark" variant="primary">
      <b-navbar-brand href="#" @click="newSecret">
        <i class="fas fa-user-secret"></i> OTS - One Time Secrets
      </b-navbar-brand>

      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav class="ml-auto">
          <b-nav-item @click="newSecret"><i class="fas fa-plus"></i> {{ $t('btn-new-secret') }}</b-nav-item>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>

    <b-container class="mt-4">

      <b-row class="justify-content-center">
        <b-col md="8">
          <b-alert v-model="showError" variant="danger" dismissible v-html="error"></b-alert>
        </b-col>
      </b-row>

      <b-row>
        <b-col>

          <b-card
            border-variant="primary"
            header-bg-variant="primary"
            header-text-variant="white"
            v-if="mode == 'create' && !secretId"
            >
            <span slot="header" v-html="$t('title-new-secret')"></span>
            <b-form-group :label="$t('label-secret-data')">
              <b-form-textarea
                id="secret"
                max-rows="25"
                rows="5"
                v-model="secret"
                >
              </b-form-textarea>
            </b-form-group>
            <b-button variant="success" @click="createSecret">{{ $t('btn-create-secret') }}</b-button>
          </b-card>

          <b-card
            border-variant="success"
            header-bg-variant="success"
            header-text-variant="white"
            v-if="mode == 'create' && secretId"
            >
            <span slot="header" v-html="$t('title-secret-created')"></span>
            <p v-html="$t('text-pre-url')"></p>
            <b-form-group>
              <b-form-input :value="secretUrl" readonly></b-form-input>
            </b-form-group>
            <p v-html="$t('text-burn-hint')"></p>
          </b-card>

          <b-card
            border-variant="primary"
            header-bg-variant="primary"
            header-text-variant="white"
            v-if="mode == 'view'"
            >
            <span slot="header" v-html="$t('title-reading-secret')"></span>
            <template v-if="!secret">
              <p v-html="$t('text-pre-reveal-hint')"></p>
              <b-button variant="success" @click=requestSecret>{{ $t('btn-reveal-secret') }}</b-button>
            </template>
            <template v-else>
              <b-form-group>
                <b-form-textarea
                  max-rows="25"
                  readonly
                  rows="5"
                  :value="secret"
                  >
                </b-form-textarea>
              </b-form-group>
              <p v-html="$t('text-hint-burned')"></p>
            </template>
          </b-card>

        </b-col>
      </b-row>

      <b-row class="mt-5">
        <b-col class="footer">
          {{ $t('text-powered-by') }} <a href="https://github.com/Luzifer/ots"><i class="fab fa-github"></i> Luzifer/OTS</a> {{ $root.version }}
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import axios from 'axios'
import AES from 'gibberish-aes/src/gibberish-aes'

export default {
  name: 'app',

  computed: {
    secretUrl() {
      return `${window.location.href}#${this.secretId}|${this.securePassword}`
    },
  },

  data() {
    return {
      error: '',
      mode: 'create',
      secret: '',
      securePassword: '',
      secretId: '',
      showError: false,
    }
  },

  methods: {
    // createSecret executes the secret creation after encrypting the secret
    createSecret() {
      this.securePassword = Math.random().toString(36)
        .substring(2)
      const secret = AES.enc(this.secret, this.securePassword)

      axios.post('api/create', { secret })
        .then(resp => {
          this.secretId = resp.data.secret_id
          this.secret = ''
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
      const hash = window.location.hash
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

  // Trigger initialization functions
  mounted() {
    window.onhashchange = this.hashLoad
    this.hashLoad()
  },
}
</script>

<style>
textarea {
  font-family: monospace;
}
.footer {
  color: #2f2f2f;
  font-size: 0.9em;
  text-align: center;
}
</style>
