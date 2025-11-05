package secatest

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Network

func (suite *testSuite) createOrUpdateNetworkV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	resource *schema.Network,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.NetworkSpec,
	expectedState schema.ResourceState,
) {
	expectedMeta.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(
		t,
		ctx,
		suite,
		stepName,
		suite.setNetworkV1StepParams,
		"CreateOrUpdateNetwork",
		resource.Metadata.Workspace,
		resource,
		func(context.Context, *schema.Network) (*stepFuncResponse[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec], error) {
			resp, err := api.CreateOrUpdateNetwork(ctx, resource)
			return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
		},
		nil,
		expectedMeta,
		suite.verifyRegionalWorkspaceResourceMetadataStep,
		expectedSpec,
		suite.verifyNetworkSpecStep,
		expectedState,
	)
}

func (suite *testSuite) getNetworkV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.NetworkSpec,
	expectedState schema.ResourceState,
) *schema.Network {
	var resp *schema.Network

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNetwork", string(wref.Workspace))
		retry := newStepResourceStateRetry(
			suite.baseDelay,
			suite.baseInterval,
			suite.maxAttempts,
			func() (schema.ResourceState, error) {
				var err error
				resp, err = api.GetNetwork(ctx, wref)
				requireNoError(sCtx, err)
				requireNotNilResponse(sCtx, resp)

				suite.requireNotNilStatus(sCtx, resp.Status)
				return *resp.Status.State, nil
			},
			func() {
				expectedMeta.Verb = http.MethodGet
				suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

				suite.verifyNetworkSpecStep(sCtx, expectedSpec, &resp.Spec)

				suite.verifyStatusStep(sCtx, expectedState, *resp.Status.State)
			},
		)
		retry.run(sCtx, "GetNetwork", expectedState)
	})
	return resp
}

func (suite *testSuite) getNetworkWithErrorV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	expectedError error,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNetwork", string(wref.Workspace))

		_, err := api.GetNetwork(ctx, wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteNetworkV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.NetworkV1, resource *schema.Network) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteNetwork", resource.Metadata.Workspace)

		err := api.DeleteNetwork(ctx, resource)
		requireNoError(sCtx, err)
	})
}

// Internet Gateway

func (suite *testSuite) createOrUpdateInternetGatewayV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	resource *schema.InternetGateway,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.InternetGatewaySpec,
	expectedState schema.ResourceState,
) {
	expectedMeta.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(
		t,
		ctx,
		suite,
		stepName,
		suite.setNetworkV1StepParams,
		"CreateOrUpdateInternetGateway",
		resource.Metadata.Workspace,
		resource,
		func(context.Context, *schema.InternetGateway) (*stepFuncResponse[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec], error) {
			resp, err := api.CreateOrUpdateInternetGateway(ctx, resource)
			return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
		},
		nil,
		expectedMeta,
		suite.verifyRegionalWorkspaceResourceMetadataStep,
		expectedSpec,
		suite.verifyInternetGatewaySpecStep,
		expectedState,
	)
}

func (suite *testSuite) getInternetGatewayV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.InternetGatewaySpec,
	expectedState schema.ResourceState,
) *schema.InternetGateway {
	var resp *schema.InternetGateway

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetInternetGateway", string(wref.Workspace))
		retry := newStepResourceStateRetry(
			suite.baseDelay,
			suite.baseInterval,
			suite.maxAttempts,
			func() (schema.ResourceState, error) {
				var err error
				resp, err = api.GetInternetGateway(ctx, wref)
				requireNoError(sCtx, err)
				requireNotNilResponse(sCtx, resp)

				suite.requireNotNilStatus(sCtx, resp.Status)
				return *resp.Status.State, nil
			},
			func() {
				expectedMeta.Verb = http.MethodGet
				suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

				suite.verifyInternetGatewaySpecStep(sCtx, expectedSpec, &resp.Spec)

				suite.verifyStatusStep(sCtx, expectedState, *resp.Status.State)
			},
		)
		retry.run(sCtx, "GetInternetGateway", expectedState)
	})
	return resp
}

