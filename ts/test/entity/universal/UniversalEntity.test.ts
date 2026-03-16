
const envlocal = __dirname + '/../../../.env.local'
require('dotenv').config({ quiet: true, path: [envlocal] })

import { test, describe } from 'node:test'
import assert from 'node:assert'


import { UniversalManager, UniversalSDK, stdutil } from '../../..'

import {
  envOverride,
} from '../../utility'


describe('UniversalEntity', async () => {

  const um = new UniversalManager({ registry: __dirname + '/../../../test/registry' })
  const sdk = um.make('voxgig-solardemo')
  const entityMap: any = sdk._config.entity


  test('instance', async () => {
    const struct = sdk.utility().struct
    const items = struct.items

    items(entityMap, (item: any[]) => {
      const name = item[1].name
      const uent = sdk.Entity(name)
      assert(null != uent)
    })
  })


  const struct = sdk.utility().struct
  const items = struct.items

  items(entityMap, (item: any[]) => {
    const entityDef = item[1]
    const entityName = entityDef.name

    test('basic-' + entityName, async () => {
      const setup = await basicSetup(um, entityMap, entityName)
      const client = setup.client
      const struct = setup.struct

      const ops = entityDef.op || {}
      const ref = entityName + '_ref01'
      const ent = client.Entity(entityName)

      let createdData: any = null

      if (ops.create) {
        createdData = await testCreate(setup, ent, entityName, ref, entityDef)
      }

      if (ops.list) {
        await testList(setup, ent, entityDef, createdData, true)
      }

      if (ops.update && createdData) {
        await testUpdate(setup, ent, entityName, entityDef, createdData)
      }

      if (ops.load && createdData) {
        await testLoad(setup, ent, entityDef, createdData)
      }

      if (ops.remove && createdData) {
        await testRemove(setup, ent, entityDef, createdData)
      }

      if (ops.list && ops.remove && createdData) {
        await testList(setup, ent, entityDef, createdData, false)
      }
    })
  })

})



function resolveIdFields(data: any, idmap: any): any {
  const out: any = { ...data }
  for (const key of Object.keys(out)) {
    if (key.endsWith('_id')) {
      const baseRef = key.substring(0, key.length - 3) + '01'
      if (null != idmap[baseRef]) {
        out[key] = idmap[baseRef]
      }
    }
  }
  return out
}


async function testCreate(
  setup: any,
  ent: any,
  entityName: string,
  ref: string,
  entityDef: any,
) {
  let reqdata = resolveIdFields(setup.data.new[entityName][ref], setup.idmap)
  const resdata = await ent.create(reqdata)
  assert(null != resdata.id)
  return resdata
}


async function testList(
  setup: any,
  ent: any,
  entityDef: any,
  createdData: any,
  shouldExist: boolean,
) {
  const struct = setup.struct
  const isempty = struct.isempty
  const select = struct.select

  const matchFields = getDefaultTargetFields(entityDef, 'list')
  const match: any = {}
  for (const field of matchFields) {
    if (field !== 'id' && createdData && null != createdData[field]) {
      match[field] = createdData[field]
    }
  }

  const list = await ent.list(match)

  if (createdData) {
    if (shouldExist) {
      assert(!isempty(select(list, { id: createdData.id })))
    }
    else {
      assert(isempty(select(list, { id: createdData.id })))
    }
  }
}


async function testUpdate(
  setup: any,
  ent: any,
  entityName: string,
  entityDef: any,
  createdData: any,
) {
  const reqdata: any = {}
  reqdata.id = createdData.id

  const matchFields = getDefaultTargetFields(entityDef, 'update')
  for (const field of matchFields) {
    if (field !== 'id' && null != createdData[field]) {
      reqdata[field] = createdData[field]
    }
  }

  const textfield = findTextField(entityDef)
  let markdef: any = null

  if (textfield) {
    markdef = { name: textfield, value: 'Mark01-' + entityName + '_ref01_' + setup.now }
    reqdata[markdef.name] = markdef.value
  }

  const resdata = await ent.update(reqdata)
  assert(resdata.id === reqdata.id)

  if (markdef) {
    assert(resdata[markdef.name] === markdef.value)
  }
}


async function testLoad(
  setup: any,
  ent: any,
  entityDef: any,
  createdData: any,
) {
  const matchFields = getDefaultTargetFields(entityDef, 'load')
  const match: any = {}
  match.id = createdData.id
  for (const field of matchFields) {
    if (field !== 'id' && null != createdData[field]) {
      match[field] = createdData[field]
    }
  }

  const resdata = await ent.load(match)
  assert(resdata.id === createdData.id)
}


