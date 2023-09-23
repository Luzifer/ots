/* eslint-disable sort-imports */
/* global version */
import Vue from 'vue'
import VueI18n from 'vue-i18n'

import BootstrapVue from 'bootstrap-vue'

import './style.scss'

import app from './app.vue'
import messages from './langs/langs.js'

Vue.use(BootstrapVue)
Vue.use(VueI18n)

const cookieSet = Object.fromEntries(document.cookie.split('; ')
  .map(el => el.split('=')
    .map(el => decodeURIComponent(el))))

const i18n = new VueI18n({
  fallbackLocale: 'en',
  locale: cookieSet.lang || navigator?.language || 'en',
  messages,
})

new Vue({
  components: { app },
  data: { version },
  el: '#app',
  i18n,
  name: 'OTSAppInterface',
  render: createElement => createElement('app'),
})
