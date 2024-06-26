package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"mybook-api/src/infrastructure/autenticacao"
	"mybook-api/src/infrastructure/security"
	"mybook-api/src/models"
	"mybook-api/src/presentation"
	"mybook-api/src/repository/followers"
	"mybook-api/src/repository/users"
	"mybook-api/src/response"
	"net/http"

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

	repository := users.UsersRepository("br")
	_, status := repository.Criar(&usuario)
	if status.Err != nil {
		response.JSON(w, status.StatusCode, status.Message)
	} else {
		response.JSON(w, status.StatusCode, usuario)
	}

}

func ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	repository := users.UsersRepository("br")
	usuarios, status := repository.Listar()
	if status.Err != nil {
		response.Erro(w, status.StatusCode, status.Err)
	} else {
		response.JSON(w, status.StatusCode, usuarios)
	}
}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	repository := users.UsersRepository("br")
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

	repository := users.UsersRepository("br")
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

	repository := users.UsersRepository("br")

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

	userRepository := users.UsersRepository("br")
	if _, status := userRepository.BuscarUsuario(id); status.Err != nil {
		response.Erro(w, http.StatusBadRequest, status.Err)
		return
	}

	if _, status := userRepository.BuscarUsuario(seguidorID); status.Err != nil {
		response.JSON(w, http.StatusBadRequest, status.Err)
		return
	}

	followersRepository := followers.FollowersRepository("br")
	followers, _, status := followersRepository.FollowUsuario(&id, &seguidorID)

	if status.Err != nil {
		response.Erro(w, status.StatusCode, status.Err)
	} else {
		response.JSON(w, status.StatusCode, presentation.NewFollowersResponse(*followers))
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

	userRepository := users.UsersRepository("br")
	if _, status := userRepository.BuscarUsuario(id); status.Err != nil {
		response.Erro(w, http.StatusBadRequest, status.Err)
		return
	}

	if _, status := userRepository.BuscarUsuario(seguidorID); status.Err != nil {
		response.JSON(w, http.StatusBadRequest, status.Err)
		return
	}

	followersRepository := followers.FollowersRepository("br")
	followers, _, status := followersRepository.UnfollowUsuario(&id, &seguidorID)

	if status.Err != nil {
		response.Erro(w, status.StatusCode, status.Err)
	} else {
		response.JSON(w, status.StatusCode, presentation.NewFollowersResponse(*followers))
	}
}

func Followers(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	followersRepository := followers.FollowersRepository("br")
	followers, status := followersRepository.FindFollowers(id)

	if status.Err != nil {
		response.Erro(w, status.StatusCode, status.Err)
	} else {
		response.JSON(w, status.StatusCode, presentation.NewFollowersResponse(followers))
	}
}

func Following(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	followersRepository := followers.FollowersRepository("br")
	following, status := followersRepository.FindFollowing(id)

	if status.Err != nil {
		response.Erro(w, status.StatusCode, status.Err)
	} else {
		response.JSON(w, status.StatusCode, presentation.NewFollowingResponse(following))
	}
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	usuarioIDNoToken, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	userID := mux.Vars(r)["id"]

	if userID != usuarioIDNoToken {
		response.Erro(w, http.StatusForbidden, errors.New("Não é possível atualizar a senha de outro usuário"))
		return
	}

	var updatePassword models.Senha
	updatePasswordRequest, err := ioutil.ReadAll(r.Body)
	if err = json.Unmarshal(updatePasswordRequest, &updatePassword); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	userRepository := users.UsersRepository("br")
	senhaSalvaNoBanco, status := userRepository.FindPassword(userID)

	if status.Err != nil {
		response.Erro(w, status.StatusCode, status.Err)
		return
	}

	if err = security.VerificarSenha(senhaSalvaNoBanco, updatePassword.Atual); err != nil {
		response.Erro(w, http.StatusUnauthorized, errors.New("A senha atual não condiz com a senha salva no banco"))
		return
	}

	senhaComHash, err := security.Hash(updatePassword.Nova)
	if err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	if status = userRepository.UpdatePassword(userID, string(senhaComHash)); err != nil {
		response.Erro(w, status.StatusCode, status.Err)
		return
	}

	response.Sucesso(w, status.StatusCode, status.Message)
}
