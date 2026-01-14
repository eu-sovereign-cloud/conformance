package params

import (
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Authorization

type AuthorizationLifeCycleParamsV1 struct {
	RoleInitial           *schema.Role
	RoleUpdated           *schema.Role
	RoleAssignmentInitial *schema.RoleAssignment
	RoleAssignmentUpdated *schema.RoleAssignment
}

type AuthorizationListParamsV1 struct {
	Roles           []schema.Role
	RoleAssignments []schema.RoleAssignment
}

// Region

type RegionListParamsV1 struct {
	Regions []schema.Region
}

// Workspace

type WorkspaceLifeCycleParamsV1 struct {
	WorkspaceInitial *schema.Workspace
	WorkspaceUpdated *schema.Workspace
}

type WorkspaceListParamsV1 struct {
	Workspaces []schema.Workspace
}

// Compute

type ComputeLifeCycleParamsV1 struct {
	Workspace       *schema.Workspace
	BlockStorage    *schema.BlockStorage
	InitialInstance *schema.Instance
	UpdatedInstance *schema.Instance
}

type ComputeListParamsV1 struct {
	Workspace    *schema.Workspace
	BlockStorage *schema.BlockStorage
	Instances    []schema.Instance
}

// Storage

type StorageLifeCycleParamsV1 struct {
	Workspace           *schema.Workspace
	BlockStorageInitial *schema.BlockStorage
	BlockStorageUpdated *schema.BlockStorage
	ImageInitial        *schema.Image
	ImageUpdated        *schema.Image
}

type StorageListParamsV1 struct {
	Workspace     *schema.Workspace
	BlockStorages []schema.BlockStorage
	Images        []schema.Image
}

// Network

type NetworkLifeCycleParamsV1 struct {
	Workspace              *schema.Workspace
	BlockStorage           *schema.BlockStorage
	Instance               *schema.Instance
	NetworkInitial         *schema.Network
	NetworkUpdated         *schema.Network
	InternetGatewayInitial *schema.InternetGateway
	InternetGatewayUpdated *schema.InternetGateway
	RouteTableInitial      *schema.RouteTable
	RouteTableUpdated      *schema.RouteTable
	SubnetInitial          *schema.Subnet
	SubnetUpdated          *schema.Subnet
	NicInitial             *schema.Nic
	NicUpdated             *schema.Nic
	PublicIpInitial        *schema.PublicIp
	PublicIpUpdated        *schema.PublicIp
	SecurityGroupInitial   *schema.SecurityGroup
	SecurityGroupUpdated   *schema.SecurityGroup
}

type NetworkListParamsV1 struct {
	Workspace        *schema.Workspace
	BlockStorage     *schema.BlockStorage
	Instance         *schema.Instance
	Networks         []schema.Network
	InternetGateways []schema.InternetGateway
	RouteTables      []schema.RouteTable
	Subnets          []schema.Subnet
	Nics             []schema.Nic
	PublicIps        []schema.PublicIp
	SecurityGroups   []schema.SecurityGroup
}

// Usage

type FoundationUsageParamsV1 struct {
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
