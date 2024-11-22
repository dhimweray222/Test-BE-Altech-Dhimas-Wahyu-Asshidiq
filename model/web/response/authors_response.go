package response

type AuthorResponse struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	BirthDate string `json:"birth_date"`
}
