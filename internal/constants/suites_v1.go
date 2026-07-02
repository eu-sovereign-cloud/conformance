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

	BlockStorageErrorV1SuiteName      SuiteName = "Storage.V1.BlockStorageError"
	ImageErrorV1SuiteName             SuiteName = "Storage.V1.ImageError"
	InternetGatewayErrorV1SuiteName   SuiteName = "Network.V1.InternetGatewayError"
	SecurityGroupRuleErrorV1SuiteName SuiteName = "Network.V1.SecurityGroupRuleError"
	SecurityGroupErrorV1SuiteName     SuiteName = "Network.V1.SecurityGroupError"
	PublicIpErrorV1SuiteName          SuiteName = "Network.V1.PublicIpError"
	NetworkErrorV1SuiteName           SuiteName = "Network.V1.NetworkError"
	RouteTableErrorV1SuiteName        SuiteName = "Network.V1.RouteTableError"
	SubnetErrorV1SuiteName            SuiteName = "Network.V1.SubnetError"
	NicErrorV1SuiteName               SuiteName = "Network.V1.NicError"

	// Constraints

	RoleConstraintsValidationV1SuiteName              SuiteName = "Authorization.V1.RoleConstraintsValidation"
	RoleAssignmentConstraintsValidationV1SuiteName    SuiteName = "Authorization.V1.RoleAssignmentConstraintsValidation"
	InstanceConstraintsValidationV1SuiteName          SuiteName = "Compute.V1.InstanceConstraintsValidation"
	WorkspaceConstraintsValidationV1SuiteName         SuiteName = "Workspace.V1.WorkspaceConstraintsValidation"
	BlockStorageConstraintsValidationV1SuiteName      SuiteName = "Storage.V1.BlockStorageConstraintsValidation"
	ImageConstraintsValidationV1SuiteName             SuiteName = "Storage.V1.ImageConstraintsValidation"
	NetworkConstraintsValidationV1SuiteName           SuiteName = "Network.V1.NetworkConstraintsValidation"
	InternetGatewayConstraintsValidationV1SuiteName   SuiteName = "Network.V1.InternetGatewayConstraintsValidation"
	PublicIpConstraintsValidationV1SuiteName          SuiteName = "Network.V1.PublicIpConstraintsValidation"
	NicConstraintsValidationV1SuiteName               SuiteName = "Network.V1.NicConstraintsValidation"
	SecurityGroupConstraintsValidationV1SuiteName     SuiteName = "Network.V1.SecurityGroupConstraintsValidation"
	SecurityGroupRuleConstraintsValidationV1SuiteName SuiteName = "Network.V1.SecurityGroupRuleConstraintsValidation"
	RouteTableConstraintsValidationV1SuiteName        SuiteName = "Network.V1.RouteTableConstraintsValidation"
	SubnetConstraintsValidationV1SuiteName            SuiteName = "Network.V1.SubnetConstraintsValidation"
)

var AllSuiteNames = []SuiteName{
	AuthorizationProviderLifeCycleV1SuiteName,
	AuthorizationProviderQueriesV1SuiteName,
	RoleLifeCycleV1SuiteName,
	RoleAssignmentLifeCycleV1SuiteName,
	RoleConstraintsValidationV1SuiteName,
	RoleAssignmentConstraintsValidationV1SuiteName,

	RegionProviderQueriesV1SuiteName,

	WorkspaceProviderLifeCycleV1SuiteName,
	WorkspaceProviderQueriesV1SuiteName,
	WorkspaceConstraintsValidationV1SuiteName,

	ComputeProviderLifeCycleV1SuiteName,
	ComputeProviderQueriesV1SuiteName,
	InstanceConstraintsValidationV1SuiteName,

	StorageProviderLifeCycleV1SuiteName,
	StorageProviderQueriesV1SuiteName,
	BlockStorageLifeCycleV1SuiteName,
	ImageLifeCycleV1SuiteName,
	BlockStorageConstraintsValidationV1SuiteName,
	ImageConstraintsValidationV1SuiteName,

	NetworkProviderLifeCycleV1SuiteName,
	NetworkProviderQueriesV1SuiteName,
	NetworkLifeCycleV1SuiteName,
	NetworkConstraintsValidationV1SuiteName,
	NetworkErrorV1SuiteName,
	SubnetLifeCycleV1SuiteName,
	SubnetConstraintsValidationV1SuiteName,
	SubnetErrorV1SuiteName,
	SecurityGroupRuleLifeCycleV1SuiteName,
	SecurityGroupRuleConstraintsValidationV1SuiteName,
	SecurityGroupRuleErrorV1SuiteName,
	SecurityGroupLifeCycleV1SuiteName,
	SecurityGroupConstraintsValidationV1SuiteName,
	SecurityGroupErrorV1SuiteName,
	InternetGatewayLifeCycleV1SuiteName,
	InternetGatewayConstraintsValidationV1SuiteName,
	InternetGatewayErrorV1SuiteName,
	PublicIpLifeCycleV1SuiteName,
	PublicIpConstraintsValidationV1SuiteName,
	PublicIpErrorV1SuiteName,
	NicLifeCycleV1SuiteName,
	NicConstraintsValidationV1SuiteName,
	NicErrorV1SuiteName,
	RouteTableLifeCycleV1SuiteName,
	RouteTableConstraintsValidationV1SuiteName,
	RouteTableErrorV1SuiteName,

	UsageFoundationProvidersV1SuiteName,
}
