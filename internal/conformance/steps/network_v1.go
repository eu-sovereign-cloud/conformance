//nolint:dupl
package steps

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/pkg/wrappers"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Sku

func (configurator *StepsConfigurator) ListNetworkSkusV1Step(
	stepName string, api secapi.NetworkV1, tref secapi.TenantReference, opts *secapi.ListOptions,
) []*schema.NetworkSku {
	var resp []*schema.NetworkSku
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetComputeV1StepParams(sCtx, "ListSkus", tref.Name)

		iter, err := api.ListSkusWithOptions(configurator.t.Context(), secapi.TenantPath{Tenant: tref.Tenant}, opts)
		requireNoError(sCtx, err)

		// Iterate through all items
		resp, err := iter.All(configurator.t.Context())
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireLenResponse(sCtx, len(resp))
	})
	return resp
}

// Network

func (configurator *StepsConfigurator) CreateOrUpdateNetworkV1Step(stepName string, api secapi.NetworkV1, resource *schema.Network,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.Network, schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec, schema.NetworkStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.CreateOrUpdateNetworkOperation,
			workspace:      resource.Metadata.Workspace,
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

func (configurator *StepsConfigurator) GetNetworkV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.NetworkSpec],
) *schema.Network {
	responseExpects.Metadata.Verb = http.MethodGet
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
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
			expectedResourceStatus: responseExpects.ResourceStatus,
		},
	)
}

func (configurator *StepsConfigurator) ListNetworkV1Step(
	stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference, opts *secapi.ListOptions,
) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "ListNetwork", string(wref.Workspace))
		iter, err := api.ListNetworksWithOptions(configurator.t.Context(), secapi.WorkspacePath{Tenant: wref.Tenant, Workspace: wref.Workspace}, opts)
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) WatchNetworkUntilDeletedV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		watchWorkspaceResourceUntilDeletedStep(configurator.t, configurator.suite,
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
	})
}

func (configurator *StepsConfigurator) DeleteNetworkV1Step(stepName string, api secapi.NetworkV1, resource *schema.Network) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetNetworkV1StepParams(sCtx, "DeleteNetwork", resource.Metadata.Workspace)

		err := api.DeleteNetwork(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Internet Gateway

func (configurator *StepsConfigurator) CreateOrUpdateInternetGatewayV1Step(stepName string, api secapi.NetworkV1, resource *schema.InternetGateway,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.InternetGateway, schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec, schema.Status]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.CreateOrUpdateInternetGatewayOperation,
			workspace:      resource.Metadata.Workspace,
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

func (configurator *StepsConfigurator) GetInternetGatewayV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InternetGatewaySpec],
) *schema.InternetGateway {
	responseExpects.Metadata.Verb = http.MethodGet
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
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
			expectedResourceStatus: responseExpects.ResourceStatus,
		},
	)
}

func (configurator *StepsConfigurator) ListInternetGatewayV1Step(
	stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference, opts *secapi.ListOptions,
) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "ListInternetGateway", wref.Name)
		iter, err := api.ListInternetGatewaysWithOptions(configurator.t.Context(), secapi.WorkspacePath{Tenant: wref.Tenant, Workspace: wref.Workspace}, opts)
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) WatchInternetGatewayUntilDeletedV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		watchWorkspaceResourceUntilDeletedStep(configurator.t, configurator.suite,
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
	})
}

func (configurator *StepsConfigurator) DeleteInternetGatewayV1Step(stepName string, api secapi.NetworkV1, resource *schema.InternetGateway) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetNetworkV1StepParams(sCtx, "DeleteInternetGateway", resource.Metadata.Workspace)

		err := api.DeleteInternetGateway(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Route Table

func (configurator *StepsConfigurator) CreateOrUpdateRouteTableV1Step(stepName string, api secapi.NetworkV1, resource *schema.RouteTable,
	responseExpects ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	createOrUpdateNetworkResourceStep(configurator.t, configurator.suite,
		createOrUpdateNetworkResourceParams[schema.RouteTable, schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec, schema.RouteTableStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  constants.CreateOrUpdateRouteTableOperation,
			workspace:      resource.Metadata.Workspace,
			network:        resource.Metadata.Network,
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

func (configurator *StepsConfigurator) GetRouteTableV1Step(stepName string, api secapi.NetworkV1, nref secapi.NetworkReference,
	responseExpects ResponseExpectsWithCondition[schema.RegionalNetworkResourceMetadata, schema.RouteTableSpec],
) *schema.RouteTable {
	responseExpects.Metadata.Verb = http.MethodGet
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
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
			expectedResourceStatus: responseExpects.ResourceStatus,
		},
	)
}

func (configurator *StepsConfigurator) ListRouteTableV1Step(
	stepName string, api secapi.NetworkV1, nref secapi.NetworkReference, opts *secapi.ListOptions,
) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "ListRouteTable", nref.Name)
		iter, err := api.ListRouteTablesWithOptions(configurator.t.Context(), secapi.NetworkPath{Tenant: nref.Tenant, Workspace: nref.Workspace, Network: nref.Network}, opts)
		requireNoError(sCtx, err)
		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) WatchRouteTableUntilDeletedV1Step(stepName string, api secapi.NetworkV1, nref secapi.NetworkReference) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		watchNetworkResourceUntilDeletedStep(configurator.t, configurator.suite,
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
	})
}

