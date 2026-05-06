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
//   - spec.userData: maxLength 65536 (InstanceSpec)
//   - spec.antiAffinityGroup: maxLength 64 (InstanceSpec)
//   - spec.sshKeys[]: maxLength 4096 (InstanceSpec)
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

	instanceSkuName := suite.config.InstanceSkus[rand.Intn(len(suite.config.InstanceSkus))]
	storageSkuName := suite.config.StorageSkus[rand.Intn(len(suite.config.StorageSkus))]
	zone := suite.config.AvailableZones[rand.Intn(len(suite.config.AvailableZones))]

	workspaceName := generators.GenerateWorkspaceName()
	blockStorageName := generators.GenerateBlockStorageName()

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
		Annotations(schema.Annotations{"description": "Block storage for conformance testing"}).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: constants.BlockStorageInitialSize,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}

	buildInstance := func(name string, labels schema.Labels, annotations schema.Annotations, spec *schema.InstanceSpec) *schema.Instance {
		instance, err := builders.NewInstanceBuilder().
			Name(name).
			Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
			Labels(labels).
			Annotations(annotations).
			Spec(spec).Build()
		if err != nil {
			t.Fatalf("Failed to build Instance: %v", err)
		}
		return instance
	}

	baseSpec := func() *schema.InstanceSpec {
		return &schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		}
	}

	baseLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}

	overLengthUserDataSpec := baseSpec()
	overLengthUserDataSpec.UserData = strings.Repeat("a", 65537)

	overLengthAntiAffinityGroupSpec := baseSpec()
	overLengthAntiAffinityGroupSpec.AntiAffinityGroup = strings.Repeat("a", 65)

	overLengthSshKeySpec := baseSpec()
	overLengthSshKeySpec.SshKeys = []string{strings.Repeat("a", 4097)}

	p := &params.InstanceConstraintsValidationV1Params{
		Workspace:    workspace,
		BlockStorage: blockStorage,
		OverLengthNameInstance: buildInstance(
			strings.Repeat("a", 129),
			baseLabels,
			schema.Annotations{"description": "Instance with over-length name"},
			baseSpec(),
		),
		InvalidPatternNameInstance: buildInstance(
			"Invalid-Name-With-Uppercase",
			baseLabels,
			schema.Annotations{"description": "Instance with non-kebab-case name"},
			baseSpec(),
		),
		OverLengthLabelValueInstance: buildInstance(
			generators.GenerateInstanceName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel, "constraint-test": strings.Repeat("x", 64)},
			schema.Annotations{"description": "Instance with over-length label value"},
			baseSpec(),
		),
		OverLengthAnnotationInstance: buildInstance(
			generators.GenerateInstanceName(),
			baseLabels,
			schema.Annotations{"description": "Instance with over-length annotation value", "long-annotation": strings.Repeat("y", 1025)},
			baseSpec(),
		),
		OverLengthUserDataInstance: buildInstance(
			generators.GenerateInstanceName(),
			baseLabels,
			schema.Annotations{"description": "Instance with over-length userData"},
			overLengthUserDataSpec,
		),
		OverLengthAntiAffinityGroupInstance: buildInstance(
			generators.GenerateInstanceName(),
			baseLabels,
			schema.Annotations{"description": "Instance with over-length antiAffinityGroup"},
			overLengthAntiAffinityGroupSpec,
		),
		OverLengthSshKeyInstance: buildInstance(
			generators.GenerateInstanceName(),
			baseLabels,
			schema.Annotations{"description": "Instance with over-length sshKey"},
			overLengthSshKeySpec,
		),
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockcompute.ConfigureInstanceConstraintsValidationV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *InstanceConstraintsValidationV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.ComputeProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalResourceMetadataKindResourceKindInstance))
	suite.ConfigureDepends(t,
		string(schema.RegionalResourceMetadataKindResourceKindWorkspace),
		string(schema.RegionalResourceMetadataKindResourceKindBlockStorage),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace

	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	expectWorkspaceAnnotations := workspace.Annotations
	expectWorkspaceExtensions := workspace.Extensions
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", t, suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Annotations:    expectWorkspaceAnnotations,
			Extensions:     expectWorkspaceExtensions,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   suite.params.Workspace.Metadata.Name,
	}
	stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, workspaceTRef,
		steps.ResponseExpectsWithCondition[schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			Labels:   expectWorkspaceLabels,
			Metadata: expectWorkspaceMeta,
			ResourceStatus: schema.WorkspaceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Block storage

	block := suite.params.BlockStorage
	expectedBlockMeta := block.Metadata
	expectedBlockSpec := &block.Spec
	expectedBlockLabels := block.Labels
	expectedBlockAnnotations := block.Annotations
	expectedBlockExtensions := block.Extensions
	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create a block storage", t, suite.Client.StorageV1, block,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:       expectedBlockMeta,
			Labels:         expectedBlockLabels,
			Annotations:    expectedBlockAnnotations,
			Extensions:     expectedBlockExtensions,
			Spec:           expectedBlockSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	blockWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(suite.params.Workspace.Metadata.Name),
		Name:      suite.params.BlockStorage.Metadata.Name,
	}
	stepsBuilder.GetBlockStorageV1Step("Get the created block storage", suite.Client.StorageV1, blockWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus]{
			Metadata: expectedBlockMeta,
			Spec:     expectedBlockSpec,
			ResourceStatus: schema.BlockStorageStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Constraint violations

	stepsBuilder.CreateOrUpdateInstanceExpectViolationV1Step(
		"Create an instance with name exceeding maxLength:128 — expect rejection",
		suite.Client.ComputeV1,
		suite.params.OverLengthNameInstance,
	)
	stepsBuilder.CreateOrUpdateInstanceExpectViolationV1Step(
		"Create an instance with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.ComputeV1,
		suite.params.InvalidPatternNameInstance,
	)
	stepsBuilder.CreateOrUpdateInstanceExpectViolationV1Step(
		"Create an instance with label value exceeding maxLength:63 — expect rejection",
		suite.Client.ComputeV1,
		suite.params.OverLengthLabelValueInstance,
	)
	stepsBuilder.CreateOrUpdateInstanceExpectViolationV1Step(
		"Create an instance with annotation value exceeding maxLength:1024 — expect rejection",
		suite.Client.ComputeV1,
		suite.params.OverLengthAnnotationInstance,
	)
	stepsBuilder.CreateOrUpdateInstanceExpectViolationV1Step(
		"Create an instance with userData exceeding maxLength:65536 — expect rejection",
		suite.Client.ComputeV1,
		suite.params.OverLengthUserDataInstance,
	)
	stepsBuilder.CreateOrUpdateInstanceExpectViolationV1Step(
		"Create an instance with antiAffinityGroup exceeding maxLength:64 — expect rejection",
		suite.Client.ComputeV1,
		suite.params.OverLengthAntiAffinityGroupInstance,
	)
	stepsBuilder.CreateOrUpdateInstanceExpectViolationV1Step(
		"Create an instance with sshKey exceeding maxLength:4096 — expect rejection",
		suite.Client.ComputeV1,
		suite.params.OverLengthSshKeyInstance,
	)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *InstanceConstraintsValidationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
