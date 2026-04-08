//nolint:dupl
package steps

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/pkg/wrappers"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
)

// Sku

func (configurator *StepsConfigurator) ListNetworkSkusV1Step(stepName string, api secapi.NetworkV1, tpath secapi.TenantPath, opts *secapi.ListOptions) {
	listTenantResourcesStep(configurator.t, configurator.suite,
		listTenantResourcesParams[schema.NetworkSku, schema.SkuResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.NetworkSku, schema.SkuResourceMetadata, secapi.TenantPath]{
				path: tpath, listOptions: opts,
				listFunc: func(ctx context.Context, path secapi.TenantPath, options *secapi.ListOptions) (*secapi.Iterator[schema.NetworkSku], error) {
					return api.ListSkusWithOptions(ctx, path, options)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkSkuV1StepParams,
			operationName:  constants.ListSkusOperation,
		},
	)
}

// Network

func (configurator *StepsConfigurator) CreateOrUpdateNetworkV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, resource *schema.Network,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		createOrUpdateWorkspaceResourceParams[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.CreateOrUpdateNetworkOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Network) (
				wrappers.ResourceWrapper[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus], error,
			) {
				resp, err := api.CreateOrUpdateNetwork(configurator.t.Context(), resource)
				return wrappers.NewNetworkWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyNetworkSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) ListNetworkV1Step(stepName string, api secapi.NetworkV1, wpath secapi.WorkspacePath, opts *secapi.ListOptions) {
	listWorkspaceResourcesStep(configurator.t, configurator.suite,
		listWorkspaceResourcesParams[schema.Network, schema.RegionalWorkspaceResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.Network, schema.RegionalWorkspaceResourceMetadata, secapi.WorkspacePath]{
				path: wpath, listOptions: opts,
				listFunc: func(ctx context.Context, path secapi.WorkspacePath, options *secapi.ListOptions) (*secapi.Iterator[schema.Network], error) {
					return api.ListNetworksWithOptions(ctx, path, options)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.ListNetworksOperation,
			workspace:      wpath.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) GetNetworkV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec],
) *schema.Network {
	responseExpects.Metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.GetNetworkOperation,
			wref:           wref,
			getValueFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (
				wrappers.ResourceWrapper[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus], error,
			) {
				resp, err := api.GetNetworkUntilState(ctx, wref, config)
				return wrappers.NewNetworkWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyNetworkSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) WatchNetworkUntilDeletedV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, wref secapi.WorkspaceReference) {
	watchWorkspaceResourceUntilDeletedStep(configurator.t.Context(), configurator.suite, stepCreator,
		watchWorkspaceResourceUntilDeletedParams{
			watchResourceUntilDeletedParams: watchResourceUntilDeletedParams[secapi.WorkspaceReference]{
				reference: wref,
				getErrorFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig) error {
					return api.WatchNetworkUntilDeleted(configurator.t.Context(), wref, config)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.GetNetworkOperation,
		},
	)
}

func (configurator *StepsConfigurator) DeleteNetworkV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, resource *schema.Network) {
	deleteWorkspaceResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		deleteWorkspaceResourceParams[schema.Network]{
			deleteResourceParams: deleteResourceParams[schema.Network]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.Network) error {
					return api.DeleteNetwork(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.DeleteNetworkOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
		},
	)
}

// Internet Gateway

func (configurator *StepsConfigurator) CreateOrUpdateInternetGatewayV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, resource *schema.InternetGateway,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		createOrUpdateWorkspaceResourceParams[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.Status]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.CreateOrUpdateInternetGatewayOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.InternetGateway) (
				wrappers.ResourceWrapper[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.Status], error,
			) {
				resp, err := api.CreateOrUpdateInternetGateway(configurator.t.Context(), resource)
				return wrappers.NewInternetGatewayWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyInternetGatewaySpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) ListInternetGatewayV1Step(stepName string, api secapi.NetworkV1, wpath secapi.WorkspacePath, opts *secapi.ListOptions) {
	listWorkspaceResourcesStep(configurator.t, configurator.suite,
		listWorkspaceResourcesParams[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, secapi.WorkspacePath]{
				path: wpath, listOptions: opts,
				listFunc: func(ctx context.Context, path secapi.WorkspacePath, options *secapi.ListOptions) (*secapi.Iterator[schema.InternetGateway], error) {
					return api.ListInternetGatewaysWithOptions(ctx, path, options)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.ListInternetGatewaysOperation,
			workspace:      wpath.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) GetInternetGatewayV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec],
) *schema.InternetGateway {
	responseExpects.Metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.Status]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.GetInternetGatewayOperation,
			wref:           wref,
			getValueFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (
				wrappers.ResourceWrapper[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.Status], error,
			) {
				resp, err := api.GetInternetGatewayUntilState(ctx, wref, config)
				return wrappers.NewInternetGatewayWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyInternetGatewaySpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) WatchInternetGatewayUntilDeletedV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, wref secapi.WorkspaceReference) {
	watchWorkspaceResourceUntilDeletedStep(configurator.t.Context(), configurator.suite, stepCreator,
		watchWorkspaceResourceUntilDeletedParams{
			watchResourceUntilDeletedParams: watchResourceUntilDeletedParams[secapi.WorkspaceReference]{
				reference: wref,
				getErrorFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig) error {
					return api.WatchInternetGatewayUntilDeleted(configurator.t.Context(), wref, config)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.GetInternetGatewayOperation,
		},
	)
}

func (configurator *StepsConfigurator) DeleteInternetGatewayV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, resource *schema.InternetGateway) {
	deleteWorkspaceResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		deleteWorkspaceResourceParams[schema.InternetGateway]{
			deleteResourceParams: deleteResourceParams[schema.InternetGateway]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.InternetGateway) error {
					return api.DeleteInternetGateway(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.DeleteInternetGatewayOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
		},
	)
}

// Route Table

func (configurator *StepsConfigurator) CreateOrUpdateRouteTableV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, resource *schema.RouteTable,
	responseExpects ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateNetworkResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		createOrUpdateNetworkResourceParams[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  constants.CreateOrUpdateRouteTableOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
			network:        secapi.NetworkID(resource.Metadata.Network),
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.RouteTable) (
				wrappers.ResourceWrapper[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus], error,
			) {
				resp, err := api.CreateOrUpdateRouteTable(configurator.t.Context(), resource)
				return wrappers.NewRouteTableWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalNetworkResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyRouteTableSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) ListRouteTableV1Step(stepName string, api secapi.NetworkV1, npath secapi.NetworkPath, opts *secapi.ListOptions) {
	listNetworkResourcesStep(configurator.t, configurator.suite,
		listNetworkResourcesParams[schema.RouteTable, schema.RegionalNetworkResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.RouteTable, schema.RegionalNetworkResourceMetadata, secapi.NetworkPath]{
				path: npath, listOptions: opts,
				listFunc: func(ctx context.Context, path secapi.NetworkPath, options *secapi.ListOptions) (*secapi.Iterator[schema.RouteTable], error) {
					return api.ListRouteTablesWithOptions(ctx, path, options)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  constants.ListRouteTablesOperation,
			workspace:      npath.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) GetRouteTableV1Step(stepName string, api secapi.NetworkV1, nref secapi.NetworkReference,
	responseExpects ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec],
) *schema.RouteTable {
	responseExpects.Metadata.Verb = http.MethodGet
	return getNetworkResourceStep(configurator.t, configurator.suite,
		getNetworkResourceParams[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  constants.GetRouteTableOperation,
			nref:           nref,
			getValueFunc: func(ctx context.Context, nref secapi.NetworkReference, config secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (
				wrappers.ResourceWrapper[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus], error,
			) {
				resp, err := api.GetRouteTableUntilState(ctx, nref, config)
				return wrappers.NewRouteTableWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalNetworkResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyRouteTableSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) WatchRouteTableUntilDeletedV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, nref secapi.NetworkReference) {
	watchNetworkResourceUntilDeletedStep(configurator.t.Context(), configurator.suite, stepCreator,
		watchNetworkResourceUntilDeletedParams{
			watchResourceUntilDeletedParams: watchResourceUntilDeletedParams[secapi.NetworkReference]{
				reference: nref,
				getErrorFunc: func(ctx context.Context, nref secapi.NetworkReference, config secapi.ResourceObserverConfig) error {
					return api.WatchRouteTableUntilDeleted(configurator.t.Context(), nref, config)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  constants.GetRouteTableOperation,
		},
	)
}

func (configurator *StepsConfigurator) DeleteRouteTableV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, resource *schema.RouteTable) {
	deleteNetworkResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		deleteNetworkResourceParams[schema.RouteTable]{
			deleteResourceParams: deleteResourceParams[schema.RouteTable]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.RouteTable) error {
					return api.DeleteRouteTable(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  constants.DeleteRouteTableOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
			network:        secapi.NetworkID(resource.Metadata.Network),
		},
	)
}

// Subnet

func (configurator *StepsConfigurator) CreateOrUpdateSubnetV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, resource *schema.Subnet,
	responseExpects ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateNetworkResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		createOrUpdateNetworkResourceParams[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  constants.CreateOrUpdateSubnetOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
			network:        secapi.NetworkID(resource.Metadata.Network),
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Subnet) (
				wrappers.ResourceWrapper[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus], error,
			) {
				resp, err := api.CreateOrUpdateSubnet(configurator.t.Context(), resource)
				return wrappers.NewSubnetWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalNetworkResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifySubnetSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) ListSubnetV1Step(stepName string, api secapi.NetworkV1, npath secapi.NetworkPath, opts *secapi.ListOptions) {
	listNetworkResourcesStep(configurator.t, configurator.suite,
		listNetworkResourcesParams[schema.Subnet, schema.RegionalNetworkResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.Subnet, schema.RegionalNetworkResourceMetadata, secapi.NetworkPath]{
				path: npath, listOptions: opts,
				listFunc: func(ctx context.Context, path secapi.NetworkPath, options *secapi.ListOptions) (*secapi.Iterator[schema.Subnet], error) {
					return api.ListSubnetsWithOptions(ctx, path, options)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  constants.ListSubnetsOperation,
			workspace:      npath.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) GetSubnetV1Step(stepName string, api secapi.NetworkV1, nref secapi.NetworkReference,
	responseExpects ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec],
) *schema.Subnet {
	responseExpects.Metadata.Verb = http.MethodGet
	return getNetworkResourceStep(configurator.t, configurator.suite,
		getNetworkResourceParams[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  constants.GetSubnetOperation,
			nref:           nref,
			getValueFunc: func(ctx context.Context, nref secapi.NetworkReference, config secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (
				wrappers.ResourceWrapper[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus], error,
			) {
				resp, err := api.GetSubnetUntilState(ctx, nref, config)
				return wrappers.NewSubnetWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalNetworkResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifySubnetSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) WatchSubnetUntilDeletedV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, nref secapi.NetworkReference) {
	watchNetworkResourceUntilDeletedStep(configurator.t.Context(), configurator.suite, stepCreator,
		watchNetworkResourceUntilDeletedParams{
			watchResourceUntilDeletedParams: watchResourceUntilDeletedParams[secapi.NetworkReference]{
				reference: nref,
				getErrorFunc: func(ctx context.Context, nref secapi.NetworkReference, config secapi.ResourceObserverConfig) error {
					return api.WatchSubnetUntilDeleted(configurator.t.Context(), nref, config)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  constants.GetSubnetOperation,
		},
	)
}

func (configurator *StepsConfigurator) DeleteSubnetV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, resource *schema.Subnet) {
	deleteNetworkResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		deleteNetworkResourceParams[schema.Subnet]{
			deleteResourceParams: deleteResourceParams[schema.Subnet]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.Subnet) error {
					return api.DeleteSubnet(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  constants.DeleteSubnetOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
			network:        secapi.NetworkID(resource.Metadata.Network),
		},
	)
}

// Public Ip

func (configurator *StepsConfigurator) CreateOrUpdatePublicIpV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, resource *schema.PublicIp,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		createOrUpdateWorkspaceResourceParams[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.CreateOrUpdatePublicIpOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.PublicIp) (
				wrappers.ResourceWrapper[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus], error,
			) {
				resp, err := api.CreateOrUpdatePublicIp(configurator.t.Context(), resource)
				return wrappers.NewPublicIpWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyPublicIpSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) ListPublicIpV1Step(stepName string, api secapi.NetworkV1, wpath secapi.WorkspacePath, opts *secapi.ListOptions) {
	listWorkspaceResourcesStep(configurator.t, configurator.suite,
		listWorkspaceResourcesParams[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, secapi.WorkspacePath]{
				path: wpath, listOptions: opts,
				listFunc: func(ctx context.Context, path secapi.WorkspacePath, options *secapi.ListOptions) (*secapi.Iterator[schema.PublicIp], error) {
					return api.ListPublicIpsWithOptions(ctx, path, options)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.ListPublicIpsOperation,
			workspace:      wpath.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) GetPublicIpV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec],
) *schema.PublicIp {
	responseExpects.Metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.GetPublicIpOperation,
			wref:           wref,
			getValueFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (
				wrappers.ResourceWrapper[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus], error,
			) {
				resp, err := api.GetPublicIpUntilState(ctx, wref, config)
				return wrappers.NewPublicIpWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyPublicIpSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) WatchPublicIpUntilDeletedV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, wref secapi.WorkspaceReference) {
	watchWorkspaceResourceUntilDeletedStep(configurator.t.Context(), configurator.suite, stepCreator,
		watchWorkspaceResourceUntilDeletedParams{
			watchResourceUntilDeletedParams: watchResourceUntilDeletedParams[secapi.WorkspaceReference]{
				reference: wref,
				getErrorFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig) error {
					return api.WatchPublicIpUntilDeleted(configurator.t.Context(), wref, config)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.GetPublicIpOperation,
		},
	)
}

func (configurator *StepsConfigurator) DeletePublicIpV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, resource *schema.PublicIp) {
	deleteWorkspaceResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		deleteWorkspaceResourceParams[schema.PublicIp]{
			deleteResourceParams: deleteResourceParams[schema.PublicIp]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.PublicIp) error {
					return api.DeletePublicIp(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.DeletePublicIpOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
		},
	)
}

// Nic

func (configurator *StepsConfigurator) CreateOrUpdateNicV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, resource *schema.Nic,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		createOrUpdateWorkspaceResourceParams[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.CreateOrUpdateNicOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Nic) (
				wrappers.ResourceWrapper[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus], error,
			) {
				resp, err := api.CreateOrUpdateNic(configurator.t.Context(), resource)
				return wrappers.NewNicWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyNicSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) ListNicV1Step(stepName string, api secapi.NetworkV1, wpath secapi.WorkspacePath, opts *secapi.ListOptions) {
	listWorkspaceResourcesStep(configurator.t, configurator.suite,
		listWorkspaceResourcesParams[schema.Nic, schema.RegionalWorkspaceResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.Nic, schema.RegionalWorkspaceResourceMetadata, secapi.WorkspacePath]{
				path: wpath, listOptions: opts,
				listFunc: func(ctx context.Context, path secapi.WorkspacePath, options *secapi.ListOptions) (*secapi.Iterator[schema.Nic], error) {
					return api.ListNicsWithOptions(ctx, path, options)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.ListNicsOperation,
			workspace:      wpath.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) GetNicV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec],
) *schema.Nic {
	responseExpects.Metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.GetNicOperation,
			wref:           wref,
			getValueFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (
				wrappers.ResourceWrapper[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus], error,
			) {
				resp, err := api.GetNicUntilState(ctx, wref, config)
				return wrappers.NewNicWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyNicSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) WatchNicUntilDeletedV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, wref secapi.WorkspaceReference) {
	watchWorkspaceResourceUntilDeletedStep(configurator.t.Context(), configurator.suite, stepCreator,
		watchWorkspaceResourceUntilDeletedParams{
			watchResourceUntilDeletedParams: watchResourceUntilDeletedParams[secapi.WorkspaceReference]{
				reference: wref,
				getErrorFunc: func(ctx context.Context, tref secapi.WorkspaceReference, config secapi.ResourceObserverConfig) error {
					return api.WatchNicUntilDeleted(configurator.t.Context(), wref, config)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.GetNicOperation,
		},
	)
}

func (configurator *StepsConfigurator) DeleteNicV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, resource *schema.Nic) {
	deleteWorkspaceResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		deleteWorkspaceResourceParams[schema.Nic]{
			deleteResourceParams: deleteResourceParams[schema.Nic]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.Nic) error {
					return api.DeleteNic(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.DeleteNicOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
		},
	)
}

// Security Group Rule

func (configurator *StepsConfigurator) CreateOrUpdateSecurityGroupRuleV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, resource *schema.SecurityGroupRule,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		createOrUpdateWorkspaceResourceParams[schema.SecurityGroupRule, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec, schema.Status]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.CreateOrUpdateSecurityGroupRuleOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.SecurityGroupRule) (
				wrappers.ResourceWrapper[schema.SecurityGroupRule, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec, schema.SecurityGroupRuleStatus], error,
			) {
				resp, err := api.CreateOrUpdateSecurityGroupRule(configurator.t.Context(), resource)
				return wrappers.NewSecurityGroupRuleWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifySecurityGroupRuleSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) ListSecurityGroupRuleV1Step(stepName string, api secapi.NetworkV1, wpath secapi.WorkspacePath, opts *secapi.ListOptions) {
	listWorkspaceResourcesStep(configurator.t, configurator.suite,
		listWorkspaceResourcesParams[schema.SecurityGroupRule, schema.RegionalWorkspaceResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.SecurityGroupRule, schema.RegionalWorkspaceResourceMetadata, secapi.WorkspacePath]{
				path: wpath, listOptions: opts,
				listFunc: func(ctx context.Context, path secapi.WorkspacePath, options *secapi.ListOptions) (*secapi.Iterator[schema.SecurityGroupRule], error) {
					return api.ListSecurityGroupRulesWithOptions(ctx, path, options)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.ListSecurityGroupRulesOperation,
			workspace:      wpath.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) GetSecurityGroupRuleV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec],
) *schema.SecurityGroupRule {
	responseExpects.Metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.SecurityGroupRule, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec, schema.Status]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.GetSecurityGroupRuleOperation,
			wref:           wref,
			getValueFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (
				wrappers.ResourceWrapper[schema.SecurityGroupRule, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec, schema.SecurityGroupRuleStatus], error,
			) {
				resp, err := api.GetSecurityGroupRuleUntilState(ctx, wref, config)
				return wrappers.NewSecurityGroupRuleWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifySecurityGroupRuleSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) WatchSecurityGroupRuleUntilDeletedV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, wref secapi.WorkspaceReference) {
	watchWorkspaceResourceUntilDeletedStep(configurator.t.Context(), configurator.suite, stepCreator,
		watchWorkspaceResourceUntilDeletedParams{
			watchResourceUntilDeletedParams: watchResourceUntilDeletedParams[secapi.WorkspaceReference]{
				reference: wref,
				getErrorFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig) error {
					return api.WatchSecurityGroupRuleUntilDeleted(configurator.t.Context(), wref, config)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.GetSecurityGroupRuleOperation,
		},
	)
}

func (configurator *StepsConfigurator) DeleteSecurityGroupRuleV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, resource *schema.SecurityGroupRule) {
	deleteWorkspaceResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		deleteWorkspaceResourceParams[schema.SecurityGroupRule]{
			deleteResourceParams: deleteResourceParams[schema.SecurityGroupRule]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.SecurityGroupRule) error {
					return api.DeleteSecurityGroupRule(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.DeleteSecurityGroupRuleOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
		},
	)
}

// Security Group

func (configurator *StepsConfigurator) CreateOrUpdateSecurityGroupV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, resource *schema.SecurityGroup,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		createOrUpdateWorkspaceResourceParams[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.CreateOrUpdateSecurityGroupOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.SecurityGroup) (
				wrappers.ResourceWrapper[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus], error,
			) {
				resp, err := api.CreateOrUpdateSecurityGroup(configurator.t.Context(), resource)
				return wrappers.NewSecurityGroupWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifySecurityGroupSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) ListSecurityGroupV1Step(stepName string, api secapi.NetworkV1, wpath secapi.WorkspacePath, opts *secapi.ListOptions) {
	listWorkspaceResourcesStep(configurator.t, configurator.suite,
		listWorkspaceResourcesParams[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, secapi.WorkspacePath]{
				path: wpath, listOptions: opts,
				listFunc: func(ctx context.Context, path secapi.WorkspacePath, options *secapi.ListOptions) (*secapi.Iterator[schema.SecurityGroup], error) {
					return api.ListSecurityGroupsWithOptions(ctx, path, options)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.ListSecurityGroupsOperation,
			workspace:      wpath.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) GetSecurityGroupV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec],
) *schema.SecurityGroup {
	responseExpects.Metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.GetSecurityGroupOperation,
			wref:           wref,
			getValueFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (
				wrappers.ResourceWrapper[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus], error,
			) {
				resp, err := api.GetSecurityGroupUntilState(ctx, wref, config)
				return wrappers.NewSecurityGroupWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifySecurityGroupSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) WatchSecurityGroupUntilDeletedV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, wref secapi.WorkspaceReference) {
	watchWorkspaceResourceUntilDeletedStep(configurator.t.Context(), configurator.suite, stepCreator,
		watchWorkspaceResourceUntilDeletedParams{
			watchResourceUntilDeletedParams: watchResourceUntilDeletedParams[secapi.WorkspaceReference]{
				reference: wref,
				getErrorFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig) error {
					return api.WatchSecurityGroupUntilDeleted(configurator.t.Context(), wref, config)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.GetSecurityGroupOperation,
		},
	)
}

func (configurator *StepsConfigurator) DeleteSecurityGroupV1Step(stepName string, stepCreator StepCreator, api secapi.NetworkV1, resource *schema.SecurityGroup) {
	deleteWorkspaceResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		deleteWorkspaceResourceParams[schema.SecurityGroup]{
			deleteResourceParams: deleteResourceParams[schema.SecurityGroup]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.SecurityGroup) error {
					return api.DeleteSecurityGroup(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.DeleteSecurityGroupOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
		},
	)
}
