package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DBStringConnection = ""
	Port               = 0
	SecretKey          []byte
	Collection         string
)

// Carregar vai inicializar as variáveis de ambiente
func Carregar() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	Port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		Port = 9000
	}

	DBStringConnection = fmt.Sprintf("%s://%s:%s",
		os.Getenv("DB_URI"),
		os.Getenv("DB_SERVER"),
		os.Getenv("DB_PORT"))

	SecretKey = []byte(os.Getenv("SECRET_KEY"))

	Collection = os.Getenv("USER_COLLECTION")

}
