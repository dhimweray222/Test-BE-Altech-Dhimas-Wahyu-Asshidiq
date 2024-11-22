package request

type AuthorRequest struct {
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	BirthDate string `json:"birth_date"`
}
