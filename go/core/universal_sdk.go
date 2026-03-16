package core

import (
	vs "github.com/voxgig/struct"
)

type UniversalSDK struct {
	Mode      string
	options   map[string]any
	initopts  map[string]any
	utility   *Utility
	Features  []Feature
	rootctx   *Context
	um        *UniversalManager
}

func NewUniversalSDK(um *UniversalManager, options map[string]any) *UniversalSDK {
	sdk := &UniversalSDK{
		Mode:     "live",
		Features: []Feature{},
		um:       um,
		initopts: options,
	}

	sdk.utility = NewUtility()

	config := makeConfig(options)

	sdk.rootctx = sdk.utility.MakeContext(map[string]any{
		"client":  sdk,
		"utility": sdk.utility,
		"config":  config,
		"options": options,
		"shared":  map[string]any{},
	}, nil)

	sdk.options = sdk.utility.MakeOptions(sdk.rootctx)

	if vs.GetPath([]any{"feature", "test", "active"}, sdk.options) == true {
		sdk.Mode = "test"
	}

	sdk.rootctx.Options = sdk.options

	// Add features from config.
	featureOpts := ToMapAny(vs.GetProp(sdk.options, "feature"))
	if featureOpts != nil {
		for _, item := range vs.Items(featureOpts) {
			fname, _ := item[0].(string)
			fopts := ToMapAny(item[1])
			if fopts != nil {
				if active, ok := fopts["active"]; ok {
					if ab, ok := active.(bool); ok && ab {
						sdk.utility.FeatureAdd(sdk.rootctx, MakeFeature(fname))
					}
				}
			}
		}
	}

	// Add extension features.
	if extend := vs.GetProp(sdk.options, "extend"); extend != nil {
		if extList, ok := extend.([]any); ok {
			for _, f := range extList {
				if feat, ok := f.(Feature); ok {
					sdk.utility.FeatureAdd(sdk.rootctx, feat)
				}
			}
		}
	}

	// Initialize features.
	for _, f := range sdk.Features {
		sdk.utility.FeatureInit(sdk.rootctx, f)
	}

	sdk.utility.FeatureHook(sdk.rootctx, "PostConstruct")

	return sdk
}

func (sdk *UniversalSDK) OptionsMap() map[string]any {
	out := vs.Clone(sdk.options)
	if om, ok := out.(map[string]any); ok {
		return om
	}
	return map[string]any{}
}

func (sdk *UniversalSDK) GetUtility() *Utility {
	return CopyUtility(sdk.utility)
}

func (sdk *UniversalSDK) GetRootCtx() *Context {
	return sdk.rootctx
}

func (sdk *UniversalSDK) Prepare(fetchargs map[string]any) (map[string]any, error) {
	utility := sdk.utility

	if fetchargs == nil {
		fetchargs = map[string]any{}
	}

	var ctrl map[string]any
	if c := vs.GetProp(fetchargs, "ctrl"); c != nil {
		if cm, ok := c.(map[string]any); ok {
			ctrl = cm
		}
	}
	if ctrl == nil {
		ctrl = map[string]any{}
	}

	ctx := utility.MakeContext(map[string]any{
		"opname": "prepare",
		"ctrl":   ctrl,
	}, sdk.rootctx)

	options := sdk.options

	path, _ := vs.GetProp(fetchargs, "path").(string)
	method, _ := vs.GetProp(fetchargs, "method").(string)
	if method == "" {
		method = "GET"
	}

	params := ToMapAny(vs.GetProp(fetchargs, "params"))
	if params == nil {
		params = map[string]any{}
	}
	query := ToMapAny(vs.GetProp(fetchargs, "query"))
	if query == nil {
		query = map[string]any{}
	}

	headers := utility.PrepareHeaders(ctx)

	base, _ := vs.GetProp(options, "base").(string)
	prefix, _ := vs.GetProp(options, "prefix").(string)
	suffix, _ := vs.GetProp(options, "suffix").(string)

	ctx.Spec = NewSpec(map[string]any{
		"base":    base,
		"prefix":  prefix,
		"suffix":  suffix,
		"path":    path,
		"method":  method,
		"params":  params,
		"query":   query,
		"headers": headers,
		"body":    vs.GetProp(fetchargs, "body"),
		"step":    "start",
	})

	// Merge user-provided headers.
	if uh := vs.GetProp(fetchargs, "headers"); uh != nil {
		if uhm, ok := uh.(map[string]any); ok {
			for k, v := range uhm {
				ctx.Spec.Headers[k] = v
			}
		}
	}

	_, err := utility.PrepareAuth(ctx)
	if err != nil {
		return nil, err
	}

	return utility.MakeFetchDef(ctx)
}

