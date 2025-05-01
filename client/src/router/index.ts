import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import NotFoundView from '../views/NotFoundView.vue'
import FrontPageView from '../views/FrontPageView.vue'

const routes = [
  {
    path: '/',
    name: 'main',
    component: FrontPageView
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView
  },
  {
    path: '/404',
    name: '404',
    component: NotFoundView
  },
  {
    path: '/:catchAll(.*)*',
    redirect: '/404'
  }
]

const router = createRouter({
  history: createWebHistory(), // Creates clean URLs without #
  routes // short for `routes: routes`
})

export default router
