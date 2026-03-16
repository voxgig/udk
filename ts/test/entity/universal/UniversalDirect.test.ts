
const envlocal = __dirname + '/../../../.env.local'
require('dotenv').config({ quiet: true, path: [envlocal] })

import { test, describe } from 'node:test'
import assert from 'node:assert'


import { UniversalManager, UniversalSDK, stdutil } from '../../..'

import {
  envOverride,
} from '../../utility'


describe('UniversalDirect', async () => {

  const um = new UniversalManager({ registry: __dirname + '/../../../test/registry' })
  const sdk = um.make('voxgig-solardemo')
  const entityMap: any = sdk._config.entity
  const struct = sdk.utility().struct
  const items = struct.items

  const live = 'TRUE' === process.env.UNIVERSAL_TEST_LIVE


  test('direct-exists', async () => {
    const setup = directSetup(um, sdk)
    assert('function' === typeof setup.client.direct)
    assert('function' === typeof setup.client.prepare)
  })


  items(entityMap, (item: any[]) => {
    const entityDef = item[1]
    const entityName = entityDef.name
    const ops = entityDef.op || {}

    const hasLoad = null != ops.load
    const hasList = null != ops.list

    if (!hasLoad && !hasList) {
      return
    }

    if (hasLoad) {
      const loadOp = ops.load
      const loadTarget = loadOp.targets?.[0]

      if (null != loadTarget) {
        test('direct-load-' + entityName, async () => {
          const loadPath = (loadTarget.parts || []).join('/')
          const loadParams = loadTarget.args?.params || []

          if (live) {
            const idmap = await resolveIdmap(um, sdk, entityName, entityMap)
            const setup = directSetup(um, sdk)

            // First list to discover a real entity ID.
            if (hasList) {
              const listTarget = ops.list.targets?.[0]
              if (null != listTarget) {
                const listPath = (listTarget.parts || []).join('/')
                const listParams = listTarget.args?.params || []

                // Try multiple parent refs to find one with child entities.
                let found: any = null
                let lparams: any = {}
                for (let t = 0; t < 3 && null == found; t++) {
                  lparams = {}
                  for (const p of listParams) {
                    const ref = p.name.replace(/_id$/, '') +
                      String(t).padStart(2, '0')
                    lparams[p.name] = idmap[ref] || ref
                  }

                  const listResult: any = await setup.client.direct({
                    path: listPath,
                    method: 'GET',
                    params: lparams,
                  })

                  assert(listResult.ok === true)
                  assert(Array.isArray(listResult.data))

                  if (listResult.data.length >= 1) {
                    found = listResult.data[0]
                  }
                }

                if (null != found) {
                  const params: any = {}
                  for (const p of loadParams) {
                    params[p.name] = found[p.name] || lparams[p.name]
                  }

                  const result: any = await setup.client.direct({
                    path: loadPath,
                    method: 'GET',
                    params,
                  })

                  assert(result.ok === true)
                  assert(result.status === 200)
                  assert(null != result.data)
                  assert(result.data.id === found.id)
                }
              }
            }
          }
          else {
            const setup = directSetup(um, sdk, { id: 'direct01' })
            const { client, calls } = setup

            const params: any = {}
            for (let i = 0; i < loadParams.length; i++) {
              params[loadParams[i].name] = 'direct0' + (i + 1)
            }

            const result: any = await client.direct({
              path: loadPath,
              method: 'GET',
              params,
            })

            assert(result.ok === true)
            assert(result.status === 200)
            assert(null != result.data)
            assert(result.data.id === 'direct01')

            assert(calls.length === 1)
            assert(calls[0].init.method === 'GET')

            for (let i = 0; i < loadParams.length; i++) {
              assert(calls[0].url.includes('direct0' + (i + 1)))
            }
          }
        })
      }
    }

    if (hasList) {
      const listOp = ops.list
      const listTarget = listOp.targets?.[0]

      if (null != listTarget) {
        test('direct-list-' + entityName, async () => {
          const listPath = (listTarget.parts || []).join('/')
          const listParams = listTarget.args?.params || []

          if (live) {
            const idmap = await resolveIdmap(um, sdk, entityName, entityMap)
            const setup = directSetup(um, sdk)

            // For entities with parent params, try each known parent
            // to find one that has child entities.
            let found = false
            const maxTries = listParams.length > 0 ? 3 : 1
            for (let t = 0; t < maxTries && !found; t++) {
              const params: any = {}
              for (const p of listParams) {
                const base = (p.name === 'id' ? entityName : p.name.replace(/_id$/, ''))
                const ref = base + String(t).padStart(2, '0')
                params[p.name] = idmap[ref] || ref
              }

              const result: any = await setup.client.direct({
                path: listPath,
                method: 'GET',
                params,
              })

              assert(result.ok === true)
              assert(result.status === 200)
              assert(Array.isArray(result.data))

              if (result.data.length >= 1) {
                found = true
              }
            }

            if (listParams.length === 0) {
              assert(found, 'expected at least one entity in list')
            }
          }
          else {
            const setup = directSetup(um, sdk, [{ id: 'direct01' }, { id: 'direct02' }])
            const { client, calls } = setup

            const params: any = {}
            for (let i = 0; i < listParams.length; i++) {
              params[listParams[i].name] = 'direct0' + (i + 1)
            }

            const result: any = await client.direct({
              path: listPath,
              method: 'GET',
              params,
            })

            assert(result.ok === true)
            assert(result.status === 200)
            assert(Array.isArray(result.data))
            assert(result.data.length === 2)

            assert(calls.length === 1)
            assert(calls[0].init.method === 'GET')

            for (let i = 0; i < listParams.length; i++) {
              assert(calls[0].url.includes('direct0' + (i + 1)))
            }
          }
        })
      }
    }
  })

})



