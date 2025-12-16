//nolint:dupl
package builders

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Network

type NetworkMetadataBuilder struct {
	*regionalWorkspaceResourceMetadataBuilder[NetworkMetadataBuilder]
}

func NewNetworkMetadataBuilder() *NetworkMetadataBuilder {
	builder := &NetworkMetadataBuilder{}
	builder.regionalWorkspaceResourceMetadataBuilder = newRegionalWorkspaceResourceMetadataBuilder(builder)
	return builder
}

func (builder *NetworkMetadataBuilder) Build() (*schema.RegionalWorkspaceResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork).build()
	if err != nil {
		return nil, err
	}

	resource := generators.GenerateNetworkResource(builder.metadata.Tenant, builder.metadata.Workspace, builder.metadata.Name)
	metadata.Resource = resource

	return metadata, nil
}

type NetworkBuilder struct {
	*regionalWorkspaceResourceBuilder[NetworkBuilder, schema.NetworkSpec]
	metadata *NetworkMetadataBuilder
	labels   schema.Labels
	spec     *schema.NetworkSpec
}

func NewNetworkBuilder() *NetworkBuilder {
	builder := &NetworkBuilder{
		metadata: NewNetworkMetadataBuilder(),
		spec:     &schema.NetworkSpec{},
	}

	builder.regionalWorkspaceResourceBuilder = newRegionalWorkspaceResourceBuilder(newRegionalWorkspaceResourceBuilderParams[NetworkBuilder, schema.NetworkSpec]{
		newGlobalResourceBuilderParams: &newGlobalResourceBuilderParams[NetworkBuilder, schema.NetworkSpec]{
			parent:        builder,
			setName:       func(name string) { builder.metadata.setName(name) },
			setProvider:   func(provider string) { builder.metadata.setProvider(provider) },
			setApiVersion: func(apiVersion string) { builder.metadata.setApiVersion(apiVersion) },
			setLabels:     func(labels schema.Labels) { builder.labels = labels },
			setSpec:       func(spec *schema.NetworkSpec) { builder.spec = spec },
		},
		setTenant:    func(tenant string) { builder.metadata.Tenant(tenant) },
		setWorkspace: func(workspace string) { builder.metadata.Workspace(workspace) },
		setRegion:    func(region string) { builder.metadata.Region(region) },
	})

	return builder
}

func (builder *NetworkBuilder) validateSpec() error {
	if err := validateRequired(builder.validator,
		builder.spec,
		builder.spec.Cidr,
		builder.spec.SkuRef,
		builder.spec.RouteTableRef,
	); err != nil {
		return err
	}
	if err := validateOneRequired(builder.validator,
		builder.spec.Cidr.Ipv4,
		builder.spec.Cidr.Ipv6,
	); err != nil {
		return err
	}

	return nil
}

func (builder *NetworkBuilder) Build() (*schema.Network, error) {
	if err := builder.validateSpec(); err != nil {
		return nil, err
	}

	metadata, err := builder.metadata.Build()
	if err != nil {
		return nil, err
	}

	return &schema.Network{
		Metadata: metadata,
		Labels:   builder.labels,
		Spec:     *builder.spec,
		Status:   &schema.NetworkStatus{},
	}, nil
}

// Internet Gateway

type InternetGatewayMetadataBuilder struct {
	*regionalWorkspaceResourceMetadataBuilder[InternetGatewayMetadataBuilder]
}

func NewInternetGatewayMetadataBuilder() *InternetGatewayMetadataBuilder {
	builder := &InternetGatewayMetadataBuilder{}
	builder.regionalWorkspaceResourceMetadataBuilder = newRegionalWorkspaceResourceMetadataBuilder(builder)
	return builder
}

func (builder *InternetGatewayMetadataBuilder) Build() (*schema.RegionalWorkspaceResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindInternetGateway).build()
	if err != nil {
		return nil, err
	}

	resource := generators.GenerateInternetGatewayResource(builder.metadata.Tenant, builder.metadata.Workspace, builder.metadata.Name)
	metadata.Resource = resource

	return metadata, nil
}

type InternetGatewayBuilder struct {
	*regionalWorkspaceResourceBuilder[InternetGatewayBuilder, schema.InternetGatewaySpec]
	metadata *InternetGatewayMetadataBuilder
	labels   schema.Labels
	spec     *schema.InternetGatewaySpec
}

