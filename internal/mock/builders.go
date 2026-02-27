package mock

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Region

// TODO Find a better package to it
func BuildProviderSpecV1() []schema.Provider {
	return []schema.Provider{
		{
			Name:    sdkconsts.AuthorizationProviderName,
			Version: sdkconsts.ApiVersion1,
			Url:     generators.GenerateRegionProviderUrl(sdkconsts.AuthorizationProviderName),
		},
		{
			Name:    sdkconsts.ComputeProviderName,
			Version: sdkconsts.ApiVersion1,
			Url:     generators.GenerateRegionProviderUrl(sdkconsts.ComputeProviderName),
		},
		{
			Name:    sdkconsts.NetworkProviderName,
			Version: sdkconsts.ApiVersion1,
			Url:     generators.GenerateRegionProviderUrl(sdkconsts.NetworkProviderName),
		},
		{
			Name:    sdkconsts.StorageProviderName,
			Version: sdkconsts.ApiVersion1,
			Url:     generators.GenerateRegionProviderUrl(sdkconsts.StorageProviderName),
		},
		{
			Name:    sdkconsts.WorkspaceProviderName,
			Version: sdkconsts.ApiVersion1,
			Url:     generators.GenerateRegionProviderUrl(sdkconsts.WorkspaceProviderName),
		},
	}
}
