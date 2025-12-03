package mock

import (
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/builders"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func newClient(mockURL string) (*wiremock.Client, error) {
	wm := wiremock.NewClient(mockURL)
	if err := wm.ResetAllScenarios(); err != nil {
		return nil, err
	}
	return wm, nil
}

// Network
func createNetworkList(wm *wiremock.Client, scenario string, params *NetworkParamsV1, nextScenario string) ([]schema.Network, error) {
	var networkList []schema.Network
	for i := range *params.Network {
		networkResource := secalib.GenerateNetworkResource(params.Tenant, params.Workspace.Name, (*params.Network)[i].Name)
		networkResponse, err := builders.NewNetworkBuilder().
			Name((*params.Network)[i].Name).
			Provider(secalib.NetworkProviderV1).
			Resource(networkResource).
			ApiVersion(secalib.ApiVersion1).
			Tenant(params.Tenant).
			Workspace(params.Workspace.Name).
			Region(params.Region).
			Labels((*params.Network)[i].InitialLabels).
			Spec((*params.Network)[i].InitialSpec).
			BuildResponse()
		if err != nil {
			return nil, err
		}

		var nextState string
		if i < len(*params.Network)-1 {
			nextState = (*params.Network)[i+1].Name
		} else {
			nextState = nextScenario
		}

		setCreatedRegionalWorkspaceResourceMetadata(networkResponse.Metadata)
		networkResponse.Status = newNetworkStatus(schema.ResourceStateCreating)
		networkResponse.Metadata.Verb = http.MethodPut

		if err := configurePutStub(wm, scenario,
			&stubConfig{
				url:          secalib.GenerateNetworkURL(params.Tenant, params.Workspace.Name, (*params.Network)[i].Name),
				params:       params,
				responseBody: networkResponse,
				currentState: (*params.Network)[i].Name,
				nextState:    nextState,
			}); err != nil {
			return nil, err
		}
		networkList = append(networkList, *networkResponse)
	}
	return networkList, nil
}

func deleteNetworkList(wm *wiremock.Client, scenario string, params *NetworkParamsV1, nextScenario string) error {
	for i := range *params.Network {
		networkUrl := secalib.GenerateNetworkURL(params.Tenant, params.Workspace.Name, (*params.Network)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeleteNetwork_" + (*params.Network)[i].Name
		} else {
			currentState = "GetDeletedNetwork_" + (*params.Network)[i-1].Name
		}

		nextState = "DeleteNetwork_" + (*params.Network)[i].Name

		// Delete the Network
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: networkUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return err
		}

		nextState = func() string {
			if i < len(*params.Network)-1 {
				return "GetDeletedNetwork_" + (*params.Network)[i].Name
			} else {
				return nextScenario
			}
		}()

		if err := configureGetNotFoundStub(wm, scenario,
			&stubConfig{url: networkUrl, params: params, currentState: "DeleteNetwork_" + (*params.Network)[i].Name, nextState: nextState}); err != nil {
			return err
		}
	}
	return nil
}

func createInternetGatewayList(wm *wiremock.Client, scenario string, params *NetworkParamsV1, nextScenario string) ([]schema.InternetGateway, error) {
	var gatewayList []schema.InternetGateway
	for i := range *params.InternetGateway {
		gatewayResource := secalib.GenerateInternetGatewayResource(params.Tenant, params.Workspace.Name, (*params.InternetGateway)[i].Name)
		gatewayResponse, err := builders.NewInternetGatewayBuilder().
			Name((*params.InternetGateway)[i].Name).
			Provider(secalib.NetworkProviderV1).
			Resource(gatewayResource).
			ApiVersion(secalib.ApiVersion1).
			Tenant(params.Tenant).
			Workspace(params.Workspace.Name).
			Region(params.Region).
			Labels((*params.InternetGateway)[i].InitialLabels).
			Spec((*params.InternetGateway)[i].InitialSpec).
			BuildResponse()
		if err != nil {
			return nil, err
		}

		var nextState string
		if i < len(*params.InternetGateway)-1 {
			nextState = (*params.InternetGateway)[i+1].Name
		} else {
			nextState = nextScenario
		}

		setCreatedRegionalWorkspaceResourceMetadata(gatewayResponse.Metadata)
		gatewayResponse.Status = newResourceStatus(schema.ResourceStateCreating)
		gatewayResponse.Metadata.Verb = http.MethodPut

		if err := configurePutStub(wm, scenario,
			&stubConfig{
				url:          secalib.GenerateInternetGatewayURL(params.Tenant, params.Workspace.Name, (*params.InternetGateway)[i].Name),
				params:       params,
				responseBody: gatewayResponse,
				currentState: (*params.InternetGateway)[i].Name,
				nextState:    nextState,
			}); err != nil {
			return nil, err
		}
		gatewayList = append(gatewayList, *gatewayResponse)
	}
	return gatewayList, nil
}

