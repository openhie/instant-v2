# CLI v2

An initialisation of the cli using cobra and viper

commands:
- project
  - ls
  - -generate, -g \<name>
  - -init, -i
  - -up, -u
  - -down, -d
  - -remove, -r
- package
  - ls
  - -generate, -g \<name>
  - -init, -i \<name>
  - -up, -u \<name>
  - -down, -d \<name>
  - -remove, -r \<name>
- stack
  - ls
  - -generate, -g \<name>
  - -init, -i \<name>
  - -up, -u \<name>
  - -down, -d \<name>
  - -remove, -r \<name>
- config
  - ls
  - -generate, -g \<name>

## Build
`go build .`

## Additional Features over V1
- Auto removes instant container and volume before performing a package action
- Support for multiple env files
- Support for predefined profiles to load in when performing a package action
- Generate a new project prompt
