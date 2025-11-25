package builders

import (
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Network

type NetworkBuilder struct {
	*resourceBuilder
	metadata *RegionalWorkspaceResourceMetadataBuilder
	spec     *schema.NetworkSpec
}

func NewNetworkBuilder() *NetworkBuilder {
	return &NetworkBuilder{
		resourceBuilder: newResourceBuilder(),
		metadata:        NewRegionalWorkspaceResourceMetadataBuilder(),
		spec:            &schema.NetworkSpec{},
	}
}

func (builder *NetworkBuilder) Name(name string) *NetworkBuilder {
	builder.metadata.Name(name)
	return builder
}

func (builder *NetworkBuilder) Provider(provider string) *NetworkBuilder {
	builder.metadata.Provider(provider)
	return builder
}

func (builder *NetworkBuilder) Resource(resource string) *NetworkBuilder {
	builder.metadata.Resource(resource)
	return builder
}

func (builder *NetworkBuilder) ApiVersion(apiVersion string) *NetworkBuilder {
	builder.metadata.ApiVersion(apiVersion)
	return builder
}

func (builder *NetworkBuilder) Tenant(tenant string) *NetworkBuilder {
	builder.metadata.Tenant(tenant)
	return builder
}

func (builder *NetworkBuilder) Workspace(workspace string) *NetworkBuilder {
	builder.metadata.Workspace(workspace)
	return builder
}

func (builder *NetworkBuilder) Region(region string) *NetworkBuilder {
	builder.metadata.Region(region)
	return builder
}

func (builder *NetworkBuilder) Spec(spec *schema.NetworkSpec) *NetworkBuilder {
	builder.spec = spec
	return builder
}

func (builder *NetworkBuilder) BuildResponse() (*schema.Network, error) {
	metadata, err := builder.metadata.Kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork).BuildResponse()
	if err != nil {
		return nil, err
	}

	// Validate the spec
	if err := builder.validator.ValidateRequired(
		builder.spec,
		builder.spec.Cidr,
		builder.spec.SkuRef,
		builder.spec.RouteTableRef,
	); err != nil {
		return nil, err
	}
	if err := builder.validator.ValidateOneRequired(
		builder.spec.Cidr.Ipv4,
		builder.spec.Cidr.Ipv6,
	); err != nil {
		return nil, err
	}

	return &schema.Network{
		Metadata: metadata,
		Spec:     *builder.spec,
		Status:   &schema.NetworkStatus{},
	}, nil
}

// Internet Gateway

type InternetGatewayBuilder struct {
	*resourceBuilder
	metadata *RegionalWorkspaceResourceMetadataBuilder
	spec     *schema.InternetGatewaySpec
}

func NewInternetGatewayBuilder() *InternetGatewayBuilder {
	return &InternetGatewayBuilder{
		resourceBuilder: newResourceBuilder(),
		metadata:        NewRegionalWorkspaceResourceMetadataBuilder(),
		spec:            &schema.InternetGatewaySpec{},
	}
}

func (builder *InternetGatewayBuilder) Name(name string) *InternetGatewayBuilder {
	builder.metadata.Name(name)
	return builder
}

func (builder *InternetGatewayBuilder) Provider(provider string) *InternetGatewayBuilder {
	builder.metadata.Provider(provider)
	return builder
}

func (builder *InternetGatewayBuilder) Resource(resource string) *InternetGatewayBuilder {
	builder.metadata.Resource(resource)
	return builder
}

func (builder *InternetGatewayBuilder) ApiVersion(apiVersion string) *InternetGatewayBuilder {
	builder.metadata.ApiVersion(apiVersion)
	return builder
}

func (builder *InternetGatewayBuilder) Tenant(tenant string) *InternetGatewayBuilder {
	builder.metadata.Tenant(tenant)
	return builder
}

func (builder *InternetGatewayBuilder) Workspace(workspace string) *InternetGatewayBuilder {
	builder.metadata.Workspace(workspace)
	return builder
}

func (builder *InternetGatewayBuilder) Region(region string) *InternetGatewayBuilder {
	builder.metadata.Region(region)
	return builder
}

