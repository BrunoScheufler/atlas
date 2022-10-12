# Atlasfile Reference

Atlasfiles are interactive configuration files, written as code. They're stored inside `.atlas` directories next to your service code or at the root of your repository.

## .atlas directory

Atlasfiles are stored within dedicated `.atlas` directories to create a layer of isolation from your service code (so you can store related files like package and dependency configurations) as well as to make it easy to find all Atlasfiles in your repository.

## The root Atlasfile

Whenever you run the Atlas CLI, it searches for an Atlasfile at the root of your repository, usually denoted by saving a file called `Atlasfile.root.go`. When searching, Atlas will jump up to a maximum of 5 levels from your current working directory.

Once found, Atlas searches for all Atlasfiles defined at and below root level, and merges them into one file. This means that you can define services, artifacts, and stacks at any level, but we recommend only storing services and artifacts close to your services and storing stacks in your root Atlasfile for clarity.

## Language support

Atlasfiles are simple binaries which launch a gRPC server to communicate with the CLI. For this reason, theoretically, all languages that support gRPC servers, are supported. Right now, Atlas has been tested with Go, but more languages and documentation will be added in the future.
