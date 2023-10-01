export interface DbBaseModel {
  update_at: string
  create_at: string
  id: number
}

export type NullableDbBaseModel = Partial<DbBaseModel>