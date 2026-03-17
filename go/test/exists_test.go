package test

import (
	"testing"

	sdk "voxgiguniversalsdk"

	jsonic "github.com/jsonicjs/jsonic/go"
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

	// Provide the API model inline using jsonic syntax.
	model, merr := jsonic.Parse(`
		main: kit: entity: foo: {
			name: foo
			active: true
			fields: [
				{name: id, type: '$STRING', active: true}
				{name: title, type: '$STRING', active: true}
			]
			op: list: {
				name: list
				targets: [{method: GET, parts: [foo], active: true}]
			}
		}
	`)
	if merr != nil {
		t.Fatalf("jsonic.Parse failed: %v", merr)
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

	// Use jsonic to define the custom model inline.
	customModel, merr := jsonic.Parse(`
		main: kit: entity: bar: {name: bar, active: true}
	`)
	if merr != nil {
		t.Fatalf("jsonic.Parse failed: %v", merr)
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