func deleteInternetGatewayList(wm *wiremock.Client, scenario string, params *NetworkParamsV1, nextScenario string) error {
	for i := range *params.InternetGateway {
		internetGatewayUrl := secalib.GenerateInternetGatewayURL(params.Tenant, params.Workspace.Name, (*params.InternetGateway)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeleteInternetGateway_" + (*params.InternetGateway)[i].Name
		} else {
			currentState = "GetDeletedInternetGateway_" + (*params.InternetGateway)[i-1].Name
		}

		nextState = "DeleteInternetGateway_" + (*params.InternetGateway)[i].Name

		// Delete the Internet Gateway
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: internetGatewayUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return err
		}

		nextState = func() string {
			if i < len(*params.InternetGateway)-1 {
				return "GetDeletedInternetGateway_" + (*params.InternetGateway)[i].Name
			} else {
				return nextScenario
			}
		}()

		if err := configureGetNotFoundStub(wm, scenario,
			&stubConfig{url: internetGatewayUrl, params: params, currentState: "DeleteInternetGateway_" + (*params.InternetGateway)[i].Name, nextState: nextState}); err != nil {
			return err
		}
	}
	return nil
}

func createRouteTableList(wm *wiremock.Client, scenario string, params *NetworkParamsV1, nextScenario string) ([]schema.RouteTable, error) {
	var routeTableList []schema.RouteTable
	for i := range *params.RouteTable {
		routeTableResource := secalib.GenerateRouteTableResource(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, (*params.RouteTable)[i].Name)
		routeTableResponse, err := builders.NewRouteTableBuilder().
			Name((*params.RouteTable)[i].Name).
			Provider(secalib.NetworkProviderV1).
			Resource(routeTableResource).
			ApiVersion(secalib.ApiVersion1).
			Tenant(params.Tenant).
			Workspace(params.Workspace.Name).
			Network((*params.Network)[0].Name).
			Region(params.Region).
			Labels((*params.RouteTable)[i].InitialLabels).
			Spec((*params.RouteTable)[i].InitialSpec).
			BuildResponse()
		if err != nil {
			return nil, err
		}
		var nextState string
		if i < len(*params.RouteTable)-1 {
			nextState = (*params.RouteTable)[i+1].Name
		} else {
			nextState = nextScenario
		}
		// Create a route table
		setCreatedRegionalNetworkResourceMetadata(routeTableResponse.Metadata)
		routeTableResponse.Status = newRouteTableStatus(schema.ResourceStateCreating)
		routeTableResponse.Metadata.Verb = http.MethodPut
		if err := configurePutStub(wm, scenario,
			&stubConfig{url: secalib.GenerateRouteTableURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, (*params.RouteTable)[i].Name), params: params, responseBody: routeTableResponse, currentState: (*params.RouteTable)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
		routeTableList = append(routeTableList, *routeTableResponse)
	}
	return routeTableList, nil
}

func deleteRouteTableList(wm *wiremock.Client, scenario string, params *NetworkParamsV1, nextScenario string) error {
	for i := range *params.RouteTable {
		routeTableUrl := secalib.GenerateRouteTableURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, (*params.RouteTable)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeleteRouteTable_" + (*params.RouteTable)[i].Name
		} else {
			currentState = "GetDeletedRouteTable_" + (*params.RouteTable)[i-1].Name
		}

		nextState = "DeleteRouteTable_" + (*params.RouteTable)[i].Name

		// Delete the RouteTable
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: routeTableUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return err
		}

		nextState = func() string {
			if i < len(*params.RouteTable)-1 {
				return "GetDeletedRouteTable_" + (*params.RouteTable)[i].Name
			} else {
				return nextScenario
			}
		}()

		if err := configureGetNotFoundStub(wm, scenario,
			&stubConfig{url: routeTableUrl, params: params, currentState: "DeleteRouteTable_" + (*params.RouteTable)[i].Name, nextState: nextState}); err != nil {
			return err
		}
	}
	return nil
}

