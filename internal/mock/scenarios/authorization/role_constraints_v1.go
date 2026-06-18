package mockauthorization

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
)

// ConfigureRoleConstraintsValidationV1 sets up mock stubs for the role constraints validation suite.
// Each role in the params targets a different constraint Validation, all returning 422 Unprocessable Entity.
func ConfigureRoleConstraintsValidationV1(scenario *mockscenarios.Scenario, p params.RoleConstraintsValidationV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	// Over-length name Validation
	overLengthNameRole := p.OverLengthNameRole
	overLengthNameURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, overLengthNameRole.Metadata.Tenant, overLengthNameRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Invalid pattern name Validation
	invalidPatternNameRole := p.InvalidPatternNameRole
	invalidPatternNameURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, invalidPatternNameRole.Metadata.Tenant, invalidPatternNameRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(invalidPatternNameURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length label value Validation
	overLengthLabelRole := p.OverLengthLabelValueRole
	overLengthLabelURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, overLengthLabelRole.Metadata.Tenant, overLengthLabelRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthLabelURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length annotation value Validation
	overLengthAnnotationRole := p.OverLengthAnnotationRole
	overLengthAnnotationURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, overLengthAnnotationRole.Metadata.Tenant, overLengthAnnotationRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthAnnotationURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length permission provider Validation
	overLengthProviderRole := p.OverLengthPermissionProviderRole
	overLengthProviderURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, overLengthProviderRole.Metadata.Tenant, overLengthProviderRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthProviderURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length permission resource Validation
	overLengthResourceRole := p.OverLengthPermissionResourceRole
	overLengthResourceURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, overLengthResourceRole.Metadata.Tenant, overLengthResourceRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthResourceURL, scenario.MockParams); err != nil {
		return err
	}

	// Over-length permission verb Validation
	overLengthVerbRole := p.OverLengthPermissionVerbRole
	overLengthVerbURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, overLengthVerbRole.Metadata.Tenant, overLengthVerbRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overLengthVerbURL, scenario.MockParams); err != nil {
		return err
	}

	// Empty permissions Validation
	emptyPermissionsRole := p.EmptyPermissionsRole
	emptyPermissionsURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, emptyPermissionsRole.Metadata.Tenant, emptyPermissionsRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(emptyPermissionsURL, scenario.MockParams); err != nil {
		return err
	}

	// Over maxItems permissions Validation
	overMaxPermissionsRole := p.OverMaxItemsPermissionsRole
	overMaxPermissionsURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, overMaxPermissionsRole.Metadata.Tenant, overMaxPermissionsRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxPermissionsURL, scenario.MockParams); err != nil {
		return err
	}

	// Empty permission provider Validation
	emptyProviderRole := p.EmptyPermissionProviderRole
	emptyProviderURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, emptyProviderRole.Metadata.Tenant, emptyProviderRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(emptyProviderURL, scenario.MockParams); err != nil {
		return err
	}

	// Empty permission resources Validation
	emptyResourcesRole := p.EmptyPermissionResourcesRole
	emptyResourcesURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, emptyResourcesRole.Metadata.Tenant, emptyResourcesRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(emptyResourcesURL, scenario.MockParams); err != nil {
		return err
	}

	// Over maxItems permission resources Validation
	overMaxResourcesRole := p.OverMaxItemsPermissionResourcesRole
	overMaxResourcesURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, overMaxResourcesRole.Metadata.Tenant, overMaxResourcesRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxResourcesURL, scenario.MockParams); err != nil {
		return err
	}

	// Empty permission resource value Validation
	emptyResourceValueRole := p.EmptyPermissionResourceValueRole
	emptyResourceValueURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, emptyResourceValueRole.Metadata.Tenant, emptyResourceValueRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(emptyResourceValueURL, scenario.MockParams); err != nil {
		return err
	}

	// Empty permission verbs Validation
	emptyVerbsRole := p.EmptyPermissionVerbsRole
	emptyVerbsURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, emptyVerbsRole.Metadata.Tenant, emptyVerbsRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(emptyVerbsURL, scenario.MockParams); err != nil {
		return err
	}

	// Over maxItems permission verbs Validation
	overMaxVerbsRole := p.OverMaxItemsPermissionVerbsRole
	overMaxVerbsURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, overMaxVerbsRole.Metadata.Tenant, overMaxVerbsRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(overMaxVerbsURL, scenario.MockParams); err != nil {
		return err
	}

	// Empty permission verb value Validation
	emptyVerbValueRole := p.EmptyPermissionVerbValueRole
	emptyVerbValueURL := generators.GenerateRoleURL(sdkconsts.AuthorizationProviderV1Name, emptyVerbValueRole.Metadata.Tenant, emptyVerbValueRole.Metadata.Name)
	if err := configurator.ConfigurePutUnprocessableEntityStub(emptyVerbValueURL, scenario.MockParams); err != nil {
		return err
	}

	return scenario.FinishConfiguration(configurator)
}
