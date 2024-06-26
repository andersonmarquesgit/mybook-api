package users

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

func UsersRepository(country string) *Repositorio {
	return &Repositorio{banco.GetDB().Collection(country + "-" + config.Collection[os.Getenv("USER_COLLECTION")])}
}

func (repositorio Repositorio) Criar(usuario *models.Usuario) (*models.Usuario, RequestStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usuario.ID = uuid.New().String()
	usuario.CriadoEm = time.Now()

	_, err := repositorio.collection.InsertOne(ctx, usuario)
	if err != nil {
		log.Fatalf("Erro ao inserir usuário no MongoDB: %v", err)
		return nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao inserir usuário no MongoDB", Err: err}
	}

	// Busca o documento atualizado
	err = repositorio.collection.FindOne(ctx, bson.M{"id": usuario.ID}).Decode(&usuario)
	if err == mongo.ErrNoDocuments {
		log.Fatalf("Erro ao buscar usuário inserido no MongoDB: %v", err)
		return nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar usuário inserido no MongoDB", Err: err}
	} else if err != nil {
		log.Fatalf("Erro ao buscar usuário no MongoDB: %v", err)
		return nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscarr usuário no MongoDB", Err: err}
	}

	log.Println("Usuário inserido com sucesso!")
	return usuario, RequestStatus{StatusCode: http.StatusCreated, Message: "Usuário inserido com sucesso!"}
}

func (repositorio Repositorio) Listar() ([]models.Usuario, RequestStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := repositorio.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Erro ao buscar usuários no MongoDB: %v", err)
		return nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar usuários no MongoDB", Err: err}
	}

	defer cursor.Close(ctx)

	var usuarios []models.Usuario
	for cursor.Next(ctx) {
		var usuario models.Usuario
		if err := cursor.Decode(&usuario); err != nil {
			log.Fatalf("Erro ao decodificar usuário: %v", err)
			return nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao decodificar usuário", Err: err}
		}
		usuarios = append(usuarios, usuario)
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("Erro durante a iteração do cursor: %v", err)
		return nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro durante a iteração do cursor", Err: err}
	}

	if usuarios == nil {
		log.Println("Nenhum usuário encontrado")
		return usuarios, RequestStatus{StatusCode: http.StatusNoContent, Message: "Nenhum usuário encontrado"}
	}

	return usuarios, RequestStatus{StatusCode: http.StatusOK}
}

func (repositorio Repositorio) BuscarUsuario(id string) (models.Usuario, RequestStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var usuario models.Usuario
	err := repositorio.collection.FindOne(ctx, bson.M{"id": id}).Decode(&usuario)
	if err == mongo.ErrNoDocuments {
		return usuario, RequestStatus{StatusCode: http.StatusNoContent, Message: "Usuário não encontrado", Err: err}

	} else if err != nil {
		log.Fatalf("Erro ao buscar usuário no MongoDB: %v", err)
		return usuario, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar usuário no MongoDB", Err: err}
	}

	return usuario, RequestStatus{StatusCode: http.StatusOK}
}

func (repositorio Repositorio) BuscarUsuarioPorEmail(email string) (models.Usuario, RequestStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var usuario models.Usuario
	err := repositorio.collection.FindOne(ctx, bson.M{"email": email}).Decode(&usuario)
	if err == mongo.ErrNoDocuments {
		return usuario, RequestStatus{StatusCode: http.StatusNoContent, Message: "Usuário não encontrado", Err: err}

	} else if err != nil {
		log.Fatalf("Erro ao buscar usuário no MongoDB: %v", err)
		return usuario, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar usuário no MongoDB", Err: err}
	}

	return usuario, RequestStatus{StatusCode: http.StatusOK}
}

func (repositorio Repositorio) Atualizar(usuario *models.Usuario) (*models.Usuario, RequestStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": usuario.ID}

	update := bson.M{
		"$set": bson.M{
			"nome":         usuario.Nome,
			"email":        usuario.Email,
			"nick":         usuario.Nick,
			"atualizadoem": time.Now(),
		},
	}

	_, err := repositorio.collection.UpdateOne(ctx, filter, update)
	if err == mongo.ErrNoDocuments {
		return usuario, RequestStatus{StatusCode: http.StatusNoContent, Message: "Usuário não encontrado", Err: err}
	} else if err != nil {
		log.Fatalf("Erro ao atualizar usuário no MongoDB: %v", err)
		return usuario, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao atualizar usuário no MongoDB", Err: err}
	}

	err = repositorio.collection.FindOne(ctx, filter).Decode(&usuario)
	if err == mongo.ErrNoDocuments {
		return usuario, RequestStatus{StatusCode: http.StatusNoContent, Message: "Usuário não encontrado", Err: err}

	} else if err != nil {
		log.Fatalf("Erro ao buscar usuário no MongoDB: %v", err)
		return usuario, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar usuário no MongoDB", Err: err}

	}

	return usuario, RequestStatus{StatusCode: http.StatusOK}
}

func (repositorio Repositorio) DeletarUsuario(id string) RequestStatus {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := repositorio.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		log.Fatalf("Erro ao deletar usuário no MongoDB: %v", err)
		return RequestStatus{StatusCode: http.StatusNoContent, Message: "Usuário não encontrado", Err: err}
	}

	if result.DeletedCount == 0 {
		return RequestStatus{StatusCode: http.StatusNoContent, Message: "Usuário não encontrado", Err: err}
	}

	return RequestStatus{StatusCode: http.StatusOK, Message: "Usuário deletado com sucesso!"}
}

func (repositorio Repositorio) FindPassword(userID string) (string, RequestStatus) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result struct {
		Senha string `bson:"senha"`
	}

	filter := bson.M{"id": userID}
	err := repositorio.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", RequestStatus{StatusCode: http.StatusNoContent, Message: "Usuário não encontrado", Err: err}
		}
		return "", RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao buscar usuário no MongoDB", Err: err}
	}

	if result.Senha == "" {
		return "", RequestStatus{StatusCode: http.StatusNotFound, Message: "Usuário não encontrado"}
	}

	return result.Senha, RequestStatus{StatusCode: http.StatusOK}
}

func (repositorio Repositorio) UpdatePassword(userID string, novaSenha string) RequestStatus {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": userID}
	update := bson.M{
		"$set": bson.M{
			"senha":        string(novaSenha),
			"atualizadoem": time.Now(),
		},
	}

	result, err := repositorio.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return RequestStatus{StatusCode: http.StatusNotFound, Message: "Usuário não encontrado", Err: err}
		}
		return RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao atualizar senha do usuário no MongoDB", Err: err}
	}

	if result.MatchedCount == 0 {
		return RequestStatus{StatusCode: http.StatusNotFound, Message: "Usuário não encontrado"}
	}

	return RequestStatus{StatusCode: http.StatusOK, Message: "Senha atualizada com sucesso"}
}
