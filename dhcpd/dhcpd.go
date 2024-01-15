package dhcp

import (
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// DHCP defines the interface for DHCP operations
type DHCP interface {
	FileExists() bool
	SetFaceAsDHCPOrRemove(networkInterface string) error
	SetFaceAsStatic(networkInterface, ip, subnet, gateway string) error
}

// dhcpImpl implements DHCP interface
type dhcpImpl struct {
	filePath string
}

// NewDHCP creates a new DHCP instance
func NewDHCP(filePath string) DHCP {
	if filePath == "" {
		filePath = "/etc/dhcpcd.conf"
	}
	return &dhcpImpl{filePath: filePath}
}

func (d *dhcpImpl) FileExists() bool {
	info, err := os.Stat(d.filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// SetFaceAsDHCPOrRemove removes an interface configuration from dhcpd.conf
func (d *dhcpImpl) SetFaceAsDHCPOrRemove(networkInterface string) error {
	if !isLinux() {
		return errors.New("RemoveInterface is only supported on Linux")
	}

	return removeInterface(d.filePath, networkInterface)
}

// SetFaceAsStatic sets a static IP for the interface in dhcpd.conf
func (d *dhcpImpl) SetFaceAsStatic(networkInterface, ip, subnet, gateway string) error {
	if !isLinux() {
		return errors.New("SetFaceAsStatic is only supported on Linux")
	}

	return setStaticIP(d.filePath, networkInterface, ip, subnet, gateway)
}

func removeInterface(filePath, networkInterface string) error {
	// Read the contents of the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	var skip bool

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "interface "+networkInterface) {
			// Start skipping lines until an empty line or a new interface block is encountered
			skip = true
			continue
		}
		if skip && (trimmedLine == "" || strings.HasPrefix(trimmedLine, "interface ")) {
			// Stop skipping lines
			skip = false
		}
		if !skip {
			newLines = append(newLines, line)
		}
	}

	// Write the updated content back to the file
	return os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644)
}

func setStaticIP(filePath, networkInterface, ip, subnet, gateway string) error {
	// First, remove any existing configuration for the interface
	err := validateIPSubnetGateway(ip, subnet, gateway)
	if err != nil {
		return err
	}
	if err := removeInterface(filePath, networkInterface); err != nil {
		return fmt.Errorf("error removing existing interface configuration: %w", err)
	}

	// Open the file in append mode
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Prepare the static IP configuration
	staticConfig := fmt.Sprintf("\n\ninterface %s\nstatic ip_address=%s/%s\nstatic routers=%s\nstatic domain_name_servers=%s\n",
		networkInterface, ip, subnet, gateway, gateway)

	// Write the new configuration to the file
	if _, err := file.WriteString(staticConfig); err != nil {
		return err
	}

	return nil
}

// validateIPSubnetGateway checks if the IP, subnet, and gateway are valid and in the same network.
func validateIPSubnetGateway(ipStr, subnetStr, gatewayStr string) error {
	// Parse the IP address
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return fmt.Errorf("invalid IP address: %s", ipStr)
	}

	// Convert subnet mask from dotted decimal to CIDR
	cidr, err := subnetMaskToCIDR(subnetStr)
	if err != nil {
		return err
	}

	// Parse the subnet
	_, subnet, err := net.ParseCIDR(ipStr + cidr)
	if err != nil {
		return fmt.Errorf("invalid subnet mask: %s", subnetStr)
	}

	// Parse the gateway IP address
	gateway := net.ParseIP(gatewayStr)
	if gateway == nil {
		return fmt.Errorf("invalid gateway IP address: %s", gatewayStr)
	}

	// Check if the gateway is in the same network as the IP
	if !subnet.Contains(gateway) {
		return fmt.Errorf("gateway IP %s is not in the same network as IP %s with subnet mask %s", gatewayStr, ipStr, subnetStr)
	}

	return nil
}

// subnetMaskToCIDR converts a dotted decimal subnet mask to CIDR notation
func subnetMaskToCIDR(mask string) (string, error) {
	maskParts := strings.Split(mask, ".")
	if len(maskParts) != 4 {
		return "", fmt.Errorf("invalid subnet mask format: %s", mask)
	}

	var cidrBits int
	for _, part := range maskParts {
		val, err := strconv.Atoi(part)
		if err != nil {
			return "", fmt.Errorf("invalid subnet mask value: %s", mask)
		}

		for i := 0; i < 8; i++ {
			if val&128 != 0 {
				cidrBits++
			}
			val <<= 1
		}
	}

	return fmt.Sprintf("/%d", cidrBits), nil
}

func isLinux() bool {
	return runtime.GOOS == "linux"
}
