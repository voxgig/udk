
import { Context } from '../types'

function transformRequest(ctx: Context) {
  const spec = ctx.spec
  const utility = ctx.utility
  const target = ctx.target
  const isfunc = utility.struct.isfunc
  const transform = utility.struct.transform

  if (spec) {
    spec.step = 'reqform'
  }

  try {
    const reqform = target.transform.req
    const reqdata = isfunc(reqform) ? reqform(ctx) : transform({
      reqdata: ctx.reqdata
    }, reqform)

    return reqdata
  }
  catch (err) {
    return utility.makeError(ctx, err)
  }
}


export {
  transformRequest
}
