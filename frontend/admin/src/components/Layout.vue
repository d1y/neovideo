<template>
  <div class="common-layout">
    <el-container class="w-full h-full">
      <el-aside class="h-full bg-[#191a23] box-border" :width="width" :style="{
        padding: isExpand ? '1rem' : '0',
        transition: 'all .12s',
      }">
        <div class="text-white pl-4">
          <p class="cursor-pointer m-2 mb-4 p-2 flex items-center" :style="{
            backgroundColor: currentPath == item.path ? '#4d70ff' : '',
            borderRadius: `4px`,
          }" v-for="item in menus" @click="$router.push(item.path)">
            <el-icon><component :is="item.icon" /></el-icon>
            <span class="ml-[12px]" v-if="isExpand">{{ item.title }}</span>
          </p>
        </div>
      </el-aside>
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import router from '@/router'
import { onMounted, computed, ref } from "vue"
import { useRoute } from "vue-router"

const route = useRoute()

const isExpand = ref<boolean>(true)

const width = computed(()=> {
  if (isExpand.value) return `240px`
  return `64px`
})

function toggleExpand() {
  isExpand.value = !isExpand.value
}

// bind cmd+b like vscode
function bindCmdB() {
  window.addEventListener('keydown', (e)=> {
    if (e.key === "b" && e.metaKey) {
      toggleExpand()
    }
  })
}

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
    icon: "HomeFilled",
  },
  {
    title: "解析源管理",
    path: "/jiexi",
    icon: "Compass",
  },
  {
    title: "苹果CMS管理",
    path: "/maccms",
    icon: "Grid",
  },
  {
    title: "爬虫任务管理",
    path: "/spider",
    icon: "Odometer",
  },
]
</script>

<style scoped>
.common-layout {
  width: 100vw;
  height: 100vh;
}
</style>
