
import { test, describe } from 'node:test'
import assert from 'node:assert'


import { UniversalManager, UniversalSDK } from '..'


describe('exists', async () => {

  test('test-mode', async () => {
    const um = new UniversalManager({ registry: __dirname + '/../test/registry' })
    const solardk = um.make('voxgig-solardemo')
    assert(null != solardk)
  })

})
