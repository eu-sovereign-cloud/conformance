package mock

import "github.com/eu-sovereign-cloud/conformance/secalib"

// Params

type AuthorizationParamsV1 struct {
	Role           *secalib.ResourceParams[secalib.RoleSpecV1]
	RoleAssignment *secalib.ResourceParams[secalib.RoleAssignmentSpecV1]
}

type WorkspaceParamsV1 struct {
	Workspace *secalib.ResourceParams[secalib.WorkspaceSpecV1]
}

type ComputeParamsV1 struct {
	Workspace    *secalib.ResourceParams[secalib.WorkspaceSpecV1]
	BlockStorage *secalib.ResourceParams[secalib.BlockStorageSpecV1]
	Instance     *secalib.ResourceParams[secalib.InstanceSpecV1]
}

type StorageParamsV1 struct {
	Workspace    *secalib.ResourceParams[secalib.WorkspaceSpecV1]
	BlockStorage *secalib.ResourceParams[secalib.BlockStorageSpecV1]
	Image        *secalib.ResourceParams[secalib.ImageSpecV1]
}

type NetworkParamsV1 struct {
	Workspace       *secalib.ResourceParams[secalib.WorkspaceSpecV1]
	BlockStorage    *secalib.ResourceParams[secalib.BlockStorageSpecV1]
	Instance        *secalib.ResourceParams[secalib.InstanceSpecV1]
	Network         *secalib.ResourceParams[secalib.NetworkSpecV1]
	InternetGateway *secalib.ResourceParams[secalib.InternetGatewaySpecV1]
	RouteTable      *secalib.ResourceParams[secalib.RouteTableSpecV1]
	Subnet          *secalib.ResourceParams[secalib.SubnetSpecV1]
	NIC             *secalib.ResourceParams[secalib.NICSpecV1]
	PublicIP        *secalib.ResourceParams[secalib.PublicIpSpecV1]
	SecurityGroup   *secalib.ResourceParams[secalib.SecurityGroupSpecV1]
}

type UsageParamsV1 struct {
	Role            *secalib.ResourceParams[secalib.RoleSpecV1]
	RoleAssignment  *secalib.ResourceParams[secalib.RoleAssignmentSpecV1]
	Workspace       *secalib.ResourceParams[secalib.WorkspaceSpecV1]
	BlockStorage    *secalib.ResourceParams[secalib.BlockStorageSpecV1]
	Image           *secalib.ResourceParams[secalib.ImageSpecV1]
	Network         *secalib.ResourceParams[secalib.NetworkSpecV1]
	InternetGateway *secalib.ResourceParams[secalib.InternetGatewaySpecV1]
	RouteTable      *secalib.ResourceParams[secalib.RouteTableSpecV1]
	Subnet          *secalib.ResourceParams[secalib.SubnetSpecV1]
	NIC             *secalib.ResourceParams[secalib.NICSpecV1]
	PublicIP        *secalib.ResourceParams[secalib.PublicIpSpecV1]
	SecurityGroup   *secalib.ResourceParams[secalib.SecurityGroupSpecV1]
	Instance        *secalib.ResourceParams[secalib.InstanceSpecV1]
}
