interface ApiResult<T> {
  message: str
  data: T
  success: bool
}