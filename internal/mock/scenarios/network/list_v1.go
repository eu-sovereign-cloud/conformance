package network

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/wiremock/go-wiremock"
)

func ConfigureListScenarioV1(scenario string, params *mock.NetworkListParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := stubs.NewStubConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, params.Tenant, params.Workspace.Name)
	blockUrl := generators.GenerateBlockStorageURL(constants.StorageProviderV1, params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceUrl := generators.GenerateInstanceURL(constants.ComputeProviderV1, params.Tenant, params.Workspace.Name, params.Instance.Name)
	networkListUrl := generators.GenerateNetworkListURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name)
	gatewayListUrl := generators.GenerateInternetGatewayListURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name)
	publicIpListUrl := generators.GeneratePublicIpListURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name)
	nicListUrl := generators.GenerateNicListURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name)
	securityGroupListUrl := generators.GenerateSecurityGroupListURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name)
	skuListUrl := generators.GenerateNetworkSkuListURL(constants.NetworkProviderV1, params.Tenant)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspaceResponse, workspaceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Create networks
	networkList, err := stubs.BulkCreateNetworksStubV1(configurator, params.GetBaseParams(), params.Workspace.Name, params.Networks)
	if err != nil {
		return nil, err
	}
	networkResponse, err := builders.NewNetworkIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).
		Items(networkList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	if err := configurator.ConfigureGetListNetworkStub(networkResponse, networkListUrl, params.GetBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	networkResponse.Items = networkList[:1]
	if err := configurator.ConfigureGetListNetworkStub(networkResponse, networkListUrl, params.GetBaseParams(), mock.PathParamsLimit("1")); err != nil {
		return nil, err
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
	networkResponse.Items = networksWithLabel(networkList)
	if err := configurator.ConfigureGetListNetworkStub(networkResponse, networkListUrl, params.GetBaseParams(), mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	networkResponse.Items = networksWithLabel(networkList)[:1]
	if err := configurator.ConfigureGetListNetworkStub(networkResponse, networkListUrl, params.GetBaseParams(), mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Test Network Skus
	// Create skus
	skusList := steps.GenerateNetworkSkusV1(params.GetBaseParams().Tenant)
	skuResponse, err := builders.NewNetworkSkuIteratorBuilder().Provider(constants.StorageProviderV1).Tenant(params.Tenant).Items(skusList).Build()
	if err != nil {
		return nil, err
	}

	// List skus
	if err := configurator.ConfigureGetListNetworkSkuStub(skuResponse, skuListUrl, params.GetBaseParams(), nil); err != nil {
		return nil, err
	}

	// List skus with limit 1
	if err := configurator.ConfigureGetListNetworkSkuStub(skuResponse, skuListUrl, params.GetBaseParams(), mock.PathParamsLimit("1")); err != nil {
		return nil, err
	}

	// Create internet gateways
	gatewayList, err := stubs.BulkCreateInternetGatewaysStubV1(configurator, params.GetBaseParams(), params.Workspace.Name, params.InternetGateways)
	if err != nil {
		return nil, err
	}
	gatewayResponse, err := builders.NewInternetGatewayIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).
		Items(gatewayList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	if err := configurator.ConfigureGetListInternetGatewayStub(gatewayResponse, gatewayListUrl, params.GetBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	gatewayResponse.Items = gatewayList[:1]
	if err := configurator.ConfigureGetListInternetGatewayStub(gatewayResponse, gatewayListUrl, params.GetBaseParams(), mock.PathParamsLimit("1")); err != nil {
		return nil, err
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
	gatewayResponse.Items = gatewayWithLabel(gatewayList)
	if err := configurator.ConfigureGetListInternetGatewayStub(gatewayResponse, gatewayListUrl, params.GetBaseParams(), mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	gatewayResponse.Items = gatewayWithLabel(gatewayList)[:1]
	if err := configurator.ConfigureGetListInternetGatewayStub(gatewayResponse, gatewayListUrl, params.GetBaseParams(), mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Get a network name from the list
	networkName := networkList[0].Metadata.Name

	// Create route tables
	routeTableList, err := stubs.BulkCreateRouteTableStubV1(configurator, params.GetBaseParams(), params.Workspace.Name, networkName, params.RouteTables)
	if err != nil {
		return nil, err
	}
	routeTableResponse, err := builders.NewRouteTableIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Network(networkName).
		Items(routeTableList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	routeTableListUrl := generators.GenerateRouteTableListURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, networkName)
	if err := configurator.ConfigureGetListRouteTableStub(routeTableResponse, routeTableListUrl, params.GetBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	routeTableResponse.Items = routeTableList[:1]
	if err := configurator.ConfigureGetListRouteTableStub(routeTableResponse, routeTableListUrl, params.GetBaseParams(), mock.PathParamsLimit("1")); err != nil {
		return nil, err
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

	routeTableResponse.Items = routeTableWithLabel(routeTableList)
	if err := configurator.ConfigureGetListRouteTableStub(routeTableResponse, routeTableListUrl, params.GetBaseParams(), mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	routeTableResponse.Items = routeTableWithLabel(routeTableList)[:1]
	if err := configurator.ConfigureGetListRouteTableStub(routeTableResponse, routeTableListUrl, params.GetBaseParams(), mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Subnet
	subnetList, err := stubs.BulkCreateSubnetsStubV1(configurator, params.GetBaseParams(), params.Workspace.Name, networkName, params.Subnets)
	if err != nil {
		return nil, err
	}
	subnetResponse, err := builders.NewSubnetIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Network(networkName).
		Items(subnetList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	subnetListUrl := generators.GenerateSubnetListURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, networkName)
	if err := configurator.ConfigureGetListSubnetStub(subnetResponse, subnetListUrl, params.GetBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	subnetResponse.Items = subnetList[:1]
	if err := configurator.ConfigureGetListSubnetStub(subnetResponse, subnetListUrl, params.GetBaseParams(), mock.PathParamsLimit("1")); err != nil {
		return nil, err
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
	subnetResponse.Items = subnetWithLabel(subnetList)
	if err := configurator.ConfigureGetListSubnetStub(subnetResponse, subnetListUrl, params.GetBaseParams(), mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	subnetResponse.Items = subnetWithLabel(subnetList)[:1]
	if err := configurator.ConfigureGetListSubnetStub(subnetResponse, subnetListUrl, params.GetBaseParams(), mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Cretae public ips
	publicIpList, err := stubs.BulkCreatePublicIpsStubV1(configurator, params.GetBaseParams(), params.Workspace.Name, params.PublicIps)
	if err != nil {
		return nil, err
	}
	publicIpResponse, err := builders.NewPublicIpIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).
		Items(publicIpList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	if err := configurator.ConfigureGetListPublicIpStub(publicIpResponse, publicIpListUrl, params.GetBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	publicIpResponse.Items = publicIpList[:1]
	if err := configurator.ConfigureGetListPublicIpStub(publicIpResponse, publicIpListUrl, params.GetBaseParams(), mock.PathParamsLimit("1")); err != nil {
		return nil, err
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
	publicIpResponse.Items = publicIpWithLabel(publicIpList)
	if err := configurator.ConfigureGetListPublicIpStub(publicIpResponse, publicIpListUrl, params.GetBaseParams(), mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	publicIpResponse.Items = publicIpWithLabel(publicIpList)[:1]
	if err := configurator.ConfigureGetListPublicIpStub(publicIpResponse, publicIpListUrl, params.GetBaseParams(), mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Create nics
	nicList, err := stubs.BulkCreateNicsStubV1(configurator, params.GetBaseParams(), params.Workspace.Name, params.Nics)
	if err != nil {
		return nil, err
	}
	nicResponse, err := builders.NewNicIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).
		Items(nicList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	if err := configurator.ConfigureGetListNicStub(nicResponse, nicListUrl, params.GetBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	nicResponse.Items = nicList[:1]
	if err := configurator.ConfigureGetListNicStub(nicResponse, nicListUrl, params.GetBaseParams(), mock.PathParamsLimit("1")); err != nil {
		return nil, err
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
	nicResponse.Items = nicWithLabel(nicList)
	if err := configurator.ConfigureGetListNicStub(nicResponse, nicListUrl, params.GetBaseParams(), mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	nicResponse.Items = nicWithLabel(nicList)[:1]
	if err := configurator.ConfigureGetListNicStub(nicResponse, nicListUrl, params.GetBaseParams(), mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Create security groups
	securityGroupList, err := stubs.BulkCreateSecurityGroupsStubV1(configurator, params.GetBaseParams(), params.Workspace.Name, params.SecurityGroups)
	if err != nil {
		return nil, err
	}
	securityGroupResponse, err := builders.NewSecurityGroupIteratorBuilder().
		Provider(constants.StorageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).
		Items(securityGroupList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	if err := configurator.ConfigureGetListSecurityGroupStub(securityGroupResponse, securityGroupListUrl, params.GetBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	securityGroupResponse.Items = securityGroupList[:1]
	if err := configurator.ConfigureGetListSecurityGroupStub(securityGroupResponse, securityGroupListUrl, params.GetBaseParams(), mock.PathParamsLimit("1")); err != nil {
		return nil, err
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

	securityGroupResponse.Items = secGroupWithLabel(securityGroupList)
	if err := configurator.ConfigureGetListSecurityGroupStub(securityGroupResponse, securityGroupListUrl, params.GetBaseParams(), mock.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	securityGroupResponse.Items = secGroupWithLabel(securityGroupList)[:1]
	if err := configurator.ConfigureGetListSecurityGroupStub(securityGroupResponse, securityGroupListUrl, params.GetBaseParams(), mock.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(params.BlockStorage.Name).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Labels(params.BlockStorage.InitialLabels).
		Spec(params.BlockStorage.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	if err := configurator.ConfigureCreateBlockStorageStub(blockResponse, blockUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Instance
	instanceResponse, err := builders.NewInstanceBuilder().
		Name(params.Instance.Name).
		Provider(constants.ComputeProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Labels(params.Instance.InitialLabels).
		Spec(params.Instance.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create an instance
	if err := configurator.ConfigureCreateInstanceStub(instanceResponse, instanceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the created instance
	if err := configurator.ConfigureGetActiveInstanceStub(instanceResponse, instanceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete

	// Delete the instance
	if err := configurator.ConfigureDeleteStub(instanceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted instance
	if err := configurator.ConfigureGetNotFoundStub(instanceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.ConfigureDeleteStub(blockUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configurator.ConfigureGetNotFoundStub(blockUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Delete the security groups
	for _, securityGroup := range params.SecurityGroups {
		securityGroupUrl := generators.GenerateSecurityGroupURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, securityGroup.Name)

		// Delete the security group
		if err := configurator.ConfigureDeleteStub(securityGroupUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted security group
		if err := configurator.ConfigureGetNotFoundStub(securityGroupUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete the nics
	for _, nic := range params.Nics {
		nicUrl := generators.GenerateNicURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, nic.Name)

		// Delete the nic
		if err := configurator.ConfigureDeleteStub(nicUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted nic
		if err := configurator.ConfigureGetNotFoundStub(nicUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete the public ips
	for _, publicIp := range params.PublicIps {
		publicIpUrl := generators.GeneratePublicIpURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, publicIp.Name)

		// Delete the public ip
		if err := configurator.ConfigureDeleteStub(publicIpUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted public ip
		if err := configurator.ConfigureGetNotFoundStub(publicIpUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete the subnets
	for _, subnet := range params.Subnets {
		subnetUrl := generators.GenerateSubnetURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, networkName, subnet.Name)

		// Delete the subnet
		if err := configurator.ConfigureDeleteStub(subnetUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted subnet
		if err := configurator.ConfigureGetNotFoundStub(subnetUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete the route tables
	for _, routeTable := range params.RouteTables {
		routeTableUrl := generators.GenerateRouteTableURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, networkName, routeTable.Name)

		// Delete the route table
		if err := configurator.ConfigureDeleteStub(routeTableUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted route table
		if err := configurator.ConfigureGetNotFoundStub(routeTableUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete the internet gateways
	for _, gateway := range params.InternetGateways {
		gatewayUrl := generators.GenerateInternetGatewayURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, gateway.Name)

		// Delete the internet gateway
		if err := configurator.ConfigureDeleteStub(gatewayUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted internet gateway
		if err := configurator.ConfigureGetNotFoundStub(gatewayUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete the networks
	for _, network := range params.Networks {
		networkUrl := generators.GenerateInternetGatewayURL(constants.NetworkProviderV1, params.Tenant, params.Workspace.Name, network.Name)

		// Delete the network
		if err := configurator.ConfigureDeleteStub(networkUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted network
		if err := configurator.ConfigureGetNotFoundStub(networkUrl, params.GetBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetNotFoundStub(workspaceUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}
	return configurator.Client, nil
}
