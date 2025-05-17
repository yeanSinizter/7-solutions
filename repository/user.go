package repository

import (
	"7-solutions/model"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectionInterface interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
}

type UsersRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, id string, user *model.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]model.User, error)
	Count(ctx context.Context) (int64, error)
}

type UserRepository struct {
	collection CollectionInterface
}

func NewUserRepository(db *mongo.Database) UsersRepository {
	return &UserRepository{collection: db.Collection("7-solutions")}
}

func NewUserRepositoryFromCollection(coll CollectionInterface) UsersRepository {
	return &UserRepository{collection: coll}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	user.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	return &user, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return &user, err
}

func (r *UserRepository) Update(ctx context.Context, id string, user *model.User) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"name": user.Name, "email": user.Email}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context) ([]model.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []model.User
	for cursor.Next(ctx) {
		var user model.User
		cursor.Decode(&user)
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}
