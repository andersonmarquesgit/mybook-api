package repository

import (
	"context"
	"log"
	"mybook-api/src/infrastructure/banco"
	"mybook-api/src/models"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RequestStatus struct {
	StatusCode int
	Message    string
	Err        error
}

func Criar(usuario *models.Usuario) (*models.Usuario, RequestStatus) {
	collection := banco.GetDB().Collection("usuarios")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usuario.ID = uuid.New().String()
	usuario.CriadoEm = time.Now()

	_, err := collection.InsertOne(ctx, usuario)
	if err != nil {
		log.Fatalf("Erro ao inserir usuário no MongoDB: %v", err)
		return nil, RequestStatus{StatusCode: http.StatusInternalServerError, Message: "Erro ao inserir usuário no MongoDB", Err: err}
	}

	// Busca o documento atualizado
	err = collection.FindOne(ctx, bson.M{"id": usuario.ID}).Decode(&usuario)
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

func Listar() ([]models.Usuario, RequestStatus) {
	collection := banco.GetDB().Collection("usuarios")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
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
