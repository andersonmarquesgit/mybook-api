package models

import "time"

type Seguidores struct {
	UserID       string    `json:"userid,omitempty"`
	Seguidores   []string  `json:"seguidores,omitempty"`
	AtualizadoEm time.Time `json:"atualizadoEm,omitempty"`
}
