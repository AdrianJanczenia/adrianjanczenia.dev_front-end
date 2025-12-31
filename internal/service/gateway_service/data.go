package gateway_service

type PageContent struct {
	Meta         Meta              `json:"meta"`
	Profile      Profile           `json:"profile"`
	Skills       []SkillGroup      `json:"skills"`
	Experience   []Job             `json:"experience"`
	Languages    []Language        `json:"languages"`
	Contact      Contact           `json:"contact"`
	Translations map[string]string `json:"translations"`
}

type Meta struct {
	Title string `json:"title"`
}

type Profile struct {
	Name     string   `json:"name"`
	Headline string   `json:"headline"`
	Location string   `json:"location"`
	About    string   `json:"about"`
	Tags     []string `json:"tags"`
}

type SkillGroup struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

type Job struct {
	Role             string   `json:"role"`
	Company          string   `json:"company"`
	Period           string   `json:"period"`
	Location         string   `json:"location"`
	Type             string   `json:"type"`
	Summary          string   `json:"summary"`
	Responsibilities []string `json:"responsibilities"`
	SkillsUsed       []string `json:"skills_used"`
}

type Language struct {
	Language    string `json:"language"`
	Proficiency string `json:"proficiency"`
}

type Contact struct {
	Email    string `json:"email"`
	Linkedin string `json:"linkedin"`
	Github   string `json:"github"`
}
