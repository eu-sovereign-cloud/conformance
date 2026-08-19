package usage

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockUsage "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/usage"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// HaMultiZoneV1TestSuite provisions two replicas of the same service, tagged with the
// same AntiAffinityGroup but placed in different zones, and starts both so they run
// concurrently - the high-availability placement pattern where a provider is expected
// to keep anti-affinity-grouped instances off the same underlying host.
type HaMultiZoneV1TestSuite struct {
	suites.RegionalTestSuite

	config *HaMultiZoneV1Config
	params *params.HaMultiZoneV1Params
}

type HaMultiZoneV1Config struct {
	RegionZones  []string
	StorageSkus  []string
	InstanceSkus []string
}

func CreateHaMultiZoneV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *HaMultiZoneV1Config) *HaMultiZoneV1TestSuite {
	suite := &HaMultiZoneV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.UsageHaMultiZoneV1SuiteName.String()
	return suite
}

func (suite *HaMultiZoneV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.UsageParentSuite)

	// Pick two distinct zones so the replicas land in different failure domains
	zoneA := suite.config.RegionZones[rand.Intn(len(suite.config.RegionZones))]
	zoneB := suite.config.RegionZones[(indexOf(suite.config.RegionZones, zoneA)+1)%len(suite.config.RegionZones)]

	storageSkuName := suite.config.StorageSkus[rand.Intn(len(suite.config.StorageSkus))]
	instanceSkuName := suite.config.InstanceSkus[rand.Intn(len(suite.config.InstanceSkus))]

	antiAffinityGroup := fmt.Sprintf("ha-group-%d", rand.Intn(math.MaxInt32))

	workspaceName := generators.GenerateWorkspaceName()

	storageSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, storageSkuName)
	instanceSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.ComputeProviderV1Name, suite.Tenant, instanceSkuName)

	blockStorageAName := generators.GenerateBlockStorageName()
	blockStorageARefObj := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, blockStorageAName)

	blockStorageBName := generators.GenerateBlockStorageName()
	blockStorageBRefObj := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, blockStorageBName)

	replicaAName := generators.GenerateInstanceName()
	replicaBName := generators.GenerateInstanceName()

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	blockStorageA, err := builders.NewBlockStorageBuilder().
		Name(blockStorageAName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: constants.BlockStorageInitialSize,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage A: %v", err)
	}

	blockStorageB, err := builders.NewBlockStorageBuilder().
		Name(blockStorageBName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: constants.BlockStorageInitialSize,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage B: %v", err)
	}

	replicaA, err := builders.NewInstanceBuilder().
		Name(replicaAName).
		Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.InstanceSpec{
			SkuRef:            *instanceSkuRefObj,
			Zone:              zoneA,
			AntiAffinityGroup: antiAffinityGroup,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageARefObj,
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Replica A Instance: %v", err)
	}

	replicaB, err := builders.NewInstanceBuilder().
		Name(replicaBName).
		Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.InstanceSpec{
			SkuRef:            *instanceSkuRefObj,
			Zone:              zoneB,
			AntiAffinityGroup: antiAffinityGroup,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageBRefObj,
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Replica B Instance: %v", err)
	}

	p := &params.HaMultiZoneV1Params{
		Workspace:     workspace,
		BlockStorageA: blockStorageA,
		BlockStorageB: blockStorageB,
		ReplicaA:      replicaA,
		ReplicaB:      replicaB,
	}
	suite.params = p
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockUsage.ConfigureHaMultiZoneV1, *p)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *HaMultiZoneV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t,
		sdkconsts.WorkspaceProviderV1Name,
		sdkconsts.StorageProviderV1Name,
		sdkconsts.ComputeProviderV1Name,
	)
	suite.ConfigureResources(t,
		string(schema.RegionalResourceMetadataKindResourceKindWorkspace),
		string(schema.RegionalResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalResourceMetadataKindResourceKindInstance),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", t, suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1,
		secapi.TenantReference{Tenant: secapi.TenantID(workspace.Metadata.Tenant), Name: workspace.Metadata.Name},
		steps.ResponseExpectsWithCondition[schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			Labels:   expectWorkspaceLabels,
			Metadata: expectWorkspaceMeta,
			ResourceStatus: schema.WorkspaceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	blockA := suite.params.BlockStorageA
	expectBlockAMeta := blockA.Metadata
	expectBlockASpec := &blockA.Spec
	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create the replica A boot volume", t, suite.Client.StorageV1, blockA,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:       expectBlockAMeta,
			Spec:           expectBlockASpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetBlockStorageV1Step("Get the replica A boot volume", suite.Client.StorageV1,
		secapi.WorkspaceReference{Tenant: secapi.TenantID(blockA.Metadata.Tenant), Workspace: secapi.WorkspaceID(blockA.Metadata.Workspace), Name: blockA.Metadata.Name},
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus]{
			Metadata: expectBlockAMeta,
			Spec:     expectBlockASpec,
			ResourceStatus: schema.BlockStorageStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	blockB := suite.params.BlockStorageB
	expectBlockBMeta := blockB.Metadata
	expectBlockBSpec := &blockB.Spec
	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create the replica B boot volume", t, suite.Client.StorageV1, blockB,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:       expectBlockBMeta,
			Spec:           expectBlockBSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	stepsBuilder.GetBlockStorageV1Step("Get the replica B boot volume", suite.Client.StorageV1,
		secapi.WorkspaceReference{Tenant: secapi.TenantID(blockB.Metadata.Tenant), Workspace: secapi.WorkspaceID(blockB.Metadata.Workspace), Name: blockB.Metadata.Name},
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus]{
			Metadata: expectBlockBMeta,
			Spec:     expectBlockBSpec,
			ResourceStatus: schema.BlockStorageStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Replica A - the response spec comparison round-trips both Zone and
	// AntiAffinityGroup, proving the API accepts and persists the grouping.
	replicaA := suite.params.ReplicaA
	expectReplicaAMeta := replicaA.Metadata
	expectReplicaASpec := &replicaA.Spec
	stepsBuilder.CreateOrUpdateInstanceV1Step("Create replica A", t, suite.Client.ComputeV1, replicaA,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:       expectReplicaAMeta,
			Spec:           expectReplicaASpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	replicaAWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(replicaA.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(replicaA.Metadata.Workspace),
		Name:      replicaA.Metadata.Name,
	}
	replicaA = stepsBuilder.GetInstanceV1Step("Get replica A", suite.Client.ComputeV1, replicaAWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			Metadata: expectReplicaAMeta,
			Spec:     expectReplicaASpec,
			ResourceStatus: schema.InstanceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	replicaB := suite.params.ReplicaB
	expectReplicaBMeta := replicaB.Metadata
	expectReplicaBSpec := &replicaB.Spec
	stepsBuilder.CreateOrUpdateInstanceV1Step("Create replica B", t, suite.Client.ComputeV1, replicaB,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:       expectReplicaBMeta,
			Spec:           expectReplicaBSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)
	replicaBWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(replicaB.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(replicaB.Metadata.Workspace),
		Name:      replicaB.Metadata.Name,
	}
	replicaB = stepsBuilder.GetInstanceV1Step("Get replica B", suite.Client.ComputeV1, replicaBWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			Metadata: expectReplicaBMeta,
			Spec:     expectReplicaBSpec,
			ResourceStatus: schema.InstanceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Bring both replicas up concurrently, as a real HA pair would run
	stepsBuilder.StartInstanceV1Step("Start replica A", suite.Client.ComputeV1, replicaA)
	replicaA = stepsBuilder.GetInstanceV1Step("Get the started replica A", suite.Client.ComputeV1, replicaAWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			Metadata: expectReplicaAMeta,
			Spec:     expectReplicaASpec,
			ResourceStatus: schema.InstanceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterStartingWithoutUpdate,
			},
		},
	)

	stepsBuilder.StartInstanceV1Step("Start replica B", suite.Client.ComputeV1, replicaB)
	replicaB = stepsBuilder.GetInstanceV1Step("Get the started replica B", suite.Client.ComputeV1, replicaBWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			Metadata: expectReplicaBMeta,
			Spec:     expectReplicaBSpec,
			ResourceStatus: schema.InstanceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterStartingWithoutUpdate,
			},
		},
	)

	// Resources deletion
	stepsBuilder.DeleteInstanceV1Step("Delete replica A", t, suite.Client.ComputeV1, replicaA)
	stepsBuilder.DeleteInstanceV1Step("Delete replica B", t, suite.Client.ComputeV1, replicaB)

	stepsBuilder.DeleteBlockStorageV1Step("Delete the replica A boot volume", t, suite.Client.StorageV1, blockA)
	stepsBuilder.DeleteBlockStorageV1Step("Delete the replica B boot volume", t, suite.Client.StorageV1, blockB)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)

	suite.FinishScenario()
}

func (suite *HaMultiZoneV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}

func indexOf(items []string, target string) int {
	for i, item := range items {
		if item == target {
			return i
		}
	}
	return 0
}
