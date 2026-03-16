"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TEST_JSON_FILE = exports.SDK = exports.UniversalSDK = exports.um = void 0;
const __1 = require("../..");
Object.defineProperty(exports, "UniversalSDK", { enumerable: true, get: function () { return __1.UniversalSDK; } });
const TEST_JSON_FILE = '../../.sdk/test/test.json';
exports.TEST_JSON_FILE = TEST_JSON_FILE;
const um = new __1.UniversalManager({ registry: __dirname + '/../../test/registry' });
exports.um = um;
const SDK = um.make('voxgig-solardemo');
exports.SDK = SDK;
//# sourceMappingURL=index.js.map