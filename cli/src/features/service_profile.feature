Feature: Test Deploy Profiles
  Scenario: Initialise Core Service in Dev Mode
    When the command "package init --profile=test-dev" is run
    Then check the CLI output is "init -t swarm --dev core client"

  Scenario: Initialise Custom Package only
    When the command "package init --profile=test-custom-package" is run with profile
      | custom-packages |
      | custom-package-test |
    Then check that the CLI added custom packages
      | directory |
      | custom-package-test |
    Then check the CLI output is "init -t swarm custom-package-test"

  Scenario: Initialise Mixed packages
    When the command "package init --profile=test-mixed-package" is run with profile
      | custom-packages |
      | custom-package-test |
    Then check that the CLI added custom packages
      | directory |
      | custom-package-test |
    Then check the CLI output is "init -t swarm core client custom-package-test"

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
      | custom-package-test  |
    Then check that the CLI added custom packages
      | directory         |
      | custom-local-package |
      | custom-package-test  |
    Then check the CLI output is "init -t swarm --dev core custom-package-test custom-local-package"

  Scenario: Initialise mixed packages in Dev Mode
    When the command "package init --profile=test-mixed-custom-package" is run with profile
      | custom-packages   |
      | custom-local-package |
      | custom-package-test  |
    Then check that the CLI added custom packages
      | directory         |
      | custom-local-package |
      | custom-package-test  |
    Then check the CLI output is "init -t swarm --dev core custom-package-test custom-local-package"
