# Architecture

The fundamental concept of Instant OpenHIE is that it can be extended to support additional use cases and workflows. This is achieved through packages. Packages conform to a specific set of rules and contain scripts and resources that allow particular applications to be spun up and configured on the platforms supported by Instant OpenHIE. The currently supported platforms are **local and remote deployments via Docker swarm**.

Packages alone, however, do not solve the whole problem. A consistent way of managing and executing these packages is necessary. Much of the work that goes into Instant OpenHIE is developing the software that enables the packages to be deployed and managed.&#x20;

## Package architecture[#](https://openhie.github.io/instant/docs/more-info/architecture#package-architecture) <a href="#package-architecture" id="package-architecture"></a>

Packages can be one of two different types. An **infrastructural package** or a **use case package**. Infrastructural packages setup and configure particular applications or sets of applications that are commonly grouped together. By themselves these packages just get the applications started and they aren't configured for a particular use case. On the other hand, use case packages depend on infrastructural packages and configure the applications set up by them and setup additional mediators that allow these applications to work together. They do this to enable a particular use case to be enacted.

You can think of use case packages as adding features to the end-user whereas infrastructural packages provide the dependencies to the use case packages. This separation allows an infrastructural package that, say implements a FHIR server to be replaced with different packages that implements a different FHIR server. As long as these packages can be configured in a standards-based way, a use case packages could work with either of these infrastructural packages. This gives users options for the applications that they wish to use.

<figure><img src="https://openhie.github.io/instant/img/instant-openhie-package-arch.png" alt=""><figcaption><p>Note: instant.json has been renamed to package-metadata.json in the latest versions of Instant</p></figcaption></figure>

Each package will contain the following sorts of technical artifacts that allow it to be spun up and down within the supported platforms:

* Bash scripts for setting up the applications required for this package’s use cases and workflows. These bash script will use docker compose files to launch container within [Docker Swarm](https://docs.docker.com/engine/swarm/).
* Configuration scripts to setup required configuration metadata
* Extensions to the test harness to test the added use cases with test data

The exact artefacts that go into a package are described in the section on [creating packages](../package/create-a-custom-package/).

## Execution architecture[#](https://openhie.github.io/instant/docs/more-info/architecture#execution-architecture) <a href="#execution-architecture" id="execution-architecture"></a>

Much of the tooling that Instant OpenHIE provides is to allow packages to be executed once they have been defined. The core principles of the execution architecture are that:

* Packages may be executed in a cross platform way (on Window, OSX or Linux)
* Complexity should be hidden from the user to make spinning up packages easy
* The software that executes the packages should be self contained so that it is easy to download and get running

The core concept of the execution architecture is that all the tooling to execute packages is contained within a single docker image - called the execution image. This enables it to run cross-platform as long as the required dependencies are installed (Docker Desktop for Windows/OSX or Docker engine and Docker compose for Linux). To enable the execution image to spin up packages using docker it is passed the docker socket file from the host so that any containers it created are in fact created on the host. The execution image gives us a consistent environment to develop for and allow us to run the scripts that execute the packages without caring about the underlying host OS, we just need to interface with the deployment platform.

An overview of the architecture is displayed below. Following that we describe how these components all works together in the sections that follow.

<figure><img src="https://openhie.github.io/instant/img/instant-openhie-arch.png" alt=""><figcaption></figcaption></figure>

### Go executable CLI app[#](https://openhie.github.io/instant/docs/more-info/architecture#go-executable-cli-app) <a href="#go-executable-cli-app" id="go-executable-cli-app"></a>

The CLI app will allow the setup and configuration of Instant OpenHIE to be easily managed by the user. It is the only executable that the user would need to download to interact with instant OpenHIE. It is written in Go so that it may compile to a platform-independent executable with 0 dependencies. It will know how to fetch the execution image from docker hub and be able to execute it at the command of the user.

The app will allow the user to execute it as an interactive CLI tool, however, it will also have the ability to launch a Web UI and API server that work together to provide the user with a Web interface for managing and executing Instant OpenHIE. Choices made in the Web UI will be sent to the API server and from there the server will issue commands by running the execution image with a particular set of parameters as shown in the diagram.

The CLI app will also be responsible for ensuring the execution image is executed with the correct configuration to ensure it works cross-platform and with the necessary user setting. For example, it will mount the docker socket file into the container so that docker containers are created on the host and it will mount various user config files for docker hub and/or AWS credentials. This enables the execution image to execute as if it were configured like the user's host system.

### The entry point script[#](https://openhie.github.io/instant/docs/more-info/architecture#the-entry-point-script) <a href="#the-entry-point-script" id="the-entry-point-script"></a>

The execution image, when executed itself, runs an execution script that is written in Typescript. The purpose of this script is to accept various parameters via the command line, discover the available packages and execute the requested action for a list of the available packages. The list of commandline options that it supports is shown as an example in the architecture diagram. This script is the heart of Instant OpenHIE as it is what actually executes and controls the packages. When the execution image is run this is the script that gets executed to perform all actions.

The execution image bundles all the dependencies that the execution script requires so that when the CLI app downloads the execution image from docker hub it is ready to go. The user will be required to have docker installed to run this image. The bundled dependencies in this image include: the docker and docker-compose clients and the cucumber executable for testing purposes.

It is also the responsibility of this script to determine dependencies between packages and ensure that the packages are spun up in dependency order, starting with the [Core package](https://openhie.github.io/instant/docs/packages/core).

### 3rd party packages[#](https://openhie.github.io/instant/docs/more-info/architecture#3rd-party-packages) <a href="#3rd-party-packages" id="3rd-party-packages"></a>

By default no packages are bundled with Instant OpenHIE. Instead it is left up to the community to provide packages for application or use cases. Package location are described in the config file for the CLI app.

Jembi has create an initial set of packages that allows us to ceate a foundational HIE using the applications that we like to use, you may use this as a base to get started. See [OpenHIM Platform](https://jembi.gitbook.io/openhim-platform/).

### Test harness[#](https://openhie.github.io/instant/docs/more-info/architecture#test-harness) <a href="#test-harness" id="test-harness"></a>

A test harness is also built into the execution image that can execute a suite of tests against the stood up infrastructure and ensure that they are functioning correctly. The test harness utilizes the Gerkin language to describe tests and the Cucumber tool to execute these. Packages are required to include a Gerkin feature file and the source code that is able to execute the features in a `features` folder. The entry point script can then execute these test scripts on demand for each package. This provides the user with a mechanism to test an instantiation of the architecture and provides a way to explore what each use case package can do.
