package secatest

import (
	"context"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"k8s.io/utils/ptr"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type FoundationUsageV1TestSuite struct {
	mixedTestSuite

	users          []string
	networkCidr    string
	publicIpsRange string
	regionZones    []string
	storageSkus    []string
	instanceSkus   []string
	networkSkus    []string
}

func (suite *FoundationUsageV1TestSuite) TestSuite(t provider.T) {
	ctx := context.Background()
	var err error
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t,
		secalib.AuthorizationProviderV1, secalib.RoleKind, secalib.RoleAssignmentKind,
		secalib.WorkspaceProviderV1, secalib.WorkspaceKind,
		secalib.StorageProviderV1, secalib.BlockStorageKind, secalib.ImageKind,
		secalib.NetworkProviderV1, secalib.NetworkKind, secalib.InternetGatewayKind, secalib.NicKind, secalib.PublicIpKind, secalib.RouteTableKind, secalib.SubnetKind, secalib.SecurityGroupKind,
		secalib.ComputeProviderV1, secalib.InstanceKind,
	)

	// Generate the subnet cidr
	subnetCidr, err := secalib.GenerateSubnetCidr(suite.networkCidr, 8, 1)
	if err != nil {
		slog.Error("Failed to generate subnet cidr", "error", err)
		return
	}

	// Generate the nic addresses
	nicAddress1, err := secalib.GenerateNicAddress(subnetCidr, 1)
	if err != nil {
		slog.Error("Failed to generate nic address", "error", err)
		return
	}

	// Generate the public ips
	publicIpAddress1, err := secalib.GeneratePublicIp(suite.publicIpsRange, 1)
	if err != nil {
		slog.Error("Failed to generate public ip", "error", err)
		return
	}

	// Select zones
	zone := suite.regionZones[rand.Intn(len(suite.regionZones))]

	// Select subs
	roleAssignmentSub := suite.users[rand.Intn(len(suite.users))]

	// Select skus
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]
	instanceSkuName := suite.instanceSkus[rand.Intn(len(suite.instanceSkus))]
	networkSkuName := suite.networkSkus[rand.Intn(len(suite.networkSkus))]

	// Generate scenario data
	workspaceName := secalib.GenerateWorkspaceName()
	workspaceResource := secalib.GenerateWorkspaceResource(suite.tenant, workspaceName)

	roleName := secalib.GenerateRoleName()
	roleResource := secalib.GenerateRoleResource(suite.tenant, roleName)

	roleAssignmentName := secalib.GenerateRoleAssignmentName()
	roleAssignmentResource := secalib.GenerateRoleAssignmentResource(suite.tenant, roleAssignmentName)

	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)
	storageSkuRefObj, err := secapi.BuildReferenceFromURN(storageSkuRef)
	if err != nil {
		t.Fatal(err)
	}

	blockStorageName := secalib.GenerateBlockStorageName()
	blockStorageResource := secalib.GenerateBlockStorageResource(suite.tenant, workspaceName, blockStorageName)
	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName)
	blockStorageRefObj, err := secapi.BuildReferenceFromURN(blockStorageRef)
	if err != nil {
		t.Fatal(err)
	}
	initialStorageSize := secalib.GenerateBlockStorageSize()

	imageName := secalib.GenerateImageName()
	imageResource := secalib.GenerateImageResource(suite.tenant, imageName)

	instanceSkuRef := secalib.GenerateSkuRef(instanceSkuName)
	instanceSkuRefObj, err := secapi.BuildReferenceFromURN(instanceSkuRef)
	if err != nil {
		t.Fatal(err)
	}

	instanceName := secalib.GenerateInstanceName()
	instanceResource := secalib.GenerateInstanceResource(suite.tenant, workspaceName, instanceName)

	networkSkuRef1 := secalib.GenerateSkuRef(networkSkuName)
	networkSkuRefObj, err := secapi.BuildReferenceFromURN(networkSkuRef1)
	if err != nil {
		t.Fatal(err)
	}

	networkName := secalib.GenerateNetworkName()
	networkResource := secalib.GenerateNetworkResource(suite.tenant, workspaceName, networkName)

	internetGatewayName := secalib.GenerateInternetGatewayName()
	internetGatewayResource := secalib.GenerateInternetGatewayResource(suite.tenant, workspaceName, internetGatewayName)
	internetGatewayRef := secalib.GenerateInternetGatewayRef(internetGatewayName)
	internetGatewayRefObj, err := secapi.BuildReferenceFromURN(internetGatewayRef)
	if err != nil {
		t.Fatal(err)
	}

	routeTableName := secalib.GenerateRouteTableName()
	routeTableResource := secalib.GenerateRouteTableResource(suite.tenant, workspaceName, networkName, routeTableName)
	routeTableRef := secalib.GenerateRouteTableRef(routeTableName)
	routeTableRefObj, err := secapi.BuildReferenceFromURN(routeTableRef)
	if err != nil {
		t.Fatal(err)
	}

	subnetName := secalib.GenerateSubnetName()
	subnetResource := secalib.GenerateSubnetResource(suite.tenant, workspaceName, networkName, subnetName)
	subnetRef := secalib.GenerateSubnetRef(subnetName)
	subnetRefObj, err := secapi.BuildReferenceFromURN(subnetRef)
	if err != nil {
		t.Fatal(err)
	}

	nicName := secalib.GenerateNicName()
	nicResource := secalib.GenerateNicResource(suite.tenant, workspaceName, nicName)

	publicIpName := secalib.GeneratePublicIpName()
	publicIpResource := secalib.GeneratePublicIpResource(suite.tenant, workspaceName, publicIpName)
	publicIpRef := secalib.GeneratePublicIpRef(publicIpName)
	publicIpRefObj, err := secapi.BuildReferenceFromURN(publicIpRef)
	if err != nil {
		t.Fatal(err)
	}

	securityGroupName := secalib.GenerateSecurityGroupName()
	securityGroupResource := secalib.GenerateSecurityGroupResource(suite.tenant, workspaceName, securityGroupName)

	// Setup mock, if configured to use
	if suite.mockEnabled {
		mockParams := &mock.FoundationUsageParamsV1{
			Params: &mock.Params{
				MockURL:   *suite.mockServerURL,
				AuthToken: suite.authToken,
				Tenant:    suite.tenant,
				Region:    suite.region,
			},
			Role: &mock.ResourceParams[schema.RoleSpec]{
				Name: roleName,
				InitialSpec: &schema.RoleSpec{
					Permissions: []schema.Permission{
						{Provider: secalib.StorageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
					},
				},
			},
			RoleAssignment: &mock.ResourceParams[schema.RoleAssignmentSpec]{
				Name: roleAssignmentName,
				InitialSpec: &schema.RoleAssignmentSpec{
					Roles: []string{roleName},
					Subs:  []string{roleAssignmentSub},
					Scopes: []schema.RoleAssignmentScope{
						{Tenants: &[]string{suite.tenant}},
					},
				},
			},
			Workspace: &mock.ResourceParams[schema.WorkspaceSpec]{
				Name: workspaceName,
				InitialLabels: schema.Labels{
					secalib.EnvLabel: secalib.EnvDevelopmentLabel,
				},
			},
			BlockStorage: &mock.ResourceParams[schema.BlockStorageSpec]{
				Name: blockStorageName,
				InitialSpec: &schema.BlockStorageSpec{
					SkuRef: *storageSkuRefObj,
					SizeGB: initialStorageSize,
				},
			},
			Image: &mock.ResourceParams[schema.ImageSpec]{
				Name: imageName,
				InitialSpec: &schema.ImageSpec{
					BlockStorageRef: *blockStorageRefObj,
					CpuArchitecture: secalib.CpuArchitectureAmd64,
				},
			},
			Network: &mock.ResourceParams[schema.NetworkSpec]{
				Name: networkName,
				InitialSpec: &schema.NetworkSpec{
					Cidr:          schema.Cidr{Ipv4: ptr.To(suite.networkCidr)},
					SkuRef:        *networkSkuRefObj,
					RouteTableRef: *routeTableRefObj,
				},
			},
			InternetGateway: &mock.ResourceParams[schema.InternetGatewaySpec]{
				Name:        internetGatewayName,
				InitialSpec: &schema.InternetGatewaySpec{EgressOnly: ptr.To(false)},
			},
			RouteTable: &mock.ResourceParams[schema.RouteTableSpec]{
				Name: routeTableName,
				InitialSpec: &schema.RouteTableSpec{
					Routes: []schema.RouteSpec{
						{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *internetGatewayRefObj},
					},
				},
			},
			Subnet: &mock.ResourceParams[schema.SubnetSpec]{
				Name: subnetName,
				InitialSpec: &schema.SubnetSpec{
					Cidr: schema.Cidr{Ipv4: &subnetCidr},
					Zone: zone,
				},
			},
			Nic: &mock.ResourceParams[schema.NicSpec]{
				Name: nicName,
				InitialSpec: &schema.NicSpec{
					Addresses:    []string{nicAddress1},
					PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
					SubnetRef:    *subnetRefObj,
				},
			},
			PublicIp: &mock.ResourceParams[schema.PublicIpSpec]{
				Name: publicIpName,
				InitialSpec: &schema.PublicIpSpec{
					Version: secalib.IpVersion4,
					Address: ptr.To(publicIpAddress1),
				},
			},
			SecurityGroup: &mock.ResourceParams[schema.SecurityGroupSpec]{
				Name: securityGroupName,
				InitialSpec: &schema.SecurityGroupSpec{
					Rules: []schema.SecurityGroupRuleSpec{{Direction: secalib.SecurityRuleDirectionIngress}},
				},
			},
			Instance: &mock.ResourceParams[schema.InstanceSpec]{
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
		wm, err := mock.ConfigFoundationUsageScenario(suite.scenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}

		suite.mockClient = wm
	}

	// Role
	// Create a role
	role := &schema.Role{
		Metadata: &schema.GlobalTenantResourceMetadata{
			Tenant: suite.tenant,
			Name:   roleName,
		},
		Spec: schema.RoleSpec{
			Permissions: []schema.Permission{
				{
					Provider:  secalib.StorageProviderV1,
					Resources: []string{imageResource},
					Verb:      []string{http.MethodGet},
				},
			},
		},
	}
	expectRoleMeta := secalib.NewGlobalTenantResourceMetadata(roleName,
		secalib.AuthorizationProviderV1,
		roleResource,
		secalib.ApiVersion1,
		secalib.RoleKind,
		suite.tenant)
	expectRoleSpec := &schema.RoleSpec{
		Permissions: []schema.Permission{
			{
				Provider:  secalib.StorageProviderV1,
				Resources: []string{imageResource},
				Verb:      []string{http.MethodGet},
			},
		},
	}
	suite.createOrUpdateRoleV1Step("Create a role", t, ctx, suite.globalClient.AuthorizationV1, role,
		expectRoleMeta, expectRoleSpec, secalib.CreatingResourceState)

	// Get the created role
	roleTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   roleName,
	}
	role = suite.getRoleV1Step("Get the created role", t, ctx, suite.globalClient.AuthorizationV1, *roleTRef,
		expectRoleMeta, expectRoleSpec, secalib.ActiveResourceState)

	// Role assignment
	// Create a role assignment
	roleAssign := &schema.RoleAssignment{
		Metadata: &schema.GlobalTenantResourceMetadata{
			Tenant: suite.tenant,
			Name:   roleAssignmentName,
		},
		Spec: schema.RoleAssignmentSpec{
			Roles:  []string{roleName},
			Subs:   []string{roleAssignmentSub},
			Scopes: []schema.RoleAssignmentScope{{Tenants: &[]string{suite.tenant}}},
		},
	}
	expectRoleAssignMeta := secalib.NewGlobalTenantResourceMetadata(roleAssignmentName,
		secalib.AuthorizationProviderV1,
		roleAssignmentResource,
		secalib.ApiVersion1,
		secalib.RoleAssignmentKind,
		suite.tenant)
	expectRoleAssignSpec := &schema.RoleAssignmentSpec{
		Roles:  []string{roleName},
		Subs:   []string{roleAssignmentSub},
		Scopes: []schema.RoleAssignmentScope{{Tenants: &[]string{suite.tenant}}},
	}
	suite.createOrUpdateRoleAssignmentV1Step("Create a role assignment", t, ctx, suite.globalClient.AuthorizationV1, roleAssign,
		expectRoleAssignMeta, expectRoleAssignSpec, secalib.CreatingResourceState)

	// Get the created role assignment
	roleAssignTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   roleAssignmentName,
	}
	roleAssign = suite.getRoleAssignmentV1Step("Get the created role assignment", t, ctx, suite.globalClient.AuthorizationV1, *roleAssignTRef,
		expectRoleAssignMeta, expectRoleAssignSpec, secalib.ActiveResourceState)

	// Workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			secalib.EnvLabel: secalib.EnvDevelopmentLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   workspaceName,
		},
	}
	expectMeta := secalib.NewRegionalResourceMetadata(workspaceName,
		secalib.WorkspaceProviderV1,
		workspaceResource,
		secalib.ApiVersion1,
		secalib.WorkspaceKind,
		suite.tenant, suite.region)
	expectLabels := schema.Labels{secalib.EnvLabel: secalib.EnvDevelopmentLabel}
	suite.createOrUpdateWorkspaceV1Step("Create a workspace", t, ctx, suite.regionalClient.WorkspaceV1, workspace,
		expectMeta, expectLabels, secalib.CreatingResourceState)

	// Get the created Workspace
	tref := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   workspaceName,
	}
	workspace = suite.getWorkspaceV1Step("Get the created workspace", t, ctx, suite.regionalClient.WorkspaceV1, *tref,
		expectMeta, expectLabels, secalib.ActiveResourceState)

	// Storage

	// Image
	image := &schema.Image{
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   imageName,
		},
		Spec: schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: secalib.CpuArchitectureAmd64,
		},
	}
	expectedImageMeta := secalib.NewRegionalResourceMetadata(imageName,
		secalib.StorageProviderV1,
		imageResource,
		secalib.ApiVersion1,
		secalib.ImageKind,
		suite.tenant,
		suite.region)
	expectedImageSpec := &schema.ImageSpec{
		BlockStorageRef: *blockStorageRefObj,
		CpuArchitecture: secalib.CpuArchitectureAmd64,
	}
	suite.createOrUpdateImageV1Step("Create an image", t, ctx, suite.regionalClient.StorageV1, image,
		expectedImageMeta, expectedImageSpec, secalib.CreatingResourceState)

	// Get the created image
	imageTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   imageName,
	}
	image = suite.getImageV1Step("Get the created image", t, ctx, suite.regionalClient.StorageV1, *imageTRef,
		expectedImageMeta, expectedImageSpec, secalib.ActiveResourceState)

	// Block storage
	// Create a block storage
	block := &schema.BlockStorage{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Name:      blockStorageName,
		},
		Spec: schema.BlockStorageSpec{
			SizeGB: initialStorageSize,
			SkuRef: *storageSkuRefObj,
		},
	}
	expectedBlockMeta := secalib.NewRegionalWorkspaceResourceMetadata(blockStorageName,
		secalib.StorageProviderV1,
		blockStorageResource,
		secalib.ApiVersion1,
		secalib.BlockStorageKind,
		suite.tenant,
		workspaceName,
		suite.region)
	expectedBlockSpec := &schema.BlockStorageSpec{
		SizeGB: initialStorageSize,
		SkuRef: *storageSkuRefObj,
	}
	suite.createOrUpdateBlockStorageV1Step("Create a block storage", t, ctx, suite.regionalClient.StorageV1, block,
		expectedBlockMeta, expectedBlockSpec, secalib.CreatingResourceState)

	// Get the created block storage
	blockWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      blockStorageName,
	}
	block = suite.getBlockStorageV1Step("Get the created block storage", t, ctx, suite.regionalClient.StorageV1, *blockWRef,
		expectedBlockMeta, expectedBlockSpec, secalib.ActiveResourceState)

	// Network
	// Create a network
	network := &schema.Network{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Name:      networkName,
		},
		Spec: schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: &suite.networkCidr},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		},
	}
	expectNetworkMeta := secalib.NewRegionalWorkspaceResourceMetadata(networkName,
		secalib.NetworkProviderV1,
		networkResource,
		secalib.ApiVersion1,
		secalib.NetworkKind,
		suite.tenant,
		workspaceName,
		suite.region,
	)
	expectNetworkSpec := &schema.NetworkSpec{
		Cidr:          schema.Cidr{Ipv4: &suite.networkCidr},
		SkuRef:        *networkSkuRefObj,
		RouteTableRef: *routeTableRefObj,
	}
	suite.createOrUpdateNetworkV1Step("Create a network", t, ctx, suite.regionalClient.NetworkV1, network,
		expectNetworkMeta, expectNetworkSpec, secalib.CreatingResourceState)

	// Get the created network
	networkWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      networkName,
	}
	suite.getNetworkV1Step("Get the created network", t, ctx, suite.regionalClient.NetworkV1, *networkWRef,
		expectNetworkMeta, expectNetworkSpec, secalib.ActiveResourceState)

	// Internet gateway

	// Create an internet gateway
	gateway := &schema.InternetGateway{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Name:      internetGatewayName,
		},
	}
	expectGatewayMeta := secalib.NewRegionalWorkspaceResourceMetadata(internetGatewayName,
		secalib.NetworkProviderV1,
		internetGatewayResource,
		secalib.ApiVersion1,
		secalib.InternetGatewayKind,
		suite.tenant,
		workspaceName,
		suite.region,
	)
	expectGatewaySpec := &schema.InternetGatewaySpec{EgressOnly: ptr.To(false)}
	suite.createOrUpdateInternetGatewayV1Step("Create a internet gateway", t, ctx, suite.regionalClient.NetworkV1, gateway,
		expectGatewayMeta, expectGatewaySpec, secalib.CreatingResourceState)

	// Get the created internet gateway
	gatewayWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      internetGatewayName,
	}
	suite.getInternetGatewayV1Step("Get the created internet gateway", t, ctx, suite.regionalClient.NetworkV1, *gatewayWRef,
		expectGatewayMeta, expectGatewaySpec, secalib.ActiveResourceState)

	// Route table

	// Create a route table
	route := &schema.RouteTable{
		Metadata: &schema.RegionalNetworkResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Network:   networkName,
			Name:      routeTableName,
		},
		Spec: schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *internetGatewayRefObj},
			},
		},
	}
	expectRouteMeta := secalib.NewRegionalNetworkResourceMetadata(routeTableName,
		secalib.NetworkProviderV1,
		routeTableResource,
		secalib.ApiVersion1,
		secalib.RouteTableKind,
		suite.tenant,
		workspaceName,
		networkName,
		suite.region)
	expectRouteSpec := &schema.RouteTableSpec{
		Routes: []schema.RouteSpec{
			{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *internetGatewayRefObj},
		},
	}
	suite.createOrUpdateRouteTableV1Step("Create a route table", t, ctx, suite.regionalClient.NetworkV1, route,
		expectRouteMeta, expectRouteSpec, secalib.CreatingResourceState)

	// Get the created route table
	routeNRef := &secapi.NetworkReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Network:   secapi.NetworkID(networkName),
		Name:      routeTableName,
	}
	suite.getRouteTableV1Step("Get the created route table", t, ctx, suite.regionalClient.NetworkV1, *routeNRef,
		expectRouteMeta, expectRouteSpec, secalib.ActiveResourceState)

	// Subnet

	// Create a subnet
	subnet := &schema.Subnet{
		Metadata: &schema.RegionalNetworkResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Network:   networkName,
			Name:      subnetName,
		},
		Spec: schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: &subnetCidr},
			Zone: zone,
		},
	}
	expectSubnetMeta := secalib.NewRegionalNetworkResourceMetadata(subnetName,
		secalib.NetworkProviderV1,
		subnetResource,
		secalib.ApiVersion1,
		secalib.SubnetKind,
		suite.tenant,
		workspaceName,
		networkName,
		suite.region,
	)
	expectSubnetSpec := &schema.SubnetSpec{
		Cidr: schema.Cidr{Ipv4: &subnetCidr},
		Zone: zone,
	}
	suite.createOrUpdateSubnetV1Step("Create a subnet", t, ctx, suite.regionalClient.NetworkV1, subnet,
		expectSubnetMeta, expectSubnetSpec, secalib.CreatingResourceState)

	// Get the created subnet
	subnetNRef := &secapi.NetworkReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Network:   secapi.NetworkID(networkName),
		Name:      subnetName,
	}
	suite.getSubnetV1Step("Get the created subnet", t, ctx, suite.regionalClient.NetworkV1, *subnetNRef,
		expectSubnetMeta, expectSubnetSpec, secalib.ActiveResourceState)

	// Security Group

	// Create a security group
	group := &schema.SecurityGroup{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Name:      securityGroupName,
		},
		Spec: schema.SecurityGroupSpec{
			Rules: []schema.SecurityGroupRuleSpec{
				{Direction: secalib.SecurityRuleDirectionIngress},
			},
		},
	}
	expectGroupMeta := secalib.NewRegionalWorkspaceResourceMetadata(securityGroupName,
		secalib.NetworkProviderV1,
		securityGroupResource,
		secalib.ApiVersion1,
		secalib.SecurityGroupKind,
		suite.tenant,
		workspaceName,
		suite.region,
	)
	expectGroupSpec := &schema.SecurityGroupSpec{
		Rules: []schema.SecurityGroupRuleSpec{
			{Direction: secalib.SecurityRuleDirectionIngress},
		},
	}
	suite.createOrUpdateSecurityGroupV1Step("Create a security group", t, ctx, suite.regionalClient.NetworkV1, group,
		expectGroupMeta, expectGroupSpec, secalib.CreatingResourceState)

	// Get the created security group
	groupWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      securityGroupName,
	}
	suite.getSecurityGroupV1Step("Get the created security group", t, ctx, suite.regionalClient.NetworkV1, *groupWRef,
		expectGroupMeta, expectGroupSpec, secalib.ActiveResourceState)

	// Public ip

	// Create a public ip
	publicIp := &schema.PublicIp{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Name:      publicIpName,
		},
		Spec: schema.PublicIpSpec{
			Address: &publicIpAddress1,
			Version: secalib.IpVersion4,
		},
	}
	expectPublicIpMeta := secalib.NewRegionalWorkspaceResourceMetadata(publicIpName,
		secalib.NetworkProviderV1,
		publicIpResource,
		secalib.ApiVersion1,
		secalib.PublicIpKind,
		suite.tenant,
		workspaceName,
		suite.region,
	)
	expectPublicIpSpec := &schema.PublicIpSpec{
		Address: &publicIpAddress1,
		Version: secalib.IpVersion4,
	}
	suite.createOrUpdatePublicIpV1Step("Create a public ip", t, ctx, suite.regionalClient.NetworkV1, publicIp,
		expectPublicIpMeta, expectPublicIpSpec, secalib.CreatingResourceState)

	// Get the created public ip
	publicIpWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      publicIpName,
	}
	suite.getPublicIpV1Step("Get the created public ip", t, ctx, suite.regionalClient.NetworkV1, *publicIpWRef,
		expectPublicIpMeta, expectPublicIpSpec, secalib.ActiveResourceState)

	// Nic

	// Create a nic
	nic := &schema.Nic{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Name:      nicName,
		},
		Spec: schema.NicSpec{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		},
	}
	expectNicMeta := secalib.NewRegionalWorkspaceResourceMetadata(nicName,
		secalib.NetworkProviderV1,
		nicResource,
		secalib.ApiVersion1,
		secalib.NicKind,
		suite.tenant, workspaceName, suite.region)
	expectNicSpec := &schema.NicSpec{
		Addresses:    []string{nicAddress1},
		PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
		SubnetRef:    *subnetRefObj,
	}
	suite.createOrUpdateNicV1Step("Create a nic", t, ctx, suite.regionalClient.NetworkV1, nic,
		expectNicMeta, expectNicSpec, secalib.CreatingResourceState)

	// Get the created nic
	nicWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      nicName,
	}
	suite.getNicV1Step("Get the created nic", t, ctx, suite.regionalClient.NetworkV1, *nicWRef,
		expectNicMeta, expectNicSpec, secalib.ActiveResourceState)

	// Instance
	// Create an instance
	instance := &schema.Instance{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.tenant,
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

	expectInstanceMeta := secalib.NewRegionalWorkspaceResourceMetadata(instanceName,
		secalib.ComputeProviderV1,
		instanceResource,
		secalib.ApiVersion1,
		secalib.InstanceKind,
		suite.tenant,
		workspaceName,
		suite.region)
	expectInstanceSpec := &schema.InstanceSpec{
		SkuRef: *instanceSkuRefObj,
		Zone:   instance.Spec.Zone,
		BootVolume: schema.VolumeReference{
			DeviceRef: *blockStorageRefObj,
		},
	}

	suite.createOrUpdateInstanceV1Step("Create an instance", t, ctx, suite.regionalClient.ComputeV1, instance, expectInstanceMeta, expectInstanceSpec, secalib.CreatingResourceState)

	// Get the created instance
	instanceWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      instanceName,
	}

	instance = suite.getInstanceV1Step("Get the created instance", t, ctx, suite.regionalClient.ComputeV1, *instanceWRef, expectInstanceMeta, expectInstanceSpec, secalib.ActiveResourceState)

	// Delete All

	// Delete instance
	suite.deleteInstanceV1Step("Delete the instance", t, ctx, suite.regionalClient.ComputeV1, instance)

	// Delete security group
	suite.deleteSecurityGroupV1Step("Delete the security group", t, ctx, suite.regionalClient.NetworkV1, group)

	// Delete Nic
	suite.deleteNicV1Step("Delete the nic", t, ctx, suite.regionalClient.NetworkV1, nic)

	// Delete PublicIp
	suite.deletePublicIpV1Step("Delete the public ip", t, ctx, suite.regionalClient.NetworkV1, publicIp)

	// Delete subnet
	suite.deleteSubnetV1Step("Delete the subnet", t, ctx, suite.regionalClient.NetworkV1, subnet)

	// Delete Route-table
	suite.deleteRouteTableV1Step("Delete the route table", t, ctx, suite.regionalClient.NetworkV1, route)

	// Delete Internet-gateway
	suite.deleteInternetGatewayV1Step("Delete the internet gateway", t, ctx, suite.regionalClient.NetworkV1, gateway)

	// Delete Network
	suite.deleteNetworkV1Step("Delete the network", t, ctx, suite.regionalClient.NetworkV1, network)

	// Delete BlockStorage
	suite.deleteBlockStorageV1Step("Delete the block storage", t, ctx, suite.regionalClient.StorageV1, block)

	// Delete Image
	suite.deleteImageV1Step("Delete the image", t, ctx, suite.regionalClient.StorageV1, image)

	// Delete Workspace
	suite.deleteWorkspaceV1Step("Delete the workspace", t, ctx, suite.regionalClient.WorkspaceV1, workspace)

	// Delete Role Assignment
	suite.deleteRoleAssignmentV1Step("Delete the role assignment", t, ctx, suite.globalClient.AuthorizationV1, roleAssign)

	// Delete Role
	suite.deleteRoleV1Step("Delete the role", t, ctx, suite.globalClient.AuthorizationV1, role)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *FoundationUsageV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
