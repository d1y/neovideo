<template>
  <div class="detail">
    <Header />

    <div class="player-wrap">
      <div class="left">
        <div ref="videoRef" v-if="!currentVideo?.embed"></div>
        <iframe v-else :src="currentVideo?.url"></iframe>
      </div>
      <div class="right">
        <h3>{{ $t('Pages.playUrl') }}</h3>
        <div class="p-list">
          <div class="p-item p-flag" :class="{ active: currentPlayFlag == item }" @click="handleChangeCurrentFlag(item)"
            v-for="item in videoFlags">{{ item }}</div>
        </div>
        <div style="margin: 12px; width: 100%;height: 2px; background-color: rgba(0,0,0,.1);"></div>
        <div class="p-list">
          <div class="p-item" :class="{ active: item.url == currentVideo?.url }" v-for="item in currentVideos"
            @click="startPlay(item)">
            {{ item.name }}
          </div>
        </div>
      </div>
    </div>

    <div class="content-wrap">
      <div class="meta-wrap">
        <h1 class="title">{{ vodData?.title }}</h1>
        <div class="info-wrap">
          <img v-lazy="'/public/' + vodData?.cover" />
          <div class="info">
            <div v-if="vodData?.director" class="info-item">{{ $t('Pages.director') }}: {{ vodData?.director }}</div>
            <div v-if="vodData?.actor" class="info-item">{{ $t('Pages.actors') }}: {{ vodData?.actor }}</div>
            <div class="category">{{ vodData?.lang }} / {{ vodData?.area }} /
              {{ vodData?.year }}
            </div>
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
import { VideoInfo, VideoVideo } from '@/api/types'
import Header from '@/views/components/header.vue'

import xgplayer from 'xgplayer'
import HlsPlugin from 'xgplayer-hls'
import 'xgplayer/dist/index.min.css';

const props = defineProps<{ id: string }>()

const videoRef = ref<any>()
let xgplayerInstance: xgplayer

const currentVideo = ref<VideoVideo>()
const currentPlayFlag = ref('')

watch(currentVideo, _ => {
  xgplayerInstance && xgplayerInstance.destroy()
})

const vodData = ref<VideoInfo>()

const videoMap = computed(() => {
  const m = new Map<string, VideoVideo[]>()
  const videos = vodData.value?.videos || []
  if (!videos.length) return m
  videos.forEach(item => {
    m.set(item.flag, item.videos)
  })
  return m
})

const videoFlags = computed(() => {
  if (videoMap.value.size == 0) return []
  return Array.from(videoMap.value.keys())
})

const currentVideos = computed(() => {
  if (videoMap.value.size == 0) return []
  const videos = videoMap.value.get(currentPlayFlag.value)
  return videos
})

onMounted(async () => {
  await getData()
  if (vodData.value && vodData.value.videos.length) {
    const vd = vodData.value.videos[0]
    const onceVideo = vd.videos[0]
    currentPlayFlag.value = vd.flag
    startPlay(onceVideo)
  }
})

const startPlay = async (video: VideoVideo) => {
  currentVideo.value = video
  if (!video.embed) {
    await nextTick()
    if (xgplayerInstance) {
      xgplayerInstance.destroy()
    }
    const $0 = videoRef.value
    xgplayerInstance = new xgplayer({
      el: $0,
      width: '100%',
      height: '420px',
      isLive: false,
      poster: vodData.value!.cover,
      pip: true,
      screenShot: true,
      cssFullscreen: true,
      videoTitle: true,
      showList: true,
      showHistory: true,
      quitMiniMode: true,
      closeVideoTouch: true,
      commonStyle: {
        progressColor: '#fff',
        playedColor: '#06aeec',
      },
      // ignores: ['replay', 'error'],
      plugins: [
        HlsPlugin,
      ],
      url: video.url,
      autoplay: true,
    })
  }
}

const getData = async function () {
  const data = await getDetail(props.id)
  vodData.value = data
}

function handleChangeCurrentFlag(flag: string) {
  if (flag == currentPlayFlag.value) return
  const idx = currentVideos.value!.findIndex(item => item.url == currentVideo.value?.url)
  if (idx <= -1) return
  currentPlayFlag.value = flag
  const curr = currentVideos.value![idx] // FIXME: null check(index out of range)
  startPlay(curr)
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

iframe {
  width: 100%;
  min-height: 42vh;
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

      .p-flag {
        font-size: 18px;
        border-radius: 12px;
        background-color: #333;
        color: #fff;
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
    flex: 1;
    display: flex;
    flex-direction: column;

    .title {
      padding-top: 16px;
      border-top: 1px solid rgba(0, 0, 0, 0.04);
    }

    .category {
      padding-bottom: 16px;
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
