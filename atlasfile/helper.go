package atlasfile

import (
	"fmt"
	"path/filepath"
)

func (a *Atlasfile) GetService(name string) *ServiceConfig {
	for i := range a.Services {
		if a.Services[i].Name == name {
			return &a.Services[i]
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

// GetDirpath returns path of .atlas directory artifact was declared in
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

func (s *StackConfig) GetService(serviceName string) *StackService {
	for i, service := range s.Services {
		if service.Name == serviceName {
			return &s.Services[i]
		}
	}
	return nil
}

func (s *StackConfig) SetContainerName(service, containerName string) {
	if s.containerNames == nil {
		s.containerNames = make(map[string]string)
	}

	s.containerNames[service] = containerName
}

func (s *StackConfig) GetContainerName(service string) string {
	return s.containerNames[service]
}

// GetDirpath returns path of .atlas directory service was declared in
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

func (c *VolumeConfig) GetVolumeNameOrHostPath(cwd string, physicalVolName string) string {
	if c.IsVolume {
		return physicalVolName
	}

	if filepath.IsAbs(c.HostPathOrVolumeName) {
		return c.HostPathOrVolumeName
	}

	return filepath.Join(cwd, c.HostPathOrVolumeName)
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
