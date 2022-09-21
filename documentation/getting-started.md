# Getting Started

## Download

### Binary

The binary may be downloaded from the [releases page of the github repo](https://github.com/openhie/package-starter-kit/releases)

Or

The binary may be download via the terminal with the following url based on your operating system

{% tabs %}
{% tab title="Linux" %}
Download the binary

```bash
curl -L https://github.com/openhie/package-starter-kit/releases/download/0.6.0/gocli-linux -o platform-linux
```

Grant the binary executable permissions

```bash
chmod +x ./platform-linux
```
{% endtab %}

{% tab title="MacOS" %}
Download the binary

```bash
curl -L https://github.com/openhie/package-starter-kit/releases/download/0.6.0/gocli-macos -o platform-macos
```

Grant the binary executable permissions

```bash
chmod +x ./platform-macos
```
{% endtab %}

{% tab title="Windows" %}
Download the binary

```bash
curl -L https://github.com/openhie/package-starter-kit/releases/download/0.6.0/gocli.exe -o platform.exe
```
{% endtab %}
{% endtabs %}

Or

The binary may be downloaded with go

Golang would need to be installed on your machine to be able to download the package with go. Follow the steps [here](https://go.dev/doc/install) to get `go` set up on your machine

Thereafter run&#x20;

```bash
go install github.com/openhie/package-starter-kit/cli@latest
```

Ensure the following lines are in your \~/.bashrc or \~/.zshrc file

```
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

The run `source ~/.bashrc` or `source ~/.zshrc` to load the changes

Now the CLI should be available in your terminal and accessable by running `cli --help`



## Set up

Create a new project directory `mkdir my-project`

