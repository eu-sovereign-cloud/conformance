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
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"k8s.io/utils/ptr"
)

type FoundationUsageV1TestSuite struct {
	suites.MixedTestSuite

	config *FoundationUsageV1Config
	params *params.FoundationUsageV1Params
}

type FoundationUsageV1Config struct {
	Users          []string
	NetworkCidr    string
	PublicIpsRange string
	RegionZones    []string
	StorageSkus    []string
	InstanceSkus   []string
	NetworkSkus    []string
}

func CreateFoundationV1TestSuite(mixedTestSuite suites.MixedTestSuite, config *FoundationUsageV1Config) *FoundationUsageV1TestSuite {
	suite := &FoundationUsageV1TestSuite{
		MixedTestSuite: mixedTestSuite,
		config:         config,
	}
	suite.ScenarioName = constants.FoundationUsageV1SuiteName
	return suite
}

func (suite *FoundationUsageV1TestSuite) BeforeAll(t provider.T) {
	var err error

	subnetCidr, err := generators.GenerateSubnetCidr(suite.config.NetworkCidr, 8, 1)
	if err != nil {
		slog.Error("Failed to generate subnet cidr", "error", err)
		return
	}

	// Generate the nic addresses
	nicAddress1, err := generators.GenerateNicAddress(subnetCidr, 1)
	if err != nil {
		slog.Error("Failed to generate nic address", "error", err)
		return
	}

	// Generate the public ips
	publicIpAddress1, err := generators.GeneratePublicIp(suite.config.PublicIpsRange, 1)
	if err != nil {
		slog.Error("Failed to generate public ip", "error", err)
		return
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

	storageSkuRefObj, err := generators.GenerateSkuRefObject(storageSkuName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	blockStorageName := generators.GenerateBlockStorageName()

	blockStorageRefObj, err := generators.GenerateBlockStorageRefObject(blockStorageName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}
	blockStorageSize := generators.GenerateBlockStorageSize()

	imageName := generators.GenerateImageName()
	imageResource := generators.GenerateImageResource(suite.Tenant, imageName)

	instanceSkuRefObj, err := generators.GenerateSkuRefObject(instanceSkuName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	instanceName := generators.GenerateInstanceName()

	networkSkuRefObj, err := generators.GenerateSkuRefObject(networkSkuName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	networkName := generators.GenerateNetworkName()

	internetGatewayName := generators.GenerateInternetGatewayName()
	internetGatewayRefObj, err := generators.GenerateInternetGatewayRefObject(internetGatewayName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	routeTableName := generators.GenerateRouteTableName()
	routeTableRefObj, err := generators.GenerateRouteTableRefObject(routeTableName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	subnetName := generators.GenerateSubnetName()
	subnetRefObj, err := generators.GenerateSubnetRefObject(subnetName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	nicName := generators.GenerateNicName()

	publicIpName := generators.GeneratePublicIpName()
	publicIpRefObj, err := generators.GeneratePublicIpRefObject(publicIpName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	securityGroupName := generators.GenerateSecurityGroupName()

	Role, err := builders.NewRoleBuilder().
		Name(roleName).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Spec(&schema.RoleSpec{
			Permissions: []schema.Permission{
				{Provider: constants.StorageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
			},
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Role: %v", err)
	}

	RoleAssignment, err := builders.NewRoleAssignmentBuilder().
		Name(roleAssignmentName).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Spec(&schema.RoleAssignmentSpec{
			Roles: []string{roleName},
			Subs:  []string{roleAssignmentSub},
			Scopes: []schema.RoleAssignmentScope{
				{Tenants: &[]string{suite.Tenant}},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build RoleAssignment: %v", err)
	}

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
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
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
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
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
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
		Provider(constants.ComputeProviderV1).ApiVersion(constants.ApiVersion1).
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
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: ptr.To(suite.config.NetworkCidr)},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Network: %v", err)
	}

	internetGateway, err := builders.NewInternetGatewayBuilder().
		Name(internetGatewayName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.InternetGatewaySpec{
			EgressOnly: ptr.To(false),
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Internet Gateway: %v", err)
	}

	routeTable, err := builders.NewRouteTableBuilder().
		Name(routeTableName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
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
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).Network(networkName).
		Spec(&schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: &subnetCidr},
			Zone: zone,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Subnet: %v", err)
	}

	nic, err := builders.NewNicBuilder().
		Name(nicName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NicSpec{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Nic: %v", err)
	}

	publicIp, err := builders.NewPublicIpBuilder().
		Name(publicIpName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.PublicIpSpec{
			Version: schema.IPVersionIPv4,
			Address: ptr.To(publicIpAddress1),
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Public IP: %v", err)
	}

	securityGroup, err := builders.NewSecurityGroupBuilder().
		Name(securityGroupName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
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
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockUsage.ConfigureFoundationScenarioV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

//nolint:dupl
func (suite *FoundationUsageV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t,
		constants.AuthorizationProviderV1,
		string(schema.GlobalTenantResourceMetadataKindResourceKindRole),
		string(schema.GlobalTenantResourceMetadataKindResourceKindRoleAssignment),
		constants.WorkspaceProviderV1,
		string(schema.RegionalResourceMetadataKindResourceKindWorkspace),
		constants.StorageProviderV1,
		string(schema.RegionalResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalResourceMetadataKindResourceKindImage),
		constants.NetworkProviderV1,
		string(schema.RegionalResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalResourceMetadataKindResourceKindInternetGateway),
		string(schema.RegionalResourceMetadataKindResourceKindInternetGateway),
		string(schema.RegionalResourceMetadataKindResourceKindNic),
		string(schema.RegionalResourceMetadataKindResourceKindNic),
		string(schema.RegionalResourceMetadataKindResourceKindPublicIP),
		string(schema.RegionalResourceMetadataKindResourceKindPublicIP),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup),
		constants.ComputeProviderV1,
		string(schema.RegionalResourceMetadataKindResourceKindInstance),
	)

	role := suite.params.Role
	roleAssign := suite.params.RoleAssignment
	workspace := suite.params.Workspace
	image := suite.params.Image
	block := suite.params.BlockStorage
	network := suite.params.Network
	gateway := suite.params.InternetGateway
	route := suite.params.RouteTable
	subnet := suite.params.Subnet
	group := suite.params.SecurityGroup
	publicIp := suite.params.PublicIp
	nic := suite.params.Nic
	instance := suite.params.Instance

	t.WithNewStep("Role", func(roleCtx provider.StepCtx) {
		roleSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, roleCtx)

		expectRoleMeta := role.Metadata
		expectRoleSpec := &role.Spec
		roleSteps.CreateOrUpdateRoleV1Step("Create", suite.GlobalClient.AuthorizationV1, role,
			steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
				Metadata:      expectRoleMeta,
				Spec:          expectRoleSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		roleTRef := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.Tenant),
			Name:   role.Metadata.Name,
		}
		role = roleSteps.GetRoleV1Step("GetCreated", suite.GlobalClient.AuthorizationV1, roleTRef,
			steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
				Metadata:      expectRoleMeta,
				Spec:          expectRoleSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	t.WithNewStep("RoleAssignment", func(raCtx provider.StepCtx) {
		raSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, raCtx)

		expectRoleAssignMeta := roleAssign.Metadata
		expectRoleAssignSpec := &roleAssign.Spec
		raSteps.CreateOrUpdateRoleAssignmentV1Step("Create", suite.GlobalClient.AuthorizationV1, roleAssign,
			steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
				Metadata:      expectRoleAssignMeta,
				Spec:          expectRoleAssignSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		roleAssignTRef := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.Tenant),
			Name:   roleAssign.Metadata.Name,
		}
		roleAssign = raSteps.GetRoleAssignmentV1Step("GetCreated", suite.GlobalClient.AuthorizationV1, roleAssignTRef,
			steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
				Metadata:      expectRoleAssignMeta,
				Spec:          expectRoleAssignSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	t.WithNewStep("Workspace", func(wsCtx provider.StepCtx) {
		wsSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, wsCtx)

		expectWorkspaceMeta := workspace.Metadata
		expectWorkspaceLabels := workspace.Labels
		wsSteps.CreateOrUpdateWorkspaceV1Step("Create", suite.RegionalClient.WorkspaceV1, workspace,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				Labels:        expectWorkspaceLabels,
				Metadata:      expectWorkspaceMeta,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		workspaceTRef := secapi.TenantReference{
			Tenant: secapi.TenantID(workspace.Metadata.Tenant),
			Name:   workspace.Metadata.Name,
		}
		workspace = wsSteps.GetWorkspaceV1Step("GetCreated", suite.RegionalClient.WorkspaceV1, workspaceTRef,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				Labels:        expectWorkspaceLabels,
				Metadata:      expectWorkspaceMeta,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	t.WithNewStep("Image", func(imgCtx provider.StepCtx) {
		imgSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, imgCtx)

		expectedImageMeta := image.Metadata
		expectedImageSpec := &image.Spec
		imgSteps.CreateOrUpdateImageV1Step("Create", suite.RegionalClient.StorageV1, image,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
				Metadata:      expectedImageMeta,
				Spec:          expectedImageSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		imageTRef := secapi.TenantReference{
			Tenant: secapi.TenantID(image.Metadata.Tenant),
			Name:   image.Metadata.Name,
		}
		image = imgSteps.GetImageV1Step("GetCreated", suite.RegionalClient.StorageV1, imageTRef,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
				Metadata:      expectedImageMeta,
				Spec:          expectedImageSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	t.WithNewStep("BlockStorage", func(bsCtx provider.StepCtx) {
		bsSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, bsCtx)

		expectedBlockMeta := block.Metadata
		expectedBlockSpec := &block.Spec
		bsSteps.CreateOrUpdateBlockStorageV1Step("Create", suite.RegionalClient.StorageV1, block,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
				Metadata:      expectedBlockMeta,
				Spec:          expectedBlockSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		blockWRef := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(block.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(block.Metadata.Workspace),
			Name:      block.Metadata.Name,
		}
		block = bsSteps.GetBlockStorageV1Step("GetCreated", suite.RegionalClient.StorageV1, blockWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
				Metadata:      expectedBlockMeta,
				Spec:          expectedBlockSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	t.WithNewStep("Network", func(nwCtx provider.StepCtx) {
		nwSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, nwCtx)

		expectNetworkMeta := network.Metadata
		expectNetworkSpec := &network.Spec
		nwSteps.CreateOrUpdateNetworkV1Step("Create", suite.RegionalClient.NetworkV1, network,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
				Metadata:      expectNetworkMeta,
				Spec:          expectNetworkSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		networkWRef := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(network.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(network.Metadata.Workspace),
			Name:      network.Metadata.Name,
		}
		nwSteps.GetNetworkV1Step("GetCreated", suite.RegionalClient.NetworkV1, networkWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
				Metadata:      expectNetworkMeta,
				Spec:          expectNetworkSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	t.WithNewStep("InternetGateway", func(igCtx provider.StepCtx) {
		igSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, igCtx)

		expectGatewayMeta := gateway.Metadata
		expectGatewaySpec := &gateway.Spec
		igSteps.CreateOrUpdateInternetGatewayV1Step("Create", suite.RegionalClient.NetworkV1, gateway,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
				Metadata:      expectGatewayMeta,
				Spec:          expectGatewaySpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		gatewayWRef := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(gateway.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(gateway.Metadata.Workspace),
			Name:      gateway.Metadata.Name,
		}
		igSteps.GetInternetGatewayV1Step("GetCreated", suite.RegionalClient.NetworkV1, gatewayWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
				Metadata:      expectGatewayMeta,
				Spec:          expectGatewaySpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	t.WithNewStep("RouteTable", func(rtCtx provider.StepCtx) {
		rtSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, rtCtx)

		expectRouteMeta := route.Metadata
		expectRouteSpec := &route.Spec
		rtSteps.CreateOrUpdateRouteTableV1Step("Create", suite.RegionalClient.NetworkV1, route,
			steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
				Metadata:      expectRouteMeta,
				Spec:          expectRouteSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		routeNRef := secapi.NetworkReference{
			Tenant:    secapi.TenantID(route.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(route.Metadata.Workspace),
			Network:   secapi.NetworkID(route.Metadata.Network),
			Name:      route.Metadata.Name,
		}
		rtSteps.GetRouteTableV1Step("GetCreated", suite.RegionalClient.NetworkV1, routeNRef,
			steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
				Metadata:      expectRouteMeta,
				Spec:          expectRouteSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	t.WithNewStep("Subnet", func(snCtx provider.StepCtx) {
		snSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, snCtx)

		expectSubnetMeta := subnet.Metadata
		expectSubnetSpec := &subnet.Spec
		snSteps.CreateOrUpdateSubnetV1Step("Create", suite.RegionalClient.NetworkV1, subnet,
			steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
				Metadata:      expectSubnetMeta,
				Spec:          expectSubnetSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		subnetNRef := secapi.NetworkReference{
			Tenant:    secapi.TenantID(subnet.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(subnet.Metadata.Workspace),
			Network:   secapi.NetworkID(subnet.Metadata.Network),
			Name:      subnet.Metadata.Name,
		}
		snSteps.GetSubnetV1Step("GetCreated", suite.RegionalClient.NetworkV1, subnetNRef,
			steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
				Metadata:      expectSubnetMeta,
				Spec:          expectSubnetSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	t.WithNewStep("SecurityGroup", func(sgCtx provider.StepCtx) {
		gSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, sgCtx)

		expectGroupMeta := group.Metadata
		expectGroupSpec := &group.Spec
		gSteps.CreateOrUpdateSecurityGroupV1Step("Create", suite.RegionalClient.NetworkV1, group,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
				Metadata:      expectGroupMeta,
				Spec:          expectGroupSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		groupWRef := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(group.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(group.Metadata.Workspace),
			Name:      group.Metadata.Name,
		}
		gSteps.GetSecurityGroupV1Step("GetCreated", suite.RegionalClient.NetworkV1, groupWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
				Metadata:      expectGroupMeta,
				Spec:          expectGroupSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	t.WithNewStep("PublicIp", func(piCtx provider.StepCtx) {
		piSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, piCtx)

		expectPublicIpMeta := publicIp.Metadata
		expectPublicIpSpec := &publicIp.Spec
		piSteps.CreateOrUpdatePublicIpV1Step("Create", suite.RegionalClient.NetworkV1, publicIp,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
				Metadata:      expectPublicIpMeta,
				Spec:          expectPublicIpSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		publicIpWRef := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(publicIp.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(publicIp.Metadata.Workspace),
			Name:      publicIp.Metadata.Name,
		}
		piSteps.GetPublicIpV1Step("GetCreated", suite.RegionalClient.NetworkV1, publicIpWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
				Metadata:      expectPublicIpMeta,
				Spec:          expectPublicIpSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	t.WithNewStep("Nic", func(nicCtx provider.StepCtx) {
		nicSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, nicCtx)

		expectNicMeta := nic.Metadata
		expectNicSpec := &nic.Spec
		nicSteps.CreateOrUpdateNicV1Step("Create", suite.RegionalClient.NetworkV1, nic,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
				Metadata:      expectNicMeta,
				Spec:          expectNicSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		nicWRef := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(nic.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(nic.Metadata.Workspace),
			Name:      nic.Metadata.Name,
		}
		nicSteps.GetNicV1Step("GetCreated", suite.RegionalClient.NetworkV1, nicWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
				Metadata:      expectNicMeta,
				Spec:          expectNicSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	t.WithNewStep("Instance", func(instCtx provider.StepCtx) {
		instSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, instCtx)

		expectInstanceMeta := instance.Metadata
		expectInstanceSpec := &instance.Spec
		instSteps.CreateOrUpdateInstanceV1Step("Create", suite.RegionalClient.ComputeV1, instance,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
				Metadata:      expectInstanceMeta,
				Spec:          expectInstanceSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		instanceWRef := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(instance.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(instance.Metadata.Workspace),
			Name:      instance.Metadata.Name,
		}
		instance = instSteps.GetInstanceV1Step("GetCreated", suite.RegionalClient.ComputeV1, instanceWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
				Metadata:      expectInstanceMeta,
				Spec:          expectInstanceSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	t.WithNewStep("Delete", func(delCtx provider.StepCtx) {
		delSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, delCtx)

		delSteps.DeleteInstanceV1Step("Instance", suite.RegionalClient.ComputeV1, instance)
		delSteps.DeleteSecurityGroupV1Step("SecurityGroup", suite.RegionalClient.NetworkV1, group)
		delSteps.DeleteNicV1Step("Nic", suite.RegionalClient.NetworkV1, nic)
		delSteps.DeletePublicIpV1Step("PublicIp", suite.RegionalClient.NetworkV1, publicIp)
		delSteps.DeleteSubnetV1Step("Subnet", suite.RegionalClient.NetworkV1, subnet)
		delSteps.DeleteRouteTableV1Step("RouteTable", suite.RegionalClient.NetworkV1, route)
		delSteps.DeleteInternetGatewayV1Step("InternetGateway", suite.RegionalClient.NetworkV1, gateway)
		delSteps.DeleteNetworkV1Step("Network", suite.RegionalClient.NetworkV1, network)
		delSteps.DeleteBlockStorageV1Step("BlockStorage", suite.RegionalClient.StorageV1, block)
		delSteps.DeleteImageV1Step("Image", suite.RegionalClient.StorageV1, image)
		delSteps.DeleteWorkspaceV1Step("Workspace", suite.RegionalClient.WorkspaceV1, workspace)
		delSteps.DeleteRoleAssignmentV1Step("RoleAssignment", suite.GlobalClient.AuthorizationV1, roleAssign)
		delSteps.DeleteRoleV1Step("Role", suite.GlobalClient.AuthorizationV1, role)
	})

	suite.FinishScenario()
}

func (suite *FoundationUsageV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
