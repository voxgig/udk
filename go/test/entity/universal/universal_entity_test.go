package universal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	sdk "voxgiguniversalsdk"
	"voxgiguniversalsdk/core"

	vs "github.com/voxgig/struct"
)

func TestUniversalEntity(t *testing.T) {
	um := sdk.NewUniversalManager(map[string]any{
		"registry": "../../../test/registry",
	})
	baseSDK := um.Make("voxgig-solardemo")
	entityMap := core.ToMapAny(vs.GetPath([]any{"entity"}, baseSDK.GetRootCtx().Config))

	t.Run("instance", func(t *testing.T) {
		for _, item := range vs.Items(entityMap) {
			entDef := core.ToMapAny(item[1])
			name, _ := vs.GetProp(entDef, "name").(string)
			uent := baseSDK.Entity(name, nil)
			if uent == nil {
				t.Errorf("Entity(%q) should not be nil", name)
			}
		}
	})

	for _, item := range vs.Items(entityMap) {
		entDef := core.ToMapAny(item[1])
		entityName, _ := vs.GetProp(entDef, "name").(string)

		t.Run("basic-"+entityName, func(t *testing.T) {
			setup := basicSetup(t, um, entityMap, entityName)
			client := setup.client

			ops := core.ToMapAny(vs.GetProp(entDef, "op"))
			if ops == nil {
				ops = map[string]any{}
			}
			ref := entityName + "_ref01"
			ent := client.Entity(entityName, nil)

			var createdData map[string]any

			if vs.GetProp(ops, "create") != nil {
				createdData = testCreate(t, setup, ent, entityName, ref, entDef)
			}

			if vs.GetProp(ops, "list") != nil {
				testList(t, setup, ent, entDef, createdData, true)
			}

			if vs.GetProp(ops, "update") != nil && createdData != nil {
				testUpdate(t, setup, ent, entityName, entDef, createdData)
			}

			if vs.GetProp(ops, "load") != nil && createdData != nil {
				testLoad(t, setup, ent, entDef, createdData)
			}

			if vs.GetProp(ops, "remove") != nil && createdData != nil {
				testRemove(t, setup, ent, entDef, createdData)
			}

			if vs.GetProp(ops, "list") != nil && vs.GetProp(ops, "remove") != nil && createdData != nil {
				testList(t, setup, ent, entDef, createdData, false)
			}
		})
	}
}

func resolveIdFields(data map[string]any, idmap map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range data {
		out[k] = v
	}
	for key := range out {
		if strings.HasSuffix(key, "_id") {
			baseRef := key[:len(key)-3] + "01"
			if idmap[baseRef] != nil {
				out[key] = idmap[baseRef]
			}
		}
	}
	return out
}

func testCreate(
	t *testing.T,
	setup *entityTestSetup,
	ent sdk.UniversalEntity,
	entityName string,
	ref string,
	entityDef map[string]any,
) map[string]any {
	t.Helper()

	newData := core.ToMapAny(vs.GetPath([]any{"new", entityName, ref}, setup.data))
	if newData == nil {
		t.Fatalf("No new data for %s/%s", entityName, ref)
	}
	reqdata := resolveIdFields(newData, setup.idmap)

	result, err := ent.Create(reqdata, nil)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	resdata := core.ToMapAny(result)
	if resdata == nil {
		t.Fatal("Create result should be a map")
	}
	if resdata["id"] == nil {
		t.Fatal("Created entity should have an id")
	}
	return resdata
}

