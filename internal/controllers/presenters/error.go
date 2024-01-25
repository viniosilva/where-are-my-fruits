package presenters

type ErrorRes struct {
	Error    string   `json:"error" example:"Not Found"`
	Message  string   `json:"message,omitempty" example:"fruit not found"`
	Messages []string `json:"messages,omitempty" example:"invalid field,invalid value"`
}