func NewInternetGatewayBuilder() *InternetGatewayBuilder {
	builder := &InternetGatewayBuilder{
		metadata: NewInternetGatewayMetadataBuilder(),
		spec:     &schema.InternetGatewaySpec{},
	}

	builder.regionalWorkspaceResourceBuilder = newRegionalWorkspaceResourceBuilder(newRegionalWorkspaceResourceBuilderParams[InternetGatewayBuilder, schema.InternetGatewaySpec]{
		newGlobalResourceBuilderParams: &newGlobalResourceBuilderParams[InternetGatewayBuilder, schema.InternetGatewaySpec]{
			parent:        builder,
			setName:       func(name string) { builder.metadata.setName(name) },
			setProvider:   func(provider string) { builder.metadata.setProvider(provider) },
			setApiVersion: func(apiVersion string) { builder.metadata.setApiVersion(apiVersion) },
			setLabels:     func(labels schema.Labels) { builder.labels = labels },
			setSpec:       func(spec *schema.InternetGatewaySpec) { builder.spec = spec },
		},
		setTenant:    func(tenant string) { builder.metadata.Tenant(tenant) },
		setWorkspace: func(workspace string) { builder.metadata.Workspace(workspace) },
		setRegion:    func(region string) { builder.metadata.Region(region) },
	})

	return builder
}

func (builder *InternetGatewayBuilder) validateSpec() error {
	if err := validateRequired(builder.validator,
		builder.spec,
		builder.spec.EgressOnly,
	); err != nil {
		return err
	}

	return nil
}

func (builder *InternetGatewayBuilder) Build() (*schema.InternetGateway, error) {
	if err := builder.validateSpec(); err != nil {
		return nil, err
	}

	metadata, err := builder.metadata.Build()
	if err != nil {
		return nil, err
	}

	return &schema.InternetGateway{
		Metadata: metadata,
		Labels:   builder.labels,
		Spec:     *builder.spec,
		Status:   &schema.InternetGatewayStatus{},
	}, nil
}

// Route Table

type RouteTableMetadataBuilder struct {
	*regionalNetworkResourceMetadataBuilder[RouteTableMetadataBuilder]
}

func NewRouteTableMetadataBuilder() *RouteTableMetadataBuilder {
	builder := &RouteTableMetadataBuilder{}
	builder.regionalNetworkResourceMetadataBuilder = newRegionalNetworkResourceMetadataBuilder(builder)
	return builder
}

func (builder *RouteTableMetadataBuilder) Build() (*schema.RegionalNetworkResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable).build()
	if err != nil {
		return nil, err
	}

	resource := generators.GenerateRouteTableResource(builder.metadata.Tenant, builder.metadata.Workspace, builder.metadata.Network, builder.metadata.Name)
	metadata.Resource = resource

	return metadata, nil
}

type RouteTableBuilder struct {
	*regionalNetworkResourceBuilder[RouteTableBuilder, schema.RouteTableSpec]
	metadata *RouteTableMetadataBuilder
	labels   schema.Labels
	spec     *schema.RouteTableSpec
}

func NewRouteTableBuilder() *RouteTableBuilder {
	builder := &RouteTableBuilder{
		metadata: NewRouteTableMetadataBuilder(),
		spec:     &schema.RouteTableSpec{},
	}

	builder.regionalNetworkResourceBuilder = newRegionalNetworkResourceBuilder(newRegionalNetworkResourceBuilderParams[RouteTableBuilder, schema.RouteTableSpec]{
		newGlobalResourceBuilderParams: &newGlobalResourceBuilderParams[RouteTableBuilder, schema.RouteTableSpec]{
			parent:        builder,
			setName:       func(name string) { builder.metadata.setName(name) },
			setProvider:   func(provider string) { builder.metadata.setProvider(provider) },
			setApiVersion: func(apiVersion string) { builder.metadata.setApiVersion(apiVersion) },
			setLabels:     func(labels schema.Labels) { builder.labels = labels },
			setSpec:       func(spec *schema.RouteTableSpec) { builder.spec = spec },
		},
		setTenant:    func(tenant string) { builder.metadata.Tenant(tenant) },
		setWorkspace: func(workspace string) { builder.metadata.Workspace(workspace) },
		setNetwork:   func(network string) { builder.metadata.Network(network) },
		setRegion:    func(region string) { builder.metadata.Region(region) },
	})

	return builder
}