async function resolveIdmap(um: any, sdk: any, entityName: string, entityMap: any): Promise<any> {
  const clientStruct = sdk.utility().struct
  const items = clientStruct.items
  const transform = clientStruct.transform

  const idEntries: string[] = []
  items(entityMap, (item: any[]) => {
    const ename = item[1].name
    for (let i = 0; i < 3; i++) {
      idEntries.push(`${ename}${String(i).padStart(2, '0')}`)
    }
  })

  let idmap = transform(
    idEntries,
    {
      '`$PACK`': ['', {
        '`$KEY`': '`$COPY`',
        '`$VAL`': ['`$FORMAT`', 'upper', '`$COPY`']
      }]
    })

  const env = envOverride({
    'UNIVERSAL_TEST_ENTID': idmap,
    'UNIVERSAL_TEST_LIVE': 'FALSE',
  })

  idmap = env['UNIVERSAL_TEST_ENTID']

  // In live mode, discover real parent entity IDs by listing parent entities.
  if ('TRUE' === process.env.UNIVERSAL_TEST_LIVE) {
    const liveClient = new UniversalSDK(um, {
      ref: 'voxgig-solardemo',
      model: sdk._options.model,
    })

    const discoveries: Promise<void>[] = []
    items(entityMap, (item: any[]) => {
      const eDef = item[1]
      const eName = eDef.name
      const listOp = eDef.op?.list
      const listTarget = listOp?.targets?.[0]
      if (null == listTarget) return

      const listParams = listTarget.args?.params || []
      if (listParams.length > 0) return // skip nested entities in discovery

      const listPath = (listTarget.parts || []).join('/')
      discoveries.push((async () => {
        const res: any = await liveClient.direct({ path: listPath, method: 'GET', params: {} })
        if (res.ok && Array.isArray(res.data)) {
          for (let i = 0; i < Math.min(res.data.length, 3); i++) {
            const ref = `${eName}${String(i).padStart(2, '0')}`
            idmap[ref] = res.data[i].id
          }
        }
      })())
    })

    await Promise.all(discoveries)
  }

  return idmap
}


function directSetup(um: any, sdk: any, mockres?: any) {
  const live = 'TRUE' === process.env.UNIVERSAL_TEST_LIVE

  if (live) {
    const client = new UniversalSDK(um, {
      ref: 'voxgig-solardemo',
      model: sdk._options.model,
    })
    return { client, calls: [] as any[], live: true }
  }

  const calls: any[] = []

  const mockFetch = async (url: string, init: any) => {
    calls.push({ url, init })
    return {
      status: 200,
      statusText: 'OK',
      headers: {},
      json: async () => (null != mockres ? mockres : { id: 'direct01' }),
    }
  }

  const client = new UniversalSDK(um, {
    ref: 'voxgig-solardemo',
    model: sdk._options.model,
    base: 'http://localhost:8080',
    system: { fetch: mockFetch },
  })

  return { client, calls, live: false }
}
