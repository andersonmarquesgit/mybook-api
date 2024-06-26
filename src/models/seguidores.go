package models

import "time"

type Seguidores struct {
	UserID       string    `json:"userid,omitempty"`
	Followers    []string  `json:"followers,omitempty"`
	AtualizadoEm time.Time `json:"atualizadoEm,omitempty"`
}

type Seguindo struct {
	UserID       string    `json:"userid,omitempty"`
	Following    []string  `json:"following,omitempty"`
	AtualizadoEm time.Time `json:"atualizadoEm,omitempty"`
}
