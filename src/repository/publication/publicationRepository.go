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

func (repositorio Repositorio) BuscarPublicacoes() ([]models.Publication, RequestStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := repositorio.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Erro ao buscar publicações no MongoDB: %v", err)
		return nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar publicações no MongoDB", Err: err}
	}

	defer cursor.Close(ctx)

	var publications []models.Publication
	for cursor.Next(ctx) {
		var publication models.Publication
		if err := cursor.Decode(&publication); err != nil {
			log.Fatalf("Erro ao decodificar usuário: %v", err)
			return nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao decodificar publicação", Err: err}
		}
		publications = append(publications, publication)
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("Erro durante a iteração do cursor: %v", err)
		return nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro durante a iteração do cursor", Err: err}
	}

	if publications == nil {
		log.Println("Nenhuma publicação encontrada")
		return publications, RequestStatus{StatusCode: http.StatusNoContent, Message: "Nenhuma publicação encontrada"}
	}

	return publications, RequestStatus{StatusCode: http.StatusOK}
}