func (configurator *StepsConfigurator) DeleteRouteTableV1Step(stepName string, api secapi.NetworkV1, resource *schema.RouteTable) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetNetworkV1StepParams(sCtx, "DeleteRouteTable", resource.Metadata.Workspace)

		err := api.DeleteRouteTable(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Subnet

func (configurator *StepsConfigurator) CreateOrUpdateSubnetV1Step(stepName string, api secapi.NetworkV1, resource *schema.Subnet,
	responseExpects ResponseExpects[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	createOrUpdateNetworkResourceStep(configurator.t, configurator.suite,
		createOrUpdateNetworkResourceParams[schema.Subnet, schema.RegionalNetworkResourceMetadata, schema.SubnetSpec, schema.SubnetStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkNetworkV1StepParams,
			operationName:  constants.CreateOrUpdateSubnetOperation,
			workspace:      resource.Metadata.Workspace,
			network:        resource.Metadata.Network,
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

func (configurator *StepsConfigurator) GetSubnetV1Step(stepName string, api secapi.NetworkV1, nref secapi.NetworkReference,
	responseExpects ResponseExpectsWithCondition[schema.RegionalNetworkResourceMetadata, schema.SubnetSpec],
) *schema.Subnet {
	responseExpects.Metadata.Verb = http.MethodGet
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
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
			expectedResourceStatus: responseExpects.ResourceStatus,
		},
	)
}

func (configurator *StepsConfigurator) ListSubnetV1Step(
	stepName string, api secapi.NetworkV1, nref secapi.NetworkReference, opts *secapi.ListOptions,
) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "ListSubnet", nref.Name)
		iter, err := api.ListSubnetsWithOptions(configurator.t.Context(), secapi.NetworkPath{Tenant: nref.Tenant, Workspace: nref.Workspace, Network: nref.Network}, opts)
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) WatchSubnetUntilDeletedV1Step(stepName string, api secapi.NetworkV1, nref secapi.NetworkReference) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		watchNetworkResourceUntilDeletedStep(configurator.t, configurator.suite,
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
	})
}

func (configurator *StepsConfigurator) DeleteSubnetV1Step(stepName string, api secapi.NetworkV1, resource *schema.Subnet) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetNetworkV1StepParams(sCtx, "DeleteSubnet", resource.Metadata.Workspace)

		err := api.DeleteSubnet(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Public Ip

func (configurator *StepsConfigurator) CreateOrUpdatePublicIpV1Step(stepName string, api secapi.NetworkV1, resource *schema.PublicIp,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.PublicIp, schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec, schema.PublicIpStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.CreateOrUpdatePublicIpOperation,
			workspace:      resource.Metadata.Workspace,
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

func (configurator *StepsConfigurator) GetPublicIpV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.PublicIpSpec],
) *schema.PublicIp {
	responseExpects.Metadata.Verb = http.MethodGet
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
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
			expectedResourceStatus: responseExpects.ResourceStatus,
		},
	)
}

func (configurator *StepsConfigurator) ListPublicIpV1Step(
	stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference, opts *secapi.ListOptions,
) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "ListPublicIp", wref.Name)
		iter, err := api.ListPublicIpsWithOptions(configurator.t.Context(), secapi.WorkspacePath{Tenant: wref.Tenant, Workspace: wref.Workspace}, opts)
		requireNoError(sCtx, err)

		// Iterate through all items
		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) WatchPublicIpUntilDeletedV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		watchWorkspaceResourceUntilDeletedStep(configurator.t, configurator.suite,
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
	})
}

