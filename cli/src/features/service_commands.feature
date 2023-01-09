Feature: Test Package Deploy Commands
  Scenario: Initialise Core Service
    When the command "package init -n=core" is run
    Then check the CLI output is "init -t swarm core"

  Scenario: Return Error From No Packages Specified
    When the command "package init" is run in error
    Then check the CLI output is "no packages selected in any of command-line/profiles, use the 'project' command for project level functions"

  Scenario: Up Core Service
    When the command "package up -n=core" is run
    Then check the CLI output is "up -t swarm core"

  Scenario: Down Core Service
    When the command "package down -n=core" is run
    Then check the CLI output is "down -t swarm core"

  Scenario: Destroy Core Service
    When the command "package destroy -n=core" is run
    Then check the CLI output is "destroy -t swarm core"
