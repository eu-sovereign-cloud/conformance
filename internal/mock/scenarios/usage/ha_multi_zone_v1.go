package mockusage

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

func ConfigureHaMultiZoneV1(scenario *mockscenarios.Scenario, p params.HaMultiZoneV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	workspace := p.Workspace
	blockA := p.BlockStorageA
	blockB := p.BlockStorageB
	replicaA := p.ReplicaA
	replicaB := p.ReplicaB

	workspaceUrl := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
	blockAUrl := generators.GenerateBlockStorageURL(sdkconsts.StorageProviderV1Name, blockA.Metadata.Tenant, blockA.Metadata.Workspace, blockA.Metadata.Name)
	blockBUrl := generators.GenerateBlockStorageURL(sdkconsts.StorageProviderV1Name, blockB.Metadata.Tenant, blockB.Metadata.Workspace, blockB.Metadata.Name)
	replicaAUrl := generators.GenerateInstanceURL(sdkconsts.ComputeProviderV1Name, replicaA.Metadata.Tenant, replicaA.Metadata.Workspace, replicaA.Metadata.Name)
	replicaAStartUrl := generators.GenerateInstanceStartURL(sdkconsts.ComputeProviderV1Name, replicaA.Metadata.Tenant, replicaA.Metadata.Workspace, replicaA.Metadata.Name)
	replicaBUrl := generators.GenerateInstanceURL(sdkconsts.ComputeProviderV1Name, replicaB.Metadata.Tenant, replicaB.Metadata.Workspace, replicaB.Metadata.Name)
	replicaBStartUrl := generators.GenerateInstanceStartURL(sdkconsts.ComputeProviderV1Name, replicaB.Metadata.Tenant, replicaB.Metadata.Workspace, replicaB.Metadata.Name)

	if err := configurator.ConfigureCreateWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateBlockStorageStub(blockA, blockAUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingBlockStorageStub(blockA, blockAUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockA, blockAUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateBlockStorageStub(blockB, blockBUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingBlockStorageStub(blockB, blockBUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveBlockStorageStub(blockB, blockBUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateInstanceStub(replicaA, replicaAUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingInstanceStub(replicaA, replicaAUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInstanceStub(replicaA, replicaAUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := configurator.ConfigureCreateInstanceStub(replicaB, replicaBUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetCreatingInstanceStub(replicaB, replicaBUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInstanceStub(replicaB, replicaBUrl, scenario.MockParams); err != nil {
		return err
	}

	// Start both replicas
	if err := configurator.ConfigureInstanceOperationStub(replicaA, replicaAStartUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInstanceStub(replicaA, replicaAUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureInstanceOperationStub(replicaB, replicaBStartUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureGetActiveInstanceStub(replicaB, replicaBUrl, scenario.MockParams); err != nil {
		return err
	}

	// Delete all
	if err := configurator.ConfigureDeleteStub(replicaAUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(replicaBUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(blockAUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(blockBUrl, scenario.MockParams); err != nil {
		return err
	}
	if err := configurator.ConfigureDeleteStub(workspaceUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
