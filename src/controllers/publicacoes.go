package controllers

import (
	"encoding/json"
	"io/ioutil"
	"mybook-api/src/infrastructure/autenticacao"
	"mybook-api/src/models"
	publications "mybook-api/src/repository/publication"
	"mybook-api/src/response"
	"net/http"

	"github.com/gorilla/mux"
)

func CriarPublicacoes(w http.ResponseWriter, r *http.Request) {
	userID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	var publication models.Publication
	publication.AuthorID = userID
	if err = json.Unmarshal(request, &publication); err != nil {
		response.JSON(w, http.StatusBadRequest, "Erro ao converter o request para struct")
		return
	}

	if err = publication.Preparar(); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	repository := publications.PublicationRepository("br")
	publicationCreated, status := repository.CriarPublicacoes(&publication)
	if status.Err != nil {
		response.JSON(w, status.StatusCode, status.Message)
	} else {
		response.JSON(w, status.StatusCode, publicationCreated)
	}

}
func BuscarPublicacoes(w http.ResponseWriter, r *http.Request) {}

func BuscarPublicacao(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	repository := publications.PublicationRepository("br")
	publication, status := repository.BuscarPublicacao(id)

	if status.Err != nil {
		response.Erro(w, status.StatusCode, status.Err)
	} else {
		response.JSON(w, status.StatusCode, publication)
	}
}

func AtualizarPublicacao(w http.ResponseWriter, r *http.Request) {}
func DeletarPublicacao(w http.ResponseWriter, r *http.Request)   {}