func (suite *testSuite) getInternetGatewayWithErrorV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	expectedError error,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetInternetGateway", string(wref.Workspace))

		_, err := api.GetInternetGateway(ctx, wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteInternetGatewayV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.NetworkV1, resource *schema.InternetGateway) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteInternetGateway", resource.Metadata.Workspace)

		err := api.DeleteInternetGateway(ctx, resource)
		requireNoError(sCtx, err)
	})
}

// Route Table

func (suite *testSuite) createOrUpdateRouteTableV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	resource *schema.RouteTable,
	expectedMeta *schema.RegionalNetworkResourceMetadata,
	expectedSpec *schema.RouteTableSpec,
	expectedState schema.ResourceState,
) {
	expectedMeta.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(
		t,
		ctx,
		suite,
		stepName,
		suite.setNetworkV1StepParams,
		"CreateOrUpdateRouteTable",
		resource.Metadata.Workspace,
		resource,
		func(context.Context, *schema.RouteTable) (*stepFuncResponse[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec], error) {
			resp, err := api.CreateOrUpdateRouteTable(ctx, resource)
			return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
		},
		nil,
		expectedMeta,
		suite.verifyRegionalNetworkResourceMetadataStep,
		expectedSpec,
		suite.verifyRouteTableSpecStep,
		expectedState,
	)
}

func (suite *testSuite) getRouteTableV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	nref secapi.NetworkReference,
	expectedMeta *schema.RegionalNetworkResourceMetadata,
	expectedSpec *schema.RouteTableSpec,
	expectedState schema.ResourceState,
) *schema.RouteTable {
	var resp *schema.RouteTable

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetRouteTable", string(nref.Workspace))
		retry := newStepResourceStateRetry(
			suite.baseDelay,
			suite.baseInterval,
			suite.maxAttempts,
			func() (schema.ResourceState, error) {
				var err error
				resp, err = api.GetRouteTable(ctx, nref)
				requireNoError(sCtx, err)
				requireNotNilResponse(sCtx, resp)

				suite.requireNotNilStatus(sCtx, resp.Status)
				return *resp.Status.State, nil
			},
			func() {
				expectedMeta.Verb = http.MethodGet
				suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

				suite.verifyRouteTableSpecStep(sCtx, expectedSpec, &resp.Spec)

				suite.verifyStatusStep(sCtx, expectedState, *resp.Status.State)
			},
		)
		retry.run(sCtx, "GetRouteTable", expectedState)
	})
	return resp
}

func (suite *testSuite) getRouteTableWithErrorV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	nref secapi.NetworkReference,
	expectedError error,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetRouteTable", string(nref.Workspace))

		_, err := api.GetRouteTable(ctx, nref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteRouteTableV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.NetworkV1, resource *schema.RouteTable) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteRouteTable", resource.Metadata.Workspace)

		err := api.DeleteRouteTable(ctx, resource)
		requireNoError(sCtx, err)
	})
}

// Subnet

func (suite *testSuite) createOrUpdateSubnetV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	resource *schema.Subnet,
	expectedMeta *schema.RegionalNetworkResourceMetadata,
	expectedSpec *schema.SubnetSpec,
	expectedState schema.ResourceState,
) {
	expectedMeta.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(
		t,
		ctx,
		suite,
		stepName,
		suite.setNetworkV1StepParams,
		"CreateOrUpdateSubnet",
		resource.Metadata.Workspace,
		resource,
		func(context.Context, *schema.Subnet) (*stepFuncResponse[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec], error) {
			resp, err := api.CreateOrUpdateSubnet(ctx, resource)
			return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
		},
		nil,
		expectedMeta,
		suite.verifyRegionalNetworkResourceMetadataStep,
		expectedSpec,
		suite.verifySubnetSpecStep,
		expectedState,
	)
}

