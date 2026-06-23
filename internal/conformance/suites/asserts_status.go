package suites

import (
	"fmt"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (suite *TestSuite) VerifyStatusStateStep(ctx provider.StepCtx, expected schema.ResourceState, actual schema.ResourceState) {
	ctx.WithNewStep("Verify status state", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected, actual, "Status state should match expected")
	})
}

func (suite *TestSuite) VerifyStatusStatesStep(ctx provider.StepCtx, expected []schema.ResourceState, actual schema.ResourceState) {
	ctx.WithNewStep("Verify status state", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Contains(expected, actual, "Status state should match expected")
	})
}

func (suite *TestSuite) VerifyStatusConditionsStep(ctx provider.StepCtx, expected []schema.StatusCondition, actual []schema.StatusCondition) {
	ctx.WithNewStep("Verify status conditions", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(len(expected), len(actual), "Status conditions length should match expected")
		for i := range expected {
			stepCtx.Require().Equal(expected[i].State, actual[i].State,
				fmt.Sprintf("Condition [%d] state should match expected", i))

			stepCtx.Require().NotEmpty(actual[i].LastTransitionAt)
			stepCtx.Require().False(actual[i].LastTransitionAt.IsZero(),
				fmt.Sprintf("Condition [%d] lastTransitionAt should not be zero", i))

			if i != 0 {
				stepCtx.Require().False(actual[i].LastTransitionAt.Before(actual[i-1].LastTransitionAt),
					fmt.Sprintf("Condition [%d] lastTransitionAt should be after to previous condition's lastTransitionAt", i))
			}

		}
	})
}

func (suite *TestSuite) VerifyStatusPowerStateStep(ctx provider.StepCtx, expected schema.InstanceStatusPowerState, actual schema.InstanceStatusPowerState) {
	ctx.WithNewStep("Verify status power state", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected, actual, "Status power state should match expected")
	})
}

func (suite *TestSuite) VerifyLabelsStep(ctx provider.StepCtx, expected schema.Labels, actual schema.Labels) {
	ctx.WithNewStep("Verify label", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_labels", expected,
			"actual_labels", actual,
		)

		stepCtx.Require().Equal(expected, actual, "Labels should match expected")
	})
	ctx.WithNewStep("Verify label constraints", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"actual_labels", actual,
		)
		for key, value := range actual {
			stepCtx.Assert().LessOrEqual(len(value), 63, fmt.Sprintf("Label %q value max length should be <= 63", key))
		}
	})
}

func (suite *TestSuite) VerifyAnnotationsStep(ctx provider.StepCtx, expected schema.Annotations, actual schema.Annotations) {
	ctx.WithNewStep("Verify annotation", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_annotations", expected,
			"actual_annotations", actual,
		)

		stepCtx.Require().Equal(expected, actual, "Annotations should match expected")
	})

	ctx.WithNewStep("Verify annotation constraints", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"actual_annotations", actual,
		)
		for key, value := range actual {
			stepCtx.Assert().LessOrEqual(len(value), 1024, fmt.Sprintf("Annotation %q value max length should be <= 1024", key))
		}
	})
}

func (suite *TestSuite) VerifyExtensionsStep(ctx provider.StepCtx, expected schema.Extensions, actual schema.Extensions) {
	ctx.WithNewStep("Verify extension", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_extensions", expected,
			"actual_extensions", actual,
		)

		stepCtx.Require().Equal(expected, actual, "Extensions should match expected")
	})
}
