package secatest

import (
	"context"
	"fmt"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Steps

type stepFuncResponse[R, M, E any] struct {
	resource *R
	labels   schema.Labels
	metadata *M
	spec     E
	state    *schema.ResourceState
}

func newStepFuncResponse[R, M, E any](resource *R, labels schema.Labels, metadata *M, spec E, state *schema.ResourceState) *stepFuncResponse[R, M, E] {
	return &stepFuncResponse[R, M, E]{
		resource: resource,
		labels:   labels,
		metadata: metadata,
		spec:     spec,
		state:    state,
	}
}

func createOrUpdateResourceStep[R, M, E any](
	t provider.T,
	ctx context.Context,
	suite *testSuite,
	stepName string,
	stepParamsFunc func(provider.StepCtx, string),
	operationName string,
	resource *R,
	createOrUpdateFunc func(context.Context, *R) (*stepFuncResponse[R, M, E], error),
	expectedLabels schema.Labels,
	expectedMetadata *M,
	verifyMetadataFunc func(provider.StepCtx, *M, *M),
	expectedSpec *E,
	verifySpecFunc func(provider.StepCtx, *E, *E),
	expectedResourceState schema.ResourceState,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		stepParamsFunc(sCtx, operationName)

		resp, err := createOrUpdateFunc(ctx, resource)
		requireNoError(sCtx, err)

		if expectedLabels != nil {
			suite.verifyLabelsStep(sCtx, expectedLabels, resp.labels)
		}

		if expectedMetadata != nil {
			verifyMetadataFunc(sCtx, expectedMetadata, resp.metadata)
		}

		if expectedSpec != nil {
			verifySpecFunc(sCtx, expectedSpec, &resp.spec)
		}

		suite.verifyStatusStep(sCtx, expectedResourceState, *resp.state)
	})
}

func createOrUpdateWorkspaceResourceStep[R, M, E any](
	t provider.T,
	ctx context.Context,
	suite *testSuite,
	stepName string,
	stepParamsFunc func(provider.StepCtx, string, string),
	operationName string,
	workspace string,
	resource *R,
	createOrUpdateFunc func(context.Context, *R) (*stepFuncResponse[R, M, E], error),
	expectedLabels schema.Labels,
	expectedMetadata *M,
	verifyMetadataFunc func(provider.StepCtx, *M, *M),
	expectedSpec *E,
	verifySpecFunc func(provider.StepCtx, *E, *E),
	expectedResourceState schema.ResourceState,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		stepParamsFunc(sCtx, operationName, workspace)

		resp, err := createOrUpdateFunc(ctx, resource)
		requireNoError(sCtx, err)

		if expectedLabels != nil {
			suite.verifyLabelsStep(sCtx, expectedLabels, resp.labels)
		}

		if expectedMetadata != nil {
			verifyMetadataFunc(sCtx, expectedMetadata, resp.metadata)
		}

		if expectedSpec != nil {
			verifySpecFunc(sCtx, expectedSpec, &resp.spec)
		}

		suite.verifyStatusStep(sCtx, expectedResourceState, *resp.state)
	})
}

func getTenantResourceStep[R, M, E any](
	t provider.T,
	ctx context.Context,
	suite *testSuite,
	stepName string,
	stepParamsFunc func(provider.StepCtx, string),
	operationName string,
	tref secapi.TenantReference,
	getFunc func(context.Context, secapi.TenantReference) (*stepFuncResponse[R, M, E], error),
	expectedLabels schema.Labels,
	expectedMetadata *M,
	verifyMetadataFunc func(provider.StepCtx, *M, *M),
	expectedSpec *E,
	verifySpecFunc func(provider.StepCtx, *E, *E),
	expectedResourceState schema.ResourceState,
) *R {
	var resp *stepFuncResponse[R, M, E]

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		stepParamsFunc(sCtx, operationName)

		retry := newStepResourceStateRetry(
			suite.baseDelay,
			suite.baseInterval,
			suite.maxAttempts,
			func() (schema.ResourceState, error) {
				var err error
				resp, err = getFunc(ctx, tref)
				requireNoError(sCtx, err)

				return *resp.state, nil
			},
			func() {
				if expectedLabels != nil {
					suite.verifyLabelsStep(sCtx, expectedLabels, resp.labels)
				}

				if expectedMetadata != nil {
					verifyMetadataFunc(sCtx, expectedMetadata, resp.metadata)
				}

				if expectedSpec != nil {
					verifySpecFunc(sCtx, expectedSpec, &resp.spec)
				}

				suite.verifyStatusStep(sCtx, expectedResourceState, *resp.state)
			},
		)
		retry.run(sCtx, operationName, expectedResourceState)
	})
	return resp.resource
}

func getWorkspaceResourceStep[R, M, E any](
	t provider.T,
	ctx context.Context,
	suite *testSuite,
	stepName string,
	stepParamsFunc func(provider.StepCtx, string, string),
	operationName string,
	wref secapi.WorkspaceReference,
	getFunc func(context.Context, secapi.WorkspaceReference) (*stepFuncResponse[R, M, E], error),
	expectedLabels schema.Labels,
	expectedMetadata *M,
	verifyMetadataFunc func(provider.StepCtx, *M, *M),
	expectedSpec *E,
	verifySpecFunc func(provider.StepCtx, *E, *E),
	expectedResourceState schema.ResourceState,
) *R {
	var resp *stepFuncResponse[R, M, E]

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		stepParamsFunc(sCtx, operationName, string(wref.Workspace))

		retry := newStepResourceStateRetry(
			suite.baseDelay,
			suite.baseInterval,
			suite.maxAttempts,
			func() (schema.ResourceState, error) {
				var err error
				resp, err = getFunc(ctx, wref)
				requireNoError(sCtx, err)

				return *resp.state, nil
			},
			func() {
				if expectedLabels != nil {
					suite.verifyLabelsStep(sCtx, expectedLabels, resp.labels)
				}

				if expectedMetadata != nil {
					verifyMetadataFunc(sCtx, expectedMetadata, resp.metadata)
				}

				if expectedSpec != nil {
					verifySpecFunc(sCtx, expectedSpec, &resp.spec)
				}

				suite.verifyStatusStep(sCtx, expectedResourceState, *resp.state)
			},
		)
		retry.run(sCtx, operationName, expectedResourceState)
	})
	return resp.resource
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

func (suite *testSuite) requireNotNilStatus(sCtx provider.StepCtx, status any) {
	sCtx.WithNewStep("Verify not nil status", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters("status", fmt.Sprintf("%v", status))
		stepCtx.Require().NotNil(status, "Should be not nil")
	})
}

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
