import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

import AuthView from './views/AuthView.vue'
import TasksView from './views/TasksView.vue'

const routes: Array<RouteRecordRaw> = [
  { path: '/', redirect: '/auth' },
  {
    path: '/auth',
    name: 'auth',
    component: AuthView,
  },
  {
    path: '/tasks',
    name: 'tasks',
    component: TasksView,
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL), // Vite usa import.meta.env
  routes,
})

export default router
