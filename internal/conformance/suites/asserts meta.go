package suites

import (
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (suite *TestSuite) VerifyGlobalTenantResourceMetadataStep(ctx provider.StepCtx, expected *schema.GlobalTenantResourceMetadata, actual *schema.GlobalTenantResourceMetadata) {
	ctx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
		stepCtx.Assert().Equal(expected.Name, actual.Name, "Name should match expected")
		stepCtx.Assert().Equal(expected.Provider, actual.Provider, "Provider should match expected")
		stepCtx.Assert().Equal(expected.Resource, actual.Resource, "Resource should match expected")
		stepCtx.Assert().Equal(expected.ApiVersion, actual.ApiVersion, "ApiVersion should match expected")
		stepCtx.Assert().Equal(expected.Verb, actual.Verb, "Verb should match expected")
		stepCtx.Assert().Equal(expected.Kind, actual.Kind, "Kind should match expected")
		stepCtx.Assert().Equal(expected.Ref, actual.Ref, "Metadata: Ref should match expected")
		stepCtx.Assert().Equal(expected.Tenant, actual.Tenant, "Tenant should match expected")

		stepCtx.Assert().LessOrEqual(len(actual.Tenant), 64, "Tenant identifier max length should be <= 64")

		suite.verifyAssertState(stepCtx)
	})
}

func (suite *TestSuite) VerifyGlobalResourceMetadataStep(ctx provider.StepCtx, expected *schema.GlobalResourceMetadata, actual *schema.GlobalResourceMetadata) {
	ctx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
		stepCtx.Assert().Equal(expected.Name, actual.Name, "Metadata: Name should match expected")
		stepCtx.Assert().Equal(expected.Provider, actual.Provider, "Metadata: Provider should match expected")
		stepCtx.Assert().Equal(expected.Resource, actual.Resource, "Metadata: Resource should match expected")
		stepCtx.Assert().Equal(expected.ApiVersion, actual.ApiVersion, "Metadata: ApiVersion should match expected")
		stepCtx.Assert().Equal(expected.Verb, actual.Verb, "Metadata: Verb should match expected")
		stepCtx.Assert().Equal(expected.Kind, actual.Kind, "Metadata: Kind should match expected")
		stepCtx.Assert().Equal(expected.Ref, actual.Ref, "Metadata: Ref should match expected")

		suite.verifyAssertState(stepCtx)
	})
}

func (suite *TestSuite) VerifyRegionalResourceMetadataStep(ctx provider.StepCtx, expected *schema.RegionalResourceMetadata, actual *schema.RegionalResourceMetadata) {
	ctx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
		stepCtx.Assert().Equal(expected.Name, actual.Name, "Metadata: Name should match expected")
		stepCtx.Assert().Equal(expected.Provider, actual.Provider, "Metadata: Provider should match expected")
		stepCtx.Assert().Equal(expected.Resource, actual.Resource, "Metadata: Resource should match expected")
		stepCtx.Assert().Equal(expected.ApiVersion, actual.ApiVersion, "Metadata: ApiVersion should match expected")
		stepCtx.Assert().Equal(expected.Verb, actual.Verb, "Metadata: Verb should match expected")
		stepCtx.Assert().Equal(expected.Kind, actual.Kind, "Metadata: Kind should match expected")
		stepCtx.Assert().Equal(expected.Ref, actual.Ref, "Metadata: Ref should match expected")
		stepCtx.Assert().Equal(expected.Tenant, actual.Tenant, "Metadata: Tenant should match expected")
		stepCtx.Assert().Equal(expected.Region, actual.Region, "Metadata: Region should match expected")

		stepCtx.Assert().LessOrEqual(len(actual.Tenant), 64, "Tenant identifier max length should be <= 64")
		stepCtx.Assert().LessOrEqual(len(actual.Region), 64, "Region identifier max length should be <= 64")

		suite.verifyAssertState(stepCtx)
	})
}

