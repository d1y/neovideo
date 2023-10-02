import http from '@/utils/http'
import { VodHome, VodItem } from "./types"

export async function getHome(): Promise<VodItem[]> {
  const data = ((await http.get<VodHome>("/vod/home")).data).data
  return data
}
