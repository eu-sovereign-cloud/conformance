package params

import (
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
	Roles           []schema.Role
	RoleAssignments []schema.RoleAssignment
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
	Regions []schema.Region
}

// Workspace

type WorkspaceProviderLifeCycleV1Params struct {
	WorkspaceInitial *schema.Workspace
	WorkspaceUpdated *schema.Workspace
}

type WorkspaceProviderQueriesV1Params struct {
	Workspaces []schema.Workspace
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
	Instances    []schema.Instance
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
	BlockStorages []schema.BlockStorage
	Images        []schema.Image
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
	OverLengthNameRole       *schema.Role
	InvalidPatternNameRole   *schema.Role
	OverLengthLabelValueRole *schema.Role
	OverLengthAnnotationRole *schema.Role
}
type RoleAssignmentConstraintsValidationV1Params struct {
	OverLengthNameRoleAssignment       *schema.RoleAssignment
	InvalidPatternNameRoleAssignment   *schema.RoleAssignment
	OverLengthLabelValueRoleAssignment *schema.RoleAssignment
	OverLengthAnnotationRoleAssignment *schema.RoleAssignment
}

type InstanceConstraintsValidationV1Params struct {
	Workspace                    *schema.Workspace
	BlockStorage                 *schema.BlockStorage
	OverLengthNameInstance       *schema.Instance
	InvalidPatternNameInstance   *schema.Instance
	OverLengthLabelValueInstance *schema.Instance
	OverLengthAnnotationInstance *schema.Instance
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
}

type ImageConstraintsValidationV1Params struct {
	Workspace                 *schema.Workspace
	BlockStorage              *schema.BlockStorage
	OverLengthNameImage       *schema.Image
	InvalidPatternNameImage   *schema.Image
	OverLengthLabelValueImage *schema.Image
	OverLengthAnnotationImage *schema.Image
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
	Networks           []schema.Network
	InternetGateways   []schema.InternetGateway
	RouteTables        []schema.RouteTable
	Subnets            []schema.Subnet
	Nics               []schema.Nic
	PublicIps          []schema.PublicIp
	SecurityGroupRules []schema.SecurityGroupRule
	SecurityGroups     []schema.SecurityGroup
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
	Workspace                   *schema.Workspace
	InternetGateway             *schema.InternetGateway
	OverLengthNameNetwork       *schema.Network
	InvalidPatternNameNetwork   *schema.Network
	OverLengthLabelValueNetwork *schema.Network
	OverLengthAnnotationNetwork *schema.Network
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
}

type NicConstraintsValidationV1Params struct {
	Workspace               *schema.Workspace
	OverLengthNameNic       *schema.Nic
	InvalidPatternNameNic   *schema.Nic
	OverLengthLabelValueNic *schema.Nic
	OverLengthAnnotationNic *schema.Nic
}

type SecurityGroupConstraintsValidationV1Params struct {
	Workspace                         *schema.Workspace
	OverLengthNameSecurityGroup       *schema.SecurityGroup
	InvalidPatternNameSecurityGroup   *schema.SecurityGroup
	OverLengthLabelValueSecurityGroup *schema.SecurityGroup
	OverLengthAnnotationSecurityGroup *schema.SecurityGroup
}

type SecurityGroupRuleConstraintsValidationV1Params struct {
	Workspace                             *schema.Workspace
	OverLengthNameSecurityGroupRule       *schema.SecurityGroupRule
	InvalidPatternNameSecurityGroupRule   *schema.SecurityGroupRule
	OverLengthLabelValueSecurityGroupRule *schema.SecurityGroupRule
	OverLengthAnnotationSecurityGroupRule *schema.SecurityGroupRule
}

type RouteTableConstraintsValidationV1Params struct {
	Workspace                      *schema.Workspace
	Network                        *schema.Network
	InternetGateway                *schema.InternetGateway
	OverLengthNameRouteTable       *schema.RouteTable
	InvalidPatternNameRouteTable   *schema.RouteTable
	OverLengthLabelValueRouteTable *schema.RouteTable
	OverLengthAnnotationRouteTable *schema.RouteTable
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
}