func (suite *TestSuite) VerifyRegionalWorkspaceResourceMetadataStep(ctx provider.StepCtx, expected *schema.RegionalWorkspaceResourceMetadata, actual *schema.RegionalWorkspaceResourceMetadata) {
	ctx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
		stepCtx.Assert().Equal(expected.Name, actual.Name, "Metadata: Name should match expected")
		stepCtx.Assert().Equal(expected.Provider, actual.Provider, "Metadata: Provider should match expected")
		stepCtx.Assert().Equal(expected.Resource, actual.Resource, "Metadata: Resource should match expected")
		stepCtx.Assert().Equal(expected.ApiVersion, actual.ApiVersion, "Metadata: ApiVersion should match expected")
		stepCtx.Assert().Equal(expected.Verb, actual.Verb, "Metadata: Verb should match expected")
		stepCtx.Assert().Equal(expected.Kind, actual.Kind, "Metadata: Kind should match expected")
		stepCtx.Assert().Equal(expected.Ref, actual.Ref, "Metadata: Ref should match expected")
		stepCtx.Assert().Equal(expected.Tenant, actual.Tenant, "Metadata: Tenant should match expected")
		stepCtx.Assert().Equal(expected.Workspace, actual.Workspace, "Metadata: Workspace should match expected")
		stepCtx.Assert().Equal(expected.Region, actual.Region, "Metadata: Region should match expected")

		stepCtx.Assert().LessOrEqual(len(actual.Tenant), 64, "Tenant identifier max length should be <= 64")
		stepCtx.Assert().LessOrEqual(len(actual.Workspace), 64, "Workspace identifier max length should be <= 64")
		stepCtx.Assert().LessOrEqual(len(actual.Region), 64, "Region identifier max length should be <= 64")

		suite.verifyAssertState(stepCtx)
	})
}

func (suite *TestSuite) VerifyRegionalNetworkResourceMetadataStep(ctx provider.StepCtx, expected *schema.RegionalNetworkResourceMetadata, actual *schema.RegionalNetworkResourceMetadata) {
	ctx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
		stepCtx.Assert().Equal(expected.Name, actual.Name, "Metadata: Name should match expected")
		stepCtx.Assert().Equal(expected.Provider, actual.Provider, "Metadata: Provider should match expected")
		stepCtx.Assert().Equal(expected.Resource, actual.Resource, "Metadata: Resource should match expected")
		stepCtx.Assert().Equal(expected.ApiVersion, actual.ApiVersion, "Metadata: ApiVersion should match expected")
		stepCtx.Assert().Equal(expected.Verb, actual.Verb, "Metadata: Verb should match expected")
		stepCtx.Assert().Equal(expected.Kind, actual.Kind, "Metadata: Kind should match expected")
		stepCtx.Assert().Equal(expected.Ref, actual.Ref, "Metadata: Ref should match expected")
		stepCtx.Assert().Equal(expected.Tenant, actual.Tenant, "Metadata: Tenant should match expected")
		stepCtx.Assert().Equal(expected.Workspace, actual.Workspace, "Metadata: Workspace should match expected")
		stepCtx.Assert().Equal(expected.Network, actual.Network, "Metadata: Network should match expected")
		stepCtx.Assert().Equal(expected.Region, actual.Region, "Metadata: Region should match expected")

		stepCtx.Assert().LessOrEqual(len(actual.Tenant), 64, "Tenant identifier max length should be <= 64")
		stepCtx.Assert().LessOrEqual(len(actual.Workspace), 64, "Workspace identifier max length should be <= 64")
		stepCtx.Assert().LessOrEqual(len(actual.Network), 64, "Network identifier max length should be <= 64")
		stepCtx.Assert().LessOrEqual(len(actual.Region), 64, "Region identifier max length should be <= 64")

		suite.verifyAssertState(stepCtx)
	})
}

func (suite *TestSuite) VerifyResponseMetadataStep(ctx provider.StepCtx, expected *schema.ResponseMetadata, actual *schema.ResponseMetadata) {
	ctx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
		stepCtx.Assert().Equal(expected.Provider, actual.Provider, "Metadata: Provider should match expected")
		stepCtx.Assert().Equal(expected.Resource, actual.Resource, "Metadata: Resource should match expected")
		stepCtx.Assert().Equal(expected.Verb, actual.Verb, "Metadata: Verb should match expected")

		suite.verifyAssertState(stepCtx)
	})
}
