package model

type URL struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ShortRef string `json:"short_ref"`
}
