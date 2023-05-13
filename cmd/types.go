package cmd

type ProductInputs struct {
	Data []ProductInputData `json"data"`
}

type ProductInputData struct {
	Name         string   `json:"name"`
	Amount       int      `json:"amount"`
	Thumbnail    string   `json:"thumbnail"`
	PrimaryImage string   `json:"primary_image"`
	OriginalURL  string   `json:"original_url"`
	Description  string   `json:"description"`
	Brand        string   `json:"brand"`
	Collections  string   `json:"collections"`
	Media        []string `json:"media"`
}

type ProductAttributes struct {
	Name            string   `json:"name"`
	Amount          int      `json:"amount"`
	Thumbnail       string   `json:"thumbnail"`
	PrimaryImage    string   `json:"primary_image"`
	ProductMetadata []string `json:"product_metadata"`
	Media           []string `json:"media"`
	OriginalURL     string   `json:"original_url"`
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
	Slug            string   `json:"slug"`
	Description     string   `json:"description"`
	GoodFor         []string `json:"goodFor"`
	AgeRange        string   `json:"age_range"`
	Featured        bool     `json:"featured"`
}

type ImageData struct {
	Name string
	Body []byte
}

type ProductRequestAttributes struct {
	Name         string   `json:"name"`
	Amount       int      `json:"amount"`
	Brand        int      `json:"brand"`
	Thumbnail    string   `json:"thumbnail"`
	PrimaryImage string   `json:"primary_image"`
	OriginalURL  string   `json:"original_url"`
	Description  string   `json:"description"`
	Media        []string `json:"media"`
	Slug         string   `json:"slug"`
}
type ProductRequest struct {
	Data *ProductRequestAttributes `json:"data"`
}

type BrandResponseAttributes struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Logo        string `json:"primary_image"`
	CoverPhoto  string `json:"cover_photo"`
}

type BrandData struct {
	Id         int                      `json:"id"`
	Attributes *BrandResponseAttributes `json:"data"`
}

type BrandResponse struct {
	Data []BrandData `json:"data"`
}

type ProductData struct {
	Id         int `json:"id"`
	Attributes ProductAttributes
}

type ProductResponse struct {
	Data ProductData `json:"data"`
}
