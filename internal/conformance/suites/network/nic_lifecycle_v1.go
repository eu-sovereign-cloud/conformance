package network

import (
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mocknetwork "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/network"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type NicLifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	config *NicLifeCycleV1Config
	params *params.NicLifeCycleV1Params
}

type NicLifeCycleV1Config struct {
	NetworkCidr    string
	PublicIpsRange string
	RegionZones    []string
	NetworkSkus    []string
}

func CreateNicLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *NicLifeCycleV1Config) *NicLifeCycleV1TestSuite {
	suite := &NicLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.NicLifeCycleV1SuiteName.String()
	return suite
}

func (suite *NicLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Network")

	workspaceName := generators.GenerateWorkspaceName()
	networkName := generators.GenerateNetworkName()
	internetGatewayName := generators.GenerateInternetGatewayName()
	routeTableName := generators.GenerateRouteTableName()
	nicName := generators.GenerateNicName()
	subnetName := generators.GenerateSubnetName()
	networkSkuName := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]

	zone1 := suite.config.RegionZones[rand.Intn(len(suite.config.RegionZones))]

	subnetCidr, err := generators.GenerateSubnetCidr(suite.config.NetworkCidr, 8, 1)
	if err != nil {
		t.Fatalf("Failed to generate subnet cidr: %v", err)
	}

	subnetRefObj := generators.GenerateSubnetRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, subnetName)

	networkSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, networkSkuName)

	internetGatewayRefObj := generators.GenerateInternetGatewayRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, internetGatewayName)

	routeTableRefObj := generators.GenerateRouteTableRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, routeTableName)

	// Generate the nic addresses
	nicAddress1, err := generators.GenerateNicAddress(subnetCidr, 1)
	if err != nil {
		t.Fatalf("Failed to generate nic address: %v", err)
	}
	nicAddress2, err := generators.GenerateNicAddress(subnetCidr, 2)
	if err != nil {
		t.Fatalf("Failed to generate nic address: %v", err)
	}

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Workspace for conformance testing",
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	network, err := builders.NewNetworkBuilder().
		Name(networkName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Network for conformance testing",
		}).
		Spec(&schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: suite.config.NetworkCidr},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	internetGateway, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Internet Gateway for conformance testing",
		}).
		Spec(&schema.InternetGatewaySpec{
			EgressOnly: false,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Internet Gateway: %v", err)
	}

	routeTable, err := builders.NewRouteTableBuilder().
		Name(routeTableName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Route Table for conformance testing",
		}).
		Spec(&schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Route Table: %v", err)
	}
	subnet, err := builders.NewSubnetBuilder().
		Name(subnetName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Subnet for conformance testing",
		}).
		Spec(&schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: subnetCidr},
			Zone: zone1,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Subnet: %v", err)
	}

	nicInitial, err := builders.NewNicBuilder().
		Name(nicName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Nic for conformance testing",
		}).
		Spec(&schema.NicSpec{
			Addresses: []string{nicAddress1},
			SubnetRef: *subnetRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Nic: %v", err)
	}

	nicUpdated, err := builders.NewNicBuilder().
		Name(nicName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Nic for conformance testing",
		}).
		Spec(&schema.NicSpec{
			Addresses: []string{nicAddress2},
			SubnetRef: *subnetRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Nic: %v", err)
	}

	params := &params.NicLifeCycleV1Params{
		Workspace:       workspace,
		Network:         network,
		InternetGateway: internetGateway,
		RouteTable:      routeTable,
		Subnet:          subnet,
		NicInitial:      nicInitial,
		NicUpdated:      nicUpdated,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mocknetwork.ConfigureNicLifecycleScenarioV1, *params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *NicLifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.NetworkProviderV1Name,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace

	// Create a workspace
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	expectWorkspaceAnnotations := workspace.Annotations
	expectWorkspaceExtensions := workspace.Extensions
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Annotations:    expectWorkspaceAnnotations,
			Extensions:     expectWorkspaceExtensions,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created Workspace
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, workspaceTRef,
		steps.ResponseExpectsWithCondition[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:      expectWorkspaceLabels,
			Annotations: expectWorkspaceAnnotations,
			Extensions:  expectWorkspaceExtensions,
			Metadata:    expectWorkspaceMeta,
			ResourceStatus: schema.Status{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)
	// Network

	// Create a network
	network := suite.params.Network
	expectNetworkMeta := network.Metadata
	expectNetworkSpec := &network.Spec
	expectNetworkLabels := network.Labels
	expectNetworkAnnotations := network.Annotations
	expectNetworkExtensions := network.Extensions
	stepsBuilder.CreateOrUpdateNetworkV1Step("Create a network", suite.Client.NetworkV1, network,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Labels:         expectNetworkLabels,
			Annotations:    expectNetworkAnnotations,
			Extensions:     expectNetworkExtensions,
			Metadata:       expectNetworkMeta,
			Spec:           expectNetworkSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Internet gateway

	// Create an internet gateway
	gateway := suite.params.InternetGateway
	expectGatewayMeta := gateway.Metadata
	expectGatewaySpec := &gateway.Spec
	expectGatewayLabels := gateway.Labels
	expectGatewayAnnotations := gateway.Annotations
	expectGatewayExtensions := gateway.Extensions
	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create a internet gateway", suite.Client.NetworkV1, gateway,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Labels:         expectGatewayLabels,
			Annotations:    expectGatewayAnnotations,
			Extensions:     expectGatewayExtensions,
			Metadata:       expectGatewayMeta,
			Spec:           expectGatewaySpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created internet gateway
	gatewayWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(gateway.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(gateway.Metadata.Workspace),
		Name:      gateway.Metadata.Name,
	}
	stepsBuilder.GetInternetGatewayV1Step("Get the created internet gateway", suite.Client.NetworkV1, gatewayWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Labels:      expectGatewayLabels,
			Annotations: expectGatewayAnnotations,
			Extensions:  expectGatewayExtensions,
			Metadata:    expectGatewayMeta,
			Spec:        expectGatewaySpec,
			ResourceStatus: schema.Status{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Route table

	// Create a route table
	route := suite.params.RouteTable
	expectRouteMeta := route.Metadata
	expectRouteSpec := &route.Spec
	expectRouteLabels := route.Labels
	expectRouteAnnotations := route.Annotations
	expectRouteExtensions := route.Extensions
	stepsBuilder.CreateOrUpdateRouteTableV1Step("Create a route table", suite.Client.NetworkV1, route,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Labels:         expectRouteLabels,
			Annotations:    expectRouteAnnotations,
			Extensions:     expectRouteExtensions,
			Metadata:       expectRouteMeta,
			Spec:           expectRouteSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created route table
	routeNRef := secapi.NetworkReference{
		Tenant:    secapi.TenantID(route.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(route.Metadata.Workspace),
		Network:   secapi.NetworkID(route.Metadata.Network),
		Name:      route.Metadata.Name,
	}
	stepsBuilder.GetRouteTableV1Step("Get the created route table", suite.Client.NetworkV1, routeNRef,
		steps.ResponseExpectsWithCondition[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Labels:      expectRouteLabels,
			Annotations: expectRouteAnnotations,
			Extensions:  expectRouteExtensions,
			Metadata:    expectRouteMeta,
			Spec:        expectRouteSpec,
			ResourceStatus: schema.Status{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Get the created network
	networkWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(network.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(network.Metadata.Workspace),
		Name:      network.Metadata.Name,
	}
	stepsBuilder.GetNetworkV1Step("Get the created network", suite.Client.NetworkV1, networkWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Labels:      expectNetworkLabels,
			Annotations: expectNetworkAnnotations,
			Extensions:  expectNetworkExtensions,
			Metadata:    expectNetworkMeta,
			Spec:        expectNetworkSpec,
			ResourceStatus: schema.Status{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Subnet

	// Create a subnet
	subnet := suite.params.Subnet
	expectSubnetMeta := subnet.Metadata
	expectSubnetSpec := &subnet.Spec
	expectSubnetLabels := subnet.Labels
	expectSubnetAnnotations := subnet.Annotations
	expectSubnetExtensions := subnet.Extensions
	stepsBuilder.CreateOrUpdateSubnetV1Step("Create a subnet", suite.Client.NetworkV1, subnet,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Labels:         expectSubnetLabels,
			Annotations:    expectSubnetAnnotations,
			Extensions:     expectSubnetExtensions,
			Metadata:       expectSubnetMeta,
			Spec:           expectSubnetSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created subnet
	subnetNRef := secapi.NetworkReference{
		Tenant:    secapi.TenantID(subnet.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(subnet.Metadata.Workspace),
		Network:   secapi.NetworkID(subnet.Metadata.Network),
		Name:      subnet.Metadata.Name,
	}
	stepsBuilder.GetSubnetV1Step("Get the created subnet", suite.Client.NetworkV1, subnetNRef,
		steps.ResponseExpectsWithCondition[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Labels:      expectSubnetLabels,
			Annotations: expectSubnetAnnotations,
			Extensions:  expectSubnetExtensions,
			Metadata:    expectSubnetMeta,
			Spec:        expectSubnetSpec,
			ResourceStatus: schema.Status{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Nic

	// Create a nic
	nic := suite.params.NicInitial
	expectNicMeta := nic.Metadata
	expectNicSpec := &nic.Spec
	expectNicLabels := nic.Labels
	expectNicAnnotations := nic.Annotations
	expectNicExtensions := nic.Extensions
	stepsBuilder.CreateOrUpdateNicV1Step("Create a nic", suite.Client.NetworkV1, nic,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Labels:         expectNicLabels,
			Annotations:    expectNicAnnotations,
			Extensions:     expectNicExtensions,
			Metadata:       expectNicMeta,
			Spec:           expectNicSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created nic
	nicWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(nic.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(nic.Metadata.Workspace),
		Name:      nic.Metadata.Name,
	}
	stepsBuilder.GetNicV1Step("Get the created nic", suite.Client.NetworkV1, nicWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Labels:      expectNicLabels,
			Annotations: expectNicAnnotations,
			Extensions:  expectNicExtensions,
			Metadata:    expectNicMeta,
			Spec:        expectNicSpec,
			ResourceStatus: schema.Status{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Update the nic
	nic = suite.params.NicUpdated
	expectNicSpec = &nic.Spec
	expectNicLabels = nic.Labels
	expectNicAnnotations = nic.Annotations
	expectNicExtensions = nic.Extensions
	stepsBuilder.CreateOrUpdateNicV1Step("Create a nic", suite.Client.NetworkV1, nic,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Labels:         expectNicLabels,
			Annotations:    expectNicAnnotations,
			Extensions:     expectNicExtensions,
			Metadata:       expectNicMeta,
			Spec:           expectNicSpec,
			ResourceStates: suites.UpdatedResourceExpectedStates,
		},
	)
	// Get the updated nic
	stepsBuilder.GetNicV1Step("Get the updated nic", suite.Client.NetworkV1, nicWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Labels:      expectNicLabels,
			Annotations: expectNicAnnotations,
			Extensions:  expectNicExtensions,
			Metadata:    expectNicMeta,
			Spec:        expectNicSpec,
			ResourceStatus: schema.Status{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterUpdating,
			},
		},
	)

	// Resources deletion
	stepsBuilder.DeleteNicV1Step("Delete the nic", suite.Client.NetworkV1, nic)
	stepsBuilder.WatchNicUntilDeletedV1Step("Watch the nic deletion", suite.Client.NetworkV1, nicWRef)

	stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", suite.Client.NetworkV1, gateway)
	stepsBuilder.WatchInternetGatewayUntilDeletedV1Step("Watch the internet gateway deletion", suite.Client.NetworkV1, gatewayWRef)

	stepsBuilder.DeleteRouteTableV1Step("Delete the route table", suite.Client.NetworkV1, route)
	stepsBuilder.WatchRouteTableUntilDeletedV1Step("Watch the route table deletion", suite.Client.NetworkV1, routeNRef)

	stepsBuilder.DeleteNetworkV1Step("Delete the network", suite.Client.NetworkV1, network)
	stepsBuilder.WatchNetworkUntilDeletedV1Step("Watch the network deletion", suite.Client.NetworkV1, networkWRef)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *NicLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
