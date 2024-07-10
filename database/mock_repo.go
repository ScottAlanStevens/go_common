package database

import (
	"github.com/stretchr/testify/mock"
)

type MockEntityRepository struct {
	mock.Mock
}

func (mock *MockEntityRepository) UpsertEntity(entity Entity) error {
	args := mock.Called(entity)
	return args.Error(0)
}

func (mock *MockEntityRepository) GetEntity(id string, entityType EntityType) (*Entity, error) {
	args := mock.Called(id, entityType)
	a := args.Get(0)
	if a == nil {
		return nil, args.Error(1)
	}
	return a.(*Entity), args.Error(1)
}

func (mock *MockEntityRepository) UpsertObject(id string, data any, entityType EntityType) error {
	args := mock.Called(id, data, entityType)
	return args.Error(0)
}
func (mock *MockEntityRepository) GetObject(id string, entityType EntityType, result any) error {
	args := mock.Called(id, entityType, result)
	return args.Error(0)
}
