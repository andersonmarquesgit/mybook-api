package followers

import (
	"context"
	"log"
	"mybook-api/src/infrastructure/banco"
	"mybook-api/src/infrastructure/config"
	"mybook-api/src/models"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repositorio struct {
	collection *mongo.Collection
}

// var Respository Repositorio

type RequestStatus struct {
	StatusCode int
	Message    string
	Err        error
}

// func init() {
// 	Respository = Repositorio{collection: banco.GetDB().Collection("br" + "-" + config.Collection[os.Getenv("FOLLOWERS_COLLECTION")])}
// }

func FollowersRepository(country string) *Repositorio {
	return &Repositorio{banco.GetDB().Collection(country + "-" + config.Collection[os.Getenv("FOLLOWERS_COLLECTION")])}
}

func (repositorio Repositorio) SeguirUsuario(userID *string, seguidorID *string) (*models.Seguidores, RequestStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"userid": userID}

	update := bson.M{
		"$addToSet": bson.M{
			"seguidores": seguidorID,
		},
		"$set": bson.M{
			"atualizadoem": time.Now(),
		},
	}

	options := options.Update().SetUpsert(true)
	result, err := repositorio.collection.UpdateOne(ctx, filter, update, options)
	if err != nil {
		return nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao atualizar seguidores do usuário no MongoDB", Err: err}
	}

	if result.MatchedCount == 0 {
		log.Printf("Nenhum documento de seguidores encontrado, criando um novo")
	} else {
		log.Printf("Seguidores do usuário %v atualizados com sucesso", *userID)
	}

	var followers models.Seguidores

	err = repositorio.collection.FindOne(ctx, filter).Decode(&followers)
	if err == mongo.ErrNoDocuments {
		return &followers, RequestStatus{StatusCode: http.StatusNoContent, Message: "Usuário não encontrado", Err: err}

	} else if err != nil {
		log.Fatalf("Erro ao buscar usuário no MongoDB: %v", err)
		return &followers, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar usuário no MongoDB", Err: err}

	}

	return &followers, RequestStatus{StatusCode: http.StatusOK}
}

// TODO refatorar repositorio de seguidores
func (repositorio Repositorio) UnfollowUsuario(userID *string, seguidorID *string) (*models.Seguidores, RequestStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"userid": userID}

	update := bson.M{
		"$pull": bson.M{
			"seguidores": seguidorID,
		},
		"$set": bson.M{
			"atualizadoem": time.Now(),
		},
	}

	result, err := repositorio.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao atualizar seguidores do usuário no MongoDB", Err: err}
	}

	if result.MatchedCount == 0 {
		log.Printf("Nenhum documento de seguidores encontrado para o usuário %v", *userID)
	} else {
		log.Printf("Seguidores do usuário %v atualizados com sucesso", *userID)
	}

	var followers models.Seguidores

	err = repositorio.collection.FindOne(ctx, filter).Decode(&followers)
	if err == mongo.ErrNoDocuments {
		return &followers, RequestStatus{StatusCode: http.StatusNoContent, Message: "Usuário não encontrado", Err: err}

	} else if err != nil {
		log.Fatalf("Erro ao buscar usuário no MongoDB: %v", err)
		return &followers, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar usuário no MongoDB", Err: err}
	}

	return &followers, RequestStatus{StatusCode: http.StatusOK}
}

func (repositorio Repositorio) FindFollowers(id string) (models.Seguidores, RequestStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var usuario models.Seguidores
	err := repositorio.collection.FindOne(ctx, bson.M{"userid": id}).Decode(&usuario)
	if err == mongo.ErrNoDocuments {
		return usuario, RequestStatus{StatusCode: http.StatusNoContent, Message: "Usuário não encontrado", Err: err}

	} else if err != nil {
		log.Fatalf("Erro ao buscar usuário no MongoDB: %v", err)
		return usuario, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar usuário no MongoDB", Err: err}
	}

	return usuario, RequestStatus{StatusCode: http.StatusOK}
}
