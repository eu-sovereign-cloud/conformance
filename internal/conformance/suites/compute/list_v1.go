package compute

import (
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockCompute "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/compute"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type ComputeListV1TestSuite struct {
	suites.RegionalTestSuite

	config *ComputeListV1Config
	params *params.ComputeListV1Params
}

type ComputeListV1Config struct {
	AvailableZones []string
	InstanceSkus   []string
	StorageSkus    []string
}

func CreateListV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *ComputeListV1Config) *ComputeListV1TestSuite {
	suite := &ComputeListV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.ComputeV1ListSuiteName
	return suite
}

func (suite *ComputeListV1TestSuite) BeforeAll(t provider.T) {
	var err error

	// Select skus
	instanceSkuName := suite.config.InstanceSkus[rand.Intn(len(suite.config.InstanceSkus))]
	storageSkuName := suite.config.StorageSkus[rand.Intn(len(suite.config.StorageSkus))]

	// Select zones
	zone := suite.config.AvailableZones[rand.Intn(len(suite.config.AvailableZones))]

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

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	blockStorage, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: blockStorageSize,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}

	instance1, err := builders.NewInstanceBuilder().
		Name(instanceName1).
		Provider(constants.ComputeProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Instance: %v", err)
	}

	instance2, err := builders.NewInstanceBuilder().
		Name(instanceName2).
		Provider(constants.ComputeProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Instance: %v", err)
	}

	instance3, err := builders.NewInstanceBuilder().
		Name(instanceName3).
		Provider(constants.ComputeProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Instance: %v", err)
	}

	instances := []schema.Instance{*instance1, *instance2, *instance3}

	params := &params.ComputeListV1Params{
		Workspace:    workspace,
		BlockStorage: blockStorage,
		Instances:    instances,
	}

	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockCompute.ConfigureListScenarioV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ComputeListV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.ComputeProviderV1, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	workspace := suite.params.Workspace
	block := suite.params.BlockStorage
	instances := suite.params.Instances

	t.WithNewStep("Workspace", func(wsCtx provider.StepCtx) {
		wsSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, wsCtx)

		expectWorkspaceMeta := workspace.Metadata
		expectWorkspaceLabels := workspace.Labels
		wsSteps.CreateOrUpdateWorkspaceV1Step("Create", suite.Client.WorkspaceV1, workspace,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				Labels:        expectWorkspaceLabels,
				Metadata:      expectWorkspaceMeta,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	})

	t.WithNewStep("BlockStorage", func(bsCtx provider.StepCtx) {
		bsSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, bsCtx)

		expectedBlockMeta := block.Metadata
		expectedBlockSpec := &block.Spec
		bsSteps.CreateOrUpdateBlockStorageV1Step("Create", suite.Client.StorageV1, block,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
				Metadata:      expectedBlockMeta,
				Spec:          expectedBlockSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	})

	t.WithNewStep("Instance", func(instCtx provider.StepCtx) {
		instSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, instCtx)

		for _, instance := range instances {
			inst := instance
			expectInstanceMeta := inst.Metadata
			expectInstanceSpec := &inst.Spec
			instSteps.CreateOrUpdateInstanceV1Step("Create", suite.Client.ComputeV1, &inst,
				steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
					Metadata:      expectInstanceMeta,
					Spec:          expectInstanceSpec,
					ResourceState: schema.ResourceStateCreating,
				},
			)
		}

		wref := secapi.WorkspaceReference{
			Name:      workspace.Metadata.Name,
			Workspace: secapi.WorkspaceID(workspace.Metadata.Name),
			Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
		}
		instSteps.GetListInstanceV1Step("ListAll", suite.Client.ComputeV1, wref, nil)
		instSteps.GetListInstanceV1Step("ListWithLimit", suite.Client.ComputeV1, wref,
			secapi.NewListOptions().WithLimit(1))
		instSteps.GetListInstanceV1Step("ListWithLabel", suite.Client.ComputeV1, wref,
			secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
		instSteps.GetListInstanceV1Step("ListWithLabelAndLimit", suite.Client.ComputeV1, wref,
			secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
	})

	t.WithNewStep("Skus", func(skuCtx provider.StepCtx) {
		skuSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, skuCtx)

		ttRef := secapi.TenantReference{Tenant: secapi.TenantID(workspace.Metadata.Tenant)}
		skuSteps.GetListSkusV1Step("ListAll", suite.Client.ComputeV1, ttRef, nil)
		skuSteps.GetListSkusV1Step("ListWithLimit", suite.Client.ComputeV1, ttRef,
			secapi.NewListOptions().WithLimit(1))
	})

	t.WithNewStep("Delete", func(delCtx provider.StepCtx) {
		delSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, delCtx)

		for _, instance := range instances {
			inst := instance
			delSteps.DeleteInstanceV1Step("Instance", suite.Client.ComputeV1, &inst)

			instanceWRef := secapi.WorkspaceReference{
				Tenant:    secapi.TenantID(inst.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(inst.Metadata.Workspace),
				Name:      inst.Metadata.Name,
			}
			delSteps.GetInstanceWithErrorV1Step("GetDeletedInstance", suite.Client.ComputeV1, instanceWRef, secapi.ErrResourceNotFound)
		}

		delSteps.DeleteBlockStorageV1Step("BlockStorage", suite.Client.StorageV1, block)

		blockWRef := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(block.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(block.Metadata.Workspace),
			Name:      block.Metadata.Name,
		}
		delSteps.GetBlockStorageWithErrorV1Step("GetDeletedBlockStorage", suite.Client.StorageV1, blockWRef, secapi.ErrResourceNotFound)

		workspaceTRef := secapi.TenantReference{
			Tenant: secapi.TenantID(workspace.Metadata.Tenant),
			Name:   workspace.Metadata.Name,
		}
		delSteps.DeleteWorkspaceV1Step("Workspace", suite.Client.WorkspaceV1, workspace)
		delSteps.GetWorkspaceWithErrorV1Step("GetDeletedWorkspace", suite.Client.WorkspaceV1, workspaceTRef, secapi.ErrResourceNotFound)
	})

	suite.FinishScenario()
}

func (suite *ComputeListV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
