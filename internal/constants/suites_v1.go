package constants

type SuiteName string

func (s SuiteName) String() string {
	return string(s)
}

const (
	// Suite Names

	AuthorizationLifeCycleV1SuiteName SuiteName = "Authorization.V1.LifeCycle"
	AuthorizationListV1SuiteName      SuiteName = "Authorization.V1.List"

	RegionListV1SuiteName SuiteName = "Region.V1.List"

	WorkspaceLifeCycleV1SuiteName SuiteName = "Workspace.V1.LifeCycle"
	WorkspaceListV1SuiteName      SuiteName = "Workspace.V1.List"
	CreateWorkspaceV1SuiteName    SuiteName = "Workspace.V1.CreateWorkspace"

	ComputeLifeCycleV1SuiteName SuiteName = "Compute.V1.LifeCycle"
	ComputeListV1SuiteName      SuiteName = "Compute.V1.List"

	StorageLifeCycleV1SuiteName   SuiteName = "Storage.V1.LifeCycle"
	StorageListV1SuiteName        SuiteName = "Storage.V1.List"
	CreateBlockStorageV1SuiteName SuiteName = "Storage.V1.CreateBlockStorage"
	UpdateBlockStorageV1SuiteName SuiteName = "Storage.V1.UpdateBlockStorage"

	NetworkLifeCycleV1SuiteName SuiteName = "Network.V1.LifeCycle"
	NetworkListV1SuiteName      SuiteName = "Network.V1.List"

	FoundationUsageV1SuiteName SuiteName = "Foundation.V1.Usage"
)

var AllSuiteNames = []SuiteName{
	AuthorizationLifeCycleV1SuiteName,
	AuthorizationListV1SuiteName,

	RegionListV1SuiteName,

	WorkspaceLifeCycleV1SuiteName,
	WorkspaceListV1SuiteName,
	CreateWorkspaceV1SuiteName,

	ComputeLifeCycleV1SuiteName,
	ComputeListV1SuiteName,

	StorageLifeCycleV1SuiteName,
	StorageListV1SuiteName,
	CreateBlockStorageV1SuiteName,
	UpdateBlockStorageV1SuiteName,

	NetworkLifeCycleV1SuiteName,
	NetworkListV1SuiteName,

	FoundationUsageV1SuiteName,
}