func (builder *InternetGatewayBuilder) Spec(spec *schema.InternetGatewaySpec) *InternetGatewayBuilder {
	builder.spec = spec
	return builder
}

func (builder *InternetGatewayBuilder) BuildResponse() (*schema.InternetGateway, error) {
	metadata, err := builder.metadata.Kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindInternetGateway).BuildResponse()
	if err != nil {
		return nil, err
	}

	// Validate the spec
	if err := builder.validator.ValidateRequired(
		builder.spec,
		builder.spec.EgressOnly,
	); err != nil {
		return nil, err
	}

	return &schema.InternetGateway{
		Metadata: metadata,
		Spec:     *builder.spec,
		Status:   &schema.InternetGatewayStatus{},
	}, nil
}

// Route Table

type RouteTableBuilder struct {
	*resourceBuilder
	metadata *RegionalNetworkResourceMetadataBuilder
	spec     *schema.RouteTableSpec
}

func NewRouteTableBuilder() *RouteTableBuilder {
	return &RouteTableBuilder{
		resourceBuilder: newResourceBuilder(),
		metadata:        NewRegionalNetworkResourceMetadataBuilder(),
		spec:            &schema.RouteTableSpec{},
	}
}

func (builder *RouteTableBuilder) Name(name string) *RouteTableBuilder {
	builder.metadata.Name(name)
	return builder
}

func (builder *RouteTableBuilder) Provider(provider string) *RouteTableBuilder {
	builder.metadata.Provider(provider)
	return builder
}

func (builder *RouteTableBuilder) Resource(resource string) *RouteTableBuilder {
	builder.metadata.Resource(resource)
	return builder
}

func (builder *RouteTableBuilder) ApiVersion(apiVersion string) *RouteTableBuilder {
	builder.metadata.ApiVersion(apiVersion)
	return builder
}

func (builder *RouteTableBuilder) Tenant(tenant string) *RouteTableBuilder {
	builder.metadata.Tenant(tenant)
	return builder
}

func (builder *RouteTableBuilder) Workspace(workspace string) *RouteTableBuilder {
	builder.metadata.Workspace(workspace)
	return builder
}

func (builder *RouteTableBuilder) Network(network string) *RouteTableBuilder {
	builder.metadata.Network(network)
	return builder
}

func (builder *RouteTableBuilder) Region(region string) *RouteTableBuilder {
	builder.metadata.Region(region)
	return builder
}

func (builder *RouteTableBuilder) Spec(spec *schema.RouteTableSpec) *RouteTableBuilder {
	builder.spec = spec
	return builder
}

func (builder *RouteTableBuilder) BuildResponse() (*schema.RouteTable, error) {
	metadata, err := builder.metadata.Kind(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable).BuildResponse()
	if err != nil {
		return nil, err
	}

	// Validate the spec
	if err := builder.validator.ValidateRequired(
		builder.spec.Routes,
	); err != nil {
		return nil, err
	}
	// Validate each route
	for _, route := range builder.spec.Routes {
		if err := builder.validator.ValidateRequired(
			route.DestinationCidrBlock,
			route.TargetRef,
		); err != nil {
			return nil, err
		}
	}

	return &schema.RouteTable{
		Metadata: metadata,
		Spec:     *builder.spec,
		Status:   &schema.RouteTableStatus{},
	}, nil
}

// Subnet

type SubnetBuilder struct {
	*resourceBuilder
	metadata *RegionalNetworkResourceMetadataBuilder
	spec     *schema.SubnetSpec
}

func NewSubnetBuilder() *SubnetBuilder {
	return &SubnetBuilder{
		resourceBuilder: newResourceBuilder(),
		metadata:        NewRegionalNetworkResourceMetadataBuilder(),
		spec:            &schema.SubnetSpec{},
	}
}

func (builder *SubnetBuilder) Name(name string) *SubnetBuilder {
	builder.metadata.Name(name)
	return builder
}

func (builder *SubnetBuilder) Provider(provider string) *SubnetBuilder {
	builder.metadata.Provider(provider)
	return builder
}

func (builder *SubnetBuilder) Resource(resource string) *SubnetBuilder {
	builder.metadata.Resource(resource)
	return builder
}

