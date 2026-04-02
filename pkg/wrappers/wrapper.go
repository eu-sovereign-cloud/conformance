package wrappers

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/types"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Resource Wrapper

type ResourceWrapper[R types.ResourceType, M types.MetadataType, S types.SpecType, T types.StatusType] interface {
	GetResource() *R
	GetLabels() schema.Labels
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

// Global Resource Wrapper

type GlobalResourceWrapper[R types.ResourceType, M types.MetadataType, S types.SpecType] interface {
	GetResource() *R
	GetMetadata() *M
	GetSpec() *S
}

type BaseGlobalResourceWrapper[R types.ResourceType] struct {
	resource *R
}

func newBaseGlobalResourceWrapper[R types.ResourceType](resource *R) *BaseGlobalResourceWrapper[R] {
	return &BaseGlobalResourceWrapper[R]{
		resource: resource,
	}
}
