package secatest

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Network

func (suite *testSuite) createOrUpdateNetworkV1Step(stepName string, t provider.T, api *secapi.NetworkV1, resource *schema.Network,
	responseExpects responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec],
) {
	responseExpects.metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(t, suite,
		createOrUpdateWorkspaceResourceParams[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setNetworkV1StepParams,
			operationName:  "CreateOrUpdateNetwork",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Network) (*stepFuncResponse[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec], error) {
				resp, err := api.CreateOrUpdateNetwork(t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyNetworkSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getNetworkV1Step(stepName string, t provider.T, api *secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec],
) *schema.Network {
	responseExpects.metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(t, suite,
		getWorkspaceResourceParams[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setNetworkV1StepParams,
			operationName:  "GetNetwork",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec], error) {
				resp, err := api.GetNetworkUntilState(ctx, wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyNetworkSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getListNetworkV1Step(
	stepName string,
	t provider.T,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListNetwork", string(wref.Workspace))
		var iter *secapi.Iterator[schema.Network]
		var err error
		if opts != nil {
			iter, err = api.ListNetworksWithFilters(t.Context(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListNetworks(t.Context(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, t, *iter)
	})
}

func (suite *testSuite) getNetworkWithErrorV1Step(stepName string, t provider.T, api *secapi.NetworkV1, wref secapi.WorkspaceReference, expectedError error) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNetwork", string(wref.Workspace))

		_, err := api.GetNetwork(t.Context(), wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteNetworkV1Step(stepName string, t provider.T, api *secapi.NetworkV1, resource *schema.Network) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteNetwork", resource.Metadata.Workspace)

		err := api.DeleteNetwork(t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Internet Gateway

func (suite *testSuite) createOrUpdateInternetGatewayV1Step(stepName string, t provider.T, api *secapi.NetworkV1, resource *schema.InternetGateway,
	responseExpects responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec],
) {
	responseExpects.metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(t, suite,
		createOrUpdateWorkspaceResourceParams[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setNetworkV1StepParams,
			operationName:  "CreateOrUpdateInternetGateway",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.InternetGateway) (*stepFuncResponse[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec], error) {
				resp, err := api.CreateOrUpdateInternetGateway(t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyInternetGatewaySpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getInternetGatewayV1Step(stepName string, t provider.T, api *secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec],
) *schema.InternetGateway {
	responseExpects.metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(t, suite,
		getWorkspaceResourceParams[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setNetworkV1StepParams,
			operationName:  "GetInternetGateway",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec], error) {
				resp, err := api.GetInternetGatewayUntilState(ctx, wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyInternetGatewaySpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getListInternetGatewayV1Step(
	stepName string,
	t provider.T,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListInternetGateway", wref.Name)
		var iter *secapi.Iterator[schema.InternetGateway]
		var err error
		if opts != nil {
			iter, err = api.ListInternetGatewaysWithFilters(t.Context(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListInternetGateways(t.Context(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, t, *iter)
	})
}

func (suite *testSuite) getInternetGatewayWithErrorV1Step(stepName string, t provider.T, api *secapi.NetworkV1, wref secapi.WorkspaceReference, expectedError error) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetInternetGateway", string(wref.Workspace))

		_, err := api.GetInternetGateway(t.Context(), wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteInternetGatewayV1Step(stepName string, t provider.T, api *secapi.NetworkV1, resource *schema.InternetGateway) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteInternetGateway", resource.Metadata.Workspace)

		err := api.DeleteInternetGateway(t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Route Table

func (suite *testSuite) createOrUpdateRouteTableV1Step(stepName string, t provider.T, api *secapi.NetworkV1, resource *schema.RouteTable,
	responseExpects responseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec],
) {
	responseExpects.metadata.Verb = http.MethodPut
	createOrUpdateNetworkResourceStep(t, suite,
		createOrUpdateNetworkResourceParams[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setNetworkNetworkV1StepParams,
			operationName:  "CreateOrUpdateRouteTable",
			workspace:      resource.Metadata.Workspace,
			network:        resource.Metadata.Network,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.RouteTable) (*stepFuncResponse[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec], error) {
				resp, err := api.CreateOrUpdateRouteTable(t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalNetworkResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyRouteTableSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getRouteTableV1Step(stepName string, t provider.T, api *secapi.NetworkV1, nref secapi.NetworkReference,
	responseExpects responseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec],
) *schema.RouteTable {
	responseExpects.metadata.Verb = http.MethodGet
	return getNetworkResourceStep(t, suite,
		getNetworkResourceParams[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setNetworkNetworkV1StepParams,
			operationName:  "GetRouteTable",
			nref:           nref,
			getFunc: func(ctx context.Context, nref secapi.NetworkReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec], error) {
				resp, err := api.GetRouteTableUntilState(ctx, nref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalNetworkResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyRouteTableSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getListRouteTableV1Step(
	stepName string,
	t provider.T,
	api *secapi.NetworkV1,
	nref secapi.NetworkReference,
	opts *secapi.ListOptions,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListRouteTable", nref.Name)
		var iter *secapi.Iterator[schema.RouteTable]
		var err error
		if opts != nil {
			iter, err = api.ListRouteTablesWithFilters(t.Context(), nref.Tenant, nref.Workspace, nref.Network, opts)
		} else {
			iter, err = api.ListRouteTables(t.Context(), nref.Tenant, nref.Workspace, nref.Network)
		}
		requireNoError(sCtx, err)
		verifyIterListStep(sCtx, t, *iter)
	})
}

func (suite *testSuite) getRouteTableWithErrorV1Step(stepName string, t provider.T, api *secapi.NetworkV1, nref secapi.NetworkReference, expectedError error) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetRouteTable", string(nref.Workspace))

		_, err := api.GetRouteTable(t.Context(), nref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteRouteTableV1Step(stepName string, t provider.T, api *secapi.NetworkV1, resource *schema.RouteTable) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteRouteTable", resource.Metadata.Workspace)

		err := api.DeleteRouteTable(t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Subnet

func (suite *testSuite) createOrUpdateSubnetV1Step(stepName string, t provider.T, api *secapi.NetworkV1, resource *schema.Subnet,
	responseExpects responseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec],
) {
	responseExpects.metadata.Verb = http.MethodPut
	createOrUpdateNetworkResourceStep(t, suite,
		createOrUpdateNetworkResourceParams[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setNetworkNetworkV1StepParams,
			operationName:  "CreateOrUpdateSubnet",
			workspace:      resource.Metadata.Workspace,
			network:        resource.Metadata.Network,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Subnet) (*stepFuncResponse[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec], error) {
				resp, err := api.CreateOrUpdateSubnet(t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalNetworkResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifySubnetSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getSubnetV1Step(stepName string, t provider.T, api *secapi.NetworkV1, nref secapi.NetworkReference,
	responseExpects responseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec],
) *schema.Subnet {
	responseExpects.metadata.Verb = http.MethodGet
	return getNetworkResourceStep(t, suite,
		getNetworkResourceParams[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setNetworkNetworkV1StepParams,
			operationName:  "GetSubnet",
			nref:           nref,
			getFunc: func(ctx context.Context, nref secapi.NetworkReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec], error) {
				resp, err := api.GetSubnetUntilState(ctx, nref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalNetworkResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifySubnetSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getListSubnetV1Step(
	stepName string,
	t provider.T,
	api *secapi.NetworkV1,
	nref secapi.NetworkReference,
	opts *secapi.ListOptions,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListSubnet", nref.Name)
		var iter *secapi.Iterator[schema.Subnet]
		var err error
		if opts != nil {
			iter, err = api.ListSubnetsWithFilters(t.Context(), nref.Tenant, nref.Workspace, nref.Network, opts)
		} else {
			iter, err = api.ListSubnets(t.Context(), nref.Tenant, nref.Workspace, nref.Network)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, t, *iter)
	})
}

func (suite *testSuite) getSubnetWithErrorV1Step(stepName string, t provider.T, api *secapi.NetworkV1, nref secapi.NetworkReference, expectedError error) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSubnet", string(nref.Workspace))

		_, err := api.GetSubnet(t.Context(), nref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteSubnetV1Step(stepName string, t provider.T, api *secapi.NetworkV1, resource *schema.Subnet) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteSubnet", resource.Metadata.Workspace)

		err := api.DeleteSubnet(t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Public Ip

func (suite *testSuite) createOrUpdatePublicIpV1Step(stepName string, t provider.T, api *secapi.NetworkV1, resource *schema.PublicIp,
	responseExpects responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec],
) {
	responseExpects.metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(t, suite,
		createOrUpdateWorkspaceResourceParams[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setNetworkV1StepParams,
			operationName:  "CreateOrUpdatePublicIp",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.PublicIp) (*stepFuncResponse[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec], error) {
				resp, err := api.CreateOrUpdatePublicIp(t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyPublicIpSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getPublicIpV1Step(stepName string, t provider.T, api *secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec],
) *schema.PublicIp {
	responseExpects.metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(t, suite,
		getWorkspaceResourceParams[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setNetworkV1StepParams,
			operationName:  "GetPublicIp",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec], error) {
				resp, err := api.GetPublicIpUntilState(ctx, wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyPublicIpSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getListPublicIpV1Step(
	stepName string,
	t provider.T,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListPublicIp", wref.Name)
		var iter *secapi.Iterator[schema.PublicIp]
		var err error
		if opts != nil {
			iter, err = api.ListPublicIpsWithFilters(t.Context(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListPublicIps(t.Context(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)

		// Iterate through all items
		verifyIterListStep(sCtx, t, *iter)
	})
}

func (suite *testSuite) getPublicIpWithErrorV1Step(stepName string, t provider.T, api *secapi.NetworkV1, wref secapi.WorkspaceReference, expectedError error) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetPublicIp", string(wref.Workspace))

		_, err := api.GetPublicIp(t.Context(), wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deletePublicIpV1Step(stepName string, t provider.T, api *secapi.NetworkV1, resource *schema.PublicIp) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeletePublicIp", resource.Metadata.Workspace)

		err := api.DeletePublicIp(t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Nic

func (suite *testSuite) createOrUpdateNicV1Step(stepName string, t provider.T, api *secapi.NetworkV1, resource *schema.Nic,
	responseExpects responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec],
) {
	responseExpects.metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(t, suite,
		createOrUpdateWorkspaceResourceParams[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setNetworkV1StepParams,
			operationName:  "CreateOrUpdateNic",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Nic) (*stepFuncResponse[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec], error) {
				resp, err := api.CreateOrUpdateNic(t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyNicSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getNicV1Step(stepName string, t provider.T, api *secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec],
) *schema.Nic {
	responseExpects.metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(t, suite,
		getWorkspaceResourceParams[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setNetworkV1StepParams,
			operationName:  "GetNic",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec], error) {
				resp, err := api.GetNicUntilState(ctx, wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyNicSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getListNicV1Step(
	stepName string,
	t provider.T,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListNic", wref.Name)
		var iter *secapi.Iterator[schema.Nic]
		var err error
		if opts != nil {
			iter, err = api.ListNicsWithFilters(t.Context(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListNics(t.Context(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, t, *iter)
	})
}

func (suite *testSuite) getNicWithErrorV1Step(stepName string, t provider.T, api *secapi.NetworkV1, wref secapi.WorkspaceReference, expectedError error) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetNic", string(wref.Workspace))

		_, err := api.GetNic(t.Context(), wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteNicV1Step(stepName string, t provider.T, api *secapi.NetworkV1, resource *schema.Nic) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteNic", resource.Metadata.Workspace)

		err := api.DeleteNic(t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Security Group

func (suite *testSuite) createOrUpdateSecurityGroupV1Step(stepName string, t provider.T, api *secapi.NetworkV1, resource *schema.SecurityGroup,
	responseExpects responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec],
) {
	responseExpects.metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(t, suite,
		createOrUpdateWorkspaceResourceParams[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setNetworkV1StepParams,
			operationName:  "CreateOrUpdateSecurityGroup",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.SecurityGroup) (*stepFuncResponse[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec], error) {
				resp, err := api.CreateOrUpdateSecurityGroup(t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifySecurityGroupSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getSecurityGroupV1Step(stepName string, t provider.T, api *secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec],
) *schema.SecurityGroup {
	responseExpects.metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(t, suite,
		getWorkspaceResourceParams[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setNetworkV1StepParams,
			operationName:  "GetSecurityGroup",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec], error) {
				resp, err := api.GetSecurityGroupUntilState(ctx, wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifySecurityGroupSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getListSecurityGroupV1Step(
	stepName string,
	t provider.T,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListSecurityGroup", wref.Name)
		var iter *secapi.Iterator[schema.SecurityGroup]
		var err error
		if opts != nil {
			iter, err = api.ListSecurityGroupsWithFilters(t.Context(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListSecurityGroups(t.Context(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, t, *iter)
	})
}

func (suite *testSuite) getSecurityGroupWithErrorV1Step(stepName string, t provider.T, api *secapi.NetworkV1, wref secapi.WorkspaceReference, expectedError error) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "GetSecurityGroup", string(wref.Workspace))

		_, err := api.GetSecurityGroup(t.Context(), wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteSecurityGroupV1Step(stepName string, t provider.T, api *secapi.NetworkV1, resource *schema.SecurityGroup) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setNetworkV1StepParams(sCtx, "DeleteSecurityGroup", resource.Metadata.Workspace)

		err := api.DeleteSecurityGroup(t.Context(), resource)
		requireNoError(sCtx, err)
	})
}
