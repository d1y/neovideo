import http from "@/shared/http";
import { JiexiTable } from "@t/jiexi";
import { ApiResult } from "@t/http";

export async function getList() {
  return (await http.get<ApiResult<JiexiTable[]>>("/maccms")).data;
}

export async function create(data: Partial<JiexiTable>) {
  return (await http.post<ApiResult<JiexiTable>>("/maccms", data)).data;
}

export async function del(id: number) {
  return (await http.delete<ApiResult<number>>(`/maccms/${id}`)).data;
}

export async function batchImport(data: string) {
  return (await http.post<ApiResult<number>>("/maccms/batch_import", { data })).data;
}
