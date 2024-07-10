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
