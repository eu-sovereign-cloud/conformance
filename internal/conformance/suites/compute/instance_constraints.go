package compute

import (
	"math/rand"
	"strings"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockcompute "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/compute"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// InstanceConstraintsValidationV1TestSuite verifies that Instance resources violating
// field constraints are rejected by the API with 422 Unprocessable Entity.
//
// Constraints tested:
//   - name: maxLength 128 (NameMetadata)
//   - name: pattern ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$ (NameMetadata)
//   - labels values: maxLength 63 (UserResourceMetadata)
//   - annotations values: maxLength 1024 (UserResourceMetadata)
type InstanceConstraintsValidationV1TestSuite struct {
	suites.RegionalTestSuite
	config *InstanceContraintsValidationV1Config

	params *params.InstanceConstraintsValidationV1Params
}

type InstanceContraintsValidationV1Config struct {
	AvailableZones []string
	InstanceSkus   []string
	StorageSkus    []string
}

func CreateInstanceConstraintsValidationV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *InstanceContraintsValidationV1Config) *InstanceConstraintsValidationV1TestSuite {
	suite := &InstanceConstraintsValidationV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.InstanceConstraintsValidationV1SuiteName.String()
	return suite
}

func (suite *InstanceConstraintsValidationV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Constraints")

	// Select skus
	instanceSkuName := suite.config.InstanceSkus[rand.Intn(len(suite.config.InstanceSkus))]
	storageSkuName := suite.config.StorageSkus[rand.Intn(len(suite.config.StorageSkus))]

	// Select zone
	zone := suite.config.AvailableZones[rand.Intn(len(suite.config.AvailableZones))]

	workspaceName := generators.GenerateWorkspaceName()
	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageSize := constants.BlockStorageInitialSize

	instanceSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.ComputeProviderV1Name, suite.Tenant, instanceSkuName)
	storageSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, storageSkuName)
	blockStorageRefObj := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, blockStorageName)

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Workspace for instance constraints violations testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	blockStorage, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Annotations(schema.Annotations{
			"description": "Block storage for conformance testing",
		}).
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

	// Instance with name exceeding maxLength: 128 (129 chars)
	overLengthName := strings.Repeat("a", 129)
	overLengthNameInstance, err := builders.NewInstanceBuilder().
		Name(overLengthName).
		Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Instance with over-length name"}).
		Spec(&schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		},
		).Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthNameInstance: %v", err)
	}

	// Instance with name violating kebab-case pattern (uppercase letters)
	invalidPatternName := "Invalid-Name-With-Uppercase"
	invalidPatternNameInstance, err := builders.NewInstanceBuilder().
		Name(invalidPatternName).
		Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Instance with non-kebab-case name"}).
		Spec(&schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		},
		).Build()
	if err != nil {
		t.Fatalf("Failed to build invalidPatternNameInstance: %v", err)
	}

	// Instance with label value exceeding maxLength: 63 (64 chars)
	overLengthLabelInstance, err := builders.NewInstanceBuilder().
		Name(generators.GenerateInstanceName()).
		Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
			"constraint-test":  strings.Repeat("x", 64),
		}).
		Annotations(schema.Annotations{"description": "Instance with over-length label value"}).
		Spec(&schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		},
		).Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthLabelInstance: %v", err)
	}

	// Instance with annotation value exceeding maxLength: 1024 (1025 chars)
	overLengthAnnotationInstance, err := builders.NewInstanceBuilder().
		Name(generators.GenerateInstanceName()).
		Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{
			"description":     "Instance with over-length annotation value",
			"long-annotation": strings.Repeat("y", 1025),
		}).
		Spec(&schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		},
		).Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthAnnotationInstance: %v", err)
	}

	p := &params.InstanceConstraintsValidationV1Params{
		Workspace:                    workspace,
		BlockStorage:                 blockStorage,
		OverLengthNameInstance:       overLengthNameInstance,
		InvalidPatternNameInstance:   invalidPatternNameInstance,
		OverLengthLabelValueInstance: overLengthLabelInstance,
		OverLengthAnnotationInstance: overLengthAnnotationInstance,
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockcompute.ConfigureInstanceConstraintsValidationV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *InstanceConstraintsValidationV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.ComputeProviderV1Name, string(schema.RegionalWorkspaceResourceMetadataKindResourceKindInstance))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Create workspace — required environment for instance operations
	workspace := suite.params.Workspace
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create workspace for test environment", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         workspace.Labels,
			Annotations:    workspace.Annotations,
			Metadata:       workspace.Metadata,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   workspace.Metadata.Name,
	}
	stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, workspaceTRef,
		steps.ResponseExpectsWithCondition[schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			Labels:   workspace.Labels,
			Metadata: workspace.Metadata,
			ResourceStatus: schema.WorkspaceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// name: maxLength 128 — must be rejected
	stepsBuilder.CreateOrUpdateInstanceExpectViolationV1Step(
		"Create an instance with name exceeding maxLength:128 — expect rejection",
		suite.Client.ComputeV1,
		suite.params.OverLengthNameInstance,
	)

	// name: pattern — must be rejected
	stepsBuilder.CreateOrUpdateInstanceExpectViolationV1Step(
		"Create an instance with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.ComputeV1,
		suite.params.InvalidPatternNameInstance,
	)

	// labels value: maxLength 63 — must be rejected
	stepsBuilder.CreateOrUpdateInstanceExpectViolationV1Step(
		"Create an instance with label value exceeding maxLength:63 — expect rejection",
		suite.Client.ComputeV1,
		suite.params.OverLengthLabelValueInstance,
	)

	// annotations value: maxLength 1024 — must be rejected
	stepsBuilder.CreateOrUpdateInstanceExpectViolationV1Step(
		"Create an instance with annotation value exceeding maxLength:1024 — expect rejection",
		suite.Client.ComputeV1,
		suite.params.OverLengthAnnotationInstance,
	)

	// Teardown workspace
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *InstanceConstraintsValidationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
