package adrianjanczenia_dev_content_service

import (
	"net/http"
)

// Client is responsible for communicating with the external content service API.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

// GetPageContent fetches the page content for a specific language.
func (c *Client) GetPageContent(lang string) (*PageContent, error) {
	if lang == "en" {
		return getEnglishMockData(), nil
	}
	return getPolishMockData(), nil

	/*
		// --- ORIGINAL HTTP REQUEST LOGIC (TEMPORARILY DISABLED) ---

		url := fmt.Sprintf("%s/content?lang=%s", c.baseURL, lang)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("User-Agent", "PortfolioFrontend/1.0")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to execute request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("api returned non-200 status: %d", resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			return nil, fmt.Errorf("unexpected content type: got %s, want application/json", contentType)
		}

		var pageContent PageContent
		if err := json.NewDecoder(resp.Body).Decode(&pageContent); err != nil {
			return nil, fmt.Errorf("failed to decode json response: %w", err)
		}

		return &pageContent, nil
	*/
}

func getPolishMockData() *PageContent {
	return &PageContent{
		Meta: Meta{Title: "Jan Kowalski - Backend Go Developer"},
		Profile: Profile{
			Name:   "Jan Kowalski",
			Role:   "Backend Go Developer",
			Status: "Otwarty na nowe wyzwania",
			Bio:    "Jestem doświadczonym inżynierem oprogramowania z pasją do budowania solidnych i wydajnych systemów backendowych. Moim głównym językiem programowania jest Go.",
		},
		Skills: []SkillGroup{
			{Key: "languages", Values: []string{"Go", "Python", "SQL"}},
			{Key: "databases", Values: []string{"PostgreSQL", "Redis", "MongoDB"}},
			{Key: "devops", Values: []string{"Docker", "Kubernetes", "GitLab CI", "Terraform"}},
		},
		Experience: []Job{
			{Role: "Backend Go Developer", Company: "Tech Corp", Period: "2022 - Obecnie", Description: "Projektowanie i implementacja mikroserwisów w Go. Praca nad systemami o wysokiej dostępności."},
			{Role: "Software Tester", Company: "Quality Solutions", Period: "2018 - 2022", Description: "Automatyzacja testów dla aplikacji webowych w Python (Selenium, Pytest)."},
		},
		Contact: Contact{Email: "jan.kowalski@email.com", LinkedIn: "jankowalski-dev"},
		Translations: Translations{
			NavAbout: "o-mnie", NavSkills: "umiejetnosci", NavExperience: "doswiadczenie", NavContact: "kontakt", NavCV: "pobierz-cv",
			HeaderAbout: "o-mnie", HeaderSkills: "technologie", HeaderExperience: "doswiadczenie", HeaderContact: "kontakt",
			ProfileKeyName: "imie", ProfileKeyRole: "rola", ProfileKeyStatus: "status", ProfileKeyBio: "bio",
			ContactKeyEmail: "email", ContactKeyLinkedIn: "linkedin",
			BackToTopLabel: "Powrót na górę", ErrorTitle: "Wystąpił błąd", ErrorMessage: "Pracujemy nad rozwiązaniem problemu. Proszę spróbować ponownie.",
		},
	}
}

func getEnglishMockData() *PageContent {
	return &PageContent{
		Meta: Meta{Title: "John Smith - Backend Go Developer"},
		Profile: Profile{
			Name:   "John Smith",
			Role:   "Backend Go Developer",
			Status: "Open for new opportunities",
			Bio:    "I am an experienced software engineer passionate about building robust and efficient backend systems. My primary programming language is Go.",
		},
		Skills: []SkillGroup{
			{Key: "languages", Values: []string{"Go", "Python", "SQL"}},
			{Key: "databases", Values: []string{"PostgreSQL", "Redis", "MongoDB"}},
			{Key: "devops", Values: []string{"Docker", "Kubernetes", "GitLab CI", "Terraform"}},
		},
		Experience: []Job{
			{Role: "Backend Go Developer", Company: "Tech Corp", Period: "2022 - Present", Description: "Designing and implementing microservices in Go. Working on high-availability and low-latency systems."},
			{Role: "Software Tester", Company: "Quality Solutions", Period: "2018 - 2022", Description: "Test automation for web applications in Python (Selenium, Pytest)."},
		},
		Contact: Contact{Email: "john.smith@email.com", LinkedIn: "johnsmith-dev"},
		Translations: Translations{
			NavAbout: "about", NavSkills: "skills", NavExperience: "experience", NavContact: "contact", NavCV: "download-cv",
			HeaderAbout: "about", HeaderSkills: "technologies", HeaderExperience: "experience", HeaderContact: "contact",
			ProfileKeyName: "name", ProfileKeyRole: "role", ProfileKeyStatus: "status", ProfileKeyBio: "bio",
			ContactKeyEmail: "email", ContactKeyLinkedIn: "linkedin",
			BackToTopLabel: "Back to top", ErrorTitle: "An error occurred", ErrorMessage: "We are working on fixing the problem. Please try again later.",
		},
	}
}
