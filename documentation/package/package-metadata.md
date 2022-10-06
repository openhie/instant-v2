# Package Metadata

A package's metadata is defined in a `package-metadata.json` file

### Example

{% code title="package-metadata.json" lineNumbers="true" %}
```json
{
  "id": "dashboard-visualiser-kibana",
  "name": "Dashboard Visualiser Kibana",
  "description": "A dashboard to interpret the data from the ElasticSearch data store",
  "type": "infrastructure",
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
