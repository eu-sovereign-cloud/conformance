package usage

import (
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	mockusage "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/usage"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"k8s.io/utils/ptr"
)

type FoundationV1TestSuite struct {
	suites.MixedTestSuite

	Users          []string
	NetworkCidr    string
	PublicIpsRange string
	RegionZones    []string
	StorageSkus    []string
	InstanceSkus   []string
	NetworkSkus    []string
}

func (suite *FoundationV1TestSuite) TestScenario(t provider.T) {
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

	var err error

	// Generate the subnet cidr
	subnetCidr, err := generators.GenerateSubnetCidr(suite.NetworkCidr, 8, 1)
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
	publicIpAddress1, err := generators.GeneratePublicIp(suite.PublicIpsRange, 1)
	if err != nil {
		slog.Error("Failed to generate public ip", "error", err)
		return
	}

	// Select zones
	zone := suite.RegionZones[rand.Intn(len(suite.RegionZones))]

	// Select subs
	roleAssignmentSub := suite.Users[rand.Intn(len(suite.Users))]

	// Select skus
	storageSkuName := suite.StorageSkus[rand.Intn(len(suite.StorageSkus))]
	instanceSkuName := suite.InstanceSkus[rand.Intn(len(suite.InstanceSkus))]
	networkSkuName := suite.NetworkSkus[rand.Intn(len(suite.NetworkSkus))]

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
	initialStorageSize := generators.GenerateBlockStorageSize()

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

	// Setup mock, if configured to use
	if suite.MockEnabled {
		mockParams := &params.FoundationUsageParamsV1{
			BaseParams: &params.BaseParams{
				Tenant: suite.Tenant,
				Region: suite.Region,
				MockParams: &mock.MockParams{
					ServerURL: *suite.MockServerURL,
					AuthToken: suite.AuthToken,
				},
			},
			Role: &params.ResourceParams[schema.RoleSpec]{
				Name: roleName,
				InitialSpec: &schema.RoleSpec{
					Permissions: []schema.Permission{
						{Provider: constants.StorageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
					},
				},
			},
			RoleAssignment: &params.ResourceParams[schema.RoleAssignmentSpec]{
				Name: roleAssignmentName,
				InitialSpec: &schema.RoleAssignmentSpec{
					Roles: []string{roleName},
					Subs:  []string{roleAssignmentSub},
					Scopes: []schema.RoleAssignmentScope{
						{Tenants: &[]string{suite.Tenant}},
					},
				},
			},
			Workspace: &params.ResourceParams[schema.WorkspaceSpec]{
				Name: workspaceName,
				InitialLabels: schema.Labels{
					constants.EnvLabel: constants.EnvDevelopmentLabel,
				},
			},
			BlockStorage: &params.ResourceParams[schema.BlockStorageSpec]{
				Name: blockStorageName,
				InitialSpec: &schema.BlockStorageSpec{
					SkuRef: *storageSkuRefObj,
					SizeGB: initialStorageSize,
				},
			},
			Image: &params.ResourceParams[schema.ImageSpec]{
				Name: imageName,
				InitialSpec: &schema.ImageSpec{
					BlockStorageRef: *blockStorageRefObj,
					CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
				},
			},
			Network: &params.ResourceParams[schema.NetworkSpec]{
				Name: networkName,
				InitialSpec: &schema.NetworkSpec{
					Cidr:          schema.Cidr{Ipv4: ptr.To(suite.NetworkCidr)},
					SkuRef:        *networkSkuRefObj,
					RouteTableRef: *routeTableRefObj,
				},
			},
			InternetGateway: &params.ResourceParams[schema.InternetGatewaySpec]{
				Name:        internetGatewayName,
				InitialSpec: &schema.InternetGatewaySpec{EgressOnly: ptr.To(false)},
			},
			RouteTable: &params.ResourceParams[schema.RouteTableSpec]{
				Name: routeTableName,
				InitialSpec: &schema.RouteTableSpec{
					Routes: []schema.RouteSpec{
						{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
					},
				},
			},
			Subnet: &params.ResourceParams[schema.SubnetSpec]{
				Name: subnetName,
				InitialSpec: &schema.SubnetSpec{
					Cidr: schema.Cidr{Ipv4: &subnetCidr},
					Zone: zone,
				},
			},
			Nic: &params.ResourceParams[schema.NicSpec]{
				Name: nicName,
				InitialSpec: &schema.NicSpec{
					Addresses:    []string{nicAddress1},
					PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
					SubnetRef:    *subnetRefObj,
				},
			},
			PublicIp: &params.ResourceParams[schema.PublicIpSpec]{
				Name: publicIpName,
				InitialSpec: &schema.PublicIpSpec{
					Version: schema.IPVersionIPv4,
					Address: ptr.To(publicIpAddress1),
				},
			},
			SecurityGroup: &params.ResourceParams[schema.SecurityGroupSpec]{
				Name: securityGroupName,
				InitialSpec: &schema.SecurityGroupSpec{
					Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress}},
				},
			},
			Instance: &params.ResourceParams[schema.InstanceSpec]{
				Name: instanceName,
				InitialSpec: &schema.InstanceSpec{
					SkuRef: *instanceSkuRefObj,
					Zone:   zone,
					BootVolume: schema.VolumeReference{
						DeviceRef: *blockStorageRefObj,
					},
				},
			},
		}
		wm, err := mockusage.ConfigureFoundationScenarioV1(suite.ScenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}

		suite.MockClient = wm
	}

	stepsBuilder := steps.NewStepsConfigurator(&suite.TestSuite, t)

	// Role

	// Create a role
	role := &schema.Role{
		Metadata: &schema.GlobalTenantResourceMetadata{
			Tenant: suite.Tenant,
			Name:   roleName,
		},
		Spec: schema.RoleSpec{
			Permissions: []schema.Permission{
				{
					Provider:  constants.StorageProviderV1,
					Resources: []string{imageResource},
					Verb:      []string{http.MethodGet},
				},
			},
		},
	}
	expectRoleMeta, err := builders.NewRoleMetadataBuilder().
		Name(roleName).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectRoleSpec := &schema.RoleSpec{
		Permissions: []schema.Permission{
			{
				Provider:  constants.StorageProviderV1,
				Resources: []string{imageResource},
				Verb:      []string{http.MethodGet},
			},
		},
	}
	stepsBuilder.CreateOrUpdateRoleV1Step("Create a role", suite.GlobalClient.AuthorizationV1, role,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			Metadata:      expectRoleMeta,
			Spec:          expectRoleSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created role
	roleTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   roleName,
	}
	role = stepsBuilder.GetRoleV1Step("Get the created role", suite.GlobalClient.AuthorizationV1, *roleTRef,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			Metadata:      expectRoleMeta,
			Spec:          expectRoleSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Role assignment

	// Create a role assignment
	roleAssign := &schema.RoleAssignment{
		Metadata: &schema.GlobalTenantResourceMetadata{
			Tenant: suite.Tenant,
			Name:   roleAssignmentName,
		},
		Spec: schema.RoleAssignmentSpec{
			Roles:  []string{roleName},
			Subs:   []string{roleAssignmentSub},
			Scopes: []schema.RoleAssignmentScope{{Tenants: &[]string{suite.Tenant}}},
		},
	}
	expectRoleAssignMeta, err := builders.NewRoleAssignmentMetadataBuilder().
		Name(roleAssignmentName).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectRoleAssignSpec := &schema.RoleAssignmentSpec{
		Roles:  []string{roleName},
		Subs:   []string{roleAssignmentSub},
		Scopes: []schema.RoleAssignmentScope{{Tenants: &[]string{suite.Tenant}}},
	}
	stepsBuilder.CreateOrUpdateRoleAssignmentV1Step("Create a role assignment", suite.GlobalClient.AuthorizationV1, roleAssign,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			Metadata:      expectRoleAssignMeta,
			Spec:          expectRoleAssignSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created role assignment
	roleAssignTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   roleAssignmentName,
	}
	roleAssign = stepsBuilder.GetRoleAssignmentV1Step("Get the created role assignment", suite.GlobalClient.AuthorizationV1, *roleAssignTRef,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			Metadata:      expectRoleAssignMeta,
			Spec:          expectRoleAssignSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Workspace

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.Tenant,
			Name:   workspaceName,
		},
	}
	expectWorkspaceMeta, err := builders.NewWorkspaceMetadataBuilder().
		Name(workspaceName).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectWorkspaceLabels := schema.Labels{constants.EnvLabel: constants.EnvDevelopmentLabel}
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.RegionalClient.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectWorkspaceLabels,
			Metadata:      expectWorkspaceMeta,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created Workspace
	tref := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   workspaceName,
	}
	workspace = stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.RegionalClient.WorkspaceV1, *tref,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectWorkspaceLabels,
			Metadata:      expectWorkspaceMeta,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Image

	// Create an image
	image := &schema.Image{
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.Tenant,
			Name:   imageName,
		},
		Spec: schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
		},
	}
	expectedImageMeta, err := builders.NewImageMetadataBuilder().
		Name(imageName).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectedImageSpec := &schema.ImageSpec{
		BlockStorageRef: *blockStorageRefObj,
		CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
	}
	stepsBuilder.CreateOrUpdateImageV1Step("Create an image", suite.RegionalClient.StorageV1, image,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			Metadata:      expectedImageMeta,
			Spec:          expectedImageSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created image
	imageTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   imageName,
	}
	image = stepsBuilder.GetImageV1Step("Get the created image", suite.RegionalClient.StorageV1, *imageTRef,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			Metadata:      expectedImageMeta,
			Spec:          expectedImageSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Block storage

	// Create a block storage
	block := &schema.BlockStorage{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.Tenant,
			Workspace: workspaceName,
			Name:      blockStorageName,
		},
		Spec: schema.BlockStorageSpec{
			SizeGB: initialStorageSize,
			SkuRef: *storageSkuRefObj,
		},
	}
	expectedBlockMeta, err := builders.NewBlockStorageMetadataBuilder().
		Name(blockStorageName).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectedBlockSpec := &schema.BlockStorageSpec{
		SizeGB: initialStorageSize,
		SkuRef: *storageSkuRefObj,
	}
	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create a block storage", suite.RegionalClient.StorageV1, block,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:      expectedBlockMeta,
			Spec:          expectedBlockSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created block storage
	blockWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      blockStorageName,
	}
	block = stepsBuilder.GetBlockStorageV1Step("Get the created block storage", suite.RegionalClient.StorageV1, *blockWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:      expectedBlockMeta,
			Spec:          expectedBlockSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Network

	// Create a network
	network := &schema.Network{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.Tenant,
			Workspace: workspaceName,
			Name:      networkName,
		},
		Spec: schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: &suite.NetworkCidr},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		},
	}
	expectNetworkMeta, err := builders.NewNetworkMetadataBuilder().
		Name(networkName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectNetworkSpec := &schema.NetworkSpec{
		Cidr:          schema.Cidr{Ipv4: &suite.NetworkCidr},
		SkuRef:        *networkSkuRefObj,
		RouteTableRef: *routeTableRefObj,
	}
	stepsBuilder.CreateOrUpdateNetworkV1Step("Create a network", suite.RegionalClient.NetworkV1, network,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:      expectNetworkMeta,
			Spec:          expectNetworkSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created network
	networkWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      networkName,
	}
	stepsBuilder.GetNetworkV1Step("Get the created network", suite.RegionalClient.NetworkV1, *networkWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:      expectNetworkMeta,
			Spec:          expectNetworkSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Internet gateway

	// Create an internet gateway
	gateway := &schema.InternetGateway{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.Tenant,
			Workspace: workspaceName,
			Name:      internetGatewayName,
		},
	}
	expectGatewayMeta, err := builders.NewInternetGatewayMetadataBuilder().
		Name(internetGatewayName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectGatewaySpec := &schema.InternetGatewaySpec{EgressOnly: ptr.To(false)}
	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create a internet gateway", suite.RegionalClient.NetworkV1, gateway,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Metadata:      expectGatewayMeta,
			Spec:          expectGatewaySpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created internet gateway
	gatewayWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      internetGatewayName,
	}
	stepsBuilder.GetInternetGatewayV1Step("Get the created internet gateway", suite.RegionalClient.NetworkV1, *gatewayWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Metadata:      expectGatewayMeta,
			Spec:          expectGatewaySpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Route table

	// Create a route table
	route := &schema.RouteTable{
		Metadata: &schema.RegionalNetworkResourceMetadata{
			Tenant:    suite.Tenant,
			Workspace: workspaceName,
			Network:   networkName,
			Name:      routeTableName,
		},
		Spec: schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
			},
		},
	}
	expectRouteMeta, err := builders.NewRouteTableMetadataBuilder().
		Name(routeTableName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Network(networkName).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectRouteSpec := &schema.RouteTableSpec{
		Routes: []schema.RouteSpec{
			{DestinationCidrBlock: constants.RouteTableDefaultDestination, TargetRef: *internetGatewayRefObj},
		},
	}
	stepsBuilder.CreateOrUpdateRouteTableV1Step("Create a route table", suite.RegionalClient.NetworkV1, route,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Metadata:      expectRouteMeta,
			Spec:          expectRouteSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created route table
	routeNRef := &secapi.NetworkReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Network:   secapi.NetworkID(networkName),
		Name:      routeTableName,
	}
	stepsBuilder.GetRouteTableV1Step("Get the created route table", suite.RegionalClient.NetworkV1, *routeNRef,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Metadata:      expectRouteMeta,
			Spec:          expectRouteSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Subnet

	// Create a subnet
	subnet := &schema.Subnet{
		Metadata: &schema.RegionalNetworkResourceMetadata{
			Tenant:    suite.Tenant,
			Workspace: workspaceName,
			Network:   networkName,
			Name:      subnetName,
		},
		Spec: schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: &subnetCidr},
			Zone: zone,
		},
	}
	expectSubnetMeta, err := builders.NewSubnetMetadataBuilder().
		Name(subnetName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Network(networkName).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectSubnetSpec := &schema.SubnetSpec{
		Cidr: schema.Cidr{Ipv4: &subnetCidr},
		Zone: zone,
	}
	stepsBuilder.CreateOrUpdateSubnetV1Step("Create a subnet", suite.RegionalClient.NetworkV1, subnet,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Metadata:      expectSubnetMeta,
			Spec:          expectSubnetSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created subnet
	subnetNRef := &secapi.NetworkReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Network:   secapi.NetworkID(networkName),
		Name:      subnetName,
	}
	stepsBuilder.GetSubnetV1Step("Get the created subnet", suite.RegionalClient.NetworkV1, *subnetNRef,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Metadata:      expectSubnetMeta,
			Spec:          expectSubnetSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Security Group

	// Create a security group
	group := &schema.SecurityGroup{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.Tenant,
			Workspace: workspaceName,
			Name:      securityGroupName,
		},
		Spec: schema.SecurityGroupSpec{
			Rules: []schema.SecurityGroupRuleSpec{
				{Direction: schema.SecurityGroupRuleDirectionIngress},
			},
		},
	}
	expectGroupMeta, err := builders.NewSecurityGroupMetadataBuilder().
		Name(securityGroupName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectGroupSpec := &schema.SecurityGroupSpec{
		Rules: []schema.SecurityGroupRuleSpec{
			{Direction: schema.SecurityGroupRuleDirectionIngress},
		},
	}
	stepsBuilder.CreateOrUpdateSecurityGroupV1Step("Create a security group", suite.RegionalClient.NetworkV1, group,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Metadata:      expectGroupMeta,
			Spec:          expectGroupSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created security group
	groupWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      securityGroupName,
	}
	stepsBuilder.GetSecurityGroupV1Step("Get the created security group", suite.RegionalClient.NetworkV1, *groupWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Metadata:      expectGroupMeta,
			Spec:          expectGroupSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Public ip

	// Create a public ip
	publicIp := &schema.PublicIp{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.Tenant,
			Workspace: workspaceName,
			Name:      publicIpName,
		},
		Spec: schema.PublicIpSpec{
			Address: &publicIpAddress1,
			Version: schema.IPVersionIPv4,
		},
	}
	expectPublicIpMeta, err := builders.NewPublicIpMetadataBuilder().
		Name(publicIpName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectPublicIpSpec := &schema.PublicIpSpec{
		Address: &publicIpAddress1,
		Version: schema.IPVersionIPv4,
	}
	stepsBuilder.CreateOrUpdatePublicIpV1Step("Create a public ip", suite.RegionalClient.NetworkV1, publicIp,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			Metadata:      expectPublicIpMeta,
			Spec:          expectPublicIpSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created public ip
	publicIpWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      publicIpName,
	}
	stepsBuilder.GetPublicIpV1Step("Get the created public ip", suite.RegionalClient.NetworkV1, *publicIpWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			Metadata:      expectPublicIpMeta,
			Spec:          expectPublicIpSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Nic

	// Create a nic
	nic := &schema.Nic{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.Tenant,
			Workspace: workspaceName,
			Name:      nicName,
		},
		Spec: schema.NicSpec{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		},
	}
	expectNicMeta, err := builders.NewNicMetadataBuilder().
		Name(nicName).
		Provider(constants.NetworkProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectNicSpec := &schema.NicSpec{
		Addresses:    []string{nicAddress1},
		PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
		SubnetRef:    *subnetRefObj,
	}
	stepsBuilder.CreateOrUpdateNicV1Step("Create a nic", suite.RegionalClient.NetworkV1, nic,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Metadata:      expectNicMeta,
			Spec:          expectNicSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created nic
	nicWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      nicName,
	}
	stepsBuilder.GetNicV1Step("Get the created nic", suite.RegionalClient.NetworkV1, *nicWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Metadata:      expectNicMeta,
			Spec:          expectNicSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Instance

	// Create an instance
	instance := &schema.Instance{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.Tenant,
			Workspace: workspaceName,
			Name:      instanceName,
		},
		Spec: schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		},
	}
	expectInstanceMeta, err := builders.NewInstanceMetadataBuilder().
		Name(instanceName).
		Provider(constants.ComputeProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectInstanceSpec := &schema.InstanceSpec{
		SkuRef: *instanceSkuRefObj,
		Zone:   instance.Spec.Zone,
		BootVolume: schema.VolumeReference{
			DeviceRef: *blockStorageRefObj,
		},
	}

	stepsBuilder.CreateOrUpdateInstanceV1Step("Create an instance", suite.RegionalClient.ComputeV1, instance,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:      expectInstanceMeta,
			Spec:          expectInstanceSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created instance
	instanceWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      instanceName,
	}

	instance = stepsBuilder.GetInstanceV1Step("Get the created instance", suite.RegionalClient.ComputeV1, *instanceWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:      expectInstanceMeta,
			Spec:          expectInstanceSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Resources deletion
	stepsBuilder.DeleteInstanceV1Step("Delete the instance", suite.RegionalClient.ComputeV1, instance)

	stepsBuilder.DeleteSecurityGroupV1Step("Delete the security group", suite.RegionalClient.NetworkV1, group)
	stepsBuilder.DeleteNicV1Step("Delete the nic", suite.RegionalClient.NetworkV1, nic)
	stepsBuilder.DeletePublicIpV1Step("Delete the public ip", suite.RegionalClient.NetworkV1, publicIp)
	stepsBuilder.DeleteSubnetV1Step("Delete the subnet", suite.RegionalClient.NetworkV1, subnet)
	stepsBuilder.DeleteRouteTableV1Step("Delete the route table", suite.RegionalClient.NetworkV1, route)
	stepsBuilder.DeleteInternetGatewayV1Step("Delete the internet gateway", suite.RegionalClient.NetworkV1, gateway)
	stepsBuilder.DeleteNetworkV1Step("Delete the network", suite.RegionalClient.NetworkV1, network)

	stepsBuilder.DeleteBlockStorageV1Step("Delete the block storage", suite.RegionalClient.StorageV1, block)
	stepsBuilder.DeleteImageV1Step("Delete the image", suite.RegionalClient.StorageV1, image)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.RegionalClient.WorkspaceV1, workspace)

	stepsBuilder.DeleteRoleAssignmentV1Step("Delete the role assignment", suite.GlobalClient.AuthorizationV1, roleAssign)
	stepsBuilder.DeleteRoleV1Step("Delete the role", suite.GlobalClient.AuthorizationV1, role)

	suite.FinishScenario()
}

func (suite *FoundationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
