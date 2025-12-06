package repo

import "shortener/internal/repo/url"

type Repo struct {
	URL url.URLInterface
}

func NewRepo(url url.URLInterface) *Repo {
	return &Repo{
		URL: url,
	}
}
