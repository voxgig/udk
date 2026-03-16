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

func (um *UniversalManager) Make(ref string) *UniversalSDK {
	model := um.ResolveModel(ref)
	sdkopts := map[string]any{
		"ref":   ref,
		"model": model,
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
