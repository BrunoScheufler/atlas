# Context

This page outlines the design decisions that were made for Atlas. Read this to understand why Atlas works the way it does and why it may or may not support your use case.

## Improving the local development experience

Atlas has been designed for local development exclusively. Other than Docker Compose which can be used for remote deployments or production environments, Atlas does not have the ambition to serve multiple environments and workflows but strongly focuses on local development. This decision was made because a lack of focus has been a problem with other tools in the past, and we want to avoid that.

## Basic principles

As you will read in the following, Atlas is a lightweight wrapper on top of Docker, adding powerful local development capabilities. Whenever possible and in line with the mission to improve the local development experience, Atlas tries to leverage existing features and tools over providing its own.

## Configuration as code

As you might have seen, [Atlasfiles](./atlasfile.md) can be written in multiple languages you may already be using, including Go and TypeScript.

If you've worked with Docker Compose or Kubernetes before, you will most likely have encountered frustrating sessions trying to make YAML work or do something it wasn't designed for (templating, dynamic values, etc.).

For Atlas, I chose to support interactive configuration files that would be evaluated just in time. Writing your configuration as code allows to pull in external systems, template values, run static code analysis and linting, and generally do anything you can do with code.

You can also include your favorite libraries or share code with other services, as you're working in the same environment, with the same stack.

And the best thing, nobody on your team has to learn a new stack!

Inherent downsides of interactive configuration include that evaluation is more expensive than just parsing a text file, but Anzu solves this by intelligently caching the result of the evaluation.

## Monorepo first

Most teams have decided to colocate multiple services in a single repository. Atlas is designed around this setup, so your experience may vary if you're trying to use it in other environments.

## Service and build configuration close to code

When working with a monorepo, you will usually create a root Atlasfile as well as individual Atlasfiles in your service directories. This makes it much easier to update all values when you're refactoring something or reviewing code, growing with your team.

This is in strong contrast to Docker Compose, where you have a single centralized Compose file which dictates both build and runtime behaviour. With Atlas, all build instructions are configured close to your service and only overrides are configured centrally when you arrange your services into stacks.

## Docker for isolated, repeatable workloads

Since most services end up packaged and delivered using OCI-compliant container images (which Docker emits), Atlas builds on that foundation. This means you can and should use the same service images across development and production.

Due to using Docker as a direct dependency, we can also use the layer caching it provides to speed up builds significantly. When your service (or one of its dependencies) has been built before, it will only need to be rebuilt when something changed (as defined in your Dockerfiles). You don't have to specify any other caching configuration yourself, it just works (and Atlas doesn't contain a line of code it doesn't need).

