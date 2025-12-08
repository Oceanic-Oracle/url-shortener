package dto

// requests

type CreateCodeURLRequest struct {
	URL string `json:"url"`
}

// responses

type CreateCodeURLResponse struct {
	Code string `json:"code"`
}
