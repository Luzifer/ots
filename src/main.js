import Vue from 'vue'
import VueI18n from 'vue-i18n'

import BootstrapVue from 'bootstrap-vue'

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import 'bootswatch/dist/flatly/bootstrap.css'

import app from './app.vue'
import messages from './langs/langs.js'

console.log(['app.js loaded',app,messages])

Vue.use(BootstrapVue)
Vue.use(VueI18n)

const i18n = new VueI18n({
  locale,
  fallbackLocale: 'en',
  messages,
})

new Vue({
  components: { app },

  data: {
    error: null,
    secret: '',
    securePassword: null,
    view: 'create',
    version,
  },

  el: '#app',
  i18n,
  render: createElement => createElement('app'),
})
