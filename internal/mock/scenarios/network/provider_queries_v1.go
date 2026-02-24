package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func ConfigureProviderQueriesV1(scenario *mockscenarios.Scenario, params *params.NetworkProviderQueriesV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := *params.Workspace
	networks := params.Networks
	internetGateways := params.InternetGateways
	routeTables := params.RouteTables
	subnets := params.Subnets
	nics := params.Nics
	publicIps := params.PublicIps
	securityGroups := params.SecurityGroups

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	networkListUrl := generators.GenerateNetworkListURL(constants.NetworkProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	gatewayListUrl := generators.GenerateInternetGatewayListURL(constants.NetworkProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	publicIpListUrl := generators.GeneratePublicIpListURL(constants.NetworkProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	nicListUrl := generators.GenerateNicListURL(constants.NetworkProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	securityGroupListUrl := generators.GenerateSecurityGroupListURL(constants.NetworkProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	skuListUrl := generators.GenerateNetworkSkuListURL(constants.NetworkProviderV1, workspace.Metadata.Tenant)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(workspace.Metadata.Name).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(workspace.Metadata.Tenant).Region(workspace.Metadata.Region).
		Labels(workspace.Labels).
		Build()
	if err != nil {
		return err
	}

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspaceResponse, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Create networks
	err = stubs.BulkCreateNetworksStubV1(configurator, scenario.MockParams, networks)
	if err != nil {
		return err
	}
	networkResponse, err := builders.NewNetworkIteratorBuilder().
		Provider(constants.NetworkProviderV1).
		Tenant(workspace.Metadata.Tenant).Workspace(workspace.Metadata.Name).
		Items(networks).
		Build()
	if err != nil {
		return err
	}

	// List
	if err := configurator.ConfigureGetListNetworkStub(networkResponse, networkListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with Limit 1
	networkResponse.Items = networks[:1]
	if err := configurator.ConfigureGetListNetworkStub(networkResponse, networkListUrl, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
		return err
	}

	// List with Label
	networksWithLabel := func(networkList []schema.Network) []schema.Network {
		var filteredNetworks []schema.Network
		for _, network := range networkList {
			if val, ok := network.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredNetworks = append(filteredNetworks, network)
			}
		}
		return filteredNetworks
	}
	networkResponse.Items = networksWithLabel(networks)
	if err := configurator.ConfigureGetListNetworkStub(networkResponse, networkListUrl, scenario.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with Limit and Label
	networkResponse.Items = networksWithLabel(networks)[:1]
	if err := configurator.ConfigureGetListNetworkStub(networkResponse, networkListUrl, scenario.MockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Test Network Skus
	// Create skus
	skusList := steps.GenerateNetworkSkusV1(workspace.Metadata.Tenant)
	skuResponse, err := builders.NewNetworkSkuIteratorBuilder().Provider(constants.StorageProviderV1).Tenant(workspace.Metadata.Tenant).Items(skusList).Build()
	if err != nil {
		return err
	}

	// List skus
	if err := configurator.ConfigureGetListNetworkSkuStub(skuResponse, skuListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List skus with limit 1
	if err := configurator.ConfigureGetListNetworkSkuStub(skuResponse, skuListUrl, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
		return err
	}

	// Create internet gateways
	err = stubs.BulkCreateInternetGatewaysStubV1(configurator, scenario.MockParams, internetGateways)
	if err != nil {
		return err
	}
	gatewayResponse, err := builders.NewInternetGatewayIteratorBuilder().
		Provider(constants.NetworkProviderV1).
		Tenant(workspace.Metadata.Tenant).Workspace(workspace.Metadata.Name).
		Items(internetGateways).
		Build()
	if err != nil {
		return err
	}

	// List
	if err := configurator.ConfigureGetListInternetGatewayStub(gatewayResponse, gatewayListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with Limit 1
	gatewayResponse.Items = internetGateways[:1]
	if err := configurator.ConfigureGetListInternetGatewayStub(gatewayResponse, gatewayListUrl, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
		return err
	}

	// List with Label
	gatewayWithLabel := func(gatewayList []schema.InternetGateway) []schema.InternetGateway {
		var filteredGateway []schema.InternetGateway
		for _, gateway := range gatewayList {
			if val, ok := gateway.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredGateway = append(filteredGateway, gateway)
			}
		}
		return filteredGateway
	}
	gatewayResponse.Items = gatewayWithLabel(internetGateways)
	if err := configurator.ConfigureGetListInternetGatewayStub(gatewayResponse, gatewayListUrl, scenario.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with Limit and Label
	gatewayResponse.Items = gatewayWithLabel(internetGateways)[:1]
	if err := configurator.ConfigureGetListInternetGatewayStub(gatewayResponse, gatewayListUrl, scenario.MockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Create route tables
	err = stubs.BulkCreateRouteTableStubV1(configurator, scenario.MockParams, routeTables)
	if err != nil {
		return err
	}
	networkName := networks[0].Metadata.Name
	routeTableResponse, err := builders.NewRouteTableIteratorBuilder().
		Provider(constants.NetworkProviderV1).
		Tenant(workspace.Metadata.Tenant).Workspace(workspace.Metadata.Name).Network(networkName).
		Items(routeTables).
		Build()
	if err != nil {
		return err
	}

	// List
	routeTableListUrl := generators.GenerateRouteTableListURL(constants.NetworkProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name, networkName)
	if err := configurator.ConfigureGetListRouteTableStub(routeTableResponse, routeTableListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with Limit 1
	routeTableResponse.Items = routeTables[:1]
	if err := configurator.ConfigureGetListRouteTableStub(routeTableResponse, routeTableListUrl, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
		return err
	}

	// List with Label
	routeTableWithLabel := func(routeTableList []schema.RouteTable) []schema.RouteTable {
		var filteredRouteTable []schema.RouteTable
		for _, routeTable := range routeTableList {
			if val, ok := routeTable.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredRouteTable = append(filteredRouteTable, routeTable)
			}
		}
		return filteredRouteTable
	}

	routeTableResponse.Items = routeTableWithLabel(routeTables)
	if err := configurator.ConfigureGetListRouteTableStub(routeTableResponse, routeTableListUrl, scenario.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with Limit and Label
	routeTableResponse.Items = routeTableWithLabel(routeTables)[:1]
	if err := configurator.ConfigureGetListRouteTableStub(routeTableResponse, routeTableListUrl, scenario.MockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Subnet
	err = stubs.BulkCreateSubnetsStubV1(configurator, scenario.MockParams, subnets)
	if err != nil {
		return err
	}
	subnetResponse, err := builders.NewSubnetIteratorBuilder().
		Provider(constants.NetworkProviderV1).
		Tenant(workspace.Metadata.Tenant).Workspace(workspace.Metadata.Name).Network(networkName).
		Items(subnets).
		Build()
	if err != nil {
		return err
	}

	// List
	subnetListUrl := generators.GenerateSubnetListURL(constants.NetworkProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name, networkName)
	if err := configurator.ConfigureGetListSubnetStub(subnetResponse, subnetListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with Limit 1
	subnetResponse.Items = subnets[:1]
	if err := configurator.ConfigureGetListSubnetStub(subnetResponse, subnetListUrl, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
		return err
	}

	// List with Label
	subnetWithLabel := func(subnetList []schema.Subnet) []schema.Subnet {
		var filteredSubnet []schema.Subnet
		for _, subnet := range subnetList {
			if val, ok := subnet.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredSubnet = append(filteredSubnet, subnet)
			}
		}
		return filteredSubnet
	}
	subnetResponse.Items = subnetWithLabel(subnets)
	if err := configurator.ConfigureGetListSubnetStub(subnetResponse, subnetListUrl, scenario.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with Limit and Label
	subnetResponse.Items = subnetWithLabel(subnets)[:1]
	if err := configurator.ConfigureGetListSubnetStub(subnetResponse, subnetListUrl, scenario.MockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Cretae public ips
	err = stubs.BulkCreatePublicIpsStubV1(configurator, scenario.MockParams, publicIps)
	if err != nil {
		return err
	}
	publicIpResponse, err := builders.NewPublicIpIteratorBuilder().
		Provider(constants.NetworkProviderV1).
		Tenant(workspace.Metadata.Tenant).Workspace(workspace.Metadata.Name).
		Items(publicIps).
		Build()
	if err != nil {
		return err
	}

	// List
	if err := configurator.ConfigureGetListPublicIpStub(publicIpResponse, publicIpListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with Limit 1
	publicIpResponse.Items = publicIps[:1]
	if err := configurator.ConfigureGetListPublicIpStub(publicIpResponse, publicIpListUrl, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
		return err
	}

	// List with Label
	publicIpWithLabel := func(publicIpList []schema.PublicIp) []schema.PublicIp {
		var filteredPublicIp []schema.PublicIp
		for _, publicIp := range publicIpList {
			if val, ok := publicIp.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredPublicIp = append(filteredPublicIp, publicIp)
			}
		}
		return filteredPublicIp
	}
	publicIpResponse.Items = publicIpWithLabel(publicIps)
	if err := configurator.ConfigureGetListPublicIpStub(publicIpResponse, publicIpListUrl, scenario.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with Limit and Label
	publicIpResponse.Items = publicIpWithLabel(publicIps)[:1]
	if err := configurator.ConfigureGetListPublicIpStub(publicIpResponse, publicIpListUrl, scenario.MockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Create nics
	err = stubs.BulkCreateNicsStubV1(configurator, scenario.MockParams, nics)
	if err != nil {
		return err
	}
	nicResponse, err := builders.NewNicIteratorBuilder().
		Provider(constants.NetworkProviderV1).
		Tenant(workspace.Metadata.Tenant).Workspace(workspace.Metadata.Name).
		Items(nics).
		Build()
	if err != nil {
		return err
	}

	// List
	if err := configurator.ConfigureGetListNicStub(nicResponse, nicListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with Limit 1
	nicResponse.Items = nics[:1]
	if err := configurator.ConfigureGetListNicStub(nicResponse, nicListUrl, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
		return err
	}

	// List with Label
	nicWithLabel := func(nicList []schema.Nic) []schema.Nic {
		var filteredNic []schema.Nic
		for _, nic := range nicList {
			if val, ok := nic.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredNic = append(filteredNic, nic)
			}
		}
		return filteredNic
	}
	nicResponse.Items = nicWithLabel(nics)
	if err := configurator.ConfigureGetListNicStub(nicResponse, nicListUrl, scenario.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with Limit and Label
	nicResponse.Items = nicWithLabel(nics)[:1]
	if err := configurator.ConfigureGetListNicStub(nicResponse, nicListUrl, scenario.MockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Create security groups
	err = stubs.BulkCreateSecurityGroupsStubV1(configurator, scenario.MockParams, securityGroups)
	if err != nil {
		return err
	}
	securityGroupResponse, err := builders.NewSecurityGroupIteratorBuilder().
		Provider(constants.NetworkProviderV1).
		Tenant(workspace.Metadata.Tenant).Workspace(workspace.Metadata.Name).
		Items(securityGroups).
		Build()
	if err != nil {
		return err
	}

	// List
	if err := configurator.ConfigureGetListSecurityGroupStub(securityGroupResponse, securityGroupListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with Limit 1
	securityGroupResponse.Items = securityGroups[:1]
	if err := configurator.ConfigureGetListSecurityGroupStub(securityGroupResponse, securityGroupListUrl, scenario.MockParams, mock.PathParamsLimit("1")); err != nil {
		return err
	}

	// List with Label
	secGroupWithLabel := func(securityGroupList []schema.SecurityGroup) []schema.SecurityGroup {
		var filteredSecurity []schema.SecurityGroup
		for _, sec := range securityGroupList {
			if val, ok := sec.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredSecurity = append(filteredSecurity, sec)
			}
		}
		return filteredSecurity
	}

	securityGroupResponse.Items = secGroupWithLabel(securityGroups)
	if err := configurator.ConfigureGetListSecurityGroupStub(securityGroupResponse, securityGroupListUrl, scenario.MockParams, mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with Limit and Label
	securityGroupResponse.Items = secGroupWithLabel(securityGroups)[:1]
	if err := configurator.ConfigureGetListSecurityGroupStub(securityGroupResponse, securityGroupListUrl, scenario.MockParams, mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Delete the security groups
	for _, securityGroup := range securityGroups {
		securityGroupUrl := generators.GenerateSecurityGroupURL(constants.NetworkProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name, securityGroup.Metadata.Name)

		// Delete the security group
		if err := configurator.ConfigureDeleteStub(securityGroupUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted security group
		if err := configurator.ConfigureGetNotFoundStub(securityGroupUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete the nics
	for _, nic := range nics {
		nicUrl := generators.GenerateNicURL(constants.NetworkProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name, nic.Metadata.Name)

		// Delete the nic
		if err := configurator.ConfigureDeleteStub(nicUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted nic
		if err := configurator.ConfigureGetNotFoundStub(nicUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete the public ips
	for _, publicIp := range publicIps {
		publicIpUrl := generators.GeneratePublicIpURL(constants.NetworkProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name, publicIp.Metadata.Name)

		// Delete the public ip
		if err := configurator.ConfigureDeleteStub(publicIpUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted public ip
		if err := configurator.ConfigureGetNotFoundStub(publicIpUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete the subnets
	for _, subnet := range subnets {
		subnetUrl := generators.GenerateSubnetURL(constants.NetworkProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name, networkName, subnet.Metadata.Name)

		// Delete the subnet
		if err := configurator.ConfigureDeleteStub(subnetUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted subnet
		if err := configurator.ConfigureGetNotFoundStub(subnetUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete the route tables
	for _, routeTable := range routeTables {
		routeTableUrl := generators.GenerateRouteTableURL(constants.NetworkProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name, networkName, routeTable.Metadata.Name)

		// Delete the route table
		if err := configurator.ConfigureDeleteStub(routeTableUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted route table
		if err := configurator.ConfigureGetNotFoundStub(routeTableUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete the internet gateways
	for _, gateway := range internetGateways {
		gatewayUrl := generators.GenerateInternetGatewayURL(constants.NetworkProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name, gateway.Metadata.Name)

		// Delete the internet gateway
		if err := configurator.ConfigureDeleteStub(gatewayUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted internet gateway
		if err := configurator.ConfigureGetNotFoundStub(gatewayUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete the networks
	for _, network := range networks {
		networkUrl := generators.GenerateNetworkURL(constants.NetworkProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name, network.Metadata.Name)

		// Delete the network
		if err := configurator.ConfigureDeleteStub(networkUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted network
		if err := configurator.ConfigureGetNotFoundStub(networkUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetNotFoundStub(workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
