package params

import (
	authorization "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.authorization.v1"
	compute "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.compute.v1"
	network "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.network.v1"
	storage "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.storage.v1"
	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Authorization

type AuthorizationProviderLifeCycleV1Params struct {
	RoleInitial           *schema.Role
	RoleUpdated           *schema.Role
	RoleAssignmentInitial *schema.RoleAssignment
	RoleAssignmentUpdated *schema.RoleAssignment
}

type AuthorizationProviderQueriesV1Params struct {
	Roles           authorization.RoleIterator
	RoleAssignments authorization.RoleAssignmentIterator
}

type RoleLifeCycleV1Params struct {
	RoleInitial *schema.Role
	RoleUpdated *schema.Role
}

type RoleAssignmentLifeCycleV1Params struct {
	RoleAssignmentInitial *schema.RoleAssignment
	RoleAssignmentUpdated *schema.RoleAssignment
}

// Region

type RegionProviderQueriesV1Params struct {
	RegionsMetadata []schema.GlobalResourceMetadata

	MockProviders []string
}

// Workspace

type WorkspaceProviderLifeCycleV1Params struct {
	WorkspaceInitial *schema.Workspace
	WorkspaceUpdated *schema.Workspace
}

type WorkspaceProviderQueriesV1Params struct {
	Workspaces workspace.WorkspaceIterator
}

// Compute

type ComputeProviderLifeCycleV1Params struct {
	Workspace       *schema.Workspace
	BlockStorage    *schema.BlockStorage
	InitialInstance *schema.Instance
	UpdatedInstance *schema.Instance
}

type ComputeProviderQueriesV1Params struct {
	Workspace    *schema.Workspace
	BlockStorage *schema.BlockStorage
	Instances    compute.InstanceIterator
}

// Storage

type StorageProviderLifeCycleV1Params struct {
	Workspace           *schema.Workspace
	BlockStorageInitial *schema.BlockStorage
	BlockStorageUpdated *schema.BlockStorage
	ImageInitial        *schema.Image
	ImageUpdated        *schema.Image
}

type StorageProviderQueriesV1Params struct {
	Workspace     *schema.Workspace
	BlockStorages storage.BlockStorageIterator
	Images        storage.ImageIterator
}

type BlockStorageLifeCycleV1Params struct {
	Workspace           *schema.Workspace
	BlockStorageInitial *schema.BlockStorage
	BlockStorageUpdated *schema.BlockStorage
}

type ImageLifeCycleV1Params struct {
	Workspace    *schema.Workspace
	BlockStorage *schema.BlockStorage
	ImageInitial *schema.Image
	ImageUpdated *schema.Image
}

// Constraints Validation

type RoleConstraintsValidationV1Params struct {
	OverLengthNameRole               *schema.Role
	InvalidPatternNameRole           *schema.Role
	OverLengthLabelValueRole         *schema.Role
	OverLengthAnnotationRole         *schema.Role
	OverLengthPermissionProviderRole *schema.Role
	OverLengthPermissionResourceRole *schema.Role
	OverLengthPermissionVerbRole     *schema.Role

	EmptyPermissionsRole                *schema.Role
	OverMaxItemsPermissionsRole         *schema.Role
	EmptyPermissionProviderRole         *schema.Role
	EmptyPermissionResourcesRole        *schema.Role
	OverMaxItemsPermissionResourcesRole *schema.Role
	EmptyPermissionResourceValueRole    *schema.Role
	EmptyPermissionVerbsRole            *schema.Role
	OverMaxItemsPermissionVerbsRole     *schema.Role
	EmptyPermissionVerbValueRole        *schema.Role
}

type RoleAssignmentConstraintsValidationV1Params struct {
	OverLengthNameRoleAssignment           *schema.RoleAssignment
	InvalidPatternNameRoleAssignment       *schema.RoleAssignment
	OverLengthLabelValueRoleAssignment     *schema.RoleAssignment
	OverLengthAnnotationRoleAssignment     *schema.RoleAssignment
	OverLengthSubRoleAssignment            *schema.RoleAssignment
	OverLengthRoleNameRoleAssignment       *schema.RoleAssignment
	OverLengthScopeTenantRoleAssignment    *schema.RoleAssignment
	OverLengthScopeRegionRoleAssignment    *schema.RoleAssignment
	OverLengthScopeWorkspaceRoleAssignment *schema.RoleAssignment

	EmptyRolesRoleAssignment                  *schema.RoleAssignment
	OverMaxItemsRolesRoleAssignment           *schema.RoleAssignment
	EmptyRoleValueRoleAssignment              *schema.RoleAssignment
	EmptySubsRoleAssignment                   *schema.RoleAssignment
	OverMaxItemsSubsRoleAssignment            *schema.RoleAssignment
	EmptySubValueRoleAssignment               *schema.RoleAssignment
	EmptyScopesRoleAssignment                 *schema.RoleAssignment
	OverMaxItemsScopesRoleAssignment          *schema.RoleAssignment
	EmptyScopeTenantValueRoleAssignment       *schema.RoleAssignment
	OverMaxItemsScopeTenantsRoleAssignment    *schema.RoleAssignment
	EmptyScopeRegionValueRoleAssignment       *schema.RoleAssignment
	OverMaxItemsScopeRegionsRoleAssignment    *schema.RoleAssignment
	EmptyScopeWorkspaceValueRoleAssignment    *schema.RoleAssignment
	OverMaxItemsScopeWorkspacesRoleAssignment *schema.RoleAssignment
}

type InstanceConstraintsValidationV1Params struct {
	Workspace                           *schema.Workspace
	BlockStorage                        *schema.BlockStorage
	OverLengthNameInstance              *schema.Instance
	InvalidPatternNameInstance          *schema.Instance
	OverLengthLabelValueInstance        *schema.Instance
	OverLengthAnnotationInstance        *schema.Instance
	OverLengthUserDataInstance          *schema.Instance
	OverLengthAntiAffinityGroupInstance *schema.Instance
	OverLengthSshKeyInstance            *schema.Instance

	OverMaxItemsSshKeysInstance     *schema.Instance
	EmptySshKeyValueInstance        *schema.Instance
	OverLengthZoneInstance          *schema.Instance
	EmptyZoneInstance               *schema.Instance
	OverMaxItemsDataVolumesInstance *schema.Instance
}

type WorkspaceConstraintsValidationV1Params struct {
	OverLengthNameWorkspace       *schema.Workspace
	InvalidPatternNameWorkspace   *schema.Workspace
	OverLengthLabelValueWorkspace *schema.Workspace
	OverLengthAnnotationWorkspace *schema.Workspace
}

type BlockStorageConstraintsValidationV1Params struct {
	Workspace                        *schema.Workspace
	OverLengthNameBlockStorage       *schema.BlockStorage
	InvalidPatternNameBlockStorage   *schema.BlockStorage
	OverLengthLabelValueBlockStorage *schema.BlockStorage
	OverLengthAnnotationBlockStorage *schema.BlockStorage
	OverMaxSizeBlockStorage          *schema.BlockStorage
	ZeroSizeBlockStorage             *schema.BlockStorage
}

type ImageConstraintsValidationV1Params struct {
	Workspace                   *schema.Workspace
	BlockStorage                *schema.BlockStorage
	OverLengthNameImage         *schema.Image
	InvalidPatternNameImage     *schema.Image
	OverLengthLabelValueImage   *schema.Image
	OverLengthAnnotationImage   *schema.Image
	InvalidCpuArchitectureImage *schema.Image
	InvalidInitializerImage     *schema.Image
	InvalidBootImage            *schema.Image
}

// Network

type NetworkProviderLifeCycleV1Params struct {
	Workspace                *schema.Workspace
	BlockStorage             *schema.BlockStorage
	Instance                 *schema.Instance
	NetworkInitial           *schema.Network
	NetworkUpdated           *schema.Network
	InternetGatewayInitial   *schema.InternetGateway
	InternetGatewayUpdated   *schema.InternetGateway
	RouteTableInitial        *schema.RouteTable
	RouteTableUpdated        *schema.RouteTable
	SubnetInitial            *schema.Subnet
	SubnetUpdated            *schema.Subnet
	NicInitial               *schema.Nic
	NicUpdated               *schema.Nic
	PublicIpInitial          *schema.PublicIp
	PublicIpUpdated          *schema.PublicIp
	SecurityGroupInitial     *schema.SecurityGroup
	SecurityGroupUpdated     *schema.SecurityGroup
	SecurityGroupRuleInitial *schema.SecurityGroupRule
	SecurityGroupRuleUpdated *schema.SecurityGroupRule
}

type NetworkProviderQueriesV1Params struct {
	Workspace          *schema.Workspace
	BlockStorage       *schema.BlockStorage
	Instance           *schema.Instance
	Networks           network.NetworkIterator
	InternetGateways   network.InternetGatewayIterator
	RouteTables        network.RouteTableIterator
	Subnets            network.SubnetIterator
	Nics               network.NicIterator
	PublicIps          network.PublicIpIterator
	SecurityGroupRules network.SecurityGroupRuleIterator
	SecurityGroups     network.SecurityGroupIterator
}

type NetworkLifeCycleV1Params struct {
	Workspace       *schema.Workspace
	NetworkInitial  *schema.Network
	NetworkUpdated  *schema.Network
	RouteTable      *schema.RouteTable
	InternetGateway *schema.InternetGateway
}

type SubnetLifeCycleV1Params struct {
	Workspace       *schema.Workspace
	Network         *schema.Network
	RouteTable      *schema.RouteTable
	InternetGateway *schema.InternetGateway
	SubnetInitial   *schema.Subnet
	SubnetUpdated   *schema.Subnet
}

type InternetGatewayLifeCycleV1Params struct {
	Workspace              *schema.Workspace
	InternetGatewayInitial *schema.InternetGateway
	InternetGatewayUpdated *schema.InternetGateway
}

type NicLifeCycleV1Params struct {
	Workspace       *schema.Workspace
	Network         *schema.Network
	InternetGateway *schema.InternetGateway
	RouteTable      *schema.RouteTable
	Subnet          *schema.Subnet
	NicInitial      *schema.Nic
	NicUpdated      *schema.Nic
}

type RouteTableLifeCycleV1Params struct {
	Workspace         *schema.Workspace
	Network           *schema.Network
	InternetGateway   *schema.InternetGateway
	RouteTableInitial *schema.RouteTable
	RouteTableUpdated *schema.RouteTable
}

type PublicIpLifeCycleV1Params struct {
	Workspace       *schema.Workspace
	PublicIpInitial *schema.PublicIp
	PublicIpUpdated *schema.PublicIp
}

type SecurityGroupLifeCycleV1Params struct {
	Workspace            *schema.Workspace
	SecurityGroupInitial *schema.SecurityGroup
	SecurityGroupUpdated *schema.SecurityGroup
}

type SecurityGroupRuleLifeCycleV1Params struct {
	Workspace                *schema.Workspace
	SecurityGroupRuleInitial *schema.SecurityGroupRule
	SecurityGroupRuleUpdated *schema.SecurityGroupRule
}

// Usage

type FoundationUsageV1Params struct {
	Role            *schema.Role
	RoleAssignment  *schema.RoleAssignment
	Workspace       *schema.Workspace
	Image           *schema.Image
	BlockStorage    *schema.BlockStorage
	Network         *schema.Network
	InternetGateway *schema.InternetGateway
	RouteTable      *schema.RouteTable
	Subnet          *schema.Subnet
	SecurityGroup   *schema.SecurityGroup
	PublicIp        *schema.PublicIp
	Nic             *schema.Nic
	Instance        *schema.Instance
}

// Network Constraints

type NetworkConstraintsValidationV1Params struct {
	Workspace                           *schema.Workspace
	OverLengthNameNetwork               *schema.Network
	InvalidPatternNameNetwork           *schema.Network
	OverLengthLabelValueNetwork         *schema.Network
	OverLengthAnnotationNetwork         *schema.Network
	OverLengthCidrIpv4Network           *schema.Network
	UnderLengthCidrIpv4Network          *schema.Network
	OverLengthCidrIpv6Network           *schema.Network
	UnderLengthCidrIpv6Network          *schema.Network
	OverLengthAdditionalCidrIpv4Network *schema.Network
	OverLengthAdditionalCidrIpv6Network *schema.Network

	UnderLengthAdditionalCidrIpv4Network *schema.Network
	UnderLengthAdditionalCidrIpv6Network *schema.Network
	OverMaxItemsAdditionalCidrsNetwork   *schema.Network
}

type InternetGatewayConstraintsValidationV1Params struct {
	Workspace                           *schema.Workspace
	OverLengthNameInternetGateway       *schema.InternetGateway
	InvalidPatternNameInternetGateway   *schema.InternetGateway
	OverLengthLabelValueInternetGateway *schema.InternetGateway
	OverLengthAnnotationInternetGateway *schema.InternetGateway
}

type PublicIpConstraintsValidationV1Params struct {
	Workspace                    *schema.Workspace
	OverLengthNamePublicIp       *schema.PublicIp
	InvalidPatternNamePublicIp   *schema.PublicIp
	OverLengthLabelValuePublicIp *schema.PublicIp
	OverLengthAnnotationPublicIp *schema.PublicIp
	OverLengthAddressPublicIp    *schema.PublicIp
}

type NicConstraintsValidationV1Params struct {
	Workspace               *schema.Workspace
	OverLengthNameNic       *schema.Nic
	InvalidPatternNameNic   *schema.Nic
	OverLengthLabelValueNic *schema.Nic
	OverLengthAnnotationNic *schema.Nic
	OverLengthAddressNic    *schema.Nic
}

type SecurityGroupConstraintsValidationV1Params struct {
	Workspace                         *schema.Workspace
	OverLengthNameSecurityGroup       *schema.SecurityGroup
	InvalidPatternNameSecurityGroup   *schema.SecurityGroup
	OverLengthLabelValueSecurityGroup *schema.SecurityGroup
	OverLengthAnnotationSecurityGroup *schema.SecurityGroup
	InvalidDirectionSecurityGroup     *schema.SecurityGroup
	InvalidVersionSecurityGroup       *schema.SecurityGroup
	InvalidProtocolSecurityGroup      *schema.SecurityGroup
	OverMaxPortFromSecurityGroup      *schema.SecurityGroup
	UnderMinPortFromSecurityGroup     *schema.SecurityGroup
	OverMaxPortToSecurityGroup        *schema.SecurityGroup
	UnderMinPortToSecurityGroup       *schema.SecurityGroup
	OverMaxPortListSecurityGroup      *schema.SecurityGroup
	UnderMinPortListSecurityGroup     *schema.SecurityGroup
	OverMaxIcmpTypeSecurityGroup      *schema.SecurityGroup
	OverMaxIcmpCodeSecurityGroup      *schema.SecurityGroup
}

type SecurityGroupRuleConstraintsValidationV1Params struct {
	Workspace                             *schema.Workspace
	OverLengthNameSecurityGroupRule       *schema.SecurityGroupRule
	InvalidPatternNameSecurityGroupRule   *schema.SecurityGroupRule
	OverLengthLabelValueSecurityGroupRule *schema.SecurityGroupRule
	OverLengthAnnotationSecurityGroupRule *schema.SecurityGroupRule
	InvalidDirectionSecurityGroupRule     *schema.SecurityGroupRule
	InvalidVersionSecurityGroupRule       *schema.SecurityGroupRule
	InvalidProtocolSecurityGroupRule      *schema.SecurityGroupRule
	OverMaxPortFromSecurityGroupRule      *schema.SecurityGroupRule
	UnderMinPortFromSecurityGroupRule     *schema.SecurityGroupRule
	OverMaxPortToSecurityGroupRule        *schema.SecurityGroupRule
	UnderMinPortToSecurityGroupRule       *schema.SecurityGroupRule
	OverMaxPortListSecurityGroupRule      *schema.SecurityGroupRule
	UnderMinPortListSecurityGroupRule     *schema.SecurityGroupRule
	OverMaxIcmpTypeSecurityGroupRule      *schema.SecurityGroupRule
	OverMaxIcmpCodeSecurityGroupRule      *schema.SecurityGroupRule
}

type RouteTableConstraintsValidationV1Params struct {
	Workspace                                *schema.Workspace
	Network                                  *schema.Network
	InternetGateway                          *schema.InternetGateway
	OverLengthNameRouteTable                 *schema.RouteTable
	InvalidPatternNameRouteTable             *schema.RouteTable
	OverLengthLabelValueRouteTable           *schema.RouteTable
	OverLengthAnnotationRouteTable           *schema.RouteTable
	OverLengthDestinationCidrBlockRouteTable *schema.RouteTable

	EmptyRoutesRouteTable                 *schema.RouteTable
	OverMaxItemsRoutesRouteTable          *schema.RouteTable
	EmptyDestinationCidrBlockRouteTable   *schema.RouteTable
	InvalidDestinationCidrBlockRouteTable *schema.RouteTable
}

type SubnetConstraintsValidationV1Params struct {
	Workspace                  *schema.Workspace
	Network                    *schema.Network
	InternetGateway            *schema.InternetGateway
	RouteTable                 *schema.RouteTable
	OverLengthNameSubnet       *schema.Subnet
	InvalidPatternNameSubnet   *schema.Subnet
	OverLengthLabelValueSubnet *schema.Subnet
	OverLengthAnnotationSubnet *schema.Subnet

	OverLengthZoneSubnet *schema.Subnet
}

// errors
type BlockStorageErrorV1Params struct {
	Workspace                        *schema.Workspace
	InvalidRegionBlockStorage        *schema.BlockStorage
	InvalidSkuBlockStorage           *schema.BlockStorage
	NonExistentWorkspaceBlockStorage *schema.BlockStorage
}

type ImageErrorV1Params struct {
	Workspace                    *schema.Workspace
	BlockStorage                 *schema.BlockStorage
	InvalidRegionImage           *schema.Image
	InvalidCpuArchitectureImage  *schema.Image
	CrossRegionBlockStorageImage *schema.Image
	NonExistentWorkspaceImage    *schema.Image
}

type InternetGatewayErrorV1Params struct {
	Workspace                           *schema.Workspace
	InvalidRegionInternetGateway        *schema.InternetGateway
	NonExistentWorkspaceInternetGateway *schema.InternetGateway
}

type SecurityGroupRuleErrorV1Params struct {
	Workspace                             *schema.Workspace
	InvalidRegionSecurityGroupRule        *schema.SecurityGroupRule
	NonExistentWorkspaceSecurityGroupRule *schema.SecurityGroupRule
}

type SecurityGroupErrorV1Params struct {
	Workspace                         *schema.Workspace
	InvalidRegionSecurityGroup        *schema.SecurityGroup
	NonExistentWorkspaceSecurityGroup *schema.SecurityGroup
	NonExistentRuleRefSecurityGroup   *schema.SecurityGroup
}

type PublicIpErrorV1Params struct {
	Workspace                    *schema.Workspace
	InvalidRegionPublicIp        *schema.PublicIp
	NonExistentWorkspacePublicIp *schema.PublicIp
	InvalidVersionPublicIp       *schema.PublicIp
}

type NetworkErrorV1Params struct {
	Workspace                   *schema.Workspace
	InvalidRegionNetwork        *schema.Network
	InvalidSkuNetwork           *schema.Network
	NonExistentWorkspaceNetwork *schema.Network
}

type RouteTableErrorV1Params struct {
	Workspace                      *schema.Workspace
	Network                        *schema.Network
	InternetGateway                *schema.InternetGateway
	InvalidRegionRouteTable        *schema.RouteTable
	NonExistentWorkspaceRouteTable *schema.RouteTable
	NonExistentNetworkRouteTable   *schema.RouteTable
	NonExistentTargetRefRouteTable *schema.RouteTable
}

type SubnetErrorV1Params struct {
	Workspace                      *schema.Workspace
	Network                        *schema.Network
	InternetGateway                *schema.InternetGateway
	RouteTable                     *schema.RouteTable
	InvalidRegionSubnet            *schema.Subnet
	InvalidZoneSubnet              *schema.Subnet
	NonExistentWorkspaceSubnet     *schema.Subnet
	NonExistentNetworkSubnet       *schema.Subnet
	NonExistentRouteTableRefSubnet *schema.Subnet
	OutsideCidrSubnet              *schema.Subnet
}

type NicErrorV1Params struct {
	Workspace                 *schema.Workspace
	Network                   *schema.Network
	InternetGateway           *schema.InternetGateway
	RouteTable                *schema.RouteTable
	Subnet                    *schema.Subnet
	InvalidRegionNic          *schema.Nic
	NonExistentWorkspaceNic   *schema.Nic
	NonExistentSubnetRefNic   *schema.Nic
	NonExistentPublicIpRefNic *schema.Nic
}
