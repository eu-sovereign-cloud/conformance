package secatest

import (
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
)

func generateSubnetCidr(networkCidr string, size int, netNum int) (string, error) {
	_, network, err := net.ParseCIDR(networkCidr)
	if err != nil {
		return "", err
	}

	subnet, err := cidr.Subnet(network, size, netNum)
	if err != nil {
		return "", err
	}

	return subnet.String(), nil
}

func generateNicAddress(subnetCidr string, hostNum int) (string, error) {
	_, network, err := net.ParseCIDR(subnetCidr)
	if err != nil {
		return "", err
	}

	ip, err := cidr.Host(network, hostNum)
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}

func generatePublicIp(publicIpRange string, hostNum int) (string, error) {
	_, network, err := net.ParseCIDR(publicIpRange)
	if err != nil {
		return "", err
	}

	ip, err := cidr.Host(network, hostNum)
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}
