export interface DbBaseModel {
  UpdatedAt: string
  CreatedAt: string
  ID: number
}

export type NullableDbBaseModel = Partial<DbBaseModel>