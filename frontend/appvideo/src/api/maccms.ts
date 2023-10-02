import http from '@/utils/http'
import { ApiResult, Data, DataVideo } from './types'

export interface MacCMSRepo {
  id: number
  created_at: string
  update_at: string
  api: string
  name: string
  last_check: string
  available: boolean
}

export enum RequestAction {
  home,
  category,
  detail,
  search
}

export async function getList() {
  const data = (await http.get<ApiResult<MacCMSRepo>>("/maccms")).data
  return data
}

export async function getDetail(mid: number | string, detailID: number | string) {
  const data = (await http.request<ApiResult<DataVideo[]>>({
    method: "post",
    url: `/maccms/proxy/${mid}`,
    data: {
      request_action: RequestAction.detail,
      ids: detailID,
    },
  })).data.data[0]
  return data
}