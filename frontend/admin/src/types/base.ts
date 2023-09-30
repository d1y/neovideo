export interface DbBaseModel {
  create_at: string
  update_at: string
  id: number
}

export type NullableDbBaseModel = Partial<DbBaseModel>