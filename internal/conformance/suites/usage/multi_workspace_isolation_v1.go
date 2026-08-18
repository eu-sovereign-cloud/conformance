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

// MultiWorkspaceIsolationV1TestSuite provisions two fully independent regional resource
// stacks (network, subnet, route table, gateway, security group, public IP, nic, block
// storage, instance) in two separate workspaces of the same tenant, all governed by a
// single Role/RoleAssignment scoped to both workspace names. It verifies that the shared
// RBAC grant works across both workspaces, and that listing instances per workspace
// only ever returns that workspace's own resources (isolation).
type MultiWorkspaceIsolationV1TestSuite struct {
	suites.MixedTestSuite

	config *MultiWorkspaceIsolationV1Config
	params *params.MultiWorkspaceIsolationV1Params
}

type MultiWorkspaceIsolationV1Config struct {
	Users          []string
	NetworkCidr    string
	PublicIpsRange string
	RegionZones    []string
	StorageSkus    []string
	InstanceSkus   []string
	NetworkSkus    []string
}

func CreateMultiWorkspaceIsolationV1TestSuite(mixedTestSuite suites.MixedTestSuite, config *MultiWorkspaceIsolationV1Config) *MultiWorkspaceIsolationV1TestSuite {
	suite := &MultiWorkspaceIsolationV1TestSuite{
		MixedTestSuite: mixedTestSuite,
		config:         config,
	}
	suite.ScenarioName = constants.UsageMultiWorkspaceIsolationV1SuiteName.String()
	return suite
}

