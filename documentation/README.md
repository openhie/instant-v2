---
description: Instantly deploy complex HIE components
layout: landing
---

# Overview

This is the second iteration of the Instant OpenHIE project. Instant OpenHIE aims to allow Health Information Exchange components to be packages up, deployed, operated and scaled via a simple CLI.

The key changes from original Instant OpenHIE are:

* A rewrite of the original CLI - the commands and parameters have changed
* Docker swarm is now the only supported target - this allow us to scale services across servers
* The entry point bash script for packages is now named `swarm.sh`

This repository houses the 2 base components:

`package-base`: A docker image with the base configuration for standing up containers and running tests on which implementations can base their images. Instant OpenHIE v2 contains no packages by default, implementers can either extend the image to add their own packages or include them at runtime with CLI flags.

`cli`: A Go CLI app for coordinating the deployment and configuration of packages.
