import { NullableDbBaseModel } from "./base"

export interface MacCMSRepo extends Readonly<NullableDbBaseModel> {
  api: string
  name: string
  last_check: string
  available: boolean
}