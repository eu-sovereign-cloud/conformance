package builders

import "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

// GlobalTenantResourceMetadataBuilder

type GlobalTenantResourceMetadataBuilder struct {
	*metadataBuilder
	metadata *schema.GlobalTenantResourceMetadata
}

func NewGlobalTenantResourceMetadataBuilder() *GlobalTenantResourceMetadataBuilder {
	return &GlobalTenantResourceMetadataBuilder{
		metadataBuilder: newMetadataBuilder(),
		metadata:        &schema.GlobalTenantResourceMetadata{},
	}
}
func (builder *GlobalTenantResourceMetadataBuilder) Name(name string) *GlobalTenantResourceMetadataBuilder {
	builder.metadata.Name = name
	return builder
}
func (builder *GlobalTenantResourceMetadataBuilder) Provider(provider string) *GlobalTenantResourceMetadataBuilder {
	builder.metadata.Provider = provider
	return builder
}
func (builder *GlobalTenantResourceMetadataBuilder) Resource(resource string) *GlobalTenantResourceMetadataBuilder {
	builder.metadata.Resource = resource
	return builder
}
func (builder *GlobalTenantResourceMetadataBuilder) ApiVersion(apiVersion string) *GlobalTenantResourceMetadataBuilder {
	builder.metadata.ApiVersion = apiVersion
	return builder
}
func (builder *GlobalTenantResourceMetadataBuilder) Kind(kind string) *GlobalTenantResourceMetadataBuilder {
	builder.metadata.Kind = schema.GlobalTenantResourceMetadataKind(kind)
	return builder
}
func (builder *GlobalTenantResourceMetadataBuilder) Tenant(tenant string) *GlobalTenantResourceMetadataBuilder {
	builder.metadata.Tenant = tenant
	return builder
}
func (builder *GlobalTenantResourceMetadataBuilder) BuildResponse() (*schema.GlobalTenantResourceMetadata, error) {
	if err := builder.validator.ValidateRequireds(
		[]any{
			builder.metadata,
			builder.metadata.Name,
			builder.metadata.Provider,
			builder.metadata.Resource,
			builder.metadata.ApiVersion,
			builder.metadata.Kind,
			builder.metadata.Tenant,
		},
	); err != nil {
		return nil, err
	}

	return builder.metadata, nil
}

// RegionalResourceMetadata

type RegionalResourceMetadataBuilder struct {
	*metadataBuilder
	metadata *schema.RegionalResourceMetadata
}

func NewRegionalResourceMetadataBuilder() *RegionalResourceMetadataBuilder {
	return &RegionalResourceMetadataBuilder{
		metadataBuilder: newMetadataBuilder(),
		metadata:        &schema.RegionalResourceMetadata{},
	}
}
func (builder *RegionalResourceMetadataBuilder) Name(name string) *RegionalResourceMetadataBuilder {
	builder.metadata.Name = name
	return builder
}
func (builder *RegionalResourceMetadataBuilder) Provider(provider string) *RegionalResourceMetadataBuilder {
	builder.metadata.Provider = provider
	return builder
}
func (builder *RegionalResourceMetadataBuilder) Resource(resource string) *RegionalResourceMetadataBuilder {
	builder.metadata.Resource = resource
	return builder
}
func (builder *RegionalResourceMetadataBuilder) ApiVersion(apiVersion string) *RegionalResourceMetadataBuilder {
	builder.metadata.ApiVersion = apiVersion
	return builder
}
func (builder *RegionalResourceMetadataBuilder) Kind(kind string) *RegionalResourceMetadataBuilder {
	builder.metadata.Kind = schema.RegionalResourceMetadataKind(kind)
	return builder
}
func (builder *RegionalResourceMetadataBuilder) Tenant(tenant string) *RegionalResourceMetadataBuilder {
	builder.metadata.Tenant = tenant
	return builder
}
func (builder *RegionalResourceMetadataBuilder) Region(region string) *RegionalResourceMetadataBuilder {
	builder.metadata.Region = region
	return builder
}
func (builder *RegionalResourceMetadataBuilder) BuildResponse() (*schema.RegionalResourceMetadata, error) {
	if err := builder.validator.ValidateRequireds(
		[]any{
			builder.metadata,
			builder.metadata.Name,
			builder.metadata.Provider,
			builder.metadata.Resource,
			builder.metadata.ApiVersion,
			builder.metadata.Kind,
			builder.metadata.Tenant,
			builder.metadata.Region,
		},
	); err != nil {
		return nil, err
	}

	return builder.metadata, nil
}

// RegionalWorkspaceResourceMetadata

type RegionalWorkspaceResourceMetadataBuilder struct {
	*metadataBuilder
	metadata *schema.RegionalWorkspaceResourceMetadata
}

