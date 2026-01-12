package params

import (
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Authorization

type AuthorizationLifeCycleParamsV1 struct {
	*mock.MockParams
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
	*mock.MockParams
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
	*mock.MockParams
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
	*BaseParams

	Workspace       *ResourceParams[schema.WorkspaceSpec]
	BlockStorage    *ResourceParams[schema.BlockStorageSpec]
	Instance        *ResourceParams[schema.InstanceSpec]
	Network         *ResourceParams[schema.NetworkSpec]
	InternetGateway *ResourceParams[schema.InternetGatewaySpec]
	RouteTable      *ResourceParams[schema.RouteTableSpec]
	Subnet          *ResourceParams[schema.SubnetSpec]
	Nic             *ResourceParams[schema.NicSpec]
	PublicIp        *ResourceParams[schema.PublicIpSpec]
	SecurityGroup   *ResourceParams[schema.SecurityGroupSpec]
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

	Role            *ResourceParams[schema.RoleSpec]
	RoleAssignment  *ResourceParams[schema.RoleAssignmentSpec]
	Workspace       *ResourceParams[schema.WorkspaceSpec]
	Image           *ResourceParams[schema.ImageSpec]
	BlockStorage    *ResourceParams[schema.BlockStorageSpec]
	Network         *ResourceParams[schema.NetworkSpec]
	InternetGateway *ResourceParams[schema.InternetGatewaySpec]
	RouteTable      *ResourceParams[schema.RouteTableSpec]
	Subnet          *ResourceParams[schema.SubnetSpec]
	SecurityGroup   *ResourceParams[schema.SecurityGroupSpec]
	PublicIp        *ResourceParams[schema.PublicIpSpec]
	Nic             *ResourceParams[schema.NicSpec]
	Instance        *ResourceParams[schema.InstanceSpec]
}
