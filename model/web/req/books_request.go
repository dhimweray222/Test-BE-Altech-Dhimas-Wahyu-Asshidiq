package request

type BookRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	AuthorId    string `json:"author_id" validate:"required"`
	PublishDate string `json:"publish_date" validate:"required"`
}
