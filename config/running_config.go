package config


// RunningConfig is the config instance kept by server while running.
// Read only in most cases.
// Note: It is better to implement RunningConfig as in memory struct
// 		 than database or file.
type RunningConfig interface {
	Host() string
	Port() uint16
	RequestOutput() bool // Whether to output request info.
	HostOnly() bool 	 // No port in built url if true.
	HostOnlyOn()		 // Set HostOnly on
	HostOnlyOff()		 // Set HostOnly off
}

type rConfig struct {
	host string
	port uint16
	//requestOutput bool
	hostOnly bool
}

func (r *rConfig) Host() string {
	return r.host
}

func (r *rConfig) Port() uint16 {
	return r.port
}

// TODO: The ability to config this.
func (r *rConfig) RequestOutput() bool {
	return true
}

func (r *rConfig) HostOnly() bool {
	return r.hostOnly
}

func (r *rConfig) HostOnlyOn() {
	r.hostOnly = true
}

func (r *rConfig) HostOnlyOff() {
	r.hostOnly = false
}