package secatest

import (
	"context"
	"errors"
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
	RegionNameA := secalib.GenerateRegionName()
	RegionNameB := secalib.GenerateRegionName()

	if suite.mockEnabled {
		wm, err := mock.CreateRegionLifecycleScenarioV1(suite.scenarioName,
			mock.RegionParamsV1{
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
						Name: RegionNameA,
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
						Name: RegionNameB,
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
			})
		if err != nil {
			t.Fatalf("Failed to create region scenario: %v", err)
		}
		suite.mockClient = wm
	}

	ctx := context.Background()
	var regionResp *schema.Region

	t.WithNewStep("List Region", func(sCtx provider.StepCtx) {
		suite.setRegionV1StepParams(sCtx, "")

		iter, err := suite.client.RegionV1.ListRegions(ctx)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, iter)

		if iter != nil {
			firstRegion, err := iter.Next(ctx)
			if err != nil {
				sCtx.Error("Failed to get next region: %v", err)
			} else if firstRegion == nil {
				sCtx.Error("No regions found")
			} else {
				regionResp = firstRegion

				requireNotNilResponse(sCtx, regionResp)
				expected := &schema.GlobalResourceMetadata{
					Name:       regionResp.Metadata.Name,
					Verb:       regionResp.Metadata.Verb,
					Resource:   regionResp.Metadata.Resource,
					ApiVersion: regionResp.Metadata.ApiVersion,
					Kind:       regionResp.Metadata.Kind,
				}
				suite.verifyGlobalResourceMetadataStep(sCtx, expected, regionResp.Metadata)

				regions, err := iter.All(ctx)
				requireNoError(sCtx, err)

				verifyRegionExists(sCtx, suite.regionName, regions)

			}
		} else {
			sCtx.Error("No regions found")
		}
	})

	t.WithNewStep("List Region with params", func(sCtx provider.StepCtx) {
		suite.setRegionV1StepParams(sCtx, "")

		labelsParams := builders.NewLabelsBuilder().
			Equals(secalib.LabelEnvKey, secalib.LabelEnvValue)

		listOptions := builders.NewListOptions().WithLimit(10).WithLabels(labelsParams)

		iter, err := suite.client.RegionV1.ListRegionsWithFilters(ctx, listOptions)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, iter)

		if iter != nil {
			firstRegion, err := iter.Next(ctx)
			if err != nil {
				sCtx.Error("Failed to get next region: %v", err)
			} else if firstRegion == nil {
				sCtx.Error("No regions found")
			} else {
				regionResp = firstRegion

				requireNotNilResponse(sCtx, regionResp)
				expected := &schema.GlobalResourceMetadata{
					Name:       regionResp.Metadata.Name,
					Verb:       regionResp.Metadata.Verb,
					Resource:   regionResp.Metadata.Resource,
					ApiVersion: regionResp.Metadata.ApiVersion,
					Kind:       regionResp.Metadata.Kind,
				}
				suite.verifyGlobalResourceMetadataStep(sCtx, expected, regionResp.Metadata)

				regions, err := iter.All(ctx)
				requireNoError(sCtx, err)

				verifyRegionExists(sCtx, suite.regionName, regions)

			}
		} else {
			sCtx.Error("No regions found")
		}
	})

	t.WithNewStep("Get Region", func(sCtx provider.StepCtx) {
		suite.setRegionV1StepParams(sCtx, "")

		regionResp, err := suite.client.RegionV1.GetRegion(ctx, suite.regionName)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, regionResp)

		expected := &schema.GlobalResourceMetadata{
			Name:       regionResp.Metadata.Name,
			Verb:       regionResp.Metadata.Verb,
			Resource:   regionResp.Metadata.Resource,
			ApiVersion: regionResp.Metadata.ApiVersion,
			Kind:       regionResp.Metadata.Kind,
		}
		suite.verifyGlobalResourceMetadataStep(sCtx, expected, regionResp.Metadata)
	})

	t.WithNewStep("Not Exist Region", func(sCtx provider.StepCtx) {
		suite.setRegionV1StepParams(sCtx, "")

		_, err := suite.client.RegionV1.GetRegion(ctx, secalib.GenerateRegionName())
		expectedError := errors.New("resource not found")
		requireError(sCtx, err, expectedError)
	})
}

func verifyRegionExists(ctx provider.StepCtx, expectedRegion string, actualRegions []*schema.Region) {
	for _, region := range actualRegions {
		if region.Metadata.Name == config.clientRegion {
			ctx.WithNewStep("Verify status", func(stepCtx provider.StepCtx) {
				stepCtx.Require().Equal(expectedRegion, region.Metadata.Name, "State should match expected")
			})
		}
	}
}
