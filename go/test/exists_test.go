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

func TestUniversalManagerRegisterAndDescribe(t *testing.T) {
	um := sdk.NewUniversalManager(map[string]any{})

	modelA, _ := jsonic.Parse(`{
		name: alpha
		main: kit: {
			info: {title: 'Alpha API', version: '1.0'}
			entity: { a1: {name: a1, active: true} }
		}
	}`)
	modelB, _ := jsonic.Parse(`{
		name: beta
		main: kit: {
			info: {title: 'Beta API', version: '2.0'}
			entity: { b1: {name: b1, active: true} }
		}
	}`)

	um.Register("beta-api", modelB.(map[string]any))
	um.Register("alpha-api", modelA.(map[string]any))

	desc := um.Describe()
	models, _ := desc["models"].([]string)
	if len(models) != 2 {
		t.Fatalf("Expected 2 models, got %d", len(models))
	}
	if models[0] != "alpha-api" || models[1] != "beta-api" {
		t.Fatalf("Expected sorted [alpha-api beta-api], got %v", models)
	}
}

func TestUniversalManagerMakeRegistersModel(t *testing.T) {
	um := sdk.NewUniversalManager(map[string]any{})

	model, _ := jsonic.Parse(`{
		name: gamma
		main: kit: entity: { g1: {name: g1, active: true} }
	}`)

	um.Make("gamma-api", map[string]any{"model": model})

	desc := um.Describe()
	models, _ := desc["models"].([]string)
	if len(models) != 1 || models[0] != "gamma-api" {
		t.Fatalf("Expected [gamma-api], got %v", models)
	}

	// Making again should not duplicate.
	um.Make("gamma-api", map[string]any{"model": model})
	desc = um.Describe()
	models, _ = desc["models"].([]string)
	if len(models) != 1 {
		t.Fatalf("Expected 1 model after second Make, got %d", len(models))
	}
}

func TestUniversalSDKModel(t *testing.T) {
	um := sdk.NewUniversalManager(map[string]any{})

	model, _ := jsonic.Parse(`{
		name: solardemo
		main: kit: {
			info: {title: 'Solar System API', version: '1.0'}
			entity: {
				planet: {
					name: planet
					active: true
					fields: [
						{name: id, type: '$STRING', active: true}
						{name: name, type: '$STRING', active: true}
					]
				}
				moon: {
					name: moon
					active: true
					fields: [
						{name: id, type: '$STRING', active: true}
						{name: name, type: '$STRING', active: true}
						{name: planet_id, type: '$STRING', active: true}
					]
				}
			}
		}
	}`)

	udk := um.Make("demo", map[string]any{"model": model})

	t.Run("model-returns-full-model", func(t *testing.T) {
		m := udk.Model()
		if m["name"] != "solardemo" {
			t.Fatalf("Expected model name 'solardemo', got %v", m["name"])
		}
	})

	t.Run("entity-model-returns-entity", func(t *testing.T) {
		em := udk.EntityModel("planet")
		if em["name"] != "planet" {
			t.Fatalf("Expected entity name 'planet', got %v", em["name"])
		}
		fields, _ := em["fields"].([]any)
		if len(fields) != 2 {
			t.Fatalf("Expected 2 fields for planet, got %d", len(fields))
		}
	})

	t.Run("entity-model-missing-returns-empty", func(t *testing.T) {
		em := udk.EntityModel("nonexistent")
		if len(em) != 0 {
			t.Fatalf("Expected empty map for missing entity, got %v", em)
		}
	})

	t.Run("describe-api", func(t *testing.T) {
		desc := udk.Describe(map[string]any{"what": "api"})
		if desc["name"] != "solardemo" {
			t.Fatalf("Expected name 'solardemo', got %v", desc["name"])
		}
		entities, _ := desc["entities"].([]string)
		if len(entities) != 2 {
			t.Fatalf("Expected 2 entities, got %d", len(entities))
		}
		if entities[0] != "moon" || entities[1] != "planet" {
			t.Fatalf("Expected sorted [moon planet], got %v", entities)
		}
		info, _ := desc["info"].(map[string]any)
		if info["title"] != "Solar System API" {
			t.Fatalf("Expected title 'Solar System API', got %v", info["title"])
		}
	})

	t.Run("describe-entity", func(t *testing.T) {
		desc := udk.Describe(map[string]any{"what": "entity", "entity": "moon"})
		if desc["name"] != "moon" {
			t.Fatalf("Expected name 'moon', got %v", desc["name"])
		}
		fields, _ := desc["fields"].([]any)
		if len(fields) != 3 {
			t.Fatalf("Expected 3 fields for moon, got %d", len(fields))
		}
	})

	t.Run("describe-entity-missing", func(t *testing.T) {
		desc := udk.Describe(map[string]any{"what": "entity", "entity": "nope"})
		if desc["name"] != "" {
			t.Fatalf("Expected empty name for missing entity, got %v", desc["name"])
		}
	})
}
