package secatest

import (
	"context"
	"fmt"
	"time"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Steps

/// Create Or Update

type createOrUpdateTenantResourceParams[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string)
	operationName         string
	resource              *R
	createOrUpdateFunc    func(context.Context, *R) (*stepFuncResponse[R, M, E], error)
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
	expectedResourceState schema.ResourceState
}

type createOrUpdateWorkspaceResourceParams[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string, string)
	operationName         string
	workspace             string
	resource              *R
	createOrUpdateFunc    func(context.Context, *R) (*stepFuncResponse[R, M, E], error)
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
	expectedResourceState schema.ResourceState
}

type createOrUpdateNetworkResourceParams[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string, string, string)
	operationName         string
	workspace             string
	network               string
	resource              *R
	createOrUpdateFunc    func(context.Context, *R) (*stepFuncResponse[R, M, E], error)
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
	expectedResourceState schema.ResourceState
}

type createOrUpdateResourceParams[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType] struct {
	resource              *R
	createOrUpdateFunc    func(context.Context, *R) (*stepFuncResponse[R, M, E], error)
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
	expectedResourceState schema.ResourceState
}

/// Get

type getTenantResourceParams[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string)
	operationName         string
	tref                  secapi.TenantReference
	getFunc               func(context.Context, secapi.TenantReference, secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[R, M, E], error)
	expectedResourceState schema.ResourceState
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
}

type getWorkspaceResourceParams[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string, string)
	operationName         string
	wref                  secapi.WorkspaceReference
	getFunc               func(context.Context, secapi.WorkspaceReference, secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[R, M, E], error)
	expectedResourceState schema.ResourceState
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
}

type getNetworkResourceParams[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType] struct {
	stepName              string
	stepParamsFunc        func(provider.StepCtx, string, string, string)
	operationName         string
	nref                  secapi.NetworkReference
	getFunc               func(context.Context, secapi.NetworkReference, secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[R, M, E], error)
	expectedResourceState schema.ResourceState
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
}

type getResourceWithObserverParams[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType, F secapi.Reference, V any] struct {
	reference             F
	observerExpectedValue V
	getFunc               func(context.Context, F, secapi.ResourceObserverConfig[V]) (*stepFuncResponse[R, M, E], error)
	expectedLabels        schema.Labels
	expectedMetadata      *M
	verifyMetadataFunc    func(provider.StepCtx, *M, *M)
	expectedSpec          *E
	verifySpecFunc        func(provider.StepCtx, *E, *E)
	expectedResourceState schema.ResourceState
}

/// Response

type stepFuncResponse[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType] struct {
	resource *R
	labels   schema.Labels
	metadata *M
	spec     E
	state    *schema.ResourceState
}

func newStepFuncResponse[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType](resource *R, labels schema.Labels, metadata *M, spec E, state *schema.ResourceState) *stepFuncResponse[R, M, E] {
	return &stepFuncResponse[R, M, E]{
		resource: resource,
		labels:   labels,
		metadata: metadata,
		spec:     spec,
		state:    state,
	}
}

type responseExpects[M secapi.MetadataType, E secapi.SpecType] struct {
	labels        schema.Labels
	metadata      *M
	spec          *E
	resourceState schema.ResourceState
}

func createOrUpdateTenantResourceStep[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType](
	t provider.T,
	suite *testSuite,
	params createOrUpdateTenantResourceParams[R, M, E],
) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)

		createOrUpdateResourceStep(t, suite, sCtx, createOrUpdateResourceParams[R, M, E]{
			resource:              params.resource,
			createOrUpdateFunc:    params.createOrUpdateFunc,
			expectedLabels:        params.expectedLabels,
			expectedMetadata:      params.expectedMetadata,
			verifyMetadataFunc:    params.verifyMetadataFunc,
			expectedSpec:          params.expectedSpec,
			verifySpecFunc:        params.verifySpecFunc,
			expectedResourceState: params.expectedResourceState,
		})
	})
}

