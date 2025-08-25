package secatest

import "github.com/ozontech/allure-go/pkg/framework/provider"

type verifyRegionalMetadataStepParams struct {
	name       string
	provider   string
	resource   string
	verb       string
	apiVersion string
	kind       string
	tenant     string
	workspace  string
	region     string
}

type verifySkuMetadataStepParams struct {
	name string
}

type verifyStatusStepParams struct {
	expectedState string
	actualState   string
}

func verifyRegionalMetadataStep(ctx provider.StepCtx, expected verifyRegionalMetadataStepParams, actual verifyRegionalMetadataStepParams) {
	ctx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_name", expected.name,
			"actual_name", actual.name,

			"expected_provider", expected.provider,
			"actual_provider", actual.provider,

			"expected_resource", expected.resource,
			"actual_resource", actual.resource,

			"expected_verb", expected.verb,
			"actual_verb", actual.verb,

			"expected_apiVersion", expected.apiVersion,
			"actual_apiVersion", actual.apiVersion,

			"expected_kind", expected.kind,
			"actual_kind", actual.kind,

			"expected_tenant", expected.tenant,
			"actual_tenant", actual.tenant,

			"expected_workspace", expected.workspace,
			"actual_workspace", actual.workspace,

			"expected_region", expected.region,
			"actual_region", actual.region,
		)

		stepCtx.Require().Equal(expected.name, actual.name, "Name should match expected")
		stepCtx.Require().Equal(expected.provider, actual.provider, "Provider should match expected")
		stepCtx.Require().Equal(expected.resource, actual.resource, "Resource should match expected")
		stepCtx.Require().Equal(expected.verb, actual.verb, "Verb should match expected")
		stepCtx.Require().Equal(expected.apiVersion, actual.apiVersion, "ApiVersion should match expected")
		stepCtx.Require().Equal(expected.kind, actual.kind, "Kind should match expected")
		stepCtx.Require().Equal(expected.tenant, actual.tenant, "Tenant should match expected")
		stepCtx.Require().Equal(expected.region, actual.region, "Region should match expected")
	})
}

func verifySkuMetadataStep(ctx provider.StepCtx, expected verifySkuMetadataStepParams, actual verifySkuMetadataStepParams) {
	ctx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_name", expected.name,
			"actual_name", actual.name,
		)
		stepCtx.Require().Equal(expected.name, actual.name, "Name should match expected")
	})
}

func verifyStatusStep(ctx provider.StepCtx, params verifyStatusStepParams) {
	ctx.WithNewStep("Verify status", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_state", params.expectedState,
			"actual_state", params.actualState,
		)
		stepCtx.Require().Equal(params.expectedState, params.actualState, "State should match expected")
	})
}
