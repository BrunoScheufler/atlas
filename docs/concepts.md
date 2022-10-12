# Concepts

## basics

In each Atlasfile, you can define artifacts, services, and stacks. Since all Atlasfiles are read and merged when
you use the Atlas CLI, there are no requirements as to how many Atlasfiles you should create, so you can simply find out
what works best for you.

## artifacts

Artifacts generate OCI-compliant container images using `docker build`. You can pass all relevant options like context,
dockerfile, and build args. Artifacts can depend on other artifacts, which means that they will be built in the correct order.

## services

Services require an image or artifact to create a container from, and can be configured with environment variables,
environment files, ports, volumes, and commands.

## stacks

Stacks assemble multiple services in a specific order, and can be started, stopped, and restarted together.
