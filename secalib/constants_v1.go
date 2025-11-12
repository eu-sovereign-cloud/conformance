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
	RoleURLV1            = UrlProvidersPrefix + AuthorizationProviderV1 + "/" + RoleResource
	RolesURLV1           = UrlProvidersPrefix + AuthorizationProviderV1 + "/" + RolesResource
	RoleAssignmentURLV1  = UrlProvidersPrefix + AuthorizationProviderV1 + "/" + RoleAssignmentResource
	RoleAssignmentsURLV1 = UrlProvidersPrefix + AuthorizationProviderV1 + "/" + RoleAssignmentsResource
	RegionsURLV1         = UrlProvidersPrefix + RegionProviderV1 + "/regions"
	RegionURLV1          = UrlProvidersPrefix + RegionProviderV1 + "/" + RegionResource
	WorkspaceURLV1       = UrlProvidersPrefix + WorkspaceProviderV1 + "/" + WorkspaceResource
	InstanceSkuURLV1     = UrlProvidersPrefix + ComputeProviderV1 + "/" + SkuResource
	InstanceURLV1        = UrlProvidersPrefix + ComputeProviderV1 + "/" + InstanceResource
	StorageSkuURLV1      = UrlProvidersPrefix + StorageProviderV1 + "/" + SkuResource
	BlockStorageURLV1    = UrlProvidersPrefix + StorageProviderV1 + "/" + BlockStorageResource
	ImageURLV1           = UrlProvidersPrefix + StorageProviderV1 + "/" + ImageResource
	NetworkSkuURLV1      = UrlProvidersPrefix + NetworkProviderV1 + "/" + SkuResource
	NetworkURLV1         = UrlProvidersPrefix + NetworkProviderV1 + "/" + NetworkResource
	InternetGatewayURLV1 = UrlProvidersPrefix + NetworkProviderV1 + "/" + InternetGatewayResource
	NicURLV1             = UrlProvidersPrefix + NetworkProviderV1 + "/" + NicResource
	PublicIpURLV1        = UrlProvidersPrefix + NetworkProviderV1 + "/" + PublicIpResource
	RouteTableURLV1      = UrlProvidersPrefix + NetworkProviderV1 + "/" + RouteTableResource
	SubnetURLV1          = UrlProvidersPrefix + NetworkProviderV1 + "/" + SubnetResource
	SecurityGroupURLV1   = UrlProvidersPrefix + NetworkProviderV1 + "/" + SecurityGroupResource
)
