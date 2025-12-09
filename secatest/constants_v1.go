package secatest

const (
	// Suite Names
	authorizationV1LifeCycleSuiteName = "Authorization.V1.LifeCycle"
	regionV1LifeCycleSuiteName        = "Region.V1.LifeCycle"
	computeV1LifeCycleSuiteName       = "Compute.V1.LifeCycle"
	networkV1LifeCycleSuiteName       = "Network.V1.LifeCycle"
	storageV1LifeCycleSuiteName       = "Storage.V1.LifeCycle"
	workspaceV1LifeCycleSuiteName     = "Workspace.V1.LifeCycle"
	foundationV1UsageSuiteName        = "Foundation.V1.Usage"

	// Test Data

	apiVersion1 = "v1"

	// Providers
	authorizationProviderV1 = "seca.authorization/" + apiVersion1
	regionProviderV1        = "seca.region/" + apiVersion1
	workspaceProviderV1     = "seca.workspace/" + apiVersion1
	computeProviderV1       = "seca.compute/" + apiVersion1
	storageProviderV1       = "seca.storage/" + apiVersion1
	networkProviderV1       = "seca.network/" + apiVersion1
)
