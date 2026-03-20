package identity

import (
	"net"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/athosbes/PeritiaGo/internal/models"
	"golang.org/x/sys/windows/registry"
)

// GetFullIdentity gathers all requested identification data for the machine.
func GetFullIdentity() models.MachineIdentity {
	hostname, _ := os.Hostname()
	currUser, _ := user.Current()

	id := models.MachineIdentity{
		Hostname:     hostname,
		CurrentUser:  currUser.Username,
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
	// Using wmic as a fallback for complex hardware info
	out, err := exec.Command("wmic", "csproduct", "get", "Vendor,Name,IdentifyingNumber,UUID", "/format:list").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				val := strings.TrimSpace(parts[1])
				switch strings.TrimSpace(parts[0]) {
				case "Vendor":
					vendor = val
				case "Name":
					model = val
				case "IdentifyingNumber":
					serial = val
				case "UUID":
					uuid = val
				}
			}
		}
	}
	return
}

func getDomainInfo() string {
	out, err := exec.Command("wmic", "computersystem", "get", "domain", "/format:list").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.Contains(line, "=") {
				return strings.TrimSpace(strings.Split(line, "=")[1])
			}
		}
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
