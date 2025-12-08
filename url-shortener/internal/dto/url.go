package dto

import "net/url"

// requests

type CreateCodeURLRequest struct {
	URL string `json:"url"`
}

// responses

type CreateCodeURLResponse struct {
	Code string `json:"code"`
}

// methods

func (cc *CreateCodeURLRequest) Validate(host string) error {
	u, err := url.ParseRequestURI(cc.URL)
	if err != nil {
		return ErrParseURL
	}

	if u.Scheme == "" || u.Host == "" {
		return ErrURLMissingSchemaOrHost
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return ErrURLUnsupportedScheme
	}

	if u.Host == host {
		return ErrURLPointsToService
	}

	return nil
}
