package secatest

import (
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/builders"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/generators"
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
	var err error
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t,
		authorizationProviderV1,
		string(schema.GlobalTenantResourceMetadataKindResourceKindRole),
		string(schema.GlobalTenantResourceMetadataKindResourceKindRoleAssignment),
		workspaceProviderV1,
		string(schema.RegionalResourceMetadataKindResourceKindWorkspace),
		storageProviderV1,
		string(schema.RegionalResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalResourceMetadataKindResourceKindImage),
		networkProviderV1,
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
		computeProviderV1,
		string(schema.RegionalResourceMetadataKindResourceKindInstance),
	)

	// Generate the subnet cidr
	subnetCidr, err := generators.GenerateSubnetCidr(suite.networkCidr, 8, 1)
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
	publicIpAddress1, err := generators.GeneratePublicIp(suite.publicIpsRange, 1)
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
	workspaceName := generators.GenerateWorkspaceName()

	roleName := generators.GenerateRoleName()

	roleAssignmentName := generators.GenerateRoleAssignmentName()

	storageSkuRef := generators.GenerateSkuRef(storageSkuName)
	storageSkuRefObj, err := secapi.BuildReferenceFromURN(storageSkuRef)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageRef := generators.GenerateBlockStorageRef(blockStorageName)
	blockStorageRefObj, err := secapi.BuildReferenceFromURN(blockStorageRef)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}
	initialStorageSize := generators.GenerateBlockStorageSize()

	imageName := generators.GenerateImageName()
	imageResource := generators.GenerateImageResource(suite.tenant, imageName)

	instanceSkuRef := generators.GenerateSkuRef(instanceSkuName)
	instanceSkuRefObj, err := secapi.BuildReferenceFromURN(instanceSkuRef)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	instanceName := generators.GenerateInstanceName()

	networkSkuRef1 := generators.GenerateSkuRef(networkSkuName)
	networkSkuRefObj, err := secapi.BuildReferenceFromURN(networkSkuRef1)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	networkName := generators.GenerateNetworkName()

	internetGatewayName := generators.GenerateInternetGatewayName()
	internetGatewayRef := generators.GenerateInternetGatewayRef(internetGatewayName)
	internetGatewayRefObj, err := secapi.BuildReferenceFromURN(internetGatewayRef)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	routeTableName := generators.GenerateRouteTableName()
	routeTableRef := generators.GenerateRouteTableRef(routeTableName)
	routeTableRefObj, err := secapi.BuildReferenceFromURN(routeTableRef)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	subnetName := generators.GenerateSubnetName()
	subnetRef := generators.GenerateSubnetRef(subnetName)
	subnetRefObj, err := secapi.BuildReferenceFromURN(subnetRef)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	nicName := generators.GenerateNicName()

	publicIpName := generators.GeneratePublicIpName()
	publicIpRef := generators.GeneratePublicIpRef(publicIpName)
	publicIpRefObj, err := secapi.BuildReferenceFromURN(publicIpRef)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	securityGroupName := generators.GenerateSecurityGroupName()

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
						{Provider: storageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
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
					envLabel: envDevelopmentLabel,
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
					CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
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
					Version: schema.IPVersionIPv4,
					Address: ptr.To(publicIpAddress1),
				},
			},
			SecurityGroup: &mock.ResourceParams[schema.SecurityGroupSpec]{
				Name: securityGroupName,
				InitialSpec: &schema.SecurityGroupSpec{
					Rules: []schema.SecurityGroupRuleSpec{{Direction: schema.SecurityGroupRuleDirectionIngress}},
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
					Provider:  storageProviderV1,
					Resources: []string{imageResource},
					Verb:      []string{http.MethodGet},
				},
			},
		},
	}
	expectRoleMeta, err := builders.NewRoleMetadataBuilder().
		Name(roleName).
		Provider(authorizationProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectRoleSpec := &schema.RoleSpec{
		Permissions: []schema.Permission{
			{
				Provider:  storageProviderV1,
				Resources: []string{imageResource},
				Verb:      []string{http.MethodGet},
			},
		},
	}
	suite.createOrUpdateRoleV1Step("Create a role", t, suite.globalClient.AuthorizationV1, role,
		responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			metadata:      expectRoleMeta,
			spec:          expectRoleSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created role
	roleTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   roleName,
	}
	role = suite.getRoleV1Step("Get the created role", t, suite.globalClient.AuthorizationV1, *roleTRef,
		responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			metadata:      expectRoleMeta,
			spec:          expectRoleSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

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
	expectRoleAssignMeta, err := builders.NewRoleAssignmentMetadataBuilder().
		Name(roleAssignmentName).
		Provider(authorizationProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectRoleAssignSpec := &schema.RoleAssignmentSpec{
		Roles:  []string{roleName},
		Subs:   []string{roleAssignmentSub},
		Scopes: []schema.RoleAssignmentScope{{Tenants: &[]string{suite.tenant}}},
	}
	suite.createOrUpdateRoleAssignmentV1Step("Create a role assignment", t, suite.globalClient.AuthorizationV1, roleAssign,
		responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			metadata:      expectRoleAssignMeta,
			spec:          expectRoleAssignSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created role assignment
	roleAssignTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   roleAssignmentName,
	}
	roleAssign = suite.getRoleAssignmentV1Step("Get the created role assignment", t, suite.globalClient.AuthorizationV1, *roleAssignTRef,
		responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			metadata:      expectRoleAssignMeta,
			spec:          expectRoleAssignSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Workspace

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			envLabel: envDevelopmentLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   workspaceName,
		},
	}
	expectWorkspaceMeta, err := builders.NewWorkspaceMetadataBuilder().
		Name(workspaceName).
		Provider(workspaceProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectWorkspaceLabels := schema.Labels{envLabel: envDevelopmentLabel}
	suite.createOrUpdateWorkspaceV1Step("Create a workspace", t, suite.regionalClient.WorkspaceV1, workspace,
		responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			labels:        expectWorkspaceLabels,
			metadata:      expectWorkspaceMeta,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created Workspace
	tref := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   workspaceName,
	}
	workspace = suite.getWorkspaceV1Step("Get the created workspace", t, suite.regionalClient.WorkspaceV1, *tref,
		responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			labels:        expectWorkspaceLabels,
			metadata:      expectWorkspaceMeta,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Image

	// Create an image
	image := &schema.Image{
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   imageName,
		},
		Spec: schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
		},
	}
	expectedImageMeta, err := builders.NewImageMetadataBuilder().
		Name(imageName).
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectedImageSpec := &schema.ImageSpec{
		BlockStorageRef: *blockStorageRefObj,
		CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
	}
	suite.createOrUpdateImageV1Step("Create an image", t, suite.regionalClient.StorageV1, image,
		responseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			metadata:      expectedImageMeta,
			spec:          expectedImageSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created image
	imageTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   imageName,
	}
	image = suite.getImageV1Step("Get the created image", t, suite.regionalClient.StorageV1, *imageTRef,
		responseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			metadata:      expectedImageMeta,
			spec:          expectedImageSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

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
	expectedBlockMeta, err := builders.NewBlockStorageMetadataBuilder().
		Name(blockStorageName).
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectedBlockSpec := &schema.BlockStorageSpec{
		SizeGB: initialStorageSize,
		SkuRef: *storageSkuRefObj,
	}
	suite.createOrUpdateBlockStorageV1Step("Create a block storage", t, suite.regionalClient.StorageV1, block,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			metadata:      expectedBlockMeta,
			spec:          expectedBlockSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created block storage
	blockWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      blockStorageName,
	}
	block = suite.getBlockStorageV1Step("Get the created block storage", t, suite.regionalClient.StorageV1, *blockWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			metadata:      expectedBlockMeta,
			spec:          expectedBlockSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

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
	expectNetworkMeta, err := builders.NewNetworkMetadataBuilder().
		Name(networkName).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectNetworkSpec := &schema.NetworkSpec{
		Cidr:          schema.Cidr{Ipv4: &suite.networkCidr},
		SkuRef:        *networkSkuRefObj,
		RouteTableRef: *routeTableRefObj,
	}
	suite.createOrUpdateNetworkV1Step("Create a network", t, suite.regionalClient.NetworkV1, network,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			metadata:      expectNetworkMeta,
			spec:          expectNetworkSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created network
	networkWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      networkName,
	}
	suite.getNetworkV1Step("Get the created network", t, suite.regionalClient.NetworkV1, *networkWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			metadata:      expectNetworkMeta,
			spec:          expectNetworkSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Internet gateway

	// Create an internet gateway
	gateway := &schema.InternetGateway{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Name:      internetGatewayName,
		},
	}
	expectGatewayMeta, err := builders.NewInternetGatewayMetadataBuilder().
		Name(internetGatewayName).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectGatewaySpec := &schema.InternetGatewaySpec{EgressOnly: ptr.To(false)}
	suite.createOrUpdateInternetGatewayV1Step("Create a internet gateway", t, suite.regionalClient.NetworkV1, gateway,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			metadata:      expectGatewayMeta,
			spec:          expectGatewaySpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created internet gateway
	gatewayWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      internetGatewayName,
	}
	suite.getInternetGatewayV1Step("Get the created internet gateway", t, suite.regionalClient.NetworkV1, *gatewayWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			metadata:      expectGatewayMeta,
			spec:          expectGatewaySpec,
			resourceState: schema.ResourceStateActive,
		},
	)

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
	expectRouteMeta, err := builders.NewRouteTableMetadataBuilder().
		Name(routeTableName).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Network(networkName).Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectRouteSpec := &schema.RouteTableSpec{
		Routes: []schema.RouteSpec{
			{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *internetGatewayRefObj},
		},
	}
	suite.createOrUpdateRouteTableV1Step("Create a route table", t, suite.regionalClient.NetworkV1, route,
		responseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			metadata:      expectRouteMeta,
			spec:          expectRouteSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created route table
	routeNRef := &secapi.NetworkReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Network:   secapi.NetworkID(networkName),
		Name:      routeTableName,
	}
	suite.getRouteTableV1Step("Get the created route table", t, suite.regionalClient.NetworkV1, *routeNRef,
		responseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			metadata:      expectRouteMeta,
			spec:          expectRouteSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

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
	expectSubnetMeta, err := builders.NewSubnetMetadataBuilder().
		Name(subnetName).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Network(networkName).Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectSubnetSpec := &schema.SubnetSpec{
		Cidr: schema.Cidr{Ipv4: &subnetCidr},
		Zone: zone,
	}
	suite.createOrUpdateSubnetV1Step("Create a subnet", t, suite.regionalClient.NetworkV1, subnet,
		responseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			metadata:      expectSubnetMeta,
			spec:          expectSubnetSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created subnet
	subnetNRef := &secapi.NetworkReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Network:   secapi.NetworkID(networkName),
		Name:      subnetName,
	}
	suite.getSubnetV1Step("Get the created subnet", t, suite.regionalClient.NetworkV1, *subnetNRef,
		responseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			metadata:      expectSubnetMeta,
			spec:          expectSubnetSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

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
				{Direction: schema.SecurityGroupRuleDirectionIngress},
			},
		},
	}
	expectGroupMeta, err := builders.NewSecurityGroupMetadataBuilder().
		Name(securityGroupName).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectGroupSpec := &schema.SecurityGroupSpec{
		Rules: []schema.SecurityGroupRuleSpec{
			{Direction: schema.SecurityGroupRuleDirectionIngress},
		},
	}
	suite.createOrUpdateSecurityGroupV1Step("Create a security group", t, suite.regionalClient.NetworkV1, group,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			metadata:      expectGroupMeta,
			spec:          expectGroupSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created security group
	groupWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      securityGroupName,
	}
	suite.getSecurityGroupV1Step("Get the created security group", t, suite.regionalClient.NetworkV1, *groupWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			metadata:      expectGroupMeta,
			spec:          expectGroupSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

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
			Version: schema.IPVersionIPv4,
		},
	}
	expectPublicIpMeta, err := builders.NewPublicIpMetadataBuilder().
		Name(publicIpName).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectPublicIpSpec := &schema.PublicIpSpec{
		Address: &publicIpAddress1,
		Version: schema.IPVersionIPv4,
	}
	suite.createOrUpdatePublicIpV1Step("Create a public ip", t, suite.regionalClient.NetworkV1, publicIp,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			metadata:      expectPublicIpMeta,
			spec:          expectPublicIpSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created public ip
	publicIpWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      publicIpName,
	}
	suite.getPublicIpV1Step("Get the created public ip", t, suite.regionalClient.NetworkV1, *publicIpWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			metadata:      expectPublicIpMeta,
			spec:          expectPublicIpSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

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
	expectNicMeta, err := builders.NewNicMetadataBuilder().
		Name(nicName).
		Provider(networkProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectNicSpec := &schema.NicSpec{
		Addresses:    []string{nicAddress1},
		PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
		SubnetRef:    *subnetRefObj,
	}
	suite.createOrUpdateNicV1Step("Create a nic", t, suite.regionalClient.NetworkV1, nic,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			metadata:      expectNicMeta,
			spec:          expectNicSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created nic
	nicWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      nicName,
	}
	suite.getNicV1Step("Get the created nic", t, suite.regionalClient.NetworkV1, *nicWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			metadata:      expectNicMeta,
			spec:          expectNicSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

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
	expectInstanceMeta, err := builders.NewInstanceMetadataBuilder().
		Name(instanceName).
		Provider(computeProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectInstanceSpec := &schema.InstanceSpec{
		SkuRef: *instanceSkuRefObj,
		Zone:   instance.Spec.Zone,
		BootVolume: schema.VolumeReference{
			DeviceRef: *blockStorageRefObj,
		},
	}

	suite.createOrUpdateInstanceV1Step("Create an instance", t, suite.regionalClient.ComputeV1, instance,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			metadata:      expectInstanceMeta,
			spec:          expectInstanceSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created instance
	instanceWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      instanceName,
	}

	instance = suite.getInstanceV1Step("Get the created instance", t, suite.regionalClient.ComputeV1, *instanceWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			metadata:      expectInstanceMeta,
			spec:          expectInstanceSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Resources deletion

	suite.deleteInstanceV1Step("Delete the instance", t, suite.regionalClient.ComputeV1, instance)

	suite.deleteSecurityGroupV1Step("Delete the security group", t, suite.regionalClient.NetworkV1, group)

	suite.deleteNicV1Step("Delete the nic", t, suite.regionalClient.NetworkV1, nic)

	suite.deletePublicIpV1Step("Delete the public ip", t, suite.regionalClient.NetworkV1, publicIp)

	suite.deleteSubnetV1Step("Delete the subnet", t, suite.regionalClient.NetworkV1, subnet)

	suite.deleteRouteTableV1Step("Delete the route table", t, suite.regionalClient.NetworkV1, route)

	suite.deleteInternetGatewayV1Step("Delete the internet gateway", t, suite.regionalClient.NetworkV1, gateway)

	suite.deleteNetworkV1Step("Delete the network", t, suite.regionalClient.NetworkV1, network)

	suite.deleteBlockStorageV1Step("Delete the block storage", t, suite.regionalClient.StorageV1, block)

	suite.deleteImageV1Step("Delete the image", t, suite.regionalClient.StorageV1, image)

	suite.deleteWorkspaceV1Step("Delete the workspace", t, suite.regionalClient.WorkspaceV1, workspace)

	suite.deleteRoleAssignmentV1Step("Delete the role assignment", t, suite.globalClient.AuthorizationV1, roleAssign)

	suite.deleteRoleV1Step("Delete the role", t, suite.globalClient.AuthorizationV1, role)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *FoundationUsageV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
