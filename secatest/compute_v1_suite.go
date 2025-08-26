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
}

type verifyInstanceSkuLabelsStepParams struct {
	architecture string
	provider     string
	tier         string
}

type verifyInstanceSkuSpecStepParams struct {
	ram  int
	vCPU int
}

type verifyInstanceSpecStepParams struct {
	skuRef        string
	zone          string
	bootDeviceRef string
}

func (suite *ComputeV1TestSuite) TestComputeV1(t provider.T) {
	t.Title("Compute Lifecycle Test")
	configureTags(t, secalib.ComputeProviderV1, secalib.InstanceSkuKind, secalib.InstanceKind)

	// TODO Export configuration
	instanceZones := []string{"zone-a", "zone-b"}

	// TODO Export to mock configuration
	skuProvider := "seca"
	instanceSkuTier := "DXS"
	instanceSkuRAM := 2048
	instanceSkuVCPU := 2
	storageSkuTier := "gold"

	// Select zones
	instanceZone1 := instanceZones[rand.Intn(len(instanceZones))]
	instanceZone2 := instanceZones[rand.Intn(len(instanceZones))]

	// TODO Dynamically create before the scenario
	workspaceName := secalib.GenerateWorkspaceName()

	// Generate scenario data
	instanceSkuName := secalib.GenerateSkuName()
	instanceSkuRef := secalib.GenerateSkuRef(instanceSkuName)

	instanceName := secalib.GenerateInstanceName()
	instanceResource := secalib.GenerateInstanceResource(suite.tenant, workspaceName, instanceName)

	storageSkuName := secalib.GenerateSkuName()
	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)

	blockStorageName := secalib.GenerateBlockStorageName()
	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName)

	storageSkuIops := secalib.GenerateStorageSkuIops()
	blockStorageSize := secalib.GenerateBlockStorageSize()
	storageSkuMinVolumeSize := secalib.GenerateStorageSkuMinVolumeSize(blockStorageSize)

	// Setup mock, if configured to use
	if suite.isMockEnabled() {
		wm, err := mock.CreateComputeLifecycleScenarioV1("Compute Lifecycle",
			mock.ComputeParamsV1{
				Params: mock.Params{
					MockURL:   suite.mockServerURL,
					AuthToken: suite.authToken,
					Tenant:    suite.tenant,
					Workspace: workspaceName,
					Region:    suite.region,
				},
				StorageSku: mock.StorageSkuParamsV1{
					Name:          storageSkuName,
					Provider:      skuProvider,
					Tier:          storageSkuTier,
					Iops:          storageSkuIops,
					StorageType:   secalib.StorageTypeLocalEphemeral,
					MinVolumeSize: storageSkuMinVolumeSize,
				},
				BlockStorage: mock.BlockStorageParamsV1{
					Name:          blockStorageName,
					SkuRef:        storageSkuRef,
					SizeGBInitial: blockStorageSize,
				},
				InstanceSku: mock.InstanceSkuParamsV1{
					Name:         instanceSkuName,
					Architecture: secalib.CpuArchitectureAmd64,
					Provider:     skuProvider,
					Tier:         instanceSkuTier,
					RAM:          instanceSkuRAM,
					VCPU:         instanceSkuVCPU,
				},
				Instance: mock.InstanceParamsV1{
					Name:          instanceName,
					SkuRef:        instanceSkuRef,
					ZoneInitial:   instanceZone1,
					ZoneUpdated:   instanceZone2,
					BootDeviceRef: blockStorageRef,
				},
			})
		if err != nil {
			slog.Error("Failed to create compute scenario", "error", err)
			return
		}
		suite.mockClient = wm
	}

	ctx := context.Background()
	var storageSkuResp *storage.StorageSku
	var blockStorageResp *storage.BlockStorage
	var instanceSkuResp *compute.InstanceSku
	var instanceResp *compute.Instance
	var err error

	// Step 1
	t.WithNewStep("Get storage sku", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetStorageSku",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   storageSkuName,
		}
		storageSkuResp, err = suite.client.StorageV1.GetSku(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, storageSkuResp)
	})

	// Step 2
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

	// Step 3
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

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.ActiveStatusState,
			actualState:   string(*blockStorageResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 4
	t.WithNewStep("Get instance sku", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetInstanceSku",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   instanceSkuName,
		}
		instanceSkuResp, err = suite.client.ComputeV1.GetSku(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceSkuResp)

		verifySkuMetadataStep(sCtx,
			verifySkuMetadataStepParams{name: instanceSkuName},
			verifySkuMetadataStepParams{name: instanceSkuResp.Metadata.Name},
		)

		verifyLabelsParams := verifyInstanceSkuLabelsStepParams{
			architecture: secalib.CpuArchitectureAmd64,
			provider:     skuProvider,
			tier:         instanceSkuTier,
		}
		verifyInstanceSkuLabelsStep(sCtx, verifyLabelsParams, *instanceSkuResp.Labels)

		verifySpecParams := verifyInstanceSkuSpecStepParams{
			ram:  instanceSkuRAM,
			vCPU: instanceSkuVCPU,
		}
		verifyInstanceSkuSpecStep(sCtx, verifySpecParams, instanceSkuResp.Spec)
	})

	var expectedMetadata verifyRegionalMetadataStepParams
	var expectedSpec verifyInstanceSpecStepParams

	// Step 5
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
				Zone:   instanceZone1,
			},
		}
		inst.Spec.BootVolume.DeviceRef = blockStorageRef

		instanceResp, err = suite.client.ComputeV1.CreateOrUpdateInstance(ctx, wref, inst, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMetadata = verifyRegionalMetadataStepParams{
			name:       instanceName,
			provider:   secalib.ComputeProviderV1,
			resource:   instanceResource,
			verb:       http.MethodPut,
			apiVersion: secalib.ApiVersion1,
			kind:       secalib.InstanceKind,
			tenant:     suite.tenant,
			region:     suite.region,
		}
		verifyInstanceZonalMetadataStep(sCtx, expectedMetadata, instanceResp.Metadata)

		expectedSpec = verifyInstanceSpecStepParams{
			skuRef:        instanceSkuRef,
			zone:          instanceZone1,
			bootDeviceRef: blockStorageRef,
		}
		verifyInstanceSpecStep(sCtx, expectedSpec, instanceResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.CreatingStatusState,
			actualState:   string(*instanceResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 6
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

		expectedMetadata.verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedMetadata, instanceResp.Metadata)

		verifyInstanceSpecStep(sCtx, expectedSpec, instanceResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.ActiveStatusState,
			actualState:   string(*instanceResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 7
	t.WithNewStep("Update instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		instanceResp.Spec.Zone = instanceZone2

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		instanceResp, err = suite.client.ComputeV1.CreateOrUpdateInstance(ctx, wref, instanceResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMetadata.verb = http.MethodPut
		verifyInstanceZonalMetadataStep(sCtx, expectedMetadata, instanceResp.Metadata)

		expectedSpec.zone = instanceZone2
		verifyInstanceSpecStep(sCtx, expectedSpec, instanceResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.UpdatingStatusState,
			actualState:   string(*instanceResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 8
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

		expectedMetadata.verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedMetadata, instanceResp.Metadata)

		verifyInstanceSpecStep(sCtx, expectedSpec, instanceResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.ActiveStatusState,
			actualState:   string(*instanceResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 9
	t.WithNewStep("Stop instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "StopInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		err = suite.client.ComputeV1.StopInstance(ctx, instanceResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 10
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

		expectedMetadata.verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedMetadata, instanceResp.Metadata)

		verifyInstanceSpecStep(sCtx, expectedSpec, instanceResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.SuspendedStatusState,
			actualState:   string(*instanceResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 11
	t.WithNewStep("Start instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "StartInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		err = suite.client.ComputeV1.StartInstance(ctx, instanceResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 12
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

		expectedMetadata.verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedMetadata, instanceResp.Metadata)

		verifyInstanceSpecStep(sCtx, expectedSpec, instanceResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.ActiveStatusState,
			actualState:   string(*instanceResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 13
	t.WithNewStep("Restart instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "Restart",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		err = suite.client.ComputeV1.RestartInstance(ctx, instanceResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 14
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

		expectedMetadata.verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedMetadata, instanceResp.Metadata)

		verifyInstanceSpecStep(sCtx, expectedSpec, instanceResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.ActiveStatusState,
			actualState:   string(*instanceResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 15
	t.WithNewStep("Delete instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		err = suite.client.ComputeV1.DeleteInstance(ctx, instanceResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 16
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
}

func (suite *ComputeV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}

func verifyInstanceZonalMetadataStep(ctx provider.StepCtx, expected verifyRegionalMetadataStepParams, metadata *compute.ZonalResourceMetadata) {
	actualMetadata := verifyRegionalMetadataStepParams{
		name:       metadata.Name,
		provider:   metadata.Provider,
		verb:       metadata.Verb,
		resource:   metadata.Resource,
		apiVersion: metadata.ApiVersion,
		kind:       string(metadata.Kind),
		tenant:     metadata.Tenant,
		workspace:  *metadata.Workspace,
		region:     metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}

func verifyInstanceSkuLabelsStep(ctx provider.StepCtx, expected verifyInstanceSkuLabelsStepParams, labels map[string]string) {
	ctx.WithNewStep("Verify labels", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_architecture", expected.architecture,
			"actual_architecture", labels[secalib.ArchitectureLabel],

			"expected_provider", expected.provider,
			"actual_provider", labels[secalib.ProviderLabel],

			"expected_tier", expected.tier,
			"actual_tier", labels[secalib.TierLabel],
		)
		stepCtx.Require().Equal(expected.architecture, labels[secalib.ArchitectureLabel], "Architecture should match expected")
		stepCtx.Require().Equal(expected.provider, labels[secalib.ProviderLabel], "Provider should match expected")
		stepCtx.Require().Equal(expected.tier, labels[secalib.TierLabel], "Tier should match expected")
	})
}

func verifyInstanceSkuSpecStep(ctx provider.StepCtx, expected verifyInstanceSkuSpecStepParams, actual *compute.InstanceSkuSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_ram", expected.ram,
			"actual_ram", actual.Ram,

			"expected_vCPU", expected.vCPU,
			"actual_vCPU", actual.VCPU,
		)
		stepCtx.Require().Equal(expected.ram, actual.Ram, "Ram should match expected")
		stepCtx.Require().Equal(expected.vCPU, actual.VCPU, "vCPU should match expected")
	})
}

func verifyInstanceSpecStep(ctx provider.StepCtx, expected verifyInstanceSpecStepParams, actual compute.InstanceSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_sizeGB", expected.skuRef,
			"actual_sizeGB", actual.SkuRef,

			"expected_skuRef", expected.zone,
			"actual_skuRef", actual.Zone,

			"expected_skuRef", expected.bootDeviceRef,
			"actual_skuRef", actual.BootVolume.DeviceRef,
		)
		stepCtx.Require().Equal(expected.skuRef, actual.SkuRef, "SkuRef should match expected")
		stepCtx.Require().Equal(expected.zone, actual.Zone, "Zone should match expected")
		stepCtx.Require().Equal(expected.bootDeviceRef, actual.BootVolume.DeviceRef, "BootVolume.DeviceRef should match expected")
	})
}
