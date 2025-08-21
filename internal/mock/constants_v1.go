package mock

const (
	// Versions
	version1 = "v1"

	// Providers
	workspaceProviderV1 = "seca.workspace/" + version1
	computeProviderV1   = "seca.compute/" + version1
	storageProviderV1   = "seca.storage/" + version1

	// Base URLs
	workspaceURLV1    = "/providers/" + workspaceProviderV1 + "/" + workspaceResource
	instanceSkuURLV1  = "/providers/" + computeProviderV1 + "/" + skuResource
	instanceURLV1     = "/providers/" + computeProviderV1 + "/" + instanceResource
	storageSkuURLV1   = "/providers/" + storageProviderV1 + "/" + skuResource
	blockStorageURLV1 = "/providers/" + storageProviderV1 + "/" + blockStorageResource
	imageURLV1        = "/providers/" + storageProviderV1 + "/" + imageResource
)
