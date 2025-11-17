package secatest

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

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
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNetwork", resource.Metadata.Workspace)

		resp, err := api.CreateOrUpdateNetwork(ctx, resource)
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
		suite.setNetworkV1StepParams(sCtx, "GetNetwork", string(wref.Workspace))

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

func (suite *testSuite) getListNetworkV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	tref secapi.TenantReference,
	wref secapi.WorkspaceReference,
	opts *builders.ListOptions,
) []*schema.Network {
	var respNext []*schema.Network
	var respAll []*schema.Network
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListNetwork", string(wref.Name))
		var iter *secapi.Iterator[schema.Network]
		var err error
		if opts != nil {
			iter, err = api.ListNetworksWithFilters(ctx, secapi.TenantID(tref.Name), secapi.WorkspaceID(wref.Name), opts)
		} else {
			iter, err = api.ListNetworks(ctx, secapi.TenantID(tref.Name), secapi.WorkspaceID(wref.Name))
		}
		requireNoError(sCtx, err)
		for {
			item, err := iter.Next(context.Background())
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			respNext = append(respNext, item)
		}
		requireNotNilResponse(sCtx, respNext)
		requireLenResponse(sCtx, len(respNext))

		respAll, err = iter.All(ctx)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, respAll)
		requireLenResponse(sCtx, len(respAll))

		compareIteratorsResponse(sCtx, len(respNext), len(respAll))
	})
	return respAll
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
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateInternetGateway", resource.Metadata.Workspace)

		resp, err := api.CreateOrUpdateInternetGateway(ctx, resource)
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
		suite.setNetworkV1StepParams(sCtx, "GetInternetGateway", string(wref.Workspace))

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

func (suite *testSuite) getListInternetGatewayV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	tref secapi.TenantReference,
	wref secapi.WorkspaceReference,
	opts *builders.ListOptions,
) []*schema.InternetGateway {
	var respNext []*schema.InternetGateway
	var respAll []*schema.InternetGateway
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListInternetGateway", string(wref.Name))
		var iter *secapi.Iterator[schema.InternetGateway]
		var err error
		if opts != nil {
			iter, err = api.ListInternetGatewaysWithFilters(ctx, secapi.TenantID(tref.Name), secapi.WorkspaceID(wref.Name), opts)
		} else {
			iter, err = api.ListInternetGateways(ctx, secapi.TenantID(tref.Name), secapi.WorkspaceID(wref.Name))
		}
		requireNoError(sCtx, err)
		for {
			item, err := iter.Next(context.Background())
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			respNext = append(respNext, item)
		}
		requireNotNilResponse(sCtx, respNext)
		requireLenResponse(sCtx, len(respNext))

		respAll, err = iter.All(ctx)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, respAll)
		requireLenResponse(sCtx, len(respAll))

		compareIteratorsResponse(sCtx, len(respNext), len(respAll))
	})
	return respAll
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
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateRouteTable", resource.Metadata.Workspace)

		resp, err := api.CreateOrUpdateRouteTable(ctx, resource)
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
		suite.setNetworkV1StepParams(sCtx, "GetRouteTable", string(nref.Workspace))

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

func (suite *testSuite) getListRouteTableV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	tref secapi.TenantReference,
	wref secapi.WorkspaceReference,
	nref secapi.NetworkReference,
	opts *builders.ListOptions,
) []*schema.RouteTable {
	var respNext []*schema.RouteTable
	var respAll []*schema.RouteTable
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListRouteTable", string(nref.Name))
		var iter *secapi.Iterator[schema.RouteTable]
		var err error
		if opts != nil {
			iter, err = api.ListRouteTablesWithFilters(ctx, secapi.TenantID(tref.Name), secapi.WorkspaceID(wref.Name), secapi.NetworkID(nref.Name), opts)
		} else {
			iter, err = api.ListRouteTables(ctx, secapi.TenantID(tref.Name), secapi.WorkspaceID(wref.Name), secapi.NetworkID(nref.Name))
		}
		requireNoError(sCtx, err)
		for {
			item, err := iter.Next(context.Background())
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			respNext = append(respNext, item)
		}
		requireNotNilResponse(sCtx, respNext)
		requireLenResponse(sCtx, len(respNext))

		respAll, err = iter.All(ctx)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, respAll)
		requireLenResponse(sCtx, len(respAll))

		compareIteratorsResponse(sCtx, len(respNext), len(respAll))
	})
	return respAll
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
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateSubnet", resource.Metadata.Workspace)

		resp, err := api.CreateOrUpdateSubnet(ctx, resource)
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
		suite.setNetworkV1StepParams(sCtx, "GetSubnet", string(nref.Workspace))

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