func testList(
	t *testing.T,
	setup *entityTestSetup,
	ent sdk.UniversalEntity,
	entityDef map[string]any,
	createdData map[string]any,
	shouldExist bool,
) {
	t.Helper()

	matchFields := getDefaultTargetFields(entityDef, "list")
	match := map[string]any{}
	for _, field := range matchFields {
		if field != "id" && createdData != nil && createdData[field] != nil {
			match[field] = createdData[field]
		}
	}

	result, err := ent.List(match, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if createdData != nil {
		list := toSlice(result)
		found := vs.Select(list, map[string]any{"id": createdData["id"]})

		if shouldExist {
			if vs.IsEmpty(found) {
				t.Fatalf("Expected created entity %v in list", createdData["id"])
			}
		} else {
			if !vs.IsEmpty(found) {
				t.Fatalf("Expected created entity %v NOT in list after remove", createdData["id"])
			}
		}
	}
}

func testUpdate(
	t *testing.T,
	setup *entityTestSetup,
	ent sdk.UniversalEntity,
	entityName string,
	entityDef map[string]any,
	createdData map[string]any,
) {
	t.Helper()

	reqdata := map[string]any{}
	reqdata["id"] = createdData["id"]

	matchFields := getDefaultTargetFields(entityDef, "update")
	for _, field := range matchFields {
		if field != "id" && createdData[field] != nil {
			reqdata[field] = createdData[field]
		}
	}

	textfield := findTextField(entityDef)
	var markName, markValue string

	if textfield != "" {
		markName = textfield
		markValue = "Mark01-" + entityName + "_ref01_" + strconv.FormatInt(setup.now, 10)
		reqdata[markName] = markValue
	}

	result, err := ent.Update(reqdata, nil)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	resdata := core.ToMapAny(result)
	if resdata == nil {
		t.Fatal("Update result should be a map")
	}
	if resdata["id"] != reqdata["id"] {
		t.Fatalf("Expected id=%v, got %v", reqdata["id"], resdata["id"])
	}

	if markName != "" {
		if resdata[markName] != markValue {
			t.Fatalf("Expected %s=%v, got %v", markName, markValue, resdata[markName])
		}
	}
}

func testLoad(
	t *testing.T,
	setup *entityTestSetup,
	ent sdk.UniversalEntity,
	entityDef map[string]any,
	createdData map[string]any,
) {
	t.Helper()

	matchFields := getDefaultTargetFields(entityDef, "load")
	match := map[string]any{}
	match["id"] = createdData["id"]
	for _, field := range matchFields {
		if field != "id" && createdData[field] != nil {
			match[field] = createdData[field]
		}
	}

	result, err := ent.Load(match, nil)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	resdata := core.ToMapAny(result)
	if resdata == nil {
		t.Fatal("Load result should be a map")
	}
	if resdata["id"] != createdData["id"] {
		t.Fatalf("Expected id=%v, got %v", createdData["id"], resdata["id"])
	}
}

func testRemove(
	t *testing.T,
	setup *entityTestSetup,
	ent sdk.UniversalEntity,
	entityDef map[string]any,
	createdData map[string]any,
) {
	t.Helper()

	matchFields := getDefaultTargetFields(entityDef, "remove")
	match := map[string]any{}
	match["id"] = createdData["id"]
	for _, field := range matchFields {
		if field != "id" && createdData[field] != nil {
			match[field] = createdData[field]
		}
	}

	_, err := ent.Remove(match, nil)
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}
}

