package core

type Profile struct {
	Name     string   `yaml:"name"`
	Packages []string `yaml:"packages"`
	EnvFiles []string `yaml:"envFiles"`
	EnvVars  []string `yaml:"envVars"`
	Dev      bool     `yaml:"dev"`
	Only     bool     `yaml:"only"`
}

type CustomPackage struct {
	Id   string `yaml:"id"`
	Path string `yaml:"path"`
}

type Config struct {
	ProjectName    string          `yaml:"projectName,omitempty"`
	Image          string          `yaml:"image,omitempty"`
	PlatformImage  string          `yaml:"platformImage,omitempty"`
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
}

type GeneratePackageSpec struct {
	Id             string
	Name           string
	Image          string
	Stack          string
	Description    string
	Type           string
	IncludeDevFile bool
	TargetPort     string
	PublishedPort  string
}
