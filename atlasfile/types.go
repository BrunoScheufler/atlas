package atlasfile

type BuildOptions struct {
	Dockerfile string            `json:"dockerfile"`
	Context    string            `json:"context"`
	BuildArgs  map[string]string `json:"build_args"`
	Target     string            `json:"target"`

	ImageName string `json:"imageName"`
	TagName   string `json:"tagName"`
}

type ArtifactRef struct {
	Name     string          `json:"name"`
	Artifact *ArtifactConfig `json:"artifact"`
}

type VolumeConfig struct {
	volName              string
	IsVolume             bool   `json:"isVolume"`
	HostPathOrVolumeName string `json:"hostPath"`
	ContainerPath        string `json:"containerPath"`
}

type PortRequest struct {
	ContainerPort int    `json:"containerPort"`
	Protocol      string `json:"protocol"`
}

type PortExpose struct {
	HostPort      int
	ContainerPort int
}

type ServiceConfig struct {
	dirpath string

	Name string `json:"name"`

	Artifact *ArtifactRef `json:"artifact"`
	Image    string       `json:"image"`

	Entrypoint []string `json:"entrypoint"`
	Command    []string `json:"command"`

	Ports []PortRequest `json:"port_requests"`

	Environment      map[string]string `json:"environment"`
	EnvironmentFiles []string          `json:"environment_files"`

	Volumes []VolumeConfig `json:"volumes"`
}

type StackService struct {
	Name              string            `json:"name"`
	ServiceName       string            `json:"serviceName"`
	Environment       map[string]string `json:"environment"`
	JoinStackNetworks []string          `json:"joinStackNetworks"`
	ExposePorts       []PortExpose      `json:"exposePorts"`
}

type StackConfig struct {
	dirpath        string
	containerNames map[string]string

	networkName string

	Name     string         `json:"name"`
	Services []StackService `json:"services"`
}

type ArtifactDependsOn struct {
	Services  []string `json:"services"`
	Artifacts []string `json:"artifacts"`
}

type ArtifactConfig struct {
	dirpath string
	Name    string `json:"name"`

	Build     BuildOptions      `json:"build"`
	DependsOn ArtifactDependsOn `json:"depends_on"`
}

type Atlasfile struct {
	dirpath   string
	Artifacts []ArtifactConfig `json:"artifacts"`
	Services  []ServiceConfig  `json:"services"`
	Stacks    []StackConfig    `json:"stacks"`
}