async function testRemove(
  setup: any,
  ent: any,
  entityDef: any,
  createdData: any,
) {
  const matchFields = getDefaultTargetFields(entityDef, 'remove')
  const match: any = {}
  match.id = createdData.id
  for (const field of matchFields) {
    if (field !== 'id' && null != createdData[field]) {
      match[field] = createdData[field]
    }
  }

  await ent.remove(match)
}



function getDefaultTargetFields(entityDef: any, opname: string): string[] {
  const op = entityDef.op?.[opname]
  if (!op) return []
  const targets = op.targets || []
  for (let i = targets.length - 1; i >= 0; i--) {
    if (!targets[i].select?.$action) {
      return targets[i].select?.exist || []
    }
  }
  return []
}


function findTextField(entityDef: any): string | null {
  for (const field of entityDef.fields || []) {
    if (field.type === '`$STRING`' && field.name !== 'id' && !field.name.endsWith('_id')) {
      return field.name
    }
  }
  return null
}


function makeEntityTestData(entityDef: any) {
  const fields = entityDef.fields || []
  const name = entityDef.name

  const data: any = {
    existing: { [name]: {} },
    new: { [name]: {} }
  }

  const idcount = 3
  const refs = Array.from({ length: idcount }, (_, i) =>
    `${name}${String(i).padStart(2, '0')}`)

  const idmapLocal = refs.reduce((a: any, ref) => (a[ref] = ref.toUpperCase(), a), {})

  let idx = 1
  for (const ref of refs) {
    const id = idmapLocal[ref]
    const ent: any = data.existing[name][id] = {}
    makeEntityTestFields(fields, idx++, ent)
    ent.id = id
  }

  const newRef = name + '_ref01'
  const newEnt: any = data.new[name][newRef] = {}
  makeEntityTestFields(fields, idx++, newEnt)
  delete newEnt.id

  return data
}


function makeEntityTestFields(fields: any[], start: number, entdata: any) {
  let num = start * fields.length * 10
  for (const field of fields) {
    entdata[field.name] =
      field.name.endsWith('_id') ?
        field.name.substring(0, field.name.length - 3).toUpperCase() + '01' :
        '`$NUMBER`' === field.type ? num :
          '`$BOOLEAN`' === field.type ? 0 === num % 2 :
            '`$OBJECT`' === field.type ? {} :
              '`$MAP`' === field.type ? {} :
                '`$ARRAY`' === field.type ? [] :
                  '`$LIST`' === field.type ? [] :
                    's' + (num.toString(16))
    num++
  }
}


async function basicSetup(um: any, entityMap: any, entityName: string, extra?: any) {
  const options: any = {}

  const allExisting: any = {}
  const allNew: any = {}

  const struct = stdutil.struct
  const items = struct.items
  const flatten = struct.flatten

  items(entityMap, (item: any[]) => {
    const entityDef = item[1]
    const testData = makeEntityTestData(entityDef)
    Object.assign(allExisting, testData.existing)
    Object.assign(allNew, testData.new)
  })

  options.entity = allExisting

  const sdk = um.make('voxgig-solardemo')
  let client = sdk.test(options, { ref: 'voxgig-solardemo', model: sdk._options.model })

  const clientStruct = client.utility().struct
  const merge = clientStruct.merge
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
    'UNIVERSAL_TEST_EXPLAIN': 'FALSE',
    'UNIVERSAL_APIKEY': 'NONE',
  })

  idmap = env['UNIVERSAL_TEST_ENTID']

  if ('TRUE' === env.UNIVERSAL_TEST_LIVE) {
    const liveopts: any = {
      ref: 'voxgig-solardemo',
      model: sdk._options.model,
      apikey: env.UNIVERSAL_APIKEY,
    }
    client = new UniversalSDK(um, null != extra ? merge([liveopts, extra]) : liveopts)

    // Discover real parent entity IDs from the live API.
    for (const item of Object.values(entityMap) as any[]) {
      const eDef = item
      const eName = eDef.name
      const listOp = eDef.op?.list
      const listTarget = listOp?.targets?.[0]
      if (null == listTarget) continue

      const listParams = listTarget.args?.params || []
      if (listParams.length > 0) continue // skip nested entities

      const listPath = (listTarget.parts || []).join('/')
      const res: any = await client.direct({ path: listPath, method: 'GET', params: {} })
      if (res.ok && Array.isArray(res.data)) {
        for (let i = 0; i < Math.min(res.data.length, 3); i++) {
          const ref = `${eName}${String(i).padStart(2, '0')}`
          idmap[ref] = res.data[i].id
        }
      }
    }
  }

  const setup = {
    idmap,
    env,
    options,
    client,
    struct: client.utility().struct,
    data: { existing: allExisting, new: allNew },
    explain: 'TRUE' === env.UNIVERSAL_TEST_EXPLAIN,
    now: Date.now(),
  }

  return setup
}
