# CLI

## Usage

```
cli <deploy command> [custom flags] <package ids>
```

{% hint style="info" %}
Use \`cli --help\` for more information on how to use the cli
{% endhint %}

{% hint style="warning" %}
CLI usage requires a [config.md](config.md "mention") file in the directory where the cli is being called
{% endhint %}

## Commands

```
init        Initializing a service
up          Start up a service that has been shut down or update a service
down        Destroy a service
destroy     Bring down a running service
help        Help info
```

## Flags

<pre><code>--dev                    Specifies the development mode in which all service ports are exposed
--only, -o               Used to specify a single service for services that have dependencies. 
                         For cases where one wants to shut down or destroy a service without affecting its dependencies
<strong>-e                       For specifying an environment variable
</strong>--env-file               For specifying the path to an environment variables file
--custom-package, -c     Specifies path or url to a custom package. Git ssh urls are supported
--image-version          The version of the project used for the deploy. Defaults to 'latest'
-t                       Specifies the target to deploy to. Options are docker, swarm (docker swarm) and k8s (kubernetes) - project dependant
-*, --*                  Unrecognised flags are passed through uninterpreted</code></pre>

## Examples

```bash
{your_binary_file} init -t=swarm --dev -e="NODE_ENV=prod" --env-file="../env.dev" -c="../customPackage1" -c="<git@github.com/customPackage2>"  interoperability-layer-openhim customPackage1_id customPackage2_id
```

```bash
{your_binary_file} down -t=docker --only elastic_analytics
```

