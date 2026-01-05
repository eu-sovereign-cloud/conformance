package conformance

const (
	// Suite Names
	AuthorizationV1LifeCycleSuiteName = "Authorization.V1.LifeCycle"
	AuthorizationV1ListSuiteName      = "Authorization.V1.List"
	RegionV1LifeCycleSuiteName        = "Region.V1.LifeCycle"
	ComputeV1LifeCycleSuiteName       = "Compute.V1.LifeCycle"
	ComputeV1ListSuiteName            = "Compute.V1.List"
	NetworkV1LifeCycleSuiteName       = "Network.V1.LifeCycle"
	NetworkV1ListSuiteName            = "Network.V1.List"
	StorageV1LifeCycleSuiteName       = "Storage.V1.LifeCycle"
	StorageV1ListSuiteName            = "Storage.V1.List"
	WorkspaceV1LifeCycleSuiteName     = "Workspace.V1.LifeCycle"
	WorkspaceV1ListSuiteName          = "Workspace.V1.List"
	FoundationV1UsageSuiteName        = "Foundation.V1.Usage"

	// Versions
	ApiVersion1 = "v1"

	// Test Data

	/// Zones
	ZoneA = "a"
	ZoneB = "b"

	/// Providers
	AuthorizationProviderV1 = "seca.authorization/" + ApiVersion1
	RegionProviderV1        = "seca.region/" + ApiVersion1
	WorkspaceProviderV1     = "seca.workspace/" + ApiVersion1
	ComputeProviderV1       = "seca.compute/" + ApiVersion1
	StorageProviderV1       = "seca.storage/" + ApiVersion1
	NetworkProviderV1       = "seca.network/" + ApiVersion1
)
