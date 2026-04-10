package mock

import (
	"slices"

	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Region

// TODO Find a better package to it
func BuildProviderSpecV1(mockProviders []string) []schema.Provider {
	var providers []schema.Provider

	if len(mockProviders) > 0 {
		if slices.Contains(mockProviders, sdkconsts.AuthorizationProviderName) {
			providers = append(providers, schema.Provider{
				Name:    sdkconsts.AuthorizationProviderName,
				Version: sdkconsts.ApiVersion1,
				Url:     generators.GenerateRegionProviderUrl(sdkconsts.AuthorizationProviderName),
			})
		}

		if slices.Contains(mockProviders, sdkconsts.ComputeProviderName) {
			providers = append(providers, schema.Provider{
				Name:    sdkconsts.ComputeProviderName,
				Version: sdkconsts.ApiVersion1,
				Url:     generators.GenerateRegionProviderUrl(sdkconsts.ComputeProviderName),
			})
		}

		if slices.Contains(mockProviders, sdkconsts.NetworkProviderName) {
			providers = append(providers, schema.Provider{
				Name:    sdkconsts.NetworkProviderName,
				Version: sdkconsts.ApiVersion1,
				Url:     generators.GenerateRegionProviderUrl(sdkconsts.NetworkProviderName),
			})
		}

		if slices.Contains(mockProviders, sdkconsts.StorageProviderName) {
			providers = append(providers, schema.Provider{
				Name:    sdkconsts.StorageProviderName,
				Version: sdkconsts.ApiVersion1,
				Url:     generators.GenerateRegionProviderUrl(sdkconsts.StorageProviderName),
			})
		}

		if slices.Contains(mockProviders, sdkconsts.WorkspaceProviderName) {
			providers = append(providers, schema.Provider{
				Name:    sdkconsts.WorkspaceProviderName,
				Version: sdkconsts.ApiVersion1,
				Url:     generators.GenerateRegionProviderUrl(sdkconsts.WorkspaceProviderName),
			})
		}
	} else {
		providers = append(providers,
			schema.Provider{
				Name:    sdkconsts.AuthorizationProviderName,
				Version: sdkconsts.ApiVersion1,
				Url:     generators.GenerateRegionProviderUrl(sdkconsts.AuthorizationProviderName),
			},
			schema.Provider{
				Name:    sdkconsts.ComputeProviderName,
				Version: sdkconsts.ApiVersion1,
				Url:     generators.GenerateRegionProviderUrl(sdkconsts.ComputeProviderName),
			},
			schema.Provider{
				Name:    sdkconsts.NetworkProviderName,
				Version: sdkconsts.ApiVersion1,
				Url:     generators.GenerateRegionProviderUrl(sdkconsts.NetworkProviderName),
			},
			schema.Provider{
				Name:    sdkconsts.StorageProviderName,
				Version: sdkconsts.ApiVersion1,
				Url:     generators.GenerateRegionProviderUrl(sdkconsts.StorageProviderName),
			},
			schema.Provider{
				Name:    sdkconsts.WorkspaceProviderName,
				Version: sdkconsts.ApiVersion1,
				Url:     generators.GenerateRegionProviderUrl(sdkconsts.WorkspaceProviderName),
			})
	}

	return providers
}
