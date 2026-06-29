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

// NicErrorV1TestSuite verifies that Nic resources with invalid references
// are rejected by the API with 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create nic with invalid region
//   - Create nic with non-existent workspace
//   - Create nic with non-existent subnetRef
//   - Create nic with non-existent publicIpRef
type NicErrorV1TestSuite struct {
	suites.RegionalTestSuite

	config *NicLifeCycleV1Config
	params *params.NicErrorV1Params
}

func CreateNicErrorV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *NicLifeCycleV1Config) *NicErrorV1TestSuite {
	suite := &NicErrorV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.NicErrorV1SuiteName.String()
	return suite
}

func (suite *NicErrorV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.NetworkParentSuite)

	workspaceName := generators.GenerateWorkspaceName()
	networkName := generators.GenerateNetworkName()
	internetGatewayName := generators.GenerateInternetGatewayName()
	routeTableName := generators.GenerateRouteTableName()
	subnetName := generators.GenerateSubnetName()
	networkSkuName := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]
	zone := suite.config.RegionZones[rand.Intn(len(suite.config.RegionZones))]

	subnetCidr, err := generators.GenerateSubnetCidr(suite.config.NetworkCidr, 8, 1)
	if err != nil {
		t.Fatalf("Failed to generate subnet cidr: %v", err)
	}

	nicAddress, err := generators.GenerateNicAddress(subnetCidr, 1)
	if err != nil {
		t.Fatalf("Failed to generate nic address: %v", err)
	}

	networkSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, networkSkuName)
	internetGatewayRefObj := generators.GenerateInternetGatewayRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, internetGatewayName)
	routeTableRefObj := generators.GenerateRouteTableRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, routeTableName)
	subnetRefObj := generators.GenerateSubnetRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, subnetName)
	nonExistentSubnetRefObj := generators.GenerateSubnetRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, "non-existent-subnet")
	nonExistentPublicIpRefObj := generators.GeneratePublicIpRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, "non-existent-pip")
	baseLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Workspace for nic error scenarios testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	network, err := builders.NewNetworkBuilder().
		Name(networkName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Network for nic error scenarios testing"}).
		Spec(&schema.NetworkSpec{
			Cidr:   schema.Cidr{Ipv4: suite.config.NetworkCidr},
			SkuRef: *networkSkuRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	internetGateway, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "InternetGateway for nic error scenarios testing"}).
		Spec(&schema.InternetGatewaySpec{EgressOnly: false}).Build()
	if err != nil {
		t.Fatalf("Failed to build InternetGateway: %v", err)
	}

	routeTable, err := builders.NewRouteTableBuilder().
		Name(routeTableName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "RouteTable for nic error scenarios testing"}).
		Spec(&schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build RouteTable: %v", err)
	}

	subnet, err := builders.NewSubnetBuilder().
		Name(subnetName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Subnet for nic error scenarios testing"}).
		Spec(&schema.SubnetSpec{
			Cidr:          schema.Cidr{Ipv4: subnetCidr},
			RouteTableRef: *routeTableRefObj,
			Zone:          zone,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Subnet: %v", err)
	}

	buildNic := func(name string, workspaceRef string, region string, subnetRef schema.Reference, publicIpRefs []schema.Reference) *schema.Nic {
		n, err := builders.NewNicBuilder().
			Name(name).
			Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceRef).Region(region).
			Labels(baseLabels).
			Annotations(schema.Annotations{"description": "Nic for error scenario testing"}).
			Spec(&schema.NicSpec{
				Addresses:    []string{nicAddress},
				SubnetRef:    subnetRef,
				PublicIpRefs: publicIpRefs,
			}).Build()
		if err != nil {
			t.Fatalf("Failed to build Nic: %v", err)
		}
		return n
	}

	p := &params.NicErrorV1Params{
		Workspace:       workspace,
		Network:         network,
		InternetGateway: internetGateway,
		RouteTable:      routeTable,
		Subnet:          subnet,

		// Invalid region — random string, everything else valid
		InvalidRegionNic: buildNic(
			generators.GenerateNicName(),
			workspaceName,
			"invalid-region",
			*subnetRefObj,
			nil,
		),

		// Non-existent workspace — workspace was never created
		NonExistentWorkspaceNic: buildNic(
			generators.GenerateNicName(),
			"non-existent-workspace",
			suite.Region,
			*subnetRefObj,
			nil,
		),

		// Non-existent subnetRef — valid workspace + region, subnet does not exist
		NonExistentSubnetRefNic: buildNic(
			generators.GenerateNicName(),
			workspaceName,
			suite.Region,
			*nonExistentSubnetRefObj,
			nil,
		),

		// Non-existent publicIpRef — valid workspace + region + subnet, public ip does not exist
		NonExistentPublicIpRefNic: buildNic(
			generators.GenerateNicName(),
			workspaceName,
			suite.Region,
			*subnetRefObj,
			[]schema.Reference{*nonExistentPublicIpRefObj},
		),
	}

	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mocknetwork.ConfigureNicErrorV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *NicErrorV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.NetworkProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNic))
	suite.ConfigureDepends(t,
		string(schema.RegionalResourceMetadataKindResourceKindWorkspace),
		string(schema.RegionalResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalResourceMetadataKindResourceKindInternetGateway),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
		string(schema.RegionalResourceMetadataKindResourceKindSubnet),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace setup
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	expectWorkspaceAnnotations := workspace.Annotations
	expectWorkspaceExtensions := workspace.Extensions

	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", t, suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Annotations:    expectWorkspaceAnnotations,
			Extensions:     expectWorkspaceExtensions,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, workspaceTRef,
		steps.ResponseExpectsWithCondition[schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			Labels:   expectWorkspaceLabels,
			Metadata: expectWorkspaceMeta,
			ResourceStatus: schema.WorkspaceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Network setup
	network := suite.params.Network
	expectNetworkMeta := network.Metadata
	expectNetworkSpec := &network.Spec
	expectNetworkLabels := network.Labels
	expectNetworkAnnotations := network.Annotations
	expectNetworkExtensions := network.Extensions

	stepsBuilder.CreateOrUpdateNetworkV1Step("Create a network", t, suite.Client.NetworkV1, network,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Labels:         expectNetworkLabels,
			Annotations:    expectNetworkAnnotations,
			Extensions:     expectNetworkExtensions,
			Metadata:       expectNetworkMeta,
			Spec:           expectNetworkSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	networkWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(network.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(network.Metadata.Workspace),
		Name:      network.Metadata.Name,
	}

	stepsBuilder.GetNetworkV1Step("Get the created network", suite.Client.NetworkV1, networkWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus]{
			Metadata: expectNetworkMeta,
			Spec:     expectNetworkSpec,
			ResourceStatus: schema.NetworkStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Internet gateway setup
	internetGat := suite.params.InternetGateway
	expectIgMeta := internetGat.Metadata
	expectIgSpec := &internetGat.Spec
	expectIgLabels := internetGat.Labels
	expectIgAnnotations := internetGat.Annotations
	expectIgExtensions := internetGat.Extensions

	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create an internet gateway", t, suite.Client.NetworkV1, internetGat,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Labels:         expectIgLabels,
			Annotations:    expectIgAnnotations,
			Extensions:     expectIgExtensions,
			Metadata:       expectIgMeta,
			Spec:           expectIgSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	internetGatWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(internetGat.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(internetGat.Metadata.Workspace),
		Name:      internetGat.Metadata.Name,
	}
	stepsBuilder.GetInternetGatewayV1Step("Get the created internet gateway", suite.Client.NetworkV1, internetGatWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.InternetGatewayStatus]{
			Metadata: expectIgMeta,
			Spec:     expectIgSpec,
			ResourceStatus: schema.InternetGatewayStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Route table setup
	route := suite.params.RouteTable
	expectRouteMeta := route.Metadata
	expectRouteSpec := &route.Spec
	expectRouteLabels := route.Labels
	expectRouteAnnotations := route.Annotations
	expectRouteExtensions := route.Extensions

	stepsBuilder.CreateOrUpdateRouteTableV1Step("Create a route table", t, suite.Client.NetworkV1, route,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Labels:         expectRouteLabels,
			Annotations:    expectRouteAnnotations,
			Extensions:     expectRouteExtensions,
			Metadata:       expectRouteMeta,
			Spec:           expectRouteSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	routeNRef := secapi.NetworkReference{
		Tenant:    secapi.TenantID(route.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(route.Metadata.Workspace),
		Network:   secapi.NetworkID(route.Metadata.Network),
		Name:      route.Metadata.Name,
	}
	stepsBuilder.GetRouteTableV1Step("Get the created route table", suite.Client.NetworkV1, routeNRef,
		steps.ResponseExpectsWithCondition[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus]{
			Metadata: expectRouteMeta,
			Spec:     expectRouteSpec,
			ResourceStatus: schema.RouteTableStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Subnet setup
	subnet := suite.params.Subnet
	expectSubnetMeta := subnet.Metadata
	expectSubnetSpec := &subnet.Spec
	expectSubnetLabels := subnet.Labels
	expectSubnetAnnotations := subnet.Annotations
	expectSubnetExtensions := subnet.Extensions

	stepsBuilder.CreateOrUpdateSubnetV1Step("Create a subnet", t, suite.Client.NetworkV1, subnet,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Labels:         expectSubnetLabels,
			Annotations:    expectSubnetAnnotations,
			Extensions:     expectSubnetExtensions,
			Metadata:       expectSubnetMeta,
			Spec:           expectSubnetSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	subnetNRef := secapi.NetworkReference{
		Tenant:    secapi.TenantID(subnet.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(subnet.Metadata.Workspace),
		Network:   secapi.NetworkID(subnet.Metadata.Network),
		Name:      subnet.Metadata.Name,
	}
	stepsBuilder.GetSubnetV1Step("Get the created subnet", suite.Client.NetworkV1, subnetNRef,
		steps.ResponseExpectsWithCondition[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus]{
			Metadata: expectSubnetMeta,
			Spec:     expectSubnetSpec,
			ResourceStatus: schema.SubnetStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Error scenarios — all must be rejected with 422
	stepsBuilder.CreateOrUpdateNicExpectViolationV1Step(
		"Create a nic with invalid region — expect rejection",
		suite.Client.NetworkV1,
		suite.params.InvalidRegionNic,
	)

	stepsBuilder.CreateOrUpdateNicExpectViolationV1Step(
		"Create a nic with non-existent workspace — expect rejection",
		suite.Client.NetworkV1,
		suite.params.NonExistentWorkspaceNic,
	)

	stepsBuilder.CreateOrUpdateNicExpectViolationV1Step(
		"Create a nic with non-existent subnetRef — expect rejection",
		suite.Client.NetworkV1,
		suite.params.NonExistentSubnetRefNic,
	)

	stepsBuilder.CreateOrUpdateNicExpectViolationV1Step(
		"Create a nic with non-existent publicIpRef — expect rejection",
		suite.Client.NetworkV1,
		suite.params.NonExistentPublicIpRefNic,
	)

	// Teardown — reverse dependency order
	stepsBuilder.DeleteSubnetV1Step("Delete the subnet", t, suite.Client.NetworkV1, subnet)
	stepsBuilder.WatchSubnetUntilDeletedV1Step("Watch the subnet deletion", t, suite.Client.NetworkV1, subnetNRef)

	stepsBuilder.DeleteRouteTableV1Step("Delete the route table", t, suite.Client.NetworkV1, route)
	stepsBuilder.WatchRouteTableUntilDeletedV1Step("Watch the route table deletion", t, suite.Client.NetworkV1, routeNRef)

	stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", t, suite.Client.NetworkV1, internetGat)
	stepsBuilder.WatchInternetGatewayUntilDeletedV1Step("Watch the internet gateway deletion", t, suite.Client.NetworkV1, internetGatWRef)

	stepsBuilder.DeleteNetworkV1Step("Delete the network", t, suite.Client.NetworkV1, network)
	stepsBuilder.WatchNetworkUntilDeletedV1Step("Watch the network deletion", t, suite.Client.NetworkV1, networkWRef)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *NicErrorV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
