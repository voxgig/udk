package universal

import (
	"fmt"
	"os"
	"strings"
	"testing"

	sdk "voxgiguniversalsdk"
	"voxgiguniversalsdk/core"

	vs "github.com/voxgig/struct"
)

func TestUniversalDirect(t *testing.T) {
	um := sdk.NewUniversalManager(map[string]any{
		"registry": "../../../test/registry",
	})
	baseSDK := um.Make("voxgig-solardemo")
	entityMap := core.ToMapAny(vs.GetPath([]any{"entity"}, baseSDK.GetRootCtx().Config))

	live := os.Getenv("UNIVERSAL_TEST_LIVE") == "TRUE"

	t.Run("direct-exists", func(t *testing.T) {
		setup := directSetup(um, baseSDK, nil)
		if setup.client == nil {
			t.Fatal("client should not be nil")
		}
	})

	for _, item := range vs.Items(entityMap) {
		entDef := core.ToMapAny(item[1])
		if entDef == nil {
			continue
		}
		entityName, _ := vs.GetProp(entDef, "name").(string)
		ops := core.ToMapAny(vs.GetProp(entDef, "op"))
		if ops == nil {
			continue
		}

		hasLoad := vs.GetProp(ops, "load") != nil
		hasList := vs.GetProp(ops, "list") != nil

		if !hasLoad && !hasList {
			continue
		}

		if hasLoad {
			loadOp := core.ToMapAny(vs.GetProp(ops, "load"))
			targets, _ := vs.GetProp(loadOp, "targets").([]any)
			if len(targets) == 0 {
				continue
			}
			loadTarget := core.ToMapAny(targets[0])
			if loadTarget == nil {
				continue
			}

			name := entityName
			t.Run("direct-load-"+name, func(t *testing.T) {
				parts, _ := vs.GetProp(loadTarget, "parts").([]any)
				loadPath := vs.Join(parts, "/", true)
				loadParams, _ := vs.GetPath([]any{"args", "params"}, loadTarget).([]any)

				if live {
					idmap := resolveIdmap(um, baseSDK, entityMap)
					setup := directSetup(um, baseSDK, nil)

					if hasList {
						listOp := core.ToMapAny(vs.GetProp(ops, "list"))
						listTargets, _ := vs.GetProp(listOp, "targets").([]any)
						if len(listTargets) > 0 {
							listTarget := core.ToMapAny(listTargets[0])
							if listTarget != nil {
								listParts, _ := vs.GetProp(listTarget, "parts").([]any)
								listPath := vs.Join(listParts, "/", true)
								listParams, _ := vs.GetPath([]any{"args", "params"}, listTarget).([]any)

								// Try multiple parent refs to find one with child entities.
								var found map[string]any
								lparams := map[string]any{}
								for tr := 0; tr < 3 && found == nil; tr++ {
									lparams = map[string]any{}
									for _, p := range listParams {
										pm := core.ToMapAny(p)
										if pm == nil {
											continue
										}
										pname, _ := vs.GetProp(pm, "name").(string)
										ref := strings.TrimSuffix(pname, "_id") +
											fmt.Sprintf("%02d", tr)
										if idmap[ref] != nil {
											lparams[pname] = idmap[ref]
										} else {
											lparams[pname] = ref
										}
									}

									listResult, err := setup.client.Direct(map[string]any{
										"path": listPath, "method": "GET", "params": lparams,
									})
									if err != nil {
										t.Fatalf("direct list failed: %v", err)
									}
									if listResult["ok"] != true {
										t.Fatalf("expected ok=true, got %v", listResult["ok"])
									}
									if dataList, ok := listResult["data"].([]any); ok && len(dataList) >= 1 {
										found = core.ToMapAny(dataList[0])
									}
								}

								if found != nil {
									params := map[string]any{}
									for _, p := range loadParams {
										pm := core.ToMapAny(p)
										if pm == nil {
											continue
										}
										pname, _ := vs.GetProp(pm, "name").(string)
										if found[pname] != nil {
											params[pname] = found[pname]
										} else if lparams[pname] != nil {
											params[pname] = lparams[pname]
										}
									}

									result, err := setup.client.Direct(map[string]any{
										"path": loadPath, "method": "GET", "params": params,
									})
									if err != nil {
										t.Fatalf("direct load failed: %v", err)
									}
									if result["ok"] != true {
										t.Fatalf("expected ok=true, got %v", result["ok"])
									}
									if core.ToInt(result["status"]) != 200 {
										t.Fatalf("expected status 200, got %v", result["status"])
									}
									if result["data"] == nil {
										t.Fatal("expected data to be non-nil")
									}
									if dm := core.ToMapAny(result["data"]); dm != nil {
										if dm["id"] != found["id"] {
											t.Fatalf("expected id=%v, got %v", found["id"], dm["id"])
										}
									}
								}
							}
						}
					}
				} else {
					setup := directSetup(um, baseSDK, map[string]any{"id": "direct01"})

					params := map[string]any{}
					for i, p := range loadParams {
						pm := core.ToMapAny(p)
						if pm != nil {
							pname, _ := vs.GetProp(pm, "name").(string)
							params[pname] = fmt.Sprintf("direct0%d", i+1)
						}
					}

					result, err := setup.client.Direct(map[string]any{
						"path": loadPath, "method": "GET", "params": params,
					})
					if err != nil {
						t.Fatalf("direct failed: %v", err)
					}
					if result["ok"] != true {
						t.Fatalf("expected ok=true, got %v", result["ok"])
					}
					if core.ToInt(result["status"]) != 200 {
						t.Fatalf("expected status 200, got %v", result["status"])
					}
					if result["data"] == nil {
						t.Fatal("expected data to be non-nil")
					}
					if dm := core.ToMapAny(result["data"]); dm != nil {
						if dm["id"] != "direct01" {
							t.Fatalf("expected data.id=direct01, got %v", dm["id"])
						}
					}

					if len(*setup.calls) != 1 {
						t.Fatalf("expected 1 call, got %d", len(*setup.calls))
					}
					call := (*setup.calls)[0]
					if initMap, ok := call["init"].(map[string]any); ok {
						if initMap["method"] != "GET" {
							t.Fatalf("expected method GET, got %v", initMap["method"])
						}
					}
					for i := range loadParams {
						url, _ := call["url"].(string)
						if !strings.Contains(url, fmt.Sprintf("direct0%d", i+1)) {
							t.Fatalf("expected url to contain direct0%d, got %v", i+1, url)
						}
					}
				}
			})
		}

		if hasList {
			listOp := core.ToMapAny(vs.GetProp(ops, "list"))
			targets, _ := vs.GetProp(listOp, "targets").([]any)
			if len(targets) == 0 {
				continue
			}
			listTarget := core.ToMapAny(targets[0])
			if listTarget == nil {
				continue
			}

			name := entityName
			t.Run("direct-list-"+name, func(t *testing.T) {
				parts, _ := vs.GetProp(listTarget, "parts").([]any)
				listPath := vs.Join(parts, "/", true)
				listParams, _ := vs.GetPath([]any{"args", "params"}, listTarget).([]any)

				if live {
					idmap := resolveIdmap(um, baseSDK, entityMap)
					setup := directSetup(um, baseSDK, nil)

					// For entities with parent params, try each known parent
					// to find one that has child entities.
					found := false
					maxTries := 1
					if len(listParams) > 0 {
						maxTries = 3
					}
					for tr := 0; tr < maxTries && !found; tr++ {
						params := map[string]any{}
						for _, p := range listParams {
							pm := core.ToMapAny(p)
							if pm == nil {
								continue
							}
							pname, _ := vs.GetProp(pm, "name").(string)
							base := pname
							if pname == "id" {
								base = entityName
							} else {
								base = strings.TrimSuffix(pname, "_id")
							}
							ref := fmt.Sprintf("%s%02d", base, tr)
							if idmap[ref] != nil {
								params[pname] = idmap[ref]
							} else {
								params[pname] = ref
							}
						}

						result, err := setup.client.Direct(map[string]any{
							"path": listPath, "method": "GET", "params": params,
						})
						if err != nil {
							t.Fatalf("direct failed: %v", err)
						}
						if result["ok"] != true {
							t.Fatalf("expected ok=true, got %v", result["ok"])
						}
						if core.ToInt(result["status"]) != 200 {
							t.Fatalf("expected status 200, got %v", result["status"])
						}
						if dataList, ok := result["data"].([]any); ok && len(dataList) >= 1 {
							found = true
						}
					}

					if len(listParams) == 0 && !found {
						t.Fatal("expected at least one entity in list")
					}
				} else {
					mockData := []any{
						map[string]any{"id": "direct01"},
						map[string]any{"id": "direct02"},
					}
					setup := directSetup(um, baseSDK, mockData)

					params := map[string]any{}
					for i, p := range listParams {
						pm := core.ToMapAny(p)
						if pm != nil {
							pname, _ := vs.GetProp(pm, "name").(string)
							params[pname] = fmt.Sprintf("direct0%d", i+1)
						}
					}

					result, err := setup.client.Direct(map[string]any{
						"path": listPath, "method": "GET", "params": params,
					})
					if err != nil {
						t.Fatalf("direct failed: %v", err)
					}
					if result["ok"] != true {
						t.Fatalf("expected ok=true, got %v", result["ok"])
					}
					if core.ToInt(result["status"]) != 200 {
						t.Fatalf("expected status 200, got %v", result["status"])
					}
					if dataList, ok := result["data"].([]any); ok {
						if len(dataList) != 2 {
							t.Fatalf("expected 2 items, got %d", len(dataList))
						}
					} else {
						t.Fatalf("expected data to be an array, got %T", result["data"])
					}

					if len(*setup.calls) != 1 {
						t.Fatalf("expected 1 call, got %d", len(*setup.calls))
					}
					call := (*setup.calls)[0]
					if initMap, ok := call["init"].(map[string]any); ok {
						if initMap["method"] != "GET" {
							t.Fatalf("expected method GET, got %v", initMap["method"])
						}
					}

					for i := range listParams {
						url, _ := call["url"].(string)
						if !strings.Contains(url, fmt.Sprintf("direct0%d", i+1)) {
							t.Fatalf("expected url to contain direct0%d, got %v", i+1, url)
						}
					}
				}
			})
		}
	}
}

