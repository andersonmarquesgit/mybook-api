package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"mybook-api/src/models"
	"mybook-api/src/repository"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserção de usuário")
	usuarioRequest, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalf("Falha ao ler o corpo da requisição: %v", err)
		createResponse(w, http.StatusBadRequest, "Falha ao ler o corpo da requisição")
		return
	}

	var usuario models.Usuario

	if err = json.Unmarshal(usuarioRequest, &usuario); err != nil {
		createResponse(w, http.StatusBadRequest, "Erro ao converter o usuário para struct")
		return
	}

	_, status := repository.Criar(&usuario)
	if status.Err != nil {
		createResponse(w, status.StatusCode, status.Message)
	} else {
		createResponse(w, status.StatusCode, usuario)
	}

}

func ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	usuarios, status := repository.Listar()
	if status.Err != nil {
		createResponse(w, status.StatusCode, status.Message)
	} else {
		createResponse(w, status.StatusCode, usuarios)
	}
}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	usuario, status := repository.BuscarUsuario(id)

	if status.Err != nil {
		createResponse(w, status.StatusCode, status.Message)
	} else {
		createResponse(w, status.StatusCode, usuario)
	}
}

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	var usuario models.Usuario
	usuario.ID = mux.Vars(r)["id"]
	usuario.AtualizadoEm = time.Now()
	usuarioRequest, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalf("Falha ao ler o corpo da requisição: %v", err)
		createResponse(w, http.StatusBadRequest, "Falha ao ler o corpo da requisição")
		return
	}

	if err = json.Unmarshal(usuarioRequest, &usuario); err != nil {
		createResponse(w, http.StatusBadRequest, "Erro ao converter o usuário para struct")
		return
	}

	_, status := repository.Atualizar(&usuario)
	if status.Err != nil {
		createResponse(w, status.StatusCode, status.Message)
	} else {
		createResponse(w, status.StatusCode, usuario)
	}
}

func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	status := repository.DeletarUsuario(id)

	createResponse(w, status.StatusCode, status.Message)

}

// Função genérica para criar uma resposta HTTP
func createResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if body != nil {
		if err := json.NewEncoder(w).Encode(body); err != nil {
			log.Printf("Erro ao codificar a resposta JSON: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Erro interno do servidor"))
		}
	} else {
		w.Write([]byte("{}"))
	}
}
