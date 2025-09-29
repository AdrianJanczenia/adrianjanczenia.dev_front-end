package adrianjanczenia_dev_content_service

// PageContent matches the JSON structure returned from Content Service API.
type PageContent struct {
	Meta         Meta         `json:"meta"`
	Profile      Profile      `json:"profile"`
	Skills       []SkillGroup `json:"skills"`
	Experience   []Job        `json:"experience"`
	Contact      Contact      `json:"contact"`
	Translations Translations `json:"translations"`
}

type Meta struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Profile struct {
	Name   string `json:"name"`
	Role   string `json:"role"`
	Status string `json:"status"`
	Bio    string `json:"bio"`
}

type SkillGroup struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

type Job struct {
	Role        string `json:"role"`
	Company     string `json:"company"`
	Period      string `json:"period"`
	Description string `json:"description"`
}

type Contact struct {
	Email    string `json:"email"`
	LinkedIn string `json:"linkedin"`
}

type Translations struct {
	NavAbout           string `json:"nav_about"`
	NavSkills          string `json:"nav_skills"`
	NavExperience      string `json:"nav_experience"`
	NavContact         string `json:"nav_contact"`
	NavCV              string `json:"nav_cv"`
	HeaderAbout        string `json:"header_about"`
	HeaderSkills       string `json:"header_skills"`
	HeaderExperience   string `json:"header_experience"`
	HeaderContact      string `json:"header_contact"`
	ProfileKeyName     string `json:"profile_key_name"`
	ProfileKeyRole     string `json:"profile_key_role"`
	ProfileKeyStatus   string `json:"profile_key_status"`
	ProfileKeyBio      string `json:"profile_key_bio"`
	ContactKeyEmail    string `json:"contact_key_email"`
	ContactKeyLinkedIn string `json:"contact_key_linkedin"`
	BackToTopLabel     string `json:"back_to_top_label"`
	ErrorTitle         string `json:"error_title"`
	ErrorMessage       string `json:"error_message"`
}
