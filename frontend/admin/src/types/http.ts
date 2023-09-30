export interface ApiResult<T> {
  message: string
  data: T
  success: boolean
}