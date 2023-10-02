import http from '@/utils/http'
import { ApiResult } from "./types"

export type VodHome = ApiResult<VodItem[]>

export interface VodItem {
  id: number
  api: string
  name: string
  data: Data
}

export interface Data {
  list_header: ListHeader
  category: Category[]
  videos: DataVideo[]
}

export interface Category {
  text: string
  id: number
}

export interface ListHeader {
  page: number
  page_count: number
  page_size: number
  record_count: number
}

export interface DataVideo {
  last: Date
  id: number
  tid: number
  name: string
  type: 'XML' | 'JSON'
  dt: string
  note: string
  desc: string
  lang: string
  area: string
  year: string
  state: string
  actor: string
  director: string
  pic: string
  dd: DD[] | null
}

export interface DD {
  flag: string
  videos: DDVideo[]
}

export interface DDVideo {
  name: string
  url: string
}

export async function getHome(): Promise<VodItem[]> {
  const data = ((await http.get<VodHome>("/vod/home")).data).data
  return data
}
