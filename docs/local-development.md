# Local Development

Atlas has been built for running code in containers, but sometimes you need to run code on your host machine. Debugging your code in your IDE of choice should not be an edge case, so Atlas supports starting and stopping individual services, and exporting service environment variables to an `.env.local` file in your service directory, so you have everything you need to start a service. No more copying environment variables around.

## Starting and stopping services

When you want to develop a service outside of Docker, simply stop the service container using `atlas stop`. Once you're done, you can start it again.

```bash
# Stop container
atlas stop -s my-stack my-service

# Work outside of Docker

# Start it again
atlas start -s my-stack my-service
```

## Printing environment variables

To start up your services, you'll often need to set a couple of environment variables like database URLs, logging settings, and more. Since you've already defined those variables in your Atlasfile, you can simply export them close to your service using an `.env.local` file. This file is read by your `.env` library of choice (which should even handle the [different file name](https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use) automatically as an override of other existing env files) and injected into the service at runtime.

```bash
# always run command in  service directory
cd services/my-service

# export env variables of my-service in my-stack
atlas env -s my-stack my-service

# Atlas will create (or overwrite!) an .env.local file
cat .env.local
# ...
```

## Overwriting environment variables to run outside of Docker

When dealing with environment variables like URLs for services and databases running in Docker, simply copying them over will not suffice as you cannot reach the same host you use with Docker's DNS. For this reason, stack services configured in your root [Atlasfiles](./atlasfile.md) include a `LocalEnvironment` map where you can pass variables that overwrite any other variables defined on the stack or service level.