func (builder *SubnetBuilder) ApiVersion(apiVersion string) *SubnetBuilder {
	builder.metadata.ApiVersion(apiVersion)
	return builder
}

func (builder *SubnetBuilder) Tenant(tenant string) *SubnetBuilder {
	builder.metadata.Tenant(tenant)
	return builder
}

func (builder *SubnetBuilder) Workspace(workspace string) *SubnetBuilder {
	builder.metadata.Workspace(workspace)
	return builder
}

func (builder *SubnetBuilder) Network(network string) *SubnetBuilder {
	builder.metadata.Network(network)
	return builder
}

func (builder *SubnetBuilder) Region(region string) *SubnetBuilder {
	builder.metadata.Region(region)
	return builder
}

func (builder *SubnetBuilder) Spec(spec *schema.SubnetSpec) *SubnetBuilder {
	builder.spec = spec
	return builder
}

func (builder *SubnetBuilder) BuildResponse() (*schema.Subnet, error) {
	metadata, err := builder.metadata.Kind(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet).BuildResponse()
	if err != nil {
		return nil, err
	}

	// Validate the spec
	if err := builder.validator.ValidateRequired(
		builder.spec,
		builder.spec.Cidr,
		builder.spec.Zone,
	); err != nil {
		return nil, err
	}

	return &schema.Subnet{
		Metadata: metadata,
		Spec:     *builder.spec,
		Status:   &schema.SubnetStatus{},
	}, nil
}

// Public Ip

type PublicIpBuilder struct {
	*resourceBuilder
	metadata *RegionalWorkspaceResourceMetadataBuilder
	spec     *schema.PublicIpSpec
}

func NewPublicIpBuilder() *PublicIpBuilder {
	return &PublicIpBuilder{
		resourceBuilder: newResourceBuilder(),
		metadata:        NewRegionalWorkspaceResourceMetadataBuilder(),
		spec:            &schema.PublicIpSpec{},
	}
}

func (builder *PublicIpBuilder) Name(name string) *PublicIpBuilder {
	builder.metadata.Name(name)
	return builder
}

func (builder *PublicIpBuilder) Provider(provider string) *PublicIpBuilder {
	builder.metadata.Provider(provider)
	return builder
}

func (builder *PublicIpBuilder) Resource(resource string) *PublicIpBuilder {
	builder.metadata.Resource(resource)
	return builder
}

func (builder *PublicIpBuilder) ApiVersion(apiVersion string) *PublicIpBuilder {
	builder.metadata.ApiVersion(apiVersion)
	return builder
}

func (builder *PublicIpBuilder) Tenant(tenant string) *PublicIpBuilder {
	builder.metadata.Tenant(tenant)
	return builder
}

func (builder *PublicIpBuilder) Workspace(workspace string) *PublicIpBuilder {
	builder.metadata.Workspace(workspace)
	return builder
}

func (builder *PublicIpBuilder) Region(region string) *PublicIpBuilder {
	builder.metadata.Region(region)
	return builder
}

func (builder *PublicIpBuilder) Spec(spec *schema.PublicIpSpec) *PublicIpBuilder {
	builder.spec = spec
	return builder
}

func (builder *PublicIpBuilder) BuildResponse() (*schema.PublicIp, error) {
	metadata, err := builder.metadata.Kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindPublicIP).BuildResponse()
	if err != nil {
		return nil, err
	}

	// Validate the spec
	if err := builder.validator.ValidateRequired(
		builder.spec,
		builder.spec.Version,
	); err != nil {
		return nil, err
	}

	return &schema.PublicIp{
		Metadata: metadata,
		Spec:     *builder.spec,
		Status:   &schema.PublicIpStatus{},
	}, nil
}

// Nic

type NicBuilder struct {
	*resourceBuilder
	metadata *RegionalWorkspaceResourceMetadataBuilder
	spec     *schema.NicSpec
}

func NewNicBuilder() *NicBuilder {
	return &NicBuilder{
		resourceBuilder: newResourceBuilder(),
		metadata:        NewRegionalWorkspaceResourceMetadataBuilder(),
		spec:            &schema.NicSpec{},
	}
}

