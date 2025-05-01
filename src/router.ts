import { createMemoryHistory, createRouter } from 'vue-router'

import AppCreate from './components/create.vue'
import AppDisplayURL from './components/display-url.vue'
import AppExplanation from './components/explanation.vue'
import AppSecretDisplay from './components/secret-display.vue'

const routes = [
  {
    component: AppCreate,
    name: 'create',
    path: '/',
  },
  {
    component: AppDisplayURL,
    name: 'display-secret-url',
    path: '/display-secret-url',
    props: route => ({
      expiresAt: route.query.expiresAt ? new Date(route.query.expiresAt) : null,
      secretId: route.query.secretId,
      securePassword: route.query.securePassword,
    }),
  },
  {
    component: AppExplanation,
    name: 'explanation',
    path: '/explanation',
  },
  {
    component: AppSecretDisplay,
    name: 'secret',
    path: '/secret',
    props: route => ({
      secretId: route.query.secretId,
      securePassword: route.query.securePassword,
    }),
  },
]

const router = createRouter({
  history: createMemoryHistory(),
  routes,
})

export default router
