import http from '@/utils/http'
import { ApiResult, IPaginationResult, VideoInfo } from "./types"

export async function getVideos(page = 1, limit = 20): Promise<IPaginationResult<VideoInfo>> {
  const data = (await http.request<ApiResult<IPaginationResult<VideoInfo>>>({
    url: "/vod/videos",
    method: "get",
    params: { page, limit }
  })).data
  return data.data
}

export async function getDetail(id: number | string): Promise<VideoInfo> {
  const data = (await http.get<ApiResult<VideoInfo>>(`/vod/video/${id}`)).data
  return data.data
}