func getDefaultTargetFields(entityDef map[string]any, opname string) []string {
	op := core.ToMapAny(vs.GetPath([]any{"op", opname}, entityDef))
	if op == nil {
		return []string{}
	}
	targets, _ := vs.GetProp(op, "targets").([]any)
	if targets == nil {
		return []string{}
	}
	for i := len(targets) - 1; i >= 0; i-- {
		tm := core.ToMapAny(targets[i])
		if tm == nil {
			continue
		}
		sel := core.ToMapAny(vs.GetProp(tm, "select"))
		if sel != nil && vs.GetProp(sel, "$action") != nil {
			continue
		}
		exist, _ := vs.GetProp(sel, "exist").([]any)
		if exist == nil {
			return []string{}
		}
		out := make([]string, 0, len(exist))
		for _, e := range exist {
			if s, ok := e.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return []string{}
}

func findTextField(entityDef map[string]any) string {
	fields, _ := vs.GetProp(entityDef, "fields").([]any)
	for _, f := range fields {
		fm := core.ToMapAny(f)
		if fm == nil {
			continue
		}
		ftype, _ := vs.GetProp(fm, "type").(string)
		fname, _ := vs.GetProp(fm, "name").(string)
		if ftype == "`$STRING`" && fname != "id" && !strings.HasSuffix(fname, "_id") {
			return fname
		}
	}
	return ""
}

func makeEntityTestData(entityDef map[string]any) map[string]any {
	fields, _ := vs.GetProp(entityDef, "fields").([]any)
	if fields == nil {
		fields = []any{}
	}
	name, _ := vs.GetProp(entityDef, "name").(string)

	existing := map[string]any{}
	newData := map[string]any{}

	idcount := 3
	refs := make([]string, idcount)
	idmapLocal := map[string]string{}
	for i := 0; i < idcount; i++ {
		ref := fmt.Sprintf("%s%02d", name, i)
		refs[i] = ref
		idmapLocal[ref] = strings.ToUpper(ref)
	}

	idx := 1
	for _, ref := range refs {
		id := idmapLocal[ref]
		ent := map[string]any{}
		makeEntityTestFields(fields, idx, ent)
		idx++
		ent["id"] = id
		existing[id] = ent
	}

	newRef := name + "_ref01"
	newEnt := map[string]any{}
	makeEntityTestFields(fields, idx, newEnt)
	delete(newEnt, "id")
	newData[newRef] = newEnt

	return map[string]any{
		"existing": map[string]any{name: existing},
		"new":      map[string]any{name: newData},
	}
}

func makeEntityTestFields(fields []any, start int, entdata map[string]any) {
	num := start * len(fields) * 10
	for _, f := range fields {
		fm := core.ToMapAny(f)
		if fm == nil {
			continue
		}
		fname, _ := vs.GetProp(fm, "name").(string)
		ftype, _ := vs.GetProp(fm, "type").(string)

		if strings.HasSuffix(fname, "_id") {
			entdata[fname] = strings.ToUpper(fname[:len(fname)-3]) + "01"
		} else if ftype == "`$NUMBER`" {
			entdata[fname] = num
		} else if ftype == "`$BOOLEAN`" {
			entdata[fname] = num%2 == 0
		} else if ftype == "`$OBJECT`" || ftype == "`$MAP`" {
			entdata[fname] = map[string]any{}
		} else if ftype == "`$ARRAY`" || ftype == "`$LIST`" {
			entdata[fname] = []any{}
		} else {
			entdata[fname] = "s" + strconv.FormatInt(int64(num), 16)
		}
		num++
	}
}

type entityTestSetup struct {
	idmap   map[string]any
	env     map[string]string
	client  *sdk.UniversalSDK
	data    map[string]any
	explain bool
	now     int64
}

func basicSetup(
	t *testing.T,
	um *sdk.UniversalManager,
	entityMap map[string]any,
	entityName string,
) *entityTestSetup {
	t.Helper()

	allExisting := map[string]any{}
	allNew := map[string]any{}

	for _, item := range vs.Items(entityMap) {
		entDef := core.ToMapAny(item[1])
		testData := makeEntityTestData(entDef)

		existMap := core.ToMapAny(vs.GetProp(testData, "existing"))
		for k, v := range existMap {
			allExisting[k] = v
		}
		newMap := core.ToMapAny(vs.GetProp(testData, "new"))
		for k, v := range newMap {
			allNew[k] = v
		}
	}

	options := map[string]any{
		"entity": allExisting,
	}

	sdkInst := um.Make("voxgig-solardemo")
	client := sdkInst.Test(options, nil)

	// Build idmap
	idEntries := []any{}
	for _, item := range vs.Items(entityMap) {
		entDef := core.ToMapAny(item[1])
		ename, _ := vs.GetProp(entDef, "name").(string)
		for i := 0; i < 3; i++ {
			idEntries = append(idEntries, fmt.Sprintf("%s%02d", ename, i))
		}
	}

	idmap := core.ToMapAny(vs.Transform(
		idEntries,
		map[string]any{
			"`$PACK`": []any{"", map[string]any{
				"`$KEY`": "`$COPY`",
				"`$VAL`": []any{"`$FORMAT`", "upper", "`$COPY`"},
			}},
		},
	))
	if idmap == nil {
		idmap = map[string]any{}
	}

	env := envOverride(map[string]string{
		"UNIVERSAL_TEST_LIVE":    "FALSE",
		"UNIVERSAL_TEST_EXPLAIN": "FALSE",
		"UNIVERSAL_APIKEY":       "NONE",
	})

	if env["UNIVERSAL_TEST_LIVE"] == "TRUE" {
		liveopts := map[string]any{
			"ref":    "voxgig-solardemo",
			"model":  um.ResolveModel("voxgig-solardemo"),
			"apikey": env["UNIVERSAL_APIKEY"],
		}
		client = sdk.NewUniversalSDK(um, liveopts)

		// Discover real parent entity IDs from the live API.
		for _, item := range vs.Items(entityMap) {
			eDef := core.ToMapAny(item[1])
			eName, _ := vs.GetProp(eDef, "name").(string)
			listOp := core.ToMapAny(vs.GetPath([]any{"op", "list"}, eDef))
			if listOp == nil {
				continue
			}
			targets, _ := vs.GetProp(listOp, "targets").([]any)
			if len(targets) == 0 {
				continue
			}
			listTarget := core.ToMapAny(targets[0])
			if listTarget == nil {
				continue
			}
			listParams, _ := vs.GetPath([]any{"args", "params"}, listTarget).([]any)
			if len(listParams) > 0 {
				continue // skip nested entities
			}
			parts, _ := vs.GetProp(listTarget, "parts").([]any)
			listPath := vs.Join(parts, "/", true)

			res, err := client.Direct(map[string]any{
				"path": listPath, "method": "GET", "params": map[string]any{},
			})
			if err == nil && res["ok"] == true {
				if dataList, ok := res["data"].([]any); ok {
					for i := 0; i < len(dataList) && i < 3; i++ {
						if dm := core.ToMapAny(dataList[i]); dm != nil {
							ref := fmt.Sprintf("%s%02d", eName, i)
							idmap[ref] = dm["id"]
						}
					}
				}
			}
		}
	}

	return &entityTestSetup{
		idmap:   idmap,
		env:     env,
		client:  client,
		data:    map[string]any{"existing": allExisting, "new": allNew},
		explain: env["UNIVERSAL_TEST_EXPLAIN"] == "TRUE",
		now:     time.Now().UnixMilli(),
	}
}

func envOverride(m map[string]string) map[string]string {
	live := os.Getenv("UNIVERSAL_TEST_LIVE")
	override := os.Getenv("UNIVERSAL_TEST_OVERRIDE")

	if live == "TRUE" || override == "TRUE" {
		for k := range m {
			if envval := os.Getenv(k); envval != "" {
				m[k] = strings.TrimSpace(envval)
			}
		}
	}

	if explain := os.Getenv("UNIVERSAL_TEST_EXPLAIN"); explain != "" {
		m["UNIVERSAL_TEST_EXPLAIN"] = explain
	}

	return m
}

// toSlice converts the result from List (which may be []any or a
// wrapped entity list) into a plain []any of maps.
func toSlice(val any) []any {
	if val == nil {
		return []any{}
	}
	if list, ok := val.([]any); ok {
		out := make([]any, 0, len(list))
		for _, item := range list {
			if ent, ok := item.(sdk.Entity); ok {
				out = append(out, ent.Data())
			} else {
				out = append(out, item)
			}
		}
		return out
	}
	return []any{}
}