func createSubnetList(wm *wiremock.Client, scenario string, params *NetworkParamsV1, nextScenario string) ([]schema.Subnet, error) {
	var subnetList []schema.Subnet
	for i := range *params.Subnet {
		subnetResource := secalib.GenerateSubnetResource(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, (*params.Subnet)[i].Name)
		subnetResponse, err := builders.NewSubnetBuilder().
			Name((*params.Subnet)[i].Name).
			Provider(secalib.NetworkProviderV1).
			Resource(subnetResource).
			ApiVersion(secalib.ApiVersion1).
			Tenant(params.Tenant).
			Workspace(params.Workspace.Name).
			Network((*params.Network)[0].Name).
			Region(params.Region).
			Labels((*params.Subnet)[i].InitialLabels).
			Spec((*params.Subnet)[i].InitialSpec).
			BuildResponse()
		if err != nil {
			return nil, err
		}
		var nextState string
		if i < len(*params.Subnet)-1 {
			nextState = (*params.Subnet)[i+1].Name
		} else {
			nextState = nextScenario
		}
		// Create a subnet
		setCreatedRegionalNetworkResourceMetadata(subnetResponse.Metadata)
		subnetResponse.Status = newSubnetStatus(schema.ResourceStateCreating)
		subnetResponse.Metadata.Verb = http.MethodPut
		if err := configurePutStub(wm, scenario,
			&stubConfig{url: secalib.GenerateSubnetURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, (*params.Subnet)[i].Name), params: params, responseBody: subnetResponse, currentState: (*params.Subnet)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
		subnetList = append(subnetList, *subnetResponse)
	}
	return subnetList, nil
}

func deleteSubnetList(wm *wiremock.Client, scenario string, params *NetworkParamsV1, nextScenario string) error {
	for i := range *params.Subnet {
		subnetUrl := secalib.GenerateSubnetURL(params.Tenant, params.Workspace.Name, (*params.Network)[0].Name, (*params.Subnet)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeleteSubnet_" + (*params.Subnet)[i].Name
		} else {
			currentState = "GetDeletedSubnet_" + (*params.Subnet)[i-1].Name
		}

		nextState = "DeleteSubnet_" + (*params.Subnet)[i].Name

		// Delete the Subnet
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: subnetUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return err
		}

		nextState = func() string {
			if i < len(*params.Subnet)-1 {
				return "GetDeletedSubnet_" + (*params.Subnet)[i].Name
			} else {
				return nextScenario
			}
		}()

		if err := configureGetNotFoundStub(wm, scenario,
			&stubConfig{url: subnetUrl, params: params, currentState: "DeleteSubnet_" + (*params.Subnet)[i].Name, nextState: nextState}); err != nil {
			return err
		}
	}
	return nil
}