func (suite *testSuite) getSubnetV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	nref secapi.NetworkReference,
	expectedMeta *schema.RegionalNetworkResourceMetadata,
	expectedSpec *schema.SubnetSpec,
	expectedState schema.ResourceState,
) *schema.Subnet {
	var resp *schema.Subnet

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSubnet", string(nref.Workspace))
		retry := newStepResourceStateRetry(
			suite.baseDelay,
			suite.baseInterval,
			suite.maxAttempts,
			func() (schema.ResourceState, error) {
				var err error
				resp, err = api.GetSubnet(ctx, nref)
				requireNoError(sCtx, err)
				requireNotNilResponse(sCtx, resp)

				suite.requireNotNilStatus(sCtx, resp.Status)
				return *resp.Status.State, nil
			},
			func() {
				expectedMeta.Verb = http.MethodGet
				suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

				suite.verifySubnetSpecStep(sCtx, expectedSpec, &resp.Spec)

				suite.verifyStatusStep(sCtx, expectedState, *resp.Status.State)
			},
		)
		retry.run(sCtx, "GetSubnet", expectedState)
	})
	return resp
}

func (suite *testSuite) getSubnetWithErrorV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	nref secapi.NetworkReference,
	expectedError error,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSubnet", string(nref.Workspace))

		_, err := api.GetSubnet(ctx, nref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteSubnetV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.NetworkV1, resource *schema.Subnet) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteSubnet", resource.Metadata.Workspace)

		err := api.DeleteSubnet(ctx, resource)
		requireNoError(sCtx, err)
	})
}

// Public Ip

func (suite *testSuite) createOrUpdatePublicIpV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	resource *schema.PublicIp,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.PublicIpSpec,
	expectedState schema.ResourceState,
) {
	expectedMeta.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(
		t,
		ctx,
		suite,
		stepName,
		suite.setNetworkV1StepParams,
		"CreateOrUpdatePublicIp",
		resource.Metadata.Workspace,
		resource,
		func(context.Context, *schema.PublicIp) (*stepFuncResponse[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec], error) {
			resp, err := api.CreateOrUpdatePublicIp(ctx, resource)
			return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
		},
		nil,
		expectedMeta,
		suite.verifyRegionalWorkspaceResourceMetadataStep,
		expectedSpec,
		suite.verifyPublicIpSpecStep,
		expectedState,
	)
}

func (suite *testSuite) getPublicIpV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.PublicIpSpec,
	expectedState schema.ResourceState,
) *schema.PublicIp {
	var resp *schema.PublicIp

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetPublicIp", string(wref.Workspace))
		retry := newStepResourceStateRetry(
			suite.baseDelay,
			suite.baseInterval,
			suite.maxAttempts,
			func() (schema.ResourceState, error) {
				var err error
				resp, err = api.GetPublicIp(ctx, wref)
				requireNoError(sCtx, err)
				requireNotNilResponse(sCtx, resp)

				suite.requireNotNilStatus(sCtx, resp.Status)
				return *resp.Status.State, nil
			},
			func() {
				expectedMeta.Verb = http.MethodGet
				suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

				suite.verifyPublicIpSpecStep(sCtx, expectedSpec, &resp.Spec)

				suite.verifyStatusStep(sCtx, expectedState, *resp.Status.State)
			},
		)
		retry.run(sCtx, "GetPublicIp", expectedState)
	})
	return resp
}

func (suite *testSuite) getPublicIpWithErrorV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	expectedError error,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetPublicIp", string(wref.Workspace))

		_, err := api.GetPublicIp(ctx, wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deletePublicIpV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.NetworkV1, resource *schema.PublicIp) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeletePublicIp", resource.Metadata.Workspace)

		err := api.DeletePublicIp(ctx, resource)
		requireNoError(sCtx, err)
	})
}

// Nic

func (suite *testSuite) createOrUpdateNicV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	resource *schema.Nic,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.NicSpec,
	expectedState schema.ResourceState,
) {
	expectedMeta.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(
		t,
		ctx,
		suite,
		stepName,
		suite.setNetworkV1StepParams,
		"CreateOrUpdateNic",
		resource.Metadata.Workspace,
		resource,
		func(context.Context, *schema.Nic) (*stepFuncResponse[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec], error) {
			resp, err := api.CreateOrUpdateNic(ctx, resource)
			return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
		},
		nil,
		expectedMeta,
		suite.verifyRegionalWorkspaceResourceMetadataStep,
		expectedSpec,
		suite.verifyNicSpecStep,
		expectedState,
	)
}

