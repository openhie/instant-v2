Feature: Test Deploy Commands
  Scenario: Initialise Core Service
    When the command "package init -n=core" is run
    Then check the CLI output is "init -t swarm core"

  Scenario: Up Core Service
    When the command "package up -n=core" is run
    Then check the CLI output is "up -t swarm core"

  Scenario: Down Core Service
    When the command "package down -n=core" is run
    Then check the CLI output is "down -t swarm core"

  Scenario: Destroy Core Service
    When the command "package destroy -n=core" is run
    Then check the CLI output is "destroy -t swarm core"
