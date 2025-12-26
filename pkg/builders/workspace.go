package builders

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// workspace

type WorkspaceMetadataBuilder struct {
	*regionalResourceMetadataBuilder[WorkspaceMetadataBuilder]
}

func NewWorkspaceMetadataBuilder() *WorkspaceMetadataBuilder {
	builder := &WorkspaceMetadataBuilder{}
	builder.regionalResourceMetadataBuilder = newRegionalResourceMetadataBuilder(builder)
	return builder
}

func (builder *WorkspaceMetadataBuilder) Build() (*schema.RegionalResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalResourceMetadataKindResourceKindWorkspace).build()
	if err != nil {
		return nil, err
	}

	resource := generators.GenerateWorkspaceResource(builder.metadata.Tenant, builder.metadata.Name)
	metadata.Resource = resource

	return metadata, nil
}

type WorkspaceBuilder struct {
	*regionalResourceBuilder[WorkspaceBuilder, schema.WorkspaceSpec]
	metadata *WorkspaceMetadataBuilder
	labels   schema.Labels
}

func NewWorkspaceBuilder() *WorkspaceBuilder {
	builder := &WorkspaceBuilder{
		metadata: NewWorkspaceMetadataBuilder(),
	}

	builder.regionalResourceBuilder = newRegionalResourceBuilder(newRegionalResourceBuilderParams[WorkspaceBuilder, schema.WorkspaceSpec]{
		newGlobalResourceBuilderParams: &newGlobalResourceBuilderParams[WorkspaceBuilder, schema.WorkspaceSpec]{
			parent:        builder,
			setName:       func(name string) { builder.metadata.setName(name) },
			setProvider:   func(provider string) { builder.metadata.setProvider(provider) },
			setApiVersion: func(apiVersion string) { builder.metadata.setApiVersion(apiVersion) },
			setLabels:     func(labels schema.Labels) { builder.labels = labels },
		},
		setTenant: func(tenant string) { builder.metadata.Tenant(tenant) },
		setRegion: func(region string) { builder.metadata.Region(region) },
	})

	return builder
}

func (builder *WorkspaceBuilder) Build() (*schema.Workspace, error) {
	metadata, err := builder.metadata.Build()
	if err != nil {
		return nil, err
	}

	return &schema.Workspace{
		Metadata: metadata,
		Labels:   builder.labels,
		Spec:     schema.WorkspaceSpec{},
		Status:   &schema.WorkspaceStatus{},
	}, nil
}
