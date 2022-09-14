package atlasfile

import (
	"fmt"
	"path/filepath"
)

func (a *Atlasfile) GetService(name string) *ServiceConfig {
	for _, service := range a.Services {
		if service.Name == name {
			return &service
		}
	}
	return nil
}

func (a *Atlasfile) GetArtifact(name string) *ArtifactConfig {
	for _, artifact := range a.Artifacts {
		if artifact.Name == name {
			return &artifact
		}
	}
	return nil
}

func (ac *ArtifactConfig) GetDirpath() string {
	return ac.dirpath
}

func (a *Atlasfile) GetStack(name string) *StackConfig {
	for _, stack := range a.Stacks {
		if stack.Name == name {
			return &stack
		}
	}
	return nil
}

func (s *StackConfig) GetNetworkName() string {
	return s.networkName
}

func (s *StackConfig) SetNetworkName(networkName string) {
	s.networkName = networkName
}

func (s *StackConfig) SetContainerName(service, containerName string) {
	if s.containerNames == nil {
		s.containerNames = make(map[string]string)
	}

	s.containerNames[service] = containerName
}

func (s *ServiceConfig) GetDirpath() string {
	return s.dirpath
}

func BuildImageName(artifact *ArtifactConfig) string {
	imageName := artifact.Build.ImageName
	if imageName == "" {
		imageName = artifact.Name
	}

	tagName := artifact.Build.TagName
	if tagName == "" {
		tagName = "latest"
	}

	return fmt.Sprintf("%s:%s", imageName, tagName)
}

func (c *VolumeConfig) GetVolumeNameOrHostPath(cwd string) string {
	if c.IsVolume {
		return c.volName
	}

	if filepath.IsAbs(c.HostPathOrVolumeName) {
		return c.HostPathOrVolumeName
	}

	return filepath.Join(cwd, c.HostPathOrVolumeName)
}

func (c *VolumeConfig) SetVolName(name string) string {
	c.volName = name
	return name
}

func GetServicePort(requests []PortRequest, port int) *PortRequest {
	for i, request := range requests {
		if request.ContainerPort == port {
			return &requests[i]
		}
	}

	return nil
}

func (a *Atlasfile) String() string {
	var str string

	if len(a.Artifacts) > 0 {
		str += "--- Artifacts ---\n"
		for i, artifact := range a.Artifacts {
			str += fmt.Sprintf("%d. %s (%s)\n", i+1, artifact.Name, artifact.dirpath)
		}
	} else {
		str += "No artifacts found\n"
	}

	if len(a.Services) > 0 {
		str += "--- Services ---\n"
		for i, service := range a.Services {
			str += fmt.Sprintf("%d. %s (%s)\n", i+1, service.Name, service.dirpath)
		}
	} else {
		str += "No services found\n"
	}

	if len(a.Stacks) > 0 {
		str += "--- Stacks ---\n"
		for i, stack := range a.Stacks {
			str += fmt.Sprintf("%d. %s (%s)\n", i+1, stack.Name, stack.dirpath)
		}
	} else {
		str += "No stacks found\n"
	}

	return str
}

func (a *Atlasfile) GetServiceImage(service *ServiceConfig) (string, error) {
	if service.Image != "" {
		return service.Image, nil
	}

	var artifact *ArtifactConfig
	if service.Artifact.Name != "" {
		artifact = a.GetArtifact(service.Artifact.Name)
	} else {
		artifact = service.Artifact.Artifact
	}

	if artifact == nil {
		return "", fmt.Errorf("could not find artifact %s", service.Artifact.Name)
	}

	return BuildImageName(artifact), nil
}
