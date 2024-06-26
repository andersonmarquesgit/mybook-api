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
	followersCollection *mongo.Collection
	followingCollection *mongo.Collection
}

type RequestStatus struct {
	StatusCode int
	Message    string
	Err        error
}

func FollowersRepository(country string) *Repositorio {
	return &Repositorio{
		followersCollection: banco.GetDB().Collection(country + "-" + config.Collection[os.Getenv("FOLLOWERS_COLLECTION")]),
		followingCollection: banco.GetDB().Collection(country + "-" + config.Collection[os.Getenv("FOLLOWING_COLLECTION")]),
	}
}

func (repositorio Repositorio) FollowUsuario(userID *string, seguidorID *string) (*models.Seguidores, *models.Seguindo, RequestStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Atualizar a coleção de seguidores
	err := repositorio.atualizarSeguidores(ctx, userID, seguidorID)
	if err != nil {
		return nil, nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao atualizar seguidores do usuário no MongoDB", Err: err}
	}

	// Atualizar a coleção de seguidos
	err = repositorio.atualizarSeguindo(ctx, userID, seguidorID)
	if err != nil {
		return nil, nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao atualizar seguidos do usuário no MongoDB", Err: err}
	}

	var followers models.Seguidores
	var following models.Seguindo

	seguidoresCollection := repositorio.followersCollection
	seguindoCollection := repositorio.followingCollection

	filter := bson.M{"userid": userID}
	filterSeguindo := bson.M{"userid": seguidorID}

	err = seguidoresCollection.FindOne(ctx, filter).Decode(&followers)
	if err == mongo.ErrNoDocuments {
		return &followers, &following, RequestStatus{StatusCode: http.StatusNoContent, Message: "Usuário não encontrado", Err: err}
	} else if err != nil {
		log.Fatalf("Erro ao buscar seguidores no MongoDB: %v", err)
		return &followers, &following, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar seguidores no MongoDB", Err: err}
	}

	err = seguindoCollection.FindOne(ctx, filterSeguindo).Decode(&following)
	if err == mongo.ErrNoDocuments {
		return &followers, &following, RequestStatus{StatusCode: http.StatusNoContent, Message: "Seguidor não encontrado", Err: err}
	} else if err != nil {
		log.Fatalf("Erro ao buscar seguidos no MongoDB: %v", err)
		return &followers, &following, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar seguidos no MongoDB", Err: err}
	}

	return &followers, &following, RequestStatus{StatusCode: http.StatusOK}
}

func (repositorio Repositorio) UnfollowUsuario(userID *string, seguidorID *string) (*models.Seguidores, *models.Seguindo, RequestStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Atualizar a coleção de seguidores
	err := repositorio.atualizarSeguidoresUnfollow(ctx, userID, seguidorID)
	if err != nil {
		return nil, nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao atualizar seguidores do usuário no MongoDB", Err: err}
	}

	// Atualizar a coleção de seguidos
	err = repositorio.atualizarSeguindoUnfollow(ctx, userID, seguidorID)
	if err != nil {
		return nil, nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao atualizar seguidos do usuário no MongoDB", Err: err}
	}

	var followers models.Seguidores
	var following models.Seguindo

	seguidoresCollection := repositorio.followersCollection
	seguindoCollection := repositorio.followingCollection // Suponha que você tenha uma referência para a coleção de seguidos

	filter := bson.M{"userid": userID}
	filterSeguindo := bson.M{"userid": seguidorID}

	err = seguidoresCollection.FindOne(ctx, filter).Decode(&followers)
	if err == mongo.ErrNoDocuments {
		return &followers, &following, RequestStatus{StatusCode: http.StatusNoContent, Message: "Usuário não encontrado", Err: err}
	} else if err != nil {
		log.Fatalf("Erro ao buscar seguidores no MongoDB: %v", err)
		return &followers, &following, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar seguidores no MongoDB", Err: err}
	}

	err = seguindoCollection.FindOne(ctx, filterSeguindo).Decode(&following)
	if err == mongo.ErrNoDocuments {
		return &followers, &following, RequestStatus{StatusCode: http.StatusNoContent, Message: "Seguidor não encontrado", Err: err}
	} else if err != nil {
		log.Fatalf("Erro ao buscar seguidos no MongoDB: %v", err)
		return &followers, &following, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar seguidos no MongoDB", Err: err}
	}

	return &followers, &following, RequestStatus{StatusCode: http.StatusOK}
}

func (repositorio Repositorio) FindFollowers(id string) (models.Seguidores, RequestStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var usuario models.Seguidores
	err := repositorio.followersCollection.FindOne(ctx, bson.M{"userid": id}).Decode(&usuario)
	if err == mongo.ErrNoDocuments {
		return usuario, RequestStatus{StatusCode: http.StatusNoContent, Message: "Usuário não encontrado", Err: err}

	} else if err != nil {
		log.Fatalf("Erro ao buscar usuário no MongoDB: %v", err)
		return usuario, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar usuário no MongoDB", Err: err}
	}

	return usuario, RequestStatus{StatusCode: http.StatusOK}
}

func (repositorio Repositorio) FindFollowing(id string) (models.Seguindo, RequestStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var usuario models.Seguindo
	err := repositorio.followingCollection.FindOne(ctx, bson.M{"userid": id}).Decode(&usuario)
	if err == mongo.ErrNoDocuments {
		return usuario, RequestStatus{StatusCode: http.StatusNoContent, Message: "Usuário não encontrado", Err: err}

	} else if err != nil {
		log.Fatalf("Erro ao buscar usuário no MongoDB: %v", err)
		return usuario, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar usuário no MongoDB", Err: err}
	}

	return usuario, RequestStatus{StatusCode: http.StatusOK}
}

func (repositorio Repositorio) atualizarSeguidores(ctx context.Context, userID *string, seguidorID *string) error {
	collection := repositorio.followersCollection

	filter := bson.M{"userid": userID}
	update := bson.M{
		"$addToSet": bson.M{
			"followers": seguidorID,
		},
		"$set": bson.M{
			"atualizadoem": time.Now(),
		},
	}

	options := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, options)
	return err
}

func (repositorio Repositorio) atualizarSeguindo(ctx context.Context, userID *string, seguindoID *string) error {
	collection := repositorio.followingCollection

	filter := bson.M{"userid": seguindoID}
	update := bson.M{
		"$addToSet": bson.M{
			"following": userID,
		},
		"$set": bson.M{
			"atualizadoem": time.Now(),
		},
	}

	options := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, options)
	return err
}

func (repositorio Repositorio) atualizarSeguindoUnfollow(ctx context.Context, userID *string, seguindoID *string) error {
	collection := repositorio.followingCollection

	filter := bson.M{"userid": seguindoID}
	update := bson.M{
		"$pull": bson.M{
			"following": userID,
		},
		"$set": bson.M{
			"atualizadoem": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (repositorio Repositorio) atualizarSeguidoresUnfollow(ctx context.Context, userID *string, seguidorID *string) error {
	collection := repositorio.followersCollection

	filter := bson.M{"userid": userID}
	update := bson.M{
		"$pull": bson.M{
			"followers": seguidorID,
		},
		"$set": bson.M{
			"atualizadoem": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
