package secatest

import (
	"context"
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi/builders"
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
						Providers:      secalib.GenerateProviderSpec(),
					},
				},
				{
					Name: regionNameA,
					InitialSpec: &schema.RegionSpec{
						AvailableZones: []string{secalib.ZoneA, secalib.ZoneB},
						Providers:      secalib.GenerateProviderSpec(),
					},
				},
				{
					Name: regionNameB,
					InitialSpec: &schema.RegionSpec{
						AvailableZones: []string{secalib.ZoneA, secalib.ZoneB},
						Providers:      secalib.GenerateProviderSpec(),
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
		},
	}
	for _, region := range regions {
		regionResource := secalib.GenerateRegionResource(region.Metadata.Name)
		expectedRegionMeta := secalib.NewGlobalResourceMetadata(region.Metadata.Name, secalib.RegionProviderV1, regionResource, secalib.ApiVersion1, secalib.RegionKind)
		suite.getRegion("Get region", t, ctx, suite.client.RegionV1, expectedRegionMeta)
	}
	// List all Regions
	suite.getListRegion("Get list region", t, ctx, suite.client.RegionV1)

	// List number of Regions defined in Limit key & labels
	labelsParams := builders.NewLabelsBuilder().
		Equals("env", "Development")
	listOptions := builders.NewListOptions().WithLimit(1).WithLabels(labelsParams)
	suite.getListRegionWithParameters("Get list region with Limit", t, ctx, suite.client.RegionV1, listOptions)

}
