# Adding your package

Packages can be added in two different ways. These are described below.

## Using a custom package config

You may define custom packages either in the [config file](../getting-started/config.md) or via a [command line flag](../getting-started/cli.md). This configuration can either be a **local path** to the package or a **github url**.

## In a custom docker image

Packages can be built into a custom docker image that you may version and push to Docker Hub as you wish. This is the image referenced in the `image` property of the config file. This image MUST be built by extending the `openhie/package-base` package which contains some key dependencies. See an example of how Jembi does this [for Jembi Platform here](https://github.com/jembi/platform/blob/main/Dockerfile#L1-L3). You simple need to add your package folder into the image.
