package mock

import (
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// Network

func configureCreateNetworkStub(configurator *scenarioConfigurator, response *schema.Network, url string, params HasParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newNetworkStatus(schema.ResourceStateCreating)
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func configureUpdateNetworkStub(configurator *scenarioConfigurator, response *schema.Network, url string, labels schema.Labels, params HasParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setNetworkState(response.Status, schema.ResourceStateUpdating)
	response.Labels = labels
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func configureGetActiveNetworkStub(configurator *scenarioConfigurator, response *schema.Network, url string, params HasParams) error {
	setNetworkState(response.Status, schema.ResourceStateActive)
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

// Internet gateway

func configureCreateInternetGatewayStub(configurator *scenarioConfigurator, response *schema.InternetGateway, url string, params HasParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newResourceStatus(schema.ResourceStateCreating)
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func configureUpdateInternetGatewayStub(configurator *scenarioConfigurator, response *schema.InternetGateway, url string, labels schema.Labels, params HasParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setResourceState(response.Status, schema.ResourceStateUpdating)
	response.Labels = labels
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func configureGetActiveInternetGatewayStub(configurator *scenarioConfigurator, response *schema.InternetGateway, url string, params HasParams) error {
	setResourceState(response.Status, schema.ResourceStateActive)
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

// Route table

func configureCreateRouteTableStub(configurator *scenarioConfigurator, response *schema.RouteTable, url string, params HasParams) error {
	setCreatedRegionalNetworkResourceMetadata(response.Metadata)
	response.Status = newRouteTableStatus(schema.ResourceStateCreating)
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func configureUpdateRouteTableStub(configurator *scenarioConfigurator, response *schema.RouteTable, url string, labels schema.Labels, params HasParams) error {
	setModifiedRegionalNetworkResourceMetadata(response.Metadata)
	setRouteTableState(response.Status, schema.ResourceStateUpdating)
	response.Labels = labels
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func configureGetActiveRouteTableStub(configurator *scenarioConfigurator, response *schema.RouteTable, url string, params HasParams) error {
	setRouteTableState(response.Status, schema.ResourceStateActive)
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

// Subnet

func configureCreateSubnetStub(configurator *scenarioConfigurator, response *schema.Subnet, url string, params HasParams) error {
	setCreatedRegionalNetworkResourceMetadata(response.Metadata)
	response.Status = newSubnetStatus(schema.ResourceStateCreating)
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func configureUpdateSubnetStub(configurator *scenarioConfigurator, response *schema.Subnet, url string, labels schema.Labels, params HasParams) error {
	setModifiedRegionalNetworkResourceMetadata(response.Metadata)
	setSubnetState(response.Status, schema.ResourceStateUpdating)
	response.Labels = labels
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func configureGetActiveSubnetStub(configurator *scenarioConfigurator, response *schema.Subnet, url string, params HasParams) error {
	setSubnetState(response.Status, schema.ResourceStateActive)
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

// Public ip

func configureCreatePublicIpStub(configurator *scenarioConfigurator, response *schema.PublicIp, url string, params HasParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newPublicIpStatus(schema.ResourceStateCreating)
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func configureUpdatePublicIpStub(configurator *scenarioConfigurator, response *schema.PublicIp, url string, labels schema.Labels, params HasParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setPublicIpState(response.Status, schema.ResourceStateUpdating)
	response.Labels = labels
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func configureGetActivePublicIpStub(configurator *scenarioConfigurator, response *schema.PublicIp, url string, params HasParams) error {
	setPublicIpState(response.Status, schema.ResourceStateActive)
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

// Nic

func configureCreateNicStub(configurator *scenarioConfigurator, response *schema.Nic, url string, params HasParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newNicStatus(schema.ResourceStateCreating)
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func configureUpdateNicStub(configurator *scenarioConfigurator, response *schema.Nic, url string, labels schema.Labels, params HasParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setNicState(response.Status, schema.ResourceStateUpdating)
	response.Labels = labels
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func configureGetActiveNicStub(configurator *scenarioConfigurator, response *schema.Nic, url string, params HasParams) error {
	setNicState(response.Status, schema.ResourceStateActive)
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

// Security group

func configureCreateSecurityGroupStub(configurator *scenarioConfigurator, response *schema.SecurityGroup, url string, params HasParams) error {
	setCreatedRegionalWorkspaceResourceMetadata(response.Metadata)
	response.Status = newSecurityGroupStatus(schema.ResourceStateCreating)
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func configureUpdateSecurityGroupStub(configurator *scenarioConfigurator, response *schema.SecurityGroup, url string, labels schema.Labels, params HasParams) error {
	setModifiedRegionalWorkspaceResourceMetadata(response.Metadata)
	setSecurityGroupState(response.Status, schema.ResourceStateUpdating)
	response.Labels = labels
	response.Metadata.Verb = http.MethodPut

	if err := configurator.configurePutStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}

func configureGetActiveSecurityGroupStub(configurator *scenarioConfigurator, response *schema.SecurityGroup, url string, params HasParams) error {
	setSecurityGroupState(response.Status, schema.ResourceStateActive)
	response.Metadata.Verb = http.MethodGet

	if err := configurator.configureGetStub(url, params, response, false); err != nil {
		return err
	}
	return nil
}
