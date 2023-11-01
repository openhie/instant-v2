---
description: Instantly deploy complex HIE components
cover: .gitbook/assets/iohie-logo-darkbg-transparent.svg
coverY: 0
layout:
  cover:
    visible: true
    size: hero
  title:
    visible: false
  description:
    visible: true
  tableOfContents:
    visible: false
  outline:
    visible: false
  pagination:
    visible: true
---

# Landing

Instant OpenHIE aims to allow Health Information Exchange components to be packaged up, deployed, operated and scaled via a simple CLI.

<div data-full-width="false">

<figure><img src=".gitbook/assets/feature.introduction.svg" alt="" width="375"><figcaption></figcaption></figure>

</div>

## Introduction

The Instant OpenHIE project aims to reduce the costs and skills required for software developers to deploy an OpenHIE architecture for quicker solution testing and as a starting point for faster production implementation and customisation.

{% content-ref url="getting-started/" %}
[getting-started](getting-started/)
{% endcontent-ref %}

<figure><img src=".gitbook/assets/feature.concepts (1).svg" alt="" width="375"><figcaption></figcaption></figure>

## Concepts

Instant OpenHIE provides an easy way to setup, explore and develop with the OpenHIE Architecture. It is a deployment tool that allows packages to be added that support multiple different use cases and workflows specified by OpenHIE. Each package contains scripts to stand up and configure applications that support these various workflows.

{% content-ref url="concepts/" %}
[concepts](concepts/)
{% endcontent-ref %}

<figure><img src=".gitbook/assets/feature.packages.svg" alt="" width="375"><figcaption></figcaption></figure>

## Packages

The fundamental concept of Instant OpenHIE is that it can be extended to support additional use cases and workflows. This is achieved through packages. No packages are included by default in Instant OpenHIE, packages are provided and maintained by the community.

{% hint style="info" %}
Jembi has developed a set of packages called the [OpenHIM Platform](https://jembi.gitbook.io/openhim-platform/) which allow you to get a foundational health information exchange setup in an instant.
{% endhint %}

{% content-ref url="package/" %}
[package](package/)
{% endcontent-ref %}

## Differences from Instant OpenHIE v1

The key changes from original Instant OpenHIE are:

* A rewrite of the original CLI - the commands and parameters have changed
* Docker swarm is now the only supported target - this allow us to scale services across servers
* A set of bash init function are included to help implementers create package deployment scripts
* The entry point bash script for packages is now named `swarm.sh`
