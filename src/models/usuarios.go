package models

import (
	"errors"
	"mybook-api/src/infrastructure/security"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

// Usuario representa um usuário utilizando a rede social
type Usuario struct {
	ID           string    `json:"id,omitempty"`
	Nome         string    `json:"nome,omitempty"`
	Nick         string    `json:"nick,omitempty"`
	Email        string    `json:"email,omitempty"`
	Senha        string    `json:"senha,omitempty"`
	Seguidores   []string  `json:"seguidores,omitempty"`
	CriadoEm     time.Time `json:"criadoEm,omitempty"`
	AtualizadoEm time.Time `json:"atualizadoEm,omitempty"`
}

// Preparar vai chamar os métodos validar e formatar
func (usuario *Usuario) Preparar() error {
	if err := usuario.validar(); err != nil {
		return err
	}

	if err := usuario.formatar(); err != nil {
		return err
	}

	return nil
}

func (usuario *Usuario) validar() error {
	if usuario.Nome == "" {
		return errors.New("O nome é obrigatório e não pode estar em branco")
	}

	if usuario.Nick == "" {
		return errors.New("O nick é obrigatório e não pode estar em branco")
	}

	if usuario.Email == "" {
		return errors.New("O e-mail é obrigatório e não pode estar em branco")
	}

	if err := checkmail.ValidateFormat(usuario.Email); err != nil {
		return errors.New("O e-mail inserido é inválido")
	}

	if usuario.Senha == "" {
		return errors.New("A senha é obrigatório e não pode estar em branco")
	}

	return nil
}

func (usuario *Usuario) formatar() error {
	usuario.Nome = strings.TrimSpace(usuario.Nome)
	usuario.Nick = strings.TrimSpace(usuario.Nick)
	usuario.Email = strings.TrimSpace(usuario.Email)

	senhaComHash, err := security.Hash(usuario.Senha)
	if err != nil {
		return err
	}

	usuario.Senha = string(senhaComHash)
	return nil
}
