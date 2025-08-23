package secalib

const (
	// Providers
	AuthorizationProviderV1 = "seca.authorization/" + ApiVersion1
	WorkspaceProviderV1     = "seca.workspace/" + ApiVersion1
	ComputeProviderV1       = "seca.compute/" + ApiVersion1
	StorageProviderV1       = "seca.storage/" + ApiVersion1

	// Base URLs
	RoleURLV1           = "/providers/" + AuthorizationProviderV1 + "/" + RoleResource
	RoleAssignmentURLV1 = "/providers/" + AuthorizationProviderV1 + "/" + RoleAssignmentResource
	WorkspaceURLV1      = "/providers/" + WorkspaceProviderV1 + "/" + WorkspaceResource
	InstanceSkuURLV1    = "/providers/" + ComputeProviderV1 + "/" + SkuResource
	InstanceURLV1       = "/providers/" + ComputeProviderV1 + "/" + InstanceResource
	StorageSkuURLV1     = "/providers/" + StorageProviderV1 + "/" + SkuResource
	BlockStorageURLV1   = "/providers/" + StorageProviderV1 + "/" + BlockStorageResource
	ImageURLV1          = "/providers/" + StorageProviderV1 + "/" + ImageResource
)
