import {createRouter, createWebHistory} from 'vue-router';

import AuthView from "@/views/AuthView.vue";
import TasksView from "@/views/TasksView.vue";

const routes = [{
    path: '/', redirect: '/auth'
}, {
    path: '/auth', name: 'auth', component: AuthView,
}, {
    path: '/tasks', name: 'tasks', component: TasksView,
}]

const router = createRouter({
    history: createWebHistory(process.env.BASE_URL), routes,
})

export default router