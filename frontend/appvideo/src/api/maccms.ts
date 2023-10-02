import http from '@/utils/http'
import { ApiResult } from './types'

export interface MacCMSRepo {
  id: number
  created_at: string
  update_at: string
  api: string
  name: string
  last_check: string
  available: boolean
}

export async function getList() {
  const data = (await http.get<ApiResult<MacCMSRepo>>("/maccms")).data
  return data
}