# What are packages?

The fundamental concept of Instant OpenHIE is that it can be extended to support additional use cases and workflows. This is achieved through packages. A package is a set of scripts and config files that are responsible for configuring a particular application, set of application or applications configured for a specific use-case.  A package may depend on from another existing package for more complicated functionality.

A package is intended to encompass a set of functionality rather than just setup generic applications. Packages are expected to configure the applications so that they may enact a particular functional role with the HIE. This may include setting up test data, necessary metadata and pre-configuring applications.

Packages can be one of two different types. An **infrastructural package** and a **use case package**. Infrastructural packages setup and configure particular applications or sets of applications that may be grouped together. By themselves, these packages only start the applications and they aren't configured for a particular use case. On the other hand, use case packages rely on infrastructural packages and configure the applications set up by them and setup additional mediators that allow applications to work together. They do this to enable a particular use case to be enacted. You can think of use case packages as adding features for the end-user whereas infrastructural packages provide the dependencies to the use case packages that enable the feature to work.

Each package will contain the following types of technical artefacts:

* Bash scripts for setting up the applications required for this package’s use cases and workflows. These bash script will use docker compose files to launch container within [Docker Swarm](https://docs.docker.com/engine/swarm/).
* Configuration scripts to setup required configuration metadata
* Extensions to the test harness to test the added use cases with test data

The below diagram shows how packages will extend off of each other to add use cases of increasing complexity.

<figure><img src="../.gitbook/assets/image.png" alt=""><figcaption></figcaption></figure>
