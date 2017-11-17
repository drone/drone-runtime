package engine

type (
	// Config defines the runtime configuration.
	Config struct {
		Version  int64      `json:"version"`
		Stages   []*Stage   `json:"pipeline"` // pipeline stages
		Networks []*Network `json:"networks"` // network definitions
		Volumes  []*Volume  `json:"volumes"`  // volume definitions
	}

	// Stage denotes a collection of one or more steps.
	Stage struct {
		Name  string  `json:"name,omitempty"`
		Alias string  `json:"alias,omitempty"`
		Steps []*Step `json:"steps,omitempty"`
	}

	// Step defines a container process.
	Step struct {
		Name         string            `json:"name"`
		Alias        string            `json:"alias,omitempty"`
		Image        string            `json:"image,omitempty"`
		Pull         bool              `json:"pull,omitempty"`
		Detached     bool              `json:"detach,omitempty"`
		Privileged   bool              `json:"privileged,omitempty"`
		WorkingDir   string            `json:"working_dir,omitempty"`
		Secrets      []*Secret         `json:"secrets,omitempty"`
		Environment  map[string]string `json:"environment,omitempty"`
		Labels       map[string]string `json:"labels,omitempty"`
		Entrypoint   []string          `json:"entrypoint,omitempty"`
		Command      []string          `json:"command,omitempty"`
		ExtraHosts   []string          `json:"extra_hosts,omitempty"`
		Volumes      []*VolumeMapping  `json:"volumes,omitempty"`
		Tmpfs        []string          `json:"tmpfs,omitempty"`
		Devices      []*DeviceMapping  `json:"devices,omitempty"`
		Networks     []*NetworkMapping `json:"networks,omitempty"`
		DNS          []string          `json:"dns,omitempty"`
		DNSSearch    []string          `json:"dns_search,omitempty"`
		MemSwapLimit int64             `json:"memswap_limit,omitempty"`
		MemLimit     int64             `json:"mem_limit,omitempty"`
		ShmSize      int64             `json:"shm_size,omitempty"`
		CPUQuota     int64             `json:"cpu_quota,omitempty"`
		CPUShares    int64             `json:"cpu_shares,omitempty"`
		CPUSet       string            `json:"cpu_set,omitempty"`
		ErrIgnore    bool              `json:"err_ignore,omitempty"`
		OnFailure    bool              `json:"on_failure,omitempty"`
		OnSuccess    bool              `json:"on_success,omitempty"`
		AuthConfig   Auth              `json:"auth_config,omitempty"`
		NetworkMode  string            `json:"network_mode,omitempty"`
		IpcMode      string            `json:"ipc_mode,omitempty"`
		Exports      []*File           `json:"exports,omitempty"`
		Sysctls      map[string]string `json:"sysctls,omitempty"`
		Backup       []*Snapshot       `json:"backup,omitempty"`
		Restore      []*Snapshot       `json:"restore,omitempty"`
	}

	// Auth defines registry authentication credentials.
	Auth struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
		Email    string `json:"email,omitempty"`
	}

	// NetworkMapping defines a container network mapping.
	NetworkMapping struct {
		Name    string   `json:"name"`
		Aliases []string `json:"aliases"`
	}

	// Network defines a container network.
	Network struct {
		Name       string            `json:"name,omitempty"`
		Driver     string            `json:"driver,omitempty"`
		DriverOpts map[string]string `json:"driver_opts,omitempty"`
	}

	// Volume defines a container volume.
	Volume struct {
		Name       string            `json:"name,omitempty"`
		Driver     string            `json:"driver,omitempty"`
		DriverOpts map[string]string `json:"driver_opts,omitempty"`
	}

	// VolumeMapping describes a volume mapping.
	VolumeMapping struct {
		Name   string `json:"name"`
		Source string `json:"source"`
		Target string `json:"target"`
	}

	// DeviceMapping describes a device mapping.
	DeviceMapping struct {
		Source string `json:"source"`
		Target string `json:"target"`
	}

	// Secret defines a runtime secret
	Secret struct {
		Name  string `json:"name,omitempty"`  // Secret name
		Value string `json:"value,omitempty"` // Secret value
		Mount string `json:"mount,omitempty"` // Secrets are mounted as a file
		Mask  bool   `json:"mask,omitempty"`  // Secrets are masked in output
	}

	// File defines a file exported from the container.
	File struct {
		Path string `json:"path,omitempty"` // File path
		Mime string `json:"mime,omitempty"` // File mime type
	}

	// FileInfo defines a file stat
	FileInfo struct {
		Name  string // File path
		Path  string // File path
		Size  int64  // File size
		Time  int64  // File time
		Mime  string // File mime type
		IsDir bool   // File is directory
	}

	// State represents the container state.
	State struct {
		ExitCode  int  // Container exit code
		Exited    bool // Container exited
		OOMKilled bool // Container is oom killed
	}

	// Snapshot defines a container volume snapshot
	Snapshot struct {
		Data   []byte `json:"data,omitempty"`
		Source string `json:"source,omitempty"`
		Target string `json:"target,omitempty"`
	}
)
