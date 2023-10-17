<!-- eslint-disable vue/no-v-html -->
<template>
  <div id="app">
    <app-navbar />

    <div class="container mt-4">
      <div
        v-if="error"
        class="row justify-content-center"
      >
        <div class="col-8">
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
        v-if="!$root.customize.disablePoweredBy"
        class="row mt-4"
      >
        <div class="col form-text text-center">
          {{ $t('text-powered-by') }}
          <a href="https://github.com/Luzifer/ots"><i class="fab fa-github" /> OTS</a>
          {{ $root.version }}
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
  },

  name: 'App',
}
</script>
