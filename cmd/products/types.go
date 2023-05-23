package products

type ProductInputs struct {
	Data []ProductInputData `json:"data"`
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

type BrandRequestAttributes struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
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
	Collections  []int    `json:"collections"`
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

type BrandsResponse struct {
	Data []BrandData `json:"data"`
}

type BrandResponse struct {
	Data BrandData `json:"data"`
}

type ProductData struct {
	Id         int `json:"id"`
	Attributes ProductAttributes
}

type ProductResponse struct {
	Data ProductData `json:"data"`
}

type CollectionResponseAttributes struct {
	Name string `json:"name"`
}

type BrandRequestBody struct {
	Data BrandRequestAttributes `json:"data"`
}

type CollectionRequestBody struct {
	Data CollectionRequestAttributes `json:"data"`
}

type CollectionData struct {
	Id         int                           `json:"id"`
	Attributes *CollectionResponseAttributes `json:"data"`
}

type CollectionResponse struct {
	Data CollectionData `json:"data"`
}

type CollectionsResponse struct {
	Data []CollectionData `json:"data"`
}

type CollectionRequestAttributes struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}
