package wrappers

import "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

type RegionWrapper struct {
	*BaseGlobalResourceWrapper[schema.Region]
}

func NewRegionWrapper(resource *schema.Region) *RegionWrapper {
	return &RegionWrapper{
		BaseGlobalResourceWrapper: newBaseGlobalResourceWrapper(resource),
	}
}

func (wrapper *RegionWrapper) GetResource() *schema.Region {
	return wrapper.resource
}

func (wrapper *RegionWrapper) GetMetadata() *schema.GlobalResourceMetadata {
	return wrapper.resource.Metadata
}

func (wrapper *RegionWrapper) GetSpec() *schema.RegionSpec {
	return &wrapper.resource.Spec
}
