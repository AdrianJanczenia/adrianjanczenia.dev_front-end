package gateway_service

type PageContent struct {
	SchemaVersion int    `json:"schema_version"`
	Meta          Meta   `json:"meta"`
	Shared        Shared `json:"shared"`
	Pages         Pages  `json:"pages"`
	UI            UI     `json:"ui"`
}

type Meta struct {
	Locale      string `json:"locale"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Shared struct {
	Navigation   Navigation   `json:"navigation"`
	Footer       Footer       `json:"footer"`
	ContactLinks ContactLinks `json:"contact_links"`
}

type Navigation struct {
	About      string `json:"about"`
	Experience string `json:"experience"`
	HowItWorks string `json:"how_it_works"`
	Contact    string `json:"contact"`
	CV         string `json:"cv"`
	Menu       string `json:"menu"`
}

type Footer struct {
	Privacy string `json:"privacy"`
}
type ContactLinks struct {
	Email    string `json:"email"`
	Github   string `json:"github"`
	Linkedin string `json:"linkedin"`
}

type TextBlock struct {
	Eyebrow string `json:"eyebrow"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Action  string `json:"action"`
}

type Pages struct {
	About      AboutPage      `json:"about"`
	Experience ExperiencePage `json:"experience"`
	HowItWorks HowItWorksPage `json:"how_it_works"`
	Contact    ContactPage    `json:"contact"`
	Privacy    PrivacyPolicy  `json:"privacy"`
}

type AboutPage struct {
	Hero           TextBlock `json:"hero"`
	CareerSnapshot TextBlock `json:"career_snapshot"`
	SiteTeaser     TextBlock `json:"site_teaser"`
	Name           string    `json:"name"`
	Role           string    `json:"role"`
	IdentityDetail string    `json:"identity_detail"`
	Tags           []string  `json:"tags"`
}

type HowItWorksPage struct {
	Eyebrow      string                       `json:"eyebrow"`
	Title        string                       `json:"title"`
	Body         string                       `json:"body"`
	UICopy       map[string]string            `json:"ui_copy"`
	ScenarioCopy map[string]map[string]string `json:"scenario_copy"`
}

type OrderedRoles struct {
	Order []string        `json:"order"`
	Items map[string]Role `json:"items"`
}
type Role struct {
	Title   string   `json:"title"`
	Period  string   `json:"period"`
	Summary string   `json:"summary"`
	Tasks   []string `json:"tasks"`
}

type OrderedCases struct {
	Order []string        `json:"order"`
	Items map[string]Case `json:"items"`
}
type Case struct {
	Category     string   `json:"category"`
	Title        string   `json:"title"`
	Summary      string   `json:"summary"`
	Problem      string   `json:"problem"`
	Contribution []string `json:"contribution"`
	Result       string   `json:"result"`
	Takeaway     string   `json:"takeaway"`
	TechContext  []string `json:"tech_context"`
}

type OrderedTechnologyGroups struct {
	Order []string                   `json:"order"`
	Items map[string]TechnologyGroup `json:"items"`
}
type TechnologyGroup struct {
	Title string           `json:"title"`
	Items []TechnologyItem `json:"items"`
}
type TechnologyItem struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type ExperiencePage struct {
	Intro             TextBlock               `json:"intro"`
	Roles             OrderedRoles            `json:"roles"`
	Cases             OrderedCases            `json:"cases"`
	TechnologyGroups  OrderedTechnologyGroups `json:"technology_groups"`
	RolesLabel        string                  `json:"roles_label"`
	CasesLabel        string                  `json:"cases_label"`
	TechnologyLabel   string                  `json:"technology_label"`
	ReadCase          string                  `json:"read_case"`
	Back              string                  `json:"back"`
	ProblemLabel      string                  `json:"problem_label"`
	ContributionLabel string                  `json:"contribution_label"`
	ResultLabel       string                  `json:"result_label"`
}

type ContactPage struct {
	Intro       TextBlock `json:"intro"`
	EmailCard   TextBlock `json:"email_card"`
	ProfileCard TextBlock `json:"profile_card"`
}

type PrivacyPolicy struct {
	Title    string           `json:"title"`
	Sections []PrivacySection `json:"sections"`
}
type PrivacySection struct {
	Header string        `json:"header"`
	Items  []PrivacyItem `json:"items"`
}
type PrivacyItem struct {
	Label string `json:"label"`
	Text  string `json:"text"`
}

type UI struct {
	CV          CVUI                 `json:"cv"`
	EmailCopied string               `json:"email_copied"`
	Errors      map[string]ErrorCopy `json:"errors"`
}
type CVUI struct {
	Title           string `json:"title"`
	Intro           string `json:"intro"`
	AccessCodeLabel string `json:"access_code_label"`
	CaptchaLabel    string `json:"captcha_label"`
	Submit          string `json:"submit"`
	Retry           string `json:"retry"`
	Close           string `json:"close"`
	ContactPrompt   string `json:"contact_prompt"`
}
type ErrorCopy struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Action  string `json:"action"`
}
