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
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "CreateOrUpdateNetwork",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Network) (
				*stepFuncResponse[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus], error,
			) {
				resp, err := api.CreateOrUpdateNetwork(configurator.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifyNetworkSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetNetworkV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec],
) *schema.Network {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetNetwork",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
				*stepFuncResponse[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus], error,
			) {
				resp, err := api.GetNetworkUntilState(ctx, wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifyNetworkSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetListNetworkV1Step(
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
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetNetwork",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference) error {
				_, err := api.GetNetwork(ctx, wref)
				return err
			},
			expectedError: expectedError,
		},
	)
}

func (configurator *StepsConfigurator) DeleteNetworkV1Step(stepName string, api secapi.NetworkV1, resource *schema.Network) {
	configurator.logStepName(stepName)
	deleteWorkspaceResourceStep(configurator.t,
		deleteWorkspaceResourceParams[schema.Network]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "DeleteNetwork",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			deleteFunc: func(ctx context.Context, r *schema.Network) error {
				return api.DeleteNetwork(ctx, r)
			},
		},
	)
}

// Internet Gateway

func (configurator *StepsConfigurator) CreateOrUpdateInternetGatewayV1Step(stepName string, api secapi.NetworkV1, resource *schema.InternetGateway,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.Status]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "CreateOrUpdateInternetGateway",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.InternetGateway) (
				*stepFuncResponse[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.Status], error,
			) {
				resp, err := api.CreateOrUpdateInternetGateway(configurator.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifyInternetGatewaySpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetInternetGatewayV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec],
) *schema.InternetGateway {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.Status]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetInternetGateway",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
				*stepFuncResponse[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.Status], error,
			) {
				resp, err := api.GetInternetGatewayUntilState(ctx, wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifyInternetGatewaySpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetListInternetGatewayV1Step(
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
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetInternetGateway",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference) error {
				_, err := api.GetInternetGateway(ctx, wref)
				return err
			},
			expectedError: expectedError,
		},
	)
}

func (configurator *StepsConfigurator) DeleteInternetGatewayV1Step(stepName string, api secapi.NetworkV1, resource *schema.InternetGateway) {
	configurator.logStepName(stepName)
	deleteWorkspaceResourceStep(configurator.t,
		deleteWorkspaceResourceParams[schema.InternetGateway]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "DeleteInternetGateway",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			deleteFunc: func(ctx context.Context, r *schema.InternetGateway) error {
				return api.DeleteInternetGateway(ctx, r)
			},
		},
	)
}

// Route Table

