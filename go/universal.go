package voxgiguniversalsdk

import (
	"voxgiguniversalsdk/core"
	"voxgiguniversalsdk/entity"
	"voxgiguniversalsdk/feature"
	_ "voxgiguniversalsdk/utility"
)

// Type aliases preserve external API.
type UniversalSDK = core.UniversalSDK
type UniversalManager = core.UniversalManager
type Context = core.Context
type Utility = core.Utility
type Feature = core.Feature
type Entity = core.Entity
type UniversalEntity = core.UniversalEntity
type FetcherFunc = core.FetcherFunc
type Spec = core.Spec
type Result = core.Result
type Response = core.Response
type Operation = core.Operation
type Control = core.Control
type UniversalError = core.UniversalError

// BaseFeature from feature package.
type BaseFeature = feature.BaseFeature

func init() {
	core.NewBaseFeatureFunc = func() core.Feature {
		return feature.NewBaseFeature()
	}
	core.NewTestFeatureFunc = func() core.Feature {
		return feature.NewTestFeature()
	}
	core.NewUniversalEntityFunc = func(client *core.UniversalSDK, name string, entopts map[string]any) core.UniversalEntity {
		return entity.NewUniversalEntity(client, name, entopts)
	}
}

// Constructor re-exports.
var NewUniversalManager = core.NewUniversalManager
var NewUniversalSDK = core.NewUniversalSDK
var NewContext = core.NewContext
var NewSpec = core.NewSpec
var NewResult = core.NewResult
var NewResponse = core.NewResponse
var NewOperation = core.NewOperation
var MakeConfig = core.MakeConfig
var NewBaseFeature = feature.NewBaseFeature
var NewTestFeature = feature.NewTestFeature
