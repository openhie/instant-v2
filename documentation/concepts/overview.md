# Overview

Instant OpenHIE provides a framework for deploying and configuring OpenHIE components for particular OpenHIE use cases and workflows. These scripts and configurations are organised into self-contained packages. Each of these packages may depend on other packages which allows highly complex infrastructure to be setup instantly by deploying a number of packages. Instant OpenHIE itself doesn't provide any packages, it is just the framework and specification for deploying packages. Packages are implemented and maintained by the community.

Each of these packages contain scripts which setup containerised applications. The scripts configure and pre-load necessary data into the containers. Docker will be used to containerise the applications which allows them to be easily deployed.

Instant OpenHIE currently supports deploying of packages to [Docker Swarm](https://docs.docker.com/engine/swarm/) to allow package to be easily setup both locally or on production server in a scalable and high available way.
