package wrappers

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

type ResourceWrapper[R types.ResourceType, M types.MetadataType, S types.SpecType, T types.StatusType] interface {
	GetResource() *R
	GetLabels() schema.Labels
	GetAnnotations() schema.Annotations
	GetExtensions() schema.Extensions
	GetMetadata() *M
	GetSpec() *S
	GetStatus() *T
}

type BaseResourceWrapper[R types.ResourceType] struct {
	resource *R
}

func newBaseResourceWrapper[R types.ResourceType](resource *R) *BaseResourceWrapper[R] {
	return &BaseResourceWrapper[R]{
		resource: resource,
	}
}
