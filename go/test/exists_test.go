package test

import (
	"testing"

	sdk "voxgiguniversalsdk"

	vs "github.com/voxgig/struct"
)

func TestUniversalManagerExists(t *testing.T) {
	um := sdk.NewUniversalManager(map[string]any{
		"registry": "./registry",
	})

	if um == nil {
		t.Fatal("UniversalManager should not be nil")
	}

	udk := um.Make("voxgig-solardemo")
	if udk == nil {
		t.Fatal("UniversalSDK should not be nil")
	}
}

func TestUniversalManagerModelMap(t *testing.T) {
	// Create a manager with no registry path.
	um := sdk.NewUniversalManager(map[string]any{})

	if um == nil {
		t.Fatal("UniversalManager should not be nil")
	}

	// Provide the API model inline as a map.
	model := map[string]any{
		"main": map[string]any{
			"kit": map[string]any{
				"entity": map[string]any{
					"foo": map[string]any{
						"name":   "foo",
						"active": true,
						"fields": []any{
							map[string]any{"name": "id", "type": "`$STRING`", "active": true},
							map[string]any{"name": "title", "type": "`$STRING`", "active": true},
						},
						"op": map[string]any{
							"list": map[string]any{
								"name": "list",
								"targets": []any{
									map[string]any{
										"method": "GET",
										"parts":  []any{"foo"},
										"active": true,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	udk := um.Make("custom-api", map[string]any{"model": model})
	if udk == nil {
		t.Fatal("UniversalSDK should not be nil with model map")
	}

	// Verify the entity config was loaded from the inline model.
	config := udk.GetRootCtx().Config
	entityName := vs.GetPath([]any{"entity", "foo", "name"}, config)
	if entityName != "foo" {
		t.Fatalf("Expected entity name 'foo', got %v", entityName)
	}
}

func TestUniversalManagerModelMapOverridesRegistry(t *testing.T) {
	// Even with a registry path, an explicit model takes precedence.
	um := sdk.NewUniversalManager(map[string]any{
		"registry": "./registry",
	})

	customModel := map[string]any{
		"main": map[string]any{
			"kit": map[string]any{
				"entity": map[string]any{
					"bar": map[string]any{
						"name":   "bar",
						"active": true,
					},
				},
			},
		},
	}

	udk := um.Make("voxgig-solardemo", map[string]any{"model": customModel})
	if udk == nil {
		t.Fatal("UniversalSDK should not be nil")
	}

	// Should have "bar" from the inline model, NOT "moon"/"planet" from registry.
	config := udk.GetRootCtx().Config
	barName := vs.GetPath([]any{"entity", "bar", "name"}, config)
	if barName != "bar" {
		t.Fatalf("Expected entity 'bar' from inline model, got %v", barName)
	}
	moonName := vs.GetPath([]any{"entity", "moon", "name"}, config)
	if moonName != nil {
		t.Fatalf("Expected no 'moon' entity when inline model is used, got %v", moonName)
	}
}
