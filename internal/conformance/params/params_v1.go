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

// Network

type NetworkProviderLifeCycleV1Params struct {
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

type NetworkProviderQueriesV1Params struct {
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
