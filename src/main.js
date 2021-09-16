import Vue from 'vue'
import VueI18n from 'vue-i18n'

import BootstrapVue from 'bootstrap-vue'

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import './style.scss'

import app from './app.vue'
import messages from './langs/langs.js'

Vue.use(BootstrapVue)
Vue.use(VueI18n)

const i18n = new VueI18n({
  locale,
  fallbackLocale: 'en',
  messages,
})

new Vue({
  el: '#app',
  components: { app },
  data: { version },
  i18n,
  render: createElement => createElement('app'),
})