func NewRegionalWorkspaceResourceMetadataBuilder() *RegionalWorkspaceResourceMetadataBuilder {
	return &RegionalWorkspaceResourceMetadataBuilder{
		metadataBuilder: newMetadataBuilder(),
		metadata:        &schema.RegionalWorkspaceResourceMetadata{},
	}
}
func (builder *RegionalWorkspaceResourceMetadataBuilder) Name(name string) *RegionalWorkspaceResourceMetadataBuilder {
	builder.metadata.Name = name
	return builder
}
func (builder *RegionalWorkspaceResourceMetadataBuilder) Provider(provider string) *RegionalWorkspaceResourceMetadataBuilder {
	builder.metadata.Provider = provider
	return builder
}
func (builder *RegionalWorkspaceResourceMetadataBuilder) Resource(resource string) *RegionalWorkspaceResourceMetadataBuilder {
	builder.metadata.Resource = resource
	return builder
}
func (builder *RegionalWorkspaceResourceMetadataBuilder) ApiVersion(apiVersion string) *RegionalWorkspaceResourceMetadataBuilder {
	builder.metadata.ApiVersion = apiVersion
	return builder
}
func (builder *RegionalWorkspaceResourceMetadataBuilder) Kind(kind string) *RegionalWorkspaceResourceMetadataBuilder {
	builder.metadata.Kind = schema.RegionalWorkspaceResourceMetadataKind(kind)
	return builder
}
func (builder *RegionalWorkspaceResourceMetadataBuilder) Tenant(tenant string) *RegionalWorkspaceResourceMetadataBuilder {
	builder.metadata.Tenant = tenant
	return builder
}
func (builder *RegionalWorkspaceResourceMetadataBuilder) Workspace(workspace string) *RegionalWorkspaceResourceMetadataBuilder {
	builder.metadata.Workspace = workspace
	return builder
}
func (builder *RegionalWorkspaceResourceMetadataBuilder) Region(region string) *RegionalWorkspaceResourceMetadataBuilder {
	builder.metadata.Region = region
	return builder
}

func (builder *RegionalWorkspaceResourceMetadataBuilder) BuildResponse() (*schema.RegionalWorkspaceResourceMetadata, error) {
	if err := builder.validator.ValidateRequireds(
		[]any{
			builder.metadata,
			builder.metadata.Name,
			builder.metadata.Provider,
			builder.metadata.Resource,
			builder.metadata.ApiVersion,
			builder.metadata.Kind,
			builder.metadata.Tenant,
			builder.metadata.Workspace,
			builder.metadata.Region,
		},
	); err != nil {
		return nil, err
	}

	return builder.metadata, nil
}

// RegionalNetworkResourceMetadata

type RegionalNetworkResourceMetadataBuilder struct {
	*metadataBuilder
	metadata *schema.RegionalNetworkResourceMetadata
}

func NewRegionalNetworkResourceMetadataBuilder() *RegionalNetworkResourceMetadataBuilder {
	return &RegionalNetworkResourceMetadataBuilder{
		metadataBuilder: newMetadataBuilder(),
		metadata:        &schema.RegionalNetworkResourceMetadata{},
	}
}
func (builder *RegionalNetworkResourceMetadataBuilder) Name(name string) *RegionalNetworkResourceMetadataBuilder {
	builder.metadata.Name = name
	return builder
}
func (builder *RegionalNetworkResourceMetadataBuilder) Provider(provider string) *RegionalNetworkResourceMetadataBuilder {
	builder.metadata.Provider = provider
	return builder
}
func (builder *RegionalNetworkResourceMetadataBuilder) Resource(resource string) *RegionalNetworkResourceMetadataBuilder {
	builder.metadata.Resource = resource
	return builder
}
func (builder *RegionalNetworkResourceMetadataBuilder) ApiVersion(apiVersion string) *RegionalNetworkResourceMetadataBuilder {
	builder.metadata.ApiVersion = apiVersion
	return builder
}
func (builder *RegionalNetworkResourceMetadataBuilder) Kind(kind string) *RegionalNetworkResourceMetadataBuilder {
	builder.metadata.Kind = schema.RegionalNetworkResourceMetadataKind(kind)
	return builder
}
func (builder *RegionalNetworkResourceMetadataBuilder) Tenant(tenant string) *RegionalNetworkResourceMetadataBuilder {
	builder.metadata.Tenant = tenant
	return builder
}
func (builder *RegionalNetworkResourceMetadataBuilder) Workspace(workspace string) *RegionalNetworkResourceMetadataBuilder {
	builder.metadata.Workspace = workspace
	return builder
}
func (builder *RegionalNetworkResourceMetadataBuilder) Network(network string) *RegionalNetworkResourceMetadataBuilder {
	builder.metadata.Network = network
	return builder
}
func (builder *RegionalNetworkResourceMetadataBuilder) Region(region string) *RegionalNetworkResourceMetadataBuilder {
	builder.metadata.Region = region
	return builder
}
func (builder *RegionalNetworkResourceMetadataBuilder) BuildResponse() (*schema.RegionalNetworkResourceMetadata, error) {
	if err := builder.validator.ValidateRequireds(
		[]any{
			builder.metadata,
			builder.metadata.Name,
			builder.metadata.Provider,
			builder.metadata.Resource,
			builder.metadata.ApiVersion,
			builder.metadata.Kind,
			builder.metadata.Tenant,
			builder.metadata.Workspace,
			builder.metadata.Network,
			builder.metadata.Region,
		},
	); err != nil {
		return nil, err
	}

	return builder.metadata, nil
}
