package builders

import (
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func newNetworkResponse(name, provider, resource, apiVersion, tenant, workspace, region string, spec *schema.NetworkSpec) (*schema.Network, error) {
	medatata, err := NewRegionalWorkspaceResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.NetworkKind).
		Tenant(tenant).
		Workspace(workspace).
		Region(region).
		BuildResponse()
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
	medatata, err := NewRegionalWorkspaceResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.InternetGatewayKind).
		Tenant(tenant).
		Workspace(workspace).
		Region(region).
		BuildResponse()
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
	medatata, err := NewRegionalNetworkResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.RouteTableKind).
		Tenant(tenant).
		Workspace(workspace).
		Network(network).
		Region(region).
		BuildResponse()
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
	medatata, err := NewRegionalNetworkResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.SubnetKind).
		Tenant(tenant).
		Workspace(workspace).
		Network(network).
		Region(region).
		BuildResponse()
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
	medatata, err := NewRegionalWorkspaceResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.PublicIpKind).
		Tenant(tenant).
		Workspace(workspace).
		Region(region).
		BuildResponse()
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
	medatata, err := NewRegionalWorkspaceResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.NicKind).
		Tenant(tenant).
		Workspace(workspace).
		Region(region).
		BuildResponse()
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
	medatata, err := NewRegionalWorkspaceResourceMetadataBuilder().
		Name(name).
		Provider(provider).
		Resource(resource).
		ApiVersion(apiVersion).
		Kind(secalib.SecurityGroupKind).
		Tenant(tenant).
		Workspace(workspace).
		Region(region).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	return &schema.SecurityGroup{
		Metadata: medatata,
		Spec:     *spec,
		Status:   &schema.SecurityGroupStatus{},
	}, nil
}
