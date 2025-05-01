import { createI18n } from "vue-i18n";

import messages from './langs/langs.js'

const cookieSet = Object.fromEntries(document.cookie.split('; ')
  .map(el => el.split('=')
    .map(el => decodeURIComponent(el))))

const i18n = createI18n({
  legacy: false,
  fallbackLocale: 'en',
  locale: cookieSet.lang || navigator?.language || 'en',
  messages,
})

export default i18n
