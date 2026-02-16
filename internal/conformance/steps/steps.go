package steps

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

type stepFuncResponse[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType] struct {
	resource *R
	labels   schema.Labels
	metadata *M
	spec     E
	status   *S
}

func newStepFuncResponse[R types.ResourceType, M types.MetadataType, E types.SpecType, S types.StatusType](
	resource *R, labels schema.Labels, metadata *M, spec E, status *S,
) *stepFuncResponse[R, M, E, S] {
	return &stepFuncResponse[R, M, E, S]{
		resource: resource,
		labels:   labels,
		metadata: metadata,
		spec:     spec,
		status:   status,
	}
}

type ResponseExpects[M types.MetadataType, E types.SpecType] struct {
	Labels        schema.Labels
	Metadata      *M
	Spec          *E
	ResourceState schema.ResourceState
}
