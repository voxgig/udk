package test

import (
	"testing"

	sdk "voxgiguniversalsdk"
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