func createPublicIpList(wm *wiremock.Client, scenario string, params *NetworkParamsV1, nextScenario string) ([]schema.PublicIp, error) {
	var publicIpList []schema.PublicIp
	for i := range *params.PublicIp {
		publicIpResource := secalib.GeneratePublicIpResource(params.Tenant, params.Workspace.Name, (*params.PublicIp)[i].Name)
		publicIpResponse, err := builders.NewPublicIpBuilder().
			Name((*params.PublicIp)[i].Name).
			Provider(secalib.NetworkProviderV1).
			Resource(publicIpResource).
			ApiVersion(secalib.ApiVersion1).
			Tenant(params.Tenant).
			Workspace(params.Workspace.Name).
			Region(params.Region).
			Labels((*params.PublicIp)[i].InitialLabels).
			Spec((*params.PublicIp)[i].InitialSpec).
			BuildResponse()
		if err != nil {
			return nil, err
		}
		var nextState string
		if i < len(*params.PublicIp)-1 {
			nextState = (*params.PublicIp)[i+1].Name
		} else {
			nextState = nextScenario
		}
		// Create a PublicIp
		setCreatedRegionalWorkspaceResourceMetadata(publicIpResponse.Metadata)
		publicIpResponse.Status = newPublicIpStatus(schema.ResourceStateCreating)
		publicIpResponse.Metadata.Verb = http.MethodPut
		if err := configurePutStub(wm, scenario,
			&stubConfig{url: secalib.GeneratePublicIpURL(params.Tenant, params.Workspace.Name, (*params.PublicIp)[i].Name), params: params, responseBody: publicIpResponse, currentState: (*params.PublicIp)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
		publicIpList = append(publicIpList, *publicIpResponse)
	}
	return publicIpList, nil
}

func deletePublicIpList(wm *wiremock.Client, scenario string, params *NetworkParamsV1, nextScenario string) error {
	for i := range *params.PublicIp {
		publicIpUrl := secalib.GeneratePublicIpURL(params.Tenant, params.Workspace.Name, (*params.PublicIp)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeletePublicIp_" + (*params.PublicIp)[i].Name
		} else {
			currentState = "GetDeletedPublicIp_" + (*params.PublicIp)[i-1].Name
		}

		nextState = "DeletePublicIp_" + (*params.PublicIp)[i].Name

		// Delete the PublicIp
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: publicIpUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return err
		}

		nextState = func() string {
			if i < len(*params.PublicIp)-1 {
				return "GetDeletedPublicIp_" + (*params.PublicIp)[i].Name
			} else {
				return nextScenario
			}
		}()

		if err := configureGetNotFoundStub(wm, scenario,
			&stubConfig{url: publicIpUrl, params: params, currentState: "DeletePublicIp_" + (*params.PublicIp)[i].Name, nextState: nextState}); err != nil {
			return err
		}
	}
	return nil
}

func createNicList(wm *wiremock.Client, scenario string, params *NetworkParamsV1, nextScenario string) ([]schema.Nic, error) {
	var nicList []schema.Nic
	for i := range *params.NIC {
		nicResource := secalib.GenerateNicResource(params.Tenant, params.Workspace.Name, (*params.NIC)[i].Name)
		nicResponse, err := builders.NewNicBuilder().
			Name((*params.NIC)[i].Name).
			Provider(secalib.NetworkProviderV1).
			Resource(nicResource).
			ApiVersion(secalib.ApiVersion1).
			Tenant(params.Tenant).
			Workspace(params.Workspace.Name).
			Region(params.Region).
			Labels((*params.NIC)[i].InitialLabels).
			Spec((*params.NIC)[i].InitialSpec).
			BuildResponse()
		if err != nil {
			return nil, err
		}

		var nextState string
		if i < len(*params.NIC)-1 {
			nextState = (*params.NIC)[i+1].Name
		} else {
			nextState = nextScenario
		}
		// Create a NIC
		setCreatedRegionalWorkspaceResourceMetadata(nicResponse.Metadata)
		nicResponse.Status = newNicStatus(schema.ResourceStateCreating)
		nicResponse.Metadata.Verb = http.MethodPut
		if err := configurePutStub(wm, scenario,
			&stubConfig{url: secalib.GenerateNicURL(params.Tenant, params.Workspace.Name, (*params.NIC)[i].Name), params: params, responseBody: nicResponse, currentState: (*params.NIC)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
		nicList = append(nicList, *nicResponse)
	}
	return nicList, nil
}

func deleteNics(wm *wiremock.Client, scenario string, params *NetworkParamsV1, nextScenario string) error {
	for i := range *params.NIC {
		nicUrl := secalib.GenerateNicURL(params.Tenant, params.Workspace.Name, (*params.NIC)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeleteNic_" + (*params.NIC)[i].Name
		} else {
			currentState = "GetDeletedNic_" + (*params.NIC)[i-1].Name
		}

		nextState = "DeleteNic_" + (*params.NIC)[i].Name

		// Delete the NIC
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: nicUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return err
		}

		// Get the deleted NIC (should return 404)
		nextState = func() string {
			if i < len(*params.NIC)-1 {
				return "GetDeletedNic_" + (*params.NIC)[i].Name
			} else {
				return nextScenario
			}
		}()

		if err := configureGetNotFoundStub(wm, scenario,
			&stubConfig{url: nicUrl, params: params, currentState: "DeleteNic_" + (*params.NIC)[i].Name, nextState: nextState}); err != nil {
			return err
		}
	}
	return nil
}

func createSecurityGroupList(wm *wiremock.Client, scenario string, params *NetworkParamsV1, nextScenario string) ([]schema.SecurityGroup, error) {
	var securityGroupList []schema.SecurityGroup
	for i := range *params.SecurityGroup {
		securityGroupResource := secalib.GenerateSecurityGroupResource(params.Tenant, params.Workspace.Name, (*params.SecurityGroup)[i].Name)
		securityGroupResponse, err := builders.NewSecurityGroupBuilder().
			Name((*params.SecurityGroup)[i].Name).
			Provider(secalib.NetworkProviderV1).
			Resource(securityGroupResource).
			ApiVersion(secalib.ApiVersion1).
			Tenant(params.Tenant).
			Workspace(params.Workspace.Name).
			Region(params.Region).
			Labels((*params.SecurityGroup)[i].InitialLabels).
			Spec((*params.SecurityGroup)[i].InitialSpec).
			BuildResponse()
		if err != nil {
			return nil, err
		}

		var nextState string
		if i < len(*params.SecurityGroup)-1 {
			nextState = (*params.SecurityGroup)[i+1].Name
		} else {
			nextState = nextScenario
		}
		// Create a security group
		setCreatedRegionalWorkspaceResourceMetadata(securityGroupResponse.Metadata)
		securityGroupResponse.Status = newSecurityGroupStatus(schema.ResourceStateCreating)
		securityGroupResponse.Metadata.Verb = http.MethodPut
		if err := configurePutStub(wm, scenario,
			&stubConfig{url: secalib.GenerateSecurityGroupURL(params.Tenant, params.Workspace.Name, (*params.SecurityGroup)[i].Name), params: params, responseBody: securityGroupResponse, currentState: (*params.SecurityGroup)[i].Name, nextState: nextState}); err != nil {
			return nil, err
		}
		securityGroupList = append(securityGroupList, *securityGroupResponse)
	}
	return securityGroupList, nil
}

func deleteSecurityGroupList(wm *wiremock.Client, scenario string, params *NetworkParamsV1, nextScenario string) error {
	for i := range *params.SecurityGroup {
		securityGroupUrl := secalib.GenerateSecurityGroupURL(params.Tenant, params.Workspace.Name, (*params.SecurityGroup)[i].Name)
		var currentState string
		var nextState string

		if i == 0 {
			currentState = "DeleteSecurityGroup_" + (*params.SecurityGroup)[i].Name
		} else {
			currentState = "GetDeletedSecurityGroup_" + (*params.SecurityGroup)[i-1].Name
		}

		nextState = "DeleteSecurityGroup_" + (*params.SecurityGroup)[i].Name

		// Delete the SecurityGroup
		if err := configureDeleteStub(wm, scenario,
			&stubConfig{url: securityGroupUrl, params: params, currentState: currentState, nextState: nextState}); err != nil {
			return err
		}

		nextState = func() string {
			if i < len(*params.SecurityGroup)-1 {
				return "GetDeletedSecurityGroup_" + (*params.SecurityGroup)[i].Name
			} else {
				return nextScenario
			}
		}()

		if err := configureGetNotFoundStub(wm, scenario,
			&stubConfig{url: securityGroupUrl, params: params, currentState: "DeleteSecurityGroup_" + (*params.SecurityGroup)[i].Name, nextState: nextState}); err != nil {
			return err
		}
	}
	return nil
}
