package mock

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"testing"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/stretchr/testify/assert"
)

func TestNetworkLifecycleScenarioV1(t *testing.T) {
	// TODO Calculate the subnet cidr from network cidr
	subnetCIDR := "10.1.0.0/24"
	subnetCIDRUpdate := "10.1.0.1/24"

	// TODO Calculate the nic cidr from subnet cidr
	nicAddress := "10.1.0.1"
	nicAddressUpdate := "10.1.0.2"

	// TODO Receive via configuration the public ip address range and calculate an ip address
	publicIPAddress := "192.168.0.1"

	// TODO Load from region
	zones := []string{"zone-a", "zone-b"}

	storageSkus := []string{"LD100"}

	// TODO Get from list instance skus endpoint
	instanceSkus := []string{"SXL"}

	// TODO Get from list network skus endpoint
	networkSkus := []string{"N1K"}
	// Select skus
	storageSkuName := storageSkus[rand.Intn(len(storageSkus))]
	instanceSkuName := instanceSkus[rand.Intn(len(instanceSkus))]
	networkSkuName := networkSkus[rand.Intn(len(networkSkus))]
	zone := zones[rand.Intn(len(zones))]

	// TODO Dynamically create before the scenario
	workspaceName := secalib.GenerateWorkspaceName()

	// Generate scenario data
	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)

	blockStorageName := secalib.GenerateBlockStorageName()
	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName)

	instanceSkuRef := secalib.GenerateSkuRef(instanceSkuName)
	instanceName := secalib.GenerateInstanceName()

	networkSkuRef := secalib.GenerateSkuRef(networkSkuName)
	networkName := secalib.GenerateNetworkName()
	networkRef := secalib.GenerateNetworkRef(networkName)

	internetGatewayName := secalib.GenerateInternetGatewayName()
	internetGatewayRef := secalib.GenerateInternetGatewayRef(internetGatewayName)

	routeTableName := secalib.GenerateRouteTableName()
	routeTableRef := secalib.GenerateRouteTableRef(routeTableName)

	subnetName := secalib.GenerateSubnetName()
	subnetRef := secalib.GenerateSubnetRef(subnetName)

	nicName := secalib.GenerateNicName()

	publicIPName := secalib.GeneratePublicIPName()
	publicIPRef := secalib.GeneratePublicIPRef(publicIPName)

	securityGroupName := secalib.GenerateSecurityGroupName()

	blockStorageSize := secalib.GenerateBlockStorageSize()

	// Generate URLs
	// networkSkuURL := WireMockURL + secalib.GenerateNetworkSkuURL(TenantName, networkSkuName)
	networkURL := WireMockURL + secalib.GenerateNetworkURL(TenantName, workspaceName, networkName)
	internetGatewayURL := WireMockURL + secalib.GenerateInternetGatewayURL(TenantName, workspaceName, internetGatewayName)
	nicURL := WireMockURL + secalib.GenerateNicURL(TenantName, workspaceName, nicName)
	publicIPURL := WireMockURL + secalib.GeneratePublicIPURL(TenantName, workspaceName, publicIPName)
	routeTableURL := WireMockURL + secalib.GenerateRouteTableURL(TenantName, workspaceName, routeTableName)
	subnetURL := WireMockURL + secalib.GenerateSubnetURL(TenantName, workspaceName, subnetName)
	securityGroupURL := WireMockURL + secalib.GenerateSecurityGroupURL(TenantName, workspaceName, securityGroupName)
	// instanceSkuURL := WireMockURL + secalib.GenerateInstanceSkuURL(TenantName, instanceSkuName)
	instanceURL := WireMockURL + secalib.GenerateInstanceURL(TenantName, workspaceName, instanceName)
	blockStorageURL := WireMockURL + secalib.GenerateBlockStorageURL(TenantName, workspaceName, blockStorageName)
	// storageSkuURL := WireMockURL + secalib.GenerateStorageSkuURL(TenantName, storageSkuName)

	networkParams := NetworkParamsV1{
		Params: &Params{
			MockURL:   WireMockURL,
			AuthToken: Token,
			Tenant:    TenantName,
			Region:    Region,
		},
		Network: &ResourceParams[secalib.NetworkSpecV1]{
			Name: networkName,

			InitialSpec: &secalib.NetworkSpecV1{
				Cidr: &secalib.NetworkSpecCIDRV1{
					Ipv4: "10.0.0.0",
				},
				SkuRef:        networkSkuRef,
				RouteTableRef: routeTableRef,
			},
			UpdatedSpec: &secalib.NetworkSpecV1{
				Cidr: &secalib.NetworkSpecCIDRV1{
					Ipv4: "10.0.0.1",
				},
				SkuRef:        networkSkuRef,
				RouteTableRef: routeTableRef,
			},
		},
		InternetGateway: &ResourceParams[secalib.InternetGatewaySpecV1]{
			Name: internetGatewayName,
			InitialSpec: &secalib.InternetGatewaySpecV1{
				EgressOnly: false,
			},
			UpdatedSpec: &secalib.InternetGatewaySpecV1{
				EgressOnly: true,
			},
		},
		RouteTable: &ResourceParams[secalib.RouteTableSpecV1]{
			Name: routeTableName,
			InitialSpec: &secalib.RouteTableSpecV1{
				LocalRef: networkRef,
				Routes: []*secalib.RouteTableRouteV1{
					{
						DestinationCidrBlock: "10.0.0.0/16",
						TargetRef:            internetGatewayRef,
					},
				},
			},
			UpdatedSpec: &secalib.RouteTableSpecV1{
				LocalRef: networkRef,
				Routes: []*secalib.RouteTableRouteV1{
					{
						DestinationCidrBlock: "10.0.0.1/16",
						TargetRef:            internetGatewayRef,
					},
				},
			},
		},
		Subnet: &ResourceParams[secalib.SubnetSpecV1]{
			Name: subnetName,
			InitialSpec: &secalib.SubnetSpecV1{
				Cidr: &secalib.SubnetSpecCIDRV1{Ipv4: subnetCIDR},
				Zone: zone,
			},
			UpdatedSpec: &secalib.SubnetSpecV1{
				Cidr: &secalib.SubnetSpecCIDRV1{Ipv4: subnetCIDRUpdate},
				Zone: zone,
			},
		},
		PublicIP: &ResourceParams[secalib.PublicIpSpecV1]{
			Name: publicIPName,
			InitialSpec: &secalib.PublicIpSpecV1{
				Version: secalib.IPVersion4,
				Address: publicIPAddress,
			},
		},
		NIC: &ResourceParams[secalib.NICSpecV1]{
			Name: nicName,
			InitialSpec: &secalib.NICSpecV1{
				Addresses:    []string{nicAddress},
				PublicIpRefs: []string{publicIPRef},
				SubnetRef:    subnetRef,
			},
			UpdatedSpec: &secalib.NICSpecV1{
				Addresses:    []string{nicAddressUpdate},
				PublicIpRefs: []string{publicIPRef},
				SubnetRef:    subnetRef,
			},
		},
		SecurityGroup: &ResourceParams[secalib.SecurityGroupSpecV1]{
			Name: securityGroupName,
			InitialSpec: &secalib.SecurityGroupSpecV1{
				Rules: []*secalib.SecurityGroupRuleV1{
					{
						Direction: secalib.SecurityRuleDirectionIngress,
					},
				},
			},
		},
		BlockStorage: &ResourceParams[secalib.BlockStorageSpecV1]{
			Name: blockStorageName,
			InitialSpec: &secalib.BlockStorageSpecV1{
				SkuRef: storageSkuRef,
				SizeGB: blockStorageSize,
			},
		},
		Instance: &ResourceParams[secalib.InstanceSpecV1]{
			Name: instanceName,
			InitialSpec: &secalib.InstanceSpecV1{
				SkuRef:        instanceSkuRef,
				Zone:          zone,
				BootDeviceRef: blockStorageRef,
			},
		},
	}

	wm, err := CreateNetworkLifecycleScenarioV1(fmt.Sprintf("Network lifecycle %d", rand.Intn(100)), networkParams)
	if err != nil {
		t.Errorf("Failed to create network scenario: %v\n", err)
		t.FailNow()
		return
	}
	if wm == nil {
		t.Errorf("Failed to create network scenario: %v\n", err)
		t.FailNow()
		return
	}
	/*
		// Continue with the test
		//Network Sku
		response, error := request("GET", networkSkuURL, Token)
		if error != nil {
			t.Errorf("Error Get NetworkSku: %v\n", error)
			t.FailNow()
			return
		}
		assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get NetworkSku")
	*/
	//Network
	// Create Network
	response, error := request("PUT", networkURL, Token)
	if error != nil {
		t.Errorf("Error Put Network: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusCreated, response.StatusCode, "Expected status code 201 for Put Network")

	// Get Network
	response, error = request("GET", networkURL, Token)
	if error != nil {
		t.Errorf("Error Get Network: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get Network with "+networkURL)

	// Update Network
	response, error = request("PUT", networkURL, Token)
	if error != nil {
		t.Errorf("Error Update Network: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Update Network with "+networkURL)

	// Get Network 2x
	response, error = request("GET", networkURL, Token)
	if error != nil {
		t.Errorf("Error Get Network: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get Network with "+networkURL)

	// internet-Gateway
	// Create internet-Gateway
	response, error = request("PUT", internetGatewayURL, Token)
	if error != nil {
		t.Errorf("Error Create internet-Gateway: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusCreated, response.StatusCode, "Expected status code 201 for Create internet-Gateway with "+internetGatewayURL)

	// Get internet-Gateway
	response, error = request("GET", internetGatewayURL, Token)
	if error != nil {
		t.Errorf("Error Get internet-Gateway: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get internet-Gateway with "+internetGatewayURL)

	// Update internet-Gateway
	response, error = request("PUT", internetGatewayURL, Token)
	if error != nil {
		t.Errorf("Error Update internet-Gateway: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Update internet-Gateway with "+internetGatewayURL)

	// Get internet-Gateway 2x
	response, error = request("GET", internetGatewayURL, Token)
	if error != nil {
		t.Errorf("Error Get internet-Gateway: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get internet-Gateway with "+internetGatewayURL)

	// Route Table
	// Create Route Table
	response, error = request("PUT", routeTableURL, Token)
	if error != nil {
		t.Errorf("Error Create Route Table: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusCreated, response.StatusCode, "Expected status code 201 for Create Route Table with "+routeTableURL)

	// Get Route Table
	response, error = request("GET", routeTableURL, Token)
	if error != nil {
		t.Errorf("Error Get Route Table: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get Route Table with "+routeTableURL)

	// Update Route Table
	response, error = request("PUT", routeTableURL, Token)
	if error != nil {
		t.Errorf("Error Update Route Table: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Update Route Table with "+routeTableURL)

	// Get Route Table 2x
	response, error = request("GET", routeTableURL, Token)
	if error != nil {
		t.Errorf("Error Get Route Table: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get Route Table with "+routeTableURL)

	// Subnet
	// Create Subnet
	response, error = request("PUT", subnetURL, Token)
	if error != nil {
		t.Errorf("Error Create Subnet: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusCreated, response.StatusCode, "Expected status code 201 for Create Subnet with "+subnetURL)
	// Get Subnet
	response, error = request("GET", subnetURL, Token)
	if error != nil {
		t.Errorf("Error Get Subnet: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get Subnet with "+subnetURL)
	// Update Subnet
	response, error = request("PUT", subnetURL, Token)
	if error != nil {
		t.Errorf("Error Update Subnet: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Update Subnet with "+subnetURL)
	// Get Subnet 2x
	response, error = request("GET", subnetURL, Token)
	if error != nil {
		t.Errorf("Error Get Subnet: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get Subnet with "+subnetURL)

	// Public-IP
	// Create Public-IP
	response, error = request("PUT", publicIPURL, Token)
	if error != nil {
		t.Errorf("Error Create Public-IP: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusCreated, response.StatusCode, "Expected status code 201 for Create Public-IP with "+publicIPURL)

	// Get Public-IP
	response, error = request("GET", publicIPURL, Token)
	if error != nil {
		t.Errorf("Error Get Public-IP: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get Public-IP with "+publicIPURL)

	// Update Public-IP
	response, error = request("PUT", publicIPURL, Token)
	if error != nil {
		t.Errorf("Error Update Public-IP: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Update Public-IP with "+publicIPURL)

	// Get Public-IP 2x
	response, error = request("GET", publicIPURL, Token)
	if error != nil {
		t.Errorf("Error Get Public-IP: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get Public-IP with "+publicIPURL)

	// NIC
	// Create NIC
	response, error = request("PUT", nicURL, Token)
	if error != nil {
		t.Errorf("Error Create NIC: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusCreated, response.StatusCode, "Expected status code 201 for Create NIC with "+nicURL)

	// Get NIC
	response, error = request("GET", nicURL, Token)
	if error != nil {
		t.Errorf("Error Get NIC: %v\n", error)
		t.FailNow()
		return
	}

	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get NIC with "+nicURL)

	// Update NIC
	response, error = request("PUT", nicURL, Token)
	if error != nil {
		t.Errorf("Error Update NIC: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Update NIC with "+nicURL)

	// Get NIC 2x
	response, error = request("GET", nicURL, Token)
	if error != nil {
		t.Errorf("Error Get NIC: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get NIC with "+nicURL)

	// Security Group
	// Create Security Group
	response, error = request("PUT", securityGroupURL, Token)
	if error != nil {
		log.Printf("Error Create Security Group: %v\n", error)
		return
	}
	assert.Equal(t, http.StatusCreated, response.StatusCode, "Expected status code 201 for Create Security Group with "+securityGroupURL)

	// Get Security Group
	response, error = request("GET", securityGroupURL, Token)
	if error != nil {
		log.Printf("Error Get Security Group: %v\n", error)
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get Security Group")

	// Update Security Group
	response, error = request("PUT", securityGroupURL, Token)
	if error != nil {
		log.Printf("Error Update Security Group: %v\n", error)
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Update Security Group")

	// Get Security Group 2x
	response, error = request("GET", securityGroupURL, Token)
	if error != nil {
		log.Printf("Error Get Security Group: %v\n", error)
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get Security Group")

	// Compute
	// Get Storage Sku
	/*
		response, error = request("GET", storageSkuURL, Token)
		if error != nil {
			log.Printf("Error Create Block Storage: %v\n", error)
			return
		}
		assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get Storage Sku WITH "+storageSkuURL)
	*/
	// Create Block Storage
	response, error = request("PUT", blockStorageURL, Token)
	if error != nil {
		log.Printf("Error Create Block Storage: %v\n", error)
		return
	}
	assert.Equal(t, http.StatusCreated, response.StatusCode, "Expected status code 201 for Create Block Storage WITH "+blockStorageURL)
	// Get Block Storage
	response, error = request("GET", blockStorageURL, Token)
	if error != nil {
		log.Printf("Error Get Block Storage: %v\n", error)
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get Block Storage WITH "+blockStorageURL)
	/*
		// Get instance sku
		response, error = request("GET", instanceSkuURL, Token)
		if error != nil {
			log.Printf("Error Get Instance SKU: %v\n", error)
			return
		}
		assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get Instance SKU WITH "+instanceSkuURL)
	*/
	// Create Instance
	response, error = request("PUT", instanceURL, Token)
	if error != nil {
		log.Printf("Error Create Instance: %v\n", error)
		return
	}
	assert.Equal(t, http.StatusCreated, response.StatusCode, "Expected status code 201 for Create Instance")
	// Get Instance
	response, error = request("GET", instanceURL, Token)
	if error != nil {
		log.Printf("Error Get Instance: %v\n", error)
		return
	}
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code 200 for Get Instance")

	// Delete Instance
	response, error = request("DELETE", instanceURL, Token)
	if error != nil {
		t.Errorf("Error Delete Instance: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusAccepted, response.StatusCode, "Expected status code 204 for Delete Instance with "+instanceURL)

	// Delete Block Storage
	response, error = request("DELETE", blockStorageURL, Token)
	if error != nil {
		t.Errorf("Error Delete Block Storage: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusAccepted, response.StatusCode, "Expected status code 204 for Delete Block Storage with "+blockStorageURL)

	// Delete Security Group
	response, error = request("DELETE", securityGroupURL, Token)
	if error != nil {
		t.Errorf("Error Delete Security Group: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusAccepted, response.StatusCode, "Expected status code 204 for Delete Security Group with "+securityGroupURL)

	// Delete Nic
	response, error = request("DELETE", nicURL, Token)
	if error != nil {
		t.Errorf("Error Delete Nic: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusAccepted, response.StatusCode, "Expected status code 204 for Delete Nic with "+nicURL)

	// Delete public ip
	response, error = request("DELETE", publicIPURL, Token)
	if error != nil {
		t.Errorf("Error Delete Public IP: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusAccepted, response.StatusCode, "Expected status code 204 for Delete Public IP with "+publicIPURL)

	// Delete subnet
	response, error = request("DELETE", subnetURL, Token)
	if error != nil {
		t.Errorf("Error Delete Subnet: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusAccepted, response.StatusCode, "Expected status code 204 for Delete Subnet with "+subnetURL)

	// Delete Route-table
	response, error = request("DELETE", routeTableURL, Token)
	if error != nil {
		t.Errorf("Error Delete Route-table: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusAccepted, response.StatusCode, "Expected status code 204 for Delete Route-table with "+routeTableURL)

	// Delete Internet-gateway
	response, error = request("DELETE", internetGatewayURL, Token)
	if error != nil {
		t.Errorf("Error Delete Internet-gateway: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusAccepted, response.StatusCode, "Expected status code 204 for Delete Internet-gateway with "+internetGatewayURL)

	// Delete Network
	response, error = request("DELETE", networkURL, Token)
	if error != nil {
		t.Errorf("Error Delete Network: %v\n", error)
		t.FailNow()
		return
	}
	assert.Equal(t, http.StatusAccepted, response.StatusCode, "Expected status code 204 for Delete Network with "+networkURL)
}

func request(method string, url string, token string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Printf("Error creating request for %s: %v\n", url, err)
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending PUT request to %s: %v\n", url, err)
		return nil, err
	}

	return resp, nil
}
