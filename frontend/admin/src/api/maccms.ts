import http from "@/shared/http"
import { JiexiTable } from "@t/jiexi"
import { ApiResult } from "@t/http"
import { MacCMSRepo } from "@/types/maccms"

export async function getList() {
  return (await http.get<ApiResult<MacCMSRepo[]>>("/maccms")).data
}

export async function create(data: Partial<JiexiTable>) {
  return (await http.post<ApiResult<MacCMSRepo>>("/maccms", data)).data
}

export async function del(id: number) {
  return (await http.delete<ApiResult<number>>(`/maccms/${id}`)).data
}

export async function batchImport(data: string) {
  return (await http.post<ApiResult<number>>("/maccms/batch_import", { data })).data
}

export async function allcheckAndSync() {
  return (await http.post<ApiResult<{
    id: number
    error: string
  }[]>>("/maccms/allcheck/sync")).data.data
}

export async function removeUnavailable() {
  return (await http.delete<ApiResult<MacCMSRepo[]>>("/maccms/allcheck/unavailable")).data.data
}

export async function checkOnceStatus(id: number) {
  return (await http.post<ApiResult<{
    successful: boolean
    message: string
  }>>(`/maccms/check/${id}`)).data.data
}