func (builder *RouteTableBuilder) validateSpec() error {
	if err := validateRequired(builder.validator,
		builder.spec.Routes,
	); err != nil {
		return err
	}

	// Validate each route
	for _, route := range builder.spec.Routes {
		if err := validateRequired(builder.validator,
			route.DestinationCidrBlock,
			route.TargetRef,
		); err != nil {
			return err
		}
	}

	return nil
}

func (builder *RouteTableBuilder) Build() (*schema.RouteTable, error) {
	if err := builder.validateSpec(); err != nil {
		return nil, err
	}

	metadata, err := builder.metadata.Build()
	if err != nil {
		return nil, err
	}

	return &schema.RouteTable{
		Metadata: metadata,
		Labels:   builder.labels,
		Spec:     *builder.spec,
		Status:   &schema.RouteTableStatus{},
	}, nil
}

// Subnet

type SubnetMetadataBuilder struct {
	*regionalNetworkResourceMetadataBuilder[SubnetMetadataBuilder]
}

func NewSubnetMetadataBuilder() *SubnetMetadataBuilder {
	builder := &SubnetMetadataBuilder{}
	builder.regionalNetworkResourceMetadataBuilder = newRegionalNetworkResourceMetadataBuilder(builder)
	return builder
}

func (builder *SubnetMetadataBuilder) Build() (*schema.RegionalNetworkResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet).build()
	if err != nil {
		return nil, err
	}

	resource := generators.GenerateSubnetResource(builder.metadata.Tenant, builder.metadata.Workspace, builder.metadata.Network, builder.metadata.Name)
	metadata.Resource = resource

	return metadata, nil
}

type SubnetBuilder struct {
	*regionalNetworkResourceBuilder[SubnetBuilder, schema.SubnetSpec]
	metadata *SubnetMetadataBuilder
	labels   schema.Labels
	spec     *schema.SubnetSpec
}

func NewSubnetBuilder() *SubnetBuilder {
	builder := &SubnetBuilder{
		metadata: NewSubnetMetadataBuilder(),
		spec:     &schema.SubnetSpec{},
	}

	builder.regionalNetworkResourceBuilder = newRegionalNetworkResourceBuilder(newRegionalNetworkResourceBuilderParams[SubnetBuilder, schema.SubnetSpec]{
		newGlobalResourceBuilderParams: &newGlobalResourceBuilderParams[SubnetBuilder, schema.SubnetSpec]{
			parent:        builder,
			setName:       func(name string) { builder.metadata.setName(name) },
			setProvider:   func(provider string) { builder.metadata.setProvider(provider) },
			setApiVersion: func(apiVersion string) { builder.metadata.setApiVersion(apiVersion) },
			setLabels:     func(labels schema.Labels) { builder.labels = labels },
			setSpec:       func(spec *schema.SubnetSpec) { builder.spec = spec },
		},
		setTenant:    func(tenant string) { builder.metadata.Tenant(tenant) },
		setWorkspace: func(workspace string) { builder.metadata.Workspace(workspace) },
		setNetwork:   func(network string) { builder.metadata.Network(network) },
		setRegion:    func(region string) { builder.metadata.Region(region) },
	})

	return builder
}

func (builder *SubnetBuilder) validateSpec() error {
	if err := validateRequired(builder.validator,
		builder.spec,
		builder.spec.Cidr,
		builder.spec.Zone,
	); err != nil {
		return err
	}

	return nil
}

func (builder *SubnetBuilder) Build() (*schema.Subnet, error) {
	if err := builder.validateSpec(); err != nil {
		return nil, err
	}

	metadata, err := builder.metadata.Build()
	if err != nil {
		return nil, err
	}

	return &schema.Subnet{
		Metadata: metadata,
		Labels:   builder.labels,
		Spec:     *builder.spec,
		Status:   &schema.SubnetStatus{},
	}, nil
}

// Public Ip

type PublicIpMetadataBuilder struct {
	*regionalWorkspaceResourceMetadataBuilder[PublicIpMetadataBuilder]
}

func NewPublicIpMetadataBuilder() *PublicIpMetadataBuilder {
	builder := &PublicIpMetadataBuilder{}
	builder.regionalWorkspaceResourceMetadataBuilder = newRegionalWorkspaceResourceMetadataBuilder(builder)
	return builder
}