type directSetupResult struct {
	client *sdk.UniversalSDK
	calls  *[]map[string]any
	live   bool
}

func directSetup(um *sdk.UniversalManager, baseSDK *sdk.UniversalSDK, mockres any) *directSetupResult {
	live := os.Getenv("UNIVERSAL_TEST_LIVE") == "TRUE"

	if live {
		client := sdk.NewUniversalSDK(um, map[string]any{
			"ref":   "voxgig-solardemo",
			"model": um.ResolveModel("voxgig-solardemo"),
		})
		return &directSetupResult{client: client, calls: &[]map[string]any{}, live: true}
	}

	calls := &[]map[string]any{}

	mockFetch := func(url string, init map[string]any) (map[string]any, error) {
		*calls = append(*calls, map[string]any{"url": url, "init": init})
		return map[string]any{
			"status":     200,
			"statusText": "OK",
			"headers":    map[string]any{},
			"json": (func() any)(func() any {
				if mockres != nil {
					return mockres
				}
				return map[string]any{"id": "direct01"}
			}),
		}, nil
	}

	client := core.NewUniversalSDK(um, map[string]any{
		"model": um.ResolveModel("voxgig-solardemo"),
		"ref":   "voxgig-solardemo",
		"base":  "http://localhost:8080",
		"system": map[string]any{
			"fetch": (func(string, map[string]any) (map[string]any, error))(mockFetch),
		},
	})

	return &directSetupResult{client: client, calls: calls, live: false}
}

func resolveIdmap(um *sdk.UniversalManager, baseSDK *sdk.UniversalSDK, entityMap map[string]any) map[string]any {
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

	// In live mode, discover real parent entity IDs.
	if os.Getenv("UNIVERSAL_TEST_LIVE") == "TRUE" {
		liveClient := sdk.NewUniversalSDK(um, map[string]any{
			"ref":   "voxgig-solardemo",
			"model": um.ResolveModel("voxgig-solardemo"),
		})

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

			res, err := liveClient.Direct(map[string]any{
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

	return idmap
}
