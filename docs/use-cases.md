# Use Cases

Atlas works best when your setup meets one or more of the following characteristics:

## Multiple services

With its [interactive configuration](./atlasfile.md), Atlas makes it easy to compose multiple services and reuse the same images with different entrypoints, dedicated images on a similar base image, or just reusing details like environment variables.

The more services you have, the more useful Atlas becomes. This doesn't mean Atlas only works for microservices, even setups with two or three services already benefit from straightforward configuration and sharing code through Artifacts.

## Multiple regions

When you need to run multiple regions or deployment targets locally, Atlas shines. Stacks are isolated collections of services, so you can run multiple regions or environments locally without any conflicts.
