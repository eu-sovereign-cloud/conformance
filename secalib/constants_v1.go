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
	RoleURLV1                = UrlProvidersPrefix + AuthorizationProviderV1 + "/" + RoleResource
	RoleListURLV1            = UrlProvidersPrefix + AuthorizationProviderV1 + "/" + RolesResource
	RoleAssignmentURLV1      = UrlProvidersPrefix + AuthorizationProviderV1 + "/" + RoleAssignmentResource
	RoleAssignmentListURLV1  = UrlProvidersPrefix + AuthorizationProviderV1 + "/" + RoleAssignmentsResource
	RegionsURLV1             = UrlProvidersPrefix + RegionProviderV1 + "/regions"
	RegionURLV1              = UrlProvidersPrefix + RegionProviderV1 + "/" + RegionResource
	WorkspaceURLV1           = UrlProvidersPrefix + WorkspaceProviderV1 + "/" + WorkspaceResource
	WorkspaceListURLV1       = UrlProvidersPrefix + WorkspaceProviderV1 + "/" + WorkspaceListResource
	InstanceSkuURLV1         = UrlProvidersPrefix + ComputeProviderV1 + "/" + SkuResource
	InstanceURLV1            = UrlProvidersPrefix + ComputeProviderV1 + "/" + InstanceResource
	InstanceListURLV1        = UrlProvidersPrefix + ComputeProviderV1 + "/" + InstanceListResource
	StorageSkuURLV1          = UrlProvidersPrefix + StorageProviderV1 + "/" + SkuResource
	StorageSkuListURLV1      = UrlProvidersPrefix + StorageProviderV1 + "/" + SkuListResource
	BlockStorageURLV1        = UrlProvidersPrefix + StorageProviderV1 + "/" + BlockStorageResource
	BlockStorageListURLV1    = UrlProvidersPrefix + StorageProviderV1 + "/" + BlockStorageListResource
	ImageURLV1               = UrlProvidersPrefix + StorageProviderV1 + "/" + ImageResource
	ImageListURLV1           = UrlProvidersPrefix + StorageProviderV1 + "/" + ImageListResource
	NetworkSkuURLV1          = UrlProvidersPrefix + NetworkProviderV1 + "/" + SkuResource
	NetworkURLV1             = UrlProvidersPrefix + NetworkProviderV1 + "/" + NetworkResource
	NetworkListURLV1         = UrlProvidersPrefix + NetworkProviderV1 + "/" + NetworkListResource
	InternetGatewayURLV1     = UrlProvidersPrefix + NetworkProviderV1 + "/" + InternetGatewayResource
	InternetGatewayListURLV1 = UrlProvidersPrefix + NetworkProviderV1 + "/" + InternetGatewayListResource
	NicURLV1                 = UrlProvidersPrefix + NetworkProviderV1 + "/" + NicResource
	NicListURLV1             = UrlProvidersPrefix + NetworkProviderV1 + "/" + NicListResource
	PublicIpURLV1            = UrlProvidersPrefix + NetworkProviderV1 + "/" + PublicIpResource
	PublicIpListURLV1        = UrlProvidersPrefix + NetworkProviderV1 + "/" + PublicIpListResource
	RouteTableURLV1          = UrlProvidersPrefix + NetworkProviderV1 + "/" + RouteTableResource
	RouteTableListURLV1      = UrlProvidersPrefix + NetworkProviderV1 + "/" + RouteTableListResource
	SubnetURLV1              = UrlProvidersPrefix + NetworkProviderV1 + "/" + SubnetResource
	SubnetListURLV1          = UrlProvidersPrefix + NetworkProviderV1 + "/" + SubnetListResource
	SecurityGroupURLV1       = UrlProvidersPrefix + NetworkProviderV1 + "/" + SecurityGroupResource
	SecurityGroupListURLV1   = UrlProvidersPrefix + NetworkProviderV1 + "/" + SecurityGroupListResource
)
