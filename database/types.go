package database

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rs/zerolog"
)

type EntityType = string

type EntityRepository struct {
	svc            dynamodb.Client
	entityTableMap map[EntityType]string
	ctx            context.Context
	logger         zerolog.Logger
}

type IEntityRepository interface {
	UpsertEntity(entity Entity) error
	GetEntity(id string, entityType EntityType) (*Entity, error)

	UpsertObject(id string, data any, entityType EntityType) error
	GetObject(id string, entityType EntityType, result any) error
}

type Entity struct {
	Id   string     `dynamodbav:"id"`
	Type EntityType `dynamodbav:"type"`
	Data string     `dynamodbav:"data"`
}
