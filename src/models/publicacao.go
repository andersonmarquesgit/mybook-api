package models

import "time"

type Publication struct {
	ID         string    `json:"id,omitempty"`
	Title      string    `json:"title,omitempty"`
	Content    string    `json:"content,omitempty"`
	AuthorID   string    `json:"authorid,omitempty"`
	AuthorNick string    `json:"authornick,omitempty"`
	Likes      uint64    `json:"likes"`
	CreatedIn  time.Time `json:"createdin,omitempty"`
}
