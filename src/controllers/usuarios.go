package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"mybook-api/src/infrastructure/autenticacao"
	"mybook-api/src/models"
	"mybook-api/src/repository"
	"mybook-api/src/response"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserção de usuário")
	usuarioRequest, err := ioutil.ReadAll(r.Body)
	//country := mux.Vars(r)["country"]

	if err != nil {
		log.Fatalf("Falha ao ler o corpo da requisição: %v", err)
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	var usuario models.Usuario

	if err = json.Unmarshal(usuarioRequest, &usuario); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	if err = usuario.Preparar(); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	repository := repository.NovoRepositorio("br")
	_, status := repository.Criar(&usuario)
	if status.Err != nil {
		response.JSON(w, status.StatusCode, status.Message)
	} else {
		response.JSON(w, status.StatusCode, usuario)
	}

}

func ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	repository := repository.NovoRepositorio("br")
	usuarios, status := repository.Listar()
	if status.Err != nil {
		response.Erro(w, status.StatusCode, status.Err)
	} else {
		response.JSON(w, status.StatusCode, usuarios)
	}
}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	repository := repository.NovoRepositorio("br")
	usuario, status := repository.BuscarUsuario(id)

	if status.Err != nil {
		response.Erro(w, status.StatusCode, status.Err)
	} else {
		response.JSON(w, status.StatusCode, usuario)
	}
}

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	var usuario models.Usuario
	usuario.ID = mux.Vars(r)["id"]
	usuario.AtualizadoEm = time.Now()
	usuarioRequest, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalf("Falha ao ler o corpo da requisição: %v", err)
		response.JSON(w, http.StatusBadRequest, "Falha ao ler o corpo da requisição")
		return
	}

	usuarioIDNoToken, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
	}

	if usuario.ID != usuarioIDNoToken {
		response.Erro(w, http.StatusForbidden, errors.New("Não é possível atualizar um usuário diferente do seu"))
		return
	}

	if err = json.Unmarshal(usuarioRequest, &usuario); err != nil {
		response.JSON(w, http.StatusBadRequest, "Erro ao converter o usuário para struct")
		return
	}

	repository := repository.NovoRepositorio("br")
	_, status := repository.Atualizar(&usuario)
	if status.Err != nil {
		response.JSON(w, status.StatusCode, status.Message)
	} else {
		response.JSON(w, status.StatusCode, usuario)
	}
}

func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	usuarioIDNoToken, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
	}

	if id != usuarioIDNoToken {
		response.Erro(w, http.StatusForbidden, errors.New("Não é possível excluir um usuário diferente do seu"))
		return
	}

	repository := repository.NovoRepositorio("br")

	status := repository.DeletarUsuario(id)

	response.JSON(w, status.StatusCode, status.Message)

}