func (builder *NicBuilder) Name(name string) *NicBuilder {
	builder.metadata.Name(name)
	return builder
}

func (builder *NicBuilder) Provider(provider string) *NicBuilder {
	builder.metadata.Provider(provider)
	return builder
}

func (builder *NicBuilder) Resource(resource string) *NicBuilder {
	builder.metadata.Resource(resource)
	return builder
}

func (builder *NicBuilder) ApiVersion(apiVersion string) *NicBuilder {
	builder.metadata.ApiVersion(apiVersion)
	return builder
}

func (builder *NicBuilder) Tenant(tenant string) *NicBuilder {
	builder.metadata.Tenant(tenant)
	return builder
}

func (builder *NicBuilder) Workspace(workspace string) *NicBuilder {
	builder.metadata.Workspace(workspace)
	return builder
}

func (builder *NicBuilder) Region(region string) *NicBuilder {
	builder.metadata.Region(region)
	return builder
}

func (builder *NicBuilder) Spec(spec *schema.NicSpec) *NicBuilder {
	builder.spec = spec
	return builder
}

func (builder *NicBuilder) BuildResponse() (*schema.Nic, error) {
	metadata, err := builder.metadata.Kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindNic).BuildResponse()
	if err != nil {
		return nil, err
	}

	// Validate the spec
	if err := builder.validator.ValidateRequired(
		builder.spec,
		builder.spec.Addresses,
		builder.spec.SubnetRef,
	); err != nil {
		return nil, err
	}

	return &schema.Nic{
		Metadata: metadata,
		Spec:     *builder.spec,
		Status:   &schema.NicStatus{},
	}, nil
}

// Security Group

type SecurityGroupBuilder struct {
	*resourceBuilder
	metadata *RegionalWorkspaceResourceMetadataBuilder
	spec     *schema.SecurityGroupSpec
}

func NewSecurityGroupBuilder() *SecurityGroupBuilder {
	return &SecurityGroupBuilder{
		resourceBuilder: newResourceBuilder(),
		metadata:        NewRegionalWorkspaceResourceMetadataBuilder(),
		spec:            &schema.SecurityGroupSpec{},
	}
}

func (builder *SecurityGroupBuilder) Name(name string) *SecurityGroupBuilder {
	builder.metadata.Name(name)
	return builder
}

func (builder *SecurityGroupBuilder) Provider(provider string) *SecurityGroupBuilder {
	builder.metadata.Provider(provider)
	return builder
}

func (builder *SecurityGroupBuilder) Resource(resource string) *SecurityGroupBuilder {
	builder.metadata.Resource(resource)
	return builder
}

func (builder *SecurityGroupBuilder) ApiVersion(apiVersion string) *SecurityGroupBuilder {
	builder.metadata.ApiVersion(apiVersion)
	return builder
}

func (builder *SecurityGroupBuilder) Tenant(tenant string) *SecurityGroupBuilder {
	builder.metadata.Tenant(tenant)
	return builder
}

func (builder *SecurityGroupBuilder) Workspace(workspace string) *SecurityGroupBuilder {
	builder.metadata.Workspace(workspace)
	return builder
}

func (builder *SecurityGroupBuilder) Region(region string) *SecurityGroupBuilder {
	builder.metadata.Region(region)
	return builder
}

func (builder *SecurityGroupBuilder) Spec(spec *schema.SecurityGroupSpec) *SecurityGroupBuilder {
	builder.spec = spec
	return builder
}

func (builder *SecurityGroupBuilder) BuildResponse() (*schema.SecurityGroup, error) {
	metadata, err := builder.metadata.Kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup).BuildResponse()
	if err != nil {
		return nil, err
	}

	// Validate the spec
	if err := builder.validator.ValidateRequired(
		builder.spec,
		builder.spec.Rules,
	); err != nil {
		return nil, err
	}
	// Validate each rule
	for _, rule := range builder.spec.Rules {
		if err := builder.validator.ValidateRequired(
			rule.Direction,
		); err != nil {
			return nil, err
		}
	}

	return &schema.SecurityGroup{
		Metadata: metadata,
		Spec:     *builder.spec,
		Status:   &schema.SecurityGroupStatus{},
	}, nil
}
