package mock

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Region

func BuildProviderSpec() []schema.Provider {
	return []schema.Provider{
		{
			Name:    authorizationProvider,
			Version: apiVersion1,
			Url:     generators.GenerateRegionProviderUrl(authorizationProvider),
		},
		{
			Name:    computeProvider,
			Version: apiVersion1,
			Url:     generators.GenerateRegionProviderUrl(computeProvider),
		},
		{
			Name:    networkProvider,
			Version: apiVersion1,
			Url:     generators.GenerateRegionProviderUrl(networkProvider),
		},
		{
			Name:    storageProvider,
			Version: apiVersion1,
			Url:     generators.GenerateRegionProviderUrl(storageProvider),
		},
		{
			Name:    workspaceProvider,
			Version: apiVersion1,
			Url:     generators.GenerateRegionProviderUrl(workspaceProvider),
		},
	}
}