func (suite *testSuite) getNicV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.NicSpec,
	expectedState schema.ResourceState,
) *schema.Nic {
	var resp *schema.Nic

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNic", string(wref.Workspace))
		retry := newStepResourceStateRetry(
			suite.baseDelay,
			suite.baseInterval,
			suite.maxAttempts,
			func() (schema.ResourceState, error) {
				var err error
				resp, err = api.GetNic(ctx, wref)
				requireNoError(sCtx, err)
				requireNotNilResponse(sCtx, resp)

				suite.requireNotNilStatus(sCtx, resp.Status)
				return *resp.Status.State, nil
			},
			func() {
				expectedMeta.Verb = http.MethodGet
				suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

				suite.verifyNicSpecStep(sCtx, expectedSpec, &resp.Spec)

				suite.verifyStatusStep(sCtx, expectedState, *resp.Status.State)
			},
		)
		retry.run(sCtx, "GetNic", expectedState)
	})
	return resp
}

func (suite *testSuite) getNicWithErrorV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	expectedError error,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNic", string(wref.Workspace))

		_, err := api.GetNic(ctx, wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteNicV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.NetworkV1, resource *schema.Nic) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteNic", resource.Metadata.Workspace)

		err := api.DeleteNic(ctx, resource)
		requireNoError(sCtx, err)
	})
}

// Security Group

func (suite *testSuite) createOrUpdateSecurityGroupV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	resource *schema.SecurityGroup,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.SecurityGroupSpec,
	expectedState schema.ResourceState,
) {
	expectedMeta.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(
		t,
		ctx,
		suite,
		stepName,
		suite.setNetworkV1StepParams,
		"CreateOrUpdateSecurityGroup",
		resource.Metadata.Workspace,
		resource,
		func(context.Context, *schema.SecurityGroup) (*stepFuncResponse[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec], error) {
			resp, err := api.CreateOrUpdateSecurityGroup(ctx, resource)
			return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
		},
		nil,
		expectedMeta,
		suite.verifyRegionalWorkspaceResourceMetadataStep,
		expectedSpec,
		suite.verifySecurityGroupSpecStep,
		expectedState,
	)
}

func (suite *testSuite) getSecurityGroupV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.SecurityGroupSpec,
	expectedState schema.ResourceState,
) *schema.SecurityGroup {
	var resp *schema.SecurityGroup

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSecurityGroup", string(wref.Workspace))
		retry := newStepResourceStateRetry(
			suite.baseDelay,
			suite.baseInterval,
			suite.maxAttempts,
			func() (schema.ResourceState, error) {
				var err error
				resp, err = api.GetSecurityGroup(ctx, wref)
				requireNoError(sCtx, err)
				requireNotNilResponse(sCtx, resp)

				suite.requireNotNilStatus(sCtx, resp.Status)
				return *resp.Status.State, nil
			},
			func() {
				expectedMeta.Verb = http.MethodGet
				suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

				suite.verifySecurityGroupSpecStep(sCtx, expectedSpec, &resp.Spec)

				suite.verifyStatusStep(sCtx, expectedState, *resp.Status.State)
			},
		)
		retry.run(sCtx, "GetSecurityGroup", expectedState)
	})
	return resp
}

func (suite *testSuite) getSecurityGroupWithErrorV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	expectedError error,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSecurityGroup", string(wref.Workspace))

		_, err := api.GetSecurityGroup(ctx, wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteSecurityGroupV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.NetworkV1, resource *schema.SecurityGroup) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteSecurityGroup", resource.Metadata.Workspace)

		err := api.DeleteSecurityGroup(ctx, resource)
		requireNoError(sCtx, err)
	})
}
