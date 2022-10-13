import { createRouter, createWebHistory } from 'vue-router'
import Main from '../views/Main'
import ModelPageId from '@/views/ModelPageId'
import Authorization from '@/views/Authorization'

const routes = [
  {
    path: '/',
    name: 'main',
    component: Main
  },
  {
    path: '/:id',
    component: ModelPageId
  },
  {
    path: '/authorization',
    name: 'authorization',
    component: Authorization,
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
