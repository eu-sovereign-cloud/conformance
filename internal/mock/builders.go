package mock

import (
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Region

// TODO Find a better package to it
func BuildProviderSpecV1() []schema.Provider {
	return []schema.Provider{
		{
			Name:    constants.AuthorizationProvider,
			Version: constants.ApiVersion1,
			Url:     generators.GenerateRegionProviderUrl(constants.AuthorizationProvider),
		},
		{
			Name:    constants.ComputeProvider,
			Version: constants.ApiVersion1,
			Url:     generators.GenerateRegionProviderUrl(constants.ComputeProvider),
		},
		{
			Name:    constants.NetworkProvider,
			Version: constants.ApiVersion1,
			Url:     generators.GenerateRegionProviderUrl(constants.NetworkProvider),
		},
		{
			Name:    constants.StorageProvider,
			Version: constants.ApiVersion1,
			Url:     generators.GenerateRegionProviderUrl(constants.StorageProvider),
		},
		{
			Name:    constants.WorkspaceProvider,
			Version: constants.ApiVersion1,
			Url:     generators.GenerateRegionProviderUrl(constants.WorkspaceProvider),
		},
	}
}
