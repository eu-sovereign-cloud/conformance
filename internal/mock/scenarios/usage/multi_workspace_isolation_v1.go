package mockusage

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

// workspaceStackURLs holds every URL used across the create, list and delete phases
// of a stack, so they only need to be derived once per stack.
type workspaceStackURLs struct {
	workspace, block, network, gateway, route, subnet, group, publicIp, nic, instance, instanceList string
}

func newWorkspaceStackURLs(workspace *schema.Workspace, stack params.WorkspaceStackV1) workspaceStackURLs {
	return workspaceStackURLs{
		workspace:    generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name),
		block:        generators.GenerateBlockStorageURL(sdkconsts.StorageProviderV1Name, stack.BlockStorage.Metadata.Tenant, stack.BlockStorage.Metadata.Workspace, stack.BlockStorage.Metadata.Name),
		network:      generators.GenerateNetworkURL(sdkconsts.NetworkProviderV1Name, stack.Network.Metadata.Tenant, stack.Network.Metadata.Workspace, stack.Network.Metadata.Name),
		gateway:      generators.GenerateInternetGatewayURL(sdkconsts.NetworkProviderV1Name, stack.InternetGateway.Metadata.Tenant, stack.InternetGateway.Metadata.Workspace, stack.InternetGateway.Metadata.Name),
		route:        generators.GenerateRouteTableURL(sdkconsts.NetworkProviderV1Name, stack.RouteTable.Metadata.Tenant, stack.RouteTable.Metadata.Workspace, stack.RouteTable.Metadata.Network, stack.RouteTable.Metadata.Name),
		subnet:       generators.GenerateSubnetURL(sdkconsts.NetworkProviderV1Name, stack.Subnet.Metadata.Tenant, stack.Subnet.Metadata.Workspace, stack.Subnet.Metadata.Network, stack.Subnet.Metadata.Name),
		group:        generators.GenerateSecurityGroupURL(sdkconsts.NetworkProviderV1Name, stack.SecurityGroup.Metadata.Tenant, stack.SecurityGroup.Metadata.Workspace, stack.SecurityGroup.Metadata.Name),
		publicIp:     generators.GeneratePublicIpURL(sdkconsts.NetworkProviderV1Name, stack.PublicIp.Metadata.Tenant, stack.PublicIp.Metadata.Workspace, stack.PublicIp.Metadata.Name),
		nic:          generators.GenerateNicURL(sdkconsts.NetworkProviderV1Name, stack.Nic.Metadata.Tenant, stack.Nic.Metadata.Workspace, stack.Nic.Metadata.Name),
		instance:     generators.GenerateInstanceURL(sdkconsts.ComputeProviderV1Name, stack.Instance.Metadata.Tenant, stack.Instance.Metadata.Workspace, stack.Instance.Metadata.Name),
		instanceList: generators.GenerateInstanceListURL(sdkconsts.ComputeProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name),
	}
}

// configureWorkspaceStackCreateV1, configureWorkspaceStackListV1 and
// configureWorkspaceStackDeleteV1 are split into three phases - rather than one
// combined per-stack function - because WireMock scenario state is a single
// sequential chain per suite. The runtime call order in TestScenario is
// create(A), create(B), list(A), list(B), delete(A), delete(B), so the mock
// registration order must follow that same interleaving across both stacks, not
// register one stack's full create+list+delete sequence before starting the next.

