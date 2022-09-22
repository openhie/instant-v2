package config

var CustomOptions CustomOption

var Cfg Config

type Config struct {
	Image                        string    `yaml:"image"`
	DefaultTargetLauncher        string    `yaml:"defaultTargetLauncher"`
	Packages                     []Package `yaml:"packages"`
	DisableKubernetes            bool      `yaml:"disableKubernetes"`
	DisableIG                    bool      `yaml:"disableIG"`
	DisableCustomTargetSelection bool      `yaml:"disableCustomTargetSelection"`
	LogPath                      string    `yaml:"logPath"`
}

type Package struct {
	Name string `yaml:"name"`
	ID   string `yaml:"id"`
}

type CustomOption struct {
	StartupAction              string
	StartupPackages            []string
	EnvVarFileLocation         string
	EnvVars                    []string
	CustomPackageFileLocations []string
	OnlyFlag                   bool
	ImageVersion               string
	TargetLauncher             string
	DevMode                    bool
}

type Params struct {
	// none, token, basic, custom
	TypeAuth  string
	Token     string
	BasicUser string
	BasicPass string
}

func init() {
	CustomOptions = CustomOption{
		StartupAction:      "init",
		EnvVarFileLocation: "",
		OnlyFlag:           false,
		ImageVersion:       "latest",
		TargetLauncher:     "docker",
		DevMode:            false,
	}
}