func createOrUpdateWorkspaceResourceStep[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType](
	t provider.T,
	suite *testSuite,
	params createOrUpdateWorkspaceResourceParams[R, M, E],
) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace)

		createOrUpdateResourceStep(t, suite, sCtx, createOrUpdateResourceParams[R, M, E]{
			resource:              params.resource,
			createOrUpdateFunc:    params.createOrUpdateFunc,
			expectedLabels:        params.expectedLabels,
			expectedMetadata:      params.expectedMetadata,
			verifyMetadataFunc:    params.verifyMetadataFunc,
			expectedSpec:          params.expectedSpec,
			verifySpecFunc:        params.verifySpecFunc,
			expectedResourceState: params.expectedResourceState,
		})
	})
}

func createOrUpdateNetworkResourceStep[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType](
	t provider.T,
	suite *testSuite,
	params createOrUpdateNetworkResourceParams[R, M, E],
) {
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, params.workspace, params.network)

		createOrUpdateResourceStep(t, suite, sCtx, createOrUpdateResourceParams[R, M, E]{
			resource:              params.resource,
			createOrUpdateFunc:    params.createOrUpdateFunc,
			expectedLabels:        params.expectedLabels,
			expectedMetadata:      params.expectedMetadata,
			verifyMetadataFunc:    params.verifyMetadataFunc,
			expectedSpec:          params.expectedSpec,
			verifySpecFunc:        params.verifySpecFunc,
			expectedResourceState: params.expectedResourceState,
		})
	})
}

func createOrUpdateResourceStep[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType](
	t provider.T,
	suite *testSuite,
	sCtx provider.StepCtx,
	params createOrUpdateResourceParams[R, M, E],
) {
	resp, err := params.createOrUpdateFunc(t.Context(), params.resource)
	requireNoError(sCtx, err)
	if params.expectedLabels != nil {
		suite.verifyLabelsStep(sCtx, params.expectedLabels, resp.labels)
	}

	if params.expectedMetadata != nil {
		params.verifyMetadataFunc(sCtx, params.expectedMetadata, resp.metadata)
	}

	if params.expectedSpec != nil {
		params.verifySpecFunc(sCtx, params.expectedSpec, &resp.spec)
	}

	suite.verifyStatusStep(sCtx, params.expectedResourceState, *resp.state)
}

func getTenantResourceStep[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType](
	t provider.T,
	suite *testSuite,
	params getTenantResourceParams[R, M, E],
) *R {
	var resp *stepFuncResponse[R, M, E]
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName)

		resp = getResourceWithObserver(t, suite, sCtx,
			getResourceWithObserverParams[R, M, E, secapi.TenantReference, schema.ResourceState]{
				reference:             params.tref,
				observerExpectedValue: params.expectedResourceState,
				getFunc:               params.getFunc,
				expectedLabels:        params.expectedLabels,
				expectedMetadata:      params.expectedMetadata,
				verifyMetadataFunc:    params.verifyMetadataFunc,
				expectedSpec:          params.expectedSpec,
				verifySpecFunc:        params.verifySpecFunc,
				expectedResourceState: params.expectedResourceState,
			},
		)
		requireNotNilResponse(sCtx, resp)
	})
	return resp.resource
}

func getWorkspaceResourceStep[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType](
	t provider.T,
	suite *testSuite,
	params getWorkspaceResourceParams[R, M, E],
) *R {
	var resp *stepFuncResponse[R, M, E]
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, string(params.wref.Workspace))

		resp = getResourceWithObserver(t, suite, sCtx,
			getResourceWithObserverParams[R, M, E, secapi.WorkspaceReference, schema.ResourceState]{
				reference:             params.wref,
				observerExpectedValue: params.expectedResourceState,
				getFunc:               params.getFunc,
				expectedLabels:        params.expectedLabels,
				expectedMetadata:      params.expectedMetadata,
				verifyMetadataFunc:    params.verifyMetadataFunc,
				expectedSpec:          params.expectedSpec,
				verifySpecFunc:        params.verifySpecFunc,
				expectedResourceState: params.expectedResourceState,
			},
		)
		requireNotNilResponse(sCtx, resp)
	})
	return resp.resource
}

