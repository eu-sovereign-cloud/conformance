package region

import (
	"context"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	mockRegion "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/region"

	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type ListV1TestSuite struct {
	suites.GlobalTestSuite

	RegionName string

	params *params.RegionListParamsV1
}

func (suite *ListV1TestSuite) BeforeAll(t provider.T) {
	var err error

	// Generate scenario Names
	regionName := generators.GenerateRegionName()
	regionName2 := generators.GenerateRegionName()
	regionName3 := generators.GenerateRegionName()

	region, err := builders.NewRegionBuilder().
		Name(regionName).
		Provider(constants.RegionProviderV1).ApiVersion(constants.ApiVersion1).
		Spec(&schema.RegionSpec{
			AvailableZones: []string{constants.ZoneA, constants.ZoneB},
			Providers:      mock.BuildProviderSpecV1(),
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Region: %v", err)
	}

	region2, err := builders.NewRegionBuilder().
		Name(regionName2).
		Provider(constants.RegionProviderV1).ApiVersion(constants.ApiVersion1).
		Spec(&schema.RegionSpec{
			AvailableZones: []string{constants.ZoneA, constants.ZoneB},
			Providers:      mock.BuildProviderSpecV1(),
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Region: %v", err)
	}

	region3, err := builders.NewRegionBuilder().
		Name(regionName3).
		Provider(constants.RegionProviderV1).ApiVersion(constants.ApiVersion1).
		Spec(&schema.RegionSpec{
			AvailableZones: []string{constants.ZoneA, constants.ZoneB},
			Providers:      mock.BuildProviderSpecV1(),
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Region: %v", err)
	}

	regions := []schema.Region{*region, *region2, *region3}

	params := &params.RegionListParamsV1{
		Regions: regions,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(&suite.TestSuite, mockRegion.ConfigureListScenarioV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ListV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.RegionProviderV1, string(schema.GlobalResourceMetadataKindResourceKindRegion))

	stepsBuilder := steps.NewStepsConfigurator(&suite.TestSuite, t)

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

func (suite *ListV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
