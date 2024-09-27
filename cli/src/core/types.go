package core

type Profile struct {
	Name     string   `yaml:"name"`
	Packages []string `yaml:"packages"`
	EnvVars  []string `yaml:"envVars,omitempty"`
	EnvFiles []string `yaml:"envFiles,omitempty"`
	Dev      bool     `yaml:"dev,omitempty"`
	Only     bool     `yaml:"only,omitempty"`
}

type CustomPackage struct {
	Id   string `yaml:"id"`
	Path string `yaml:"path"`
}

type Config struct {
	ProjectName    string          `yaml:"projectName,omitempty"`
	Image          string          `yaml:"image,omitempty"`
	LogPath        string          `yaml:"logPath,omitempty"`
	Packages       []string        `yaml:"packages,omitempty"`
	CustomPackages []CustomPackage `yaml:"customPackages,omitempty"`
	Profiles       []Profile       `yaml:"profiles,omitempty"`
}

type PackageSpec struct {
	EnvironmentVariables []string
	DeployCommand        string
	Packages             []string
	IsDev                bool
	IsOnly               bool
	CustomPackages       []CustomPackage
	ImageVersion         string
	TargetLauncher       string
	Concurrency          string
}

type GeneratePackageSpec struct {
	Id             string
	Name           string
	Image          string
	Description    string
	Type           string
	IncludeDevFile bool
	TargetPort     string
	PublishedPort  string
}
