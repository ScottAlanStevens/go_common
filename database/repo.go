package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/smithy-go"
	"github.com/rs/zerolog"
)

const (
	ENTITY_ID   = "id"
	ENTITY_TYPE = "type"
)

func NewEntityRepository(ctx context.Context, logger zerolog.Logger, entityTableMap map[EntityType]string) (IEntityRepository, error) {
	config, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	svc := dynamodb.NewFromConfig(config)

	return &EntityRepository{
		svc:            *svc,
		entityTableMap: entityTableMap,
		ctx:            ctx,
		logger:         logger.With().Str("service", "EntityRepository").Logger(),
	}, nil
}

func (c *EntityRepository) UpsertEntity(entity Entity) error {

	attributesMap, err := attributevalue.MarshalMap(entity)
	if err != nil {
		return fmt.Errorf("failed to marshal record, %w", err)
	}

	_, err = c.svc.PutItem(c.ctx, &dynamodb.PutItemInput{
		TableName: aws.String(c.entityTableMap[entity.Type]),
		Item:      attributesMap,
	})

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case "ConditionalCheckFailedException":
				return apiErr
			case "ResourceNotFoundException": // Table not found
				return apiErr
			default:
				c.logger.Error().Msgf("%s: %s", apiErr.ErrorCode(), apiErr.ErrorMessage())
			}
		} else {
			return err
		}
	}

	return err
}

func (c *EntityRepository) GetEntity(id string, entityType EntityType) (*Entity, error) {

	key := map[string]types.AttributeValue{
		ENTITY_ID: &types.AttributeValueMemberS{Value: id},
		// ENTITY_TYPE: &types.AttributeValueMemberS{Value: entityType},
	}

	tableName := c.entityTableMap[entityType]
	result, err := c.svc.GetItem(c.ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve record, %w", err)
	}

	if result == nil || result.Item == nil {
		return nil, nil
	}

	entity := Entity{}

	err = attributevalue.UnmarshalMap(result.Item, &entity)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal record, %w", err)
	}

	return &entity, nil
}

func (c *EntityRepository) UpsertObject(id string, data any, entityType EntityType) error {
	entity, err := MarshalObject(data, id, entityType)
	if err != nil {
		return err
	}

	return c.UpsertEntity(*entity)
}

func (c *EntityRepository) GetObject(id string, entityType EntityType, result any) error {

	entity, err := c.GetEntity(id, entityType)
	if err != nil {
		return err
	}

	if entity == nil {
		return nil
	}

	err = UnmarshalObject(*entity, &result)
	return err
}
