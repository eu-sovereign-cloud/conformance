package mock

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Region

func BuildProviderSpecV1() []schema.Provider {
	return []schema.Provider{
		{
			Name:    authorizationProvider,
			Version: ApiVersion1,
			Url:     generators.GenerateRegionProviderUrl(authorizationProvider),
		},
		{
			Name:    computeProvider,
			Version: ApiVersion1,
			Url:     generators.GenerateRegionProviderUrl(computeProvider),
		},
		{
			Name:    networkProvider,
			Version: ApiVersion1,
			Url:     generators.GenerateRegionProviderUrl(networkProvider),
		},
		{
			Name:    storageProvider,
			Version: ApiVersion1,
			Url:     generators.GenerateRegionProviderUrl(storageProvider),
		},
		{
			Name:    workspaceProvider,
			Version: ApiVersion1,
			Url:     generators.GenerateRegionProviderUrl(workspaceProvider),
		},
	}
}
