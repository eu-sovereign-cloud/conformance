package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/network"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestNetworkV1Suites(t *testing.T) {
	regionalTestSuite := suites.CreateRegionalTestSuite(config.Parameters, config.Clients)

	// Provider LifeCycle Suite
	providerLifeCycleSuite := network.CreateProviderLifeCycleV1TestSuite(regionalTestSuite,
		&network.ProviderLifeCycleV1Config{
			NetworkCidr:    config.Parameters.ScenariosCidr,
			PublicIpsRange: config.Parameters.ScenariosPublicIps,
			RegionZones:    config.Clients.RegionZones,
			InstanceSkus:   config.Clients.InstanceSkus,
			StorageSkus:    config.Clients.StorageSkus,
			NetworkSkus:    config.Clients.NetworkSkus,
		},
	)
	if providerLifeCycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, providerLifeCycleSuite)
	}

	// Provider Queries Suite
	providerQueriesSuite := network.CreateProviderQueriesV1TestSuite(regionalTestSuite,
		&network.ProviderQueriesV1Config{
			NetworkCidr:    config.Parameters.ScenariosCidr,
			PublicIpsRange: config.Parameters.ScenariosPublicIps,
			RegionZones:    config.Clients.RegionZones,
			InstanceSkus:   config.Clients.InstanceSkus,
			StorageSkus:    config.Clients.StorageSkus,
			NetworkSkus:    config.Clients.NetworkSkus,
		},
	)
	if providerQueriesSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, providerQueriesSuite)
	}

	// Network Lifecycle Suite
	networkLifecycleSuite := network.CreateNetworkLifeCycleV1TestSuite(regionalTestSuite,
		&network.NetworkLifeCycleV1Config{
			NetworkCidr: config.Parameters.ScenariosCidr,
			NetworkSkus: config.Clients.NetworkSkus,
		},
	)
	if networkLifecycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, networkLifecycleSuite)
	}

	// Nic Lifecycle Suite
	nicLifecycleSuite := network.CreateNicLifeCycleV1TestSuite(regionalTestSuite,
		&network.NicLifeCycleV1Config{
			NetworkCidr:    config.Parameters.ScenariosCidr,
			PublicIpsRange: config.Parameters.ScenariosPublicIps,
			RegionZones:    config.Clients.RegionZones,
			NetworkSkus:    config.Clients.NetworkSkus,
		},
	)
	if nicLifecycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, nicLifecycleSuite)
	}

	// Route Table Lifecycle Suite
	routeTableLifecycleSuite := network.CreateRouteTableLifeCycleV1TestSuite(regionalTestSuite,
		&network.RouteTableLifeCycleV1Config{
			NetworkCidr:    config.Parameters.ScenariosCidr,
			PublicIpsRange: config.Parameters.ScenariosPublicIps,
			RegionZones:    config.Clients.RegionZones,
			NetworkSkus:    config.Clients.NetworkSkus,
		},
	)
	if routeTableLifecycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, routeTableLifecycleSuite)
	}

	// Internet Gateway Lifecycle Suite
	internetGatewayLifecycleSuite := network.CreateInternetGatewayLifeCycleV1TestSuite(regionalTestSuite)
	if internetGatewayLifecycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, internetGatewayLifecycleSuite)
	}

	// Subnet Lifecycle Suite
	subnetLifecycleSuite := network.CreateSubnetLifeCycleV1TestSuite(regionalTestSuite,
		&network.SubnetLifeCycleV1Config{
			NetworkCidr: config.Parameters.ScenariosCidr,
			RegionZones: config.Clients.RegionZones,
			NetworkSkus: config.Clients.NetworkSkus,
		},
	)
	if subnetLifecycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, subnetLifecycleSuite)
	}

	// Public IP Lifecycle Suite
	publicIpLifecycleSuite := network.CreatePublicIpLifeCycleV1TestSuite(regionalTestSuite,
		&network.PublicIpLifeCycleV1Config{
			PublicIpsRange: config.Parameters.ScenariosPublicIps,
			NetworkSkus:    config.Clients.NetworkSkus,
		},
	)
	if publicIpLifecycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, publicIpLifecycleSuite)
	}

	// Security Group Rule Lifecycle Suite
	securityGroupRuleLifecycleSuite := network.CreateSecurityGroupRuleLifeCycleV1TestSuite(regionalTestSuite)
	if securityGroupRuleLifecycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, securityGroupRuleLifecycleSuite)
	}

	// Security Group Lifecycle Suite
	securityGroupLifecycleSuite := network.CreateSecurityGroupLifeCycleV1TestSuite(regionalTestSuite)
	if securityGroupLifecycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, securityGroupLifecycleSuite)
	}

	// Network Constraints Suite
	networkConstraintsSuite := network.CreateNetworkConstraintsValidationV1TestSuite(regionalTestSuite,
		&network.NetworkLifeCycleV1Config{
			NetworkCidr: config.Parameters.ScenariosCidr,
			NetworkSkus: config.Clients.NetworkSkus,
		},
	)
	if networkConstraintsSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, networkConstraintsSuite)
	}

	// Internet Gateway Constraints Suite
	internetGatewayConstraintsSuite := network.CreateInternetGatewayConstraintsValidationV1TestSuite(regionalTestSuite)
	if internetGatewayConstraintsSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, internetGatewayConstraintsSuite)
	}

	// Public IP Constraints Suite
	publicIpConstraintsSuite := network.CreatePublicIpConstraintsValidationV1TestSuite(regionalTestSuite,
		&network.PublicIpLifeCycleV1Config{
			PublicIpsRange: config.Parameters.ScenariosPublicIps,
			NetworkSkus:    config.Clients.NetworkSkus,
		},
	)
	if publicIpConstraintsSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, publicIpConstraintsSuite)
	}

	// Nic Constraints Suite
	nicConstraintsSuite := network.CreateNicConstraintsValidationV1TestSuite(regionalTestSuite,
		&network.NicLifeCycleV1Config{
			NetworkCidr:    config.Parameters.ScenariosCidr,
			PublicIpsRange: config.Parameters.ScenariosPublicIps,
			RegionZones:    config.Clients.RegionZones,
			NetworkSkus:    config.Clients.NetworkSkus,
		},
	)
	if nicConstraintsSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, nicConstraintsSuite)
	}

	// Security Group Constraints Suite
	securityGroupConstraintsSuite := network.CreateSecurityGroupConstraintsValidationV1TestSuite(regionalTestSuite)
	if securityGroupConstraintsSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, securityGroupConstraintsSuite)
	}

	// Security Group Rule Constraints Suite
	securityGroupRuleConstraintsSuite := network.CreateSecurityGroupRuleConstraintsValidationV1TestSuite(regionalTestSuite)
	if securityGroupRuleConstraintsSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, securityGroupRuleConstraintsSuite)
	}

	// Route Table Constraints Suite
	routeTableConstraintsSuite := network.CreateRouteTableConstraintsValidationV1TestSuite(regionalTestSuite,
		&network.RouteTableLifeCycleV1Config{
			NetworkCidr:    config.Parameters.ScenariosCidr,
			PublicIpsRange: config.Parameters.ScenariosPublicIps,
			RegionZones:    config.Clients.RegionZones,
			NetworkSkus:    config.Clients.NetworkSkus,
		},
	)
	if routeTableConstraintsSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, routeTableConstraintsSuite)
	}

	// Subnet Constraints Suite
	subnetConstraintsSuite := network.CreateSubnetConstraintsValidationV1TestSuite(regionalTestSuite,
		&network.SubnetLifeCycleV1Config{
			NetworkCidr: config.Parameters.ScenariosCidr,
			RegionZones: config.Clients.RegionZones,
			NetworkSkus: config.Clients.NetworkSkus,
		},
	)
	if subnetConstraintsSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, subnetConstraintsSuite)
	}
}
