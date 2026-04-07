package steps

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type StepResponseExpects[M types.MetadataType, E types.SpecType] struct {
	Labels         schema.Labels
	Metadata       *M
	Spec           *E
	ResourceStates []schema.ResourceState
}

type StepCreator interface {
	WithNewStep(stepName string, step func(sCtx provider.StepCtx), params ...*allure.Parameter)
}
