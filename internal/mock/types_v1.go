package mock

import "github.com/eu-sovereign-cloud/conformance/secalib"

// Params

type AuthorizationParamsV1 struct {
	*Params
	Role           *ResourceParams[secalib.RoleSpecV1]
	RoleAssignment *ResourceParams[secalib.RoleAssignmentSpecV1]
}

func (p AuthorizationParamsV1) getParams() *Params { return p.Params }

type RegionParamsV1 struct {
	*Params
	Regions []ResourceParams[secalib.RegionSpecV1]
}

func (p RegionParamsV1) getParams() *Params { return p.Params }

type WorkspaceParamsV1 struct {
	*Params
	Workspace *ResourceParams[secalib.WorkspaceSpecV1]
}

func (p WorkspaceParamsV1) getParams() *Params { return p.Params }

type ComputeParamsV1 struct {
	*Params
	Workspace    *ResourceParams[secalib.WorkspaceSpecV1]
	BlockStorage *ResourceParams[secalib.BlockStorageSpecV1]
	Instance     *ResourceParams[secalib.InstanceSpecV1]
}

func (p ComputeParamsV1) getParams() *Params { return p.Params }

type StorageParamsV1 struct {
	*Params
	Workspace    *ResourceParams[secalib.WorkspaceSpecV1]
	BlockStorage *ResourceParams[secalib.BlockStorageSpecV1]
	Image        *ResourceParams[secalib.ImageSpecV1]
}

func (p StorageParamsV1) getParams() *Params { return p.Params }

type NetworkParamsV1 struct {
	*Params

	Workspace       *ResourceParams[secalib.WorkspaceSpecV1]
	BlockStorage    *ResourceParams[secalib.BlockStorageSpecV1]
	Instance        *ResourceParams[secalib.InstanceSpecV1]
	Network         *ResourceParams[secalib.NetworkSpecV1]
	InternetGateway *ResourceParams[secalib.InternetGatewaySpecV1]
	RouteTable      *ResourceParams[secalib.RouteTableSpecV1]
	Subnet          *ResourceParams[secalib.SubnetSpecV1]
	NIC             *ResourceParams[secalib.NICSpecV1]
	PublicIP        *ResourceParams[secalib.PublicIpSpecV1]
	SecurityGroup   *ResourceParams[secalib.SecurityGroupSpecV1]
}

func (p NetworkParamsV1) getParams() *Params { return p.Params }

type UsageParamsV1 struct {
	*Params
	Workspace     *WorkspaceParamsV1
	Storage       *StorageParamsV1
	Compute       *ComputeParamsV1
	Network       *NetworkParamsV1
	Authorization *AuthorizationParamsV1
}

func (p UsageParamsV1) getParams() *Params { return p.Params }
