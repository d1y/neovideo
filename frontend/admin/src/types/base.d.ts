// golang like
type nil = null
type int = number
type str = string
type bool = boolean
// ===========

interface DbBaseModel {
  create_at: str
  update_at: str
  id: int
}

type NullableDbBaseModel = Partial<DbBaseModel>