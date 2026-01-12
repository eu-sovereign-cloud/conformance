package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"

	"github.com/wiremock/go-wiremock"
)

func ConfigureLifecycleScenarioV1(scenario string, mockParams *mock.MockParams, suiteParams *params.NetworkLifeCycleParamsV1) (*wiremock.Client, error) {
	scenarios.LogScenarioMocking(scenario)

	workspace := *suiteParams.Workspace
	blockStorage := *suiteParams.BlockStorage
	instance := *suiteParams.Instance
	network := *suiteParams.NetworkInitial
	gateway := *suiteParams.InternetGatewayInitial
	nic := *suiteParams.NicInitial
	publicIp := *suiteParams.PublicIpInitial
	routeTable := *suiteParams.RouteTableInitial
	subnet := *suiteParams.SubnetInitial
	securityGroup := *suiteParams.SecurityGroupInitial
	configurator, err := stubs.NewStubConfigurator(scenario, mockParams)
	if err != nil {
		return nil, err
	}

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(constants.WorkspaceProviderV1, workspace.Metadata.Tenant, workspace.Metadata.Name)
	blockUrl := generators.GenerateBlockStorageURL(constants.StorageProviderV1, blockStorage.Metadata.Tenant, blockStorage.Metadata.Workspace, blockStorage.Metadata.Name)
	instanceUrl := generators.GenerateInstanceURL(constants.ComputeProviderV1, instance.Metadata.Tenant, instance.Metadata.Workspace, instance.Metadata.Name)
	networkUrl := generators.GenerateNetworkURL(constants.NetworkProviderV1, network.Metadata.Tenant, network.Metadata.Workspace, network.Metadata.Name)
	gatewayUrl := generators.GenerateInternetGatewayURL(constants.NetworkProviderV1, gateway.Metadata.Tenant, gateway.Metadata.Workspace, gateway.Metadata.Name)
	nicUrl := generators.GenerateNicURL(constants.NetworkProviderV1, nic.Metadata.Tenant, nic.Metadata.Workspace, nic.Metadata.Name)
	publicIpUrl := generators.GeneratePublicIpURL(constants.NetworkProviderV1, publicIp.Metadata.Tenant, publicIp.Metadata.Workspace, publicIp.Metadata.Name)
	routeUrl := generators.GenerateRouteTableURL(constants.NetworkProviderV1, routeTable.Metadata.Tenant, routeTable.Metadata.Workspace, routeTable.Metadata.Network, routeTable.Metadata.Name)
	subnetUrl := generators.GenerateSubnetURL(constants.NetworkProviderV1, subnet.Metadata.Tenant, subnet.Metadata.Workspace, subnet.Metadata.Network, subnet.Metadata.Name)
	groupUrl := generators.GenerateSecurityGroupURL(constants.NetworkProviderV1, securityGroup.Metadata.Tenant, securityGroup.Metadata.Workspace, securityGroup.Metadata.Name)

	// Workspace
	workspaceResponse, err := builders.NewWorkspaceBuilder().
		Name(workspace.Metadata.Name).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(workspace.Metadata.Tenant).Region(workspace.Metadata.Region).
		Labels(workspace.Labels).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspaceResponse, workspaceUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspaceResponse, workspaceUrl, mockParams); err != nil {
		return nil, err
	}

	// Network
	networkResponse, err := builders.NewNetworkBuilder().
		Name(network.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(network.Metadata.Tenant).Workspace(network.Metadata.Workspace).Region(network.Metadata.Region).
		Spec(&network.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a network
	if err := configurator.ConfigureCreateNetworkStub(networkResponse, networkUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created network
	if err := configurator.ConfigureGetActiveNetworkStub(networkResponse, networkUrl, mockParams); err != nil {
		return nil, err
	}

	// Update the network
	networkResponse.Spec = *&suiteParams.NetworkUpdated.Spec
	if err := configurator.ConfigureUpdateNetworkStub(networkResponse, networkUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the updated network
	if err := configurator.ConfigureGetActiveNetworkStub(networkResponse, networkUrl, mockParams); err != nil {
		return nil, err
	}

	// Internet gateway
	gatewayResponse, err := builders.NewInternetGatewayBuilder().
		Name(gateway.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(gateway.Metadata.Tenant).Workspace(gateway.Metadata.Workspace).Region(gateway.Metadata.Region).
		Spec(&gateway.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create an internet gateway
	if err := configurator.ConfigureCreateInternetGatewayStub(gatewayResponse, gatewayUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created internet gateway
	if err := configurator.ConfigureGetActiveInternetGatewayStub(gatewayResponse, gatewayUrl, mockParams); err != nil {
		return nil, err
	}

	// Update the internet gateway
	gatewayResponse.Spec = *&suiteParams.InternetGatewayUpdated.Spec
	if err := configurator.ConfigureUpdateInternetGatewayStub(gatewayResponse, gatewayUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the updated internet gateway
	if err := configurator.ConfigureGetActiveInternetGatewayStub(gatewayResponse, gatewayUrl, mockParams); err != nil {
		return nil, err
	}

	// Route table
	routeResponse, err := builders.NewRouteTableBuilder().
		Name(routeTable.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(routeTable.Metadata.Tenant).Workspace(routeTable.Metadata.Workspace).Network(routeTable.Metadata.Network).Region(routeTable.Metadata.Region).
		Spec(&routeTable.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a route table
	if err := configurator.ConfigureCreateRouteTableStub(routeResponse, routeUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created route table
	if err := configurator.ConfigureGetActiveRouteTableStub(routeResponse, routeUrl, mockParams); err != nil {
		return nil, err
	}

	// Update the route table
	routeResponse.Spec = *&suiteParams.RouteTableUpdated.Spec
	if err := configurator.ConfigureUpdateRouteTableStub(routeResponse, routeUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the updated route table
	if err := configurator.ConfigureGetActiveRouteTableStub(routeResponse, routeUrl, mockParams); err != nil {
		return nil, err
	}

	// Subnet
	subnetResponse, err := builders.NewSubnetBuilder().
		Name(subnet.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(subnet.Metadata.Tenant).Workspace(subnet.Metadata.Workspace).Network(subnet.Metadata.Network).Region(subnet.Metadata.Region).
		Spec(&subnet.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a subnet
	if err := configurator.ConfigureCreateSubnetStub(subnetResponse, subnetUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created subnet
	if err := configurator.ConfigureGetActiveSubnetStub(subnetResponse, subnetUrl, mockParams); err != nil {
		return nil, err
	}

	// Update the subnet
	subnetResponse.Spec = *&suiteParams.SubnetUpdated.Spec
	if err := configurator.ConfigureUpdateSubnetStub(subnetResponse, subnetUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the updated subnet
	if err := configurator.ConfigureGetActiveSubnetStub(subnetResponse, subnetUrl, mockParams); err != nil {
		return nil, err
	}

	// Public ip
	publicIpResponse, err := builders.NewPublicIpBuilder().
		Name(publicIp.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(publicIp.Metadata.Tenant).Workspace(publicIp.Metadata.Workspace).Region(publicIp.Metadata.Region).
		Spec(&publicIp.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a public ip
	if err := configurator.ConfigureCreatePublicIpStub(publicIpResponse, publicIpUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created public ip
	if err := configurator.ConfigureGetActivePublicIpStub(publicIpResponse, publicIpUrl, mockParams); err != nil {
		return nil, err
	}

	// Update the public ip
	publicIpResponse.Spec = *&suiteParams.PublicIpUpdated.Spec
	if err := configurator.ConfigureUpdatePublicIpStub(publicIpResponse, publicIpUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the updated public ip
	if err := configurator.ConfigureGetActivePublicIpStub(publicIpResponse, publicIpUrl, mockParams); err != nil {
		return nil, err
	}

	// Nic
	nicResponse, err := builders.NewNicBuilder().
		Name(nic.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(nic.Metadata.Tenant).Workspace(nic.Metadata.Workspace).Region(nic.Metadata.Region).
		Spec(&nic.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a nic
	if err := configurator.ConfigureCreateNicStub(nicResponse, nicUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created nic
	if err := configurator.ConfigureGetActiveNicStub(nicResponse, nicUrl, mockParams); err != nil {
		return nil, err
	}

	// Update the nic
	nicResponse.Spec = *&suiteParams.NicUpdated.Spec
	if err := configurator.ConfigureUpdateNicStub(nicResponse, nicUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the updated nic
	if err := configurator.ConfigureGetActiveNicStub(nicResponse, nicUrl, mockParams); err != nil {
		return nil, err
	}

	// Security group
	groupResponse, err := builders.NewSecurityGroupBuilder().
		Name(securityGroup.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(securityGroup.Metadata.Tenant).Workspace(securityGroup.Metadata.Workspace).Region(securityGroup.Metadata.Region).
		Spec(&securityGroup.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a security group
	if err := configurator.ConfigureCreateSecurityGroupStub(groupResponse, groupUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created security group
	if err := configurator.ConfigureGetActiveSecurityGroupStub(groupResponse, groupUrl, mockParams); err != nil {
		return nil, err
	}

	// Update the security group
	groupResponse.Spec = *&suiteParams.SecurityGroupUpdated.Spec
	if err := configurator.ConfigureUpdateSecurityGroupStub(groupResponse, groupUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the updated security group
	if err := configurator.ConfigureGetActiveSecurityGroupStub(groupResponse, groupUrl, mockParams); err != nil {
		return nil, err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(blockStorage.Metadata.Name).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(blockStorage.Metadata.Tenant).Workspace(blockStorage.Metadata.Workspace).Region(blockStorage.Metadata.Region).
		Spec(&blockStorage.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create a block storage
	if err := configurator.ConfigureCreateBlockStorageStub(blockResponse, blockUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created block storage
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockResponse, blockUrl, mockParams); err != nil {
		return nil, err
	}

	// Instance
	instanceResponse, err := builders.NewInstanceBuilder().
		Name(instance.Metadata.Name).
		Provider(constants.ComputeProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(instance.Metadata.Tenant).Workspace(instance.Metadata.Workspace).Region(instance.Metadata.Region).
		Spec(&instance.Spec).
		Build()
	if err != nil {
		return nil, err
	}

	// Create an instance
	if err := configurator.ConfigureCreateInstanceStub(instanceResponse, instanceUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the created instance
	if err := configurator.ConfigureGetActiveInstanceStub(instanceResponse, instanceUrl, mockParams); err != nil {
		return nil, err
	}
	// Delete the instance
	if err := configurator.ConfigureDeleteStub(instanceUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted instance
	if err := configurator.ConfigureGetNotFoundStub(instanceUrl, mockParams); err != nil {
		return nil, err
	}

	// Delete the block storage
	if err := configurator.ConfigureDeleteStub(blockUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted block storage
	if err := configurator.ConfigureGetNotFoundStub(blockUrl, mockParams); err != nil {
		return nil, err
	}

	// Delete the security group
	if err := configurator.ConfigureDeleteStub(groupUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted security group
	if err := configurator.ConfigureGetNotFoundStub(groupUrl, mockParams); err != nil {
		return nil, err
	}

	// Delete the nic
	if err := configurator.ConfigureDeleteStub(nicUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted nic
	if err := configurator.ConfigureGetNotFoundStub(nicUrl, mockParams); err != nil {
		return nil, err
	}

	// Delete the public ip
	if err := configurator.ConfigureDeleteStub(publicIpUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted public ip
	if err := configurator.ConfigureGetNotFoundStub(publicIpUrl, mockParams); err != nil {
		return nil, err
	}

	// Delete the subnet
	if err := configurator.ConfigureDeleteStub(subnetUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted subnet
	if err := configurator.ConfigureGetNotFoundStub(subnetUrl, mockParams); err != nil {
		return nil, err
	}

	// Delete the route table
	if err := configurator.ConfigureDeleteStub(routeUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted route table
	if err := configurator.ConfigureGetNotFoundStub(routeUrl, mockParams); err != nil {
		return nil, err
	}

	// Delete the internet gateway
	if err := configurator.ConfigureDeleteStub(gatewayUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted internet gateway
	if err := configurator.ConfigureGetNotFoundStub(gatewayUrl, mockParams); err != nil {
		return nil, err
	}

	// Delete the network
	if err := configurator.ConfigureDeleteStub(networkUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted network
	if err := configurator.ConfigureGetNotFoundStub(networkUrl, mockParams); err != nil {
		return nil, err
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceUrl, mockParams); err != nil {
		return nil, err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetNotFoundStub(workspaceUrl, mockParams); err != nil {
		return nil, err
	}

	return configurator.Client, nil
}
