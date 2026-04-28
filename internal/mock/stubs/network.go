package stubs

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/wiremock"
	network "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.network.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Network

func (configurator *Configurator) ConfigureCreateNetworkStub(response *schema.Network, url string, params wiremock.MockParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newNetworkStatus(schema.ResourceStatePending)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureUpdateNetworkStub(response *schema.Network, url string, params wiremock.MockParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = beforeUpdateNetworkStatus()
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetCreatingNetworkStub(response *schema.Network, url string, params wiremock.MockParams) error {
	setNetworkState(response.Status, schema.ResourceStateCreating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetActiveNetworkStub(response *schema.Network, url string, params wiremock.MockParams) error {
	setNetworkState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetUpdatingNetworkStub(response *schema.Network, url string, params wiremock.MockParams) error {
	setNetworkState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetDeletingNetworkStub(response *schema.Network, url string, params wiremock.MockParams) error {
	setNetworkState(response.Status, schema.ResourceStateDeleting)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureListNetworkStub(response *network.NetworkIterator, url string, params wiremock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetWithPathStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Internet gateway

func (configurator *Configurator) ConfigureCreateInternetGatewayStub(response *schema.InternetGateway, url string, params wiremock.MockParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newResourceStatus(schema.ResourceStatePending)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureUpdateInternetGatewayStub(response *schema.InternetGateway, url string, params wiremock.MockParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = beforeUpdateResourceStatus()
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetCreatingInternetGatewayStub(response *schema.InternetGateway, url string, params wiremock.MockParams) error {
	setResourceState(response.Status, schema.ResourceStateCreating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetActiveInternetGatewayStub(response *schema.InternetGateway, url string, params wiremock.MockParams) error {
	setResourceState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetUpdatingInternetGatewayStub(response *schema.InternetGateway, url string, params wiremock.MockParams) error {
	setResourceState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetDeletingInternetGatewayStub(response *schema.InternetGateway, url string, params wiremock.MockParams) error {
	setResourceState(response.Status, schema.ResourceStateDeleting)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureListInternetGatewayStub(response *network.InternetGatewayIterator, url string, params wiremock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetWithPathStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Route table

func (configurator *Configurator) ConfigureCreateRouteTableStub(response *schema.RouteTable, url string, params wiremock.MockParams) error {
	setCreatedRegionalNetworkResourceMetadata(response.Metadata)
	response.Status = newRouteTableStatus(schema.ResourceStatePending)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureUpdateRouteTableStub(response *schema.RouteTable, url string, params wiremock.MockParams) error {
	setModifiedRegionalNetworkResourceMetadata(response.Metadata)
	response.Status = beforeUpdateRouteTableStatus()
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetCreatingRouteTableStub(response *schema.RouteTable, url string, params wiremock.MockParams) error {
	setRouteTableState(response.Status, schema.ResourceStateCreating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetActiveRouteTableStub(response *schema.RouteTable, url string, params wiremock.MockParams) error {
	setRouteTableState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetUpdatingRouteTableStub(response *schema.RouteTable, url string, params wiremock.MockParams) error {
	setRouteTableState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetDeletingRouteTableStub(response *schema.RouteTable, url string, params wiremock.MockParams) error {
	setRouteTableState(response.Status, schema.ResourceStateDeleting)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureListRouteTableStub(response *network.RouteTableIterator, url string, params wiremock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetWithPathStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Subnet

func (configurator *Configurator) ConfigureCreateSubnetStub(response *schema.Subnet, url string, params wiremock.MockParams) error {
	setCreatedRegionalNetworkResourceMetadata(response.Metadata)
	response.Status = newSubnetStatus(schema.ResourceStatePending)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureUpdateSubnetStub(response *schema.Subnet, url string, params wiremock.MockParams) error {
	setModifiedRegionalNetworkResourceMetadata(response.Metadata)
	response.Status = beforeUpdateSubnetStatus()
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetCreatingSubnetStub(response *schema.Subnet, url string, params wiremock.MockParams) error {
	setSubnetState(response.Status, schema.ResourceStateCreating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetActiveSubnetStub(response *schema.Subnet, url string, params wiremock.MockParams) error {
	setSubnetState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetUpdatingSubnetStub(response *schema.Subnet, url string, params wiremock.MockParams) error {
	setSubnetState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetDeletingSubnetStub(response *schema.Subnet, url string, params wiremock.MockParams) error {
	setSubnetState(response.Status, schema.ResourceStateDeleting)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureListSubnetStub(response *network.SubnetIterator, url string, params wiremock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetWithPathStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Public ip

func (configurator *Configurator) ConfigureCreatePublicIpStub(response *schema.PublicIp, url string, params wiremock.MockParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newPublicIpStatus(schema.ResourceStatePending)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureUpdatePublicIpStub(response *schema.PublicIp, url string, params wiremock.MockParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = beforeUpdatePublicIpStatus()
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetCreatingPublicIpStub(response *schema.PublicIp, url string, params wiremock.MockParams) error {
	setPublicIpState(response.Status, schema.ResourceStateCreating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetActivePublicIpStub(response *schema.PublicIp, url string, params wiremock.MockParams) error {
	setPublicIpState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetUpdatingPublicIpStub(response *schema.PublicIp, url string, params wiremock.MockParams) error {
	setPublicIpState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetDeletingPublicIpStub(response *schema.PublicIp, url string, params wiremock.MockParams) error {
	setPublicIpState(response.Status, schema.ResourceStateDeleting)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureListPublicIpStub(response *network.PublicIpIterator, url string, params wiremock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetWithPathStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Nic

func (configurator *Configurator) ConfigureCreateNicStub(response *schema.Nic, url string, params wiremock.MockParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newNicStatus(schema.ResourceStatePending)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureUpdateNicStub(response *schema.Nic, url string, params wiremock.MockParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = beforeUpdateNicStatus()
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetCreatingNicStub(response *schema.Nic, url string, params wiremock.MockParams) error {
	setNicState(response.Status, schema.ResourceStateCreating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetActiveNicStub(response *schema.Nic, url string, params wiremock.MockParams) error {
	setNicState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetUpdatingNicStub(response *schema.Nic, url string, params wiremock.MockParams) error {
	setNicState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetDeletingNicStub(response *schema.Nic, url string, params wiremock.MockParams) error {
	setNicState(response.Status, schema.ResourceStateDeleting)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureListNicStub(response *network.NicIterator, url string, params wiremock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetWithPathStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Security group rule

func (configurator *Configurator) ConfigureCreateSecurityGroupRuleStub(response *schema.SecurityGroupRule, url string, params wiremock.MockParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newSecurityGroupRuleStatus(schema.ResourceStatePending)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureUpdateSecurityGroupRuleStub(response *schema.SecurityGroupRule, url string, params wiremock.MockParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = beforeUpdateSecurityGroupRuleStatus()

	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetCreatingSecurityGroupRuleStub(response *schema.SecurityGroupRule, url string, params wiremock.MockParams) error {
	setSecurityGroupRuleState(response.Status, schema.ResourceStateCreating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetActiveSecurityGroupRuleStub(response *schema.SecurityGroupRule, url string, params wiremock.MockParams) error {
	setSecurityGroupRuleState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetUpdatingSecurityGroupRuleStub(response *schema.SecurityGroupRule, url string, params wiremock.MockParams) error {
	setSecurityGroupRuleState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetDeletingSecurityGroupRuleStub(response *schema.SecurityGroupRule, url string, params wiremock.MockParams) error {
	setSecurityGroupRuleState(response.Status, schema.ResourceStateDeleting)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Security group

func (configurator *Configurator) ConfigureCreateSecurityGroupStub(response *schema.SecurityGroup, url string, params wiremock.MockParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newSecurityGroupStatus(schema.ResourceStatePending)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureUpdateSecurityGroupStub(response *schema.SecurityGroup, url string, params wiremock.MockParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = beforeUpdateSecurityGroupStatus()
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetCreatingSecurityGroupStub(response *schema.SecurityGroup, url string, params wiremock.MockParams) error {
	setSecurityGroupState(response.Status, schema.ResourceStateCreating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetActiveSecurityGroupStub(response *schema.SecurityGroup, url string, params wiremock.MockParams) error {
	setSecurityGroupState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetUpdatingSecurityGroupStub(response *schema.SecurityGroup, url string, params wiremock.MockParams) error {
	setSecurityGroupState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetDeletingSecurityGroupStub(response *schema.SecurityGroup, url string, params wiremock.MockParams) error {
	setSecurityGroupState(response.Status, schema.ResourceStateDeleting)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureListSecurityGroupRuleStub(response *network.SecurityGroupRuleIterator, url string, params wiremock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetWithPathStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureListSecurityGroupStub(response *network.SecurityGroupIterator, url string, params wiremock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetWithPathStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureListNetworkSkuStub(response *network.SkuIterator, url string, params wiremock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetWithPathStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
