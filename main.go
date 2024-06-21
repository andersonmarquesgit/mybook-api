package main

import (
	"fmt"
	"log"
	"mybook-api/src/infrastructure/banco"
	"mybook-api/src/infrastructure/config"
	"mybook-api/src/infrastructure/router"
	"net/http"
)

// Usado apenas para gerar a secret key que uso no .env
// func init() {
// 	chave := make([]byte, 64)

// 	if _, erro := rand.Read(chave); erro != nil {
// 		log.Fatal(erro)
// 	}

// 	stringBase64 := base64.StdEncoding.EncodeToString(chave)
// 	fmt.Println(stringBase64)
// }

func main() {
	config.Carregar()
	banco.ConectarMongoDB()
	defer banco.DesconectarMongoDB()

	log.Printf("Escutando na porta %d", config.Port)

	r := router.Gerar()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))

}
