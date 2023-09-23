import http from '@/shared/http'

export async function getList() {
  return (await http.get<ApiResult<JiexiTable[]>>("/jiexi")).data
}

export async function create(data: Partial<JiexiTable>) {
  return (await http.post<ApiResult<JiexiTable>>("/jiexi", data)).data
}

export async function del(id: int){
  return (await http.delete<ApiResult<int>>(`/jiexi/${id}`)).data
}

export async function batchImport(data: str) {
  return (await http.post<ApiResult<int>>("/jiexi/batch_import", {data})).data
}