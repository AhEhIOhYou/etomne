import { createRouter, createWebHistory } from 'vue-router'
import Main from '../views/Main'
import ModelPageId from '@/views/ModelPageId'
import Authorization from '@/views/Authorization'
import Registration from '@/views/Registration'
import UploadModel from '@/views/UploadModel'

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
  },
  {
    path: '/registration',
    name: 'registration',
    component: Registration,
  },
  {
    path: '/uploadmodel',
    name: 'uploadmodel',
    component: UploadModel,
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
