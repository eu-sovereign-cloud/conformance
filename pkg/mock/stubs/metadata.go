package stubs

import (
	"time"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// GlobalResourceMetadata

func setCreatedGlobalResourceMetadata(metadata *schema.GlobalResourceMetadata) {
	metadata.CreatedAt = time.Now()
	metadata.LastModifiedAt = time.Now()
	metadata.ResourceVersion = 1
}

// GlobalTenantResourceMetadata

func setCreatedGlobalTenantResourceMetadata(metadata *schema.GlobalTenantResourceMetadata) {
	metadata.CreatedAt = time.Now()
	metadata.LastModifiedAt = time.Now()
	metadata.ResourceVersion = 1
}

func setModifiedGlobalTenantResourceMetadata(metadata *schema.GlobalTenantResourceMetadata) {
	metadata.LastModifiedAt = time.Now()
	metadata.ResourceVersion = metadata.ResourceVersion + 1
}

// RegionalResourceMetadata

func setCreatedRegionalResourceMetadata(metadata *schema.RegionalResourceMetadata) {
	metadata.CreatedAt = time.Now()
	metadata.LastModifiedAt = time.Now()
	metadata.ResourceVersion = 1
}

func setModifiedRegionalResourceMetadata(metadata *schema.RegionalResourceMetadata) {
	metadata.LastModifiedAt = time.Now()
	metadata.ResourceVersion = metadata.ResourceVersion + 1
}

// RegionalWorkspaceResourceMetadata

func setCreatedRegionalWorkspaceResourceMetadata(metadata *schema.RegionalWorkspaceResourceMetadata) {
	metadata.CreatedAt = time.Now()
	metadata.LastModifiedAt = time.Now()
	metadata.ResourceVersion = 1
}

func setModifiedRegionalWorkspaceResourceMetadata(metadata *schema.RegionalWorkspaceResourceMetadata) {
	metadata.LastModifiedAt = time.Now()
	metadata.ResourceVersion = metadata.ResourceVersion + 1
}

// RegionalNetworkResourceMetadata

func setCreatedRegionalNetworkResourceMetadata(metadata *schema.RegionalNetworkResourceMetadata) {
	metadata.CreatedAt = time.Now()
	metadata.LastModifiedAt = time.Now()
	metadata.ResourceVersion = 1
}

func setModifiedRegionalNetworkResourceMetadata(metadata *schema.RegionalNetworkResourceMetadata) {
	metadata.LastModifiedAt = time.Now()
	metadata.ResourceVersion = metadata.ResourceVersion + 1
}
