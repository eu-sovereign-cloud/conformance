package constants

type OperationName string

const (
	// Authorization

	CreateOrUpdateRoleOperation           OperationName = "CreateOrUpdateRole"
	GetRoleOperation                      OperationName = "GetRole"
	CreateOrUpdateRoleAssignmentOperation OperationName = "CreateOrUpdate"
	GetRoleAssignmentOperation            OperationName = "Get"

	// Workspace
	CreateOrUpdateWorkspaceOperation OperationName = "CreateOrUpdateWorkspace"
	GetWorkspaceOperation            OperationName = "GetWorkspace"

	// Storage
	CreateOrUpdateBlockStorageOperation OperationName = "CreateOrUpdateBlockStorage"
	GetBlockStorageOperation            OperationName = "GetBlockStorage"
	CreateOrUpdateImageOperation        OperationName = "CreateOrUpdateImage"
	GetImageOperation                   OperationName = "GetImage"

	// Compute
	CreateOrUpdateInstanceOperation OperationName = "CreateOrUpdateInstance"
	GetInstanceOperation            OperationName = "GetInstance"

	// Network
	CreateOrUpdateNetworkOperation           OperationName = "CreateOrUpdateNetwork"
	GetNetworkOperation                      OperationName = "GetNetwork"
	CreateOrUpdateInternetGatewayOperation   OperationName = "CreateOrUpdateInternetGateway"
	GetInternetGatewayOperation              OperationName = "GetInternetGateway"
	CreateOrUpdateRouteTableOperation        OperationName = "CreateOrUpdateRouteTable"
	GetRouteTableOperation                   OperationName = "GetRouteTable"
	CreateOrUpdateSubnetOperation            OperationName = "CreateOrUpdateSubnet"
	GetSubnetOperation                       OperationName = "GetSubnet"
	CreateOrUpdatePublicIpOperation          OperationName = "CreateOrUpdatePublicIp"
	GetPublicIpOperation                     OperationName = "GetPublicIp"
	CreateOrUpdateNicOperation               OperationName = "CreateOrUpdateNic"
	GetNicOperation                          OperationName = "GetNic"
	CreateOrUpdateSecurityGroupRuleOperation OperationName = "CreateOrUpdateSecurityGroupRule"
	GetSecurityGroupRuleOperation            OperationName = "GetSecurityGroupRule"
	CreateOrUpdateSecurityGroupOperation     OperationName = "CreateOrUpdateSecurityGroup"
	GetSecurityGroupOperation                OperationName = "GetSecurityGroup"
)
