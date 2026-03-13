package identity

import (
	"net"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// GetMachineUUID retrieves the MachineGuid from the registry, which is a unique ID for the Windows installation.
func GetMachineUUID() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE)
	if err != nil {
		return "Unknown-UUID"
	}
	defer k.Close()

	guid, _, err := k.GetStringValue("MachineGuid")
	if err != nil {
		return "Unknown-UUID"
	}
	return guid
}

// GetMACAddress returns the MAC address of the first non-loopback active network interface.
func GetMACAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "Unknown-MAC"
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.HardwareAddr != nil && !strings.Contains(iface.Name, "loopback") {
			addr := iface.HardwareAddr.String()
			return strings.ReplaceAll(strings.ReplaceAll(addr, ":", ""), "-", "")
		}
	}

	return "Unknown-MAC"
}
