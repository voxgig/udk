package utility

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	sdk "voxgiguniversalsdk"

	vs "github.com/voxgig/struct"
)

var cachedTestSpec map[string]any

func loadTestSpec(t *testing.T) map[string]any {
	t.Helper()
	if cachedTestSpec != nil {
		return cachedTestSpec
	}
	data, err := os.ReadFile("../../../.sdk/test/test.json")
	if err != nil {
		t.Fatalf("Failed to load test.json: %v", err)
	}
	var spec map[string]any
	if err := json.Unmarshal(data, &spec); err != nil {
		t.Fatalf("Failed to parse test.json: %v", err)
	}
	cachedTestSpec = spec
	return spec
}

func getSpec(spec map[string]any, keys ...string) map[string]any {
	var cur any = spec
	for _, key := range keys {
		if m, ok := cur.(map[string]any); ok {
			cur = m[key]
		} else {
			return nil
		}
	}
	if m, ok := cur.(map[string]any); ok {
		return m
	}
	return nil
}

type RunSubject func(entry map[string]any) (any, error)

func runset(t *testing.T, testspec map[string]any, subject RunSubject) {
	t.Helper()
	set, ok := testspec["set"].([]any)
	if !ok {
		return
	}

	for i, e := range set {
		entry, ok := e.(map[string]any)
		if !ok {
			continue
		}

		mark := ""
		if m := entry["mark"]; m != nil {
			mark = fmt.Sprintf(" (mark=%v)", m)
		}

		result, err := subject(entry)

		expectedErr := entry["err"]

		if err != nil {
			if expectedErr != nil {
				errMsg := err.Error()
				if expStr, ok := expectedErr.(string); ok {
					if !strings.Contains(errMsg, expStr) {
						t.Errorf("entry %d%s: error mismatch: got %q, want contains %q",
							i, mark, errMsg, expStr)
					}
				} else if expBool, ok := expectedErr.(bool); ok && expBool {
					// err: true means any error is acceptable
				}
				if matchSpec, ok := entry["match"].(map[string]any); ok {
					resultMap := map[string]any{
						"in":  entry["in"],
						"out": jsonNormalize(result),
						"err": map[string]any{"message": err.Error()},
					}
					matchDeep(t, i, mark, matchSpec, resultMap, "")
				}
				continue
			}
			t.Errorf("entry %d%s: unexpected error: %v", i, mark, err)
			continue
		}

		if expectedErr != nil {
			t.Errorf("entry %d%s: expected error containing %q but got result: %v",
				i, mark, expectedErr, jsonStr(result))
			continue
		}

		matched := false
		if matchSpec, ok := entry["match"].(map[string]any); ok {
			resultMap := map[string]any{
				"in":  entry["in"],
				"out": jsonNormalize(result),
			}
			if args := entry["args"]; args != nil {
				resultMap["args"] = args
			} else if entry["in"] != nil {
				resultMap["args"] = []any{entry["in"]}
			}
			if ctxData := entry["ctx"]; ctxData != nil {
				resultMap["ctx"] = ctxData
			}
			matchDeep(t, i, mark, matchSpec, resultMap, "")
			matched = true
		}

		expectedOut := entry["out"]
		if expectedOut == nil && matched {
			continue
		}
		if expectedOut != nil {
			normResult := jsonNormalize(result)
			normExpected := jsonNormalize(expectedOut)
			if !reflect.DeepEqual(normResult, normExpected) {
				t.Errorf("entry %d%s: output mismatch:\n  got:  %v\n  want: %v",
					i, mark, jsonStr(normResult), jsonStr(normExpected))
			}
		}
	}
}

func jsonNormalize(val any) any {
	if val == nil {
		return nil
	}
	j, err := json.Marshal(val)
	if err != nil {
		return val
	}
	var out any
	json.Unmarshal(j, &out)
	return out
}

func jsonStr(val any) string {
	j, err := json.Marshal(val)
	if err != nil {
		return fmt.Sprintf("%v", val)
	}
	return string(j)
}

