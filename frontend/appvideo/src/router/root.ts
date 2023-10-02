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
      },
      {
        path: 'dianshiju',
        name: 'dianshiju',
        component: () => import('@/views/home/dianshiju.vue'),
      },
      {
        path: 'zongyi',
        name: 'zongyi',
        component: () => import('@/views/home/zongyi.vue'),
      },
      {
        path: 'dongman',
        name: 'dongman',
        component: () => import('@/views/home/dongman.vue'),
      },
      {
        path: 'zhuanti',
        name: 'zhuanti',
        component: () => import('@/views/home/zhuanti.vue'),
      },
    ],
  },
  {
    path: '/detail/:id',
    name: 'detail',
    component: () => import('@/views/detail.vue'),
    props: (route) => ({
      id: route.params.id,
      mid: route.query.mid as string,
    })
  },
]

export default constantRouterMap
