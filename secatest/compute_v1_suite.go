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
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type ComputeV1TestSuite struct {
	regionalTestSuite

	availableZones []string
	instanceSkus   []string
	storageSkus    []string
}

func (suite *ComputeV1TestSuite) TestComputeV1(t provider.T) {
	slog.Info("Starting Compute Lifecycle Test")

	t.Title("Compute Lifecycle Test")
	configureTags(t, secalib.ComputeProviderV1, secalib.InstanceKind)

	// Select skus
	instanceSkuName := suite.instanceSkus[rand.Intn(len(suite.instanceSkus))]
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]

	// Select zones
	initialInstanceZone := suite.availableZones[rand.Intn(len(suite.availableZones))]
	updatedInstanceZone := suite.availableZones[rand.Intn(len(suite.availableZones))]

	// TODO Dynamically create before the scenario
	workspaceName := secalib.GenerateWorkspaceName()

	// Generate scenario data
	instanceSkuRef := secalib.GenerateSkuRef(instanceSkuName)

	instanceName := secalib.GenerateInstanceName()
	instanceResource := secalib.GenerateInstanceResource(suite.tenant, workspaceName, instanceName)

	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)

	blockStorageName := secalib.GenerateBlockStorageName()
	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName)

	blockStorageSize := secalib.GenerateBlockStorageSize()

	// Setup mock, if configured to use
	if suite.isMockEnabled() {
		wm, err := mock.CreateComputeLifecycleScenarioV1("Compute Lifecycle",
			mock.ComputeParamsV1{
				Params: &mock.Params{
					MockURL:   suite.mockServerURL,
					AuthToken: suite.authToken,
					Tenant:    suite.tenant,
					Workspace: workspaceName,
					Region:    suite.region,
				},
				BlockStorage: &mock.ResourceParams[secalib.BlockStorageSpecV1]{
					Name: blockStorageName,
					InitialSpec: &secalib.BlockStorageSpecV1{
						SkuRef: storageSkuRef,
						SizeGB: blockStorageSize,
					},
				},
				Instance: &mock.ResourceParams[secalib.InstanceSpecV1]{
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
			})
		if err != nil {
			slog.Error("Failed to create compute scenario", "error", err)
			return
		}
		suite.mockClient = wm
	}

	ctx := context.Background()
	var blockStorageResp *storage.BlockStorage
	var instanceResp *compute.Instance
	var err error

	// Step 1
	t.WithNewStep("Create block storage", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateBlockStorage",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		bo := &storage.BlockStorage{
			Spec: storage.BlockStorageSpec{
				SizeGB: blockStorageSize,
				SkuRef: storageSkuRef,
			},
		}
		blockStorageResp, err = suite.client.StorageV1.CreateOrUpdateBlockStorage(ctx, wref, bo, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockStorageResp)
	})

	// Step 2
	t.WithNewStep("Get created block storage", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetBlockStorage",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		blockStorageResp, err = suite.client.StorageV1.GetBlockStorage(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockStorageResp)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*blockStorageResp.Status.State)},
		)
	})

	var expectedMetadata *secalib.Metadata
	var expectedSpec *secalib.InstanceSpecV1

	// Step 3
	t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		inst := &compute.Instance{
			Spec: compute.InstanceSpec{
				SkuRef: instanceSkuRef,
				Zone:   initialInstanceZone,
			},
		}
		inst.Spec.BootVolume.DeviceRef = blockStorageRef

		instanceResp, err = suite.client.ComputeV1.CreateOrUpdateInstance(ctx, wref, inst, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMetadata = &secalib.Metadata{
			Name:       instanceName,
			Provider:   secalib.ComputeProviderV1,
			Resource:   instanceResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.InstanceKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyInstanceZonalMetadataStep(sCtx, expectedMetadata, instanceResp.Metadata)

		expectedSpec = &secalib.InstanceSpecV1{
			SkuRef:        instanceSkuRef,
			Zone:          initialInstanceZone,
			BootDeviceRef: blockStorageRef,
		}
		verifyInstanceSpecStep(sCtx, expectedSpec, &instanceResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	// Step 4
	t.WithNewStep("Get created instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		instanceResp, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMetadata.Verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedMetadata, instanceResp.Metadata)

		verifyInstanceSpecStep(sCtx, expectedSpec, &instanceResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	// Step 5
	t.WithNewStep("Update instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		instanceResp, err = suite.client.ComputeV1.CreateOrUpdateInstance(ctx, wref, instanceResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMetadata.Verb = http.MethodPut
		verifyInstanceZonalMetadataStep(sCtx, expectedMetadata, instanceResp.Metadata)

		expectedSpec.Zone = updatedInstanceZone
		verifyInstanceSpecStep(sCtx, expectedSpec, &instanceResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	// Step 6
	t.WithNewStep("Get updated instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		instanceResp, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMetadata.Verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedMetadata, instanceResp.Metadata)

		verifyInstanceSpecStep(sCtx, expectedSpec, &instanceResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	// Step 7
	t.WithNewStep("Stop instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "StopInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		err = suite.client.ComputeV1.StopInstance(ctx, instanceResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 8
	t.WithNewStep("Get stopped instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		instanceResp, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMetadata.Verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedMetadata, instanceResp.Metadata)

		verifyInstanceSpecStep(sCtx, expectedSpec, &instanceResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.SuspendedStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	// Step 9
	t.WithNewStep("Start instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "StartInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		err = suite.client.ComputeV1.StartInstance(ctx, instanceResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 10
	t.WithNewStep("Get started instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		instanceResp, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMetadata.Verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedMetadata, instanceResp.Metadata)

		verifyInstanceSpecStep(sCtx, expectedSpec, &instanceResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	// Step 11
	t.WithNewStep("Restart instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "Restart",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		err = suite.client.ComputeV1.RestartInstance(ctx, instanceResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 12
	t.WithNewStep("Get restarted instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		instanceResp, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMetadata.Verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedMetadata, instanceResp.Metadata)

		verifyInstanceSpecStep(sCtx, expectedSpec, &instanceResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*instanceResp.Status.State)},
		)
	})

	// Step 13
	t.WithNewStep("Delete instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		err = suite.client.ComputeV1.DeleteInstance(ctx, instanceResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 14
	t.WithNewStep("Get deleted instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
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
		stepCtx.WithNewParameters(
			"expected_sizeGB", expected.SkuRef,
			"actual_sizeGB", actual.SkuRef,

			"expected_skuRef", expected.Zone,
			"actual_skuRef", actual.Zone,

			"expected_skuRef", expected.BootDeviceRef,
			"actual_skuRef", actual.BootVolume.DeviceRef,
		)
		stepCtx.Require().Equal(expected.SkuRef, actual.SkuRef, "SkuRef should match expected")
		stepCtx.Require().Equal(expected.Zone, actual.Zone, "Zone should match expected")
		stepCtx.Require().Equal(expected.BootDeviceRef, actual.BootVolume.DeviceRef, "BootVolume.DeviceRef should match expected")
	})
}
