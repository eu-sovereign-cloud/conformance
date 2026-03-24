package region

import (
	"context"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	mockRegion "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/region"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type ProviderQueriesV1TestSuite struct {
	suites.GlobalTestSuite

	Regions []string

	params *params.RegionProviderQueriesV1Params
}

func CreateProviderQueriesV1TestSuite(globalTestSuite suites.GlobalTestSuite, clientRegion string, additionalRegions []string) *ProviderQueriesV1TestSuite {
	suite := &ProviderQueriesV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		Regions:         append([]string{clientRegion}, additionalRegions...),
	}
	suite.ScenarioName = constants.RegionProviderQueriesV1SuiteName.String()
	return suite
}

func (suite *ProviderQueriesV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Region")

	// Generate scenario Names
	var regions []schema.Region
	for _, region := range suite.Regions {

		resource, err := builders.NewRegionBuilder().
			Name(region).
			Provider(sdkconsts.RegionProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Spec(&schema.RegionSpec{
				AvailableZones: []string{constants.ZoneA, constants.ZoneB},
				Providers:      mock.BuildProviderSpecV1(),
			}).
			Build()
		if err != nil {
			t.Fatalf("Failed to build Region: %v", err)
		}

		regions = append(regions, *resource)
	}

	params := &params.RegionProviderQueriesV1Params{
		Regions: regions,
	}
	suite.params = params
	err := suites.SetupMockIfEnabled(suite.TestSuite, mockRegion.ConfigureProviderQueriesV1, *params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ProviderQueriesV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.RegionProviderV1Name, string(schema.GlobalResourceMetadataKindResourceKindRegion))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	ctx := context.Background()

	// Test List iterator's (Next and All) for Regions and verify both responses have the same length
	regions := stepsBuilder.ListRegionsV1Step("List all regions", ctx, suite.Client.RegionV1)

	// Call Get Region and verify response
	expectedRegionMeta := suite.params.Regions[0].Metadata
	stepsBuilder.GetRegionV1Step("Get region "+regions[0].Metadata.Name, ctx, suite.Client.RegionV1, expectedRegionMeta)

	suite.FinishScenario()
}

func (suite *ProviderQueriesV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
