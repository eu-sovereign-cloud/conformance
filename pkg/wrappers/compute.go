package wrappers

import "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

type InstanceWrapper struct {
	*BaseResourceWrapper[schema.Instance]
}

func NewInstanceWrapper(resource *schema.Instance) *InstanceWrapper {
	return &InstanceWrapper{
		BaseResourceWrapper: newBaseResourceWrapper(resource),
	}
}

func (wrapper *InstanceWrapper) GetResource() *schema.Instance {
	return wrapper.resource
}

func (wrapper *InstanceWrapper) GetLabels() schema.Labels {
	return wrapper.resource.Labels
}

func (wrapper *InstanceWrapper) GetAnnotations() schema.Annotations {
	return wrapper.resource.Annotations
}

func (wrapper *InstanceWrapper) GetExtensions() schema.Extensions {
	return wrapper.resource.Extensions
}

func (wrapper *InstanceWrapper) GetMetadata() *schema.RegionalWorkspaceResourceMetadata {
	return wrapper.resource.Metadata
}

func (wrapper *InstanceWrapper) GetSpec() *schema.InstanceSpec {
	return &wrapper.resource.Spec
}

func (wrapper *InstanceWrapper) GetStatus() *schema.InstanceStatus {
	return wrapper.resource.Status
}
