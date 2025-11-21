package mock

import (
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Authorization

func newRoleResponse(name, provider, resource, apiVersion, tenant string, spec *schema.RoleSpec) *schema.Role {
	return &schema.Role{
		Metadata: secalib.NewGlobalTenantResourceMetadata(name, provider, resource, apiVersion, secalib.RoleKind, tenant),
		Spec:     *spec,
		Status:   &schema.RoleStatus{},
	}
}

func newRoleAssignmentResponse(name, provider, resource, apiVersion, tenant string, spec *schema.RoleAssignmentSpec) *schema.RoleAssignment {
	return &schema.RoleAssignment{
		Metadata: secalib.NewGlobalTenantResourceMetadata(name, provider, resource, apiVersion, secalib.RoleAssignmentKind, tenant),
		Spec:     *spec,
		Status:   &schema.RoleAssignmentStatus{},
	}
}

// Region
func newRegionResponse(name, provider, resource, apiVersion string, spec *schema.RegionSpec) *schema.Region {
	return &schema.Region{
		Metadata: secalib.NewGlobalResourceMetadata(name, provider, resource, apiVersion, secalib.RegionKind),
		Spec:     *spec,
	}
}

// Workspace

func newWorkspaceResponse(name, provider, resource, apiVersion, tenant, region string, labels schema.Labels) *schema.Workspace {
	return &schema.Workspace{
		Metadata: secalib.NewRegionalResourceMetadata(name, provider, resource, apiVersion, secalib.WorkspaceKind, tenant, region),
		Labels:   labels,
		Spec:     schema.WorkspaceSpec{},
		Status:   &schema.WorkspaceStatus{},
	}
}

// Storage

func newBlockStorageResponse(name, provider, resource, apiVersion, tenant, workspace, region string, labels schema.Labels, spec *schema.BlockStorageSpec) *schema.BlockStorage {
	return &schema.BlockStorage{
		Metadata: secalib.NewRegionalWorkspaceResourceMetadata(name, provider, resource, apiVersion, secalib.BlockStorageKind, tenant, workspace, region),
		Labels:   labels,
		Spec:     *spec,
		Status:   &schema.BlockStorageStatus{},
	}
}

func newImageResponse(name, provider, resource, apiVersion, tenant, region string, label *schema.Labels, spec *schema.ImageSpec) *schema.Image {
	return &schema.Image{
		Metadata: secalib.NewRegionalResourceMetadata(name, provider, resource, apiVersion, secalib.ImageKind, tenant, region),
		Labels:   *label,
		Spec:     *spec,
		Status:   &schema.ImageStatus{},
	}
}

// Compute

func newInstanceResponse(name, provider, resource, apiVersion, tenant, workspace, region string, label *schema.Labels, spec *schema.InstanceSpec) *schema.Instance {
	return &schema.Instance{
		Metadata: secalib.NewRegionalWorkspaceResourceMetadata(name, provider, resource, apiVersion, secalib.InstanceKind, tenant, workspace, region),
		Labels:   *label,
		Spec:     *spec,
		Status:   &schema.InstanceStatus{},
	}
}

// Network

func newNetworkResponse(name, provider, resource, apiVersion, tenant, workspace, region string, label *schema.Labels, spec *schema.NetworkSpec) *schema.Network {
	return &schema.Network{
		Metadata: secalib.NewRegionalWorkspaceResourceMetadata(name, provider, resource, apiVersion, secalib.NetworkKind, tenant, workspace, region),
		Spec:     *spec,
		Status:   &schema.NetworkStatus{},
	}
}

func newInternetGatewayResponse(name, provider, resource, apiVersion, tenant, workspace, region string, spec *schema.InternetGatewaySpec) *schema.InternetGateway {
	return &schema.InternetGateway{
		Metadata: secalib.NewRegionalWorkspaceResourceMetadata(name, provider, resource, apiVersion, secalib.InternetGatewayKind, tenant, workspace, region),
		Spec:     *spec,
		Status:   &schema.InternetGatewayStatus{},
	}
}

func newRouteTableResponse(name, provider, resource, apiVersion, tenant, workspace, network, region string, spec *schema.RouteTableSpec) *schema.RouteTable {
	return &schema.RouteTable{
		Metadata: secalib.NewRegionalNetworkResourceMetadata(name, provider, resource, apiVersion, secalib.RouteTableKind, tenant, workspace, network, region),
		Spec:     *spec,
		Status:   &schema.RouteTableStatus{},
	}
}

func newSubnetResponse(name, provider, resource, apiVersion, tenant, workspace, network, region string, spec *schema.SubnetSpec) *schema.Subnet {
	return &schema.Subnet{
		Metadata: secalib.NewRegionalNetworkResourceMetadata(name, provider, resource, apiVersion, secalib.SubnetKind, tenant, workspace, network, region),
		Spec:     *spec,
		Status:   &schema.SubnetStatus{},
	}
}

func newPublicIpResponse(name, provider, resource, apiVersion, tenant, workspace, region string, spec *schema.PublicIpSpec) *schema.PublicIp {
	return &schema.PublicIp{
		Metadata: secalib.NewRegionalWorkspaceResourceMetadata(name, provider, resource, apiVersion, secalib.PublicIpKind, tenant, workspace, region),
		Spec:     *spec,
		Status:   &schema.PublicIpStatus{},
	}
}

func newNicResponse(name, provider, resource, apiVersion, tenant, workspace, region string, spec *schema.NicSpec) *schema.Nic {
	return &schema.Nic{
		Metadata: secalib.NewRegionalWorkspaceResourceMetadata(name, provider, resource, apiVersion, secalib.NicKind, tenant, workspace, region),
		Spec:     *spec,
		Status:   &schema.NicStatus{},
	}
}

func newSecurityGroupResponse(name, provider, resource, apiVersion, tenant, workspace, region string, spec *schema.SecurityGroupSpec) *schema.SecurityGroup {
	return &schema.SecurityGroup{
		Metadata: secalib.NewRegionalWorkspaceResourceMetadata(name, provider, resource, apiVersion, secalib.SecurityGroupKind, tenant, workspace, region),
		Spec:     *spec,
		Status:   &schema.SecurityGroupStatus{},
	}
}
