
import { Result, Context } from '../types'


import { clean } from './CleanUtility'

import { clone, delprop } from './StructUtility'


function makeError(ctx: Context, err?: any) {

  ctx = ctx || {}
  const op = ctx.op || {}
  op.name = op.name || 'unknown operation'


  const result = ctx.result || new Result({})
  result.ok = false

  const reserr = result.err

  err = undefined === err ? reserr : err
  err = err || ctx.error('unknown', 'unknown error')

  const errmsg = err.message || 'unknown error'
  const msg = 'UniversalSDK: ' + op.name + ': ' + errmsg
  err.message = clean(ctx, msg)

  if (result.err) {
    delprop(result, 'err')
  }

  const spec = ctx.spec || {}

  if (ctx.ctrl.explain) {
    ctx.ctrl.explain.err = {
      ...clone({ err }).err,
      message: err.message,
      stack: err.stack,
    }
  }

  err.result = clean(ctx, result)
  err.spec = clean(ctx, spec)

  ctx.ctrl.err = err

  if (false === ctx.ctrl.throw) {
    return result.resdata
  }
  else {
    throw err
  }
}


export {
  makeError
}
