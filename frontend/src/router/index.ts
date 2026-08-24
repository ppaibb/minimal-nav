import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import Index from '../pages/Index.vue'
import Admin from '../pages/Admin.vue'
import AiDocs from '../pages/AiDocs.vue'

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
    path: '/docs/ai',
    name: 'AiDocs',
    component: AiDocs,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
