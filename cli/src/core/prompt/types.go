package prompt

type generatePackagePromptResponse struct {
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
