# Config

A single file may be used for configuration

{% hint style="warning" %}
The \`config.yml\` file should be located at the root of the project and is required for cli usage.
{% endhint %}

The config schema may be found [here](../schema/config.schema.json)

To load this schema into a vscode project include this setting in your `.vscode/settings.json` file

{% code title=".vscode/settings.json" %}
```json
{
  "yaml.schemas": {
    "https://raw.githubusercontent.com/openhie/package-starter-kit/main/schema/config.schema.json": "config.yml"
  }
}
```
{% endcode %}
