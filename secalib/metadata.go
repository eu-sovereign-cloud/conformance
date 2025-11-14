package secalib

import (
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

type metadataBuilder struct {
	validator *Validator
}

func newMetadataBuilder() *metadataBuilder {
	return &metadataBuilder{
		validator: newValidator(),
	}
}

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
func (b *GlobalTenantResourceMetadataBuilder) Name(name string) *GlobalTenantResourceMetadataBuilder {
	b.metadata.Name = name
	return b
}
func (b *GlobalTenantResourceMetadataBuilder) Provider(provider string) *GlobalTenantResourceMetadataBuilder {
	b.metadata.Provider = provider
	return b
}
func (b *GlobalTenantResourceMetadataBuilder) Resource(resource string) *GlobalTenantResourceMetadataBuilder {
	b.metadata.Resource = resource
	return b
}
func (b *GlobalTenantResourceMetadataBuilder) ApiVersion(apiVersion string) *GlobalTenantResourceMetadataBuilder {
	b.metadata.ApiVersion = apiVersion
	return b
}
func (b *GlobalTenantResourceMetadataBuilder) Kind(kind string) *GlobalTenantResourceMetadataBuilder {
	b.metadata.Kind = schema.GlobalTenantResourceMetadataKind(kind)
	return b
}
func (b *GlobalTenantResourceMetadataBuilder) Tenant(tenant string) *GlobalTenantResourceMetadataBuilder {
	b.metadata.Tenant = tenant
	return b
}
func (b *GlobalTenantResourceMetadataBuilder) Build() (*schema.GlobalTenantResourceMetadata, error) {
	if err := b.validator.ValidateRequireds(
		[]any{
			b.metadata,
			b.metadata.Name,
			b.metadata.Provider,
			b.metadata.Resource,
			b.metadata.ApiVersion,
			b.metadata.Kind,
			b.metadata.Tenant,
		},
	); err != nil {
		return nil, err
	}

	return b.metadata, nil
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
func (b *RegionalResourceMetadataBuilder) Name(name string) *RegionalResourceMetadataBuilder {
	b.metadata.Name = name
	return b
}
func (b *RegionalResourceMetadataBuilder) Provider(provider string) *RegionalResourceMetadataBuilder {
	b.metadata.Provider = provider
	return b
}
func (b *RegionalResourceMetadataBuilder) Resource(resource string) *RegionalResourceMetadataBuilder {
	b.metadata.Resource = resource
	return b
}
func (b *RegionalResourceMetadataBuilder) ApiVersion(apiVersion string) *RegionalResourceMetadataBuilder {
	b.metadata.ApiVersion = apiVersion
	return b
}
func (b *RegionalResourceMetadataBuilder) Kind(kind string) *RegionalResourceMetadataBuilder {
	b.metadata.Kind = schema.RegionalResourceMetadataKind(kind)
	return b
}
func (b *RegionalResourceMetadataBuilder) Tenant(tenant string) *RegionalResourceMetadataBuilder {
	b.metadata.Tenant = tenant
	return b
}
func (b *RegionalResourceMetadataBuilder) Region(region string) *RegionalResourceMetadataBuilder {
	b.metadata.Region = region
	return b
}
func (b *RegionalResourceMetadataBuilder) Build() (*schema.RegionalResourceMetadata, error) {
	if err := b.validator.ValidateRequireds(
		[]any{
			b.metadata,
			b.metadata.Name,
			b.metadata.Provider,
			b.metadata.Resource,
			b.metadata.ApiVersion,
			b.metadata.Kind,
			b.metadata.Tenant,
			b.metadata.Region,
		},
	); err != nil {
		return nil, err
	}

	return b.metadata, nil
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
func (b *RegionalWorkspaceResourceMetadataBuilder) Name(name string) *RegionalWorkspaceResourceMetadataBuilder {
	b.metadata.Name = name
	return b
}
func (b *RegionalWorkspaceResourceMetadataBuilder) Provider(provider string) *RegionalWorkspaceResourceMetadataBuilder {
	b.metadata.Provider = provider
	return b
}
func (b *RegionalWorkspaceResourceMetadataBuilder) Resource(resource string) *RegionalWorkspaceResourceMetadataBuilder {
	b.metadata.Resource = resource
	return b
}
func (b *RegionalWorkspaceResourceMetadataBuilder) ApiVersion(apiVersion string) *RegionalWorkspaceResourceMetadataBuilder {
	b.metadata.ApiVersion = apiVersion
	return b
}
func (b *RegionalWorkspaceResourceMetadataBuilder) Kind(kind string) *RegionalWorkspaceResourceMetadataBuilder {
	b.metadata.Kind = schema.RegionalWorkspaceResourceMetadataKind(kind)
	return b
}
func (b *RegionalWorkspaceResourceMetadataBuilder) Tenant(tenant string) *RegionalWorkspaceResourceMetadataBuilder {
	b.metadata.Tenant = tenant
	return b
}
func (b *RegionalWorkspaceResourceMetadataBuilder) Workspace(workspace string) *RegionalWorkspaceResourceMetadataBuilder {
	b.metadata.Workspace = workspace
	return b
}
func (b *RegionalWorkspaceResourceMetadataBuilder) Region(region string) *RegionalWorkspaceResourceMetadataBuilder {
	b.metadata.Region = region
	return b
}

func (b *RegionalWorkspaceResourceMetadataBuilder) Build() (*schema.RegionalWorkspaceResourceMetadata, error) {
	if err := b.validator.ValidateRequireds(
		[]any{
			b.metadata,
			b.metadata.Name,
			b.metadata.Provider,
			b.metadata.Resource,
			b.metadata.ApiVersion,
			b.metadata.Kind,
			b.metadata.Tenant,
			b.metadata.Workspace,
			b.metadata.Region,
		},
	); err != nil {
		return nil, err
	}

	return b.metadata, nil
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
func (b *RegionalNetworkResourceMetadataBuilder) Name(name string) *RegionalNetworkResourceMetadataBuilder {
	b.metadata.Name = name
	return b
}
func (b *RegionalNetworkResourceMetadataBuilder) Provider(provider string) *RegionalNetworkResourceMetadataBuilder {
	b.metadata.Provider = provider
	return b
}
func (b *RegionalNetworkResourceMetadataBuilder) Resource(resource string) *RegionalNetworkResourceMetadataBuilder {
	b.metadata.Resource = resource
	return b
}
func (b *RegionalNetworkResourceMetadataBuilder) ApiVersion(apiVersion string) *RegionalNetworkResourceMetadataBuilder {
	b.metadata.ApiVersion = apiVersion
	return b
}
func (b *RegionalNetworkResourceMetadataBuilder) Kind(kind string) *RegionalNetworkResourceMetadataBuilder {
	b.metadata.Kind = schema.RegionalNetworkResourceMetadataKind(kind)
	return b
}
func (b *RegionalNetworkResourceMetadataBuilder) Tenant(tenant string) *RegionalNetworkResourceMetadataBuilder {
	b.metadata.Tenant = tenant
	return b
}
func (b *RegionalNetworkResourceMetadataBuilder) Workspace(workspace string) *RegionalNetworkResourceMetadataBuilder {
	b.metadata.Workspace = workspace
	return b
}
func (b *RegionalNetworkResourceMetadataBuilder) Network(network string) *RegionalNetworkResourceMetadataBuilder {
	b.metadata.Network = network
	return b
}
func (b *RegionalNetworkResourceMetadataBuilder) Region(region string) *RegionalNetworkResourceMetadataBuilder {
	b.metadata.Region = region
	return b
}
func (b *RegionalNetworkResourceMetadataBuilder) Build() (*schema.RegionalNetworkResourceMetadata, error) {
	if err := b.validator.ValidateRequireds(
		[]any{
			b.metadata,
			b.metadata.Name,
			b.metadata.Provider,
			b.metadata.Resource,
			b.metadata.ApiVersion,
			b.metadata.Kind,
			b.metadata.Tenant,
			b.metadata.Workspace,
			b.metadata.Network,
			b.metadata.Region,
		},
	); err != nil {
		return nil, err
	}

	return b.metadata, nil
}
