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

{% hint style="warning" %}
* Packages in a profile must be specified in either the customPackages or packages section
{% endhint %}

## Launching packages

Once a config has been defined with 1 or more packages, you may launch or stop packages by using the [`./instant package <init|up|down|remove> -n <package_id>` command](cli.md#package).

## Launching projects

Instead of individual packages, you can also launch everything define in your project config by using the [`./instant project <init|up|down|remove>` command](cli.md#project).
