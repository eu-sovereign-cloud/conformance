//nolint:dupl
package builders

import (
	"fmt"

	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	network "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.network.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// NetworkSku

// / NetworkSkuIteratorBuilder
type NetworkSkuIteratorBuilder struct {
	*tenantResponseMetadataBuilder[NetworkSkuIteratorBuilder]

	items []schema.NetworkSku
}

func NewNetworkSkuIteratorBuilder() *NetworkSkuIteratorBuilder {
	builder := &NetworkSkuIteratorBuilder{}
	builder.tenantResponseMetadataBuilder = newTenantResponseMetadataBuilder(builder)
	return builder
}

func (builder *NetworkSkuIteratorBuilder) Items(items []schema.NetworkSku) *NetworkSkuIteratorBuilder {
	builder.items = items
	return builder
}

func (builder *NetworkSkuIteratorBuilder) Build() (*network.SkuIterator, error) {
	err := builder.validate()
	if err != nil {
		return nil, err
	}

	builder.metadata.Resource = generators.GenerateSkuListResource(builder.tenant)

	return &network.SkuIterator{
		Metadata: *builder.metadata,
		Items:    builder.items,
	}, nil
}

// Network

/// NetworkMetadataBuilder

type NetworkMetadataBuilder struct {
	*regionalWorkspaceResourceMetadataBuilder[NetworkMetadataBuilder]
}

func NewNetworkMetadataBuilder() *NetworkMetadataBuilder {
	builder := &NetworkMetadataBuilder{}
	builder.regionalWorkspaceResourceMetadataBuilder = newRegionalWorkspaceResourceMetadataBuilder(builder)
	return builder
}

func (builder *NetworkMetadataBuilder) Build() (*schema.RegionalWorkspaceResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork).
		Resource(generators.GenerateNetworkResource(builder.metadata.Tenant, builder.metadata.Workspace, builder.metadata.Name)).
		Ref(generators.GenerateNetworkRef(builder.metadata.Name)).
		build()
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

/// NetworkBuilder

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
		field("spec", builder.spec),
		field("spec.Cidr", builder.spec.Cidr),
		field("spec.SkuRef", builder.spec.SkuRef),
		field("spec.RouteTableRef", builder.spec.RouteTableRef),
	); err != nil {
		return err
	}
	if err := validateOneRequired(builder.validator,
		field("spec.Cidr.Ipv4", builder.spec.Cidr.Ipv4),
		field("spec.Cidr.Ipv6", builder.spec.Cidr.Ipv6),
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

/// NetworkIteratorBuilder

type NetworkIteratorBuilder struct {
	*workspaceResponseMetadataBuilder[NetworkIteratorBuilder]

	items []schema.Network
}

func NewNetworkIteratorBuilder() *NetworkIteratorBuilder {
	builder := &NetworkIteratorBuilder{}
	builder.workspaceResponseMetadataBuilder = newWorkspaceResponseMetadataBuilder(builder)
	return builder
}

func (builder *NetworkIteratorBuilder) Items(items []schema.Network) *NetworkIteratorBuilder {
	builder.items = items
	return builder
}

func (builder *NetworkIteratorBuilder) Build() (*network.NetworkIterator, error) {
	err := builder.validate()
	if err != nil {
		return nil, err
	}

	builder.metadata.Resource = generators.GenerateNetworkListResource(builder.tenant, builder.workspace)

	return &network.NetworkIterator{
		Metadata: *builder.metadata,
		Items:    builder.items,
	}, nil
}

// InternetGateway

/// InternetGatewayMetadataBuilder

type InternetGatewayMetadataBuilder struct {
	*regionalWorkspaceResourceMetadataBuilder[InternetGatewayMetadataBuilder]
}

func NewInternetGatewayMetadataBuilder() *InternetGatewayMetadataBuilder {
	builder := &InternetGatewayMetadataBuilder{}
	builder.regionalWorkspaceResourceMetadataBuilder = newRegionalWorkspaceResourceMetadataBuilder(builder)
	return builder
}

func (builder *InternetGatewayMetadataBuilder) Build() (*schema.RegionalWorkspaceResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindInternetGateway).
		Resource(generators.GenerateInternetGatewayResource(builder.metadata.Tenant, builder.metadata.Workspace, builder.metadata.Name)).
		Ref(generators.GenerateInternetGatewayRef(builder.metadata.Name)).
		build()
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

/// InternetGatewayBuilder

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
		field("spec", builder.spec),
		field("spec.EgressOnly", builder.spec.EgressOnly),
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

/// InternetGatewayIteratorBuilder

type InternetGatewayIteratorBuilder struct {
	*workspaceResponseMetadataBuilder[InternetGatewayIteratorBuilder]

	items []schema.InternetGateway
}

func NewInternetGatewayIteratorBuilder() *InternetGatewayIteratorBuilder {
	builder := &InternetGatewayIteratorBuilder{}
	builder.workspaceResponseMetadataBuilder = newWorkspaceResponseMetadataBuilder(builder)
	return builder
}

func (builder *InternetGatewayIteratorBuilder) Items(items []schema.InternetGateway) *InternetGatewayIteratorBuilder {
	builder.items = items
	return builder
}

func (builder *InternetGatewayIteratorBuilder) Build() (*network.InternetGatewayIterator, error) {
	err := builder.validate()
	if err != nil {
		return nil, err
	}

	builder.metadata.Resource = generators.GenerateInternetGatewayListResource(builder.tenant, builder.workspace)

	return &network.InternetGatewayIterator{
		Metadata: *builder.metadata,
		Items:    builder.items,
	}, nil
}

// RouteTable

/// RouteTableMetadataBuilder

type RouteTableMetadataBuilder struct {
	*regionalNetworkResourceMetadataBuilder[RouteTableMetadataBuilder]
}

func NewRouteTableMetadataBuilder() *RouteTableMetadataBuilder {
	builder := &RouteTableMetadataBuilder{}
	builder.regionalNetworkResourceMetadataBuilder = newRegionalNetworkResourceMetadataBuilder(builder)
	return builder
}

func (builder *RouteTableMetadataBuilder) Build() (*schema.RegionalNetworkResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable).
		Resource(generators.GenerateRouteTableResource(builder.metadata.Tenant, builder.metadata.Workspace, builder.metadata.Network, builder.metadata.Name)).
		Ref(generators.GenerateRouteTableRef(builder.metadata.Name)).
		build()
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

/// RouteTableBuilder

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
		field("spec.Routes", builder.spec.Routes),
	); err != nil {
		return err
	}

	// Validate each route
	for i, route := range builder.spec.Routes {
		if err := validateRequired(builder.validator,
			field(fmt.Sprintf("spec.Routes[%d].DestinationCidrBlock", i), route.DestinationCidrBlock),
			field(fmt.Sprintf("spec.Routes[%d].TargetRef", i), route.TargetRef),
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

/// RouteTableIteratorBuilder

type RouteTableIteratorBuilder struct {
	*networkResponseMetadataBuilder[RouteTableIteratorBuilder]

	items []schema.RouteTable
}

func NewRouteTableIteratorBuilder() *RouteTableIteratorBuilder {
	builder := &RouteTableIteratorBuilder{}
	builder.networkResponseMetadataBuilder = newNetworkResponseMetadataBuilder(builder)
	return builder
}

func (builder *RouteTableIteratorBuilder) Items(items []schema.RouteTable) *RouteTableIteratorBuilder {
	builder.items = items
	return builder
}

func (builder *RouteTableIteratorBuilder) Build() (*network.RouteTableIterator, error) {
	err := builder.validate()
	if err != nil {
		return nil, err
	}

	builder.metadata.Resource = generators.GenerateRouteTableListResource(builder.tenant, builder.workspace, builder.network)

	return &network.RouteTableIterator{
		Metadata: *builder.metadata,
		Items:    builder.items,
	}, nil
}

// Subnet

/// SubnetMetadataBuilder

type SubnetMetadataBuilder struct {
	*regionalNetworkResourceMetadataBuilder[SubnetMetadataBuilder]
}

func NewSubnetMetadataBuilder() *SubnetMetadataBuilder {
	builder := &SubnetMetadataBuilder{}
	builder.regionalNetworkResourceMetadataBuilder = newRegionalNetworkResourceMetadataBuilder(builder)
	return builder
}

func (builder *SubnetMetadataBuilder) Build() (*schema.RegionalNetworkResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet).
		Resource(generators.GenerateSubnetResource(builder.metadata.Tenant, builder.metadata.Workspace, builder.metadata.Network, builder.metadata.Name)).
		Ref(generators.GenerateSubnetRef(builder.metadata.Name)).
		build()
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

/// SubnetBuilder

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
		field("spec", builder.spec),
		field("spec.Cidr", builder.spec.Cidr),
		field("spec.Zone", builder.spec.Zone),
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

/// SubnetIteratorBuilder

type SubnetIteratorBuilder struct {
	*networkResponseMetadataBuilder[SubnetIteratorBuilder]

	items []schema.Subnet
}

func NewSubnetIteratorBuilder() *SubnetIteratorBuilder {
	builder := &SubnetIteratorBuilder{}
	builder.networkResponseMetadataBuilder = newNetworkResponseMetadataBuilder(builder)
	return builder
}

func (builder *SubnetIteratorBuilder) Items(items []schema.Subnet) *SubnetIteratorBuilder {
	builder.items = items
	return builder
}

func (builder *SubnetIteratorBuilder) Build() (*network.SubnetIterator, error) {
	err := builder.validate()
	if err != nil {
		return nil, err
	}

	builder.metadata.Resource = generators.GenerateSubnetListResource(builder.tenant, builder.workspace, builder.network)

	return &network.SubnetIterator{
		Metadata: *builder.metadata,
		Items:    builder.items,
	}, nil
}

// PublicIp

/// PublicIpMetadataBuilder

type PublicIpMetadataBuilder struct {
	*regionalWorkspaceResourceMetadataBuilder[PublicIpMetadataBuilder]
}

func NewPublicIpMetadataBuilder() *PublicIpMetadataBuilder {
	builder := &PublicIpMetadataBuilder{}
	builder.regionalWorkspaceResourceMetadataBuilder = newRegionalWorkspaceResourceMetadataBuilder(builder)
	return builder
}

func (builder *PublicIpMetadataBuilder) Build() (*schema.RegionalWorkspaceResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindPublicIP).
		Resource(generators.GeneratePublicIpResource(builder.metadata.Tenant, builder.metadata.Workspace, builder.metadata.Name)).
		Ref(generators.GeneratePublicIpRef(builder.metadata.Name)).
		build()
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

/// PublicIpBuilder

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
		field("spec", builder.spec),
		field("spec.Version", builder.spec.Version),
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

/// PublicIpIteratorBuilder

type PublicIpIteratorBuilder struct {
	*workspaceResponseMetadataBuilder[PublicIpIteratorBuilder]

	items []schema.PublicIp
}

func NewPublicIpIteratorBuilder() *PublicIpIteratorBuilder {
	builder := &PublicIpIteratorBuilder{}
	builder.workspaceResponseMetadataBuilder = newWorkspaceResponseMetadataBuilder(builder)
	return builder
}

func (builder *PublicIpIteratorBuilder) Items(items []schema.PublicIp) *PublicIpIteratorBuilder {
	builder.items = items
	return builder
}

func (builder *PublicIpIteratorBuilder) Build() (*network.PublicIpIterator, error) {
	err := builder.validate()
	if err != nil {
		return nil, err
	}

	builder.metadata.Resource = generators.GeneratePublicIpListResource(builder.tenant, builder.workspace)

	return &network.PublicIpIterator{
		Metadata: *builder.metadata,
		Items:    builder.items,
	}, nil
}

// Nic

/// NicMetadataBuilder

type NicMetadataBuilder struct {
	*regionalWorkspaceResourceMetadataBuilder[NicMetadataBuilder]
}

func NewNicMetadataBuilder() *NicMetadataBuilder {
	builder := &NicMetadataBuilder{}
	builder.regionalWorkspaceResourceMetadataBuilder = newRegionalWorkspaceResourceMetadataBuilder(builder)
	return builder
}

func (builder *NicMetadataBuilder) Build() (*schema.RegionalWorkspaceResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindNic).
		Resource(generators.GenerateNicResource(builder.metadata.Tenant, builder.metadata.Workspace, builder.metadata.Name)).
		Ref(generators.GenerateNicRef(builder.metadata.Name)).
		build()
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

/// NicBuilder

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
		field("spec", builder.spec),
		field("spec.Addresses", builder.spec.Addresses),
		field("spec.SubnetRef", builder.spec.SubnetRef),
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

/// NicIteratorBuilder

type NicIteratorBuilder struct {
	*workspaceResponseMetadataBuilder[NicIteratorBuilder]

	items []schema.Nic
}

func NewNicIteratorBuilder() *NicIteratorBuilder {
	builder := &NicIteratorBuilder{}
	builder.workspaceResponseMetadataBuilder = newWorkspaceResponseMetadataBuilder(builder)
	return builder
}

func (builder *NicIteratorBuilder) Items(items []schema.Nic) *NicIteratorBuilder {
	builder.items = items
	return builder
}

func (builder *NicIteratorBuilder) Build() (*network.NicIterator, error) {
	err := builder.validate()
	if err != nil {
		return nil, err
	}

	builder.metadata.Resource = generators.GenerateNicListResource(builder.tenant, builder.workspace)

	return &network.NicIterator{
		Metadata: *builder.metadata,
		Items:    builder.items,
	}, nil
}

// SecurityGroupRule

/// SecurityGroupRuleMetadataBuilder

type SecurityGroupRuleMetadataBuilder struct {
	*regionalWorkspaceResourceMetadataBuilder[SecurityGroupRuleMetadataBuilder]
}

func NewSecurityGroupRuleMetadataBuilder() *SecurityGroupRuleMetadataBuilder {
	builder := &SecurityGroupRuleMetadataBuilder{}
	builder.regionalWorkspaceResourceMetadataBuilder = newRegionalWorkspaceResourceMetadataBuilder(builder)
	return builder
}

func (builder *SecurityGroupRuleMetadataBuilder) Build() (*schema.RegionalWorkspaceResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroupRule).
		Resource(generators.GenerateSecurityGroupRuleResource(builder.metadata.Tenant, builder.metadata.Workspace, builder.metadata.Name)).
		Ref(generators.GenerateSecurityGroupRuleRef(builder.metadata.Name)).
		build()
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

/// SecurityGroupRuleBuilder

type SecurityGroupRuleBuilder struct {
	*regionalWorkspaceResourceBuilder[SecurityGroupRuleBuilder, schema.SecurityGroupRuleSpec]
	metadata *SecurityGroupRuleMetadataBuilder
	labels   schema.Labels
	spec     *schema.SecurityGroupRuleSpec
}

func NewSecurityGroupRuleBuilder() *SecurityGroupRuleBuilder {
	builder := &SecurityGroupRuleBuilder{
		metadata: NewSecurityGroupRuleMetadataBuilder(),
		spec:     &schema.SecurityGroupRuleSpec{},
	}

	builder.regionalWorkspaceResourceBuilder = newRegionalWorkspaceResourceBuilder(newRegionalWorkspaceResourceBuilderParams[SecurityGroupRuleBuilder, schema.SecurityGroupRuleSpec]{
		newGlobalResourceBuilderParams: &newGlobalResourceBuilderParams[SecurityGroupRuleBuilder, schema.SecurityGroupRuleSpec]{
			parent:        builder,
			setName:       func(name string) { builder.metadata.setName(name) },
			setProvider:   func(provider string) { builder.metadata.setProvider(provider) },
			setApiVersion: func(apiVersion string) { builder.metadata.setApiVersion(apiVersion) },
			setLabels:     func(labels schema.Labels) { builder.labels = labels },
			setSpec:       func(spec *schema.SecurityGroupRuleSpec) { builder.spec = spec },
		},
		setTenant:    func(tenant string) { builder.metadata.Tenant(tenant) },
		setWorkspace: func(workspace string) { builder.metadata.Workspace(workspace) },
		setRegion:    func(region string) { builder.metadata.Region(region) },
	})

	return builder
}

func (builder *SecurityGroupRuleBuilder) validateSpec() error {
	if err := validateRequired(builder.validator,
		field("spec", builder.spec),
		field("spec.Direction", builder.spec.Direction),
	); err != nil {
		return err
	}

	return nil
}

func (builder *SecurityGroupRuleBuilder) Build() (*schema.SecurityGroupRule, error) {
	if err := builder.validateSpec(); err != nil {
		return nil, err
	}

	metadata, err := builder.metadata.Build()
	if err != nil {
		return nil, err
	}

	return &schema.SecurityGroupRule{
		Metadata: metadata,
		Labels:   builder.labels,
		Spec:     *builder.spec,
		Status:   &schema.Status{},
	}, nil
}

// SecurityGroup

/// SecurityGroupMetadataBuilder

type SecurityGroupMetadataBuilder struct {
	*regionalWorkspaceResourceMetadataBuilder[SecurityGroupMetadataBuilder]
}

func NewSecurityGroupMetadataBuilder() *SecurityGroupMetadataBuilder {
	builder := &SecurityGroupMetadataBuilder{}
	builder.regionalWorkspaceResourceMetadataBuilder = newRegionalWorkspaceResourceMetadataBuilder(builder)
	return builder
}

func (builder *SecurityGroupMetadataBuilder) Build() (*schema.RegionalWorkspaceResourceMetadata, error) {
	metadata, err := builder.kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup).
		Resource(generators.GenerateSecurityGroupResource(builder.metadata.Tenant, builder.metadata.Workspace, builder.metadata.Name)).
		Ref(generators.GenerateSecurityGroupRef(builder.metadata.Name)).
		build()
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

/// SecurityGroupBuilder

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
		field("spec", builder.spec),
		field("spec.Rules", builder.spec.Rules),
	); err != nil {
		return err
	}

	// Validate each rule
	for i, rule := range *builder.spec.Rules {
		if err := validateRequired(builder.validator,
			field(fmt.Sprintf("spec.Rules[%d].Direction", i), rule.Direction),
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

/// SecurityGroupIteratorBuilder

type SecurityGroupIteratorBuilder struct {
	*workspaceResponseMetadataBuilder[SecurityGroupIteratorBuilder]

	items []schema.SecurityGroup
}

func NewSecurityGroupIteratorBuilder() *SecurityGroupIteratorBuilder {
	builder := &SecurityGroupIteratorBuilder{}
	builder.workspaceResponseMetadataBuilder = newWorkspaceResponseMetadataBuilder(builder)
	return builder
}

func (builder *SecurityGroupIteratorBuilder) Items(items []schema.SecurityGroup) *SecurityGroupIteratorBuilder {
	builder.items = items
	return builder
}

func (builder *SecurityGroupIteratorBuilder) Build() (*network.SecurityGroupIterator, error) {
	err := builder.validate()
	if err != nil {
		return nil, err
	}

	builder.metadata.Resource = generators.GenerateSecurityGroupListResource(builder.tenant, builder.workspace)

	return &network.SecurityGroupIterator{
		Metadata: *builder.metadata,
		Items:    builder.items,
	}, nil
}
