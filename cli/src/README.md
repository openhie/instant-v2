# CLI v2

An initialisation of the cli using cobra and viper with the official go sdk

## Prerequisites

- Minimum Go version of 1.18

## Commands
For a full breakdown of all functionality:

`cli --help`

## Build
`go build .`

## Generate autocomplete/suggestions
Generate the completion script with:

`cli completion bash|fish|powershell|zsh > /tmp/completion`

Load in the completion script with:

`source /tmp/completion`

## Additional Features over V1
- Refactored structure allows for easier contribution and maintainability
- Auto removes instant container and volume before performing a package action
- Support for multiple env files
- Support for predefined profiles to load in when performing a package action
- Generate a new project prompt
- Generate a new package prompt
- Autocomplete/Autosuggestions based off the config file

<!-- TODO: docs for tests -->
