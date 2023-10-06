<template>
  <div class="detail">
    <Header />

    <div class="player-wrap">
      <div class="left">
        <video id="player-box"></video>
      </div>
      <div class="right">
        <h3>{{ $t('Pages.playUrl') }}</h3>
        <template v-for="item in (vodData?.videos || [])">
          <div class="p-list">
            <div class="p-item" :class="{ active: currentPlayFlag == item.flag }">{{ item.flag }}</div>
          </div>
          <div style="margin: 12px; width: 100%;height: 2px; background-color: rgba(0,0,0,.1);"></div>
          <div class="p-list">
            <div class="p-item" :class="{ active: i.url == currentLink }" v-for="i in item.videos" @click="startPlay(i.url)">
              {{ i.name }}
            </div>
          </div>
        </template>
      </div>
    </div>

    <div class="content-wrap">
      <div class="meta-wrap">
        <h1 class="title">{{ vodData?.title }}</h1>
        <div class="category"
          >{{ vodData?.lang }} / {{ vodData?.area }} /
          {{ vodData?.year }}
        </div>
        <div class="info-wrap">
          <img v-lazy="vodData?.real_cover" />
          <div class="info">
            <div v-if="vodData?.director" class="info-item">{{ $t('Pages.director') }}: {{ vodData?.director }}</div>
            <div v-if="vodData?.actor" class="info-item">{{ $t('Pages.actors') }}: {{ vodData?.actor }}</div>
          </div>
        </div>
      </div>
      <h3 style="height: 40px; border-bottom: 1px solid rgba(0, 0, 0, 0.04)">内容详情</h3>
      <div v-if="vodData?.desc" v-html="vodData?.desc"></div>
      <div v-else>暂无内容简介</div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { getDetail } from '@/api/vod'
import { VideoInfo } from '@/api/types'
import Header from '@/views/components/header.vue'

import TCPlayer from 'tcplayer.js'
import 'tcplayer.js/dist/tcplayer.min.css'

const props = defineProps<{ id: string }>()

const currentLink = ref('')
const currentPlayFlag = ref('')

const vodData = ref<VideoInfo>()

let player: any = undefined
onMounted(async () => {
  player = TCPlayer('player-box', {})
  await getData()
  if (vodData.value && vodData.value.videos.length) {
    const vd = vodData.value.videos[0]
    const link = vd.videos[0].url
    currentPlayFlag.value = vd.flag
    startPlay(link)
  }
})

const startPlay = (link) => {
  currentLink.value = link
  player.src(currentLink.value)
}

const getData = async function () {
  const data = await getDetail(props.id)
  vodData.value = data
}
</script>
<style scoped lang="less">
@media screen and (min-width: 1px) and (max-width: 768px) {
  // 播放器相关
  #player-box {
    height: 220px !important;
  }

  .video-item {
    width: calc((100% - 2 * 16px) / 3) !important;
  }
}

@media screen and (min-width: 768px) {
  // 播放器相关
  #player-box {
    height: 480px !important;
  }

  // 分集
  .p-item {
    width: calc((100% - 3 * 10px) / 4) !important;
  }
}

dl,
ol,
ul {
  margin-bottom: 0px;
  margin-top: 0;
}

.detail {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.player-wrap {
  margin-top: 70px;
  width: 100%;
  max-width: 1024px;
  display: flex;
  flex-direction: column;
  gap: 16px;

  .left {
    #player-box {
      width: 100%;
      background-color: #000;
    }
  }

  .right {
    width: 100%;
    height: auto;
    padding: 16px;
    border-radius: 4px;
    background-color: rgba(0, 0, 0, 0.04);

    .p-list {
      max-height: calc(100% - 25px - 32px);
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
      overflow-y: auto;

      .p-item {
        cursor: pointer;
        align-items: center;
        justify-content: center;
        width: calc((100% - 2 * 10px) / 3);
        padding: 8px 0;
        border-radius: 2px;
        color: #0c0d0f;
        background-color: rgba(0, 0, 0, 0.04);
        font-weight: inherit;
        font-family: inherit;
        text-align: center;
        text-overflow: ellipsis;
        overflow: hidden;
        white-space: nowrap;
      }

      .active {
        color: #fe3355;
      }
    }
  }
}

.content-wrap {
  padding-left: 12px;
  padding-right: 12px;
  width: 100%;
  max-width: 1024px;
  display: flex;
  flex-direction: column;
  gap: 16px;

  .meta-wrap {
    //padding-left: 12px;
    //padding-right: 12px;
    flex: 1;
    display: flex;
    flex-direction: column;

    .title {
      padding-top: 16px;
      border-top: 1px solid rgba(0, 0, 0, 0.04);
    }

    .category {
      padding-bottom: 16px;
      border-bottom: 1px solid rgba(0, 0, 0, 0.04);
      color: #606266;
      font-size: 14px;
    }

    .info-wrap {
      margin-top: 16px;
      display: flex;
      flex-direction: row;
      gap: 16px;

      img {
        width: 100px;
        height: 120px;
        background-size: cover;
        object-fit: cover;
      }

      .info {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 12px;
        color: #606266;
        font-size: 14px;

        .line {
          width: 100%;
          text-overflow: ellipsis;
          overflow: hidden;
          white-space: nowrap;
        }
      }
    }
  }

  .lvideo-list {
    min-height: 200px;
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
}
</style>
