import { createRouter, createWebHistory } from 'vue-router'
import root from './root'

const router = createRouter({
  history: createWebHistory(),
  routes: root,
})

router.beforeEach(async (to, from, next) => {
  next()
})

router.afterEach((_to) => {
  // 回到顶部
  document.getElementById('html')?.scrollTo(0, 0)
})

export default router
