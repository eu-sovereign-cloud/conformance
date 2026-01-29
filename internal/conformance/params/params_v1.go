package params

import (
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Authorization

type AuthorizationLifeCycleV1Params struct {
	RoleInitial           *schema.Role
	RoleUpdated           *schema.Role
	RoleAssignmentInitial *schema.RoleAssignment
	RoleAssignmentUpdated *schema.RoleAssignment
}

type AuthorizationListV1Params struct {
	Roles           []schema.Role
	RoleAssignments []schema.RoleAssignment
}

// Region

type RegionListV1Params struct {
	Regions []schema.Region
}

// Workspace

type WorkspaceLifeCycleV1Params struct {
	WorkspaceInitial *schema.Workspace
	WorkspaceUpdated *schema.Workspace
}

type WorkspaceListV1Params struct {
	Workspaces []schema.Workspace
}

type CreateWorkspaceV1Params struct {
	Workspace *schema.Workspace
}

// Compute

type ComputeLifeCycleV1Params struct {
	Workspace       *schema.Workspace
	BlockStorage    *schema.BlockStorage
	InitialInstance *schema.Instance
	UpdatedInstance *schema.Instance
}

type ComputeListV1Params struct {
	Workspace    *schema.Workspace
	BlockStorage *schema.BlockStorage
	Instances    []schema.Instance
}

// Storage

type StorageLifeCycleV1Params struct {
	Workspace           *schema.Workspace
	BlockStorageInitial *schema.BlockStorage
	BlockStorageUpdated *schema.BlockStorage
	ImageInitial        *schema.Image
	ImageUpdated        *schema.Image
}

type StorageListV1Params struct {
	Workspace     *schema.Workspace
	BlockStorages []schema.BlockStorage
	Images        []schema.Image
}

type CreateBlockStorageV1Params struct {
	Workspace    *schema.Workspace
	BlockStorage *schema.BlockStorage
}

type UpdateBlockStorageV1Params struct {
	Workspace           *schema.Workspace
	BlockStorageInitial *schema.BlockStorage
	BlockStorageUpdated *schema.BlockStorage
}

// Network

type NetworkLifeCycleV1Params struct {
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

type NetworkListV1Params struct {
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
