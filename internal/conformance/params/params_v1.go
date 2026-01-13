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
	*BaseParams
	Roles           []ResourceParams[schema.RoleSpec]
	RoleAssignments []ResourceParams[schema.RoleAssignmentSpec]
}

// Region

type RegionListParamsV1 struct {
	*BaseParams
	Regions []ResourceParams[schema.RegionSpec]
}

// Workspace

type WorkspaceLifeCycleParamsV1 struct {
	WorkspaceInitial *schema.Workspace
	WorkspaceUpdated *schema.Workspace
}

type WorkspaceListParamsV1 struct {
	*BaseParams
	Workspaces []ResourceParams[schema.WorkspaceSpec]
}

// Compute

type ComputeLifeCycleParamsV1 struct {
	Workspace       *schema.Workspace
	BlockStorage    *schema.BlockStorage
	InitialInstance *schema.Instance
	UpdatedInstance *schema.Instance
}

type ComputeListParamsV1 struct {
	*BaseParams
	Workspace    *ResourceParams[schema.WorkspaceSpec]
	BlockStorage *ResourceParams[schema.BlockStorageSpec]
	Instances    []ResourceParams[schema.InstanceSpec]
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
	*BaseParams
	Workspace     *ResourceParams[schema.WorkspaceSpec]
	BlockStorages []ResourceParams[schema.BlockStorageSpec]
	Images        []ResourceParams[schema.ImageSpec]
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
	*BaseParams

	Workspace        *ResourceParams[schema.WorkspaceSpec]
	BlockStorage     *ResourceParams[schema.BlockStorageSpec]
	Instance         *ResourceParams[schema.InstanceSpec]
	Networks         []ResourceParams[schema.NetworkSpec]
	InternetGateways []ResourceParams[schema.InternetGatewaySpec]
	RouteTables      []ResourceParams[schema.RouteTableSpec]
	Subnets          []ResourceParams[schema.SubnetSpec]
	Nics             []ResourceParams[schema.NicSpec]
	PublicIps        []ResourceParams[schema.PublicIpSpec]
	SecurityGroups   []ResourceParams[schema.SecurityGroupSpec]
}

// Usage

type FoundationUsageParamsV1 struct {
	*BaseParams

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
