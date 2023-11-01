# Quick start

To get started with Instant OpenHIE you will first need to download the CLI tool.

The latest binary may be downloaded from the [releases page of the github repo](https://github.com/openhie/package-starter-kit/releases)

Or

The binary may be download via the terminal with the following url based on your operating system

{% tabs %}
{% tab title="Linux" %}
Download the binary

```bash
curl -L https://github.com/openhie/instant-v2/releases/latest/download/instant-linux -o instant
```

Grant the binary executable permissions

```bash
chmod +x ./instant
```
{% endtab %}

{% tab title="MacOS" %}
Download the binary

```bash
curl -L https://github.com/openhie/instant-v2/releases/latest/download/instant-macos -o instant
```

Grant the binary executable permissions

```bash
chmod +x ./instant
```
{% endtab %}

{% tab title="Windows" %}
For Windows it is recommend to install the [WSL](https://learn.microsoft.com/en-us/windows/wsl/install) tools and continue to follow the linux instructurions.
{% endtab %}
{% endtabs %}

To test that the binary works, run the executable with no commands to see the help text.

```
$ ./instant
A cli to assist with package deployment and management

Usage:
  cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  package     Package level commands
  project     Project level commands
  version     Print the CLI version

Flags:
  -h, --help   help for cli

Use "cli [command] --help" for more information about a command.
```

Next, you would want to configure which packages Instant can deploy for your particular needs.

Instant doesn't ship with any default packages to deploy. Packages are expected to be created by the community and in time there will be many option available. Jembi has curated a set of packages that we commonly use to help implementer to get started with a foundational set of health information exchange components. To get started with that pre-configured package set see the [Jembi Platform docs](https://jembi.gitbook.io/jembi-platform/).

Otherwise, you may create your own config for your [own set of packages](../package/create-a-custom-package/). Continue to config section to find out how.
