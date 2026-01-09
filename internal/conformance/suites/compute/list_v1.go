package compute

import (
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	mockcompute "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/compute"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type ListV1TestSuite struct {
	suites.RegionalTestSuite

	AvailableZones []string
	InstanceSkus   []string
	StorageSkus    []string

	params *params.ComputeListParamsV1
}

func (suite *ListV1TestSuite) BeforeAll(t provider.T) {
	var err error

	// Select skus
	instanceSkuName := suite.InstanceSkus[rand.Intn(len(suite.InstanceSkus))]
	storageSkuName := suite.StorageSkus[rand.Intn(len(suite.StorageSkus))]

	// Select zones
	initialInstanceZone := suite.AvailableZones[rand.Intn(len(suite.AvailableZones))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	instanceSkuRefObj, err := generators.GenerateSkuRefObject(instanceSkuName)
	if err != nil {
		t.Fatalf("Failed to build instanceSkuRef to URN: %v", err)
	}

	instanceName1 := generators.GenerateInstanceName()
	instanceName2 := generators.GenerateInstanceName()
	instanceName3 := generators.GenerateInstanceName()

	storageSkuRefObj, err := generators.GenerateSkuRefObject(storageSkuName)
	if err != nil {
		t.Fatalf("Failed to build storageSkuRef to URN: %v", err)
	}

	blockStorageName := generators.GenerateBlockStorageName()

	blockStorageRefObj, err := generators.GenerateBlockStorageRefObject(blockStorageName)
	if err != nil {
		t.Fatalf("Failed to build blockStorageRef to URN: %v", err)
	}

	blockStorageSize := generators.GenerateBlockStorageSize()

	params := &params.ComputeListParamsV1{
		BaseParams: &params.BaseParams{
			Tenant: suite.Tenant,
			Region: suite.Region,
			MockParams: &mock.MockParams{
				ServerURL: *suite.MockServerURL,
				AuthToken: suite.AuthToken,
			},
		},
		Workspace: &params.ResourceParams[schema.WorkspaceSpec]{
			Name: workspaceName,
			InitialLabels: schema.Labels{
				constants.EnvLabel: constants.EnvConformanceLabel,
			},
		},
		BlockStorage: &params.ResourceParams[schema.BlockStorageSpec]{
			Name: blockStorageName,
			InitialLabels: map[string]string{
				constants.EnvLabel: constants.EnvConformanceLabel,
			},
			InitialSpec: &schema.BlockStorageSpec{
				SkuRef: *storageSkuRefObj,
				SizeGB: blockStorageSize,
			},
		},
		Instances: []params.ResourceParams[schema.InstanceSpec]{
			{
				Name: instanceName1,
				InitialLabels: map[string]string{
					constants.EnvLabel: constants.EnvConformanceLabel,
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
					constants.EnvLabel: constants.EnvConformanceLabel,
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
					constants.EnvLabel: constants.EnvConformanceLabel,
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
	suite.params = params

	err = suites.SetupMock(&suite.TestSuite, mockcompute.ConfigureListScenarioV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ListV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.ComputeProviderV1, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	stepsBuilder := steps.NewStepsConfigurator(&suite.TestSuite, t)

	// Workspace

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.Tenant,
			Name:   suite.params.Workspace.Name,
		},
	}

	expectWorkspaceMeta, err := builders.NewWorkspaceMetadataBuilder().
		Name(suite.params.Workspace.Name).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectWorkspaceLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}

	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectWorkspaceLabels,
			Metadata:      expectWorkspaceMeta,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Block storage

	// Create a block storage
	block := &schema.BlockStorage{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.Tenant,
			Workspace: suite.params.Workspace.Name,
			Name:      suite.params.BlockStorage.Name,
		},
		Spec: schema.BlockStorageSpec{
			SizeGB: suite.params.BlockStorage.InitialSpec.SizeGB,
			SkuRef: suite.params.BlockStorage.InitialSpec.SkuRef,
		},
	}

	expectedBlockMeta, err := builders.NewBlockStorageMetadataBuilder().
		Name(suite.params.BlockStorage.Name).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(suite.params.Workspace.Name).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectedBlockSpec := &schema.BlockStorageSpec{
		SizeGB: suite.params.BlockStorage.InitialSpec.SizeGB,
		SkuRef: suite.params.BlockStorage.InitialSpec.SkuRef,
	}
	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create a block storage", suite.Client.StorageV1, block, steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
		Metadata:      expectedBlockMeta,
		Spec:          expectedBlockSpec,
		ResourceState: schema.ResourceStateCreating,
	})

	// Instance

	// Create instances
	instances := &[]schema.Instance{
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: suite.params.Workspace.Name,
				Name:      suite.params.Instances[0].Name,
				Resource:  generators.GenerateInstanceResource(suite.Tenant, suite.params.Workspace.Name, suite.params.Instances[0].Name),
			},
			Labels: map[string]string{
				constants.EnvLabel: constants.EnvConformanceLabel,
			},
			Spec: schema.InstanceSpec{
				SkuRef: suite.params.Instances[0].InitialSpec.SkuRef,
				Zone:   suite.params.Instances[0].InitialSpec.Zone,
				BootVolume: schema.VolumeReference{
					DeviceRef: suite.params.Instances[0].InitialSpec.BootVolume.DeviceRef,
				},
			},
		},
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: suite.params.Workspace.Name,
				Name:      suite.params.Instances[1].Name,
				Resource:  generators.GenerateInstanceResource(suite.Tenant, suite.params.Workspace.Name, suite.params.Instances[1].Name),
			},
			Labels: map[string]string{
				constants.EnvLabel: constants.EnvConformanceLabel,
			},
			Spec: schema.InstanceSpec{
				SkuRef: suite.params.Instances[1].InitialSpec.SkuRef,
				Zone:   suite.params.Instances[1].InitialSpec.Zone,
				BootVolume: schema.VolumeReference{
					DeviceRef: suite.params.Instances[1].InitialSpec.BootVolume.DeviceRef,
				},
			},
		},
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: suite.params.Workspace.Name,
				Name:      suite.params.Instances[2].Name,
				Resource:  generators.GenerateInstanceResource(suite.Tenant, suite.params.Workspace.Name, suite.params.Instances[2].Name),
			},
			Labels: map[string]string{
				constants.EnvLabel: constants.EnvConformanceLabel,
			},
			Spec: schema.InstanceSpec{
				SkuRef: suite.params.Instances[2].InitialSpec.SkuRef,
				Zone:   suite.params.Instances[2].InitialSpec.Zone,
				BootVolume: schema.VolumeReference{
					DeviceRef: suite.params.Instances[2].InitialSpec.BootVolume.DeviceRef,
				},
			},
		},
	}
	// Create instances
	for _, instance := range *instances {
		expectInstanceMeta, err := builders.NewInstanceMetadataBuilder().
			Name(instance.Metadata.Name).
			Provider(constants.ComputeProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(suite.Tenant).Workspace(suite.params.Workspace.Name).Region(suite.Region).
			Build()
		if err != nil {
			t.Fatalf("Failed to build Metadata: %v", err)
		}

		expectInstanceSpec := &schema.InstanceSpec{
			SkuRef: instance.Spec.SkuRef,
			Zone:   instance.Spec.Zone,
			BootVolume: schema.VolumeReference{
				DeviceRef: instance.Spec.BootVolume.DeviceRef,
			},
		}
		stepsBuilder.CreateOrUpdateInstanceV1Step("Create an instance", suite.Client.ComputeV1, &instance,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
				Metadata:      expectInstanceMeta,
				Spec:          expectInstanceSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}

	wref := secapi.WorkspaceReference{
		Workspace: secapi.WorkspaceID(suite.params.Workspace.Name),
		Tenant:    secapi.TenantID(suite.Tenant),
	}
	// List instances
	stepsBuilder.GetListInstanceV1Step("List instances", suite.Client.ComputeV1, wref, nil)

	// List instances with limit
	stepsBuilder.GetListInstanceV1Step("Get list of instances", suite.Client.ComputeV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List Instances with Label
	stepsBuilder.GetListInstanceV1Step("Get list of instances", suite.Client.ComputeV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List Instances with Limit and label
	stepsBuilder.GetListInstanceV1Step("Get list of instances", suite.Client.ComputeV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// SKUS
	// List SKUS
	stepsBuilder.GetListSkusV1Step("List skus", suite.Client.ComputeV1, secapi.TenantReference{Tenant: secapi.TenantID(suite.Tenant)}, nil)

	// List SKUS with limit
	stepsBuilder.GetListSkusV1Step("Get list of skus", suite.Client.ComputeV1, secapi.TenantReference{Tenant: secapi.TenantID(suite.Tenant)},
		secapi.NewListOptions().WithLimit(1))

	// Delete all instances
	for _, instance := range *instances {
		stepsBuilder.DeleteInstanceV1Step("Delete the instance", suite.Client.ComputeV1, &instance)

		// Get the deleted instance
		instanceWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.Tenant),
			Workspace: secapi.WorkspaceID(suite.params.Workspace.Name),
			Name:      instance.Metadata.Name,
		}
		stepsBuilder.GetInstanceWithErrorV1Step("Get the deleted instance", suite.Client.ComputeV1, *instanceWRef, secapi.ErrResourceNotFound)
	}

	// Delete the block storage
	stepsBuilder.DeleteBlockStorageV1Step("Delete the block storage", suite.Client.StorageV1, block)

	// Get the deleted block storage
	blockWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(suite.params.Workspace.Name),
		Name:      suite.params.BlockStorage.Name,
	}
	stepsBuilder.GetBlockStorageWithErrorV1Step("Get the deleted block storage", suite.Client.StorageV1, *blockWRef, secapi.ErrResourceNotFound)

	// Delete the workspace
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)

	// Get the deleted workspace
	workspaceTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   suite.params.Workspace.Name,
	}
	stepsBuilder.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	suite.FinishScenario()
}

func (suite *ListV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
