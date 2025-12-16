package mock

import (
	network "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.network.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Network

func (configurator *scenarioConfigurator) configureCreateNetworkStub(response *schema.Network, url string, params HasParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newNetworkStatus(schema.ResourceStateCreating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateNetworkStub(response *schema.Network, url string, params HasParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setNetworkState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveNetworkStub(response *schema.Network, url string, params HasParams) error {
	setNetworkState(response.Status, schema.ResourceStateActive)
	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetListNetworkStub(response *network.NetworkIterator, url string, params HasParams, pathParams map[string]string) error {
	if err := configurator.configureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Internet gateway

func (configurator *scenarioConfigurator) configureCreateInternetGatewayStub(response *schema.InternetGateway, url string, params HasParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newResourceStatus(schema.ResourceStateCreating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateInternetGatewayStub(response *schema.InternetGateway, url string, params HasParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setResourceState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveInternetGatewayStub(response *schema.InternetGateway, url string, params HasParams) error {
	setResourceState(response.Status, schema.ResourceStateActive)
	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetListInternetGatewayStub(response *network.InternetGatewayIterator, url string, params HasParams, pathParams map[string]string) error {
	if err := configurator.configureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Route table

func (configurator *scenarioConfigurator) configureCreateRouteTableStub(response *schema.RouteTable, url string, params HasParams) error {
	setCreatedRegionalNetworkResourceMetadata(response.Metadata)
	response.Status = newRouteTableStatus(schema.ResourceStateCreating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateRouteTableStub(response *schema.RouteTable, url string, params HasParams) error {
	setModifiedRegionalNetworkResourceMetadata(response.Metadata)
	setRouteTableState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveRouteTableStub(response *schema.RouteTable, url string, params HasParams) error {
	setRouteTableState(response.Status, schema.ResourceStateActive)
	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetListRouteTableStub(response *network.RouteTableIterator, url string, params HasParams, pathParams map[string]string) error {
	if err := configurator.configureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Subnet

func (configurator *scenarioConfigurator) configureCreateSubnetStub(response *schema.Subnet, url string, params HasParams) error {
	setCreatedRegionalNetworkResourceMetadata(response.Metadata)
	response.Status = newSubnetStatus(schema.ResourceStateCreating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateSubnetStub(response *schema.Subnet, url string, params HasParams) error {
	setModifiedRegionalNetworkResourceMetadata(response.Metadata)
	setSubnetState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveSubnetStub(response *schema.Subnet, url string, params HasParams) error {
	setSubnetState(response.Status, schema.ResourceStateActive)
	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetListSubnetStub(response *network.SubnetIterator, url string, params HasParams, pathParams map[string]string) error {
	if err := configurator.configureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Public ip

func (configurator *scenarioConfigurator) configureCreatePublicIpStub(response *schema.PublicIp, url string, params HasParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newPublicIpStatus(schema.ResourceStateCreating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdatePublicIpStub(response *schema.PublicIp, url string, params HasParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setPublicIpState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActivePublicIpStub(response *schema.PublicIp, url string, params HasParams) error {
	setPublicIpState(response.Status, schema.ResourceStateActive)
	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetListPublicIpStub(response *network.PublicIpIterator, url string, params HasParams, pathParams map[string]string) error {
	if err := configurator.configureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Nic

func (configurator *scenarioConfigurator) configureCreateNicStub(response *schema.Nic, url string, params HasParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newNicStatus(schema.ResourceStateCreating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateNicStub(response *schema.Nic, url string, params HasParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setNicState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveNicStub(response *schema.Nic, url string, params HasParams) error {
	setNicState(response.Status, schema.ResourceStateActive)
	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetListNicStub(response *network.NicIterator, url string, params HasParams, pathParams map[string]string) error {
	if err := configurator.configureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Security group

func (configurator *scenarioConfigurator) configureCreateSecurityGroupStub(response *schema.SecurityGroup, url string, params HasParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newSecurityGroupStatus(schema.ResourceStateCreating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureUpdateSecurityGroupStub(response *schema.SecurityGroup, url string, params HasParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setSecurityGroupState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.configurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetActiveSecurityGroupStub(response *schema.SecurityGroup, url string, params HasParams) error {
	setSecurityGroupState(response.Status, schema.ResourceStateActive)
	if err := configurator.configureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) configureGetListSecurityGroupStub(response *network.SecurityGroupIterator, url string, params HasParams, pathParams map[string]string) error {
	if err := configurator.configureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
