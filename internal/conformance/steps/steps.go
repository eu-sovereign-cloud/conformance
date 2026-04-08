package steps

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type ResponseExpects[M types.MetadataType, E types.SpecType] struct {
	Labels         schema.Labels
	Annotations    schema.Annotations
	Extensions     schema.Extensions
	Metadata       *M
	Spec           *E
	ResourceStates []schema.ResourceState
}

type ResponseExpectsWithCondition[M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	Labels         schema.Labels
	Extensions     schema.Extensions
	Annotations    schema.Annotations
	Metadata       *M
	Spec           *E
	ResourceStatus S
}

type StepCreator interface {
	WithNewStep(stepName string, step func(sCtx provider.StepCtx), params ...*allure.Parameter)
}
