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
(0, node_test_1.describe)('UniversalEntity', async () => {
    const um = new __1.UniversalManager({ registry: __dirname + '/../../../test/registry' });
    const sdk = um.make('voxgig-solardemo');
    const entityMap = sdk._config.entity;
    (0, node_test_1.test)('instance', async () => {
        const struct = sdk.utility().struct;
        const items = struct.items;
        items(entityMap, (item) => {
            const name = item[1].name;
            const uent = sdk.Entity(name);
            (0, node_assert_1.default)(null != uent);
        });
    });
    const struct = sdk.utility().struct;
    const items = struct.items;
    items(entityMap, (item) => {
        const entityDef = item[1];
        const entityName = entityDef.name;
        (0, node_test_1.test)('basic-' + entityName, async () => {
            const setup = basicSetup(um, entityMap, entityName);
            const client = setup.client;
            const struct = setup.struct;
            const ops = entityDef.op || {};
            const ref = entityName + '_ref01';
            const ent = client.Entity(entityName);
            let createdData = null;
            if (ops.create) {
                createdData = await testCreate(setup, ent, entityName, ref, entityDef);
            }
            if (ops.list) {
                await testList(setup, ent, entityDef, createdData, true);
            }
            if (ops.update && createdData) {
                await testUpdate(setup, ent, entityName, entityDef, createdData);
            }
            if (ops.load && createdData) {
                await testLoad(setup, ent, entityDef, createdData);
            }
            if (ops.remove && createdData) {
                await testRemove(setup, ent, entityDef, createdData);
            }
            if (ops.list && ops.remove && createdData) {
                await testList(setup, ent, entityDef, createdData, false);
            }
        });
    });
});
function resolveIdFields(data, idmap) {
    const out = { ...data };
    for (const key of Object.keys(out)) {
        if (key.endsWith('_id')) {
            const baseRef = key.substring(0, key.length - 3) + '01';
            if (null != idmap[baseRef]) {
                out[key] = idmap[baseRef];
            }
        }
    }
    return out;
}
async function testCreate(setup, ent, entityName, ref, entityDef) {
    let reqdata = resolveIdFields(setup.data.new[entityName][ref], setup.idmap);
    const resdata = await ent.create(reqdata);
    (0, node_assert_1.default)(null != resdata.id);
    return resdata;
}
async function testList(setup, ent, entityDef, createdData, shouldExist) {
    const struct = setup.struct;
    const isempty = struct.isempty;
    const select = struct.select;
    const matchFields = getDefaultTargetFields(entityDef, 'list');
    const match = {};
    for (const field of matchFields) {
        if (field !== 'id' && createdData && null != createdData[field]) {
            match[field] = createdData[field];
        }
    }
    const list = await ent.list(match);
    if (createdData) {
        if (shouldExist) {
            (0, node_assert_1.default)(!isempty(select(list, { id: createdData.id })));
        }
        else {
            (0, node_assert_1.default)(isempty(select(list, { id: createdData.id })));
        }
    }
}
async function testUpdate(setup, ent, entityName, entityDef, createdData) {
    const reqdata = {};
    reqdata.id = createdData.id;
    const matchFields = getDefaultTargetFields(entityDef, 'update');
    for (const field of matchFields) {
        if (field !== 'id' && null != createdData[field]) {
            reqdata[field] = createdData[field];
        }
    }
    const textfield = findTextField(entityDef);
    let markdef = null;
    if (textfield) {
        markdef = { name: textfield, value: 'Mark01-' + entityName + '_ref01_' + setup.now };
        reqdata[markdef.name] = markdef.value;
    }
    const resdata = await ent.update(reqdata);
    (0, node_assert_1.default)(resdata.id === reqdata.id);
    if (markdef) {
        (0, node_assert_1.default)(resdata[markdef.name] === markdef.value);
    }
}
async function testLoad(setup, ent, entityDef, createdData) {
    const matchFields = getDefaultTargetFields(entityDef, 'load');
    const match = {};
    match.id = createdData.id;
    for (const field of matchFields) {
        if (field !== 'id' && null != createdData[field]) {
            match[field] = createdData[field];
        }
    }
    const resdata = await ent.load(match);
    (0, node_assert_1.default)(resdata.id === createdData.id);
}
async function testRemove(setup, ent, entityDef, createdData) {
    const matchFields = getDefaultTargetFields(entityDef, 'remove');
    const match = {};
    match.id = createdData.id;
    for (const field of matchFields) {
        if (field !== 'id' && null != createdData[field]) {
            match[field] = createdData[field];
        }
    }
    await ent.remove(match);
}
function getDefaultTargetFields(entityDef, opname) {
    const op = entityDef.op?.[opname];
    if (!op)
        return [];
    const targets = op.targets || [];
    for (let i = targets.length - 1; i >= 0; i--) {
        if (!targets[i].select?.$action) {
            return targets[i].select?.exist || [];
        }
    }
    return [];
}
function findTextField(entityDef) {
    for (const field of entityDef.fields || []) {
        if (field.type === '`$STRING`' && field.name !== 'id' && !field.name.endsWith('_id')) {
            return field.name;
        }
    }
    return null;
}
function makeEntityTestData(entityDef) {
    const fields = entityDef.fields || [];
    const name = entityDef.name;
    const data = {
        existing: { [name]: {} },
        new: { [name]: {} }
    };
    const idcount = 3;
    const refs = Array.from({ length: idcount }, (_, i) => `${name}${String(i).padStart(2, '0')}`);
    const idmapLocal = refs.reduce((a, ref) => (a[ref] = ref.toUpperCase(), a), {});
    let idx = 1;
    for (const ref of refs) {
        const id = idmapLocal[ref];
        const ent = data.existing[name][id] = {};
        makeEntityTestFields(fields, idx++, ent);
        ent.id = id;
    }
    const newRef = name + '_ref01';
    const newEnt = data.new[name][newRef] = {};
    makeEntityTestFields(fields, idx++, newEnt);
    delete newEnt.id;
    return data;
}
function makeEntityTestFields(fields, start, entdata) {
    let num = start * fields.length * 10;
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
                                        's' + (num.toString(16));
        num++;
    }
}
function basicSetup(um, entityMap, entityName, extra) {
    const options = {};
    const allExisting = {};
    const allNew = {};
    const struct = __1.stdutil.struct;
    const items = struct.items;
    const flatten = struct.flatten;
    items(entityMap, (item) => {
        const entityDef = item[1];
        const testData = makeEntityTestData(entityDef);
        Object.assign(allExisting, testData.existing);
        Object.assign(allNew, testData.new);
    });
    options.entity = allExisting;
    const sdk = um.make('voxgig-solardemo');
    let client = sdk.test(options, { ref: 'voxgig-solardemo', model: sdk._options.model });
    const clientStruct = client.utility().struct;
    const merge = clientStruct.merge;
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
        'UNIVERSAL_TEST_EXPLAIN': 'FALSE',
        'UNIVERSAL_APIKEY': 'NONE',
    });
    idmap = env['UNIVERSAL_TEST_ENTID'];
    if ('TRUE' === env.UNIVERSAL_TEST_LIVE) {
        const liveopts = {
            ref: 'voxgig-solardemo',
            model: sdk._options.model,
            apikey: env.UNIVERSAL_APIKEY,
        };
        client = new __1.UniversalSDK(um, null != extra ? merge([liveopts, extra]) : liveopts);
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
    };
    return setup;
}
//# sourceMappingURL=UniversalEntity.test.js.map