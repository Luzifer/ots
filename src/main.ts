import { createApp, h } from 'vue'

import './style.scss'
import '@fortawesome/fontawesome-free/css/all.css' // All FA free icons
import appView from './app.vue'

import i18n from './i18n.ts'
import router from './router.ts'

const app = createApp({
  name: 'OTS',
  render() { return h(appView) },
})

app.use(i18n)
app.use(router)

app.mount('#app')
