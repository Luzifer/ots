/* eslint-disable sort-imports */
/* global version */

import Vue from 'vue'
import VueI18n from 'vue-i18n'
import VueRouter from 'vue-router'

import './style.scss'

import app from './app.vue'
import messages from './langs/langs.js'
import router from './router'

Vue.use(VueI18n)
Vue.use(VueRouter)

const cookieSet = Object.fromEntries(document.cookie.split('; ')
  .map(el => el.split('=')
    .map(el => decodeURIComponent(el))))

const i18n = new VueI18n({
  fallbackLocale: 'en',
  locale: cookieSet.lang || navigator?.language || 'en',
  messages,
})

Vue.mixin({
  beforeRouteLeave(_to, _from, next) {
    // Before leaving the component, reset the errors the component displayed
    this.$emit('error', null)
    next()
  },
})

new Vue({
  components: { app },

  computed: {
    isSecureEnvironment() {
      return Boolean(window.crypto.subtle)
    },
  },

  data: {
    customize: {},
    theme: 'auto',
    version,
  },

  el: '#app',
  i18n,

  methods: {
    navigate(to) {
      this.$router.replace(to)
        .catch(err => {
          if (VueRouter.isNavigationFailure(err, VueRouter.NavigationFailureType.duplicated)) {
            // Hide duplicate nav errors
            return
          }
          throw err
        })
    },
  },

  mounted() {
    this.customize = window.OTSCustomize
    this.theme = window.getThemeFromStorage()

    window.matchMedia('(prefers-color-scheme: light)').addEventListener('change', () => {
      window.refreshTheme()
    })
  },

  name: 'OTS',
  render: createElement => createElement('app'),
  router,

  watch: {
    theme(to) {
      window.setTheme(to)
    },
  },
})
