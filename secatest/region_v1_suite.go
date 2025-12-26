package secatest

import (
	"context"
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
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
	configureTags(t, regionProviderV1, string(schema.GlobalResourceMetadataKindResourceKindRegion))

	// Generate scenario Names
	regionNameA := generators.GenerateRegionName()
	regionNameB := generators.GenerateRegionName()
	regionNameC := generators.GenerateRegionName()
	// Configure Mock if enabled
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
						AvailableZones: []string{zoneA, zoneB},
						Providers:      mock.BuildProviderSpec(),
					},
				},
				{
					Name: regionNameA,
					InitialSpec: &schema.RegionSpec{
						AvailableZones: []string{zoneA, zoneB},
						Providers:      mock.BuildProviderSpec(),
					},
				},
				{
					Name: regionNameB,
					InitialSpec: &schema.RegionSpec{
						AvailableZones: []string{zoneA, zoneB},
						Providers:      mock.BuildProviderSpec(),
					},
				},
				{
					Name: regionNameC,
					InitialSpec: &schema.RegionSpec{
						AvailableZones: []string{zoneA, zoneB},
						Providers:      mock.BuildProviderSpec(),
					},
				},
			},
		}

		wm, err := mock.ConfigRegionLifecycleScenarioV1(suite.scenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	ctx := context.Background()

	// Test List iterator's (Next and All) for Regions and verify both responses have the same length
	regions := suite.listRegionsV1Step("List all regions", t, ctx, suite.client.RegionV1)

	// Call Get Region and verify response
	expectedRegionMeta, err := builders.NewRegionMetadataBuilder().
		Name(regions[0].Metadata.Name).
		Provider(regionProviderV1).ApiVersion(apiVersion1).
		Build()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}

	suite.getRegionV1Step("Get region "+regions[0].Metadata.Name, t, ctx, suite.client.RegionV1, expectedRegionMeta)
}

func (suite *RegionV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
