package secatest

import (
	"fmt"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func verifyGlobalMetadataStep(ctx provider.StepCtx, expected *secalib.Metadata, actual *secalib.Metadata) {
	ctx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_name", expected.Name,
			"actual_name", actual.Name,

			"expected_provider", expected.Provider,
			"actual_provider", actual.Provider,

			"expected_resource", expected.Resource,
			"actual_resource", actual.Resource,

			"expected_verb", expected.Verb,
			"actual_verb", actual.Verb,

			"expected_api_version", expected.ApiVersion,
			"actual_api_version", actual.ApiVersion,

			"expected_kind", expected.Kind,
			"actual_kind", actual.Kind,

			"expected_tenant", expected.Tenant,
			"actual_tenant", actual.Tenant,
		)

		stepCtx.Require().Equal(expected.Name, actual.Name, "Name should match expected")
		stepCtx.Require().Equal(expected.Provider, actual.Provider, "Provider should match expected")
		stepCtx.Require().Equal(expected.Resource, actual.Resource, "Resource should match expected")
		stepCtx.Require().Equal(expected.Verb, actual.Verb, "Verb should match expected")
		stepCtx.Require().Equal(expected.ApiVersion, actual.ApiVersion, "ApiVersion should match expected")
		stepCtx.Require().Equal(expected.Kind, actual.Kind, "Kind should match expected")
	})
}

func verifyRegionalMetadataStep(ctx provider.StepCtx, expected *secalib.Metadata, actual *secalib.Metadata) {
	ctx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_name", expected.Name,
			"actual_name", actual.Name,

			"expected_provider", expected.Provider,
			"actual_provider", actual.Provider,

			"expected_resource", expected.Resource,
			"actual_resource", actual.Resource,

			"expected_verb", expected.Verb,
			"actual_verb", actual.Verb,

			"expected_api_version", expected.ApiVersion,
			"actual_api_version", actual.ApiVersion,

			"expected_kind", expected.Kind,
			"actual_kind", actual.Kind,

			"expected_tenant", expected.Tenant,
			"actual_tenant", actual.Tenant,

			"expected_workspace", expected.Workspace,
			"actual_workspace", actual.Workspace,

			"expected_region", expected.Region,
			"actual_region", actual.Region,
		)

		stepCtx.Require().Equal(expected.Name, actual.Name, "Name should match expected")
		stepCtx.Require().Equal(expected.Provider, actual.Provider, "Provider should match expected")
		stepCtx.Require().Equal(expected.Resource, actual.Resource, "Resource should match expected")
		stepCtx.Require().Equal(expected.Verb, actual.Verb, "Verb should match expected")
		stepCtx.Require().Equal(expected.ApiVersion, actual.ApiVersion, "ApiVersion should match expected")
		stepCtx.Require().Equal(expected.Kind, actual.Kind, "Kind should match expected")
		stepCtx.Require().Equal(expected.Tenant, actual.Tenant, "Tenant should match expected")
		stepCtx.Require().Equal(expected.Region, actual.Region, "Region should match expected")
	})
}

func verifyStatusStep(ctx provider.StepCtx, expected *secalib.Status, actual *secalib.Status) {
	ctx.WithNewStep("Verify status", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_state", expected.State,
			"actual_state", actual.State,
		)
		stepCtx.Require().Equal(expected.State, actual.State, "State should match expected")
	})
}

func verifyLabelStep(ctx provider.StepCtx, expected *[]secalib.Label, actual *[]secalib.Label) {
	ctx.WithNewStep("Verify label", func(stepCtx provider.StepCtx) {
		for i := 0; i < len(*expected); i++ {
			expectedValue := (*expected)[i].Value
			found := false
			for _, label := range *actual {
				if label.Value == expectedValue {
					found = true
					break
				}
			}
			stepCtx.Require().True(
				found,
				fmt.Sprintf("Expected label with value '%s' not found in actual labels.", expectedValue),
			)
		}
	})
}
