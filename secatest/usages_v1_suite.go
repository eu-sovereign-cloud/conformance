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
	compute "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.compute.v1"
	network "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.network.v1"
	storage "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.storage.v1"
	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/google/uuid"

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

func (suite *UsagesV1TestSuite) generateUsagesParams() (*secalib.UsageParamsV1, error) {
	workspaceName := secalib.GenerateWorkspaceName()
	roleName := secalib.GenerateRoleName()
	roleAssignmentName := secalib.GenerateRoleAssignmentName()
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]
	blockStorageName := secalib.GenerateBlockStorageName()
	imageName := secalib.GenerateImageName()
	instanceSkuName := suite.instanceSkus[rand.Intn(len(suite.instanceSkus))]
	instanceName := secalib.GenerateInstanceName()
	networkSkuName := suite.networkSkus[rand.Intn(len(suite.networkSkus))]
	networkName := secalib.GenerateNetworkName()
	internetGatewayName := secalib.GenerateInternetGatewayName()
	routeTableName := secalib.GenerateRouteTableName()
	subnetName := secalib.GenerateSubnetName()
	publicIPName := secalib.GeneratePublicIPName()
	securityGroupName := secalib.GenerateSecurityGroupName()

	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)
	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName)
	instanceSkuRef := secalib.GenerateSkuRef(instanceSkuName)
	networkSkuRef := secalib.GenerateSkuRef(networkSkuName)
	networkRef := secalib.GenerateNetworkRef(networkName)
	internetGatewayRef := secalib.GenerateInternetGatewayRef(internetGatewayName)
	routeTableRef := secalib.GenerateRouteTableRef(routeTableName)
	subnetRef := secalib.GenerateSubnetRef(subnetName)
	nicName := secalib.GenerateNicName()
	publicIPRef := secalib.GeneratePublicIPRef(publicIPName)

	imageResource := secalib.GenerateImageResource(suite.tenant, imageName)

	// Random data
	blockStorageSize := secalib.GenerateBlockStorageSize()
	zone := suite.regionZones[rand.Intn(len(suite.regionZones))]
	roleAssignmentSub := suite.users[rand.Intn(len(suite.users))]

	subnetCidr, err := secalib.GenerateSubnetCidr(suite.networkCidr, 8, 1)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate subnet cidr: %w", err)
	}

	nicAddress, err := secalib.GenerateNicAddress(subnetCidr, 1)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate nic address: %w", err)
	}

	publicIpAddress, err := secalib.GeneratePublicIp(suite.publicIpsRange, 1)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate public ip: %w", err)
	}

	return &secalib.UsageParamsV1{
		Workspace: &secalib.ResourceParams[secalib.WorkspaceSpecV1]{
			Name: workspaceName,
			InitialSpec: &secalib.WorkspaceSpecV1{
				Labels: &[]secalib.Label{},
			},
		},
		Role: &secalib.ResourceParams[secalib.RoleSpecV1]{
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
		RoleAssignment: &secalib.ResourceParams[secalib.RoleAssignmentSpecV1]{
			Name: roleAssignmentName,
			InitialSpec: &secalib.RoleAssignmentSpecV1{
				Roles:  []string{roleName},
				Subs:   []string{roleAssignmentSub},
				Scopes: []*secalib.RoleAssignmentSpecScopeV1{{Tenants: []string{suite.tenant}}},
			},
		},
		BlockStorage: &secalib.ResourceParams[secalib.BlockStorageSpecV1]{
			Name: blockStorageName,
			InitialSpec: &secalib.BlockStorageSpecV1{
				SkuRef: storageSkuRef,
				SizeGB: blockStorageSize,
			},
		},
		Image: &secalib.ResourceParams[secalib.ImageSpecV1]{
			Name: imageName,
			InitialSpec: &secalib.ImageSpecV1{
				BlockStorageRef: blockStorageRef,
				CpuArchitecture: secalib.CpuArchitectureAmd64,
			},
		},
		Network: &secalib.ResourceParams[secalib.NetworkSpecV1]{
			Name: networkName,
			InitialSpec: &secalib.NetworkSpecV1{
				Cidr:          &secalib.NetworkSpecCIDRV1{Ipv4: suite.networkCidr},
				SkuRef:        networkSkuRef,
				RouteTableRef: routeTableRef,
			},
		},
		InternetGateway: &secalib.ResourceParams[secalib.InternetGatewaySpecV1]{
			Name:        internetGatewayName,
			InitialSpec: &secalib.InternetGatewaySpecV1{EgressOnly: false},
		},
		RouteTable: &secalib.ResourceParams[secalib.RouteTableSpecV1]{
			Name: routeTableName,
			InitialSpec: &secalib.RouteTableSpecV1{
				LocalRef: networkRef,
				Routes: []*secalib.RouteTableRouteV1{
					{DestinationCidrBlock: routeTableDefaultDestination, TargetRef: internetGatewayRef},
				},
			},
		},
		Subnet: &secalib.ResourceParams[secalib.SubnetSpecV1]{
			Name: subnetName,
			InitialSpec: &secalib.SubnetSpecV1{
				Cidr: &secalib.SubnetSpecCIDRV1{Ipv4: subnetCidr},
				Zone: zone,
			},
		},
		Nic: &secalib.ResourceParams[secalib.NICSpecV1]{
			Name: nicName,
			InitialSpec: &secalib.NICSpecV1{
				Addresses:    []string{nicAddress},
				PublicIpRefs: []string{publicIPRef},
				SubnetRef:    subnetRef,
			},
		},
		PublicIp: &secalib.ResourceParams[secalib.PublicIpSpecV1]{
			Name: publicIPName,
			InitialSpec: &secalib.PublicIpSpecV1{
				Version: secalib.IPVersion4,
				Address: publicIpAddress,
			},
		},
		SecurityGroup: &secalib.ResourceParams[secalib.SecurityGroupSpecV1]{
			Name: securityGroupName,
			InitialSpec: &secalib.SecurityGroupSpecV1{
				Rules: []*secalib.SecurityGroupRuleV1{{Direction: secalib.SecurityRuleDirectionIngress}},
			},
		},
		Instance: &secalib.ResourceParams[secalib.InstanceSpecV1]{
			Name: instanceName,
			InitialSpec: &secalib.InstanceSpecV1{
				SkuRef:        instanceSkuRef,
				Zone:          zone,
				BootDeviceRef: blockStorageRef,
			},
		},
	}, nil
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

	params, err := suite.generateUsagesParams()
	if err != nil {
		t.Fatalf("Failed to generate usages params: %v", err)
	}

	// Resource URIs
	workspaceResource := secalib.GenerateWorkspaceResource(suite.tenant, params.Workspace.Name)
	roleResource := secalib.GenerateRoleResource(suite.tenant, params.Role.Name)
	roleAssignmentResource := secalib.GenerateRoleAssignmentResource(suite.tenant, params.RoleAssignment.Name)
	imageResource := secalib.GenerateImageResource(suite.tenant, params.Image.Name)
	blockStorageResource := secalib.GenerateBlockStorageResource(suite.tenant, params.Workspace.Name, params.BlockStorage.Name)
	instanceResource := secalib.GenerateInstanceResource(suite.tenant, params.Workspace.Name, params.Instance.Name)
	networkResource := secalib.GenerateNetworkResource(suite.tenant, params.Workspace.Name, params.Network.Name)
	internetGatewayResource := secalib.GenerateInternetGatewayResource(suite.tenant, params.Workspace.Name, params.InternetGateway.Name)
	routeTableResource := secalib.GenerateRouteTableResource(suite.tenant, params.Workspace.Name, params.RouteTable.Name)
	subnetResource := secalib.GenerateSubnetResource(suite.tenant, params.Workspace.Name, params.Subnet.Name)
	nicResource := secalib.GenerateNicResource(suite.tenant, params.Workspace.Name, params.Nic.Name)
	publicIPResource := secalib.GeneratePublicIPResource(suite.tenant, params.Workspace.Name, params.PublicIp.Name)
	securityGroupResource := secalib.GenerateSecurityGroupResource(suite.tenant, params.Workspace.Name, params.SecurityGroup.Name)

	// Setup mock, if configured to use
	if suite.isMockEnabled() {
		scenarios := mock.NewUsageScenariosV1(suite.authToken, suite.tenant, suite.region, suite.mockServerURL)

		id := uuid.New().String()

		wm, err := scenarios.ConfigureUsageScenario(id, params)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}

		suite.mockClient = wm
	}
	ctx := context.Background()	

	// Role
	var roleResp *authorization.Role
	var expectedRoleMeta *secalib.Metadata

	t.WithNewStep("Create role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateRole")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Role.Name,
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
			Name:       params.Role.Name,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleKind,
			Tenant:     suite.tenant,
		}
		verifyAuthorizationMetadataStep(sCtx, expectedRoleMeta, roleResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*roleResp.Status.State)},
		)
	})

	t.WithNewStep("Get created role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRole")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Role.Name,
		}
		roleResp, err = suite.globalClient.AuthorizationV1.GetRole(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, roleResp)

		expectedRoleMeta.Verb = http.MethodGet
		verifyAuthorizationMetadataStep(sCtx, expectedRoleMeta, roleResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*roleResp.Status.State)},
		)
	})
	// Role assignment
	var assignResp *authorization.RoleAssignment
	var expectedAssignMeta *secalib.Metadata

	t.WithNewStep("Create role assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateRoleAssignment")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.RoleAssignment.Name,
		}
		assign := &authorization.RoleAssignment{
			Spec: authorization.RoleAssignmentSpec{
				Roles:  params.RoleAssignment.InitialSpec.Roles,
				Subs:   params.RoleAssignment.InitialSpec.Subs,
				Scopes: []authorization.RoleAssignmentScope{{Tenants: &[]string{suite.tenant}}},
			},
		}
		assignResp, err = suite.globalClient.AuthorizationV1.CreateOrUpdateRoleAssignment(ctx, tref, assign, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, assignResp)

		expectedAssignMeta = &secalib.Metadata{
			Name:       params.RoleAssignment.Name,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleAssignmentResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleAssignmentKind,
			Tenant:     suite.tenant,
		}
		verifyAuthorizationMetadataStep(sCtx, expectedAssignMeta, assignResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*assignResp.Status.State)},
		)
	})

	t.WithNewStep("Get created role assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRoleAssignment")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.RoleAssignment.Name,
		}
		assignResp, err = suite.globalClient.AuthorizationV1.GetRoleAssignment(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, assignResp)

		expectedAssignMeta.Verb = http.MethodGet
		verifyAuthorizationMetadataStep(sCtx, expectedAssignMeta, assignResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*assignResp.Status.State)},
		)
	})

	// Workspace
	var workspaceResp *workspace.Workspace
	var expectedMeta *secalib.Metadata

	t.WithNewStep("Create workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "CreateOrUpdateWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Workspace.Name,
		}
		ws := &workspace.Workspace{}
		workspaceResp, err = suite.regionalClient.WorkspaceV1.CreateOrUpdateWorkspace(ctx, tref, ws, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, workspaceResp)

		expectedMeta = &secalib.Metadata{
			Name:       params.Workspace.Name,
			Provider:   secalib.WorkspaceProviderV1,
			Resource:   workspaceResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.WorkspaceKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
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
			Name:   params.Workspace.Name,
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

	// Image
	var imageResp *storage.Image
	var expectedImageMeta *secalib.Metadata

	t.WithNewStep("Create image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateImage", "")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Image.Name,
		}
		img := &storage.Image{
			Spec: storage.ImageSpec{
				BlockStorageRef: params.Image.InitialSpec.BlockStorageRef,
				CpuArchitecture: secalib.CpuArchitectureAmd64,
			},
		}
		imageResp, err = suite.regionalClient.StorageV1.CreateOrUpdateImage(ctx, tref, img, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMeta = &secalib.Metadata{
			Name:       params.Image.Name,
			Provider:   secalib.StorageProviderV1,
			Resource:   imageResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.ImageKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*imageResp.Status.State)},
		)
	})

	t.WithNewStep("Get created image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetImage", "")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Image.Name,
		}
		imageResp, err = suite.regionalClient.StorageV1.GetImage(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMeta.Verb = http.MethodGet
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*imageResp.Status.State)},
		)
	})

	// Block storage
	var blockResp *storage.BlockStorage
	var expectedBlockMeta *secalib.Metadata

	t.WithNewStep("Create block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateBlockStorage", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.BlockStorage.Name,
		}
		bo := &storage.BlockStorage{
			Spec: storage.BlockStorageSpec{
				SizeGB: params.BlockStorage.InitialSpec.SizeGB,
				SkuRef: params.BlockStorage.InitialSpec.SkuRef,
			},
		}
		blockResp, err = suite.regionalClient.StorageV1.CreateOrUpdateBlockStorage(ctx, wref, bo, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMeta = &secalib.Metadata{
			Name:       params.BlockStorage.Name,
			Provider:   secalib.StorageProviderV1,
			Resource:   blockStorageResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.BlockStorageKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyStorageZonalMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
	})

	t.WithNewStep("Get created block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetBlockStorage", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.BlockStorage.Name,
		}
		blockResp, err = suite.regionalClient.StorageV1.GetBlockStorage(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMeta.Verb = http.MethodGet
		verifyStorageZonalMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
	})
	// Network
	var networkResp *network.Network
	var expectedNetworkMeta *secalib.Metadata

	t.WithNewStep("Create network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNetwork", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Network.Name,
		}
		net := &network.Network{
			Spec: network.NetworkSpec{
				Cidr:          network.Cidr{Ipv4: &suite.networkCidr},
				SkuRef:        params.Network.InitialSpec.SkuRef,
				RouteTableRef: params.Network.InitialSpec.RouteTableRef,
			},
		}
		networkResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateNetwork(ctx, wref, net, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMeta = &secalib.Metadata{
			Name:       params.Network.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   networkResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.NetworkKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
	})

	t.WithNewStep("Get created network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNetwork", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Network.Name,
		}
		networkResp, err = suite.regionalClient.NetworkV1.GetNetwork(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, networkResp)

		expectedNetworkMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedNetworkMeta, networkResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*networkResp.Status.State)},
		)
	})

	// Internet gateway
	var gatewayResp *network.InternetGateway
	var expectedGatewayMeta *secalib.Metadata

	t.WithNewStep("Create internet gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateInternetGateway", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.InternetGateway.Name,
		}
		gtw := &network.InternetGateway{}
		gatewayResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateInternetGateway(ctx, wref, gtw, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMeta = &secalib.Metadata{
			Name:       params.InternetGateway.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   internetGatewayResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InternetGatewayKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
	})

	t.WithNewStep("Get created internet gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetInternetGateway", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.InternetGateway.Name,
		}
		gatewayResp, err = suite.regionalClient.NetworkV1.GetInternetGateway(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, gatewayResp)

		expectedGatewayMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedGatewayMeta, gatewayResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*gatewayResp.Status.State)},
		)
	})

	// Route table
	var routeResp *network.RouteTable
	var expectedRouteMeta *secalib.Metadata

	t.WithNewStep("Create route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateRouteTable", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.RouteTable.Name,
		}
		route := &network.RouteTable{}
		routeResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateRouteTable(ctx, wref, route, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMeta = &secalib.Metadata{
			Name:       params.RouteTable.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   routeTableResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RouteTableKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
	})

	t.WithNewStep("Get created route table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetRouteTable", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.RouteTable.Name,
		}
		routeResp, err = suite.regionalClient.NetworkV1.GetRouteTable(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, routeResp)

		expectedRouteMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedRouteMeta, routeResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*routeResp.Status.State)},
		)
	})

	// Subnet
	var subnetResp *network.Subnet
	var expectedSubnetMeta *secalib.Metadata

	t.WithNewStep("Create subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateSubnet", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Subnet.Name,
		}
		sub := &network.Subnet{
			Spec: network.SubnetSpec{
				Cidr: network.Cidr{Ipv4: &params.Subnet.InitialSpec.Cidr.Ipv4},
				Zone: params.Subnet.InitialSpec.Zone,
			},
		}
		subnetResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateSubnet(ctx, wref, sub, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMeta = &secalib.Metadata{
			Name:       params.Subnet.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   subnetResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.SubnetKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkZonalMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
	})

	t.WithNewStep("Get created subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNetwork", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Subnet.Name,
		}
		subnetResp, err = suite.regionalClient.NetworkV1.GetSubnet(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, subnetResp)

		expectedSubnetMeta.Verb = http.MethodGet
		verifyNetworkZonalMetadataStep(sCtx, expectedSubnetMeta, subnetResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*subnetResp.Status.State)},
		)
	})

	// Security Group
	var groupResp *network.SecurityGroup
	var expectedGroupMeta *secalib.Metadata

	t.WithNewStep("Create security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateSecurityGroup", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.SecurityGroup.Name,
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
			Name:       params.SecurityGroup.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   securityGroupResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.SecurityGroupKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*groupResp.Status.State)},
		)
	})

	t.WithNewStep("Get created security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSecurityGroup", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.SecurityGroup.Name,
		}
		groupResp, err = suite.regionalClient.NetworkV1.GetSecurityGroup(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, groupResp)

		expectedGroupMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedGroupMeta, groupResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*groupResp.Status.State)},
		)
	})

	// Public ip
	var publicIpResp *network.PublicIp
	var expectedIpMeta *secalib.Metadata

	t.WithNewStep("Create public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdatePublicIp", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.PublicIp.Name,
		}
		ip := &network.PublicIp{
			Spec: network.PublicIpSpec{
				Address: &params.PublicIp.InitialSpec.Address,
				Version: network.IPVersion(params.PublicIp.InitialSpec.Version),
			},
		}
		publicIpResp, err = suite.regionalClient.NetworkV1.CreateOrUpdatePublicIp(ctx, wref, ip, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta = &secalib.Metadata{
			Name:       params.PublicIp.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   publicIPResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.PublicIPKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkRegionalMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	t.WithNewStep("Get created public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetPublicIp", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.PublicIp.Name,
		}
		publicIpResp, err = suite.regionalClient.NetworkV1.GetPublicIp(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, publicIpResp)

		expectedIpMeta.Verb = http.MethodGet
		verifyNetworkRegionalMetadataStep(sCtx, expectedIpMeta, publicIpResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*publicIpResp.Status.State)},
		)
	})

	// Nic
	var nicResp *network.Nic
	var expectedNicMeta *secalib.Metadata

	t.WithNewStep("Create nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNic", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Nic.Name,
		}
		nic := &network.Nic{
			Spec: network.NicSpec{
				Addresses:    params.Nic.InitialSpec.Addresses,
				PublicIpRefs: &[]interface{}{params.Nic.InitialSpec.PublicIpRefs},
				SubnetRef:    params.Nic.InitialSpec.SubnetRef,
			},
		}
		nicResp, err = suite.regionalClient.NetworkV1.CreateOrUpdateNic(ctx, wref, nic, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMeta = &secalib.Metadata{
			Name:       params.Nic.Name,
			Provider:   secalib.NetworkProviderV1,
			Resource:   nicResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.NicKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyNetworkZonalMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
	})

	t.WithNewStep("Get created nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNic", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Nic.Name,
		}
		nicResp, err = suite.regionalClient.NetworkV1.GetNic(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, nicResp)

		expectedNicMeta.Verb = http.MethodGet
		verifyNetworkZonalMetadataStep(sCtx, expectedNicMeta, nicResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*nicResp.Status.State)},
		)
	})

	// Instance
	var instanceResp *compute.Instance

	t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "CreateOrUpdateInstance", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Instance.Name,
		}
		inst := &compute.Instance{
			Spec: compute.InstanceSpec{
				SkuRef: params.Instance.InitialSpec.SkuRef,
				Zone:   params.Instance.InitialSpec.Zone,
			},
		}
		inst.Spec.BootVolume.DeviceRef = params.Instance.InitialSpec.BootDeviceRef

		instanceResp, err = suite.regionalClient.ComputeV1.CreateOrUpdateInstance(ctx, wref, inst, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMeta = &secalib.Metadata{
			Name:       params.Instance.Name,
			Provider:   secalib.ComputeProviderV1,
			Resource:   instanceResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyInstanceZonalMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	t.WithNewStep("Get created instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Instance.Name,
		}
		instanceResp, err = suite.regionalClient.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMeta.Verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	// DELETE ALL
	// Delete instance
	t.WithNewStep("Delete instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "DeleteInstance", params.Workspace.Name)

		err = suite.regionalClient.ComputeV1.DeleteInstance(ctx, instanceResp, nil)
		requireNoError(sCtx, err)
	})

	// Delete security group
	t.WithNewStep("Delete security group", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteSecurityGroup", params.Workspace.Name)

		err = suite.regionalClient.NetworkV1.DeleteSecurityGroup(ctx, groupResp, nil)
		requireNoError(sCtx, err)
	})

	// Delete Nic
	t.WithNewStep("Delete Nic", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteNic", params.Workspace.Name)

		err = suite.regionalClient.NetworkV1.DeleteNic(ctx, nicResp, nil)
		requireNoError(sCtx, err)
	})

	// Delete PublicIP
	t.WithNewStep("Delete Public ip", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeletePublicIp", params.Workspace.Name)

		err = suite.regionalClient.NetworkV1.DeletePublicIp(ctx, publicIpResp, nil)
		requireNoError(sCtx, err)
	})

	// Delete subnet
	t.WithNewStep("Delete Subnet", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteSubnet", params.Workspace.Name)

		err = suite.regionalClient.NetworkV1.DeleteSubnet(ctx, subnetResp, nil)
		requireNoError(sCtx, err)
	})

	// Delete Route-table
	t.WithNewStep("Delete Route-table", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteRouteTable", params.Workspace.Name)

		err = suite.regionalClient.NetworkV1.DeleteRouteTable(ctx, routeResp, nil)
		requireNoError(sCtx, err)
	})

	// Delete Internet-gateway
	t.WithNewStep("Delete Internet-gateway", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteInternetGateway", params.Workspace.Name)

		err = suite.regionalClient.NetworkV1.DeleteInternetGateway(ctx, gatewayResp, nil)
		requireNoError(sCtx, err)
	})

	// Delete Network
	t.WithNewStep("Delete Network", func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteNetwork", params.Workspace.Name)

		err = suite.regionalClient.NetworkV1.DeleteNetwork(ctx, networkResp, nil)
		requireNoError(sCtx, err)
	})

	// Delete BlockStorage
	t.WithNewStep("Delete BlockStorage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "DeleteBlockStorage", params.Workspace.Name)

		err = suite.regionalClient.StorageV1.DeleteBlockStorage(ctx, blockResp, nil)
		requireNoError(sCtx, err)
	})

	// Delete Image
	t.WithNewStep("Delete Image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "DeleteImage", "")

		err = suite.regionalClient.StorageV1.DeleteImage(ctx, imageResp, nil)
		requireNoError(sCtx, err)
	})

	// Delete Workspace
	t.WithNewStep("Delete Workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "DeleteWorkspace")

		err = suite.regionalClient.WorkspaceV1.DeleteWorkspace(ctx, workspaceResp, nil)
		requireNoError(sCtx, err)
	})

	// Delete Role Assignment
	t.WithNewStep("Delete Role Assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteRoleAssignment")

		err = suite.globalClient.AuthorizationV1.DeleteRoleAssignment(ctx, assignResp, nil)
		requireNoError(sCtx, err)
	})

	// Delete Role
	t.WithNewStep("Delete Role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteRole")

		err = suite.globalClient.AuthorizationV1.DeleteRole(ctx, roleResp, nil)
		requireNoError(sCtx, err)
	})

	slog.Info("Finishing Usages Lifecycle Test")
}

func (suite *UsagesV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
