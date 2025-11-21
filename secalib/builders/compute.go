package builders

import (
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Instance

type InstanceBuilder struct {
	*resourceBuilder
	metadata *RegionalWorkspaceResourceMetadataBuilder
	spec     *schema.InstanceSpec
}

func NewInstanceBuilder() *InstanceBuilder {
	return &InstanceBuilder{
		resourceBuilder: newResourceBuilder(),
		metadata:        NewRegionalWorkspaceResourceMetadataBuilder(),
		spec:            &schema.InstanceSpec{},
	}
}

func (builder *InstanceBuilder) Name(name string) *InstanceBuilder {
	builder.metadata.Name(name)
	return builder
}

func (builder *InstanceBuilder) Provider(provider string) *InstanceBuilder {
	builder.metadata.Provider(provider)
	return builder
}

func (builder *InstanceBuilder) Resource(resource string) *InstanceBuilder {
	builder.metadata.Resource(resource)
	return builder
}

func (builder *InstanceBuilder) ApiVersion(apiVersion string) *InstanceBuilder {
	builder.metadata.ApiVersion(apiVersion)
	return builder
}

func (builder *InstanceBuilder) Tenant(tenant string) *InstanceBuilder {
	builder.metadata.Tenant(tenant)
	return builder
}

func (builder *InstanceBuilder) Workspace(workspace string) *InstanceBuilder {
	builder.metadata.Workspace(workspace)
	return builder
}

func (builder *InstanceBuilder) Region(region string) *InstanceBuilder {
	builder.metadata.Region(region)
	return builder
}

func (builder *InstanceBuilder) Spec(spec *schema.InstanceSpec) *InstanceBuilder {
	builder.spec = spec
	return builder
}

func (builder *InstanceBuilder) BuildResponse() (*schema.Instance, error) {
	medatata, err := builder.metadata.Kind(secalib.InstanceKind).BuildResponse()
	if err != nil {
		return nil, err
	}

	// Validate the spec
	if err := builder.validator.ValidateRequired(
		builder.spec,
		builder.spec.SkuRef,
		builder.spec.Zone,
		builder.spec.BootVolume,
		builder.spec.BootVolume.DeviceRef,
	); err != nil {
		return nil, err
	}

	return &schema.Instance{
		Metadata: medatata,
		Labels:   schema.Labels{},
		Spec:     *builder.spec,
		Status:   &schema.InstanceStatus{},
	}, nil
}
