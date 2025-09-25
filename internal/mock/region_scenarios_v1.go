package mock

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/eu-sovereign-cloud/conformance/secalib"

	"github.com/wiremock/go-wiremock"
)

func CreateRegionLifecycleScenarioV1(scenario string, params RegionParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to Region Lifecycle Scenario")

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}
	var regionsResponse []*resourceResponse[secalib.RegionSpecV1]

	for _, region := range params.Regions {
		regionResource := secalib.GenerateRegionResource(region.Name)
		regionUrl := secalib.GenerateRegionURL(region.Name)

		regionsResponse = append(regionsResponse, &resourceResponse[secalib.RegionSpecV1]{
			Metadata: &secalib.Metadata{
				Name:           region.Name,
				Provider:       secalib.RegionProviderV1,
				Resource:       regionResource,
				ApiVersion:     secalib.ApiVersion1,
				Kind:           secalib.RegionKind,
				Tenant:         params.Tenant,
				Verb:           http.MethodGet,
				CreatedAt:      time.Now().Format(time.RFC3339),
				LastModifiedAt: time.Now().Format(time.RFC3339),
			},
			Spec: &secalib.RegionSpecV1{
				AvailableZones: region.InitialSpec.AvailableZones,
				Providers:      region.InitialSpec.Providers,
			},
		})

		// Region
		regionResponse := &resourceResponse[secalib.RegionSpecV1]{
			Metadata: &secalib.Metadata{
				Name:           region.Name,
				Provider:       secalib.RegionProviderV1,
				Resource:       regionResource,
				ApiVersion:     secalib.ApiVersion1,
				Kind:           secalib.RegionKind,
				Tenant:         params.Tenant,
				Verb:           http.MethodGet,
				CreatedAt:      time.Now().Format(time.RFC3339),
				LastModifiedAt: time.Now().Format(time.RFC3339),
			},
			Spec: &secalib.RegionSpecV1{
				AvailableZones: region.InitialSpec.AvailableZones,
				Providers:      region.InitialSpec.Providers,
			},
		}
		// Get Region
		if err := configureGetWithoutStateStub(wm, scenario, stubConfig{
			url:        regionUrl,
			params:     params,
			response:   regionResponse,
			template:   regionResponseTemplateV1,
			httpStatus: http.StatusOK,
		}); err != nil {
			return nil, err
		}

	}

	// List Region

	if err := configureGetWithoutStateStub(wm, scenario, stubConfig{
		url:        secalib.RegionsURLV1,
		params:     params,
		response:   regionsResponse,
		template:   regionListResponseTemplateV1,
		httpStatus: http.StatusOK,
	}); err != nil {
		return nil, err
	}

	return wm, nil
}
