package usage

import (
	"log/slog"
	"math/rand"
	"net/http"

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

type FoundationProvidersV1TestSuite struct {
	suites.MixedTestSuite

	config *FoundationProvidersV1Config
	params *params.FoundationUsageV1Params
}

type FoundationProvidersV1Config struct {
	Users          []string
	NetworkCidr    string
	PublicIpsRange string
	RegionZones    []string
	StorageSkus    []string
	InstanceSkus   []string
	NetworkSkus    []string
}

func CreateFoundationProvidersV1TestSuite(mixedTestSuite suites.MixedTestSuite, config *FoundationProvidersV1Config) *FoundationProvidersV1TestSuite {
	suite := &FoundationProvidersV1TestSuite{
		MixedTestSuite: mixedTestSuite,
		config:         config,
	}
	suite.ScenarioName = constants.UsageFoundationProvidersV1SuiteName.String()
	return suite
}

func (suite *FoundationProvidersV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.UsageParentSuite)

	subnetCidr, err := generators.GenerateSubnetCidr(suite.config.NetworkCidr, 8, 1)
	if err != nil {
		slog.Error("Failed to generate subnet cidr", "error", err)
		t.FailNow()
	}

	// Generate the nic addresses
	nicAddress1, err := generators.GenerateNicAddress(subnetCidr, 1)
	if err != nil {
		slog.Error("Failed to generate nic address", "error", err)
		t.FailNow()
	}

	// Generate the public ips
	publicIpAddress1, err := generators.GeneratePublicIp(suite.config.PublicIpsRange, 1)
	if err != nil {
		slog.Error("Failed to generate public ip", "error", err)
		t.FailNow()
	}

	// Select zones
	zone := suite.config.RegionZones[rand.Intn(len(suite.config.RegionZones))]

	// Select subs
	roleAssignmentSub := suite.config.Users[rand.Intn(len(suite.config.Users))]

	// Select skus
	storageSkuName := suite.config.StorageSkus[rand.Intn(len(suite.config.StorageSkus))]
	instanceSkuName := suite.config.InstanceSkus[rand.Intn(len(suite.config.InstanceSkus))]
	networkSkuName := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	roleName := generators.GenerateRoleName()

	roleAssignmentName := generators.GenerateRoleAssignmentName()

	storageSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, storageSkuName)

	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageRefObj := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, blockStorageName)
	blockStorageSize := constants.BlockStorageInitialSize

	imageName := generators.GenerateImageName()
	imageResource := generators.GenerateImageResource(imageName)

	instanceSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.ComputeProviderV1Name, suite.Tenant, instanceSkuName)
	instanceName := generators.GenerateInstanceName()

	networkSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, networkSkuName)
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

	securityGroupName := generators.GenerateSecurityGroupName()

	Role, err := builders.NewRoleBuilder().
		Name(roleName).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Spec(&schema.RoleSpec{
			Permissions: []schema.Permission{
				{Provider: sdkconsts.StorageProviderV1Name, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
			},
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Role: %v", err)
	}

	RoleAssignment, err := builders.NewRoleAssignmentBuilder().
		Name(roleAssignmentName).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Spec(&schema.RoleAssignmentSpec{
			Roles: []string{roleName},
			Subs:  []string{roleAssignmentSub},
			Scopes: []schema.RoleAssignmentScope{
				{Tenants: []string{suite.Tenant}},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build RoleAssignment: %v", err)
	}

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
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
			SizeGB: blockStorageSize,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}

	image, err := builders.NewImageBuilder().
		Name(imageName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Spec(&schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Image: %v", err)
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
		},
		).Build()
	if err != nil {
		t.Fatalf("Failed to build Instance: %v", err)
	}

	network, err := builders.NewNetworkBuilder().
		Name(networkName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: suite.config.NetworkCidr},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	internetGateway, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
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
			Cidr: schema.Cidr{Ipv4: subnetCidr},
			Zone: zone,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Subnet: %v", err)
	}

	nic, err := builders.NewNicBuilder().
		Name(nicName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NicSpec{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: []schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Nic: %v", err)
	}

	publicIp, err := builders.NewPublicIpBuilder().
		Name(publicIpName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.PublicIpSpec{
			Version: schema.IPVersionIPv4,
			Address: publicIpAddress1,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Public IP: %v", err)
	}

	securityGroup, err := builders.NewSecurityGroupBuilder().
		Name(securityGroupName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.SecurityGroupSpec{
			Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress}},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Security Group: %v", err)
	}

	params := &params.FoundationUsageV1Params{
		Role:            Role,
		RoleAssignment:  RoleAssignment,
		Workspace:       workspace,
		BlockStorage:    blockStorage,
		Image:           image,
		Instance:        instance,
		Network:         network,
		InternetGateway: internetGateway,
		RouteTable:      routeTable,
		Subnet:          subnet,
		Nic:             nic,
		PublicIp:        publicIp,
		SecurityGroup:   securityGroup,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockUsage.ConfigureFoundationScenarioV1, *params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *FoundationProvidersV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t,
		sdkconsts.AuthorizationProviderV1Name,
		sdkconsts.WorkspaceProviderV1Name,
		sdkconsts.StorageProviderV1Name,
		sdkconsts.ComputeProviderV1Name,
		sdkconsts.NetworkProviderV1Name,
	)
	suite.ConfigureResources(t,
		string(schema.GlobalTenantResourceMetadataKindResourceKindRole),
		string(schema.GlobalTenantResourceMetadataKindResourceKindRoleAssignment),
		string(schema.RegionalResourceMetadataKindResourceKindWorkspace),
		string(schema.RegionalResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalResourceMetadataKindResourceKindImage),
		string(schema.RegionalResourceMetadataKindResourceKindInstance),
		string(schema.RegionalResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalResourceMetadataKindResourceKindInternetGateway),
		string(schema.RegionalResourceMetadataKindResourceKindNic),
		string(schema.RegionalResourceMetadataKindResourceKindPublicIP),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroupRule),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup),
	)

	stepsConfigurator := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Authorization

	// Create a role
	role := suite.params.Role
	expectRoleMeta := role.Metadata
	expectRoleSpec := &role.Spec
	stepsConfigurator.CreateOrUpdateRoleV1Step("Create a role", t, suite.GlobalClient.AuthorizationV1, role,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			Metadata:       expectRoleMeta,
			Spec:           expectRoleSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created role
	roleTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   role.Metadata.Name,
	}
	role = stepsConfigurator.GetRoleV1Step("Get the created role", suite.GlobalClient.AuthorizationV1, roleTRef,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			Metadata:       expectRoleMeta,
			Spec:           expectRoleSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Create a role assignment
	roleAssign := suite.params.RoleAssignment
	expectRoleAssignMeta := roleAssign.Metadata
	expectRoleAssignSpec := &roleAssign.Spec
	stepsConfigurator.CreateOrUpdateRoleAssignmentV1Step("Create a role assignment", t, suite.GlobalClient.AuthorizationV1, roleAssign,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			Metadata:       expectRoleAssignMeta,
			Spec:           expectRoleAssignSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created role assignment
	roleAssignTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   roleAssign.Metadata.Name,
	}
	roleAssign = stepsConfigurator.GetRoleAssignmentV1Step("Get the created role assignment", suite.GlobalClient.AuthorizationV1, roleAssignTRef,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			Metadata:       expectRoleAssignMeta,
			Spec:           expectRoleAssignSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Workspace

	// Create a workspace
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	stepsConfigurator.CreateOrUpdateWorkspaceV1Step("Create a workspace", t, suite.RegionalClient.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created Workspace
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	workspace = stepsConfigurator.GetWorkspaceV1Step("Get the created workspace", suite.RegionalClient.WorkspaceV1, workspaceTRef,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Storage

	// Create an image
	image := suite.params.Image
	expectedImageMeta := image.Metadata
	expectedImageSpec := &image.Spec
	stepsConfigurator.CreateOrUpdateImageV1Step("Create an image", t, suite.RegionalClient.StorageV1, image,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			Metadata:       expectedImageMeta,
			Spec:           expectedImageSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created image
	imageTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(image.Metadata.Tenant),
		Name:   image.Metadata.Name,
	}
	image = stepsConfigurator.GetImageV1Step("Get the created image", suite.RegionalClient.StorageV1, imageTRef,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			Metadata:       expectedImageMeta,
			Spec:           expectedImageSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Create a block storage
	block := suite.params.BlockStorage
	expectedBlockMeta := block.Metadata
	expectedBlockSpec := &block.Spec
	stepsConfigurator.CreateOrUpdateBlockStorageV1Step("Create a block storage", t, suite.RegionalClient.StorageV1, block,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
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
	block = stepsConfigurator.GetBlockStorageV1Step("Get the created block storage", suite.RegionalClient.StorageV1, blockWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:       expectedBlockMeta,
			Spec:           expectedBlockSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Network

	// Create a network
	network := suite.params.Network
	expectNetworkMeta := network.Metadata
	expectNetworkSpec := &network.Spec
	stepsConfigurator.CreateOrUpdateNetworkV1Step("Create a network", t, suite.RegionalClient.NetworkV1, network,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:       expectNetworkMeta,
			Spec:           expectNetworkSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Create an internet gateway
	gateway := suite.params.InternetGateway
	expectGatewayMeta := gateway.Metadata
	expectGatewaySpec := &gateway.Spec
	stepsConfigurator.CreateOrUpdateInternetGatewayV1Step("Create a internet gateway", t, suite.RegionalClient.NetworkV1, gateway,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
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
	stepsConfigurator.GetInternetGatewayV1Step("Get the created internet gateway", suite.RegionalClient.NetworkV1, gatewayWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Metadata:       expectGatewayMeta,
			Spec:           expectGatewaySpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Create a route table
	route := suite.params.RouteTable
	expectRouteMeta := route.Metadata
	expectRouteSpec := &route.Spec
	stepsConfigurator.CreateOrUpdateRouteTableV1Step("Create a route table", t, suite.RegionalClient.NetworkV1, route,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
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
	stepsConfigurator.GetRouteTableV1Step("Get the created route table", suite.RegionalClient.NetworkV1, routeNRef,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Metadata:       expectRouteMeta,
			Spec:           expectRouteSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Get the created network
	networkWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(network.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(network.Metadata.Workspace),
		Name:      network.Metadata.Name,
	}
	stepsConfigurator.GetNetworkV1Step("Get the created network", suite.RegionalClient.NetworkV1, networkWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:       expectNetworkMeta,
			Spec:           expectNetworkSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Create a subnet
	subnet := suite.params.Subnet
	expectSubnetMeta := subnet.Metadata
	expectSubnetSpec := &subnet.Spec
	stepsConfigurator.CreateOrUpdateSubnetV1Step("Create a subnet", t, suite.RegionalClient.NetworkV1, subnet,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
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
	stepsConfigurator.GetSubnetV1Step("Get the created subnet", suite.RegionalClient.NetworkV1, subnetNRef,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Metadata:       expectSubnetMeta,
			Spec:           expectSubnetSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Create a security group
	group := suite.params.SecurityGroup
	expectGroupMeta := group.Metadata
	expectGroupSpec := &group.Spec
	stepsConfigurator.CreateOrUpdateSecurityGroupV1Step("Create a security group", t, suite.RegionalClient.NetworkV1, group,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
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
	stepsConfigurator.GetSecurityGroupV1Step("Get the created security group", suite.RegionalClient.NetworkV1, groupWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Metadata:       expectGroupMeta,
			Spec:           expectGroupSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Create a public ip
	publicIp := suite.params.PublicIp
	expectPublicIpMeta := publicIp.Metadata
	expectPublicIpSpec := &publicIp.Spec
	stepsConfigurator.CreateOrUpdatePublicIpV1Step("Create a public ip", t, suite.RegionalClient.NetworkV1, publicIp,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
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
	stepsConfigurator.GetPublicIpV1Step("Get the created public ip", suite.RegionalClient.NetworkV1, publicIpWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			Metadata:       expectPublicIpMeta,
			Spec:           expectPublicIpSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Create a nic
	nic := suite.params.Nic
	expectNicMeta := nic.Metadata
	expectNicSpec := &nic.Spec
	stepsConfigurator.CreateOrUpdateNicV1Step("Create a nic", t, suite.RegionalClient.NetworkV1, nic,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
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
	stepsConfigurator.GetNicV1Step("Get the created nic", suite.RegionalClient.NetworkV1, nicWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Metadata:       expectNicMeta,
			Spec:           expectNicSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Compute

	// Create an instance
	instance := suite.params.Instance
	expectInstanceMeta := instance.Metadata
	expectInstanceSpec := &instance.Spec
	stepsConfigurator.CreateOrUpdateInstanceV1Step("Create an instance", t, suite.RegionalClient.ComputeV1, instance,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
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
	instance = stepsConfigurator.GetInstanceV1Step("Get the created instance", suite.RegionalClient.ComputeV1, instanceWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:       expectInstanceMeta,
			Spec:           expectInstanceSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Resources deletion
	stepsConfigurator.DeleteInstanceV1Step("Delete the instance", t, suite.RegionalClient.ComputeV1, instance)

	stepsConfigurator.DeleteSecurityGroupV1Step("Delete the security group", t, suite.RegionalClient.NetworkV1, group)
	stepsConfigurator.DeleteNicV1Step("Delete the nic", t, suite.RegionalClient.NetworkV1, nic)
	stepsConfigurator.DeletePublicIpV1Step("Delete the public ip", t, suite.RegionalClient.NetworkV1, publicIp)
	stepsConfigurator.DeleteSubnetV1Step("Delete the subnet", t, suite.RegionalClient.NetworkV1, subnet)
	stepsConfigurator.DeleteRouteTableV1Step("Delete the route table", t, suite.RegionalClient.NetworkV1, route)
	stepsConfigurator.DeleteInternetGatewayV1Step("Delete the internet gateway", t, suite.RegionalClient.NetworkV1, gateway)
	stepsConfigurator.DeleteNetworkV1Step("Delete the network", t, suite.RegionalClient.NetworkV1, network)

	stepsConfigurator.DeleteBlockStorageV1Step("Delete the block storage", t, suite.RegionalClient.StorageV1, block)
	stepsConfigurator.DeleteImageV1Step("Delete the image", t, suite.RegionalClient.StorageV1, image)

	stepsConfigurator.DeleteWorkspaceV1Step("Delete the workspace", t, suite.RegionalClient.WorkspaceV1, workspace)

	stepsConfigurator.DeleteRoleAssignmentV1Step("Delete the role assignment", t, suite.GlobalClient.AuthorizationV1, roleAssign)
	stepsConfigurator.DeleteRoleV1Step("Delete the role", t, suite.GlobalClient.AuthorizationV1, role)

	suite.FinishScenario()
}

func (suite *FoundationProvidersV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