func getNetworkResourceStep[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType](
	t provider.T,
	suite *testSuite,
	params getNetworkResourceParams[R, M, E],
) *R {
	var resp *stepFuncResponse[R, M, E]
	t.WithNewStep(params.stepName, func(sCtx provider.StepCtx) {
		params.stepParamsFunc(sCtx, params.operationName, string(params.nref.Workspace), string(params.nref.Network))

		resp = getResourceWithObserver(t, suite, sCtx,
			getResourceWithObserverParams[R, M, E, secapi.NetworkReference, schema.ResourceState]{
				reference:             params.nref,
				observerExpectedValue: params.expectedResourceState,
				getFunc:               params.getFunc,
				expectedLabels:        params.expectedLabels,
				expectedMetadata:      params.expectedMetadata,
				verifyMetadataFunc:    params.verifyMetadataFunc,
				expectedSpec:          params.expectedSpec,
				verifySpecFunc:        params.verifySpecFunc,
				expectedResourceState: params.expectedResourceState,
			},
		)
		requireNotNilResponse(sCtx, resp)
	})
	return resp.resource
}

func getResourceWithObserver[R secapi.ResourceType, M secapi.MetadataType, E secapi.SpecType, F secapi.Reference, V any](
	t provider.T,
	suite *testSuite,
	sCtx provider.StepCtx,
	params getResourceWithObserverParams[R, M, E, F, V],
) *stepFuncResponse[R, M, E] {
	config := secapi.ResourceObserverConfig[V]{
		ExpectedValue: params.observerExpectedValue,
		Delay:         time.Duration(suite.baseDelay) * time.Second,
		Interval:      time.Duration(suite.baseInterval) * time.Second,
		MaxAttempts:   suite.maxAttempts,
	}

	resp, err := params.getFunc(t.Context(), params.reference, config)
	requireNoError(sCtx, err)
	requireNotNilResponse(sCtx, resp)

	if params.expectedLabels != nil {
		suite.verifyLabelsStep(sCtx, params.expectedLabels, resp.labels)
	}

	if params.expectedMetadata != nil {
		params.verifyMetadataFunc(sCtx, params.expectedMetadata, resp.metadata)
	}

	if params.expectedSpec != nil {
		params.verifySpecFunc(sCtx, params.expectedSpec, &resp.spec)
	}

	suite.verifyStatusStep(sCtx, params.expectedResourceState, *resp.state)

	return resp
}

// Metadata

func (suite *testSuite) verifyGlobalTenantResourceMetadataStep(ctx provider.StepCtx, expected *schema.GlobalTenantResourceMetadata, actual *schema.GlobalTenantResourceMetadata) {
	ctx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.Name, actual.Name, "Name should match expected")
		stepCtx.Require().Equal(expected.Provider, actual.Provider, "Provider should match expected")
		stepCtx.Require().Equal(expected.Resource, actual.Resource, "Resource should match expected")
		stepCtx.Require().Equal(expected.ApiVersion, actual.ApiVersion, "ApiVersion should match expected")
		stepCtx.Require().Equal(expected.Verb, actual.Verb, "Verb should match expected")
		stepCtx.Require().Equal(expected.Kind, actual.Kind, "Kind should match expected")
		stepCtx.Require().Equal(expected.Tenant, actual.Tenant, "Tenant should match expected")
	})
}

func (suite *testSuite) verifyRegionalResourceMetadataStep(ctx provider.StepCtx, expected *schema.RegionalResourceMetadata, actual *schema.RegionalResourceMetadata) {
	ctx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.Name, actual.Name, "Name should match expected")
		stepCtx.Require().Equal(expected.Provider, actual.Provider, "Provider should match expected")
		stepCtx.Require().Equal(expected.Resource, actual.Resource, "Resource should match expected")
		stepCtx.Require().Equal(expected.ApiVersion, actual.ApiVersion, "ApiVersion should match expected")
		stepCtx.Require().Equal(expected.Verb, actual.Verb, "Verb should match expected")
		stepCtx.Require().Equal(expected.Kind, actual.Kind, "Kind should match expected")
		stepCtx.Require().Equal(expected.Tenant, actual.Tenant, "Tenant should match expected")
		stepCtx.Require().Equal(expected.Region, actual.Region, "Region should match expected")
	})
}

