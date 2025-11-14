package mock

import (
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Authorization

func newRoleResponse(name, provider, resource, apiVersion, tenant string, spec *schema.RoleSpec) (*schema.Role, error) {
	medatata, err := secalib.NewGlobalTenantResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.RoleKind).
		Tenant(tenant).
		Build()
	if err != nil {
		return nil, err
	}

	return &schema.Role{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.RoleStatus{},
	}, nil
}

func newRoleAssignmentResponse(name, provider, resource, apiVersion, tenant string, spec *schema.RoleAssignmentSpec) (*schema.RoleAssignment, error) {
	medatata, err := secalib.NewGlobalTenantResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.RoleAssignmentKind).
		Tenant(tenant).
		Build()
	if err != nil {
		return nil, err
	}

	return &schema.RoleAssignment{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.RoleAssignmentStatus{},
	}, nil
}

// Workspace

func newWorkspaceResponse(name, provider, resource, apiVersion, tenant, region string, labels schema.Labels) (*schema.Workspace, error) {
	medatata, err := secalib.NewRegionalResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.WorkspaceKind).
		Tenant(tenant).
		Region(region).
		Build()
	if err != nil {
		return nil, err
	}

	return &schema.Workspace{
		Metadata: medatata,
		Labels:   labels,
		Spec:     schema.WorkspaceSpec{},
		Status:   &schema.WorkspaceStatus{},
	}, nil
}

// Storage

func newBlockStorageResponse(name, provider, resource, apiVersion, tenant, workspace, region string, spec *schema.BlockStorageSpec) (*schema.BlockStorage, error) {
	medatata, err := secalib.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.BlockStorageKind).
		Tenant(tenant).
		Workspace(workspace).
		Region(region).
		Build()
	if err != nil {
		return nil, err
	}

	return &schema.BlockStorage{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.BlockStorageStatus{},
	}, nil
}

func newImageResponse(name, provider, resource, apiVersion, tenant, region string, spec *schema.ImageSpec) (*schema.Image, error) {
	medatata, err := secalib.NewRegionalResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.ImageKind).
		Tenant(tenant).
		Region(region).
		Build()
	if err != nil {
		return nil, err
	}

	return &schema.Image{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.ImageStatus{},
	}, nil
}

// Compute

func newInstanceResponse(name, provider, resource, apiVersion, tenant, workspace, region string, spec *schema.InstanceSpec) (*schema.Instance, error) {
	medatata, err := secalib.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.InstanceKind).
		Tenant(tenant).
		Workspace(workspace).
		Region(region).
		Build()
	if err != nil {
		return nil, err
	}

	return &schema.Instance{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.InstanceStatus{},
	}, nil
}

// Network

func newNetworkResponse(name, provider, resource, apiVersion, tenant, workspace, region string, spec *schema.NetworkSpec) (*schema.Network, error) {
	medatata, err := secalib.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.NetworkKind).
		Tenant(tenant).
		Workspace(workspace).
		Region(region).
		Build()
	if err != nil {
		return nil, err
	}

	return &schema.Network{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.NetworkStatus{},
	}, nil
}

func newInternetGatewayResponse(name, provider, resource, apiVersion, tenant, workspace, region string, spec *schema.InternetGatewaySpec) (*schema.InternetGateway, error) {
	medatata, err := secalib.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.InternetGatewayKind).
		Tenant(tenant).
		Workspace(workspace).
		Region(region).
		Build()
	if err != nil {
		return nil, err
	}

	return &schema.InternetGateway{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.InternetGatewayStatus{},
	}, nil
}

func newRouteTableResponse(name, provider, resource, apiVersion, tenant, workspace, network, region string, spec *schema.RouteTableSpec) (*schema.RouteTable, error) {
	medatata, err := secalib.NewRegionalNetworkResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.RouteTableKind).
		Tenant(tenant).
		Workspace(workspace).
		Network(network).
		Region(region).
		Build()
	if err != nil {
		return nil, err
	}

	return &schema.RouteTable{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.RouteTableStatus{},
	}, nil
}

func newSubnetResponse(name, provider, resource, apiVersion, tenant, workspace, network, region string, spec *schema.SubnetSpec) (*schema.Subnet, error) {
	medatata, err := secalib.NewRegionalNetworkResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.SubnetKind).
		Tenant(tenant).
		Workspace(workspace).
		Network(network).
		Region(region).
		Build()
	if err != nil {
		return nil, err
	}

	return &schema.Subnet{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.SubnetStatus{},
	}, nil
}

func newPublicIpResponse(name, provider, resource, apiVersion, tenant, workspace, region string, spec *schema.PublicIpSpec) (*schema.PublicIp, error) {
	medatata, err := secalib.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.PublicIpKind).
		Tenant(tenant).
		Workspace(workspace).
		Region(region).
		Build()
	if err != nil {
		return nil, err
	}

	return &schema.PublicIp{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.PublicIpStatus{},
	}, nil
}

func newNicResponse(name, provider, resource, apiVersion, tenant, workspace, region string, spec *schema.NicSpec) (*schema.Nic, error) {
	medatata, err := secalib.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.NicKind).
		Tenant(tenant).
		Workspace(workspace).
		Region(region).
		Build()
	if err != nil {
		return nil, err
	}

	return &schema.Nic{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.NicStatus{},
	}, nil
}

func newSecurityGroupResponse(name, provider, resource, apiVersion, tenant, workspace, region string, spec *schema.SecurityGroupSpec) (*schema.SecurityGroup, error) {
	medatata, err := secalib.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.SecurityGroupKind).
		Tenant(tenant).
		Workspace(workspace).
		Region(region).
		Build()
	if err != nil {
		return nil, err
	}

	return &schema.SecurityGroup{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.SecurityGroupStatus{},
	}, nil
}
