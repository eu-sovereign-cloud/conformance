package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	network "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.network.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigNetworkLifecycleScenarioV1(scenario string, params *NetworkParamsV1) (*wiremock.Client, error) {
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
	if err := configurator.configureCreateWorkspaceStub(workspaceResponse, workspaceUrl, params); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.configureGetActiveWorkspaceStub(workspaceResponse, workspaceUrl, params); err != nil {
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
	if err := configurator.configureCreateNetworkStub(networkResponse, networkUrl, params); err != nil {
		return nil, err
	}

	// Get the created network
	if err := configurator.configureGetActiveNetworkStub(networkResponse, networkUrl, params); err != nil {
		return nil, err
	}

	// Update the network
	networkResponse.Spec = *params.Network.UpdatedSpec
	if err := configurator.configureUpdateNetworkStub(networkResponse, networkUrl, params); err != nil {
		return nil, err
	}

	// Get the updated network
	if err := configurator.configureGetActiveNetworkStub(networkResponse, networkUrl, params); err != nil {
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
	if err := configurator.configureCreateInternetGatewayStub(gatewayResponse, gatewayUrl, params); err != nil {
		return nil, err
	}

	// Get the created internet gateway
	if err := configurator.configureGetActiveInternetGatewayStub(gatewayResponse, gatewayUrl, params); err != nil {
		return nil, err
	}

	// Update the internet gateway
	gatewayResponse.Spec = *params.InternetGateway.UpdatedSpec
	if err := configurator.configureUpdateInternetGatewayStub(gatewayResponse, gatewayUrl, params); err != nil {
		return nil, err
	}

	// Get the updated internet gateway
	if err := configurator.configureGetActiveInternetGatewayStub(gatewayResponse, gatewayUrl, params); err != nil {
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
	if err := configurator.configureCreateRouteTableStub(routeResponse, routeUrl, params); err != nil {
		return nil, err
	}

	// Get the created route table
	if err := configurator.configureGetActiveRouteTableStub(routeResponse, routeUrl, params); err != nil {
		return nil, err
	}

	// Update the route table
	routeResponse.Spec = *params.RouteTable.UpdatedSpec
	if err := configurator.configureUpdateRouteTableStub(routeResponse, routeUrl, params); err != nil {
		return nil, err
	}

	// Get the updated route table
	if err := configurator.configureGetActiveRouteTableStub(routeResponse, routeUrl, params); err != nil {
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
	if err := configurator.configureCreateSubnetStub(subnetResponse, subnetUrl, params); err != nil {
		return nil, err
	}

	// Get the created subnet
	if err := configurator.configureGetActiveSubnetStub(subnetResponse, subnetUrl, params); err != nil {
		return nil, err
	}

	// Update the subnet
	subnetResponse.Spec = *params.Subnet.UpdatedSpec
	if err := configurator.configureUpdateSubnetStub(subnetResponse, subnetUrl, params); err != nil {
		return nil, err
	}

	// Get the updated subnet
	if err := configurator.configureGetActiveSubnetStub(subnetResponse, subnetUrl, params); err != nil {
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
	if err := configurator.configureCreatePublicIpStub(publicIpResponse, publicIpUrl, params); err != nil {
		return nil, err
	}

	// Get the created public ip
	if err := configurator.configureGetActivePublicIpStub(publicIpResponse, publicIpUrl, params); err != nil {
		return nil, err
	}

	// Update the public ip
	publicIpResponse.Spec = *params.PublicIp.UpdatedSpec
	if err := configurator.configureUpdatePublicIpStub(publicIpResponse, publicIpUrl, params); err != nil {
		return nil, err
	}

	// Get the updated public ip
	if err := configurator.configureGetActivePublicIpStub(publicIpResponse, publicIpUrl, params); err != nil {
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
	if err := configurator.configureCreateNicStub(nicResponse, nicUrl, params); err != nil {
		return nil, err
	}

	// Get the created nic
	if err := configurator.configureGetActiveNicStub(nicResponse, nicUrl, params); err != nil {
		return nil, err
	}

	// Update the nic
	nicResponse.Spec = *params.Nic.UpdatedSpec
	if err := configurator.configureUpdateNicStub(nicResponse, nicUrl, params); err != nil {
		return nil, err
	}

	// Get the updated nic
	if err := configurator.configureGetActiveNicStub(nicResponse, nicUrl, params); err != nil {
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
	if err := configurator.configureCreateSecurityGroupStub(groupResponse, groupUrl, params); err != nil {
		return nil, err
	}

	// Get the created security group
	if err := configurator.configureGetActiveSecurityGroupStub(groupResponse, groupUrl, params); err != nil {
		return nil, err
	}

	// Update the security group
	groupResponse.Spec = *params.SecurityGroup.UpdatedSpec
	if err := configurator.configureUpdateSecurityGroupStub(groupResponse, groupUrl, params); err != nil {
		return nil, err
	}

	// Get the updated security group
	if err := configurator.configureGetActiveSecurityGroupStub(groupResponse, groupUrl, params); err != nil {
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
	if err := configurator.configureCreateBlockStorageStub(blockResponse, blockUrl, params); err != nil {
		return nil, err
	}

	// Get the created block storage
	if err := configurator.configureGetActiveBlockStorageStub(blockResponse, blockUrl, params); err != nil {
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
	if err := configurator.configureCreateInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}

	// Get the created instance
	if err := configurator.configureGetActiveInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}
	// Delete the instance
	if err := configurator.configureDeleteStub(instanceUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted instance
	if err := configurator.configureGetNotFoundStub(instanceUrl, params); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.configureDeleteStub(blockUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configurator.configureGetNotFoundStub(blockUrl, params); err != nil {
		return nil, err
	}

	// Delete the security group
	if err := configurator.configureDeleteStub(groupUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted security group
	if err := configurator.configureGetNotFoundStub(groupUrl, params); err != nil {
		return nil, err
	}

	// Delete the nic
	if err := configurator.configureDeleteStub(nicUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted nic
	if err := configurator.configureGetNotFoundStub(nicUrl, params); err != nil {
		return nil, err
	}

	// Delete the public ip
	if err := configurator.configureDeleteStub(publicIpUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted public ip
	if err := configurator.configureGetNotFoundStub(publicIpUrl, params); err != nil {
		return nil, err
	}

	// Delete the subnet
	if err := configurator.configureDeleteStub(subnetUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted subnet
	if err := configurator.configureGetNotFoundStub(subnetUrl, params); err != nil {
		return nil, err
	}

	// Delete the route table
	if err := configurator.configureDeleteStub(routeUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted route table
	if err := configurator.configureGetNotFoundStub(routeUrl, params); err != nil {
		return nil, err
	}

	// Delete the internet gateway
	if err := configurator.configureDeleteStub(gatewayUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted internet gateway
	if err := configurator.configureGetNotFoundStub(gatewayUrl, params); err != nil {
		return nil, err
	}

	// Delete the network
	if err := configurator.configureDeleteStub(networkUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted network
	if err := configurator.configureGetNotFoundStub(networkUrl, params); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.configureDeleteStub(workspaceUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.configureGetNotFoundStub(workspaceUrl, params); err != nil {
		return nil, err
	}

	return configurator.client, nil
}

func ConfigNetworkListLifecycleScenarioV1(scenario string, params *NetworkListParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := newScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(workspaceProviderV1, params.Tenant, params.Workspace.Name)
	blockUrl := generators.GenerateBlockStorageURL(storageProviderV1, params.Tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceUrl := generators.GenerateInstanceURL(computeProviderV1, params.Tenant, params.Workspace.Name, params.Instance.Name)

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
	if err := configurator.configureCreateWorkspaceStub(workspaceResponse, workspaceUrl, params); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.configureGetActiveWorkspaceStub(workspaceResponse, workspaceUrl, params); err != nil {
		return nil, err
	}

	// Network
	var networkList []schema.Network
	for _, network := range *params.Network {
		networkUrl := generators.GenerateNetworkURL(networkProviderV1, params.Tenant, params.Workspace.Name, network.Name)

		networkResponse, err := builders.NewNetworkBuilder().
			Name(network.Name).
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
			Spec(network.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a network
		if err := configurator.configureCreateNetworkStub(networkResponse, networkUrl, params); err != nil {
			return nil, err
		}
		networkList = append(networkList, *networkResponse)
	}
	// List
	networkListUrl := generators.GenerateNetworkListURL(networkProviderV1, params.Tenant, params.Workspace.Name)
	networkResponse := &network.NetworkIterator{
		Metadata: schema.ResponseMetadata{
			Provider: networkProviderV1,
			Resource: generators.GenerateNetworkListResource(params.Tenant, params.Workspace.Name),
			Verb:     http.MethodGet,
		},
		Items: networkList,
	}

	if err := configurator.configureGetListNetworkStub(networkResponse, networkListUrl, params, nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	networkResponse.Items = networkList[:1]
	if err := configurator.configureGetListNetworkStub(networkResponse, networkListUrl, params, pathParamsLimit("1")); err != nil {
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
	if err := configurator.configureGetListNetworkStub(networkResponse, networkListUrl, params, pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	networkResponse.Items = networksWithLabel(networkList)[:1]
	if err := configurator.configureGetListNetworkStub(networkResponse, networkListUrl, params, pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Internet Gateway
	var gatewayList []schema.InternetGateway
	for _, gateway := range *params.InternetGateway {
		gatewayUrl := generators.GenerateInternetGatewayURL(networkProviderV1, params.Tenant, params.Workspace.Name, gateway.Name)
		gatewayResponse, err := builders.NewInternetGatewayBuilder().
			Name(gateway.Name).
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
			Spec(gateway.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a InternetGateway
		if err := configurator.configureCreateInternetGatewayStub(gatewayResponse, gatewayUrl, params); err != nil {
			return nil, err
		}
		gatewayList = append(gatewayList, *gatewayResponse)
	}

	// List
	gatewayListUrl := generators.GenerateInternetGatewayListURL(networkProviderV1, params.Tenant, params.Workspace.Name)
	gatewayResource := generators.GenerateInternetGatewayListResource(params.Tenant, params.Workspace.Name)
	gatewayResponse := &network.InternetGatewayIterator{
		Metadata: schema.ResponseMetadata{
			Provider: networkProviderV1,
			Resource: gatewayResource,
			Verb:     http.MethodGet,
		},
		Items: gatewayList,
	}
	// List
	if err := configurator.configureGetListInternetGatewayStub(gatewayResponse, gatewayListUrl, params, nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	gatewayResponse.Items = gatewayList[:1]
	if err := configurator.configureGetListInternetGatewayStub(gatewayResponse, gatewayListUrl, params, pathParamsLimit("1")); err != nil {
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
	if err := configurator.configureGetListInternetGatewayStub(gatewayResponse, gatewayListUrl, params, pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	gatewayResponse.Items = gatewayWithLabel(gatewayList)[:1]
	if err := configurator.configureGetListInternetGatewayStub(gatewayResponse, gatewayListUrl, params, pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Route Table
	var routeTableList []schema.RouteTable
	networkGeneric := networkList[:1]
	for _, routeTable := range *params.RouteTable {
		routeTableUrl := generators.GenerateRouteTableURL(networkProviderV1, params.Tenant, params.Workspace.Name, networkGeneric[0].Metadata.Name, routeTable.Name)
		routeTableResponse, err := builders.NewRouteTableBuilder().
			Name(routeTable.Name).
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(params.Tenant).Workspace(params.Workspace.Name).Network(networkGeneric[0].Metadata.Name).Region(params.Region).
			Spec(routeTable.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a RouteTable
		if err := configurator.configureCreateRouteTableStub(routeTableResponse, routeTableUrl, params); err != nil {
			return nil, err
		}
		routeTableList = append(routeTableList, *routeTableResponse)
	}

	routeTableListUrl := generators.GenerateRouteTableListURL(networkProviderV1, params.Tenant, params.Workspace.Name, networkGeneric[0].Metadata.Name)
	routeTableResource := generators.GenerateRouteTableListResource(params.Tenant, params.Workspace.Name, networkGeneric[0].Metadata.Name)
	routeTableResponse := &network.RouteTableIterator{
		Metadata: schema.ResponseMetadata{
			Provider: networkProviderV1,
			Resource: routeTableResource,
			Verb:     http.MethodGet,
		},
		Items: routeTableList,
	}

	//List
	if err := configurator.configureGetListRouteTableStub(routeTableResponse, routeTableListUrl, params, nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	routeTableResponse.Items = routeTableList[:1]
	if err := configurator.configureGetListRouteTableStub(routeTableResponse, routeTableListUrl, params, pathParamsLimit("1")); err != nil {
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
	if err := configurator.configureGetListRouteTableStub(routeTableResponse, routeTableListUrl, params, pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	routeTableResponse.Items = routeTableWithLabel(routeTableList)[:1]
	if err := configurator.configureGetListRouteTableStub(routeTableResponse, routeTableListUrl, params, pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Subnet
	var subnetList []schema.Subnet
	for _, subnet := range *params.Subnet {
		subnetUrl := generators.GenerateSubnetURL(networkProviderV1, params.Tenant, params.Workspace.Name, networkGeneric[0].Metadata.Name, subnet.Name)
		subnetResponse, err := builders.NewSubnetBuilder().
			Name(subnet.Name).
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(params.Tenant).Workspace(params.Workspace.Name).Network(networkGeneric[0].Metadata.Name).Region(params.Region).
			Spec(subnet.InitialSpec).
			Build()

		if err != nil {
			return nil, err
		}
		// Create a RouteTable
		if err := configurator.configureCreateSubnetStub(subnetResponse, subnetUrl, params); err != nil {
			return nil, err
		}
		subnetList = append(subnetList, *subnetResponse)
	}

	subnetListUrl := generators.GenerateSubnetListURL(networkProviderV1, params.Tenant, params.Workspace.Name, networkGeneric[0].Metadata.Name)
	subnetResource := generators.GenerateSubnetListResource(params.Tenant, params.Workspace.Name, networkGeneric[0].Metadata.Name)
	subnetResponse := &network.SubnetIterator{
		Metadata: schema.ResponseMetadata{
			Provider: networkProviderV1,
			Resource: subnetResource,
			Verb:     http.MethodGet,
		},
		Items: subnetList,
	}
	//List
	if err := configurator.configureGetListSubnetStub(subnetResponse, subnetListUrl, params, nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	subnetResponse.Items = subnetList[:1]
	if err := configurator.configureGetListSubnetStub(subnetResponse, subnetListUrl, params, pathParamsLimit("1")); err != nil {
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
	if err := configurator.configureGetListSubnetStub(subnetResponse, subnetListUrl, params, pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	subnetResponse.Items = subnetWithLabel(subnetList)[:1]
	if err := configurator.configureGetListSubnetStub(subnetResponse, subnetListUrl, params, pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// Public IP
	var publicIpList []schema.PublicIp
	for _, publicIp := range *params.PublicIp {
		publicIpUrl := generators.GeneratePublicIpURL(networkProviderV1, params.Tenant, params.Workspace.Name, publicIp.Name)
		publicIpResponse, err := builders.NewPublicIpBuilder().
			Name(publicIp.Name).
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
			Spec(publicIp.InitialSpec).
			Build()

		if err != nil {
			return nil, err
		}
		// Create a Public IP
		if err := configurator.configureCreatePublicIpStub(publicIpResponse, publicIpUrl, params); err != nil {
			return nil, err
		}
		publicIpList = append(publicIpList, *publicIpResponse)
	}

	//List
	publicIpListUrl := generators.GeneratePublicIpListURL(networkProviderV1, params.Tenant, params.Workspace.Name)
	publicIpResource := generators.GeneratePublicIpListResource(params.Tenant, params.Workspace.Name)
	publicIpResponse := &network.PublicIpIterator{
		Metadata: schema.ResponseMetadata{
			Provider: networkProviderV1,
			Resource: publicIpResource,
			Verb:     http.MethodGet,
		},
		Items: publicIpList,
	}

	if err := configurator.configureGetListPublicIpStub(publicIpResponse, publicIpListUrl, params, nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	publicIpResponse.Items = publicIpList[:1]
	if err := configurator.configureGetListPublicIpStub(publicIpResponse, publicIpListUrl, params, pathParamsLimit("1")); err != nil {
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
	if err := configurator.configureGetListPublicIpStub(publicIpResponse, publicIpListUrl, params, pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	publicIpResponse.Items = publicIpWithLabel(publicIpList)[:1]
	if err := configurator.configureGetListPublicIpStub(publicIpResponse, publicIpListUrl, params, pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// NIC
	var nicList []schema.Nic
	for _, nic := range *params.Nic {
		nicUrl := generators.GenerateNicURL(networkProviderV1, params.Tenant, params.Workspace.Name, nic.Name)
		nicResponse, err := builders.NewNicBuilder().
			Name(nic.Name).
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
			Spec(nic.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a Nic
		if err := configurator.configureCreateNicStub(nicResponse, nicUrl, params); err != nil {
			return nil, err
		}
		nicList = append(nicList, *nicResponse)
	}

	//List
	nicListUrl := generators.GenerateNicListURL(networkProviderV1, params.Tenant, params.Workspace.Name)
	nicResource := generators.GenerateNicListResource(params.Tenant, params.Workspace.Name)
	nicResponse := &network.NicIterator{
		Metadata: schema.ResponseMetadata{
			Provider: networkProviderV1,
			Resource: nicResource,
			Verb:     http.MethodGet,
		},
		Items: nicList,
	}

	if err := configurator.configureGetListNicStub(nicResponse, nicListUrl, params, nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	nicResponse.Items = nicList[:1]
	if err := configurator.configureGetListNicStub(nicResponse, nicListUrl, params, pathParamsLimit("1")); err != nil {
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
	if err := configurator.configureGetListNicStub(nicResponse, nicListUrl, params, pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	nicResponse.Items = nicWithLabel(nicList)[:1]
	if err := configurator.configureGetListNicStub(nicResponse, nicListUrl, params, pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// SecurityGroup
	var securityGroupList []schema.SecurityGroup
	for _, securityGroup := range *params.SecurityGroup {
		securityGroupUrl := generators.GenerateSecurityGroupURL(networkProviderV1, params.Tenant, params.Workspace.Name, securityGroup.Name)
		securityGroupResponse, err := builders.NewSecurityGroupBuilder().
			Name(securityGroup.Name).
			Provider(networkProviderV1).ApiVersion(apiVersion1).
			Tenant(params.Tenant).Workspace(params.Workspace.Name).Region(params.Region).
			Spec(securityGroup.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}
		// Create a SecurityGroup
		if err := configurator.configureCreateSecurityGroupStub(securityGroupResponse, securityGroupUrl, params); err != nil {
			return nil, err
		}
		securityGroupList = append(securityGroupList, *securityGroupResponse)
	}
	//List
	securityGroupListUrl := generators.GenerateSecurityGroupListURL(networkProviderV1, params.Tenant, params.Workspace.Name)
	securityGroupResource := generators.GenerateSecurityGroupListResource(params.Tenant, params.Workspace.Name)
	securityGroupResponse := &network.SecurityGroupIterator{
		Metadata: schema.ResponseMetadata{
			Provider: networkProviderV1,
			Resource: securityGroupResource,
			Verb:     http.MethodGet,
		},
		Items: securityGroupList,
	}
	if err := configurator.configureGetListSecurityGroupStub(securityGroupResponse, securityGroupListUrl, params, nil); err != nil {
		return nil, err
	}

	// List with Limit 1
	securityGroupResponse.Items = securityGroupList[:1]
	if err := configurator.configureGetListSecurityGroupStub(securityGroupResponse, nicListUrl, params, pathParamsLimit("1")); err != nil {
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
	if err := configurator.configureGetListSecurityGroupStub(securityGroupResponse, nicListUrl, params, pathParamsLabel(generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
		return nil, err
	}

	// List with Limit and Label
	securityGroupResponse.Items = secGroupWithLabel(securityGroupList)[:1]
	if err := configurator.configureGetListSecurityGroupStub(securityGroupResponse, nicListUrl, params, pathParamsLimitAndLabel("1", generators.EnvLabel, generators.EnvConformanceLabel)); err != nil {
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
	if err := configurator.configureCreateBlockStorageStub(blockResponse, blockUrl, params); err != nil {
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
	if err := configurator.configureCreateInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}

	// Get the created instance
	if err := configurator.configureGetActiveInstanceStub(instanceResponse, instanceUrl, params); err != nil {
		return nil, err
	}

	// Delete

	// Delete the instance
	if err := configurator.configureDeleteStub(instanceUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted instance
	if err := configurator.configureGetNotFoundStub(instanceUrl, params); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.configureDeleteStub(blockUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configurator.configureGetNotFoundStub(blockUrl, params); err != nil {
		return nil, err
	}
	for _, securityGroup := range *params.SecurityGroup {
		securityGroupUrl := generators.GenerateSecurityGroupURL(networkProviderV1, params.Tenant, params.Workspace.Name, securityGroup.Name)
		// Delete the security group
		if err := configurator.configureDeleteStub(securityGroupUrl, params); err != nil {
			return nil, err
		}

		// Get the deleted security group
		if err := configurator.configureGetNotFoundStub(securityGroupUrl, params); err != nil {
			return nil, err
		}
	}

	// Delete NIC
	for _, nic := range *params.Nic {
		nicUrl := generators.GenerateNicURL(networkProviderV1, params.Tenant, params.Workspace.Name, nic.Name)
		// Delete NIC
		if err := configurator.configureDeleteStub(nicUrl, params); err != nil {
			return nil, err
		}

		// Get the deleted NIC
		if err := configurator.configureGetNotFoundStub(nicUrl, params); err != nil {
			return nil, err
		}
	}

	// Public IP
	for _, publicIp := range *params.PublicIp {
		publicIpUrl := generators.GeneratePublicIpURL(networkProviderV1, params.Tenant, params.Workspace.Name, publicIp.Name)

		// Delete Public IP
		if err := configurator.configureDeleteStub(publicIpUrl, params); err != nil {
			return nil, err
		}

		// Get the deleted Public IP
		if err := configurator.configureGetNotFoundStub(publicIpUrl, params); err != nil {
			return nil, err
		}
	}

	// Subnet
	for _, subnet := range *params.Subnet {
		subnetUrl := generators.GenerateSubnetURL(networkProviderV1, params.Tenant, params.Workspace.Name, networkGeneric[0].Metadata.Name, subnet.Name)

		// Delete Subnet
		if err := configurator.configureDeleteStub(subnetUrl, params); err != nil {
			return nil, err
		}

		// Get the deleted Subnet
		if err := configurator.configureGetNotFoundStub(subnetUrl, params); err != nil {
			return nil, err
		}
	}

	// Route Table
	for _, routeTable := range *params.RouteTable {
		routeTableUrl := generators.GenerateRouteTableURL(networkProviderV1, params.Tenant, params.Workspace.Name, networkGeneric[0].Metadata.Name, routeTable.Name)

		// Delete Route Table
		if err := configurator.configureDeleteStub(routeTableUrl, params); err != nil {
			return nil, err
		}

		// Get the deleted Route Table
		if err := configurator.configureGetNotFoundStub(routeTableUrl, params); err != nil {
			return nil, err
		}
	}

	// Internet Gateway
	for _, gateway := range *params.InternetGateway {
		gatewayUrl := generators.GenerateInternetGatewayURL(networkProviderV1, params.Tenant, params.Workspace.Name, gateway.Name)

		// Delete Internet Gateway
		if err := configurator.configureDeleteStub(gatewayUrl, params); err != nil {
			return nil, err
		}

		// Get the deleted Internet Gateway
		if err := configurator.configureGetNotFoundStub(gatewayUrl, params); err != nil {
			return nil, err
		}
	}

	// Network
	for _, network := range *params.Network {
		networkUrl := generators.GenerateInternetGatewayURL(networkProviderV1, params.Tenant, params.Workspace.Name, network.Name)

		// Delete Network
		if err := configurator.configureDeleteStub(networkUrl, params); err != nil {
			return nil, err
		}

		// Get the deleted Network
		if err := configurator.configureGetNotFoundStub(networkUrl, params); err != nil {
			return nil, err
		}
	}

	// Delete the workspace
	if err := configurator.configureDeleteStub(workspaceUrl, params); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.configureGetNotFoundStub(workspaceUrl, params); err != nil {
		return nil, err
	}
	return configurator.client, nil
}