func (suite *testSuite) verifyRegionalWorkspaceResourceMetadataStep(ctx provider.StepCtx, expected *schema.RegionalWorkspaceResourceMetadata, actual *schema.RegionalWorkspaceResourceMetadata) {
	ctx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.Name, actual.Name, "Name should match expected")
		stepCtx.Require().Equal(expected.Provider, actual.Provider, "Provider should match expected")
		stepCtx.Require().Equal(expected.Resource, actual.Resource, "Resource should match expected")
		stepCtx.Require().Equal(expected.ApiVersion, actual.ApiVersion, "ApiVersion should match expected")
		stepCtx.Require().Equal(expected.Verb, actual.Verb, "Verb should match expected")
		stepCtx.Require().Equal(expected.Kind, actual.Kind, "Kind should match expected")
		stepCtx.Require().Equal(expected.Tenant, actual.Tenant, "Tenant should match expected")
		stepCtx.Require().Equal(expected.Workspace, actual.Workspace, "Workspace should match expected")
		stepCtx.Require().Equal(expected.Region, actual.Region, "Region should match expected")
	})
}

func (suite *testSuite) verifyRegionalNetworkResourceMetadataStep(ctx provider.StepCtx, expected *schema.RegionalNetworkResourceMetadata, actual *schema.RegionalNetworkResourceMetadata) {
	ctx.WithNewStep("Verify metadata", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.Name, actual.Name, "Name should match expected")
		stepCtx.Require().Equal(expected.Provider, actual.Provider, "Provider should match expected")
		stepCtx.Require().Equal(expected.Resource, actual.Resource, "Resource should match expected")
		stepCtx.Require().Equal(expected.ApiVersion, actual.ApiVersion, "ApiVersion should match expected")
		stepCtx.Require().Equal(expected.Verb, actual.Verb, "Verb should match expected")
		stepCtx.Require().Equal(expected.Kind, actual.Kind, "Kind should match expected")
		stepCtx.Require().Equal(expected.Tenant, actual.Tenant, "Tenant should match expected")
		stepCtx.Require().Equal(expected.Workspace, actual.Workspace, "Workspace should match expected")
		stepCtx.Require().Equal(expected.Network, actual.Network, "Network should match expected")
		stepCtx.Require().Equal(expected.Region, actual.Region, "Region should match expected")
	})
}

// Status

func (suite *testSuite) verifyStatusStep(ctx provider.StepCtx, expected schema.ResourceState, actual schema.ResourceState) {
	ctx.WithNewStep("Verify status state", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected, actual, "Status state should match expected")
	})
}

func (suite *testSuite) verifyLabelsStep(ctx provider.StepCtx, expected schema.Labels, actual schema.Labels) {
	ctx.WithNewStep("Verify label", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_labels", expected,
			"actual_labels", actual,
		)

		stepCtx.Require().Equal(expected, actual, "Labels should match expected")
	})
}

// Specs

// Authorization Specs

func (suite *testSuite) verifyRoleSpecStep(ctx provider.StepCtx, expected *schema.RoleSpec, actual *schema.RoleSpec) {
	ctx.WithNewStep("Verify RoleSpec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(len(expected.Permissions), len(actual.Permissions),
			"Permissions list length should match expected")

		for i := 0; i < len(expected.Permissions); i++ {
			expectedPerm := expected.Permissions[i]
			actualPerm := actual.Permissions[i]
			stepCtx.Require().Equal(expectedPerm.Provider, actualPerm.Provider,
				fmt.Sprintf("Permission [%d] provider should match expected", i))
			stepCtx.Require().Equal(expectedPerm.Resources, actualPerm.Resources,
				fmt.Sprintf("Permission [%d] resources should match expected", i))
			stepCtx.Require().Equal(expectedPerm.Verb, actualPerm.Verb,
				fmt.Sprintf("Permission [%d] verb should match expected", i))
		}
	})
}

