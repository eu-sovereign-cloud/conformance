package constants

type SuiteName string

func (s SuiteName) String() string {
	return string(s)
}

const (
	// Suite Names

	AuthorizationProviderLifeCycleV1SuiteName SuiteName = "Authorization.V1.ProviderLifeCycle"
	AuthorizationProviderQueriesV1SuiteName   SuiteName = "Authorization.V1.ProviderQueries"
	RoleLifeCycleV1SuiteName                  SuiteName = "Authorization.V1.RoleLifeCycle"
	RoleAssignmentLifeCycleV1SuiteName        SuiteName = "Authorization.V1.RoleAssignmentLifeCycle"

	RegionProviderQueriesV1SuiteName SuiteName = "Region.V1.ProviderQueries"

	WorkspaceProviderLifeCycleV1SuiteName SuiteName = "Workspace.V1.ProviderLifeCycle"
	WorkspaceProviderQueriesV1SuiteName   SuiteName = "Workspace.V1.ProviderQueries"

	ComputeProviderLifeCycleV1SuiteName SuiteName = "Compute.V1.ProviderLifeCycle"
	ComputeProviderQueriesV1SuiteName   SuiteName = "Compute.V1.ProviderQueries"

	StorageProviderLifeCycleV1SuiteName SuiteName = "Storage.V1.ProviderLifeCycle"
	StorageProviderQueriesV1SuiteName   SuiteName = "Storage.V1.ProviderQueries"
	BlockStorageLifeCycleV1SuiteName    SuiteName = "Storage.V1.BlockStorageLifeCycle"
	ImageLifeCycleV1SuiteName           SuiteName = "Storage.V1.ImageLifeCycle"

	NetworkProviderLifeCycleV1SuiteName SuiteName = "Network.V1.ProviderLifeCycle"
	NetworkProviderQueriesV1SuiteName   SuiteName = "Network.V1.ProviderQueries"

	UsageFoundationProvidersV1SuiteName SuiteName = "Usage.V1.FoundationProviders"
)

var AllSuiteNames = []SuiteName{
	AuthorizationProviderLifeCycleV1SuiteName,
	AuthorizationProviderQueriesV1SuiteName,
	RoleLifeCycleV1SuiteName,
	RoleAssignmentLifeCycleV1SuiteName,

	RegionProviderQueriesV1SuiteName,

	WorkspaceProviderLifeCycleV1SuiteName,
	WorkspaceProviderQueriesV1SuiteName,

	ComputeProviderLifeCycleV1SuiteName,
	ComputeProviderQueriesV1SuiteName,

	StorageProviderLifeCycleV1SuiteName,
	StorageProviderQueriesV1SuiteName,
	BlockStorageLifeCycleV1SuiteName,
	ImageLifeCycleV1SuiteName,

	NetworkProviderLifeCycleV1SuiteName,
	NetworkProviderQueriesV1SuiteName,

	UsageFoundationProvidersV1SuiteName,
}
