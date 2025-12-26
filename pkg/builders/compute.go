//nolint:dupl
package builders

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	compute "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.compute.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// InstanceSku

/// InstanceSkuIteratorBuilder

type InstanceSkuIteratorBuilder struct {
	*tenantResponseMetadataBuilder[InstanceSkuIteratorBuilder]

	items []schema.InstanceSku
}

func NewInstanceSkuIteratorBuilder() *InstanceSkuIteratorBuilder {
	builder := &InstanceSkuIteratorBuilder{}
	builder.tenantResponseMetadataBuilder = newTenantResponseMetadataBuilder(builder)
	return builder
}

func (builder *InstanceSkuIteratorBuilder) Items(items []schema.InstanceSku) *InstanceSkuIteratorBuilder {
	builder.items = items
	return builder
}

func (builder *InstanceSkuIteratorBuilder) Build() (*compute.SkuIterator, error) {
	err := builder.validate()
	if err != nil {
		return nil, err
	}

	builder.metadata.Resource = generators.GenerateSkuListResource(builder.tenant)

	return &compute.SkuIterator{
		Metadata: *builder.metadata,
		Items:    builder.items,
	}, nil
}

// Instance

/// InstanceMetadataBuilder

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

/// InstanceBuilder

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

/// InstanceIteratorBuilder

type InstanceIteratorBuilder struct {
	*workspaceResponseMetadataBuilder[InstanceIteratorBuilder]

	items []schema.Instance
}

func NewInstanceIteratorBuilder() *InstanceIteratorBuilder {
	builder := &InstanceIteratorBuilder{}
	builder.workspaceResponseMetadataBuilder = newWorkspaceResponseMetadataBuilder(builder)
	return builder
}

func (builder *InstanceIteratorBuilder) Items(items []schema.Instance) *InstanceIteratorBuilder {
	builder.items = items
	return builder
}

func (builder *InstanceIteratorBuilder) Build() (*compute.InstanceIterator, error) {
	err := builder.validate()
	if err != nil {
		return nil, err
	}

	builder.metadata.Resource = generators.GenerateInstanceListResource(builder.tenant, builder.workspace)

	return &compute.InstanceIterator{
		Metadata: *builder.metadata,
		Items:    builder.items,
	}, nil
}
