# Package Metadata

A package's metadata is defined in a `package-metadata.json` file. This file allows a package to be detected as an Instant OpenHIE package and gives key metadata that Instant OpenHIE needs.

### Example

{% code title="package-metadata.json" lineNumbers="true" %}
```json
{
  "id": "dashboard-visualiser-kibana",
  "name": "Dashboard Visualiser Kibana",
  "description": "A dashboard to interpret the data from the ElasticSearch data store",
  "type": "infrastructure|use-case",
  "version": "0.0.1",
  "dependencies": ["analytics-datastore-elastic-search"],
  "environmentVariables": {
    "KIBANA_INSTANCES": 1,
    "KIBANA_USERNAME": "elastic",
    "KIBANA_PASSWORD": "dev_password_only",
    "KIBANA_SSL": "true",
    "KIBANA_MEMORY_LIMIT": "3G",
    "KIBANA_MEMORY_RESERVE": "500M"
  }
}
```
{% endcode %}

{% hint style="warning" %}
Environment Variables in the `environmentVariables` section will be used as defaults, but will be overridden by matching variables in an environment variable file
{% endhint %}
