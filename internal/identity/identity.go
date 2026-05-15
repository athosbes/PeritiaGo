package identity

import (
	"net"
	"os"
	"os/user"
	"strings"

	"github.com/athosbes/PeritiaGo/internal/capture"
	"github.com/athosbes/PeritiaGo/internal/models"
	"golang.org/x/sys/windows/registry"
)

// Win32_ComputerSystemProduct represents the WMI class for computer hardware.
type Win32_ComputerSystemProduct struct {
	Vendor            string
	Name              string
	IdentifyingNumber string
	UUID              string
}

// Win32_ComputerSystem represents the WMI class for computer system info.
type Win32_ComputerSystem struct {
	Domain string
}

// GetFullIdentity gathers all requested identification data for the machine.
func GetFullIdentity() models.MachineIdentity {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-hostname"
	}

	currUser, err := user.Current()
	username := "unknown-user"
	if err == nil {
		username = currUser.Username
	}

	id := models.MachineIdentity{
		Hostname:     hostname,
		CurrentUser:  username,
		MachineGUID:  GetMachineUUID(),
		IPAddresses:  getIPAddresses(),
		MACAddresses: getMACAddresses(),
	}

	// Capture OS Info
	id.OSName, id.OSVersion, id.OSBuild = getOSInfo()

	// Capture Hardware Info
	id.Manufacturer, id.Model, id.SerialNumber, id.BIOSUUID = getHardwareInfo()

	// Capture Domain/Workgroup
	id.Domain = getDomainInfo()

	return id
}

// GetMachineUUID retrieves the MachineGuid from the registry.
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

func getOSInfo() (name, version, build string) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err == nil {
		defer k.Close()
		name, _, _ = k.GetStringValue("ProductName")
		version, _, _ = k.GetStringValue("DisplayVersion")
		if version == "" {
			version, _, _ = k.GetStringValue("ReleaseId")
		}
		build, _, _ = k.GetStringValue("CurrentBuild")
	}
	return
}

func getHardwareInfo() (vendor, model, serial, uuid string) {
	var dst []Win32_ComputerSystemProduct
	query := "SELECT Vendor, Name, IdentifyingNumber, UUID FROM Win32_ComputerSystemProduct"
	err := capture.QueryWMI(query, &dst)
	if err == nil && len(dst) > 0 {
		vendor = dst[0].Vendor
		model = dst[0].Name
		serial = dst[0].IdentifyingNumber
		uuid = dst[0].UUID
	}
	return
}

func getDomainInfo() string {
	var dst []Win32_ComputerSystem
	query := "SELECT Domain FROM Win32_ComputerSystem"
	err := capture.QueryWMI(query, &dst)
	if err == nil && len(dst) > 0 {
		return dst[0].Domain
	}
	return "WORKGROUP"
}

func getIPAddresses() []string {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ips = append(ips, ipnet.IP.String())
				}
			}
		}
	}
	return ips
}

func getMACAddresses() map[string]string {
	macs := make(map[string]string)
	ifcs, err := net.Interfaces()
	if err == nil {
		for _, ifc := range ifcs {
			if ifc.HardwareAddr != nil {
				macs[ifc.Name] = ifc.HardwareAddr.String()
			}
		}
	}
	return macs
}

// GetMACAddress returns the MAC address of the first non-loopback active network interface. (Kept for compatibility)
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
