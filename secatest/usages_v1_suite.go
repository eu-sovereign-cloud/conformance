package secatest

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	authorization "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.authorization.v1"
	network "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.network.v1"
	storage "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.storage.v1"
	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type UsagesV1TestSuite struct {
	complexTestSuite

	users          []string
	networkCidr    string
	publicIpsRange string
	regionZones    []string
	storageSkus    []string
	instanceSkus   []string
	networkSkus    []string
}

func (suite *UsagesV1TestSuite) TestUsagesV1(t provider.T) {
	slog.Info("Starting Usages V1 Test")

	t.Title("Usages Lifecycle Test")
	configureTags(t, secalib.NetworkProviderV1, secalib.NetworkKind, secalib.InternetGatewayKind, secalib.NicKind, secalib.PublicIPKind, secalib.RouteTableKind,
		secalib.SubnetKind, secalib.SecurityGroupKind)

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

	networkSkuRef1 := secalib.GenerateSkuRef(networkSkuName1)
	networkSkuRef2 := secalib.GenerateSkuRef(networkSkuName2)

	networkName := secalib.GenerateNetworkName()
	networkResource := secalib.GenerateNetworkResource(suite.tenant, workspaceName, networkName)
	networkRef := secalib.GenerateNetworkRef(networkName)

	internetGatewayName := secalib.GenerateInternetGatewayName()
	internetGatewayResource := secalib.GenerateInternetGatewayResource(suite.tenant, workspaceName, internetGatewayName)
	internetGatewayRef := secalib.GenerateInternetGatewayRef(internetGatewayName)

	routeTableName := secalib.GenerateRouteTableName()
	routeTableResource := secalib.GenerateRouteTableResource(suite.tenant, workspaceName, routeTableName)
	routeTableRef := secalib.GenerateRouteTableRef(routeTableName)

	subnetName := secalib.GenerateSubnetName()
	subnetResource := secalib.GenerateSubnetResource(suite.tenant, workspaceName, subnetName)
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
				MockURL:   suite.mockServerURL,
				AuthToken: suite.authToken,
				Tenant:    suite.tenant,
				Region:    suite.region,
			},
			Authorization: &mock.AuthorizationParamsV1{
				Params: &mock.Params{
					MockURL:   suite.mockServerURL,
					AuthToken: suite.authToken,
					Tenant:    suite.tenant,
				},
				Role: &mock.ResourceParams[secalib.RoleSpecV1]{
					Name: roleName,
					InitialSpec: &secalib.RoleSpecV1{
						Permissions: []*secalib.RoleSpecPermissionV1{
							{
								Provider:  secalib.StorageProviderV1,
								Resources: []string{imageResource},
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
			slog.Error("Failed to create wiremock scenario", "error", err)
			return
		}

		suite.mockClient = wm
	}
	ctx := context.Background()
	var roleResp *authorization.Role
	var assignResp *authorization.RoleAssignment

	// Role
	var expectedRoleMeta *secalib.Metadata
	var expectedRoleSpec *secalib.RoleSpecV1

	t.WithNewStep("Create role", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateRole",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleName,
		}
		role := &authorization.Role{
			Spec: authorization.RoleSpec{
				Permissions: []authorization.Permission{
					{
						Provider:  secalib.StorageProviderV1,
						Resources: []string{imageResource},
						Verb:      []string{http.MethodGet},
					},
				},
			},
		}
		roleResp, err = suite.globalClient.AuthorizationV1.CreateOrUpdateRole(ctx, tref, role, nil)
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
		sCtx.WithNewParameters(
			operationStepParameter, "GetRole",
			tenantStepParameter, suite.tenant,
		)

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
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateRoleAssignment",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleAssignmentName,
		}
		assign := &authorization.RoleAssignment{
			Spec: authorization.RoleAssignmentSpec{
				Roles:  []string{roleName},
				Subs:   []string{roleAssignmentSub1},
				Scopes: []authorization.RoleAssignmentScope{{Tenants: &[]string{suite.tenant}}},
			},
		}
		assignResp, err = suite.globalClient.AuthorizationV1.CreateOrUpdateRoleAssignment(ctx, tref, assign, nil)
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
		sCtx.WithNewParameters(
			operationStepParameter, "GetRoleAssignment",
			tenantStepParameter, suite.tenant,
		)

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
	var resp *workspace.Workspace
	var expectedMeta *secalib.Metadata
	t.WithNewStep("Create workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateWorkspace",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspaceName,
		}
		ws := &workspace.Workspace{}
		resp, err = suite.regionalClient.WorkspaceV1.CreateOrUpdateWorkspace(ctx, tref, ws, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta = &secalib.Metadata{
			Name:       workspaceName,
			Provider:   secalib.WorkspaceProviderV1,
			Resource:   workspaceResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.WorkspaceKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyWorkspaceMetadataStep(sCtx, expectedMeta, resp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*resp.Status.State)},
		)
	})

	t.WithNewStep("Get created workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetWorkspace",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspaceName,
		}
		resp, err = suite.regionalClient.WorkspaceV1.GetWorkspace(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta.Verb = http.MethodGet
		verifyWorkspaceMetadataStep(sCtx, expectedMeta, resp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*resp.Status.State)},
		)
	})

	//Storage

	// Image
	var imageResp *storage.Image
	var expectedImageMeta *secalib.Metadata
	var expectedImageSpec *secalib.ImageSpecV1

	t.WithNewStep("Create image", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateImage",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   imageName,
		}
		img := &storage.Image{
			Spec: storage.ImageSpec{
				BlockStorageRef: blockStorageRef,
				CpuArchitecture: secalib.CpuArchitectureAmd64,
			},
		}
		imageResp, err = suite.regionalClient.StorageV1.CreateOrUpdateImage(ctx, tref, img, nil)
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
			Region:     suite.region,
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
		sCtx.WithNewParameters(
			operationStepParameter, "GetImage",
			tenantStepParameter, suite.tenant,
		)

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
	var blockResp *storage.BlockStorage
	var expectedBlockMeta *secalib.Metadata
	var expectedBlockSpec *secalib.BlockStorageSpecV1

	t.WithNewStep("Create block storage", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateBlockStorage",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		bo := &storage.BlockStorage{
			Spec: storage.BlockStorageSpec{
				SizeGB: blockStorageSize,
				SkuRef: storageSkuRef,
			},
		}
		blockResp, err = suite.regionalClient.StorageV1.CreateOrUpdateBlockStorage(ctx, wref, bo, nil)
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
			Region:     suite.region,
		}

		verifyStorageZonalMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

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
		sCtx.WithNewParameters(
			operationStepParameter, "GetBlockStorage",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		blockResp, err = suite.regionalClient.StorageV1.GetBlockStorage(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMeta.Verb = http.MethodGet
		verifyStorageZonalMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		verifyBlockStorageSpecStep(sCtx, expectedBlockSpec, blockResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
	})
	// Network
	var networkResp *network.Network
	var expectedNetworkMeta *secalib.Metadata
	var expectedNetworkSpec *secalib.NetworkSpecV1

	t.WithNewStep("Create network", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			// TODO Create a method to define the step parameters
			// TODO Add the provider prefix in each operation
			operationStepParameter, "CreateOrUpdateNetwork",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      networkName,
		}
		net := &network.Network{
			Spec: network.NetworkSpec{
				Cidr:          network.Cidr{Ipv4: &suite.networkCidr},
				SkuRef:        networkSkuRef1,
				RouteTableRef: routeTableRef,
			},
		}
		networkResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateNetwork(ctx, wref, net, nil)
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
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		expectedNetworkSpec = &secalib.NetworkSpecV1{
			Cidr:          &secalib.NetworkSpecCIDRV1{Ipv4: suite.networkCidr},
			SkuRef:        networkSkuRef1,
			RouteTableRef: routeTableRef,
		}
		verifyNetworkSpecStep(sCtx, expectedNetworkSpec, networkResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
	})

	t.WithNewStep("Get created network", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetNetwork",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      networkName,
		}
		networkResp, err = suite.regionalClient.NetworkV1.GetNetwork(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		verifyNetworkSpecStep(sCtx, expectedNetworkSpec, networkResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
	})

	// Internet gateway
	var gatewayResp *network.InternetGateway
	var expectedGatewayMeta *secalib.Metadata
	var expectedGatewaySpec *secalib.InternetGatewaySpecV1

	t.WithNewStep("Create internet gateway", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateInternetGateway",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      internetGatewayName,
		}
		gtw := &network.InternetGateway{}
		gatewayResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateInternetGateway(ctx, wref, gtw, nil)
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
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		expectedGatewaySpec = &secalib.InternetGatewaySpecV1{
			EgressOnly: false,
		}
		verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, gatewayResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
	})

	t.WithNewStep("Get created internet gateway", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetInternetGateway",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      internetGatewayName,
		}
		gatewayResp, err = suite.regionalClient.NetworkV1.GetInternetGateway(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		verifyInternetGatewaySpecStep(sCtx, expectedGatewaySpec, gatewayResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
	})

	// Route table
	var routeResp *network.RouteTable
	var expectedRouteMeta *secalib.Metadata
	var expectedRouteSpec *secalib.RouteTableSpecV1

	t.WithNewStep("Create route table", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateRouteTable",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      routeTableName,
		}
		route := &network.RouteTable{}
		routeResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateRouteTable(ctx, wref, route, nil)
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
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		expectedRouteSpec = &secalib.RouteTableSpecV1{
			LocalRef: networkRef,
			Routes: []*secalib.RouteTableRouteV1{
				{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: internetGatewayRef},
			},
		}
		verifyRouteTableSpecStep(sCtx, expectedRouteSpec, routeResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
	})

	t.WithNewStep("Get created route table", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetRouteTable",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      routeTableName,
		}
		routeResp, err = suite.regionalClient.NetworkV1.GetRouteTable(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		verifyRouteTableSpecStep(sCtx, expectedRouteSpec, routeResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
	})

	// Subnet
	var subnetResp *network.Subnet
	var expectedSubnetMeta *secalib.Metadata
	var expectedSubnetSpec *secalib.SubnetSpecV1

	t.WithNewStep("Create subnet", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateSubnet",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      subnetName,
		}
		sub := &network.Subnet{
			Spec: network.SubnetSpec{
				Cidr: network.Cidr{Ipv4: &subnetCidr},
				Zone: zone1,
			},
		}
		subnetResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateSubnet(ctx, wref, sub, nil)
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
			Region:     suite.region,
		}
		verifyNetworkZonalMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		expectedSubnetSpec = &secalib.SubnetSpecV1{
			Cidr: &secalib.SubnetSpecCIDRV1{Ipv4: subnetCidr},
			Zone: zone1,
		}
		verifySubNetSpecStep(sCtx, expectedSubnetSpec, subnetResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
	})

	t.WithNewStep("Get created subnet", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetNetwork",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      subnetName,
		}
		subnetResp, err = suite.regionalClient.NetworkV1.GetSubnet(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMeta.Verb = http.MethodGet
		verifyNetworkZonalMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		verifySubNetSpecStep(sCtx, expectedSubnetSpec, subnetResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
	})

	// Security Group
	var groupResp *network.SecurityGroup
	var expectedGroupMeta *secalib.Metadata
	var expectedGroupSpec *secalib.SecurityGroupSpecV1

	t.WithNewStep("Create security group", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateSecurityGroup",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      securityGroupName,
		}
		group := &network.SecurityGroup{
			Spec: network.SecurityGroupSpec{
				Rules: []network.SecurityGroupRuleSpec{
					{
						Direction: secalib.SecurityRuleDirectionIngress,
					},
				},
			},
		}
		groupResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateSecurityGroup(ctx, wref, group, nil)
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
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		expectedGroupSpec = &secalib.SecurityGroupSpecV1{
			Rules: []*secalib.SecurityGroupRuleV1{{Direction: secalib.SecurityRuleDirectionIngress}},
		}
		verifySecurityGroupSpecStep(sCtx, expectedGroupSpec, groupResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*groupResp.Status.State)},
		)
	})

	t.WithNewStep("Get created security group", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetSecurityGroup",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      securityGroupName,
		}
		groupResp, err = suite.regionalClient.NetworkV1.GetSecurityGroup(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, groupResp)

		expectedGroupMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		verifySecurityGroupSpecStep(sCtx, expectedGroupSpec, groupResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*groupResp.Status.State)},
		)
	})

	// Public ip
	var publicIpResp *network.PublicIp
	var expectedIpMeta *secalib.Metadata
	var expectedIpSpec *secalib.PublicIpSpecV1

	t.WithNewStep("Create public ip", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdatePublicIp",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      publicIPName,
		}
		ip := &network.PublicIp{
			Spec: network.PublicIpSpec{
				Address: &publicIpAddress1,
				Version: secalib.IPVersion4,
			},
		}
		publicIpResp, err = suite.regionalClient.NetworkV1.CreateOrUpdatePublicIp(ctx, wref, ip, nil)
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
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		expectedIpSpec = &secalib.PublicIpSpecV1{
			Version: secalib.IPVersion4,
			Address: publicIpAddress1,
		}
		verifyPublicIpSpecStep(sCtx, expectedIpSpec, publicIpResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	t.WithNewStep("Get created public ip", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetPublicIP",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      publicIPName,
		}
		publicIpResp, err = suite.regionalClient.NetworkV1.GetPublicIp(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		verifyPublicIpSpecStep(sCtx, expectedIpSpec, publicIpResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	// Nic
	var nicResp *network.Nic
	var expectedNicMeta *secalib.Metadata
	var expectedNicSpec *secalib.NICSpecV1

	t.WithNewStep("Create nic", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateNic",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      nicName,
		}
		nic := &network.Nic{
			Spec: network.NicSpec{
				Addresses:    []string{nicAddress1},
				PublicIpRefs: &[]interface{}{publicIPRef},
				SubnetRef:    subnetRef,
			},
		}
		nicResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateNic(ctx, wref, nic, nil)
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
			Region:     suite.region,
		}
		verifyNetworkZonalMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		expectedNicSpec = &secalib.NICSpecV1{
			Addresses:    []string{nicAddress1},
			PublicIpRefs: []string{publicIPRef},
			SubnetRef:    subnetRef,
		}
		verifyNicSpecStep(sCtx, expectedNicSpec, nicResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
	})

	t.WithNewStep("Get created nic", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetNic",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      nicName,
		}
		nicResp, err = suite.regionalClient.NetworkV1.GetNic(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMeta.Verb = http.MethodGet
		verifyNetworkZonalMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		verifyNicSpecStep(sCtx, expectedNicSpec, nicResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
	})
	slog.Info("Finishing Usages Lifecycle Test")
}

func (suite *UsagesV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
