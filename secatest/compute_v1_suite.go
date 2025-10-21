package secatest

import (
	"context"
	"log/slog"
	"math/rand"

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

	// Workspace

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			secalib.EnvLabel: secalib.EnvDevelopmentLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   workspaceName,
		},
	}
	suite.createOrUpdateWorkspaceV1Step("Create a workspace", t, ctx, suite.client.WorkspaceV1, workspace, nil, nil, secalib.CreatingResourceState)

	// Get the created Workspace
	workspaceTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   workspaceName,
	}
	suite.getWorkspaceV1Step("Get the created workspace", t, ctx, suite.client.WorkspaceV1, *workspaceTRef, nil, nil, secalib.ActiveResourceState)

	// Block storage

	// Create a block storage
	block := &schema.BlockStorage{
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
	suite.createOrUpdateBlockStorageV1Step("Create a block storage", t, ctx, suite.client.StorageV1, block,
		nil, nil, secalib.CreatingResourceState)

	// Get the created block storage
	blockWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      blockStorageName,
	}
	suite.getBlockStorageV1Step("Get the created block storage", t, ctx, suite.client.StorageV1, blockWRef,
		nil, nil, secalib.ActiveResourceState)

	// Instance

	// Create an instance
	instance := &schema.Instance{
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
	expectInstanceMeta := secalib.NewRegionalWorkspaceResourceMetadata(instanceName,
		secalib.ComputeProviderV1,
		instanceResource,
		secalib.ApiVersion1,
		secalib.InstanceKind,
		suite.tenant,
		workspaceName,
		suite.region)
	expectInstanceSpec := instance.Spec
	suite.createOrUpdateInstanceV1Step("Create an instance", t, ctx, suite.client.ComputeV1, instance,
		expectInstanceMeta, &expectInstanceSpec, secalib.CreatingResourceState)

	// Get the created instance
	instanceWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      instanceName,
	}
	instance = suite.getInstanceV1Step("Get the created instance", t, ctx, suite.client.ComputeV1, instanceWRef,
		expectInstanceMeta, &expectInstanceSpec, secalib.ActiveResourceState)

	// Update the instance
	instance.Spec.Zone = updatedInstanceZone
	expectInstanceSpec = instance.Spec
	suite.createOrUpdateInstanceV1Step("Update the instance", t, ctx, suite.client.ComputeV1, instance,
		expectInstanceMeta, &expectInstanceSpec, secalib.UpdatingResourceState)

	// Get the updated instance
	instance = suite.getInstanceV1Step("Get the updated instance", t, ctx, suite.client.ComputeV1, instanceWRef,
		expectInstanceMeta, &expectInstanceSpec, secalib.ActiveResourceState)

	// Stop the instance
	suite.stopInstanceV1Step("Stop the instance", t, ctx, suite.client.ComputeV1, instance)

	// Get the stoped instance
	instance = suite.getInstanceV1Step("Get the updated instance", t, ctx, suite.client.ComputeV1, instanceWRef,
		expectInstanceMeta, &expectInstanceSpec, secalib.SuspendedResourceState)

	// Start the instance
	suite.startInstanceV1Step("Start the instance", t, ctx, suite.client.ComputeV1, instance)

	// Get the started instance
	instance = suite.getInstanceV1Step("Get the started instance", t, ctx, suite.client.ComputeV1, instanceWRef,
		expectInstanceMeta, &expectInstanceSpec, secalib.ActiveResourceState)

	// Restart the instance
	suite.restartInstanceV1Step("Restart the instance", t, ctx, suite.client.ComputeV1, instance)

	// Get the restarted instance
	// TODO Find an away to assert if the instance is restarted
	instance = suite.getInstanceV1Step("Get the updated instance", t, ctx, suite.client.ComputeV1, instanceWRef,
		expectInstanceMeta, &expectInstanceSpec, secalib.ActiveResourceState)

	// Delete the instance
	suite.deleteInstanceV1Step("Delete the instance", t, ctx, suite.client.ComputeV1, instance)

	// Get the deleted instance
	suite.getInstanceWithErrorV1Step("Get the deleted instance", t, ctx, suite.client.ComputeV1, instanceWRef, secapi.ErrResourceNotFound)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *ComputeV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
