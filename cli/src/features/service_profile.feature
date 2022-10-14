Feature: Test Deploy Profiles
  Scenario: Initialise Core Service in Dev Mode
    When the command "package init --profile=dev" is run
    Then check the CLI output is "init -t swarm --dev core client"

  Scenario: Initialise Custom Package only
    When the command "package init --profile=custom-package" is run with profile
      | custom-packages |
      | disi-on-platform |
    Then check that the CLI added custom packages
      | directory |
      | disi-on-platform |
    Then check the CLI output is "init -t swarm disi-on-platform"

  Scenario: Initialise Mixed packages
    When the command "package init --profile=mixed-package" is run with profile
      | custom-packages |
      | disi-on-platform |
    Then check that the CLI added custom packages
      | directory |
      | disi-on-platform |
    Then check the CLI output is "init -t swarm core client disi-on-platform"

  Scenario: Initialise Custom package from local path in Dev Mode
    When the command "package init --profile=local-custom-package" is run with profile
      | custom-packages   |
      | cares-on-platform |
    Then check that the CLI added custom packages
      | directory         |
      | cares-on-platform |
    Then check the CLI output is "init -t swarm --dev cares-on-platform"

  Scenario: Initialise mixed packages in Dev Mode
    When the command "package init --profile=mixed-custom-package" is run with profile
      | custom-packages   |
      | cares-on-platform |
      | disi-on-platform  |
    Then check that the CLI added custom packages
      | directory         |
      | cares-on-platform |
      | disi-on-platform  |
    Then check the CLI output is "init -t swarm --dev core disi-on-platform cares-on-platform"

  Scenario: Initialise mixed packages in Dev Mode
    When the command "package init --profile=mixed-custom-package" is run with profile
      | custom-packages   |
      | cares-on-platform |
      | disi-on-platform  |
    Then check that the CLI added custom packages
      | directory         |
      | cares-on-platform |
      | disi-on-platform  |
    Then check the CLI output is "init -t swarm --dev core disi-on-platform cares-on-platform"