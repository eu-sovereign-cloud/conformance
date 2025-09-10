package secatest

import (
	"context"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	compute "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.compute.v1"
	storage "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.storage.v1"
	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/google/uuid"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type ComputeV1TestSuite struct {
	regionalTestSuite

	availableZones []string
	instanceSkus   []string
	storageSkus    []string
}

func (suite *ComputeV1TestSuite) generateLifecycleParams() *secalib.ComputeLifeCycleParamsV1 {
	workspaceName := secalib.GenerateWorkspaceName()
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]
	blockStorageName := secalib.GenerateBlockStorageName()
	instanceSkuName := suite.instanceSkus[rand.Intn(len(suite.instanceSkus))]
	instanceName := secalib.GenerateInstanceName()

	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)
	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName)
	instanceSkuRef := secalib.GenerateSkuRef(instanceSkuName)

	// Random data
	blockStorageSize := secalib.GenerateBlockStorageSize()
	initialInstanceZone := suite.availableZones[rand.Intn(len(suite.availableZones))]
	updatedInstanceZone := suite.availableZones[rand.Intn(len(suite.availableZones))]

	return &secalib.ComputeLifeCycleParamsV1{
		Workspace: &secalib.ResourceParams[secalib.WorkspaceSpecV1]{
			Name: workspaceName,
			InitialSpec: &secalib.WorkspaceSpecV1{
				Labels: &[]secalib.Label{
					{
						Name:  secalib.EnvLabel,
						Value: secalib.EnvDevelopmentLabel,
					},
				},
			},
		},
		BlockStorage: &secalib.ResourceParams[secalib.BlockStorageSpecV1]{
			Name: blockStorageName,
			InitialSpec: &secalib.BlockStorageSpecV1{
				SkuRef: storageSkuRef,
				SizeGB: blockStorageSize,
			},
		},
		Instance: &secalib.ResourceParams[secalib.InstanceSpecV1]{
			Name: instanceName,
			InitialSpec: &secalib.InstanceSpecV1{
				SkuRef:        instanceSkuRef,
				Zone:          initialInstanceZone,
				BootDeviceRef: blockStorageRef,
			},
			UpdatedSpec: &secalib.InstanceSpecV1{
				SkuRef:        instanceSkuRef,
				Zone:          updatedInstanceZone,
				BootDeviceRef: blockStorageRef,
			},
		},
	}
}

