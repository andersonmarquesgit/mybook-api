package banco

import (
	"context"
	"log"
	"mybook-api/src/infrastructure/config"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConectarMongoDB() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.DBStringConnection))
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}

	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Erro ao pingar o MongoDB: %v", err)
	}

	log.Println("Conectado ao MongoDB!")
}

func DesconectarMongoDB() {
	if Client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Disconnect(ctx); err != nil {
		log.Fatalf("Erro ao desconectar do MongoDB: %v", err)
	}

	log.Println("Desconectado do MongoDB!")
}

func GetDB() *mongo.Database {
	if Client == nil {
		log.Fatal("Cliente MongoDB não inicializado")
	}
	return Client.Database(os.Getenv("DB_NAME"))
}
