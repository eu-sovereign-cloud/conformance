package mock

import (
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Params

type ClientsInitParams struct {
	*Params
}

func (p ClientsInitParams) getParams() *Params { return p.Params }

type AuthorizationParamsV1 struct {
	*Params
	Role           *ResourceParams[schema.RoleSpec]
	RoleAssignment *ResourceParams[schema.RoleAssignmentSpec]
}

func (p AuthorizationParamsV1) getParams() *Params { return p.Params }

type RegionParamsV1 struct {
	*Params
	Regions []ResourceParams[schema.RegionSpec]
}

func (p RegionParamsV1) getParams() *Params { return p.Params }

type WorkspaceParamsV1 struct {
	*Params
	Workspace *ResourceParams[schema.WorkspaceSpec]
}

func (p WorkspaceParamsV1) getParams() *Params { return p.Params }

type WorkspaceListParamsV1 struct {
	*Params
	Workspace *[]ResourceParams[schema.WorkspaceSpec]
}

func (p WorkspaceListParamsV1) getParams() *Params { return p.Params }

type ComputeParamsV1 struct {
	*Params
	Workspace    *ResourceParams[schema.WorkspaceSpec]
	BlockStorage *ResourceParams[schema.BlockStorageSpec]
	Instance     *ResourceParams[schema.InstanceSpec]
}

func (p ComputeParamsV1) getParams() *Params { return p.Params }

type StorageParamsV1 struct {
	*Params
	Workspace    *ResourceParams[schema.WorkspaceSpec]
	BlockStorage *ResourceParams[schema.BlockStorageSpec]
	Image        *ResourceParams[schema.ImageSpec]
}

func (p StorageParamsV1) getParams() *Params { return p.Params }

type NetworkParamsV1 struct {
	*Params

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

func (p NetworkParamsV1) getParams() *Params { return p.Params }

type FoundationUsageParamsV1 struct {
	*Params

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

func (p FoundationUsageParamsV1) getParams() *Params { return p.Params }