func (suite *ComputeV1TestSuite) TestComputeLifeCycleV1(t provider.T) {
	slog.Info("Starting Compute Lifecycle Test")

	t.Title("Compute Lifecycle Test")
	configureTags(t, secalib.ComputeProviderV1, secalib.InstanceKind)

	params := suite.generateLifecycleParams()
	instanceResource := secalib.GenerateInstanceResource(suite.tenant, params.Workspace.Name, params.Instance.Name)

	// Setup mock, if configured to use
	if suite.isMockEnabled() {
		scenarios := mock.NewComputeScenariosV1(suite.authToken, suite.tenant, suite.region, suite.mockServerURL)

		id := uuid.New().String()
		wm, err := scenarios.ConfigureLifecycleScenario(id, params)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	ctx := context.Background()	
	var err error

	// Workspace
	var workResp *workspace.Workspace

	t.WithNewStep("Create workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "CreateOrUpdateWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Workspace.Name,
		}
		ws := &workspace.Workspace{}
		workResp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, tref, ws, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, workResp)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*workResp.Status.State)},
		)
	})

	t.WithNewStep("Get created workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "GetWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Workspace.Name,
		}
		workResp, err = suite.client.WorkspaceV1.GetWorkspace(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, workResp)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*workResp.Status.State)},
		)
	})

	// Block storage
	var blockStorageResp *storage.BlockStorage

	t.WithNewStep("Create block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateBlockStorage", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.BlockStorage.Name,
		}
		bo := &storage.BlockStorage{
			Spec: storage.BlockStorageSpec{
				SizeGB: params.BlockStorage.InitialSpec.SizeGB,
				SkuRef: params.BlockStorage.InitialSpec.SkuRef,
			},
		}
		blockStorageResp, err = suite.client.StorageV1.CreateOrUpdateBlockStorage(ctx, wref, bo, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockStorageResp)
	})

	t.WithNewStep("Get created block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetBlockStorage", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.BlockStorage.Name,
		}
		blockStorageResp, err = suite.client.StorageV1.GetBlockStorage(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockStorageResp)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*blockStorageResp.Status.State)},
		)
	})

	// Instance
	var instanceResp *compute.Instance
	var expectedMeta *secalib.Metadata

	t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "CreateOrUpdateInstance", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Instance.Name,
		}
		inst := &compute.Instance{
			Spec: compute.InstanceSpec{
				SkuRef: params.Instance.InitialSpec.SkuRef,
				Zone:   params.Instance.InitialSpec.Zone,
			},
		}
		inst.Spec.BootVolume.DeviceRef = params.Instance.InitialSpec.BootDeviceRef

		instanceResp, err = suite.client.ComputeV1.CreateOrUpdateInstance(ctx, wref, inst, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMeta = &secalib.Metadata{
			Name:       params.Instance.Name,
			Provider:   secalib.ComputeProviderV1,
			Resource:   instanceResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyInstanceZonalMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		verifyInstanceSpecStep(sCtx, params.Instance.InitialSpec, &instanceResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	t.WithNewStep("Get created instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Instance.Name,
		}
		instanceResp, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMeta.Verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		verifyInstanceSpecStep(sCtx, params.Instance.InitialSpec, &instanceResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	t.WithNewStep("Update instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "CreateOrUpdateInstance", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Instance.Name,
		}
		instanceResp, err = suite.client.ComputeV1.CreateOrUpdateInstance(ctx, wref, instanceResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMeta.Verb = http.MethodPut
		verifyInstanceZonalMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		verifyInstanceSpecStep(sCtx, params.Instance.UpdatedSpec, &instanceResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	t.WithNewStep("Get updated instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Instance.Name,
		}
		instanceResp, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMeta.Verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		verifyInstanceSpecStep(sCtx, params.Instance.UpdatedSpec, &instanceResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	t.WithNewStep("Stop instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "StopInstance", params.Workspace.Name)

		err = suite.client.ComputeV1.StopInstance(ctx, instanceResp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get stopped instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Instance.Name,
		}
		instanceResp, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMeta.Verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.SuspendedStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	t.WithNewStep("Start instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "StartInstance", params.Workspace.Name)

		err = suite.client.ComputeV1.StartInstance(ctx, instanceResp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get started instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Instance.Name,
		}
		instanceResp, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMeta.Verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	t.WithNewStep("Restart instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "RestartInstance", params.Workspace.Name)

		err = suite.client.ComputeV1.RestartInstance(ctx, instanceResp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get restarted instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Instance.Name,
		}
		instanceResp, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMeta.Verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	t.WithNewStep("Delete instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "DeleteInstance", params.Workspace.Name)

		err = suite.client.ComputeV1.DeleteInstance(ctx, instanceResp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.Instance.Name,
		}
		_, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	slog.Info("Finishing Compute Lifecycle Test")
}

func (suite *ComputeV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}

func verifyInstanceZonalMetadataStep(ctx provider.StepCtx, expected *secalib.Metadata, metadata *compute.ZonalResourceMetadata) {
	actualMetadata := &secalib.Metadata{
		Name:       metadata.Name,
		Provider:   metadata.Provider,
		Verb:       metadata.Verb,
		Resource:   metadata.Resource,
		ApiVersion: metadata.ApiVersion,
		Kind:       string(metadata.Kind),
		Tenant:     metadata.Tenant,
		Workspace:  *metadata.Workspace,
		Region:     metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}

func verifyInstanceSpecStep(ctx provider.StepCtx, expected *secalib.InstanceSpecV1, actual *compute.InstanceSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.SkuRef, actual.SkuRef, "SkuRef should match expected")
		stepCtx.Require().Equal(expected.Zone, actual.Zone, "Zone should match expected")
		stepCtx.Require().Equal(expected.BootDeviceRef, actual.BootVolume.DeviceRef, "BootVolume.DeviceRef should match expected")
	})
}
