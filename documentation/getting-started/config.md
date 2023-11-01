# Config

Instant config describes which packages to make available in the CLI and a set of profiles which allow you to easily operate on a group of packages and apply env var files to configure then.

To generate a config file use the [`./instant project generate` command](cli.md#project).

A single file may be used for configuration

{% hint style="warning" %}
The `config.yml` file should be located at the root of the project, or pointed to using the `--config-file` command
{% endhint %}

A reference config file looks like this:

<pre class="language-yaml"><code class="lang-yaml"><strong>projectName: platform
</strong>image: jembi/platform
logPath: /tmp/logs

packages:
    - &#x3C;&#x3C;package-id>>

customPackages:
    - id: &#x3C;&#x3C;custom-package-id>>
      path: &#x3C;&#x3C;custom-package-path>>

profiles:
    - name: &#x3C;&#x3C;profile-name>>
      packages:
        - &#x3C;&#x3C;profile-package-id-1>>
        - &#x3C;&#x3C;profile-package-id-2>>
      envFiles:
        - &#x3C;&#x3C;env-file-1>>
        - &#x3C;&#x3C;env-file-2>>
      dev: false
      only: false
</code></pre>

* projectName - gives this project configuration a name
* image - defines the Docker image to use during deployment. This image should container the packages you wish to launch if you are not using customPackages. A default image with no packages included can be found at `openhie/package-base:latest`
* logPath - gives a location to put log of the CLI output, useful for debugging
* packages - lists the package ids that you expect to exist in the image
* customPackages - lists packages that are not in the image. The path can either point to a file system location or a github url.
* profiles - lists a number of profiles that are defind for this project. A profile is a group of packages, config and env var files that can be operated on (i.e. launched) together.&#x20;
  * name - a profile name
  * packages - list of package ids that form part of this profile
  * envFiles - a lsit of env var file to apply when operating on these packages to configure then to do what you want
  * dev - launches the profile is dev mode, which is an instruction to packages to start in dev mode which usually mean to expose more ports than they usually would for development and debugging reasons
  * only - instructs the profile to operate on only the packages listed and any package dependencies will be ignored.

{% hint style="info" %}
* Packages listed in a profile must be specified in either the customPackages or packages section
{% endhint %}

## Launching individual packages

Once a config has been defined with 1 or more packages, you may launch or stop packages by using the [`./instant package <init|up|down|remove> -n <package_id>` command](cli.md#package).

## Launching a profile

If you want to launch all the packages listed in a profile with the specified configuration and env vars applied, use the [`./instant package <init|up|down|remove> -p <profile_name>` command](cli.md#package).

## Launching entire projects

Instead of individual packages, you can also launch everything define in your project config by using the [`./instant project <init|up|down|remove>` command](cli.md#project).

## Changing the default banner displayed in the CLI

The CLI banner refers to the uppermost display when running the CLI. This banner may be overwritten by including a banner.txt file at the root of the project that includes the new banner in ascii format. A tool that may be used to generate this banner may be found [here](https://manytools.org/hacker-tools/ascii-banner/)
