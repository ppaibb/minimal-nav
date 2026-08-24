import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import Index from '../pages/Index.vue'
import Admin from '../pages/Admin.vue'
import Docs from '../pages/Docs.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: Index,
  },
  {
    path: '/admin',
    name: 'Admin',
    component: Admin,
  },
  {
    path: '/docs/:id?',
    name: 'Docs',
    component: Docs,
  },
  {
    path: '/docs/ai',
    name: 'AiDocs',
    component: Docs,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
