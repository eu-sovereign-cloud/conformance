//nolint:dupl
package steps

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Network

func (builder *Builder) CreateOrUpdateNetworkV1Step(stepName string, api *secapi.NetworkV1, resource *schema.Network,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(builder.t, builder.suite,
		createOrUpdateWorkspaceResourceParams[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetNetworkV1StepParams,
			operationName:  "CreateOrUpdateNetwork",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Network) (*stepFuncResponse[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec], error) {
				resp, err := api.CreateOrUpdateNetwork(builder.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyNetworkSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetNetworkV1Step(stepName string, api *secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec],
) *schema.Network {
	responseExpects.Metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(builder.t, builder.suite,
		getWorkspaceResourceParams[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetNetworkV1StepParams,
			operationName:  "GetNetwork",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec], error) {
				resp, err := api.GetNetworkUntilState(ctx, wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyNetworkSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetListNetworkV1Step(
	stepName string,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListNetwork", string(wref.Workspace))
		var iter *secapi.Iterator[schema.Network]
		var err error
		if opts != nil {
			iter, err = api.ListNetworksWithFilters(builder.t.Context(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListNetworks(builder.t.Context(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, builder.t, *iter)
	})
}

func (builder *Builder) GetNetworkWithErrorV1Step(stepName string, api *secapi.NetworkV1, wref secapi.WorkspaceReference, expectedError error) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetNetworkV1StepParams(sCtx, "GetNetwork", string(wref.Workspace))

		_, err := api.GetNetwork(builder.t.Context(), wref)
		requireError(sCtx, err, expectedError)
	})
}

func (builder *Builder) DeleteNetworkV1Step(stepName string, api *secapi.NetworkV1, resource *schema.Network) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetNetworkV1StepParams(sCtx, "DeleteNetwork", resource.Metadata.Workspace)

		err := api.DeleteNetwork(builder.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Internet Gateway

func (builder *Builder) CreateOrUpdateInternetGatewayV1Step(stepName string, api *secapi.NetworkV1, resource *schema.InternetGateway,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(builder.t, builder.suite,
		createOrUpdateWorkspaceResourceParams[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetNetworkV1StepParams,
			operationName:  "CreateOrUpdateInternetGateway",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.InternetGateway) (*stepFuncResponse[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec], error) {
				resp, err := api.CreateOrUpdateInternetGateway(builder.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyInternetGatewaySpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetInternetGatewayV1Step(stepName string, api *secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec],
) *schema.InternetGateway {
	responseExpects.Metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(builder.t, builder.suite,
		getWorkspaceResourceParams[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetNetworkV1StepParams,
			operationName:  "GetInternetGateway",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec], error) {
				resp, err := api.GetInternetGatewayUntilState(ctx, wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyInternetGatewaySpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetListInternetGatewayV1Step(
	stepName string,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListInternetGateway", wref.Name)
		var iter *secapi.Iterator[schema.InternetGateway]
		var err error
		if opts != nil {
			iter, err = api.ListInternetGatewaysWithFilters(builder.t.Context(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListInternetGateways(builder.t.Context(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, builder.t, *iter)
	})
}

func (builder *Builder) GetInternetGatewayWithErrorV1Step(stepName string, api *secapi.NetworkV1, wref secapi.WorkspaceReference, expectedError error) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetNetworkV1StepParams(sCtx, "GetInternetGateway", string(wref.Workspace))

		_, err := api.GetInternetGateway(builder.t.Context(), wref)
		requireError(sCtx, err, expectedError)
	})
}

func (builder *Builder) DeleteInternetGatewayV1Step(stepName string, api *secapi.NetworkV1, resource *schema.InternetGateway) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetNetworkV1StepParams(sCtx, "DeleteInternetGateway", resource.Metadata.Workspace)

		err := api.DeleteInternetGateway(builder.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Route Table

func (builder *Builder) CreateOrUpdateRouteTableV1Step(stepName string, api *secapi.NetworkV1, resource *schema.RouteTable,
	responseExpects ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateNetworkResourceStep(builder.t, builder.suite,
		createOrUpdateNetworkResourceParams[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetNetworkNetworkV1StepParams,
			operationName:  "CreateOrUpdateRouteTable",
			workspace:      resource.Metadata.Workspace,
			network:        resource.Metadata.Network,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.RouteTable) (*stepFuncResponse[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec], error) {
				resp, err := api.CreateOrUpdateRouteTable(builder.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalNetworkResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyRouteTableSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetRouteTableV1Step(stepName string, api *secapi.NetworkV1, nref secapi.NetworkReference,
	responseExpects ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec],
) *schema.RouteTable {
	responseExpects.Metadata.Verb = http.MethodGet
	return getNetworkResourceStep(builder.t, builder.suite,
		getNetworkResourceParams[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetNetworkNetworkV1StepParams,
			operationName:  "GetRouteTable",
			nref:           nref,
			getFunc: func(ctx context.Context, nref secapi.NetworkReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec], error) {
				resp, err := api.GetRouteTableUntilState(ctx, nref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalNetworkResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyRouteTableSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetListRouteTableV1Step(
	stepName string,
	api *secapi.NetworkV1,
	nref secapi.NetworkReference,
	opts *secapi.ListOptions,
) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListRouteTable", nref.Name)
		var iter *secapi.Iterator[schema.RouteTable]
		var err error
		if opts != nil {
			iter, err = api.ListRouteTablesWithFilters(builder.t.Context(), nref.Tenant, nref.Workspace, nref.Network, opts)
		} else {
			iter, err = api.ListRouteTables(builder.t.Context(), nref.Tenant, nref.Workspace, nref.Network)
		}
		requireNoError(sCtx, err)
		verifyIterListStep(sCtx, builder.t, *iter)
	})
}

func (builder *Builder) GetRouteTableWithErrorV1Step(stepName string, api *secapi.NetworkV1, nref secapi.NetworkReference, expectedError error) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetNetworkV1StepParams(sCtx, "GetRouteTable", string(nref.Workspace))

		_, err := api.GetRouteTable(builder.t.Context(), nref)
		requireError(sCtx, err, expectedError)
	})
}

func (builder *Builder) DeleteRouteTableV1Step(stepName string, api *secapi.NetworkV1, resource *schema.RouteTable) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetNetworkV1StepParams(sCtx, "DeleteRouteTable", resource.Metadata.Workspace)

		err := api.DeleteRouteTable(builder.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Subnet

func (builder *Builder) CreateOrUpdateSubnetV1Step(stepName string, api *secapi.NetworkV1, resource *schema.Subnet,
	responseExpects ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateNetworkResourceStep(builder.t, builder.suite,
		createOrUpdateNetworkResourceParams[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetNetworkNetworkV1StepParams,
			operationName:  "CreateOrUpdateSubnet",
			workspace:      resource.Metadata.Workspace,
			network:        resource.Metadata.Network,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Subnet) (*stepFuncResponse[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec], error) {
				resp, err := api.CreateOrUpdateSubnet(builder.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalNetworkResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifySubnetSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetSubnetV1Step(stepName string, api *secapi.NetworkV1, nref secapi.NetworkReference,
	responseExpects ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec],
) *schema.Subnet {
	responseExpects.Metadata.Verb = http.MethodGet
	return getNetworkResourceStep(builder.t, builder.suite,
		getNetworkResourceParams[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetNetworkNetworkV1StepParams,
			operationName:  "GetSubnet",
			nref:           nref,
			getFunc: func(ctx context.Context, nref secapi.NetworkReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec], error) {
				resp, err := api.GetSubnetUntilState(ctx, nref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalNetworkResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifySubnetSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetListSubnetV1Step(
	stepName string,
	api *secapi.NetworkV1,
	nref secapi.NetworkReference,
	opts *secapi.ListOptions,
) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListSubnet", nref.Name)
		var iter *secapi.Iterator[schema.Subnet]
		var err error
		if opts != nil {
			iter, err = api.ListSubnetsWithFilters(builder.t.Context(), nref.Tenant, nref.Workspace, nref.Network, opts)
		} else {
			iter, err = api.ListSubnets(builder.t.Context(), nref.Tenant, nref.Workspace, nref.Network)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, builder.t, *iter)
	})
}

func (builder *Builder) GetSubnetWithErrorV1Step(stepName string, api *secapi.NetworkV1, nref secapi.NetworkReference, expectedError error) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetNetworkV1StepParams(sCtx, "GetSubnet", string(nref.Workspace))

		_, err := api.GetSubnet(builder.t.Context(), nref)
		requireError(sCtx, err, expectedError)
	})
}

func (builder *Builder) DeleteSubnetV1Step(stepName string, api *secapi.NetworkV1, resource *schema.Subnet) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetNetworkV1StepParams(sCtx, "DeleteSubnet", resource.Metadata.Workspace)

		err := api.DeleteSubnet(builder.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Public Ip

func (builder *Builder) CreateOrUpdatePublicIpV1Step(stepName string, api *secapi.NetworkV1, resource *schema.PublicIp,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(builder.t, builder.suite,
		createOrUpdateWorkspaceResourceParams[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetNetworkV1StepParams,
			operationName:  "CreateOrUpdatePublicIp",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.PublicIp) (*stepFuncResponse[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec], error) {
				resp, err := api.CreateOrUpdatePublicIp(builder.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyPublicIpSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetPublicIpV1Step(stepName string, api *secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec],
) *schema.PublicIp {
	responseExpects.Metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(builder.t, builder.suite,
		getWorkspaceResourceParams[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetNetworkV1StepParams,
			operationName:  "GetPublicIp",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec], error) {
				resp, err := api.GetPublicIpUntilState(ctx, wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyPublicIpSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetListPublicIpV1Step(
	stepName string,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListPublicIp", wref.Name)
		var iter *secapi.Iterator[schema.PublicIp]
		var err error
		if opts != nil {
			iter, err = api.ListPublicIpsWithFilters(builder.t.Context(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListPublicIps(builder.t.Context(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)

		// Iterate through all items
		verifyIterListStep(sCtx, builder.t, *iter)
	})
}

func (builder *Builder) GetPublicIpWithErrorV1Step(stepName string, api *secapi.NetworkV1, wref secapi.WorkspaceReference, expectedError error) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetNetworkV1StepParams(sCtx, "GetPublicIp", string(wref.Workspace))

		_, err := api.GetPublicIp(builder.t.Context(), wref)
		requireError(sCtx, err, expectedError)
	})
}

func (builder *Builder) DeletePublicIpV1Step(stepName string, api *secapi.NetworkV1, resource *schema.PublicIp) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetNetworkV1StepParams(sCtx, "DeletePublicIp", resource.Metadata.Workspace)

		err := api.DeletePublicIp(builder.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Nic

func (builder *Builder) CreateOrUpdateNicV1Step(stepName string, api *secapi.NetworkV1, resource *schema.Nic,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(builder.t, builder.suite,
		createOrUpdateWorkspaceResourceParams[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetNetworkV1StepParams,
			operationName:  "CreateOrUpdateNic",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Nic) (*stepFuncResponse[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec], error) {
				resp, err := api.CreateOrUpdateNic(builder.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyNicSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetNicV1Step(stepName string, api *secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec],
) *schema.Nic {
	responseExpects.Metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(builder.t, builder.suite,
		getWorkspaceResourceParams[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetNetworkV1StepParams,
			operationName:  "GetNic",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec], error) {
				resp, err := api.GetNicUntilState(ctx, wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyNicSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetListNicV1Step(
	stepName string,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListNic", wref.Name)
		var iter *secapi.Iterator[schema.Nic]
		var err error
		if opts != nil {
			iter, err = api.ListNicsWithFilters(builder.t.Context(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListNics(builder.t.Context(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, builder.t, *iter)
	})
}

func (builder *Builder) GetNicWithErrorV1Step(stepName string, api *secapi.NetworkV1, wref secapi.WorkspaceReference, expectedError error) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetNetworkV1StepParams(sCtx, "GetNic", string(wref.Workspace))

		_, err := api.GetNic(builder.t.Context(), wref)
		requireError(sCtx, err, expectedError)
	})
}

func (builder *Builder) DeleteNicV1Step(stepName string, api *secapi.NetworkV1, resource *schema.Nic) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetNetworkV1StepParams(sCtx, "DeleteNic", resource.Metadata.Workspace)

		err := api.DeleteNic(builder.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Security Group

func (builder *Builder) CreateOrUpdateSecurityGroupV1Step(stepName string, api *secapi.NetworkV1, resource *schema.SecurityGroup,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(builder.t, builder.suite,
		createOrUpdateWorkspaceResourceParams[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetNetworkV1StepParams,
			operationName:  "CreateOrUpdateSecurityGroup",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.SecurityGroup) (*stepFuncResponse[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec], error) {
				resp, err := api.CreateOrUpdateSecurityGroup(builder.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifySecurityGroupSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetSecurityGroupV1Step(stepName string, api *secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec],
) *schema.SecurityGroup {
	responseExpects.Metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(builder.t, builder.suite,
		getWorkspaceResourceParams[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetNetworkV1StepParams,
			operationName:  "GetSecurityGroup",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec], error) {
				resp, err := api.GetSecurityGroupUntilState(ctx, wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifySecurityGroupSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetListSecurityGroupV1Step(
	stepName string,
	api *secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListSecurityGroup", wref.Name)
		var iter *secapi.Iterator[schema.SecurityGroup]
		var err error
		if opts != nil {
			iter, err = api.ListSecurityGroupsWithFilters(builder.t.Context(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListSecurityGroups(builder.t.Context(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, builder.t, *iter)
	})
}

func (builder *Builder) GetSecurityGroupWithErrorV1Step(stepName string, api *secapi.NetworkV1, wref secapi.WorkspaceReference, expectedError error) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetNetworkV1StepParams(sCtx, "GetSecurityGroup", string(wref.Workspace))

		_, err := api.GetSecurityGroup(builder.t.Context(), wref)
		requireError(sCtx, err, expectedError)
	})
}

func (builder *Builder) DeleteSecurityGroupV1Step(stepName string, api *secapi.NetworkV1, resource *schema.SecurityGroup) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetNetworkV1StepParams(sCtx, "DeleteSecurityGroup", resource.Metadata.Workspace)

		err := api.DeleteSecurityGroup(builder.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (builder *Builder) GetListNetworkSkusV1Step(
	stepName string,
	api *secapi.NetworkV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) []*schema.NetworkSku {
	var resp []*schema.NetworkSku

	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetComputeV1StepParams(sCtx, "ListSkus", tref.Name)

		var iter *secapi.Iterator[schema.NetworkSku]
		var err error
		if opts != nil {
			iter, err = api.ListSkusWithFilters(builder.t.Context(), tref.Tenant, opts)
		} else {
			iter, err = api.ListSkus(builder.t.Context(), tref.Tenant)
		}
		requireNoError(sCtx, err)

		// Iterate through all items
		resp, err := iter.All(builder.t.Context())
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireLenResponse(sCtx, len(resp))
	})
	return resp
}
