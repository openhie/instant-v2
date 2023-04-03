# Script

A `swarm.sh` file acts as an entry point to your package and runs within the instant container during deploy.

Two arguments are passed by default into the swarm script\
\- $1 is the action type ( init|up|down|destroy )\
\- $2 is the MODE in which it is run (dev)

Due to this script running in the instant container, all references made to files within the package folder would need to be prefixed with the `PACKAGE_PATH` variable

To supply config option to your package, make use of env vars which will be made available to this script and therefore to any docker command that you execute (so you may use env vars in your compose files for example). There are varies option on the CLI via flags or the config file to supply env var files or env vars themselves.

As a coding standard we encourage use of the [Shell Style Guide](https://google.github.io/styleguide/shellguide.html)

Should you use VS Code for editing we suggest the `pinage404.bash-extension-pack`&#x20;

### Example

{% code title="swarm.sh" lineNumbers="true" %}
```bash
#!/bin/bash

readonly ACTION=$1
readonly MODE=$2

PACKAGE_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")" || exit
  pwd -P
)
readonly PACKAGE_PATH

main() {
  if [[ "${MODE}" == "dev" ]]; then
    echo "Running Dashboard Visualiser Kibana package in DEV mode"
    kibana_dev_compose_param="-c ${PACKAGE_PATH}/docker-compose.dev.yml"
  else
    echo "Running Dashboard Visualiser Kibana package in PROD mode"
    kibana_dev_compose_param=""
  fi

  if [[ "${ACTION}" == "init" ]] || [[ "${ACTION}" == "up" ]]; then
    docker stack deploy -c "${PACKAGE_PATH}"/docker-compose.yml $kibana_dev_compose_param instant
  elif [[ "${ACTION}" == "down" ]]; then
    docker service scale instant_dashboard-visualiser-kibana=0
  elif [[ "${ACTION}" == "destroy" ]]; then
    docker service remove instant_dashboard-visualiser-kibana
  else
    echo "Valid options are: init, up, down, or destroy"
  fi
}

main "$@"
```
{% endcode %}

#### Breakdown

* Lines 2 & 3 extract the two arguments that instant provides to this script during any deploy command involving this package ie. the ACTION and MODE respectively.
* Lines 6-9 retrieve the current path to this swarm.sh file which should exist at the root of your package. This path may then be used to reference any files within the package eg. docker-compose.yml.
* Lines 12-32 define and execute a main block within which all our executable code will run. This is done to preserve the scope of any locally scoped variables used in the script.
* Lines 13-19 specify the path to a docker-compose file that contains override configs specifically for dev use. This variable then gets used on line 22 to override the default configs of the main docker-compose file.
* Lines 21-29 define logic based on the ACTION parameter
* Line 22 deploys all services specified in the docker compose files provided and assign them to the `instant` stack. The dev docker compose file is only used if the MODE argument was received as "dev".
* Line 24 scales down the specified services to 0. This stops all containers but keeps their volumes and configs intact which will allow you to perform maintenance without losing data.
* Line 26 removes the service which will also stop and remove any of it's containers. Volume removal may also occur here