func (configurator *StepsConfigurator) DeletePublicIpV1Step(stepName string, api secapi.NetworkV1, resource *schema.PublicIp) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetNetworkV1StepParams(sCtx, "DeletePublicIp", resource.Metadata.Workspace)

		err := api.DeletePublicIp(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Nic

func (configurator *StepsConfigurator) CreateOrUpdateNicV1Step(stepName string, api secapi.NetworkV1, resource *schema.Nic,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.Nic, schema.RegionalWorkspaceResourceMetadata, schema.NicSpec, schema.NicStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.CreateOrUpdateNicOperation,
			workspace:      resource.Metadata.Workspace,
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

func (configurator *StepsConfigurator) GetNicV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.NicSpec],
) *schema.Nic {
	responseExpects.Metadata.Verb = http.MethodGet
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
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
			expectedResourceStatus: responseExpects.ResourceStatus,
		},
	)
}

func (configurator *StepsConfigurator) ListNicV1Step(
	stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference, opts *secapi.ListOptions,
) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "ListNic", wref.Name)
		iter, err := api.ListNicsWithOptions(configurator.t.Context(), secapi.WorkspacePath{Tenant: wref.Tenant, Workspace: wref.Workspace}, opts)
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) WatchNicUntilDeletedV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		watchWorkspaceResourceUntilDeletedStep(configurator.t, configurator.suite,
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
	})
}

func (configurator *StepsConfigurator) DeleteNicV1Step(stepName string, api secapi.NetworkV1, resource *schema.Nic) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetNetworkV1StepParams(sCtx, "DeleteNic", resource.Metadata.Workspace)

		err := api.DeleteNic(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Security Group Rule

func (configurator *StepsConfigurator) CreateOrUpdateSecurityGroupRuleV1Step(stepName string, api secapi.NetworkV1, resource *schema.SecurityGroupRule,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.SecurityGroupRule, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec, schema.Status]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.CreateOrUpdateSecurityGroupRuleOperation,
			workspace:      resource.Metadata.Workspace,
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

func (configurator *StepsConfigurator) GetSecurityGroupRuleV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupRuleSpec],
) *schema.SecurityGroupRule {
	responseExpects.Metadata.Verb = http.MethodGet
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
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
			expectedResourceStatus: responseExpects.ResourceStatus,
		},
	)
}

func (configurator *StepsConfigurator) ListSecurityGroupRuleV1Step(
	stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference, opts *secapi.ListOptions,
) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "ListSecurityGroupRule", wref.Name)
		iter, err := api.ListSecurityGroupRulesWithOptions(configurator.t.Context(), secapi.WorkspacePath{Tenant: wref.Tenant, Workspace: wref.Workspace}, opts)
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) WatchSecurityGroupRuleUntilDeletedV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		watchWorkspaceResourceUntilDeletedStep(configurator.t, configurator.suite,
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
	})
}

func (configurator *StepsConfigurator) DeleteSecurityGroupRuleV1Step(stepName string, api secapi.NetworkV1, resource *schema.SecurityGroupRule) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetNetworkV1StepParams(sCtx, "DeleteSecurityGroupRule", resource.Metadata.Workspace)

		err := api.DeleteSecurityGroupRule(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Security Group

func (configurator *StepsConfigurator) CreateOrUpdateSecurityGroupV1Step(stepName string, api secapi.NetworkV1, resource *schema.SecurityGroup,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.SecurityGroup, schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec, schema.SecurityGroupStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetNetworkV1StepParams,
			operationName:  constants.CreateOrUpdateSecurityGroupOperation,
			workspace:      resource.Metadata.Workspace,
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

func (configurator *StepsConfigurator) GetSecurityGroupV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.SecurityGroupSpec],
) *schema.SecurityGroup {
	responseExpects.Metadata.Verb = http.MethodGet
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
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
			expectedResourceStatus: responseExpects.ResourceStatus,
		},
	)
}

func (configurator *StepsConfigurator) ListSecurityGroupV1Step(
	stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference, opts *secapi.ListOptions,
) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "ListSecurityGroup", wref.Name)
		iter, err := api.ListSecurityGroupsWithOptions(configurator.t.Context(), secapi.WorkspacePath{Tenant: wref.Tenant, Workspace: wref.Workspace}, opts)
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) WatchSecurityGroupUntilDeletedV1Step(stepName string, api secapi.NetworkV1, wref secapi.WorkspaceReference) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		watchWorkspaceResourceUntilDeletedStep(configurator.t, configurator.suite,
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
	})
}

func (configurator *StepsConfigurator) DeleteSecurityGroupV1Step(stepName string, api secapi.NetworkV1, resource *schema.SecurityGroup) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetNetworkV1StepParams(sCtx, "DeleteSecurityGroup", resource.Metadata.Workspace)

		err := api.DeleteSecurityGroup(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}
