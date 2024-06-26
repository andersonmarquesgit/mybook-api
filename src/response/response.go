package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON retorna uma resposta JSON para a requisição
func JSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if body != nil {
		if err := json.NewEncoder(w).Encode(body); err != nil {
			log.Printf("Erro ao codificar a resposta JSON: %v", err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(statusCode)
		w.Write([]byte("{}"))
	}
}

func Erro(w http.ResponseWriter, statusCode int, err error) {
	JSON(w, statusCode, struct {
		Erro string `json:"erro"`
	}{
		Erro: err.Error(),
	})
}

// Sucesso retorna uma resposta JSON com mensagem de sucesso
func Sucesso(w http.ResponseWriter, statusCode int, message string) {
	JSON(w, statusCode, struct {
		Mensagem string `json:"mensagem"`
	}{
		Mensagem: message,
	})
}
