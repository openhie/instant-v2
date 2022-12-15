Feature: Test Deploy Options Commands
  Scenario: Initialise Core Service in Dev Mode
    When the command "package init -n=core --dev" is run
    Then check the CLI output is "init -t swarm --dev core"

  Scenario: Initialise Core Service With Custom config.yml and .env File
    When the command "package init -n=core --config=features/test-conf/config.yml --env-file=features/test-conf/.env.test" is run
    Then check the CLI output is "init -t swarm core"

  Scenario: Down only Core Service
    When the command "package down -n=core --only" is run
    Then check the CLI output is "down -t swarm --only core"

  Scenario: Initialise Template Custom Service With Custom .env File
    When the command "package init -n=template -c=https://github.com/jembi/instant-openhie-template-package.git --env-file=features/test-conf/.env.test" is run
    Then check that the CLI added custom packages
      | directory |
      | instant-openhie-template-package |
    Then check the CLI output is "init -t swarm template"

  Scenario: Initialise Multiple Services
    When the command "package init -n=covid19immunization -n=client -n=covid19surveillance --custom-path=https://github.com/jembi/covid19-immunization-tracking-package.git -c=https://github.com/jembi/who-covid19-surveillance-package.git" is run
    Then check that the CLI added custom packages
      | directory |
      | covid19-immunization-tracking-package |
      | who-covid19-surveillance-package |
    Then check the CLI output is "init -t swarm covid19immunization client covid19surveillance"
