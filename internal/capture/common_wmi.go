package capture

import (
	"github.com/yusufpapurcu/wmi"
)

// QueryWMI is a generic helper to run a WMI query and populate a struct slice.
// This replaces the deprecated 'wmic' CLI calls.
func QueryWMI(query string, dst interface{}) error {
	return wmi.Query(query, dst)
}

// QueryWMIWithNamespace allows querying specific namespaces like ROOT\CIMV2.
func QueryWMIWithNamespace(query string, dst interface{}, namespace string) error {
	return wmi.QueryNamespace(query, dst, namespace)
}
