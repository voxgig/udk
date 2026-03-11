

import { UniversalManager, UniversalSDK, stdutil } from '../..'


const TEST_JSON_FILE = '../../.sdk/test/test.json'


const um = new UniversalManager({ registry: __dirname + '/../../test/registry' })
const SDK = um.make('voxgig-solardemo')


export {
  SDK,
  TEST_JSON_FILE,
}
