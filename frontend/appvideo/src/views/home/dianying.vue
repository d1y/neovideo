<template>
  <div class="layout-content">
    <div class="category-layout">
      <div class="category__line">
        <li class="category__first-ele">分类:</li>
        <ul class="category__list category__sub">
          <li class="category__list__item category__sub__item" :class="{
            active: currentCategory?.id == item.id,
          }" @click="hdlClickCategory(item)" v-for="item in category" >
          {{ item.text }}
          </li>
        </ul>
      </div>
    </div>

    <div class="lvideo-list">
      <div v-if="!currentVideos.length">
        <h1>暂无数据 :(</h1>
      </div>
      <a v-else class="video-item"  :href="handleDetail(subItem.id, id)" v-for="subItem in currentVideos">
        <div class="cover-wrap">
          <img v-lazy="subItem.pic" />
          <span class="remarks">{{ subItem.desc }}</span>
        </div>
        <div class="meta-wrap">
          <div class="title">{{ subItem.name }}</div>
          <div class="info">{{ subItem.last }}更新</div>
        </div>
      </a>
    </div>

    <div class="page-wrap" v-if="show">
      <div class="page-item" :class="{
        disable: !isPrev,
      }" @click="hdlPageChange(false)">上页</div>
      <div>{{ text }}</div>
      <div class="page-item" :class="{
        disable: !isNext,
      }" @click="hdlPageChange(true)">下页</div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { Category, Data, DataVideo } from '@/api/types'
import useVods from '@/store/modules/useVods'
import { watch } from 'vue'
import { useRoute } from 'vue-router'
import usePagination from '@/composition/usePagination'
import * as maccmsApi from '@/api/maccms'

const props = defineProps<{
  id: string | number
}>()

const currentPage = ref<number>(1)

const { isPrev, isNext, show, text } = usePagination(currentPage, computed(()=> currentData.value?.list_header))

const { getCategoryByID, loadVodHomeDataWithApi } = useVods()
const { category: categoryCol } = storeToRefs(useVods())

const route = useRoute()
watch(()=> route.query, p=> {
  currentCategory.value = null
  const id = +(p.id as string)
  init(id)
})

const currentCategory = ref<Category | null>(null)
const currentData = ref<Data>()

const category = computed(()=> {
  const val = categoryCol.value.get(+props.id)
  return val
})

const currentVideos = computed<DataVideo[]>(()=> {
  if (!currentData.value) return []
  return currentData.value.videos || []
})

const handleDetail = (vod_id: number, mid: number | string) => {
  return `/detail/${vod_id}?mid=${mid}`
}

async function getData() {
  const data = await maccmsApi.getHomeWithPageAndCategory(+props.id, currentPage.value, currentCategory.value?.id || -1)
  currentData.value = data
}

async function init(id?: number) {
  id = id ? id : +props.id
  const val = getCategoryByID(id)
  currentCategory.value = val.length ? val[0] : null
  getData()
}

async function hdlClickCategory(item: any) {
  currentPage.value = 1
  currentCategory.value = item
  await getData()
}

async function hdlPageChange(next: boolean) {
  if (next) currentPage.value++
  else currentPage.value--
  await getData()
}

onMounted(async ()=> {
  await loadVodHomeDataWithApi()
  await init()
})

</script>

<style scoped lang="less">
@media screen and (min-width: 1px) and (max-width: 768px) {
  .video-item {
    width: calc((100% - 2 * 16px) / 3) !important;
  }

  .category__list {
    white-space: nowrap !important;
    overflow-x: auto !important;
  }
}

.layout-content {
  width: 100%;
  padding: 1px 0px;

  .category-layout {
    .category__line {
      display: flex;
      font-size: 14px;

      .category__first-ele {
        line-height: 32px;
        height: 32px;
        flex-shrink: 0;
        padding-right: 12px;
        border-radius: 2px;
      }

      li {
        list-style: none;
      }

      .category__sub {
        line-height: normal;
        margin-bottom: 12px;
      }

      .category__list {
        font-size: 0;
        position: relative;

        .category__list__item {
          display: inline-block;
          cursor: pointer;
          font-size: 16px;
          line-height: 22px;
          color: #606266;
        }

        .category__sub__item {
          font-size: 14px;
          line-height: 32px;
          margin-bottom: 2px;
          color: #0c0d0f;
          padding: 0 8px;
          border-radius: 2px;
        }

        .active {
          background-color: rgba(0, 0, 0, 0.04);
          color: #fe3355;
          font-weight: 500;
        }
      }

      .category__list::-webkit-scrollbar {
        display: none;
      }
    }
  }

  .lvideo-list {
    min-height: 200px;
    margin-top: 12px;
    display: flex;
    flex-wrap: wrap;
    gap: 16px;

    .video-item {
      width: calc((100% - 3 * 16px) / 4);
      aspect-ratio: 3/5;
      min-height: 120px;

      .cover-wrap {
        position: relative;
        width: 100%;
        height: 85%;

        img {
          border-radius: 4px;
          overflow: hidden;
          // todo 修改默认图
          background: url(../images/load.gif) no-repeat;
          // background-color: #e6f2f5;
          width: 100%;
          height: 100%;
          background-size: cover;
          object-fit: cover;
        }

        .remarks {
          position: absolute;
          right: 4px;
          bottom: 1px;
          color: #fff;
          font-size: 12px;
        }
      }

      .meta-wrap {
        .title {
          text-align: center;
          text-overflow: ellipsis;
          overflow: hidden;
          white-space: nowrap;
        }

        .info {
          display: none;
        }
      }
    }
  }

  .page-wrap {
    margin-top: 16px;
    display: flex;
    flex-direction: row;
    gap: 16px;
    justify-content: center; /* 水平居中 */
    align-items: center; /* 垂直居中 */

    .page-item {
      user-select: none;
      cursor: pointer;
      padding: 0 16px;
      height: 30px;
      line-height: 30px;
      border-radius: 4px;
      background-color: #eee;
      color: #000;
      font-size: 12px;
      text-align: center;
    }

    .disable {
      color: grey;
    }
  }
}
</style>
