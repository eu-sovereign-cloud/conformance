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
	RegionsURLV1         = urlProvidersPrefix + RegionProviderV1 + "/regions"

	roleURLV1            = urlProvidersPrefix + AuthorizationProviderV1 + "/" + roleResource
	roleAssignmentURLV1  = urlProvidersPrefix + AuthorizationProviderV1 + "/" + roleAssignmentResource	
	regionURLV1          = urlProvidersPrefix + RegionProviderV1 + "/" + regionResource
	workspaceURLV1       = urlProvidersPrefix + WorkspaceProviderV1 + "/" + workspaceResource
	instanceSkuURLV1     = urlProvidersPrefix + ComputeProviderV1 + "/" + skuResource
	instanceURLV1        = urlProvidersPrefix + ComputeProviderV1 + "/" + instanceResource
	storageSkuURLV1      = urlProvidersPrefix + StorageProviderV1 + "/" + skuResource
	blockStorageURLV1    = urlProvidersPrefix + StorageProviderV1 + "/" + blockStorageResource
	imageURLV1           = urlProvidersPrefix + StorageProviderV1 + "/" + imageResource
	networkSkuURLV1      = urlProvidersPrefix + NetworkProviderV1 + "/" + skuResource
	networkURLV1         = urlProvidersPrefix + NetworkProviderV1 + "/" + networkResource
	internetGatewayURLV1 = urlProvidersPrefix + NetworkProviderV1 + "/" + internetGatewayResource
	nicURLV1             = urlProvidersPrefix + NetworkProviderV1 + "/" + nicResource
	publicIpURLV1        = urlProvidersPrefix + NetworkProviderV1 + "/" + publicIpResource
	routeTableURLV1      = urlProvidersPrefix + NetworkProviderV1 + "/" + routeTableResource
	subnetURLV1          = urlProvidersPrefix + NetworkProviderV1 + "/" + subnetResource
	securityGroupURLV1   = urlProvidersPrefix + NetworkProviderV1 + "/" + securityGroupResource
)
