package wrappers

import "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

type NetworkWrapper struct {
	*BaseResourceWrapper[schema.Network]
}

func NewNetworkWrapper(resource *schema.Network) *NetworkWrapper {
	return &NetworkWrapper{
		BaseResourceWrapper: newBaseResourceWrapper(resource),
	}
}

func (wrapper *NetworkWrapper) GetResource() *schema.Network {
	return wrapper.resource
}

func (wrapper *NetworkWrapper) GetLabels() schema.Labels {
	return wrapper.resource.Labels
}

func (wrapper *NetworkWrapper) GetAnnotations() schema.Annotations {
	return wrapper.resource.Annotations
}

func (wrapper *NetworkWrapper) GetExtensions() schema.Extensions {
	return wrapper.resource.Extensions
}

func (wrapper *NetworkWrapper) GetMetadata() *schema.RegionalWorkspaceResourceMetadata {
	return wrapper.resource.Metadata
}

func (wrapper *NetworkWrapper) GetSpec() *schema.NetworkSpec {
	return &wrapper.resource.Spec
}

func (wrapper *NetworkWrapper) GetStatus() *schema.NetworkStatus {
	return wrapper.resource.Status
}

type InternetGatewayWrapper struct {
	*BaseResourceWrapper[schema.InternetGateway]
}

func NewInternetGatewayWrapper(resource *schema.InternetGateway) *InternetGatewayWrapper {
	return &InternetGatewayWrapper{
		BaseResourceWrapper: newBaseResourceWrapper(resource),
	}
}

func (wrapper *InternetGatewayWrapper) GetResource() *schema.InternetGateway {
	return wrapper.resource
}

func (wrapper *InternetGatewayWrapper) GetLabels() schema.Labels {
	return wrapper.resource.Labels
}

func (wrapper *InternetGatewayWrapper) GetAnnotations() schema.Annotations {
	return wrapper.resource.Annotations
}

func (wrapper *InternetGatewayWrapper) GetExtensions() schema.Extensions {
	return wrapper.resource.Extensions
}

func (wrapper *InternetGatewayWrapper) GetMetadata() *schema.RegionalWorkspaceResourceMetadata {
	return wrapper.resource.Metadata
}

func (wrapper *InternetGatewayWrapper) GetSpec() *schema.InternetGatewaySpec {
	return &wrapper.resource.Spec
}

func (wrapper *InternetGatewayWrapper) GetStatus() *schema.InternetGatewayStatus {
	return wrapper.resource.Status
}

type RouteTableWrapper struct {
	*BaseResourceWrapper[schema.RouteTable]
}

func NewRouteTableWrapper(resource *schema.RouteTable) *RouteTableWrapper {
	return &RouteTableWrapper{
		BaseResourceWrapper: newBaseResourceWrapper(resource),
	}
}

func (wrapper *RouteTableWrapper) GetResource() *schema.RouteTable {
	return wrapper.resource
}

func (wrapper *RouteTableWrapper) GetLabels() schema.Labels {
	return wrapper.resource.Labels
}

func (wrapper *RouteTableWrapper) GetAnnotations() schema.Annotations {
	return wrapper.resource.Annotations
}

func (wrapper *RouteTableWrapper) GetExtensions() schema.Extensions {
	return wrapper.resource.Extensions
}

func (wrapper *RouteTableWrapper) GetMetadata() *schema.RegionalNetworkResourceMetadata {
	return wrapper.resource.Metadata
}

func (wrapper *RouteTableWrapper) GetSpec() *schema.RouteTableSpec {
	return &wrapper.resource.Spec
}

func (wrapper *RouteTableWrapper) GetStatus() *schema.RouteTableStatus {
	return wrapper.resource.Status
}

type SubnetWrapper struct {
	*BaseResourceWrapper[schema.Subnet]
}

func NewSubnetWrapper(resource *schema.Subnet) *SubnetWrapper {
	return &SubnetWrapper{
		BaseResourceWrapper: newBaseResourceWrapper(resource),
	}
}

func (wrapper *SubnetWrapper) GetResource() *schema.Subnet {
	return wrapper.resource
}

func (wrapper *SubnetWrapper) GetLabels() schema.Labels {
	return wrapper.resource.Labels
}

func (wrapper *SubnetWrapper) GetMetadata() *schema.RegionalNetworkResourceMetadata {
	return wrapper.resource.Metadata
}

func (wrapper *SubnetWrapper) GetAnnotations() schema.Annotations {
	return wrapper.resource.Annotations
}

func (wrapper *SubnetWrapper) GetExtensions() schema.Extensions {
	return wrapper.resource.Extensions
}

func (wrapper *SubnetWrapper) GetSpec() *schema.SubnetSpec {
	return &wrapper.resource.Spec
}

func (wrapper *SubnetWrapper) GetStatus() *schema.SubnetStatus {
	return wrapper.resource.Status
}

type PublicIpWrapper struct {
	*BaseResourceWrapper[schema.PublicIp]
}

