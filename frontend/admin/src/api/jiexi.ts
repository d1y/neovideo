import http from '@/shared/http'
import { JiexiTable } from '@t/jiexi'
import { ApiResult } from '@t/http'

export async function getList() {
  return (await http.get<ApiResult<JiexiTable[]>>("/jiexi")).data
}

export async function create(data: Partial<JiexiTable>) {
  return (await http.post<ApiResult<JiexiTable>>("/jiexi", data)).data
}

export async function del(id: number){
  return (await http.delete<ApiResult<number>>(`/jiexi/${id}`)).data
}

export async function batchImport(data: string) {
  return (await http.post<ApiResult<number>>("/jiexi/batch_import", {data})).data
}