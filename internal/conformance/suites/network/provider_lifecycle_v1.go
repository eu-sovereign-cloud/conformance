package network

import (
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockNetwork "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/network"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type ProviderLifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	config *ProviderLifeCycleV1Config
	params *params.NetworkProviderLifeCycleV1Params
}

type ProviderLifeCycleV1Config struct {
	NetworkCidr    string
	PublicIpsRange string
	RegionZones    []string
	StorageSkus    []string
	InstanceSkus   []string
	NetworkSkus    []string
}

func CreateProviderLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *ProviderLifeCycleV1Config) *ProviderLifeCycleV1TestSuite {
	suite := &ProviderLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.NetworkProviderLifeCycleV1SuiteName.String()
	return suite
}

func (suite *ProviderLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Network")

	// Generate the subnet cidr
	subnetCidr, err := generators.GenerateSubnetCidr(suite.config.NetworkCidr, 8, 1)
	if err != nil {
		t.Fatalf("Failed to generate subnet cidr: %v", err)
	}

	// Generate the nic addresses
	nicAddress1, err := generators.GenerateNicAddress(subnetCidr, 1)
	if err != nil {
		t.Fatalf("Failed to generate nic address: %v", err)
	}
	nicAddress2, err := generators.GenerateNicAddress(subnetCidr, 2)
	if err != nil {
		t.Fatalf("Failed to generate nic address: %v", err)
	}

	// Generate the public ips
	publicIpAddress1, err := generators.GeneratePublicIp(suite.config.PublicIpsRange, 1)
	if err != nil {
		t.Fatalf("Failed to generate public ip: %v", err)
	}
	publicIpAddress2, err := generators.GeneratePublicIp(suite.config.PublicIpsRange, 2)
	if err != nil {
		t.Fatalf("Failed to generate public ip: %v", err)
	}

	// Select zones
	zone1 := suite.config.RegionZones[rand.Intn(len(suite.config.RegionZones))]
	zone2 := suite.config.RegionZones[rand.Intn(len(suite.config.RegionZones))]

	// Select skus
	storageSkuName := suite.config.StorageSkus[rand.Intn(len(suite.config.StorageSkus))]
	instanceSkuName := suite.config.InstanceSkus[rand.Intn(len(suite.config.InstanceSkus))]
	networkSkuName1 := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]
	networkSkuName2 := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	storageSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, storageSkuName)

	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageRefObj := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, blockStorageName)

	instanceSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.ComputeProviderV1Name, suite.Tenant, instanceSkuName)
	instanceName := generators.GenerateInstanceName()
	instanceRefObj := generators.GenerateInstanceRefObject(sdkconsts.ComputeProviderV1Name, suite.Tenant, workspaceName, instanceName)

	networkSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, networkSkuName1)
	networkSkuRef2Obj := generators.GenerateSkuRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, networkSkuName2)

	networkName := generators.GenerateNetworkName()

	internetGatewayName := generators.GenerateInternetGatewayName()
	internetGatewayRefObj := generators.GenerateInternetGatewayRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, internetGatewayName)

	routeTableName := generators.GenerateRouteTableName()
	routeTableRefObj := generators.GenerateRouteTableRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, routeTableName)

	subnetName := generators.GenerateSubnetName()
	subnetRefObj := generators.GenerateSubnetRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, subnetName)

	nicName := generators.GenerateNicName()

	publicIpName := generators.GeneratePublicIpName()
	publicIpRefObj := generators.GeneratePublicIpRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, publicIpName)

	securityGroupRuleName := generators.GenerateSecurityGroupRuleName()
	securityGroupRuleRefObj := generators.GenerateSecurityGroupRuleRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, securityGroupRuleName)

	securityGroupName := generators.GenerateSecurityGroupName()

	blockStorageSize := constants.BlockStorageInitialSize

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

	blockStorage, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Block Storage for conformance testing",
		}).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: blockStorageSize,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}

	instance, err := builders.NewInstanceBuilder().
		Name(instanceName).
		Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Instance for conformance testing",
		}).
		Spec(&schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone1,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		},
		).Build()
	if err != nil {
		t.Fatalf("Failed to build Instance: %v", err)
	}

	networkInitial, err := builders.NewNetworkBuilder().
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

	networkUpdated, err := builders.NewNetworkBuilder().
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
			SkuRef:        *networkSkuRef2Obj,
			RouteTableRef: *routeTableRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	internetGatInitial, err := builders.NewInternetGatewayBuilder().
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

	internetGatUpdated, err := builders.NewInternetGatewayBuilder().
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
			EgressOnly: true,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Internet Gateway: %v", err)
	}

	routeTableInitial, err := builders.NewRouteTableBuilder().
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

	routeTableUpdated, err := builders.NewRouteTableBuilder().
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
				{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *instanceRefObj},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Route Table: %v", err)
	}

	subnetInitial, err := builders.NewSubnetBuilder().
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

	subnetUpdated, err := builders.NewSubnetBuilder().
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
			Zone: zone2,
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
			Addresses:    []string{nicAddress1},
			PublicIpRefs: []schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
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
			Addresses:    []string{nicAddress2},
			PublicIpRefs: []schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Nic: %v", err)
	}

	publicIpInitial, err := builders.NewPublicIpBuilder().
		Name(publicIpName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Public IP for conformance testing",
		}).
		Spec(&schema.PublicIpSpec{
			Version: schema.IPVersionIPv4,
			Address: publicIpAddress1,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Public IP: %v", err)
	}

	publicIpUpdated, err := builders.NewPublicIpBuilder().
		Name(publicIpName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Public IP for conformance testing",
		}).
		Spec(&schema.PublicIpSpec{
			Version: schema.IPVersionIPv4,
			Address: publicIpAddress2,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Public IP: %v", err)
	}

	securityGroupRuleInitial, err := builders.NewSecurityGroupRuleBuilder().
		Name(securityGroupRuleName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Security Group Rule for conformance testing",
		}).
		Spec(&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionIngress}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group Rule: %v", err)
	}

	securityGroupRuleUpdated, err := builders.NewSecurityGroupRuleBuilder().
		Name(securityGroupRuleName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Security Group Rule for conformance testing",
		}).
		Spec(&schema.SecurityGroupRuleSpec{Direction: schema.SecurityGroupRuleDirectionEgress}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group Rule: %v", err)
	}

	securityGroupInitial, err := builders.NewSecurityGroupBuilder().
		Name(securityGroupName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Security Group for conformance testing",
		}).
		Spec(&schema.SecurityGroupSpec{
			RuleRefs: []schema.Reference{*securityGroupRuleRefObj},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group: %v", err)
	}

	securityGroupUpdated, err := builders.NewSecurityGroupBuilder().
		Name(securityGroupName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Security Group for conformance testing",
		}).
		Spec(&schema.SecurityGroupSpec{
			Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress}},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group: %v", err)
	}

	params := &params.NetworkProviderLifeCycleV1Params{
		Workspace:                workspace,
		BlockStorage:             blockStorage,
		Instance:                 instance,
		NetworkInitial:           networkInitial,
		NetworkUpdated:           networkUpdated,
		InternetGatewayInitial:   internetGatInitial,
		InternetGatewayUpdated:   internetGatUpdated,
		RouteTableInitial:        routeTableInitial,
		RouteTableUpdated:        routeTableUpdated,
		SubnetInitial:            subnetInitial,
		SubnetUpdated:            subnetUpdated,
		NicInitial:               nicInitial,
		NicUpdated:               nicUpdated,
		PublicIpInitial:          publicIpInitial,
		PublicIpUpdated:          publicIpUpdated,
		SecurityGroupInitial:     securityGroupInitial,
		SecurityGroupUpdated:     securityGroupUpdated,
		SecurityGroupRuleInitial: securityGroupRuleInitial,
		SecurityGroupRuleUpdated: securityGroupRuleUpdated,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockNetwork.ConfigureProviderLifecycleScenarioV1, *params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ProviderLifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.NetworkProviderV1Name,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindInternetGateway),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindNic),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindPublicIP),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup),
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
		steps.ResponseExpectsWithCondition[schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			Labels:   expectWorkspaceLabels,
			Metadata: expectWorkspaceMeta,
			ResourceStatus: schema.WorkspaceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Network

	// Create a network
	network := suite.params.NetworkInitial
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

	networkWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(network.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(network.Metadata.Workspace),
		Name:      network.Metadata.Name,
	}

	// Internet gateway

	// Create an internet gateway
	gateway := suite.params.InternetGatewayInitial
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

	gatewayWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(gateway.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(gateway.Metadata.Workspace),
		Name:      gateway.Metadata.Name,
	}

	// Get the created internet gateway
	stepsBuilder.GetInternetGatewayV1Step("Get the created internet gateway", suite.Client.NetworkV1, gatewayWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.InternetGatewayStatus]{
			Metadata: expectGatewayMeta,
			Spec:     expectGatewaySpec,
			ResourceStatus: schema.InternetGatewayStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Update the internet gateway
	gateway = suite.params.InternetGatewayUpdated
	expectGatewaySpec.EgressOnly = gateway.Spec.EgressOnly
	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Update the internet gateway", suite.Client.NetworkV1, gateway,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Labels:         expectGatewayLabels,
			Annotations:    expectGatewayAnnotations,
			Extensions:     expectGatewayExtensions,
			Metadata:       expectGatewayMeta,
			Spec:           expectGatewaySpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the updated internet gateway
	stepsBuilder.GetInternetGatewayV1Step("Get the updated internet gateway", suite.Client.NetworkV1, gatewayWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.InternetGatewayStatus]{
			Metadata: expectGatewayMeta,
			Spec:     expectGatewaySpec,
			ResourceStatus: schema.InternetGatewayStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterUpdating,
			},
		},
	)

	// Route table

	// Create a route table
	route := suite.params.RouteTableInitial
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

	routeNRef := secapi.NetworkReference{
		Tenant:    secapi.TenantID(route.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(route.Metadata.Workspace),
		Network:   secapi.NetworkID(route.Metadata.Network),
		Name:      route.Metadata.Name,
	}

	// Get the created route table
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

	// Get the created network
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

	// Update the route table
	route = suite.params.RouteTableUpdated
	expectRouteSpec.Routes = route.Spec.Routes
	expectRouteLabels = route.Labels
	expectRouteAnnotations = route.Annotations
	expectRouteExtensions = route.Extensions
	stepsBuilder.CreateOrUpdateRouteTableV1Step("Update the route table", suite.Client.NetworkV1, route,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Labels:         expectRouteLabels,
			Annotations:    expectRouteAnnotations,
			Extensions:     expectRouteExtensions,
			Metadata:       expectRouteMeta,
			Spec:           expectRouteSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the updated route table
	stepsBuilder.GetRouteTableV1Step("Get the updated route table", suite.Client.NetworkV1, routeNRef,
		steps.ResponseExpectsWithCondition[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus]{
			Metadata: expectRouteMeta,
			Spec:     expectRouteSpec,
			ResourceStatus: schema.RouteTableStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterUpdating,
			},
		},
	)

	// Network

	// Update the network
	network = suite.params.NetworkUpdated
	expectNetworkSpec.SkuRef = network.Spec.SkuRef
	expectNetworkLabels = network.Labels
	expectNetworkAnnotations = network.Annotations
	expectNetworkExtensions = network.Extensions
	stepsBuilder.CreateOrUpdateNetworkV1Step("Update the network", suite.Client.NetworkV1, network,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Labels:         expectNetworkLabels,
			Annotations:    expectNetworkAnnotations,
			Extensions:     expectNetworkExtensions,
			Metadata:       expectNetworkMeta,
			Spec:           expectNetworkSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the updated network
	stepsBuilder.GetNetworkV1Step("Get the updated network", suite.Client.NetworkV1, networkWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus]{
			Metadata: expectNetworkMeta,
			Spec:     expectNetworkSpec,
			ResourceStatus: schema.NetworkStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterUpdating,
			},
		},
	)

	// Subnet

	// Create a subnet
	subnet := suite.params.SubnetInitial
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

	subnetNRef := secapi.NetworkReference{
		Tenant:    secapi.TenantID(subnet.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(subnet.Metadata.Workspace),
		Network:   secapi.NetworkID(subnet.Metadata.Network),
		Name:      subnet.Metadata.Name,
	}

	// Get the created subnet
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

	// Update the subnet
	subnet = suite.params.SubnetUpdated
	expectSubnetSpec.Zone = subnet.Spec.Zone
	expectSubnetLabels = subnet.Labels
	expectSubnetAnnotations = subnet.Annotations
	expectSubnetExtensions = subnet.Extensions
	stepsBuilder.CreateOrUpdateSubnetV1Step("Update the subnet", suite.Client.NetworkV1, subnet,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Labels:         expectSubnetLabels,
			Annotations:    expectSubnetAnnotations,
			Extensions:     expectSubnetExtensions,
			Metadata:       expectSubnetMeta,
			Spec:           expectSubnetSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the updated subnet
	stepsBuilder.GetSubnetV1Step("Get the updated subnet", suite.Client.NetworkV1, subnetNRef,
		steps.ResponseExpectsWithCondition[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus]{
			Metadata: expectSubnetMeta,
			Spec:     expectSubnetSpec,
			ResourceStatus: schema.SubnetStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterUpdating,
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

	nicWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(nic.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(nic.Metadata.Workspace),
		Name:      nic.Metadata.Name,
	}

	// Get the created nic
	stepsBuilder.GetNicV1Step("Get the created nic", suite.Client.NetworkV1, nicWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus]{
			Metadata: expectNicMeta,
			Spec:     expectNicSpec,
			ResourceStatus: schema.NicStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Update the nic
	nic = suite.params.NicUpdated
	expectNicSpec.Addresses = nic.Spec.Addresses
	expectNicLabels = nic.Labels
	expectNicAnnotations = nic.Annotations
	expectNicExtensions = nic.Extensions
	stepsBuilder.CreateOrUpdateNicV1Step("Update the nic", suite.Client.NetworkV1, nic,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Labels:         expectNicLabels,
			Annotations:    expectNicAnnotations,
			Extensions:     expectNicExtensions,
			Metadata:       expectNicMeta,
			Spec:           expectNicSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the updated nic
	stepsBuilder.GetNicV1Step("Get the updated nic", suite.Client.NetworkV1, nicWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus]{
			Metadata: expectNicMeta,
			Spec:     expectNicSpec,
			ResourceStatus: schema.NicStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterUpdating,
			},
		},
	)
	// Public ip

	// Create a public ip
	publicIp := suite.params.PublicIpInitial
	expectPublicIpMeta := publicIp.Metadata
	expectPublicIpSpec := &publicIp.Spec
	expectPublicIpLabels := publicIp.Labels
	expectPublicIpAnnotations := publicIp.Annotations
	expectPublicIpExtensions := publicIp.Extensions
	stepsBuilder.CreateOrUpdatePublicIpV1Step("Create a public ip", suite.Client.NetworkV1, publicIp,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			Labels:         expectPublicIpLabels,
			Annotations:    expectPublicIpAnnotations,
			Extensions:     expectPublicIpExtensions,
			Metadata:       expectPublicIpMeta,
			Spec:           expectPublicIpSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created public ip
	publicIpWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(publicIp.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(publicIp.Metadata.Workspace),
		Name:      publicIp.Metadata.Name,
	}
	stepsBuilder.GetPublicIpV1Step("Get the created public ip", suite.Client.NetworkV1, publicIpWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus]{
			Metadata: expectPublicIpMeta,
			Spec:     expectPublicIpSpec,
			ResourceStatus: schema.PublicIpStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Update the public ip
	publicIp = suite.params.PublicIpUpdated
	expectPublicIpSpec.Address = publicIp.Spec.Address
	expectPublicIpLabels = publicIp.Labels
	expectPublicIpAnnotations = publicIp.Annotations
	expectPublicIpExtensions = publicIp.Extensions
	stepsBuilder.CreateOrUpdatePublicIpV1Step("Update the public ip", suite.Client.NetworkV1, publicIp,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			Labels:         expectPublicIpLabels,
			Annotations:    expectPublicIpAnnotations,
			Extensions:     expectPublicIpExtensions,
			Metadata:       expectPublicIpMeta,
			Spec:           expectPublicIpSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the updated public ip
	stepsBuilder.GetPublicIpV1Step("Get the updated public ip", suite.Client.NetworkV1, publicIpWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus]{
			Metadata: expectPublicIpMeta,
			Spec:     expectPublicIpSpec,
			ResourceStatus: schema.PublicIpStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterUpdating,
			},
		},
	)

	// Security Group Rule

	// Create a security group rule
	rule := suite.params.SecurityGroupRuleInitial
	expectRuleMeta := rule.Metadata
	expectRuleSpec := &rule.Spec
	expectRuleLabels := rule.Labels
	expectRuleAnnotations := rule.Annotations
	expectRuleExtensions := rule.Extensions
	stepsBuilder.CreateOrUpdateSecurityGroupRuleV1Step("Create a security group rule", suite.Client.NetworkV1, rule,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec]{
			Labels:         expectRuleLabels,
			Annotations:    expectRuleAnnotations,
			Extensions:     expectRuleExtensions,
			Metadata:       expectRuleMeta,
			Spec:           expectRuleSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created security group rule
	ruleWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(rule.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(rule.Metadata.Workspace),
		Name:      rule.Metadata.Name,
	}
	stepsBuilder.GetSecurityGroupRuleV1Step("Get the created security group rule", suite.Client.NetworkV1, ruleWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec, schema.SecurityGroupRuleStatus]{
			Metadata: expectRuleMeta,
			Spec:     expectRuleSpec,
			ResourceStatus: schema.SecurityGroupRuleStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Update the security group rule
	rule = suite.params.SecurityGroupRuleUpdated
	expectRuleSpec.Direction = rule.Spec.Direction
	expectRuleLabels = rule.Labels
	expectRuleAnnotations = rule.Annotations
	expectRuleExtensions = rule.Extensions
	stepsBuilder.CreateOrUpdateSecurityGroupRuleV1Step("Update the security group rule", suite.Client.NetworkV1, rule,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec]{
			Labels:         expectRuleLabels,
			Annotations:    expectRuleAnnotations,
			Extensions:     expectRuleExtensions,
			Metadata:       expectRuleMeta,
			Spec:           expectRuleSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the updated security group rule
	stepsBuilder.GetSecurityGroupRuleV1Step("Get the updated security group rule", suite.Client.NetworkV1, ruleWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec, schema.SecurityGroupRuleStatus]{
			Metadata: expectRuleMeta,
			Spec:     expectRuleSpec,
			ResourceStatus: schema.SecurityGroupRuleStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterUpdating,
			},
		},
	)

	// Security Group

	// Create a security group
	group := suite.params.SecurityGroupInitial
	expectGroupMeta := group.Metadata
	expectGroupSpec := &group.Spec
	expectGroupLabels := group.Labels
	expectGroupAnnotations := group.Annotations
	expectGroupExtensions := group.Extensions
	stepsBuilder.CreateOrUpdateSecurityGroupV1Step("Create a security group", suite.Client.NetworkV1, group,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Labels:         expectGroupLabels,
			Annotations:    expectGroupAnnotations,
			Extensions:     expectGroupExtensions,
			Metadata:       expectGroupMeta,
			Spec:           expectGroupSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created security group
	groupWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(group.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(group.Metadata.Workspace),
		Name:      group.Metadata.Name,
	}
	stepsBuilder.GetSecurityGroupV1Step("Get the created security group", suite.Client.NetworkV1, groupWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus]{
			Metadata: expectGroupMeta,
			Spec:     expectGroupSpec,
			ResourceStatus: schema.SecurityGroupStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Update the security group

	group.Spec.Rules = []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress}}
	expectGroupSpec = &group.Spec
	expectGroupLabels = group.Labels
	expectGroupAnnotations = group.Annotations
	expectGroupExtensions = group.Extensions
	stepsBuilder.CreateOrUpdateSecurityGroupV1Step("Update the security group", suite.Client.NetworkV1, group,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Labels:         expectGroupLabels,
			Annotations:    expectGroupAnnotations,
			Extensions:     expectGroupExtensions,
			Metadata:       expectGroupMeta,
			Spec:           expectGroupSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the updated security group
	stepsBuilder.GetSecurityGroupV1Step("Get the updated security group", suite.Client.NetworkV1, groupWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus]{
			Metadata: expectGroupMeta,
			Spec:     expectGroupSpec,
			ResourceStatus: schema.SecurityGroupStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterUpdating,
			},
		},
	)

	// Block storage

	// Create a block storage
	block := suite.params.BlockStorage
	expectedBlockMeta := block.Metadata
	expectedBlockSpec := &block.Spec
	expectedBlockLabels := block.Labels
	expectedBlockAnnotations := block.Annotations
	expectedBlockExtensions := block.Extensions
	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create a block storage", suite.Client.StorageV1, block,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Labels:         expectedBlockLabels,
			Annotations:    expectedBlockAnnotations,
			Extensions:     expectedBlockExtensions,
			Metadata:       expectedBlockMeta,
			Spec:           expectedBlockSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created block storage
	blockWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(block.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(block.Metadata.Workspace),
		Name:      block.Metadata.Name,
	}
	stepsBuilder.GetBlockStorageV1Step("Get the created block storage", suite.Client.StorageV1, blockWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus]{
			Metadata: expectedBlockMeta,
			Spec:     expectedBlockSpec,
			ResourceStatus: schema.BlockStorageStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Instance

	// Create an instance
	instance := suite.params.Instance
	expectInstanceMeta := instance.Metadata
	expectInstanceSpec := &instance.Spec
	expectInstanceLabels := instance.Labels
	expectInstanceAnnotations := instance.Annotations
	expectInstanceExtensions := instance.Extensions
	stepsBuilder.CreateOrUpdateInstanceV1Step("Create an instance", suite.Client.ComputeV1, instance,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Labels:         expectInstanceLabels,
			Annotations:    expectInstanceAnnotations,
			Extensions:     expectInstanceExtensions,
			Metadata:       expectInstanceMeta,
			Spec:           expectInstanceSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created instance
	instanceWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(instance.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(instance.Metadata.Workspace),
		Name:      instance.Metadata.Name,
	}
	instance = stepsBuilder.GetInstanceV1Step("Get the created instance", suite.Client.ComputeV1, instanceWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			Metadata: expectInstanceMeta,
			Spec:     expectInstanceSpec,
			ResourceStatus: schema.InstanceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Resources deletion

	stepsBuilder.DeleteInstanceV1Step("Delete the instance", suite.Client.ComputeV1, instance)
	stepsBuilder.WatchInstanceUntilDeletedV1Step("Watch the instance deletion", suite.Client.ComputeV1, instanceWRef)

	stepsBuilder.DeleteBlockStorageV1Step("Delete the block storage", suite.Client.StorageV1, block)
	stepsBuilder.WatchBlockStorageUntilDeletedV1Step("Watch the block storage deletion", suite.Client.StorageV1, blockWRef)

	stepsBuilder.DeleteSecurityGroupRuleV1Step("Delete the security group rule", suite.Client.NetworkV1, rule)
	stepsBuilder.WatchSecurityGroupRuleUntilDeletedV1Step("Watch the security group rule deletion", suite.Client.NetworkV1, ruleWRef)

	stepsBuilder.DeleteSecurityGroupV1Step("Delete the security group", suite.Client.NetworkV1, group)
	stepsBuilder.WatchSecurityGroupUntilDeletedV1Step("Watch the security group deletion", suite.Client.NetworkV1, groupWRef)

	stepsBuilder.DeleteNicV1Step("Delete the nic", suite.Client.NetworkV1, nic)
	stepsBuilder.WatchNicUntilDeletedV1Step("Watch the nic deletion", suite.Client.NetworkV1, nicWRef)

	stepsBuilder.DeletePublicIpV1Step("Delete the public ip", suite.Client.NetworkV1, publicIp)
	stepsBuilder.WatchPublicIpUntilDeletedV1Step("Watch the public ip deletion", suite.Client.NetworkV1, publicIpWRef)

	stepsBuilder.DeleteSubnetV1Step("Delete the subnet", suite.Client.NetworkV1, subnet)
	stepsBuilder.WatchSubnetUntilDeletedV1Step("Watch the subnet deletion", suite.Client.NetworkV1, subnetNRef)

	stepsBuilder.DeleteRouteTableV1Step("Delete the route table", suite.Client.NetworkV1, route)
	stepsBuilder.WatchRouteTableUntilDeletedV1Step("Watch the route table deletion", suite.Client.NetworkV1, routeNRef)

	stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", suite.Client.NetworkV1, gateway)
	stepsBuilder.WatchInternetGatewayUntilDeletedV1Step("Watch the internet gateway deletion", suite.Client.NetworkV1, gatewayWRef)

	stepsBuilder.DeleteNetworkV1Step("Delete the network", suite.Client.NetworkV1, network)
	stepsBuilder.WatchNetworkUntilDeletedV1Step("Watch the network deletion", suite.Client.NetworkV1, networkWRef)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *ProviderLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
