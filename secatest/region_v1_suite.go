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
	"k8s.io/utils/ptr"
)

type RegionV1TestSuite struct {
	globalTestSuite
	regionName string
}

func (suite *RegionV1TestSuite) TestRegionV1(t provider.T) {
	slog.Info("Starting Region Lifecycle Test")

	t.Title("Region Lifecycle Test")
	configureTags(t, secalib.RegionKind)
	RegionNameA := secalib.GenerateRegionName()
	RegionNameB := secalib.GenerateRegionName()
	if suite.isMockEnabled() {
		wm, err := mock.CreateRegionLifecycleScenarioV1("Region Lifecycle",
			mock.RegionParamsV1{
				Params: &mock.Params{
					MockURL:   *suite.mockServerURL,
					AuthToken: suite.authToken,
					Tenant:    suite.tenant,
				},
				Regions: []mock.ResourceParams[secalib.RegionSpecV1]{
					{
						Name: suite.regionName,
						InitialSpec: &secalib.RegionSpecV1{
							AvailableZones: []string{"zone-a", "zone-b"},
							Providers: []secalib.Providers{
								{
									Name:    secalib.AuthorizationProvider,
									Version: secalib.ApiVersion1,
									URL:     secalib.GenerateRegionProviderUrl(secalib.AuthorizationProvider),
								},
								{
									Name:    secalib.ComputeProvider,
									Version: secalib.ApiVersion1,
									URL:     secalib.GenerateRegionProviderUrl(secalib.ComputeProvider),
								},
								{
									Name:    secalib.NetworkProvider,
									Version: secalib.ApiVersion1,
									URL:     secalib.GenerateRegionProviderUrl(secalib.NetworkProvider),
								},
								{
									Name:    secalib.StorageProvider,
									Version: secalib.ApiVersion1,
									URL:     secalib.GenerateRegionProviderUrl(secalib.StorageProvider),
								},
								{
									Name:    secalib.WorkspaceProvider,
									Version: secalib.ApiVersion1,
									URL:     secalib.GenerateRegionProviderUrl(secalib.WorkspaceProvider),
								},
							},
						},
					},
					{
						Name: RegionNameA,
						InitialSpec: &secalib.RegionSpecV1{
							AvailableZones: []string{"zone-a", "zone-b"},
							Providers: []secalib.Providers{
								{
									Name:    secalib.AuthorizationProvider,
									Version: secalib.ApiVersion1,
									URL:     secalib.GenerateRegionProviderUrl(secalib.AuthorizationProviderV1),
								},
								{
									Name:    secalib.ComputeProvider,
									Version: secalib.ApiVersion1,
									URL:     secalib.GenerateRegionProviderUrl(secalib.ComputeProviderV1),
								},
								{
									Name:    secalib.NetworkProvider,
									Version: secalib.ApiVersion1,
									URL:     secalib.GenerateRegionProviderUrl(secalib.NetworkProviderV1),
								},
								{
									Name:    secalib.StorageProvider,
									Version: secalib.ApiVersion1,
									URL:     secalib.GenerateRegionProviderUrl(secalib.StorageProviderV1),
								},
								{
									Name:    secalib.WorkspaceProvider,
									Version: secalib.ApiVersion1,
									URL:     secalib.GenerateRegionProviderUrl(secalib.WorkspaceProviderV1),
								},
							},
						},
					},
					{
						Name: RegionNameB,
						InitialSpec: &secalib.RegionSpecV1{
							AvailableZones: []string{"zone-a", "zone-b"},
							Providers: []secalib.Providers{
								{
									Name:    secalib.AuthorizationProvider,
									Version: secalib.ApiVersion1,
									URL:     secalib.GenerateRegionProviderUrl(secalib.AuthorizationProviderV1),
								},
								{
									Name:    secalib.ComputeProvider,
									Version: secalib.ApiVersion1,
									URL:     secalib.GenerateRegionProviderUrl(secalib.ComputeProviderV1),
								},
								{
									Name:    secalib.NetworkProvider,
									Version: secalib.ApiVersion1,
									URL:     secalib.GenerateRegionProviderUrl(secalib.NetworkProviderV1),
								},
								{
									Name:    secalib.StorageProvider,
									Version: secalib.ApiVersion1,
									URL:     secalib.GenerateRegionProviderUrl(secalib.StorageProviderV1),
								},
								{
									Name:    secalib.WorkspaceProvider,
									Version: secalib.ApiVersion1,
									URL:     secalib.GenerateRegionProviderUrl(secalib.WorkspaceProviderV1),
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
				expectedMetadata := &secalib.Metadata{
					Name:       regionResp.Metadata.Name,
					Verb:       regionResp.Metadata.Verb,
					Resource:   regionResp.Metadata.Resource,
					ApiVersion: regionResp.Metadata.ApiVersion,
					Kind:       string(regionResp.Metadata.Kind),
					Tenant:     regionResp.Metadata.Tenant,
				}
				verifyRegionMetadataStep(sCtx, expectedMetadata, regionResp.Metadata)

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
		limit := 1
		labels := builders.NewLabelsBuilder().
			Equals("env", "test").
			Build()

		iter, err := suite.client.RegionV1.ListRegionsWithFilters(ctx, ptr.To(limit), ptr.To(labels))
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
				expectedMetadata := &secalib.Metadata{
					Name:       regionResp.Metadata.Name,
					Verb:       regionResp.Metadata.Verb,
					Resource:   regionResp.Metadata.Resource,
					ApiVersion: regionResp.Metadata.ApiVersion,
					Kind:       string(regionResp.Metadata.Kind),
					Tenant:     regionResp.Metadata.Tenant,
				}
				verifyRegionMetadataStep(sCtx, expectedMetadata, regionResp.Metadata)

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

		expectedMetadata := &secalib.Metadata{
			Name:       regionResp.Metadata.Name,
			Verb:       regionResp.Metadata.Verb,
			Resource:   regionResp.Metadata.Resource,
			ApiVersion: regionResp.Metadata.ApiVersion,
			Kind:       string(regionResp.Metadata.Kind),
			Tenant:     regionResp.Metadata.Tenant,
		}
		verifyRegionMetadataStep(sCtx, expectedMetadata, regionResp.Metadata)
	})

	t.WithNewStep("Not Exist Region", func(sCtx provider.StepCtx) {
		suite.setRegionV1StepParams(sCtx, "")

		_, err := suite.client.RegionV1.GetRegion(ctx, secalib.GenerateRegionName())
		expectedError := errors.New("resource not found")
		requireError(sCtx, err, expectedError)
	})
}

func verifyRegionMetadataStep(ctx provider.StepCtx, expected *secalib.Metadata, actual *schema.GlobalResourceMetadata) {
	actualMetadata := &secalib.Metadata{
		Name:       actual.Name,
		Verb:       actual.Verb,
		Resource:   actual.Resource,
		ApiVersion: actual.ApiVersion,
		Kind:       string(actual.Kind),
		Tenant:     actual.Tenant,
	}
	verifyGlobalMetadataStep(ctx, expected, actualMetadata)
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
