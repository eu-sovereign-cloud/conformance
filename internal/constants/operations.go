package constants

type OperationName string

const (
	// Skus
	ListSkusOperation OperationName = "ListSkus"
	GetSkuOperation   OperationName = "GetSku"

	// Regions
	ListRegionsOperation OperationName = "ListRegions"
	GetRegionOperation   OperationName = "GetRegion"

	// Authorization
	ListRolesOperation          OperationName = "ListRoles"
	GetRoleOperation            OperationName = "GetRole"
	CreateOrUpdateRoleOperation OperationName = "CreateOrUpdateRole"
	DeleteRoleOperation         OperationName = "DeleteRole"

	ListRoleAssignmentsOperation          OperationName = "ListRoleAssignments"
	GetRoleAssignmentOperation            OperationName = "GetRoleAssignment"
	CreateOrUpdateRoleAssignmentOperation OperationName = "CreateOrUpdateRoleAssignment"
	DeleteRoleAssignmentOperation         OperationName = "DeleteRoleAssignment"

	// Workspace
	ListWorkspacesOperation          OperationName = "ListWorkspace"
	GetWorkspaceOperation            OperationName = "GetWorkspace"
	CreateOrUpdateWorkspaceOperation OperationName = "CreateOrUpdateWorkspace"
	DeleteWorkspaceOperation         OperationName = "DeleteWorkspace"

	// Storage
	ListBlockStorageOperation           OperationName = "ListBlockStorage"
	GetBlockStorageOperation            OperationName = "GetBlockStorage"
	CreateOrUpdateBlockStorageOperation OperationName = "CreateOrUpdateBlockStorage"
	DeleteBlockStorageOperation         OperationName = "DeleteBlockStorage"

	ListImagesOperation          OperationName = "ListImage"
	GetImageOperation            OperationName = "GetImage"
	CreateOrUpdateImageOperation OperationName = "CreateOrUpdateImage"
	DeleteImageOperation         OperationName = "DeleteImage"

	// Compute
	ListInstancesOperation          OperationName = "ListInstances"
	GetInstanceOperation            OperationName = "GetInstance"
	CreateOrUpdateInstanceOperation OperationName = "CreateOrUpdateInstance"
	StartInstanceOperation          OperationName = "StartInstance"
	StopInstanceOperation           OperationName = "StopInstance"
	RestartInstanceOperation        OperationName = "RestartInstance"
	DeleteInstanceOperation         OperationName = "DeleteInstance"

	// Network
	ListNetworksOperation          OperationName = "ListNetworks"
	GetNetworkOperation            OperationName = "GetNetwork"
	CreateOrUpdateNetworkOperation OperationName = "CreateOrUpdateNetwork"
	DeleteNetworkOperation         OperationName = "DeleteNetwork"

	ListInternetGatewaysOperation          OperationName = "ListInternetGateways"
	GetInternetGatewayOperation            OperationName = "GetInternetGateway"
	CreateOrUpdateInternetGatewayOperation OperationName = "CreateOrUpdateInternetGateway"
	DeleteInternetGatewayOperation         OperationName = "DeleteInternetGateway"

	ListRouteTablesOperation          OperationName = "ListRouteTables"
	GetRouteTableOperation            OperationName = "GetRouteTable"
	CreateOrUpdateRouteTableOperation OperationName = "CreateOrUpdateRouteTable"
	DeleteRouteTableOperation         OperationName = "DeleteRouteTable"

	ListSubnetsOperation          OperationName = "ListSubnets"
	GetSubnetOperation            OperationName = "GetSubnet"
	CreateOrUpdateSubnetOperation OperationName = "CreateOrUpdateSubnet"
	DeleteSubnetOperation         OperationName = "DeleteSubnet"

	ListPublicIpsOperation          OperationName = "ListPublicIps"
	GetPublicIpOperation            OperationName = "GetPublicIp"
	CreateOrUpdatePublicIpOperation OperationName = "CreateOrUpdatePublicIp"
	DeletePublicIpOperation         OperationName = "DeletePublicIp"

	ListNicsOperation          OperationName = "ListNics"
	GetNicOperation            OperationName = "GetNic"
	CreateOrUpdateNicOperation OperationName = "CreateOrUpdateNic"
	DeleteNicOperation         OperationName = "DeleteNic"

	ListSecurityGroupRulesOperation          OperationName = "ListSecurityGroupRules"
	GetSecurityGroupRuleOperation            OperationName = "GetSecurityGroupRule"
	CreateOrUpdateSecurityGroupRuleOperation OperationName = "CreateOrUpdateSecurityGroupRule"
	DeleteSecurityGroupRuleOperation         OperationName = "DeleteSecurityGroupRule"

	ListSecurityGroupsOperation          OperationName = "ListSecurityGroups"
	GetSecurityGroupOperation            OperationName = "GetSecurityGroup"
	CreateOrUpdateSecurityGroupOperation OperationName = "CreateOrUpdateSecurityGroup"
	DeleteSecurityGroupOperation         OperationName = "DeleteSecurityGroup"
)
