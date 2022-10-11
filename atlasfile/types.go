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

type ContainerRestarts string

const (
	ContainerRestartsAlways        = "always"
	ContainerRestartsOnFailure     = "on-failure"
	ContainerRestartsUnlessStopped = "unless-stopped"
	ContainerRestartsNo            = "no"
)

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

	Restart ContainerRestarts `json:"restart"`

	Interactive bool `json:"interactive"`
	TTY         bool `json:"tty"`
}

type StackService struct {
	Name        string `json:"name"`
	ServiceName string `json:"serviceName"`

	// Environment overwrites environment variables specified in ServiceConfig.Environment and ServiceConfig.EnvironmentFiles
	Environment map[string]string `json:"environment"`

	JoinStackNetworks []string     `json:"joinStackNetworks"`
	ExposePorts       []PortExpose `json:"exposePorts"`

	// LocalEnvironment specifies variables that overwrite Environment, ServiceConfig.Environment and ServiceConfig.EnvironmentFiles
	// when running atlas env (usually URLs that should be rewritten to localhost when running a service outside of Docker)
	LocalEnvironment map[string]string `json:"localEnvironment"`
}

type StackConfig struct {
	dirpath        string
	containerNames map[string]string

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
