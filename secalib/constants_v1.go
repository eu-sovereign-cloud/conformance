package secalib

const (
	// Providers
	AuthorizationProviderV1 = "seca.authorization/" + ApiVersion1
	RegionProviderV1        = "seca.region/" + ApiVersion1
	WorkspaceProviderV1     = "seca.workspace/" + ApiVersion1
	ComputeProviderV1       = "seca.compute/" + ApiVersion1
	StorageProviderV1       = "seca.storage/" + ApiVersion1
	NetworkProviderV1       = "seca.network/" + ApiVersion1

	// Base URLs
	urlProvidersPrefixV1 = "/providers/"

	RoleURLV1            = urlProvidersPrefixV1 + AuthorizationProviderV1 + "/" + RoleResource
	RoleAssignmentURLV1  = urlProvidersPrefixV1 + AuthorizationProviderV1 + "/" + RoleAssignmentResource
	RegionsURLV1         = urlProvidersPrefixV1 + RegionProviderV1 + "/regions"
	RegionURLV1          = urlProvidersPrefixV1 + RegionProviderV1 + "/" + RegionResource
	WorkspaceURLV1       = urlProvidersPrefixV1 + WorkspaceProviderV1 + "/" + WorkspaceResource
	InstanceSkuURLV1     = urlProvidersPrefixV1 + ComputeProviderV1 + "/" + SkuResource
	InstanceURLV1        = urlProvidersPrefixV1 + ComputeProviderV1 + "/" + InstanceResource
	StorageSkuURLV1      = urlProvidersPrefixV1 + StorageProviderV1 + "/" + SkuResource
	BlockStorageURLV1    = urlProvidersPrefixV1 + StorageProviderV1 + "/" + BlockStorageResource
	ImageURLV1           = urlProvidersPrefixV1 + StorageProviderV1 + "/" + ImageResource
	NetworkSkuURLV1      = urlProvidersPrefixV1 + NetworkProviderV1 + "/" + SkuResource
	NetworkURLV1         = urlProvidersPrefixV1 + NetworkProviderV1 + "/" + NetworkResource
	InternetGatewayURLV1 = urlProvidersPrefixV1 + NetworkProviderV1 + "/" + InternetGatewayResource
	NicURLV1             = urlProvidersPrefixV1 + NetworkProviderV1 + "/" + NicResource
	PublicIpURLV1        = urlProvidersPrefixV1 + NetworkProviderV1 + "/" + PublicIpResource
	RouteTableURLV1      = urlProvidersPrefixV1 + NetworkProviderV1 + "/" + RouteTableResource
	SubnetURLV1          = urlProvidersPrefixV1 + NetworkProviderV1 + "/" + SubnetResource
	SecurityGroupURLV1   = urlProvidersPrefixV1 + NetworkProviderV1 + "/" + SecurityGroupResource
)
