Feature: Test Deploy Profiles
  Scenario: Initialise Core Service in Dev Mode
    When the command "package init --profile=test-dev" is run
    Then check the CLI output is "init -t swarm --dev core client"

  Scenario: Initialise Custom Package only
    When the command "package init --profile=test-custom-package" is run with profile
      | custom-packages |
      | disi-on-platform |
    Then check that the CLI added custom packages
      | directory |
      | disi-on-platform |
    Then check the CLI output is "init -t swarm disi-on-platform"

  Scenario: Initialise Mixed packages
    When the command "package init --profile=test-mixed-package" is run with profile
      | custom-packages |
      | disi-on-platform |
    Then check that the CLI added custom packages
      | directory |
      | disi-on-platform |
    Then check the CLI output is "init -t swarm core client disi-on-platform"

  Scenario: Initialise Custom package from local path in Dev Mode
    When the command "package init --profile=test-local-custom-package" is run with profile
      | custom-packages   |
      | custom-local-package |
    Then check that the CLI added custom packages
      | directory         |
      | custom-local-package |
    Then check the CLI output is "init -t swarm --dev custom-local-package"

  Scenario: Initialise mixed packages in Dev Mode
    When the command "package init --profile=test-mixed-custom-package" is run with profile
      | custom-packages   |
      | custom-local-package |
      | disi-on-platform  |
    Then check that the CLI added custom packages
      | directory         |
      | custom-local-package |
      | disi-on-platform  |
    Then check the CLI output is "init -t swarm --dev core disi-on-platform custom-local-package"

  Scenario: Initialise mixed packages in Dev Mode
    When the command "package init --profile=test-mixed-custom-package" is run with profile
      | custom-packages   |
      | custom-local-package |
      | disi-on-platform  |
    Then check that the CLI added custom packages
      | directory         |
      | custom-local-package |
      | disi-on-platform  |
    Then check the CLI output is "init -t swarm --dev core disi-on-platform custom-local-package"
