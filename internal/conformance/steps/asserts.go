package steps

import (
	"fmt"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func requireNoError(sCtx provider.StepCtx, err error) {
	sCtx.WithNewStep("Verify no error", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters("error", fmt.Sprintf("%v", err))
		stepCtx.Require().NoError(err, "Should not return an error")
	})
}

func requirePreConditionFailedError(sCtx provider.StepCtx, err error) {
	sCtx.WithNewStep("Verify error returned", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters("error", fmt.Sprintf("%v", err))
		stepCtx.Require().Error(err, "Should return an error")
		stepCtx.Require().Equal(err.Error(), secapi.ErrRequestPreconditionFailed.Error(), "expected error message")
	})
}

func requireNotNilResponse(sCtx provider.StepCtx, resp any) {
	sCtx.WithNewStep("Verify not nil response", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters("response", fmt.Sprintf("%v", resp))
		stepCtx.Require().NotNil(resp, "Should be not nil")
	})
}

func requireNotEmptyResponse[R types.ResourceType](sCtx provider.StepCtx, resp []*R) {
	sCtx.WithNewStep("Verify response length", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters("response", fmt.Sprintf("%v", resp))
		stepCtx.Require().NotNil(resp, "Should be not nil")
		stepCtx.Require().GreaterOrEqual(len(resp), 1, "Should have length greater than 1")
	})
}

func requireValidResponseMetadata(sCtx provider.StepCtx, metadata schema.ResponseMetadata, expects schema.ResponseMetadata) {
	sCtx.WithNewStep("Verify response metadata", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"provider", metadata.Provider,
			"resource", metadata.Resource,
			"verb", metadata.Verb,
			"expectedProvider", expects.Provider,
			"expectedResource", expects.Resource,
			"expectedVerb", expects.Verb,
		)
		stepCtx.Require().NotEmpty(metadata.Provider, "ResponseMetadata.Provider should not be empty")
		stepCtx.Require().NotEmpty(metadata.Resource, "ResponseMetadata.Resource should not be empty")
		stepCtx.Require().NotEmpty(metadata.Verb, "ResponseMetadata.Verb should not be empty")
		if expects.Provider != "" {
			stepCtx.Require().Equal(expects.Provider, metadata.Provider, "ResponseMetadata.Provider mismatch")
		}
		if expects.Resource != "" {
			stepCtx.Require().Equal(expects.Resource, metadata.Resource, "ResponseMetadata.Resource mismatch")
		}
		if expects.Verb != "" {
			stepCtx.Require().Equal(expects.Verb, metadata.Verb, "ResponseMetadata.Verb mismatch")
		}
	})
}