func NewPublicIpWrapper(resource *schema.PublicIp) *PublicIpWrapper {
	return &PublicIpWrapper{
		BaseResourceWrapper: newBaseResourceWrapper(resource),
	}
}

func (wrapper *PublicIpWrapper) GetResource() *schema.PublicIp {
	return wrapper.resource
}

func (wrapper *PublicIpWrapper) GetLabels() schema.Labels {
	return wrapper.resource.Labels
}

func (wrapper *PublicIpWrapper) GetAnnotations() schema.Annotations {
	return wrapper.resource.Annotations
}

func (wrapper *PublicIpWrapper) GetExtensions() schema.Extensions {
	return wrapper.resource.Extensions
}

func (wrapper *PublicIpWrapper) GetMetadata() *schema.RegionalWorkspaceResourceMetadata {
	return wrapper.resource.Metadata
}

func (wrapper *PublicIpWrapper) GetSpec() *schema.PublicIpSpec {
	return &wrapper.resource.Spec
}

func (wrapper *PublicIpWrapper) GetStatus() *schema.PublicIpStatus {
	return wrapper.resource.Status
}

type NicWrapper struct {
	*BaseResourceWrapper[schema.Nic]
}

func NewNicWrapper(resource *schema.Nic) *NicWrapper {
	return &NicWrapper{
		BaseResourceWrapper: newBaseResourceWrapper(resource),
	}
}

func (wrapper *NicWrapper) GetResource() *schema.Nic {
	return wrapper.resource
}

func (wrapper *NicWrapper) GetLabels() schema.Labels {
	return wrapper.resource.Labels
}

func (wrapper *NicWrapper) GetMetadata() *schema.RegionalWorkspaceResourceMetadata {
	return wrapper.resource.Metadata
}

func (wrapper *NicWrapper) GetAnnotations() schema.Annotations {
	return wrapper.resource.Annotations
}

func (wrapper *NicWrapper) GetExtensions() schema.Extensions {
	return wrapper.resource.Extensions
}

func (wrapper *NicWrapper) GetSpec() *schema.NicSpec {
	return &wrapper.resource.Spec
}

func (wrapper *NicWrapper) GetStatus() *schema.NicStatus {
	return wrapper.resource.Status
}

type SecurityGroupRuleWrapper struct {
	*BaseResourceWrapper[schema.SecurityGroupRule]
}

func NewSecurityGroupRuleWrapper(resource *schema.SecurityGroupRule) *SecurityGroupRuleWrapper {
	return &SecurityGroupRuleWrapper{
		BaseResourceWrapper: newBaseResourceWrapper(resource),
	}
}

func (wrapper *SecurityGroupRuleWrapper) GetResource() *schema.SecurityGroupRule {
	return wrapper.resource
}

func (wrapper *SecurityGroupRuleWrapper) GetLabels() schema.Labels {
	return wrapper.resource.Labels
}

func (wrapper *SecurityGroupRuleWrapper) GetAnnotations() schema.Annotations {
	return wrapper.resource.Annotations
}

func (wrapper *SecurityGroupRuleWrapper) GetExtensions() schema.Extensions {
	return wrapper.resource.Extensions
}

func (wrapper *SecurityGroupRuleWrapper) GetMetadata() *schema.RegionalWorkspaceResourceMetadata {
	return wrapper.resource.Metadata
}

func (wrapper *SecurityGroupRuleWrapper) GetSpec() *schema.SecurityGroupRuleSpec {
	return &wrapper.resource.Spec
}

func (wrapper *SecurityGroupRuleWrapper) GetStatus() *schema.SecurityGroupRuleStatus {
	return wrapper.resource.Status
}

type SecurityGroupWrapper struct {
	*BaseResourceWrapper[schema.SecurityGroup]
}

func NewSecurityGroupWrapper(resource *schema.SecurityGroup) *SecurityGroupWrapper {
	return &SecurityGroupWrapper{
		BaseResourceWrapper: newBaseResourceWrapper(resource),
	}
}

func (wrapper *SecurityGroupWrapper) GetResource() *schema.SecurityGroup {
	return wrapper.resource
}

func (wrapper *SecurityGroupWrapper) GetLabels() schema.Labels {
	return wrapper.resource.Labels
}

func (wrapper *SecurityGroupWrapper) GetAnnotations() schema.Annotations {
	return wrapper.resource.Annotations
}

func (wrapper *SecurityGroupWrapper) GetExtensions() schema.Extensions {
	return wrapper.resource.Extensions
}

func (wrapper *SecurityGroupWrapper) GetMetadata() *schema.RegionalWorkspaceResourceMetadata {
	return wrapper.resource.Metadata
}

func (wrapper *SecurityGroupWrapper) GetSpec() *schema.SecurityGroupSpec {
	return &wrapper.resource.Spec
}

func (wrapper *SecurityGroupWrapper) GetStatus() *schema.SecurityGroupStatus {
	return wrapper.resource.Status
}
