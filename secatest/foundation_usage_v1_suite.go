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
	zone1 := suite.regionZones[rand.Intn(len(suite.regionZones))]

	// Select subs
	roleAssignmentSub1 := suite.users[rand.Intn(len(suite.users))]

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

	imageName := secalib.GenerateImageName()
	imageResource := secalib.GenerateImageResource(suite.tenant, imageName)

	instanceSkuRef := secalib.GenerateSkuRef(instanceSkuName)
	instanceSkuRefObj, err := secapi.BuildReferenceFromURN(instanceSkuRef)
	if err != nil {
		t.Fatal(err)
	}

	instanceName := secalib.GenerateInstanceName()
	instanceResource := secalib.GenerateInstanceResource(suite.tenant, workspaceName, instanceName)

	networkSkuRef := secalib.GenerateSkuRef(networkSkuName)
	networkSkuRefObj, err := secapi.BuildReferenceFromURN(networkSkuRef)
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

	blockStorageSize := secalib.GenerateBlockStorageSize()

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
					Subs:  []string{roleAssignmentSub1},
					Scopes: []schema.RoleAssignmentScope{
						{Tenants: &[]string{suite.tenant}},
					},
				},
			},
			Workspace: &mock.ResourceParams[schema.WorkspaceSpec]{
				Name: workspaceName,
			},
			BlockStorage: &mock.ResourceParams[schema.BlockStorageSpec]{
				Name: blockStorageName,
				InitialSpec: &schema.BlockStorageSpec{
					SkuRef: *storageSkuRefObj,
					SizeGB: blockStorageSize,
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
					Zone: zone1,
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
					Zone:   zone1,
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
	var roleResp *schema.Role
	var expectedRoleMeta *schema.GlobalTenantResourceMetadata
	var expectedRoleSpec *schema.RoleSpec

	t.WithNewStep("Create role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateRole")

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
		roleResp, err = suite.globalClient.AuthorizationV1.CreateOrUpdateRole(ctx, role)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, roleResp)

		expectedRoleMeta = secalib.NewGlobalTenantResourceMetadata(roleName, secalib.AuthorizationProviderV1, roleResource, secalib.ApiVersion1, secalib.RoleKind,
			suite.tenant)
		expectedRoleMeta.Verb = http.MethodPut
		suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedRoleMeta, roleResp.Metadata)

		expectedRoleSpec = &schema.RoleSpec{
			Permissions: []schema.Permission{
				{
					Provider:  secalib.StorageProviderV1,
					Resources: []string{imageResource},
					Verb:      []string{http.MethodGet},
				},
			},
		}
		suite.verifyRoleSpecStep(sCtx, expectedRoleSpec, &roleResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingResourceState, *roleResp.Status.State)
	})

	t.WithNewStep("Get created role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRole")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleName,
		}
		roleResp, err = suite.globalClient.AuthorizationV1.GetRole(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, roleResp)

		expectedRoleMeta.Verb = http.MethodGet
		suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedRoleMeta, roleResp.Metadata)

		suite.verifyRoleSpecStep(sCtx, expectedRoleSpec, &roleResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *roleResp.Status.State)
	})

	// Role assignment
	var assignResp *schema.RoleAssignment
	var expectedAssignMeta *schema.GlobalTenantResourceMetadata
	var expectedAssignSpec *schema.RoleAssignmentSpec

	t.WithNewStep("Create role assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateRoleAssignment")

		assign := &schema.RoleAssignment{
			Metadata: &schema.GlobalTenantResourceMetadata{
				Tenant: suite.tenant,
				Name:   roleAssignmentName,
			},
			Spec: schema.RoleAssignmentSpec{
				Roles:  []string{roleName},
				Subs:   []string{roleAssignmentSub1},
				Scopes: []schema.RoleAssignmentScope{{Tenants: &[]string{suite.tenant}}},
			},
		}
		assignResp, err = suite.globalClient.AuthorizationV1.CreateOrUpdateRoleAssignment(ctx, assign)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, assignResp)

		expectedAssignMeta = secalib.NewGlobalTenantResourceMetadata(roleAssignmentName, secalib.AuthorizationProviderV1, roleAssignmentResource, secalib.ApiVersion1, secalib.RoleAssignmentKind,
			suite.tenant)
		expectedAssignMeta.Verb = http.MethodPut
		suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedAssignMeta, assignResp.Metadata)

		expectedAssignSpec = &schema.RoleAssignmentSpec{
			Roles:  []string{roleName},
			Subs:   []string{roleAssignmentSub1},
			Scopes: []schema.RoleAssignmentScope{{Tenants: &[]string{suite.tenant}}},
		}
		suite.verifyRoleAssignmentSpecStep(sCtx, expectedAssignSpec, &assignResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingResourceState, *assignResp.Status.State)
	})

	t.WithNewStep("Get created role assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRoleAssignment")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleAssignmentName,
		}
		assignResp, err = suite.globalClient.AuthorizationV1.GetRoleAssignment(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, assignResp)

		expectedAssignMeta.Verb = http.MethodGet
		suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedAssignMeta, assignResp.Metadata)

		suite.verifyRoleAssignmentSpecStep(sCtx, expectedAssignSpec, &assignResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *assignResp.Status.State)
	})

	// Workspace
	var workspaceResp *schema.Workspace
	var expectedWorkspaceMeta *schema.RegionalResourceMetadata

	t.WithNewStep("Create workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "CreateOrUpdateWorkspace")

		ws := &schema.Workspace{
			Metadata: &schema.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   workspaceName,
			},
		}
		workspaceResp, err = suite.regionalClient.WorkspaceV1.CreateOrUpdateWorkspace(ctx, ws)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, workspaceResp)

		expectedWorkspaceMeta = secalib.NewRegionalResourceMetadata(workspaceName, secalib.WorkspaceProviderV1, workspaceResource, secalib.ApiVersion1, secalib.WorkspaceKind,
			suite.tenant, suite.region)
		expectedWorkspaceMeta.Verb = http.MethodPut
		suite.verifyRegionalResourceMetadataStep(sCtx, expectedWorkspaceMeta, workspaceResp.Metadata)

		suite.verifyStatusStep(sCtx, secalib.CreatingResourceState, *workspaceResp.Status.State)
	})

	t.WithNewStep("Get created workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "GetWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspaceName,
		}
		workspaceResp, err = suite.regionalClient.WorkspaceV1.GetWorkspace(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, workspaceResp)

		expectedWorkspaceMeta.Verb = http.MethodGet
		suite.verifyRegionalResourceMetadataStep(sCtx, expectedWorkspaceMeta, workspaceResp.Metadata)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *workspaceResp.Status.State)
	})

	// Image
	var imageResp *schema.Image
	var expectedImageMeta *schema.RegionalResourceMetadata
	var expectedImageSpec *schema.ImageSpec

	t.WithNewStep("Create image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateImage", "")

		img := &schema.Image{
			Metadata: &schema.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   imageName,
			},
			Spec: schema.ImageSpec{
				BlockStorageRef: *blockStorageRefObj,
				CpuArchitecture: secalib.CpuArchitectureAmd64,
			},
		}
		imageResp, err = suite.regionalClient.StorageV1.CreateOrUpdateImage(ctx, img)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMeta = secalib.NewRegionalResourceMetadata(imageName, secalib.StorageProviderV1, imageResource, secalib.ApiVersion1, secalib.ImageKind,
			suite.tenant, suite.region)
		expectedImageMeta.Verb = http.MethodPut
		suite.verifyRegionalResourceMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

		expectedImageSpec = &schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: secalib.CpuArchitectureAmd64,
		}
		suite.verifyImageSpecStep(sCtx, expectedImageSpec, &imageResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingResourceState, *imageResp.Status.State)
	})

	t.WithNewStep("Get created image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetImage", "")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   imageName,
		}
		imageResp, err = suite.regionalClient.StorageV1.GetImage(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMeta.Verb = http.MethodGet
		suite.verifyRegionalResourceMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

		suite.verifyImageSpecStep(sCtx, expectedImageSpec, &imageResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *imageResp.Status.State)
	})

	// Block storage
	var blockResp *schema.BlockStorage
	var expectedBlockMeta *schema.RegionalWorkspaceResourceMetadata
	var expectedBlockSpec *schema.BlockStorageSpec

	t.WithNewStep("Create block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateBlockStorage", workspaceName)

		bo := &schema.BlockStorage{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      blockStorageName,
			},
			Spec: schema.BlockStorageSpec{
				SizeGB: blockStorageSize,
				SkuRef: *storageSkuRefObj,
			},
		}
		blockResp, err = suite.regionalClient.StorageV1.CreateOrUpdateBlockStorage(ctx, bo)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMeta = secalib.NewRegionalWorkspaceResourceMetadata(blockStorageName, secalib.StorageProviderV1, blockStorageResource, secalib.ApiVersion1, secalib.BlockStorageKind,
			suite.tenant, workspaceName, suite.region)
		expectedBlockMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		expectedBlockSpec = &schema.BlockStorageSpec{
			SizeGB: blockStorageSize,
			SkuRef: *storageSkuRefObj,
		}
		suite.verifyBlockStorageSpecStep(sCtx, expectedBlockSpec, &blockResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingResourceState, *blockResp.Status.State)
	})

	t.WithNewStep("Get created block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetBlockStorage", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		blockResp, err = suite.regionalClient.StorageV1.GetBlockStorage(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		suite.verifyBlockStorageSpecStep(sCtx, expectedBlockSpec, &blockResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *blockResp.Status.State)
	})

	// Network
	var networkResp *schema.Network
	var expectedNetworkMeta *schema.RegionalWorkspaceResourceMetadata
	var expectedNetworkSpec *schema.NetworkSpec

	t.WithNewStep("Create network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNetwork", workspaceName)

		net := &schema.Network{
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
		networkResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateNetwork(ctx, net)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMeta = secalib.NewRegionalWorkspaceResourceMetadata(networkName, secalib.NetworkProviderV1, networkResource, secalib.ApiVersion1, secalib.NetworkKind,
			suite.tenant, workspaceName, suite.region)
		expectedNetworkMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		expectedNetworkSpec = &schema.NetworkSpec{
			Cidr:          schema.Cidr{Ipv4: &suite.networkCidr},
			SkuRef:        *networkSkuRefObj,
			RouteTableRef: *routeTableRefObj,
		}
		suite.verifyNetworkSpecStep(sCtx, expectedNetworkSpec, &networkResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingResourceState, *networkResp.Status.State)
	})

	t.WithNewStep("Get created network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNetwork", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      networkName,
		}
		networkResp, err = suite.regionalClient.NetworkV1.GetNetwork(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		suite.verifyNetworkSpecStep(sCtx, expectedNetworkSpec, &networkResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *networkResp.Status.State)
	})

	// Internet gateway
	var gatewayResp *schema.InternetGateway
	var expectedGatewayMeta *schema.RegionalWorkspaceResourceMetadata
	var expectedGatewaySpec *schema.InternetGatewaySpec

	t.WithNewStep("Create internet gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateInternetGateway", workspaceName)

		gtw := &schema.InternetGateway{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      internetGatewayName,
			},
		}
		gatewayResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateInternetGateway(ctx, gtw)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMeta = secalib.NewRegionalWorkspaceResourceMetadata(internetGatewayName, secalib.NetworkProviderV1, internetGatewayResource, secalib.ApiVersion1, secalib.InternetGatewayKind,
			suite.tenant, workspaceName, suite.region)
		expectedGatewayMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		expectedGatewaySpec = &schema.InternetGatewaySpec{EgressOnly: ptr.To(false)}
		suite.verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, &gatewayResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingResourceState, *gatewayResp.Status.State)
	})

	t.WithNewStep("Get created internet gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetInternetGateway", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      internetGatewayName,
		}
		gatewayResp, err = suite.regionalClient.NetworkV1.GetInternetGateway(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		suite.verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, &gatewayResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *gatewayResp.Status.State)
	})

	// Route table
	var routeResp *schema.RouteTable
	var expectedRouteMeta *schema.RegionalNetworkResourceMetadata
	var expectedRouteSpec *schema.RouteTableSpec

	t.WithNewStep("Create route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateRouteTable", workspaceName)

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
		routeResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateRouteTable(ctx, route)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMeta = secalib.NewRegionalNetworkResourceMetadata(routeTableName, secalib.NetworkProviderV1, routeTableResource, secalib.ApiVersion1, secalib.RouteTableKind,
			suite.tenant, workspaceName, networkName, suite.region)
		expectedRouteMeta.Verb = http.MethodPut
		suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		expectedRouteSpec = &schema.RouteTableSpec{
			Routes: []schema.RouteSpec{
				{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: *internetGatewayRefObj},
			},
		}
		suite.verifyRouteTableSpecStep(sCtx, expectedRouteSpec, &routeResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingResourceState, *routeResp.Status.State)
	})

	t.WithNewStep("Get created route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetRouteTable", workspaceName)

		nref := secapi.NetworkReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Network:   secapi.NetworkID(networkName),
			Name:      routeTableName,
		}
		routeResp, err = suite.regionalClient.NetworkV1.GetRouteTable(ctx, nref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMeta.Verb = http.MethodGet
		suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		suite.verifyRouteTableSpecStep(sCtx, expectedRouteSpec, &routeResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *routeResp.Status.State)
	})

	// Subnet
	var subnetResp *schema.Subnet
	var expectedSubnetMeta *schema.RegionalNetworkResourceMetadata
	var expectedSubnetSpec *schema.SubnetSpec

	t.WithNewStep("Create subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateSubnet", workspaceName)

		sub := &schema.Subnet{
			Metadata: &schema.RegionalNetworkResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Network:   networkName,
				Name:      subnetName,
			},
			Spec: schema.SubnetSpec{
				Cidr: schema.Cidr{Ipv4: &subnetCidr},
				Zone: zone1,
			},
		}
		subnetResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateSubnet(ctx, sub)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMeta = secalib.NewRegionalNetworkResourceMetadata(subnetName, secalib.NetworkProviderV1, subnetResource, secalib.ApiVersion1, secalib.SubnetKind,
			suite.tenant, workspaceName, networkName, suite.region)
		expectedSubnetMeta.Verb = http.MethodPut
		suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		expectedSubnetSpec = &schema.SubnetSpec{
			Cidr: schema.Cidr{Ipv4: &subnetCidr},
			Zone: zone1,
		}
		suite.verifySubNetSpecStep(sCtx, expectedSubnetSpec, &subnetResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingResourceState, *subnetResp.Status.State)
	})

	t.WithNewStep("Get created subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNetwork", workspaceName)

		wref := secapi.NetworkReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Network:   secapi.NetworkID(networkName),
			Name:      subnetName,
		}
		subnetResp, err = suite.regionalClient.NetworkV1.GetSubnet(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMeta.Verb = http.MethodGet
		suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		suite.verifySubNetSpecStep(sCtx, expectedSubnetSpec, &subnetResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *subnetResp.Status.State)
	})

	// Security Group
	var groupResp *schema.SecurityGroup
	var expectedGroupMeta *schema.RegionalWorkspaceResourceMetadata
	var expectedGroupSpec *schema.SecurityGroupSpec

	t.WithNewStep("Create security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateSecurityGroup", workspaceName)

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
		groupResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateSecurityGroup(ctx, group)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, groupResp)

		expectedGroupMeta = secalib.NewRegionalWorkspaceResourceMetadata(securityGroupName, secalib.NetworkProviderV1, securityGroupResource, secalib.ApiVersion1, secalib.SecurityGroupKind,
			suite.tenant, workspaceName, suite.region)
		expectedGroupMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		expectedGroupSpec = &schema.SecurityGroupSpec{
			Rules: []schema.SecurityGroupRuleSpec{
				{Direction: secalib.SecurityRuleDirectionIngress},
			},
		}
		suite.verifySecurityGroupSpecStep(sCtx, expectedGroupSpec, &groupResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingResourceState, *groupResp.Status.State)
	})

	t.WithNewStep("Get created security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSecurityGroup", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      securityGroupName,
		}
		groupResp, err = suite.regionalClient.NetworkV1.GetSecurityGroup(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, groupResp)

		expectedGroupMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		suite.verifySecurityGroupSpecStep(sCtx, expectedGroupSpec, &groupResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *groupResp.Status.State)
	})

	// Public ip
	var publicIpResp *schema.PublicIp
	var expectedIpMeta *schema.RegionalWorkspaceResourceMetadata
	var expectedIpSpec *schema.PublicIpSpec

	t.WithNewStep("Create public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdatePublicIp", workspaceName)

		ip := &schema.PublicIp{
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
		publicIpResp, err = suite.regionalClient.NetworkV1.CreateOrUpdatePublicIp(ctx, ip)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta = secalib.NewRegionalWorkspaceResourceMetadata(publicIpName, secalib.NetworkProviderV1, publicIpResource, secalib.ApiVersion1, secalib.PublicIpKind,
			suite.tenant, workspaceName, suite.region)
		expectedIpMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		expectedIpSpec = &schema.PublicIpSpec{
			Address: &publicIpAddress1,
			Version: secalib.IpVersion4,
		}
		suite.verifyPublicIpSpecStep(sCtx, expectedIpSpec, &publicIpResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingResourceState, *publicIpResp.Status.State)
	})

	t.WithNewStep("Get created public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetPublicIp", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      publicIpName,
		}
		publicIpResp, err = suite.regionalClient.NetworkV1.GetPublicIp(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		suite.verifyPublicIpSpecStep(sCtx, expectedIpSpec, &publicIpResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *publicIpResp.Status.State)
	})

	// Nic
	var nicResp *schema.Nic
	var expectedNicMeta *schema.RegionalWorkspaceResourceMetadata
	var expectedNicSpec *schema.NicSpec

	t.WithNewStep("Create nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNic", workspaceName)

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
		nicResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateNic(ctx, nic)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMeta = secalib.NewRegionalWorkspaceResourceMetadata(nicName, secalib.NetworkProviderV1, nicResource, secalib.ApiVersion1, secalib.NicKind,
			suite.tenant, workspaceName, suite.region)
		expectedNicMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		expectedNicSpec = &schema.NicSpec{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: &[]schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
		}
		suite.verifyNicSpecStep(sCtx, expectedNicSpec, &nicResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingResourceState, *nicResp.Status.State)
	})

	t.WithNewStep("Get created nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNic", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      nicName,
		}
		nicResp, err = suite.regionalClient.NetworkV1.GetNic(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		suite.verifyNicSpecStep(sCtx, expectedNicSpec, &nicResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *nicResp.Status.State)
	})

	// Instance
	var instanceResp *schema.Instance
	var expectedInstanceMeta *schema.RegionalWorkspaceResourceMetadata
	var expectedInstanceSpec *schema.InstanceSpec

	t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "CreateOrUpdateInstance", workspaceName)

		inst := &schema.Instance{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      instanceName,
			},
			Spec: schema.InstanceSpec{
				SkuRef: *instanceSkuRefObj,
				Zone:   zone1,
				BootVolume: schema.VolumeReference{
					DeviceRef: *blockStorageRefObj,
				},
			},
		}

		instanceResp, err = suite.regionalClient.ComputeV1.CreateOrUpdateInstance(ctx, inst)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedInstanceMeta = secalib.NewRegionalWorkspaceResourceMetadata(instanceName, secalib.ComputeProviderV1, instanceResource, secalib.ApiVersion1, secalib.InstanceKind,
			suite.tenant, workspaceName, suite.region)
		expectedInstanceMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedInstanceMeta, instanceResp.Metadata)

		expectedInstanceSpec = &schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone1,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		}
		suite.verifyInstanceSpecStep(sCtx, expectedInstanceSpec, &instanceResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingResourceState, *instanceResp.Status.State)
	})

	t.WithNewStep("Get created instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		instanceResp, err = suite.regionalClient.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedInstanceMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedInstanceMeta, instanceResp.Metadata)

		suite.verifyInstanceSpecStep(sCtx, expectedInstanceSpec, &instanceResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *instanceResp.Status.State)
	})

	// Delete All

	// Delete instance
	t.WithNewStep("Delete instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "DeleteInstance", workspaceName)

		err = suite.regionalClient.ComputeV1.DeleteInstance(ctx, instanceResp)
		requireNoError(sCtx, err)
	})

	// Delete security group
	t.WithNewStep("Delete security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteSecurityGroup", workspaceName)

		err = suite.regionalClient.NetworkV1.DeleteSecurityGroup(ctx, groupResp)
		requireNoError(sCtx, err)
	})

	// Delete Nic
	t.WithNewStep("Delete Nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteNic", workspaceName)

		err = suite.regionalClient.NetworkV1.DeleteNic(ctx, nicResp)
		requireNoError(sCtx, err)
	})

	// Delete PublicIp
	t.WithNewStep("Delete Public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeletePublicIp", workspaceName)

		err = suite.regionalClient.NetworkV1.DeletePublicIp(ctx, publicIpResp)
		requireNoError(sCtx, err)
	})

	// Delete subnet
	t.WithNewStep("Delete Subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteSubnet", workspaceName)

		err = suite.regionalClient.NetworkV1.DeleteSubnet(ctx, subnetResp)
		requireNoError(sCtx, err)
	})

	// Delete Route-table
	t.WithNewStep("Delete Route-table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteRouteTable", workspaceName)

		err = suite.regionalClient.NetworkV1.DeleteRouteTable(ctx, routeResp)
		requireNoError(sCtx, err)
	})

	// Delete Internet-gateway
	t.WithNewStep("Delete Internet-gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteInternetGateway", workspaceName)

		err = suite.regionalClient.NetworkV1.DeleteInternetGateway(ctx, gatewayResp)
		requireNoError(sCtx, err)
	})

	// Delete Network
	t.WithNewStep("Delete Network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteNetwork", workspaceName)

		err = suite.regionalClient.NetworkV1.DeleteNetwork(ctx, networkResp)
		requireNoError(sCtx, err)
	})

	// Delete BlockStorage
	t.WithNewStep("Delete BlockStorage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "DeleteBlockStorage", workspaceName)

		err = suite.regionalClient.StorageV1.DeleteBlockStorage(ctx, blockResp)
		requireNoError(sCtx, err)
	})

	// Delete Image
	t.WithNewStep("Delete Image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "DeleteImage", "")

		err = suite.regionalClient.StorageV1.DeleteImage(ctx, imageResp)
		requireNoError(sCtx, err)
	})

	// Delete Workspace
	t.WithNewStep("Delete Workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "DeleteWorkspace")

		err = suite.regionalClient.WorkspaceV1.DeleteWorkspace(ctx, workspaceResp)
		requireNoError(sCtx, err)
	})

	// Delete Role Assignment
	t.WithNewStep("Delete Role Assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteRoleAssignment")

		err = suite.globalClient.AuthorizationV1.DeleteRoleAssignment(ctx, assignResp)
		requireNoError(sCtx, err)
	})

	// Delete Role
	t.WithNewStep("Delete Role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteRole")

		err = suite.globalClient.AuthorizationV1.DeleteRole(ctx, roleResp)
		requireNoError(sCtx, err)
	})

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *FoundationUsageV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
