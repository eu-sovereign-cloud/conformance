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

	roleURLV1                = urlProvidersPrefix + AuthorizationProviderV1 + "/" + roleResource
	roleListURLV1            = urlProvidersPrefix + AuthorizationProviderV1 + "/" + rolesResource
	roleAssignmentURLV1      = urlProvidersPrefix + AuthorizationProviderV1 + "/" + roleAssignmentResource
	roleAssignmentListURLV1  = urlProvidersPrefix + AuthorizationProviderV1 + "/" + roleAssignmentsResource
	regionURLV1              = urlProvidersPrefix + RegionProviderV1 + "/" + regionResource
	RegionsURLV1             = urlProvidersPrefix + RegionProviderV1 + "/regions"
	workspaceURLV1           = urlProvidersPrefix + WorkspaceProviderV1 + "/" + workspaceResource
	workspaceListURLV1       = urlProvidersPrefix + WorkspaceProviderV1 + "/" + workspaceListResource
	instanceSkuURLV1         = urlProvidersPrefix + ComputeProviderV1 + "/" + skuResource
	instanceURLV1            = urlProvidersPrefix + ComputeProviderV1 + "/" + instanceResource
	instanceListURLV1        = urlProvidersPrefix + ComputeProviderV1 + "/" + instanceListResource
	storageSkuURLV1          = urlProvidersPrefix + StorageProviderV1 + "/" + skuResource
	storageSkuListURLV1      = urlProvidersPrefix + StorageProviderV1 + "/" + skuListResource
	blockStorageURLV1        = urlProvidersPrefix + StorageProviderV1 + "/" + blockStorageResource
	blockStorageListURLV1    = urlProvidersPrefix + StorageProviderV1 + "/" + blockStorageListResource
	imageURLV1               = urlProvidersPrefix + StorageProviderV1 + "/" + imageResource
	imageListURLV1           = urlProvidersPrefix + StorageProviderV1 + "/" + imageListResource
	networkSkuURLV1          = urlProvidersPrefix + NetworkProviderV1 + "/" + skuResource
	networkURLV1             = urlProvidersPrefix + NetworkProviderV1 + "/" + networkResource
	networkListURLV1         = urlProvidersPrefix + NetworkProviderV1 + "/" + networkListResource
	internetGatewayURLV1     = urlProvidersPrefix + NetworkProviderV1 + "/" + internetGatewayResource
	internetGatewayListURLV1 = urlProvidersPrefix + NetworkProviderV1 + "/" + internetGatewayListResource
	nicURLV1                 = urlProvidersPrefix + NetworkProviderV1 + "/" + nicResource
	nicListURLV1             = urlProvidersPrefix + NetworkProviderV1 + "/" + nicListResource
	publicIpURLV1            = urlProvidersPrefix + NetworkProviderV1 + "/" + publicIpResource
	publicIpListURLV1        = urlProvidersPrefix + NetworkProviderV1 + "/" + publicIpListResource
	routeTableURLV1          = urlProvidersPrefix + NetworkProviderV1 + "/" + routeTableResource
	routeTableListURLV1      = urlProvidersPrefix + NetworkProviderV1 + "/" + routeTableListResource
	subnetURLV1              = urlProvidersPrefix + NetworkProviderV1 + "/" + subnetResource
	subnetListURLV1          = urlProvidersPrefix + NetworkProviderV1 + "/" + subnetListResource
	securityGroupURLV1       = urlProvidersPrefix + NetworkProviderV1 + "/" + securityGroupResource
	securityGroupListURLV1   = urlProvidersPrefix + NetworkProviderV1 + "/" + securityGroupListResource
)
