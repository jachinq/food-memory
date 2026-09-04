import { createRouter, createWebHistory } from 'vue-router'
import AppLayout from '../layouts/AppLayout.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/unlock', component: () => import('../views/Unlock.vue') },
    {
      path: '/',
      component: AppLayout,
      children: [
        { path: '', name: 'home', component: () => import('../views/Home.vue') },
        { path: 'dishes', name: 'dishes', component: () => import('../views/DishList.vue') },
        { path: 'dishes/new', name: 'dish-new', component: () => import('../views/DishForm.vue') },
        { path: 'dishes/:id', name: 'dish-detail', component: () => import('../views/DishDetail.vue') },
        { path: 'dishes/:id/edit', name: 'dish-edit', component: () => import('../views/DishForm.vue') },
        { path: 'dishes/:id/records/new', name: 'record-new', component: () => import('../views/CookRecordForm.vue') },
        { path: 'recook', name: 'recook', component: () => import('../views/RecookPlan.vue') },
      ],
    },
  ],
  scrollBehavior() {
    return { top: 0 }
  },
})

export default router
