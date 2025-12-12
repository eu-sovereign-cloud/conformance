package builders

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Instance

type InstanceMetadataBuilder struct {
	*regionalWorkspaceResourceMetadataBuilder[InstanceMetadataBuilder]
}

func NewInstanceMetadataBuilder() *InstanceMetadataBuilder {
	builder := &InstanceMetadataBuilder{}
	builder.regionalWorkspaceResourceMetadataBuilder = newRegionalWorkspaceResourceMetadataBuilder(builder)
	return builder
}

func (builder *InstanceMetadataBuilder) Build() (*schema.RegionalWorkspaceResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindInstance).build()
	if err != nil {
		return nil, err
	}

	resource := generators.GenerateInstanceResource(builder.metadata.Tenant, builder.metadata.Workspace, builder.metadata.Name)
	metadata.Resource = resource

	return metadata, nil
}

type InstanceListMetadataBuilder struct {
	*responseMetadataBuilder
}

func NewInstanceListMetadataBuilder() *InstanceListMetadataBuilder {
	return &InstanceListMetadataBuilder{
		responseMetadataBuilder: NewResponseMetadataBuilder(),
	}
}

func (builder *InstanceListMetadataBuilder) Build(tenant, workspace, name string) (*schema.ResponseMetadata, error) {
	metadata, err := builder.build()
	if err != nil {
		return nil, err
	}

	metadata.Resource = generators.GenerateInstanceResource(tenant, workspace, name)
	return metadata, nil
}

type InstanceBuilder struct {
	*regionalWorkspaceResourceBuilder[InstanceBuilder, schema.InstanceSpec]
	metadata *InstanceMetadataBuilder
	labels   schema.Labels
	spec     *schema.InstanceSpec
}

func NewInstanceBuilder() *InstanceBuilder {
	builder := &InstanceBuilder{
		metadata: NewInstanceMetadataBuilder(),
		spec:     &schema.InstanceSpec{},
	}

	builder.regionalWorkspaceResourceBuilder = newRegionalWorkspaceResourceBuilder(newRegionalWorkspaceResourceBuilderParams[InstanceBuilder, schema.InstanceSpec]{
		newGlobalResourceBuilderParams: &newGlobalResourceBuilderParams[InstanceBuilder, schema.InstanceSpec]{
			parent:        builder,
			setName:       func(name string) { builder.metadata.setName(name) },
			setProvider:   func(provider string) { builder.metadata.setProvider(provider) },
			setApiVersion: func(apiVersion string) { builder.metadata.setApiVersion(apiVersion) },
			setLabels:     func(labels schema.Labels) { builder.labels = labels },
			setSpec:       func(spec *schema.InstanceSpec) { builder.spec = spec },
		},
		setTenant:    func(tenant string) { builder.metadata.Tenant(tenant) },
		setWorkspace: func(workspace string) { builder.metadata.Workspace(workspace) },
		setRegion:    func(region string) { builder.metadata.Region(region) },
	})

	return builder
}

func (builder *InstanceBuilder) validateSpec() error {
	if err := validateRequired(builder.validator,
		builder.spec,
		builder.spec.SkuRef,
		builder.spec.Zone,
		builder.spec.BootVolume,
		builder.spec.BootVolume.DeviceRef,
	); err != nil {
		return err
	}

	return nil
}

func (builder *InstanceBuilder) Build() (*schema.Instance, error) {
	if err := builder.validateSpec(); err != nil {
		return nil, err
	}

	metadata, err := builder.metadata.Build()
	if err != nil {
		return nil, err
	}

	return &schema.Instance{
		Metadata: metadata,
		Labels:   builder.labels,
		Spec:     *builder.spec,
		Status:   &schema.InstanceStatus{},
	}, nil
}
