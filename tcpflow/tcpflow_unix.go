// +build darwin freebsd

package tcpflow

// GetHostFlows gets host flows.
func GetHostFlows() (HostFlows, error) {
	return GetHostFlowsByProcfs()
}
