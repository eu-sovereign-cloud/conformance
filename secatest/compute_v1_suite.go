package secatest

import (
	"context"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type ComputeV1TestSuite struct {
	regionalTestSuite

	availableZones []string
	instanceSkus   []string
	storageSkus    []string
}

func (suite *ComputeV1TestSuite) TestSuite(t provider.T) {
	ctx := context.Background()
	var err error
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t, secalib.ComputeProviderV1, secalib.InstanceKind)

	// Select skus
	instanceSkuName := suite.instanceSkus[rand.Intn(len(suite.instanceSkus))]
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]

	// Select zones
	initialInstanceZone := suite.availableZones[rand.Intn(len(suite.availableZones))]
	updatedInstanceZone := suite.availableZones[rand.Intn(len(suite.availableZones))]

	// Generate scenario data
	workspaceName := secalib.GenerateWorkspaceName()

	instanceSkuRef := secalib.GenerateSkuRef(instanceSkuName)
	instanceSkuRefObj, err := secapi.BuildReferenceFromURN(instanceSkuRef)
	if err != nil {
		t.Fatalf("Failed to build instanceSkuRef to URN: %v", err)
	}

	instanceName := secalib.GenerateInstanceName()
	instanceResource := secalib.GenerateInstanceResource(suite.tenant, workspaceName, instanceName)

	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)
	storageSkuRefObj, err := secapi.BuildReferenceFromURN(storageSkuRef)
	if err != nil {
		t.Fatalf("Failed to build storageSkuRef to URN: %v", err)
	}

	blockStorageName := secalib.GenerateBlockStorageName()
	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName)
	blockStorageRefObj, err := secapi.BuildReferenceFromURN(blockStorageRef)
	if err != nil {
		t.Fatalf("Failed to build blockStorageRef to URN: %v", err)
	}

	blockStorageSize := secalib.GenerateBlockStorageSize()

	// Setup mock, if configured to use
	if suite.mockEnabled {
		mockParams := &mock.ComputeParamsV1{
			Params: &mock.Params{
				MockURL:   *suite.mockServerURL,
				AuthToken: suite.authToken,
				Tenant:    suite.tenant,
				Region:    suite.region,
			},
			Workspace: &mock.ResourceParams[schema.WorkspaceSpec]{
				Name: workspaceName,
				InitialLabels: schema.Labels{
					secalib.EnvLabel: secalib.EnvDevelopmentLabel,
				},
			},
			BlockStorage: &mock.ResourceParams[schema.BlockStorageSpec]{
				Name: blockStorageName,
				InitialSpec: &schema.BlockStorageSpec{
					SkuRef: *storageSkuRefObj,
					SizeGB: blockStorageSize,
				},
			},
			Instance: &mock.ResourceParams[schema.InstanceSpec]{
				Name: instanceName,
				InitialSpec: &schema.InstanceSpec{
					SkuRef: *instanceSkuRefObj,
					Zone:   initialInstanceZone,
					BootVolume: schema.VolumeReference{
						DeviceRef: *blockStorageRefObj,
					},
				},
				UpdatedSpec: &schema.InstanceSpec{
					SkuRef: *instanceSkuRefObj,
					Zone:   updatedInstanceZone,
					BootVolume: schema.VolumeReference{
						DeviceRef: *blockStorageRefObj,
					},
				},
			},
		}
		wm, err := mock.ConfigComputeLifecycleScenarioV1(suite.scenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	var workspaceResp *schema.Workspace
	var blockStorageResp *schema.BlockStorage

	t.WithNewStep("Create workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "CreateOrUpdateWorkspace")

		ws := &schema.Workspace{
			Metadata: &schema.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   workspaceName,
			},
		}
		workspaceResp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, ws)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, workspaceResp)

		suite.verifyStatusStep(sCtx, secalib.CreatingResourceState, *workspaceResp.Status.State)
	})

	t.WithNewStep("Get created workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "GetWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspaceName,
		}
		workspaceResp, err = suite.client.WorkspaceV1.GetWorkspace(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, workspaceResp)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *workspaceResp.Status.State)
	})

	t.WithNewStep("Create block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateBlockStorage", workspaceName)

		bo := &schema.BlockStorage{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      blockStorageName,
			},
			Spec: schema.BlockStorageSpec{
				SizeGB: blockStorageSize,
				SkuRef: *storageSkuRefObj,
			},
		}
		blockStorageResp, err = suite.client.StorageV1.CreateOrUpdateBlockStorage(ctx, bo)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockStorageResp)
	})

	t.WithNewStep("Get created block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetBlockStorage", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		blockStorageResp, err = suite.client.StorageV1.GetBlockStorage(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockStorageResp)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *blockStorageResp.Status.State)
	})

	// Instance
	var instanceResp *schema.Instance
	var expectedMeta *schema.RegionalWorkspaceResourceMetadata
	var expectedSpec *schema.InstanceSpec

	t.WithNewStep("Create instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "CreateOrUpdateInstance", workspaceName)

		inst := &schema.Instance{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      instanceName,
			},
			Spec: schema.InstanceSpec{
				SkuRef: *instanceSkuRefObj,
				Zone:   initialInstanceZone,
				BootVolume: schema.VolumeReference{
					DeviceRef: *blockStorageRefObj,
				},
			},
		}
		instanceResp, err = suite.client.ComputeV1.CreateOrUpdateInstance(ctx, inst)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMeta = secalib.NewRegionalWorkspaceResourceMetadata(instanceName, secalib.ComputeProviderV1, instanceResource, secalib.ApiVersion1, secalib.InstanceKind,
			suite.tenant, workspaceName, suite.region)
		expectedMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		expectedSpec = &schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   initialInstanceZone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		}
		suite.verifyInstanceSpecStep(sCtx, expectedSpec, &instanceResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingResourceState, *instanceResp.Status.State)
	})

	t.WithNewStep("Get created instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		instanceResp, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		suite.verifyInstanceSpecStep(sCtx, expectedSpec, &instanceResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *instanceResp.Status.State)
	})

	t.WithNewStep("Update instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "CreateOrUpdateInstance", workspaceName)

		instanceResp, err = suite.client.ComputeV1.CreateOrUpdateInstance(ctx, instanceResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		expectedSpec.Zone = updatedInstanceZone
		suite.verifyInstanceSpecStep(sCtx, expectedSpec, &instanceResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.UpdatingResourceState, *instanceResp.Status.State)
	})

	t.WithNewStep("Get updated instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		instanceResp, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		suite.verifyInstanceSpecStep(sCtx, expectedSpec, &instanceResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *instanceResp.Status.State)
	})

	t.WithNewStep("Stop instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "StopInstance", workspaceName)

		err = suite.client.ComputeV1.StopInstance(ctx, instanceResp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get stopped instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		instanceResp, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		suite.verifyInstanceSpecStep(sCtx, expectedSpec, &instanceResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.SuspendedResourceState, *instanceResp.Status.State)
	})

	t.WithNewStep("Start instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "StartInstance", workspaceName)

		err = suite.client.ComputeV1.StartInstance(ctx, instanceResp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get started instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		instanceResp, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		suite.verifyInstanceSpecStep(sCtx, expectedSpec, &instanceResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *instanceResp.Status.State)
	})

	t.WithNewStep("Restart instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "RestartInstance", workspaceName)

		err = suite.client.ComputeV1.RestartInstance(ctx, instanceResp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get restarted instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		instanceResp, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, instanceResp)

		expectedMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, instanceResp.Metadata)

		suite.verifyInstanceSpecStep(sCtx, expectedSpec, &instanceResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveResourceState, *instanceResp.Status.State)
	})

	t.WithNewStep("Delete instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "DeleteInstance", workspaceName)

		err = suite.client.ComputeV1.DeleteInstance(ctx, instanceResp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted instance", func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instanceName,
		}
		_, err = suite.client.ComputeV1.GetInstance(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *ComputeV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
