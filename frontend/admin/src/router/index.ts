import { createRouter, createWebHashHistory } from "vue-router";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: "/",
      component: () => import("@/components/Layout.vue"),
      redirect: "/dashboard",
      children: [
        {
          path: "dashboard",
          component: () => import("@/views/Dashboard.vue"),
        },
        {
          path: "jiexi",
          component: () => import("@/views/Jiexi.vue"),
        },
        {
          path: "maccms",
          component: () => import("@/views/Maccms.vue"),
        },
        {
          path: "spider",
          component: () => import("@/views/Spider.vue"),
        },
      ],
    },
  ],
});

export default router;
