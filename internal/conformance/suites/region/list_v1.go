package region

import (
	"context"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/region"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type RegionV1TestSuite struct {
	suites.GlobalTestSuite

	RegionName string
}

func (suite *RegionV1TestSuite) TestListScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.RegionProviderV1, string(schema.GlobalResourceMetadataKindResourceKindRegion))

	// Generate scenario Names
	regionNameA := generators.GenerateRegionName()
	regionNameB := generators.GenerateRegionName()
	regionNameC := generators.GenerateRegionName()
	// Configure Mock if enabled
	if suite.MockEnabled {

		mockParams := &mock.RegionListParamsV1{
			BaseParams: &mock.BaseParams{
				MockURL:   *suite.MockServerURL,
				AuthToken: suite.AuthToken,
				Tenant:    suite.Tenant,
			},
			Regions: []mock.ResourceParams[schema.RegionSpec]{
				{
					Name: suite.RegionName,
					InitialSpec: &schema.RegionSpec{
						AvailableZones: []string{constants.ZoneA, constants.ZoneA},
						Providers:      mock.BuildProviderSpecV1(),
					},
				},
				{
					Name: regionNameA,
					InitialSpec: &schema.RegionSpec{
						AvailableZones: []string{constants.ZoneA, constants.ZoneA},
						Providers:      mock.BuildProviderSpecV1(),
					},
				},
				{
					Name: regionNameB,
					InitialSpec: &schema.RegionSpec{
						AvailableZones: []string{constants.ZoneA, constants.ZoneA},
						Providers:      mock.BuildProviderSpecV1(),
					},
				},
				{
					Name: regionNameC,
					InitialSpec: &schema.RegionSpec{
						AvailableZones: []string{constants.ZoneA, constants.ZoneA},
						Providers:      mock.BuildProviderSpecV1(),
					},
				},
			},
		}

		wm, err := region.ConfigureListScenarioV1(suite.ScenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.MockClient = wm
	}

	stepsBuilder := steps.NewBuilder(&suite.TestSuite, t)

	ctx := context.Background()

	// Test List iterator's (Next and All) for Regions and verify both responses have the same length
	regions := stepsBuilder.ListRegionsV1Step("List all regions", ctx, suite.Client.RegionV1)

	// Call Get Region and verify response
	expectedRegionMeta, err := builders.NewRegionMetadataBuilder().
		Name(regions[0].Metadata.Name).
		Provider(constants.RegionProviderV1).ApiVersion(constants.ApiVersion1).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}

	stepsBuilder.GetRegionV1Step("Get region "+regions[0].Metadata.Name, ctx, suite.Client.RegionV1, expectedRegionMeta)

	suite.FinishScenario()
}

func (suite *RegionV1TestSuite) AfterEach(t provider.T) {
	suite.ResetAllScenarios()
}
