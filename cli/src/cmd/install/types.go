package install

type params struct {
	// none, token, basic, custom
	TypeAuth  string
	Token     string
	BasicUser string
	BasicPass string
}

type indexJSON struct {
	IndexVersion int32      `json:"index-version"`
	Files        []filesRep `json:"files"`
}

type filesRep struct {
	Filename     string `json:"filename"`
	ResourceType string `json:"resourceType"`
	Id           string `json:"id"`
	Url          string `json:"url"`
	Version      string `json:"version"`
	Kind         string `json:"kind"`
	Type         string `json:"type"`
}