func (sdk *UniversalSDK) Direct(fetchargs map[string]any) (map[string]any, error) {
	utility := sdk.utility

	fetchdef, err := sdk.Prepare(fetchargs)
	if err != nil {
		return map[string]any{"ok": false, "err": err}, nil
	}

	if fetchargs == nil {
		fetchargs = map[string]any{}
	}

	var ctrl map[string]any
	if c := vs.GetProp(fetchargs, "ctrl"); c != nil {
		if cm, ok := c.(map[string]any); ok {
			ctrl = cm
		}
	}
	if ctrl == nil {
		ctrl = map[string]any{}
	}

	ctx := utility.MakeContext(map[string]any{
		"opname": "direct",
		"ctrl":   ctrl,
	}, sdk.rootctx)

	url, _ := fetchdef["url"].(string)
	fetched, fetchErr := utility.Fetcher(ctx, url, fetchdef)

	if fetchErr != nil {
		return map[string]any{"ok": false, "err": fetchErr}, nil
	}

	if fetched == nil {
		return map[string]any{
			"ok":  false,
			"err": ctx.MakeError("direct_no_response", "response: undefined"),
		}, nil
	}

	if fm, ok := fetched.(map[string]any); ok {
		status := ToInt(vs.GetProp(fm, "status"))
		var jsonData any
		if jf := vs.GetProp(fm, "json"); jf != nil {
			if f, ok := jf.(func() any); ok {
				jsonData = f()
			}
		}

		return map[string]any{
			"ok":      status >= 200 && status < 300,
			"status":  status,
			"headers": vs.GetProp(fm, "headers"),
			"data":    jsonData,
		}, nil
	}

	return map[string]any{"ok": false, "err": ctx.MakeError("direct_invalid", "invalid response type")}, nil
}

func (sdk *UniversalSDK) Entity(name string, entopts map[string]any) UniversalEntity {
	return NewUniversalEntityFunc(sdk, name, entopts)
}

func (sdk *UniversalSDK) Test(testopts map[string]any, sdkopts map[string]any) *UniversalSDK {
	if sdkopts == nil {
		sdkopts = map[string]any{}
	}
	sdkopts = vs.Clone(sdkopts).(map[string]any)

	if testopts == nil {
		testopts = map[string]any{}
	}
	testopts = vs.Clone(testopts).(map[string]any)
	testopts["active"] = true

	vs.SetPath(sdkopts, []any{"feature", "test"}, testopts)

	// Carry over model/ref from original SDK so entity config is available.
	if sdkopts["model"] == nil && sdk.initopts != nil {
		if model := vs.GetProp(sdk.initopts, "model"); model != nil {
			sdkopts["model"] = model
		}
	}
	if sdkopts["ref"] == nil && sdk.initopts != nil {
		if ref := vs.GetProp(sdk.initopts, "ref"); ref != nil {
			sdkopts["ref"] = ref
		}
	}

	testSDK := NewUniversalSDK(sdk.um, sdkopts)
	testSDK.Mode = "test"

	return testSDK
}

func (sdk *UniversalSDK) Tester(testopts map[string]any, sdkopts map[string]any) *UniversalSDK {
	return sdk.Test(testopts, sdkopts)
}

// makeConfig builds the SDK config by merging the base config with
// entity definitions from the model.
func makeConfig(options map[string]any) map[string]any {
	config := MakeConfig()

	model := ToMapAny(vs.GetProp(options, "model"))
	if model != nil {
		entity := vs.GetPath([]any{"main", "kit", "entity"}, model)
		if entity != nil {
			if em, ok := entity.(map[string]any); ok {
				config["entity"] = em
			}
		}
	}

	return config
}

// NewUniversalEntityFunc is set by the entity package via init().
var NewUniversalEntityFunc func(client *UniversalSDK, name string, entopts map[string]any) UniversalEntity
