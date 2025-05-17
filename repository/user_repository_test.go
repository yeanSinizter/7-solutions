package repository_test

import (
	"7-solutions/model"
	"7-solutions/repository"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter)
	raw := args.Get(0).(bson.Raw)
	sr := mongo.NewSingleResultFromDocument(raw, nil, nil)
	return sr
}

func (m *MockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (m *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (m *MockCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func TestUserRepository_Create(t *testing.T) {
	mockColl := new(MockCollection)
	repo := repository.NewUserRepositoryFromCollection(mockColl)

	user := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	mockColl.On("InsertOne", mock.Anything, mock.AnythingOfType("*model.User")).Return(&mongo.InsertOneResult{}, nil)

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
	mockColl.AssertExpectations(t)
}

func TestUserRepository_GetByID(t *testing.T) {
	mockColl := new(MockCollection)
	repo := repository.NewUserRepositoryFromCollection(mockColl)

	user := model.User{
		ID:    primitive.NewObjectID(),
		Name:  "Test User",
		Email: "test@example.com",
	}

	raw, _ := bson.Marshal(user)

	mockColl.On("FindOne", mock.Anything, mock.Anything).Return(bson.Raw(raw))

	result, err := repo.GetByID(context.Background(), user.ID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, user.Email, result.Email)
	mockColl.AssertExpectations(t)
}

func TestUserRepository_Update(t *testing.T) {
	mockColl := new(MockCollection)
	repo := repository.NewUserRepositoryFromCollection(mockColl)

	userID := primitive.NewObjectID()
	user := &model.User{
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	mockColl.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).
		Return(&mongo.UpdateResult{MatchedCount: 1}, nil)

	err := repo.Update(context.Background(), userID.Hex(), user)
	assert.NoError(t, err)
	mockColl.AssertExpectations(t)
}

func TestUserRepository_Delete(t *testing.T) {
	mockColl := new(MockCollection)
	repo := repository.NewUserRepositoryFromCollection(mockColl)

	userID := primitive.NewObjectID()

	mockColl.On("DeleteOne", mock.Anything, mock.Anything).
		Return(&mongo.DeleteResult{DeletedCount: 1}, nil)

	err := repo.Delete(context.Background(), userID.Hex())
	assert.NoError(t, err)
	mockColl.AssertExpectations(t)
}

func TestUserRepository_Count(t *testing.T) {
	mockColl := new(MockCollection)
	repo := repository.NewUserRepositoryFromCollection(mockColl)

	mockColl.On("CountDocuments", mock.Anything, mock.Anything).
		Return(int64(5), nil)

	count, err := repo.Count(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	mockColl.AssertExpectations(t)
}