func configureWorkspaceStackCreateV1(configurator *stubs.Configurator, mockParams mock.MockParams, urls workspaceStackURLs, workspace *schema.Workspace, stack params.WorkspaceStackV1) error {
	if err := configurator.ConfigureCreateWorkspaceStub(workspace, urls.workspace, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspace, urls.workspace, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspace, urls.workspace, mockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateBlockStorageStub(stack.BlockStorage, urls.block, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingBlockStorageStub(stack.BlockStorage, urls.block, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveBlockStorageStub(stack.BlockStorage, urls.block, mockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateNetworkStub(stack.Network, urls.network, mockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateInternetGatewayStub(stack.InternetGateway, urls.gateway, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingInternetGatewayStub(stack.InternetGateway, urls.gateway, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInternetGatewayStub(stack.InternetGateway, urls.gateway, mockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateRouteTableStub(stack.RouteTable, urls.route, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingRouteTableStub(stack.RouteTable, urls.route, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRouteTableStub(stack.RouteTable, urls.route, mockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureGetCreatingNetworkStub(stack.Network, urls.network, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNetworkStub(stack.Network, urls.network, mockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateSubnetStub(stack.Subnet, urls.subnet, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingSubnetStub(stack.Subnet, urls.subnet, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSubnetStub(stack.Subnet, urls.subnet, mockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateSecurityGroupStub(stack.SecurityGroup, urls.group, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingSecurityGroupStub(stack.SecurityGroup, urls.group, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveSecurityGroupStub(stack.SecurityGroup, urls.group, mockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreatePublicIpStub(stack.PublicIp, urls.publicIp, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingPublicIpStub(stack.PublicIp, urls.publicIp, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActivePublicIpStub(stack.PublicIp, urls.publicIp, mockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateNicStub(stack.Nic, urls.nic, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingNicStub(stack.Nic, urls.nic, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveNicStub(stack.Nic, urls.nic, mockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateInstanceStub(stack.Instance, urls.instance, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingInstanceStub(stack.Instance, urls.instance, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInstanceStub(stack.Instance, urls.instance, mockParams); err != nil {
		return err
	}

	return nil
}

func configureWorkspaceStackListV1(configurator *stubs.Configurator, mockParams mock.MockParams, urls workspaceStackURLs, stack params.WorkspaceStackV1) error {
	instances := stack.Instances
	return configurator.ConfigureListInstanceStub(&instances, urls.instanceList, mockParams, nil)
}

func configureWorkspaceStackDeleteV1(configurator *stubs.Configurator, mockParams mock.MockParams, urls workspaceStackURLs) error {
	if err := configurator.ConfigureDeleteStub(urls.instance, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(urls.nic, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(urls.publicIp, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(urls.group, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(urls.subnet, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(urls.route, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(urls.gateway, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(urls.network, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(urls.block, mockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(urls.workspace, mockParams); err != nil {
		return err
	}
	return nil
}

func ConfigureMultiWorkspaceIsolationV1(scenario *mockscenarios.Scenario, p params.MultiWorkspaceIsolationV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	role := p.Role
	roleAssignment := p.RoleAssignment

	roleUrl := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, role.Metadata.Tenant, role.Metadata.Name)
	roleAssignUrl := generators.GenerateRoleAssignmentURL(sdkconsts.AuthorizationProviderV1Name, roleAssignment.Metadata.Tenant, roleAssignment.Metadata.Name)

	if err := configurator.ConfigureCreateRoleStub(role, roleUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingRoleStub(role, roleUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRoleStub(role, roleUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateRoleAssignmentStub(roleAssignment, roleAssignUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingRoleAssignmentStub(roleAssignment, roleAssignUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveRoleAssignmentStub(roleAssignment, roleAssignUrl, scenario.MockParams); err != nil {
		return err
	}

	urlsA := newWorkspaceStackURLs(p.WorkspaceA, p.StackA)
	urlsB := newWorkspaceStackURLs(p.WorkspaceB, p.StackB)

	// Matches TestScenario: createStack(A), createStack(B), List(A), List(B),
	// deleteStack(A), deleteStack(B)
	if err := configureWorkspaceStackCreateV1(configurator, scenario.MockParams, urlsA, p.WorkspaceA, p.StackA); err != nil {
		return err
	}
	if err := configureWorkspaceStackCreateV1(configurator, scenario.MockParams, urlsB, p.WorkspaceB, p.StackB); err != nil {
		return err
	}

	if err := configureWorkspaceStackListV1(configurator, scenario.MockParams, urlsA, p.StackA); err != nil {
		return err
	}
	if err := configureWorkspaceStackListV1(configurator, scenario.MockParams, urlsB, p.StackB); err != nil {
		return err
	}

	if err := configureWorkspaceStackDeleteV1(configurator, scenario.MockParams, urlsA); err != nil {
		return err
	}
	if err := configureWorkspaceStackDeleteV1(configurator, scenario.MockParams, urlsB); err != nil {
		return err
	}

	if err := configurator.ConfigureDeleteStub(roleAssignUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(roleUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
