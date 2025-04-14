import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

import LoginView from './views/LoginView.vue'
import TasksView from './views/TasksView.vue'

const routes: Array<RouteRecordRaw> = [
  { path: '/', redirect: '/login' },
  {
    path: '/login',
    name: 'login',
    component: LoginView,
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
