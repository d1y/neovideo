import { RouteRecordRaw } from 'vue-router'

// 路由表
const constantRouterMap: RouteRecordRaw[] = [
  // ************* 前台路由 **************
  {
    path: '/',
    name: 'home',
    redirect: '/index',
    component: () => import('@/views/home.vue'),
    children: [
      {
        path: 'index',
        name: 'index',
        component: () => import('@/views/home/index.vue'),
      },
      {
        path: 'search',
        name: 'search',
        component: () => import('@/views/home/search.vue'),
      },
      {
        path: 'dianying',
        name: 'dianying',
        component: () => import('@/views/home/dianying.vue'),
        props: (route) => ({
          id: route.query.id,
        }),
      },
      {
        path: 'zhuanti',
        name: 'zhuanti',
        component: () => import('@/views/home/zhuanti.vue')
      },
    ],
  },
  {
    path: '/detail/:id',
    name: 'detail',
    component: () => import('@/views/detail.vue'),
    props: (route) => ({
      id: route.params.id
    })
  },
]

export default constantRouterMap
