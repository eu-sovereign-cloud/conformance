package mock

import (
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Params

type ClientsInitParams struct {
	*BaseParams
}

type AuthorizationLifeCycleParamsV1 struct {
	*BaseParams
	Role           *ResourceParams[schema.RoleSpec]
	RoleAssignment *ResourceParams[schema.RoleAssignmentSpec]
}

type AuthorizationListParamsV1 struct {
	*BaseParams
	Roles           []ResourceParams[schema.RoleSpec]
	RoleAssignments []ResourceParams[schema.RoleAssignmentSpec]
}

type RegionListParamsV1 struct {
	*BaseParams
	Regions []ResourceParams[schema.RegionSpec]
}

type WorkspaceLifeCycleParamsV1 struct {
	*BaseParams
	Workspace *ResourceParams[schema.WorkspaceSpec]
}

type WorkspaceListParamsV1 struct {
	*BaseParams
	Workspaces []ResourceParams[schema.WorkspaceSpec]
}

type ComputeLifeCycleParamsV1 struct {
	*BaseParams
	Workspace    *ResourceParams[schema.WorkspaceSpec]
	BlockStorage *ResourceParams[schema.BlockStorageSpec]
	Instance     *ResourceParams[schema.InstanceSpec]
}

type ComputeListParamsV1 struct {
	*BaseParams
	Workspace    *ResourceParams[schema.WorkspaceSpec]
	BlockStorage *ResourceParams[schema.BlockStorageSpec]
	Instances    []ResourceParams[schema.InstanceSpec]
}

type StorageLifeCycleParamsV1 struct {
	*BaseParams
	Workspace    *ResourceParams[schema.WorkspaceSpec]
	BlockStorage *ResourceParams[schema.BlockStorageSpec]
	Image        *ResourceParams[schema.ImageSpec]
}

type StorageListParamsV1 struct {
	*BaseParams
	Workspace     *ResourceParams[schema.WorkspaceSpec]
	BlockStorages []ResourceParams[schema.BlockStorageSpec]
	Images        []ResourceParams[schema.ImageSpec]
}

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
