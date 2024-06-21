package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"mybook-api/src/infrastructure/autenticacao"
	"mybook-api/src/infrastructure/security"
	"mybook-api/src/models"
	"mybook-api/src/repository"
	"mybook-api/src/response"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Login de usuário")
	userLoginRequest, err := ioutil.ReadAll(r.Body)

	if err != nil {
		response.Erro(w, http.StatusUnprocessableEntity, err)
	}

	var usuarioLogin models.Usuario
	if err = json.Unmarshal(userLoginRequest, &usuarioLogin); err != nil {
		response.Erro(w, http.StatusBadRequest, err)
		return
	}

	repository := repository.NovoRepositorio("br")
	usuarioDoBanco, status := repository.BuscarUsuarioPorEmail(usuarioLogin.Email)

	if status.Err != nil {
		response.Erro(w, status.StatusCode, status.Err)
	}

	if err = security.VerificarSenha(usuarioDoBanco.Senha, usuarioLogin.Senha); err != nil {
		response.Erro(w, http.StatusUnauthorized, err)
		return
	}

	token, _ := autenticacao.CriarToken(usuarioDoBanco.ID)
	response.JSON(w, http.StatusCreated, token)

}
