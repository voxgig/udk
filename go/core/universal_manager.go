package core

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type UniversalManager struct {
	options map[string]any
	utility *Utility
}

func NewUniversalManager(options map[string]any) *UniversalManager {
	return &UniversalManager{
		options: options,
		utility: NewUtility(),
	}
}

func (um *UniversalManager) Options() map[string]any {
	return um.options
}

func (um *UniversalManager) Make(ref string, opts ...map[string]any) *UniversalSDK {
	sdkopts := map[string]any{
		"ref": ref,
	}

	// Merge caller-supplied options (including an optional "model" key).
	for _, o := range opts {
		for k, v := range o {
			sdkopts[k] = v
		}
	}

	// Only resolve from registry when no model was provided.
	if sdkopts["model"] == nil {
		sdkopts["model"] = um.ResolveModel(ref)
	}

	return NewUniversalSDK(um, sdkopts)
}

func (um *UniversalManager) ResolveModel(ref string) map[string]any {
	registry, _ := um.options["registry"].(string)
	modelpath := filepath.Join(registry, "local", ref+".json")

	data, err := os.ReadFile(modelpath)
	if err != nil {
		return map[string]any{}
	}

	var model map[string]any
	if err := json.Unmarshal(data, &model); err != nil {
		return map[string]any{}
	}

	return model
}
