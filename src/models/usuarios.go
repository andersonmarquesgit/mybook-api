package models

import "time"

// Usuario representa um usuário utilizando a rede social
type Usuario struct {
	ID           string    `json:"id,omitempty"`
	Nome         string    `json:"nome,omitempty"`
	Nick         string    `json:"nick,omitempty"`
	Email        string    `json:"email,omitempty"`
	Senha        string    `json:"senha,omitempty"`
	CriadoEm     time.Time `json:"criadoEm,omitempty"`
	AtualizadoEm time.Time `json:"atualizadoEm,omitempty"`
}
