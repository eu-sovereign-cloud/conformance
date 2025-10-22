package secalib

const (
	// Providers
	AuthorizationProviderV1 = "seca.authorization/" + ApiVersion1
	WorkspaceProviderV1     = "seca.workspace/" + ApiVersion1
	ComputeProviderV1       = "seca.compute/" + ApiVersion1
	StorageProviderV1       = "seca.storage/" + ApiVersion1
	NetworkProviderV1       = "seca.network/" + ApiVersion1

	// Base URLs
	RoleURLV1            = urlProvidersPrefix + AuthorizationProviderV1 + "/" + RoleResource
	RoleAssignmentURLV1  = urlProvidersPrefix + AuthorizationProviderV1 + "/" + RoleAssignmentResource
	WorkspaceURLV1       = urlProvidersPrefix + WorkspaceProviderV1 + "/" + WorkspaceResource
	InstanceSkuURLV1     = urlProvidersPrefix + ComputeProviderV1 + "/" + SkuResource
	InstanceURLV1        = urlProvidersPrefix + ComputeProviderV1 + "/" + InstanceResource
	StorageSkuURLV1      = urlProvidersPrefix + StorageProviderV1 + "/" + SkuResource
	BlockStorageURLV1    = urlProvidersPrefix + StorageProviderV1 + "/" + BlockStorageResource
	ImageURLV1           = urlProvidersPrefix + StorageProviderV1 + "/" + ImageResource
	NetworkSkuURLV1      = urlProvidersPrefix + NetworkProviderV1 + "/" + SkuResource
	NetworkURLV1         = urlProvidersPrefix + NetworkProviderV1 + "/" + NetworkResource
	InternetGatewayURLV1 = urlProvidersPrefix + NetworkProviderV1 + "/" + InternetGatewayResource
	NicURLV1             = urlProvidersPrefix + NetworkProviderV1 + "/" + NicResource
	PublicIpURLV1        = urlProvidersPrefix + NetworkProviderV1 + "/" + PublicIpResource
	RouteTableURLV1      = urlProvidersPrefix + NetworkProviderV1 + "/" + RouteTableResource
	SubnetURLV1          = urlProvidersPrefix + NetworkProviderV1 + "/" + SubnetResource
	SecurityGroupURLV1   = urlProvidersPrefix + NetworkProviderV1 + "/" + SecurityGroupResource
)
