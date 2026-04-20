package storage

import (
	"math/rand"
	"strings"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockstorage "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/storage"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// BlockStorageConstraintsV1TestSuite verifies that BlockStorage resources violating
// field constraints are rejected by the API with 422 Unprocessable Entity.
//
// Constraints tested:
//   - name: maxLength 128 (NameMetadata)
//   - name: pattern ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$ (NameMetadata)
//   - labels values: maxLength 63 (UserResourceMetadata)
//   - annotations values: maxLength 1024 (UserResourceMetadata)
type BlockStorageConstraintsV1TestSuite struct {
	suites.RegionalTestSuite

	StorageSkus []string

	params *params.BlockStorageConstraintsViolationsV1Params
}

func CreateBlockStorageConstraintsV1TestSuite(regionalTestSuite suites.RegionalTestSuite, storageSkus []string) *BlockStorageConstraintsV1TestSuite {
	suite := &BlockStorageConstraintsV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		StorageSkus:       storageSkus,
	}
	suite.ScenarioName = constants.BlockStorageConstraintsV1SuiteName.String()
	return suite
}

func (suite *BlockStorageConstraintsV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Constraints")

	// Select sku
	storageSkuName := suite.StorageSkus[rand.Intn(len(suite.StorageSkus))]

	workspaceName := generators.GenerateWorkspaceName()

	storageSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, storageSkuName)

	storageSize := constants.BlockStorageInitialSize

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Workspace for block storage constraints violations testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	// BlockStorage with name exceeding maxLength: 128 (129 chars)
	overLengthName := strings.Repeat("a", 129)
	overLengthNameBlockStorage, err := builders.NewBlockStorageBuilder().
		Name(overLengthName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "BlockStorage with over-length name"}).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: storageSize,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthNameBlockStorage: %v", err)
	}

	// BlockStorage with name violating kebab-case pattern (uppercase letters)
	invalidPatternName := "Invalid-Name-With-Uppercase"
	invalidPatternNameBlockStorage, err := builders.NewBlockStorageBuilder().
		Name(invalidPatternName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "BlockStorage with non-kebab-case name"}).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: storageSize,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build invalidPatternNameBlockStorage: %v", err)
	}

	// BlockStorage with label value exceeding maxLength: 63 (64 chars)
	overLengthLabelBlockStorage, err := builders.NewBlockStorageBuilder().
		Name(generators.GenerateBlockStorageName()).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
			"constraint-test":  strings.Repeat("x", 64),
		}).
		Annotations(schema.Annotations{"description": "BlockStorage with over-length label value"}).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: storageSize,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthLabelBlockStorage: %v", err)
	}

	// BlockStorage with annotation value exceeding maxLength: 1024 (1025 chars)
	overLengthAnnotationBlockStorage, err := builders.NewBlockStorageBuilder().
		Name(generators.GenerateBlockStorageName()).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{
			"description":     "BlockStorage with over-length annotation value",
			"long-annotation": strings.Repeat("y", 1025),
		}).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: storageSize,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthAnnotationBlockStorage: %v", err)
	}

	p := &params.BlockStorageConstraintsViolationsV1Params{
		Workspace:                        workspace,
		OverLengthNameBlockStorage:       overLengthNameBlockStorage,
		InvalidPatternNameBlockStorage:   invalidPatternNameBlockStorage,
		OverLengthLabelValueBlockStorage: overLengthLabelBlockStorage,
		OverLengthAnnotationBlockStorage: overLengthAnnotationBlockStorage,
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockstorage.ConfigureBlockStorageConstraintsV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *BlockStorageConstraintsV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.StorageProviderV1Name, string(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	workspace := suite.params.Workspace
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   workspace.Metadata.Name,
	}

	// Create workspace
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create workspace for test environment", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         workspace.Labels,
			Annotations:    workspace.Annotations,
			Metadata:       workspace.Metadata,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
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
	stepsBuilder.CreateOrUpdateBlockStorageExpectViolationV1Step(
		"Create a block storage with name exceeding maxLength:128 — expect rejection",
		suite.Client.StorageV1,
		suite.params.OverLengthNameBlockStorage,
	)

	// name: pattern — must be rejected
	stepsBuilder.CreateOrUpdateBlockStorageExpectViolationV1Step(
		"Create a block storage with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.StorageV1,
		suite.params.InvalidPatternNameBlockStorage,
	)

	// labels value: maxLength 64 — must be rejected
	stepsBuilder.CreateOrUpdateBlockStorageExpectViolationV1Step(
		"Create a block storage with label value exceeding maxLength:64 — expect rejection",
		suite.Client.StorageV1,
		suite.params.OverLengthLabelValueBlockStorage,
	)

	// annotations value: maxLength 1024 — must be rejected
	stepsBuilder.CreateOrUpdateBlockStorageExpectViolationV1Step(
		"Create a block storage with annotation value exceeding maxLength:1024 — expect rejection",
		suite.Client.StorageV1,
		suite.params.OverLengthAnnotationBlockStorage,
	)

	// Teardown workspace
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *BlockStorageConstraintsV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
