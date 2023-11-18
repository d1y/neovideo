<template>
  <div class="common-layout">
    <el-container class="w-full h-full">
      <el-aside class="h-full bg-[#191a23] box-border p-4" width="240px">
        <div class="flex items-center justify-center mb-4">
          <h1 class="inline-flex text-white font-bold text-2xl">后台管理</h1>
        </div>
        <div class="text-white pl-4">
          <p class="cursor-pointer m-2 mb-4 p-2" :style="{
            backgroundColor: currentPath == item.path ? '#4d70ff' : '',
            borderRadius: `4px`,
          }" v-for="item in menus" @click="$router.push(item.path)">{{ item.title }}</p>
        </div>
      </el-aside>
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import router from '@/router'
import { onMounted } from "vue"
import { useRoute } from "vue-router"

const route = useRoute()

const currentPath = ref('')
router.afterEach((to)=> {
  currentPath.value = to.path
})
onMounted(()=> {
  currentPath.value = route.path
})
const menus = [
  {
    title: "系统面板",
    path: "/dashboard",
  },
  {
    title: "解析源管理",
    path: "/jiexi",
  },
  {
    title: "苹果CMS管理",
    path: "/maccms",
  },
  {
    title: "爬虫任务管理",
    path: "/spider",
  },
]
</script>

<style scoped>
.common-layout {
  width: 100vw;
  height: 100vh;
}
</style>
