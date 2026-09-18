
import { Context } from '../types'


function prepareParams(ctx: Context) {
  const utility = ctx.utility
  const findparam = utility.param


  const target = ctx.target

  let params = target.args.params

  params = params || []

  let out: any = {}
  for (let pd of params) {
    let val = findparam(ctx, pd)
    if (null != val) {
      out[pd.name] = val
    }
  }


  return out
}


export {
  prepareParams
}
