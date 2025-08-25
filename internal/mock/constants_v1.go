package mock

const (
	// Versions
	version1 = "v1"

	// Providers
	workspaceProviderV1     = "seca.workspace/" + version1
	computeProviderV1       = "seca.compute/" + version1
	storageProviderV1       = "seca.storage/" + version1
	authorizationProviderV1 = "seca.authorization/" + version1
	networkProviderV1       = "seca.network/" + version1

	// Base URLs
	workspaceURLV1       = "/providers/" + workspaceProviderV1 + "/" + workspaceResource
	instanceSkuURLV1     = "/providers/" + computeProviderV1 + "/" + instanceSkuResource
	instanceURLV1        = "/providers/" + computeProviderV1 + "/" + instanceResource
	storageSkuURLV1      = "/providers/" + storageProviderV1 + "/" + storageSkuResource
	blockStorageURLV1    = "/providers/" + storageProviderV1 + "/" + blockStorageResource
	imageV1URLV1         = "/providers/" + storageProviderV1 + "/" + imageResource
	rolesURLV1           = "/providers/" + authorizationProviderV1 + "/" + rolesResource
	roleAssignmentURLV1  = "/providers/" + authorizationProviderV1 + "/" + roleAssignmentResource
	networkURLV1         = "/providers/" + networkProviderV1 + "/" + networkResource
	networkSkuURLV1      = "/providers/" + networkProviderV1 + "/" + networkSkuResource
	internetGatewayURLV1 = "/providers/" + networkProviderV1 + "/" + internetGatewayResource
	nicURLV1             = "/providers/" + networkProviderV1 + "/" + nicResource
	publicIPURLV1        = "/providers/" + networkProviderV1 + "/" + publicIPResource
	routeTableURLV1      = "/providers/" + networkProviderV1 + "/" + routeTableResource
	subnetURLV1          = "/providers/" + networkProviderV1 + "/" + subnetResource
	securityGroupURLV1   = "/providers/" + networkProviderV1 + "/" + SecurityGroupResource
)
