package secatest

import (
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
)

func generateSubnet(networkCidr string, size int) (string, error) {
	_, network, err := net.ParseCIDR(networkCidr)
	if err != nil {
		return "", err
	}

	subnet, _ := cidr.NextSubnet(network, size)

	return subnet.String(), nil
}

func generateIPAddress(subnetCidr string, num int) (string, error) {
	_, network, err := net.ParseCIDR(subnetCidr)
	if err != nil {
		return "", err
	}

	ip, err := cidr.Host(network, num)
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}
