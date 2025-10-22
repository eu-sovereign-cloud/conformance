package secatest

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
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
	instance *schema.Network,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.NetworkSpec,
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateNetwork")

		resp, err := api.CreateOrUpdateNetwork(ctx, instance)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodPut
			suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifyNetworkSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
}

func (suite *testSuite) getNetworkV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.NetworkSpec,
	expectedStatusState string,
) *schema.Network {
	var resp *schema.Network
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetNetwork")

		resp, err = api.GetNetwork(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodGet
			suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifyNetworkSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
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
		suite.setAuthorizationV1StepParams(sCtx, "GetNetwork")

		_, err := api.GetNetwork(ctx, wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteNetworkV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.NetworkV1, instance *schema.Network) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteNetwork")

		err := api.DeleteNetwork(ctx, instance)
		requireNoError(sCtx, err)
	})
}

// Internet Gateway

func (suite *testSuite) createOrUpdateInternetGatewayV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	InternetGateway *schema.InternetGateway,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.InternetGatewaySpec,
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateInternetGateway")

		resp, err := api.CreateOrUpdateInternetGateway(ctx, InternetGateway)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodPut
			suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifyInternetGatewaySpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
}

func (suite *testSuite) getInternetGatewayV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.InternetGatewaySpec,
	expectedStatusState string,
) *schema.InternetGateway {
	var resp *schema.InternetGateway
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetInternetGateway")

		resp, err = api.GetInternetGateway(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodGet
			suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifyInternetGatewaySpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
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
		suite.setAuthorizationV1StepParams(sCtx, "GetInternetGateway")

		_, err := api.GetInternetGateway(ctx, wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteInternetGatewayV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.NetworkV1, InternetGateway *schema.InternetGateway) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteInternetGateway")

		err := api.DeleteInternetGateway(ctx, InternetGateway)
		requireNoError(sCtx, err)
	})
}

// Route Table

func (suite *testSuite) createOrUpdateRouteTableV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	RouteTable *schema.RouteTable,
	expectedMeta *schema.RegionalNetworkResourceMetadata,
	expectedSpec *schema.RouteTableSpec,
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateRouteTable")

		resp, err := api.CreateOrUpdateRouteTable(ctx, RouteTable)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodPut
			suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifyRouteTableSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
}

func (suite *testSuite) getRouteTableV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	nref secapi.NetworkReference,
	expectedMeta *schema.RegionalNetworkResourceMetadata,
	expectedSpec *schema.RouteTableSpec,
	expectedStatusState string,
) *schema.RouteTable {
	var resp *schema.RouteTable
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRouteTable")

		resp, err = api.GetRouteTable(ctx, nref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodGet
			suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifyRouteTableSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
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
		suite.setAuthorizationV1StepParams(sCtx, "GetRouteTable")

		_, err := api.GetRouteTable(ctx, nref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteRouteTableV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.NetworkV1, RouteTable *schema.RouteTable) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteRouteTable")

		err := api.DeleteRouteTable(ctx, RouteTable)
		requireNoError(sCtx, err)
	})
}

// Subnet

func (suite *testSuite) createOrUpdateSubnetV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	Subnet *schema.Subnet,
	expectedMeta *schema.RegionalNetworkResourceMetadata,
	expectedSpec *schema.SubnetSpec,
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateSubnet")

		resp, err := api.CreateOrUpdateSubnet(ctx, Subnet)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodPut
			suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifySubnetSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
}

func (suite *testSuite) getSubnetV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	nref secapi.NetworkReference,
	expectedMeta *schema.RegionalNetworkResourceMetadata,
	expectedSpec *schema.SubnetSpec,
	expectedStatusState string,
) *schema.Subnet {
	var resp *schema.Subnet
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetSubnet")

		resp, err = api.GetSubnet(ctx, nref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodGet
			suite.verifyRegionalNetworkResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifySubnetSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
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
		suite.setAuthorizationV1StepParams(sCtx, "GetSubnet")

		_, err := api.GetSubnet(ctx, nref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteSubnetV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.NetworkV1, Subnet *schema.Subnet) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteSubnet")

		err := api.DeleteSubnet(ctx, Subnet)
		requireNoError(sCtx, err)
	})
}

// Public Ip

func (suite *testSuite) createOrUpdatePublicIpV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	PublicIp *schema.PublicIp,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.PublicIpSpec,
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdatePublicIp")

		resp, err := api.CreateOrUpdatePublicIp(ctx, PublicIp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodPut
			suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifyPublicIpSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
}

func (suite *testSuite) getPublicIpV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.PublicIpSpec,
	expectedStatusState string,
) *schema.PublicIp {
	var resp *schema.PublicIp
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetPublicIp")

		resp, err = api.GetPublicIp(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodGet
			suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifyPublicIpSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
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
		suite.setAuthorizationV1StepParams(sCtx, "GetPublicIp")

		_, err := api.GetPublicIp(ctx, wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deletePublicIpV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.NetworkV1, PublicIp *schema.PublicIp) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeletePublicIp")

		err := api.DeletePublicIp(ctx, PublicIp)
		requireNoError(sCtx, err)
	})
}

// Nic

func (suite *testSuite) createOrUpdateNicV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	Nic *schema.Nic,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.NicSpec,
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateNic")

		resp, err := api.CreateOrUpdateNic(ctx, Nic)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodPut
			suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifyNicSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
}

func (suite *testSuite) getNicV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.NicSpec,
	expectedStatusState string,
) *schema.Nic {
	var resp *schema.Nic
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetNic")

		resp, err = api.GetNic(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodGet
			suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifyNicSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
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
		suite.setAuthorizationV1StepParams(sCtx, "GetNic")

		_, err := api.GetNic(ctx, wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteNicV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.NetworkV1, Nic *schema.Nic) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteNic")

		err := api.DeleteNic(ctx, Nic)
		requireNoError(sCtx, err)
	})
}

// Security Group

func (suite *testSuite) createOrUpdateSecurityGroupV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	SecurityGroup *schema.SecurityGroup,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.SecurityGroupSpec,
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateSecurityGroup")

		resp, err := api.CreateOrUpdateSecurityGroup(ctx, SecurityGroup)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodPut
			suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifySecurityGroupSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
}

func (suite *testSuite) getSecurityGroupV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.SecurityGroupSpec,
	expectedStatusState string,
) *schema.SecurityGroup {
	var resp *schema.SecurityGroup
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetSecurityGroup")

		resp, err = api.GetSecurityGroup(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodGet
			suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifySecurityGroupSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
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
		suite.setAuthorizationV1StepParams(sCtx, "GetSecurityGroup")

		_, err := api.GetSecurityGroup(ctx, wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteSecurityGroupV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.NetworkV1, SecurityGroup *schema.SecurityGroup) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteSecurityGroup")

		err := api.DeleteSecurityGroup(ctx, SecurityGroup)
		requireNoError(sCtx, err)
	})
}
