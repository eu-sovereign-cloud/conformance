package mock

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/builders"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigClientsInitScenario(params *ClientsInitParams) (*wiremock.Client, error) {
	slog.Info("Configuring mock to ClientsInit scenario")

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	resource := secalib.GenerateRegionResource(params.Region)
	url := secalib.GenerateRegionURL(params.Region)

	headers := map[string]string{
		authorizationHttpHeaderKey: authorizationHttpHeaderValuePrefix + params.AuthToken,
	}

	spec := &schema.RegionSpec{
		AvailableZones: []string{secalib.ZoneA, secalib.ZoneB},
		Providers: []schema.Provider{
			{
				Name:    secalib.AuthorizationProvider,
				Version: secalib.ApiVersion1,
				Url:     secalib.GenerateRegionProviderUrl(secalib.AuthorizationProvider),
			},
			{
				Name:    secalib.ComputeProvider,
				Version: secalib.ApiVersion1,
				Url:     secalib.GenerateRegionProviderUrl(secalib.ComputeProvider),
			},
			{
				Name:    secalib.NetworkProvider,
				Version: secalib.ApiVersion1,
				Url:     secalib.GenerateRegionProviderUrl(secalib.NetworkProvider),
			},
			{
				Name:    secalib.StorageProvider,
				Version: secalib.ApiVersion1,
				Url:     secalib.GenerateRegionProviderUrl(secalib.StorageProvider),
			},
			{
				Name:    secalib.WorkspaceProvider,
				Version: secalib.ApiVersion1,
				Url:     secalib.GenerateRegionProviderUrl(secalib.WorkspaceProvider),
			},
		},
	}

	response, err := builders.NewRegionBuilder().
		Name(params.Region).
		Provider(secalib.RegionProviderV1).
		Resource(resource).
		ApiVersion(secalib.ApiVersion1).
		Spec(spec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	if err := configureGetStub(wm, "ClientsInit",
		&stubConfig{url: url, params: params, pathParams: headers, responseBody: response}); err != nil {
		return nil, err
	}

	return wm, nil
}
