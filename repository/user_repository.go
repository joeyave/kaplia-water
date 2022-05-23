package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type UserRepository struct {
	MongoClient *mongo.Client
}

type User struct {
	ID   int64  `bson:"_id"`
	Name string `bson:"name"`
	Role string `bson:"role"`
}

const AdminRole = "admin"

func (r *UserRepository) FindAll(ctx context.Context) ([]*User, error) {
	collection := r.MongoClient.Database(os.Getenv("MONGODB_DATABASE_NAME")).Collection("users")

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var users []*User
	err = cursor.All(context.TODO(), &users)
	return users, err
}

func (r *UserRepository) FindManyByRole(ctx context.Context, role string) ([]*User, error) {
	collection := r.MongoClient.Database(os.Getenv("MONGODB_DATABASE_NAME")).Collection("users")

	cursor, err := collection.Find(ctx, bson.D{{"role", role}})
	if err != nil {
		return nil, err
	}

	var users []*User
	err = cursor.All(context.TODO(), &users)
	return users, err
}

func (r *UserRepository) FindOneByID(ctx context.Context, ID int64) (*User, error) {
	collection := r.MongoClient.Database(os.Getenv("MONGODB_DATABASE_NAME")).Collection("users")

	result := collection.FindOne(ctx, bson.D{{"_id", ID}})
	if result.Err() != nil {
		return nil, result.Err()
	}

	var user *User
	err := result.Decode(&user)
	return user, err
}

func (r *UserRepository) UpdateOne(ctx context.Context, user *User) (*User, error) {
	collection := r.MongoClient.Database(os.Getenv("MONGODB_DATABASE_NAME")).Collection("users")

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	after := options.After
	upsert := true
	opts := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	result := collection.FindOneAndUpdate(context.TODO(), filter, update, &opts)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var updatedUser *User
	err := result.Decode(&updatedUser)
	return updatedUser, err
}
