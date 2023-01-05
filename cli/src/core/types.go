package core

type Profile struct {
	Name     string   `yaml:"name"`
	EnvFiles []string `yaml:"envFiles"`
	Dev      bool     `yaml:"dev"`
	Only     bool     `yaml:"only"`
	Packages []string `yaml:"packages"`
}

type CustomPackage struct {
	Id   string `yaml:"id"`
	Path string `yaml:"path"`
}

type Config struct {
	Image          string          `yaml:"image"`
	LogPath        string          `yaml:"logPath"`
	Packages       []string        `yaml:"packages"`
	CustomPackages []CustomPackage `yaml:"customPackages"`
	Profiles       []Profile       `yaml:"profiles"`
	ProjectName    string          `yaml:"projectName"`
	PlatformImage  string          `yaml:"platformImage"`
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
