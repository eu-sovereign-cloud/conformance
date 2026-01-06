package mock

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigureNetworkLifecycleScenarioV1(scenario string, params *NetworkLifeCycleParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(workspaceProviderV1, params.Tenant, params.Workspace.Name)
	blockUrl := generators.GenerateBlockStorageURL(storageProviderV1, params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceUrl := generators.GenerateInstanceURL(computeProviderV1, params.Tenant, params.Workspace.Name, params.Instance.Name)
	networkUrl := generators.GenerateNetworkURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.Network.Name)
	gatewayUrl := generators.GenerateInternetGatewayURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.InternetGateway.Name)
	nicUrl := generators.GenerateNicURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.Nic.Name)
	publicIpUrl := generators.GeneratePublicIpURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.PublicIp.Name)
	routeUrl := generators.GenerateRouteTableURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.Network.Name, params.RouteTable.Name)
	subnetUrl := generators.GenerateSubnetURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.Network.Name, params.Subnet.Name)
	groupUrl := generators.GenerateSecurityGroupURL(networkProviderV1, params.Tenant, params.Workspace.Name, params.SecurityGroup.Name)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(workspaceProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	if err := configurator.configureCreateWorkspaceStub(workspaceResponse, workspaceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.configureGetActiveWorkspaceStub(workspaceResponse, workspaceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Network
	networkResponse, err := builders.NewNetworkBuilder().
		Name(params.Network.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.Network.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a network
	if err := configurator.configureCreateNetworkStub(networkResponse, networkUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the created network
	if err := configurator.configureGetActiveNetworkStub(networkResponse, networkUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Update the network
	networkResponse.Spec = *params.Network.UpdatedSpec
	if err := configurator.configureUpdateNetworkStub(networkResponse, networkUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the updated network
	if err := configurator.configureGetActiveNetworkStub(networkResponse, networkUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Internet gateway
	gatewayResponse, err := builders.NewInternetGatewayBuilder().
		Name(params.InternetGateway.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.InternetGateway.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create an internet gateway
	if err := configurator.configureCreateInternetGatewayStub(gatewayResponse, gatewayUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the created internet gateway
	if err := configurator.configureGetActiveInternetGatewayStub(gatewayResponse, gatewayUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Update the internet gateway
	gatewayResponse.Spec = *params.InternetGateway.UpdatedSpec
	if err := configurator.configureUpdateInternetGatewayStub(gatewayResponse, gatewayUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the updated internet gateway
	if err := configurator.configureGetActiveInternetGatewayStub(gatewayResponse, gatewayUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Route table
	routeResponse, err := builders.NewRouteTableBuilder().
		Name(params.RouteTable.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Network(params.Network.Name).Region(params.Region).
		Spec(params.RouteTable.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a route table
	if err := configurator.configureCreateRouteTableStub(routeResponse, routeUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the created route table
	if err := configurator.configureGetActiveRouteTableStub(routeResponse, routeUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Update the route table
	routeResponse.Spec = *params.RouteTable.UpdatedSpec
	if err := configurator.configureUpdateRouteTableStub(routeResponse, routeUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the updated route table
	if err := configurator.configureGetActiveRouteTableStub(routeResponse, routeUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Subnet
	subnetResponse, err := builders.NewSubnetBuilder().
		Name(params.Subnet.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Network(params.Network.Name).Region(params.Region).
		Spec(params.Subnet.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a subnet
	if err := configurator.configureCreateSubnetStub(subnetResponse, subnetUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the created subnet
	if err := configurator.configureGetActiveSubnetStub(subnetResponse, subnetUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Update the subnet
	subnetResponse.Spec = *params.Subnet.UpdatedSpec
	if err := configurator.configureUpdateSubnetStub(subnetResponse, subnetUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the updated subnet
	if err := configurator.configureGetActiveSubnetStub(subnetResponse, subnetUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Public ip
	publicIpResponse, err := builders.NewPublicIpBuilder().
		Name(params.PublicIp.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.PublicIp.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a public ip
	if err := configurator.configureCreatePublicIpStub(publicIpResponse, publicIpUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the created public ip
	if err := configurator.configureGetActivePublicIpStub(publicIpResponse, publicIpUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Update the public ip
	publicIpResponse.Spec = *params.PublicIp.UpdatedSpec
	if err := configurator.configureUpdatePublicIpStub(publicIpResponse, publicIpUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the updated public ip
	if err := configurator.configureGetActivePublicIpStub(publicIpResponse, publicIpUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Nic
	nicResponse, err := builders.NewNicBuilder().
		Name(params.Nic.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.Nic.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a nic
	if err := configurator.configureCreateNicStub(nicResponse, nicUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the created nic
	if err := configurator.configureGetActiveNicStub(nicResponse, nicUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Update the nic
	nicResponse.Spec = *params.Nic.UpdatedSpec
	if err := configurator.configureUpdateNicStub(nicResponse, nicUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the updated nic
	if err := configurator.configureGetActiveNicStub(nicResponse, nicUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Security group
	groupResponse, err := builders.NewSecurityGroupBuilder().
		Name(params.SecurityGroup.Name).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.SecurityGroup.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a security group
	if err := configurator.configureCreateSecurityGroupStub(groupResponse, groupUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the created security group
	if err := configurator.configureGetActiveSecurityGroupStub(groupResponse, groupUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Update the security group
	groupResponse.Spec = *params.SecurityGroup.UpdatedSpec
	if err := configurator.configureUpdateSecurityGroupStub(groupResponse, groupUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the updated security group
	if err := configurator.configureGetActiveSecurityGroupStub(groupResponse, groupUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(params.BlockStorage.Name).
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.BlockStorage.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	if err := configurator.configureCreateBlockStorageStub(blockResponse, blockUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the created block storage
	if err := configurator.configureGetActiveBlockStorageStub(blockResponse, blockUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Instance
	instanceResponse, err := builders.NewInstanceBuilder().
		Name(params.Instance.Name).
		Provider(computeProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Spec(params.Instance.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create an instance
	if err := configurator.configureCreateInstanceStub(instanceResponse, instanceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the created instance
	if err := configurator.configureGetActiveInstanceStub(instanceResponse, instanceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}
	// Delete the instance
	if err := configurator.configureDeleteStub(instanceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted instance
	if err := configurator.configureGetNotFoundStub(instanceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.configureDeleteStub(blockUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configurator.configureGetNotFoundStub(blockUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the security group
	if err := configurator.configureDeleteStub(groupUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted security group
	if err := configurator.configureGetNotFoundStub(groupUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the nic
	if err := configurator.configureDeleteStub(nicUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted nic
	if err := configurator.configureGetNotFoundStub(nicUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the public ip
	if err := configurator.configureDeleteStub(publicIpUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted public ip
	if err := configurator.configureGetNotFoundStub(publicIpUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the subnet
	if err := configurator.configureDeleteStub(subnetUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted subnet
	if err := configurator.configureGetNotFoundStub(subnetUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the route table
	if err := configurator.configureDeleteStub(routeUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted route table
	if err := configurator.configureGetNotFoundStub(routeUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the internet gateway
	if err := configurator.configureDeleteStub(gatewayUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted internet gateway
	if err := configurator.configureGetNotFoundStub(gatewayUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the network
	if err := configurator.configureDeleteStub(networkUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted network
	if err := configurator.configureGetNotFoundStub(networkUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.configureDeleteStub(workspaceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.configureGetNotFoundStub(workspaceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	return configurator.client, nil
}

func ConfigureNetworkListScenarioV1(scenario string, params *NetworkListParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(workspaceProviderV1, params.Tenant, params.Workspace.Name)
	blockUrl := generators.GenerateBlockStorageURL(storageProviderV1, params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceUrl := generators.GenerateInstanceURL(computeProviderV1, params.Tenant, params.Workspace.Name, params.Instance.Name)
	networkListUrl := generators.GenerateNetworkListURL(networkProviderV1, params.Tenant, params.Workspace.Name)
	gatewayListUrl := generators.GenerateInternetGatewayListURL(networkProviderV1, params.Tenant, params.Workspace.Name)
	publicIpListUrl := generators.GeneratePublicIpListURL(networkProviderV1, params.Tenant, params.Workspace.Name)
	nicListUrl := generators.GenerateNicListURL(networkProviderV1, params.Tenant, params.Workspace.Name)
	securityGroupListUrl := generators.GenerateSecurityGroupListURL(networkProviderV1, params.Tenant, params.Workspace.Name)
	skuListUrl := generators.GenerateNetworkSkuListURL(networkProviderV1, params.Tenant)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(params.Workspace.Name).
		Provider(workspaceProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Region(params.Region).
		Labels(params.Workspace.InitialLabels).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	if err := configurator.configureCreateWorkspaceStub(workspaceResponse, workspaceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Create networks
	networkList, err := bulkCreateNetworksStubV1(configurator, params.getBaseParams(), params.Workspace.Name, params.Networks)
	if err != nil {
		return nil, err
	}
	networkResponse, err := builders.NewNetworkIteratorBuilder().
		Provider(storageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).
		Items(networkList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	if err := configurator.configureGetListNetworkStub(networkResponse, networkListUrl, params.getBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	networkResponse.Items = networkList[:1]
	if err := configurator.configureGetListNetworkStub(networkResponse, networkListUrl, params.getBaseParams(), pathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List with Label
	networksWithLabel := func(networkList []schema.Network) []schema.Network {
		var filteredNetworks []schema.Network
		for _, network := range networkList {
			if val, ok := network.Labels[generators.EnvLabel]; ok && val == generators.EnvConformanceLabel {
				filteredNetworks = append(filteredNetworks, network)
			}
		}
		return filteredNetworks
	}
	networkResponse.Items = networksWithLabel(networkList)
	if err := configurator.configureGetListNetworkStub(networkResponse, networkListUrl, params.getBaseParams(), pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	networkResponse.Items = networksWithLabel(networkList)[:1]
	if err := configurator.configureGetListNetworkStub(networkResponse, networkListUrl, params.getBaseParams(), pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Test Network Skus
	// Create skus
	skusList := generateNetworkSkusV1(params.getBaseParams().Tenant)
	skuResponse, err := builders.NewNetworkSkuIteratorBuilder().Provider(storageProviderV1).Tenant(params.Tenant).Items(skusList).Build()
	if err != nil {
		return nil, err
	}

	// List skus
	if err := configurator.configureGetListNetworkSkuStub(skuResponse, skuListUrl, params.getBaseParams(), nil); err != nil {
		return nil, err
	}

	// List skus with limit 1
	if err := configurator.configureGetListNetworkSkuStub(skuResponse, skuListUrl, params.getBaseParams(), pathParamsLimit("1")); err != nil {
		return nil, err
	}

	// Create internet gateways
	gatewayList, err := bulkCreateInternetGatewaysStubV1(configurator, params.getBaseParams(), params.Workspace.Name, params.InternetGateways)
	if err != nil {
		return nil, err
	}
	gatewayResponse, err := builders.NewInternetGatewayIteratorBuilder().
		Provider(storageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).
		Items(gatewayList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	if err := configurator.configureGetListInternetGatewayStub(gatewayResponse, gatewayListUrl, params.getBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	gatewayResponse.Items = gatewayList[:1]
	if err := configurator.configureGetListInternetGatewayStub(gatewayResponse, gatewayListUrl, params.getBaseParams(), pathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List with Label
	gatewayWithLabel := func(gatewayList []schema.InternetGateway) []schema.InternetGateway {
		var filteredGateway []schema.InternetGateway
		for _, gateway := range gatewayList {
			if val, ok := gateway.Labels[generators.EnvLabel]; ok && val == generators.EnvConformanceLabel {
				filteredGateway = append(filteredGateway, gateway)
			}
		}
		return filteredGateway
	}
	gatewayResponse.Items = gatewayWithLabel(gatewayList)
	if err := configurator.configureGetListInternetGatewayStub(gatewayResponse, gatewayListUrl, params.getBaseParams(), pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	gatewayResponse.Items = gatewayWithLabel(gatewayList)[:1]
	if err := configurator.configureGetListInternetGatewayStub(gatewayResponse, gatewayListUrl, params.getBaseParams(), pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Get a network name from the list
	networkName := networkList[0].Metadata.Name

	// Create route tables
	routeTableList, err := bulkCreateRouteTableStubV1(configurator, params.getBaseParams(), params.Workspace.Name, networkName, params.RouteTables)
	if err != nil {
		return nil, err
	}
	routeTableResponse, err := builders.NewRouteTableIteratorBuilder().
		Provider(storageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Network(networkName).
		Items(routeTableList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	routeTableListUrl := generators.GenerateRouteTableListURL(networkProviderV1, params.Tenant, params.Workspace.Name, networkName)
	if err := configurator.configureGetListRouteTableStub(routeTableResponse, routeTableListUrl, params.getBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	routeTableResponse.Items = routeTableList[:1]
	if err := configurator.configureGetListRouteTableStub(routeTableResponse, routeTableListUrl, params.getBaseParams(), pathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List with Label
	routeTableWithLabel := func(routeTableList []schema.RouteTable) []schema.RouteTable {
		var filteredRouteTable []schema.RouteTable
		for _, routeTable := range routeTableList {
			if val, ok := routeTable.Labels[generators.EnvLabel]; ok && val == generators.EnvConformanceLabel {
				filteredRouteTable = append(filteredRouteTable, routeTable)
			}
		}
		return filteredRouteTable
	}

	routeTableResponse.Items = routeTableWithLabel(routeTableList)
	if err := configurator.configureGetListRouteTableStub(routeTableResponse, routeTableListUrl, params.getBaseParams(), pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	routeTableResponse.Items = routeTableWithLabel(routeTableList)[:1]
	if err := configurator.configureGetListRouteTableStub(routeTableResponse, routeTableListUrl, params.getBaseParams(), pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Subnet
	subnetList, err := bulkCreateSubnetsStubV1(configurator, params.getBaseParams(), params.Workspace.Name, networkName, params.Subnets)
	if err != nil {
		return nil, err
	}
	subnetResponse, err := builders.NewSubnetIteratorBuilder().
		Provider(storageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Network(networkName).
		Items(subnetList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	subnetListUrl := generators.GenerateSubnetListURL(networkProviderV1, params.Tenant, params.Workspace.Name, networkName)
	if err := configurator.configureGetListSubnetStub(subnetResponse, subnetListUrl, params.getBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	subnetResponse.Items = subnetList[:1]
	if err := configurator.configureGetListSubnetStub(subnetResponse, subnetListUrl, params.getBaseParams(), pathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List with Label
	subnetWithLabel := func(subnetList []schema.Subnet) []schema.Subnet {
		var filteredSubnet []schema.Subnet
		for _, subnet := range subnetList {
			if val, ok := subnet.Labels[generators.EnvLabel]; ok && val == generators.EnvConformanceLabel {
				filteredSubnet = append(filteredSubnet, subnet)
			}
		}
		return filteredSubnet
	}
	subnetResponse.Items = subnetWithLabel(subnetList)
	if err := configurator.configureGetListSubnetStub(subnetResponse, subnetListUrl, params.getBaseParams(), pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	subnetResponse.Items = subnetWithLabel(subnetList)[:1]
	if err := configurator.configureGetListSubnetStub(subnetResponse, subnetListUrl, params.getBaseParams(), pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Cretae public ips
	publicIpList, err := bulkCreatePublicIpsStubV1(configurator, params.getBaseParams(), params.Workspace.Name, params.PublicIps)
	if err != nil {
		return nil, err
	}
	publicIpResponse, err := builders.NewPublicIpIteratorBuilder().
		Provider(storageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).
		Items(publicIpList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	if err := configurator.configureGetListPublicIpStub(publicIpResponse, publicIpListUrl, params.getBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	publicIpResponse.Items = publicIpList[:1]
	if err := configurator.configureGetListPublicIpStub(publicIpResponse, publicIpListUrl, params.getBaseParams(), pathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List with Label
	publicIpWithLabel := func(publicIpList []schema.PublicIp) []schema.PublicIp {
		var filteredPublicIp []schema.PublicIp
		for _, publicIp := range publicIpList {
			if val, ok := publicIp.Labels[generators.EnvLabel]; ok && val == generators.EnvConformanceLabel {
				filteredPublicIp = append(filteredPublicIp, publicIp)
			}
		}
		return filteredPublicIp
	}
	publicIpResponse.Items = publicIpWithLabel(publicIpList)
	if err := configurator.configureGetListPublicIpStub(publicIpResponse, publicIpListUrl, params.getBaseParams(), pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	publicIpResponse.Items = publicIpWithLabel(publicIpList)[:1]
	if err := configurator.configureGetListPublicIpStub(publicIpResponse, publicIpListUrl, params.getBaseParams(), pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Create nics
	nicList, err := bulkCreateNicsStubV1(configurator, params.getBaseParams(), params.Workspace.Name, params.Nics)
	if err != nil {
		return nil, err
	}
	nicResponse, err := builders.NewNicIteratorBuilder().
		Provider(storageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).
		Items(nicList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	if err := configurator.configureGetListNicStub(nicResponse, nicListUrl, params.getBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	nicResponse.Items = nicList[:1]
	if err := configurator.configureGetListNicStub(nicResponse, nicListUrl, params.getBaseParams(), pathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List with Label
	nicWithLabel := func(nicList []schema.Nic) []schema.Nic {
		var filteredNic []schema.Nic
		for _, nic := range nicList {
			if val, ok := nic.Labels[generators.EnvLabel]; ok && val == generators.EnvConformanceLabel {
				filteredNic = append(filteredNic, nic)
			}
		}
		return filteredNic
	}
	nicResponse.Items = nicWithLabel(nicList)
	if err := configurator.configureGetListNicStub(nicResponse, nicListUrl, params.getBaseParams(), pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	nicResponse.Items = nicWithLabel(nicList)[:1]
	if err := configurator.configureGetListNicStub(nicResponse, nicListUrl, params.getBaseParams(), pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Create security groups
	securityGroupList, err := bulkCreateSecurityGroupsStubV1(configurator, params.getBaseParams(), params.Workspace.Name, params.SecurityGroups)
	if err != nil {
		return nil, err
	}
	securityGroupResponse, err := builders.NewSecurityGroupIteratorBuilder().
		Provider(storageProviderV1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).
		Items(securityGroupList).
		Build()
	if err != nil {
		return nil, err
	}

	// List
	if err := configurator.configureGetListSecurityGroupStub(securityGroupResponse, securityGroupListUrl, params.getBaseParams(), nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	securityGroupResponse.Items = securityGroupList[:1]
	if err := configurator.configureGetListSecurityGroupStub(securityGroupResponse, securityGroupListUrl, params.getBaseParams(), pathParamsLimit("1")); err != nil {
		return nil, err
	}

	// List with Label
	secGroupWithLabel := func(securityGroupList []schema.SecurityGroup) []schema.SecurityGroup {
		var filteredSecurity []schema.SecurityGroup
		for _, sec := range securityGroupList {
			if val, ok := sec.Labels[generators.EnvLabel]; ok && val == generators.EnvConformanceLabel {
				filteredSecurity = append(filteredSecurity, sec)
			}
		}
		return filteredSecurity
	}

	securityGroupResponse.Items = secGroupWithLabel(securityGroupList)
	if err := configurator.configureGetListSecurityGroupStub(securityGroupResponse, securityGroupListUrl, params.getBaseParams(), pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	securityGroupResponse.Items = secGroupWithLabel(securityGroupList)[:1]
	if err := configurator.configureGetListSecurityGroupStub(securityGroupResponse, securityGroupListUrl, params.getBaseParams(), pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(params.BlockStorage.Name).
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Labels(params.BlockStorage.InitialLabels).
		Spec(params.BlockStorage.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	if err := configurator.configureCreateBlockStorageStub(blockResponse, blockUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Instance
	instanceResponse, err := builders.NewInstanceBuilder().
		Name(params.Instance.Name).
		Provider(computeProviderV1).ApiVersion(apiVersion1).
		Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
		Labels(params.Instance.InitialLabels).
		Spec(params.Instance.InitialSpec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create an instance
	if err := configurator.configureCreateInstanceStub(instanceResponse, instanceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the created instance
	if err := configurator.configureGetActiveInstanceStub(instanceResponse, instanceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete

	// Delete the instance
	if err := configurator.configureDeleteStub(instanceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted instance
	if err := configurator.configureGetNotFoundStub(instanceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.configureDeleteStub(blockUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configurator.configureGetNotFoundStub(blockUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Delete the security groups
	for _, securityGroup := range params.SecurityGroups {
		securityGroupUrl := generators.GenerateSecurityGroupURL(networkProviderV1, params.Tenant, params.Workspace.Name, securityGroup.Name)

		// Delete the security group
		if err := configurator.configureDeleteStub(securityGroupUrl, params.getBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted security group
		if err := configurator.configureGetNotFoundStub(securityGroupUrl, params.getBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete the nics
	for _, nic := range params.Nics {
		nicUrl := generators.GenerateNicURL(networkProviderV1, params.Tenant, params.Workspace.Name, nic.Name)

		// Delete the nic
		if err := configurator.configureDeleteStub(nicUrl, params.getBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted nic
		if err := configurator.configureGetNotFoundStub(nicUrl, params.getBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete the public ips
	for _, publicIp := range params.PublicIps {
		publicIpUrl := generators.GeneratePublicIpURL(networkProviderV1, params.Tenant, params.Workspace.Name, publicIp.Name)

		// Delete the public ip
		if err := configurator.configureDeleteStub(publicIpUrl, params.getBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted public ip
		if err := configurator.configureGetNotFoundStub(publicIpUrl, params.getBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete the subnets
	for _, subnet := range params.Subnets {
		subnetUrl := generators.GenerateSubnetURL(networkProviderV1, params.Tenant, params.Workspace.Name, networkName, subnet.Name)

		// Delete the subnet
		if err := configurator.configureDeleteStub(subnetUrl, params.getBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted subnet
		if err := configurator.configureGetNotFoundStub(subnetUrl, params.getBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete the route tables
	for _, routeTable := range params.RouteTables {
		routeTableUrl := generators.GenerateRouteTableURL(networkProviderV1, params.Tenant, params.Workspace.Name, networkName, routeTable.Name)

		// Delete the route table
		if err := configurator.configureDeleteStub(routeTableUrl, params.getBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted route table
		if err := configurator.configureGetNotFoundStub(routeTableUrl, params.getBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete the internet gateways
	for _, gateway := range params.InternetGateways {
		gatewayUrl := generators.GenerateInternetGatewayURL(networkProviderV1, params.Tenant, params.Workspace.Name, gateway.Name)

		// Delete the internet gateway
		if err := configurator.configureDeleteStub(gatewayUrl, params.getBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted internet gateway
		if err := configurator.configureGetNotFoundStub(gatewayUrl, params.getBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete the networks
	for _, network := range params.Networks {
		networkUrl := generators.GenerateInternetGatewayURL(networkProviderV1, params.Tenant, params.Workspace.Name, network.Name)

		// Delete the network
		if err := configurator.configureDeleteStub(networkUrl, params.getBaseParams()); err != nil {
			return nil, err
		}

		// Get the deleted network
		if err := configurator.configureGetNotFoundStub(networkUrl, params.getBaseParams()); err != nil {
			return nil, err
		}
	}

	// Delete the workspace
	if err := configurator.configureDeleteStub(workspaceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.configureGetNotFoundStub(workspaceUrl, params.getBaseParams()); err != nil {
		return nil, err
	}
	return configurator.client, nil
}
