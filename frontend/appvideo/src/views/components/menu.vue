<template>
  <div class="menu-wrap">
    <div class="category__list">
      <div class="category__list__item" :class="currentMenu === item.id ? 'active' : ''" v-for="item in menuData" @click="hdlClick(item.id)">{{ item.label }} </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import useVods from '@/store/modules/useVods'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const { setCurrentCategory } = useVods()
const { menus, currentCategory: currentMenu }= storeToRefs(useVods())

const menuData = computed<{label: string, id: number}[]>(()=> {
  return [
  { label: '首页', id: -1 },
    ...menus.value,
  ]
})

function hdlClick(id: number | number) {
  setCurrentCategory(id)
  if (id == -1) {
    router.replace({
      path: '/index',
    })
    return
  }
  router.push({
    path: '/dianying',
    query: {
      id,
    }
  })
}

onMounted(()=> {
  const p = route.path
  const q = route.query.id as string
  if (p == "/index") {
    setCurrentCategory(-1)
  } else {
    q && setCurrentCategory(+q)
  }
})
</script>

<style scoped lang="less">
.menu-wrap {
  margin-top: 64px;
  padding: 20px 0 0;

  .category__list {
    width: 100%;
    margin-bottom: 39px;
    font-size: 0;
    position: relative;
    display: flex;
    flex-direction: row;
    gap: 20px;
    overflow-x: auto;
    overflow-y: hidden;

    .category__list__item {
      display: inline-block;
      cursor: pointer;
      font-size: 16px;
      line-height: 22px;
      color: #606266;
      text-align: center;
      //text-overflow: ellipsis;
      //overflow: hidden;
      white-space: nowrap;
    }

    .active {
      color: #0c0d0f;
      font-weight: 500;
      position: relative;
    }

    .active:after {
      position: absolute;
      top: 34px;
      content: ' ';
      width: 20px;
      left: 50%;
      transform: translateX(-50%);
      height: 3px;
      background-color: #0c0d0f;
    }

    li {
      list-style: none;
    }
  }

  .category__list:after {
    content: ' ';
    display: block;
    height: 1px;
    background-color: rgba(12, 13, 15, 0.06);
    position: absolute;
    top: 37px;
    left: 0;
    right: 0;
  }
}
</style>