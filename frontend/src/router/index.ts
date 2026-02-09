import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/select',
      component: () => import('../views/SelectCar.vue'),
    },
    {
      path: '/car/:id',
      component: () => import('../views/CarMap.vue'),
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/select',
    },
  ],
})

export default router
