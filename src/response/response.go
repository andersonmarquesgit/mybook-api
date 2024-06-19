package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON retornar uma resposta Json para a requisição
func JSON(w http.ResponseWriter, statusCode int, body interface{}) {
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

func Erro(w http.ResponseWriter, statusCode int, err error) {
	JSON(w, statusCode, struct {
		Erro string `json:"erro"`
	}{
		Erro: err.Error(),
	})
}
