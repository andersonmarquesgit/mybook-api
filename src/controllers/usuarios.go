package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"mybook-api/src/infrastructure/autenticacao"
	"mybook-api/src/models"
	"mybook-api/src/repository/followers"
	repository "mybook-api/src/repository/users"
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

	repository := repository.UsersRepository("br")
	_, status := repository.Criar(&usuario)
	if status.Err != nil {
		response.JSON(w, status.StatusCode, status.Message)
	} else {
		response.JSON(w, status.StatusCode, usuario)
	}

}

func ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	repository := repository.UsersRepository("br")
	usuarios, status := repository.Listar()
	if status.Err != nil {
		response.Erro(w, status.StatusCode, status.Err)
	} else {
		response.JSON(w, status.StatusCode, usuarios)
	}
}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	repository := repository.UsersRepository("br")
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

	repository := repository.UsersRepository("br")
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

	repository := repository.UsersRepository("br")

	status := repository.DeletarUsuario(id)

	response.JSON(w, status.StatusCode, status.Message)

}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	seguidorID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
	}

	id := mux.Vars(r)["id"]

	if id == seguidorID {
		response.Erro(w, http.StatusForbidden, errors.New("Não é possível seguir você mesmo"))
		return
	}

	userRepository := repository.UsersRepository("br")
	if _, status := userRepository.BuscarUsuario(id); status.Err != nil {
		response.Erro(w, http.StatusBadRequest, status.Err)
		return
	}

	if _, status := userRepository.BuscarUsuario(seguidorID); status.Err != nil {
		response.JSON(w, http.StatusBadRequest, status.Err)
		return
	}

	followersRepository := followers.FollowersRepository("br")
	followers, status := followersRepository.SeguirUsuario(&id, &seguidorID)

	if status.Err != nil {
		response.Erro(w, status.StatusCode, status.Err)
	} else {
		response.JSON(w, status.StatusCode, *followers)
	}
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	seguidorID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
	}

	id := mux.Vars(r)["id"]

	if id == seguidorID {
		response.Erro(w, http.StatusForbidden, errors.New("Não é possível deixar de seguir você mesmo"))
		return
	}

	userRepository := repository.UsersRepository("br")
	if _, status := userRepository.BuscarUsuario(id); status.Err != nil {
		response.Erro(w, http.StatusBadRequest, status.Err)
		return
	}

	if _, status := userRepository.BuscarUsuario(seguidorID); status.Err != nil {
		response.JSON(w, http.StatusBadRequest, status.Err)
		return
	}

	followersRepository := followers.FollowersRepository("br")
	followers, status := followersRepository.UnfollowUsuario(&id, &seguidorID)

	if status.Err != nil {
		response.Erro(w, status.StatusCode, status.Err)
	} else {
		response.JSON(w, status.StatusCode, *followers)
	}
}

func Followers(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	followersRepository := followers.FollowersRepository("br")
	usuario, status := followersRepository.FindFollowers(id)

	if status.Err != nil {
		response.Erro(w, status.StatusCode, status.Err)
	} else {
		response.JSON(w, status.StatusCode, usuario)
	}
}
