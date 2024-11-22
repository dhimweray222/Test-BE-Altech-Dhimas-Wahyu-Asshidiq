package response

type BookResponse struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	AuthorId    string `json:"author_id"`
	PublishDate string `json:"publish_date"`
	Description string `json:"description"`
	AuthorName  string `json:"author_name"`
}
