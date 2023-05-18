package mentions

type MentionInput struct {
	Category    string `json:"category"`
	ProductName string `json:"productName"`
	Publisher   string `json:"publisher"`
	DateTime    string `json:"datetime"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Rating      string `json:"rating"`
}

type MentionInputs struct {
	Data []MentionInput `json:"data"`
}

type NewMentionRequest struct {
	Data NewMentionData `json:"data"`
}

type NewMentionData struct {
	Category           string `json:"category"`
	ProductName        string `json:"product_name"`
	Publisher          string `json:"publisher"`
	MentionPublishedAt string `json:"mention_published_at"`
	URL                string `json:"url"`
	Title              string `json:"title"`
	Product            int    `json:"product"`
}
