package secatest

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
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
	configureTags(t, computeV1Provider, instanceSkuKind, instanceSkuKind)

	// TODO Export to configuration or create before the test
	workspaceName := "workspace-1"
	skuProvider := "seca"
	instanceSkuTier := "DXS"
	instanceSkuRAM := 2048
	instanceSkuVCPU := 2
	instanceZoneA := "zone-a"
	instanceZoneB := "zone-b"
	storageSkuTier := "gold"

	// Generate scenario data
	// TODO Create generator functions
	instanceSkuName := fmt.Sprintf("sku-%d", rand.Intn(math.MaxInt32))
	instanceSkuRef := fmt.Sprintf(skuRef, instanceSkuName)

	instanceName := fmt.Sprintf("instance-%d", rand.Intn(math.MaxInt32))
	instanceResource := fmt.Sprintf(instanceResource, suite.tenant, workspaceName, instanceName)

	storageSkuName := fmt.Sprintf("sku-%d", rand.Intn(math.MaxInt32))
	storageSkuRef := fmt.Sprintf(skuRef, storageSkuName)

	blockStorageName := fmt.Sprintf("disk-%d", rand.Intn(math.MaxInt32))
	blockStorageRef := fmt.Sprintf(blockStoragesRef, blockStorageName)

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
					Iops:          storageSkuIops500,
					StorageType:   storageTypeLocalEphemeral,
					MinVolumeSize: storageSkuMinVolumeSize10,
				},
				BlockStorage: mock.BlockStorageParamsV1{
					Name:          blockStorageName,
					SkuRef:        storageSkuRef,
					CreatedSizeGB: blockStorageSizeGB100,
				},
				InstanceSku: mock.InstanceSkuParamsV1{
					Name:         instanceSkuName,
					Architecture: cpuArchitectureAmd64,
					Provider:     skuProvider,
					Tier:         instanceSkuTier,
					RAM:          instanceSkuRAM,
					VCPU:         instanceSkuVCPU,
				},
				Instance: mock.InstanceParamsV1{
					Name:          instanceName,
					SkuRef:        fmt.Sprintf(skuRef, instanceSkuName),
					CreatedZone:   instanceZoneA,
					UpdatedZone:   instanceZoneB,
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

		bo := &storage.BlockStorage{
			Metadata: &storage.ZonalResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: &workspaceName,
				Name:      blockStorageName,
			},
			Spec: storage.BlockStorageSpec{
				SizeGB: blockStorageSizeGB100,
				SkuRef: storageSkuRef,
			},
		}
		blockStorageResp, err = suite.client.StorageV1.CreateOrUpdateBlockStorage(ctx, bo)
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
			expectedState: activeStatusState,
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
			architecture: cpuArchitectureAmd64,
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

	expectedInstanceMetadata := verifyRegionalMetadataStepParams{
		name:       instanceName,
		provider:   computeV1Provider,
		resource:   instanceResource,
		apiVersion: version1,
		kind:       instanceKind,
		tenant:     suite.tenant,
		region:     suite.region,
	}

	// Step 5
	t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateInstance",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		inst := &compute.Instance{
			Metadata: &compute.ZonalResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: &workspaceName,
				Name:      instanceName,
			},
			Spec: compute.InstanceSpec{
				SkuRef: instanceSkuRef,
				Zone:   instanceZoneA,
			},
		}
		inst.Spec.BootVolume.DeviceRef = blockStorageRef

		instanceResp, err = suite.client.ComputeV1.CreateOrUpdateInstance(ctx, inst)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedInstanceMetadata.verb = http.MethodPut
		verifyInstanceZonalMetadataStep(sCtx, expectedInstanceMetadata, instanceResp.Metadata)

		verifySpecParams := verifyInstanceSpecStepParams{
			skuRef:        instanceSkuRef,
			zone:          instanceZoneA,
			bootDeviceRef: blockStorageRef,
		}
		verifyInstanceSpecStep(sCtx, verifySpecParams, instanceResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: creatingStatusState,
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

		expectedInstanceMetadata.verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedInstanceMetadata, instanceResp.Metadata)

		verifySpecParams := verifyInstanceSpecStepParams{
			skuRef:        instanceSkuRef,
			zone:          instanceZoneA,
			bootDeviceRef: blockStorageRef,
		}
		verifyInstanceSpecStep(sCtx, verifySpecParams, instanceResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: activeStatusState,
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

		instanceResp.Spec.Zone = instanceZoneB
		instanceResp, err = suite.client.ComputeV1.CreateOrUpdateInstance(ctx, instanceResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedInstanceMetadata.verb = http.MethodPut
		verifyInstanceZonalMetadataStep(sCtx, expectedInstanceMetadata, instanceResp.Metadata)

		verifySpecParams := verifyInstanceSpecStepParams{
			skuRef:        instanceSkuRef,
			zone:          instanceZoneB,
			bootDeviceRef: blockStorageRef,
		}
		verifyInstanceSpecStep(sCtx, verifySpecParams, instanceResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: updatingStatusState,
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

		expectedInstanceMetadata.verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedInstanceMetadata, instanceResp.Metadata)

		verifySpecParams := verifyInstanceSpecStepParams{
			skuRef:        instanceSkuRef,
			zone:          instanceZoneB,
			bootDeviceRef: blockStorageRef,
		}
		verifyInstanceSpecStep(sCtx, verifySpecParams, instanceResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: activeStatusState,
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

		err = suite.client.ComputeV1.StopInstance(ctx, instanceResp)
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

		expectedInstanceMetadata.verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedInstanceMetadata, instanceResp.Metadata)

		verifySpecParams := verifyInstanceSpecStepParams{
			skuRef:        instanceSkuRef,
			zone:          instanceZoneB,
			bootDeviceRef: blockStorageRef,
		}
		verifyInstanceSpecStep(sCtx, verifySpecParams, instanceResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: suspendedStatusState,
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

		err = suite.client.ComputeV1.StartInstance(ctx, instanceResp)
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

		expectedInstanceMetadata.verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedInstanceMetadata, instanceResp.Metadata)

		verifySpecParams := verifyInstanceSpecStepParams{
			skuRef:        instanceSkuRef,
			zone:          instanceZoneB,
			bootDeviceRef: blockStorageRef,
		}
		verifyInstanceSpecStep(sCtx, verifySpecParams, instanceResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: activeStatusState,
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

		err = suite.client.ComputeV1.RestartInstance(ctx, instanceResp)
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

		expectedInstanceMetadata.verb = http.MethodGet
		verifyInstanceZonalMetadataStep(sCtx, expectedInstanceMetadata, instanceResp.Metadata)

		verifySpecParams := verifyInstanceSpecStepParams{
			skuRef:        instanceSkuRef,
			zone:          instanceZoneB,
			bootDeviceRef: blockStorageRef,
		}
		verifyInstanceSpecStep(sCtx, verifySpecParams, instanceResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: activeStatusState,
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

		err = suite.client.ComputeV1.DeleteInstance(ctx, instanceResp)
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
		instanceResp, err = suite.client.ComputeV1.GetInstance(ctx, wref)
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
			"actual_architecture", labels[architectureLabel],

			"expected_provider", expected.provider,
			"actual_provider", labels[providerLabel],

			"expected_tier", expected.tier,
			"actual_tier", labels[tierLabel],
		)
		stepCtx.Require().Equal(expected.architecture, labels[architectureLabel], "Architecture should match expected")
		stepCtx.Require().Equal(expected.provider, labels[providerLabel], "Provider should match expected")
		stepCtx.Require().Equal(expected.tier, labels[tierLabel], "Tier should match expected")
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
