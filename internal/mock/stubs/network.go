package stubs

import (
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	network "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.network.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Network

func (configurator *scenarioConfigurator) ConfigureCreateNetworkStub(response *schema.Network, url string, params *mock.BaseParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newNetworkStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureUpdateNetworkStub(response *schema.Network, url string, params *mock.BaseParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setNetworkState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetActiveNetworkStub(response *schema.Network, url string, params *mock.BaseParams) error {
	setNetworkState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetListNetworkStub(response *network.NetworkIterator, url string, params *mock.BaseParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Internet gateway

func (configurator *scenarioConfigurator) ConfigureCreateInternetGatewayStub(response *schema.InternetGateway, url string, params *mock.BaseParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newResourceStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureUpdateInternetGatewayStub(response *schema.InternetGateway, url string, params *mock.BaseParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setResourceState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetActiveInternetGatewayStub(response *schema.InternetGateway, url string, params *mock.BaseParams) error {
	setResourceState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetListInternetGatewayStub(response *network.InternetGatewayIterator, url string, params *mock.BaseParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Route table

func (configurator *scenarioConfigurator) ConfigureCreateRouteTableStub(response *schema.RouteTable, url string, params *mock.BaseParams) error {
	setCreatedRegionalNetworkResourceMetadata(response.Metadata)
	response.Status = newRouteTableStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureUpdateRouteTableStub(response *schema.RouteTable, url string, params *mock.BaseParams) error {
	setModifiedRegionalNetworkResourceMetadata(response.Metadata)
	setRouteTableState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetActiveRouteTableStub(response *schema.RouteTable, url string, params *mock.BaseParams) error {
	setRouteTableState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetListRouteTableStub(response *network.RouteTableIterator, url string, params *mock.BaseParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Subnet

func (configurator *scenarioConfigurator) ConfigureCreateSubnetStub(response *schema.Subnet, url string, params *mock.BaseParams) error {
	setCreatedRegionalNetworkResourceMetadata(response.Metadata)
	response.Status = newSubnetStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureUpdateSubnetStub(response *schema.Subnet, url string, params *mock.BaseParams) error {
	setModifiedRegionalNetworkResourceMetadata(response.Metadata)
	setSubnetState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetActiveSubnetStub(response *schema.Subnet, url string, params *mock.BaseParams) error {
	setSubnetState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetListSubnetStub(response *network.SubnetIterator, url string, params *mock.BaseParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Public ip

func (configurator *scenarioConfigurator) ConfigureCreatePublicIpStub(response *schema.PublicIp, url string, params *mock.BaseParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newPublicIpStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureUpdatePublicIpStub(response *schema.PublicIp, url string, params *mock.BaseParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setPublicIpState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetActivePublicIpStub(response *schema.PublicIp, url string, params *mock.BaseParams) error {
	setPublicIpState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetListPublicIpStub(response *network.PublicIpIterator, url string, params *mock.BaseParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Nic

func (configurator *scenarioConfigurator) ConfigureCreateNicStub(response *schema.Nic, url string, params *mock.BaseParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newNicStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureUpdateNicStub(response *schema.Nic, url string, params *mock.BaseParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setNicState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetActiveNicStub(response *schema.Nic, url string, params *mock.BaseParams) error {
	setNicState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetListNicStub(response *network.NicIterator, url string, params *mock.BaseParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Security group

func (configurator *scenarioConfigurator) ConfigureCreateSecurityGroupStub(response *schema.SecurityGroup, url string, params *mock.BaseParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newSecurityGroupStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureUpdateSecurityGroupStub(response *schema.SecurityGroup, url string, params *mock.BaseParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setSecurityGroupState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetActiveSecurityGroupStub(response *schema.SecurityGroup, url string, params *mock.BaseParams) error {
	setSecurityGroupState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetListSecurityGroupStub(response *network.SecurityGroupIterator, url string, params *mock.BaseParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *scenarioConfigurator) ConfigureGetListNetworkSkuStub(response *network.SkuIterator, url string, params *mock.BaseParams, pathParams map[string]string) error {
	response.Metadata.Verb = http.MethodGet

	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
