<!-- eslint-disable vue/no-v-html -->
<template>
  <div id="app">
    <app-navbar />

    <div class="container mt-4">
      <div
        v-if="error"
        class="row justify-content-center"
      >
        <div class="col-12 col-md-8">
          <div
            class="alert alert-danger"
            role="alert"
            v-html="error"
          />
        </div>
      </div>

      <div class="row">
        <div class="col">
          <router-view @error="displayError" />
        </div>
      </div>

      <div
        class="row mt-4"
      >
        <div class="col form-text text-center">
          <span
            v-if="!$root.customize.disablePoweredBy"
            class="mx-2"
          >
            {{ $t('text-powered-by') }}
            <a href="https://github.com/Luzifer/ots"><i class="fab fa-github" /> OTS</a>
            {{ $root.version }}
          </span>
          <span
            v-for="link in $root.customize.footerLinks"
            :key="link.url"
            class="mx-2"
          >
            <a :href="link.url">{{ link.name }}</a>
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import AppNavbar from './components/navbar.vue'

export default {
  components: {
    AppNavbar,
  },

  created() {
    this.$root.navigate('/')
  },

  data() {
    return {
      error: '',
    }
  },

  methods: {
    displayError(error) {
      this.error = error
    },

    // hashLoad reacts on a changed window hash an starts the diplaying of the secret
    hashLoad() {
      const hash = decodeURIComponent(window.location.hash)
      if (hash.length === 0) {
        return
      }

      const parts = hash.substring(1).split('|')
      const secretId = parts[0]
      let securePassword = null

      if (parts.length === 2) {
        securePassword = parts[1]
      }

      this.$root.navigate({
        path: '/secret',
        query: {
          secretId,
          securePassword,
        },
      })
    },
  },

  // Trigger initialization functions
  mounted() {
    window.onhashchange = this.hashLoad
    this.hashLoad()

    if (!this.$root.isSecureEnvironment) {
      this.error = this.$t('alert-insecure-environment')
    }
  },

  name: 'App',
}
</script>
