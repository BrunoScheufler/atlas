# atlas: better local development

## background

Local development involves running as many services as possible locally, mostly using containers for reproducible builds and isolated environments. While it's easy to configure basic setups using Docker Compose, a higher number of services or more complex dependency relationships can quickly exceed its capabilities.

Atlas is built to handle complex local development environments, based largely on Docker containers. Each service is defined in an Atlasfile, which is usually written as code to make configuration straightforward (reusing environment variables or fetching data from external systems like IaC is just a few lines of code away), and composed in a stack. Atlas simply detects the dependency graph, builds all required artifacts (Docker images for services), and starts the necessary containers.

## features

- **Artifact Graph**: All required artifacts are collected and built in the most efficient order, leveraging layer caching and parallel builds.
- **Atlasfiles**: Atlasfiles can be written in Go, Node.js, TOML, and potentially any other language supporting gRPC.
- **Services**: Services are defined close to the relevant code
- **Stacks**: Stacks can define multiple services and overwrite configuration where needed

## installation

Download the Atlas CLI from the [releases](https://github.com/brunoscheufler/atlas/releases) page.

## getting started

### creating a root Atlasfile

In the root directory of your repository, you can create a `Atlasfile.root.go` file. This is structurally the same as any other Atlasfile, but tells Atlas to stop looking for further Atlasfiles in parent directories. This is relevant, because you can use the Atlas CLI in any layer of your project directory structure, and it will include all detected Atlasfiles from the root file downwards.

```go
package main

import (
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/sdk/atlas-sdk-go"
	"os"
)

func main() {
	err := sdk.Start(&atlasfile.Atlasfile{
		Services: []atlasfile.ServiceConfig{
			{
				Name:  "global-db",
				Image: "postgres:14",
				Ports: []atlasfile.PortRequest{{5432, "tcp"}},
				Environment: map[string]string{
					"POSTGRES_USER":     "directory",
					"POSTGRES_PASSWORD": "directory",
					"POSTGRES_DB":       "directory",
					"PGDATA":            "/var/lib/postgresql/test-data",
				},
				Volumes: []atlasfile.VolumeConfig{
					{IsVolume: true, HostPathOrVolumeName: "postgres", ContainerPath: "/var/lib/postgresql/test-data"},
				},
			},
			{
				Name:  "localstack",
				Image: "localstack/localstack:0.13.3",
				Ports: []atlasfile.PortRequest{{4566, "tcp"}},
				Environment: map[string]string{
					"SERVICES":          "s3,sqs,kinesis,firehose",
					"DEFAULT_REGION":    "eu-central-1",
					"EDGE_PORT":         "4566",
					"HOSTNAME_EXTERNAL": "localstack",
				},
				Volumes: []atlasfile.VolumeConfig{
					{HostPathOrVolumeName: "./localstack", ContainerPath: "/docker-entrypoint-initaws.d"},
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("could not start atlasfile: %s", err.Error())
		os.Exit(1)
	}
}
```

### adding a service

In the directory of your service, you can create another Atlasfile, outlining how your service should be built, which ports it offers, and which environment variables (or env files) it receives by default.

```go
package main

import (
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/sdk/atlas-sdk-go"
	"os"
)

func main() {
	err := sdk.Start(&atlasfile.Atlasfile{
		Services: []atlasfile.ServiceConfig{
			{
				Name: "api",
				Artifact: &atlasfile.ArtifactRef{
					Artifact: &atlasfile.ArtifactConfig{
						Name: "api",
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("could not start atlasfile: %s", err.Error())
		os.Exit(1)
	}
}
```

As you can see, defining your service can range from simple cases to more complex ones. You don't necessarily need to build an image, you could also just pass an `Image` right away. If you want to reuse the same base image for multiple service entrypoints, just declare the artifact once and refer to it by name:

```go
package main

import (
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/sdk/atlas-sdk-go"
	"os"
)

func main() {
	err := sdk.Start(&atlasfile.Atlasfile{
		Artifacts: []atlasfile.ArtifactConfig{
			{
				Name: "base",
			},
		},
		Services: []atlasfile.ServiceConfig{
			{
				Name: "api",
				Artifact: &atlasfile.ArtifactRef{
					Name: "base",
				},
				Command:          []string{"--server"},
				EnvironmentFiles: []string{".env"},
			},
			{
				Name: "worker",
				Artifact: &atlasfile.ArtifactRef{
					Name: "base",
				},
				Command:          []string{"--worker"},
				EnvironmentFiles: []string{".env"},
			},
		},
	})
	if err != nil {
		fmt.Printf("could not start atlasfile: %s", err.Error())
		os.Exit(1)
	}
}
```

### defining a stack

Last but not least, Atlas allows you to define groups of services that should be provisioned together, as part of a stack.

```go
package main

import (
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/brunoscheufler/atlas/sdk/atlas-sdk-go"
	"os"
)

func main() {
	err := sdk.Start(&atlasfile.Atlasfile{
		Stacks: []atlasfile.StackConfig{
			{
				Name: "regional",
				Services: []atlasfile.StackService{
					{
						Name:        "api",
					},
					{
						Name:        "worker",
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("could not start atlasfile: %s", err.Error())
		os.Exit(1)
	}
}
```

### launching the stack

Simply run

```bash
atlas up
```

Atlas will build all artifacts in order, after which all configured service containers are launched in dedicated Docker networks.

## reference

## Atlasfile

### basics

In each Atlasfile, you can define artifacts, services, and stacks. Since all Atlasfiles are read and merged at whenever you use the Atlas CLI, there are no requirements as to how many Atlasfiles you should create, so you can simply find out what works best for you.

### artifacts

Artifacts generate OCI-compliant container images using `docker build`. You can pass all relevant options like context, dockerfile, and build args.

### services

Services require an image or artifact to create a container from, and can be configured with environment variables, environment files, ports, volumes, and commands.

### stacks

Stacks assemble multiple services in a specific order, and can be started, stopped, and restarted together.
