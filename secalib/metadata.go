package secalib

import (
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func NewGlobalTenantResourceMetadata(name, provider, resource, apiVersion, kind, tenant string) *schema.GlobalTenantResourceMetadata {
	return &schema.GlobalTenantResourceMetadata{
		Name:       name,
		Provider:   provider,
		Resource:   resource,
		ApiVersion: apiVersion,
		Kind:       schema.GlobalTenantResourceMetadataKind(kind),
		Tenant:     tenant,
	}
}

func NewRegionalResourceMetadata(name, provider, resource, apiVersion, kind, tenant, region string) *schema.RegionalResourceMetadata {
	return &schema.RegionalResourceMetadata{
		Name:       name,
		Provider:   provider,
		Resource:   resource,
		ApiVersion: apiVersion,
		Kind:       schema.RegionalResourceMetadataKind(kind),
		Tenant:     tenant,
		Region:     region,
	}
}

func NewRegionalWorkspaceResourceMetadata(name, provider, resource, apiVersion, kind, tenant, workspace, region string) *schema.RegionalWorkspaceResourceMetadata {
	return &schema.RegionalWorkspaceResourceMetadata{
		Name:       name,
		Provider:   provider,
		Resource:   resource,
		ApiVersion: apiVersion,
		Kind:       schema.RegionalWorkspaceResourceMetadataKind(kind),
		Tenant:     tenant,
		Workspace:  workspace,
		Region:     region,
	}
}

func NewRegionalNetworkResourceMetadata(name, provider, resource, apiVersion, kind, tenant, workspace, network, region string) *schema.RegionalNetworkResourceMetadata {
	return &schema.RegionalNetworkResourceMetadata{
		Name:       name,
		Provider:   provider,
		Resource:   resource,
		ApiVersion: apiVersion,
		Kind:       schema.RegionalNetworkResourceMetadataKind(kind),
		Tenant:     tenant,
		Workspace:  workspace,
		Network:    network,
		Region:     region,
	}
}
