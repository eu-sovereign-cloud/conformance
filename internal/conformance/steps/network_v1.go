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

func (configurator *StepsConfigurator) CreateOrUpdateNetworkV1Step(stepName string, api secapi.NetworkV1, resource *schema.Network,
	responseExpects StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus]{
			createOrUpdateResourceParams: createOrUpdateResourceParams[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus]{
				resource: resource,
				createOrUpdateFunc: func(context.Context, *schema.Network) (
					*createOrUpdateStepFuncResponse[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus], error,
				) {
					if resp, err := api.CreateOrUpdateNetwork(configurator.t.Context(), resource); err == nil {
						return newCreateOrUpdateStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyNetworkSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "CreateOrUpdateNetwork",
			workspace:      resource.Metadata.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) GetNetworkV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec],
) *schema.Network {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus]{
			getResourceWithObserverParams: getResourceWithObserverParams[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus, secapi.WorkspaceReference, schema.ResourceState]{
				reference: wref,
				getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
					*getStepFuncResponse[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus], error,
				) {
					if resp, err := api.GetNetworkUntilState(ctx, wref, config); err == nil {
						return newGetStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyNetworkSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetNetwork",
		},
	)
}

func (configurator *StepsConfigurator) ListNetworkV1Step(
	stepName string,
	api secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) {
	configurator.logStepName(stepName)
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListNetwork", string(wref.Workspace))
		var iter *secapi.Iterator[schema.Network]
		var err error
		if opts != nil {
			iter, err = api.ListNetworksWithFilters(configurator.t.Context(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListNetworks(configurator.t.Context(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) GetNetworkWithErrorV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference, expectedError error) {
	configurator.logStepName(stepName)
	getWorkspaceResourceWithErrorStep(configurator.t,
		getWorkspaceResourceWithErrorParams{
			getResourceWithErrorParams: getResourceWithErrorParams[secapi.WorkspaceReference]{
				reference: wref,
				getFunc: func(ctx context.Context, wref secapi.WorkspaceReference) error {
					_, err := api.GetNetwork(ctx, wref)
					return err
				},
				expectedError: expectedError,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetNetwork",
		},
	)
}

func (configurator *StepsConfigurator) DeleteNetworkV1Step(stepName string, api secapi.NetworkV1, resource *schema.Network) {
	configurator.logStepName(stepName)
	deleteWorkspaceResourceStep(configurator.t,
		deleteWorkspaceResourceParams[schema.Network]{
			deleteResourceParams: deleteResourceParams[schema.Network]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.Network) error {
					return api.DeleteNetwork(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "DeleteNetwork",
			workspace:      resource.Metadata.Workspace,
		},
	)
}

// Internet Gateway

func (configurator *StepsConfigurator) CreateOrUpdateInternetGatewayV1Step(stepName string, api secapi.NetworkV1, resource *schema.InternetGateway,
	responseExpects StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.Status]{
			createOrUpdateResourceParams: createOrUpdateResourceParams[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.Status]{
				resource: resource,
				createOrUpdateFunc: func(context.Context, *schema.InternetGateway) (
					*createOrUpdateStepFuncResponse[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.Status], error,
				) {
					if resp, err := api.CreateOrUpdateInternetGateway(configurator.t.Context(), resource); err == nil {
						return newCreateOrUpdateStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyInternetGatewaySpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "CreateOrUpdateInternetGateway",
			workspace:      resource.Metadata.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) GetInternetGatewayV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec],
) *schema.InternetGateway {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.Status]{
			getResourceWithObserverParams: getResourceWithObserverParams[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.Status, secapi.WorkspaceReference, schema.ResourceState]{
				reference: wref,
				getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
					*getStepFuncResponse[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.Status], error,
				) {
					if resp, err := api.GetInternetGatewayUntilState(ctx, wref, config); err == nil {
						return newGetStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyInternetGatewaySpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetInternetGateway",
		},
	)
}

func (configurator *StepsConfigurator) ListInternetGatewayV1Step(
	stepName string,
	api secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) {
	configurator.logStepName(stepName)
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListInternetGateway", wref.Name)
		var iter *secapi.Iterator[schema.InternetGateway]
		var err error
		if opts != nil {
			iter, err = api.ListInternetGatewaysWithFilters(configurator.t.Context(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListInternetGateways(configurator.t.Context(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) GetInternetGatewayWithErrorV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference, expectedError error) {
	configurator.logStepName(stepName)
	getWorkspaceResourceWithErrorStep(configurator.t,
		getWorkspaceResourceWithErrorParams{
			getResourceWithErrorParams: getResourceWithErrorParams[secapi.WorkspaceReference]{
				reference: wref,
				getFunc: func(ctx context.Context, wref secapi.WorkspaceReference) error {
					_, err := api.GetInternetGateway(ctx, wref)
					return err
				},
				expectedError: expectedError,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetInternetGateway",
		},
	)
}

func (configurator *StepsConfigurator) DeleteInternetGatewayV1Step(stepName string, api secapi.NetworkV1, resource *schema.InternetGateway) {
	configurator.logStepName(stepName)
	deleteWorkspaceResourceStep(configurator.t,
		deleteWorkspaceResourceParams[schema.InternetGateway]{
			deleteResourceParams: deleteResourceParams[schema.InternetGateway]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.InternetGateway) error {
					return api.DeleteInternetGateway(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "DeleteInternetGateway",
			workspace:      resource.Metadata.Workspace,
		},
	)
}

// Route Table

func (configurator *StepsConfigurator) CreateOrUpdateRouteTableV1Step(stepName string, api secapi.NetworkV1, resource *schema.RouteTable,
	responseExpects StepResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateNetworkResourceStep(configurator.t, configurator.suite,
		createOrUpdateNetworkResourceParams[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus]{
			createOrUpdateResourceParams: createOrUpdateResourceParams[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus]{
				resource: resource,
				createOrUpdateFunc: func(context.Context, *schema.RouteTable) (
					*createOrUpdateStepFuncResponse[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus], error,
				) {
					if resp, err := api.CreateOrUpdateRouteTable(configurator.t.Context(), resource); err == nil {
						return newCreateOrUpdateStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalNetworkResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyRouteTableSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  "CreateOrUpdateRouteTable",
			workspace:      resource.Metadata.Workspace,
			network:        resource.Metadata.Network,
		},
	)
}

func (configurator *StepsConfigurator) GetRouteTableV1Step(stepName string, api secapi.NetworkV1, nref secapi.NetworkReference,
	responseExpects StepResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec],
) *schema.RouteTable {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getNetworkResourceStep(configurator.t, configurator.suite,
		getNetworkResourceParams[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus]{
			getResourceWithObserverParams: getResourceWithObserverParams[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus, secapi.NetworkReference, schema.ResourceState]{
				reference: nref,
				getFunc: func(ctx context.Context, nref secapi.NetworkReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
					*getStepFuncResponse[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus], error,
				) {
					if resp, err := api.GetRouteTableUntilState(ctx, nref, config); err == nil {
						return newGetStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalNetworkResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyRouteTableSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  "GetRouteTable",
		},
	)
}

func (configurator *StepsConfigurator) ListRouteTableV1Step(
	stepName string,
	api secapi.NetworkV1,
	nref secapi.NetworkReference,
	opts *secapi.ListOptions,
) {
	configurator.logStepName(stepName)
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListRouteTable", nref.Name)
		var iter *secapi.Iterator[schema.RouteTable]
		var err error
		if opts != nil {
			iter, err = api.ListRouteTablesWithFilters(configurator.t.Context(), nref.Tenant, nref.Workspace, nref.Network, opts)
		} else {
			iter, err = api.ListRouteTables(configurator.t.Context(), nref.Tenant, nref.Workspace, nref.Network)
		}
		requireNoError(sCtx, err)
		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) GetRouteTableWithErrorV1Step(stepName string, api secapi.NetworkV1, nref secapi.NetworkReference, expectedError error) {
	configurator.logStepName(stepName)
	getNetworkResourceWithErrorStep(configurator.t,
		getNetworkResourceWithErrorParams{
			getResourceWithErrorParams: getResourceWithErrorParams[secapi.NetworkReference]{
				reference: nref,
				getFunc: func(ctx context.Context, nref secapi.NetworkReference) error {
					_, err := api.GetRouteTable(ctx, nref)
					return err
				},
				expectedError: expectedError,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  "GetRouteTable",
		},
	)
}

func (configurator *StepsConfigurator) DeleteRouteTableV1Step(stepName string, api secapi.NetworkV1, resource *schema.RouteTable) {
	configurator.logStepName(stepName)
	deleteNetworkResourceStep(configurator.t,
		deleteNetworkResourceParams[schema.RouteTable]{
			deleteResourceParams: deleteResourceParams[schema.RouteTable]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.RouteTable) error {
					return api.DeleteRouteTable(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  "DeleteRouteTable",
			workspace:      resource.Metadata.Workspace,
			network:        resource.Metadata.Network,
		},
	)
}

// Subnet

func (configurator *StepsConfigurator) CreateOrUpdateSubnetV1Step(stepName string, api secapi.NetworkV1, resource *schema.Subnet,
	responseExpects StepResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateNetworkResourceStep(configurator.t, configurator.suite,
		createOrUpdateNetworkResourceParams[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus]{
			createOrUpdateResourceParams: createOrUpdateResourceParams[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus]{
				resource: resource,
				createOrUpdateFunc: func(context.Context, *schema.Subnet) (
					*createOrUpdateStepFuncResponse[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus], error,
				) {
					if resp, err := api.CreateOrUpdateSubnet(configurator.t.Context(), resource); err == nil {
						return newCreateOrUpdateStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalNetworkResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifySubnetSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  "CreateOrUpdateSubnet",
			workspace:      resource.Metadata.Workspace,
			network:        resource.Metadata.Network,
		},
	)
}

func (configurator *StepsConfigurator) GetSubnetV1Step(stepName string, api secapi.NetworkV1, nref secapi.NetworkReference,
	responseExpects StepResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec],
) *schema.Subnet {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getNetworkResourceStep(configurator.t, configurator.suite,
		getNetworkResourceParams[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus]{
			getResourceWithObserverParams: getResourceWithObserverParams[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus, secapi.NetworkReference, schema.ResourceState]{
				reference: nref,
				getFunc: func(ctx context.Context, nref secapi.NetworkReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
					*getStepFuncResponse[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus], error,
				) {
					if resp, err := api.GetSubnetUntilState(ctx, nref, config); err == nil {
						return newGetStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalNetworkResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifySubnetSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  "GetSubnet",
		},
	)
}

func (configurator *StepsConfigurator) ListSubnetV1Step(
	stepName string,
	api secapi.NetworkV1,
	nref secapi.NetworkReference,
	opts *secapi.ListOptions,
) {
	configurator.logStepName(stepName)
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListSubnet", nref.Name)
		var iter *secapi.Iterator[schema.Subnet]
		var err error
		if opts != nil {
			iter, err = api.ListSubnetsWithFilters(configurator.t.Context(), nref.Tenant, nref.Workspace, nref.Network, opts)
		} else {
			iter, err = api.ListSubnets(configurator.t.Context(), nref.Tenant, nref.Workspace, nref.Network)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) GetSubnetWithErrorV1Step(stepName string, api secapi.NetworkV1, nref secapi.NetworkReference, expectedError error) {
	configurator.logStepName(stepName)
	getNetworkResourceWithErrorStep(configurator.t,
		getNetworkResourceWithErrorParams{
			getResourceWithErrorParams: getResourceWithErrorParams[secapi.NetworkReference]{
				reference: nref,
				getFunc: func(ctx context.Context, nref secapi.NetworkReference) error {
					_, err := api.GetSubnet(ctx, nref)
					return err
				},
				expectedError: expectedError,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  "GetSubnet",
		},
	)
}

func (configurator *StepsConfigurator) DeleteSubnetV1Step(stepName string, api secapi.NetworkV1, resource *schema.Subnet) {
	configurator.logStepName(stepName)
	deleteNetworkResourceStep(configurator.t,
		deleteNetworkResourceParams[schema.Subnet]{
			deleteResourceParams: deleteResourceParams[schema.Subnet]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.Subnet) error {
					return api.DeleteSubnet(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  "DeleteSubnet",
			workspace:      resource.Metadata.Workspace,
			network:        resource.Metadata.Network,
		},
	)
}

// Public Ip

func (configurator *StepsConfigurator) CreateOrUpdatePublicIpV1Step(stepName string, api secapi.NetworkV1, resource *schema.PublicIp,
	responseExpects StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus]{
			createOrUpdateResourceParams: createOrUpdateResourceParams[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus]{
				resource: resource,
				createOrUpdateFunc: func(context.Context, *schema.PublicIp) (
					*createOrUpdateStepFuncResponse[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus], error,
				) {
					if resp, err := api.CreateOrUpdatePublicIp(configurator.t.Context(), resource); err == nil {
						return newCreateOrUpdateStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyPublicIpSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "CreateOrUpdatePublicIp",
			workspace:      resource.Metadata.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) GetPublicIpV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec],
) *schema.PublicIp {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus]{
			getResourceWithObserverParams: getResourceWithObserverParams[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus, secapi.WorkspaceReference, schema.ResourceState]{
				reference: wref,
				getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
					*getStepFuncResponse[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus], error,
				) {
					if resp, err := api.GetPublicIpUntilState(ctx, wref, config); err == nil {
						return newGetStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyPublicIpSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetPublicIp",
		},
	)
}

func (configurator *StepsConfigurator) ListPublicIpV1Step(
	stepName string,
	api secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) {
	configurator.logStepName(stepName)
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListPublicIp", wref.Name)
		var iter *secapi.Iterator[schema.PublicIp]
		var err error
		if opts != nil {
			iter, err = api.ListPublicIpsWithFilters(configurator.t.Context(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListPublicIps(configurator.t.Context(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)

		// Iterate through all items
		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) GetPublicIpWithErrorV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference, expectedError error) {
	configurator.logStepName(stepName)
	getWorkspaceResourceWithErrorStep(configurator.t,
		getWorkspaceResourceWithErrorParams{
			getResourceWithErrorParams: getResourceWithErrorParams[secapi.WorkspaceReference]{
				reference: wref,
				getFunc: func(ctx context.Context, wref secapi.WorkspaceReference) error {
					_, err := api.GetPublicIp(ctx, wref)
					return err
				},
				expectedError: expectedError,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetPublicIp",
		},
	)
}

func (configurator *StepsConfigurator) DeletePublicIpV1Step(stepName string, api secapi.NetworkV1, resource *schema.PublicIp) {
	configurator.logStepName(stepName)
	deleteWorkspaceResourceStep(configurator.t,
		deleteWorkspaceResourceParams[schema.PublicIp]{
			deleteResourceParams: deleteResourceParams[schema.PublicIp]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.PublicIp) error {
					return api.DeletePublicIp(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "DeletePublicIp",
			workspace:      resource.Metadata.Workspace,
		},
	)
}

// Nic

func (configurator *StepsConfigurator) CreateOrUpdateNicV1Step(stepName string, api secapi.NetworkV1, resource *schema.Nic,
	responseExpects StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus]{
			createOrUpdateResourceParams: createOrUpdateResourceParams[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus]{
				resource: resource,
				createOrUpdateFunc: func(context.Context, *schema.Nic) (
					*createOrUpdateStepFuncResponse[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus], error,
				) {
					if resp, err := api.CreateOrUpdateNic(configurator.t.Context(), resource); err == nil {
						return newCreateOrUpdateStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyNicSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "CreateOrUpdateNic",
			workspace:      resource.Metadata.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) GetNicV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec],
) *schema.Nic {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus]{
			getResourceWithObserverParams: getResourceWithObserverParams[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus, secapi.WorkspaceReference, schema.ResourceState]{
				reference: wref,
				getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
					*getStepFuncResponse[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus], error,
				) {
					if resp, err := api.GetNicUntilState(ctx, wref, config); err == nil {
						return newGetStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyNicSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetNic",
		},
	)
}

func (configurator *StepsConfigurator) ListNicV1Step(
	stepName string,
	api secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) {
	configurator.logStepName(stepName)
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListNic", wref.Name)
		var iter *secapi.Iterator[schema.Nic]
		var err error
		if opts != nil {
			iter, err = api.ListNicsWithFilters(configurator.t.Context(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListNics(configurator.t.Context(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) GetNicWithErrorV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference, expectedError error) {
	configurator.logStepName(stepName)
	getWorkspaceResourceWithErrorStep(configurator.t,
		getWorkspaceResourceWithErrorParams{
			getResourceWithErrorParams: getResourceWithErrorParams[secapi.WorkspaceReference]{
				reference: wref,
				getFunc: func(ctx context.Context, wref secapi.WorkspaceReference) error {
					_, err := api.GetNic(ctx, wref)
					return err
				},
				expectedError: expectedError,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetNic",
		},
	)
}

func (configurator *StepsConfigurator) DeleteNicV1Step(stepName string, api secapi.NetworkV1, resource *schema.Nic) {
	configurator.logStepName(stepName)
	deleteWorkspaceResourceStep(configurator.t,
		deleteWorkspaceResourceParams[schema.Nic]{
			deleteResourceParams: deleteResourceParams[schema.Nic]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.Nic) error {
					return api.DeleteNic(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "DeleteNic",
			workspace:      resource.Metadata.Workspace,
		},
	)
}

// Security Group

func (configurator *StepsConfigurator) CreateOrUpdateSecurityGroupV1Step(stepName string, api secapi.NetworkV1, resource *schema.SecurityGroup,
	responseExpects StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus]{
			createOrUpdateResourceParams: createOrUpdateResourceParams[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus]{
				resource: resource,
				createOrUpdateFunc: func(context.Context, *schema.SecurityGroup) (
					*createOrUpdateStepFuncResponse[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus], error,
				) {
					if resp, err := api.CreateOrUpdateSecurityGroup(configurator.t.Context(), resource); err == nil {
						return newCreateOrUpdateStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifySecurityGroupSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "CreateOrUpdateSecurityGroup",
			workspace:      resource.Metadata.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) GetSecurityGroupV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec],
) *schema.SecurityGroup {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus]{
			getResourceWithObserverParams: getResourceWithObserverParams[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus, secapi.WorkspaceReference, schema.ResourceState]{
				reference: wref,
				getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
					*getStepFuncResponse[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus], error,
				) {
					if resp, err := api.GetSecurityGroupUntilState(ctx, wref, config); err == nil {
						return newGetStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifySecurityGroupSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetSecurityGroup",
		},
	)
}

func (configurator *StepsConfigurator) ListSecurityGroupV1Step(
	stepName string,
	api secapi.NetworkV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) {
	configurator.logStepName(stepName)
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListSecurityGroup", wref.Name)
		var iter *secapi.Iterator[schema.SecurityGroup]
		var err error
		if opts != nil {
			iter, err = api.ListSecurityGroupsWithFilters(configurator.t.Context(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListSecurityGroups(configurator.t.Context(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) GetSecurityGroupWithErrorV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference, expectedError error) {
	configurator.logStepName(stepName)
	getWorkspaceResourceWithErrorStep(configurator.t,
		getWorkspaceResourceWithErrorParams{
			getResourceWithErrorParams: getResourceWithErrorParams[secapi.WorkspaceReference]{
				reference: wref,
				getFunc: func(ctx context.Context, wref secapi.WorkspaceReference) error {
					_, err := api.GetSecurityGroup(ctx, wref)
					return err
				},
				expectedError: expectedError,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetSecurityGroup",
		},
	)
}

func (configurator *StepsConfigurator) DeleteSecurityGroupV1Step(stepName string, api secapi.NetworkV1, resource *schema.SecurityGroup) {
	configurator.logStepName(stepName)
	deleteWorkspaceResourceStep(configurator.t,
		deleteWorkspaceResourceParams[schema.SecurityGroup]{
			deleteResourceParams: deleteResourceParams[schema.SecurityGroup]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.SecurityGroup) error {
					return api.DeleteSecurityGroup(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "DeleteSecurityGroup",
			workspace:      resource.Metadata.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) ListNetworkSkusV1Step(
	stepName string,
	api secapi.NetworkV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) []*schema.NetworkSku {
	var resp []*schema.NetworkSku
	configurator.logStepName(stepName)
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetNetworkV1StepParams(sCtx, "ListSkus", tref.Name)

		var iter *secapi.Iterator[schema.NetworkSku]
		var err error
		if opts != nil {
			iter, err = api.ListSkusWithFilters(configurator.t.Context(), tref.Tenant, opts)
		} else {
			iter, err = api.ListSkus(configurator.t.Context(), tref.Tenant)
		}
		requireNoError(sCtx, err)

		// Iterate through all items
		resp, err := iter.All(configurator.t.Context())
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireNotEmptyResponse(sCtx, resp)
	})
	return resp
}
