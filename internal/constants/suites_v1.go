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

	NetworkProviderLifeCycleV1SuiteName   SuiteName = "Network.V1.ProviderLifeCycle"
	NetworkProviderQueriesV1SuiteName     SuiteName = "Network.V1.ProviderQueries"
	NetworkLifeCycleV1SuiteName           SuiteName = "Network.V1.NetworkLifeCycle"
	SubnetLifeCycleV1SuiteName            SuiteName = "Network.V1.SubnetLifeCycle"
	SecurityGroupLifeCycleV1SuiteName     SuiteName = "Network.V1.SecurityGroupLifeCycle"
	SecurityGroupRuleLifeCycleV1SuiteName SuiteName = "Network.V1.SecurityGroupRuleLifeCycle"
	InternetGatewayLifeCycleV1SuiteName   SuiteName = "Network.V1.InternetGatewayLifeCycle"
	PublicIpLifeCycleV1SuiteName          SuiteName = "Network.V1.PublicIpLifeCycle"
	NicLifeCycleV1SuiteName               SuiteName = "Network.V1.NicLifeCycle"
	RouteTableLifeCycleV1SuiteName        SuiteName = "Network.V1.RouteTableLifeCycle"

	UsageFoundationProvidersV1SuiteName SuiteName = "Usage.V1.FoundationProviders"

	// Constraints

	RoleConstraintsV1SuiteName              SuiteName = "Authorization.V1.RoleConstraints"
	RoleAssignmentConstraintsV1SuiteName    SuiteName = "Authorization.V1.RoleAssignmentConstraints"
	InstanceConstraintsV1SuiteName          SuiteName = "Compute.V1.InstanceConstraints"
	WorkspaceConstraintsV1SuiteName         SuiteName = "Workspace.V1.WorkspaceConstraints"
	BlockStorageConstraintsV1SuiteName      SuiteName = "BlockStorageConstraintsV1"
	ImageConstraintsV1SuiteName             SuiteName = "ImageConstraintsV1"
	NetworkConstraintsV1SuiteName           SuiteName = "NetworkConstraintsV1"
	InternetGatewayConstraintsV1SuiteName   SuiteName = "InternetGatewayConstraintsV1"
	PublicIpConstraintsV1SuiteName          SuiteName = "PublicIpConstraintsV1"
	NicConstraintsV1SuiteName               SuiteName = "NicConstraintsV1"
	SecurityGroupConstraintsV1SuiteName     SuiteName = "SecurityGroupConstraintsV1"
	SecurityGroupRuleConstraintsV1SuiteName SuiteName = "SecurityGroupRuleConstraintsV1"
	RouteTableConstraintsV1SuiteName        SuiteName = "RouteTableConstraintsV1"
	SubnetConstraintsV1SuiteName            SuiteName = "SubnetConstraintsV1"
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
	NetworkLifeCycleV1SuiteName,
	SubnetLifeCycleV1SuiteName,
	SecurityGroupRuleLifeCycleV1SuiteName,
	SecurityGroupLifeCycleV1SuiteName,
	InternetGatewayLifeCycleV1SuiteName,
	PublicIpLifeCycleV1SuiteName,
	NicLifeCycleV1SuiteName,
	RouteTableLifeCycleV1SuiteName,

	UsageFoundationProvidersV1SuiteName,

	RoleConstraintsV1SuiteName,
	RoleAssignmentConstraintsV1SuiteName,
	InstanceConstraintsV1SuiteName,
	WorkspaceConstraintsV1SuiteName,
	BlockStorageConstraintsV1SuiteName,
	ImageConstraintsV1SuiteName,
	NetworkConstraintsV1SuiteName,
	InternetGatewayConstraintsV1SuiteName,
	PublicIpConstraintsV1SuiteName,
	NicConstraintsV1SuiteName,
	SecurityGroupConstraintsV1SuiteName,
	SecurityGroupRuleConstraintsV1SuiteName,
	RouteTableConstraintsV1SuiteName,
	SubnetConstraintsV1SuiteName,
}