func (builder *PublicIpMetadataBuilder) Build() (*schema.RegionalWorkspaceResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindPublicIP).build()
	if err != nil {
		return nil, err
	}

	resource := generators.GeneratePublicIpResource(builder.metadata.Tenant, builder.metadata.Workspace, builder.metadata.Name)
	metadata.Resource = resource

	return metadata, nil
}

type PublicIpBuilder struct {
	*regionalWorkspaceResourceBuilder[PublicIpBuilder, schema.PublicIpSpec]
	metadata *PublicIpMetadataBuilder
	labels   schema.Labels
	spec     *schema.PublicIpSpec
}

func NewPublicIpBuilder() *PublicIpBuilder {
	builder := &PublicIpBuilder{
		metadata: NewPublicIpMetadataBuilder(),
		spec:     &schema.PublicIpSpec{},
	}

	builder.regionalWorkspaceResourceBuilder = newRegionalWorkspaceResourceBuilder(newRegionalWorkspaceResourceBuilderParams[PublicIpBuilder, schema.PublicIpSpec]{
		newGlobalResourceBuilderParams: &newGlobalResourceBuilderParams[PublicIpBuilder, schema.PublicIpSpec]{
			parent:        builder,
			setName:       func(name string) { builder.metadata.setName(name) },
			setProvider:   func(provider string) { builder.metadata.setProvider(provider) },
			setApiVersion: func(apiVersion string) { builder.metadata.setApiVersion(apiVersion) },
			setLabels:     func(labels schema.Labels) { builder.labels = labels },
			setSpec:       func(spec *schema.PublicIpSpec) { builder.spec = spec },
		},
		setTenant:    func(tenant string) { builder.metadata.Tenant(tenant) },
		setWorkspace: func(workspace string) { builder.metadata.Workspace(workspace) },
		setRegion:    func(region string) { builder.metadata.Region(region) },
	})

	return builder
}

func (builder *PublicIpBuilder) validateSpec() error {
	if err := validateRequired(builder.validator,
		builder.spec,
		builder.spec.Version,
	); err != nil {
		return err
	}

	return nil
}

func (builder *PublicIpBuilder) Build() (*schema.PublicIp, error) {
	if err := builder.validateSpec(); err != nil {
		return nil, err
	}

	metadata, err := builder.metadata.Build()
	if err != nil {
		return nil, err
	}

	return &schema.PublicIp{
		Metadata: metadata,
		Labels:   builder.labels,
		Spec:     *builder.spec,
		Status:   &schema.PublicIpStatus{},
	}, nil
}

// Nic

type NicMetadataBuilder struct {
	*regionalWorkspaceResourceMetadataBuilder[NicMetadataBuilder]
}

func NewNicMetadataBuilder() *NicMetadataBuilder {
	builder := &NicMetadataBuilder{}
	builder.regionalWorkspaceResourceMetadataBuilder = newRegionalWorkspaceResourceMetadataBuilder(builder)
	return builder
}

func (builder *NicMetadataBuilder) Build() (*schema.RegionalWorkspaceResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindNic).build()
	if err != nil {
		return nil, err
	}

	resource := generators.GenerateNicResource(builder.metadata.Tenant, builder.metadata.Workspace, builder.metadata.Name)
	metadata.Resource = resource

	return metadata, nil
}

type NicBuilder struct {
	*regionalWorkspaceResourceBuilder[NicBuilder, schema.NicSpec]
	metadata *NicMetadataBuilder
	labels   schema.Labels
	spec     *schema.NicSpec
}

func NewNicBuilder() *NicBuilder {
	builder := &NicBuilder{
		metadata: NewNicMetadataBuilder(),
		spec:     &schema.NicSpec{},
	}

	builder.regionalWorkspaceResourceBuilder = newRegionalWorkspaceResourceBuilder(newRegionalWorkspaceResourceBuilderParams[NicBuilder, schema.NicSpec]{
		newGlobalResourceBuilderParams: &newGlobalResourceBuilderParams[NicBuilder, schema.NicSpec]{
			parent:        builder,
			setName:       func(name string) { builder.metadata.setName(name) },
			setProvider:   func(provider string) { builder.metadata.setProvider(provider) },
			setApiVersion: func(apiVersion string) { builder.metadata.setApiVersion(apiVersion) },
			setLabels:     func(labels schema.Labels) { builder.labels = labels },
			setSpec:       func(spec *schema.NicSpec) { builder.spec = spec },
		},
		setTenant:    func(tenant string) { builder.metadata.Tenant(tenant) },
		setWorkspace: func(workspace string) { builder.metadata.Workspace(workspace) },
		setRegion:    func(region string) { builder.metadata.Region(region) },
	})

	return builder
}

