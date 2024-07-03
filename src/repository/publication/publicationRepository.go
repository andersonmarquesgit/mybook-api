package publications

import (
	"context"
	"log"
	"mybook-api/src/infrastructure/banco"
	"mybook-api/src/infrastructure/config"
	"mybook-api/src/models"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositorio struct {
	collection *mongo.Collection
}

type RequestStatus struct {
	StatusCode int
	Message    string
	Err        error
}

func PublicationRepository(country string) *Repositorio {
	return &Repositorio{banco.GetDB().Collection(country + "-" + config.Collection[os.Getenv("PUBLICATIONS_COLLECTION")])}
}

func (repositorio Repositorio) CriarPublicacoes(publication *models.Publication) (*models.Publication, RequestStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	publication.ID = uuid.New().String()
	publication.CreatedIn = time.Now()

	_, err := repositorio.collection.InsertOne(ctx, publication)
	if err != nil {
		log.Fatalf("Erro ao inserir publicação no MongoDB: %v", err)
		return nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao inserir publicação no MongoDB", Err: err}
	}

	// Busca o documento atualizado
	err = repositorio.collection.FindOne(ctx, bson.M{"id": publication.ID}).Decode(&publication)
	if err == mongo.ErrNoDocuments {
		log.Fatalf("Erro ao buscar publicação inserida no MongoDB: %v", err)
		return nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar publicação inserido no MongoDB", Err: err}
	} else if err != nil {
		log.Fatalf("Erro ao buscar publicação no MongoDB: %v", err)
		return nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar publicação no MongoDB", Err: err}
	}

	log.Println("Publicação inserida com sucesso!")
	return publication, RequestStatus{StatusCode: http.StatusCreated, Message: "Publicação inserida com sucesso!"}
}

func (repositorio Repositorio) BuscarPublicacao(id string) (models.Publication, RequestStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var publication models.Publication

	err := repositorio.collection.FindOne(ctx, bson.M{"id": id}).Decode(&publication)
	if err == mongo.ErrNoDocuments {
		return publication, RequestStatus{StatusCode: http.StatusNoContent, Message: "Publicação não encontrada", Err: err}

	} else if err != nil {
		log.Fatalf("Erro ao buscar publicação no MongoDB: %v", err)
		return publication, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar publicação no MongoDB", Err: err}
	}

	return publication, RequestStatus{StatusCode: http.StatusOK}

}
