package usage

import (
	"log/slog"
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockUsage "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/usage"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// PrivateSecureWorkspaceV1TestSuite provisions a fully isolated instance: a network
// routed only through an egress-only internet gateway (no inbound path from the
// internet), a restrictive security group, and a nic with no public IP attached at
// all. Unlike the other usage scenarios, no PublicIp resource is ever created here -
// that absence is the point of the scenario.
type PrivateSecureWorkspaceV1TestSuite struct {
	suites.RegionalTestSuite

	config *PrivateSecureWorkspaceV1Config
	params *params.PrivateSecureWorkspaceV1Params
}

type PrivateSecureWorkspaceV1Config struct {
	NetworkCidr  string
	RegionZones  []string
	StorageSkus  []string
	InstanceSkus []string
	NetworkSkus  []string
}

func CreatePrivateSecureWorkspaceV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *PrivateSecureWorkspaceV1Config) *PrivateSecureWorkspaceV1TestSuite {
	suite := &PrivateSecureWorkspaceV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.UsagePrivateSecureWorkspaceV1SuiteName.String()
	return suite
}

func (suite *PrivateSecureWorkspaceV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.UsageParentSuite)

	subnetCidr, err := generators.GenerateSubnetCidr(suite.config.NetworkCidr, 8, 1)
	if err != nil {
		slog.Error("Failed to generate subnet cidr", "error", err)
		t.FailNow()
	}
	nicAddress, err := generators.GenerateNicAddress(subnetCidr, 1)
	if err != nil {
		slog.Error("Failed to generate nic address", "error", err)
		t.FailNow()
	}

	zone := suite.config.RegionZones[rand.Intn(len(suite.config.RegionZones))]
	storageSkuName := suite.config.StorageSkus[rand.Intn(len(suite.config.StorageSkus))]
	instanceSkuName := suite.config.InstanceSkus[rand.Intn(len(suite.config.InstanceSkus))]
	networkSkuName := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]

	workspaceName := generators.GenerateWorkspaceName()

	storageSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, storageSkuName)
	instanceSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.ComputeProviderV1Name, suite.Tenant, instanceSkuName)
	networkSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, networkSkuName)

	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageRefObj := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, blockStorageName)

	networkName := generators.GenerateNetworkName()

	internetGatewayName := generators.GenerateInternetGatewayName()
	internetGatewayRefObj := generators.GenerateInternetGatewayRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, internetGatewayName)

	routeTableName := generators.GenerateRouteTableName()
	routeTableRefObj := generators.GenerateRouteTableRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, routeTableName)

	subnetName := generators.GenerateSubnetName()
	subnetRefObj := generators.GenerateSubnetRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, networkName, subnetName)

	securityGroupName := generators.GenerateSecurityGroupName()

	nicName := generators.GenerateNicName()

	instanceName := generators.GenerateInstanceName()

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Annotations(schema.Annotations{
			"description": "Fully isolated workspace with no internet ingress",
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	blockStorage, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: constants.BlockStorageInitialSize,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}

	network, err := builders.NewNetworkBuilder().
		Name(networkName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NetworkSpec{
			Cidr:   schema.Cidr{Ipv4: suite.config.NetworkCidr},
			SkuRef: *networkSkuRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	// Egress-only: outbound traffic is allowed, nothing can be initiated from the internet
	internetGateway, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.InternetGatewaySpec{
			EgressOnly: true,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Internet Gateway: %v", err)
	}

	routeTable, err := builders.NewRouteTableBuilder().
		Name(routeTableName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
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
		Spec(&schema.SubnetSpec{
			Cidr:          schema.Cidr{Ipv4: subnetCidr},
			RouteTableRef: *routeTableRefObj,
			Zone:          zone,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Subnet: %v", err)
	}

	// Restrictive baseline: explicit ingress and egress entries, no open ingress rule
	securityGroup, err := builders.NewSecurityGroupBuilder().
		Name(securityGroupName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.SecurityGroupSpec{
			Rules: []schema.SecurityGroupRuleSpec{
				{Direction: schema.SecurityGroupRuleDirectionIngress},
				{Direction: schema.SecurityGroupRuleDirectionEgress},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group: %v", err)
	}

	// No PublicIpRefs - this nic is only reachable from inside the network
	nic, err := builders.NewNicBuilder().
		Name(nicName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NicSpec{
			Addresses: []string{nicAddress},
			SubnetRef: *subnetRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Nic: %v", err)
	}

	instance, err := builders.NewInstanceBuilder().
		Name(instanceName).
		Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Instance: %v", err)
	}

	p := &params.PrivateSecureWorkspaceV1Params{
		Workspace:       workspace,
		BlockStorage:    blockStorage,
		Network:         network,
		InternetGateway: internetGateway,
		RouteTable:      routeTable,
		Subnet:          subnet,
		SecurityGroup:   securityGroup,
		Nic:             nic,
		Instance:        instance,
	}
	suite.params = p
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockUsage.ConfigurePrivateSecureWorkspaceV1, *p)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *PrivateSecureWorkspaceV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t,
		sdkconsts.WorkspaceProviderV1Name,
		sdkconsts.StorageProviderV1Name,
		sdkconsts.ComputeProviderV1Name,
		sdkconsts.NetworkProviderV1Name,
	)
	suite.ConfigureResources(t,
		string(schema.RegionalResourceMetadataKindResourceKindWorkspace),
		string(schema.RegionalResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalResourceMetadataKindResourceKindInstance),
		string(schema.RegionalResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalResourceMetadataKindResourceKindInternetGateway),
		string(schema.RegionalResourceMetadataKindResourceKindNic),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	expectWorkspaceAnnotations := workspace.Annotations
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", t, suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Annotations:    expectWorkspaceAnnotations,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1,
		secapi.TenantReference{Tenant: secapi.TenantID(workspace.Metadata.Tenant), Name: workspace.Metadata.Name},
		steps.ResponseExpectsWithCondition[schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			Labels:   expectWorkspaceLabels,
			Metadata: expectWorkspaceMeta,
			ResourceStatus: schema.WorkspaceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	block := suite.params.BlockStorage
	expectBlockMeta := block.Metadata
	expectBlockSpec := &block.Spec
	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create a block storage", t, suite.Client.StorageV1, block,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:       expectBlockMeta,
			Spec:           expectBlockSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetBlockStorageV1Step("Get the created block storage", suite.Client.StorageV1,
		secapi.WorkspaceReference{Tenant: secapi.TenantID(block.Metadata.Tenant), Workspace: secapi.WorkspaceID(block.Metadata.Workspace), Name: block.Metadata.Name},
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus]{
			Metadata: expectBlockMeta,
			Spec:     expectBlockSpec,
			ResourceStatus: schema.BlockStorageStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	network := suite.params.Network
	expectNetworkMeta := network.Metadata
	expectNetworkSpec := &network.Spec
	stepsBuilder.CreateOrUpdateNetworkV1Step("Create a network", t, suite.Client.NetworkV1, network,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:       expectNetworkMeta,
			Spec:           expectNetworkSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// The response spec comparison here round-trips EgressOnly: true, proving the
	// provider accepts and persists an internet gateway with no inbound path.
	gateway := suite.params.InternetGateway
	expectGatewayMeta := gateway.Metadata
	expectGatewaySpec := &gateway.Spec
	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create an egress-only internet gateway", t, suite.Client.NetworkV1, gateway,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Metadata:       expectGatewayMeta,
			Spec:           expectGatewaySpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetInternetGatewayV1Step("Get the created internet gateway", suite.Client.NetworkV1,
		secapi.WorkspaceReference{Tenant: secapi.TenantID(gateway.Metadata.Tenant), Workspace: secapi.WorkspaceID(gateway.Metadata.Workspace), Name: gateway.Metadata.Name},
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.InternetGatewayStatus]{
			Metadata: expectGatewayMeta,
			Spec:     expectGatewaySpec,
			ResourceStatus: schema.InternetGatewayStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	route := suite.params.RouteTable
	expectRouteMeta := route.Metadata
	expectRouteSpec := &route.Spec
	stepsBuilder.CreateOrUpdateRouteTableV1Step("Create a route table", t, suite.Client.NetworkV1, route,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Metadata:       expectRouteMeta,
			Spec:           expectRouteSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetRouteTableV1Step("Get the created route table", suite.Client.NetworkV1,
		secapi.NetworkReference{Tenant: secapi.TenantID(route.Metadata.Tenant), Workspace: secapi.WorkspaceID(route.Metadata.Workspace), Network: secapi.NetworkID(route.Metadata.Network), Name: route.Metadata.Name},
		steps.ResponseExpectsWithCondition[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus]{
			Metadata: expectRouteMeta,
			Spec:     expectRouteSpec,
			ResourceStatus: schema.RouteTableStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	stepsBuilder.GetNetworkV1Step("Get the created network", suite.Client.NetworkV1,
		secapi.WorkspaceReference{Tenant: secapi.TenantID(network.Metadata.Tenant), Workspace: secapi.WorkspaceID(network.Metadata.Workspace), Name: network.Metadata.Name},
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus]{
			Metadata: expectNetworkMeta,
			Spec:     expectNetworkSpec,
			ResourceStatus: schema.NetworkStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	subnet := suite.params.Subnet
	expectSubnetMeta := subnet.Metadata
	expectSubnetSpec := &subnet.Spec
	stepsBuilder.CreateOrUpdateSubnetV1Step("Create a subnet", t, suite.Client.NetworkV1, subnet,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Metadata:       expectSubnetMeta,
			Spec:           expectSubnetSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetSubnetV1Step("Get the created subnet", suite.Client.NetworkV1,
		secapi.NetworkReference{Tenant: secapi.TenantID(subnet.Metadata.Tenant), Workspace: secapi.WorkspaceID(subnet.Metadata.Workspace), Network: secapi.NetworkID(subnet.Metadata.Network), Name: subnet.Metadata.Name},
		steps.ResponseExpectsWithCondition[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus]{
			Metadata: expectSubnetMeta,
			Spec:     expectSubnetSpec,
			ResourceStatus: schema.SubnetStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	group := suite.params.SecurityGroup
	expectGroupMeta := group.Metadata
	expectGroupSpec := &group.Spec
	stepsBuilder.CreateOrUpdateSecurityGroupV1Step("Create a restrictive security group", t, suite.Client.NetworkV1, group,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Metadata:       expectGroupMeta,
			Spec:           expectGroupSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetSecurityGroupV1Step("Get the created security group", suite.Client.NetworkV1,
		secapi.WorkspaceReference{Tenant: secapi.TenantID(group.Metadata.Tenant), Workspace: secapi.WorkspaceID(group.Metadata.Workspace), Name: group.Metadata.Name},
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus]{
			Metadata: expectGroupMeta,
			Spec:     expectGroupSpec,
			ResourceStatus: schema.SecurityGroupStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// No PublicIp resource is ever created in this scenario - the nic's spec
	// comparison below proves PublicIpRefs comes back empty.
	nic := suite.params.Nic
	expectNicMeta := nic.Metadata
	expectNicSpec := &nic.Spec
	stepsBuilder.CreateOrUpdateNicV1Step("Create a private nic", t, suite.Client.NetworkV1, nic,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Metadata:       expectNicMeta,
			Spec:           expectNicSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetNicV1Step("Get the created private nic", suite.Client.NetworkV1,
		secapi.WorkspaceReference{Tenant: secapi.TenantID(nic.Metadata.Tenant), Workspace: secapi.WorkspaceID(nic.Metadata.Workspace), Name: nic.Metadata.Name},
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus]{
			Metadata: expectNicMeta,
			Spec:     expectNicSpec,
			ResourceStatus: schema.NicStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	instance := suite.params.Instance
	expectInstanceMeta := instance.Metadata
	expectInstanceSpec := &instance.Spec
	stepsBuilder.CreateOrUpdateInstanceV1Step("Create the private instance", t, suite.Client.ComputeV1, instance,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:       expectInstanceMeta,
			Spec:           expectInstanceSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	instanceWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(instance.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(instance.Metadata.Workspace),
		Name:      instance.Metadata.Name,
	}
	instance = stepsBuilder.GetInstanceV1Step("Get the created private instance", suite.Client.ComputeV1, instanceWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			Metadata: expectInstanceMeta,
			Spec:     expectInstanceSpec,
			ResourceStatus: schema.InstanceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Resources deletion - reverse dependency order
	stepsBuilder.DeleteInstanceV1Step("Delete the private instance", t, suite.Client.ComputeV1, instance)
	stepsBuilder.DeleteNicV1Step("Delete the private nic", t, suite.Client.NetworkV1, nic)
	stepsBuilder.DeleteSecurityGroupV1Step("Delete the security group", t, suite.Client.NetworkV1, group)
	stepsBuilder.DeleteSubnetV1Step("Delete the subnet", t, suite.Client.NetworkV1, subnet)
	stepsBuilder.DeleteRouteTableV1Step("Delete the route table", t, suite.Client.NetworkV1, route)
	stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", t, suite.Client.NetworkV1, gateway)
	stepsBuilder.DeleteNetworkV1Step("Delete the network", t, suite.Client.NetworkV1, network)
	stepsBuilder.DeleteBlockStorageV1Step("Delete the block storage", t, suite.Client.StorageV1, block)
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)

	suite.FinishScenario()
}

func (suite *PrivateSecureWorkspaceV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