func (suite *testSuite) verifyRoleAssignmentSpecStep(ctx provider.StepCtx, expected *schema.RoleAssignmentSpec, actual *schema.RoleAssignmentSpec) {
	ctx.WithNewStep("Verify RoleAssignmentSpec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.Roles, actual.Roles, "Roles provider should match expected")
		stepCtx.Require().Equal(expected.Subs, actual.Subs, "Subs should match expected")
		stepCtx.Require().Equal(len(expected.Scopes), len(actual.Scopes), "Scope list length should match expected")

		for i := 0; i < len(expected.Scopes); i++ {
			expectedScope := expected.Scopes[i]
			actualScope := actual.Scopes[i]

			if actualScope.Tenants != nil && len(*actualScope.Tenants) > 0 {
				stepCtx.Require().Equal(expectedScope.Tenants, actualScope.Tenants, fmt.Sprintf("Scope [%d] tenants should match expected", i))
			}
			if actualScope.Regions != nil && len(*actualScope.Regions) > 0 {
				stepCtx.Require().Equal(expectedScope.Regions, actualScope.Regions, fmt.Sprintf("Scope [%d] regions should match expected", i))
			}
			if actualScope.Workspaces != nil && len(*actualScope.Workspaces) > 0 {
				stepCtx.Require().Equal(expectedScope.Workspaces, actualScope.Workspaces, fmt.Sprintf("Scope [%d] workspaces should match expected", i))
			}
		}
	})
}

// Storage Specs

func (suite *testSuite) verifyBlockStorageSpecStep(ctx provider.StepCtx, expected *schema.BlockStorageSpec, actual *schema.BlockStorageSpec) {
	ctx.WithNewStep("Verify BlockStorageSpec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.SizeGB, actual.SizeGB, "SizeGB should match expected")
		stepCtx.Require().Equal(expected.SkuRef, actual.SkuRef, "SkuRef should match expected")
	})
}

func (suite *testSuite) verifyImageSpecStep(ctx provider.StepCtx, expected *schema.ImageSpec, actual *schema.ImageSpec) {
	ctx.WithNewStep("Verify ImageSpec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.BlockStorageRef, actual.BlockStorageRef, "BlockStorageRef should match expected")
		stepCtx.Require().Equal(expected.CpuArchitecture, actual.CpuArchitecture, "CpuArchitecture should match expected")
	})
}

// Compute Specs

func (suite *testSuite) verifyInstanceSpecStep(ctx provider.StepCtx, expected *schema.InstanceSpec, actual *schema.InstanceSpec) {
	ctx.WithNewStep("Verify InstanceSpec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.SkuRef, actual.SkuRef, "SkuRef should match expected")
		stepCtx.Require().Equal(expected.Zone, actual.Zone, "Zone should match expected")
		stepCtx.Require().Equal(expected.BootVolume.DeviceRef, actual.BootVolume.DeviceRef, "BootVolume.DeviceRef should match expected")
	})
}

// Network Specs

func (suite *testSuite) verifyNetworkSpecStep(ctx provider.StepCtx, expected *schema.NetworkSpec, actual *schema.NetworkSpec) {
	ctx.WithNewStep("Verify NetworkSpec", func(stepCtx provider.StepCtx) {
		if actual.Cidr.Ipv4 != nil {
			stepCtx.Require().Equal(expected.Cidr.Ipv4, actual.Cidr.Ipv4, "Cidr.Ipv4 should match expected")
		}
		if actual.Cidr.Ipv6 != nil {
			stepCtx.Require().Equal(expected.Cidr.Ipv6, actual.Cidr.Ipv6, "Cidr.Ipv6 should match expected")
		}

		stepCtx.Require().Equal(expected.SkuRef, actual.SkuRef, "SkuRef should match expected")
		stepCtx.Require().Equal(expected.RouteTableRef, actual.RouteTableRef, "RouteTableRef should match expected")
	})
}

