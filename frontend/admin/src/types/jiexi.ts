import { NullableDbBaseModel } from "./base"

export interface JiexiTable extends Readonly<NullableDbBaseModel> {
  name: string
  url: string
  note: string
}