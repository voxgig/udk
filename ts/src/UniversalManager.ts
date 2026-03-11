


import { inspect } from 'node:util'
import * as Fs from 'node:fs'
import Path from 'node:path'


import type {
  UniversalOptions
} from './types'



import { Utility } from './utility/Utility'


import { UniversalSDK } from './UniversalSDK'


const stdutil = new Utility()


class UniversalManager {
  _options: UniversalOptions

  _utility = new Utility()

  constructor(options: Partial<UniversalOptions>) {

    // TODO: validation
    this._options = options as UniversalOptions

  }


  options() {
    return this._utility.struct.clone(this._options)
  }


  utility() {
    return this._utility.struct.clone(this._utility)
  }


  make(ref: string): UniversalSDK {
    const model = this.resolveModel(ref)
    const udk = new UniversalSDK(this, { ref, model })
    return udk
  }


  resolveModel(ref: string) {
    const modelpath = Path.join(this._options.registry, 'local', ref + '.json')
    const modelsrc = Fs.readFileSync(modelpath).toString()
    const model = JSON.parse(modelsrc)
    // console.log('resolveModel', modelpath, model)
    return model
  }


  toJSON() {
    return { name: 'Universal' }
  }

  toString() {
    return 'Universal ' + this._utility.struct.jsonify(this.toJSON())
  }

  [inspect.custom]() {
    return this.toString()
  }

}




export {
  stdutil,

  UniversalManager,
  UniversalSDK,
}


