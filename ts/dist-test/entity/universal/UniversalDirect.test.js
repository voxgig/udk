"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const envlocal = __dirname + '/../../../.env.local';
require('dotenv').config({ quiet: true, path: [envlocal] });
const node_test_1 = require("node:test");
const node_assert_1 = __importDefault(require("node:assert"));
const __1 = require("../../..");
const utility_1 = require("../../utility");
(0, node_test_1.describe)('UniversalDirect', async () => {
    const um = new __1.UniversalManager({ registry: __dirname + '/../../../test/registry' });
    const sdk = um.make('voxgig-solardemo');
    const entityMap = sdk._config.entity;
    const struct = sdk.utility().struct;
    const items = struct.items;
    const live = 'TRUE' === process.env.UNIVERSAL_TEST_LIVE;
    (0, node_test_1.test)('direct-exists', async () => {
        const setup = directSetup(um, sdk);
        (0, node_assert_1.default)('function' === typeof setup.client.direct);
        (0, node_assert_1.default)('function' === typeof setup.client.prepare);
    });
    items(entityMap, (item) => {
        const entityDef = item[1];
        const entityName = entityDef.name;
        const ops = entityDef.op || {};
        const hasLoad = null != ops.load;
        const hasList = null != ops.list;
        if (!hasLoad && !hasList) {
            return;
        }
        if (hasLoad) {
            const loadOp = ops.load;
            const loadTarget = loadOp.targets?.[0];
            if (null != loadTarget) {
                (0, node_test_1.test)('direct-load-' + entityName, async () => {
                    const loadPath = (loadTarget.parts || []).join('/');
                    const loadParams = loadTarget.args?.params || [];
                    if (live) {
                        const idmap = resolveIdmap(um, sdk, entityName, entityMap);
                        const setup = directSetup(um, sdk);
                        // First list to discover a real entity ID.
                        if (hasList) {
                            const listTarget = ops.list.targets?.[0];
                            if (null != listTarget) {
                                const listPath = (listTarget.parts || []).join('/');
                                const listParams = listTarget.args?.params || [];
                                const lparams = {};
                                for (const p of listParams) {
                                    const ref = p.name.replace(/_id$/, '') + '01';
                                    lparams[p.name] = idmap[ref] || ref;
                                }
                                const listResult = await setup.client.direct({
                                    path: listPath,
                                    method: 'GET',
                                    params: lparams,
                                });
                                (0, node_assert_1.default)(listResult.ok === true);
                                (0, node_assert_1.default)(Array.isArray(listResult.data));
                                (0, node_assert_1.default)(listResult.data.length >= 1);
                                const found = listResult.data[0];
                                const params = {};
                                for (const p of loadParams) {
                                    params[p.name] = found[p.name] || lparams[p.name];
                                }
                                const result = await setup.client.direct({
                                    path: loadPath,
                                    method: 'GET',
                                    params,
                                });
                                (0, node_assert_1.default)(result.ok === true);
                                (0, node_assert_1.default)(result.status === 200);
                                (0, node_assert_1.default)(null != result.data);
                                (0, node_assert_1.default)(result.data.id === found.id);
                            }
                        }
                    }
                    else {
                        const setup = directSetup(um, sdk, { id: 'direct01' });
                        const { client, calls } = setup;
                        const params = {};
                        for (let i = 0; i < loadParams.length; i++) {
                            params[loadParams[i].name] = 'direct0' + (i + 1);
                        }
                        const result = await client.direct({
                            path: loadPath,
                            method: 'GET',
                            params,
                        });
                        (0, node_assert_1.default)(result.ok === true);
                        (0, node_assert_1.default)(result.status === 200);
                        (0, node_assert_1.default)(null != result.data);
                        (0, node_assert_1.default)(result.data.id === 'direct01');
                        (0, node_assert_1.default)(calls.length === 1);
                        (0, node_assert_1.default)(calls[0].init.method === 'GET');
                        for (let i = 0; i < loadParams.length; i++) {
                            (0, node_assert_1.default)(calls[0].url.includes('direct0' + (i + 1)));
                        }
                    }
                });
            }
        }
        if (hasList) {
            const listOp = ops.list;
            const listTarget = listOp.targets?.[0];
            if (null != listTarget) {
                (0, node_test_1.test)('direct-list-' + entityName, async () => {
                    const listPath = (listTarget.parts || []).join('/');
                    const listParams = listTarget.args?.params || [];
                    if (live) {
                        const idmap = resolveIdmap(um, sdk, entityName, entityMap);
                        const params = {};
                        for (const p of listParams) {
                            const ref = (p.name === 'id' ? entityName : p.name.replace(/_id$/, '')) + '01';
                            params[p.name] = idmap[ref] || ref;
                        }
                        const setup = directSetup(um, sdk);
                        const result = await setup.client.direct({
                            path: listPath,
                            method: 'GET',
                            params,
                        });
                        (0, node_assert_1.default)(result.ok === true);
                        (0, node_assert_1.default)(result.status === 200);
                        (0, node_assert_1.default)(Array.isArray(result.data));
                        (0, node_assert_1.default)(result.data.length >= 1);
                    }
                    else {
                        const setup = directSetup(um, sdk, [{ id: 'direct01' }, { id: 'direct02' }]);
                        const { client, calls } = setup;
                        const params = {};
                        for (let i = 0; i < listParams.length; i++) {
                            params[listParams[i].name] = 'direct0' + (i + 1);
                        }
                        const result = await client.direct({
                            path: listPath,
                            method: 'GET',
                            params,
                        });
                        (0, node_assert_1.default)(result.ok === true);
                        (0, node_assert_1.default)(result.status === 200);
                        (0, node_assert_1.default)(Array.isArray(result.data));
                        (0, node_assert_1.default)(result.data.length === 2);
                        (0, node_assert_1.default)(calls.length === 1);
                        (0, node_assert_1.default)(calls[0].init.method === 'GET');
                        for (let i = 0; i < listParams.length; i++) {
                            (0, node_assert_1.default)(calls[0].url.includes('direct0' + (i + 1)));
                        }
                    }
                });
            }
        }
    });
});
function resolveIdmap(um, sdk, entityName, entityMap) {
    const clientStruct = sdk.utility().struct;
    const items = clientStruct.items;
    const transform = clientStruct.transform;
    const idEntries = [];
    items(entityMap, (item) => {
        const ename = item[1].name;
        for (let i = 0; i < 3; i++) {
            idEntries.push(`${ename}${String(i).padStart(2, '0')}`);
        }
    });
    let idmap = transform(idEntries, {
        '`$PACK`': ['', {
                '`$KEY`': '`$COPY`',
                '`$VAL`': ['`$FORMAT`', 'upper', '`$COPY`']
            }]
    });
    const env = (0, utility_1.envOverride)({
        'UNIVERSAL_TEST_ENTID': idmap,
        'UNIVERSAL_TEST_LIVE': 'FALSE',
    });
    return env['UNIVERSAL_TEST_ENTID'];
}
function directSetup(um, sdk, mockres) {
    const live = 'TRUE' === process.env.UNIVERSAL_TEST_LIVE;
    if (live) {
        const client = new __1.UniversalSDK(um, {
            ref: 'voxgig-solardemo',
            model: sdk._options.model,
        });
        return { client, calls: [], live: true };
    }
    const calls = [];
    const mockFetch = async (url, init) => {
        calls.push({ url, init });
        return {
            status: 200,
            statusText: 'OK',
            headers: {},
            json: async () => (null != mockres ? mockres : { id: 'direct01' }),
        };
    };
    const client = new __1.UniversalSDK(um, {
        ref: 'voxgig-solardemo',
        model: sdk._options.model,
        base: 'http://localhost:8080',
        system: { fetch: mockFetch },
    });
    return { client, calls, live: false };
}
//# sourceMappingURL=UniversalDirect.test.js.map