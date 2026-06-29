package network

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	pgkscenarios "github.com/eu-sovereign-cloud/conformance/pkg/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/mock/stubs"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
)

func ConfigureProviderQueriesV1(scenario *pgkscenarios.Scenario, params params.NetworkProviderQueriesV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := params.Workspace
	networks := params.Networks
	internetGateways := params.InternetGateways
	routeTables := params.RouteTables
	subnets := params.Subnets
	nics := params.Nics
	publicIps := params.PublicIps
	securityGroupRules := params.SecurityGroupRules
	securityGroups := params.SecurityGroups

	// Generate URLs
	workspaceUrl := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	networkListUrl := generators.GenerateNetworkListURL(sdkconsts.NetworkProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	gatewayListUrl := generators.GenerateInternetGatewayListURL(sdkconsts.NetworkProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	publicIpListUrl := generators.GeneratePublicIpListURL(sdkconsts.NetworkProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	nicListUrl := generators.GenerateNicListURL(sdkconsts.NetworkProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	securityGroupRuleListUrl := generators.GenerateSecurityGroupRuleListURL(sdkconsts.NetworkProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	securityGroupListUrl := generators.GenerateSecurityGroupListURL(sdkconsts.NetworkProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	skuListUrl := generators.GenerateNetworkSkuListURL(sdkconsts.NetworkProviderV1Name, workspace.Metadata.Tenant)

	// Workspace

	// Create a workspace
	if err := configurator.ConfigureCreateWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Create networks
	err = stubs.BulkCreateNetworksStubV1(configurator, scenario.MockParams, networks.Items)
	if err != nil {
		return err
	}
	networkResponse := &params.Networks

	// List
	if err := configurator.ConfigureListNetworkStub(networkResponse, networkListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with limit 1
	networkResponse.Items = networks.Items[:1]
	if err := configurator.ConfigureListNetworkStub(networkResponse, networkListUrl, scenario.MockParams, scenarios.PathParamsLimit("1")); err != nil {
		return err
	}

	// List with label
	networksWithLabel := func(networkList []schema.Network) []schema.Network {
		var filteredNetworks []schema.Network
		for _, network := range networkList {
			if val, ok := network.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredNetworks = append(filteredNetworks, network)
			}
		}
		return filteredNetworks
	}
	networkResponse.Items = networksWithLabel(networks.Items)
	if err := configurator.ConfigureListNetworkStub(networkResponse, networkListUrl, scenario.MockParams, scenarios.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with limit and label
	networkResponse.Items = networksWithLabel(networks.Items)[:1]
	if err := configurator.ConfigureListNetworkStub(networkResponse, networkListUrl, scenario.MockParams, scenarios.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Test Network Skus

	// Create skus
	skusList := steps.GenerateNetworkSkusV1(secapi.TenantID(workspace.Metadata.Tenant))
	skuResponse, err := builders.NewNetworkSkuIteratorBuilder().Provider(sdkconsts.NetworkProviderV1Name).Tenant(workspace.Metadata.Tenant).Items(skusList).Build()
	if err != nil {
		return err
	}

	// List skus
	if err := configurator.ConfigureListNetworkSkuStub(skuResponse, skuListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List skus with limit 1
	if err := configurator.ConfigureListNetworkSkuStub(skuResponse, skuListUrl, scenario.MockParams, scenarios.PathParamsLimit("1")); err != nil {
		return err
	}

	// Create internet gateways
	err = stubs.BulkCreateInternetGatewaysStubV1(configurator, scenario.MockParams, internetGateways.Items)
	if err != nil {
		return err
	}
	gatewayResponse := &params.InternetGateways

	// List
	if err := configurator.ConfigureListInternetGatewayStub(gatewayResponse, gatewayListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with limit 1
	gatewayResponse.Items = internetGateways.Items[:1]
	if err := configurator.ConfigureListInternetGatewayStub(gatewayResponse, gatewayListUrl, scenario.MockParams, scenarios.PathParamsLimit("1")); err != nil {
		return err
	}

	// List with label
	gatewayWithLabel := func(gatewayList []schema.InternetGateway) []schema.InternetGateway {
		var filteredGateway []schema.InternetGateway
		for _, gateway := range gatewayList {
			if val, ok := gateway.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredGateway = append(filteredGateway, gateway)
			}
		}
		return filteredGateway
	}
	gatewayResponse.Items = gatewayWithLabel(internetGateways.Items)
	if err := configurator.ConfigureListInternetGatewayStub(gatewayResponse, gatewayListUrl, scenario.MockParams, scenarios.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with limit and label
	gatewayResponse.Items = gatewayWithLabel(internetGateways.Items)[:1]
	if err := configurator.ConfigureListInternetGatewayStub(gatewayResponse, gatewayListUrl, scenario.MockParams, scenarios.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Create route tables
	err = stubs.BulkCreateRouteTableStubV1(configurator, scenario.MockParams, routeTables.Items)
	if err != nil {
		return err
	}
	networkName := networks.Items[0].Metadata.Name
	routeTableResponse := &params.RouteTables

	// List
	routeTableListUrl := generators.GenerateRouteTableListURL(sdkconsts.NetworkProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name, networkName)
	if err := configurator.ConfigureListRouteTableStub(routeTableResponse, routeTableListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with limit 1
	routeTableResponse.Items = routeTables.Items[:1]
	if err := configurator.ConfigureListRouteTableStub(routeTableResponse, routeTableListUrl, scenario.MockParams, scenarios.PathParamsLimit("1")); err != nil {
		return err
	}

	// List with label
	routeTableWithLabel := func(routeTableList []schema.RouteTable) []schema.RouteTable {
		var filteredRouteTable []schema.RouteTable
		for _, routeTable := range routeTableList {
			if val, ok := routeTable.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredRouteTable = append(filteredRouteTable, routeTable)
			}
		}
		return filteredRouteTable
	}

	routeTableResponse.Items = routeTableWithLabel(routeTables.Items)
	if err := configurator.ConfigureListRouteTableStub(routeTableResponse, routeTableListUrl, scenario.MockParams, scenarios.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with limit and label
	routeTableResponse.Items = routeTableWithLabel(routeTables.Items)[:1]
	if err := configurator.ConfigureListRouteTableStub(routeTableResponse, routeTableListUrl, scenario.MockParams, scenarios.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Subnet
	err = stubs.BulkCreateSubnetsStubV1(configurator, scenario.MockParams, subnets.Items)
	if err != nil {
		return err
	}
	subnetResponse := &params.Subnets

	// List
	subnetListUrl := generators.GenerateSubnetListURL(sdkconsts.NetworkProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name, networkName)
	if err := configurator.ConfigureListSubnetStub(subnetResponse, subnetListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with limit 1
	subnetResponse.Items = subnets.Items[:1]
	if err := configurator.ConfigureListSubnetStub(subnetResponse, subnetListUrl, scenario.MockParams, scenarios.PathParamsLimit("1")); err != nil {
		return err
	}

	// List with label
	subnetWithLabel := func(subnetList []schema.Subnet) []schema.Subnet {
		var filteredSubnet []schema.Subnet
		for _, subnet := range subnetList {
			if val, ok := subnet.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredSubnet = append(filteredSubnet, subnet)
			}
		}
		return filteredSubnet
	}
	subnetResponse.Items = subnetWithLabel(subnets.Items)
	if err := configurator.ConfigureListSubnetStub(subnetResponse, subnetListUrl, scenario.MockParams, scenarios.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with limit and label
	subnetResponse.Items = subnetWithLabel(subnets.Items)[:1]
	if err := configurator.ConfigureListSubnetStub(subnetResponse, subnetListUrl, scenario.MockParams, scenarios.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Cretae public ips
	err = stubs.BulkCreatePublicIpsStubV1(configurator, scenario.MockParams, publicIps.Items)
	if err != nil {
		return err
	}
	publicIpResponse := &params.PublicIps

	// List
	if err := configurator.ConfigureListPublicIpStub(publicIpResponse, publicIpListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with limit 1
	publicIpResponse.Items = publicIps.Items[:1]
	if err := configurator.ConfigureListPublicIpStub(publicIpResponse, publicIpListUrl, scenario.MockParams, scenarios.PathParamsLimit("1")); err != nil {
		return err
	}

	// List with label
	publicIpWithLabel := func(publicIpList []schema.PublicIp) []schema.PublicIp {
		var filteredPublicIp []schema.PublicIp
		for _, publicIp := range publicIpList {
			if val, ok := publicIp.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredPublicIp = append(filteredPublicIp, publicIp)
			}
		}
		return filteredPublicIp
	}
	publicIpResponse.Items = publicIpWithLabel(publicIps.Items)
	if err := configurator.ConfigureListPublicIpStub(publicIpResponse, publicIpListUrl, scenario.MockParams, scenarios.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with limit and label
	publicIpResponse.Items = publicIpWithLabel(publicIps.Items)[:1]
	if err := configurator.ConfigureListPublicIpStub(publicIpResponse, publicIpListUrl, scenario.MockParams, scenarios.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Create nics
	err = stubs.BulkCreateNicsStubV1(configurator, scenario.MockParams, nics.Items)
	if err != nil {
		return err
	}
	nicResponse := &params.Nics

	// List
	if err := configurator.ConfigureListNicStub(nicResponse, nicListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with limit 1
	nicResponse.Items = nics.Items[:1]
	if err := configurator.ConfigureListNicStub(nicResponse, nicListUrl, scenario.MockParams, scenarios.PathParamsLimit("1")); err != nil {
		return err
	}

	// List with label
	nicWithLabel := func(nicList []schema.Nic) []schema.Nic {
		var filteredNic []schema.Nic
		for _, nic := range nicList {
			if val, ok := nic.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredNic = append(filteredNic, nic)
			}
		}
		return filteredNic
	}
	nicResponse.Items = nicWithLabel(nics.Items)
	if err := configurator.ConfigureListNicStub(nicResponse, nicListUrl, scenario.MockParams, scenarios.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with limit and label
	nicResponse.Items = nicWithLabel(nics.Items)[:1]
	if err := configurator.ConfigureListNicStub(nicResponse, nicListUrl, scenario.MockParams, scenarios.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Create security group rules
	err = stubs.BulkCreateSecurityGroupRulesStubV1(configurator, scenario.MockParams, securityGroupRules.Items)
	if err != nil {
		return err
	}
	securityGroupRuleResponse := &params.SecurityGroupRules

	// List
	if err := configurator.ConfigureListSecurityGroupRuleStub(securityGroupRuleResponse, securityGroupRuleListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with limit 1
	securityGroupRuleResponse.Items = securityGroupRules.Items[:1]
	if err := configurator.ConfigureListSecurityGroupRuleStub(securityGroupRuleResponse, securityGroupRuleListUrl, scenario.MockParams, scenarios.PathParamsLimit("1")); err != nil {
		return err
	}

	// List with label
	secGroupRuleWithLabel := func(securityGroupRuleList []schema.SecurityGroupRule) []schema.SecurityGroupRule {
		var filteredSecurityRule []schema.SecurityGroupRule
		for _, sec := range securityGroupRuleList {
			if val, ok := sec.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredSecurityRule = append(filteredSecurityRule, sec)
			}
		}
		return filteredSecurityRule
	}

	securityGroupRuleResponse.Items = secGroupRuleWithLabel(securityGroupRules.Items)
	if err := configurator.ConfigureListSecurityGroupRuleStub(securityGroupRuleResponse, securityGroupRuleListUrl, scenario.MockParams, scenarios.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with limit and label
	securityGroupRuleResponse.Items = secGroupRuleWithLabel(securityGroupRules.Items)[:1]
	if err := configurator.ConfigureListSecurityGroupRuleStub(securityGroupRuleResponse, securityGroupRuleListUrl, scenario.MockParams, scenarios.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Create security groups
	err = stubs.BulkCreateSecurityGroupsStubV1(configurator, scenario.MockParams, securityGroups.Items)
	if err != nil {
		return err
	}
	securityGroupResponse := &params.SecurityGroups

	// List
	if err := configurator.ConfigureListSecurityGroupStub(securityGroupResponse, securityGroupListUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// List with limit 1
	securityGroupResponse.Items = securityGroups.Items[:1]
	if err := configurator.ConfigureListSecurityGroupStub(securityGroupResponse, securityGroupListUrl, scenario.MockParams, scenarios.PathParamsLimit("1")); err != nil {
		return err
	}

	// List with label
	secGroupWithLabel := func(securityGroupList []schema.SecurityGroup) []schema.SecurityGroup {
		var filteredSecurity []schema.SecurityGroup
		for _, sec := range securityGroupList {
			if val, ok := sec.Labels[constants.EnvLabel]; ok && val == constants.EnvConformanceLabel {
				filteredSecurity = append(filteredSecurity, sec)
			}
		}
		return filteredSecurity
	}

	securityGroupResponse.Items = secGroupWithLabel(securityGroups.Items)
	if err := configurator.ConfigureListSecurityGroupStub(securityGroupResponse, securityGroupListUrl, scenario.MockParams, scenarios.PathParamsLabel(constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// List with limit and label
	securityGroupResponse.Items = secGroupWithLabel(securityGroups.Items)[:1]
	if err := configurator.ConfigureListSecurityGroupStub(securityGroupResponse, securityGroupListUrl, scenario.MockParams, scenarios.PathParamsLimitAndLabel("1", constants.EnvLabel, constants.EnvConformanceLabel)); err != nil {
		return err
	}

	// Delete the security group rules
	for _, securityGroupRule := range securityGroupRules.Items {
		securityGroupRuleUrl := generators.GenerateSecurityGroupRuleURL(sdkconsts.NetworkProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name, securityGroupRule.Metadata.Name)

		// Delete the security group rule
		if err := configurator.ConfigureDeleteStub(securityGroupRuleUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted security group rule
		if err := configurator.ConfigureGetDeletingSecurityGroupRuleStub(&securityGroupRule, securityGroupRuleUrl, scenario.MockParams); err != nil {
			return err
		}
		if err := configurator.ConfigureGetNotFoundStub(securityGroupRuleUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete the security groups
	for _, securityGroup := range securityGroups.Items {
		securityGroupUrl := generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name, securityGroup.Metadata.Name)

		// Delete the security group
		if err := configurator.ConfigureDeleteStub(securityGroupUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted security group
		if err := configurator.ConfigureGetDeletingSecurityGroupStub(&securityGroup, securityGroupUrl, scenario.MockParams); err != nil {
			return err
		}
		if err := configurator.ConfigureGetNotFoundStub(securityGroupUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete the nics
	for _, nic := range nics.Items {
		nicUrl := generators.GenerateNicURL(sdkconsts.NetworkProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name, nic.Metadata.Name)

		// Delete the nic
		if err := configurator.ConfigureDeleteStub(nicUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted nic
		if err := configurator.ConfigureGetDeletingNicStub(&nic, nicUrl, scenario.MockParams); err != nil {
			return err
		}
		if err := configurator.ConfigureGetNotFoundStub(nicUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete the public ips
	for _, publicIp := range publicIps.Items {
		publicIpUrl := generators.GeneratePublicIpURL(sdkconsts.NetworkProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name, publicIp.Metadata.Name)

		// Delete the public ip
		if err := configurator.ConfigureDeleteStub(publicIpUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted public ip
		if err := configurator.ConfigureGetDeletingPublicIpStub(&publicIp, publicIpUrl, scenario.MockParams); err != nil {
			return err
		}
		if err := configurator.ConfigureGetNotFoundStub(publicIpUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete the subnets
	for _, subnet := range subnets.Items {
		subnetUrl := generators.GenerateSubnetURL(sdkconsts.NetworkProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name, networkName, subnet.Metadata.Name)

		// Delete the subnet
		if err := configurator.ConfigureDeleteStub(subnetUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted subnet
		if err := configurator.ConfigureGetDeletingSubnetStub(&subnet, subnetUrl, scenario.MockParams); err != nil {
			return err
		}
		if err := configurator.ConfigureGetNotFoundStub(subnetUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete the route tables
	for _, routeTable := range routeTables.Items {
		routeTableUrl := generators.GenerateRouteTableURL(sdkconsts.NetworkProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name, networkName, routeTable.Metadata.Name)

		// Delete the route table
		if err := configurator.ConfigureDeleteStub(routeTableUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted route table
		if err := configurator.ConfigureGetDeletingRouteTableStub(&routeTable, routeTableUrl, scenario.MockParams); err != nil {
			return err
		}
		if err := configurator.ConfigureGetNotFoundStub(routeTableUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete the internet gateways
	for _, gateway := range internetGateways.Items {
		gatewayUrl := generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name, gateway.Metadata.Name)

		// Delete the internet gateway
		if err := configurator.ConfigureDeleteStub(gatewayUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted internet gateway
		if err := configurator.ConfigureGetDeletingInternetGatewayStub(&gateway, gatewayUrl, scenario.MockParams); err != nil {
			return err
		}
		if err := configurator.ConfigureGetNotFoundStub(gatewayUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete the networks
	for _, network := range networks.Items {
		networkUrl := generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name, network.Metadata.Name)

		// Delete the network
		if err := configurator.ConfigureDeleteStub(networkUrl, scenario.MockParams); err != nil {
			return err
		}

		// Get the deleted network
		if err := configurator.ConfigureGetDeletingNetworkStub(&network, networkUrl, scenario.MockParams); err != nil {
			return err
		}
		if err := configurator.ConfigureGetNotFoundStub(networkUrl, scenario.MockParams); err != nil {
			return err
		}
	}

	// Delete the workspace
	if err := configurator.ConfigureDeleteStub(workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	// Get the deleted workspace
	if err := configurator.ConfigureGetDeletingWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetNotFoundStub(workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
