package secatest

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type UsagesV1TestSuite struct {
	mixedTestSuite

	users          []string
	networkCidr    string
	publicIpsRange string
	regionZones    []string
	storageSkus    []string
	instanceSkus   []string
	networkSkus    []string

	networkTestSuite *NetworkV1TestSuite
}

func (suite *UsagesV1TestSuite) TestUsagesV1(t provider.T) {
	slog.Info("Starting Usages V1 Test")

	t.Title("Usages Lifecycle Test")
	configureTags(t,
		secalib.AuthorizationProviderV1, secalib.RoleKind, secalib.RoleAssignmentKind,
		secalib.WorkspaceProviderV1, secalib.WorkspaceKind,
		secalib.StorageProviderV1, secalib.BlockStorageKind, secalib.ImageKind,
		secalib.NetworkProviderV1, secalib.NetworkKind, secalib.InternetGatewayKind, secalib.NicKind, secalib.PublicIPKind, secalib.RouteTableKind, secalib.SubnetKind, secalib.SecurityGroupKind,
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
	networkSkuName1 := suite.networkSkus[rand.Intn(len(suite.networkSkus))]
	networkSkuName2 := suite.networkSkus[rand.Intn(len(suite.networkSkus))]

	// Generate scenario data
	workspaceName := secalib.GenerateWorkspaceName()
	workspaceResource := secalib.GenerateWorkspaceResource(suite.tenant, workspaceName)

	roleName := secalib.GenerateRoleName()
	roleResource := secalib.GenerateRoleResource(suite.tenant, roleName)

	roleAssignmentName := secalib.GenerateRoleAssignmentName()
	roleAssignmentResource := secalib.GenerateRoleAssignmentResource(suite.tenant, roleAssignmentName)

	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)

	blockStorageName := secalib.GenerateBlockStorageName()
	blockStorageResource := secalib.GenerateBlockStorageResource(suite.tenant, workspaceName, blockStorageName)
	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName)

	imageName := secalib.GenerateImageName()
	imageResource := secalib.GenerateImageResource(suite.tenant, imageName)

	instanceSkuRef := secalib.GenerateSkuRef(instanceSkuName)
	instanceName := secalib.GenerateInstanceName()
	instanceResource := secalib.GenerateInstanceResource(suite.tenant, workspaceName, instanceName)

	networkSkuRef1 := secalib.GenerateSkuRef(networkSkuName1)
	networkSkuRef2 := secalib.GenerateSkuRef(networkSkuName2)

	networkName := secalib.GenerateNetworkName()
	networkResource := secalib.GenerateNetworkResource(suite.tenant, workspaceName, networkName)
	networkRef := secalib.GenerateNetworkRef(networkName)

	internetGatewayName := secalib.GenerateInternetGatewayName()
	internetGatewayResource := secalib.GenerateInternetGatewayResource(suite.tenant, workspaceName, internetGatewayName)
	internetGatewayRef := secalib.GenerateInternetGatewayRef(internetGatewayName)

	routeTableName := secalib.GenerateRouteTableName()
	routeTableResource := secalib.GenerateRouteTableResource(suite.tenant, workspaceName, networkName, routeTableName)
	routeTableRef := secalib.GenerateRouteTableRef(routeTableName)

	subnetName := secalib.GenerateSubnetName()
	subnetResource := secalib.GenerateSubnetResource(suite.tenant, workspaceName, networkName, subnetName)
	subnetRef := secalib.GenerateSubnetRef(subnetName)

	nicName := secalib.GenerateNicName()
	nicResource := secalib.GenerateNicResource(suite.tenant, workspaceName, nicName)

	publicIPName := secalib.GeneratePublicIPName()
	publicIPResource := secalib.GeneratePublicIPResource(suite.tenant, workspaceName, publicIPName)
	publicIPRef := secalib.GeneratePublicIPRef(publicIPName)

	securityGroupName := secalib.GenerateSecurityGroupName()
	securityGroupResource := secalib.GenerateSecurityGroupResource(suite.tenant, workspaceName, securityGroupName)

	blockStorageSize := secalib.GenerateBlockStorageSize()

	// Setup mock, if configured to use
	if suite.isMockEnabled() {
		wm, err := mock.CreateUsageScenario(fmt.Sprintf("Usages V1 Test %d", rand.Intn(100)), mock.UsageParamsV1{
			Params: &mock.Params{
				MockURL:   *suite.mockServerURL,
				AuthToken: suite.authToken,
				Tenant:    suite.tenant,
				Region:    suite.region,
			},
			Authorization: &mock.AuthorizationParamsV1{
				Role: &mock.ResourceParams[secalib.RoleSpecV1]{
					Name: roleName,
					InitialSpec: &secalib.RoleSpecV1{
						Permissions: []*secalib.RoleSpecPermissionV1{
							{
								Provider:  secalib.StorageProviderV1,
								Resources: []string{imageResource},
								Verb:      []string{http.MethodGet},
							},
						},
					},
				},
				RoleAssignment: &mock.ResourceParams[secalib.RoleAssignmentSpecV1]{
					Name: roleAssignmentName,
					InitialSpec: &secalib.RoleAssignmentSpecV1{
						Roles:  []string{roleName},
						Subs:   []string{roleAssignmentSub1},
						Scopes: []*secalib.RoleAssignmentSpecScopeV1{{Tenants: []string{suite.tenant}}},
					},
				},
			},
			Workspace: &mock.WorkspaceParamsV1{
				Workspace: &mock.ResourceParams[secalib.WorkspaceSpecV1]{
					Name: workspaceName,
				},
			},
			Storage: &mock.StorageParamsV1{
				BlockStorage: &mock.ResourceParams[secalib.BlockStorageSpecV1]{
					Name: blockStorageName,
					InitialSpec: &secalib.BlockStorageSpecV1{
						SkuRef: storageSkuRef,
						SizeGB: blockStorageSize,
					},
				},
				Image: &mock.ResourceParams[secalib.ImageSpecV1]{
					Name: imageName,
					InitialSpec: &secalib.ImageSpecV1{
						BlockStorageRef: blockStorageRef,
						CpuArchitecture: secalib.CpuArchitectureAmd64,
					},
				},
			},
			Network: &mock.NetworkParamsV1{
				Network: &mock.ResourceParams[secalib.NetworkSpecV1]{
					Name: networkName,
					InitialSpec: &secalib.NetworkSpecV1{
						Cidr:          &secalib.NetworkSpecCIDRV1{Ipv4: suite.networkCidr},
						SkuRef:        networkSkuRef1,
						RouteTableRef: routeTableRef,
					},
					UpdatedSpec: &secalib.NetworkSpecV1{
						Cidr:          &secalib.NetworkSpecCIDRV1{Ipv4: suite.networkCidr},
						SkuRef:        networkSkuRef2,
						RouteTableRef: routeTableRef,
					},
				},
				InternetGateway: &mock.ResourceParams[secalib.InternetGatewaySpecV1]{
					Name:        internetGatewayName,
					InitialSpec: &secalib.InternetGatewaySpecV1{EgressOnly: false},
				},
				RouteTable: &mock.ResourceParams[secalib.RouteTableSpecV1]{
					Name: routeTableName,
					InitialSpec: &secalib.RouteTableSpecV1{
						LocalRef: networkRef,
						Routes: []*secalib.RouteTableRouteV1{
							{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: internetGatewayRef},
						},
					},
				},
				Subnet: &mock.ResourceParams[secalib.SubnetSpecV1]{
					Name: subnetName,
					InitialSpec: &secalib.SubnetSpecV1{
						Cidr: &secalib.SubnetSpecCIDRV1{Ipv4: subnetCidr},
						Zone: zone1,
					},
				},
				NIC: &mock.ResourceParams[secalib.NICSpecV1]{
					Name: nicName,
					InitialSpec: &secalib.NICSpecV1{
						Addresses:    []string{nicAddress1},
						PublicIpRefs: []string{publicIPRef},
						SubnetRef:    subnetRef,
					},
				},
				PublicIP: &mock.ResourceParams[secalib.PublicIpSpecV1]{
					Name: publicIPName,
					InitialSpec: &secalib.PublicIpSpecV1{
						Version: secalib.IPVersion4,
						Address: publicIpAddress1,
					},
				},
				SecurityGroup: &mock.ResourceParams[secalib.SecurityGroupSpecV1]{
					Name: securityGroupName,
					InitialSpec: &secalib.SecurityGroupSpecV1{
						Rules: []*secalib.SecurityGroupRuleV1{{Direction: secalib.SecurityRuleDirectionIngress}},
					},
				},
			},
			Compute: &mock.ComputeParamsV1{
				Instance: &mock.ResourceParams[secalib.InstanceSpecV1]{
					Name: instanceName,
					InitialSpec: &secalib.InstanceSpecV1{
						SkuRef:        instanceSkuRef,
						Zone:          zone1,
						BootDeviceRef: blockStorageRef,
					},
				},
			},
		})
		if err != nil {
			t.Fatalf("Failed to create wiremock scenario: %v", err)
		}

		suite.mockClient = wm
	}
	ctx := context.Background()
	var roleResp *schema.Role
	var assignResp *schema.RoleAssignment

	// Role
	var expectedRoleMeta *secalib.Metadata
	var expectedRoleSpec *secalib.RoleSpecV1

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

		expectedRoleMeta = &secalib.Metadata{
			Name:       roleName,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleKind,
			Tenant:     suite.tenant,
		}
		verifyAuthorizationMetadataStep(sCtx, expectedRoleMeta, roleResp.Metadata)

		expectedRoleSpec = &secalib.RoleSpecV1{
			Permissions: []*secalib.RoleSpecPermissionV1{
				{
					Provider:  secalib.StorageProviderV1,
					Resources: []string{imageResource},
					Verb:      []string{http.MethodGet},
				},
			},
		}
		verifyRoleSpecStep(sCtx, expectedRoleSpec, &roleResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*roleResp.Status.State)},
		)
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
		verifyAuthorizationMetadataStep(sCtx, expectedRoleMeta, roleResp.Metadata)

		verifyRoleSpecStep(sCtx, expectedRoleSpec, &roleResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*roleResp.Status.State)},
		)
	})
	// Role assignment
	var expectedAssignMeta *secalib.Metadata
	var expectedAssignSpec *secalib.RoleAssignmentSpecV1

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

		expectedAssignMeta = &secalib.Metadata{
			Name:       roleAssignmentName,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleAssignmentResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleAssignmentKind,
			Tenant:     suite.tenant,
		}
		verifyAuthorizationMetadataStep(sCtx, expectedAssignMeta, assignResp.Metadata)

		expectedAssignSpec = &secalib.RoleAssignmentSpecV1{
			Roles:  []string{roleName},
			Subs:   []string{roleAssignmentSub1},
			Scopes: []*secalib.RoleAssignmentSpecScopeV1{{Tenants: []string{suite.tenant}}},
		}
		verifyRoleAssignmentSpecStep(sCtx, expectedAssignSpec, &assignResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*assignResp.Status.State)},
		)
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
		verifyAuthorizationMetadataStep(sCtx, expectedAssignMeta, assignResp.Metadata)

		verifyRoleAssignmentSpecStep(sCtx, expectedAssignSpec, &assignResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*assignResp.Status.State)},
		)
	})

	// Workspace
	var workspaceResp *schema.Workspace
	var expectedMeta *secalib.Metadata
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

		expectedMeta = &secalib.Metadata{
			Name:       workspaceName,
			Provider:   secalib.WorkspaceProviderV1,
			Resource:   workspaceResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.WorkspaceKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}
		verifyWorkspaceMetadataStep(sCtx, expectedMeta, workspaceResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*workspaceResp.Status.State)},
		)
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

		expectedMeta.Verb = http.MethodGet
		verifyWorkspaceMetadataStep(sCtx, expectedMeta, workspaceResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*workspaceResp.Status.State)},
		)
	})

	// Storage

	// Image
	var imageResp *schema.Image
	var expectedImageMeta *secalib.Metadata
	var expectedImageSpec *secalib.ImageSpecV1

	t.WithNewStep("Create image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateImage", "")

		blockStorageURN, err := secapi.BuildReferenceFromURN(blockStorageRef)
		if err != nil {
			t.Fatal(err)
		}

		img := &schema.Image{
			Metadata: &schema.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   imageName,
			},
			Spec: schema.ImageSpec{
				BlockStorageRef: *blockStorageURN,
				CpuArchitecture: secalib.CpuArchitectureAmd64,
			},
		}
		imageResp, err = suite.regionalClient.StorageV1.CreateOrUpdateImage(ctx, img)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMeta = &secalib.Metadata{
			Name:       imageName,
			Provider:   secalib.StorageProviderV1,
			Resource:   imageResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.ImageKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

		expectedImageSpec = &secalib.ImageSpecV1{
			BlockStorageRef: blockStorageRef,
			CpuArchitecture: secalib.CpuArchitectureAmd64,
		}
		verifyImageSpecStep(sCtx, expectedImageSpec, imageResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*imageResp.Status.State)},
		)
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
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

		verifyImageSpecStep(sCtx, expectedImageSpec, imageResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*imageResp.Status.State)},
		)
	})

	// Block storage
	var blockResp *schema.BlockStorage
	var expectedBlockMeta *secalib.Metadata
	var expectedBlockSpec *secalib.BlockStorageSpecV1

	t.WithNewStep("Create block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateBlockStorage", workspaceName)

		storageSkuURN, err := secapi.BuildReferenceFromURN(storageSkuRef)
		if err != nil {
			t.Fatal(err)
		}

		bo := &schema.BlockStorage{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      blockStorageName,
			},
			Spec: schema.BlockStorageSpec{
				SizeGB: blockStorageSize,
				SkuRef: *storageSkuURN,
			},
		}
		blockResp, err = suite.regionalClient.StorageV1.CreateOrUpdateBlockStorage(ctx, bo)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMeta = &secalib.Metadata{
			Name:       blockStorageName,
			Provider:   secalib.StorageProviderV1,
			Resource:   blockStorageResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.BlockStorageKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}

		verifyStorageWorkspaceMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		expectedBlockSpec = &secalib.BlockStorageSpecV1{
			SizeGB: blockStorageSize,
			SkuRef: storageSkuRef,
		}
		verifyBlockStorageSpecStep(sCtx, expectedBlockSpec, blockResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
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
		verifyStorageWorkspaceMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		verifyBlockStorageSpecStep(sCtx, expectedBlockSpec, blockResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
	})
	// Network
	var networkResp *schema.Network
	var expectedNetworkMeta *secalib.Metadata
	var expectedNetworkSpec *secalib.NetworkSpecV1

	t.WithNewStep("Create network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNetwork", workspaceName)

		networkSkuURN, err := secapi.BuildReferenceFromURN(networkSkuRef1)
		if err != nil {
			t.Fatal(err)
		}

		routeTableURN, err := secapi.BuildReferenceFromURN(routeTableRef)
		if err != nil {
			t.Fatal(err)
		}

		net := &schema.Network{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      networkName,
			},
			Spec: schema.NetworkSpec{
				Cidr:          schema.Cidr{Ipv4: &suite.networkCidr},
				SkuRef:        *networkSkuURN,
				RouteTableRef: *routeTableURN,
			},
		}
		networkResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateNetwork(ctx, net)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMeta = &secalib.Metadata{
			Name:       networkName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   networkResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.NetworkKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}
		suite.networkTestSuite.verifyWorkspaceMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		expectedNetworkSpec = &secalib.NetworkSpecV1{
			Cidr:          &secalib.NetworkSpecCIDRV1{Ipv4: suite.networkCidr},
			SkuRef:        networkSkuRef1,
			RouteTableRef: routeTableRef,
		}
		suite.networkTestSuite.verifyNetworkSpecStep(sCtx, expectedNetworkSpec, networkResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
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
		suite.networkTestSuite.verifyWorkspaceMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		suite.networkTestSuite.verifyNetworkSpecStep(sCtx, expectedNetworkSpec, networkResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
	})

	// Internet gateway
	var gatewayResp *schema.InternetGateway
	var expectedGatewayMeta *secalib.Metadata
	var expectedGatewaySpec *secalib.InternetGatewaySpecV1

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

		expectedGatewayMeta = &secalib.Metadata{
			Name:       internetGatewayName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   internetGatewayResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InternetGatewayKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}
		suite.networkTestSuite.verifyWorkspaceMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		expectedGatewaySpec = &secalib.InternetGatewaySpecV1{
			EgressOnly: false,
		}
		suite.networkTestSuite.verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, gatewayResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
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
		suite.networkTestSuite.verifyWorkspaceMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		suite.networkTestSuite.verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, gatewayResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
	})

	// Route table
	var routeResp *schema.RouteTable
	var expectedRouteMeta *secalib.Metadata
	var expectedRouteSpec *secalib.RouteTableSpecV1

	t.WithNewStep("Create route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateRouteTable", workspaceName)

		route := &schema.RouteTable{
			Metadata: &schema.RegionalNetworkResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Network:   networkName,
				Name:      routeTableName,
			},
		}
		routeResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateRouteTable(ctx, route)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMeta = &secalib.Metadata{
			Name:       routeTableName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   routeTableResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RouteTableKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}
		suite.networkTestSuite.verifyNetworkMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		expectedRouteSpec = &secalib.RouteTableSpecV1{
			LocalRef: networkRef,
			Routes: []*secalib.RouteTableRouteV1{
				{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: internetGatewayRef},
			},
		}
		suite.networkTestSuite.verifyRouteTableSpecStep(sCtx, expectedRouteSpec, routeResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
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
		suite.networkTestSuite.verifyNetworkMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		suite.networkTestSuite.verifyRouteTableSpecStep(sCtx, expectedRouteSpec, routeResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
	})

	// Subnet
	var subnetResp *schema.Subnet
	var expectedSubnetMeta *secalib.Metadata
	var expectedSubnetSpec *secalib.SubnetSpecV1

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

		expectedSubnetMeta = &secalib.Metadata{
			Name:       subnetName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   subnetResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.SubnetKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}
		suite.networkTestSuite.verifyNetworkMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		expectedSubnetSpec = &secalib.SubnetSpecV1{
			Cidr: &secalib.SubnetSpecCIDRV1{Ipv4: subnetCidr},
			Zone: zone1,
		}
		suite.networkTestSuite.verifySubNetSpecStep(sCtx, expectedSubnetSpec, subnetResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
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
		suite.networkTestSuite.verifyNetworkMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		suite.networkTestSuite.verifySubNetSpecStep(sCtx, expectedSubnetSpec, subnetResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
	})

	// Security Group
	var groupResp *schema.SecurityGroup
	var expectedGroupMeta *secalib.Metadata
	var expectedGroupSpec *secalib.SecurityGroupSpecV1

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
					{
						Direction: secalib.SecurityRuleDirectionIngress,
					},
				},
			},
		}
		groupResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateSecurityGroup(ctx, group)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, groupResp)

		expectedGroupMeta = &secalib.Metadata{
			Name:       securityGroupName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   securityGroupResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.SecurityGroupKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}
		suite.networkTestSuite.verifyWorkspaceMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		expectedGroupSpec = &secalib.SecurityGroupSpecV1{
			Rules: []*secalib.SecurityGroupRuleV1{{Direction: secalib.SecurityRuleDirectionIngress}},
		}
		suite.networkTestSuite.verifySecurityGroupSpecStep(sCtx, expectedGroupSpec, groupResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*groupResp.Status.State)},
		)
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
		suite.networkTestSuite.verifyWorkspaceMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		suite.networkTestSuite.verifySecurityGroupSpecStep(sCtx, expectedGroupSpec, groupResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*groupResp.Status.State)},
		)
	})

	// Public ip
	var publicIpResp *schema.PublicIp
	var expectedIpMeta *secalib.Metadata
	var expectedIpSpec *secalib.PublicIpSpecV1

	t.WithNewStep("Create public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdatePublicIp", workspaceName)

		ip := &schema.PublicIp{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      publicIPName,
			},
			Spec: schema.PublicIpSpec{
				Address: &publicIpAddress1,
				Version: secalib.IPVersion4,
			},
		}
		publicIpResp, err = suite.regionalClient.NetworkV1.CreateOrUpdatePublicIp(ctx, ip)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta = &secalib.Metadata{
			Name:       publicIPName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   publicIPResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.PublicIPKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}
		suite.networkTestSuite.verifyWorkspaceMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		expectedIpSpec = &secalib.PublicIpSpecV1{
			Version: secalib.IPVersion4,
			Address: publicIpAddress1,
		}
		suite.networkTestSuite.verifyPublicIpSpecStep(sCtx, expectedIpSpec, publicIpResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	t.WithNewStep("Get created public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetPublicIp", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      publicIPName,
		}
		publicIpResp, err = suite.regionalClient.NetworkV1.GetPublicIp(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta.Verb = http.MethodGet
		suite.networkTestSuite.verifyWorkspaceMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		suite.networkTestSuite.verifyPublicIpSpecStep(sCtx, expectedIpSpec, publicIpResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	// Nic
	var nicResp *schema.Nic
	var expectedNicMeta *secalib.Metadata
	var expectedNicSpec *secalib.NICSpecV1

	t.WithNewStep("Create nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNic", workspaceName)

		publicIPURN, err := secapi.BuildReferenceFromURN(publicIPRef)
		if err != nil {
			t.Fatal(err)
		}

		subnetURN, err := secapi.BuildReferenceFromURN(subnetRef)
		if err != nil {
			t.Fatal(err)
		}

		nic := &schema.Nic{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      nicName,
			},
			Spec: schema.NicSpec{
				Addresses:    []string{nicAddress1},
				PublicIpRefs: &[]schema.Reference{*publicIPURN},
				SubnetRef:    *subnetURN,
			},
		}
		nicResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateNic(ctx, nic)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMeta = &secalib.Metadata{
			Name:       nicName,
			Provider:   secalib.NetworkProviderV1,
			Resource:   nicResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.NicKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}
		suite.networkTestSuite.verifyWorkspaceMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		expectedNicSpec = &secalib.NICSpecV1{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: []string{publicIPRef},
			SubnetRef:    subnetRef,
		}
		suite.networkTestSuite.verifyNicSpecStep(sCtx, expectedNicSpec, nicResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
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
		suite.networkTestSuite.verifyWorkspaceMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		suite.networkTestSuite.verifyNicSpecStep(sCtx, expectedNicSpec, nicResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
	})

	// Instance
	var instanceResp *schema.Instance
	var expectedSpec *secalib.InstanceSpecV1

	t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "CreateOrUpdateInstance", workspaceName)

		instanceSkuURN, err := secapi.BuildReferenceFromURN(instanceSkuRef)
		if err != nil {
			t.Fatal(err)
		}

		blockStorageURN, err := secapi.BuildReferenceFromURN(blockStorageRef)
		if err != nil {
			t.Fatal(err)
		}

		inst := &schema.Instance{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      instanceName,
			},
			Spec: schema.InstanceSpec{
				SkuRef: *instanceSkuURN,
				Zone:   zone1,
			},
		}
		inst.Spec.BootVolume.DeviceRef = *blockStorageURN

		instanceResp, err = suite.regionalClient.ComputeV1.CreateOrUpdateInstance(ctx, inst)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMeta = &secalib.Metadata{
			Name:       instanceName,
			Provider:   secalib.ComputeProviderV1,
			Resource:   instanceResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}
		verifyInstanceMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		expectedSpec = &secalib.InstanceSpecV1{
			SkuRef:        instanceSkuRef,
			Zone:          zone1,
			BootDeviceRef: blockStorageRef,
		}
		verifyInstanceSpecStep(sCtx, expectedSpec, &instanceResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
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

		expectedMeta.Verb = http.MethodGet
		verifyInstanceMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		verifyInstanceSpecStep(sCtx, expectedSpec, &instanceResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	// DELETE ALL
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

	// Delete PublicIP
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

	slog.Info("Finishing Usages Lifecycle Test")
}

func (suite *UsagesV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