func (suite *testSuite) getListSubnetV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	tref secapi.TenantReference,
	wref secapi.WorkspaceReference,
	nref secapi.NetworkReference,
	opts *builders.ListOptions,
) []*schema.Subnet {
	var respNext []*schema.Subnet
	var respAll []*schema.Subnet
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListSubnet", string(nref.Name))
		var iter *secapi.Iterator[schema.Subnet]
		var err error
		if opts != nil {
			iter, err = api.ListSubnetsWithFilters(ctx, secapi.TenantID(tref.Name), secapi.WorkspaceID(wref.Name), secapi.NetworkID(nref.Name), opts)
		} else {
			iter, err = api.ListSubnets(ctx, secapi.TenantID(tref.Name), secapi.WorkspaceID(wref.Name), secapi.NetworkID(nref.Name))
		}
		requireNoError(sCtx, err)
		for {
			item, err := iter.Next(context.Background())
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			respNext = append(respNext, item)
		}
		requireNotNilResponse(sCtx, respNext)
		requireLenResponse(sCtx, len(respNext))

		respAll, err = iter.All(ctx)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, respAll)
		requireLenResponse(sCtx, len(respAll))

		compareIteratorsResponse(sCtx, len(respNext), len(respAll))
	})
	return respAll
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
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdatePublicIp", resource.Metadata.Workspace)

		resp, err := api.CreateOrUpdatePublicIp(ctx, resource)
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
		suite.setNetworkV1StepParams(sCtx, "GetPublicIp", string(wref.Workspace))

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

func (suite *testSuite) getListPublicIpV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	tref secapi.TenantReference,
	wref secapi.WorkspaceReference,
	opts *builders.ListOptions,
) []*schema.PublicIp {
	var respNext []*schema.PublicIp
	var respAll []*schema.PublicIp
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListPublicIp", string(wref.Name))
		var iter *secapi.Iterator[schema.PublicIp]
		var err error
		if opts != nil {
			iter, err = api.ListPublicIpsWithFilters(ctx, secapi.TenantID(tref.Name), secapi.WorkspaceID(wref.Name), opts)
		} else {
			iter, err = api.ListPublicIps(ctx, secapi.TenantID(tref.Name), secapi.WorkspaceID(wref.Name))
		}
		requireNoError(sCtx, err)
		for {
			item, err := iter.Next(context.Background())
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			respNext = append(respNext, item)
		}
		requireNotNilResponse(sCtx, respNext)
		requireLenResponse(sCtx, len(respNext))

		respAll, err = iter.All(ctx)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, respAll)
		requireLenResponse(sCtx, len(respAll))

		compareIteratorsResponse(sCtx, len(respNext), len(respAll))
	})
	return respAll
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
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateNic", resource.Metadata.Workspace)

		resp, err := api.CreateOrUpdateNic(ctx, resource)
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
		suite.setNetworkV1StepParams(sCtx, "GetNic", string(wref.Workspace))

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

func (suite *testSuite) getListNicV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	tref secapi.TenantReference,
	wref secapi.WorkspaceReference,
	opts *builders.ListOptions,
) []*schema.Nic {
	var respNext []*schema.Nic
	var respAll []*schema.Nic
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListNic", string(wref.Name))
		var iter *secapi.Iterator[schema.Nic]
		var err error
		if opts != nil {
			iter, err = api.ListNicsWithFilters(ctx, secapi.TenantID(tref.Name), secapi.WorkspaceID(wref.Name), opts)
		} else {
			iter, err = api.ListNics(ctx, secapi.TenantID(tref.Name), secapi.WorkspaceID(wref.Name))
		}
		requireNoError(sCtx, err)
		for {
			item, err := iter.Next(context.Background())
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			respNext = append(respNext, item)
		}
		requireNotNilResponse(sCtx, respNext)
		requireLenResponse(sCtx, len(respNext))

		respAll, err = iter.All(ctx)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, respAll)
		requireLenResponse(sCtx, len(respAll))

		compareIteratorsResponse(sCtx, len(respNext), len(respAll))
	})
	return respAll
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
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "CreateOrUpdateSecurityGroup", resource.Metadata.Workspace)

		resp, err := api.CreateOrUpdateSecurityGroup(ctx, resource)
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
		suite.setNetworkV1StepParams(sCtx, "GetSecurityGroup", string(wref.Workspace))

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

func (suite *testSuite) getListSecurityGroupV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.NetworkV1,
	tref secapi.TenantReference,
	wref secapi.WorkspaceReference,
	opts *builders.ListOptions,
) []*schema.SecurityGroup {
	var respNext []*schema.SecurityGroup
	var respAll []*schema.SecurityGroup
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListSecurityGroup", string(wref.Name))
		var iter *secapi.Iterator[schema.SecurityGroup]
		var err error
		if opts != nil {
			iter, err = api.ListSecurityGroupsWithFilters(ctx, secapi.TenantID(tref.Name), secapi.WorkspaceID(wref.Name), opts)
		} else {
			iter, err = api.ListSecurityGroups(ctx, secapi.TenantID(tref.Name), secapi.WorkspaceID(wref.Name))
		}
		requireNoError(sCtx, err)
		for {
			item, err := iter.Next(context.Background())
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			respNext = append(respNext, item)
		}
		requireNotNilResponse(sCtx, respNext)
		requireLenResponse(sCtx, len(respNext))

		respAll, err = iter.All(ctx)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, respAll)
		requireLenResponse(sCtx, len(respAll))

		compareIteratorsResponse(sCtx, len(respNext), len(respAll))
	})
	return respAll
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
