"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const node_test_1 = require("node:test");
const node_assert_1 = __importDefault(require("node:assert"));
const __1 = require("..");
(0, node_test_1.describe)('exists', async () => {
    (0, node_test_1.test)('test-mode', async () => {
        const um = new __1.UniversalManager({ registry: __dirname + '/../test/registry' });
        const solardk = um.make('voxgig-solardemo');
        (0, node_assert_1.default)(null != solardk);
    });
});
//# sourceMappingURL=exists.test.js.map