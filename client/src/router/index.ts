import { createWebHistory, createRouter } from 'vue-router'

import { useAuthStore } from '../stores/auth.ts'

import HomeView from '../views/HomeView.vue'
import AboutView from '../views/AboutView.vue'
import StatusView from '../views/StatusView.vue'
import DashboardView from '../views/DashboardView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'


const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView
  },
  {
    path: '/about',
    name: 'about',
    component: AboutView
  },
  {
    path: '/status',
    name: 'status',
    component: StatusView,
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: DashboardView,
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView
  },
  {
    path: '/register',
    name: 'register',
    component: RegisterView
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// TODO: Тут так же можно сделать логику GuestOnly, где login и homePage доступны только гостям
router.beforeEach((to, _from, next) => {
  const auth = useAuthStore()

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    next('/login') // Используем next() для редиректа
  } else {
    next() // Всегда вызывайте next()!
  }
})

router.beforeEach(async (to, _from) => {
  const auth = useAuthStore()

  if (to.meta.requiresAuth && !auth.isAuthenticated && to.name !== 'Login') {
    {
      return { name: 'Login' }
    }
  }
})

export default router
