package secatest

import (
	"context"
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type RegionV1TestSuite struct {
	globalTestSuite
	regionName string
}

func (suite *RegionV1TestSuite) TestSuite(t provider.T) {
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t, secalib.RegionKind)

	// Generate scenario Names
	regionNameA := secalib.GenerateRegionName()
	regionNameB := secalib.GenerateRegionName()

	if suite.mockEnabled {

		mockParams := &mock.RegionParamsV1{
			Params: &mock.Params{
				MockURL:   *suite.mockServerURL,
				AuthToken: suite.authToken,
				Tenant:    suite.tenant,
			},
			Regions: []mock.ResourceParams[schema.RegionSpec]{
				{
					Name: suite.regionName,
					InitialSpec: &schema.RegionSpec{
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
					},
				},
				{
					Name: regionNameA,
					InitialSpec: &schema.RegionSpec{
						AvailableZones: []string{secalib.ZoneA, secalib.ZoneB},
						Providers: []schema.Provider{
							{
								Name:    secalib.AuthorizationProvider,
								Version: secalib.ApiVersion1,
								Url:     secalib.GenerateRegionProviderUrl(secalib.AuthorizationProviderV1),
							},
							{
								Name:    secalib.ComputeProvider,
								Version: secalib.ApiVersion1,
								Url:     secalib.GenerateRegionProviderUrl(secalib.ComputeProviderV1),
							},
							{
								Name:    secalib.NetworkProvider,
								Version: secalib.ApiVersion1,
								Url:     secalib.GenerateRegionProviderUrl(secalib.NetworkProviderV1),
							},
							{
								Name:    secalib.StorageProvider,
								Version: secalib.ApiVersion1,
								Url:     secalib.GenerateRegionProviderUrl(secalib.StorageProviderV1),
							},
							{
								Name:    secalib.WorkspaceProvider,
								Version: secalib.ApiVersion1,
								Url:     secalib.GenerateRegionProviderUrl(secalib.WorkspaceProviderV1),
							},
						},
					},
				},
				{
					Name: regionNameB,
					InitialSpec: &schema.RegionSpec{
						AvailableZones: []string{secalib.ZoneA, secalib.ZoneB},
						Providers: []schema.Provider{
							{
								Name:    secalib.AuthorizationProvider,
								Version: secalib.ApiVersion1,
								Url:     secalib.GenerateRegionProviderUrl(secalib.AuthorizationProviderV1),
							},
							{
								Name:    secalib.ComputeProvider,
								Version: secalib.ApiVersion1,
								Url:     secalib.GenerateRegionProviderUrl(secalib.ComputeProviderV1),
							},
							{
								Name:    secalib.NetworkProvider,
								Version: secalib.ApiVersion1,
								Url:     secalib.GenerateRegionProviderUrl(secalib.NetworkProviderV1),
							},
							{
								Name:    secalib.StorageProvider,
								Version: secalib.ApiVersion1,
								Url:     secalib.GenerateRegionProviderUrl(secalib.StorageProviderV1),
							},
							{
								Name:    secalib.WorkspaceProvider,
								Version: secalib.ApiVersion1,
								Url:     secalib.GenerateRegionProviderUrl(secalib.WorkspaceProviderV1),
							},
						},
					},
				},
			},
		}

		wm, err := mock.CreateRegionLifecycleScenarioV1(suite.scenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to create region scenario: %v", err)
		}
		suite.mockClient = wm
	}

	ctx := context.Background()

	regions := []schema.Region{
		{
			Metadata: &schema.GlobalResourceMetadata{
				Name: suite.regionName,
			},
			Spec: schema.RegionSpec{
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
			},
		},
		{
			Metadata: &schema.GlobalResourceMetadata{
				Name: regionNameA,
			},
			Spec: schema.RegionSpec{
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
			},
		},
		{
			Metadata: &schema.GlobalResourceMetadata{
				Name: regionNameB,
			},
			Spec: schema.RegionSpec{
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
			},
		},
	}
	for _, region := range regions {
		regionResource := secalib.GenerateRegionResource(region.Metadata.Name)
		expectedRegionMeta := secalib.NewGlobalResourceMetadata(region.Metadata.Name, secalib.RegionProviderV1, regionResource, secalib.ApiVersion1, secalib.RegionKind)
		expectedRegionSpec := region.Spec
		suite.getRegion("Get region", t, ctx, suite.client.RegionV1, expectedRegionMeta, &expectedRegionSpec)
	}

	suite.getListRegion("Get list region", t, ctx, suite.client.RegionV1, regions)
}
