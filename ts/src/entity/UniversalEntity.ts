
import { inspect } from 'node:util'

import {
  UniversalSDK,
} from '../UniversalSDK'

import {
  Utility
} from '../utility/Utility'

import type {
  Operation,
  Context,
  Control,
} from '../types'


// TODO: needs Entity superclass
class UniversalEntity {
  name: string

  #client: UniversalSDK
  #utility: Utility
  #entopts: any
  #data: any
  #match: any

  _entctx: Context

  constructor(client: UniversalSDK, name: string, entopts: any) {
    // super()

    this.name = name

    entopts = entopts || {}
    entopts.active = false !== entopts.active

    this.#client = client
    this.#entopts = entopts
    this.#utility = client.utility()
    this.#data = {}
    this.#match = {}

    const makeContext = this.#utility.makeContext

    this._entctx = makeContext({
      entity: this,
      entopts,
    }, client._rootctx)

    const featureHook = this.#utility.featureHook
    featureHook(this._entctx, 'PostConstructEntity')
  }

  entopts() {
    return { ...this.#entopts }
  }

  client() {
    return this.#client
  }

  make() {
    return new UniversalEntity(this.#client, this.name, this.entopts())
  }


  data(this: any, data?: any) {
    const struct = this.#utility.struct
    const featureHook = this.#utility.featureHook

    if (null != data) {
      this.#data = struct.clone(data)
      featureHook(this._entctx, 'SetData')
    }

    featureHook(this._entctx, 'GetData')
    let out = struct.clone(this.#data)

    return out
  }


  match(match?: any) {
    const struct = this.#utility.struct
    const featureHook = this.#utility.featureHook

    if (null != match) {
      this.#match = struct.clone(match)
      featureHook(this._entctx, 'SetMatch')
    }

    featureHook(this._entctx, 'GetMatch')
    let out = struct.clone(this.#match)

    return out
  }


  toJSON() {
    const struct = this.#utility.struct
    return struct.merge([{}, struct.getdef(this.#data, {}), { $entity: 'Universal' }])
  }

  toString() {
    return 'UniversalEntity: ' + this.name + ': ' + this.#utility.struct.jsonify(this.#data)
  }

  [inspect.custom]() {
    return this.toString()
  }



  async load(this: any, reqmatch?: any, ctrl?: Control) {

    const utility = this.#utility

    const {
      makeContext,
      done,
      error,
      featureHook,
      makeTarget,
      makeRequest,
      makeResponse,
      makeResult,
      makeSpec,
    } = utility

    let fres: Promise<any> | undefined = undefined

    let ctx: Context = makeContext({
      opname: 'load',
      ctrl,
      match: this.#match,
      data: this.#data,
      reqmatch
    }, this._entctx)

    try {

      fres = featureHook(ctx, 'PreSelection')
      if (fres instanceof Promise) { await fres }

      ctx.out.target = makeTarget(ctx)
      if (ctx.out.target instanceof Error) {
        return error(ctx, ctx.out.target)
      }



      fres = featureHook(ctx, 'PreSpec')
      if (fres instanceof Promise) { await fres }

      ctx.out.spec = makeSpec(ctx)
      if (ctx.out.spec instanceof Error) {
        return error(ctx, ctx.out.spec)
      }



      fres = featureHook(ctx, 'PreRequest')
      if (fres instanceof Promise) { await fres }

      ctx.out.request = await makeRequest(ctx)
      if (ctx.out.request instanceof Error) {
        return error(ctx, ctx.out.request)
      }



      fres = featureHook(ctx, 'PreResponse')
      if (fres instanceof Promise) { await fres }

      ctx.out.response = await makeResponse(ctx)
      if (ctx.out.response instanceof Error) {
        return error(ctx, ctx.out.response)
      }



      fres = featureHook(ctx, 'PreResult')
      if (fres instanceof Promise) { await fres }

      ctx.out.result = await makeResult(ctx)
      if (ctx.out.result instanceof Error) {
        return error(ctx, ctx.out.result)
      }



      fres = featureHook(ctx, 'PreDone')
      if (fres instanceof Promise) { await fres }

      if (null != ctx.result) {
        if (null != ctx.result.resmatch) {
          this.#match = ctx.result.resmatch
        }

        if (null != ctx.result.resdata) {
          this.#data = ctx.result.resdata
        }
      }

      return done(ctx)
    }
    catch (err: any) {

      fres = featureHook(ctx, 'PreUnexpected')
      if (fres instanceof Promise) { await fres }

      err = this.#unexpected(ctx, err)

      if (err) {
        throw err
      }
      else {
        return undefined
      }
    }
  }



  async list(this: any, reqmatch?: any, ctrl?: Control) {

    const utility = this.#utility

    const {
      makeContext,
      done,
      error,
      featureHook,
      makeTarget,
      makeRequest,
      makeResponse,
      makeResult,
      makeSpec,
    } = utility

    let fres: Promise<any> | undefined = undefined

    let ctx: Context = makeContext({
      opname: 'list',
      ctrl,
      match: this.#match,
      data: this.#data,
      reqmatch
    }, this._entctx)

    try {

      fres = featureHook(ctx, 'PreSelection')
      if (fres instanceof Promise) { await fres }

      ctx.out.target = makeTarget(ctx)
      if (ctx.out.target instanceof Error) {
        return error(ctx, ctx.out.target)
      }



      fres = featureHook(ctx, 'PreSpec')
      if (fres instanceof Promise) { await fres }

      ctx.out.spec = makeSpec(ctx)
      if (ctx.out.spec instanceof Error) {
        return error(ctx, ctx.out.spec)
      }



      fres = featureHook(ctx, 'PreRequest')
      if (fres instanceof Promise) { await fres }

      ctx.out.request = await makeRequest(ctx)
      if (ctx.out.request instanceof Error) {
        return error(ctx, ctx.out.request)
      }



      fres = featureHook(ctx, 'PreResponse')
      if (fres instanceof Promise) { await fres }

      ctx.out.response = await makeResponse(ctx)
      if (ctx.out.response instanceof Error) {
        return error(ctx, ctx.out.response)
      }



      fres = featureHook(ctx, 'PreResult')
      if (fres instanceof Promise) { await fres }

      ctx.out.result = await makeResult(ctx)
      if (ctx.out.result instanceof Error) {
        return error(ctx, ctx.out.result)
      }



      fres = featureHook(ctx, 'PreDone')
      if (fres instanceof Promise) { await fres }

      if (null != ctx.result) {
        if (null != ctx.result.resmatch) {
          this.#match = ctx.result.resmatch
        }
      }

      return done(ctx)
    }
    catch (err: any) {

      fres = featureHook(ctx, 'PreUnexpected')
      if (fres instanceof Promise) { await fres }

      err = this.#unexpected(ctx, err)

      if (err) {
        throw err
      }
      else {
        return undefined
      }
    }
  }



  async create(this: any, reqdata?: any, ctrl?: Control) {

    const utility = this.#utility
    const {
      makeContext,
      done,
      error,
      featureHook,
      makeTarget,
      makeRequest,
      makeResponse,
      makeResult,
      makeSpec,
    } = utility

    let fres: Promise<any> | undefined = undefined

    let ctx: Context = makeContext({
      opname: 'create',
      ctrl,
      match: this.#match,
      data: this.#data,
      reqdata
    }, this._entctx)

    try {

      fres = featureHook(ctx, 'PreSelection')
      if (fres instanceof Promise) { await fres }

      ctx.out.target = makeTarget(ctx)
      if (ctx.out.target instanceof Error) {
        return error(ctx, ctx.out.target)
      }



      fres = featureHook(ctx, 'PreSpec')
      if (fres instanceof Promise) { await fres }

      ctx.out.spec = makeSpec(ctx)
      if (ctx.out.spec instanceof Error) {
        return error(ctx, ctx.out.spec)
      }



      fres = featureHook(ctx, 'PreRequest')
      if (fres instanceof Promise) { await fres }

      ctx.out.request = await makeRequest(ctx)
      if (ctx.out.request instanceof Error) {
        return error(ctx, ctx.out.request)
      }



      fres = featureHook(ctx, 'PreResponse')
      if (fres instanceof Promise) { await fres }

      ctx.out.response = await makeResponse(ctx)
      if (ctx.out.response instanceof Error) {
        return error(ctx, ctx.out.response)
      }



      fres = featureHook(ctx, 'PreResult')
      if (fres instanceof Promise) { await fres }

      ctx.out.result = await makeResult(ctx)
      if (ctx.out.result instanceof Error) {
        return error(ctx, ctx.out.result)
      }



      fres = featureHook(ctx, 'PreDone')
      if (fres instanceof Promise) { await fres }

      if (null != ctx.result) {
        if (null != ctx.result.resdata) {
          this.#data = ctx.result.resdata
        }
      }

      return done(ctx)
    }
    catch (err: any) {

      fres = featureHook(ctx, 'PreUnexpected')
      if (fres instanceof Promise) { await fres }

      err = this.#unexpected(ctx, err)

      if (err) {
        throw err
      }
      else {
        return undefined
      }
    }
  }



  async update(this: any, reqdata?: any, ctrl?: Control) {

    const utility = this.#utility

    const {
      makeContext,
      done,
      error,
      featureHook,
      makeTarget,
      makeRequest,
      makeResponse,
      makeResult,
      makeSpec,
    } = utility

    let fres: Promise<any> | undefined = undefined

    let ctx: Context = makeContext({
      opname: 'update',
      ctrl,
      match: this.#match,
      data: this.#data,
      reqdata
    }, this._entctx)

    try {


      fres = featureHook(ctx, 'PreSelection')
      if (fres instanceof Promise) { await fres }

      ctx.out.target = makeTarget(ctx)
      if (ctx.out.target instanceof Error) {
        return error(ctx, ctx.out.target)
      }



      fres = featureHook(ctx, 'PreSpec')
      if (fres instanceof Promise) { await fres }

      ctx.out.spec = makeSpec(ctx)
      if (ctx.out.spec instanceof Error) {
        return error(ctx, ctx.out.spec)
      }



      fres = featureHook(ctx, 'PreRequest')
      if (fres instanceof Promise) { await fres }

      ctx.out.request = await makeRequest(ctx)
      if (ctx.out.request instanceof Error) {
        return error(ctx, ctx.out.request)
      }



      fres = featureHook(ctx, 'PreResponse')
      if (fres instanceof Promise) { await fres }

      ctx.out.response = await makeResponse(ctx)
      if (ctx.out.response instanceof Error) {
        return error(ctx, ctx.out.response)
      }



      fres = featureHook(ctx, 'PreResult')
      if (fres instanceof Promise) { await fres }

      ctx.out.result = await makeResult(ctx)
      if (ctx.out.result instanceof Error) {
        return error(ctx, ctx.out.result)
      }



      fres = featureHook(ctx, 'PreDone')
      if (fres instanceof Promise) { await fres }

      if (null != ctx.result) {
        if (null != ctx.result.resmatch) {
          this.#match = ctx.result.resmatch
        }

        if (null != ctx.result.resdata) {
          this.#data = ctx.result.resdata
        }
      }

      return done(ctx)
    }
    catch (err: any) {

      fres = featureHook(ctx, 'PreUnexpected')
      if (fres instanceof Promise) { await fres }

      err = this.#unexpected(ctx, err)

      if (err) {
        throw err
      }
      else {
        return undefined
      }
    }
  }



  async remove(this: any, reqmatch?: any, ctrl?: Control) {

    const utility = this.#utility

    const {
      makeContext,
      done,
      error,
      featureHook,
      makeTarget,
      makeRequest,
      makeResponse,
      makeResult,
      makeSpec,
    } = utility

    let fres: Promise<any> | undefined = undefined

    let ctx: Context = makeContext({
      opname: 'remove',
      ctrl,
      match: this.#match,
      data: this.#data,
      reqmatch
    }, this._entctx)

    try {


      fres = featureHook(ctx, 'PreTarget')
      if (fres instanceof Promise) { await fres }

      ctx.out.target = makeTarget(ctx)
      if (ctx.out.target instanceof Error) {
        return error(ctx, ctx.out.target)
      }



      fres = featureHook(ctx, 'PreSpec')
      if (fres instanceof Promise) { await fres }

      ctx.out.spec = makeSpec(ctx)
      if (ctx.out.spec instanceof Error) {
        return error(ctx, ctx.out.spec)
      }



      fres = featureHook(ctx, 'PreRequest')
      if (fres instanceof Promise) { await fres }

      ctx.out.request = await makeRequest(ctx)
      if (ctx.out.request instanceof Error) {
        return error(ctx, ctx.out.request)
      }



      fres = featureHook(ctx, 'PreResponse')
      if (fres instanceof Promise) { await fres }

      ctx.out.response = await makeResponse(ctx)
      if (ctx.out.response instanceof Error) {
        return error(ctx, ctx.out.response)
      }



      fres = featureHook(ctx, 'PreResult')
      if (fres instanceof Promise) { await fres }

      ctx.out.result = await makeResult(ctx)
      if (ctx.out.result instanceof Error) {
        return error(ctx, ctx.out.result)
      }



      fres = featureHook(ctx, 'PreDone')
      if (fres instanceof Promise) { await fres }

      if (null != ctx.result) {
        if (null != ctx.result.resmatch) {
          this.#match = ctx.result.resmatch
        }

        if (null != ctx.result.resdata) {
          this.#data = ctx.result.resdata
        }
      }

      return done(ctx)
    }
    catch (err: any) {

      fres = featureHook(ctx, 'PreSelection')
      if (fres instanceof Promise) { await fres }

      err = this.#unexpected(ctx, err)

      if (err) {
        throw err
      }
      else {
        return undefined
      }
    }
  }



  #unexpected(this: any, ctx: Context, err: any) {
    const clean = this.#utility.clean
    const struct = this.#utility.struct

    const delprop = struct.delprop
    const clone = struct.clone
    const merge = struct.merge

    const ctrl = ctx.ctrl

    ctrl.err = err

    if (ctrl.explain) {
      ctx.ctrl.explain = clean(ctx, ctx.ctrl.explain)
      delprop(ctx.ctrl.explain.result, 'err')

      if (null != ctx.result && null != ctx.result.err) {
        ctrl.explain.err = clean(ctx, merge([
          clone({ err: ctx.result.err }).err,
          {
            message: ctx.result.err.message,
            stack: ctx.result.err.stack,
          }]))
      }

      const cleanerr = clean(ctx, merge([
        clone({ err }).err,
        {
          message: err.message,
          stack: err.stack,
        }]))

      if (null == ctrl.explain.err) {
        ctrl.explain.err = cleanerr
      }
      else if (ctrl.explain.err.message != cleanerr.message) {
        ctrl.explain.unexpected = cleanerr
      }
    }

    if (false === ctrl.throw) {
      return undefined
    }

    return err
  }

}


export {
  UniversalEntity
}
