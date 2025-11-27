package secatest

import (
	"log/slog"
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/builders"
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
	var err error
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t, secalib.ComputeProviderV1, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	// Select skus
	instanceSkuName := suite.instanceSkus[rand.Intn(len(suite.instanceSkus))]
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]

	// Select zones
	initialInstanceZone := suite.availableZones[rand.Intn(len(suite.availableZones))]
	updatedInstanceZone := suite.availableZones[rand.Intn(len(suite.availableZones))]

	// Generate scenario data
	workspaceName := secalib.GenerateWorkspaceName()
	workspaceResource := secalib.GenerateWorkspaceResource(suite.tenant, workspaceName)
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
	blockStorageResource := secalib.GenerateBlockStorageResource(suite.tenant, workspaceName, blockStorageName)
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
			Instance: &[]mock.ResourceParams[schema.InstanceSpec]{
				{
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
	expectWorkspaceMeta, err := builders.NewRegionalResourceMetadataBuilder().
		Name(workspaceName).
		Provider(secalib.WorkspaceProviderV1).
		Resource(workspaceResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalResourceMetadataKindResourceKindWorkspace).
		Tenant(suite.tenant).
		Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectWorkspaceLabels := schema.Labels{secalib.EnvLabel: secalib.EnvDevelopmentLabel}

	suite.createOrUpdateWorkspaceV1Step("Create a workspace", t, suite.client.WorkspaceV1, workspace,
		responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			labels:        expectWorkspaceLabels,
			metadata:      expectWorkspaceMeta,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created Workspace
	workspaceTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   workspaceName,
	}
	suite.getWorkspaceV1Step("Get the created workspace", t, suite.client.WorkspaceV1, *workspaceTRef,
		responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			labels:        expectWorkspaceLabels,
			metadata:      expectWorkspaceMeta,
			resourceState: schema.ResourceStateActive,
		},
	)
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
	expectedBlockMeta, err := builders.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(blockStorageName).
		Provider(secalib.StorageProviderV1).
		Resource(blockStorageResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage).
		Tenant(suite.tenant).
		Workspace(workspaceName).
		Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectedBlockSpec := &schema.BlockStorageSpec{
		SizeGB: blockStorageSize,
		SkuRef: *storageSkuRefObj,
	}
	suite.createOrUpdateBlockStorageV1Step("Create a block storage", t, suite.client.StorageV1, block,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			metadata:      expectedBlockMeta,
			spec:          expectedBlockSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created block storage
	blockWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      blockStorageName,
	}
	suite.getBlockStorageV1Step("Get the created block storage", t, suite.client.StorageV1, *blockWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			metadata:      expectedBlockMeta,
			spec:          expectedBlockSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

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
	expectInstanceMeta, err := builders.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(instanceName).
		Provider(secalib.ComputeProviderV1).
		Resource(instanceResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindInstance).
		Tenant(suite.tenant).
		Workspace(workspaceName).
		Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectInstanceSpec := &schema.InstanceSpec{
		SkuRef: *instanceSkuRefObj,
		Zone:   initialInstanceZone,
		BootVolume: schema.VolumeReference{
			DeviceRef: *blockStorageRefObj,
		},
	}
	suite.createOrUpdateInstanceV1Step("Create an instance", t, suite.client.ComputeV1, instance,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			metadata:      expectInstanceMeta,
			spec:          expectInstanceSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created instance
	instanceWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      instanceName,
	}
	instance = suite.getInstanceV1Step("Get the created instance", t, suite.client.ComputeV1, *instanceWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			metadata:      expectInstanceMeta,
			spec:          expectInstanceSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Update the instance
	instance.Spec.Zone = updatedInstanceZone
	expectInstanceSpec.Zone = instance.Spec.Zone
	suite.createOrUpdateInstanceV1Step("Update the instance", t, suite.client.ComputeV1, instance,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			metadata:      expectInstanceMeta,
			spec:          expectInstanceSpec,
			resourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated instance
	instance = suite.getInstanceV1Step("Get the updated instance", t, suite.client.ComputeV1, *instanceWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			metadata:      expectInstanceMeta,
			spec:          expectInstanceSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Stop the instance
	suite.stopInstanceV1Step("Stop the instance", t, suite.client.ComputeV1, instance)

	// Get the stoped instance
	instance = suite.getInstanceV1Step("Get the updated instance", t, suite.client.ComputeV1, *instanceWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			metadata:      expectInstanceMeta,
			spec:          expectInstanceSpec,
			resourceState: schema.ResourceStateSuspended,
		},
	)

	// Start the instance
	suite.startInstanceV1Step("Start the instance", t, suite.client.ComputeV1, instance)

	// Get the started instance
	instance = suite.getInstanceV1Step("Get the started instance", t, suite.client.ComputeV1, *instanceWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			metadata:      expectInstanceMeta,
			spec:          expectInstanceSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Restart the instance
	suite.restartInstanceV1Step("Restart the instance", t, suite.client.ComputeV1, instance)

	// Get the restarted instance
	// TODO Find an away to assert if the instance is restarted
	instance = suite.getInstanceV1Step("Get the updated instance", t, suite.client.ComputeV1, *instanceWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			metadata:      expectInstanceMeta,
			spec:          expectInstanceSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Delete the instance
	suite.deleteInstanceV1Step("Delete the instance", t, suite.client.ComputeV1, instance)

	// Get the deleted instance
	suite.getInstanceWithErrorV1Step("Get the deleted instance", t, suite.client.ComputeV1, *instanceWRef, secapi.ErrResourceNotFound)

	// Delete the block storage
	suite.deleteBlockStorageV1Step("Delete the block storage", t, suite.client.StorageV1, block)

	// Get the deleted block storage
	suite.getBlockStorageWithErrorV1Step("Get the deleted block storage", t, suite.client.StorageV1, *blockWRef, secapi.ErrResourceNotFound)

	// Delete the workspace
	suite.deleteWorkspaceV1Step("Delete the workspace", t, suite.client.WorkspaceV1, workspace)

	// Get the deleted workspace
	suite.getWorkspaceWithErrorV1Step("Get the deleted workspace", t, suite.client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *ComputeV1TestSuite) TestListSuite(t provider.T) {
	var err error
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t, secalib.ComputeProviderV1, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	// Select skus
	instanceSkuName := suite.instanceSkus[rand.Intn(len(suite.instanceSkus))]
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]

	// Select zones
	initialInstanceZone := suite.availableZones[rand.Intn(len(suite.availableZones))]

	// Generate scenario data
	workspaceName := secalib.GenerateWorkspaceName()
	workspaceResource := secalib.GenerateWorkspaceResource(suite.tenant, workspaceName)
	instanceSkuRef := secalib.GenerateSkuRef(instanceSkuName)
	instanceSkuRefObj, err := secapi.BuildReferenceFromURN(instanceSkuRef)
	if err != nil {
		t.Fatalf("Failed to build instanceSkuRef to URN: %v", err)
	}

	instanceName1 := secalib.GenerateInstanceName()
	instanceResource1 := secalib.GenerateInstanceResource(suite.tenant, workspaceName, instanceName1)
	instanceName2 := secalib.GenerateInstanceName()
	instanceResource2 := secalib.GenerateInstanceResource(suite.tenant, workspaceName, instanceName2)
	instanceName3 := secalib.GenerateInstanceName()
	instanceResource3 := secalib.GenerateInstanceResource(suite.tenant, workspaceName, instanceName3)

	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)
	storageSkuRefObj, err := secapi.BuildReferenceFromURN(storageSkuRef)
	if err != nil {
		t.Fatalf("Failed to build storageSkuRef to URN: %v", err)
	}

	blockStorageName := secalib.GenerateBlockStorageName()
	blockStorageResource := secalib.GenerateBlockStorageResource(suite.tenant, workspaceName, blockStorageName)
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
					secalib.EnvLabel: secalib.EnvConformanceLabel,
				},
			},
			BlockStorage: &mock.ResourceParams[schema.BlockStorageSpec]{
				Name: blockStorageName,
				InitialLabels: map[string]string{
					secalib.EnvLabel: secalib.EnvConformanceLabel,
				},
				InitialSpec: &schema.BlockStorageSpec{
					SkuRef: *storageSkuRefObj,
					SizeGB: blockStorageSize,
				},
			},
			Instance: &[]mock.ResourceParams[schema.InstanceSpec]{
				{
					Name: instanceName1,
					InitialLabels: map[string]string{
						secalib.EnvLabel: secalib.EnvConformanceLabel,
					},
					InitialSpec: &schema.InstanceSpec{
						SkuRef: *instanceSkuRefObj,
						Zone:   initialInstanceZone,
						BootVolume: schema.VolumeReference{
							DeviceRef: *blockStorageRefObj,
						},
					},
				},
				{
					Name: instanceName2,
					InitialLabels: map[string]string{
						secalib.EnvLabel: secalib.EnvConformanceLabel,
					},
					InitialSpec: &schema.InstanceSpec{
						SkuRef: *instanceSkuRefObj,
						Zone:   initialInstanceZone,
						BootVolume: schema.VolumeReference{
							DeviceRef: *blockStorageRefObj,
						},
					},
				},
				{
					Name: instanceName3,
					InitialLabels: map[string]string{
						secalib.EnvLabel: secalib.EnvConformanceLabel,
					},
					InitialSpec: &schema.InstanceSpec{
						SkuRef: *instanceSkuRefObj,
						Zone:   initialInstanceZone,
						BootVolume: schema.VolumeReference{
							DeviceRef: *blockStorageRefObj,
						},
					},
				},
			},
		}
		wm, err := mock.ConfigComputeListLifecycleScenarioV1(suite.scenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	// Workspace

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			secalib.EnvLabel: secalib.EnvConformanceLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   workspaceName,
		},
	}

	expectWorkspaceMeta, err := builders.NewRegionalResourceMetadataBuilder().
		Name(workspaceName).
		Provider(secalib.WorkspaceProviderV1).
		Resource(workspaceResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalResourceMetadataKindResourceKindWorkspace).
		Tenant(suite.tenant).
		Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectWorkspaceLabels := schema.Labels{secalib.EnvLabel: secalib.EnvConformanceLabel}

	suite.createOrUpdateWorkspaceV1Step("Create a workspace", t, suite.client.WorkspaceV1, workspace,
		responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			labels:        expectWorkspaceLabels,
			metadata:      expectWorkspaceMeta,
			resourceState: schema.ResourceStateCreating,
		},
	)

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

	expectedBlockMeta, err := builders.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(blockStorageName).
		Provider(secalib.StorageProviderV1).
		Resource(blockStorageResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage).
		Tenant(suite.tenant).
		Workspace(workspaceName).
		Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectedBlockSpec := &schema.BlockStorageSpec{
		SizeGB: blockStorageSize,
		SkuRef: *storageSkuRefObj,
	}
	suite.createOrUpdateBlockStorageV1Step("Create a block storage", t, suite.client.StorageV1, block, responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
		metadata:      expectedBlockMeta,
		spec:          expectedBlockSpec,
		resourceState: schema.ResourceStateCreating,
	},
	)

	// Instance

	// Create an instance
	instances := &[]schema.Instance{
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      instanceName1,
				Resource:  instanceResource1,
			},
			Spec: schema.InstanceSpec{
				SkuRef: *instanceSkuRefObj,
				Zone:   initialInstanceZone,
				BootVolume: schema.VolumeReference{
					DeviceRef: *blockStorageRefObj,
				},
			},
		},
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      instanceName2,
				Resource:  instanceResource2,
			},
			Spec: schema.InstanceSpec{
				SkuRef: *instanceSkuRefObj,
				Zone:   initialInstanceZone,
				BootVolume: schema.VolumeReference{
					DeviceRef: *blockStorageRefObj,
				},
			},
		},
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      instanceName3,
				Resource:  instanceResource3,
			},
			Spec: schema.InstanceSpec{
				SkuRef: *instanceSkuRefObj,
				Zone:   initialInstanceZone,
				BootVolume: schema.VolumeReference{
					DeviceRef: *blockStorageRefObj,
				},
			},
		},
	}
	// Create instances
	for _, instance := range *instances {
		instanceResource := secalib.GenerateInstanceResource(suite.tenant, workspaceName, instance.Metadata.Name)
		expectInstanceMeta, err := builders.NewRegionalWorkspaceResourceMetadataBuilder().
			Name(instance.Metadata.Name).
			Provider(secalib.ComputeProviderV1).
			Resource(instanceResource).
			ApiVersion(secalib.ApiVersion1).
			Kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindInstance).
			Tenant(suite.tenant).
			Workspace(workspaceName).
			Region(suite.region).
			BuildResponse()
		if err != nil {
			t.Fatalf("Failed to build metadata: %v", err)
		}

		expectInstanceSpec := &schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   initialInstanceZone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		}
		suite.createOrUpdateInstanceV1Step("Create an instance", t, suite.client.ComputeV1, &instance,
			responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
				metadata:      expectInstanceMeta,
				spec:          expectInstanceSpec,
				resourceState: schema.ResourceStateCreating,
			},
		)
	}

	tref := secapi.TenantReference{Tenant: secapi.TenantID(suite.tenant)}
	wref := secapi.WorkspaceReference{Workspace: secapi.WorkspaceID(workspaceName)}
	// List instances
	suite.getListInstanceV1Step("List instances", t, suite.client.ComputeV1, tref, wref, nil)

	// List instances with limit
	suite.getListInstanceV1Step("Get list of instances", t, suite.client.ComputeV1, tref, wref,
		builders.NewListOptions().WithLimit(1))

	// List Instances with Label
	suite.getListInstanceV1Step("Get list of instances", t, suite.client.ComputeV1, tref, wref,
		builders.NewListOptions().WithLabels(builders.NewLabelsBuilder().
			Equals(secalib.EnvLabel, secalib.EnvConformanceLabel)))

	// List Instances with Limit and label
	suite.getListInstanceV1Step("Get list of instances", t, suite.client.ComputeV1, tref, wref,
		builders.NewListOptions().WithLimit(1).WithLabels(builders.NewLabelsBuilder().
			Equals(secalib.EnvLabel, secalib.EnvConformanceLabel)))

	// SKUS
	// List SKUS
	suite.getListSkusV1Step("List skus", t, suite.client.ComputeV1, secapi.TenantReference{Tenant: secapi.TenantID(suite.tenant)}, nil)

	// List SKUS with limit
	suite.getListSkusV1Step("Get list of skus", t, suite.client.ComputeV1, secapi.TenantReference{Tenant: secapi.TenantID(suite.tenant)},
		builders.NewListOptions().WithLimit(1))

	// List SKUS with Label
	suite.getListSkusV1Step("Get list of skus", t, suite.client.ComputeV1, secapi.TenantReference{Tenant: secapi.TenantID(suite.tenant)},
		builders.NewListOptions().WithLabels(builders.NewLabelsBuilder().
			Equals(secalib.EnvLabel, secalib.EnvConformanceLabel)))

	// List SKUS with Limit and label
	suite.getListSkusV1Step("Get list of skus", t, suite.client.ComputeV1, secapi.TenantReference{Tenant: secapi.TenantID(suite.tenant)},
		builders.NewListOptions().WithLabels(builders.NewLabelsBuilder().
			Equals(secalib.EnvLabel, secalib.EnvConformanceLabel)))

	// Resources deletion

	// Delete all instances
	for _, instance := range *instances {
		suite.deleteInstanceV1Step("Delete the instance", t, suite.client.ComputeV1, &instance)

		// Get the deleted instance
		instanceWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      instance.Metadata.Name,
		}
		suite.getInstanceWithErrorV1Step("Get the deleted instance", t, suite.client.ComputeV1, *instanceWRef, secapi.ErrResourceNotFound)
	}

	// Delete the block storage
	suite.deleteBlockStorageV1Step("Delete the block storage", t, suite.client.StorageV1, block)

	// Get the deleted block storage
	blockWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      blockStorageName,
	}
	suite.getBlockStorageWithErrorV1Step("Get the deleted block storage", t, suite.client.StorageV1, *blockWRef, secapi.ErrResourceNotFound)

	// Delete the workspace
	suite.deleteWorkspaceV1Step("Delete the workspace", t, suite.client.WorkspaceV1, workspace)

	// Get the deleted workspace
	workspaceTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   workspaceName,
	}
	suite.getWorkspaceWithErrorV1Step("Get the deleted workspace", t, suite.client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *ComputeV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