func (configurator *StepsConfigurator) CreateOrUpdateRouteTableV1Step(stepName string, api secapi.NetworkV1, resource *schema.RouteTable,
	responseExpects ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateNetworkResourceStep(configurator.t, configurator.suite,
		createOrUpdateNetworkResourceParams[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  "CreateOrUpdateRouteTable",
			workspace:      resource.Metadata.Workspace,
			network:        resource.Metadata.Network,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.RouteTable) (
				*stepFuncResponse[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus], error,
			) {
				resp, err := api.CreateOrUpdateRouteTable(configurator.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalNetworkResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifyRouteTableSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetRouteTableV1Step(stepName string, api secapi.NetworkV1, nref secapi.NetworkReference,
	responseExpects ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec],
) *schema.RouteTable {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getNetworkResourceStep(configurator.t, configurator.suite,
		getNetworkResourceParams[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  "GetRouteTable",
			nref:           nref,
			getFunc: func(ctx context.Context, nref secapi.NetworkReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
				*stepFuncResponse[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus], error,
			) {
				resp, err := api.GetRouteTableUntilState(ctx, nref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalNetworkResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifyRouteTableSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetListRouteTableV1Step(
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
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  "GetRouteTable",
			nref:           nref,
			getFunc: func(ctx context.Context, nref secapi.NetworkReference) error {
				_, err := api.GetRouteTable(ctx, nref)
				return err
			},
			expectedError: expectedError,
		},
	)
}

func (configurator *StepsConfigurator) DeleteRouteTableV1Step(stepName string, api secapi.NetworkV1, resource *schema.RouteTable) {
	configurator.logStepName(stepName)
	deleteNetworkResourceStep(configurator.t,
		deleteNetworkResourceParams[schema.RouteTable]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  "DeleteRouteTable",
			workspace:      resource.Metadata.Workspace,
			network:        resource.Metadata.Network,
			resource:       resource,
			deleteFunc: func(ctx context.Context, r *schema.RouteTable) error {
				return api.DeleteRouteTable(ctx, r)
			},
		},
	)
}

// Subnet

func (configurator *StepsConfigurator) CreateOrUpdateSubnetV1Step(stepName string, api secapi.NetworkV1, resource *schema.Subnet,
	responseExpects ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateNetworkResourceStep(configurator.t, configurator.suite,
		createOrUpdateNetworkResourceParams[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  "CreateOrUpdateSubnet",
			workspace:      resource.Metadata.Workspace,
			network:        resource.Metadata.Network,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Subnet) (
				*stepFuncResponse[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus], error,
			) {
				resp, err := api.CreateOrUpdateSubnet(configurator.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalNetworkResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifySubnetSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetSubnetV1Step(stepName string, api secapi.NetworkV1, nref secapi.NetworkReference,
	responseExpects ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec],
) *schema.Subnet {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getNetworkResourceStep(configurator.t, configurator.suite,
		getNetworkResourceParams[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  "GetSubnet",
			nref:           nref,
			getFunc: func(ctx context.Context, nref secapi.NetworkReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
				*stepFuncResponse[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus], error,
			) {
				resp, err := api.GetSubnetUntilState(ctx, nref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalNetworkResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifySubnetSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetListSubnetV1Step(
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
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  "GetSubnet",
			nref:           nref,
			getFunc: func(ctx context.Context, nref secapi.NetworkReference) error {
				_, err := api.GetSubnet(ctx, nref)
				return err
			},
			expectedError: expectedError,
		},
	)
}

func (configurator *StepsConfigurator) DeleteSubnetV1Step(stepName string, api secapi.NetworkV1, resource *schema.Subnet) {
	configurator.logStepName(stepName)
	deleteNetworkResourceStep(configurator.t,
		deleteNetworkResourceParams[schema.Subnet]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  "DeleteSubnet",
			workspace:      resource.Metadata.Workspace,
			network:        resource.Metadata.Network,
			resource:       resource,
			deleteFunc: func(ctx context.Context, r *schema.Subnet) error {
				return api.DeleteSubnet(ctx, r)
			},
		},
	)
}

// Public Ip

func (configurator *StepsConfigurator) CreateOrUpdatePublicIpV1Step(stepName string, api secapi.NetworkV1, resource *schema.PublicIp,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "CreateOrUpdatePublicIp",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.PublicIp) (
				*stepFuncResponse[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus], error,
			) {
				resp, err := api.CreateOrUpdatePublicIp(configurator.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifyPublicIpSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetPublicIpV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec],
) *schema.PublicIp {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetPublicIp",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
				*stepFuncResponse[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus], error,
			) {
				resp, err := api.GetPublicIpUntilState(ctx, wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifyPublicIpSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetListPublicIpV1Step(
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
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetPublicIp",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference) error {
				_, err := api.GetPublicIp(ctx, wref)
				return err
			},
			expectedError: expectedError,
		},
	)
}

func (configurator *StepsConfigurator) DeletePublicIpV1Step(stepName string, api secapi.NetworkV1, resource *schema.PublicIp) {
	configurator.logStepName(stepName)
	deleteWorkspaceResourceStep(configurator.t,
		deleteWorkspaceResourceParams[schema.PublicIp]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "DeletePublicIp",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			deleteFunc: func(ctx context.Context, r *schema.PublicIp) error {
				return api.DeletePublicIp(ctx, r)
			},
		},
	)
}

// Nic

func (configurator *StepsConfigurator) CreateOrUpdateNicV1Step(stepName string, api secapi.NetworkV1, resource *schema.Nic,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "CreateOrUpdateNic",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Nic) (
				*stepFuncResponse[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus], error,
			) {
				resp, err := api.CreateOrUpdateNic(configurator.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifyNicSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetNicV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec],
) *schema.Nic {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetNic",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
				*stepFuncResponse[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus], error,
			) {
				resp, err := api.GetNicUntilState(ctx, wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifyNicSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetListNicV1Step(
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
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetNic",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference) error {
				_, err := api.GetNic(ctx, wref)
				return err
			},
			expectedError: expectedError,
		},
	)
}

func (configurator *StepsConfigurator) DeleteNicV1Step(stepName string, api secapi.NetworkV1, resource *schema.Nic) {
	configurator.logStepName(stepName)
	deleteWorkspaceResourceStep(configurator.t,
		deleteWorkspaceResourceParams[schema.Nic]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "DeleteNic",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			deleteFunc: func(ctx context.Context, r *schema.Nic) error {
				return api.DeleteNic(ctx, r)
			},
		},
	)
}

// Security Group

func (configurator *StepsConfigurator) CreateOrUpdateSecurityGroupV1Step(stepName string, api secapi.NetworkV1, resource *schema.SecurityGroup,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "CreateOrUpdateSecurityGroup",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.SecurityGroup) (
				*stepFuncResponse[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus], error,
			) {
				resp, err := api.CreateOrUpdateSecurityGroup(configurator.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifySecurityGroupSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetSecurityGroupV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec],
) *schema.SecurityGroup {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetSecurityGroup",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
				*stepFuncResponse[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus], error,
			) {
				resp, err := api.GetSecurityGroupUntilState(ctx, wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifySecurityGroupSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetListSecurityGroupV1Step(
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
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "GetSecurityGroup",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference) error {
				_, err := api.GetSecurityGroup(ctx, wref)
				return err
			},
			expectedError: expectedError,
		},
	)
}

func (configurator *StepsConfigurator) DeleteSecurityGroupV1Step(stepName string, api secapi.NetworkV1, resource *schema.SecurityGroup) {
	configurator.logStepName(stepName)
	deleteWorkspaceResourceStep(configurator.t,
		deleteWorkspaceResourceParams[schema.SecurityGroup]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  "DeleteSecurityGroup",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			deleteFunc: func(ctx context.Context, r *schema.SecurityGroup) error {
				return api.DeleteSecurityGroup(ctx, r)
			},
		},
	)
}

func (configurator *StepsConfigurator) GetListNetworkSkusV1Step(
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
		requireLenResponse(sCtx, len(resp))
	})
	return resp
}
