"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TEST_JSON_FILE = exports.SDK = void 0;
const __1 = require("../..");
const TEST_JSON_FILE = '../../.sdk/test/test.json';
exports.TEST_JSON_FILE = TEST_JSON_FILE;
const um = new __1.UniversalManager({ registry: __dirname + '/../../test/registry' });
const SDK = um.make('voxgig-solardemo');
exports.SDK = SDK;
//# sourceMappingURL=index.js.map