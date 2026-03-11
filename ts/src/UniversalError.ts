
import { Context } from './Context'


class UniversalError extends Error {

  isUniversalError = true

  sdk = 'Universal'

  code: string
  ctx: Context

  constructor(code: string, msg: string, ctx: Context) {
    super(msg)
    this.code = code
    this.ctx = ctx
  }

}

export {
  UniversalError
}

