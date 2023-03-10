# Insant OpenHIE - v2

This is the second iteration of the Instant OpenHIE project. Instant OpenHIE aims to allow Health Information Exchange components to be packages up, deployed, operated and scaled via a simple CLI.

The key changes from original Instant OpenHIE are:

* A rewrite of the original CLI - the commands and parameters have changed
* Docker swarm is now the only supported target - this allow us to scale services across servers
* The entry point bash script for packages is now named `swarm.sh`

This repository houses the 2 base components:

`package-base`: A docker image with the base configuration for standing up containers and running tests on which implementations can base their images. Instant OpenHIE v2 contains no packages by default, implementers can either extend the image to add their own packages or include them at runtime with CLI flags.

`cli`: A Go CLI app for coordinating the deployment and configuration of packages.

For more info see the [documentation](https://jembi.gitbook.io/instant-v2/).

## Developers

### package-base releases

Docker image will be built on tag using github actions and pushed to dockerhub repo `openhie/package-base`
Tag format: '0.0.1'

### cli

This is a Go CLI app and is provided as a native binary for the AMD64 architecture on Windows, macOS, and Linux.

#### Dev prerequisites

* Install go, [see here](https://golang.org/doc/install). For Ubuntu you might want to use the go snap package, [see here](https://snapcraft.io/install/go/ubuntu).
* Add go binaries to you system \$PATH, on ubuntu: Add `export PATH=$PATH:$HOME/go/bin` to the end of your ~/.bashrc file. To use this change immediately source it: `source ~/.bashrc`
* Install dependencies, run this from the cli/src folder: `go mod tidy`

#### Running

For development, run the app using `go run .` from the `cli/src` folder.

#### Testing

To run the unit tests, ensure you're in the cli directory, then do:

```bash
./unit-test.sh
```

The [godog library](https://github.com/cucumber/godog), which provides us with the [Cucumber Framework](https://cucumber.io/) is used for functional tests of the Go CLI. To run the tests, you'll need to [install godog](https://github.com/cucumber/godog#install).

Then, navigate to the `cli/src` root folder and run the command below.

```bash
go install github.com/cucumber/godog/cmd/godog@v0.12.5 && go mod tidy
godog run
```

> These functional tests can only be run when the `instant` docker volume does not exist on the machine being used for running the tests. These tests make use of the mock image of the instant project `jembi/go-cli-test-image`.

#### Building

```sh
bash ./buildreleases.sh
```

#### Deploying

If any changes have been made to the Go CLI, update the version in `cli/src/cmd/version/version`

To build releases, create an instant tag and a release, the GitHub actions will build the code after creation of the release and add the binaries to the assets of the release if the GO CLI build succeeds.

Check [here](https://github.com/openhie/instant/actions/new) to review the output of the build and status of the binary deploy.

#### Schema

Should you wish to view documentation of the available configurations files, include the provided schemas in your project.

To include the provided schemas in your project, add the following to your `.vscode/settings.json` file:

```json
"json.schemas": [
    {
      "fileMatch": ["package-metadata.json"],
      "url": "https://raw.githubusercontent.com/openhie/package-starter-kit/main/schema/package-metadata.schema.json"
    }
],
"yaml.schemas": {
    "https://raw.githubusercontent.com/openhie/package-starter-kit/main/schema/config.schema.json": "config.yml"
},
```

Alternately, adhere the below example config.yml file showcases the available fields:

```yml
image: jembi/platform
logPath: /tmp/logs

packages:
  - analytics-datastore-elastic-search
  - dashboard-visualiser-kibana
  - data-mapper-logstash

customPackages:
  - id: disi-on-platform
    path: "git@github.com:jembi/disi-on-platform.git"

profiles:
  - name: dev
    packages:
      - analytics-datastore-elastic-search
      - dashboard-visualiser-kibana
    envFiles:
      - .env.dev
    dev: true
    only: true
```