func (builder *NicBuilder) validateSpec() error {
	if err := validateRequired(builder.validator,
		builder.spec,
		builder.spec.Addresses,
		builder.spec.SubnetRef,
	); err != nil {
		return err
	}

	return nil
}

func (builder *NicBuilder) Build() (*schema.Nic, error) {
	if err := builder.validateSpec(); err != nil {
		return nil, err
	}

	metadata, err := builder.metadata.Build()
	if err != nil {
		return nil, err
	}

	return &schema.Nic{
		Metadata: metadata,
		Labels:   builder.labels,
		Spec:     *builder.spec,
		Status:   &schema.NicStatus{},
	}, nil
}

// Security Group

type SecurityGroupMetadataBuilder struct {
	*regionalWorkspaceResourceMetadataBuilder[SecurityGroupMetadataBuilder]
}

func NewSecurityGroupMetadataBuilder() *SecurityGroupMetadataBuilder {
	builder := &SecurityGroupMetadataBuilder{}
	builder.regionalWorkspaceResourceMetadataBuilder = newRegionalWorkspaceResourceMetadataBuilder(builder)
	return builder
}

func (builder *SecurityGroupMetadataBuilder) Build() (*schema.RegionalWorkspaceResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup).build()
	if err != nil {
		return nil, err
	}

	resource := generators.GenerateSecurityGroupResource(builder.metadata.Tenant, builder.metadata.Workspace, builder.metadata.Name)
	metadata.Resource = resource

	return metadata, nil
}

type SecurityGroupBuilder struct {
	*regionalWorkspaceResourceBuilder[SecurityGroupBuilder, schema.SecurityGroupSpec]
	metadata *SecurityGroupMetadataBuilder
	labels   schema.Labels
	spec     *schema.SecurityGroupSpec
}

func NewSecurityGroupBuilder() *SecurityGroupBuilder {
	builder := &SecurityGroupBuilder{
		metadata: NewSecurityGroupMetadataBuilder(),
		spec:     &schema.SecurityGroupSpec{},
	}

	builder.regionalWorkspaceResourceBuilder = newRegionalWorkspaceResourceBuilder(newRegionalWorkspaceResourceBuilderParams[SecurityGroupBuilder, schema.SecurityGroupSpec]{
		newGlobalResourceBuilderParams: &newGlobalResourceBuilderParams[SecurityGroupBuilder, schema.SecurityGroupSpec]{
			parent:        builder,
			setName:       func(name string) { builder.metadata.setName(name) },
			setProvider:   func(provider string) { builder.metadata.setProvider(provider) },
			setApiVersion: func(apiVersion string) { builder.metadata.setApiVersion(apiVersion) },
			setLabels:     func(labels schema.Labels) { builder.labels = labels },
			setSpec:       func(spec *schema.SecurityGroupSpec) { builder.spec = spec },
		},
		setTenant:    func(tenant string) { builder.metadata.Tenant(tenant) },
		setWorkspace: func(workspace string) { builder.metadata.Workspace(workspace) },
		setRegion:    func(region string) { builder.metadata.Region(region) },
	})

	return builder
}

func (builder *SecurityGroupBuilder) validateSpec() error {
	if err := validateRequired(builder.validator,
		builder.spec,
		builder.spec.Rules,
	); err != nil {
		return err
	}

	// Validate each rule
	for _, rule := range builder.spec.Rules {
		if err := validateRequired(builder.validator,
			rule.Direction,
		); err != nil {
			return err
		}
	}

	return nil
}

func (builder *SecurityGroupBuilder) Build() (*schema.SecurityGroup, error) {
	if err := builder.validateSpec(); err != nil {
		return nil, err
	}

	metadata, err := builder.metadata.Build()
	if err != nil {
		return nil, err
	}

	return &schema.SecurityGroup{
		Metadata: metadata,
		Labels:   builder.labels,
		Spec:     *builder.spec,
		Status:   &schema.SecurityGroupStatus{},
	}, nil
}