func matchDeep(t *testing.T, entryIdx int, mark string, check any, base any, path string) {
	t.Helper()

	if check == nil {
		return
	}

	checkMap, isMap := check.(map[string]any)
	checkList, isList := check.([]any)

	if isMap {
		for key, checkVal := range checkMap {
			childPath := path + "." + key
			var baseVal any
			if baseMap, ok := base.(map[string]any); ok {
				baseVal = baseMap[key]
			}
			matchDeep(t, entryIdx, mark, checkVal, baseVal, childPath)
		}
	} else if isList {
		for i, checkVal := range checkList {
			childPath := fmt.Sprintf("%s[%d]", path, i)
			var baseVal any
			if baseList, ok := base.([]any); ok && i < len(baseList) {
				baseVal = baseList[i]
			}
			matchDeep(t, entryIdx, mark, checkVal, baseVal, childPath)
		}
	} else {
		checkStr, isStr := check.(string)
		if isStr && checkStr == "__EXISTS__" {
			if base == nil {
				t.Errorf("entry %d%s: match %s: expected value to exist but got nil",
					entryIdx, mark, path)
			}
			return
		}
		if isStr && checkStr == "__UNDEF__" {
			if base != nil {
				t.Errorf("entry %d%s: match %s: expected nil but got %v",
					entryIdx, mark, path, base)
			}
			return
		}

		normCheck := jsonNormalize(check)
		normBase := jsonNormalize(base)

		if !reflect.DeepEqual(normCheck, normBase) {
			if isStr && checkStr != "" {
				baseStr := vs.Stringify(base)
				if strings.Contains(strings.ToLower(baseStr), strings.ToLower(checkStr)) {
					return
				}
			}
			t.Errorf("entry %d%s: match %s: got %v, want %v",
				entryIdx, mark, path, jsonStr(normBase), jsonStr(normCheck))
		}
	}
}

// makeCtxFromMap creates a Context from a JSON test entry's ctx or args map.
func makeCtxFromMap(ctxmap map[string]any, client *sdk.UniversalSDK, utility *sdk.Utility) *sdk.Context {
	if ctxmap == nil {
		ctxmap = map[string]any{}
	}

	// Extract opname from op map if present.
	if opMap, ok := ctxmap["op"].(map[string]any); ok {
		if opname, ok := opMap["name"].(string); ok {
			ctxmap["opname"] = opname
		}
	}

	ctx := sdk.NewContext(ctxmap, nil)

	if client != nil {
		ctx.Client = client
		ctx.Utility = utility
	}
	if ctx.Options == nil && client != nil {
		ctx.Options = client.OptionsMap()
	}

	// Handle op fields from JSON map (targets, alias, etc.)
	if opMap, ok := ctxmap["op"].(map[string]any); ok {
		ctx.Op = sdk.NewOperation(opMap)
	}

	// Handle spec from JSON map
	if specMap, ok := ctxmap["spec"].(map[string]any); ok {
		ctx.Spec = sdk.NewSpec(specMap)
	}

	// Handle result from JSON map
	if resMap, ok := ctxmap["result"].(map[string]any); ok {
		ctx.Result = sdk.NewResult(resMap)
		if errMap, ok := resMap["err"].(map[string]any); ok {
			if msg, ok := errMap["message"].(string); ok {
				ctx.Result.Err = &sdk.UniversalError{Msg: msg}
			}
		}
	}

	// Handle response from JSON map
	if respMap, ok := ctxmap["response"].(map[string]any); ok {
		// Check for native wrapper (test.json format)
		if nativeMap, ok := respMap["native"].(map[string]any); ok {
			// Map native fields: reason -> statusText
			mappedResp := map[string]any{
				"status":     nativeMap["status"],
				"statusText": nativeMap["reason"],
			}
			if body := nativeMap["body"]; body != nil {
				mappedResp["body"] = body
			}
			if headers, ok := nativeMap["headers"].(map[string]any); ok {
				lowerHeaders := map[string]any{}
				for k, v := range headers {
					lowerHeaders[strings.ToLower(k)] = v
				}
				mappedResp["headers"] = lowerHeaders
			}
			ctx.Response = sdk.NewResponse(mappedResp)
			if body := nativeMap["body"]; body != nil {
				bodyCopy := body
				ctx.Response.JsonFunc = func() any { return bodyCopy }
				ctx.Response.Body = bodyCopy
			}
		} else {
			ctx.Response = sdk.NewResponse(respMap)
			if body := respMap["body"]; body != nil {
				bodyCopy := body
				ctx.Response.JsonFunc = func() any { return bodyCopy }
			}
			if headers, ok := respMap["headers"].(map[string]any); ok {
				lowerHeaders := map[string]any{}
				for k, v := range headers {
					lowerHeaders[strings.ToLower(k)] = v
				}
				ctx.Response.Headers = lowerHeaders
			}
		}
	}

	return ctx
}

func fixctx(ctx *sdk.Context, client *sdk.UniversalSDK) {
	if ctx != nil && ctx.Client != nil && ctx.Options == nil {
		ctx.Options = ctx.Client.OptionsMap()
	}
}

// entityListToData extracts data maps from a list of Entity objects.
func entityListToData(list []any) []any {
	var out []any
	for _, item := range list {
		if ent, ok := item.(sdk.Entity); ok {
			d := ent.Data()
			if dm, ok := d.(map[string]any); ok {
				out = append(out, dm)
			}
		} else if m, ok := item.(map[string]any); ok {
			out = append(out, m)
		}
	}
	if out == nil {
		out = []any{}
	}
	return out
}
