

import * as Fs from 'node:fs'


import { UniversalSDK } from './UniversalSDK'

import { Target } from './Target'
import { Context } from './Context'
import { Control } from './Control'
import { Operation } from './Operation'
import { Response } from './Response'
import { Result } from './Result'
import { Spec } from './Spec'


type FSType = typeof Fs


type UniversalOptions = {
  fs: FSType,
  registry: string,
}



type FeatureOptions = Record<string, any> | {
  active: boolean
}


interface Feature {
  version: string
  name: string
  active: boolean

  init: (ctx: Context, options: FeatureOptions) => void | Promise<any>

  PostConstruct: (this: UniversalSDK, ctx: Context) => void | Promise<any>
  PostConstructEntity: (this: UniversalSDK, ctx: Context) => void | Promise<any>
  SetData: (this: UniversalSDK, ctx: Context) => void | Promise<any>
  GetData: (this: UniversalSDK, ctx: Context) => void | Promise<any>
  GetMatch: (this: UniversalSDK, ctx: Context) => void | Promise<any>

  PreTarget: (this: UniversalSDK, ctx: Context) => void | Promise<any>
  PreSpec: (this: UniversalSDK, ctx: Context) => void | Promise<any>
  PreRequest: (this: UniversalSDK, ctx: Context) => void | Promise<any>
  PreResponse: (this: UniversalSDK, ctx: Context) => void | Promise<any>
  PreResult: (this: UniversalSDK, ctx: Context) => void | Promise<any>
  PostOperation: (this: UniversalSDK, ctx: Context) => void | Promise<any>
}


export {
  Target,
  Context,
  Control,
  Operation,
  Response,
  Result,
  Spec,
}


export type {
  Feature,
  FeatureOptions,

  UniversalOptions,
}
