package mockusage

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigurePrivateSecureWorkspaceV1(scenario *mockscenarios.Scenario, p params.PrivateSecureWorkspaceV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := p.Workspace
	blockStorage := p.BlockStorage
	network := p.Network
	gateway := p.InternetGateway
	routeTable := p.RouteTable
	subnet := p.Subnet
	securityGroup := p.SecurityGroup
	nic := p.Nic
	instance := p.Instance

	workspaceUrl := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	blockUrl := generators.GenerateBlockStorageURL(sdkconsts.StorageProviderV1Name, blockStorage.Metadata.Tenant, blockStorage.Metadata.Workspace, blockStorage.Metadata.Name)
	networkUrl := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, network.Metadata.Tenant, network.Metadata.Workspace, network.Metadata.Name)
	gatewayUrl := generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, gateway.Metadata.Tenant, gateway.Metadata.Workspace, gateway.Metadata.Name)
	routeUrl := generators.GenerateRouteTableURL(sdkconsts.NetworkProviderV1Name, routeTable.Metadata.Tenant, routeTable.Metadata.Workspace, routeTable.Metadata.Network, routeTable.Metadata.Name)
	subnetUrl := generators.GenerateSubnetURL(sdkconsts.NetworkProviderV1Name, subnet.Metadata.Tenant, subnet.Metadata.Workspace, subnet.Metadata.Network, subnet.Metadata.Name)
	groupUrl := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, securityGroup.Metadata.Tenant, securityGroup.Metadata.Workspace, securityGroup.Metadata.Name)
	nicUrl := generators.GenerateNicURL(sdkconsts.NetworkProviderV1Name, nic.Metadata.Tenant, nic.Metadata.Workspace, nic.Metadata.Name)
	instanceUrl := generators.GenerateInstanceURL(sdkconsts.ComputeProviderV1Name, instance.Metadata.Tenant, instance.Metadata.Workspace, instance.Metadata.Name)

	// Create

	if err := configurator.ConfigureCreateWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateBlockStorageStub(blockStorage, blockUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingBlockStorageStub(blockStorage, blockUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockStorage, blockUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateNetworkStub(network, networkUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateInternetGatewayStub(gateway, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingInternetGatewayStub(gateway, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInternetGatewayStub(gateway, gatewayUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateRouteTableStub(routeTable, routeUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingRouteTableStub(routeTable, routeUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRouteTableStub(routeTable, routeUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureGetCreatingNetworkStub(network, networkUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNetworkStub(network, networkUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateSubnetStub(subnet, subnetUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingSubnetStub(subnet, subnetUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSubnetStub(subnet, subnetUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateSecurityGroupStub(securityGroup, groupUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingSecurityGroupStub(securityGroup, groupUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSecurityGroupStub(securityGroup, groupUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateNicStub(nic, nicUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingNicStub(nic, nicUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNicStub(nic, nicUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateInstanceStub(instance, instanceUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingInstanceStub(instance, instanceUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInstanceStub(instance, instanceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete all

	if err := configurator.ConfigureDeleteStub(instanceUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(nicUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(groupUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(subnetUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(routeUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(gatewayUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(networkUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(blockUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