// buildWorkspaceStack builds one fully independent regional resource stack (network
// topology + storage + a single instance) for the given workspace, using netNum to
// derive a non-overlapping subnet and public IP from the shared scenario ranges.
func (suite *MultiWorkspaceIsolationV1TestSuite) buildWorkspaceStack(t provider.T, workspaceName string, netNum int) params.WorkspaceStackV1 {
	zone := suite.config.RegionZones[rand.Intn(len(suite.config.RegionZones))]
	storageSkuName := suite.config.StorageSkus[rand.Intn(len(suite.config.StorageSkus))]
	instanceSkuName := suite.config.InstanceSkus[rand.Intn(len(suite.config.InstanceSkus))]
	networkSkuName := suite.config.NetworkSkus[rand.Intn(len(suite.config.NetworkSkus))]

	subnetCidr, err := generators.GenerateSubnetCidr(suite.config.NetworkCidr, 8, netNum)
	if err != nil {
		slog.Error("Failed to generate subnet cidr", "error", err)
		t.FailNow()
	}
	nicAddress, err := generators.GenerateNicAddress(subnetCidr, 1)
	if err != nil {
		slog.Error("Failed to generate nic address", "error", err)
		t.FailNow()
	}
	publicIpAddress, err := generators.GeneratePublicIp(suite.config.PublicIpsRange, netNum)
	if err != nil {
		slog.Error("Failed to generate public ip", "error", err)
		t.FailNow()
	}

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

	publicIpName := generators.GeneratePublicIpName()
	publicIpRefObj := generators.GeneratePublicIpRefObject(sdkconsts.NetworkProviderV1Name, suite.Tenant, workspaceName, publicIpName)

	nicName := generators.GenerateNicName()

	instanceName := generators.GenerateInstanceName()

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
			Cidr:          schema.Cidr{Ipv4: subnetCidr},
			RouteTableRef: *routeTableRefObj,
			Zone:          zone,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Subnet: %v", err)
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

	publicIp, err := builders.NewPublicIpBuilder().
		Name(publicIpName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.PublicIpSpec{
			Version: schema.IPVersionIPv4,
			Address: publicIpAddress,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Public IP: %v", err)
	}

	nic, err := builders.NewNicBuilder().
		Name(nicName).
		Provider(sdkconsts.NetworkProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.NicSpec{
			Addresses:    []string{nicAddress},
			PublicIpRefs: []schema.Reference{*publicIpRefObj},
			SubnetRef:    *subnetRefObj,
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

	instanceIterator, err := builders.NewInstanceIteratorBuilder().
		Provider(sdkconsts.ComputeProviderV1Name).
		Tenant(suite.Tenant).Workspace(workspaceName).
		Items([]schema.Instance{*instance}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build InstanceIterator: %v", err)
	}

	return params.WorkspaceStackV1{
		BlockStorage:    blockStorage,
		Network:         network,
		InternetGateway: internetGateway,
		RouteTable:      routeTable,
		Subnet:          subnet,
		SecurityGroup:   securityGroup,
		PublicIp:        publicIp,
		Nic:             nic,
		Instance:        instance,
		Instances:       *instanceIterator,
	}
}

func (suite *MultiWorkspaceIsolationV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.UsageParentSuite)

	roleAssignmentSub := suite.config.Users[rand.Intn(len(suite.config.Users))]

	workspaceAName := generators.GenerateWorkspaceName()
	workspaceBName := generators.GenerateWorkspaceName()

	roleName := generators.GenerateRoleName()
	roleAssignmentName := generators.GenerateRoleAssignmentName()

	// RBAC: a single grant scoped to both workspaces, read-only across every resource
	// kind used by the two independent stacks.
	role, err := builders.NewRoleBuilder().
		Name(roleName).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Spec(&schema.RoleSpec{
			Permissions: []schema.Permission{
				{Provider: sdkconsts.WorkspaceProviderV1Name, Resources: []string{generators.GenerateWorkspaceListResource()}, Verb: []string{http.MethodGet}},
				{Provider: sdkconsts.NetworkProviderV1Name, Resources: []string{generators.GenerateNetworkListResource()}, Verb: []string{http.MethodGet}},
				{Provider: sdkconsts.ComputeProviderV1Name, Resources: []string{generators.GenerateInstanceListResource()}, Verb: []string{http.MethodGet}},
				{Provider: sdkconsts.StorageProviderV1Name, Resources: []string{generators.GenerateBlockStorageListResource()}, Verb: []string{http.MethodGet}},
			},
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Role: %v", err)
	}

	roleAssignment, err := builders.NewRoleAssignmentBuilder().
		Name(roleAssignmentName).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Spec(&schema.RoleAssignmentSpec{
			Roles: []string{roleName},
			Subs:  []string{roleAssignmentSub},
			Scopes: []schema.RoleAssignmentScope{
				{Workspaces: []string{workspaceAName, workspaceBName}},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build RoleAssignment: %v", err)
	}

	workspaceA, err := builders.NewWorkspaceBuilder().
		Name(workspaceAName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace A: %v", err)
	}

	workspaceB, err := builders.NewWorkspaceBuilder().
		Name(workspaceBName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace B: %v", err)
	}

	stackA := suite.buildWorkspaceStack(t, workspaceAName, 1)
	stackB := suite.buildWorkspaceStack(t, workspaceBName, 2)

	p := &params.MultiWorkspaceIsolationV1Params{
		Role:           role,
		RoleAssignment: roleAssignment,
		WorkspaceA:     workspaceA,
		WorkspaceB:     workspaceB,
		StackA:         stackA,
		StackB:         stackB,
	}
	suite.params = p
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockUsage.ConfigureMultiWorkspaceIsolationV1, *p)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

// createStack runs the create+get steps for one workspace's full resource stack and
// returns the observed instance, so it can be used for the isolation List check.
func (suite *MultiWorkspaceIsolationV1TestSuite) createStack(t provider.T, stepsBuilder *steps.StepsConfigurator, label string, workspace *schema.Workspace, stack params.WorkspaceStackV1) *schema.Instance {
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create workspace "+label, t, suite.RegionalClient.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetWorkspaceV1Step("Get workspace "+label, suite.RegionalClient.WorkspaceV1,
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

	block := stack.BlockStorage
	expectBlockMeta := block.Metadata
	expectBlockSpec := &block.Spec
	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create block storage "+label, t, suite.RegionalClient.StorageV1, block,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:       expectBlockMeta,
			Spec:           expectBlockSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetBlockStorageV1Step("Get block storage "+label, suite.RegionalClient.StorageV1,
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

	network := stack.Network
	expectNetworkMeta := network.Metadata
	expectNetworkSpec := &network.Spec
	stepsBuilder.CreateOrUpdateNetworkV1Step("Create network "+label, t, suite.RegionalClient.NetworkV1, network,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			Metadata:       expectNetworkMeta,
			Spec:           expectNetworkSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	gateway := stack.InternetGateway
	expectGatewayMeta := gateway.Metadata
	expectGatewaySpec := &gateway.Spec
	stepsBuilder.CreateOrUpdateInternetGatewayV1Step("Create internet gateway "+label, t, suite.RegionalClient.NetworkV1, gateway,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			Metadata:       expectGatewayMeta,
			Spec:           expectGatewaySpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetInternetGatewayV1Step("Get internet gateway "+label, suite.RegionalClient.NetworkV1,
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

	route := stack.RouteTable
	expectRouteMeta := route.Metadata
	expectRouteSpec := &route.Spec
	stepsBuilder.CreateOrUpdateRouteTableV1Step("Create route table "+label, t, suite.RegionalClient.NetworkV1, route,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			Metadata:       expectRouteMeta,
			Spec:           expectRouteSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetRouteTableV1Step("Get route table "+label, suite.RegionalClient.NetworkV1,
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

	stepsBuilder.GetNetworkV1Step("Get network "+label, suite.RegionalClient.NetworkV1,
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

	subnet := stack.Subnet
	expectSubnetMeta := subnet.Metadata
	expectSubnetSpec := &subnet.Spec
	stepsBuilder.CreateOrUpdateSubnetV1Step("Create subnet "+label, t, suite.RegionalClient.NetworkV1, subnet,
		steps.ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			Metadata:       expectSubnetMeta,
			Spec:           expectSubnetSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetSubnetV1Step("Get subnet "+label, suite.RegionalClient.NetworkV1,
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

	group := stack.SecurityGroup
	expectGroupMeta := group.Metadata
	expectGroupSpec := &group.Spec
	stepsBuilder.CreateOrUpdateSecurityGroupV1Step("Create security group "+label, t, suite.RegionalClient.NetworkV1, group,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			Metadata:       expectGroupMeta,
			Spec:           expectGroupSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetSecurityGroupV1Step("Get security group "+label, suite.RegionalClient.NetworkV1,
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

	publicIp := stack.PublicIp
	expectPublicIpMeta := publicIp.Metadata
	expectPublicIpSpec := &publicIp.Spec
	stepsBuilder.CreateOrUpdatePublicIpV1Step("Create public ip "+label, t, suite.RegionalClient.NetworkV1, publicIp,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			Metadata:       expectPublicIpMeta,
			Spec:           expectPublicIpSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetPublicIpV1Step("Get public ip "+label, suite.RegionalClient.NetworkV1,
		secapi.WorkspaceReference{Tenant: secapi.TenantID(publicIp.Metadata.Tenant), Workspace: secapi.WorkspaceID(publicIp.Metadata.Workspace), Name: publicIp.Metadata.Name},
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus]{
			Metadata: expectPublicIpMeta,
			Spec:     expectPublicIpSpec,
			ResourceStatus: schema.PublicIpStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	nic := stack.Nic
	expectNicMeta := nic.Metadata
	expectNicSpec := &nic.Spec
	stepsBuilder.CreateOrUpdateNicV1Step("Create nic "+label, t, suite.RegionalClient.NetworkV1, nic,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			Metadata:       expectNicMeta,
			Spec:           expectNicSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetNicV1Step("Get nic "+label, suite.RegionalClient.NetworkV1,
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

	instance := stack.Instance
	expectInstanceMeta := instance.Metadata
	expectInstanceSpec := &instance.Spec
	stepsBuilder.CreateOrUpdateInstanceV1Step("Create instance "+label, t, suite.RegionalClient.ComputeV1, instance,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:       expectInstanceMeta,
			Spec:           expectInstanceSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	instance = stepsBuilder.GetInstanceV1Step("Get instance "+label, suite.RegionalClient.ComputeV1,
		secapi.WorkspaceReference{Tenant: secapi.TenantID(instance.Metadata.Tenant), Workspace: secapi.WorkspaceID(instance.Metadata.Workspace), Name: instance.Metadata.Name},
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			Metadata: expectInstanceMeta,
			Spec:     expectInstanceSpec,
			ResourceStatus: schema.InstanceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	return instance
}

// deleteStack runs the deletion steps for one workspace's full resource stack in
// reverse dependency order. stack.Instance must hold the observed instance (the
// value returned by createStack), not the pre-create fixture.
func (suite *MultiWorkspaceIsolationV1TestSuite) deleteStack(t provider.T, stepsBuilder *steps.StepsConfigurator, label string, workspace *schema.Workspace, stack params.WorkspaceStackV1) {
	stepsBuilder.DeleteInstanceV1Step("Delete instance "+label, t, suite.RegionalClient.ComputeV1, stack.Instance)
	stepsBuilder.DeleteNicV1Step("Delete nic "+label, t, suite.RegionalClient.NetworkV1, stack.Nic)
	stepsBuilder.DeletePublicIpV1Step("Delete public ip "+label, t, suite.RegionalClient.NetworkV1, stack.PublicIp)
	stepsBuilder.DeleteSecurityGroupV1Step("Delete security group "+label, t, suite.RegionalClient.NetworkV1, stack.SecurityGroup)
	stepsBuilder.DeleteSubnetV1Step("Delete subnet "+label, t, suite.RegionalClient.NetworkV1, stack.Subnet)
	stepsBuilder.DeleteRouteTableV1Step("Delete route table "+label, t, suite.RegionalClient.NetworkV1, stack.RouteTable)
	stepsBuilder.DeleteInternetGatewayV1Step("Delete internet gateway "+label, t, suite.RegionalClient.NetworkV1, stack.InternetGateway)
	stepsBuilder.DeleteNetworkV1Step("Delete network "+label, t, suite.RegionalClient.NetworkV1, stack.Network)
	stepsBuilder.DeleteBlockStorageV1Step("Delete block storage "+label, t, suite.RegionalClient.StorageV1, stack.BlockStorage)
	stepsBuilder.DeleteWorkspaceV1Step("Delete workspace "+label, t, suite.RegionalClient.WorkspaceV1, workspace)
}

func (suite *MultiWorkspaceIsolationV1TestSuite) TestScenario(t provider.T) {
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
		string(schema.RegionalResourceMetadataKindResourceKindInstance),
		string(schema.RegionalResourceMetadataKindResourceKindNetwork),
		string(schema.RegionalResourceMetadataKindResourceKindInternetGateway),
		string(schema.RegionalResourceMetadataKindResourceKindNic),
		string(schema.RegionalResourceMetadataKindResourceKindPublicIP),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindRoutingTable),
		string(schema.RegionalNetworkResourceMetadataKindResourceKindSubnet),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindSecurityGroup),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Authorization - one grant, scoped to both workspaces

	role := suite.params.Role
	expectRoleMeta := role.Metadata
	expectRoleSpec := &role.Spec
	stepsBuilder.CreateOrUpdateRoleV1Step("Create a role", t, suite.GlobalClient.AuthorizationV1, role,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			Metadata:       expectRoleMeta,
			Spec:           expectRoleSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	role = stepsBuilder.GetRoleV1Step("Get the created role", suite.GlobalClient.AuthorizationV1,
		secapi.TenantReference{Tenant: secapi.TenantID(suite.Tenant), Name: role.Metadata.Name},
		steps.ResponseExpectsWithCondition[schema.GlobalTenantResourceMetadata, schema.RoleSpec, schema.RoleStatus]{
			Metadata: expectRoleMeta,
			Spec:     expectRoleSpec,
			ResourceStatus: schema.RoleStatus{
				State: schema.ResourceStateActive,
			},
		},
	)

	roleAssign := suite.params.RoleAssignment
	expectRoleAssignMeta := roleAssign.Metadata
	expectRoleAssignSpec := &roleAssign.Spec
	stepsBuilder.CreateOrUpdateRoleAssignmentV1Step("Create a role assignment scoped to both workspaces", t, suite.GlobalClient.AuthorizationV1, roleAssign,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			Metadata:       expectRoleAssignMeta,
			Spec:           expectRoleAssignSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	roleAssign = stepsBuilder.GetRoleAssignmentV1Step("Get the created role assignment", suite.GlobalClient.AuthorizationV1,
		secapi.TenantReference{Tenant: secapi.TenantID(suite.Tenant), Name: roleAssign.Metadata.Name},
		steps.ResponseExpectsWithCondition[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec, schema.RoleAssignmentStatus]{
			Metadata: expectRoleAssignMeta,
			Spec:     expectRoleAssignSpec,
			ResourceStatus: schema.RoleAssignmentStatus{
				State: schema.ResourceStateActive,
			},
		},
	)

	// Two fully independent workspace stacks, provisioned under the same shared grant

	workspaceA := suite.params.WorkspaceA
	workspaceB := suite.params.WorkspaceB
	stackA := suite.params.StackA
	stackB := suite.params.StackB

	instanceA := suite.createStack(t, stepsBuilder, "A", workspaceA, stackA)
	instanceB := suite.createStack(t, stepsBuilder, "B", workspaceB, stackB)
	stackA.Instance = instanceA
	stackB.Instance = instanceB

	// Isolation check: listing instances scoped to workspace A must return only
	// workspace A's instance, never workspace B's (and vice versa).
	wpathA := secapi.WorkspacePath{Tenant: secapi.TenantID(workspaceA.Metadata.Tenant), Workspace: secapi.WorkspaceID(workspaceA.Metadata.Name)}
	stepsBuilder.ListInstanceV1Step("List instances in workspace A - expect only workspace A's instance", suite.RegionalClient.ComputeV1, wpathA, nil,
		steps.ListResponseExpects[schema.Instance]{
			Metadata: &stackA.Instances.Metadata,
			Items:    []schema.Instance{*instanceA},
		},
	)

	wpathB := secapi.WorkspacePath{Tenant: secapi.TenantID(workspaceB.Metadata.Tenant), Workspace: secapi.WorkspaceID(workspaceB.Metadata.Name)}
	stepsBuilder.ListInstanceV1Step("List instances in workspace B - expect only workspace B's instance", suite.RegionalClient.ComputeV1, wpathB, nil,
		steps.ListResponseExpects[schema.Instance]{
			Metadata: &stackB.Instances.Metadata,
			Items:    []schema.Instance{*instanceB},
		},
	)

	// Resources deletion - each stack torn down independently, then the shared RBAC grant

	suite.deleteStack(t, stepsBuilder, "A", workspaceA, stackA)
	suite.deleteStack(t, stepsBuilder, "B", workspaceB, stackB)

	stepsBuilder.DeleteRoleAssignmentV1Step("Delete the role assignment", t, suite.GlobalClient.AuthorizationV1, roleAssign)
	stepsBuilder.DeleteRoleV1Step("Delete the role", t, suite.GlobalClient.AuthorizationV1, role)

	suite.FinishScenario()
}

func (suite *MultiWorkspaceIsolationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
