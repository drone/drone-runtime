package engine

type (
	// Metadata provides execution metadata.
	Metadata struct {
		UID       string            `json:"uid"`
		Namespace string            `json:"namespace"`
		Name      string            `json:"name"`
		Labels    map[string]string `json:"labels"`
	}

	// Spec provides the pipeline spec. This provides the
	// required instructions for reproducable pipeline
	// execution.
	Spec struct {
		Metadata Metadata  `json:"metadata"`
		Secrets  []*Secret `json:"secrets"`
		Steps    []*Step   `json:"steps"`
		Files    []*File   `json:"files"`

		// Docker-specific settings. These settings are
		// only used by the Docker and Kubernetes runtime
		// drivers.
		Docker *DockerConfig `json:"docker,omitempty"`

		// Qemu-specific settings. These settings are only
		// used by the qemu runtime driver.
		Qemu *QemuConfig `json:"qemu,omitempty"`

		// VMWare Fusion settings. These settings are only
		// used by the VMWare runtime driver.
		Fusion *FusionConfig `json:"fusion,omitempty"`
	}

	// Step defines a pipeline step.
	Step struct {
		Metadata     Metadata          `json:"metadata"`
		Detach       bool              `json:"detach"`
		DependsOn    []string          `json:"depends_on"`
		Devices      []*VolumeDevice   `json:"devices"`
		Docker       *DockerStep       `json:"docker"`
		Envs         map[string]string `json:"environment"`
		Files        []*FileMount      `json:"files"`
		IgnoreErr    bool              `json:"ignore_err"`
		IgnoreStdout bool              `json:"ignore_stderr"`
		IgnoreStderr bool              `json:"ignore_stdout"`
		Resources    *Resources        `json:"resources"`
		RunPolicy    RunPolicy         `json:"run_policy"`
		Secrets      []string          `json:"secrets"`
		Volumes      []*VolumeMount    `json:"volumes"`
		WorkingDir   string            `json:"working_dir"`
	}

	// DockerAuth defines dockerhub authentication credentials.
	DockerAuth struct {
		Address  string `json:"address"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// DockerConfig configures a Docker-based pipeline.
	DockerConfig struct {
		Auths   []*DockerAuth `json:"auths"`
		Volumes []*Volume     `json:"volumes"`
	}

	// DockerStep configures a docker step.
	DockerStep struct {
		Args       []string   `json:"args"`
		Command    []string   `json:"command"`
		Image      string     `json:"image"`
		Networks   []string   `json:"networks"`
		Ports      []*Port    `json:"ports"`
		Privileged bool       `json:"privileged"`
		PullPolicy PullPolicy `json:"pull_policy"`
	}

	// File defines a file that should be uploaded or
	// mounted somewhere in the step container or virtual
	// machine prior to command execution.
	File struct {
		Name string
		Data []byte
	}

	// FileMount defines how a file resource should be
	// mounted or included in the runtime environment.
	FileMount struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Mode int64  `json:"mode"`
	}

	// FusionConfig configures a VMWare Fusion-based pipeline.
	FusionConfig struct {
		Image string `json:"image"`
	}

	// Platform defines the target platform.
	Platform struct {
		OS      string `json:"os"`
		Arch    string `json:"arch"`
		Variant string `json:"variant"`
		Version string `json:"version"`
	}

	// Port represents a network port in a single container.
	Port struct {
		Port     int    `json:"port"`
		Host     int    `json:"host"`
		Protocol string `json:"protocol"`
	}

	// QemuConfig configures a Qemu-based pipeline.
	QemuConfig struct {
		Image string `json:"image"`
	}

	// Resources describes the compute resource
	// requirements.
	Resources struct {
		// Limits describes the maximum amount of compute
		// resources allowed.
		Limits *ResourceObject `json:"limits"`

		// Requests describes the minimum amount of
		// compute resources required.
		Requests *ResourceObject `json:"requests"`
	}

	// ResourceObject describes compute resource
	// requirements.
	ResourceObject struct {
		CPU    int64 `json:"cpu"`
		Memory int64 `json:"memory"`
	}

	// Secret represents a secret variable.
	Secret struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	// State represents the container state.
	State struct {
		ExitCode  int  // Container exit code
		Exited    bool // Container exited
		OOMKilled bool // Container is oom killed
	}

	// Volume that can be mounted by containers.
	Volume struct {
		Metadata Metadata        `json:"metadata"`
		EmptyDir *VolumeEmptyDir `json:"temp"`
		HostPath *VolumeHostPath `json:"host"`
	}

	// VolumeDevice describes a mapping of a raw block
	// device within a container.
	VolumeDevice struct {
		Name       string `json:"name"`
		DevicePath string `json:"path"`
	}

	// VolumeMount describes a mounting of a Volume
	// within a container.
	VolumeMount struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}

	// VolumeEmptyDir mounts a temporary directory from the
	// host node's filesystem into the container. This can
	// be used as a shared scratch space.
	VolumeEmptyDir struct {
		Medium    string `json:"medium"`
		SizeLimit int64  `json:"size_limit"`
	}

	// VolumeHostPath mounts a file or directory from the
	// host node's filesystem into your container.
	VolumeHostPath struct {
		Path string `json:"path"`
	}
)
