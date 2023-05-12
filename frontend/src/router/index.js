import { createRouter, createWebHistory } from 'vue-router'
import Main from '../views/Main'
import Authorization from '@/views/Authorization'
import Registration from '@/views/Registration'
import UploadModel from '@/views/UploadModel'
import EditModel from '@/views/EditModel'

const routes = [
  {
    path: '/',
    name: 'main',
    component: Main
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
  },
  {
    path: '/:id',
    component: EditModel
  },
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