func (suite *testSuite) verifyInternetGatewaySpecStep(ctx provider.StepCtx, expected *schema.InternetGatewaySpec, actual *schema.InternetGatewaySpec) {
	ctx.WithNewStep("Verify InternetGatewaySpec", func(stepCtx provider.StepCtx) {
		if actual.EgressOnly != nil {
			stepCtx.Require().Equal(expected.EgressOnly, actual.EgressOnly, "EgressOnly should match expected")
		}
	})
}

func (suite *testSuite) verifyRouteTableSpecStep(ctx provider.StepCtx, expected *schema.RouteTableSpec, actual *schema.RouteTableSpec) {
	ctx.WithNewStep("Verify RouteTableSpec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(len(expected.Routes), len(actual.Routes), "Route list length should match expected")
		for i := 0; i < len(expected.Routes); i++ {
			expectedRoute := expected.Routes[i]
			actualRoute := actual.Routes[i]
			stepCtx.Require().Equal(expectedRoute.DestinationCidrBlock, actualRoute.DestinationCidrBlock, fmt.Sprintf("Route [%d] DestinationCidrBlock should match expected", i))
			stepCtx.Require().Equal(expectedRoute.TargetRef, actualRoute.TargetRef, fmt.Sprintf("Route [%d] TargetRef should match expected", i))
		}
	})
}

func (suite *testSuite) verifySubnetSpecStep(ctx provider.StepCtx, expected *schema.SubnetSpec, actual *schema.SubnetSpec) {
	ctx.WithNewStep("Verify SubnetSpec", func(stepCtx provider.StepCtx) {
		if actual.Cidr.Ipv4 != nil {
			stepCtx.Require().Equal(expected.Cidr.Ipv4, actual.Cidr.Ipv4, "Cidr.Ipv4 should match expected")
		}
		if actual.Cidr.Ipv6 != nil {
			stepCtx.Require().Equal(expected.Cidr.Ipv6, actual.Cidr.Ipv6, "Cidr.Ipv6 should match expected")
		}
		stepCtx.Require().Equal(expected.Zone, actual.Zone, "Zone should match expected")
	})
}

func (suite *testSuite) verifyPublicIpSpecStep(ctx provider.StepCtx, expected *schema.PublicIpSpec, actual *schema.PublicIpSpec) {
	ctx.WithNewStep("Verify PublicIpSpec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.Version, actual.Version, "Version should match expected")
		if actual.Address != nil {
			stepCtx.Require().Equal(expected.Address, actual.Address, "Address should match expected")
		}
	})
}

func (suite *testSuite) verifyNicSpecStep(ctx provider.StepCtx, expected *schema.NicSpec, actual *schema.NicSpec) {
	ctx.WithNewStep("Verify NicSpec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.Addresses, actual.Addresses, "Addresses should match expected")
		if actual.PublicIpRefs != nil {
			stepCtx.Require().Equal(expected.PublicIpRefs, actual.PublicIpRefs, "PublicIpRefs should match expected")
		}
		stepCtx.Require().Equal(expected.SubnetRef, actual.SubnetRef, "SubnetRef should match expected")
	})
}

func (suite *testSuite) verifySecurityGroupSpecStep(ctx provider.StepCtx, expected *schema.SecurityGroupSpec, actual *schema.SecurityGroupSpec) {
	ctx.WithNewStep("Verify SecurityGroupSpec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(len(expected.Rules), len(actual.Rules), "Rule list length should match expected")
		for i := 0; i < len(expected.Rules); i++ {
			expectedRule := expected.Rules[i]
			actualRule := actual.Rules[i]
			stepCtx.Require().Equal(expectedRule.Direction, actualRule.Direction, fmt.Sprintf("Rule [%d] Direction should match expected", i))
		}
	})
}
