---
description: >-
  The CLI is supported for Linux, MacOS, and partially for Windows, although
  functionality is not guaranteed for Windows.
---

# CLI

## Usage

{% hint style="warning" %}
The CLI must be named 'instant-linux' for Linux, and 'instant-macos' for MacOS for autocomplete to work
{% endhint %}

## Main Commands

```
completion    Generate the autocompletion script for the specified shell
package       Package level commands
project       Project level commands
help          Help about any command
```

## Sub Commands

### package

The package sub command includes commands:

```
init          Initialize a package with relevant configs, volumes and setup
up            Stand a package back up after it has been brought down
down          Bring a package down without removing volumes or configs
remove        Remove everything related to a package (volumes, configs, etc)
generate      Generate a new package
```

The package level commands, as shown, are there to control packages within a project, as well as generate the skeleton for a new package.

Each of these sub commands accept the following flags:

```
Flags:
      --config string         config file (default is $WORKING_DIR/config.yaml)
  -c, --custom-path strings   Path(s) to custom package(s)
  -d, --dev dev               For development related functionality (Passes dev as the second argument to your swarm file)
      --env-file strings      env file
  -e, --env-var strings       Env var(s) to set or overwrite
  -h, --help                  help for down
  -n, --name strings          The name(s) of the package(s)
  -o, --only                  Ignore package dependencies
  -p, --profile string        The profile name to load parameters from (defined in config.yml)
```

E.g. `./instant package init -n interoperability-layer-openhim`

For information about flags associated to any one of the package commands, do `instant-linux package [command] --help`

{% hint style="info" %}
After generating a new package, remember to add the package ID to the config file
{% endhint %}

{% hint style="warning" %}
* Packages in a project can only be started if included in the config file
* Command line arguments like `--dev` and `--only` will overwrite those specified in the config file profiles when using that particular profile
* Env vars in `--profile` env var files are appended to by env var files specified in the command line, or overwritten by the command line env var files if there are conflicting env vars
* Custom packages in a profile must be specified in the customPackages section of the config file
{% endhint %}

### project

The project sub command includes commands:

```
init          Initialize all packages in a project
up            Up all packages in the project
down          Down all packages in the project
destroy       Destroy all packages in the project
generate      Generate a new project
```

The project level commands, as shown, are there to simultaneously perform commands on all packages in a project, as well as generate the config file for a new project, in the desired format.

Each of these sub commands accept the following flags:

```
Flags:
      --config string         config file (default is $WORKING_DIR/config.yaml)
  -c, --custom-path strings   Path(s) to custom package(s)
  -d, --dev dev               For development related functionality (Passes dev as the second argument to your swarm file)
      --env-file strings      env file
  -e, --env-var strings       Env var(s) to set or overwrite
  -h, --help                  help for destroy
  -o, --only                  Ignore package dependencies
```

For information about flags associated to any one of the project commands, do `instant-linux project [command] --help`

### completion

The completion sub command includes commands:

```
bash          Generate the autocompletion script for bash
zsh           Generate the autocompletion script for zsh
```

The project level commands, as shown, are there to enable autocomplete for the specified shell.

{% hint style="warning" %}
Remember to reload your shell after generating the autocomplete script
{% endhint %}
