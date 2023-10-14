<template>
  <div>
    <!-- <template v-for="item in data"> -->
    <div class="header">
      <h3>{{ "demo" }}</h3>
      <a class="more">更多</a>
    </div>

    <div class="lvideo-list">
      <a class="video-item" :href="handleDetail(subItem.id)" v-for="subItem in data">
        <div class="cover-wrap">
          <img v-lazy="'/public/' + subItem.cover" />
          <span class="remarks">{{ subItem.real_type }}</span>
        </div>
        <div class="meta-wrap">
          <div class="title">{{ subItem.title }}</div>
          <div class="info">{{ subItem.real_time }}更新</div>
        </div>
      </a>
    </div>
    <!-- </template> -->
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { getVideos } from '@/api/vod'
import { VideoInfo } from '@/api/types'

const data = ref<VideoInfo[]>()
const getData = async function () {
  const vodhome = await getVideos()
  data.value = vodhome.Records
}
onMounted(getData)
function handleDetail(mid: number | string) {
  return `/detail/${mid}`
}
</script>

<style scoped lang="less">
@media screen and (min-width: 1px) and (max-width: 768px) {
  .video-item {
    width: calc((100% - 2 * 16px) / 3) !important;
  }
}

a:hover {
  color: #0c0d0f;
}

.header {
  margin-top: 32px;
  display: flex;
  flex-direction: row;
  align-items: center;
  /* 垂直居中 */
  justify-content: space-between;

  /* 两端对齐 */
  h3 {
    line-height: 30px;
    font-weight: 700;
    text-align: center;
    margin-bottom: 0;
  }

  .more {
    cursor: pointer;
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
        background-color: #e6f2f5;
        width: 100%;
        height: 100%;
        background-size: cover;
        object-fit: cover;
      }

      .remarks {
        position: absolute;
        right: 2px;
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
</style>