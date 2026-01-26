package stubs

import (
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	network "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.network.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Network

func (configurator *stubConfigurator) ConfigureCreateNetworkStub(response *schema.Network, url string, params *mock.MockParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newNetworkStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureUpdateNetworkStub(response *schema.Network, url string, params *mock.MockParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setNetworkState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetActiveNetworkStub(response *schema.Network, url string, params *mock.MockParams) error {
	setNetworkState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetListNetworkStub(response *network.NetworkIterator, url string, params *mock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Internet gateway

func (configurator *stubConfigurator) ConfigureCreateInternetGatewayStub(response *schema.InternetGateway, url string, params *mock.MockParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newResourceStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureUpdateInternetGatewayStub(response *schema.InternetGateway, url string, params *mock.MockParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setResourceState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetActiveInternetGatewayStub(response *schema.InternetGateway, url string, params *mock.MockParams) error {
	setResourceState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetListInternetGatewayStub(response *network.InternetGatewayIterator, url string, params *mock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Route table

func (configurator *stubConfigurator) ConfigureCreateRouteTableStub(response *schema.RouteTable, url string, params *mock.MockParams) error {
	setCreatedRegionalNetworkResourceMetadata(response.Metadata)
	response.Status = newRouteTableStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureUpdateRouteTableStub(response *schema.RouteTable, url string, params *mock.MockParams) error {
	setModifiedRegionalNetworkResourceMetadata(response.Metadata)
	setRouteTableState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetActiveRouteTableStub(response *schema.RouteTable, url string, params *mock.MockParams) error {
	setRouteTableState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetListRouteTableStub(response *network.RouteTableIterator, url string, params *mock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Subnet

func (configurator *stubConfigurator) ConfigureCreateSubnetStub(response *schema.Subnet, url string, params *mock.MockParams) error {
	setCreatedRegionalNetworkResourceMetadata(response.Metadata)
	response.Status = newSubnetStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureUpdateSubnetStub(response *schema.Subnet, url string, params *mock.MockParams) error {
	setModifiedRegionalNetworkResourceMetadata(response.Metadata)
	setSubnetState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetActiveSubnetStub(response *schema.Subnet, url string, params *mock.MockParams) error {
	setSubnetState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetListSubnetStub(response *network.SubnetIterator, url string, params *mock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Public ip

func (configurator *stubConfigurator) ConfigureCreatePublicIpStub(response *schema.PublicIp, url string, params *mock.MockParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newPublicIpStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureUpdatePublicIpStub(response *schema.PublicIp, url string, params *mock.MockParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setPublicIpState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetActivePublicIpStub(response *schema.PublicIp, url string, params *mock.MockParams) error {
	setPublicIpState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetListPublicIpStub(response *network.PublicIpIterator, url string, params *mock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Nic

func (configurator *stubConfigurator) ConfigureCreateNicStub(response *schema.Nic, url string, params *mock.MockParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newNicStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureUpdateNicStub(response *schema.Nic, url string, params *mock.MockParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setNicState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetActiveNicStub(response *schema.Nic, url string, params *mock.MockParams) error {
	setNicState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetListNicStub(response *network.NicIterator, url string, params *mock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

// Security group

func (configurator *stubConfigurator) ConfigureCreateSecurityGroupStub(response *schema.SecurityGroup, url string, params *mock.MockParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newSecurityGroupStatus(schema.ResourceStateCreating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureUpdateSecurityGroupStub(response *schema.SecurityGroup, url string, params *mock.MockParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setSecurityGroupState(response.Status, schema.ResourceStateUpdating)
	if err := configurator.ConfigurePutStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetActiveSecurityGroupStub(response *schema.SecurityGroup, url string, params *mock.MockParams) error {
	setSecurityGroupState(response.Status, schema.ResourceStateActive)
	if err := configurator.ConfigureGetStub(url, params, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetListSecurityGroupStub(response *network.SecurityGroupIterator, url string, params *mock.MockParams, pathParams map[string]string) error {
	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetListNetworkSkuStub(response *network.SkuIterator, url string, params *mock.MockParams, pathParams map[string]string) error {
	response.Metadata.Verb = http.MethodGet

	if err := configurator.ConfigureGetListStub(url, params, pathParams, func(verb string) { response.Metadata.Verb = verb }, response); err != nil {
		return err
	}
	return nil
}
