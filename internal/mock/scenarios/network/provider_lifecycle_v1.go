package mocknetwork

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
)

func ConfigureProviderLifecycleScenarioV1(scenario *mockscenarios.Scenario, params *params.NetworkProviderLifeCycleV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := *params.Workspace
	blockStorage := *params.BlockStorage
	instance := *params.Instance
	network := *params.NetworkInitial
	gateway := *params.InternetGatewayInitial
	nic := *params.NicInitial
	publicIp := *params.PublicIpInitial
	routeTable := *params.RouteTableInitial
	subnet := *params.SubnetInitial
	securityGroup := *params.SecurityGroupInitial

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
		return err
	}

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspaceResponse, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created workspace
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspaceResponse, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspaceResponse, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Network
	networkResponse, err := builders.NewNetworkBuilder().
		Name(network.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(network.Metadata.Tenant).Workspace(network.Metadata.Workspace).Region(network.Metadata.Region).
		Spec(&network.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a network
	if err := configurator.ConfigureCreateNetworkStub(networkResponse, networkUrl, scenario.MockParams); err != nil {
		return err
	}

	// Internet gateway
	gatewayResponse, err := builders.NewInternetGatewayBuilder().
		Name(gateway.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(gateway.Metadata.Tenant).Workspace(gateway.Metadata.Workspace).Region(gateway.Metadata.Region).
		Spec(&gateway.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create an internet gateway
	if err := configurator.ConfigureCreateInternetGatewayStub(gatewayResponse, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}

	// Route table
	routeResponse, err := builders.NewRouteTableBuilder().
		Name(routeTable.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(routeTable.Metadata.Tenant).Workspace(routeTable.Metadata.Workspace).Network(routeTable.Metadata.Network).Region(routeTable.Metadata.Region).
		Spec(&routeTable.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a route table
	if err := configurator.ConfigureCreateRouteTableStub(routeResponse, routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Subnet
	subnetResponse, err := builders.NewSubnetBuilder().
		Name(subnet.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(subnet.Metadata.Tenant).Workspace(subnet.Metadata.Workspace).Network(subnet.Metadata.Network).Region(subnet.Metadata.Region).
		Spec(&subnet.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a subnet
	if err := configurator.ConfigureCreateSubnetStub(subnetResponse, subnetUrl, scenario.MockParams); err != nil {
		return err
	}

	// Nic
	nicResponse, err := builders.NewNicBuilder().
		Name(nic.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(nic.Metadata.Tenant).Workspace(nic.Metadata.Workspace).Region(nic.Metadata.Region).
		Spec(&nic.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a nic
	if err := configurator.ConfigureCreateNicStub(nicResponse, nicUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created internet gateway
	if err := configurator.ConfigureGetCreatingInternetGatewayStub(gatewayResponse, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInternetGatewayStub(gatewayResponse, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the internet gateway
	gatewayResponse.Spec = params.InternetGatewayUpdated.Spec
	if err := configurator.ConfigureUpdateInternetGatewayStub(gatewayResponse, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated internet gateway
	if err := configurator.ConfigureGetUpdatingInternetGatewayStub(gatewayResponse, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInternetGatewayStub(gatewayResponse, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created route table
	if err := configurator.ConfigureGetCreatingRouteTableStub(routeResponse, routeUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRouteTableStub(routeResponse, routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the route table
	routeResponse.Spec = params.RouteTableUpdated.Spec
	if err := configurator.ConfigureUpdateRouteTableStub(routeResponse, routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated route table
	if err := configurator.ConfigureGetUpdatingRouteTableStub(routeResponse, routeUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRouteTableStub(routeResponse, routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created subnet
	if err := configurator.ConfigureGetCreatingSubnetStub(subnetResponse, subnetUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSubnetStub(subnetResponse, subnetUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the subnet
	subnetResponse.Spec = params.SubnetUpdated.Spec
	if err := configurator.ConfigureUpdateSubnetStub(subnetResponse, subnetUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated subnet
	if err := configurator.ConfigureGetUpdatingSubnetStub(subnetResponse, subnetUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSubnetStub(subnetResponse, subnetUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created nic
	if err := configurator.ConfigureGetCreatingNicStub(nicResponse, nicUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNicStub(nicResponse, nicUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the nic
	nicResponse.Spec = params.NicUpdated.Spec
	if err := configurator.ConfigureUpdateNicStub(nicResponse, nicUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated nic
	if err := configurator.ConfigureGetUpdatingNicStub(nicResponse, nicUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNicStub(nicResponse, nicUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created network
	if err := configurator.ConfigureGetCreatingNetworkStub(networkResponse, networkUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNetworkStub(networkResponse, networkUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the network
	networkResponse.Spec = params.NetworkUpdated.Spec
	if err := configurator.ConfigureUpdateNetworkStub(networkResponse, networkUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated network
	if err := configurator.ConfigureGetUpdatingNetworkStub(networkResponse, networkUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNetworkStub(networkResponse, networkUrl, scenario.MockParams); err != nil {
		return err
	}

	// Public ip
	publicIpResponse, err := builders.NewPublicIpBuilder().
		Name(publicIp.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(publicIp.Metadata.Tenant).Workspace(publicIp.Metadata.Workspace).Region(publicIp.Metadata.Region).
		Spec(&publicIp.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a public ip
	if err := configurator.ConfigureCreatePublicIpStub(publicIpResponse, publicIpUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created public ip
	if err := configurator.ConfigureGetCreatingPublicIpStub(publicIpResponse, publicIpUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActivePublicIpStub(publicIpResponse, publicIpUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the public ip
	publicIpResponse.Spec = params.PublicIpUpdated.Spec
	if err := configurator.ConfigureUpdatePublicIpStub(publicIpResponse, publicIpUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated public ip
	if err := configurator.ConfigureGetUpdatingPublicIpStub(publicIpResponse, publicIpUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActivePublicIpStub(publicIpResponse, publicIpUrl, scenario.MockParams); err != nil {
		return err
	}

	// Security group
	groupResponse, err := builders.NewSecurityGroupBuilder().
		Name(securityGroup.Metadata.Name).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(securityGroup.Metadata.Tenant).Workspace(securityGroup.Metadata.Workspace).Region(securityGroup.Metadata.Region).
		Spec(&securityGroup.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a security group
	if err := configurator.ConfigureCreateSecurityGroupStub(groupResponse, groupUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created security group
	if err := configurator.ConfigureGetCreatingSecurityGroupStub(groupResponse, groupUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSecurityGroupStub(groupResponse, groupUrl, scenario.MockParams); err != nil {
		return err
	}

	// Update the security group
	groupResponse.Spec = params.SecurityGroupUpdated.Spec
	if err := configurator.ConfigureUpdateSecurityGroupStub(groupResponse, groupUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the updated security group
	if err := configurator.ConfigureGetUpdatingSecurityGroupStub(groupResponse, groupUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSecurityGroupStub(groupResponse, groupUrl, scenario.MockParams); err != nil {
		return err
	}

	// Block storage
	blockResponse, err := builders.NewBlockStorageBuilder().
		Name(blockStorage.Metadata.Name).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(blockStorage.Metadata.Tenant).Workspace(blockStorage.Metadata.Workspace).Region(blockStorage.Metadata.Region).
		Spec(&blockStorage.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create a block storage
	if err := configurator.ConfigureCreateBlockStorageStub(blockResponse, blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created block storage
	if err := configurator.ConfigureGetCreatingBlockStorageStub(blockResponse, blockUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockResponse, blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Instance
	instanceResponse, err := builders.NewInstanceBuilder().
		Name(instance.Metadata.Name).
		Provider(constants.ComputeProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(instance.Metadata.Tenant).Workspace(instance.Metadata.Workspace).Region(instance.Metadata.Region).
		Spec(&instance.Spec).
		Build()
	if err != nil {
		return err
	}

	// Create an instance
	if err := configurator.ConfigureCreateInstanceStub(instanceResponse, instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the created instance
	if err := configurator.ConfigureGetCreatingInstanceStub(instanceResponse, instanceUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInstanceStub(instanceResponse, instanceUrl, scenario.MockParams); err != nil {
		return err
	}
	// Delete the instance
	if err := configurator.ConfigureDeleteStub(instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted instance
	if err := configurator.ConfigureGetNotFoundStub(instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the block storage
	if err := configurator.ConfigureDeleteStub(blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted block storage
	if err := configurator.ConfigureGetNotFoundStub(blockUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the security group
	if err := configurator.ConfigureDeleteStub(groupUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted security group
	if err := configurator.ConfigureGetNotFoundStub(groupUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the public ip
	if err := configurator.ConfigureDeleteStub(publicIpUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted public ip
	if err := configurator.ConfigureGetNotFoundStub(publicIpUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the internet gateway
	if err := configurator.ConfigureDeleteStub(gatewayUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted internet gateway
	if err := configurator.ConfigureGetNotFoundStub(gatewayUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the route table
	if err := configurator.ConfigureDeleteStub(routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted route table
	if err := configurator.ConfigureGetNotFoundStub(routeUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the subnet
	if err := configurator.ConfigureDeleteStub(subnetUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted subnet
	if err := configurator.ConfigureGetNotFoundStub(subnetUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the nic
	if err := configurator.ConfigureDeleteStub(nicUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted nic
	if err := configurator.ConfigureGetNotFoundStub(nicUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete the network
	if err := configurator.ConfigureDeleteStub(networkUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted network
	if err := configurator.ConfigureGetNotFoundStub(networkUrl, scenario.MockParams); err != nil {
		return err
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
