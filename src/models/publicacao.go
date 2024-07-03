package models

import (
	"errors"
	"strings"
	"time"
)

type Publication struct {
	ID         string    `json:"id,omitempty"`
	Title      string    `json:"title,omitempty"`
	Content    string    `json:"content,omitempty"`
	AuthorID   string    `json:"authorid,omitempty"`
	AuthorNick string    `json:"authornick,omitempty"`
	Likes      uint64    `json:"likes"`
	CreatedIn  time.Time `json:"createdin,omitempty"`
}

func (publication *Publication) Preparar() error {
	if err := publication.validar(); err != nil {
		return err
	}

	if err := publication.formatar(); err != nil {
		return err
	}

	return nil
}

func (publication *Publication) validar() error {
	if publication.Title == "" {
		return errors.New("O título é obrigatório e não pode estar em branco")
	}

	if publication.Content == "" {
		return errors.New("O conteúdo é obrigatório e não pode estar em branco")
	}

	if publication.AuthorID == "" {
		return errors.New("O autor é obrigatório e não pode estar em branco")
	}

	return nil
}

func (publication *Publication) formatar() error {
	publication.Title = strings.TrimSpace(publication.Title)
	publication.Content = strings.TrimSpace(publication.Content)
	return nil
}
